package svc

import (
	"zeroChat/apps/user/models"
	"zeroChat/apps/user/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	UserModel models.UsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config: c,

		UserModel: models.NewUsersModel(sqlConn, c.Cache),
	}
}
