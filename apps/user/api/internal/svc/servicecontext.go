package svc

import (
	"zeroChat/apps/user/api/internal/config"
	"zeroChat/apps/user/rpc/userclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		User: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
