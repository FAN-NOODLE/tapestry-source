<template>
  <el-dialog :visible="visible" width="30%" @close="close">
    <template slot="title">
      <h4>自动批准设置</h4>
    </template>
    <el-form :model="form" ref="form" :rules="rules" label-width="70px">
      <el-form-item label="自动批准" prop="">
        <el-switch
          v-model="form.auto_status"
          active-color="#13ce66"
          inactive-color="#C0C4CC"
        ></el-switch>
        <span v-if="form.auto_status">已开启</span>
        <span v-if="!form.auto_status">已关闭</span>
        <el-time-picker
          :disabled="!form.auto_status"
          prefix-icon="el-icon-date"
          v-model="form.auto_at"
          value-format="HH:mm:ss"
          format="HH:mm:ss"
          placeholder="请选择自动批准时间"
        ></el-time-picker>
        <div class="tips">
          系统将在每天的hh-mm-ss批量对上次自动审批之后到本次自动审批之前时间产生的“待审批”申请单执行“批准”操作，若有需要单独操作的申请单，请在设置的自动时间之前操作。
        </div>
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button size="small" @click="close">取 消</el-button>
      <el-button size="small" type="primary" @click="done">保 存</el-button>
    </div>
  </el-dialog>
</template>

<script lang="ts">
import { Component, Vue, Prop, Emit } from "vue-property-decorator";
import Api from "@/api/wisdom/Withdrawal";

@Component
export default class extends Vue {
  @Prop() visible!: boolean;
  @Prop({ default: 1 }) type!: number;
  private form = {
    auto_status: false,
    auto_at: ""
  };
  private rules = {};

  private mounted() {
    this.getDetail();
  }

  private getDetail() {
    Api.getSetting({ mem_type: this.type }).then(res => {
      this.form.auto_status = res.data.auto_status == 0 ? false : true;
      this.form.auto_at = res.data.auto_at;
    });
  }

  @Emit()
  private close() {
    this.form.auto_status = false;
    this.form.auto_at = "";
  }

  @Emit()
  private done() {
    return {
      auto_status: this.form.auto_status === true ? 1 : 0,
      auto_at: this.form.auto_at,
      mem_type: this.type
    };
  }
}
</script>

<style scoped lang="scss">
/deep/.el-dialog__header {
  h4 {
    margin: 0;
  }
  .select-num {
    font-size: 14px;
    font-weight: normal;
    margin-left: 30px;
  }
}
.tips {
  font-size: 12px;
  color: #999;
  line-height: 16px;
  margin-bottom: 10px;
}
</style>
