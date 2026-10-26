import commonRequest from '../utils/common-request'


/**
 * @description 获取素材标签
 */
const getMaterialTagListApi = commonRequest({
  url: 'source/category/list',
  method: 'POST'
})

/**
 * @description 获取上传素材权限
 */
const getUploadAuthApi = commonRequest({
  url: 'source/upload-auth',
  method: 'POST'
})

/**
 * @description 获取素材列表
 */
const getMaterialListApi = commonRequest({
  url: 'source/select-list',
  method: 'POST'
})

/**
 * @description 获取素材分组列表
 */
const getMaterialGroupListApi = commonRequest({
  url: 'source/group-relation/get-list-for-store',
  method: 'POST'
})

/**
 * @description 获取素材压缩结果
 */
const getCompressInfoApi = commonRequest({
  url: 'source/compress',
  method: 'POST'
})

/**
 * @description 保存素材
 */
const uploadMediaApi = commonRequest({
  url: 'source/upload-media',
  method: 'POST'
})

/**
 * @description 获取上传凭证
 */
const getCosTokenApi = commonRequest({
  url: 'source/credential',
  method: 'POST'
})

/**
 * @description 获取素材数量
 */
const getMaterialCountApi = commonRequest({
  url: 'source/get-new-num',
  method: 'POST'
})

/**
 * @description 数据上报
 */
const reportTimesApi = commonRequest({
  url: 'source/report-times',
  method: 'POST'
})

/**
 * @description 压缩图片
 */
const compressFileApi = commonRequest({
  url: 'source/compress-file',
  method: 'POST'
})

/**
 * @description 获取mediaInfo
 */
const getMediaInfoApi = commonRequest({
  url: 'source/get-media-info',
  method: 'POST'
})

/**
 * @description 获取新素材数量
 */
const getMaterialNewNumApi = commonRequest({
  url: 'source/get-new-num',
  method: 'POST'
})

/**
 * @description 首次进入
 */
const checkFirstApi = commonRequest({
  url: 'source/check-first',
  method: 'POST'
})
/**
 * @description 修改视频时间
 */
const makeCoverApi = commonRequest({
  url: 'source/make-cover',
  method: 'POST'
})
/**
 * @description 获取品牌素材
 */
const getBrandMaterialListApi = commonRequest({
  url: 'source/brand/search-list',
  method: 'POST'
})
/**
 * @description 获取品牌素材
 */
const getCollectMaterialListApi = commonRequest({
  // url: 'source/brand/search-list',
  url: 'source/brand/collect/search-list',
  method: 'POST'
})
/**
 * @description 获取品牌素材
 */
const getBrandMaterialDetailApi = commonRequest({
  url: 'source/brand/get-cache-info',
  method: 'POST'
})
/**
 * @description 获取收藏列表
 */
const getCollectListApi = commonRequest({
  url: 'source/collect-select-list',
  method: 'POST'
})
/**
 * @description 收藏
 */
const collectApi = commonRequest({
  url: 'source/collect',
  method: 'POST'
})
/**
 * @description 获取AI识别关联商品列表
 */
const getTalkieImageSearchApi = commonRequest({
  url: 'source/brand/talkie-image-search',
  method: 'POST'
})


export {
  getMaterialTagListApi,
  getUploadAuthApi,
  getMaterialListApi,
  getMaterialGroupListApi,
  getCompressInfoApi,
  uploadMediaApi,
  getCosTokenApi,
  getMaterialCountApi,
  reportTimesApi,
  compressFileApi,
  getMediaInfoApi,
  getMaterialNewNumApi,
  checkFirstApi,
  makeCoverApi,
  getBrandMaterialListApi,
  getBrandMaterialDetailApi,
  getCollectListApi,
  collectApi,
  getCollectMaterialListApi,
  getTalkieImageSearchApi
}
