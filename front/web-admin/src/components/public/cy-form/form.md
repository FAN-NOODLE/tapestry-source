# 表单组件

> 说明

- 除以下几个不支持和新增验证属性，其余同原生 el-form 组件相同
- 原生 inline 不支持
- 原生 hide-required-asterisk 不支持
- 原生 show-message 不支持
- 原生 inline-message 不支持

## 新增 form-item 内置校验参数

| 参数名称 | 说明                 | 类型            | 可选值 | 默认值 |
| -------- | -------------------- | --------------- | ------ | ------ |
| required | 是否必填             | Boolean         | -      | false  |
| email    | 校验邮箱             | Boolean         | -      | false  |
| phone    | 校验手机号           | Boolean         | -      | false  |
| idcard   | 校验身份证           | Boolean         | -      | false  |
| number   | 校验数字             | Boolean         | -      | false  |
| url      | 校验 http 链接       | Boolean         | -      | false  |
| zh       | 校验中文             | Boolean         | -      | false  |
| en       | 校验英文             | Boolean         | -      | false  |
| ennum    | 校验英文+数字        | Boolean         | -      | false  |
| en_num   | 校验英文+数字+下划线 | Boolean         | -      | false  |
| min      | 校验最小长度         | Number / String | -      | -      |
| max      | 校验最大长度         | Number / String | -      | -      |

- trigger='change' 默认校验为 blur 函数，若需其他需要传入，仅对 required 有效
- 支持自定义校验

## 原生 ElementUi 参数可用

| 参数                    | 说明                                                                                      | 类型    | 可选值                | 默认值 |
| ----------------------- | ----------------------------------------------------------------------------------------- | ------- | --------------------- | ------ |
| model                   | 表单数据对象                                                                              | object  | —                     | —      |
| rules                   | 表单验证规则                                                                              | object  | —                     | —      |
| label-position          | 表单域标签的位置，如果值为 left 或者 right 时，则需要设置 `label-width`                   | string  | right/left/top        | right  |
| label-width             | 表单域标签的宽度，例如 '50px'。作为 Form 直接子元素的 form-item 会继承该值。支持 `auto`。 | string  | —                     | —      |
| label-suffix            | 表单域标签的后缀                                                                          | string  | —                     | —      |
| status-icon             | 是否在输入框中显示校验结果反馈图标                                                        | boolean | —                     | false  |
| validate-on-rule-change | 是否在 `rules` 属性改变后立即触发一次验证                                                 | boolean | —                     | true   |
| size                    | 用于控制该表单内组件的尺寸                                                                | string  | medium / small / mini | —      |
| disabled                | 是否禁用该表单内的所有组件。若设置为 true，则表单内组件上的 disabled 属性不再生效         | boolean | —                     | false  |

## 组件方法说明

`特别说明`

`1.validateField、clearValidate 方法仅支持string类型参数`

| 方法名        | 说明                                                                                                                                                                 | 参数                                                                       |
| ------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------- |
| validate      | 对整个表单进行校验的方法，参数为一个回调函数。该回调函数会在校验结束后被调用，并传入两个参数：是否校验成功和未通过校验的字段。若不传入回调函数，则会返回一个 promise | Function(callback: Function(boolean, object))                              |
| validateField | 对部分表单字段进行校验的方法                                                                                                                                         | Function(props: array \| string, callback: Function(errorMessage: string)) |
| resetFields   | 对整个表单进行重置，将所有字段值重置为初始值并移除校验结果                                                                                                           | —                                                                          |
| clearValidate | 移除表单项的校验结果。传入待移除的表单项的 prop 属性或者 prop 组成的数组，如不传则移除整个表单的校验结果                                                             | Function(props: array \| string)                                           |

## 组件事件说明

| 事件名称 | 说明                   | 回调参数                                                   |
| -------- | ---------------------- | ---------------------------------------------------------- |
| validate | 任一表单项被校验后触发 | 被校验的表单项 prop 值，校验是否通过，错误消息（如果存在） |

## Form-Item 组件参数说明

`特别说明`

`1.form-item的rules本期暂不支持`

`2.required、min、max 请使用form-item内置校验`

| 参数        | 说明                                                                         | 类型    | 可选值                            | 默认值 |
| ----------- | ---------------------------------------------------------------------------- | ------- | --------------------------------- | ------ |
| prop        | 表单域 model 字段，在使用 validate、resetFields 方法的情况下，该属性是必填的 | string  | 传入 Form 组件的 `model` 中的字段 | —      |
| label       | 标签文本                                                                     | string  | —                                 | —      |
| label-width | 表单域标签的的宽度，例如 '50px'。支持 `auto`。                               | string  | —                                 | —      |
| required    | 是否必填，如不设置，则会根据校验规则自动生成                                 | boolean | —                                 | false  |
| error       | 表单域验证错误信息, 设置该值会使表单验证状态变为`error`，并显示该错误信息    | string  | —                                 | —      |  |  |

### Form-Item 组件方法说明

| 方法名        | 说明                                                 | 参数 |
| ------------- | ---------------------------------------------------- | ---- |
| resetField    | 对该表单项进行重置，将其值重置为初始值并移除校验结果 | -    |
| clearValidate | 移除该表单项的校验结果                               | -    |

