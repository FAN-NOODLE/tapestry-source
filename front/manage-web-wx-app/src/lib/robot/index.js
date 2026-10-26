import {
  setMdContent,
  creatMessage,
  sendMessageSync
} from './message'
import dayjs from 'dayjs'
import { getNetworkType, getSystem } from './system'

const systemInfo = getSystem()
/**
 * 上报错误信息
 * @param {Wxapp.Report} data 参数
 */
export const report = async data => {
  try {
    const _data = []
    // 标题
    _data.push({
      key: 'title',
      h: 4,
      text: [
        {
          key: 'color',
          color: 'warning',
          text: data.title
        },
        '\n'
      ]
    })
    // 系统信息
    _data.push(
      {
        text: [
          '手机型号：',
          {
            key: 'color',
            color: 'comment',
            text: systemInfo.model
          },
          '\n'
        ]
      },
      {
        text: [
          '操作系统：',
          {
            key: 'color',
            color: 'comment',
            text: systemInfo.system
          },
          '\n'
        ]
      },
      {
        text: [
          '微信版本：',
          {
            key: 'color',
            color: 'comment',
            text: systemInfo.version
          },
          '\n'
        ]
      },
      {
        text: [
          '客户端平台：',
          {
            key: 'color',
            color: 'comment',
            text: systemInfo.platform
          },
          '\n'
        ]
      },
      {
        text: [
          '基础库版本：',
          {
            key: 'color',
            color: 'comment',
            text: systemInfo.SDKVersion
          },
          '\n'
        ]
      },
    )
    // 网络状态
    const { networkType } = await getNetworkType()
    _data.push({
      text: [
        '网络状态：',
        {
          key: 'color',
          color: 'comment',
          text: networkType
        },
        '\n'
      ]
    })
    // 错误信息
    if (data.url) {
      _data.push(
        {
          text: [
            '请求地址：',
            {
              key: 'color',
              color: 'comment',
              text: data.url || ''
            },
            '\n'
          ]
        },
        {
          text: [
            '请求时间：',
            {
              key: 'color',
              color: 'comment',
              text: dayjs(data.time).format('YYYY-MM-DD HH:mm:ss')
            },
            '\n'
          ]
        },
        {
          text: [
            '请求头：',
            {
              key: 'color',
              color: 'comment',
              text: JSON.stringify(data.header || {})
            },
            '\n'
          ]
        },
        {
          text: [
            '请求参数：',
            {
              key: 'color',
              color: 'comment',
              text: JSON.stringify(data.data || {})
            },
            '\n'
          ]
        },
      )
    } else {
      _data.push({
        text: [
          '上报时间：',
          {
            key: 'color',
            color: 'comment',
            text: dayjs().format('YYYY-MM-DD HH:mm:ss')
          },
          '\n'
        ]
      })
    }
    _data.push({
      text: [
        '错误信息：',
        {
          key: 'color',
          color: 'warning',
          text: JSON.stringify(data.errMsg || {})
        },
        '\n'
      ]
    })
    // 发送信息
    sendMessageSync(
      creatMessage({
        msgtype: 'markdown',
        content: setMdContent(_data)
      }),
    )
  } catch (error) {
    console.log(error)
  }
}
