# 安装依赖管理工具'yarn'
npm install -g yarn

# 安装依赖
yarn install
或
npm install

# 启动 dev 服务，并编译打包测试环境 生成 ./dist 目录
# 开发环境
yarn dev
# 线上环境
yarn prod

# 打包正式环境 生成 ./build 目录
yarn build
# 打包正式环境 生成 ./build 目录
npm run build

# 新建文件
1. yarn create-page
2. 输入页面名称？  譬如：customer  or  customer/customerList
3. 是否需要配置文件(.json)？ Y/N

# 开发环境用聚襄小程序

### 测试小程序
    账号：xiaomili2020@126.com    密码：REDACTED_FOR_SOURCE_SHARING

### 百利助手小程序提审腾讯所需
    账号：bailisy_BL2@163.com  密码：REDACTED_FOR_SOURCE_SHARING

### 智慧导购测试提审腾讯所需
    账号：13720452952 验证码：1234

### 智慧导购测试beta提审腾讯所需
    账号：13720452952 验证码：1234

### 畅鲟智慧导购生产环境提审腾讯所需
    账号：13227766759 验证码：1234

### 畅鲟智慧导购beta环境提审腾讯所需
    账号：13227766759 验证码：1234

### 百利提审腾讯所需
    账号：18502987082  验证码：1234

### 代码库：
    https://gitlab-ce.k8s.tools.vchangyi.com/crs/crs-manage-wxapp
    
### wiki地址
https://shop-api-dev-crs.vchangyi.com/gw-assistant/doc/app    


### 关于对外页面地址


### 页面地址变更

| 变更前 | 变更后 |
|------|------|
| pages/customer-allocation/customer-allocation | packages/allocation/pages/customer-list/customer-list |
| pages/customer-allocation/customer-select-list/customer-select-list | packages/allocation/pages/customer-select-list/customer-select-list |
| pages/customer-allocation/member-select-list/member-select-list | packages/allocation/pages/member-select-list/member-select-list |
| pages/allocation-record/allocation-record | packages/allocation/pages/record/record |
| pages/allocation-record-detail/allocation-record-detail | packages/allocation/pages/record-detail/record-detail |
| pages/member-search/member-search | packages/customer/pages/member-search/member-search |
"pages/member-search-result/member-search-result | packages/customer/pages/member-search-result/member-search-result

### 小程序提交发布
```
# uat环境
yarn upload:test [版本号]

# prod环境
yarn upload:prod [版本号]
```