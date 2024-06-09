package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
	"zeroChat/apps/im/ws/internal/config"
	"zeroChat/apps/im/ws/internal/handler"
	"zeroChat/apps/im/ws/internal/svc"
	"zeroChat/apps/im/ws/websocket"
	"zeroChat/pkg/configserver"
	"zeroChat/pkg/constants"
	"zeroChat/pkg/ctxdata"
)

var configFile = flag.String("f", "etc/im.yaml", "the config file")
var wg sync.WaitGroup

func main() {
	flag.Parse()

	var c config.Config

	err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "192.168.117.24:3379",
		ProjectKey:     "2f5bb7747efda0546636fb385a3fa593",
		Namespace:      "im",
		Configs:        "im-ws.yaml",
		ConfigFilePath: "./etc/conf",
		LogLevel:       "DEBUG",
	})).MustLoad(&c, func(bytes []byte) error {
		var c config.Config
		configserver.LoadFromJsonBytes(bytes, &c)

		wg.Add(1)
		go func(c config.Config) {
			defer wg.Done()
			Run(c)
		}(c)
		return nil
	})
	if err != nil {
		panic(err)
	}

	wg.Add(1)
	go func(c config.Config) {
		defer wg.Done()

		Run(c)
	}(c)

	wg.Wait()

}

func Run(c config.Config) {
	if err := c.SetUp(); err != nil {
		panic(err)
	}
	ctx := svc.NewServiceContext(c)
	// 设置服务认证的token
	token, err := ctxdata.GetJwtToken(c.JwtAuth.AccessSecret, time.Now().Unix(), 3153600000, fmt.Sprintf("ws:%s", time.Now().Unix()))
	if err != nil {
		panic(err)
	}

	opts := []websocket.ServerOptions{
		websocket.WithServerAuthentication(handler.NewJwtAuth(ctx)),
		websocket.WithServerDiscover(websocket.NewRedisDiscover(http.Header{
			"Authorization": []string{token},
		}, constants.REDIS_DISCOVER_SRV, c.Redisx)),
	}
	srv := websocket.NewServer(c.ListenOn, opts...)
	defer srv.Stop()

	handler.RegisterHandlers(srv, ctx)

	fmt.Println("start websocket server at ", c.ListenOn, " ..... ")
	srv.Start()
}
