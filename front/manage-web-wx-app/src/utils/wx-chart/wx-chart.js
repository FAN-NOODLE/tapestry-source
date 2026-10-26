const defaultOpts = {
  devicePixelRatio: 1,
  renderer: 'canvas',
  width: 'auto',
  height: 'auto'
}

const NODE_STYLE_COMMON = {
  type: '', // 模块类型
  name: '', // 模块名称
  show: false,
  left: 'auto', // 组件容器的左侧距离。默认自适应。
  top: 'auto', // 组件容器的顶部距离。默认自适应。
  right: 'auto', // 组件容器的右侧距离。默认自适应。
  bottom: 'auto', // 组件容器的底部距离。默认自适应。
  width: 'auto', // 组件容器的宽度。默认自适应。
  height: 'auto', // 组件容器的高度。默认自适应。
  padding: [0], // 组件容器的内边距，默认0
  margin: [0], // 组件容器的外边距，默认0
  backgroundColor: 'white', // 组件容器的背景色，默认white
  boxSizing: '', // 盒模型
  borderColor: 'black', // 容器边框颜色
  borderType: 'solid', // 边框类型
  borderWidth: 0, // 边框宽度
  shadowBlur: 0, // 阴影模糊
  shadowColor: '', // 阴影颜色
  shadowOffsetX: 0, // 阴影偏移x
  shadowOffsetY: 0, // 阴影偏移y
  borderRadius: 0 // 组件容器的圆角，默认0.
}


const NODE_STYLE_TEXT = Object.assign(NODE_STYLE_COMMON, {
  type: 'text', // 模块类型
  content: '', // 绘制内容
  ellipsis: 0, // 文本省略 默认0,不省略
  margin: [0], // 相邻节点间距
  textStyle: {
    lineHeight: 28, // 文本的行高，默认28。
    baseLine: 'top', // 文本的对齐方式y，默认居左。
    align: 'left', // 文本的对齐方式x，默认居左。
    color: '#37374B', // 文字颜色
    fontSize: '28', // 文字大小
    fontWeight: 'normal', // 字重
    fontFamily: 'PingFangSC-Regular', // 字体
    font: 'normal normal 14px PingFangSC' // 文字斜体，字重，大小，字体
  }
})
const NODE_STYLE_IMAGE = Object.assign(NODE_STYLE_COMMON, {
  type: 'image', // 模块类型
  src: '', // 绘制内容
  imageStyle: {
    top: 0,
    left: 0,
    baseLine: 'middle',
    align: 'center',
    width: 0,
    height: 0
  }
})

const CALC_VALUE = {
  pt: 0, // padding-top
  pr: 0, // padding-right
  pb: 0, // padding-bottom
  pl: 0, // padding-left
  x: 0, // 绘制图片的左上角 x 坐标位置
  y: 0, // 绘制图片的左上角 y 坐标位置
  w: 0, // 绘制图片的宽度
  h: 0, // 绘制图片的高度
  dx: 0, // 绘制容器的左上角 x 坐标位置
  dy: 0, // 绘制容器的左上角 y 坐标位置
  dh: 0, // 绘制容器的高度
  dw: 0, // 绘制容器的宽度
  r: 0, // 绘制容器圆角
  r1: 0, // 绘制容器圆角
  r2: 0, // 绘制容器圆角
  r3: 0, // 绘制容器圆角
  r4: 0, // 绘制容器圆角
  sx: 0, // 绘制阴影x
  sy: 0, // 绘制阴影y
  sb: 0 // 绘制阴影blur
}

const setEllipsis = str => {
  return str.substr(0, str.length - 2) + '...'
}

