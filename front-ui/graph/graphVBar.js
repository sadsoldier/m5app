import _ from 'lodash'
import {Graph, toRadians} from './utils'

let ceilToMax = num => {
  num = Math.floor(num)
  let maxS = ('' + num)
  let mlen = maxS.length
  return Math.ceil(num / parseInt(_.padEnd('1', mlen, '0'))) * parseInt(_.padEnd('1', mlen, '0'))
}

class GraphVBar extends Graph {
  constructor (element, options) {
    super(element, options)
    this.line_height = this.options.line_height || 7
    this.line_offset = this.options.line_offset || 14
    this.line_gap = this.options.line_gap || 7
    this.axis_height = this.options.axis_height || 40
    this.left_text_offset = this.options.left_text_offset || 15
    this.numMoney = this.options.numMoney || [
      [3, 'тыс'],
      [6, 'млн'],
      [9, 'млрд']
    ]
    this.init()
    this.draw()
  }

  toTextMoney (value) {
    let text = ''
    let num = value
    let l = ('' + value).length - 1
    _.each(this.numMoney, ([len, txt], index, arr) => {
      let next = arr[index + 1]
      if (l >= len && ((next && l <= next[0]) || !next)) {
        text = txt
        num = Math.floor(value / parseInt(_.padEnd('1', len + 1, '0')))
      }
    })
    return [num, text]
  }

  getSize () {
    let {width} = this.$el.parentNode.getBoundingClientRect()
    this.width = width
    let size = _.get(this.data, 'labels', []).length
    let lineOffset = size === 1 ? this.line_offset * 1.4 : this.line_offset
    let height = _.get(this.data, 'data.length', 0) * (size * this.line_height + (size - 1) * this.line_gap + lineOffset * 2) + this.axis_height
    this.height = height

    this.scaleCanvas(width, height)

    this.rects = []
  }

  set data (value) {
    let prevData = _.cloneDeep(this._data)
    this._data = _.cloneDeep(value)
    this.getSize()
    this.animate(prevData, _.cloneDeep(value))
  }

  get data () {
    return this._data
  }

  getLabel (data, item, line, value) {
    let fn = _.get(this, 'data.labelFn')
    return fn ? (fn(data, item, line, value) || '') : ''
  }

