# 下拉选择组件

> 说明：UI 订制下拉框，满足 UI 需求，在 Element-ui 基础上订制宽度 size,和带有输入框的筛选功能

## 组件参数说明

> 额外新增参数

| 参数         | 类型    | 说明                       | 可选值                     | 默认值 |
| ------------ | ------- | -------------------------- | -------------------------- | ------ |
| width-size   | String  | 宽度大小型号               | medium(400px)/small(280px) | small  |
| search-input | Boolean | 下拉框是否展示筛选的输入框 | —                          | false  |

> Element-UI 原声参数

| 参数                  | 说明                                                                           | 类型                      | 可选值            | 默认值     |
| --------------------- | ------------------------------------------------------------------------------ | ------------------------- | ----------------- | ---------- |
| value / v-model       | 绑定值                                                                         | boolean / string / number | —                 | —          |
| multiple              | 是否多选                                                                       | boolean                   | —                 | false      |
| disabled              | 是否禁用                                                                       | boolean                   | —                 | false      |
| value-key             | 作为 value 唯一标识的键名，绑定值为对象类型时必填                              | string                    | —                 | value      |
| size                  | 输入框尺寸                                                                     | string                    | medium/small/mini | —          |
| clearable             | 是否可以清空选项                                                               | boolean                   | —                 | false      |
| collapse-tags         | 多选时是否将选中值按文字的形式展示                                             | boolean                   | —                 | false      |
| multiple-limit        | 多选时用户最多可以选择的项目数，为 0 则不限制                                  | number                    | —                 | 0          |
| name                  | select input 的 name 属性                                                      | string                    | —                 | —          |
| autocomplete          | select input 的 autocomplete 属性                                              | string                    | —                 | off        |
| placeholder           | 占位符                                                                         | string                    | —                 | 请选择     |
| filterable            | 是否可搜索                                                                     | boolean                   | —                 | false      |
| allow-create          | 是否允许用户创建新条目，需配合 `filterable` 使用                               | boolean                   | —                 | false      |
| filter-method         | 自定义搜索方法                                                                 | function                  | —                 | —          |
| remote                | 是否为远程搜索                                                                 | boolean                   | —                 | false      |
| remote-method         | 远程搜索方法                                                                   | function                  | —                 | —          |
| loading               | 是否正在从远程获取数据                                                         | boolean                   | —                 | false      |
| loading-text          | 远程加载时显示的文字                                                           | string                    | —                 | 加载中     |
| no-match-text         | 搜索条件无匹配时显示的文字，也可以使用`slot="empty"`设置                       | string                    | —                 | 无匹配数据 |
| no-data-text          | 选项为空时显示的文字，也可以使用`slot="empty"`设置                             | string                    | —                 | 无数据     |
| popper-class          | Select 下拉框的类名                                                            | string                    | —                 | —          |
| reserve-keyword       | 多选且可搜索时，是否在选中一个选项后保留当前的搜索关键词                       | boolean                   | —                 | false      |
| default-first-option  | 在输入框按下回车，选择第一个匹配项。需配合 `filterable` 或 `remote` 使用       | boolean                   | -                 | false      |
| popper-append-to-body | 是否将弹出框插入至 body 元素。在弹出框的定位出现问题时，可将该属性设置为 false | boolean                   | -                 | true       |
| automatic-dropdown    | 对于不可搜索的 Select，是否在输入框获得焦点后自动弹出选项菜单                  | boolean                   | -                 | false      |

## 组件事件说明

| 事件名称        | 说明                                     | 回调参数                      |
| --------------- | ---------------------------------------- | ----------------------------- |
| @change         | 选中值发生变化时触发                     | 目前的选中值                  |
| @visible-change | 下拉框出现/隐藏时触发                    | 出现则为 true，隐藏则为 false |
| @remove-tag     | 多选模式下移除 tag 时触发                | 移除的 tag 值                 |
| @clear          | 可清空的单选模式下用户点击清空按钮时触发 | —                             |
| @blur           | 当 input 失去焦点时触发                  | (event: Event)                |
| @focus          | 当 input 获得焦点时触发                  | (event: Event)                |

## 选择插槽

