<template>
  <el-dialog
    title="离开页面提示"
    :visible.sync="isCancelDialog"
    @close="goCloseFn"
    width="420px"
    append-to-body
  >
    <span
      ><i aria-hidden="true" class="el-icon-warning"></i>
      检测到未保存的内容，是否保存至下次新增继续编辑？</span
    >
    <span slot="footer" class="dialog-footer">
      <el-button @click="goBackFn">直接离开</el-button>
      <el-button type="primary" @click="saveGoBackFn">保存并离开</el-button>
    </span>
  </el-dialog>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
@Component({})
export default class extends Vue {
  @Prop({ default: () => "1" }) isAppendBody?: string;

  /************data**************/
  private isCancelDialog = true;
  /************methods**************/
  //直接离开
  private goBackFn() {
    this.$emit("createGoBack", "cancel");
  }
  //保存并离开
  private saveGoBackFn() {
    this.$emit("createGoBack", "save");
  }
  //关闭
  private goCloseFn() {
    this.$emit("createGoBack", "close");
  }
}
</script>
<style lang="scss" scoped>
/deep/ {
  .el-dialog {
    border-radius: 10px;
  }
  .el-dialog__header {
    padding: 30px 30px 0;
  }
  .el-dialog__body {
    padding: 20px 30px;
  }
  .el-dialog__footer {
    padding: 0 30px 30px !important;
  }
  .el-dialog__headerbtn .el-dialog__close {
    font-size: 20px;
  }
}
.el-range-editor {
  width: 280px !important;
}
.in-item {
  width: 100%;
}
</style>
