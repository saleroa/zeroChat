package websocket

import "github.com/gorilla/websocket"

// 路由

type Route struct {
	Method  string
	Handler HandlerFunc
}

type HandlerFunc func(srv *Server, conn *websocket.Conn, msg *Message)
