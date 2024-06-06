package group

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"zeroChat/apps/social/api/internal/logic/group"
	"zeroChat/apps/social/api/internal/svc"
	"zeroChat/apps/social/api/internal/types"
)

// 申请进群
func GroupPutInHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GroupPutInRep
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := group.NewGroupPutInLogic(r.Context(), svcCtx)
		resp, err := l.GroupPutIn(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
