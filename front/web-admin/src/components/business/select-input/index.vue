<template>
  <div class="select-input">
    <el-input
      v-model="value"
      size="medium"
      @change="onChangeInput"
      @blur="onBlurInput"
      @focus="onFocusInput"
      @clear="onClearInput"
      :placeholder="placeholder"
      :disabled="disabled"
    >
      <el-select
        class="select"
        slot="prepend"
        v-model="selectVal"
        @change="selectChange"
      >
        <el-option
          v-for="item in list"
          :key="item.key"
          :label="item.name"
          :value="item.key"
        ></el-option>
      </el-select>
    </el-input>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop, Watch } from "vue-property-decorator";
interface IListItem {
  key: string | number;
  name: string;
}
@Component({
  model: {
    prop: "value",
    event: "change"
  }
})
export default class SelectInput extends Vue {
  @Prop() value!: string; //绑定值
  @Prop() list!: Array<IListItem>;
  @Prop() placeholder?: string;
  @Prop({ default: false }) disabled!: boolean;

  private selectVal: number|string = null;

  @Watch("list", {deep: true})
  onChangeList(val: any) {
    this.selectVal = val[0].key
  }

  mounted() {
    if(this.list[0]) {
      this.selectVal = this.list[0].key
    }
  }

  private onChangeInput(val: any) {
    this.$emit("change", val);
  }
  // 在 Input 失去焦点时触发
  private onBlurInput() {
    this.$emit("blur");
  }
  // 在 Input 获得焦点时触发
  private onFocusInput() {
    this.$emit("focus");
  }
  // 在点击由 clearable 属性生成的清空按钮时触发
  private onClearInput() {
    this.$emit("clear");
  }
  private selectChange() {
    this.value = ''
    this.$emit("change", '');
    this.$emit("selectChange", this.selectVal);
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
  public reset() {
    this.selectVal = this.list[0].key;
  }
}
</script>
<style lang="scss" scoped>
.select-input {
  width: 280px;
  .select {
    width: 110px !important;
  }
  /deep/ {
    .el-input-group__prepend {
      background: #ffffff;
    }
  }
}
</style>
