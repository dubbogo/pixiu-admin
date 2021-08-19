<!--Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.-->

<!--
 * @Description: 子应用菜单栏列表
 * @Author: jiangyihong
 * @Date: 2020-07-22 17:24:13
 * @LastEditors: jiangyihong
-->
<template>
  <div class="custom-menu" v-if="fromType">
    <!-- <div class="menu-title">内容管理</div>
    <div class="menu-search">
      <el-input
        placeholder="搜索"
        v-model="keyWord">
        <i slot="suffix" class="el-input__icon el-icon-search"></i>
      </el-input>
      <div class="menu-search-container" v-show="searchBox">
        <div class="menu-search-item" v-for="item in searchList"
         :key="item.id" @click="handleSearchMenuClick(item)">{{item.title}}</div>
        <div class="menu-search--empty" v-show="searchList.length == 0">暂无数据</div>
      </div>
    </div> -->
    <el-menu
      :default-active="defaultActive"
      class="custom-menu-container"
      @open="handleOpen"
      @close="handleClose">
      <div v-for="(item, index) in menuList" :key="index">
        <el-submenu v-if="item.children.length >= 0" :index="item.id">
          <template slot="title">{{item.name}}</template>
          <el-menu-item :index="String(child.id)" @click="handleMenuClick(child)"
            :data-index="child.id"
            v-for="(child, childIdx) in item.children"
            :key="childIdx">{{child.name}}</el-menu-item>
        </el-submenu>
        <!-- <el-menu-item v-if="item.children.length === 0"
          :index="item.index"
          :data-index="item.index"
          @click="handleMenuClick(item)">{{item.title}}</el-menu-item> -->
      </div>
    </el-menu>
  </div>
</template>

<script>
import {getLocationSearchObj} from '@/utils/common.js'
import {mapMutations, mapGetters} from 'vuex'
import {clearLocalStorage} from '@/utils/auth.js'
const debounce = require('lodash/debounce')
import {menuList} from '@/api/menu-config.js'
export default {
  name: 'CustomMenu',
  data () {
    return {
      defaultActive: '', // 激活菜单项index值
      menuList: menuList, // 菜单列表
      searchList: [], // 搜索到的菜单列表(子菜单)
      keyWord: '', // 搜索关键字
      searchBox: false, // 搜索内容区域，显示隐藏
      fromType: true, // true:显示原有内容， false:隐藏部分内容 是否嵌入到其他项目中，需要对应隐藏部分内容(菜单，header)
    }
  },
  watch: {
    keyWord (val) {
      if (val !== '' && val !== undefined) {
        this.searchBox = true
        this.handleSearch()
      } else {
        this.searchBox = false
      }
    }
  },
  methods: {
    ...mapMutations(['setOrgId', 'setCurrentMenu']),
    // 清除分页信息
    clearPagination () {
      clearLocalStorage('picTxtIdx')
      clearLocalStorage('picTxtSize')
    },
    // 搜索过滤菜单内容
    handleSearch: debounce(function () {
      this.$get('/menu/search', {
        keyWord: this.keyWord
      }).then(res => {
        if (res.code === 200) {
          this.searchList = res.data
        }
      }).catch(err => {
        console.log(err)
      })
    }, 200),
    // 点击搜索后的菜单，滚动定位到已选项
    handleScrollPosition () {
      let el = document.querySelector('.custom-menu-container')
      let top = el.scrollTop
      let current = document.querySelector(`[data-index=${this.defaultActive}]`)
      let y = current.getBoundingClientRect().y
      top = top + y - 60 - 56
      el.scrollTop = top
    },
    // 搜索菜单项点击
    handleSearchMenuClick (item) {
      this.keyWord = ''
      this.searchBox = false
      this.searchList = []
      this.setCurrentMenu(item)
      this.handleMenuClick(item)
      this.defaultActive = item.index
      setTimeout(() => {
        this.handleScrollPosition()
      }, 300)
    },
    // 切换页面路由
    switchRouter (item) {
      this.setCurrentMenu(item)
      this.clearPagination()
      this.$router.push({
        path: item.componentName,
        query: {
          id: new Date().getTime()
        }
      })
      let el = document.querySelector('.page-content .el-scrollbar__wrap')
      if (el) {
        el.scrollTop = 0
      }
    },
    // 点击菜单项
    handleMenuClick (item) {
      this.switchRouter(item)
    },
    // 菜单栏展开
    handleOpen(key, keyPath) {
      // console.log(key, keyPath)
    },
    // 菜单栏收起
    handleClose(key, keyPath) {
      // console.log(key, keyPath)
    },
    // 初始化菜单列表
    initMenu () {
      // 0:默认框架页内加载(左侧菜单栏也需要对应)，1外部工具打开(第三方应用程序)，
      // 2新窗口打开，3覆盖本窗口打开, 4云听后台的内容嵌入，(需要先登录,暂已改成默认管理员登录)
      // 5 提示暂未开通 6应用中路由跳转
      console.log(this.menuList)
      let findShow = this.menuList.find(item => item.children.length>0)
      if(findShow){
        this.initDefaultActive()
      }else{
        // this.$router.push({
        //   path:'/403'
        // })
      }
      return;
      // var type = current.openType
    },
    // 设置激活的菜单项
    initDefaultActive () {
      let index, current
      // 如果设置了激活某一项，则优先激活显示某个路由
      let params = getLocationSearchObj()
      if (params.menuId !== undefined) {
        let list = this.menuList
        list = list.map(item => {
          return [item].concat(item.children)
        })
        list = list.flat(Infinity)
        current = list.find(item => item.id == params.menuId)
        index = current.id
      } else {
        // 使用第一个菜单项，如果有子菜单，则用子菜单的第一激活
        for (let i = 0; i < this.menuList.length; i++) {
          let item = this.menuList[i]
          if (item.children && item.children.length > 0) {
            current = item.children[0]
            index = current.id
            break
          }
        }
      }
      if (current) {
        this.defaultActive = index
        this.switchRouter(current)
      }
    },
    init() {
      let params = getLocationSearchObj()
      this.fromType = params.from == 1 ? false : true // 判断是否隐藏左侧菜单 true是不隐藏 false的隐藏
    },
  },
  created () {
    this.initMenu() 
    this.init()
  }
}
</script>

<style lang="less" scoped>
.custom-menu {
  height: 100vh;
}
.menu-title {
  height: 60px;
  line-height: 60px;
  font-size: 15px;
  font-weight: 600;
  color: #000000;
  padding-left: 20px;
  border-bottom: 1px solid #dcdfe6;
  box-sizing: border-box;
}
.custom-menu-container {
  height: calc(100% - 60px - 56px);
  overflow-y: auto;
  overflow-x: hidden;
}
.menu-search {
  position: relative;
  height: 56px;
  padding: 12px 8px;
  box-sizing: border-box;
}
.menu-search-container {
  position: absolute;
  top: 45px;
  left: 8px;
  width: 134px;
  min-height: 178px;
  max-height: 300px;
  border-radius: 4px;
  border: 1px solid #DCDFE6;
  z-index: 200;
  background: #FFF;
  overflow-x: hidden;
  overflow-y: auto;
  box-sizing: border-box;
}
.menu-search-item {
  height: 40px;
  line-height: 40px;
  padding-left: 30px;
  min-width: auto;
  cursor: pointer;
  box-sizing: border-box;
  &:hover {
    background-color: #E6E6E6 !important;
    color: #303133 !important;
  }
}
.menu-search--empty {
  text-align: center;
  padding-top: 40px;
  color: #909399;
}
</style>
