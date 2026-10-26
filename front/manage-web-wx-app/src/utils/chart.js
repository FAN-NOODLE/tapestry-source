class CanvasShape {
  constructor (canvas, count, width, height) {
    console.log('CanvasShape', width, height)
    this.canvas = canvas
    this.count = count
    this.WIDTH = width
    this.HEIGHT = height
    this.scale = 0.3
    this.time = 500
    this.radiusArr = [
      {
        radius: width / 2 - 4,
        dash: [],
        fillStyle: '#ffffff',
        strokeStyle: '#FFCD02',
        lineWidth: 1
      },
      {
        radius: width / 2 * 80 / 100,
        dash: [2, 2],
        fillStyle: 'rgba(255, 245, 203, 0.21)',
        strokeStyle: 'rgba(255, 205, 2, 1)'
      },
      {
        radius: width / 2 * 60 / 100,
        dash: [2, 2],
        fillStyle: 'rgba(255, 244, 201, 0.50)',
        strokeStyle: 'rgba(255, 205, 2, 1)'
      },
      {
        radius: width / 2 * 40 / 100,
        dash: [2, 2],
        fillStyle: 'rgba(255, 243, 193, 0.80)',
        strokeStyle: 'rgba(255, 205, 2, 1)'
      },
      {
        radius: width / 2 * 20 / 100,
        dash: [2, 2],
        fillStyle: 'rgba(253, 235, 157, 1)',
        strokeStyle: 'rgba(255, 205, 2, 1)'
      }
    ]
  }

  // 初始化
  init () {
    const ratio = wx.getSystemInfoSync().pixelRatio
    this.canvas.width = this.WIDTH * ratio
    this.canvas.height = this.HEIGHT * ratio
    this.ctx = this.canvas.getContext('2d')
    this.ctx.transform(ratio, 0, 0, -ratio, this.WIDTH * ratio / 2, this.HEIGHT * ratio / 2)
    console.log('canvas初始化完成', this.canvas)
  }

  // 根据半径和角度计算点
  getPoint (v, r) {
    const x = Math.cos(Math.PI / 2 - 2 * Math.PI / this.count * v) * r
    const y = Math.sin(Math.PI / 2 - 2 * Math.PI / this.count * v) * r

    return { x, y }
  }

  // 根据点的个数计算形状的顶点
  getPointArr (r) {
    const arr = []

    for (let i = 0; i < this.count; i++) {
      arr.push(this.getPoint(i, r))
    }

    return arr
  }

  // 根据点画形状
  drawShape (ctx, { dash = [], points, fillStyle, strokeStyle, lineWidth = 1 }) {
    ctx.save()
    ctx.beginPath()
    ctx.setLineDash(dash)
    ctx.moveTo(points[0].x, points[0].y)

    points.reduce((pre, nex) => {
      if (!pre) return nex
      ctx.lineTo(nex.x, nex.y)
      return nex
    }, null)

    ctx.closePath()
    ctx.strokeStyle = strokeStyle
    ctx.fillStyle = fillStyle
    ctx.lineWidth = lineWidth
    ctx.fill()
    ctx.stroke()
    ctx.restore()
  }

  // 根据点画圆
  drawCircle (ctx, { dash = [], point, r, fillStyle, strokeStyle, lineWidth = 1 }) {
    ctx.save()
    ctx.beginPath()
    ctx.strokeStyle = strokeStyle
    ctx.fillStyle = fillStyle
    ctx.setLineDash(dash)
    ctx.lineWidth = lineWidth
    ctx.arc(point.x, point.y, r, 0, Math.PI * 2, true)
    ctx.fill()
    ctx.stroke()
    ctx.restore()
  }


  // 画动态五边形和四边形
  drawDynamicShape (ctx) {
    if (this.levelArr.some(i => i === undefined)) return
    const nowTime = new Date().getTime()
    const pointArr = this.levelArr.map((i, idx) => this.getPoint(idx, this.radiusArr[i].radius))

    ctx.save()
    ctx.scale(this.scale, this.scale)
    const config = {
      fillStyle: 'rgba(255, 205, 2, 0.21)',
      points: pointArr,
      strokeStyle: '#FFCD02',
      lineWidth: 2
    }
    this.drawShape(ctx, config)
    ctx.restore()
    if (this.scale < 1 && (nowTime - this.startTime) <= this.time) {
      this.scale = Math.min(1, 0.3 + (nowTime - this.startTime) * (1 - 0.3) / this.time)
      this.canvas.requestAnimationFrame(() => {
        this.animation()
      })
    } else {
      pointArr.forEach(i => this.drawCircle(ctx, { point: i, r: 3, fillStyle: '#ffffff', strokeStyle: '#FFCD02' }))
    }
  }

  // 动画函数
  animation () {
    this.ctx.clearRect(-this.WIDTH / 2, -this.HEIGHT / 2, this.WIDTH, this.HEIGHT)
    this.drawBackground()
    this.drawDynamicShape(this.ctx)
  }

  // 启动函数
  draw (level) {
    this.levelArr = level
    this.drawBackground()
    this.startTime = new Date().getTime()
    this.animation()
  }
}

