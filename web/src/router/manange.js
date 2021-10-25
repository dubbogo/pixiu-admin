import Vue from 'vue'
import Router from 'vue-router'


Vue.use(Router)

export default new Router({
  routes: [{
      path: '/',
      name: 'login',
      component:() => import('@/views/Manage.vue'),
      children: [
        {
          path: 'DashBoard',
          name: 'DashBoard',
          component: () => import('@/views/dashboard/index.vue'),
        },
      ],
    }
  ]
})
