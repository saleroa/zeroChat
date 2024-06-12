package group

import (
	"context"
	"zeroChat/apps/social/api/internal/svc"
	"zeroChat/apps/social/api/internal/types"
	"zeroChat/apps/social/rpc/social"
	"zeroChat/pkg/constants"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupUserOnlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 群在线用户
func NewGroupUserOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUserOnlineLogic {
	return &GroupUserOnlineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupUserOnlineLogic) GroupUserOnline(req *types.GroupUserOnlineReq) (resp *types.GroupUserOnlineResp, err error) {
	// todo: add your logic here and delete this line

	groupUsers, err := l.svcCtx.Social.GroupUsers(l.ctx, &social.GroupUsersReq{
		GroupId: req.GroupId,
	})

	if err != nil || len(groupUsers.List) == 0 {
		return &types.GroupUserOnlineResp{}, err
	}

	// 查询缓存中的在线用户
	uids := make([]string, 0, len(groupUsers.List))
	for _, user := range groupUsers.List {
		uids = append(uids, user.UserId)
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
	return &types.GroupUserOnlineResp{
		OnlineList: resOnlineList,
	}, nil
}
