/**
 * 监听器
 * @param obj  对象
 * @param fn   回调
 * @param deep 是否深度监听
 * @returns {*}
 */
export default function watch (obj, fn, deep) {
  if (!obj || typeof obj !== 'object' || obj === null) {
    return obj
  }

  const definePro = (obj, key, value) => {
    Object.defineProperty(obj, key, {
      get () {
        return value
      },
      set (newVal) {
        fn(key, value, newVal)
        value = newVal
      }
    })
    deep && watch(value, fn)
  }

  for (let i in obj) {
    definePro(obj, i, obj[i])
  }
};
