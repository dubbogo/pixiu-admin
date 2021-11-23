package core

import (
	"fmt"
)

import (
	"go.uber.org/zap"
)

import (
	"github.com/dubbogo/pixiu-admin/pkg/config"
	"github.com/dubbogo/pixiu-admin/pkg/global"
	"github.com/dubbogo/pixiu-admin/pkg/initialize"
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
	Welcome DUBBOGO-PIXIU-ADMIN
	Default doc address: http://127.0.0.1%s/swagger/index.html
	Default running address: http://127.0.0.1:8080
`, address)

	global.LOG.Error(s.ListenAndServe().Error())
}
