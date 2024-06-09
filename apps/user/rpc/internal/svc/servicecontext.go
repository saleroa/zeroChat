package svc

import (
	"time"
	"zeroChat/apps/user/models"
	"zeroChat/apps/user/rpc/internal/config"
	"zeroChat/pkg/constants"
	"zeroChat/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	UserModel models.UsersModel
	*redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:    c,
		Redis:     redis.MustNewRedis(c.Redisx),
		UserModel: models.NewUsersModel(sqlConn, c.Cache),
	}
}

func (svc *ServiceContext) SetRootToken() error {
	// 生成jwt
	systemToken, err := ctxdata.GetJwtToken(svc.Config.Jwt.AccessSecret, time.Now().Unix(), 999999999, constants.SYSTEM_ROOT_UID)
	if err != nil {
		return err
	}
	// 写入到redis

	return svc.Redis.Set(constants.REDIS_SYSTEM_ROOT_TOKEN, systemToken)
}
