package group

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"zeroChat/apps/social/api/internal/logic/group"
	"zeroChat/apps/social/api/internal/svc"
	"zeroChat/apps/social/api/internal/types"
)

// 创群
func CreateGroupHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GroupCreateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := group.NewCreateGroupLogic(r.Context(), svcCtx)
		resp, err := l.CreateGroup(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
