
<template>
    <div v-bind:id="chartId" v-bind:style="{ height: height }"></div>
</template>

<script>
import * as am4core from "@amcharts/amcharts4/core"
import * as am4charts from "@amcharts/amcharts4/charts"

import am4themes_animated from "@amcharts/amcharts4/themes/animated"

function randomString() {
    return Math.random().toString(36).replace(/[^a-z]+/g, '').substr(0, 6)
}


function getLineSeries(legendName, dataColNum) {

    let lineSeries = new am4charts.LineSeries()
    lineSeries.dataFields.valueY = dataColNum;
    lineSeries.dataFields.dateX = "label"
    lineSeries.name = legendName
    lineSeries.strokeWidth = 2
    lineSeries.tensionY = 0.92
    lineSeries.tensionX = 0.92

    let circleBullet = new am4charts.CircleBullet()
    circleBullet.circle.stroke = am4core.color("#fff")
    circleBullet.circle.strokeWidth = 2
    circleBullet.tooltipText = "Value: [bold]{" + dataColNum + "}[/]"
    lineSeries.bullets.push(circleBullet)

    let labelBullet = new am4charts.LabelBullet()
    labelBullet.label.text = "{" + dataColNum + "}"
    labelBullet.label.dy = -20
    lineSeries.bullets.push(labelBullet)

    return lineSeries
}


export default {
    props: {
        data: {
            type: Object,
            default: []
        },
        height: {
            type: String,
            default: "30em"
        }
    },
    data() {
        return {
            chartId: "bi-chart",
        }
    },
    methods: {
        draw() {

            let collection = []
            let legend = []

            this.data.forEach((com, i) => {
                legend[i] = com.comName

                com.data.forEach((point, n) => {
                    let found = false
                    collection.forEach((colElem, l) => {
                        if (colElem.label == point.label) {
                            collection[l][i] = Math.round(point.value * 100) / 100
                            found = true
                        }
                    })
                    if (!found) {
                        let elem = {
                        }
                        elem[i] = Math.round(point.value * 100) / 100
                        elem["label"] = point.label
                        collection.push(elem)
                    }
                })
            })

            //console.log(collection)
            //console.log(legend)

            am4core.useTheme(am4themes_animated)
            let chart = am4core.create(this.chartId, am4charts.XYChart)



            chart.colors.step = 8
            chart.numberFormatter.numberFormat = "##.##"

            chart.legend = new am4charts.Legend()
            chart.cursor = new am4charts.XYCursor()

            chart.data = collection

            let dateAxis = new am4charts.DateAxis()
            let valueAxis = new am4charts.ValueAxis()

            dateAxis = chart.xAxes.push(dateAxis)
            dateAxis.title.text = "Месяц"
            dateAxis.title.fontWeight = "bold"
            dateAxis.dateFormats.setKey("month", "MMMM yyyy")

            valueAxis = chart.yAxes.push(valueAxis)
            valueAxis.title.text = "Уровень сервиса"
            valueAxis.title.fontWeight = "bold"

            //valueAxis.min = 0
            //valueAxis.max = 1.2
            //valueAxis.strictMinMax = false

            chart.dateFormatter.dateFormat = "yyyy-MM-dd"
            chart.logo.height = -15000

            legend.forEach((item, i) => {
                let legendName = legend[i]
                let dataColNum = i

                let lineSeries = chart.series.push(new am4charts.LineSeries())

                lineSeries.dataFields.valueY = dataColNum;
                lineSeries.dataFields.dateX = "label"
                lineSeries.name = legendName
                lineSeries.strokeWidth = 2
                lineSeries.tensionY = 0.92
                lineSeries.tensionX = 0.92

                let circleBullet = lineSeries.bullets.push(new am4charts.CircleBullet())
                circleBullet.circle.stroke = am4core.color("#fff")
                circleBullet.circle.strokeWidth = 2
                //circleBullet.tooltipText = legendName + ": [bold]{" + dataColNum + "}[/]"

                let labelBullet = lineSeries.bullets.push(new am4charts.LabelBullet())
                labelBullet.label.text = "{" + dataColNum + "}"
                labelBullet.label.dy = -10
            })

            //chart.scrollbarX = new am4charts.XYChartScrollbar()
            //chart.scrollbarX.parent = chart.bottomAxesContainer
            //chart.scrollbarX.thumb.minWidth = 50
        }
    },

    created() {
    },

    mounted() {
        this.chartId += randomString()
        this.draw()
    }
}
</script>
