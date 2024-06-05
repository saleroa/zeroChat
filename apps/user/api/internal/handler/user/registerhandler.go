package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"zeroChat/apps/user/api/internal/logic/user"
	"zeroChat/apps/user/api/internal/svc"
	"zeroChat/apps/user/api/internal/types"
)

// 用户注册
func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewRegisterLogic(r.Context(), svcCtx)
		resp, err := l.Register(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
