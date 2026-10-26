<template>
  <el-dialog
    title="离开页面提示"
    :visible.sync="isCancelDialog"
    width="420px"
    @close="CloseFn"
    append-to-body
  >
    <span
      ><i aria-hidden="true" class="el-icon-warning"></i>
      你确认要放弃编辑吗？放弃后编辑的信息将不会被保存</span
    >
    <span slot="footer" class="dialog-footer">
      <el-button @click="CloseFn">取消</el-button>
      <el-button type="primary" @click="saveGoBackFn">确认</el-button>
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
  //确认
  private saveGoBackFn() {
    this.$emit("editGoBack", "cancel");
  }
  //关闭
  private CloseFn() {
    this.$emit("editGoBack", "close");
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
