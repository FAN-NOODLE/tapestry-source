import commonRequest from '../utils/common-request'


/**
 * @description 获取素材标签
 */
const getInventoryListApi = commonRequest({
  url: 'inventory/list',
  method: 'POST'
})
/**
 * @description: 判断是否是首次打开
 */
const getFirstOpenApi = commonRequest({
  url: 'common/check-first-open',
  method: 'POST'
})
export {
  getInventoryListApi,
  getFirstOpenApi
}
