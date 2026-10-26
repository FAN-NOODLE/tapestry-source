import { request } from '../utils/request'

/** 查看邀请函详情 */
const inviteDetail = (data) =>
  request(
    {
      url: 'invite/get-record-info',
      method: 'POST'
    },
    data
  )

/** 查看邀请函预览详情 */
const inviteInfo = (data) =>
  request(
    {
      url: 'invite/info',
      method: 'POST'
    },
    data
  )

/** 创建邀请函订单 */
const createOrder = (data) =>
  request(
    {
      url: 'invite/create-order',
      method: 'POST'
    },
    data
  )

export { inviteDetail, createOrder, inviteInfo }
