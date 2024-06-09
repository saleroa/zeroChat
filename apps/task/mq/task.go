package main

import (
	"flag"
	"fmt"
	"sync"
	"zeroChat/apps/task/mq/internal/config"
	"zeroChat/apps/task/mq/internal/handler"
	"zeroChat/apps/task/mq/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/task.yaml", "the config file")
var wg sync.WaitGroup

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	Run(c)
	//err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
	//	ETCDEndpoints:  "192.168.117.24:3379",
	//	ProjectKey:     "2f5bb7747efda0546636fb385a3fa593",
	//	Namespace:      "task",
	//	Configs:        "task-mq.yaml",
	//	ConfigFilePath: "./etc/conf",
	//	LogLevel:       "DEBUG",
	//})).MustLoad(&c, func(bytes []byte) error {
	//	var c config.Config
	//	configserver.LoadFromJsonBytes(bytes, &c)
	//
	//	wg.Add(1)
	//	go func(c config.Config) {
	//		defer wg.Done()
	//
	//		Run(c)
	//	}(c)
	//	return nil
	//})
	//if err != nil {
	//	panic(err)
	//}
	//
	//wg.Add(1)
	//go func(c config.Config) {
	//	defer wg.Done()
	//
	//	Run(c)
	//}(c)
	//
	//wg.Wait()
}

func Run(c config.Config) {
	if err := c.SetUp(); err != nil {
		panic(err)
	}
	ctx := svc.NewServiceContext(c)

	listen := handler.NewListen(ctx)

	// service 是接口，有 start 和 stop 的方法
	serviceGroup := service.NewServiceGroup()
	// 启动 kafka 的消费者
	for _, s := range listen.Services() {
		serviceGroup.Add(s)
	}
	fmt.Println("Starting mqueue at ...")
	// 启动 group 里的所有 service
	serviceGroup.Start()
}
