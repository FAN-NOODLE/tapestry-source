<template>
  <el-dialog
    title="离开页面提示"
    :visible.sync="isCancelDialog"
    width="420px"
    @close="goCloseFn"
  >
    <span
      ><i aria-hidden="true" class="el-icon-warning"></i>
      检测到未保存的内容，是否保存到下次新增继续编辑？</span
    >
    <span slot="footer" class="dialog-footer">
      <el-button @click="goBackFn">直接离开</el-button>
      <el-button type="primary" @click="saveGoBackFn">保存并离开</el-button>
    </span>
  </el-dialog>
</template>
<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
@Component({})
export default class extends Vue {
  /************data**************/
  private isCancelDialog = true;
  /************methods**************/
  //直接离开
  private goBackFn() {
    this.$emit("goBack", "cancel");
  }
  //保存并离开
  private saveGoBackFn() {
    this.$emit("goBack", "save");
  }
  //关闭
  private goCloseFn() {
    this.$emit("goBack", "close");
  }
}
</script>
<style lang="scss" scoped>
.el-range-editor {
  width: 280px !important;
}
.in-item {
  width: 100%;
}
</style>
