/**
 * 函数防抖
 * @param   delay     延迟           type Number
 * @param   immediate 是否立即触发一次 type Boolean
 */
const debounce = (delay = 300, immediate = false) => {
  return (target, key, descriptor) => {
    const fn = descriptor.value
    let timer = null
    let now = true
    let result

    descriptor.value = (...args) => {
      clearTimeout(timer)
      if (immediate) {
        if (now) {
          result = fn.apply(target, args)
          now = false
        } else {
          timer = setTimeout(() => {
            fn.apply(target, args)
          }, delay)
        }
      } else {
        timer = setTimeout(() => {
          fn.apply(target, arguments)
        }, delay)
      }
      return result
    }
  }
}
/**
 * 函数节流
 * @param delay     延迟            type Number
 * @param limit     必触发时间限制   type Number
 */
const throttle = (delay = 300, limit = 1000) => {
  return (target, key, descriptor) => {
    const fn = descriptor.value
    let timer = null
    let previous = null
    let result

    descriptor.value = (...args) => {
      let now = +new Date()
      !previous && (previous = now)
      if (now - previous > limit) {
        result = fn.apply(target, args)
        previous = now
      } else {
        clearTimeout(timer)
        timer = setTimeout(() => {
          fn.apply(target, args)
        }, delay)
      }
      return result
    }
    return descriptor
  }
}

export {
  debounce,
  throttle
}

