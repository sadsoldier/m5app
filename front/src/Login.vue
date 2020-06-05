<template>
    <div>
        <nav class="navbar navbar-expand-sm sticky-top navbar-dark bg-dark mb-4">
            <span class="navbar-text ml-3">
                <i class="fab fa-old-republic fa-lg"></i>
                M5 Login
            </span>

            <ul class="nav justify-content-end ml-auto">
            </ul>
        </nav>

        <div class="container">
            <div class="row justify-content-center">
                <div class="col-8 col-sm-6 col-md-4 border p-4 mt-sm-5 ml-3 mr-3">

                    <form accept-charset="UTF-8" v-on:submit.prevent="login">

                        <div class="form-group">
                            <label for="username">Login name:</label>
                            <input id="username" class="form-control" type="text" v-model="username"/>
                        </div>

                        <div class="form-group">
                            <label for="password">Password:</label>
                            <input id="password" class="form-control" type="password" v-model="password"/>
                        </div>

                        <div class="text-center">
                            <button class="btn btn-primary btn-sm" type="submit">Submit</button>
                        </div>

                    </form>

                </div>
            </div>
        </div>

        <div class="container">
            <div class="row">
                <div class="col mt-5">
                    <hr class="justify-content-sm-center" />
                    <div class="text-center">
                        <small>made by <a href="http://wiki.unix7.org">oleg borodin</a></small>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import axios from 'axios'

export default {
    data() {
        return {
            username: "",
            password: ""
        };
    },
    methods: {
        login() {
            axios
                .post('/api/v1/login', {
                    username: this.username,
                    password: this.password
                })
                .then(response => {
                    if (response.data.error == false) {
                        this.$store.dispatch('login')
                        localStorage.setItem('token', response.data.result.token)
                        this.$router.push('/')
                    }
                })
                .catch(err => {})
        }
    }
}
</script>
