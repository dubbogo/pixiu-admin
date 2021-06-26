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