  draw (from, to, progress) {
    this.rects = []
    let data = to || this.data
    let datasets = _.get(data, 'data') || []

    let dataFrom = from
    let datasetsFrom = _.get(dataFrom, 'data') || []

    const LEFT_TEXT_FONT = '14px Roboto'
    const LEFT_OFFSET = _.max(_.map(datasets, item => this.getTextWidth(item.title, LEFT_TEXT_FONT))) + this.left_text_offset

    let table = {
      x: LEFT_OFFSET,
      y: this.height - this.axis_height,
      height: this.height - this.axis_height,
      width: this.width - LEFT_OFFSET
    }

    this.drawLine(table.x, table.y, table.x + table.width, table.y, '#e9e9e9', 1)
    this.drawLine(table.x, table.y, table.x, table.y - table.height, '#e9e9e9', 2, [4, 5])

    let max = ceilToMax(_.max(_.map(_.map(datasets, 'values'), arrs => _.max(_.map(arrs, _.sum)))))
    let maxFrom = ceilToMax(_.max(_.map(_.map(datasetsFrom, 'values'), arrs => _.max(_.map(arrs, _.sum)))))
    if (maxFrom) {
      max = maxFrom + (max - maxFrom) * progress
    }

    let outline = _.get(data, 'outline')
    if (outline) {
      let limitMax = ceilToMax(_.max(_.map(_.map(datasets, 'limits'), arrs => _.max(_.map(arrs, _.sum)))))
      let limitMaxFrom = ceilToMax(_.max(_.map(_.map(datasetsFrom, 'limits'), arrs => _.max(_.map(arrs, _.sum))))) || 0
      if (limitMaxFrom) {
        limitMax = limitMaxFrom + (limitMax - limitMaxFrom) * progress
      }
      max = Math.max(max, limitMax)
    }

    const COLS = 5
    const COL_WIDTH = (table.width - 24) / COLS
    const COL_VALUE_HEIGHT = Math.floor((max || 1000) / COLS)
    _.range(COLS).forEach(index => {
      this.drawLine(table.x + COL_WIDTH * (index + 1), table.y, table.x + COL_WIDTH * (index + 1), table.y - table.height, '#e9e9e9', 2, [4, 5])
      let [num, text] = this.toTextMoney(COL_VALUE_HEIGHT * (index + 1))
      this.drawText(table.x + COL_WIDTH * (index + 1) + 3, table.y + 28, num, '#bbc0c4', '500 normal 18px Roboto', 'center', 'bottom')
      this.drawText(table.x + COL_WIDTH * (index + 1) + 3, table.y + 28, `${text} ₽`, '#bbc0c4', '500 normal 12px Roboto', 'center', 'top')
    })

    this.drawText(table.x, table.y + 26, '0', '#bbc0c4', '500 normal 14px Roboto', 'center', 'bottom')

    const ROW_INNER_LINES = _.max(_.map(datasets, item => _.get(item, 'values', []).length))
    const LINE_OFFSET_CORRECTED = ROW_INNER_LINES === 1 ? this.line_offset * 1.4 : this.line_offset
    const ROW_HEIGHT = ROW_INNER_LINES * this.line_height + (ROW_INNER_LINES - 1) * this.line_gap + 2 * LINE_OFFSET_CORRECTED
    _.each(datasets, (item, index) => {
      let start = {
        x: table.x,
        y: table.y - table.height + ROW_HEIGHT * index
      }
      let valuesCount = _.get(item, 'values.length')
      this.drawText(start.x - 15, start.y + ROW_HEIGHT / 2 + (valuesCount ? 3 : 0), item.title, '#393939', '14px Roboto', 'right', 'center')

      const VALUES_MAX_WIDTH = COL_WIDTH * COLS
      _.each(item.values, (line, lineIndex) => {
        let lastX = 0
        _.each(line, (valueTo, valueIndex) => {
          let value = valueTo
          let valueFrom = _.get(datasetsFrom, `[${index}].values[${lineIndex}][${valueIndex}]`) || 0
          let percent = value / max
          if (to) {
            value = valueFrom + (value - valueFrom) * progress
            percent = value / max
          }
          let width = percent * VALUES_MAX_WIDTH
          let valuePos = {
            x: start.x + lastX,
            y: start.y + LINE_OFFSET_CORRECTED + lineIndex * (this.line_height + this.line_gap)
          }
          let color = _.get(this.data, `colors[${lineIndex}][${valueIndex}]`)
          if (!outline) {
            this.drawRoundedRect(valuePos.x, valuePos.y, valuePos.x + width, valuePos.y + this.line_height, color, null, valueIndex === line.length - 1)
            this.rects.push({
              x1: valuePos.x - 3,
              y1: valuePos.y - 3,
              x2: valuePos.x + width + 3,
              y2: valuePos.y + this.line_height + 3,
              label: this.getLabel(data, item, line, valueTo)
            })
          } else {
            let limits = _.get(item, `limits[${lineIndex}]`)
            let limitMax = _.sum(limits)
            if (to) {
              let limitsFrom = _.get(datasetsFrom, `[${index}].limits[${lineIndex}]`)
              let limitMaxFrom = _.sum(limitsFrom)
              limitMax = limitMaxFrom + (limitMax - limitMaxFrom) * progress
            }
            let lastLimitX = 0
            let valueInLimit = value
            let limitPos = _.cloneDeep(valuePos)
            _.each(limits, (limit, limitIndex) => {
              let limitFrom = _.get(datasetsFrom, `[${index}].limits[${lineIndex}][${limitIndex}]`) || 0
              if (to) {
                limit = limitFrom + (limit - limitFrom) * progress
              }
              let limitWidth = limit / max * VALUES_MAX_WIDTH
              color = _.get(this.data, `colors[${lineIndex}][${limitIndex}]`)
              limitPos.x += lastLimitX
              if (valueInLimit <= limit) {
                this.drawRoundedRect(limitPos.x, limitPos.y, limitPos.x + limitWidth, limitPos.y + this.line_height, null, color, limitIndex === limits.length - 1)
                if (valueInLimit) {
                  width = valueInLimit / limit * limitWidth
                  this.drawRoundedRect(limitPos.x, limitPos.y, limitPos.x + width, limitPos.y + this.line_height, color, null, limitIndex === limits.length - 1)
                }
                if (limitIndex === limits.length - 1) {
                  this.rects.push({
                    x1: valuePos.x - 3,
                    y1: valuePos.y - 3,
                    x2: limitPos.x + limitWidth + 3,
                    y2: limitPos.y + this.line_height + 3,
                    label: this.getLabel(data, item, line, valueTo)
                  })
                }
                valueInLimit = 0
              } else {
                this.drawRoundedRect(limitPos.x, limitPos.y, limitPos.x + limitWidth, limitPos.y + this.line_height, color, null, value < limitMax && limitIndex === limits.length - 1)
                if (value > limitMax && limitIndex === limits.length - 1) {
                  limitPos.x += limitWidth
                  width = (value - limitMax) / max * VALUES_MAX_WIDTH
                  color = _.last(_.get(this.data, `colors[${lineIndex}]`))
                  this.drawRoundedRect(limitPos.x, limitPos.y, limitPos.x + width, limitPos.y + this.line_height, color, null, true)
                  this.rects.push({
                    x1: valuePos.x - 3,
                    y1: valuePos.y - 3,
                    x2: limitPos.x + width + 3,
                    y2: limitPos.y + this.line_height + 3,
                    label: this.getLabel(data, item, line, valueTo)
                  })
                }
                valueInLimit -= limit
              }
              lastLimitX += limitWidth
            })
          }
          lastX += width
        })
      })
    })
    // _.each(this.rects, rect => {
    //   this.drawRoundedRect(rect.x1, rect.y1, rect.x2, rect.y2, null, '#f00')
    // })
  }

  drawRoundedRect (x1, y1, x2, y2, color, stroke = false, rounded = false) {
    let ctx = this.ctx
    let backupOpts = _.pick(ctx, ['fillStyle', 'strokeStyle'])
    ctx.beginPath()
    ctx.moveTo(x1, y1)
    ctx.lineTo(x1, y2)
    ctx.lineTo(x2, y2)
    let r = (y2 - y1) / 2
    if (rounded) {
      ctx.arc(x2, y2 - r, r, toRadians(90), toRadians(-90), true)
    } else {
      ctx.lineTo(x2, y1)
    }
    ctx.closePath()
    if (color) {
      ctx.fillStyle = color
      ctx.fill()
    }
    if (stroke) {
      ctx.strokeStyle = stroke
      ctx.stroke()
    }
    _.each(backupOpts, (value, key) => {
      ctx[key] = value
    })
  }

  animate (from, to) {
    this._animate(progress => {
      this.clear()
      this.draw(from, to, progress)
    }, 1000)
  }
}

export default GraphVBar
