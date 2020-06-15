import _ from 'lodash'
import raf from './raf'
import FontFaceObserver from 'fontfaceobserver'

let toRadians = deg => deg * Math.PI / 180
let distance = (x1, y1, x2, y2) => Math.sqrt(Math.pow(x1 - x2) + Math.pow(y1 - y2))

let defaultTiming = timeFraction => timeFraction
let animate = function (draw, duration, timing = defaultTiming) {
  var start = performance.now()

  raf(function animate (time) {
    var timeFraction = (time - start) / duration
    if (timeFraction < 0) timeFraction = 0
    if (timeFraction > 1) timeFraction = 1

    var progress = timing(timeFraction)

    draw(progress)

    if (timeFraction < 1) {
        raf(animate)
    }
  })
}

let animations = {
  makeEaseInOut (timing) {
    return function (timeFraction) {
      if (timeFraction < 0.5) {
        return timing(2 * timeFraction) / 2
      } else {
        return (2 - timing(2 * (1 - timeFraction))) / 2
      }
    }
  },
  circ (timeFraction) {
    return 1 - Math.sin(Math.acos(timeFraction))
  }
}

class Vector2 {
  constructor () {
    let x = 0
    let y = 0
    let angle
    let length
    let args = [].slice.call(arguments)
    if (args.length === 2) {
      angle = args[0]
      length = args[1]
    }
    if (args.length === 4) {
      x = args[0]
      y = args[1]
      angle = args[2]
      length = args[3]
    }
    this.sx = x
    this.sy = y
    this._angle = angle
    this._length = length
    this.calculate()
  }

  calculate () {
    let radians = toRadians(this._angle)
    this.x = this.sx + Math.cos(radians) * this._length
    this.y = this.sy + Math.sin(radians) * this._length
  }

  add (vector) {
    this.calculate()
    vector.calculate()
    let x = this.x + vector.x
    let y = this.y + vector.y

    let angle = Math.arc(x, y)
    let length = distance(this.cx, this.cy, x, y)
    return new Vector2(this.cx, this.cy, angle, length)
  }

  multiply (value) {
    this.length *= value
  }

  get length () {
    return this._length
  }

  set length (value) {
    this._length = value
    if (length !== value) {
      this.calculate()
    }
  }
}

let graphs = []

window.addEventListener('resize', function () {
  raf(() => {
    _.each(graphs, graph => {
      graph.getSize()
      graph.redraw()
    })
  })
})

function scaleCanvas (canvas, context, width, height) {
  // assume the device pixel ratio is 1 if the browser doesn't specify it
  const devicePixelRatio = window.devicePixelRatio || 1

  // determine the 'backing store ratio' of the canvas context
  const backingStoreRatio = (
    context.webkitBackingStorePixelRatio ||
    context.mozBackingStorePixelRatio ||
    context.msBackingStorePixelRatio ||
    context.oBackingStorePixelRatio ||
    context.backingStorePixelRatio || 1
  )

  // determine the actual ratio we want to draw at
  const ratio = devicePixelRatio / backingStoreRatio

  if (devicePixelRatio !== backingStoreRatio) {
    // set the 'real' canvas size to the higher width/height
    canvas.width = width * ratio
    canvas.height = height * ratio

    // ...then scale it back down with CSS
    canvas.style.width = width + 'px'
    canvas.style.height = height + 'px'
  } else {
    // this is a normal 1:1 device; just scale it simply
    canvas.width = width
    canvas.height = height
    canvas.style.width = ''
    canvas.style.height = ''
  }

  // scale the drawing context so everything will work at the higher ratio
  context.scale(ratio, ratio)
}

class Graph {
  constructor (element, options) {
    this.$el = element
    this.options = options
    // this.init()
    // this.draw()
  }

  init () {
    let canvas = this.canvas = this.$el
    this.ctx = canvas.getContext('2d')
    this.getSize()
    graphs.push(this)
  }

  _animate () {
    return animate.apply(this, arguments)
  }

  draw () {

  }

  scaleCanvas (width, height) {
    scaleCanvas(this.canvas, this.ctx, width, height)
  }

  getSize () {
    let {width, height} = this.$el.parentNode.getBoundingClientRect()
    this.width = width
    this.height = height

    this.scaleCanvas(width, height)

    this.center = {
      x: this.width / 2,
      y: this.height / 2
    }

    this.size = width / 2
  }

  redraw () {
    this.clear()
    this.draw()
  }

  clear () {
    let ctx = this.ctx
    ctx.fillStyle = '#fff'
    ctx.fillRect(0, 0, this.width, this.height)
  }

  getTextWidth (text, font) {
    if (font) {
      this.ctx.font = font
    }
    return this.ctx.measureText(text).width
  }

  destroy () {
    _.remove(graphs, this)
  }

  /**
   * draws text
   * @param {string} text
   * @param {string} color rgb color
   * @param {string} font bold, size, font family
   */
  drawText (x, y, text, color, font, align = 'left', baseline = 'top') {
    let ctx = this.ctx
    if (font) {
      let fontResolver = new FontFaceObserver(_.last(font.split(' ')))
      fontResolver.load(null, 10000)
        .catch(() => {})
        .finally(() => {
          drawText(ctx, ...arguments)
        })
    } else {
      drawText(ctx, ...arguments)
    }
  }

