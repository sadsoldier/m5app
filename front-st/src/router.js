
import Vue from 'vue'
import VueRouter from 'vue-router'

import store from './store.js'

import Stat from './Stat.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/', name: 'login', component: Stat
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
