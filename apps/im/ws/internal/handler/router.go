package handler

import (
	"zeroChat/apps/im/ws/internal/handler/user"
	"zeroChat/apps/im/ws/internal/svc"
	"zeroChat/apps/im/ws/websocket"
)

func RegisterHandlers(srv *websocket.Server, svc *svc.ServiceContext) {
	srv.AddRouters([]websocket.Route{
		{
			Method:  "user.online",
			Handler: user.OnLine(svc),
		},
	})
}
