import {Graph, Vector2} from './utils'

class GraphMeter extends Graph {
  constructor (element, options) {
    super(element, options)
    this.grdColors = this.options.grdColors || {first: ['#ffab59', '#ff4b0a'], second: ['#ccc', '#aaa']}
    this.angleStart = this.options.angleStart || -210
    this.angleEnd = this.options.angleEnd || 30
    this.sectorCount = this.options.sectorCount || 8
    this.init()
    this.draw()
  }
  getSize () {
    super.getSize()

    let ctx = this.ctx

    let radius = this.radius = this.size * 0.9
    let innerRadius = this.innerRadius = radius * 0.74

    let {x, y} = this.center = {
      x: this.width / 2,
      y: this.width / 2
    }

    let grd = this.gradientOrange = ctx.createRadialGradient(x, y, innerRadius, x, y, radius)
    grd.addColorStop(0, this.grdColors.first[0])
    grd.addColorStop(1, this.grdColors.first[1])

    let grd2 = this.gradientGray = ctx.createRadialGradient(x, y, innerRadius, x, y, radius)
    grd2.addColorStop(0, this.grdColors.second[0])
    grd2.addColorStop(1, this.grdColors.second[1])
  }

  set percent (value) {
    this.animate(this.percent, value)
    this._percent = value
  }

  get percent () {
    return this._percent || 0
  }

  draw (value) {
    this.drawOuterLine()
    let percent = value || this.percent / 100
    let sectorPercent = 1 / this.sectorCount
    for (let i = 0; i < this.sectorCount; i++) {
      this.drawInnerSector(i, (percent - sectorPercent * i) / sectorPercent)
    }
    this.drawGUI(percent * 100)
  }

  animate (from, to) {
    this._animate(progress => {
      this.clear()
      this.draw((from + (to - from) * progress) / 100)
    }, 1000)
  }

  drawOuterLine () {
    let radius = this.size
    let innerRadius = this.radius
    let {x, y} = this.center

    this.drawSector(x, y, innerRadius, radius, this.angleStart, this.angleEnd, '#333')
  }

  drawInnerSector (index, value) {
    const sectorAngle = Math.abs(this.angleStart - this.angleEnd) / this.sectorCount
    let {radius, innerRadius} = this
    let {x, y} = this.center
    let angleDiff = 0.5
    let sectorAngleStart = this.angleStart + sectorAngle * index + (index === 0 ? 0 : angleDiff)
    let sectorAngleEnd = sectorAngleStart + sectorAngle - (index === this.sectorCount - 1 ? 0 : angleDiff) - (index === 0 ? 0 : angleDiff)

    if (value > 0 && value < 1) {
      let delimiterAngle = sectorAngleStart + (sectorAngleEnd - sectorAngleStart) * value
      this.drawSector(x, y, innerRadius, radius, sectorAngleStart, delimiterAngle, this.gradientOrange)
      this.drawSector(x, y, innerRadius, radius, delimiterAngle, sectorAngleEnd, this.gradientGray)
    }
    if (value >= 1) {
      this.drawSector(x, y, innerRadius, radius, sectorAngleStart, sectorAngleEnd, this.gradientOrange)
    }
    if (value <= 0) {
      this.drawSector(x, y, innerRadius, radius, sectorAngleStart, sectorAngleEnd, this.gradientGray)
    }
  }

  drawGUI (percent) {
    let d = this.width * 0.02
    let color = '#ababab'
    let font = '500 normal 14px Roboto'

    let point1 = new Vector2(this.center.x, this.center.y, this.angleStart, this.size * 0.9)
    point1.length = point1.length * 0.875
    this.drawText(point1.x + d, point1.y + d, '0%', color, font)

    let point2 = new Vector2(this.center.x, this.center.y, this.angleEnd, this.size * 0.9)
    point2.length = point2.length * 0.875
    this.drawText(point2.x - d, point2.y + d, '100%', color, font, 'right')
  }
}

export default GraphMeter
