<template>
  <div>
    <el-form
      ref="myform"
      :model="model"
      :rules="selfRules"
      :inline="inline"
      :labelPosition="labelPosition"
      :labelWidth="labelWidth"
      :statusIcon="statusIcon"
      :validateOnRuleChange="validateOnRuleChange"
      :labelSuffix="labelSuffix"
      :disabled="disabled"
      :size="size"
    >
      <slot></slot>
    </el-form>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import { cloneDeep } from "@/utils/tools";
import Bus from "@/components/bus";
import { regrexFun } from "./regexp";
interface IRulesInterface {
  [propName: string]: any;
}
@Component({
  name: "cy-form"
})
export default class CyForm extends Vue {
  @Prop() model!: object;
  @Prop() rules!: object;
  @Prop({ default: false }) inline?: boolean;
  @Prop({ default: "148px" }) labelWidth?: string;
  @Prop({ default: "left" }) labelPosition?: string;
  @Prop({ default: false, type: Boolean }) statusIcon?: boolean;
  @Prop({ default: true, type: Boolean }) validateOnRuleChange?: boolean;
  @Prop({ default: "medium" }) size?: string;
  @Prop() labelSuffix?: string;
  @Prop({ default: false, type: Boolean }) disabled?: boolean;
  private selfRules: IRulesInterface = {};
  private originRules: IRulesInterface = {};
  // 暴露方法-> 校验
  public validate(callback: Function) {
    return (this.$refs.myform as any).validate(callback);
  }
  // 暴露方法-> 校验单个
  public validateField(checkname: string | Array<string>, callback: Function) {
    return (this.$refs.myform as any).validateField(checkname, callback);
  }
  // 暴露方法-> 重置
  public resetFields() {
    return (this.$refs.myform as any).resetFields();
  }
  created() {
    this.selfRules = this.rules || {};
    this.originRules = cloneDeep(this.rules);
    this.onChangeRules();
  }
  // private beforeMount() {
  //   }
  private mounted() {
    Bus.$emit("setWidth", this.labelWidth);
  }
  onChangeRules() {
    const _this = this;
    Bus.$on("changeRules", (innerRules: any) => {
      // 初始化
      _this.selfRules[innerRules.ruleName] ||
        (_this.selfRules[innerRules.ruleName] = []);
      // 生成规则
      _this.selfRules[innerRules.ruleName] = regrexFun(
        _this.originRules[innerRules.ruleName] || [],
        innerRules
      );
    });
  }
}
</script>

<style lang="scss" scoped></style>
