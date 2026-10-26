# 级联选择器组件

> 说明：UI 订制级联选择器，满足设计需求，在 Element-ui 基础上增加宽度 size

## 组件参数说明

> 额外新增参数

| 参数名称   | 说明     | 类型   | 可选值                       | 默认值 |
| ---------- | -------- | ------ | ---------------------------- | ------ |
| width-size | 宽度大小 | string | small(280px) / medium(400px) | small  |

> Element-UI 原生参数

| 参数            | 说明                                                                                                | 类型                    | 可选值                | 默认值    |
| --------------- | --------------------------------------------------------------------------------------------------- | ----------------------- | --------------------- | --------- |
| value / v-model | 选中项绑定值                                                                                        | -                       | —                     | —         |
| options         | 可选项数据源，键名可通过 `Props` 属性配置                                                           | array                   | —                     | —         |
| props           | 配置选项，具体见下表                                                                                | object                  | —                     | —         |
| size            | 尺寸                                                                                                | string                  | medium / small / mini | —         |
| placeholder     | 输入框占位文本                                                                                      | string                  | —                     | 请选择    |
| disabled        | 是否禁用                                                                                            | boolean                 | —                     | false     |
| clearable       | 是否支持清空选项                                                                                    | boolean                 | —                     | false     |
| show-all-levels | 输入框中是否显示选中值的完整路径                                                                    | boolean                 | —                     | true      |
| collapse-tags   | 多选模式下是否折叠 Tag                                                                              | boolean                 | -                     | false     |
| separator       | 选项分隔符                                                                                          | string                  | —                     | 斜杠' / ' |
| filterable      | 是否可搜索选项                                                                                      | boolean                 | —                     | —         |
| filter-method   | 自定义搜索逻辑，第一个参数是节点`node`，第二个参数是搜索关键词`keyword`，通过返回布尔值表示是否命中 | function(node, keyword) | -                     | -         |
| debounce        | 搜索关键词输入的去抖延迟，毫秒                                                                      | number                  | —                     | 300       |
| before-filter   | 筛选之前的钩子，参数为输入的值，若返回 false 或者返回 Promise 且被 reject，则停止筛选               | function(value)         | —                     | —         |
| popper-class    | 自定义浮层类名                                                                                      | string                  | —                     | —         |

## 组件事件说明

| 事件名称       | 说明                          | 回调参数                      |
| -------------- | ----------------------------- | ----------------------------- |
| change         | 当选中节点变化时触发          | 选中节点的值                  |
| expand-change  | 当展开节点发生变化时触发      | 各父级选项值组成的数组        |
| blur           | 当失去焦点时触发              | (event: Event)                |
| focus          | 当获得焦点时触发              | (event: Event)                |
| visible-change | 下拉框出现/隐藏时触发         | 出现则为 true，隐藏则为 false |
| remove-tag     | 在多选模式下，移除 Tag 时触发 | 移除的 Tag 对应的节点的值     |

## 组件方法说明

| 方法名          | 说明           | 参数                                          |
| --------------- | -------------- | --------------------------------------------- |
| getCheckedNodes | 获取选中的节点 | (leafOnly) 是否只是叶子节点，默认值为 `false` |

## 组件插槽说明

| 名称  | 说明                                                                            |
| ----- | ------------------------------------------------------------------------------- |
| -     | 自定义备选项的节点内容，参数为 { node, data }，分别为当前节点的 Node 对象和数据 |
| empty | 无匹配选项时的内容                                                              |

## props 相应的参数说明

| 参数          | 说明                                                                                               | 类型                                                                                     | 可选值        | 默认值     |
| ------------- | -------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ------------- | ---------- |
| expandTrigger | 次级菜单的展开方式                                                                                 | string                                                                                   | click / hover | 'click'    |
| multiple      | 是否多选                                                                                           | boolean                                                                                  | -             | false      |
| checkStrictly | 是否严格的遵守父子节点不互相关联                                                                   | boolean                                                                                  | -             | false      |
| emitPath      | 在选中节点改变时，是否返回由该节点所在的各级菜单的值所组成的数组，若设置 false，则只返回该节点的值 | boolean                                                                                  | -             | true       |
| lazy          | 是否动态加载子节点，需与 lazyLoad 方法结合使用                                                     | boolean                                                                                  | -             | false      |
| lazyLoad      | 加载动态数据的方法，仅在 lazy 为 true 时有效                                                       | function(node, resolve)，`node`为当前点击的节点，`resolve`为数据加载完成的回调(必须调用) | -             | -          |
| value         | 指定选项的值为选项对象的某个属性值                                                                 | string                                                                                   | —             | 'value'    |
| label         | 指定选项标签为选项对象的某个属性值                                                                 | string                                                                                   | —             | 'label'    |
| children      | 指定选项的子选项为选项对象的某个属性值                                                             | string                                                                                   | —             | 'children' |
| disabled      | 指定选项的禁用为选项对象的某个属性值                                                               | string                                                                                   | —             | 'disabled' |
| leaf          | 指定选项的叶子节点的标志位为选项对象的某个属性值                                                   | string                                                                                   | —             | 'leaf'     |

