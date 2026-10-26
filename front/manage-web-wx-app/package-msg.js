var https = require('https')

/**
 * @description: 判断是否是字符串
 * @param {any} str
 */
const isString = (str) => {
  return typeof str === 'string'
}
/**
 * @description: 判断是否是数组
 */
const isArray = (arr) => {
  return Array.isArray(arr)
}
/**
 * @description: 判断是否是对象
 * @param {any} obj
 * @return {Boolean}
*/
const isObject = (obj) => {
  return obj !== null && obj instanceof Object
}

/**
 * 机器人发送消息
 * @param {object} data
 * @returns
 */
var sendMessageSync = async data => {
  var robotKey = '4b933656-20a3-48a6-a4e6-c68696a02127'
  if (!robotKey) {
    throw new Error('缺少环境变量 WECHAT_ROBOT_KEY')
  }

  return new Promise((resolve, reject) => {
    var wxReq = https.request({
      hostname: 'qyapi.weixin.qq.com',
      port: 443,
      path: '/cgi-bin/webhook/send?key=' + encodeURIComponent(robotKey),
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      }
    }, res => {
      // 如果状态码不是 200 就输出状态码
      if (res.statusCode !== 200) console.log(res.statusCode)

      res.on('data', d => {
        // 把接收到的内容转为字符串在控制台输出
        console.log(d.toString())
        resolve(d.toString())
      })
    })
    wxReq.write(JSON.stringify(data))
    wxReq.end()
    wxReq.on('error', error => {
      console.log('error', error)
      reject(error)
    })
  })
}

/**
 * 创建消息
 * @param {object|string} data
 * @param {string} data.msgtype 消息类型 当前自定义机器人支持文本（text）、markdown（markdown）、图片（image）、图文（news）四种消息类型。
 * @param {object} data.text 文本类型
 * @param {string} data.content 文本内容，最长不超过2048个字节，必须是utf8编码
 * @param {object} data.markdown markdown类型
 * @param {string} data.content markdown内容，最长不超过4096个字节，必须是utf8编码
 * @param {[string]} data.mentioned_list userid的列表，提醒群中的指定成员(@某个成员)，@all表示提醒所有人，如果开发者获取不到userid，可以使用mentioned_mobile_list
 * @param {[string]} data.mentioned_mobile_list 手机号列表，提醒手机号对应的群成员(@某个成员)，@all表示提醒所有人
 * @param {object} data.image 图片类型
 * @param {string} data.base64 图片内容的base64编码
 * @param {string} data.md5 图片内容（base64编码前）的md5值
 * @param {object} data.news 图文类型
 * @param {[object]} data.articles 图文消息，一个图文消息支持1到8条图文
 * @param {string} data.title 标题，不超过128个字节，超过会自动截断
 * @param {string} data.description 描述，不超过512个字节，超过会自动截断
 * @param {string} data.url 点击后跳转的链接。
 * @param {string} data.picurl 图文消息的图片链接，支持JPG、PNG格式，较好的效果为大图 1068*455，小图150*150。
 * @returns
 */

var creatMessage = data => {
  data = {
    msgtype: data.msgtype,
    content: data.content
  }
  var _data = {}
  if (data instanceof Object) {
    _data.msgtype = data.msgtype
    switch (data.msgtype) {
      case 'text':
        _data[data.msgtype] = {
          content: data.content
        }
        break
      case 'markdown':
        _data[data.msgtype] = {
          content: data.content
        }
        break
      case 'image':
        _data[data.msgtype] = {
          base64: data.base64,
          md5: data.md5
        }
        break
      case 'news':
        _data[data.msgtype] = {
          articles: data.articles
        }
        break

      default:
        _data.msgtype = 'text'
        _data.text = JSON.stringify(data)
        break
    }
    // userid的列表，提醒群中的指定成员
    if (data.mentioned_list) {
      _data[data.msgtype].mentioned_list = data.mentioned_list
    }
    // 手机号列表，提醒手机号对应的群成员
    if (data.mentioned_mobile_list) {
      _data[data.msgtype].mentioned_mobile_list = data.mentioned_mobile_list
    }
  } else {
    _data.msgtype = 'text'
    _data.text = {
      content: data || '请输入消息内容'
    }
  }
  return _data
}

/**
 * 设置markdown文字颜色
 * @param {object} data 文本内容
 * @param {string} data.text 文本内容
 * @param {string} data.color 文本颜色 info-绿色；comment-灰色；warning-橙红色
 * @returns
 */
var setMdFontColor = data => {
  return `<font color="${data.color}">${data.text}</font>`
}

/**
 * 文本引用
 * @param {string} text 文本内容
 * @returns
 */
var setMdQuote = text => `> ${text}\n`

/**
 * 文本标题
 * @param {object} data 文本内容
 * @param {string} data.text 文本内容
 * @param {*} data.h 标题级别
 * @returns
 */
var setMdTitle = data => {
  var arr = []
  for (let i = 0; i < data.h; i++) {
    if (i >= 6) break
    arr.push('#')
  }

  return `${arr.join('')} ${data.text}\n`
}

/**
 * 文本链接
 * @param {object} data 数据
 * @param {string} data.text 文本内容
 * @param {string} data.url 跳转链接
 * @returns
 */
var setMdLink = data => {
  return `[${data.text}](${data.url})`
}

/**
 * 文本加粗
 * @param {string} text 文本内容
 * @returns
 */
var setMdFontBold = text => `**${text}**`

/**
 * 设置代码段
 * @param {string} text 文本内容
 * @returns
 */
var setMdCode = text => `\`${text}\``

/**
 * 设置md文本
 * @param {object} data
 * @param {string} data.text 文本内容
 * @param {string} data.key 文本类型 title:标题；link:链接；quote:引用；bold:加粗；color:颜色；code:代码段
 * @param {string} data.h 文本标题：标题级别
 * @param {string} data.url 文本链接：跳转地址
 * @param {string} data.color 文本颜色：颜色关键字
 * @returns
 */
var setMdText = data => {
  if (typeof data === 'string') return data
  let text = ''
  switch (data.key) {
    case 'title':
      text = setMdTitle(data)
      break
    case 'link':
      text = setMdLink(data)
      break
    case 'quote':
      text = setMdQuote(data.text)
      break
    case 'bold':
      text = setMdFontBold(data.text)
      break
    case 'color':
      text = setMdFontColor(data)
      break
    case 'code':
      text = setMdCode(data.text)
      break

    default:
      text = data.text || JSON.stringify(data)
      break
  }

  return text
}

var setMdContent = data => {
  let _data = []
  if (isArray(data)) {
    _data = data
  } else if (isString(data) || isObject(data)) {
    _data.push(data)
  }
  if (_data.length === 0) {
    console.log('数据内容不能为空')
    return false
  }
  for (let i = 0; i < _data.length; i++) {
    const item = _data[i]
    if (isString(item)) {
      // 如果是字符串，直接处理数据
      _data[i] = setMdText(item)
    } else if (isArray(item)) {
      // 如果是数组，直接进入递归
      _data[i] = setMdContent(data)
    } else if (isObject(item)) {
      // 如果是对象，判断是否需要递归
      if (isArray(item.text) && item.text.length > 0) {
        // 有子元素，递归
        _data[i].text = setMdContent(item.text)
        _data[i] = setMdText(item)
      } else {
        // 无子元素，直接处理数据
        _data[i] = setMdText(item)
      }
    } else {
      _data[i] = false
    }
  }
  return _data.join('')
}
module.exports = {
  creatMessage: creatMessage,
  sendMessageSync: sendMessageSync,
  setMdContent: setMdContent
}
