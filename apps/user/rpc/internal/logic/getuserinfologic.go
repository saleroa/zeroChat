package logic

import (
	"context"
	"errors"

	"zeroChat/apps/user/models"
	"zeroChat/apps/user/rpc/internal/svc"
	"zeroChat/apps/user/rpc/user"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

var ErrUserNotFound = errors.New("no such user")

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoReq) (*user.GetUserInfoResp, error) {
	// todo: add your logic here and delete this line

	userEntity, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == models.ErrNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	var res user.UserEntity
	// 三方库实现复制
	copier.Copy(&res, userEntity)
	return &user.GetUserInfoResp{
		User: &res,
	}, nil
}
