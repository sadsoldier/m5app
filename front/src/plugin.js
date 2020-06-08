export default {
    install (Vue, options) {
        Vue.mixin({
            created() {
                console.log('Hello from created hook!')
            }
        })
    }
}
