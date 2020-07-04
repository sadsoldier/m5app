
<template>
    <div v-bind:id="id" v-bind:style="style"></div>
</template>

<script>

import * as am4core from "@amcharts/amcharts4/core"
import * as am4charts from "@amcharts/amcharts4/charts"
import am4themes_animated from "@amcharts/amcharts4/themes/animated"

export default {
    data() {
        return {
            id: "bi-chart",
            style: {
                height: "60em"
            },
            data: [],
        }
    },
    methods: {
        draw() {
            console.log(this.id)
        }
    },
    created() {
        this.id = "bi-chart-" + Math.random().toString(36).replace(/[^a-z]+/g, '').substr(0, 6)
    },
    mounted() {
        this.draw()
    }
}
</script>
