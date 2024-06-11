package svc

import (
	"zeroChat/apps/im/rpc/imclient"
	"zeroChat/apps/social/api/internal/config"
	"zeroChat/apps/social/rpc/socialclient"
	"zeroChat/apps/user/rpc/userclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	ImRpc   imclient.Im
	UserRpc userclient.User
	Social  socialclient.Social
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		UserRpc: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		ImRpc:   imclient.NewIm(zrpc.MustNewClient(c.ImRpc)),
		Social:  socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
	}
}