|   name  | 说明                |
| ------- | ------------------- |
| —       | Option 组件列表     |
| prefix  | Select 组件头部内容 |
| empty   | 无选项时的列表      |

## Option 属性

| 参数     | 说明                                      | 类型                 | 可选值 | 默认值 |
| -------- | ----------------------------------------- | -------------------- | ------ | ------ |
| value    | 选项的值                                  | string/number/object | —      | —      |
| label    | 选项的标签，若不设置则默认与 `value` 相同 | string/number        | —      | —      |
| disabled | 是否禁用该选项                            | boolean              | —      | false  |

## 组件方法说明

| 方法名 | 说明                            | 参数 |
| ------ | ------------------------------- | ---- |
| focus  | 使 input 获取焦点               | -    |
| blur   | 使 input 失去焦点，并隐藏下拉框 | -    |

#### 效果图

<f-image alt="" src="https://img-crs.vchangyi.com/组件说明样式资源/管理后台/ts/select.png" width="500px" />

#### 代码参考

```html
<template>
  <div class="counter">
    <div class="template">
      <cy-select
        v-model="eat"
        size="medium"
        placeholder="宽度为280高为36的下拉选择"
        @change="changeEat"
      >
        <cy-option label="冰淇淋" value="1"></cy-option>
        <cy-option label="可乐" value="2"></cy-option>
        <cy-option label="热巧克力" value="3"></cy-option>
      </cy-select>
    </div>
    <div class="template">
      <cy-select
        v-model="eat"
        clearable
        width-size="medium"
        placeholder="宽为400高为36的下拉选择"
        @change="changeEat"
      >
        <cy-option label="冰淇淋" value="1"></cy-option>
        <cy-option label="可乐" value="2"></cy-option>
        <cy-option label="热巧克力" value="3"></cy-option>
      </cy-select>
    </div>
    <div class="template">
      <cy-select
        v-model="eat"
        clearable
        width-size="small"
        size="mini"
        placeholder="宽为280高为32的下拉选择"
        @change="changeEat"
      >
        <cy-option label="冰淇淋" value="1"></cy-option>
        <cy-option label="可乐" value="2"></cy-option>
        <cy-option label="热巧克力" value="3"></cy-option>
      </cy-select>
    </div>
    <div class="template">
      <cy-select
        v-model="eat"
        size="medium"
        placeholder="下拉带搜索功能(280)"
        search-input
        @change="changeEat"
      >
        <cy-option label="冰淇淋" value="1"></cy-option>
        <cy-option label="可乐" value="2"></cy-option>
        <cy-option label="热巧克力" value="3"></cy-option>
      </cy-select>
    </div>
    <div class="template">
      <cy-select
        v-model="eat"
        clearable
        width-size="medium"
        placeholder="下拉带搜索功能(400)"
        search-input
        @change="changeEat"
      >
        <cy-option label="冰淇淋" value="1"></cy-option>
        <cy-option label="可乐" value="2"></cy-option>
        <cy-option label="热巧克力" value="3"></cy-option>
      </cy-select>
    </div>
    <div class="template">
      <cy-select
        v-model="eat"
        clearable
        width-size="medium"
        size="mini"
        disabled
        placeholder="被禁用"
        @change="changeEat"
      >
        <cy-option label="冰淇淋" value="1"></cy-option>
        <cy-option label="可乐" value="2"></cy-option>
        <cy-option label="热巧克力" value="3"></cy-option>
      </cy-select>
    </div>
    <div class="template">
      <cy-select
        v-model="eats"
        @change="changeEat"
        multiple
        collapse-tags
        search-input
        placeholder="多选下拉"
      >
        <cy-option label="冰淇淋" value="1"></cy-option>
        <cy-option label="可乐" value="2"></cy-option>
        <cy-option label="热巧克力" value="3"></cy-option>
      </cy-select>
    </div>
  </div>
</template>

<script lang="ts">
  import { Vue, Component } from "vue-property-decorator"
  @Component
  export default class Counter extends Vue {
    private eat = ""
    private eats = ""
    private changeEat() {
      console.log(this.eat, this.eats)
    }
  }
</script>
<style lang="scss" scoped>
  .template {
    padding: 15px 30px;
  }
</style>
```
