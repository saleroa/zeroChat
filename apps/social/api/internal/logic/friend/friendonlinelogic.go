package friend

import (
	"context"
	"zeroChat/apps/social/api/internal/svc"
	"zeroChat/apps/social/api/internal/types"
	"zeroChat/apps/social/rpc/social"
	"zeroChat/pkg/constants"
	"zeroChat/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendOnlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 在线好友
func NewFriendOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendOnlineLogic {
	return &FriendOnlineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendOnlineLogic) FriendOnline(req *types.FriendOnlineReq) (resp *types.FriendOnlineResp, err error) {
	// todo: add your logic here and delete this line

	uid := ctxdata.GetUId(l.ctx)
	friendlist, err := l.svcCtx.Social.FriendList(l.ctx, &social.FriendListReq{
		UserId: uid,
	})
	if err != nil || len(friendlist.List) == 0 {
		return &types.FriendOnlineResp{}, err
	}

	// 查询缓存中的在线用户
	uids := make([]string, 0, len(friendlist.List))
	for _, friend := range friendlist.List {
		uids = append(uids, friend.FriendUid)
	}
	onlines, err := l.svcCtx.Redis.Hgetall(constants.REDIS_ONLINE_USER)
	if err != nil {
		return nil, err
	}

	resOnlineList := make(map[string]bool, len(uids))
	for _, s := range uids {
		if _, ok := onlines[s]; ok {
			resOnlineList[s] = true
		} else {
			resOnlineList[s] = false
		}
	}
	return &types.FriendOnlineResp{
		OnlineList: resOnlineList,
	}, nil
}
