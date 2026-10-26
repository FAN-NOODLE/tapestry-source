class Countdown {
  /**
   * 是否不是空对象
   * @param obj
   * @returns {boolean}
   */
  static isEmpty (obj) {
    let result
    if (typeof obj === 'object' && obj !== null) {
      result = Object.keys(obj).length > 0
    }
    return !result
  }

  /**
   * 时间格式兼容
   * @param time 时间
   * @returns {number}
   */
  static fixTime (time) {
    return time
      ? String(time).length === 13
        ? parseInt(String(Number(time) / 1000))
        : Number(time)
      : ''
  }

  /**
   * 时分秒转化
   * @param time 时间戳
   * @returns {{hour: *, minute: *, second: *}}
   */
  translateTime (time) {
    let _time = Countdown.fixTime(time)
    const heaven = Math.floor(_time / 3600 / 24)
    const hour = this.useHeaven
      ? Math.floor((_time / 3600) % 24)
      : Math.floor(_time / 3600)
    const minute = Math.floor((_time / 60) % 60)
    const second = Math.floor(_time % 60)
    const ms = this.msArr[this.msIndex]
    this.msIndex++

    return {
      heaven,
      hour: hour < 10 ? '0' + hour : String(hour),
      minute: minute < 10 ? '0' + minute : String(minute),
      second: second < 10 ? '0' + second : String(second),
      ms
    }
  }

  /**
   * 初始化
   */
  constructor (config = {}) {
    const { useMs = 1, useHeaven, useNowTime } = config
    this.useMs = useMs
    this.useHeaven = useHeaven
    this.useNowTime = useNowTime
    let index = useMs
    this.msArr = []
    while (index >= 0) {
      this.msArr.push(String(index - 1))
      index--
    }
    this.msIndex = 0
    this.init()
  }

  /**
   * 根据开始结束转换时分秒
   * @param config  开始结束时间配置对象
   * @param cb      回调函数
   * @returns {{hour: string, minute: string, second: string}}
   */
  start (config, cb) {
    this.init()
    const isEmptyObjectOld = Countdown.isEmpty(config)
    const isEmptyObjectNew = Countdown.isEmpty(this.config)
    if (isEmptyObjectOld) {
      this.error = 'params error!'
      return this
    }

    config = !isEmptyObjectNew ? this.config : config
    isEmptyObjectNew && (this.config = config)
    this.cb = cb

    const delay = 1000 / this.useMs

    // fix time format
    Object.assign(config, {
      startTime: Countdown.fixTime(config.startTime),
      endTime: Countdown.fixTime(config.endTime),
      nowTime: Countdown.fixTime(config.nowTime)
    })
    // default config
    const { endTime } = config
    // condition resolve
    clearInterval(this.timer)
    this.timer = setInterval(() => {
      const nowTime =
        config.nowTime || parseInt(String(new Date().getTime() / 1000))
      if (this.msIndex === this.useMs) {
        this.msIndex = 0
      }
      const param = this.translateTime(endTime - nowTime)
      if (nowTime - 1 >= endTime) {
        clearInterval(this.timer)
        cb && cb(null)
        this.callback && this.callback(null)
      } else {
        cb && cb(param)
        this.callback && this.callback(param)
      }
      if (this.useNowTime) {
        if (this.useMs === 1) {
          this.config.nowTime += 1
        } else if (this.msIndex === this.msArr.length - 1) {
          this.config.nowTime += 1
        }
      }
    }, delay)
    return this
  }

  /**
   * 停止计时
   */
  stop () {
    clearInterval(this.timer)
    this.init()
    return this
  }

  /**
   * 重启计时
   * @param config
   * @param cb
   */
  restart (config, cb) {
    this.stop()
    const callback = cb || this.cb
    this.start({ ...this.config, ...config }, callback)
    return this
  }

  // 初始化内部变量
  init () {
    this.config = {}
    this.cb = null
    this.error = ''
  }

  // 语法糖
  then (callback) {
    this.callback = callback
    if (this.useNowTime) {
      this.config.nowTime += 1
    }
    this.start(this.config)
    return this
  }

  // 语法糖
  finally (callback) {
    this.callback = callback
    this.restart(this.config)
    return this
  }

  // 捕获错误
  catch (cb) {
    cb(this.error)
  }
}

export default Countdown
