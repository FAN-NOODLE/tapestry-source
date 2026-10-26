## Footer 底部安全区兼容组件

### 介绍
组件说明

### Props

| 参数 | 说明 | 类型 | 可选值 | 默认值 |
| ----- | ----- | ----- | ----- | ----- |
| fixed | 是否固定在底部 | boolean | true\|false | true |
| backgroundColor | 背景颜色 | string | -- | #ffffff |

### Slots

| name  | 说明 |
| ----- | ----- |
| default | 默认插槽 |


### 外部样式类

| 类名  | 说明 |
| ----- | ----- |
| custom-class | 根结点样式类 |

### 使用示例
基础用法
```html
<mp-footer>
	<view class="footer">
		<mp-button type="primary" block>按钮</mp-button>
	</view>
</mp-footer>
```
自定义背景色
```html
<mp-footer background-color="#000">
	<view class="footer">
		<mp-button type="primary" block>按钮</mp-button>
	</view>
</mp-footer>
```
