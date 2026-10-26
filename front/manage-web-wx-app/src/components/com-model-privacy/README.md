## ModelPrivacy 用户隐私协议

### 介绍
用来处理微信隐私协议逻辑

### Props

| 参数 | 说明 | 类型 | 可选值 | 默认值 |
| ----- | ----- | ----- | ----- | ----- |
| openType | 打开方式, get: 方案一主动获取；listening: 监听页面api | string | listening\|get | listening |

### Events

| name  | 说明 |
| ----- | ----- |
| result | 用户对于隐私协议的同意或者拒绝，e.detail.value |

### 外部样式类

| 类名  | 说明 |
| ----- | ----- |
| custom-class | 根结点样式类 |

### 使用示例

```json
{
	"usingComponents": {
		"com-model-privacy": "components/business/com-model-privacy/index"
	},
	"__usePrivacyCheck__": true // 20230915之后可以不用配置这个参数
}
```

使用监听的方式，可以监听页面内相关的隐私操作
```html
<com-model-privacy></com-model-privacy>
```

使用主动获取方式，可以监听页面内相关的隐私操作
```html
<com-model-privacy open-type="get"></com-model-privacy>
```

获取同意或者不同意的结果
```html
<com-model-privacy bind:result="handlePrivacyResult"></com-model-privacy>
```

