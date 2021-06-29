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

package logic

import (
	"errors"
	"regexp"
	"strconv"
)

import (
	fc "github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config/ratelimit"
	gxetcd "github.com/dubbogo/gost/database/kv/etcd/v3"
	perrors "github.com/pkg/errors"
	clientv3 "go.etcd.io/etcd/client/v3"
)

import (
	"github.com/dubbogo/pixiu-admin/pkg/common/yaml"
	"github.com/dubbogo/pixiu-admin/pkg/config"
	"github.com/dubbogo/pixiu-admin/pkg/logger"
)

const Base = "base"
const Resources = "resources"
const Method = "method"
const ResourceId = "resourceId"
const MethodId = "methodId"
const PluginGroup = "pluginGroup"
const Plugin = "plugin"
const Filter = "filter"
const Ratelimit = "ratelimit"

const ErrID = -1

// BizGetBaseInfo get base info
func BizGetBaseInfo() (*config.BaseInfo, error) {
	content, err := config.Client.Get(getRootPath(Base))
	if err != nil {
		logger.Errorf("BizGetBaseInfo err, %v\n", err)
		return nil, perrors.WithMessage(err, "BizGetBaseInfo error")
	}
	data := &config.BaseInfo{}
	_ = yaml.UnmarshalYML([]byte(content), data)

	return data, nil
}

// BizSetBaseInfo create or modify base info
func BizSetBaseInfo(info *config.BaseInfo, created bool) error {
	// validate the api config

	data, _ := yaml.MarshalYML(info)

	if created {
		setErr := config.Client.Put(getRootPath(Base), string(data))
		if setErr != nil {
			logger.Warnf("BizSetBaseInfo create error, %v\n", setErr)
			return perrors.WithMessage(setErr, "BizSetBaseInfo error")
		}
	} else {
		setErr := config.Client.Update(getRootPath(Base), string(data))
		if setErr != nil {
			logger.Warnf("BizSetBaseInfo update error, %v\n", setErr)
			return perrors.WithMessage(setErr, "BizSetBaseInfo error")
		}
	}

	return nil
}

// BizGetResourceList get resource list
func BizGetResourceList() ([]fc.Resource, error) {
	kList, vList, err := config.Client.GetChildrenKVList(getRootPath(Resources))
	if err != nil {
		logger.Errorf("BizGetResourceList err, %v\n", err)
		return nil, perrors.WithMessage(err, "BizGetResourceList error")
	}
	var ret []fc.Resource

	for i, k := range kList {
		// only handle resource, filter method
		re := getCheckResourceRegexp()
		if m := re.Match([]byte(k)); !m {
			continue
		}
		v := vList[i]
		res := &fc.Resource{}
		err := yaml.UnmarshalYML([]byte(v), res)
		if err != nil {
			logger.Errorf("UnmarshalYML err, %v\n", err)
		}
		ret = append(ret, *res)
	}

	return ret, nil
}

// BizGetResourceDetail get resource detail
func BizGetResourceDetail(id string) (string, error) {
	key := getResourceKey(id)
	detail, err := config.Client.Get(key)
	if err != nil {
		logger.Errorf("BizGetResourceDetail err, %v\n", err)
		return "", perrors.WithMessage(err, "BizGetResourceDetail error")
	}
	return detail, nil
}

// BizSetResourceInfo create resource
func BizSetResourceInfo(res *fc.Resource, created bool) error {

	if created {
		// 备份 method
		methods := res.Methods
		res.Methods = nil

		res.ID = getResourceId()
		if res.ID == ErrID {
			logger.Warnf("can't get id from etcd")
			return perrors.New("BizSetResourceInfo error can't get id from etcd")
		}
		data, _ := yaml.MarshalYML(res)

		setErr := config.Client.Create(getResourceKey(strconv.Itoa(res.ID)), string(data))
		if setErr != nil {
			logger.Warnf("Create etcd error, %v\n", setErr)
			return perrors.WithMessage(setErr, "BizSetResourceInfo error")
		}

		for _, m := range methods {
			m.ResourcePath = res.Path
		}
		// 创建 methods
		BizBatchCreateResourceMethod(strconv.Itoa(res.ID), methods)
	} else {
		data, _ := yaml.MarshalYML(res)
		setErr := config.Client.Update(getResourceKey(strconv.Itoa(res.ID)), string(data))
		if setErr != nil {
			logger.Warnf("update etcd error, %v\n", setErr)
			return perrors.WithMessage(setErr, "BizSetResourceInfo error")
		}
	}
	return nil
}

// BizDeleteResourceInfo delete resource
func BizDeleteResourceInfo(id string) error {
	key := getResourceKey(id)
	// delete all key with prefix to delete method key
	config.Client.GetRawClient().Delete(config.Client.GetCtx(), key, clientv3.WithPrefix())
	err := config.Client.Delete(key)
	if err != nil {
		logger.Warnf("BizDeleteResourceInfo, %v\n", err)
		return perrors.WithMessage(err, "BizDeleteResourceInfo error")
	}
	return nil
}

