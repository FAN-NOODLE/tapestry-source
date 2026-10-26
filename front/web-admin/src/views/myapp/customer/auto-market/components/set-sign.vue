<template>
  <el-dialog
    title="新建短信签名"
    :visible.sync="dialogVisible"
    width="520px"
    @close="closeDialog"
  >
    <cy-form
      class="main-box"
      ref="cform"
      :model="ruleForm"
      :rules="rules"
      label-width="58px"
    >
      <!--签名-->
      <cy-form-item
        label="签名"
        prop="signName"
        @keyup.native="trimFn"
        required
        desc="建议内容为：店铺名称、公众号名称或代表发送者身份的签名，一些特殊性词汇将导致发送短信失败！内容限制长度3-8个字！"
      >
        <cy-input
          v-model="ruleForm.signName"
          maxlength="8"
          width-size="medium"
          placeholder="请输入至少3位的签名"
          show-word-limit
        ></cy-input>
      </cy-form-item>
    </cy-form>
    <div class="dialog-footer">
      <el-button size="small" @click="closeDialog">取 消</el-button>
      <el-button size="small" type="primary" @click="submitFn">提 交</el-button>
    </div>
  </el-dialog>
</template>
<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
import API from "@/api/customer";
@Component({})
export default class extends Vue {
  /************data**************/
  private dialogVisible = true;
  private rules = {
    signName: [
      {
        validator: this.validateself,
        trigger: "blur"
      }
    ]
  };
  private validateself(rule: any, value: any, callback: any) {
    value.length < 3 ? callback(new Error("请输入至少3位的签名")) : callback();
  }
  private ruleForm = {
    signName: ""
  };
  /************methods**************/
  //input 不输入空格
  private trimFn() {
    this.ruleForm.signName = this.ruleForm.signName.replace(/^\s+|\s+$/gm, "");
  }
  //保存
  private submitFn() {
    if (!this.ruleForm.signName) {
      this.$message.warning("请完善信息");
    }
    (this.$refs["cform"] as any).validate((valid: boolean) => {
      if (valid) {
        API.smsCreateSignApi({ sign_name: this.ruleForm.signName }).then(
          (res: any) => {
            if (res.code == 200) {
              this.$emit("closeSign", "create");
            } else {
              this.$message.error(res.message);
            }
          }
        );
      }
    });
  }
  // 关闭弹框
  private closeDialog() {
    this.$emit("closeSign");
  }
}
</script>
<style lang="scss" scoped>
.el-range-editor {
  width: 280px !important;
}
.in-item {
  width: 100%;
  margin-bottom: 0;
}
.dialog-footer {
  text-align: right;
  margin-top: 17px;
}
</style>
