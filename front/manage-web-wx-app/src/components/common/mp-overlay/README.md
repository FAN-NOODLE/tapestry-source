## Overlay 遮罩层

### 介绍
创建一个遮罩层，用于强调特定的页面元素，并阻止用户进行其他操作。

### Props
参数 | 说明 | 类型 | 默认值
| ----- | ----- | ----- | ----- |
show | 是否展示遮罩层 | boolean | false
z-index | z-index 层级 | string | number | 1
duration | 动画时长，单位秒 | string | number | 0.3
class-name | 自定义类名 | string | -
custom-style | 自定义样式 | string | -
lock-scroll v1.7.3 | 是否锁定背景滚动，锁定时蒙层里的内容也将无法滚动 | boolean | true
root-portal v1.10.14 | 是否从页面中脱离出来，用于解决各种 fixed 失效问题，微信基础库 >= 2.25.2  | boolean | false


### Slots

| name  | 说明 |
| ----- | ----- |
| - | 默认插槽，用于在遮罩层上方嵌入内容 |

### Events
事件名 | 说明 | 回调参数
| ----- | ----- | ----- |
bind:click | 点击时触发 | -

### 外部样式类

| 类名  | 说明 |
| ----- | ----- |
custom-class | 根节点样式类

### 使用示例
