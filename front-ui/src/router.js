
import Vue from 'vue'
import VueRouter from 'vue-router'

import store from './store.js'

import Home from './Home.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/', name: 'login', component: Home
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
