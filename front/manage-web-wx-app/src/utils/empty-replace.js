import { isEmpty } from './tools'

export default function emptyReplace (obj) {
  if (!isEmpty(obj)) return obj
  let newObj = obj instanceof Array ? [] : {}
  for (let key in obj) {
    if (obj.hasOwnProperty(key)) {
      newObj[key] = !isEmpty(obj) ? emptyReplace(obj[key]) : obj[key] === undefined || obj[key] === '' ? '--' : obj[key]
    }
  }
  return newObj
}
