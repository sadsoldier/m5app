<template>
    <span class="mx-2">{{ time }}</span>
</template>

<script>

export default {
    components: {
    },
    data() {
        return {
            timer: null,
            time: ""
        }
    },
    methods: {
        update() {
            var date = new Date()
            var hour = date.getHours()
            var min = date.getMinutes()
            var sec = date.getSeconds()
            this.time = hour + ":" + min + ":" + sec
        }
    },
    computed: {
    },
    created() {
        this.timer = setInterval(this.update, 1000)
    },
    mounted() {
    },
    beforeDestroy () {
        clearInterval(this.timer)
    }
}
</script>