// 四边形和五变形
export class PolygonShape extends CanvasShape {
  constructor (canvas, count, width, height) {
    super(canvas, count, width, height)
    super.init()
    this.radiusArr = this.radiusArr.map(i => ({ ...i, points: this.getPointArr(i.radius) }))
  }

  // 画背景
  drawBackground () {
    this.radiusArr.forEach(i => this.drawShape(this.ctx, i))
  }
}

// 三星雷达图
export class CircleShape extends CanvasShape {
  constructor (canvas, count, width, height) {
    super(canvas, count, width, height)
    super.init()
  }

  // 画背景
  drawBackground () {
    this.radiusArr.forEach(i => this.drawCircle(this.ctx, {...i, r: i.radius, point: {x: 0, y: 0}}))
    const pointArr = this.getPointArr(this.radiusArr[0].radius)
    pointArr.forEach(i => {
      this.ctx.beginPath()
      this.ctx.moveTo(0, 0)
      this.ctx.lineTo(i.x, i.y)
      this.ctx.strokeStyle = 'rgba(255, 205, 2, 1)'
      this.ctx.stroke()
    })
  }
}

// 简单雷达图
export class SimpleShape extends CanvasShape {
  constructor (canvas, count, width, height) {
    super(canvas, count, width, height)
    super.init()
    this.radiusArr = [
      {
        radius: width / 2 * 80 / 100,
        dash: [],
        fillStyle: '#ffffff',
        strokeStyle: '#FFCD02',
        lineWidth: 1
      },
      {
        radius: width / 2 * 60 / 100,
        dash: [2, 2],
        fillStyle: 'rgba(255, 245, 203, 0.21)',
        strokeStyle: 'rgba(255, 205, 2, 1)'
      },
      {
        radius: width / 2 * 40 / 100,
        dash: [],
        fillStyle: 'rgba(255, 205, 2, 0.20)',
        strokeStyle: 'rgba(255, 205, 2, 1)',
        lineWidth: 2
      }
    ].map(i => ({ ...i, points: this.getPointArr(i.radius) }))
  }

  // 画背景
  drawBackground () {
    console.log('this.radiusArr', this.radiusArr)
    this.radiusArr.forEach(i => this.drawShape(this.ctx, i))
  }

  // 添加文字
  drawText () {
    this.ctx.save()
    this.ctx.transform(1, 0, 0, -1, 0, 0)
    this.ctx.fillStyle = 'rgba(255, 205, 2, 1)'
    this.ctx.font = 'bold 20px AlibabaPuHuiTiH'
    this.ctx.textAlign = 'center'
    this.ctx.textBaseline = 'middle'
    this.ctx.fillText('SS', 0, 0)
    this.ctx.restore()
  }

  // 画线
  drawLine () {
    const pointS = this.getPoint(this.level, this.WIDTH / 2 * 40 / 100)
    const pointE = this.getPoint(this.level, this.WIDTH / 2 * 90 / 100)

    this.ctx.save()
    this.ctx.moveTo(pointS.x, pointS.y)
    this.ctx.lineTo(pointE.x, pointE.y)
    this.ctx.strokeStyle = 'rgba(255, 205, 2, 1)'
    this.ctx.lineWidth = 2
    this.ctx.stroke()
    this.ctx.restore()

    this.drawCircle(this.ctx, { r: 3, point: pointE, fillStyle: '#ffffff', strokeStyle: 'rgba(255, 205, 2, 1)', lineWidth: 2 })
  }

