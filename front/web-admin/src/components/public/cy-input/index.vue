<template>
  <div
    class="cy-input"
    :class="{
      'cy-input-mini': widthSize === 'mini',
      'cy-input-small': widthSize === 'small',
      'cy-input-medium': widthSize === 'medium',
      'cy-text-area': type === 'textarea'
    }"
  >
    <el-input
      v-model="myValue"
      ref="input"
      :type="type"
      :maxlength="maxlength"
      :minlength="minlength"
      :show-word-limit="showWordLimit"
      :placeholder="placeholder"
      :clearable="clearable"
      :show-password="showPassword"
      :disabled="disabled"
      :size="size"
      :prefix-icon="prefixIcon"
      :suffix-icon="suffixIcon"
      :rows="rows"
      :autosize="autosize"
      :autocomplete="autocomplete"
      :name="name"
      :readonly="readonly"
      :autofocus="autofocus"
      @change="onChangeInput"
      @blur="onBlurInput"
      @focus="onFocusInput"
      @clear="onClearInput"
    >
      <template slot="prepend">
        <slot name="prepend"></slot>
      </template>
      <template slot="append">
        <slot name="append"></slot>
      </template>
    </el-input>
  </div>
</template>

<script lang="ts">
interface IAutosize {
  minRows?: number;
  maxRows?: number;
}
import { Component, Vue, Prop, Watch } from "vue-property-decorator";
@Component({})
export default class CyInput extends Vue {
  @Prop({ default: "small" }) widthSize?: string; //非element-ui参数，由于UI定宽，仅有mini宽度为140px small宽度为280px和medium宽度为400px
  @Prop({ default: "text" }) type?: string; //类型 text，textarea 和其他 原生 input 的 type 值
  @Prop() value!: string | number; //绑定值
  @Prop() maxlength?: number; //最大输入长度
  @Prop() minlength?: number; //最小输入长度
  @Prop({ default: false, type: Boolean }) showWordLimit?: boolean; //是否显示输入字数统计
  @Prop() placeholder?: string;
  @Prop({ default: false, type: Boolean }) clearable?: boolean; //是否可清空
  @Prop({ default: false, type: Boolean }) showPassword?: boolean; //是否显示切换密码图标
  @Prop({ default: false, type: Boolean }) disabled?: boolean; //是否禁用
  @Prop({ default: "small" }) size!: string; //高度size
  @Prop() prefixIcon?: string; //输入框头部图标
  @Prop() suffixIcon?: string; //输入框尾部图标
  @Prop({ default: 2 }) rows?: number; //输入框行数，只对 type="textarea" 有效
  @Prop({ default: false }) autosize?: boolean | IAutosize; //自适应内容高度，只对 type="textarea" 有效，可传入对象，如，{ minRows: 2, maxRows: 6 }
  @Prop({ default: "off" }) autocomplete?: string; //原生属性，自动补全
  @Prop() name?: string; //原生属性
  @Prop({ default: false, type: Boolean }) readonly!: boolean; //原生属性，是否只读
  @Prop({ default: false, type: Boolean }) autofocus!: boolean; //原生属性，自动获取焦点
  @Watch("value")
  onChangeValue(val: any) {
    this.myValue = val;
  }
  @Watch("myValue")
  onChangeMyValue(val: any) {
    this.$emit("input", val);
  }
  /**
   * @name data数据
   * */
  private myValue = this.value;
  /**
   * @name 方法
   * */
  // 仅在输入框失去焦点或用户按下回车时触发
  private onChangeInput(val: any) {
    this.$emit("change", val);
  }
  // 在 Input 失去焦点时触发
  private onBlurInput(e: MouseEvent) {
    this.$emit("blur", e);
  }
  // 在 Input 获得焦点时触发
  private onFocusInput(e: MouseEvent) {
    this.$emit("focus", e);
  }
  // 在点击由 clearable 属性生成的清空按钮时触发
  private onClearInput() {
    this.$emit("clear");
  }
  // element的获取焦点方法
  public focus() {
    (this.$refs["input"] as any).focus();
  }
  // element的失去焦点方法
  public blur() {
    (this.$refs["input"] as any).blur();
  }
  // element提供选中 input 中的文字
  public select() {
    (this.$refs["input"] as any).select();
  }
}
</script>
<style lang="scss" scoped>
.cy-input-mini {
  width: 140px !important;
}
.cy-input-small {
  width: 280px !important;
}
.cy-input-medium {
  width: 400px !important;
}
.cy-input {
  display: inline-block;
  position: relative;
  // 将element样式 修改为 UI 规范样式
  /deep/ {
    .el-input__icon {
      font-weight: 500 !important;
    }
    .el-input__count {
      color: $color-placeholder !important;
    }
    .el-input-group__prepend {
      padding: 0 15px !important;
    }
    .el-input-group__append {
      padding: 0 15px !important;
    }
    .el-input--mini {
      height: 32px !important;
      line-height: 32px !important;
      .el-input__inner {
        @extend .el-input--mini;
      }
      .el-input-group__prepend,
      .el-input-group__append {
        line-height: 30px !important;
        height: 30px !important;
      }
    }
    .el-input--small {
      height: 36px !important;
      line-height: 36px !important;
      .el-input__inner {
        @extend .el-input--small;
      }
      .el-input-group__prepend,
      .el-input-group__append {
        line-height: 34px !important;
        height: 34px !important;
      }
    }
    .el-input--medium {
      height: 40px !important;
      line-height: 40px !important;
      .el-input__inner {
        @extend .el-input--medium;
      }
      .el-input-group__prepend,
      .el-input-group__append {
        line-height: 38px !important;
        height: 38px !important;
      }
    }
    .el-textarea__inner {
      height: 125px;
    }
  }
}
.cy-text-area {
  .el-input--small,
  .el-input--mini,
  .el-input--medium {
    height: auto !important;
  }
}
</style>