#### 效果图

<f-image alt="" src="https://img-crs.vchangyi.com/组件说明样式资源/管理后台/ts/form.png" width="500px" />

#### 代码示例

```html
<template>
  <div class="counter">
    <layout-form type="mobile">
      <template slot="content">
        <cy-form ref="cform" :rules="rules" :model="ruleForm">
          <cy-form-item
            label="测试form-one测试form-one"
            prop="name"
            number
            required
          >
            <cy-input
              v-model="ruleForm.name"
              clearable
              width-size="medium"
            ></cy-input>
          </cy-form-item>

          <cy-form-item
            label="手机"
            prop="phone"
            desc="这是对输入的描述"
            required
            phone
            number
            max="11"
          >
            <cy-input v-model="ruleForm.phone" clearable></cy-input>
          </cy-form-item>
          <cy-form-item label="邮箱" prop="email" desc="这是对邮箱的描述" email>
            <cy-input
              v-model="ruleForm.email"
              width-size="medium"
              clearable
            ></cy-input>
          </cy-form-item>
          <cy-form-item
            label="收费方式"
            prop="fee"
            desc="这是对兴趣的描述"
            required
            trigger="change"
          >
            <el-switch
              v-model="ruleForm.fee"
              active-text="按月付费"
              inactive-text="按年付费"
            ></el-switch>
          </cy-form-item>
          <cy-form-item label="身份证" prop="idcard" required idcard>
            <cy-input v-model="ruleForm.idcard" width-size="medium"></cy-input>
          </cy-form-item>
          <cy-form-item
            label="兴趣"
            prop="faverate"
            desc="这是对兴趣的描述"
            required
            trigger="change"
          >
            <el-checkbox-group v-model="ruleForm.faverate">
              <el-checkbox label="美食/餐厅线上活动" name="type"></el-checkbox>
              <el-checkbox label="地推活动" name="type"></el-checkbox>
              <el-checkbox label="线下主题活动" name="type"></el-checkbox>
              <el-checkbox label="单纯品牌曝光" name="type"></el-checkbox>
            </el-checkbox-group>
          </cy-form-item>
          <cy-form-item
            label="兴趣"
            prop="sex"
            desc="性别咯"
            required
            trigger="change"
          >
            <el-radio-group disabled v-model="ruleForm.sex">
              <el-radio label="男" value="1"></el-radio>
              <el-radio label="女" value="2"></el-radio>
            </el-radio-group>
          </cy-form-item>
          <cy-form-item prop="date" label="日期" required trigger="change">
            <el-date-picker
              placeholder="选择日期"
              v-model="ruleForm.date"
            ></el-date-picker>
          </cy-form-item>

          <cy-form-item prop="date" label="日期区间" required trigger="change">
            <el-date-picker
              type="daterange"
              v-model="ruleForm.dateVal"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
            ></el-date-picker>
          </cy-form-item>
          <cy-form-item prop="date" label="时间" required trigger="change">
            <el-time-picker
              placeholder="选择时间"
              v-model="ruleForm.time"
            ></el-time-picker>
          </cy-form-item>
          <cy-form-item prop="date" label="时间" required trigger="change">
            <el-date-picker
              placeholder="选择时间"
              type="datetimerange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              v-model="ruleForm.times"
            ></el-date-picker>
          </cy-form-item>
          <cy-form-item label="活动形式" prop="action" required max="10">
            <cy-input type="textarea" v-model="ruleForm.action"></cy-input>
          </cy-form-item>

          <cy-form-item label="上传">
            <upload-image
              class="upload-template"
              :urls="ruleForm.avator"
              @change="picSuccess"
            ></upload-image>
          </cy-form-item>
        </cy-form>
      </template>
    </layout-form>
  </div>
</template>

<script lang="ts">
  import { Vue, Component } from "vue-property-decorator"
  import UploadImage from "@/components/business/upload-image/index.vue"
  @Component({
    components: {
      UploadImage,
    },
  })
  export default class FormDemo extends Vue {
    private ruleForm = {
      name: "",
      phone: "",
      email: "",
      faverate: ["美食/餐厅线上活动"],
      fee: false,
      sex: "",
      idcard: "",
      action: "",
      date: "",
      dateVal: [],
      times: [],
      time: "",
      timeVal: [new Date(2016, 9, 10, 8, 40), new Date(2016, 9, 10, 9, 40)],
      select: "",
      avator: "",
    }
    private options = [
      {
        value: "选项1",
        label: "黄金糕",
      },
      {
        value: "选项2",
        label: "双皮奶",
      },
      {
        value: "选项3",
        label: "蚵仔煎",
      },
      {
        value: "选项4",
        label: "龙须面",
      },
      {
        value: "选项5",
        label: "北京烤鸭",
      },
    ]
    private rules = {
      name: [
        {
          validator: this.validateself,
          trigger: "blur",
        },
      ],
    }
    private validateself(rule: any, value: any, callback: any) {
      value.length > 5 ? callback(new Error("五字不行")) : callback()
    }
    private picSuccess(res: any) {
      this.$set(this.ruleForm, "avator", res)
    }
  }
</script>
```