  // 启动函数
  draw (level) {
    console.log('level', level)
    this.level = level
    this.drawBackground()
    this.drawText()
    this.drawLine()
  }
}

export class CircleOneShape extends CanvasShape {
  constructor (canvas, count, width, height) {
    super(canvas, count, width, height)
    super.init()
  }

  drawBackground () {
    this.drawCircle(this.ctx, { r: this.WIDTH / 2 - 2, point: { x: 0, y: 0 }, fillStyle: 'rgba(255, 245, 203, 0.21)', strokeStyle: 'rgba(255, 205, 2, 1)', lineWidth: 2 })
    this.drawCircle(this.ctx, { r: this.WIDTH / 2 - 2 - 15, point: { x: 0, y: 0 }, fillStyle: 'rgba(255, 245, 203, 0.21)', strokeStyle: 'rgba(255, 205, 2, 1)', lineWidth: 2 })
  }

  drawText () {
    this.ctx.save()
    this.ctx.transform(1, 0, 0, -1, 0, 0)
    this.ctx.fillStyle = 'rgba(255, 205, 2, 1)'
    this.ctx.font = 'bold 20px AlibabaPuHuiTiH'
    this.ctx.textAlign = 'center'
    this.ctx.textBaseline = 'middle'
    this.ctx.fillText(this.level, 0, 0)
    this.ctx.restore()
  }

  // 启动函数
  draw (level) {
    this.level = level
    this.drawBackground()
    if (!this.level) return
    this.drawText()
  }
}

export class CircleOneLineShape extends CanvasShape {
  constructor (canvas, count, width, height) {
    super(canvas, count, width, height)
    super.init()
  }

  // 画背景
  drawBackground () {
    this.radiusArr.forEach(i => this.drawCircle(this.ctx, {...i, r: i.radius, point: {x: 0, y: 0}}))

    if (this.level.some(i => i === undefined)) return

    this.ctx.beginPath()
    this.ctx.moveTo(0, this.radiusArr[4].radius)
    this.ctx.lineTo(0, this.radiusArr[this.level[0]].radius)
    this.ctx.strokeStyle = 'rgba(255, 205, 2, 1)'
    this.ctx.stroke()
    this.drawCircle(this.ctx, { r: 3, point: { x: 0, y: this.radiusArr[this.level[0]].radius }, fillStyle: 'rgba(255, 255, 255, 1)', strokeStyle: 'rgba(255, 205, 2, 1)', lineWidth: 1 })
  }

  // 启动函数
  draw (level) {
    this.level = level
    this.drawBackground()
  }
}

export class CircleTwoShape extends CanvasShape {
  constructor (canvas, count, width, height) {
    super(canvas, count, width, height)
    super.init()
  }

  // 画背景
  drawBackground () {
    this.radiusArr.forEach(i => this.drawCircle(this.ctx, {...i, r: i.radius, point: {x: 0, y: 0}}))

    if (this.level.some(i => i === undefined)) return

    this.ctx.beginPath()
    this.ctx.moveTo(-this.radiusArr[4].radius, 0)
    this.ctx.lineTo(-this.radiusArr[this.level[0]].radius, 0)
    this.ctx.strokeStyle = 'rgba(255, 205, 2, 1)'
    this.ctx.stroke()
    this.drawCircle(this.ctx, { r: 3, point: { x: -this.radiusArr[this.level[0]].radius, y: 0 }, fillStyle: 'rgba(255, 255, 255, 1)', strokeStyle: 'rgba(255, 205, 2, 1)', lineWidth: 1 })

    this.ctx.beginPath()
    this.ctx.moveTo(this.radiusArr[4].radius, 0)
    this.ctx.lineTo(this.radiusArr[this.level[1]].radius, 0)
    this.ctx.strokeStyle = 'rgba(255, 205, 2, 1)'
    this.ctx.stroke()
    this.drawCircle(this.ctx, { r: 3, point: { x: this.radiusArr[this.level[1]].radius, y: 0 }, fillStyle: 'rgba(255, 255, 255, 1)', strokeStyle: 'rgba(255, 205, 2, 1)', lineWidth: 1 })
  }

  // 启动函数
  draw (level) {
    this.level = level
    this.drawBackground()
  }
}
