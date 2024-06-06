package group

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"zeroChat/apps/social/api/internal/logic/group"
	"zeroChat/apps/social/api/internal/svc"
	"zeroChat/apps/social/api/internal/types"
)

// 成员列表列表
func GroupUserListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GroupUserListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := group.NewGroupUserListLogic(r.Context(), svcCtx)
		resp, err := l.GroupUserList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
