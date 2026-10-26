<template>
  <el-dialog
    :title="pageTitle"
    :visible.sync="dialogVisible"
    width="560px"
    @close="closeDialog"
  >
    <cy-form
      class="main-box"
      ref="cform"
      :rules="rules"
      :model="ruleForm"
      label-width="90px"
    >
      <!--手机号-->
      <cy-form-item
        label="手机号"
        prop="param"
        required
        v-if="curPopObj.type == 'member_mobile'"
        phone
      >
        <cy-input
          v-model="ruleForm.param"
          maxlength="11"
          placeholder="请输入导购手机号"
          show-word-limit
          @keyup.native="checkSpace()"
          @input="_checkInt(ruleForm, 'param', ruleForm.param)"
        ></cy-input>
      </cy-form-item>
      <!--微信昵称-->
      <cy-form-item
        label="微信昵称"
        prop="param"
        required
        v-if="curPopObj.type == 'member_name'"
      >
        <cy-input
          v-model="ruleForm.param"
          clearable
          maxlength="10"
          placeholder="请输入导购微信昵称"
          show-word-limit
          @keyup.native="checkSpace()"
        ></cy-input>
      </cy-form-item>
      <!--日期-->
      <cy-form-item
        label="类型"
        prop="selectData"
        required
        v-if="curPopObj.type == 'data'"
      >
        <el-radio-group v-model="ruleForm.selectData" @change="changeDateFn">
          <el-radio label="1">指定日期</el-radio>
          <el-radio label="2">日期时间段</el-radio>
        </el-radio-group>
      </cy-form-item>
      <cy-form-item
        prop="param"
        label="日期"
        required
        v-if="curPopObj.type == 'data' && ruleForm.selectData == '1'"
        trigger="change"
      >
        <el-date-picker
          v-model="ruleForm.param"
          type="date"
          placeholder="请选择日期"
          value-format="yyyy-MM-dd"
        >
        </el-date-picker>
      </cy-form-item>
      <cy-form-item
        prop="param"
        label="日期"
        required
        v-if="curPopObj.type == 'data' && ruleForm.selectData == '2'"
        trigger="change"
      >
        <el-date-picker
          type="daterange"
          v-model="ruleForm.param"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
          value-format="yyyy-MM-dd"
          range-separator="至"
          @blur="blurDateFn"
        ></el-date-picker>
      </cy-form-item>

      <!--商品名称-->
      <cy-form-item
        label="商品名称"
        prop="param"
        required
        v-if="curPopObj.type == 'goods_name'"
      >
        <cy-input
          v-model="ruleForm.param"
          maxlength="16"
          placeholder="请输入商品名称"
          show-word-limit
          @keyup.native="checkSpace()"
        ></cy-input>
      </cy-form-item>
      <!--设置数字-->
      <cy-form-item
        label="数字"
        prop="param"
        required
        number
        v-if="curPopObj.type == 'number'"
      >
        <cy-input
          v-model="ruleForm.param"
          clearable
          maxlength="10"
          placeholder="请输入数字"
          show-word-limit
          @keyup.native="checkSpace()"
          @input="_checkInt(ruleForm, 'param', ruleForm.param)"
        ></cy-input>
      </cy-form-item>
      <!--设置量词-->
      <cy-form-item
        label="量词"
        prop="param"
        required
        v-if="curPopObj.type == 'classifier'"
      >
        <cy-select width-size="small" size="small" v-model="ruleForm.param">
          <cy-option
            v-for="(item, index) in unitArr"
            :key="index"
            :label="item"
            :value="item"
            placeholder="请选择量词"
          ></cy-option>
        </cy-select>
      </cy-form-item>
    </cy-form>
    <span slot="footer" class="dialog-footer">
      <el-button size="small" @click="closeDialog">取 消</el-button>
      <el-button size="small" type="primary" @click="submitFn">保 存</el-button>
    </span>
  </el-dialog>
</template>
<script lang="ts">
import { Component, Prop, Mixins } from "vue-property-decorator";
import CheckNum from "@/mixins/NumberInput";
@Component({})
export default class extends Mixins(CheckNum) {
  @Prop() curPopObj?: any;
  //classifier 量词
  //number 数字
  //goods_name 商品名称
  //data 日期
  //member_name 员工微信昵称
  //member_mobile 员工手机号
  private unitArr = [
    "片",
    "件",
    "个",
    "只",
    "张",
    "条",
    "颗",
    "盒",
    "台",
    "批",
    "根",
    "头",
    "匹",
    "辆",
    "包"
  ]; //量词
  // private selectData = "1"; //选择日期
  private pageTitle = ""; //弹框title
  private checkGroupName = (rule: any, value: any, callback: any) => {
    callback();
  };
  private checkName = (rule: any, value: any, callback: any) => {
    if (!value) {
      return callback(new Error("请输入标签名称"));
    }
    callback();
  };
  /************data**************/
  private dialogVisible = true;
  private rules = {};
  private ruleForm = {
    param: "",
    selectData: "1"
  };
  // [("赋值员工手机号", "赋值微信昵称", "赋值日期", "赋值商品名称","设置数字","设置量词")];
  /************methods**************/
  // 禁止输入空格
  private checkSpace() {
    this.ruleForm.param = this.ruleForm.param.replace(/^\s+|\s+$/gm, "");
  }
  //保存
  private submitFn() {
    if (!this.ruleForm.param) {
      this.$message.warning("请完善信息");
    }
    (this.$refs.cform as any).validate((valid: boolean) => {
      if (valid) {
        this.$emit("closeDialog", {
          type: this.curPopObj.name,
          val: this.ruleForm.param,
          num: this.curPopObj.type,
          selectData: this.ruleForm.selectData
        });
      }
    });
  }
  private changeDateFn() {
    this.ruleForm.param = "";
  }
  private blurDateFn(val: any) {
    this.ruleForm.param = val.value[0] + "至" + val.value[1];
  }
  // 关闭弹框
  private closeDialog() {
    this.$emit("closeDialog");
  }
  private created() {
    this.pageTitle = "设置" + this.curPopObj.title;
    const title = this.curPopObj.name.split("%").join("");
    if (title != this.curPopObj.title) {
      this.ruleForm.param = this.curPopObj.name.split("%").join("");
      if (this.curPopObj.selectData) {
        this.ruleForm.selectData = this.curPopObj.selectData;
      }
    } else {
      this.ruleForm.param = "";
    }
  }
}
</script>
<style lang="scss" scoped>
/deep/.el-range-editor {
  width: 280px !important;
}
</style>
