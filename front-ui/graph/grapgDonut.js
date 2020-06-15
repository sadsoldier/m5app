import * as am4charts from '@amcharts/amcharts4/charts'
import * as am4core from '@amcharts/amcharts4/core'
import _ from 'lodash'

class GraphDonut {
  constructor (element, options) {
    this.element = element
    this.options = options || {}
    this.colors = this.options.colors || ['#ff4b0a', '#898989', '#bbb', 'orange']
    this.defaultColor = this.options.defaultColor || '#f4f4f4'
    this.angleStart = this.options.angleStart || -90
    this.angleEnd = this.options.angleEnd || 270
    this.FONT = this.options.font || '400 normal 14px Roboto'
    this.chart = null
    this.total = 0
    this.draw()
  }

  set data (value) {
    let totalpercent = 0
    let _data = _.map(value, (item, index) => _.extend({}, item, {
      color: this.colors[index % this.colors.length],
      value: item.value,
      percent: (item.value / this.total * 100),
      textPercent: (item.value / this.total) * 100
    }))

    _.forEach(_data, function (item) {
      if (item.percent < 3 && item.percent !== 0) {
        item.percent = 3
      }
      totalpercent += item.percent
    })
    if (totalpercent > 100) {
      _.maxBy(_data, 'percent').percent -= (totalpercent - 100)
    }
    this._data = _data
    this.chart.validateData()
  }

  get data () {
    return this._data
  }

  draw () {
    this.chart = am4core.create(this.element, am4charts.PieChart)
    // this.chart.responsive.enabled = true
    this.chart.resizable = true

    this.chart.innerRadius = am4core.percent(45)
    this.chart.startAngle = this.angleStart
    this.chart.endAngle = this.angleEnd
    this.chart.numberFormatter.numberFormat = '#.##'
    this.chart.data = this._data || {}

    let pieSeries = this.chart.series.push(new am4charts.PieSeries())

    pieSeries.dataFields.value = 'percent'
    pieSeries.dataFields.category = 'title'

    pieSeries.alignLabels = true

    pieSeries.slices.template.cursorTooltipEnabled = false

    pieSeries.slices.template.tooltipText = ''
    pieSeries.slices.template.states.getKey('hover').properties.scale = 1
    pieSeries.slices.template.states.getKey('active').properties.shiftRadius = 0
    pieSeries.slices.template.propertyFields.fill = 'color'

    pieSeries.labels.template.text = '{textPercent}%'
    pieSeries.labels.template.bent = true

    pieSeries.ticks.template.fill = 'color'

    // let categoryAxis = this.chart.xAxes.push(new am4charts.CategoryAxis())
    // categoryAxis.renderer.grid.template.location = 0

    window.addEventListener('resize', () => {
      this.chart.validatePosition()
    })
  }

  destroy () {
    if (this.chart) {
      this.chart.dispose()
    }
  }
}

export default GraphDonut
