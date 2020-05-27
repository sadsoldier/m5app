
import Vue from 'vue'
import VueI18n from 'vue-i18n'


import App from './App.vue'

import router from './router.js'
import store from './store'

import '@/app.scss'

import fontawesome from '@fortawesome/fontawesome-free'
import { library } from '@fortawesome/fontawesome-svg-core'
import { fas } from '@fortawesome/free-solid-svg-icons'
import { far } from '@fortawesome/free-regular-svg-icons'

library.add(fas)
library.add(far)

Vue.config.productionTip = false

Vue.use(VueI18n)

const i18n = new VueI18n({
    locale: 'ru',
    messages: {
        en: {
            "stock": "Stock"
        },
        ru: {
            "stock": "Склад"
        }
    }
})

new Vue({
    i18n,
    router,
    store,
    render: h => h(App)
}).$mount('#app')
