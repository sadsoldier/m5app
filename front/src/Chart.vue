<template>
    <div>
        <div v-bind:id="chartId" v-bind:style="{ height: height }"></div>
    </div>
</template>

<script>

import * as am4core from "@amcharts/amcharts4/core"
import * as am4charts from "@amcharts/amcharts4/charts"
import am4themes_animated from "@amcharts/amcharts4/themes/animated"

import am4lang_ru_RU from "./ChartLocale.js"

export default {
    props: {
        data: {
            type: Array,
            default() {
                return []
            }
        },
        height: {
            type: String,
            default: "34em"
        }
    },
    data() {
        return {
            chartId: "chart-xxx",
            chart: {}
        }
    },
    watch: {
        data() {
            this.draw()
        }
    },
    methods: {
        setId() {
            let randomStr = Math.random().toString(36).replace(/[^a-z]+/g, '').substr(0, 6)
            this.chartId = "chart-" + randomStr
        },
        draw() {
            let collection = []     // массив-таблица объектов данных для chart
            let legend = []         // масcив легенд для серий данных

            if (this.data.lenght == 0) return

            // де-сериализация данных в массив-таблицу
            this.data.forEach((company, i) => {
                legend[i] = company.fullName

                company.slData.forEach((point, n) => {
                    let found = false
                    collection.forEach((colElem, l) => {

                        if (colElem.label == point.label) {
                            collection[l][i] = point.value
                            found = true
                        }
                    })
                    if (!found) {
                        let elem = {}
                        elem[i] = point.value
                        elem["label"] = point.label
                        collection.push(elem)
                    }
                })
            })

            // сортировка массива-таблицы по полю метка-дата
            function sortCollection(a, b) {
                if (a["label"] === b["label"]) {
                    return 0
                }
                else {
                    return (a["label"] < b["label"]) ? -1 : 1
                }
            }
            collection.sort(sortCollection)

            // массив цветов для line series
            let colors = [
                am4core.color('#fe7a42'),
                am4core.color('#3c3c3c'),
                am4core.color('#ffbda2'),
                am4core.color('#8c919c'),
                am4core.color('#d3d3d3'),
                am4core.color('#ffb89a')
            ]

            am4core.useTheme(am4themes_animated)

            // объект графика
            this.chart = am4core.create(this.chartId, am4charts.XYChart)

            this.chart.colors.list = colors
            this.chart.numberFormatter.numberFormat = "##.##"
            this.chart.language.locale = am4lang_ru_RU
            this.chart.legend = new am4charts.Legend()
            this.chart.data = collection
            this.chart.dateFormatter.dateFormat = "yyyy-MM-dd"
            this.chart.logo.height = -15000

            // ось дат
            let dateAxis = this.chart.xAxes.push(new am4charts.DateAxis())
            dateAxis.title.fontWeight = "bold"
            dateAxis.dateFormats.setKey("month", "MMMM")
            dateAxis.renderer.labels.template.fill = am4core.color("#909498")
            dateAxis.renderer.minGridDistance = 30
            dateAxis.renderer.labels.template.fontSize = 12

            // ось зла
            let valueAxis = this.chart.yAxes.push(new am4charts.ValueAxis())
            valueAxis.renderer.grid.template.strokeDasharray = "3,3"
            valueAxis.renderer.grid.template.disabled = true
            valueAxis.renderer.labels.template.disabled = true
            valueAxis.renderer.labels.template.fontSize = 12

            let valueStep = 0.05
            valueAxis.min = 0.85
            valueAxis.max = 1.15
            valueAxis.strictMinMax = true

            // линии значений
            for (let i = valueAxis.min + valueStep; i < valueAxis.max; i = i + valueStep) {
                var range = valueAxis.axisRanges.create()
                range.value = i

                range.grid.stroke = am4core.color("#909498")
                range.grid.strokeWidth = 0
                range.grid.strokeOpacity = 1
                range.grid.strokeDasharray = "3,3"

                range.label.fill = am4core.color("#909498")
                range.label.text = "{value}"
                range.label.verticalCenter = "middle"
                range.label.horizontalCenter = "left"
            }

            // диапазоны на ось значений
            let range1 = valueAxis.axisRanges.create()
            range1.value = 0.9
            range1.grid.stroke = am4core.color("#8c919c")
            range1.grid.strokeWidth = 1
            range1.grid.strokeOpacity = 1
            range1.grid.strokeDasharray = "3,3"
            range1.label.text = "Мин. SLА"
            range1.label.inside = true
            range1.label.verticalCenter = "bottom"
            range1.label.dy = 10
            range1.label.fill = am4core.color("#8c919c")

            let range2 = valueAxis.axisRanges.create()
            range2.value = 1
            range2.grid.stroke = am4core.color("#8c919c")
            range2.grid.strokeWidth = 1
            range2.grid.strokeOpacity = 1
            range2.grid.strokeDasharray = "3,3"
            range2.label.text = "Норма"
            range2.label.inside = true
            range2.label.verticalCenter = "bottom"
            range2.label.dy = 10
            range2.label.fill = am4core.color("#8c919c")

            let range3 = valueAxis.axisRanges.create()
            range3.value = 1.1
            range3.grid.stroke = am4core.color("#8c919c")
            range3.grid.strokeWidth = 1
            range3.grid.strokeOpacity = 1
            range3.grid.strokeDasharray = "3,3"
            range3.label.text = "Макс. SLА"
            range3.label.inside = true
            range3.label.verticalCenter = "bottom"
            range3.label.dy = 10
            range3.label.fill = am4core.color("#8c919c")

            this.chart.legend.contentAlign = "left"
            this.chart.legend.fontSize = 12

            legend.forEach((item, i) => {
                let legendName = legend[i]
                let dataColNum = i

                // серии данных
                let series = this.chart.series.push(new am4charts.LineSeries())

                series.dataFields.valueY = dataColNum;
                series.dataFields.dateX = "label"
                series.name = legendName
                series.strokeWidth = 2
                //series.tensionY = 0.92
                //series.tensionX = 0.92

                // пимпочка на значение
                let bullet = series.bullets.push(new am4charts.CircleBullet())
                bullet.circle.stroke = am4core.color("#fff")
                bullet.circle.strokeWidth = 2

                // метка на значение
                let label = series.bullets.push(new am4charts.LabelBullet())
                label.label.text = "{" + dataColNum + "}"
                label.label.dy = -10
                label.fontSize = 12

            })
        }
    },
    created() {
    },
    beforeMount() {
        this.setId()
    },
    mounted() {
        this.draw()
    },
    destroyed() {
    }
}
</script>
