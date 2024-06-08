package main

import (
	"flag"
	"fmt"
	"zeroChat/apps/im/ws/internal/config"
	"zeroChat/apps/im/ws/internal/handler"
	"zeroChat/apps/im/ws/internal/svc"
	"zeroChat/apps/im/ws/websocket"

	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/im.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 就是配置中 serviceconf 的作用，启动系统监听，日志
	if err := c.SetUp(); err != nil {
		panic(err)
	}

	ctx := svc.NewServiceContext(c)
	srv := websocket.NewServer(c.ListenOn,
		websocket.WithServerAuthentication(handler.NewJwtAuth(ctx)))

	defer srv.Stop()

	handler.RegisterHandlers(srv, ctx)
	fmt.Println("start websocket server at :", c.ListenOn, "....")
	srv.Start()

}
