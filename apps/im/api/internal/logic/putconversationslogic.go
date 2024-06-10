package logic

import (
	"context"

	"zeroChat/apps/im/api/internal/svc"
	"zeroChat/apps/im/api/internal/types"
	"zeroChat/apps/im/rpc/im"
	"zeroChat/apps/im/rpc/imclient"
	"zeroChat/pkg/ctxdata"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type PutConversationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPutConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PutConversationsLogic {
	return &PutConversationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PutConversationsLogic) PutConversations(req *types.PutConversationsReq) (resp *types.PutConversationsResp, err error) {
	// todo: add your logic here and delete this line

	uid := ctxdata.GetUId(l.ctx)
	var list map[string]*im.Conversation
	copier.Copy(&list, req.ConversationList)

	data, err := l.svcCtx.PutConversations(l.ctx, &imclient.PutConversationsReq{
		UserId:           uid,
		ConversationList: list,
	})
	if err != nil {
		return nil, err
	}

	var res types.PutConversationsResp
	copier.Copy(&res, data)

	return &res, err
}
