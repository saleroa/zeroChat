package handler

import (
	"zeroChat/apps/im/ws/internal/handler/conversation"
	"zeroChat/apps/im/ws/internal/handler/push"
	"zeroChat/apps/im/ws/internal/handler/user"
	"zeroChat/apps/im/ws/internal/svc"
	"zeroChat/apps/im/ws/websocket"
)

func RegisterHandlers(srv *websocket.Server, svc *svc.ServiceContext) {
	srv.AddRoutes([]websocket.Route{
		{
			Method:  "user.online",
			Handler: user.OnLine(svc),
		},
		{
			Method:  "conversation.chat",
			Handler: conversation.Chat(svc),
		},
		{
			Method:  "conversation.markChat",
			Handler: conversation.MarkRead(svc),
		},
		{
			Method:  "push",
			Handler: push.Push(svc),
		},
	})
}
