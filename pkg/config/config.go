package config

import (
	etcdv3 "github.com/dubbogo/gost/database/kv/etcd/v3"
	"github.com/dubbogo/pixiu-admin/pkg/common/yaml"
	"github.com/dubbogo/pixiu-admin/pkg/logger"
	perrors "github.com/pkg/errors"
	"strings"
	"time"
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

func (a *AdminBootstrap) GetAddress() string {
	return a.Server.Address
}

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

type BaseInfo struct {
	Name           string `json:"name" yaml:"name"`
	Description    string `json:"description" yaml:"description"`
	PluginFilePath string `json:"pluginFilePath" yaml:"pluginFilePath"`
}

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

func CloseEtcdClient() {
	if Client != nil {
		Client.Close()
		Client = nil
	}
}
