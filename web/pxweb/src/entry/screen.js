import Vue from 'vue'
import App from '@/views/Screen.vue'
import router from '@/router/screen'
import store from '@/store/store'
import '@/styles/icons/iconfont.css'
import '@/styles/common.scss'
import '@/styles/normalize.css'
import '@/utils/directives'
import {loadding,updatePublic} from '@/utils/dialogUtils'
Vue.config.productionTip = false

import '@/element-variables.scss'
import ElementUI from 'element-ui'

Vue.mixin(loadding)
Vue.mixin(updatePublic)
Vue.use(ElementUI, {
  size: 'small'
})
new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
