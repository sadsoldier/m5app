
import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

const store = new Vuex.Store({
    strict: false,
    state: {
        isLogin: false
    },
    mutations: {
        login(state) {
            state.isLogin = true
        },
        logout(state) {
            state.isLogin = false
        }
    },
    actions: {
         login(context) {
            context.commit('login');
        },
        logout(context) {
            context.commit('logout');
        }
    },
    modules: {
    },
    getters: {
        isLogin: (state) => {
            return state.isLogin
        }
    }
})

export default store
