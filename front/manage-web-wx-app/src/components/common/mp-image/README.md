## ComponentName 组件名称

### 介绍
组件说明

### Props

| 属性 | 类型 | 默认值 | 是否必填 | 说明 | 最低版本 |
| ----- | ----- | ----- | ----- | ----- | ----- |
| src | string | -- | 否 | 图片资源地址 | 1.0.0 |
| mode | string | scaleToFill | 否 | 图片裁剪、缩放的模式 | 1.0.0 |
| show-menu-by-longpress | boolean | false | 否 | 长按图片显示发送给朋友、收藏、保存图片、搜一搜、打开名片/前往群聊/打开小程序（若图片中包含对应二维码或小程序码）的菜单。 | 2.7.0 |
| binderror | eventhandle | -- | 否 | 当错误发生时触发，event.detail = {errMsg} | 1.0.0 |
| bindload | eventhandle | -- | 否 | 当图片载入完毕时触发，event.detail = {height, width} | 1.0.0 |
| preview | boolean | false | 否 | 是否预览 | 1.0.0 |
| previewList | Array | [] | 否 | 预览图片列表 | 1.0.0 |
| viewport | object | {bottom:0} | 否 | 懒加载限制区域 | 1.0.0 |
| customStyle | string | '' | 否 | 自定义样式 | 1.0.0 |
| useErrorSlot | boolean | false | 否 | 是否使用 error 插槽 | 1.0.0 |
| useLoadingSlot | boolean | false | 否 | 是否使用 loading 插槽 | 1.0.0 |
| round | boolean | false | 否 | 是否显示为圆形 | 1.0.0 |
| radius | string | 0 | 否 | 圆角大小，默认单位为rpx | 1.0.0 |

### Skyline 特有属性
| 属性 | 类型 | 默认值 | 是否必填 | 说明 |
| ----- | ----- | ----- | ----- | ----- |
| fade-in | boolean | false | 否 | 是否渐显 |

### WebView 特有属性
| 属性 | 类型 | 默认值 | 是否必填 | 说明 | 最低版本 |
| ----- | ----- | ----- | ----- | ----- | ----- |
| webp | boolean | false | 否 | 默认不解析 webP 格式，只支持网络资源 | 2.9.0 |
| lazy-load | boolean | false | 否 | 图片懒加载，在即将进入一定范围（上下三屏）时才开始加载 | 1.5.0 |

### mode 可选值
| 合法值 | 说明 |
| ----- | ----- |
| scaleToFill | 缩放模式，不保持纵横比缩放图片，使图片的宽高完全拉伸至填满 image 元素 |
| aspectFit | 缩放模式，保持纵横比缩放图片，使图片的长边能完全显示出来。也就是说，可以完整地将图片显示出来。 |
| aspectFill | 缩放模式，保持纵横比缩放图片，只保证图片的短边能完全显示出来。也就是说，图片通常只在水平或垂直方向是完整的，另一个方向将会发生截取。 |
| widthFix | 缩放模式，宽度不变，高度自动变化，保持原图宽高比不变 |
| heightFix | 缩放模式，高度不变，宽度自动变化，保持原图宽高比不变 |
| top | 裁剪模式，不缩放图片，只显示图片的顶部区域。仅 Webview 支持。 |
| bottom | 裁剪模式，不缩放图片，只显示图片的底部区域。仅 Webview 支持。 |
| center | 裁剪模式，不缩放图片，只显示图片的中间区域。仅 Webview 支持。 |
| left | 裁剪模式，不缩放图片，只显示图片的左边区域。仅 Webview 支持。 |
| right | 裁剪模式，不缩放图片，只显示图片的右边区域。仅 Webview 支持。 |
| top left | 裁剪模式，不缩放图片，只显示图片的左上边区域。仅 Webview 支持。 |
| top right | 裁剪模式，不缩放图片，只显示图片的右上边区域。仅 Webview 支持。 |
| bottom left | 裁剪模式，不缩放图片，只显示图片的左下边区域。仅 Webview 支持。 |
| bottom right | 裁剪模式，不缩放图片，只显示图片的右下边区域。仅 Webview 支持。 |

### Events

| 名称  | 说明 |
| ----- | ----- |
| load | 加载成功 |
| error | 加载失败 |
| click | 点击事件 |

### Slot

| 名称  | 说明 |
| ----- | ----- |
| loading | 自定义加载中的提示内容 |
| error | 自定义加载失败时的提示内容 |

### 外部样式类

| 类名  | 说明 |
| ----- | ----- |
| custom-class | 根节点样式类 |
| image-class | 图片样式类 |
| loading-class | loading 样式类 |
| error-class | error 样式类 |


### 用例展示

```js
Page({
	data: {
    title: 'image',
		image: 'https://img1.baidu.com/it/u=1308526554,393095433&fm=253&fmt=auto&app=120&f=JPEG?w=800&h=800',
		list: ['scaleToFill', 'aspectFit', 'aspectFill', 'widthFix', 'top', 'bottom', 'center', 'left', 'right', 'top left', 'top right', 'bottom left', 'bottom right'],
  },
})
```
基础用法
```html
<mp-image custom-class="image" src="{{image}}"></mp-image>
```
图片填充模式
```html
<view class="list">
	<block wx:for="{{list}}" wx:key="index">
		<view class="item">
			<mp-image custom-class="item-image" src="{{image}}" mode="{{item}}"></mp-image>
			<view class="item-title">{{item}}</view>
		</view>
	</block>
</view>
```
圆形图片
```html
<view class="list">
	<block wx:for="{{list}}" wx:key="index">
		<view class="item">
			<mp-image custom-class="item-image" src="{{image}}" round mode="{{item}}"></mp-image>
			<view class="item-title">{{item}}</view>
		</view>
	</block>
</view>
```
图片圆角
```html
<view class="list">
	<block wx:for="{{list}}" wx:key="index">
		<view class="item">
			<mp-image custom-class="item-image" src="{{image}}" radius="20rpx" mode="{{item}}"></mp-image>
			<view class="item-title">{{item}}</view>
		</view>
	</block>
</view>
```
图片长按识别
```html
<mp-image custom-class="image" show-menu-by-longpress src="{{image}}"></mp-image>
```
图片预览
```html
<mp-image custom-class="image" preview src="{{image}}"></mp-image>
```
图片加载失败
```html
<mp-image custom-class="image" src="{{'image'}}"></mp-image>
```
图片懒加载
```html
<mp-image custom-class="image" lazy-load src="{{image}}"></mp-image>
```