  drawLine (x1, y1, x2, y2, color, width = 1, dash) {
    let ctx = this.ctx
    let backupOpts = _.pick(ctx, ['strokeStyle', 'lineWidth'])
    let bLineDash = ctx.getLineDash()

    ctx.beginPath()
    ctx.moveTo(x1, y1)
    ctx.lineTo(x2, y2)
    ctx.closePath()
    if (color) {
      ctx.strokeStyle = color
    }
    if (width) {
      ctx.lineWidth = width
    }
    if (dash) {
      ctx.setLineDash(dash)
    }
    ctx.stroke()

    _.each(backupOpts, (value, key) => {
      ctx[key] = value
    })
    ctx.setLineDash(bLineDash)
  }

  drawRect (x1, y1, x2, y2, color) {
    let ctx = this.ctx
    let backupOpts = _.pick(ctx, ['fillStyle'])
    ctx.beginPath()
    ctx.moveTo(x1, y1)
    ctx.lineTo(x1, y2)
    ctx.lineTo(x2, y2)
    ctx.lineTo(x2, y1)
    ctx.closePath()
    if (color) {
      ctx.fillStyle = color
    }
    ctx.fill()
    _.each(backupOpts, (value, key) => {
      ctx[key] = value
    })
  }

  drawSector () {
    drawSector(this.ctx, ...arguments)
  }

  drawSectorRounded () {
    drawSectorRounded(this.ctx, ...arguments)
  }

  drawCircle () {
    drawCircle(this.ctx, ...arguments)
  }
}

let drawText = (ctx, x, y, text, color, font, align = 'left', baseline = 'top') => {
  let backupOpts = _.pick(ctx, ['fillStyle', 'textBaseline', 'textAlign'])
  if (font) {
    ctx.font = font
  }
  if (align) {
    ctx.textAlign = align
  }
  if (baseline) {
    ctx.textBaseline = baseline
  }
  if (color) {
    ctx.fillStyle = color
  }
  ctx.fillText(text, x, y)
  _.each(backupOpts, (value, key) => {
    ctx[key] = value
  })
}

let drawSector = (ctx, x, y, innerRadius, outerRadius, angleStart, angleEnd, fillColor, strokeColor) => {
  let point = new Vector2(x, y, angleStart, outerRadius)
  ctx.beginPath()
  ctx.moveTo(point.x, point.y)

  ctx.arc(x, y, outerRadius, toRadians(angleStart), toRadians(angleEnd))
  ctx.arc(x, y, innerRadius, toRadians(angleEnd), toRadians(angleStart), true)

  ctx.closePath()
  if (fillColor) {
    ctx.fillStyle = fillColor
    ctx.fill()
  }
  if (strokeColor) {
    ctx.strokeColor = strokeColor
    ctx.stroke()
  }
}

let drawSectorRounded = (ctx, x, y, innerRadius, outerRadius, angleStart, angleEnd, fillColor, strokeColor) => {
  let point = new Vector2(x, y, angleStart, innerRadius)
  ctx.beginPath()
  ctx.moveTo(point.x, point.y)

  let roundRadius = (outerRadius - innerRadius) / 2
  let minAngle = Math.asin(roundRadius / (outerRadius - roundRadius)) * 60
  let angle = angleEnd - angleStart

  let startCenter = new Vector2(x, y, angleStart, (outerRadius + innerRadius) / 2)
  let endCenter = new Vector2(x, y, angleEnd, (outerRadius + innerRadius) / 2)
  if (angle < 360 - minAngle) {
    ctx.arc(startCenter.x, startCenter.y, roundRadius, toRadians(angleStart + 180), toRadians(angleStart))
  }
  ctx.arc(x, y, outerRadius, toRadians(angleStart), toRadians(angleEnd))

  if (angle < 360 - minAngle) {
    ctx.arc(endCenter.x, endCenter.y, roundRadius, toRadians(angleEnd), toRadians(angleEnd + 180))
  }
  ctx.arc(x, y, innerRadius, toRadians(angleEnd), toRadians(angleStart), true)

  if (angle >= 360 - minAngle) {
    ctx.arc(startCenter.x, startCenter.y, roundRadius, Math.PI, Math.PI * 3)
    ctx.arc(endCenter.x, endCenter.y, roundRadius, 0, Math.PI * 2)
  }

  ctx.closePath()
  if (fillColor) {
    ctx.fillStyle = fillColor
    ctx.fill()
  }
  if (strokeColor) {
    ctx.strokeColor = strokeColor
    ctx.stroke()
  }
}

let drawCircle = (ctx, x, y, radius, fillColor, strokeColor) => {
  ctx.beginPath()
  ctx.arc(x, y, radius, 0, Math.PI * 2)
  ctx.closePath()
  if (fillColor) {
    ctx.fillStyle = fillColor
    ctx.fill()
  }
  if (strokeColor) {
    ctx.strokeColor = strokeColor
    ctx.stroke()
  }
}

export {
  toRadians,
  Vector2,
  drawSector,
  Graph,
  animate,
  animations
}
