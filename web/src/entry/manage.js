import Vue from 'vue'
import App from '@/views/Manage.vue'
import router from '@/router/manange'
import store from '@/store/store'
import '@/plugins/element.js'
import '@/styles/icons/iconfont.css'
import '@/styles/common.scss'
import '@/styles/normalize.css'
import '@/permission'
import '@/utils/directives'
import rem from '@/utils/rem'
import {loadding,updatePublic} from '@/utils/dialogUtils'
Vue.config.productionTip = false

Vue.prototype.$rem = rem
Vue.mixin(loadding)
Vue.mixin(updatePublic)
new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
