<template>
  <div class="login">
    <div class="login-black-box">
      <div class="login-box">
        <el-tabs v-model="active">
          <!-- 账号验证码 -->
          <el-tab-pane v-if="active==='1'" name="1" label="重置密码">
            <el-form :model="formData" ref="form" :rules="formRules">
              <div class="text">密码将会被重置为账号绑定的手机号</div>
              <el-form-item prop="account">
                <el-input
                  v-model="formData.account"
                  maxlength="50"
                  @blur="getNewEnterpriseList"
                  placeholder="请输入账号"
                ></el-input>
              </el-form-item>

              <el-form-item
                v-if="enterpriseList.length > 1"
                prop="enterprise_id"
              >
                <el-select
                  v-model="formData.enterprise_id"
                  placeholder="请选择企业"
                >
                  <el-option
                    v-for="(item, index) in enterpriseList"
                    :key="index"
                    :label="item.name"
                    :value="item.enterprise_id"
                  ></el-option>
                </el-select>
              </el-form-item>

              <el-form-item prop="accountCode">
                <div class="flex">
                  <el-input
                    type="text"
                    v-model="formData.accountCode"
                    placeholder="请输入验证码"
                    @keyup.enter.native="handleReset()"
                  ></el-input>
                  <el-image fit="cover"
                    :src="verifyData.base64_img"
                    @click="getValidate"
                  ></el-image>
                </div>
              </el-form-item>
              <el-button type="primary" :loading="loading" @click="handleReset()">重 置</el-button>

            </el-form>
          </el-tab-pane>
          <!-- 短信验证码 -->
          <el-tab-pane v-if="active==='2'" name="2" label="短信验证">
            <el-form :model="formData" ref="form" :rules="formRules">
              <div class="text">验证码将发给到注册账号手机，请注意查收</div>
              <el-form-item prop="numberCode">
                <div class="flex">
                  <el-input 
                    v-model.number="formData.numberCode"
                    placeholder="请输入6位验证码"
                    maxlength="6"
                    @keyup.enter.native="handleValidateCode()"
                  ></el-input>
                  <div v-if="!is_code" class="send-vcode-btn" @click="handleSendCode()">获取验证码</div>
                  <div v-else class="send-vcode-btn"> {{time}}秒后再次发送</div>
                </div>
              </el-form-item>
              <el-button type="primary" :loading="loading" @click="handleValidateCode()">提 交</el-button>
            </el-form>
          </el-tab-pane>
        </el-tabs>
      </div>
    </div>
  </div>
