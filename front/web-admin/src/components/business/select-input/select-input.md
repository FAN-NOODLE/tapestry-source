# 附带下拉选择的输入框组件

> 说明：满足UI需求，二次封装的附带下拉选择的输入框组件

## 组件参数说明

|           参数         | 类型          | 说明          | 可选值                           | 默认值  |  是否必填   |
|------------------- |-------------- |---------------|--------------------------------  |-------- | ------- |
| v-model / value       | String  | 绑定值 |     -    | - |    是    |
| list     |            Array(obj)> obj: {key: 绑定对应value值,name: 绑定的label值}           | 下拉框的option数组 | - | - |   是  |
| placeholder       | String  | 占位符文案 |     -    | - |    否    |
| disabled       | Boolean  | 是否禁用 |     -    |   false  |    否    |

## 组件事件说明

| 事件名称 | 说明 | 回调参数 |
|---------|---------|---------|
| @change | 仅在输入框失去焦点或用户按下回车时触发 | 输入框目前的内容 |
| @selectChange |  下拉框选中内容变化回调  |  目前选中值  | 
| blur | 在 Input 失去焦点时触发 | (event: Event) |
| focus | 在 Input 获得焦点时触发 | (event: Event) |
| input | 在 Input 值改变时触发 | (value: string \| number) |

## 组件方法说明
| 方法名称 | 说明 | 参数 |
| ---- | ---- | ---- |
| focus | 使 input 获取焦点 | — |
| blur | 使 input 失去焦点 | — |
| select | 选中 input 中的文字 | — |

#### 效果图
<f-image src="https://img-crs.vchangyi.com/组件说明样式资源/管理后台/ts/select-input.png" width="500px" />

#### 代码参考

```html
  <template>
  <div class="counter">
    <div class="template">
      <!-- 需要单独引入 SelectInput -->
      <select-input
        v-model="text"
        :list="list"
        @selectChange="changed"
        @change="changeText"
        placeholder="请输入啊"
      ></select-input>
    </div>
  </div>
</template>

<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
import SelectInput from "@/components/business/select-input/index.vue";
@Component({
  name: "Counter",
  components: {
    SelectInput
  }
})
export default class Counter extends Vue {
  private text = "";
  private list = [
    { key: 1, name: "老坛胡说" },
    { key: 2, name: "徐大虾咯" }
  ];
  private changed(val: any) {
    console.log(val);
  }
  private changeText() {
    console.log(this.text);
  }
}
</script>
<style lang="scss" scoped>
.template {
  padding: 15px 30px;
}
</style>


```