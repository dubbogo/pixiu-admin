package utils

import (
	"os"
)

import (
	"go.uber.org/zap"
)

import (
	"github.com/dubbogo/pixiu-admin/pkg/global"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateDir(dirs ...string) (err error) {
	for _, v := range dirs {
		exist, err := PathExists(v)
		if err != nil {
			return err
		}
		if !exist {
			global.LOG.Debug("create directory" + v)
			if err := os.MkdirAll(v, os.ModePerm); err != nil {
				global.LOG.Error("create directory"+v, zap.Any(" error:", err))
				return err
			}
		}
	}
	return err
}
