package friend

import (
	"net/http"
	"zeroChat/apps/social/api/internal/logic/friend"
	"zeroChat/apps/social/api/internal/svc"
	"zeroChat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 在线好友
func FriendOnlineHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FriendOnlineReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := friend.NewFriendOnlineLogic(r.Context(), svcCtx)
		resp, err := l.FriendOnline(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
