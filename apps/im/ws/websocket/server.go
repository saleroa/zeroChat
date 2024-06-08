package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

// im 服务唯一的 server
type Server struct {
	sync.RWMutex   // 保证 conn 的 map 的读写不冲突
	connToUser     map[*Conn]string
	userToConn     map[string]*Conn
	authentication Authentication
	routes         map[string]HandlerFunc
	addr           string
	patten         string
	upgrader       websocket.Upgrader
	opt            *serverOption
	logx.Logger
}

func NewServer(addr string, opts ...ServerOptions) *Server {

	opt := newServerOption(opts...)

	return &Server{
		routes:         make(map[string]HandlerFunc),
		addr:           addr,
		patten:         opt.patten,
		opt:            &opt,
		upgrader:       websocket.Upgrader{},
		Logger:         logx.WithContext(context.Background()),
		authentication: opt.Authentication,
		connToUser:     make(map[*Conn]string),
		userToConn:     make(map[string]*Conn),
	}
}

// im 请求的处理函数
func (s *Server) ServerWs(w http.ResponseWriter, r *http.Request) {

	// 错误的捕捉恢复，recover
	defer func() {
		if r := recover(); r != nil {
			s.Errorf("server handler ws recover err %v", r)
		}
	}()

	conn := NewConn(s, w, r)
	if conn == nil {
		return
	}
	// conn, err := s.upgrader.Upgrade(w, r, nil)
	// if err != nil {
	// 	s.Errorf("upgrade err %v", err)
	//     return
	// }

	// auth 认证不通过
	if !s.authentication.Auth(w, r) {
		// write message 是 conn 的方法，写到某个连接里
		s.Send(&Message{FrameType: FrameData, Data: "不具备访问权限"}, conn)
		// conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprint("不具备访问权限")))
		conn.Close()
		return
	}

	// 连接对象获取请求
	// method
	s.addConn(conn, r)
	s.handlerConn(conn)

}

// 根据连接对象执行任务处理
func (s *Server) handlerConn(conn *Conn) {

	for {
		// 获取请求消息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("websocket conn read message err %v", err)
			s.Close(conn)
			return
		}

		// 解析消息
		var message Message
		if err = json.Unmarshal(msg, &message); err != nil {
			s.Errorf("json unmarshall err %v, msg %v", err, string(msg))
			s.Close(conn)
			return
		}

		switch message.FrameType {
		case FramePing:
			s.Send(&Message{FrameType: FramePing}, conn)
		case FrameData:
			// 根据method分发路由并执行
			if handler, ok := s.routes[message.Method]; ok {
				handler(s, conn, &message)
			} else {
				s.Send(&Message{FrameType: FrameData, Data: fmt.Sprintf("不存在执行的方法 %v", message.Method)}, conn)
				// conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("不存在执行的方法 %v", message.Method)))
			}
		}

	}

}

func (s *Server) AddRouters(rs []Route) {
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

func (s *Server) addConn(conn *Conn, req *http.Request) {
	uid := s.authentication.UserId(req)

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()
	// 如果有重复登录
	if c := s.userToConn[uid]; c != nil {
		c.Close()
	}

	s.connToUser[conn] = uid
	s.userToConn[uid] = conn
}

func (s *Server) GetConn(uid string) *Conn {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

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
		return
	}

	delete(s.connToUser, conn)
	delete(s.userToConn, uid)
	conn.Close()
}

// 根据 id 查询连接，发送消息
func (s *Server) SendByUserId(msg interface{}, sendIds ...string) error {
	if len(sendIds) == 0 {
		return nil
	}

	return s.Send(msg, s.GetConns(sendIds...)...)

}

// 根据连接发送消息
func (s *Server) Send(msg interface{}, conns ...*Conn) error {
	if len(conns) == 0 {
		return nil
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for _, conn := range conns {
		if err = conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return err
		}
	}
	return nil
}
