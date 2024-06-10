// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"zeroChat/apps/im/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/chatlog/readRecords",
				Handler: getChatLogReadRecordsHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.JwtAuth.AccessSecret),
		rest.WithPrefix("/v1/im"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/chatlog",
				Handler: getChatLogHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/setup/conversation",
				Handler: setUpUserConversationHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/conversation",
				Handler: getConversationsHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/conversation",
				Handler: putConversationsHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.JwtAuth.AccessSecret),
		rest.WithPrefix("/v1/im"),
	)
}
