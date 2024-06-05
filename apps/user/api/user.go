package main

import (
	"flag"
	"fmt"

	"zeroChat/apps/user/api/internal/config"
	"zeroChat/apps/user/api/internal/handler"
	"zeroChat/apps/user/api/internal/svc"
	"zeroChat/pkg/resultx"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	httpx.SetErrorHandlerCtx(resultx.ErrHandler(c.Name))
	httpx.SetOkHandler(resultx.OkHandler)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
