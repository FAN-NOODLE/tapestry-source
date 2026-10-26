/**
 * 不要尝试修改页面栈，会导致路由以及页面状态错误。
 * 不要在 App.onLaunch 的时候调用 getCurrentPages()，此时page 还没有生成
 **/
// export function getCurrentPage () {
//   let pages = getCurrentPages() // 当前页面栈
//   // 当前页面为页面栈的最后一个元素
//   let prevPage = pages[pages.length - 1] // 当前页面
//   console.log('getCurrentPage', prevPage)
//   return prevPage
// }


/* 获取当前页url */
export const getCurrentPageUrl = () => {
  // eslint-disable-next-line
  const pages = getCurrentPages()    // 获取加载的页面
  const currentPage = pages[pages.length - 1] // 获取当前页面的对象
  const url = '/' + currentPage.route // 当前页面url
  return url
}

/* 获取当前页带参数的url */
export const getCurrentPageUrlWithArgs = () => {
  // eslint-disable-next-line
  const pages = getCurrentPages()    // 获取加载的页面
  const currentPage = pages[pages.length - 1] // 获取当前页面的对象
  const url = currentPage.route // 当前页面url
  const options = currentPage.options // 如果要获取url中所带的参数可以查看options

  // 拼接url的参数
  let urlWithArgs = url + '?'
  for (const key in options) {
    const value = options[key]
    urlWithArgs += key + '=' + value + '&'
  }
  urlWithArgs = urlWithArgs.substring(0, urlWithArgs.length - 1)

  return '/' + urlWithArgs
}
