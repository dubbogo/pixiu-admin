package logic

import (
	"errors"
	gxetcd "github.com/dubbogo/gost/database/kv/etcd/v3"
	"strconv"
)

import (
	fc "github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	perrors "github.com/pkg/errors"
)

import (
	"github.com/dubbogo/pixiu-admin/pkg/common/yaml"
	"github.com/dubbogo/pixiu-admin/pkg/config"
	"github.com/dubbogo/pixiu-admin/pkg/logger"
)

const Base = "base"
const Resources = "Resources"
const ResourceId = "ResourceId"
const MethodId = "MethodId"
const PluginGroup = "pluginGroup"
const ErrID = -1

// BizSetBaseInfo business layer create base info
func BizSetBaseInfo(info *config.BaseInfo, created bool) error {
	// validate the api config

	data, _ := yaml.MarshalYML(info)

	if created {
		setErr := config.Client.Create(getRootPath(Base), string(data))

		if setErr != nil {
			logger.Warnf("update etcd error, %v\n", setErr)
			return perrors.WithMessage(setErr, "BizSetBaseInfo error")
		}
	} else {
		setErr := config.Client.Update(getRootPath(Base), string(data))

		if setErr != nil {
			logger.Warnf("update etcd error, %v\n", setErr)
			return perrors.WithMessage(setErr, "BizSetBaseInfo error")
		}
	}

	return nil
}

// BizGetBaseInfo business layer get base info
func BizGetBaseInfo() (*config.BaseInfo, error) {
	content, err := config.Client.Get(getRootPath(Base))
	if err != nil {
		logger.Errorf("GetBaseInfo err, %v\n", err)
		return nil, perrors.WithMessage(err, "BizGetBaseInfo error")
	}
	data := &config.BaseInfo{}
	_ = yaml.UnmarshalYML([]byte(content), data)

	return data, nil
}

// BizGetResourceDetail business layer get resource detail
func BizGetResourceDetail(id string) (string, error) {
	key := getResourceKey(id)
	detail, err := config.Client.Get(key)
	if err != nil {
		logger.Errorf("BizGetResourceDetail err, %v\n", err)
		return "", perrors.WithMessage(err, "BizGetResourceDetail error")
	}
	return detail, nil
}

// BizGetMethodDetail business layer get method detail
func BizGetMethodDetail(resourceId string, methodId string) (string, error) {
	key := getMethodKey(resourceId, methodId)
	detail, err := config.Client.Get(key)
	if err != nil {
		logger.Errorf("BizGetResourceDetail err, %v\n", err)
		return "", perrors.WithMessage(err, "BizGetResourceDetail error")
	}
	return detail, nil
}

// BizGetResourceList business layer get resource list
func BizGetResourceList() ([]fc.Resource, error) {
	_, vList, err := config.Client.GetChildrenKVList(getRootPath(Resources))
	if err != nil {
		logger.Errorf("GetResourceList err, %v\n", err)
		return nil, perrors.WithMessage(err, "BizGetResourceList error")
	}

	var ret []fc.Resource
	for _, v := range vList {
		res := &fc.Resource{}
		err := yaml.UnmarshalYML([]byte(v), res)
		if err != nil {
			logger.Errorf("UnmarshalYML err, %v\n", err)
		}
		ret = append(ret, *res)
	}

	return ret, nil
}

// BizGetMethodList business layer get method list
func BizGetMethodList(path string) ([]fc.Method, error) {
	key := getResourceMethodPrefixKey(path)

	_, vList, err := config.Client.GetChildrenKVList(key)
	if err != nil {
		logger.Errorf("GetResourceList err, %v\n", err)
		return nil, perrors.WithMessage(err, "BizGetResourceList error")
	}

	var ret []fc.Method
	for _, v := range vList {
		res := &fc.Method{}
		err := yaml.UnmarshalYML([]byte(v), res)
		if err != nil {
			logger.Errorf("UnmarshalYML err, %v\n", err)
		}
		ret = append(ret, *res)
	}

	return ret, nil
}

