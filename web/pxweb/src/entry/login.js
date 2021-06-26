import Vue from 'vue'
import App from '@/App.vue'
import router from '@/router/router'
import store from '@/store/store'
import '@/plugins/element.js'
import '@/styles/icons/iconfont.css'
import '@/styles/common.scss'
import '@/styles/normalize.css'
import '@/permission'
import '@/utils/directives'
import {loadding,updatePublic} from '@/utils/dialogUtils'
import '@/element-variables.scss'
import ElementUI from 'element-ui'
Vue.config.productionTip = false
import Router from 'vue-router'
 
const originalPush = Router.prototype.push
Router.prototype.push = function push(location) {
  return originalPush.call(this, location).catch(err => err)
}
Vue.use(ElementUI, {
  size: 'small'
})
Vue.mixin(loadding)
Vue.mixin(updatePublic)
new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
