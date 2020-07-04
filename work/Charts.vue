<template>
    <layout>
        <div id="stats">

            <div class="containe-fluid">

                <div class="row">
                    <div class="col col-12 w-100">
                        <div class="card mt-3">
                            <div class="card-body">
                                <h4>SL страховых компаний</h4>
                                <div id="chart1"></div>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="row">
                    <div class="col">
                        <div class="card my-3">
                            <div class="card-body">
                                <h4>Первичная обработка претензии</h4>
                                <div id="chart2"></div>
                            </div>
                        </div>
                    </div>

                    <div class="col">
                        <div class="card my-3">
                            <div class="card-body">
                                <h4>Вторичная обработка претензии</h4>
                                <div id="chart3"></div>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="row">
                    <div class="col">
                        <div class="card my-3">
                            <div class="card-body">
                                <h4>Жизненный цикл претензии</h4>
                                <div id="chart4"></div>
                            </div>
                        </div>
                    </div>

                    <div class="col">
                        <div class="card my-3">
                            <div class="card-body">
                                <h4>Количество обращений к ТК до завершения претензии</h4>
                                <div id="chart5"></div>
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
        },
        draw(data, id) {
            let collection = []     // массив-таблица объектов данных для chart
            let legend = []         // масcив легенд для серий данных

            if (data.lenght == 0) return

            // де-сериализация данных в массив-таблицу
            data.forEach((company, i) => {
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
            chart = am4core.create(id, am4charts.XYChart)

            chart.colors.list = colors
            chart.numberFormatter.numberFormat = "##.##"
            chart.language.locale = am4lang_ru_RU
            chart.legend = new am4charts.Legend()
            chart.data = collection
            chart.dateFormatter.dateFormat = "yyyy-MM-dd"
            chart.logo.height = -15000

            // ось дат
            let dateAxis = chart.xAxes.push(new am4charts.DateAxis())
            dateAxis.title.fontWeight = "bold"
            dateAxis.dateFormats.setKey("month", "MMMM")
            dateAxis.renderer.labels.template.fill = am4core.color("#909498")
            dateAxis.renderer.minGridDistance = 30
            dateAxis.renderer.labels.template.fontSize = 12

            // ось зла
            let valueAxis = chart.yAxes.push(new am4charts.ValueAxis())
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

            chart.legend.contentAlign = "left"
            chart.legend.fontSize = 12

            legend.forEach((item, i) => {
                let legendName = legend[i]
                let dataColNum = i

                // серии данных
                let series = chart.series.push(new am4charts.LineSeries())

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
    mounted() {
        this.initStats()
        this.fetchData()
        this.draw(this.stats.comsByIntegral, "chart1")
    }
}
</script>