// BizGetPluginGroupList business layer get plugin group list
func BizGetPluginGroupList() ([]fc.PluginsGroup, error) {
	key := getPluginGroupPrefixKey()

	_, vList, err := config.Client.GetChildrenKVList(key)
	if err != nil {
		logger.Errorf("GetResourceList err, %v\n", err)
		return nil, perrors.WithMessage(err, "BizGetResourceList error")
	}

	var ret []fc.PluginsGroup
	for _, v := range vList {
		res := &fc.PluginsGroup{}
		err := yaml.UnmarshalYML([]byte(v), res)
		if err != nil {
			logger.Errorf("UnmarshalYML err, %v\n", err)
		}
		ret = append(ret, *res)
	}

	return ret, nil
}

// BizGetPluginGroupDetail business layer get plugin group detail
func BizGetPluginGroupDetail(name string) (string, error) {
	key := getPluginGroupKey(name)
	detail, err := config.Client.Get(key)
	if err != nil {
		logger.Errorf("BizGetResourceDetail err, %v\n", err)
		return "", perrors.WithMessage(err, "BizGetResourceDetail error")
	}
	return detail, nil
}

func getResourceKey(path string) string {
	return getRootPath(Resources) + "/" + path
}

func getPluginGroupKey(name string) string {
	return getPluginGroupPrefixKey() + "/" + name
}

func getPluginGroupPrefixKey() string {
	return getRootPath(PluginGroup)
}

func getResourceMethodPrefixKey(path string) string {
	return getResourceKey(path) + "/" + "Method"
}

func getMethodKey(path string, method string) string {
	return getResourceMethodPrefixKey(path) + "/" + method
}

// BizSetResourceInfo business layer create resource
func BizSetResourceInfo(res *fc.Resource, created bool) error {

	// 备份 method
	methods := res.Methods
	res.Methods = nil
	// 创建 resource

	if created {
		// 填充 id
		res.Id = getResourceId()
		if res.Id == ErrID {
			logger.Warnf("can't get id from etcd")
			return perrors.New("BizSetResourceInfo error can't get id from etcd")
		}
		data, _ := yaml.MarshalYML(res)

		setErr := config.Client.Create(getResourceKey(strconv.Itoa(res.Id)), string(data))

		if setErr != nil {
			logger.Warnf("Create etcd error, %v\n", setErr)
			return perrors.WithMessage(setErr, "BizSetResourceInfo error")
		}
	} else {
		data, _ := yaml.MarshalYML(res)
		setErr := config.Client.Update(getResourceKey(strconv.Itoa(res.Id)), string(data))

		if setErr != nil {
			logger.Warnf("update etcd error, %v\n", setErr)
			return perrors.WithMessage(setErr, "BizSetResourceInfo error")
		}
	}

	// 创建 methods
	BizCreateResourceMethod(strconv.Itoa(res.Id), methods)

	return nil
}

// BizSetPluginGroupInfo create plugin group
func BizSetPluginGroupInfo(res *fc.PluginsGroup, created bool) error {

	data, _ := yaml.MarshalYML(res)
	if created {
		setErr := config.Client.Create(getPluginGroupKey(res.GroupName), string(data))

		if setErr != nil {
			logger.Warnf("create etcd error, %v\n", setErr)
			return perrors.WithMessage(setErr, "BizSetPluginGroupInfo error")
		}
	} else {
		setErr := config.Client.Update(getPluginGroupKey(res.GroupName), string(data))

		if setErr != nil {
			logger.Warnf("update etcd error, %v\n", setErr)
			return perrors.WithMessage(setErr, "BizSetPluginGroupInfo error")
		}
	}

	return nil
}

// BizDeleteResourceInfo business layer delete resource
func BizDeleteResourceInfo(id string) error {
	key := getResourceKey(id)
	err := config.Client.Delete(key)
	if err != nil {
		logger.Warnf("BizDeleteResourceInfo, %v\n", err)
		return perrors.WithMessage(err, "BizDeleteResourceInfo error")
	}
	return nil
}

