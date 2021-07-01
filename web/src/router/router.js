import Vue from 'vue'
import Router from 'vue-router'


Vue.use(Router)

export default new Router({
  routes: [
    // { path: '/login',
    //   name: 'login',
    //   component:() => import('@/views/Login.vue'),
    //   children: [
    //
    //   ],
    // },
    {
      path: '/',
      name: 'DashBoard',
      redirect: '/manage',
      component: () => import('@/views/dashboard/manage/index.vue'),
      children: [
        {
          path: 'manage',
          component: () => import('@/views/dashboard/manage/index.vue')
        },
        {
          path: 'Overview',
          component: () => import('@/views/dashboard/manage/Overview.vue')
        },
        {
          path: 'Mapping',
          component: () => import('@/views/dashboard/manage/Mapping.vue')
        },
        {
          path: 'Plug',
          component: () => import('@/views/dashboard/manage/Plug.vue')
        },
        {
          path: 'RateLimiter',
          component: () => import('@/views/dashboard/manage/RateLimiter.vue')
        },
      ],
    },
  ]
})
