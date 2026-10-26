<template lang="pug">
  //- 当非超管用户首次扫码会弹出
  .account-bind
    .bind-content
      .content-header
        .header-title 绑定已有账号
        img.close(src="@/assets/close-icon.png" @click="close")
        .header-tip 绑定账号后，企业微信和手机号均可登录畅鲟
      .content-form
        el-form(ref="bindForm" :model="form_data" :rules="rules") 
          el-row
            el-col
              el-form-item(prop="mobile" :error="phoneErrorMsg")
                el-input(placeholder="请输入手机号" v-model="form_data.mobile" @input="checkInt(form_data,'mobile',form_data.mobile)") 
          el-row
            el-col(:span="17")
              el-form-item(prop="code" :error="codeErrorMsg") 
                el-input(placeholder="请输入验证码" v-model="form_data.code" @input="checkInt(form_data,'code',form_data.code)")
            el-col(:span="6" :offset="1")  
              el-form-item 
                el-button(@click="getCode",v-if="!is_code",type="primary" plain) 获取验证码
                el-button.code-tip(v-else disabled type="primary" plain) {{second}}s后可重发 
      el-button.btn(@click="login") 登录        
</template>
<script>
import Api from "@/api/public/Common";
import CheckNum from "@/mixins/CheckNum";
import { _limitFrequence } from "@/utils/LimitFunction";
export default {
  props: {
    type: {
      type: String
    },
    token: {
      type: String
    },
    auth_code: {
      type: String
    },
    corp_id: {
      type: String
    },
    epId: {
      type: [String, Number]
    },
    work_wechat_login_token: {
      type: String
    }
  },
  mixins: [CheckNum],
  data() {
    const checkMobile = (rule, value, callback) => {
      const reg = /^1[3456789]\d{9}$/;
      if (!reg.test(value)) {
        callback(new Error("请填写正确的手机号"));
      } else {
        callback();
      }
    };
    return {
      phoneErrorMsg: "",
      codeErrorMsg: "",
      form_data: {
        mobile: "", //手机号
        code: "" //短信验证码
      },
      rules: {
        mobile: [
          {
            required: true,
            message: "请输入手机号",
            trigger: ["blur", "change"]
          },
          { validator: checkMobile, trigger: "blur" }
        ],
        code: [
          {
            required: true,
            message: "请输入验证码",
            trigger: ["blur", "change"]
          }
        ]
      },
      is_code: false, //是否正在获取验证码倒计时
      second: 60,
      timer: null
    };
  },
  methods: {
    // 获取验证码，校验手机号，不校验验证码
    getCode() {
      this.$refs["bindForm"].validate(() => {
        this.$refs["bindForm"].clearValidate("code");
        const reg = /^1[23456789]\d{9}$/;
        if (this.form_data.mobile === "" || !reg.test(this.form_data.mobile)) {
          return false;
        } else {
          const params = {
            event: "word_wechat_bind_account",
            enterprise_id: this.epId,
            mobile: this.form_data.mobile,
            length: 6
          };
          this.sendCode(params);
        }
      });
    },
    restart() {
      this.sendCode = _limitFrequence(function(params) {
        Api.sendCodeSms(params).then(res => {
          if (res.code === 200) {
            this.$message.success("验证码发送成功");
            this.is_code = true;
            this.timer = setInterval(() => {
              if (this.second > 1) {
                this.second--;
              } else {
                this.restart();
                clearInterval(this.timer);
                this.second = 60;
                this.is_code = false;
              }
            }, 1000);
          } else {
            self.is_code = false;
            this.$message.error(res.message);
          }
        });
      });
    },
    // 点击登录按钮
    login() {
      let self = this;
      this.phoneErrorMsg = "";
      this.codeErrorMsg = "";
      this.$refs["bindForm"].validate(valid => {
        if (!valid) {
          return false;
        } else {
          const params = {
            type: this.type,
            work_wechat_login_token: this.work_wechat_login_token,
            bind_account: {
              mobile: this.form_data.mobile,
              code: this.form_data.code
            }
          };
          if (this.type === "scan_code") {
            params.auth_code = this.auth_code;
          } else if (this.type === "work_wechat") {
            params.token = this.token;
            params.corpid = this.corp_id;
          }
          Api.workWechatLogin(params).then(res => {
            if (res.code === 200) {
              const { code, tmp_token } = res.data;
              localStorage.setItem("TOKEN", JSON.stringify(tmp_token));
              switch (code) {
                case "success":
                  this.$message({
                    type: "success",
                    message: "登录成功!"
                  });
                  setTimeout(function() {
                    self.$emit("success", res.data);
                  }, 1000);
                  break;
                case "verify_code_error":
                  this.$nextTick(() => {
                    this.codeErrorMsg = "验证码有误，请重新输入";
                  });
                  this.$emit("error", res.data);
                  break;
                case "admin_not_found":
                  this.$nextTick(() => {
                    this.phoneErrorMsg = "该手机号未注册，请联系管理员开通";
                  });
                  this.$emit("error", res.data);
                  break;
              }
            } else {
              this.$message.error(res.message);
            }
          });
        }
      });
    },
    // 关闭按钮回调
    close() {
      this.$emit("close");
    }
  },
  mounted() {
    this.restart();
  }
};
</script>
<style lang="sass" scoped>
.account-bind
  width: 100%
  height: 100vh
  position: fixed
  z-index: 100
  background-color: rgba(0,0,0,0.7)
  .bind-content
    width: 460px
    height: 309px
    position: absolute
    top: 0
    left: 0
    right: 0
    bottom: 0
    margin: auto
    border-radius: 10px
    background-color: #FFFFFF
    box-sizing: border-box
    padding: 30px 30px 0 30px
    .content-header
      margin-bottom: 30px
      position: relative
      &:after
        clear: both
        content: ""
        display: block
      .header-title
        font-size: 18px
        line-height: 26px
        color: #292A30
        font-weight: 600
      .header-tip
        color: #909399
        font-size: 14px
        line-height: 20px
        margin-top: 5px
      .close
        position: absolute
        width: 14px
        height: 14px
        float: right
        cursor: pointer
        margin-top: 6px
        right: 0
        top: 0
    /deep/.el-form-item
      margin-bottom: 29px
    /deep/.el-form-item__content
      line-height: 36px
    /deep/.el-input__inner
      height: 36px
      line-height: 36px
    /deep/.el-button--primary.is-plain
      height: 36px
      width: 100px
      padding: 0 0
      text-align: center
      color: #708DF5
      background: #FFFFFF
      border-color: #708DF5
    .btn
      width: 100%
      color: #FFFFFF
      background-color: #708DF5
      border-color: #708DF5
    .btn:hover, .btn:focus
      color: #FFFFFF
      background-color: #708DF5
      border-color: #708DF5
    .content-form
      .code-tip
        color: #c0c4cc
        border-color: #c0c4cc
</style>