// BizDeleteMethodInfo business layer delete method
func BizDeleteMethodInfo(resourceId string, methodId string) error {
	key := getMethodKey(resourceId, methodId)
	err := config.Client.Delete(key)
	if err != nil {
		logger.Warnf("BizDeleteMethodInfo, %v\n", err)
		return perrors.WithMessage(err, "BizDeleteMethodInfo error")
	}
	return nil
}

// BizDeletePluginGroupInfo business layer delete plugin group
func BizDeletePluginGroupInfo(name string) error {
	key := getPluginGroupKey(name)
	err := config.Client.Delete(key)
	if err != nil {
		logger.Warnf("BizDeletePluginGroupInfo, %v\n", err)
		return perrors.WithMessage(err, "BizDeletePluginGroupInfo error")
	}
	return nil
}

// BizCreateResourceMethod batch create method below specific path
func BizCreateResourceMethod(resourceId string, methods []fc.Method) error {

	if len(methods) == 0 {
		return nil
	}

	var kList, vList []string

	for _, method := range methods {
		method.Id = getMethodId()
		if method.Id == ErrID {
			logger.Warnf("can't get id from etcd")
			continue
		}
		kList = append(kList, getMethodKey(resourceId, strconv.Itoa(method.Id)))
		data, _ := yaml.MarshalYML(method)
		vList = append(vList, string(data))
	}

	err := config.Client.BatchCreate(kList, vList)
	if err != nil {
		logger.Warnf("update etcd error, %v\n", err)
		return perrors.WithMessage(err, "BizCreateResourceMethod error")
	}
	return nil
}

// BizSetResourceMethod batch create method below specific path
func BizSetResourceMethod(resourceId string, method *fc.Method, created bool) error {

	if created {

		method.Id = getMethodId()
		key := getMethodKey(resourceId, strconv.Itoa(method.Id))

		if method.Id == ErrID {
			logger.Warnf("can't get id from etcd")
			return perrors.New("BizSetResourceMethod error can't get id from etcd")
		}
		data, _ := yaml.MarshalYML(method)

		err := config.Client.Create(key, string(data))
		if err != nil {
			logger.Warnf("BizSetResourceMethod etcd error, %v\n", err)
			return perrors.WithMessage(err, "BizSetResourceMethod error")
		}
	} else {
		data, _ := yaml.MarshalYML(method)
		key := getMethodKey(resourceId, strconv.Itoa(method.Id))
		err := config.Client.Update(key, string(data))
		if err != nil {
			logger.Warnf("BizSetResourceMethod etcd error, %v\n", err)
			return perrors.WithMessage(err, "BizSetResourceMethod error")
		}
	}

	return nil
}

func getResourceId() int {
	return loopGetId(getRootPath(ResourceId))
}

func getMethodId() int {
	return loopGetId(getRootPath(MethodId))
}

func loopGetId(k string) int {

	for true {

		rawClient := config.Client.GetRawClient()
		if rawClient == nil {
			logger.Error("getResourceId etcd client is null")
			return ErrID
		}

		resp, err := rawClient.Get(config.Client.GetCtx(), k)
		if err != nil {
			return ErrID
		}

		var val string
		var rev int64
		if len(resp.Kvs) != 0 {
			val = string(resp.Kvs[0].Value)
			rev = resp.Kvs[0].ModRevision
		} else {
			val = "0"
			rev = 0
		}
		id, err := strconv.Atoi(val)

		if err != nil {
			logger.Error("getResourceId Atoi error, %v\n", err)
			return ErrID
		}

		id += 1

		if rev == 0 {
			err = config.Client.Create(k, strconv.Itoa(id))
		} else {
			err = config.Client.UpdateWithRev(k, strconv.Itoa(id), rev)
		}

		if err != nil {
			if !errors.Is(err, gxetcd.ErrCompareFail) {
				logger.Error("getResourceId UpdateWithRev error, %v\n", err)
				return ErrID
			}
			logger.Info("retry get id")
		} else {
			return id
		}
	}
	return ErrID
}

func getRootPath(key string) string {
	return config.Bootstrap.GetPath() + "/" + key
}
