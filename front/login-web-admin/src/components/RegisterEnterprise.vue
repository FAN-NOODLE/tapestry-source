<template lang="pug">
  .register-box
    .register-enterprise
      register(v-if="registe" @close="close" @sureClose="sureClose")
      .register-content
        .content-header
          .left-box
              span.title 申请试用
              span.line
              span.content-tip 填写信息
          .right-box
            span.account 已有账号？
            el-button(type="text",@click="goLogin") 登录
        .content-form
          el-form(ref="registerForm" label-position="right" label-width="80px" :model="form_data" :rules="rules" size="mini")
            el-form-item(label="试用版本" class='version' required )
              .rad-box.box1(:class="first?'is-checked':''" @click="currcheckedFir")
                img.rad-img1(src="@/assets/business.png")
                .rad-box-rig
                  .tit 电商版本
                  .desc 适合仅线上商城运营的商户
              .rad-box(:class="checked?'is-checked':''" @click="currchecked")
                img.rad-img1(src="@/assets/sell.png")
                .rad-box-rig
                  .tit 零售版本
                  .desc 适合做线上商城/线下门店联运的商户
            el-form-item(label="企业名称" prop="company_name")
              el-input.content-input(v-model.trim="form_data.company_name" placeholder="请输入企业名称" maxlength="50")
            el-form-item(label="行业" prop="industry_id")
              el-select.content-input(v-model="form_data.industry_id" filterable placeholder="请选择企业所在行业")
                el-option(v-for="item in industry_list" :key="item.id" :label="item.name" :value="item.id")   
            el-form-item(label="邮箱" prop="email")
              el-input.content-input(v-model.trim="form_data.email" placeholder="请输入企业邮箱（将用作你的登录账号）")
            el-form-item(label="姓名" prop="iptname")
              el-input.content-input(v-model.trim="form_data.iptname" placeholder="请输入联系人姓名" maxlength="20")
            el-form-item(label="手机号" prop="mobile")
              el-input.content-input(v-model.trim="form_data.mobile" onkeyup = "value=value.replace(/[^\\d]/g,'')" placeholder="请输入登录手机号" maxlength="11")
            el-form-item(label="验证码" prop="idecode")
              el-input(v-model.trim="form_data.idecode" class="code-ipt" onkeyup = "value=value.replace(/[^\\d]/g,'')" placeholder="请输入4位验证码")
              el-button(@click="sendCode",v-if="!is_code",type="primary" plain) 获取验证码
              el-button.width-112(v-else disabled type="primary" plain) {{second}}s后可重发 
            el-form-item()
              el-checkbox(v-model="form_data.checkbox",class="check")
              span 我已阅读并同意
                el-button(type="text",@click="goRegister") 《用户注册协议》
        .bot-box
          el-button.confirm-btn(type="primary" @click="submit") 开始试用 
            i(class="el-icon-loading el-icon--right" v-if="loading")
</template>

