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
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config/ratelimit"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
)

import (
	"github.com/dubbogo/pixiu-admin/pkg/common/yaml"
	"github.com/dubbogo/pixiu-admin/pkg/config"
	"github.com/dubbogo/pixiu-admin/pkg/logger"
	logic "github.com/dubbogo/pixiu-admin/pkg/logic"
)

// Version admin version
const Version = "0.1.0"
const OK = "10001"
const ERR = "10002"
const RETRY = "10003"

const ResourceId = "resourceId"
const MethodId = "methodId"

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
	err := r.Run(config.Bootstrap.GetAddress())
	if err != nil {
		logger.Errorf(err.Error())
	}
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
	r.POST("/config/api/resource", CreateResourceInfo)
	r.PUT("/config/api/resource", ModifyResourceInfo)
	r.DELETE("/config/api/resource", DeleteResourceInfo)

	r.GET("/config/api/resource/method/list", GetMethodList)
	r.GET("/config/api/resource/method/detail", GetMethodDetail)
	r.POST("/config/api/resource/method", CreateMethodInfo)
	r.PUT("/config/api/resource/method", ModifyMethodInfo)
	r.DELETE("/config/api/resource/method", DeleteMethodInfo)

	r.GET("/config/api/plugin_group/list", GetPluginGroupList)
	r.GET("/config/api/plugin_group/detail", GetPluginGroupDetail)
	r.POST("/config/api/plugin_group", CreatePluginGroup)
	r.PUT("/config/api/plugin_group", ModifyPluginGroup)
	r.DELETE("/config/api/plugin_group", DeletePluginGroup)

	r.GET("/config/api/plugin/ratelimit", GetPluginRatelimitDetail)
	r.POST("/config/api/plugin/ratelimit", CreatePluginRatelimit)
	r.PUT("/config/api/plugin/ratelimit", ModifyPluginRatelimit)
	r.DELETE("/config/api/plugin/ratelimit", DeletePluginRatelimit)

	return r
}

// GetBaseInfo get pixiu base info such as name,desc
func GetBaseInfo(c *gin.Context) {
	conf, err := logic.BizGetBaseInfo()
	if err != nil {
		c.JSON(http.StatusOK, WithError(err))
		return
	}
	data, _ := yaml.MarshalYML(conf)
	c.JSON(http.StatusOK, WithRet(string(data)))
}

// SetBaseInfo modify pixiu base info such as name,desc
func SetBaseInfo(c *gin.Context) {
	body := c.PostForm("content")

	baseInfo := &config.BaseInfo{}
	err := yaml.UnmarshalYML([]byte(body), baseInfo)

	if err != nil {
		logger.Warnf("read body err, %v\n", err)
		c.JSON(http.StatusOK, WithError(err))
		return
	}

	setErr := logic.BizSetBaseInfo(baseInfo, true)

	if setErr != nil {
		c.JSON(http.StatusOK, WithError(setErr))
		return
	}
	c.JSON(http.StatusOK, WithRet("success"))
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

// GetResourceDetail get resource detail with yml
func GetResourceDetail(c *gin.Context) {
	id := c.Query(ResourceId)
	res, err := logic.BizGetResourceDetail(id)
	if err != nil {
		c.JSON(http.StatusOK, WithError(err))
		return
	}
	c.JSON(http.StatusOK, WithRet(res))
}

// CreateResourceInfo create resource
func CreateResourceInfo(c *gin.Context) {
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
		c.JSON(http.StatusOK, WithError(setErr))
		return
	}
	c.JSON(http.StatusOK, WithRet("Success"))
}

// ModifyResourceInfo modify resource
func ModifyResourceInfo(c *gin.Context) {
	id := c.Query(ResourceId)
	body := c.PostForm("content")

	res := &fc.Resource{}
	err := yaml.UnmarshalYML([]byte(body), res)

	if err != nil {
		logger.Warnf("read body err, %v\n", err)
		c.JSON(http.StatusOK, WithError(err))
		return
	}

	if id != "" {
		res.ID, err = strconv.Atoi(id)
		if err != nil {
			logger.Warnf("resourceID not number err, %v\n", err)
			c.JSON(http.StatusOK, WithError(err))
			return
		}
	}

	setErr := logic.BizSetResourceInfo(res, false)

	if setErr != nil {
		c.JSON(http.StatusOK, WithError(setErr))
		return
	}
	c.JSON(http.StatusOK, WithRet("Success"))
}

// DeleteResourceInfo delete resource
func DeleteResourceInfo(c *gin.Context) {
	id := c.Query(ResourceId)
	err := logic.BizDeleteResourceInfo(id)

	if err != nil {
		c.JSON(http.StatusOK, WithError(err))
		return
	}

	c.JSON(http.StatusOK, WithRet("Success"))
}

// GetMethodList get all method list below one resource
func GetMethodList(c *gin.Context) {
	resourceId := c.Query(ResourceId)

	res, err := logic.BizGetMethodList(resourceId)
	if err != nil {
		c.JSON(http.StatusOK, WithError(err))
		return
	}
	data, _ := json.Marshal(res)
	c.JSON(http.StatusOK, WithRet(string(data)))
}

// GetMethodDetail get method detail with yml
func GetMethodDetail(c *gin.Context) {
	resourceId := c.Query(ResourceId)
	methodId := c.Query(MethodId)
	res, err := logic.BizGetMethodDetail(resourceId, methodId)
	if err != nil {
		c.JSON(http.StatusOK, WithError(err))
		return
	}
	c.JSON(http.StatusOK, WithRet(res))
}

