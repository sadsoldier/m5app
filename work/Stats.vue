<template>
    <layout>
        <div id="stats">

            <div class="containe-fluid">

                <div class="row">
                    <div class="col col-12 w-100">
                        <div class="card mt-3">
                            <div class="card-body">
                                <h4>SL страховых компаний</h4>
                                <chart v-bind:data="stats.comsByFirstProcessing"/>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="row">
                    <div class="col">
                        <div class="card my-3">
                            <div class="card-body">
                                <h4>Первичная обработка претензии</h4>
                                <chart v-bind:data="stats.comsByFirstProcessing"/>
                            </div>
                        </div>
                    </div>

                    <div class="col">
                        <div class="card my-3">
                            <div class="card-body">
                                <h4>Вторичная обработка претензии</h4>
                                <chart v-bind:data="stats.comsBySecondProcessing"/>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="row">
                    <div class="col">
                        <div class="card my-3">
                            <div class="card-body">
                                <h4>Жизненный цикл претензии</h4>
                                <chart v-bind:data="stats.comsByLifecicle"/>
                            </div>
                        </div>
                    </div>

                    <div class="col">
                        <div class="card my-3">
                            <div class="card-body">
                                <h4>Количество обращений к ТК до завершения претензии</h4>
                                <chart v-bind:data="stats.comsByAppealcount"/>
                            </div>
                        </div>
                    </div>

                </div>
            </div>

        </div>
    </layout>
</template>

<script>
import Layout from './Layout.vue'
import Chart from './Chart.vue'

export default {
    components: {
        Layout,
        Chart
    },
    data() {
        return {
            stats: {}
        };
    },
    methods: {
        initStats() {
            this.stats = {
                comsByIntegral: [],
                comsByFirstProcessing: [],
                comsBySecondProcessing: [],
                comsByLifecicle: [],
                comsByAppealcount: []
            }
        },
        fetchData() {
           this.$http
                .post('/api/v1/stats/get', {
                    mspName: "DellinMSP",
                    year: "2020"
                })
                .then(response => {
                    if (response.data.error == false) {
                        console.log(response.data.result)
                        this.stats = response.data.result
                    }
                })
                .catch(err => {})
        }
    },
    created() {
    },
    mounted() {
        this.initStats()
        this.fetchData()
    }
}
</script>
