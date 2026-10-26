## Observer 元素曝光

### 介绍
监听内容是否在可视区域，使用`createIntersectionObserver`实现

### Props

| 参数 | 说明 | 类型 | 可选值 | 默认值 |
| ----- | ----- | ----- | ----- | ----- |
| thresholds | 一个数值数组，包含所有阈值。 | Array\<Number\> | -- | [0, 1] |
| viewport | 可视区域范围 | Object | -- | { bottom:0 } |

### Slots

| name  | 说明 |
| ----- | ----- |
| default | 展示内容 |

### Events

| name  | 说明 |
| ----- | ----- |
| change | 进入/离开展示区域触发 |

### 使用示例

查看demo`packages/demo/pages/observer/index`