const isX = n => {
  if (Object.prototype.toString.call(n) === '[object Object]') return 0 // [object Object]
  if (Object.prototype.toString.call(n) === '[object String]') return 1 // [object String]
  if (Object.prototype.toString.call(n) === '[object Number]') return 2 // [object Number]
  if (Object.prototype.toString.call(n) === '[object Boolean]') return 3 // [object Boolean]
  if (Object.prototype.toString.call(n) === '[object Undefined]') return 4 // [object Undefined]
  if (Object.prototype.toString.call(n) === '[object Null]') return 5 // [object Null]
  if (Object.prototype.toString.call(n) === '[object Function]') return 6 // [object Function]
  if (Object.prototype.toString.call(n) === '[object Date]') return 7 // [object Date]
  if (Object.prototype.toString.call(n) === '[object Array]') return 8 // [object Array]
  if (Object.prototype.toString.call(n) === '[object RegExp]') return 9 // [object RegExp]
  if (Object.prototype.toString.call(n) === '[object Error]') return 10 // [object Error]
}
// 遇到相同元素级属性，以后者（main）为准
// 不返还新Object，而是main改变

const mergeJSON = (minor, main) => {
  for (var key in minor) {
    if (main[key] === undefined) {
      // 不冲突的，直接赋值
      main[key] = minor[key]
      continue
    }

    // 冲突了，如果是Object，看看有么有不冲突的属性
    // 不是Object 则以main为主，忽略即可。故不需要else
    if (isX(minor[key]) === 0) {
      // arguments.callee 递归调用，并且与函数名解耦
      mergeJSON(minor[key], main[key])
    }
  }
}
const fsm = wx.getFileSystemManager() // 文件管理器
const TEMP_NAME = 'temporary_file_'
const TEMP_DIR = wx.env.USER_DATA_PATH
/**
 * 1、先确定大小，再确定位置
 * 2、先确定内容大小，再确定容器大小。
 * 3、先确定容器位置，再确定内容位置
 * 4、先绘制阴影，再绘制背景
 * 5、先绘制容器，再绘制内容
 * 6、先计算，再排版
 * 7、先下载，再排版
 */
class WxChart {
  // constructor () {}
  static getWxChart () {
    if (!WxChart.wxchart) {
      WxChart.wxchart = new WxChart()
    }

    return WxChart.wxchart
  }

  init (canvas, opts = defaultOpts) {
    this.canvas = canvas
    this.ctx = canvas.ctx
    this.width = opts.width
    this.height = opts.height
    this.dpr = opts.devicePixelRatio
    this.designWidth = opts.designWidth
    return this
  }

  async setOptions (options) {
    // const { canvas, width, height, designWidth } = this;


    /* 1、处理所有type=row的数据 */
    // console.log('1、处理所有type=row的数据')
    options.series = this.initSeries(options.series)

    /* 2、下载所有网络图片 */
    // console.log('2、下载所有网络图片')
    for (let i = 0; i < options.series.length; i++) {
      const item = options.series[i]

      if (item.type === 'image') {
        if (!item.src) {
          return Promise.reject(new Error('图片被程序猿拐跑啦，请再试一次'))
        }

        let downloadResult = await this.downloadImage(item.src)

        if (downloadResult) {
          // console.log('图片下载成功 => ', downloadResult)
          options.series[i].src = downloadResult
        } else {
          return Promise.reject(new Error(`图片下载失败 => ${item.src}`))
        }
      }
    }

    this.options = options


    /* 3、开始画 */
    // console.log('3、开始画')

    return this.draw()
  }

  initSeries (data) {
    let dpr = this.dpr
    let result = data
      .map(item => {
        if (item.type === 'section') {
          let children = item.children
          // 获取子元素的新宽度
          children = children.map(child => {
            if (child.type === 'text') {
              child.textStyle = this.setTextStyle(child.textStyle)
              child.width = this.getTextWidth(child) / dpr
            }
            return child
          })

          // 下一个的位置=上一个位置+上一个宽度+padding
          children.reduce(
            (prev, item) => {
              item.left = parseFloat(prev.width) + parseFloat(prev.left)
              return item
            },
            {
              width: 0,
              left: item.left,
              padding: []
            }
          )
          item = children
        } else if (item.type === 'text') {
          item.textStyle = this.setTextStyle(item.textStyle)
        }
        return item
      })
      .reduce(function (a, b) {
        return a.concat(b)
      }, [])
    return result
  }

