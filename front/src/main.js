import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import axios from 'axios'
import elementUI from 'element-ui'
import locale from 'element-ui/lib/locale/lang/en'
import 'element-ui/lib/theme-chalk/index.css'
import './plugins/element.js'
import './assets/css/global.css'
import './assets/font/iconfont.css'
import VueCropper from 'vue-cropper'

Vue.use(VueCropper)

/* 阻止启动生产信息 */
Vue.config.productionTip = false

/* 全局使用ElementUI */
Vue.use(elementUI, { locale })

/* axios */
axios.defaults.baseURL = 'http://127.0.0.1:10000/'

/* 拦截器，在header中添加token */
axios.interceptors.request.use(config => {
  config.headers.Authorization = window.sessionStorage.getItem('entry-token')
  return config
})
Vue.prototype.$http = axios

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
