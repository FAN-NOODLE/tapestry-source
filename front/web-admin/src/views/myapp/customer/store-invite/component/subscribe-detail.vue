<template>
  <div class="contaniner">
    <el-dialog
      class="dialog-box"
      width="700px"
      min-height="900px"
      title="预约详情"
      :visible.sync="dialogVisible"
      :close-on-click-modal="false"
      @close="closeDialog"
    >
      <div class="dialog-body">
        <el-form ref="form" label-width="100px">
          <el-form-item label="顾客姓名：">顾客姓名8888</el-form-item>
          <el-form-item label="邀约时间：">{{
            1234567890 | formatDate("yyyy-mm-dd hh:mm:ss")
          }}</el-form-item>
          <el-form-item label="到店时间："
            >2022-04-20 10:00 - 11:00</el-form-item
          >
          <el-form-item label="服务名称：">皮革管理</el-form-item>
          <el-form-item label="邀约导购：">邀约导购</el-form-item>
          <el-form-item label="邀约门店：">邀约导购</el-form-item>
          <el-form-item label="门店地址："
            >门店地址门店地址门店地址</el-form-item
          >
          <el-form-item label="导购备注："
            >导购备注导购备注导购备注</el-form-item
          >
          <el-form-item label="客户留言：">----</el-form-item>
        </el-form>
      </div>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop, Emit } from "vue-property-decorator";

@Component
export default class extends Vue {
  @Prop({ default: "选择归属门店" }) title?: string; // 右侧选择标题

  /*** data ***/
  private dialogVisible = true; // 控制模态框显示/隐藏
  // 关闭模态框的回调
  @Emit("close")
  private closeDialog() {
    this.$emit("close");
  }
}
</script>

<style lang="scss" scoped>
/deep/.el-dialog {
  max-height: 700px !important;
  .el-form-item {
    margin-bottom: 0;
  }
  .el-form-item__label {
    color: gray;
  }
}
</style>
