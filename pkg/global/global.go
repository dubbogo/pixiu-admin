package global

import (
	"github.com/dubbogo/pixiu-admin/pkg/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	VP     *viper.Viper
	DB     *gorm.DB
	CONFIG config.Server
	LOG    *zap.Logger
)
