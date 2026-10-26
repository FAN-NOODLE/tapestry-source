const getAttr = (obj, str) => {
  if (!str) {
    console.log('无属性参数')
    return obj
  }
  let arr = str.split('.')
  let newObj = Object.assign({}, obj)
  arr.find(item => {
    if (newObj.hasOwnProperty(item)) {
      if (typeof newObj[item] === 'object') {
        newObj = Object.assign({}, newObj[item])
        return false
      } else {
        newObj = newObj[item]
        return true
      }
    } else {
      newObj = false
      return true
    }
  })
  if (!newObj) {
    console.info('请确认传参是否正确，第一参数为原对象， 第二参数为属性字符串以.隔开')
  }
  return newObj
}


if (wx) {
  wx.getAttr = getAttr
}

export { getAttr }
