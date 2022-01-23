package core

import (
	"fmt"
	"net/http"
)

import (
	"go.uber.org/zap"

	"v.marlon.life/toolkit/util"
)

import (
	"github.com/dubbogo/pixiu-admin/pkg/config"
	"github.com/dubbogo/pixiu-admin/pkg/global"
	"github.com/dubbogo/pixiu-admin/pkg/initialize"
)

var (
	helperInfo = `
	Welcome DUBBOGO-PIXIU-ADMIN
	Default doc address: http://127.0.0.1%s/swagger/index.html
	Default running address: http://127.0.0.1:8080
`
)

type server interface {
	ListenAndServe() error
}

// RunServer start server
func RunServer() {
	// load config
	global.VP = Viper()
	global.LOG = Zap()

	config.InitEtcdClient()
	router := initialize.Routers()

	address := fmt.Sprintf(":%d", global.CONFIG.System.Addr)

	s := initServer(address, router)

	var wg util.WaitGroupWrapper

	wg.AddAndRun(func() {
		global.LOG.Info("server run success on ", zap.String("address", address))
		fmt.Printf(helperInfo, address)

		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.LOG.Error(err.Error())
		}
	})

	wg.AddAndRun(func() {
		global.LOG.Info("xDS server run success on :18000")
		if err := StartxDsServer(); err != nil {
			global.LOG.Error(err.Error())
		}
	})

	wg.Wait()
}
