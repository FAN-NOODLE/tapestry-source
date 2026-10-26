/**
 * 数组转树形结构
 * @param {array} list 被转换的数组
 * @return array
 */
export function groupArr(list, field) {
  let fieldList = []
  let groupList = []
  // 获取数组指定字段值，并去重
  list.map((element) => {
    element.checked = false
    if (fieldList.indexOf(element[field]) === -1) {
      fieldList.push(element[field])
    }
  })
  for (let i = 0; i < fieldList.length; i++) {
    // 过滤出匹配到的数据
    let arr = list.filter((element) => {
      return element[field] === fieldList[i]
    })
    groupList.push({
      city: arr[0].city || '',
      checkedNum: 0,
      allChecked: false,
      isCollapse: true,
      children: arr
    })
  }
  // todo: 根据城市排序，城市为空的排在最后
  groupList.sort((a, b) => {
    const cityA = a.city || ''
    const cityB = b.city || ''
    if (cityA === '') {
      return 1
    } else if (cityB === '') {
      return -1
    } else {
      return cityA.localeCompare(cityB)
    }
  })
  return groupList
}
