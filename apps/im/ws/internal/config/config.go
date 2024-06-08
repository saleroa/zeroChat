package config

import "github.com/zeromicro/go-zero/core/service"

type Config struct {
	// 日志，性能监听
	service.ServiceConf

	ListenOn string

	JwtAuth struct {
		AccessSecret string
	}
}
