/**
 * 请求重试工具类
 * @param fn 定义的接口
 * @param times 重试次数
 */
class Retry {
  constructor (fn, times = 4) {
    this.wxRequest = fn // 封装的业务接口
    this.count = 0 // 请求次数
    this.maxCount = times // 请求尝试次数
    this.params = {} // 请求参数
  }
}
/**
 * 执行请求的函数
 * @param  {...any} args
 * @returns
 */
Retry.prototype.request = function (params) {
  this.params = params
  return new Promise((resolve, reject) => {
    this.dealRequest(resolve, reject)
  })
}
Retry.prototype.dealRequest = function (resolve, reject) {
  this.wxRequest(this.params).then(res => {
    console.log('retry-res', res)
    resolve(res)
  }).catch(err => {
    console.log('retry-err', err)
    if (this.count < this.maxCount) {
      setTimeout(() => {
        this.count++
        this.dealRequest(resolve, reject)
      }, 500)
    } else {
      reject(err)
    }
  })
}

export default Retry
