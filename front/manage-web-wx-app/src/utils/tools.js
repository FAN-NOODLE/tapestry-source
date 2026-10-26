/**
 * 简单深拷贝
 * @param obj
 * @returns {any}
 */
const simpleClone = obj => JSON.parse(JSON.stringify(obj))
const { globalData } = getApp()
/**
 * 深拷贝
 * @param obj 被拷贝对象
 * @returns {*}
 */
const cloneDeep = obj => {
  if (obj === null || typeof obj !== 'object') return obj
  let newObj = obj instanceof Array ? [] : {}
  for (let key in obj) {
    if (obj.hasOwnProperty(key)) {
      newObj[key] = typeof obj[key] === 'object' ? cloneDeep(obj[key]) : obj[key]
    }
  }
  return newObj
}

/**
 * 函数记忆缓存
 * @param fn    未装饰函数
 * @param isObj 参数是否是一个对象
 * @returns {*} 装饰后函数
 */
const memorize = (fn, isObj) => {
  const {
    globalData: { cache }
  } = getApp()
  return function () {
    const params = isObj ? arguments[0] : arguments
    const key = `${fn.name}-${JSON.stringify(params)}`
    if (cache.has(key)) {
      return cache.get(key)
    } else {
      const result = isObj ? fn(params) : fn.apply(null, params)
      result && cache.set(key, result)
      return result
    }
  }
}

/**
 * 获取数组中极值项
 * @param arr
 * @param type
 * @returns {number}
 */
const getArrMax = (arr, type = 'max') => (type === 'max' ? Math.max(...arr) : Math.min(...arr))

/**
 * 获取数组中极值项索引
 * @param arr
 * @param type
 * @returns {*}
 */
const getArrMaxIndex = (arr, type = 'max') => arr.indexOf(getArrMax(arr, type))

/**
 * 判断对象是否为空
 * @param obj
 * @returns {boolean} 不是对象返回都是false
 */
const isEmpty = obj => typeof obj === 'object' && obj !== null && Object.keys(obj).length === 0

/**
 * 防抖
 * @param fn 原函数
 * @param delay 延迟
 * @param immediate 是否立即执行
 * @returns {function(): *} 包装后函数
 */
const debounce = (fn, delay = 300, immediate = true) => {
  let timer = null
  let now = true
  let result
  return function () {
    clearTimeout(timer)
    if (immediate) {
      if (now) {
        now = false
        // eslint-disable-next-line
        result = fn.apply(this, arguments)
      } else {
        timer = setTimeout(() => {
          // eslint-disable-next-line
          fn.apply(this, arguments)
        }, delay)
      }
    } else {
      timer = setTimeout(() => {
        // eslint-disable-next-line
        fn.apply(this, arguments)
      }, delay)
    }
    return result
  }
}

/**
 * 节流
 * @param fn 原函数
 * @param limit 时间限制
 * @returns {function(): *} 包装后函数
 */
const throttle = (fn, limit = 500) => {
  let previous = 0
  let result

  return function () {
    let now = +new Date()
    // !previous && (previous = now)
    if (now - previous > limit) {
      // eslint-disable-next-line
      result = fn.apply(this, arguments)
      previous = now
    }
    return result
  }
}

/**
 * 文本宽度计算
 * @param text  文本
 * @param fontSize 字号
 * @returns {number}
 */
const measureText = (text, fontSize = 10) => {
  let width = 0
  String(text)
    .split('')
    .forEach(item => {
      if (/[a-zA-Z]/.test(item)) {
        width += 7
      } else if (/[0-9]/.test(item)) {
        width += 5.5
      } else if (/\./.test(item)) {
        width += 2.7
      } else if (/-/.test(item)) {
        width += 3.25
      } else if (/[\u4e00-\u9fa5]/.test(item)) {
        // 中文匹配
        width += 10
      } else if (/\(|\)/.test(item)) {
        width += 3.73
      } else if (/\s/.test(item)) {
        width += 2.5
      } else if (/%/.test(item)) {
        width += 8
      } else {
        width += 10
      }
    })
  return (width * fontSize) / 10
}

