import {Graph} from './utils'
import _ from 'lodash'

class GraphRadial extends Graph {
  constructor (element, options) {
    super(element, options)
    this.grdColors = this.options.grdColors || {first: ['#ff4b0a', '#ffab59', '#ff4b0a'], second: ['#333', '#ababab', '#333']}
    this.start_angle = this.options.start_angle || -90
    this.gradation = this.options.gradation || 0.09
    this.default_color = this.options.default_color || '#f3f3f3'
    this.init()
    this.draw()
  }
  getSize () {
    super.getSize()

    let ctx = this.ctx
    this.step = this.size * this.gradation
    this.radius = this.size - this.step

    let {x, y} = this.center = {
      x: this.width / 2,
      y: this.width / 2
    }

    var grd = this.gradientOrange = ctx.createLinearGradient(x - this.size, y, x + this.size, y)
    grd.addColorStop(0, this.grdColors.first[0])
    grd.addColorStop(0.5, this.grdColors.first[1])
    grd.addColorStop(1, this.grdColors.first[2])

    var grd2 = this.gradientGray = ctx.createLinearGradient(x - this.size, y, x + this.size, y)
    grd2.addColorStop(0, this.grdColors.second[0])
    grd2.addColorStop(0.5, this.grdColors.second[1])
    grd2.addColorStop(1, this.grdColors.second[2])
  }

  set data (value) {
    this.animate(_.cloneDeep(this._data), _.cloneDeep(value))
    this._data = value
  }

  get data () {
    return this._data || 0
  }

  draw (from, to, progress = 1) {
    const STEP = this.step
    let data = to || this.data
    let total = data.total
    let totalFrom = _.get(from, 'total') || 0
    if (totalFrom) {
      total = totalFrom + (total - totalFrom) * progress
    }
    let values = _.map(data.values, (value, index) => {
      if (to) {
        let valueFrom = _.get(from, `values[${index}]`) || 0
        value = valueFrom + (value - valueFrom) * progress
      }
      if (value > total) {
        value = total
      }
      return (this.start_angle + value / total * 360) % 360
    })
    let colors = [this.gradientOrange, this.gradientGray]
    let bgColors = ['#dfdfdf', this.default_color]

    let count = 2

    let radius = this.size
    for (let i = 0; i < count; i++) {
      let outerRadius = radius
      let innerRadius = outerRadius - STEP
      this.drawSectorRounded(this.center.x, this.center.y, innerRadius, outerRadius, 0, 360, bgColors[i])
      this.drawSectorRounded(this.center.x, this.center.y, innerRadius, outerRadius, -90, values[i], colors[i])
      radius = innerRadius - STEP * 0.8
    }
  }

  animate (from, to) {
    this._animate(progress => {
      this.clear()
      this.draw(from, to, progress)
    }, 1000)
  }
}

export default GraphRadial
