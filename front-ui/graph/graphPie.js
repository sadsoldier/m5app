import {Graph, Vector2} from './utils'
import _ from 'lodash'
import fp from 'lodash/fp'

class GraphPie extends Graph {
  constructor (element, options) {
    super(element, options)
    this.colors = this.options.colors || ['#ff4b0a', '#898989', '#bbb', 'orange']
    this.defaultColor = this.options.defaultColor || '#f4f4f4'
    this.angleStart = this.options.angleStart || -90
    this.angleEnd = this.options.angleEnd || 270
    this.FONT = this.options.font || '400 normal 14px Roboto'
    this.init()
    this.draw()
  }
  getSize () {
    super.getSize()

    let radius = this.radius = this.size * 0.6
    this.innerRadius = radius * 0.7

    this.center = {
      x: this.width / 2,
      y: this.width / 2
    }
  }

  set data (value) {
    let _data = _.map(value, (item, index) => _.extend({}, item, {color: this.colors[index % this.colors.length]}))
    this.animate(this._data, _data)
    this._data = _data
  }

  get data () {
    return this._data
  }

  draw (from, to, progress) {
    let {x, y} = this.center
    let {radius, innerRadius} = this
    let sum = _.reduce(to || this.data, (s, item) => s + item.value, 0)
    let sumFrom = _.reduce(from, (s, item) => s + item.value, 0)
    if (to) {
      sum = sumFrom + (sum - sumFrom) * progress
    }

    let lastAngle = this.angleStart
    let lastHeight = 0
    let lastPercent = 0
    fp.flow(
      fp.thru(arr => {
        arr = _.map(arr, (item, index, arr) => {
          let startValue = _.get(from, `[${index}].value`) || 0
          let endValue = _.get(to, `[${index}].value`) || 0
          let percent = Number((to ? (startValue + (endValue - startValue) * progress) / sum : item.value / sum)).toFixed(4)
          return {
            percent,
            color: item.color
          }
        })
        let lastIndex = arr.length - 1
        if (~lastIndex) {
          let prevPercents = _.reduce(arr, (s, item, index) => s + (index !== lastIndex ? Number(item.percent) : 0), 0)
          arr[lastIndex].percent = Number(1 - prevPercents).toFixed(4)
        }
        return arr
      }),
      fp.each(({percent, color}) => {
        let sectorAngleStart = lastAngle
        let sectorAngleEnd = lastAngle = lastAngle + percent * 360
        this.drawSector(x, y, innerRadius, radius, sectorAngleStart, sectorAngleEnd, color)
        let centerAngle = (sectorAngleStart + sectorAngleEnd) / 2
        let percentText = Math.floor(Number(percent * 100 * 100).toFixed(4)) / 100
        if (percentText < 0.5) {
          percentText = ''
        } else {
          percentText += '%'
        }
        let textWidth = this.getTextWidth(percentText, this.FONT)
        if (percent > 5) {
          textWidth += 15
        }
        let textHeight = 14
        let textRadius = Math.max(textWidth / 2, textHeight)
        let difference = Number(lastPercent) + Number(percent)
        if ((difference < 0.08) && (Math.abs(lastHeight - textRadius) < textHeight)) {
          if (lastHeight < 14) textRadius += 30
          else textRadius -= 20
        }
        lastHeight = textRadius
        lastPercent = percent
        let point = new Vector2(x, y, centerAngle, radius + 7 + textRadius)
        this.drawText(point.x, point.y, percentText, '#000', this.FONT, 'center', 'middle')
      })
    )(to || this.data)

    if (lastAngle < this.angleEnd) {
      this.drawSector(x, y, innerRadius, radius, lastAngle, this.angleEnd, this.defaultColor)
    }
  }

  animate (from, to) {
    this._animate(progress => {
      this.clear()
      this.draw(from, to, progress)
    }, 1000)
  }
}
export default GraphPie