/**
 * 函数执行时间计算
 * @param fn 原函数
 * @returns {Function}
 */
const time = fn => {
  return function () {
    fn.apply(this, arguments)
  }
}

/**
 * 判断两个对象相等
 * @param obj
 * @param other
 * @returns {boolean}
 */
const isEqual = (obj, other) => JSON.stringify(obj) === JSON.stringify(other)

/**
 * 金额格式化
 * @param money
 * @returns {string}
 */
const processMoney = money => {
  try {
    return money > 0 ? (money / 100).toFixed(2) : '0.00'
  } catch (e) {
    console.log(e)
    return money
  }
}

/**
 * 限制执行次数
 * @param fn
 * @param frequency
 * @returns {function(): *}
 */
const limitFrequency = (fn, frequency = 1) => {
  let _frequency = 0
  let result

  return function () {
    if (_frequency < frequency) {
      _frequency++
      result = fn.apply(this, arguments)
    } else {
      console.log('超出限制')
    }
    return result
  }
}

/**
 * 日期转换
 * @param time
 * @param format
 * @returns {{nowTime: string, nowYear: number, nowMonth: number, nowDate: number, nowHour: number, nowMinutes: number, nowSeconds: number, nowDayOfWeek: number}}
 */
const toDate = (time, format) => {
  // 兼容time传入 "2021.11.23 18:34" ios转换失败
  if (typeof time === 'string') {
    time = time.replace(/\./g, '-')
  }
  if (time.toString().length < 13) {
    time = parseInt(time) * 1000
  }
  let t = new Date(time)
  let nowYear = t.getFullYear() // 年
  let nowMonth = t.getMonth() + 1 // 月
  let nowDate = t.getDate() // 日
  let nowHour = t.getHours() // 时
  let nowMinutes = t.getMinutes() // 分
  let nowSeconds = t.getSeconds() // 秒
  let nowDayOfWeek = t.getDay() // 今天本周的第几天

  let tf = function (i) {
    return (i < 10 ? '0' : '') + i
  }
  // 当前时间
  let nowTime = function (format) {
    return format.replace(/yyyy|MM|dd|hh|mm|ss/g, function (a) {
      let str = ''
      switch (a) {
        case 'yyyy':
          str = tf(nowYear)
          break
        case 'MM':
          str = tf(nowMonth)
          break
        case 'dd':
          str = tf(nowDate)
          break
        case 'hh':
          str = tf(nowHour)
          break
        case 'mm':
          str = tf(nowMinutes)
          break
        case 'ss':
          str = tf(nowSeconds)
          break
      }
      return str
    })
  }
  return {
    nowTime: format ? nowTime(format) : '',
    nowYear: nowYear,
    nowMonth: nowMonth,
    nowDate: nowDate,
    nowHour: nowHour,
    nowMinutes: nowMinutes,
    nowSeconds: nowSeconds,
    nowDayOfWeek: nowDayOfWeek
  }
}

/**
 * 千分位/电话号切割
 * @param num 数字
 * @param mode 模式
 * @returns {string}
 */
const split = (num, mode) => {
  try {
    const base = mode === 'space' ? '$1 ' : '$1,'
    const limit = mode === 'space' ? 4 : 3
    const source = String(num).split('.') // 按小数点分成2部分
    source[0] = source[0].replace(new RegExp(`(\\d)(?=(\\d{${limit}})+$)`, 'ig'), base) // 只将整数部分进行都好分割
    return source.join('.') // 再将小数部分合并进来
  } catch (e) {
    console.log(e)
    return num
  }
}
/**
 * 遍历tree
 * @param {*} tree 整个树结构
 * @param {*} action 回调
 */
const mapTree = (tree, action) => {
  if (!tree) return false
  action(tree)
  if (!tree.children) return false
  tree.children.forEach(item => mapTree(item, action))
}