  /* 获取文字宽度 */
  getTextWidth (data, str) {
    let ctx = this.ctx
    let content = str || data.content

    // 设置文字样式
    // console.log('文字样式 => ', data.textStyle.font)
    ctx.save()
    ctx.font = data.textStyle.font
    let width = ctx.measureText(content).width
    ctx.restore()
    return width
  }

  async draw () {
    let ctx = this.ctx

    let premise = this.options.series.map(item => {
      if (item.type === 'image') {
        return this.drawImage(item)
      } else if (item.type === 'text') {
        return this.drawArea(item)
      } else if (item.type === 'section') {
        return this.drawSection(item)
      }
    })
    let res = await Promise.all(premise)
    console.log(res)
    return new Promise(resolve =>
      ctx.draw(false, () =>
        setTimeout(() => {
          this.clear()
          resolve(this)
        }, 100)
      )
    )
  }

  /* 设置text样式 */
  setTextStyle (style) {
    // TODO 文字颜色
    let { fontSize, fontFamily, fontWeight, font } = style
    fontSize = `${Math.ceil(this.dataDpr(fontSize))}px`
    font = ['normal', fontWeight, fontSize, fontFamily].join(' ')
    return Object.assign(style, {
      fontSize,
      font
    })
  }

  // 设置item样式
  setItemStyle (style) {
    switch (style.type) {
      case 'text':
        mergeJSON(NODE_STYLE_TEXT, style)
        break
      case 'image':
        mergeJSON(NODE_STYLE_IMAGE, style)
        break

      default:
        break
    }
    return style
  }

  /* 计算文本框绘制高度 */
  calcFunDh (calc, data) {
    let { dh, h, pt, pb } = calc
    let { height } = data

    if (height && isX(height) === 2) {
      dh = this.dataDpr(height)
    } else {
      dh = h + pt + pb
    }

    calc.dh = parseInt(dh)
    return calc
  }

  /* 计算文本框绘制宽度 */
  calcFunDw (calc, data) {
    let { dw, w, pr, pl } = calc
    let { width } = data

    if (width && isX(width) === 2) {
      dw = parseInt(this.dataDpr(width))
    } else {
      dw = w + pr + pl
    }
    calc.dw = dw
    return calc
  }

  /* 计算文本框绘制位置 */
  calcFunDxDy (calc, data) {
    let cw = this.width // 画布的宽度
    let ch = this.height // 画布的高多
    let { dw, dh, dx, dy } = calc
    let { top, left } = data

    /* 计算文本框绘制位置y by top */
    switch (top) {
      case 'top':
        dy = 0
        break
      case 'middle':
        dy = (ch - dh) / 2
        break
      case 'bottom':
        dy = ch - dh
        break
      case 'auto':
        break
      default:
        dy = this.dataDpr(top)
        break
    }

    /* 计算文本框绘制位置x by left */
    switch (left) {
      case 'left':
        dx = 0
        break
      case 'center':
        dx = (cw - dw) / 2
        break
      case 'right':
        dx = cw - dw
        break
      case 'auto':
        break
      default:
        dx = this.dataDpr(left)
        break
    }

    calc.dx = parseInt(dx)
    calc.dy = parseInt(dy)
    return calc
  }

  /* 计算文本绘制位置 */
  calcFunXY (calc, data) {
    let { pt, pr, pb, pl, dw, dh, dx, dy, x, y, w, h } = calc
    let { type } = data
    let { baseLine, align } = data.textStyle || data.imageStyle

    let yxxxx = type === 'text' ? h / 2 : 0
    let xyyyy = type === 'text' ? w / 2 : 0
    /* 计算文本绘制位置y by top */
    switch (baseLine) {
      case 'top':
        y = dy + pt + yxxxx
        break
      case 'middle':
        y = dy + (dh - h) / 2 + yxxxx
        break
      case 'bottom':
        y = dy + dh - pb + yxxxx
        break
      case 'auto':
        y = dy + pt + yxxxx
        break
      default:
        y = dy + pt + yxxxx
        break
    }
    /* 计算文本绘制位置x by align */
    switch (align) {
      case 'left':
        x = dx + pl
        break
      case 'center':
        x = dx + (dw - w) / 2 + xyyyy
        break
      case 'right':
        x = dx + dw - pr
        break
      case 'auto':
        x = dx + pl
        break
      default:
        break
    }

    calc.x = parseInt(x)
    calc.y = parseInt(y)
    return calc
  }

