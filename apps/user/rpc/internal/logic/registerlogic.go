package logic

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"zeroChat/apps/user/models"
	"zeroChat/apps/user/rpc/internal/svc"
	"zeroChat/apps/user/rpc/user"
	"zeroChat/pkg/ctxdata"
	"zeroChat/pkg/encrypt"
	"zeroChat/pkg/wuid"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrPhoneIsRegister = errors.New("phone number already registered")
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	// todo: add your logic here and delete this line

	// 验证用户是否注册，验证手机号
	userEntity, err := l.svcCtx.UserModel.FindByPhone(l.ctx, in.Phone)
	if err != nil && err != models.ErrNotFound {
		return nil, err
	}
	if userEntity != nil {
		return nil, ErrPhoneIsRegister
	}
	// 定义用户信息
	userEntity = &models.Users{
		Id:       wuid.GenUid(l.svcCtx.Config.Mysql.DataSource),
		Avatar:   in.Avatar,
		Nickname: in.Nickname,
		Phone:    in.Phone,
		Sex: sql.NullInt64{
			Int64: int64(in.Sex),
			Valid: true,
		},
	}
	// 密码长度校验
	if len(in.Password) > 0 {
		genPass, err := encrypt.GenPasswordHash([]byte(in.Password))
		if err != nil {
			return nil, err
		}
		userEntity.Password = sql.NullString{
			String: string(genPass),
			Valid:  true,
		}
	}

	// 插入用户
	_, err = l.svcCtx.UserModel.Insert(l.ctx, userEntity)
	if err != nil {
		return nil, err
	}

	// 生成 token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userEntity.Id)
	if err != nil {
		return nil, err
	}

	return &user.RegisterResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
