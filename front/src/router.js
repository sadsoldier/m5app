
import Vue from 'vue'
import VueRouter from 'vue-router'

import store from './store.js'

import Counter from './Counter.vue'
import About from './About.vue'
import Stock from './Stock.vue'
import Login from './Login.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/login', name: 'login', component: Login
  },
  {
    path: '/counter', name: 'counter', component: Counter, meta: { requiresAuth: true }
  },
  {
    path: '/about', name: 'about', component: About, meta: { requiresAuth: true }
  },
  {
    path: '/', name: 'stock', component: Stock, meta: { requiresAuth: true }
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

router.beforeEach((to, from, next) => {
    if(to.matched.some((record) => { return record.meta.requiresAuth } )) {
        if (store.getters.isLogin) {
            next()
            return
        }
        next('/login')
    } else {
        next()
    }
})

export default router