/**
 * 查找tree所有父级
 * @param {*} _item 当前项
 * @param {*} treeArr tree数据
 * @param {*} fn 页面回调函数
 */
const findParent = (_item, treeArr, fn) => {
  treeArr.forEach(tree => {
    mapTree(tree, item => {
      if (_item.pid === item.id) {
        fn(item)
        item.pid && findParent(item, treeArr, fn)
      }
    })
  })
}

/**
 * 模拟Promise.all方法
 * @param {*} promiseList 要执行的promise集合
 * @returns {Promise}
 */
const promiseAll = promiseList => {
  return new Promise((resolve, reject) => {
    let results = []
    let counter = 0
    for (let i = 0; i < promiseList.length; i++) {
      const promise = Promise.resolve(promiseList[i])
      // eslint-disable-next-line promise/catch-or-return
      promise
        .then(result => {
          counter++
          results[i] = result
          if (counter === promiseList.length) {
            resolve(results)
          }
        })
        .catch(e => {
          reject(e)
        })
    }
  })
}

/**
 * 模拟Promise.allSettled方法
 * @param {*} promiseList 要执行的promise集合
 * @returns {Promise}
 */
const promiseAllSettled = promiseList => {
  return new Promise(resolve => {
    let results = []
    let counter = 0
    for (let i = 0; i < promiseList.length; i++) {
      const promise = Promise.resolve(promiseList[i])
      // eslint-disable-next-line promise/catch-or-return
      promise
        .then(result => {
          results[i] = result
        })
        .catch(e => {
          console.error(e)
          results[i] = e
        })
        .finally(() => {
          counter++
          if (counter === promiseList.length) {
            resolve(results)
          }
        })
    }
  })
}

/**
 * 手机号脱敏处理
 * @param {string} phone 手机号
 * @returns {string}
 */
const handleMobile = phone => {
  if (!phone) return ''
  const reg = /(\d{3})\d*(\d{4})/
  return phone.replace(reg, '$1****$2')
}

/**
 * 截断字符串
 * @param {string} text 原始字符串
 * @param {number} length 要截断的长度
 * @param {string} postfix 截断后的后缀
 * @returns {string}
 */
const truncateText = (text, length, postfix = '...') => {
  text = text || ''
  if (text.length > length) {
    return text.substring(0, length) + postfix
  }
  return text
}

/**
 * 要解析的查询字符串
 * @param {string} queryString 查询字符串
 * @returns {object}
 */
const toQueryParams = queryString => {
  let result = {}
  queryString
    .substring(queryString.indexOf('?') + 1)
    .split('&')
    .map(x => x.split('='))
    .forEach(item => {
      let key = item[0]
      let value = item[1]
      if (key && value) {
        result[key] = value
      }
    })
  return result
}

/**
 * 把JSON对象解析为查询字符串
 * @param {object} queryParams 要解析的对象
 * @returns {string}
 */
const toQueryString = queryParams => {
  let result = Object.keys(queryParams || {})
    .sort()
    .map(key => {
      if (queryParams[key]) {
        return `${key}=${queryParams[key]}`
      } else {
        return ''
      }
    })
  return result.join('&')
}

/**
 * 给URL添加查询参数
 * @param {string} url 要添加参数的url
 * @param {object} params 要添加的参数
 * @returns {string}
 */
const setQueryParams = (url, params) => {
  const index = url.indexOf('?')
  return `${url.substring(0, index > -1 ? index : url.length)}?${toQueryString(params)}`
}

/**
 * 等待指定时间后执行
 * @param {number} milliseconds 要等待的毫秒数
 * @returns {Promise}
 */
const sleep = milliseconds => {
  return new Promise(resolve => {
    setTimeout(resolve, milliseconds)
  })
}

const processWXAvatarUrl = url => {
  if (url.indexOf('thirdwx.qlogo.cn') > -1) {
    return url.replace('thirdwx.qlogo.cn', 'wx.qlogo.cn')
  }
  return url
}

