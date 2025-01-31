package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/threading"

	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type AckType int

const (
	NoAck    AckType = iota
	OnlyAck          // ack 一次
	RigorAck         // 严格模式
)

func (t AckType) ToString() string {
	switch t {
	case OnlyAck:
		return "OnlyAck"
	case RigorAck:
		return "RigorAck"
	}

	return "NoAck"
}

type Server struct {
	sync.RWMutex

	*threading.TaskRunner

	opt            *serverOption
	authentication Authentication

	routes     map[string]HandlerFunc
	addr       string
	patten     string
	listenOn   string
	discover   Discover
	connToUser map[*Conn]string
	userToConn map[string]*Conn

	upgrader websocket.Upgrader
	logx.Logger
}

func NewServer(addr string, opts ...ServerOptions) *Server {
	opt := newServerOptions(opts...)

	s := &Server{
		routes: make(map[string]HandlerFunc),
		addr:   addr,
		patten: opt.patten,
		opt:    &opt,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},

		authentication: opt.Authentication,

		connToUser: make(map[*Conn]string),
		userToConn: make(map[string]*Conn),

		listenOn:   FigureOutListenOn(addr),
		Logger:     logx.WithContext(context.Background()),
		TaskRunner: threading.NewTaskRunner(opt.concurrency),
	}

	// 存在服务发现，采用分布式im通信的时候; 默认不做任何处理
	s.discover.Register(s.listenOn)

	return s
}

// ws 请求的入口函数
func (s *Server) ServerWs(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			s.Errorf("server handler ws recover err %v", r)
		}
	}()

	conn := NewConn(s, w, r)
	if conn == nil {
		return
	}

	if !s.authentication.Auth(w, r) {
		s.Send(&Message{FrameType: FrameData, Data: "不具备访问权限"}, conn)
		conn.Close()
		return
	}

	// 记录连接
	s.addConn(conn, r)

	// 处理连接，主要逻辑部分
	go s.handlerConn(conn)
}

// 根据连接对象执行任务处理
func (s *Server) handlerConn(conn *Conn) {

	uids := s.GetUsers(conn)
	conn.Uid = uids[0]

	// 如果存在服务发现则进行注册；默认不做任何处理
	s.discover.BoundUser(conn.Uid)
	// 处理任务
	go s.handlerWrite(conn)

	if s.isAck(nil) {
		go s.serverAck(conn)
	}

	for {
		// 获取请求消息
		_, msg, err := conn.ReadMessage()
		fmt.Println("new msg ", string(msg), err)
		if err != nil {
			s.Errorf("websocket conn read message err %v", err)
			s.Close(conn)
			return
		}
		// 解析消息
		var message Message
		if err = json.Unmarshal(msg, &message); err != nil {
			s.Errorf("json unmarshal err %v, msg %v", err, string(msg))
			continue
		}

		// 回复一个 ack 后再进行消息处理
		// 依据消息进行处理
		if s.isAck(&message) {
			// if message.FrameType == FrameAck
			s.Infof("conn message read ack msg %v", message)
			conn.appendMsgMq(&message)
		} else {
			// message 是直接发送到连接的
			conn.message <- &message
		}
	}
}

// 判断是否启动ack
func (s *Server) isAck(message *Message) bool {
	if message == nil {
		return s.opt.ack != NoAck
	}
	return s.opt.ack != NoAck && message.FrameType != FrameNoAck && message.FrameType != FrameTranspond
}

// 任务的处理
func (s *Server) handlerWrite(conn *Conn) {
	for {
		// 判断连接是否关闭
		select {
		case <-conn.done:
			// 连接关闭
			return
		case message := <-conn.message:
			switch message.FrameType {
			case FramePing:
				s.Send(&Message{FrameType: FramePing}, conn)
			case FrameData:
				// 根据请求的method分发路由并执行
				if handler, ok := s.routes[message.Method]; ok {
					handler(s, conn, message)
				} else {
					s.Send(&Message{FrameType: FrameData, Data: fmt.Sprintf("不存在执行的方法 %v 请检查", message.Method)}, conn)
				}
			}

			// ack 消息的清理
			if s.isAck(message) {
				conn.messageMu.Lock()
				delete(conn.readMessageSeq, message.Id)
				conn.messageMu.Unlock()
			}
		}
	}
}

