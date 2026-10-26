export function isObj(x) {
  var type = typeof x
  return x !== null && (type === 'object' || type === 'function')
}