// BizGetMethodList get method list
func BizGetMethodList(resourceId string) ([]fc.Method, error) {
	key := getResourceMethodPrefixKey(resourceId)

	_, vList, err := config.Client.GetChildrenKVList(key)
	if err != nil {
		logger.Errorf("BizGetMethodList err, %v\n", err)
		return nil, perrors.WithMessage(err, "BizGetMethodList error")
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

// BizGetMethodDetail get method detail
func BizGetMethodDetail(resourceId string, methodId string) (string, error) {
	key := getMethodKey(resourceId, methodId)
	detail, err := config.Client.Get(key)
	if err != nil {
		logger.Errorf("BizGetMethodDetail err, %v\n", err)
		return "", perrors.WithMessage(err, "BizGetMethodDetail error")
	}
	return detail, nil
}

// BizBatchCreateResourceMethod batch create method below one resource
func BizBatchCreateResourceMethod(resourceId string, methods []fc.Method) error {

	if len(methods) == 0 {
		return nil
	}

	var kList, vList []string

	for _, method := range methods {
		method.ID = getMethodId()
		if method.ID == ErrID {
			logger.Warnf("can't get id from etcd")
			continue
		}
		kList = append(kList, getMethodKey(resourceId, strconv.Itoa(method.ID)))
		data, _ := yaml.MarshalYML(method)
		vList = append(vList, string(data))
	}

	err := config.Client.BatchCreate(kList, vList)
	if err != nil {
		logger.Warnf("update etcd error, %v\n", err)
		return perrors.WithMessage(err, "BizBatchCreateResourceMethod error")
	}
	return nil
}

// BizSetResourceMethod create or update method below specific path
func BizSetResourceMethod(resourceId string, method *fc.Method, created bool) error {

	if created {
		method.ID = getMethodId()
		key := getMethodKey(resourceId, strconv.Itoa(method.ID))

		if method.ID == ErrID {
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
		key := getMethodKey(resourceId, strconv.Itoa(method.ID))
		err := config.Client.Update(key, string(data))
		if err != nil {
			logger.Warnf("BizSetResourceMethod etcd error, %v\n", err)
			return perrors.WithMessage(err, "BizSetResourceMethod error")
		}
	}

	return nil
}

// BizDeleteMethodInfo delete method
func BizDeleteMethodInfo(resourceId string, methodId string) error {
	key := getMethodKey(resourceId, methodId)
	err := config.Client.Delete(key)
	if err != nil {
		logger.Warnf("BizDeleteMethodInfo, %v\n", err)
		return perrors.WithMessage(err, "BizDeleteMethodInfo error")
	}
	return nil
}

// BizGetPluginGroupList get plugin group list
func BizGetPluginGroupList() ([]fc.PluginsGroup, error) {
	key := getPluginGroupPrefixKey()

	_, vList, err := config.Client.GetChildrenKVList(key)
	if err != nil {
		logger.Errorf("BizGetPluginGroupList err, %v\n", err)
		return nil, perrors.WithMessage(err, "BizGetPluginGroupList error")
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

// BizGetPluginGroupDetail get plugin group detail
func BizGetPluginGroupDetail(name string) (string, error) {
	key := getPluginGroupKey(name)
	detail, err := config.Client.Get(key)
	if err != nil {
		logger.Errorf("BizGetPluginGroupDetail err, %v\n", err)
		return "", perrors.WithMessage(err, "BizGetPluginGroupDetail error")
	}
	return detail, nil
}

// BizSetPluginGroupInfo create or update plugin group
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

// BizDeletePluginGroupInfo delete plugin group
func BizDeletePluginGroupInfo(name string) error {
	key := getPluginGroupKey(name)
	err := config.Client.Delete(key)
	if err != nil {
		logger.Warnf("BizDeletePluginGroupInfo, %v\n", err)
		return perrors.WithMessage(err, "BizDeletePluginGroupInfo error")
	}
	return nil
}

// BizGetPluginGroupDetail get plugin group detail
func BizGetPluginRatelimitConfig() (string, error) {
	key := getPluginRatelimitKey()
	detail, err := config.Client.Get(key)
	if err != nil {
		logger.Errorf("BizGetPluginRatelimitConfig err, %v\n", err)
		return "", perrors.WithMessage(err, "BizGetPluginRatelimitConfig error")
	}
	return detail, nil
}

// BizSetPluginGroupInfo create or update plugin group
func BizSetPluginRatelimitInfo(res *ratelimit.Config, created bool) error {

	data, _ := yaml.MarshalYML(res)
	if created {
		setErr := config.Client.Create(getPluginRatelimitKey(), string(data))
		if setErr != nil {
			logger.Warnf("create etcd error, %v\n", setErr)
			return perrors.WithMessage(setErr, "BizSetPluginRatelimitInfo error")
		}
	} else {
		setErr := config.Client.Update(getPluginRatelimitKey(), string(data))
		if setErr != nil {
			logger.Warnf("update etcd error, %v\n", setErr)
			return perrors.WithMessage(setErr, "BizSetPluginRatelimitInfo error")
		}
	}

	return nil
}

// BizDeletePluginRatelimit delete plugin ratelimit config
func BizDeletePluginRatelimit() error {
	key := getPluginRatelimitKey()
	err := config.Client.Delete(key)
	if err != nil {
		logger.Warnf("BizDeletePluginRatelimit, %v\n", err)
		return perrors.WithMessage(err, "BizDeletePluginRatelimit error")
	}
	return nil
}

func getResourceKey(path string) string {
	return getRootPath(Resources) + "/" + path
}

func getPluginRatelimitKey() string {
	return getFilterPrefixKey() + "/" + Ratelimit
}

func getPluginGroupKey(name string) string {
	return getPluginGroupPrefixKey() + "/" + name
}

func getPluginGroupPrefixKey() string {
	return getRootPath(PluginGroup)
}

func getFilterPrefixKey() string {
	return getRootPath(Filter)
}

func getResourceMethodPrefixKey(path string) string {
	return getResourceKey(path) + "/" + Method
}

func getMethodKey(path string, method string) string {
	return getResourceMethodPrefixKey(path) + "/" + method
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
			logger.Error("GetId etcd client is null")
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
			logger.Error("GetId Atoi error, %v\n", err)
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
				logger.Error("GetId UpdateWithRev error, %v\n", err)
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

func getCheckResourceRegexp() *regexp.Regexp {
	return regexp.MustCompile(".+/Resources/[^/]+/?$")
}
