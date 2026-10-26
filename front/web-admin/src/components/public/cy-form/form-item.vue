<template>
  <div class="in-item">
    <el-form-item
      ref="myformitem"
      :prop="prop"
      :labelPosition="labelPosition"
      :showMessage="showMessage"
      :inlineMessage="inlineMessage"
      :error="error"
    >
      <div
        class="label"
        :class="required ? 'require-label' : ''"
        :style="'max-width:' + (labelWidth - 20)"
        slot="label"
      >
        <span>{{ label }}</span>
      </div>
      <slot></slot>
      <template slot="error" slot-scope="error">
        <div class="error">{{ error.error }}</div>
      </template>
    </el-form-item>
    <div
      class="desc m-t-5 color-text-second"
      :style="'margin-left:' + labelWidth"
      v-if="desc"
    >
      {{ desc }}
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop, Watch } from "vue-property-decorator";
import Bus from "@/components/bus";
@Component({
  name: "cy-form-item",
})
export default class CyFormItem extends Vue {
  @Prop() prop?: string;
  @Prop() label?: string;
  @Prop({ default: "right" }) labelPosition!: string;
  @Prop({ default: true }) showMessage?: boolean;
  @Prop({ default: false, type: Boolean }) inlineMessage?: boolean;
  @Prop() desc?: string;
  @Prop() error?: string;
  @Prop({ default: "blur" }) trigger?: string;
  // 内置校验规则
  @Prop({ default: false, type: Boolean }) required?: boolean;
  @Prop({ default: false, type: Boolean }) email?: boolean;
  @Prop({ default: false, type: Boolean }) phone?: boolean;
  @Prop({ default: false, type: Boolean }) idcard?: boolean;
  @Prop({ default: false, type: Boolean }) number?: boolean;
  @Prop({ default: false, type: Boolean }) url?: boolean;
  @Prop({ default: false, type: Boolean }) zh?: boolean;
  @Prop({ default: false, type: Boolean }) en?: boolean;
  @Prop({ default: false, type: Boolean }) ennum?: boolean;
  @Prop({ default: false, type: Boolean }) en_num?: boolean;
  @Prop({ default: false, type: Boolean }) trims?: boolean;
  @Prop() min?: number | string;
  @Prop() max?: number | string;

  private labelWidth = "128px";
  public validate() {
    (this.$refs["myformitem"] as any).validate((valid: any) => {
      console.log(valid);
    });
  }
  beforeCreate() {
    Bus.$on("setWidth", (width: string) => {
      this.labelWidth = width;
    });
  }
  private created() {
    // 若设置内置校验则派发事件到外层form
    this.changeRules();
  }
  // private beforeMount() {}
  private changeRules() {
    const rulesArr = [
      "required",
      "min",
      "max",
      "email",
      "phone",
      "idcard",
      "number",
      "url",
      "zh",
      "en",
      "ennum",
      "en_num",
      "trims",
    ];
    const params = {};
    rulesArr.forEach((item) => {
      if (this.prop) {
        params[item] = this[item];
      }
    });
    Bus.$emit("changeRules", {
      ruleName: this.prop,
      params,
      trigger: this.trigger,
      label: this.label,
    });
  }

  @Watch("required")
  changeRequired() {
    this.changeRules();
  }
  @Watch("email")
  changeEmail() {
    this.changeRules();
  }
  @Watch("phone")
  changePhone() {
    this.changeRules();
  }
  @Watch("idcard")
  changeIdcard() {
    this.changeRules();
  }
  @Watch("number")
  changeNumber() {
    this.changeRules();
  }
  @Watch("url")
  changeUrl() {
    this.changeRules();
  }
  @Watch("min")
  changeMin() {
    this.changeRules();
  }
  @Watch("max")
  changeMax() {
    this.changeRules();
  }
}
</script>

<style lang="scss" scoped>
/deep/ {
  .el-form-item__label {
    padding-right: 20px;
    &:before {
      content: "" !important;
      margin-right: 0 !important;
    }
  }
  .el-form-item {
    margin-bottom: 0;
  }
  .el-form-item__label {
    line-height: 20px;
    position: relative;
    .label {
      text-align: right;
    }
  }
}
.in-item {
  margin-bottom: 30px;
  width: 528px;
}
.label {
  max-height: 40px;
  line-height: 20px !important;
  min-width: 20px;
  max-width: 108px;
  padding-left: 10px;
  position: absolute;
  right: 20px;
  top: 20px;
  transform: translateY(-50%);
  &.require-label {
    &:before {
      content: "*";
      color: $color-danger;
      position: absolute;
      left: 2px;
      top: 50%;
      transform: translateY(-45%);
    }
  }
}
.error {
  padding-top: 0;
  margin-top: 5px;
  line-height: 18px;
  max-height: 36px;
  color: $color-danger;
  font-size: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
}
.desc {
  margin-left: 128px;
  font-size: 12px;
  line-height: 18px;
}
</style>
