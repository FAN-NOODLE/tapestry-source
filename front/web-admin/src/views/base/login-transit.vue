<template>
  <div class="login-transite">
    <div class="login-box">
      <el-tabs v-model="loginType">
        <!-- 账号密码自由登录 -->
        <el-tab-pane
          label="输入密码登录"
          name="freedom"
        >
          <div
            class="login-form"
            v-show="login_show_module == 1"
          >
            <el-input
              class="m-t-20 login-input"
              v-model="info.account"
              size="medium"
              placeholder="请输入账号"
            ></el-input>
            <el-input
              class="m-t-30 login-input"
              v-model="info.password"
              size="medium"
              type="password"
              placeholder="请输入密码"
            ></el-input>
            <div class="code-input">
              <el-input
                class="m-t-20 login-input"
                type="text"
                v-model="info.code"
                placeholder="请输入图形验证码"
                maxlength="6"
              ></el-input>
              <el-image fit="cover"
                class="image-validate"
                :src="verifyData.base64_img"
              ></el-image>
            </div>
          </div>
          <div
            class="sms-box"
            v-show="login_show_module == 2"
          >
            <div class="text-wraper">
              验证码将发给到注册账号手机，请注意查收
            </div>
            <el-input
              type="text"
              class="m-t-20 login-input"
              v-model="info.numberCode"
              maxlength="6"
              placeholder="请输入手机验证码"
            ></el-input>
            <div class="send-vcode-btn">
              验证码（911800）
            </div>
          </div>
        </el-tab-pane>
        <!-- 选自企业自动登录 -->
        <el-tab-pane
          label="企业自动登录"
          name="wisdom"
        >
          <div class="login-form">
            <cy-select
              class="m-t-20 login-select"
              v-model="enterprise"
              placeholder="请选择想要登录的企业"
              @change="changeEnterPrise"
              search-input
              :options-list="AccountList"
              :options="{ label: 'name', value: 'name' }"
            >
            </cy-select>
            <p
              class="tips tips-one"
              v-if="!enterprise"
            >请先选择要登录的企业</p>
            <p
              class="tips tips-two"
              v-else
            >
              {{ `您即将以${enterprise}超管的身份进行登录` }}
            </p>
          </div>
        </el-tab-pane>
        <el-button
          v-if="login_show_module == 1"
          class="login-button"
          type="primary"
          size="medium"
          @click="login"
        >登录</el-button>
        <el-button
          v-if="login_show_module == 2"
          class="login-button"
          type="primary"
          size="medium"
          @click="letsGoLogin"
        >提交</el-button>
      </el-tabs>
    </div>
  </div>
</template>
<script lang="ts">
interface AccountItem {
  name: string;
  account: string;
  pswd: string;
}
import { Component, Vue, Watch } from "vue-property-decorator";
import AccountList from "@/utils/account";
import { IUserInfo } from "@/types";
import UserModule from "@/store/modules/user";
import PerModule from "@/store/modules/permission";
// import LoginAboat from "@/utils/LoginManager";
import Api from "@/api/public/public";
import JSEncrypt from 'jsencrypt'
@Component({})
export default class extends Vue {
  @Watch("loginType")
  onChangeLoginType() {
    this.enterprise = "";
    this.$set(this.info, "account", "");
    this.$set(this.info, "password", "");
    this.$set(this.info, "code", "");
  }
  //   get USER() {
  //     return UserModule.getUserInfo;
  //   }
  private loginType = "freedom"; //登录类型
  private enterprise = ""; //企业登录选中的企业名字
  private info = {
    password: "",
    account: "",
    code: "",
    mobile: "",
    numberCode: ""
  };
  private verifyData = {
    code_key: "",
    base64_img: ""
  };
  private phoneCode = {
    is_code: false,
    time: 60
  };
  private login_show_module = 1; // 1： 账号登录；2：短信验证
  private AccountList: Array<AccountItem> = [];
  mounted() {
    this.AccountList = AccountList;
    this.getValidate();
  }
  // methods
  private changeEnterPrise() {
    if (this.enterprise) {
      const { pswd, account } = this.AccountList.filter(
        item => item.name === this.enterprise
      )[0];
      Object.assign(this.info, { account, password: pswd });
    }
  }

