import { createModuleRequest as wxworkModuleRequest } from './wxwork-request'

// const systemInfo = wx.getSystemInfoSync()
const request = wxworkModuleRequest()
const commonRequest = ({ url, method, apiRoot, apiVersion }) => (data, app) =>
  request(
    url,
    {
      method,
      data,
      apiRoot,
      apiVersion
    },
    app
  )

export default commonRequest
