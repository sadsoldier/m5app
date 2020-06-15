import { Graph } from './utils'
import _ from 'lodash'

const ceilToMax = num => {
  num = Math.floor(num)
  let maxS = '' + num
  let mlen = maxS.length
  return (
    Math.ceil(num / parseInt(_.padEnd('1', mlen, '0'))) *
      parseInt(_.padEnd('1', mlen, '0')) || 0
  )
}

class GraphBar extends Graph {
  constructor (element, options) {
    super(element, options)
    this.colors = this.options.colors || ['#ff4b0a', '#e4e4e4']
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
    super.getSize()

    this.center = {
      x: this.width / 2,
      y: this.width / 2
    }
  }

  set data (value) {
    this.animate(this._data, value)
    this._data = value
  }

  get data () {
    return this._data
  }

  draw (from, to, progress) {
    let data = to || this.data
    let table = {
      x: 50,
      y: this.height,
      height: this.height,
      width: this.width - 100
    }
    let floatHeight = 6
    this.drawLine(
      table.x,
      table.y,
      table.x,
      table.y - table.height,
      '#e9e9e9',
      1
    )
    this.drawLine(
      table.x,
      table.y,
      table.x + table.width,
      table.y,
      '#e9e9e9',
      1
    )

    let max = ceilToMax(_.max(_.map(_.map(data, 'values'), _.sum)))
    let maxFrom = ceilToMax(_.max(_.map(_.map(from, 'values'), _.sum)))
    if (maxFrom) {
      max = maxFrom + (max - maxFrom) * progress
    }

    const ROWS = 4
    let rowHeight = table.height / (ROWS + 1)
    let rowValueHeight = Math.floor(max / ROWS)
    _.range(ROWS).forEach(index => {
      this.drawLine(
        table.x,
        table.y - rowHeight * (index + 1),
        table.x + table.width,
        table.y - rowHeight * (index + 1),
        '#e9e9e9',
        2,
        [4, 5]
      )
      let [num, text] = this.toTextMoney(rowValueHeight * (index + 1))
      this.drawText(
        table.x - 10,
        table.y - rowHeight * (index + 1) + 3,
        num,
        '#bbc0c4',
        '500 normal 18px Roboto',
        'right',
        'bottom'
      )
      this.drawText(
        table.x - 10,
        table.y - rowHeight * (index + 1) + 3,
        `${text} ₽`,
        '#bbc0c4',
        '500 normal 12px Roboto',
        'right',
        'top'
      )
    })

    const COLUMNS = 3
    let colWidth = table.width / COLUMNS
    _.range(COLUMNS).forEach(index => {
      this.drawLine(
        table.x + colWidth * (index + 1),
        table.y,
        table.x + colWidth * (index + 1),
        table.y - table.height,
        '#e9e9e9'
      )
      let prevHeight = floatHeight
      let col = _.get(data, `[${index}]`) || {}
      let colFrom = _.get(from, `[${index}]`) || {}
      let values = _.get(col, 'values') || []
      let valuesFrom = _.get(colFrom, 'values') || []
      let sum = _.sum(values)
      let sumFrom = _.sum(valuesFrom)
      if (to) {
        sum = sumFrom + (sum - sumFrom) * progress
      }
      let labels = col.label
      _.each(values, (value, vindex) => {
        if (to) {
          let valueFrom = valuesFrom[vindex] || 0
          value = valueFrom + (value - valueFrom) * progress
        }
        let height = ((table.height / (ROWS + 1)) * ROWS * value) / max
        let percent = Math.floor((value / sum) * 100)
        this.drawRect(
          table.x + colWidth * index + colWidth / 4,
          table.y - prevHeight,
          table.x + colWidth * index + (colWidth * 3) / 4,
          table.y - prevHeight - height,
          this.colors[vindex]
        )
        if (vindex === 0 && height > 24) {
          this.drawText(
            table.x + colWidth * index + colWidth / 2,
            table.y - 8,
            `${percent}%`,
            '#fff',
            '500 normal 16px Roboto',
            'center',
            'bottom'
          )
        }
        if (vindex === values.length - 1) {
          _.each(labels, (label, lindex) => {
            this.drawText(
              table.x + colWidth * index + colWidth / 2,
              table.y -
                prevHeight -
                height -
                12 -
                (values.length - lindex - 1) * 14,
              label,
              '#ababab',
              '500 normal 14px Roboto',
              'center',
              'bottom'
            )
          })
        }
        prevHeight += height
      })
    })

    this.drawText(
      table.x - 10,
      table.y,
      '0',
      '#bbc0c4',
      '500 normal 12px Roboto',
      'right',
      'bottom'
    )
  }

  animate (from, to) {
    this._animate(progress => {
      this.clear()
      this.draw(from, to, progress)
    }, 1000)
  }
}

export default GraphBar
