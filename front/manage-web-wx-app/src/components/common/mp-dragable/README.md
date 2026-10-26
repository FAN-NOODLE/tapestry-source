## Dragable 拖拽组件

### 介绍
解决小程序图片拖拽排序问题

### Props

| 参数 | 说明 | 类型 | 可选值 | 默认值 |
| ----- | ----- | ----- | ----- | ----- |
| list | 数据列表 | Array | -- | [] |
| selector | 元素选择器 | String | -- | '.drag-item' |
| disabled | 禁止拖拽 | Boolean | true\|false | false |

### Slots

| name  | 说明 |
| ----- | ----- |
| default | 默认插槽 |
| movable-view | 拖拽元素 |

### Events

| name  | 说明 |
| ----- | ----- |
| dragBefore | 拖拽前事件 |
| change | 拖拽完成 |


### 外部样式类

| 类名  | 说明 |
| ----- | ----- |
| custom-class | 根结点样式类 |

### 使用示例

查看demo`packages/demo/pages/drag/index`
