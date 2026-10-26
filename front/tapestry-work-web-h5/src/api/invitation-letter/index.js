
import http from "@/utils/axios";

/** 导购--邀请函模板列表 */
export const getInviteConfigList = (data)=> http('invite/config-list',{
  method: "POST",
  data
})

/** 导购--邀请函模板详情 */
export const getInvitationInfo = (data)=> http('invite/config-detail',{
  method: "POST",
  data
})

/** 发送邀请函 */
export const sendInvite = (data)=> http('invite/create-record',{
  method: "POST",
  data
})
/** 查看邀请函 */
export const recordInfo = (data)=> http('invite/record-info',{
  method: "POST",
  data
})
// 获取用户信息
export const userInfoApi = (data) => http('customer/user-simple-info', {
  method: 'POST',
  data
})
// 删除邀请记录
export const deleteInvited = (data) => http('invite/record-delete', {
  method: 'POST',
  data
})