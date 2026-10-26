<template>
  <el-dialog
    title="离开页面提示"
    :visible.sync="isCancelDialog"
    width="420px"
    @close="editGoCloseFn"
  >
    <span
      ><i aria-hidden="true" class="el-icon-warning"></i>
      你确认放弃编辑吗？放弃后编辑的信息将不会保存</span
    >
    <span slot="footer" class="dialog-footer">
      <el-button @click="goBackFn">取消</el-button>
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
  //取消
  private goBackFn() {
    this.$emit("goEditBack", "close");
  }
  //确认
  private saveGoBackFn() {
    this.$emit("goEditBack", "cancel");
  }
  //关闭
  private editGoCloseFn() {
    this.$emit("goEditBack", "close");
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
