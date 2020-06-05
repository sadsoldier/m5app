
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
                    <router-link to="/" class="nav-link">{{ $t('stock') }}</router-link>
                </li>

                <li class="nav-item">
                    <router-link to="/counter" class="nav-link">Counter</router-link>
                </li>

                <li class="nav-item">
                    <router-link to="/about" class="nav-link">About</router-link>
                </li>
            </ul>

            <ul class="nav justify-content-end ml-auto mr-3">
                <li class="nav-item">
                    <div v-on:click="logout"><i class="fas fa-sign-out-alt fa-lg"></i></div>
                </li>
            </ul>

        </nav>

        <div class="container-fluid">
            <div class="row">
                <div class="col col-10 col-sm-10 col-md-9 col-lg-8 mt-4 mx-auto">

                    <slot/>

                </div>
            </div>
        </div>

    </div>
</template>

<script>
export default {
    methods: {
        logout() {
            localStorage.removeItem('token')
            this.$store.dispatch('logout')
            this.$router.push('/login')
        }
    }
}
</script>
