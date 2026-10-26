### 介绍

---

crs-cy-components 组件库基于 element-ui 库二次封装，包含组件如下

```js
public
├─cy-cascader
├─cy-form
├─cy-input
├─cy-select
├─select-input
├─double-select
├─cy-table
├─cy-upload
├─upload-image
└─search-item
└─preview
```

### Install

---

````js
// git+ssh 私有化安装
npm install git+ssh://git@gitlab-ce.k8s.tools.vchangyi.com:32201/common/vue-toolbox.git#v1.0.14
````

### Quick Start

---

```js
import Vue from "vue"
import Element from "element-ui"
import CyComponents from "crs-cy-components"
import "crs-cy-components/dist/components.css" // 单独引入样式
 
Vue.use(Element)
Vue.use(CyComponents)
 
// 如果是 TS 需要加模块声明 
declare module "crs-cy-components";

```

### use

---

```js
<template>
  <div id="app">
    <cy-upload></cy-upload>
  </div>
</template>
```

### Tag

```js
// 通过 git tag 管理组件库版本，如要升级组件库，请务必加版本号
current tag v1.1.12
2020/11/19  v1.0.0
```

