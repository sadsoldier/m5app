
<template>
    <layout>
        <div class="about">
            <h4><i class="fas fa-address-card"></i> Home</h4>

            <el-button type="primary" v-on:click="update" size="mini">Update</el-button>

            <el-row type="flex" class="row-bg">
                <el-col :span="12">
                        <div id="chart"></div>
                </el-col>
            </el-row>
        </div>
    </layout>
</template>

<style>
#chart {
  height: 30em;
}
</style>

<script>
import * as am4core from "@amcharts/amcharts4/core"
import * as am4charts from "@amcharts/amcharts4/charts"
import am4themes_animated from "@amcharts/amcharts4/themes/animated"

import Layout from './Layout.vue'

function getA() {
    var data = []
    let value = 50
    for(var i = 0; i < 300; i++){
        var date = new Date()
        date.setHours(0,0,0,0)
        date.setDate(i)
        value -= Math.round((Math.random() < 0.5 ? 1 : -1) * Math.random() * 10)
        data.push({date: date, value: value})
    }
    return data
}

export default {
    components: {
        Layout
    },
    data() {
        return {
            data: []
        }
    },
    methods: {
        update() {
            this.getData()
            this.draw()
        },
        getData() {
            this.data = []
            let value = 50
            for(var i = 0; i < 4; i++){
                var date = new Date()
                date.setHours(0,0,0,0)
                date.setDate(i)
                value -= Math.round((Math.random() < 0.5 ? 1 : -1) * Math.random() * 10)
                this.data.push({date: date, value: value})
            }
        },
        draw() {
            //let chart = am4core.create("chart-12", am4charts.XYChart);
            //chart.data = this.data

            //let dateAxis = chart.xAxes.push(new am4charts.DateAxis())
            //dateAxis.renderer.minGridDistance = 60

            //let valueAxis = chart.yAxes.push(new am4charts.ValueAxis())

            //let series = chart.series.push(new am4charts.LineSeries())
            //series.dataFields.valueY = "value"
            //series.dataFields.dateX = "date"
            //series.tooltipText = "{value}"

            //series.tooltip.pointerOrientation = "vertical";

            //chart.cursor = new am4charts.XYCursor();
            //chart.cursor.snapToSeries = series;
            //chart.cursor.xAxis = dateAxis;

            //chart.scrollbarY = new am4core.Scrollbar();
            //chart.scrollbarX = new am4core.Scrollbar();

            var chart = am4core.create("chart", am4charts.XYChart);
            chart.data = this.data

            var dateAxis = chart.xAxes.push(new am4charts.DateAxis());

            var valueAxis = chart.yAxes.push(new am4charts.ValueAxis());

            var lineSeries = chart.series.push(new am4charts.LineSeries());
            lineSeries.dataFields.valueY = "value";
            lineSeries.dataFields.dateX = "date";
            lineSeries.name = "Sales";
            lineSeries.strokeWidth = 3;

            var circleBullet = lineSeries.bullets.push(new am4charts.CircleBullet());
            circleBullet.circle.stroke = am4core.color("#fff");
            circleBullet.circle.strokeWidth = 2;
            circleBullet.tooltipText = "Tip value: [bold]{value}[/]";

            var labelBullet = lineSeries.bullets.push(new am4charts.LabelBullet());
            labelBullet.label.text = "{value}";
            labelBullet.label.dy = -20;

            chart.logo.height = -15000
        }

    },

    mounted() {
        this.getData()
        this.draw()
    }
}
</script>