// ack 机制中第二个部分，服务端对接受到的信息进行 ack
func (s *Server) serverAck(conn *Conn) {

	send := func(msg *Message, conn *Conn) error {
		err := s.Send(msg, conn)
		if err == nil {
			return nil
		}

		s.Errorf("message ack OnlyAck send err %v message %v", err, msg)
		conn.readMessage[0].errCount++
		conn.messageMu.Unlock()

		tempDelay := time.Duration(200*conn.readMessage[0].errCount) * time.Microsecond
		if max := 1 * time.Second; tempDelay > max {
			tempDelay = max
		}

		time.Sleep(tempDelay)
		return err
	}

	for {
		select {
		case <-conn.done:
			// 关闭了连接
			s.Infof("close message ack uid %v", conn.Uid)
			return
		default:
		}

		conn.messageMu.Lock()
		if len(conn.readMessage) == 0 {
			conn.messageMu.Unlock()
			// 没有消息可以睡眠100 毫秒 -- 目的是让程序缓一缓
			time.Sleep(100 * time.Microsecond)
			continue
		}

		// 取出队列中的第一个数据
		message := conn.readMessage[0]
		if message.errCount > s.opt.sendErrCount {
			s.Infof("conn send fail, message %v, ackType %v, maxSendErrCount %v", message, s.opt.ack.ToString(), s.opt.sendErrCount)
			conn.messageMu.Unlock()
			// todo：因为发送消息多次错误，而选择放弃消息
			// msg 和 seq 都删除
			delete(conn.readMessageSeq, message.Id)
			conn.readMessage = conn.readMessage[1:]
			continue
		}

		// 根据ack的确认策略选择合适的处理方式
		switch s.opt.ack {

		// 只需要服务端 ack 一次
		case OnlyAck:
			if err := send(&Message{
				FrameType: FrameAck,
				AckSeq:    message.AckSeq + 1,
				Id:        message.Id,
			}, conn); err != nil {
				continue
			}
			// 只回答, 向客户端发送ack
			// 只删除 msg ，不删 seq
			conn.readMessage = conn.readMessage[1:]
			conn.messageMu.Unlock()
			conn.message <- message
			s.Infof("message ack OnlyAck send success mid %v", message.Id)

		// 严格 ack
		case RigorAck:
			if message.AckSeq == 0 {
				// 还未发送过确认信息
				conn.readMessage[0].AckSeq++
				conn.readMessage[0].ackTime = time.Now()
				if err := send(&Message{
					FrameType: FrameAck,
					AckSeq:    message.AckSeq,
					Id:        message.Id,
				}, conn); err != nil {
					continue
				}

				conn.messageMu.Unlock()
				s.Infof("message ack RigorAck send mid %v, seq %v, time %v", message.Id, message.AckSeq, message.ackTime.Unix())
				continue
			}

			// 已经发送过序号了，需要等待客户端返回确认
			msgSeq := conn.readMessageSeq[message.Id]

			// 1.客户端确认成功,可以处理业务了
			if msgSeq.AckSeq > message.AckSeq {

				// 只删除 msg ，不删 seq
				conn.readMessage = conn.readMessage[1:]
				conn.messageMu.Unlock()
				conn.message <- message
				s.Infof("message ack RigorAck success mid %v ", message.Id)
				continue
			}

			// 2. 没有处理成功，判断有没有超时
			val := s.opt.ackTimeout - time.Since(message.ackTime)

			if !message.ackTime.IsZero() && val <= 0 {
				// 2.2  超时了，放弃消息
				s.Errorf("message ack RigorAck fail mid %v, time %v because timeout", message.Id, message.ackTime)
				delete(conn.readMessageSeq, message.Id)
				conn.readMessage = conn.readMessage[1:]
				conn.messageMu.Unlock()
				continue
			}

			// 2.1 未超过，重新发送
			conn.messageMu.Unlock()
			s.Send(&Message{
				FrameType: FrameAck,
				Id:        message.Id,
				AckSeq:    message.AckSeq,
			}, conn)
			// 睡眠一定的时间
			time.Sleep(3 * time.Second)
		}
	}
}

// server basic functions

func (s *Server) addConn(conn *Conn, req *http.Request) {
	uid := s.authentication.UserId(req)

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	// 验证用户是否之前登入过
	if c := s.userToConn[uid]; c != nil {
		// 关闭之前的连接
		c.Close()
	}

	s.connToUser[conn] = uid
	s.userToConn[uid] = conn
}

func (s *Server) GetConn(uid string) *Conn {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	fmt.Println(s.userToConn)
	return s.userToConn[uid]
}

func (s *Server) GetConns(uids ...string) []*Conn {
	if len(uids) == 0 {
		return nil
	}

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	res := make([]*Conn, 0, len(uids))
	for _, uid := range uids {
		res = append(res, s.userToConn[uid])
	}
	return res
}

func (s *Server) GetUsers(conns ...*Conn) []string {

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	var res []string
	if len(conns) == 0 {
		// 获取全部
		res = make([]string, 0, len(s.connToUser))
		for _, uid := range s.connToUser {
			res = append(res, uid)
		}
	} else {
		// 获取部分
		res = make([]string, 0, len(conns))
		for _, conn := range conns {
			res = append(res, s.connToUser[conn])
		}
	}

	return res
}

func (s *Server) Close(conn *Conn) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	uid := s.connToUser[conn]
	if uid == "" {
		// 已经被关闭
		return
	}

	delete(s.connToUser, conn)
	delete(s.userToConn, uid)

	conn.Close()
}

func (s *Server) SendByUserId(msg interface{}, sendIds ...string) error {
	if len(sendIds) == 0 {
		return nil
	}

	return s.Send(msg, s.GetConns(sendIds...)...)
}

func (s *Server) Send(msg interface{}, conns ...*Conn) error {
	if len(conns) == 0 {
		return nil
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) AddRoutes(rs []Route) {
	for _, r := range rs {
		s.routes[r.Method] = r.Handler
	}
}

func (s *Server) Start() {
	http.HandleFunc(s.patten, s.ServerWs)
	s.Info(http.ListenAndServe(s.addr, nil))
}

func (s *Server) Stop() {
	fmt.Println("停止服务")
}
