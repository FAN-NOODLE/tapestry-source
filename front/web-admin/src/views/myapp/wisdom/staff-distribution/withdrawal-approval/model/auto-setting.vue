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
        <div class="tips">
          首次设置后，自动批准将在下一个设置时刻开始执行，仅处理设置时间到执行时间期间产生的待审批申请单；之后将在每天同一时刻自动执行上次自动批准之后产生的“待审批”申请单。
        </div>
        <el-time-picker
          :disabled="!form.auto_status"
          prefix-icon="el-icon-date"
          v-model="form.auto_at"
          value-format="HH:mm:ss"
          format="HH:mm:ss"
          placeholder="请选择自动批准时间"
        ></el-time-picker>
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button size="small" @click="close">取 消</el-button>
      <el-button size="small" type="primary" @click="done">确 认</el-button>
    </div>
  </el-dialog>
</template>

<script lang="ts">
import { Component, Vue, Prop, Emit } from "vue-property-decorator";
import Api from "@/api/wisdom/Withdrawal";

@Component
export default class extends Vue {
  @Prop() visible!: boolean;

  private form = {
    auto_status: false,
    auto_at: ""
  };
  private rules = {};

  private mounted() {
    this.getDetail();
  }

  private getDetail() {
    Api.getSetting().then(res => {
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
      auto_at: this.form.auto_at
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
