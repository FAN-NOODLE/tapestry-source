import commonRequest from '../utils/common-request'


// 群成员列表
const chatGroupMembers = commonRequest({
  url: 'chat-group/members_info',
  method: 'POST'
})

// 社群画像列表
const chatIndexList = commonRequest({
  url: 'chat-group/index_list',
  method: 'POST'
})

// // 获取数据更新时间
// const getRefreshTime = commonRequest({
//   url: 'chat-group/data-update-time',
//   method: 'POST'
// })

export { chatGroupMembers, chatIndexList }
