



<template>
    <div class="home">

        <h4><i class="fas fa-calculator"></i> Counter</h4>

        <button class="btn btn-sm btn-primary" v-on:click="dec">
            <i class="fa fa-minus"></i>
        </button>

        <span class="mx-3">
            counter: {{ counter }}
        </span>

        <button class="btn btn-sm btn-primary" v-on:click="add">
            <i class="fa fa-plus"></i>
        </button>

    </div>
</template>

<script>
export default {
    data() {
        return {
            counter: 0
        };
    },
    methods: {
        add() {
            this.counter += 1
        },
        dec() {
            this.counter -= 1
        }
    },
    computed: {
    },
    created() {
        console.log("created")
    },
    mounted() {
        console.log("mounted")
    }
}
</script>
