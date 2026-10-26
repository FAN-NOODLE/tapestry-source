<template>
  <el-option :value="value" :label="label" :disabled="disabled"></el-option>
</template>
<script lang="ts">
interface INormalObj {
  [propName: string]: any;
}
import { Component, Vue, Prop } from "vue-property-decorator";

@Component({})
export default class CyOption extends Vue {
  @Prop() value?: string | number | INormalObj; // 选项的值
  @Prop() label?: string | number; //选项的标签，若不设置则默认与 value 相同
  @Prop({ default: false, type: Boolean }) disabled!: boolean; //是否禁用该选项
}
</script>
