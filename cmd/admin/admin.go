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
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"
)

import (
	fc "github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
)

import (
	"github.com/dubbogo/pixiu-admin/pkg/common/yaml"
	"github.com/dubbogo/pixiu-admin/pkg/config"
	"github.com/dubbogo/pixiu-admin/pkg/logger"
	logic "github.com/dubbogo/pixiu-admin/pkg/logic"
)

var (
	cmdStart = cli.Command{
		Name:  "start",
		Usage: "start dubbogo pixiu admin",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "config, c",
				Usage:  "Load configuration from `FILE`",
				EnvVar: "PROXY_ADMIN_CONFIG",
				Value:  "configs/admin_config.yaml",
			},
		},
		Action: func(c *cli.Context) error {
			configPath := c.String("config")
			_, err := config.LoadAPIConfigFromFile(configPath)
			if err != nil {
				logger.Errorf("load admin config  error:%+v", err)
				return err
			}
			Start()
			return nil
		},
	}

	cmdStop = cli.Command{
		Name:  "stop",
		Usage: "stop dubbogo pixiu admin",
		Action: func(c *cli.Context) error {
			Stop()
			return nil
		},
	}
)

// Version admin version
const Version = "0.1.0"

func newAdminApp(startCmd *cli.Command) *cli.App {
	app := cli.NewApp()
	app.Name = "dubbogo pixiu admin"
	app.Version = Version
	app.Compiled = time.Now()
	app.Copyright = "(c) " + strconv.Itoa(time.Now().Year()) + " Dubbogo"
	app.Usage = "Dubbogo pixiu admin"
	app.Flags = cmdStart.Flags

	// commands
	app.Commands = []cli.Command{
		cmdStart,
		cmdStop,
	}

	// action
	app.Action = func(c *cli.Context) error {
		if c.NumFlags() == 0 {
			return cli.ShowAppHelp(c)
		}
		return startCmd.Action.(func(c *cli.Context) error)(c)
	}

	return app
}

func main() {
	app := newAdminApp(&cmdStart)
	// ignore error so we don't exit non-zero and break gfmrun README example tests
	_ = app.Run(os.Args)
}

