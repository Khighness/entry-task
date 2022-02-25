import Vue from 'vue'
import VueRouter from 'vue-router'

import Login from '@/components/Login.vue'
import Register from '@/components/Register.vue'
import Home from '@/components/Home'
import Profile from '@/components/Profile'
import UpdatePass from '@/components/UpdatePass'
import FeedBack from '@/components/FeedBack.vue'

Vue.use(VueRouter)

const router = new VueRouter({
  routes: [
    {
      path: '/',
      redirect: '/login'
    }, {
      path: '/login',
      component: Login
    }, {
      path: '/register',
      component: Register
    }, {
      name: 'home',
      path: '/home',
      component: Home,
      redirect: '/profile',
      children: [
        { path: '/profile', component: Profile },
        { path: '/updatePass', component: UpdatePass },
        { path: '/feedback', component: FeedBack }
      ]
    }]
})

// 挂载路由导航守卫
router.beforeEach((to, from, next) => {
  // to代表将要访问的路径
  // from代表从哪个路径跳转
  // next是一个函数，表示放行
  // next()放行 next('/login')强制跳转到/login页面
  if (to.path === '/login') { return next() }
  if (to.path === '/register') { return next() }
  // 获取token
  const token = window.sessionStorage.getItem('entry-token')
  if (!token) { return next('/login') }
  next()
})

export default router
