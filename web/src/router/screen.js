import Vue from 'vue'
import Router from 'vue-router'


Vue.use(Router)

export default new Router({
  routes: [{
      path:'/',
      component: () => import('@/views/screen/dataView'),
      redirect: '/dataView',
      name: 'dataView',
      hidden: true,
      children: [{
          path: 'dataView',
          component: () => import('@/views/screen/dataView')
      }]
    },
  ]
})
