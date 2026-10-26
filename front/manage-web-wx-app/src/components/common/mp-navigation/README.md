## Navigation 自定义导航栏

### 介绍
为页面提供导航功能，常用于页面顶部。

### Props
参数 | 说明 | 类型 | 默认值
| ----- | ----- | ----- | ----- |
title | 标题 | string | ''
left-text | 左侧文案 | string | ''
right-text | 右侧文案 | string | ''
left-arrow | 是否显示左侧箭头 | boolean | false
fixed | 是否固定在顶部 | boolean | false
placeholder | 固定在顶部时是否开启占位 | boolean | false
z-index | 元素 z-index | number | 1
custom-style | 根节点自定义样式 | string | -
safe-area-inset-top | 是否留出顶部安全距离（状态栏高度） | boolean | true
opacity | 状态栏透明度 | number | 1
backgroundColor | 状态栏背景颜色 | string | #fff
useBeforeClick | 使用左侧返回点击之前的事件 | boolean | false
autoBack | 自动执行内部返回逻辑 | boolean | true



### Slots
名称 | 说明
| ----- | ----- |
title | 自定义标题
left | 自定义左侧区域内容
right | 自定义右侧区域内容

### Events
事件名 | 说明 | 参数
| ----- | ----- | ----- |
bind:before-click | 点击左侧按钮之前触发 | -
bind:click-left | 点击左侧按钮时触发 | -
bind:click-right | 点击右侧按钮时触发 | -

### Methods

| name  | 说明 |
| ----- | ----- |
| scroll | 滚动渐变 |


### 外部样式类
类名 | 说明
| ----- | ----- |
custom-class | 根节点样式类
title-class | 标题样式类

### 使用示例

查看demo`packages/demo/pages/navigation/index`
