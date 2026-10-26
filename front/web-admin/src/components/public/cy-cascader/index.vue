<template>
  <div
    class="cy-cascader"
    :class="{
      'cy-input-small': widthSize === 'small',
      'cy-input-medium': widthSize === 'medium'
    }"
  >
    <el-cascader
      v-model="myValue"
      ref="cascader"
      :options="options"
      :props="props"
      :size="size"
      :placeholder="placeholder"
      :disabled="disabled"
      :clearable="clearable"
      :show-all-levels="showAllLevels"
      :collapse-tags="collapseTags"
      :separator="separator"
      :filterable="filterable"
      :filter-method="filterMethod"
      :debounce="debounce"
      :before-filter="beforeFilter"
      :popper-class="popperClass"
      @change="onChange"
      @expand-change="onExpandChange"
      @blur="onBlur"
      @visible-change="onVisibleChange"
      @remove-tag="onRemoveTag"
    >
      <slot></slot>
      <template slot="empty">
        <slot name="empty"></slot>
      </template>
    </el-cascader>
  </div>
</template>
<script lang="ts">
const fn = () => {
  return;
};
interface INormalObj {
  [propName: string]: any;
}
interface IPropface {
  expandTrigger?: string;
  multiple?: boolean;
  checkStrictly?: boolean;
  emitPath?: boolean;
  lazy?: boolean;
  lazyLoad?: Function;
  value?: string;
  label?: string;
  children?: string;
  disabled?: string;
  leaf?: string;
}
import { Component, Vue, Prop, Watch } from "vue-property-decorator";
@Component
export default class CyCascader extends Vue {
  @Prop({ default: "small" }) widthSize?: string; //非element-ui参数，由于UI规范定宽 small宽度为280px和medium宽度为400px
  @Prop() value!: any; //选中项绑定值
  @Prop() options?: Array<INormalObj>; //可选项数据源，键名可通过 Props 属性配置
  @Prop() props?: IPropface; //配置选项
  @Prop({ default: "small" }) size?: string; //尺寸
  @Prop({ default: "请选择" }) placeholder?: string; //输入框占位文本
  @Prop({ default: false, type: Boolean }) disabled!: boolean; //是否禁用
  @Prop({ default: false, type: Boolean }) clearable!: boolean; //是否支持清空选项
  @Prop({ default: true, type: Boolean }) showAllLevels!: boolean; //输入框中是否显示选中值的完整路径
  @Prop({ default: false, type: Boolean }) collapseTags!: boolean; //多选模式下是否折叠Tag
  @Prop({ default: "/" }) separator?: string; //尺寸
  @Prop({ default: false, type: Boolean }) filterable!: boolean; //是否可搜索选项
  @Prop({ default: fn }) filterMethod?: Function; //自定义搜索逻辑，第一个参数是节点node，第二个参数是搜索关键词keyword，通过返回布尔值表示是否命中
  @Prop({ default: 300 }) debounce?: number; //搜索关键词输入的去抖延迟，毫秒
  @Prop({ default: fn }) beforeFilter?: Function; //筛选之前的钩子，参数为输入的值，若返回 false 或者返回 Promise 且被 reject，则停止筛选
  @Prop() popperClass?: string; //自定义浮层类名
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
  private onChange(val: object) {
    const timeOut = setTimeout(() => {
      this.$emit("change", val);
      clearTimeout(timeOut);
    });
  }
  private onExpandChange(val: Array<any>) {
    const timeOut = setTimeout(() => {
      this.$emit("expand-change", val);
      clearTimeout(timeOut);
    });
  }
  private onBlur(event: MouseEvent) {
    const timeOut = setTimeout(() => {
      this.$emit("blur", event);
      clearTimeout(timeOut);
    });
  }
  private onFocus(event: MouseEvent) {
    const timeOut = setTimeout(() => {
      this.$emit("focus", event);
      clearTimeout(timeOut);
    });
  }
  private onVisibleChange(val: boolean) {
    const timeOut = setTimeout(() => {
      this.$emit("visible-change", val);
      clearTimeout(timeOut);
    });
  }
  private onRemoveTag(val: any) {
    const timeOut = setTimeout(() => {
      this.$emit("remove-tag", val);
      clearTimeout(timeOut);
    });
  }
  public getCheckedNodes() {
    (this.$refs["cascader"] as any).clearCheckedNodes();
  }
}
</script>
<style lang="scss" scoped>
.cy-input-small {
  width: 280px !important;
}
.cy-input-medium {
  width: 400px !important;
}
.cy-cascader {
  display: inline-block;
  position: relative;
  .el-cascader {
    width: 100% !important;
    .el-input--mini {
      height: 32px !important;
      line-height: 32px !important;
      .el-input__inner {
        .el-input__inner {
          @extend .el-input--mini;
        }
      }
    }
    .el-input--small {
      height: 36px !important;
      line-height: 36px !important;
      .el-input__inner {
        .el-input__inner {
          @extend .el-input--small;
        }
      }
    }
    .el-input--medium {
      height: 36px !important;
      line-height: 36px !important;
      .el-input__inner {
        .el-input__inner {
          @extend .el-input--medium;
        }
      }
    }
  }
}
</style>
