
<template>
    <layout>
        <div class="about">
            <h3>BI</h3>

            <el-row type="flex" class="row-bg">
                <el-col>
                    <h5>Уровень сервиса страхования</h5>
                    <sla-chart v-bind:data="slData" height="30em"/>
                </el-col>
            </el-row>

            <el-row type="flex" class="row-bg">

                <el-col :span="12">
                    <h5>Уровень A</h5>
                    <sla-chart v-bind:data="slData"/>
                </el-col>

                <el-col :span="12">
                    <h5>Уровень Б</h5>
                    <sla-chart v-bind:data="slData"/>
                </el-col>

            </el-row>

            <el-row type="flex" class="row-bg">

                <el-col :span="12">
                    <h5>Уровень C</h5>
                    <sla-chart v-bind:data="slData"/>
                </el-col>

                <el-col :span="12">
                    <h5>Уровень D</h5>
                    <sla-chart v-bind:data="slData"/>
                </el-col>

            </el-row>

        </div>
    </layout>
</template>

<script>

import Layout from './Layout.vue'
import SlaChart from './SlaChart.vue'

export default {
    components: {
        Layout,
        SlaChart
    },
    data() {
        return {
            slData: [
                {
                    comName: "Компания А",
                    data: [
                        {
                            label: "2012-09-01",
                            value: 1.1
                        },
                        {
                            label: "2012-10-01",
                            value: 1.07
                        },
                        {
                            label: "2012-11-01",
                            value: 1.171
                        },
                        {
                            label: "2012-12-01",
                            value: 0.97
                        }
                    ]
                },
                {
                    comName: "Компания Б",
                    data: [
                        {
                            label: "2012-09-01",
                            value: 1.11
                        },
                        {
                            label: "2012-10-01",
                            value: 1.05
                        },
                        {
                            label: "2012-11-01",
                            value: 0.96
                        },
                        {
                            label: "2012-12-01",
                            value: 1.16
                        }
                    ]
                },
                {
                    comName: "Компания С",
                    data: [
                        {
                            label: "2012-09-01",
                            value: 1.07
                        },
                        {
                            label: "2012-10-01",
                            value: 1.09
                        },
                        {
                            label: "2012-11-01",
                            value: 0.93
                        },
                        {
                            label: "2012-12-01",
                            value: 1.04
                        }
                    ]
                },
                {
                    comName: "Компания D",
                    data: [
                        {
                            label: "2012-09-01",
                            value: 1.06
                        },
                        {
                            label: "2012-10-01",
                            value: 0.97
                        },
                        {
                            label: "2012-11-01",
                            value: 1.1
                        },
                        {
                            label: "2012-12-01",
                            value: 1.09
                        }
                    ]
                }


            ]
        }
    }
}
</script>
