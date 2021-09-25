/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"github.com/dubbogo/pixiu-admin/cmd/admin/controller/configInfo"
	"os"
	"os/signal"
	"strconv"
	"time"
)

import (
	"github.com/gin-gonic/gin"

	"github.com/spf13/cobra"
)

import (
	"github.com/dubbogo/pixiu-admin/cmd/admin/controller"
	"github.com/dubbogo/pixiu-admin/cmd/admin/controller/account"
	"github.com/dubbogo/pixiu-admin/cmd/admin/controller/auth"
	"github.com/dubbogo/pixiu-admin/pkg/config"
	"github.com/dubbogo/pixiu-admin/pkg/logger"
)

var (
	configPath    string
	apiConfigPath string
)

var (
	rootCmd = &cobra.Command{
		Use:   "dubbogo pixiu admin",
		Short: "Dubbogo pixiu admin is the control panel of pixiu gateway.",
		Long: "dubbgo pixiu admin is used to manage the visual interface of dubbogo pixiu, supporting login, user management, \n" +
			"plugin management, service configuration, API key management, interface authority management \n" +
			"(appKey authorization, interface authority, online and offline). \n" +
			"(c) " + strconv.Itoa(time.Now().Year()) + " Dubbogo",
		Version: controller.Version,
		PreRun: func(cmd *cobra.Command, args []string) {
			initDefaultValue()
		},
		Run: func(cmd *cobra.Command, args []string) {
			_, err := config.LoadAPIConfigFromFile(configPath)
			if err != nil {
				logger.Errorf("load admin config  error:%+v", err)
			}
			Start()
			// gracefully shutdown
			sigint := make(chan os.Signal, 1)
			signal.Notify(sigint, os.Interrupt)
			<-sigint
			Stop()
		},
	}
)

// Start start init etcd client and start admin http server
func Start() {
	config.InitEtcdClient()
	r := SetupRouter()
	err := r.Run(config.Bootstrap.GetAddress())
	if err != nil {
		logger.Errorf(err.Error())
	}
}

func Stop() {
	config.CloseEtcdClient()
}

// init Init startCmd
func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", os.Getenv("DUBBOGO_PIXIU_CONFIG"), "Load configuration from `FILE`")
	rootCmd.PersistentFlags().StringVarP(&apiConfigPath, "api-config", "a", os.Getenv("DUBBOGO_PIXIU_API_CONFIG"), "Load api configuration from `FILE`")

}

func getRootCmd() *cobra.Command {
	return rootCmd
}

func initDefaultValue() {
	if configPath == "" {
		configPath = "configs/admin_config.yaml"
	}

	if apiConfigPath == "" {
		apiConfigPath = "configs/api_config.yaml"
	}
}

// main admin run method
func main() {
	app := getRootCmd()

	// ignore error so we don't exit non-zero and break gfmrun README example tests
	_ = app.Execute()
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Guest router
	r.POST("/login", account.Login)
	r.POST("/register", account.Register)

	// auth router
	taR := r.Group("/", auth.JWTAuth())

	// The following router needs to check the token
	{
		// user router
		taR.POST("/user/logout", account.Logout)
		taR.POST("/user/password/edit", account.EditPassword)
		taR.POST("/user/getInfo", account.GetUserInfo)
		taR.POST("/user/getUserRole", account.GetUserRole)
		taR.POST("/user/checkIsAdmin", account.CheckUserIsAdmin)

		taR.GET("/config/api/base", configInfo.GetBaseInfo)
		taR.POST("/config/api/base/", configInfo.SetBaseInfo)
		taR.PUT("/config/api/base/", configInfo.SetBaseInfo)

		taR.GET("/config/api/resource/list", configInfo.GetResourceList)
		taR.GET("/config/api/resource/detail", configInfo.GetResourceDetail)
		taR.POST("/config/api/resource", configInfo.CreateResourceInfo)
		taR.PUT("/config/api/resource", configInfo.ModifyResourceInfo)
		taR.DELETE("/config/api/resource", configInfo.DeleteResourceInfo)

		taR.GET("/config/api/resource/method/list", configInfo.GetMethodList)
		taR.GET("/config/api/resource/method/detail", configInfo.GetMethodDetail)
		taR.POST("/config/api/resource/method", configInfo.CreateMethodInfo)
		taR.PUT("/config/api/resource/method", configInfo.ModifyMethodInfo)
		taR.DELETE("/config/api/resource/method", configInfo.DeleteMethodInfo)

		taR.GET("/config/api/plugin_group/list", configInfo.GetPluginGroupList)
		taR.GET("/config/api/plugin_group/detail", configInfo.GetPluginGroupDetail)
		taR.POST("/config/api/plugin_group", configInfo.CreatePluginGroup)
		taR.PUT("/config/api/plugin_group", configInfo.ModifyPluginGroup)
		taR.DELETE("/config/api/plugin_group", configInfo.DeletePluginGroup)

		taR.GET("/config/api/plugin/ratelimit", configInfo.GetPluginRatelimitDetail)
		taR.POST("/config/api/plugin/ratelimit", configInfo.CreatePluginRatelimit)
		taR.PUT("/config/api/plugin/ratelimit", configInfo.ModifyPluginRatelimit)
		taR.DELETE("/config/api/plugin/ratelimit", configInfo.DeletePluginRatelimit)
	}

	return r
}

