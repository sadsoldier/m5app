
<template>
    <div id="layout">

        <nav class="navbar navbar-expand-sm sticky-top navbar-dark bg-dark">
            <button class="btn m-0 p-0">
                <div class="mr-3">
                    <i class="fa fa-bars fa-lg"></i>
                </div>
            </button>

            <div class="navbar-brand">
                <i class="fab fa-old-republic fa-lg"></i> V3
            </div>

            <ul class="navbar-nav mr-auto">
                <li class="nav-item">
                    <router-link to="/" class="nav-link">Stats</router-link>
                </li>
            </ul>

            <ul class="nav justify-content-end ml-auto mr-3">
                <li class="nav-item">
                    <span class="mx-2"><clock /></span>
                </li>

                <li class="nav-item">
                    <div v-on:click="logout"><i class="fas fa-sign-out-alt fa-lg"></i></div>
                </li>
            </ul>

        </nav>

        <div class="container-fluid">
            <div class="row">
                <div class="col col-12 col-sm-12 col-md-12 col-lg-12 col-xl-12 mt-4">

                    <slot/>

                </div>
            </div>
        </div>

    </div>
</template>

<script>
import Clock from './Clock.vue'

export default {
    components: {
        Clock
    },
    methods: {
        logout() {
            delete this.$http.defaults.headers.common["Authorization"]
            localStorage.removeItem('token')
            this.$store.dispatch('logout')
            this.$router.push('/login')
        }
    }
}
</script>
