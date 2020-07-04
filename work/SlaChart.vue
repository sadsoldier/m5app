
<template>
    <div v-bind:id="id" v-bind:style="style"></div>
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
    props: [ "data" ],
    data() {
        return {
            id: "bi-chart-x",
            style: {
                height: "30em"
            }
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
            let chart = am4core.create(this.id, am4charts.XYChart)

            chart.colors.step = 8
            chart.numberFormatter.numberFormat = "##.##"
            chart.legend = new am4charts.Legend()
            chart.cursor = new am4charts.XYCursor()

            chart.data = collection

            let dateAxis = new am4charts.DateAxis()
            dateAxis.title.text = "Месяц"
            dateAxis.title.fontWeight = "bold"

            let valueAxis = new am4charts.ValueAxis()
            valueAxis.title.text = "Уровень сервиса"
            valueAxis.title.fontWeight = "bold"

            chart.xAxes.push(dateAxis)
            chart.yAxes.push(valueAxis)

            //valueAxis.renderer.grid.template.disabled = true
            //valueAxis.renderer.labels.template.disabled = true
            //valueAxis.renderer.grid.template.strokeOpacity = 1;
            //valueAxis.renderer.grid.template.stroke = am4core.color("#A0CA92");
            //valueAxis.renderer.grid.template.strokeWidth = 2;

            //var range = valueAxis.axisRanges.create()
            //range.label.text = "{value}"

            valueAxis.min = 0
            //valueAxis.max = 1.2
            //valueAxis.strictMinMax = false

            //valueAxis.renderer.grid.template.disabled = true
            //valueAxis.renderer.labels.template.disabled = true

            //function createGrid(value) {
              //var range = valueAxis.axisRanges.create();
              //range.value = value;
              //range.label.text = "{value}";
            //}

            //createGrid(0.5);
            //createGrid(0.9);
            //createGrid(1.0);
            //createGrid(1.1);
            //createGrid(1.2);

            chart.dateFormatter.dateFormat = "yyyy-MM-dd"

            legend.forEach((item, i) => {
                let legendName = legend[i]
                let dataColNum = i

                console.log(i)

                let lineSeries = new am4charts.LineSeries()
                lineSeries.dataFields.valueY = dataColNum;
                lineSeries.dataFields.dateX = "label"
                lineSeries.name = legendName
                lineSeries.strokeWidth = 2
                lineSeries.tensionY = 0.92
                lineSeries.tensionX = 0.92
                //lineSeries.strokeDasharray = "3,3"


                let circleBullet = new am4charts.CircleBullet()
                circleBullet.circle.stroke = am4core.color("#fff")
                circleBullet.circle.strokeWidth = 2
                circleBullet.tooltipText = legendName + ": [bold]{" + dataColNum + "}[/]"
                lineSeries.bullets.push(circleBullet)

                let labelBullet = new am4charts.LabelBullet()
                labelBullet.label.text = "{" + dataColNum + "}"
                labelBullet.label.dy = -10
                lineSeries.bullets.push(labelBullet)

                chart.series.push(lineSeries)

                //let lineSeries = getLineSeries(legend[i], i)
                //chart.series.push(getLineSeries(legend[i], i))
            })
        }
    },
    created() {
    },
    mounted() {
        //this.id += randomString()
        this.draw()
    }
}
</script>
