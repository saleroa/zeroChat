package main

import (
	"flag"
	"fmt"
	"zeroChat/apps/im/ws/internal/config"
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
	srv := websocket.NewServer(c.ListenOn)

	fmt.Println("start websocket server at :", c.ListenOn, "....")
	srv.Start()

}
