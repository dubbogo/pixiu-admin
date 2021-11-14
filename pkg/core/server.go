package core

import (
	"fmt"
)

import (
	"github.com/dubbogo/pixiu-admin/pkg/config"
	"github.com/dubbogo/pixiu-admin/pkg/global"
	"github.com/dubbogo/pixiu-admin/pkg/initialize"

	"go.uber.org/zap"
)

type server interface {
	ListenAndServe() error
}

func RunServer() {

	global.VP = Viper()

	global.LOG = Zap()

	config.InitEtcdClient()

	router := initialize.Routers()

	address := fmt.Sprintf(":%d", global.CONFIG.System.Addr)

	s := initServer(address, router)

	global.LOG.Info("server run success on ", zap.String("address", address))

	fmt.Printf(`
	欢迎使用 DUBBOGO-PIXIU-ADMIN
	默认自动化文档地址:http://127.0.0.1%s/swagger/index.html
	默认前端文件运行地址:http://127.0.0.1:8080
`, address)

	global.LOG.Error(s.ListenAndServe().Error())
}
