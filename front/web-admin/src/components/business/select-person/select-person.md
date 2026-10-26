# 选择员工组件


## 组件参数说明

|           参数         | 类型          | 说明          | 可选值                           | 默认值  |  是否必填   |
| ---------------------- | ------------ | ------------- | -------------------------------- | ------- | --------- |
|    visible             | boolean      | 显示隐藏组件   |  true/ false                     |  无     | 是         |
| selectedData           | Array<guideNameListType> （参考type 中类型）| 选中列表       | 无                               | []      | 否         |
| appendToBody           | boolean      | 是否添加到body 元素中, 解决el-dialog 弹窗弹出弹窗导致的遮罩问题 | true/ false | false | 否|
| title                  | string       | 自定义标题     | 无                               | 添加员工（支持多选）| 否|


## 组件事件说明

| 事件名称 | 说明 | 回调参数 |
|---------|---------|---------|
| @submit | 点击确定之后返回值（选中员工数据） | callBackType （参考type 中类型） |


#### 代码参考

```html
<template>
  <div class="example">
      <!-- 选择员工 -->
      <select-person
        :selectedData="selectPerson.selectedData"
        :visible.sync="selectPerson.visible"
        @submit="selectPerson.submit"
        :appendToBody="true"
      ></select-person>
  </div>
</template>

<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
import SelectInput from "@/components/business/select-person/index.vue";
import { callBackType } from "@/components/business/select-person/type"

@Component({
  components: {
    SelectPerson
  }
})
export default class Example extends Vue {
  
  // 选择员工
  private selectPerson = {
    visible: false, // 判断显示隐藏
    selectedData: [], // 选中ids
    submit: this.handleSubmit
  };

  /**
   * 选择员工回调
   */
  private handleSubmit(data: callBackType) {
    console.log(data);
  }


}
</script>

```

#### 注意， 为了保存原始数据（选中的部门信息，以及选中的人员信息）， 需要在编辑/新增页保存成功后 调用保存原始数据接口， 选择修改时，获取原始数据接口并且赋值给selectedData， 可以参考@store/persons.ts 