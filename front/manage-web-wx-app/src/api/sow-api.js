import commonRequest from '../utils/common-request'

// 【播种】- 列表
const getSowTaskListApi = commonRequest({
  url: 'customer/sow-task/list',
  method: 'POST'
})

// 【播种】- 创建
const createSowTaskApi = commonRequest({
  url: 'customer/sow-task/new-create',
  method: 'POST'
})

// 【播种】- 删除
const deleteSowTaskApi = commonRequest({
  url: 'customer/sow-task/delete',
  method: 'POST'
})

// 【播种】- 编辑
const updateSowTaskApi = commonRequest({
  url: 'customer/sow-task/update',
  method: 'POST'
})

// 【播种】- 完成
const finishedSowTaskApi = commonRequest({
  url: 'customer/sow-task/finished',
  method: 'POST'
})
// 【播种】- 完成
const sowTempListApi = commonRequest({
  url: 'customer/sow-task/template-list',
  method: 'POST'
})

export {
  getSowTaskListApi,
  createSowTaskApi,
  updateSowTaskApi,
  finishedSowTaskApi,
  deleteSowTaskApi,
  sowTempListApi
}