<script>
import Register from "@/components/Register";
import Api from "@/api/public/Common";
import { checkMobile, checkEmail } from "@/utils/FormCheck.js";
// import { debuglog } from "util";
export default {
  props: {
    corp_name: {
      type: String
    },
    auth_code: {
      type: String
    },
    corp_id: {
      type: String
    },
    auth_enterprise: {
      type: Object
    }
  },
  components: {
    Register
  },
  data() {
    return {
      checked: false,
      first: false,
      second: 60,
      is_code: false,
      registe: false,
      loading: false,
      form_data: {
        version: 2,
        company_name: "",
        iptname: "",
        mobile: "",
        email: "",
        industry_id: "",
        idecode: "",
        checkbox: false // 是否同意注册协议
      },
      rules: {
        company_name: [
          {
            required: true,
            message: "请输入公司名称",
            trigger: ["blur"]
          }
        ],
        industry_id: [
          {
            required: true,
            message: "请选择所属行业",
            trigger: "change"
          }
        ],
        iptname: [
          {
            required: true,
            message: "请输入姓名",
            trigger: ["blur", "change"]
          }
        ],
        mobile: [
          {
            required: true,
            message: "请输入手机号",
            trigger: "blur"
          },
          { validator: checkMobile, trigger: "blur" }
        ],
        email: [
          {
            required: true,
            message: "请输入邮箱",
            trigger: "blur"
          },
          { validator: checkEmail, trigger: "blur" }
        ],
        idecode: [
          {
            required: true,
            message: "请输入验证码",
            trigger: "blur"
          }
        ]
      },
      industry_list: [] //行业列表
    };
  },
  methods: {
    goLogin() {
      // 点击登录！
      this.$emit("openLogin", 1);
    },
    // 点击用户协议
    goRegister() {
      this.registe = true;
    },
    // 点击零售版本
    currchecked() {
      this.checked = true;
      if (this.first) {
        this.first = false;
      }
    },
    // 点击电商版本
    currcheckedFir() {
      this.first = true;
      if (this.checked) {
        this.checked = false;
      }
    },
    // 发送验证码
    sendCode() {
      let self = this;
      const reg = /^1[3456789]\d{9}$/;
      if (self.form_data.mobile === "" || !reg.test(self.form_data.mobile)) {
        return false;
      } else {
        this.getMsg();
      }
      // this.$refs["registerForm"].validate(valid => {
      //   this.$refs["registerForm"].clearValidate(["idecode"]);
      //   console.log(valid, "1111111111111");
      //   if (self.form_data.mobile === "" || !reg.test(self.form_data.mobile)) {
      //     return false;
      //   } else {
      //     this.getMsg();
      //   }
      // });
    },
    // 获取验证码接口
    getMsg() {
      let params = {
        mobile: this.form_data.mobile,
        type: 0
      };
      Api.loginCode(params).then(() => {
        this.$message.success("验证码发送成功");
        this.is_code = true;
        this.timer = setInterval(() => {
          if (this.second > 1) {
            this.second--;
          } else {
            clearInterval(this.timer);
            this.second = 60;
            this.is_code = false;
          }
        }, 1000);
      });
    },
    submit() {
      if (!this.checked && !this.first) {
        this.$message.error("请选择试用版本");
        return;
      }
      this.$refs["registerForm"].validate(valid => {
        if (!valid) {
          return;
        }
        if (this.form_data.checkbox != true) {
          this.$message.error("请阅读并同意用户注册协议");
          return;
        } else {
          const params = {
            version: this.checked ? 2 : 3,
            corp_id: this.corp_id || this.auth_enterprise.corp_id,
            name: this.form_data.company_name,
            industry_id: this.form_data.industry_id,
            email: this.form_data.email,
            contact: this.form_data.iptname,
            mobile: this.form_data.mobile,
            verification_code: this.form_data.idecode,
            memUserid:
              this.$route.query.memUserid || this.auth_enterprise.memUserid
          };
          this.loading = true;
          Api.autoAccess(params).then(res => {
            const { code, data } = res;
            // localStorage.setItem("TOKEN", JSON.stringify(tmp_token));
            if (code === 200) {
              this.loading = false;
              this.$emit("success", data);
            } else {
              this.$message.error(res.message);
            }
          });
        }
      });
    },
    close() {
      this.registe = false;
    },
    sureClose() {
      this.form_data.checkbox = true;
      this.registe = false;
    },
    // 获取行业列表
    getIndustryList() {
      const params = {
        page: 1,
        size: 100
      };
      Api.industryList(params).then(res => {
        // this.industry_list = res.data.list;
        this.industry_list = res.list;
      });
    }
  },
  mounted() {
    // this.$set(this.form_data, "name", this.auth_enterprise.epName);
    this.form_data.company_name = this.corp_name;
    this.getIndustryList();
  }
};
</script>

<style lang="sass" scoped>
/deep/ .el-form-item--mini .el-form-item__error
  padding-top: 5px
/deep/ .el-button--text
  color: #708DF5
/deep/ .el-button--text:focus, .el-button--text:hover
  color: rgba(112,141,245,0.8)
/deep/ .el-button--primary.is-plain
  color: #606266
  border: 1px solid #DCDFE6
  background: transparent
/deep/ .el-button--primary.is-plain:hover, .el-button--primary.is-plain:focus
  color: #708DF5
  border-color: #d4ddfc
  background-color: #f1f4fe
/deep/ .el-button--primary.is-plain.is-disabled:hover
  background-color: #F5F7FA
  border-color: #E4E7ED
  color: #C0C4CC
  cursor: not-allowed
