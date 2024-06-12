package svc

import (
	"zeroChat/apps/user/api/internal/config"
	"zeroChat/apps/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	Redis  *redis.Redis
	userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Redis:  redis.MustNewRedis(c.Redisx),
		User:   userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
