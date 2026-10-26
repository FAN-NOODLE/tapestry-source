## Transition 动画

### 介绍
使元素从一种样式逐渐变化为另一种样式的效果。

### Props
参数 | 说明 | 类型 | 默认值
| ----- | ----- | ----- | ----- |
name | 动画类型 | string | fade
show | 是否展示组件 | boolean | true
duration | 动画时长，单位为毫秒 | number | object | 300
custom-style | 自定义样式 | string | -

### Events
事件名 | 说明 | 参数
| ----- | ----- | ----- |
bind:before-enter | 进入前触发 | -
bind:enter | 进入中触发 | -
bind:after-enter | 进入后触发 | -
bind:before-leave | 离开前触发 | -
bind:leave | 离开中触发 | -
bind:after-leave | 离开后触发 | -


### 外部样式类

类名 | 说明
| ----- | ----- |
custom-class | 根节点样式类
enter-class | 定义进入过渡的开始状态。在元素被插入之前生效，在元素被插入之后的下一帧移除。
enter-active-class | 定义进入过渡生效时的状态。在整个进入过渡的阶段中应用，在元素被插入之前生效，在过渡/动画完成之后移除。这个类可以被用来定义进入过渡的过程时间，延迟和曲线函数。
enter-to-class | 定义进入过渡的结束状态。在元素被插入之后下一帧生效 (与此同时 enter-class 被移除)，在过渡/动画完成之后移除。
leave-class | 定义离开过渡的开始状态。在离开过渡被触发时立刻生效，下一帧被移除。
leave-active-class | 定义离开过渡生效时的状态。在整个离开过渡的阶段中应用，在离开过渡被触发时立刻生效，在过渡/动画完成之后移除。这个类可以被用来定义离开过渡的过程时间，延迟和曲线函数。
leave-to-class | 定义离开过渡的结束状态。在离开过渡被触发之后下一帧生效 (与此同时 leave-class 被删除)，在过渡/动画完成之后移除。

### 使用示例
查看demo`packages/demo/pages/transition/index`