// Start start init etcd client and start admin http server
func Start() {
	config.InitEtcdClient()
	r := SetupRouter()
	r.Run(config.Bootstrap.GetAddress()) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func Stop() {
	config.CloseEtcdClient()
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/config/api/base", GetBaseInfo)
	r.POST("/config/api/base/", SetBaseInfo)
	r.PUT("/config/api/base/", SetBaseInfo)
	r.GET("/config/api/resource/list", GetResourceList)
	r.GET("/config/api/resource/detail", GetResourceDetail)
	r.POST("/config/api/resource", SetResourceInfo)
	r.PUT("/config/api/resource", ModifyResourceInfo)
	r.DELETE("/config/api/resource", DeleteResourceInfo)
	r.GET("/config/api/resource/method/list", GetMethodList)
	r.GET("/config/api/resource/method/detail", GetMethodDetail)
	r.POST("/config/api/resource/method", SetMethodInfo)
	r.PUT("/config/api/resource/method", ModifyMethodInfo)
	r.DELETE("/config/api/resource/method", DeleteMethodInfo)
	r.GET("/config/api/plugin_group/list", GetPluginGroupList)
	r.GET("/config/api/plugin_group/detail", GetPluginGroupDetail)
	r.POST("/config/api/plugin_group", CreatePluginGroup)
	r.PUT("/config/api/plugin_group", ModifyPluginGroup)
	r.DELETE("/config/api/plugin_group", DeletePluginGroup)

	return r
}

const OK = "10001"
const ERR = "10002"

func WithError(err error) config.RetData {
	return config.RetData{ERR, err.Error()}
}

func WithRet(data interface{}) config.RetData {
	return config.RetData{OK, data}
}

// GetResourceList get all resource list
func GetResourceList(c *gin.Context) {
	res, err := logic.BizGetResourceList()
	if err != nil {
		c.JSON(http.StatusOK, WithError(err))
		return
	}
	data, _ := json.Marshal(res)
	c.JSON(http.StatusOK, WithRet(string(data)))
}

// GetResourceList get all resource list
func GetMethodList(c *gin.Context) {
	path := c.Query("path")

	res, err := logic.BizGetMethodList(path)
	if err != nil {
		c.JSON(http.StatusOK, WithError(err))
		return
	}
	data, _ := json.Marshal(res)
	c.JSON(http.StatusOK, WithRet(string(data)))
}

// SetBaseInfo handle create base info http request
func SetBaseInfo(c *gin.Context) {
	body := c.PostForm("content")

	baseInfo := &config.BaseInfo{}
	err := yaml.UnmarshalYML([]byte(body), baseInfo)

	if err != nil {
		logger.Warnf("read body err, %v\n", err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	setErr := logic.BizSetBaseInfo(baseInfo, true)

	if setErr != nil {
		c.String(http.StatusOK, err.Error())
	}
	c.String(http.StatusOK, "success")
}

// GetBaseInfo business layer get base info
func GetBaseInfo(c *gin.Context) {
	config, err := logic.BizGetBaseInfo()
	if err != nil {
		c.JSON(http.StatusOK, WithError(err))
		return
	}
	data, _ := yaml.MarshalYML(config)
	c.JSON(http.StatusOK, WithRet(string(data)))
}

// ModifyBaseInfo handle modify base info http request
func ModifyBaseInfo(c *gin.Context) {
	body := c.PostForm("content")

	baseInfo := &config.BaseInfo{}
	err := yaml.UnmarshalYML([]byte(body), baseInfo)

	if err != nil {
		logger.Warnf("read body err, %v\n", err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	setErr := logic.BizSetBaseInfo(baseInfo, false)

	if setErr != nil {
		c.String(http.StatusOK, err.Error())
	}
	c.String(http.StatusOK, "success")
}

// GetResourceDetail get resource detail with yml
func GetResourceDetail(c *gin.Context) {
	id := c.Query("resourceId")
	res, err := logic.BizGetResourceDetail(id)
	if err != nil {
		c.JSON(http.StatusOK, WithError(err))
		return
	}
	c.JSON(http.StatusOK, WithRet(res))
}

// GetMethodDetail get method detail with yml
func GetMethodDetail(c *gin.Context) {
	resourceId := c.Query("resourceId")
	methodId := c.Query("methodId")
	res, err := logic.BizGetMethodDetail(resourceId, methodId)
	if err != nil {
		c.JSON(http.StatusOK, WithError(err))
		return
	}
	c.JSON(http.StatusOK, WithRet(res))
}

// SetResourceInfo create resource
func SetResourceInfo(c *gin.Context) {
	body := c.PostForm("content")

	res := &fc.Resource{}
	err := yaml.UnmarshalYML([]byte(body), res)

	if err != nil {
		logger.Warnf("read body err, %v\n", err)
		c.JSON(http.StatusOK, WithError(err))
		return
	}

	setErr := logic.BizSetResourceInfo(res, true)

	if setErr != nil {
		c.JSON(http.StatusOK, WithError(err))
	}
	c.JSON(http.StatusOK, WithRet("Success"))
}

// ModifyResourceInfo modify resource
func ModifyResourceInfo(c *gin.Context) {
	body := c.PostForm("content")

	res := &fc.Resource{}
	err := yaml.UnmarshalYML([]byte(body), res)

	if err != nil {
		logger.Warnf("read body err, %v\n", err)
		c.JSON(http.StatusOK, WithError(err))
		return
	}

	setErr := logic.BizSetResourceInfo(res, false)

	if setErr != nil {
		c.JSON(http.StatusOK, WithError(err))
	}
	c.JSON(http.StatusOK, WithRet("Success"))
}

// DeleteResourceInfo delete resource
func DeleteResourceInfo(c *gin.Context) {
	id := c.Query("resourceId")
	err := logic.BizDeleteResourceInfo(id)

	if err != nil {
		c.JSON(http.StatusOK, WithError(err))
	}
	c.JSON(http.StatusOK, WithRet("Success"))
}

// DeleteResourceInfo delete method
func DeleteMethodInfo(c *gin.Context) {
	resourceId := c.Query("resourceId")
	methodId := c.Query("methodId")
	err := logic.BizDeleteMethodInfo(resourceId, methodId)

	if err != nil {
		c.JSON(http.StatusOK, WithError(err))
	}
	c.JSON(http.StatusOK, WithRet("Success"))
}

// SetMethodInfo create method
func SetMethodInfo(c *gin.Context) {
	body := c.PostForm("content")
	resourceId := c.Query("resourceId")

	res := &fc.Method{}
	err := yaml.UnmarshalYML([]byte(body), res)

	if err != nil {
		logger.Warnf("read body err, %v\n", err)
		c.JSON(http.StatusOK, WithError(err))
		return
	}

	setErr := logic.BizSetResourceMethod(resourceId, res, true)

	if setErr != nil {
		c.JSON(http.StatusOK, WithError(err))
	}
	c.JSON(http.StatusOK, WithRet("Success"))
}

// ModifyMethodInfo modify method
func ModifyMethodInfo(c *gin.Context) {
	body := c.PostForm("content")
	resourceId := c.Query("resourceId")

	res := &fc.Method{}
	err := yaml.UnmarshalYML([]byte(body), res)

	if err != nil {
		logger.Warnf("read body err, %v\n", err)
		c.JSON(http.StatusOK, WithError(err))
		return
	}

	setErr := logic.BizSetResourceMethod(resourceId, res, false)

	if setErr != nil {
		c.JSON(http.StatusOK, WithError(err))
	}
	c.JSON(http.StatusOK, WithRet("Success"))
}

// GetPluginGroupList get plugin group list
func GetPluginGroupList(c *gin.Context) {
	res, err := logic.BizGetPluginGroupList()
	if err != nil {
		c.JSON(http.StatusOK, WithError(err))
		return
	}
	data, _ := json.Marshal(res)
	c.JSON(http.StatusOK, WithRet(string(data)))
}

// GetPluginGroupDetail get plugin group list
func GetPluginGroupDetail(c *gin.Context) {
	name := c.Query("name")

	res, err := logic.BizGetPluginGroupDetail(name)
	if err != nil {
		c.JSON(http.StatusOK, WithError(err))
		return
	}
	c.JSON(http.StatusOK, WithRet(res))
}

// CreatePluginGroup create plugin group
func CreatePluginGroup(c *gin.Context) {
	body := c.PostForm("content")

	res := &fc.PluginsGroup{}
	err := yaml.UnmarshalYML([]byte(body), res)

	if err != nil {
		logger.Warnf("read body err, %v\n", err)
		c.JSON(http.StatusOK, WithError(err))
		return
	}

	setErr := logic.BizSetPluginGroupInfo(res, true)

	if setErr != nil {
		c.JSON(http.StatusOK, WithError(err))
	}
	c.JSON(http.StatusOK, WithRet("Success"))
}

// ModifyPluginGroup modify plugin group
func ModifyPluginGroup(c *gin.Context) {
	body := c.PostForm("content")

	res := &fc.PluginsGroup{}
	err := yaml.UnmarshalYML([]byte(body), res)

	if err != nil {
		logger.Warnf("read body err, %v\n", err)
		c.JSON(http.StatusOK, WithError(err))
		return
	}

	setErr := logic.BizSetPluginGroupInfo(res, true)

	if setErr != nil {
		c.JSON(http.StatusOK, WithError(err))
	}
	c.JSON(http.StatusOK, WithRet("Success"))
}

// DeletePluginGroup delete plugin group
func DeletePluginGroup(c *gin.Context) {
	name := c.Query("name")
	err := logic.BizDeletePluginGroupInfo(name)

	if err != nil {
		c.JSON(http.StatusOK, WithError(err))
	}
	c.JSON(http.StatusOK, WithRet("Success"))
}
