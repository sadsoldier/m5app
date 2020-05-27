
import Vue from 'vue'
import VueRouter from 'vue-router'
import Counter from './Counter.vue'
import About from './About.vue'
import Stock from './Stock.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/counter', name: 'counter', component: Counter
  },
  {
    path: '/about', name: 'about', component: About
  },
  {
    path: '/', name: 'stock', component: Stock
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
