package websocket

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// 封装过后的连接对象
// 特点之一
type Conn struct {
	idleMu sync.Mutex

	Uid string

	*websocket.Conn
	s *Server

	idle time.Time

	// ## ack 机制
	messageMu sync.Mutex
	// 接收消息处理队列，采用数组存储可保障数据的顺序
	readMessage []*Message
	// 记录ack机制中的消息处理结果与进展
	readMessageSeq map[string]*Message
	// 消息通道，在ack验证完成后将消息投递与writeHandler处理
	message chan *Message

	// ## 心跳检测机制
	// 最大心跳检测时间
	maxConnectionIdle time.Duration
	done              chan struct{}
}

func NewConn(s *Server, w http.ResponseWriter, r *http.Request) *Conn {
	var respHeader http.Header
	if protocol := r.Header.Get("Sec-Websocket-Protocol"); protocol != "" {
		respHeader = http.Header{
			"Sec-Websocket-Protocol": []string{protocol},
		}
	}
	c, err := s.upgrader.Upgrade(w, r, respHeader)
	if err != nil {
		s.Errorf("upgrade err %v", err)
		return nil
	}

	conn := &Conn{
		Conn:              c,
		s:                 s,
		idle:              time.Now(),
		maxConnectionIdle: s.opt.maxConnectionIdle,
		readMessage:       make([]*Message, 0, 2),
		readMessageSeq:    make(map[string]*Message, 2),

		message: make(chan *Message, 1), // 减少阻塞，保证顺序
		done:    make(chan struct{}),
	}

	go conn.keepalive()
	return conn
}

// 将需要 ack 的消息添加到队列中
// 关于 ack 机制，去看看 ack 实现逻辑图
func (c *Conn) appendMsgMq(msg *Message) {
	c.messageMu.Lock()
	defer c.messageMu.Unlock()
	// 只有seq 和 msg 同时存在，并且 最新seq 大于 之前 seq 才有效

	// 已经有消息的记录，该消息已经有ack的确认
	if m, ok := c.readMessageSeq[msg.Id]; ok {
		if len(c.readMessage) == 0 {
			// 队列中没有该消息
			return
		}

		// msg.AckSeq > m.AckSeq，最新的序号一定要大于之前的序号
		if m.AckSeq >= msg.AckSeq {
			// 没有进行ack的确认, 重复
			return
		}

		c.readMessageSeq[msg.Id] = msg
		return
	}

	// 没有进行 ack 认证，就发送了 ack
	// 避免客户端重复发送多余的ack消息
	if msg.FrameType == FrameAck {
		return
	}

	c.readMessage = append(c.readMessage, msg)
	c.readMessageSeq[msg.Id] = msg

}

// 心跳检测机制
func (c *Conn) keepalive() {
	idleTimer := time.NewTimer(c.maxConnectionIdle)
	defer func() {
		idleTimer.Stop()
	}()

	for {
		select {
		case <-idleTimer.C:
			c.idleMu.Lock()
			idle := c.idle
			if idle.IsZero() { // The connection is non-idle.
				c.idleMu.Unlock()
				idleTimer.Reset(c.maxConnectionIdle)
				continue
			}
			val := c.maxConnectionIdle - time.Since(idle)
			c.idleMu.Unlock()
			if val <= 0 {
				// The connection has been idle for a duration of keepalive.MaxConnectionIdle or more.
				// Gracefully close the connection.
				c.s.Close(c)
				return
			}
			idleTimer.Reset(val)
		case <-c.done:
			return
		}
	}
}

func (c *Conn) ReadMessage() (messageType int, p []byte, err error) {
	messageType, p, err = c.Conn.ReadMessage()

	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	// 读取消息的时候将idle设置为0
	c.idle = time.Time{}
	return
}

func (c *Conn) WriteMessage(messageType int, data []byte) error {
	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	// 方法是并发不安全，所以枷锁
	err := c.Conn.WriteMessage(messageType, data)
	// 已经读取完消息，将 idle 设置为 now
	c.idle = time.Now()
	return err
}

func (c *Conn) Close() error {
	select {
	case <-c.done:
	default:
		close(c.done)
	}

	return c.Conn.Close()
}
