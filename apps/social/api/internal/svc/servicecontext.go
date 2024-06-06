package svc

import (
	"zeroChat/apps/social/api/internal/config"
	"zeroChat/apps/social/rpc/socialclient"
	"zeroChat/apps/user/rpc/userclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	UserRpc userclient.User
	Social  socialclient.Social
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		UserRpc: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),

		Social: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
	}
}