// DeleteResourceInfo delete method
func DeleteMethodInfo(c *gin.Context) {
	resourceId := c.Query(ResourceId)
	methodId := c.Query(MethodId)
	err := logic.BizDeleteMethodInfo(resourceId, methodId)

	if err != nil {
		c.JSON(http.StatusOK, WithError(err))
		return
	}
	c.JSON(http.StatusOK, WithRet("Success"))
}

// CreateMethodInfo create method
func CreateMethodInfo(c *gin.Context) {
	body := c.PostForm("content")
	resourceId := c.Query(ResourceId)

	res := &fc.Method{}
	err := yaml.UnmarshalYML([]byte(body), res)

	if err != nil {
		logger.Warnf("read body err, %v\n", err)
		c.JSON(http.StatusOK, WithError(err))
		return
	}

	resource, err := getResourceDetail(resourceId)

	if err != nil {
		logger.Warnf("CreateMethodInfo can't query resource  err, %v\n", err)
		c.JSON(http.StatusOK, WithError(err))
		return
	}
	res.ResourcePath = resource.Path

	setErr := logic.BizSetResourceMethod(resourceId, res, true)

	if setErr != nil {
		c.JSON(http.StatusOK, WithError(setErr))
		return
	}
	c.JSON(http.StatusOK, WithRet("Success"))
}

func getResourceDetail(id string) (*fc.Resource, error) {
	res, err := logic.BizGetResourceDetail(id)
	if err != nil {
		return nil, err
	}

	resource := &fc.Resource{}
	err = yaml.UnmarshalYML([]byte(res), resource)

	if err != nil {
		return nil, err
	}

	return resource, nil
}

// ModifyMethodInfo modify method
func ModifyMethodInfo(c *gin.Context) {
	body := c.PostForm("content")
	resourceId := c.Query(ResourceId)
	methodId := c.Query(MethodId)

	res := &fc.Method{}
	err := yaml.UnmarshalYML([]byte(body), res)

	if err != nil {
		logger.Warnf("read body err, %v\n", err)
		c.JSON(http.StatusOK, WithError(err))
		return
	}

	if methodId != "" {
		res.ID, err = strconv.Atoi(methodId)
		if err != nil {
			logger.Warnf("methodID not number err, %v\n", err)
			c.JSON(http.StatusOK, WithError(err))
			return
		}
	}

	resource, err := getResourceDetail(resourceId)

	if err != nil {
		logger.Warnf("CreateMethodInfo can't query resource  err, %v\n", err)
		c.JSON(http.StatusOK, WithError(err))
		return
	}
	res.ResourcePath = resource.Path

	setErr := logic.BizSetResourceMethod(resourceId, res, false)

	if setErr != nil {
		c.JSON(http.StatusOK, WithError(setErr))
		return
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

// GetPluginGroupDetail get plugin group detail
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
		c.JSON(http.StatusOK, WithError(setErr))
		return
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

	setErr := logic.BizSetPluginGroupInfo(res, false)

	if setErr != nil {
		c.JSON(http.StatusOK, WithError(setErr))
		return
	}
	c.JSON(http.StatusOK, WithRet("Success"))
}

// DeletePluginGroup delete plugin group
func DeletePluginGroup(c *gin.Context) {
	name := c.Query("name")
	err := logic.BizDeletePluginGroupInfo(name)

	if err != nil {
		c.JSON(http.StatusOK, WithError(err))
		return
	}
	c.JSON(http.StatusOK, WithRet("Success"))
}

// GetPluginRatelimitDetail get plugin ratelimit detail
func GetPluginRatelimitDetail(c *gin.Context) {
	res, err := logic.BizGetPluginRatelimitConfig()
	if err != nil {
		c.JSON(http.StatusOK, WithError(err))
		return
	}
	c.JSON(http.StatusOK, WithRet(res))
}

// CreatePluginRatelimit create plugin ratelimit conf
func CreatePluginRatelimit(c *gin.Context) {
	body := c.PostForm("content")

	res := &ratelimit.Config{}
	err := yaml.UnmarshalYML([]byte(body), res)

	if err != nil {
		logger.Warnf("read body err, %v\n", err)
		c.JSON(http.StatusOK, WithError(err))
		return
	}

	setErr := logic.BizSetPluginRatelimitInfo(res, true)

	if setErr != nil {
		c.JSON(http.StatusOK, WithError(setErr))
		return
	}
	c.JSON(http.StatusOK, WithRet("Success"))
}

// ModifyPluginRatelimit create plugin ratelimit config
func ModifyPluginRatelimit(c *gin.Context) {
	body := c.PostForm("content")

	res := &ratelimit.Config{}
	err := yaml.UnmarshalYML([]byte(body), res)

	if err != nil {
		logger.Warnf("read body err, %v\n", err)
		c.JSON(http.StatusOK, WithError(err))
		return
	}

	setErr := logic.BizSetPluginRatelimitInfo(res, false)

	if setErr != nil {
		c.JSON(http.StatusOK, WithError(setErr))
		return
	}
	c.JSON(http.StatusOK, WithRet("Success"))
}

// BizDeletePluginRatelimit delete plugin ratelimit config
func DeletePluginRatelimit(c *gin.Context) {
	err := logic.BizDeletePluginRatelimit()

	if err != nil {
		c.JSON(http.StatusOK, WithError(err))
		return
	}
	c.JSON(http.StatusOK, WithRet("Success"))
}

// WithError transform err to RetData
func WithError(err error) config.RetData {
	return config.RetData{ERR, err.Error()}
}

// WithRet transform data to RetData
func WithRet(data interface{}) config.RetData {
	return config.RetData{OK, data}
}
