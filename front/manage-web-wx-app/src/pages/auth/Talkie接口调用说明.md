### 跳转地址
- appid: wxbd40b04da240ef2b(测试)，wx19ba52f691c4d289(正式)
- 页面地址: pages/auth/index
### 跳转参数
|参数名|参数类型|默认值|必填|示例|参数说明|
|-|-|-|-|-|-|
|target|枚举| - |是|materialList-素材列表 <br/><span style="color: red;">goodsMaterialDetail-【新】商品素材详情|目标页面</span>
|query|字符串编码| - |是|当target为materialList时：encodeURIComponent('from=TALK&current=1&sku_item={skucode或者skucode_颜色}') <br/> <span style="color: red;">当target为goodsMaterialDetail时：encodeURIComponent('from=TALK&upc_code={upc_code}')</span>|查询参数，from固定为TALK，current固定值1，sku_item搜索商品, upc_code商品upc_code|


### 示例
- 跳转素材列表
  ```javascript
  wx.navigateToMiniProgram({
    appId: 'wxbd40b04da240ef2b',
    path: `pages/auth/index?target=materialList&query=${encodeURIComponent('from=TALK&current=1&sku_item=1234')}`,
  })
  ```
- 跳转商品素材详情
  ```javascript
  wx.navigateToMiniProgram({
    appId: 'wxbd40b04da240ef2b',
    path: `pages/auth/index?target=goodsMaterialDetail&query=${encodeURIComponent('from=TALK&upc_code=198685374179')}`,
  })
  ```