  /* 计算容器圆角 */
  calcFunBr (calc, data) {
    if (isX(data.borderRadius) === 8) {
      calc.r1 = this.dataDpr(data.borderRadius[0])
      calc.r2 = this.dataDpr(data.borderRadius[1])
      calc.r3 = this.dataDpr(data.borderRadius[2])
      calc.r4 = this.dataDpr(data.borderRadius[3])
    } else {
      calc.r = this.dataDpr(data.borderRadius)
      calc.r1 = this.dataDpr(data.borderRadius)
      calc.r2 = this.dataDpr(data.borderRadius)
      calc.r3 = this.dataDpr(data.borderRadius)
      calc.r4 = this.dataDpr(data.borderRadius)
    }
    // console.log('计算圆角 => ', calc)
    return calc
  }
  /* 计算文本框边框宽度 */
  calcFunBw (calc, data) {
    calc.bw = this.dataDpr(data.borderWidth)
    return calc
  }
  /* 计算容器阴影 */
  calcFunSxyb (calc, data) {
    calc.sx = this.dataDpr(data.shadowOffsetX)
    calc.sy = this.dataDpr(data.shadowOffsetY)
    calc.sb = this.dataDpr(data.shadowBlur)
    return calc
  }
  /* 计算padding */
  calcFunPtpbl (calc, data) {
    let { padding } = data
    switch (padding.length) {
      case 1:
        calc.pt = calc.pr = calc.pb = calc.pl = this.dataDpr(padding[0])
        break
      case 2:
        calc.pt = calc.pb = this.dataDpr(padding[0])
        calc.pr = calc.pl = this.dataDpr(padding[1])
        break
      case 3:
        calc.pt = this.dataDpr(padding[0])
        calc.pb = this.dataDpr(padding[2])
        calc.pr = calc.pl = this.dataDpr(padding[1])
        break
      case 4:
        calc.pt = this.dataDpr(padding[0])
        calc.pb = this.dataDpr(padding[2])
        calc.pl = this.dataDpr(padding[1])
        calc.pr = this.dataDpr(padding[3])
        break

      default:
        break
    }
    return calc
  }
  /* 计算文字宽高 */
  calcFunTwh (calc, data) {
    let ctx = this.ctx
    // let { content } = data
    let { font, width, lineHeight } = data.textStyle
    ctx.font = font // 设置文字样式
    // console.log("ctx.measureText(content).width", ctx.measureText(content).width, content);
    calc.w = this.dataDpr(width) // 文本宽度
    calc.h = this.dataDpr(lineHeight) // 文本的高度
    return calc
  }
  /* 计算图片尺寸 */
  calcFunIwh (calc, data) {
    let { width, height, imageStyle } = data
    calc.w = this.dataDpr(imageStyle.width || width)
    calc.h = this.dataDpr(imageStyle.height || height)
    return calc
  }

