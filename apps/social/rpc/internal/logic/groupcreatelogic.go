package logic

import (
	"context"

	"zeroChat/apps/social/models"
	"zeroChat/apps/social/rpc/internal/svc"
	"zeroChat/apps/social/rpc/social"
	"zeroChat/pkg/constants"
	"zeroChat/pkg/wuid"
	"zeroChat/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type GroupCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupCreateLogic {
	return &GroupCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 群要求
func (l *GroupCreateLogic) GroupCreate(in *social.GroupCreateReq) (*social.GroupCreateResp, error) {
	// todo: add your logic here and delete this line

	groups := &models.Groups{
		Id:         wuid.GenUid(l.svcCtx.Config.Mysql.Datasource),
		Name:       in.Name,
		Icon:       in.Icon,
		CreatorUid: in.CreatorUid,
		//IsVerify:   true,
		IsVerify: false,
	}

	err := l.svcCtx.GroupsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		_, err := l.svcCtx.GroupsModel.Insert(l.ctx, session, groups)

		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "insert group err %v req %v", err, in)
		}

		_, err = l.svcCtx.GroupMembersModel.Insert(l.ctx, session, &models.GroupMembers{
			GroupId:   groups.Id,
			UserId:    in.CreatorUid,
			RoleLevel: int(constants.CreatorGroupRoleLevel),
		})
		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "insert group member err %v req %v", err, in)
		}
		return nil
	})

	return &social.GroupCreateResp{
		Id: groups.Id,
	}, err
}
