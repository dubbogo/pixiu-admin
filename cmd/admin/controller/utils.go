package controller

import "github.com/dubbogo/pixiu-admin/pkg/config"

/* controller 常量与方法*/

// Version admin version
const Version = "0.1.0"
const OK = "10001"
const ERR = "10002"
const RETRY = "10003"

const ResourceId = "resourceId"
const MethodId = "methodId"


// WithError transform err to RetData
func WithError(err error) config.RetData {
	return config.RetData{ERR, err.Error()}
}

// WithRet transform data to RetData
func WithRet(data interface{}) config.RetData {
	return config.RetData{OK, data}
}