  /* 计算文字数组 */
  calcFunTa (calc, data) {
    let { ellipsis, content } = data
    let { pr, pl, dw, bw } = calc
    let curw = dw - pl - pr + bw
    let result = []
    content
      .split('')
      .map(item => {
        return {
          value: item,
          size: this.getTextWidth(data, item)
        }
      })
      .reduce(
        (prev, current, index, arr) => {
          if (arr.length - 1 === index) {
            result.push({
              value: prev.value + current.value,
              size: prev.size + current.size
            })
          }
          if (prev.size >= curw) {
            result.push(prev)
            return current
          } else {
            return {
              value: prev.value + current.value,
              size: prev.size + current.size
            }
          }
        },
        { value: '', size: 0 }
      )
    // 限制行数>0并且限制行数>总行数 直接返回显示 获者不显示行数
    if ((ellipsis > 0 && ellipsis >= result.length) || ellipsis === 0) {
      return result
    } else {
      result = result.slice(0, ellipsis)
      result[ellipsis - 1].value = setEllipsis(result[ellipsis - 1].value)
      return result
    }
  }
  /* 重新计算容器位置与大小by borderwidth */
  calcFunReset (calc) {
    let { dx, dy, dh, dw, r, r1, r2, r3, r4, bw } = calc
    calc.dx = dx + bw / 2
    calc.dy = dy + bw / 2
    calc.dw = dw - bw
    calc.dh = dh - bw
    calc.r = r - bw / 2
    calc.r1 = r1 - bw / 2
    calc.r2 = r2 - bw / 2
    calc.r3 = r3 - bw / 2
    calc.r4 = r4 - bw / 2
    return calc
  }
  /* 尺寸转化 */
  dataDpr (number) {
    let dpr = this.dpr
    return dpr * number
  }
  /* 绘制路径 */
  drawPath (ctx, calc) {
    let { dx, dy, dh, dw, r1, r2, r3, r4 } = calc
    // console.log('绘制路径 => ', calc)
    let grid = [0, 0] // [x, y];

    grid[0] = dx + r1 // 更新坐标x
    grid[1] = dy // 更新坐标y
    ctx.moveTo(grid[0], grid[1]) // 移动到起始位置

    grid[0] = dx + dw - r1 // 更新坐标x
    grid[1] = dy // 更新坐标y
    ctx.lineTo(grid[0], grid[1]) // 绘制顶部边框

    grid[0] = dx + dw - r2 // 更新坐标x
    grid[1] = dy + r2 // 更新坐标y
    ctx.arc(grid[0], grid[1], r2, (3 / 2) * Math.PI, 2 * Math.PI) // 绘制右上圆角

    grid[0] = dx + dw // 更新坐标x
    grid[1] = dy + dh - r2 // 更新坐标y
    ctx.lineTo(grid[0], grid[1]) // 绘制右边框

    grid[0] = dx + dw - r3 // 更新坐标x
    grid[1] = dy + dh - r3 // 更新坐标y
    ctx.arc(grid[0], grid[1], r3, 0, (1 / 2) * Math.PI) // 绘制右下圆角

    grid[0] = dx + r3 // 更新坐标x
    grid[1] = dy + dh // 更新坐标y
    ctx.lineTo(grid[0], grid[1]) // 绘制下边框

    grid[0] = dx + r4 // 更新坐标x
    grid[1] = dy + dh - r4 // 更新坐标y
    ctx.arc(grid[0], grid[1], r4, (1 / 2) * Math.PI, Math.PI) // 绘制左下圆角

    grid[0] = dx // 更新坐标x
    grid[1] = dy + r4 // 更新坐标y
    ctx.lineTo(grid[0], grid[1]) // 绘制左边框

    grid[0] = dx + r1 // 更新坐标x
    grid[1] = dy + r1 // 更新坐标y
    ctx.arc(grid[0], grid[1], r1, Math.PI, (3 / 2) * Math.PI) // 绘制左上圆角

    return ctx
  }
  /* 绘制阴影 */
  drawShadow (ctx, calc, data) {
    let { sx, sy, sb } = calc
    let { shadowColor } = data
    ctx.setShadow(sx, sy, sb, shadowColor)
  }
  /* 绘制边框 */
  drawBorder (ctx, calc, data) {
    let { borderColor } = data
    calc = this.calcFunReset(calc, data)
    ctx.save() // 保存绘图上下文
    ctx.beginPath() // 开始新的路径
    ctx.setLineWidth(calc.bw) // 设置线条宽度
    this.drawPath(ctx, calc, data)
    ctx.setStrokeStyle(borderColor) // 设置线条填充色
    ctx.stroke() // 填充线条
    ctx.restore() // 恢复绘图上下文
  }
  /* 绘制view */
  drawView (ctx, calc, data) {
    let { show, backgroundColor, backgroundGradient, top, left, width, height } = data
    ctx.save() // 保存之前画的内容
    ctx.beginPath() // 开始新的路径

    this.drawPath(ctx, calc, data) // 绘制路径

    if (show) {
      this.drawShadow(ctx, calc, data) // 填充阴影
    }
    /* 填充背景 */
    if (backgroundGradient) {
      if (backgroundGradient === 'linear') {
        const grd = ctx.createLinearGradient(top, left, width, height)
        grd.addColorStop(0, backgroundColor[0])
        grd.addColorStop(1, backgroundColor[1])
        ctx.setFillStyle(grd)
      }
    } else {
      ctx.setFillStyle(backgroundColor) // 设置背景填充颜色
    }
    ctx.fill() // 根据路径填充一个矩形背景和阴影
    ctx.restore() // 还原所有内容

    /* 绘制边框 */
    if (calc.bw) {
      this.drawBorder(ctx, calc, data)
    }
  }
  /* 绘制文字 */
  drawText (ctx, calc, data, content) {
    let { align } = data
    let { color, font, textDecoration } = data.textStyle
    let { pr, pl, dw, x, y, bw } = calc

    ctx.save() // 保存绘图上下文
    ctx.font = font // 设置文字样式
    ctx.setFillStyle(color) // 设置填充颜色
    ctx.setTextAlign(align) // 水平位置 以fillText方法中的dx为参考标准
    ctx.setTextBaseline('middle') // 竖直位置 以fillText方法中的dy为参考标准
    ctx.fillText(content, x, y) // 在画布上绘制被填充的文本
    // 绘制中划线
    if (textDecoration === 'line-through') {
      this.drawLine(ctx, x + pl, y, dw - pl - pr + bw, color)
    }
    ctx.restore() // 恢复绘图上下文
  }
  drawLine (ctx, lx, ly, length, color, direction = 'horizontal') {
    // 保存之前画的内容
    // ctx.save();
    // 设置填充颜色
    ctx.setStrokeStyle(color)
    ctx.beginPath()
    ctx.setLineWidth(1)
    if (direction === 'horizontal') {
      // 水平
      ctx.moveTo(lx, ly)
      ctx.lineTo(lx + length, ly)
    } else if (direction === 'vertical') {
      // 竖直
      ctx.moveTo(lx, ly)
      ctx.lineTo(lx, ly + length)
    }
    ctx.stroke()
    // 还原所有内容
    // ctx.restore();
  }
  /* 绘制区域 */
  async drawArea (data) {
    return new Promise((resolve, reject) => {
      try {
        let ctx = this.ctx
        data = this.setItemStyle(data) // 文字框样式
        // data.textStyle = this.setTextStyle(data.textStyle); // 文字样式
        let calc = CALC_VALUE
        /* 计算文本宽高 */
        this.calcFunTwh(calc, data)
        /* 计算padding */
        this.calcFunPtpbl(calc, data)
        /* 计算圆角 */
        this.calcFunBr(calc, data)
        /* 计算边框宽度 */
        this.calcFunBw(calc, data)
        /* 计算阴影 */
        this.calcFunSxyb(calc, data)
        /* 计算高度 by padding */
        this.calcFunDh(calc, data)
        /* 计算宽度 by padding */
        this.calcFunDw(calc, data)
        /* 计算容器绘制位置 by padding */
        this.calcFunDxDy(calc, data)
        /* 计算文本绘制位置 by padding */
        this.calcFunXY(calc, data)
        // console.log(calc);
        /* 绘制容器 */
        this.drawView(ctx, calc, data)


        let contentArr = this.calcFunTa(calc, data)

        for (let i = 0; i < contentArr.length; i++) {
          this.drawText(ctx, calc, data, contentArr[i].value)
          calc.y += calc.h
        }

        // console.log(`绘制完成 => ${data.name}`)
        resolve(true)
      } catch (error) {
        console.log(`绘制失败 => ${data.name}`)
        reject(error)
      }
    })
  }
  /* 绘制图片 */
  async drawImage (data) {
    return new Promise((resolve, reject) => {
      try {
        let ctx = this.ctx
        data = this.setItemStyle(data)
        let calc = CALC_VALUE
        /* 绘制图片的宽高 */
        this.calcFunIwh(calc, data)
        /* 计算padding */
        this.calcFunPtpbl(calc, data)
        /* 计算圆角 */
        this.calcFunBr(calc, data)
        /* 计算边框宽度 */
        this.calcFunBw(calc, data)
        /* 计算阴影 */
        this.calcFunSxyb(calc, data)
        /* 计算高度 by padding */
        this.calcFunDh(calc, data)
        /* 计算宽度 by padding */
        this.calcFunDw(calc, data)
        /* 计算容器绘制位置 by padding */
        this.calcFunDxDy(calc, data)
        /* 计算内容的绘制位置 */
        this.calcFunXY(calc, data)
        /* 绘制容器 */
        this.drawView(ctx, calc, data)
        /* 绘制图片 */
        ctx.save() // 保存之前画的内容
        ctx.beginPath() // 开始新的路径
        this.drawPath(ctx, calc, data) // 绘制路径
        ctx.clip()
        ctx.drawImage(data.src, calc.x, calc.y, calc.w, calc.h)
        ctx.restore()
        // console.log(`绘制完成 => ${data.name}：${data.src}`)
        resolve(true)
      } catch (error) {
        console.log(`绘制失败 => ${data.name}：${data.src}`)
        reject(error)
      }
    })
  }
  /* 绘制区块 */
  async drawSection (data) {
    console.log(data)
  }

