
<template>
    <layout>
        <div id="home">

            <h4><i class="fas fa-calculator"></i> Counter</h4>

        </div>
    </layout>
</template>

<script>
import Layout from './Layout.vue'
import axios from 'axios'

export default {
    components: {
        Layout
    },
    data() {
        return {
        };
    },
    methods: {
    },
    computed: {
    },
    created() {
    },
    mounted() {
       this.$http
                .get('/api/v1/hello')
                .then(response => {
                    if (response.data.error == false) {
                        console.log(response.data.result)
                    }
                })
                .catch(err => {})
    }
}
</script>