// // 一段时间的时间戳-例昨天、七天内、30天时间戳
const timeForMat = count => {
  // 拼接时间
  const time1 = new Date()
  const Y1 = time1.getFullYear()
  const M1 = time1.getMonth() + 1 > 10 ? time1.getMonth() + 1 : '0' + (time1.getMonth() + 1)
  const D1 = time1.getDate() > 10 ? time1.getDate() : '0' + time1.getDate()
  const timer1 = globalData.isIos ? Y1 + '/' + M1 + '/' + D1 : Y1 + '-' + M1 + '-' + D1 // 当前时间
  const time2 = new Date()
  time2.setTime(time2.getTime() - 24 * 60 * 60 * 1000 * count)
  const Y2 = time2.getFullYear()
  const M2 = time2.getMonth() + 1 > 10 ? time2.getMonth() + 1 : '0' + (time2.getMonth() + 1)
  const D2 = time2.getDate() > 10 ? time2.getDate() : '0' + time2.getDate()
  const timer2 = globalData.isIos ? Y2 + '/' + M2 + '/' + D2 : Y2 + '-' + M2 + '-' + D2 // 之前的7天或者30天
  return {
    t1: new Date(timer1 + ' 23:59:59').getTime() / 1000,
    t2: new Date(timer2 + ' 00:00:00').getTime() / 1000
  }
}


/**
 * 时间格式化
 *
 * @param {*} timestamp
 * @param {*} fmt
 * @returns
 */
const formatDate = (timestamp, fmt) => {
  if (!fmt) {
    fmt = 'yyyy-MM-dd hh:mm:ss'
  }
  let date = new Date(timestamp)
  let o = {
    'M+': date.getMonth() + 1, // 月份
    'd+': date.getDate(), // 日
    'h+': date.getHours(), // 小时
    'm+': date.getMinutes(), // 分
    's+': date.getSeconds(), // 秒
    'q+': Math.floor((date.getMonth() + 3) / 3), // 季度
    S: date.getMilliseconds() // 毫秒
  }
  if (/(y+)/.test(fmt)) {
    fmt = fmt.replace(RegExp.$1, (date.getFullYear() + '').substr(4 - RegExp.$1.length))
  }

  for (let k in o) {
    if (new RegExp('(' + k + ')').test(fmt)) {
      fmt = fmt.replace(
        RegExp.$1,

        RegExp.$1.length === 1 ? o[k] : ('00' + o[k]).substr(('' + o[k]).length)
      )
    }
  }

  return fmt
}

/** 数字转千位符 */
function formatWithArray (number) {
  // 转为字符串，并按照.拆分
  const arr = (number + '').split('.')
  // 整数部分再拆分
  const int = arr[0].split('')
  // 小数部分
  const fraction = arr[1] || ''
  // 返回的变量
  let r = ''
  int.reverse().forEach(function (v, i) {
    // 非第一位并且是位值是3的倍数，添加“,”
    if (i !== 0 && i % 3 === 0) {
      r = v + ',' + r
    } else {
      // 正常添加字符(这是好写法)
      r = v + r
    }
  })
  // 整数部分和小数部分拼接
  return r + (fraction ? '.' + fraction : '')
}
/**
 * @description: 对象转url参数
 * @param {Object} obj
 */
function objToUrlParams (obj) {
  const params = []

  Object.keys(obj).forEach((key) => {
    let value = obj[key]
    // 如果值为undefined我们将其置空
    if (typeof value === 'undefined') {
      value = ''
    }
    // 对于需要编码的文本（比如说中文）我们要进行编码
    params.push([key, encodeURIComponent(value)].join('='))
  })

  return params.join('&')
}

export { time, split, toDate, isEmpty, isEqual, debounce, throttle, memorize, cloneDeep, getArrMax, simpleClone, measureText, processMoney, getArrMaxIndex, limitFrequency, findParent, promiseAll, promiseAllSettled, handleMobile, truncateText, toQueryParams, toQueryString, setQueryParams, sleep, processWXAvatarUrl, timeForMat, formatDate, formatWithArray, objToUrlParams }
