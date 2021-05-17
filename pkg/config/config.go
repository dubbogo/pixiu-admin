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

package config

import (
	"strings"
	"time"
)

import (
	etcdv3 "github.com/dubbogo/gost/database/kv/etcd/v3"
	perrors "github.com/pkg/errors"
)

import (
	"github.com/dubbogo/pixiu-admin/pkg/common/yaml"
	"github.com/dubbogo/pixiu-admin/pkg/logger"
)

var (
	Client    *etcdv3.Client
	Bootstrap *AdminBootstrap
)

// AdminBootstrap admin bootstrap config
type AdminBootstrap struct {
	Server     ServerConfig `yaml:"server" json:"server" mapstructure:"server"`
	EtcdConfig EtcdConfig   `yaml:"etcd" json:"etcd" mapstructure:"etcd"`
}

// GetAddress get etcd server address
func (a *AdminBootstrap) GetAddress() string {
	return a.Server.Address
}

// GetPath get etcd config root path
func (a *AdminBootstrap) GetPath() string {
	return a.EtcdConfig.Path
}

// ServerConfig admin http server config
type ServerConfig struct {
	Address string `yaml:"address" json:"address" mapstructure:"address"`
}

// EtcdConfig admin etcd client config
type EtcdConfig struct {
	Address string `yaml:"address" json:"admin" mapstructure:"admin"`
	Path    string `yaml:"path" json:"path" mapstructure:"path"`
}

// BaseInfo base info
type BaseInfo struct {
	Name           string `json:"name" yaml:"name"`
	Description    string `json:"description" yaml:"description"`
	PluginFilePath string `json:"pluginFilePath" yaml:"pluginFilePath"`
}

// RetData response data
type RetData struct {
	Code string      `json:"code" yaml:"code"`
	Data interface{} `json:"data" yaml:"data"`
}

// LoadAPIConfigFromFile load config from file
func LoadAPIConfigFromFile(path string) (*AdminBootstrap, error) {
	if len(path) == 0 {
		return nil, perrors.Errorf("Config file not specified")
	}
	adminBootstrap := &AdminBootstrap{}
	err := yaml.UnmarshalYMLConfig(path, adminBootstrap)
	if err != nil {
		return nil, perrors.Errorf("unmarshalYmlConfig error %v", perrors.WithStack(err))
	}
	Bootstrap = adminBootstrap
	return adminBootstrap, nil
}

// InitEtcdClient init etcd client
func InitEtcdClient() {
	newClient, err := etcdv3.NewConfigClientWithErr(
		etcdv3.WithName(etcdv3.RegistryETCDV3Client),
		etcdv3.WithTimeout(20*time.Second),
		etcdv3.WithEndpoints(strings.Split(Bootstrap.EtcdConfig.Address, ",")...),
	)

	if err != nil {
		logger.Errorf("update etcd error, %v\n", err)
		return
	}

	Client = newClient
}

// CloseEtcdClient close etcd client
func CloseEtcdClient() {
	if Client != nil {
		Client.Close()
		Client = nil
	}
}
