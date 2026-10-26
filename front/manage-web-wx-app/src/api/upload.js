import commonRequest from '../utils/common-request'

// 获取cos token
export const getCosUploadTokenApi = commonRequest({
  url: 'upload',
  method: 'POST'
})
