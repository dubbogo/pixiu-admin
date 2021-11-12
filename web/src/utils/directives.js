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
import Vue from 'vue';
import { getLocalStorage } from '@/utils/auth'

//权限功能点检验方法
Vue.prototype.$_has = function({rights}) {
//debugger
    let operatorInfo = getLocalStorage('operatorInfo');      
    let permission = operatorInfo.rights ? operatorInfo.rights.split(',') : [];
    let resources = [];
    let result = false;

    //提取权限数据
    if(Array.isArray(rights)){
        rights.forEach((e) => {
        resources = resources.concat([e]);
      })
    }else{
      resources = resources.concat([rights.split(',')]);
    }
    //校验权限
    resources.map((p) => {
      if(typeof p != 'string'){
        p = p.toString();
      }
      if(permission.includes(p)){
        return result = true;
      }
    })

    return result

}
//账号权限指令
Vue.directive('has',{
    bind: (el, binding) => {},
    inserted:(el, binding) => {
      
      if(!Vue.prototype.$_has(binding.value)) {
        el.parentNode.removeChild(el);
      }
    }
  })