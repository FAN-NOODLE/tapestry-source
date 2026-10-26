<template>
  <el-dialog
    title="驳回确认"
    width="500px"
    label-position="top"
    :visible.sync="dialog_visible"
    @close="closeDialog"
  >
    <div class="content">
      <el-form label-position="left" label-width="250">
        <el-form-item label="你确认要审批驳回吗？请将驳回原因告知分销合伙人吧">
          <el-input
            placeholder="请输入驳回原因（必填）"
            v-model="reject_reason"
            type="textarea"
            size="small"
            maxlength="200"
            :autosize="{ minRows: 6 }"
            show-word-limit="show-word-limit"
          ></el-input>
        </el-form-item>
      </el-form>
      <div class="button">
        <el-button size="small" @click="closeDialog">取消</el-button>
        <el-button size="small" type="primary" @click="save">确定</el-button>
      </div>
    </div>
  </el-dialog>
</template>

<script lang="ts">
import { Component, Vue, Prop, Emit, Watch } from "vue-property-decorator";
import Api from "@/api/wisdom/Withdrawal";

@Component
export default class extends Vue {
  @Prop({ default: false }) visible!: boolean; // 弹框显示/隐藏
  @Prop({ default: () => [] }) case_id?: Array<any>;

  private dialog_visible = false;
  private reject_reason = "";

  @Watch("visible")
  onChangeVisible() {
    this.dialog_visible = true;
  }

  // 关闭弹框
  @Emit("close")
  private closeDialog() {
    this.dialog_visible = false;
    this.reject_reason = "";
  }

  // 点击保存
  // @Emit("success")
  private save() {
    if (!this.reject_reason) {
      this.$message.warning("请输入驳回原因");
      return;
    }
    const params = {
      action: "refuse",
      ids: this.case_id,
      reason: this.reject_reason
    };
    Api.operationStatus(params).then(() => {
      this.$message.success("驳回成功!");
      this.$emit("success");
      this.closeDialog();
    });
  }
}
</script>

<style lang="scss" scoped>
/deep/.el-dialog__body {
  padding-top: 0px;
  padding-bottom: 20px;
}
/deep/.el-dialog__header {
  padding-bottom: 0px;
}
.button {
  text-align: right;
}
</style>