  // 验证账号密码
  private async login() {
    let tempAccount = "";
    if (this.info.account.indexOf("@") < 0) {
      tempAccount = `${this.info.account}`;
    } else {
      tempAccount = this.info.account;
    }
    const res: any = await Api.login({
      login_type: 1,
      account: tempAccount,
      password: this.encryptPassword(this.$md5(this.info.password + "")),
      code: this.info.code,
      code_key: this.verifyData.code_key,
      enterprise_id: 3363603667487424 // 本地调试登录
    });
    if (res.code == 200) {
      this.login_show_module = 2;
      this.info.mobile = res.data.token;
    } else {
      this.$message.error(res.message);
    }
  }
  encryptPassword(password) {
      
    // 2. 准备公钥
    const publicKey = `-----BEGIN PUBLIC KEY-----
        MIIBCgKCAQEAuCyFjsc9bNmWJDGiT2HlFz97My22ZwKWRQ8gLxhckVvQLoRbR24F
        609EuF1N7U9ioP0QJV+zhRgNGGmHJLsCd1QyaT+5yaNBUS2hyDfDwF3w7G971Quv
        CRK2WGtQ/E1GiiavMvKI5IOioidzpSo9uzt3CyDsfY3ieNNUkkZrMtf0ay1mYQh7
        qg46IJEWpnR/lhxunlb/JSaet+ZfeRKnHdXH0K60/hA88BMWzxO7Uj3DQjIYoyfj
        QlG9YvoKa53VrRIEMP1I4QoUXayhBDR61sHClKHy6dbDOz/UZ1XIdh+MUZt0Wjcn
        u+Oh8DTcM/L1D2eAXtM044bBrGgrNfVekQIDAQAB
        -----END PUBLIC KEY-----`.replace(/\n\s/g, '');
    const encrypt = new JSEncrypt();
    
    encrypt.setPublicKey(publicKey);
    const encrypted = encrypt.encrypt(password);
    
    // 4. 组合传输数据
    return encrypted;
  }

  // 手机验证码登录
  private async letsGoLogin() {
    const {
      data: { hash, token, enterprise_list }
    } = await Api.login({
      login_type: 2,
      account: `${this.info.mobile}`,
      password: this.info.numberCode,
      enterprise_id: 3363603667487424 // 本地调试登录
    });
    localStorage.setItem("TOKEN", JSON.stringify(token));
    localStorage.setItem("LOGIN_HASH", JSON.stringify(hash));
    this.getEnterpriseToken(enterprise_list[0]);
  }
  //通过token登陆所选企业
  private async getEnterpriseToken(target_enterprise: any) {
    localStorage.setItem("REAL_HASH", JSON.stringify(target_enterprise.hash));
    const {
      data: {
        hash,
        enterprise_admin_info,
        token,
        permissions_list,
        permissions_tree,
        role_list
      }
    } = await Api.getEnterpriseToken({
      admin_id: target_enterprise.admin_id,
      enterprise_id: target_enterprise.enterprise_id || 3363603667487424
    });
    localStorage.setItem("LOGIN_HASH", JSON.stringify(hash));
    localStorage.setItem("ENTERPRISE_ID", enterprise_admin_info.enterprise_id);
    localStorage.setItem("TOKEN", JSON.stringify(token));
    localStorage.setItem("KEYS", JSON.stringify(permissions_list));
    PerModule.setPermission(permissions_tree); //存储权限树
    PerModule.setRoles(role_list); //存储角色list
    this.authorizationNew(enterprise_admin_info);
  }
  private authorizationNew(res: IUserInfo) {
    UserModule.setUserInfo(res);
    this.$router.push({
      name: "myapp_home"
    });
  }
  // 获取账号密码的验证码
  private async getValidate() {
    const { code, data }: any = await Api.validate();
    if (code === 200) {
      const { code_key, base64_img } = data;
      this.verifyData.code_key = code_key;
      this.verifyData.base64_img = base64_img;
    }
  }
}
</script>
<style lang="scss" scoped>
::v-deep {
  .el-tabs {
    width: 346px;
    margin: 0 auto;
    margin-top: 60px;

    .el-tabs__nav {
      transform: translateX(30px) !important;
    }

    .el-tabs__nav-wrap {
      &:after {
        height: 0;
      }
    }

    .login-input .el-input__inner {
      border-left: none !important;
      border-right: none !important;
      border-top: none !important;
      border-radius: 0 !important;
    }
  }
}

.login-transite {
  width: 100%;
  height: 100vh;
  background-image: url("https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/login/login-bac.png");
  background-repeat: no-repeat;
  background-size: cover;
  overflow: hidden;
  position: relative;

  .login-box {
    width: 546px;
    height: 435px;
    background-color: #ffffff;
    position: absolute;
    left: 0;
    right: 0;
    top: 0;
    bottom: 0;
    margin: auto;
    border-radius: 8px;
    overflow: hidden;

    .login-form {
      // height: 170px;
      margin-bottom: 30px;
    }

    .login-select {
      width: 334px !important;
    }

    .login-button {
      width: 346px;
    }

    .tips {
      margin-top: 40px;
    }

    .tips-one {
      color: $color-warn;
    }

    .tips-two {
      color: $color-success;
    }
  }
}

.code-input {
  // float: left;
  position: relative;
  width: 100%;
}

.image-validate,
.code-btn {
  float: left;
  height: 45px;
  width: 100px;
  margin-top: 15px;
}

.image-validate {
  position: absolute;
  right: 0;
  bottom: 1px;
}

.sms-box {
  padding-bottom: 70px;

  .text-wraper {
    padding: 20px;
  }

  /deep/input[type="number"] {
    &::-webkit-outer-spin-button,
    &::-webkit-inner-spin-button {
      -webkit-appearance: none;
    }

    -moz-appearance: textfield;
  }

  .send-vcode-btn {
    position: absolute;
    top: 86px;
    right: 10px;
    color: #c0c4cc;
    font-weight: 600;
    cursor: pointer;
    z-index: 99;
  }
}
</style>