/deep/ .el-button--primary.is-plain.is-disabled
  background-color: #F5F7FA
  border-color: #E4E7ED
  color: #C0C4CC
  cursor: not-allowed
@media screen and (max-width: 1680px)
  .register-content
    // transform: scale(0.8)
.box1
  margin-right: 10px
.rad-box
  width: 400px
  height: 45px
  cursor: pointer
  border-radius: 4px
  border: 1px solid #DCDFE6
  display: flex
  padding: 15px
  &.is-checked , &.web-checked
    position: relative
    border: 1px solid #5B89FE
    &:after
      content: ""
      display: block
      height: 20px
      width: 20px
      background-image: url(~@/assets/verson-chose.png)
      background-size: 100%
      position: absolute
      right: 0
      bottom: 0
  .rad-img1
    display: block
    width: 45px
    height: 45px
    border-radius: 4px
  .rad-box-rig
    font-size: 14px
    margin-left: 10px
    line-height: initial
    .tit
      color: #606266
    .desc
      color: #909399
      overflow: hidden
      text-overflow: ellipsis
      display: -webkit-box
      -webkit-box-orient: vertical
      -webkit-line-clamp: 2
.code-ipt
  width: 290px !important
  margin-right: 25px
  /deep/ .el-input__inner
    width: 290px !important
    margin-right: 30px
.top-box
  display: flex
.bot-box
  text-align: left
  margin-left: 75px
  button
    width: 400px
    height: 36px
    padding-top: 10px
    background: #708DF5
    box-shadow: 0 11px 21px 0 rgba(98, 143, 255, 0.32)
  button:hover
    background: rgba(112,141,245,0.8)
  button:active
    background: #5471DF
/deep/.el-form-item
  margin-bottom: 30px !important
/deep/.el-input--mini .el-input__inner
  height: 36px
  line-height: 36px
.register-box
  width: 1005px
  position: fixed
  left: 0
  top: 0
  bottom: 0
  right: 0
  margin: auto
  margin-top: 100px
  margin-bottom: 70px
  overflow-y: auto
  overflow-x: hidden
  background-color: white
  border-radius: 10px
// 滚动条整体样式
::-webkit-scrollbar
  width: 20px
  height: 8px
// 滚动条里面小方块
::-webkit-scrollbar-thumb
  border-radius: 12px
  border: 6px solid rgba(0, 0, 0, 0)
  box-shadow: 8px 0 0 #909399 inset
::-webkit-scrollbar-thumb:hover
  box-shadow: 8px 0 0 #909399 inset
.register-enterprise
  width: 1006px
  height: 750px
  padding: 40px 50px 50px 50px
  // position: absolute
  // left: 0
  // top: 0
  // bottom: 0
  // right: 0
  // margin: auto
  z-index: 100
  border-radius: 10px
  background-color: white
  box-sizing: border-box
  /deep/.el-checkbox
    margin-right: 5px !important
  .register-content
    width: 100%
    height: 100%
    .content-header
      margin-bottom: 20px
      display: flex
      justify-content: space-between
      align-items: center
      .left-box
        display: flex
        align-items: center
        .line
          display: inline-block
          width: 15px
          height: 1px
          margin: 0 10px
          background: #909399
      .title
        color: #303133
        font-weight: 600
        font-size: 24px
    .content-tip
      font-size: 18px
      color: #909399
      margin-left: 4px
    .content-form
      /deep/ .el-input__inner
        width: 400px
        border-color: #DCDFE6
        font-size: 14px
      /deep/ .el-form-item.is-error .el-input__inner
        border-color: #F56C6C
      /deep/ .el-input__inner:focus
        border-color: #708DF5
      /deep/ .el-form-item__content
        display: flex
      /deep/ .el-checkbox__input.is-checked .el-checkbox__inner
        background-color: #708DF5
        border-color: #708DF5
      .change-group
        line-height: 36px
        margin-top: 6px
      .question
        width: 14px
        height: 14px
        display: inline-block
        margin-left: 5px
        vertical-align: top
      .content-input
        // width: 340px !important
      .confirm-btn
        width: 400px !important
        background-color: #1989FA
        color: #FFFFFF
        border-color: #1989FA
      .btn:hover,.btn:focus
        color: #FFFFFF
        background-color: #1989FA
        border-color: #1989FA
</style>
