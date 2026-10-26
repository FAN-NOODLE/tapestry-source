/**
 * @description: 判断是否是字符串
 * @param {any} str
 */
export const isString = (str) => {
  return typeof str === 'string'
}
/**
 * @description: 判断是否是数组
 */
export const isArray = (arr) => {
  return Array.isArray(arr)
}
/**
 * @description: 判断是否是对象
 * @param {any} obj
 * @return {Boolean}
*/
export const isObject = (obj) => {
  return obj !== null && obj instanceof Object
}