#### 效果图

<f-image alt="" src="https://img-crs.vchangyi.com/组件说明样式资源/管理后台/ts/cascader.png" width="500px" />

#### 代码示例

```html
<template>
  <div class="counter">
    <div class="template">
      <cy-cascader
        v-model="cases"
        :options="options"
        :props="{ expandTrigger: 'hover' }"
        @change="changeCases"
        placeholder="宽度small级联选择器(hover触发子级)"
        clearable
      ></cy-cascader>
    </div>
    <div class="template">
      <cy-cascader
        width-size="medium"
        v-model="cases"
        :options="options"
        :props="{ checkStrictly: true }"
        @change="changeCases"
        placeholder="宽度medium级联选择器(带单选)"
        clearable
      ></cy-cascader>
    </div>
  </div>
</template>

<script lang="ts">
  import { Vue, Component } from "vue-property-decorator"
  @Component
  export default class extends Vue {
    private cases = []
    private options = [
      {
        value: "zhinan",
        label: "指南",
        children: [
          {
            value: "shejiyuanze",
            label: "设计原则",
            children: [
              {
                value: "yizhi",
                label: "一致",
              },
              {
                value: "fankui",
                label: "反馈",
              },
              {
                value: "xiaolv",
                label: "效率",
              },
              {
                value: "kekong",
                label: "可控",
              },
            ],
          },
          {
            value: "daohang",
            label: "导航",
            children: [
              {
                value: "cexiangdaohang",
                label: "侧向导航",
              },
              {
                value: "dingbudaohang",
                label: "顶部导航",
              },
            ],
          },
        ],
      },
      {
        value: "zujian",
        label: "组件",
        children: [
          {
            value: "basic",
            label: "Basic",
            children: [
              {
                value: "layout",
                label: "Layout 布局",
              },
              {
                value: "color",
                label: "Color 色彩",
              },
              {
                value: "typography",
                label: "Typography 字体",
              },
              {
                value: "icon",
                label: "Icon 图标",
              },
              {
                value: "button",
                label: "Button 按钮",
              },
            ],
          },
          {
            value: "form",
            label: "Form",
            children: [
              {
                value: "radio",
                label: "Radio 单选框",
              },
              {
                value: "checkbox",
                label: "Checkbox 多选框",
              },
              {
                value: "input",
                label: "Input 输入框",
              },
              {
                value: "input-number",
                label: "InputNumber 计数器",
              },
              {
                value: "select",
                label: "Select 选择器",
              },
              {
                value: "cascader",
                label: "Cascader 级联选择器",
              },
              {
                value: "switch",
                label: "Switch 开关",
              },
              {
                value: "slider",
                label: "Slider 滑块",
              },
              {
                value: "time-picker",
                label: "TimePicker 时间选择器",
              },
              {
                value: "date-picker",
                label: "DatePicker 日期选择器",
              },
              {
                value: "datetime-picker",
                label: "DateTimePicker 日期时间选择器",
              },
              {
                value: "upload",
                label: "Upload 上传",
              },
              {
                value: "rate",
                label: "Rate 评分",
              },
              {
                value: "form",
                label: "Form 表单",
              },
            ],
          },
          {
            value: "data",
            label: "Data",
            children: [
              {
                value: "table",
                label: "Table 表格",
              },
              {
                value: "tag",
                label: "Tag 标签",
              },
              {
                value: "progress",
                label: "Progress 进度条",
              },
              {
                value: "tree",
                label: "Tree 树形控件",
              },
              {
                value: "pagination",
                label: "Pagination 分页",
              },
              {
                value: "badge",
                label: "Badge 标记",
              },
            ],
          },
        ],
      },
    ]
    private changeCases() {
      console.log(this.cases)
    }
  }
</script>
<style lang="scss" scoped>
  .template {
    padding: 15px 30px;
  }
</style>
```
