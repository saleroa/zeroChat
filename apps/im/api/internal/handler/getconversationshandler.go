package handler

import (
	"net/http"

	"zeroChat/apps/im/api/internal/logic"
	"zeroChat/apps/im/api/internal/svc"
	"zeroChat/apps/im/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func getConversationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetConversationsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetConversationsLogic(r.Context(), svcCtx)
		resp, err := l.GetConversations(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