  /* 下载图片 */
  async downloadImage (url) {
    let regexp = /^https.*?/
    let regexpBase64 = /^data:image*?/
    if (regexp.test(url)) {
      return new Promise((resolve, reject) => {
        wx.downloadFile({
          url: url, // 仅为示例，并非真实的资源
          success (res) {
            if (res.statusCode === 200) {
              resolve(res.tempFilePath)
            } else {
              reject(res)
            }
          },
          fail (err) {
            console.error(err)
            reject(err)
          }
        })
      })
    } else if (regexpBase64.test(url)) {
      let FILE_BASE_NAME = `${TEMP_NAME}${Date.now().toString(36)}` // 自定义文件名
      let buffer = wx.base64ToArrayBuffer(url.replace(/data:image\/(png|jpeg|jpg);base64,/, '')) // 转为Buffer数据
      const filePath = `${TEMP_DIR}/${FILE_BASE_NAME}.png` // 文件路径
      return new Promise((resolve, reject) => {
        fsm.writeFile({
          // 把图片从本来拿到
          filePath,
          data: buffer,
          encoding: 'binary',
          success (res) {
            if (res.errMsg === 'writeFile:ok') {
              resolve(filePath)
            }
          },
          fail (err) {
            reject(err)
          }
        })
      })
    } else {
      return url
    }
  }
  // 清空临时文件
  clear () {
    fsm.readdir({
      // 获取文件列表
      dirPath: TEMP_DIR,
      success (res) {
        res.files
          .filter(item => item.indexOf(TEMP_NAME) !== -1)
          .forEach(el => {
            let name = `${TEMP_DIR}/${el}`
            // 删除时要注意文件名一定要和存的时候一样,不然会报没有unlink无文件权限问题
            fsm.unlink({
              filePath: name,
              fail (e) {
                console.log('readdir文件删除失败：', e)
              },
              success (succ) {
                console.log('readdir文件删除成功：', succ)
              }
            })
          })
      }
    })
  }
}

export default WxChart.getWxChart()