</template>
<script>
import Api from "@/api/public/public";
export default {
  data() {
    return {
      is_code: false, // 是否发送验证码
      time: 60, // 5测试 正常60s
      interval: null,
      loading: false,
      active: "1",
      verifyData: {
        code_key: "",
        base64_img: ""
      },
      formData: {
        account: "",
        enterprise_id: "",
        accountCode: "",
        numberCode: "",
        mobile: ""
      },
      formRules: {
        account: [
          {
            required: true,
            message: "请输入账号",
            trigger: "blur"
          }
        ],
        enterprise_id: [
          {
            required: true,
            message: "请选择企业",
            trigger: "blur"
          }
        ],
        accountCode: [
          {
            required: true,
            message: "请输入验证码",
            trigger: "blur"
          }
        ],
        numberCode: [
          {
            required: true,
            message: "请输入短信验证码",
            trigger: "blur"
          }
        ]
      },
      enterpriseList: []
    };
  },
  created() {
    this.getValidate()
  },
  methods: {
    /**
     * 获取验证码
     */
    getValidate() {
      Api.validate()
        .then(res => {
          if (res.code === 200) {
            this.verifyData = res.data;
          }
        })
        .catch(err => {
          console.log(err);
        });
    },
    /**
     * 新的根据用户获取相关企业列表
     */
    getNewEnterpriseList() {
      const params = {
        account: this.formData.account
      };
      Api.getEnterpriseList(params).then(res => {
        this.enterpriseList = res.data.enterprise_admin_list;
        if (this.enterpriseList && this.enterpriseList.length === 1) {
          this.formData.enterprise_id =
            this.enterpriseList[0].enterprise_id || '';
        } else if (!this.enterpriseList.length) {
          this.formData.enterprise_id = '';
          this.$message.error('该账号无权限，请联系管理员')
        }
      });
    },
    /**
     * 重置
     */
    handleReset() {
      this.$refs["form"].validate(valid => {
        if(valid) {
          this.loading = true;
          const params = {
            account: this.formData.account,
            enterprise_id: this.formData.enterprise_id || 0,
            code: this.formData.accountCode,
            code_key: this.verifyData.code_key
          };
          
          Api.checkResetPassword(params).then(res => {
            if(res.code === 200) {
              this.loading = false;
              this.active = '2'
              this.formData = {
                ...this.formData,
                ...res.data
              }
              this.handleSendCode()
            } else {
              this.loading = false;
            }
          
          }).catch(() => {
            this.loading = false;
          })
        } else {
          this.loading = false;
        }
      });
    },
    /** 验证码验证 */
    handleValidateCode() {
      const params = {
        account: this.formData.account,
        enterprise_id: this.formData.enterprise_id || 0,
        mobile: this.formData.mobile,
        code: this.formData.numberCode,
      };
      Api.resetPassword(params).then(res => {
        if(res.code === 200) {
          this.$message({
            type: 'success',
            message: '密码已重置，请重新登录'
          })
          window.location.href = `${window.location.origin}/login`;
        } else {
          this.$message.error(res.message)
        }
      })
    },
    /**
     * 获取验证码
     */
    handleSendCode() {
      let params = {
        type: 1,
        mobile: this.formData.mobile,
        enterprise_id: this.formData.enterprise_id || 0
      };
      this.loading = true
      Api.sendCode(params).then(
        res => {
          if (res.code === 200) {
            this.$message({
              type: "success",
              message: "验证码发送成功",
              duration: 1000
            });
            this.loading = false
            this.setInterval()
          } else {
            this.loading = false
            this.$message.error(res.message);
          }
        },
        () => {
          this.loading = false
        }
      );
    },
    /**
     * 提交
     */
    setInterval() {
      this.is_code=true
      this.interval = setInterval(() => {
        this.time--;
        if (this.time == 0) {
          // 停止定时器
          clearInterval(this.interval);
          this.is_code = false;
          this.time = 60;
          return;
        }
      }, 1000);
    }
  }
};
</script>

<style lang="scss" scoped>
.login {
  width: 100vw;
  height: 100vh;
  min-width: 1200px;
  min-height: 680px;
  overflow: auto;
  background-image: url('~@/assets/images/login-bac.png'), url('~@/assets/images/bg_mask.png');
  background-repeat: no-repeat;
  background-size: cover, 1338px 880px;
  background-position: center;
  position: relative;
}
.login-black-box{
  display: flex;
  -ms-flex-pack: justify;
  justify-content: space-between;
  position: relative;
  width: 1080px;
  // height: 605px;
  top: 50%;
  left: 50%;
  transform: translate(-50%,-50%)
}
// 登录模块
.login-box{
  width: 520px;
  // height: 605px;
  padding: 60px 60px;
  border-radius: 10px;
  background: #fff;
  z-index: 50;
  box-sizing: border-box;
  margin: 0 auto;

  .flex {
    display: flex;

    .el-image {
      height: 40px;
      width: 115px
    }
  }
  .el-button {
    height: 46px;
    font-size: 18px;
    display: block;
    width: 100%;
    box-shadow: 0 11px 21px 0 rgba(0, 0, 0, 0.2);
    margin-top: 100px;
  }
  .send-vcode-btn{
    color: #c0c4cc;
    font-weight: 600;
    cursor: pointer;
    z-index: 99;
    white-space: nowrap;
  }
  .tips {
    color: #999;
    margin-top: 10px;
  }
}
::v-deep {
  .el-tabs__item {
    font-size: 24px;
  }
  .el-tabs__nav-wrap {
    &::after {
      display: none;
    }
  }
  .el-tabs__active-bar {
    padding: 0 !important;
    width: 100px !important;
  }

  .el-form-item {
    margin-top: 30px;
  }

  .el-input {
    font-size: 16px;
  }

  .el-select {
    width: 100%;
  }

  .el-input__inner {
    padding: 0;
    border: none;
    border-bottom: solid 1px #E4E7EC;
    border-radius: 0px;
    &:focus {
      border-color: #000000
    }
  }
  .el-form-item.is-error .el-input__inner {
    border-color: #F56C6C
  }
  .el-form-item.is-success .el-input__inner:focus {
    border-color: #0000
  }
} 
    
</style>
