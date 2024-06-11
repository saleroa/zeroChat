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

// 将需要 ack 的消息添加到队列
// ack 机制模仿的 tcp 握手，这是第一次发送
// 并且客户端的 ack 也是发向这里
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

	// 这俩只有在这里写
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

// // ack 机制中第三个 part, 接收到客户端的 ack 然后处理
// func (s *Server) connAck(conn *Conn) {

// 	for {
// 		select {
// 		case <-conn.done:
// 			// 关闭了连接
// 			s.Infof("close message ack uid %v", conn.Uid)
// 			return
// 		default:
// 		}
//         // 只有严格ack的时候才会用到服务端的 ack
// 		if s.opt.ack == RigorAck {
//            // 取出 websocket 中的 acktype 的消息

// 		   // 检查是否超时，是否有效

// 		   // 如果有效，返回 ack

// 		}

// 		// 取出队列中的第一个数据
// 		message := conn.readMessage[0]
// 		if message.errCount > s.opt.sendErrCount {
// 			s.Infof("conn send fail, message %v, ackType %v, maxSendErrCount %v", message, s.opt.ack.ToString(), s.opt.sendErrCount)
// 			conn.messageMu.Unlock()
// 			// todo：因为发送消息多次错误，而选择放弃消息
// 			delete(conn.readMessageSeq, message.Id)
// 			conn.readMessage = conn.readMessage[1:]
// 			continue
// 		}

// 		// 根据ack的确认策略选择合适的处理方式
// 		switch s.opt.ack {
// 		case OnlyAck:
// 			if err := send(&Message{
// 				FrameType: FrameAck,
// 				AckSeq:    message.AckSeq + 1,
// 				Id:        message.Id,
// 			}, conn); err != nil {
// 				continue
// 			}
// 			// 只回答, 向客户端发送ack
// 			conn.readMessage = conn.readMessage[1:]
// 			conn.messageMu.Unlock()
// 			conn.message <- message
// 			s.Infof("message ack OnlyAck send success mid %v", message.Id)
// 		case RigorAck:
// 			if message.AckSeq == 0 {
// 				// 还未发送过确认信息
// 				conn.readMessage[0].AckSeq++
// 				conn.readMessage[0].ackTime = time.Now()
// 				if err := send(&Message{
// 					FrameType: FrameAck,
// 					AckSeq:    message.AckSeq,
// 					Id:        message.Id,
// 				}, conn); err != nil {
// 					continue
// 				}

// 				conn.messageMu.Unlock()
// 				s.Infof("message ack RigorAck send mid %v, seq %v, time %v", message.Id, message.AckSeq, message.ackTime.Unix())
// 				continue
// 			}

// 			// 已经发送过序号了，需要等待客户端返回确认
// 			msgSeq := conn.readMessageSeq[message.Id]
// 			if msgSeq.AckSeq > message.AckSeq {
// 				// 客户端确认成功,可以处理业务了
// 				conn.readMessage = conn.readMessage[1:]
// 				conn.messageMu.Unlock()
// 				conn.message <- message
// 				s.Infof("message ack RigorAck success mid %v ", message.Id)
// 				continue
// 			}

// 			// 很显然没有处理成功，先看看有没有超时
// 			val := s.opt.ackTimeout - time.Since(message.ackTime)

// 			if !message.ackTime.IsZero() && val <= 0 {
// 				// todo: 超时了，可以选择断开与客户端的连接,但实际具体细节处理仍然还需自己结合业务完善，此处选择放弃该消息
// 				s.Errorf("message ack RigorAck fail mid %v, time %v because timeout", message.Id, message.ackTime)
// 				delete(conn.readMessageSeq, message.Id)
// 				conn.readMessage = conn.readMessage[1:]
// 				conn.messageMu.Unlock()
// 				continue
// 			}

// 			conn.messageMu.Unlock()
// 			if val > 0 && val > 300*time.Microsecond {
// 				if err := send(&Message{
// 					FrameType: FrameAck,
// 					AckSeq:    message.AckSeq,
// 					Id:        message.Id,
// 				}, conn); err != nil {
// 					continue
// 				}
// 			}
// 			// 没有超时，我们让程序等等
// 			time.Sleep(300 * time.Microsecond)
// 		}
// 	}
// }
