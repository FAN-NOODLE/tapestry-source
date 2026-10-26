<template lang="pug">
  .login
    //- .logo-click(@click="goOfficial")
    //- .copy-right Copyright ©2014-2020 畅移（上海）信息科技有限公司
    img.bg-mask(src="@/assets/bg_mask.png")

    //- 非超管用户首次扫码登录绑定手机号弹窗
    account-bind(
      v-if="bind_show"
      :type="wxtype"
      :token="wx_token"
      :work_wechat_login_token="work_wechat_login_token"
      :auth_code="auth_code"
      :corp_id="corpid"
      :epId="epId"
      @close="closeBind" 
      @success="successBind" 
      @error="errorBind"
    )
    //- 试用过期出现的弹窗
    account-overdue( 
      v-if="overdue_show" 
      @close="closeOverdue"
    )
    //- 申请试用填写信息表单组件
    register-enterprise(
      v-if="register_show"
      :type="wxtype"
      :token="wx_token"
      :auth_code="auth_code"
      :corp_id="corpid"
      :auth_enterprise="auth_enterprise"
      :corp_name="corp_name"
      :work_wechat_login_token="work_wechat_login_token"
      @close="closeRegister" 
      @success="successRegister",
      @openLogin="closeRegister",
    )
    //- 防刷滑块验证
    slide-bar.slide-bar(
      v-if="valid_show"
      @success="validSuccess" 
      @close="closeValid" 
      @fail="validFail"
    )
    //- 选择登录企业弹窗
    //- select-enterprise(
    //-   v-if="is_show_select_ent"
    //-   :show="is_show_select_ent" 
    //-   :listData="enterpriseList" 
    //-   :on-select="selectEntListenner" 
    //-   :on-close="()=>{is_show_select_ent=false}"
    //- )
    //- 根据产品需求，如由第三方后台进入或扫码进入，默认不展示账号密码登录框（normal_login字段判断是否常规登录）
    .login-black-box(v-show="!register_show")
      //- 登录模块
      .login-box(v-show="show_module == 1 && normal_login && !register_show && !overdue_show && !bind_show")
        el-tabs(v-model="login_type")
          el-tab-pane.tab1-height(:label="labelName",name="password")
            el-form(:model="info" v-if="login_type === 'password'" ref="info" :rules="rule" form="mini")
              .account-info(v-show="login_show_module == 1")
                el-form-item.mobile(prop="account")
                  el-input(
                    class="input-with-select"
                    v-model.trim="info.account"
                    maxlength=50
                    placeholder="请输入账号")
                el-form-item(prop="password" v-if="!is_show")     
                  el-input(
                    type="password"
                    v-model.trim="info.password"
                    placeholder="请输入密码"
                    @keyup.enter.native="goLogin()")
                el-form-item(prop="password" v-else)
                  el-input(
                    type="text"
                    v-model="info.password"
                    placeholder="请输入密码"
                    @keyup.enter.native="goLogin()")
                  .see(@click="seePass",v-if="info.password && is_show")
                  .hide(@click="seePass",v-if="info.password && !is_show")
                //- el-form-item(prop="enterprise" v-if="newEnterpriseList.length > 1")
                //-   el-select(v-model="info.enterprise_id" placeholder="请选择企业" class="select-enterprise")
                //-     el-option(v-for="(item,index) in newEnterpriseList" :label="item.name" :value="item.enterprise_id")
                el-form-item(prop="accountCode")
                  el-input(
                    type="text"
                    v-model="info.accountCode"
                    placeholder="请输入验证码"
                    @keyup.enter.native="goLogin()")
                  .code-image(@click="getValidate")
                    el-image(class="image-validate" :src="verifyData.base64_img")
              //- 短信验证
              .sms-box(v-show="login_show_module == 2")
                .text 验证码将发给到注册账号手机，请注意查收
                el-form-item
                  el-input(
                    type="number"
                    v-model="info.numberCode"
                    placeholder="请输入6位验证码"
                    oninput="if(value.length>6)value=value.slice(0,6)"
                    @keyup.enter.native="testCode()")
                  .send-vcode-btn(@click="sendCode()" v-if="!is_code") 获取验证码
                  .send-vcode-btn(v-else) {{time}}秒后再次发送
              //- 新用户强制修改密码
              .edit-pass(v-show="login_show_module == 3")
                el-form(:model="passInfo" ref="passInfo" :rules="passRule" form="mini")
                  el-form-item(prop="originalPass")    
                    el-input(
                      type="password"
                      v-model.trim="passInfo.originalPass"
                      placeholder="原密码")
                  el-form-item(prop="newPass")    
                    el-input(
                      type="password"
                      v-model.trim="passInfo.newPass"
                      placeholder="新密码")
                  el-form-item(prop="confirmPass")    
                    el-input(
                      type="password"
                      v-model.trim="passInfo.confirmPass"
                      placeholder="确认密码"
                      @keyup.enter.native="editLogin()")
          //- 双因素登录,取消手机验证码登录，改为第二步短信验证
          //- el-tab-pane.tab2-height(label="手机验证码登录",name="code")
          //-   el-form(:model="info" v-if="login_type === 'code'" ref="info" :rules="rule" form="mini")
          //-     el-form-item.mobile(prop="mobile")
          //-       el-input(
          //-         class="input-with-select"
          //-         v-model.trim="info.mobile"
          //-         @input="checkInt(info,'mobile',info.mobile)"
          //-         maxlength=11
          //-         placeholder="请输入手机号")       
          //-     el-row.code-row-box(:class="checkCodeError ? 'red-line':''")
          //-       el-col(:span="16")
          //-         el-form-item(prop="code")
          //-           el-input(
          //-             type="text"
          //-             v-model="info.code"
          //-             placeholder="请输入6位验证码"
          //-             maxlength=6
          //-             @input="checkInt(info,'code',info.code)"
          //-             @keyup.enter.native="goLogin(login_type)")
          //-       el-col.code-box(:span="8",:offset="2")
          //-         el-button.code-btn.hover-color(@click="sendCode",v-if="!is_code",type="primary" plain) 获取验证码
          //-         el-button.code-btn.code-tip(v-else disabled type="primary" plain) 重新获取{{time}}s
          //-     el-form-item(prop="numberCode")
          //-       el-input(
          //-         type="text"
          //-         v-model="info.numberCode"
          //-         placeholder="请输入验证码"
          //-         @keyup.enter.native="goLogin(login_type)")
          //-       .code-image(@click="getValidate")
          //-         el-image(class="image-validate" :src="verifyData.base64_img")
        //- .save-seven-box
        //-   <el-checkbox v-model="sevenChecked" @change="changeSeven(sevenChecked)">7天内自动登录</el-checkbox>

        el-button(class="btn" v-if="login_show_module == 1" :loading="is_logining" @click="goLogin()") 登录
        el-button(class="btn" v-else-if="login_show_module == 2" :loading="is_logining" @click="testCode()") 提交
        el-button(class="btn" v-else-if="login_show_module == 3" :loading="is_logining" @click="editLogin()") 提交
        .tip-box
          el-button(v-show="login_show_module != 3" type='text' @click="handleForget()") 忘记密码？
          //- .cut-line
          //- .register-small-btn(@click="changeModule(2)") 免费注册
        //- .enterprise-wechat-btn(@click="qrCodeLogin")
        //-   .btn-wechat
        //-   .btn-text 企业微信登录
      //- 注册模块
      .register-box(v-show="show_module == 2 && !register_show")
        .register-top-box
          .title-left 
            .left-first 申请试用
            .left-line —
            .left-second 扫码安装
          .title-right(@click="changeModule(1)")
            .right-first 已有账号？
            .right-second 登录
        .register-center-box
          .center-first 申请试用需要企业微信管理员扫码安装
          .center-second 未注册企业微信的用户请先注册企业微信
        .wechat-btn
          .btn-wechat
          .btn-text 去企业微信扫码安装
          a.go-accredit(:href="qyInstallCodeUrl" target="_black")
        .wechat-btn-register
          .btn-wechat
          .btn-text 注册企业微信
          a.go-accredit(:href="qyInstallCodeUrl02" target="_black")
      //- //- 合作企业展示
      //- .cooperative-enterprise-box(v-show="(normal_login || show_module == 2) && !register_show")
      //-   .cooperative-title 他们都选择了畅移
      //-   .cooperative-tip 全面降低客户资产流失风险、管理社交电商流量、锁定超级用户，实现零售新增长！
      //-   .cooperative-img-box
      //-     img.item-img(:src="item.url" v-for="(item, index) in cooperative")
</template>
<script>
import { mapActions, mapGetters } from "vuex";
import Api from "@/api/public/Common";
import CreateKeys from "@/utils/ParseRole";
import CreateKeysNew from "@/utils/ParseRoleNew";
import Ciphertext from "@/utils/Ciphertext";
import { getCookie, setCookie, delCookie } from "@/utils/Cookie";
import AreaCode from "@/plugins/AreaCode";
import SlideBar from "@/components/SliderBar";
// import SelectEnterprise from "@/components/SelectEnterprise";
import RegisterEnterprise from "@/components/RegisterEnterprise";
import AccountOverdue from "@/components/AccountOverdue";
import AccountBind from "@/components/AccountBind";
import CheckNum from "@/mixins/CheckNum";
import JSEncrypt from 'jsencrypt'

const enterpriseMap = {
  'tapestry-web-uat.coach.com.cn': 3363603667487424,
  'tapestry-web-uat.stuartweitzman.cn': 3373738709328576,
  'tapestry-web-uat.katespade.cn': 3376048056335040,
  'we.coach.com.cn': 3360610327843968,
  'we.stuartweitzman.cn': 3374939108972352,
  'we.katespade.cn': 3383774945595776

}

export default {
  name: "login",
  components: {
    SlideBar,
    // SelectEnterprise,
    RegisterEnterprise,
    AccountOverdue,
    AccountBind
  },
  mixins: [CheckNum],
  data() {
    const checkMobile = (rule, value, callback) => {
      const reg = /^1[23456789]\d{9}$/;
      if (this.login_type === "code" && !value) {
        return callback(new Error("请输入手机号"));
      }
      if (this.login_type === "code" && !reg.test(value)) {
        return callback(new Error("手机号格式有误，请重新输入"));
      }
      callback();
    };
    const checkCode = (rule, value, callback) => {
      if (this.login_type === "code" && !value) {
        this.checkCodeError = true;
        return callback(new Error("请输入6验证码"));
      } else if (value.length < 6) {
        return callback(new Error("请输入6位验证码"));
      } else {
        this.checkCodeError = false;
      }
      callback();
    };
    const checAccountCode = (rule, value, callback) => {
      if (this.login_type !== "code" && !value) {
        // this.checkCodeError = true;
        return callback(new Error("请输入验证码"));
      } else {
        // this.checkCodeError = false;
      }
      callback();
    };
    const checNumberCode = (rule, value, callback) => {
      if (this.login_type === "code" && !value) {
        // this.checkCodeError = true;
        return callback(new Error("请输入验证码"));
      } else {
        // this.checkCodeError = false;
      }
      callback();
    };
    const checkAccount = (rule, value, callback) => {
      if (this.login_type === "password" && !value) {
        return callback(new Error("请输入登录账号"));
      }

      callback();
    };
    const checkOldPassword = (rule, value, callback) => {
      if (!value) {
        return callback(new Error("请输入登录密码"));
      }
      callback();
    };
    const checkPassword = (rule, value, callback) => {
      let reg = /^\S*(?=\S{8,})(?=\S*\d)(?=\S*[A-Z])(?=\S*[a-z])/;
      if (!value) {
        return callback(new Error("请输入登录密码"));
      }
      if (this.passInfo.originalPass === this.passInfo.newPass) {
        return callback(new Error("新密码与旧密码一致"));
      }
      if (!reg.test(value)) {
        return callback(
          new Error("密码至少包含8个字符，一个大写字母，一个小写字母和一个数字")
        );
      }
      callback();
    };
    const confirmPass = (rule, value, callback) => {
      let reg = /^\S*(?=\S{8,})(?=\S*\d)(?=\S*[A-Z])(?=\S*[a-z])/;
      if (!value) {
        return callback(new Error("请输入确认密码"));
      }
      if (this.passInfo.newPass !== value) {
        return callback(new Error("确认密码与新密码保持一致"));
      }
      if (!reg.test(value)) {
        return callback(
          new Error("密码至少包含8个字符，一个大写字母，一个小写字母和一个数字")
        );
      }
      callback();
    };
    return {
      labelName: "账号密码登录",
      info: {
        account: "", //账号
        password: "", //密码
        code: "", //验证码
        mobile: "", //手机号
        accountCode: "",
        enterprise_id: enterpriseMap[window.location.host], //企业标识
        numberCode: "" // 手机验证码
      },
      passInfo: {
        originalPass: "", // 原密码
        newPass: "", // 新密码
        confirmPass: "" // 确认密码
      },
      rule: {
        mobile: [{ validator: checkMobile, trigger: ["blur"] }],
        code: [{ validator: checkCode }],
        account: [{ validator: checkAccount }],
        password: [{ validator: checkOldPassword }],
        accountCode: [{ validator: checAccountCode }],
        numberCode: [{ validator: checNumberCode }]
      },
      passRule: {
        originalPass: [{ validator: checkOldPassword }],
        newPass: [{ validator: checkPassword }],
        confirmPass: [{ validator: confirmPass }]
      },
      show_module: 1, // 1： 显示登陆模块；2：显示注册模块
      login_show_module: 1, // 1： 账号登录；2：短信验证; 3: 新用户强制修改密码
      remember: false,
      areaJson: AreaCode,
      area: "+86",
      hash: "",
      token: "",
      login_type: "password",
      sevenChecked: false, // 七天记住
      is_show: false, // 是否展示密码
      is_code: false, // 是否发送验证码
      time: 60, // 5测试 正常60s
      valid_show: false,
      is_logining: false, // 控制登录按钮loading显隐
      is_show_select_ent: false,
      overdue_show: false, //试用版过期弹窗展示
      register_show: false, //注册试用企业的弹窗是否展示
      bind_show: false, //绑定已有账号弹窗是否展示
      normal_login: false, //正常的账号登录
      work_wechat_login_token: "",
      auth_enterprise: {},
      epId: "",
      qyInstallCodeUrl: "",
      qyInstallCodeUrl02: "",
      checkCodeError: false,
      corp_name: "",
      cooperative: [
        {
          url: require("@/assets/cooperative/logo_01.png"),
          opacity: "opacity-4"
        },
        {
          url: require("@/assets/cooperative/logo_02.png"),
          opacity: "opacity-5"
        },
        {
          url: require("@/assets/cooperative/logo_03.png"),
          opacity: "opacity-5"
        },
        {
          url: require("@/assets/cooperative/logo_04.png"),
          opacity: "opacity-4"
        },
        {
          url: require("@/assets/cooperative/logo_05.png"),
          opacity: "opacity-5"
        },
        {
          url: require("@/assets/cooperative/logo_06.png"),
          opacity: "opacity-35"
        },
        {
          url: require("@/assets/cooperative/logo_07.png"),
          opacity: "opacity-4"
        },
        {
          url: require("@/assets/cooperative/logo_08.png"),
          opacity: "opacity-4"
        },
        {
          url: require("@/assets/cooperative/logo_09.png"),
          opacity: "opacity-6"
        },
        {
          url: require("@/assets/cooperative/logo_10.png"),
          opacity: "opacity-55"
        },
        {
          url: require("@/assets/cooperative/logo_11.png"),
          opacity: "opacity-45"
        },
        {
          url: require("@/assets/cooperative/logo_12.png"),
          opacity: "opacity-65"
        },
        {
          url: require("@/assets/cooperative/logo_13.png"),
          opacity: "opacity-6"
        },
        {
          url: require("@/assets/cooperative/logo_14.png"),
          opacity: "opacity-55"
        },
        {
          url: require("@/assets/cooperative/logo_15.png"),
          opacity: "opacity-5"
        },
        {
          url: require("@/assets/cooperative/logo_16.png"),
          opacity: "opacity-65"
        }
      ],
      enterpriseList: [],
      tempEnterpriseList: [],
      newEnterpriseList: [],
      isFirstLogin: false, // 是否为首次登录
      verifyData: {
        code_key: "",
        base64_img: ""
      }
    };
  },
  computed: {
    ...mapGetters("MENU", {
      MENU_TOP: "getMenuTop"
    })
  },
  methods: {
    ...mapActions("USER", ["setUserInfo"]),
    ...mapActions("PERMISSION", ["setPermission", "setRoleList"]),
    // 点击logo区域，去官网
    goOfficial() {
      window.location.href = "https://www.vchangyi.com/";
    },
    // 登录注册模块切换
    changeModule(type) {
      this.show_module = type;
      this.checkCodeError = false;
      if (type === 2) {
        this.getWxInstallUrl();
      }
    },
    // 获取去企业微信扫码安装按钮二维码地址
    getWxInstallUrl() {
      const origin = window.location.origin;
      const back_url =
        origin === "https://www.vchangyi.com"
          ? `${origin}/login?installGuide=1`
          : `${origin}?installGuide=1`;
      Api.getWxInstallUrlApi({ back_url }).then(res => {
        this.qyInstallCodeUrl = res.url;
        this.qyInstallCodeUrl02 = res.qy_register_url;
      });
    },
    // 点击七天登录
    changeSeven(item) {
      this.sevenChecked = item;
      if (!item) {
        delCookie("autoAdmin");
      }
    },

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
    },

    // 登录
    goLogin() {
      if (this.newEnterpriseList.length > 1 && !this.info.enterprise_id) {
        this.$message.error("请先选择企业");
        return false;
      }
      this.checkCodeError = false;
      this.$refs["info"].validate(async valid => {
        if (!valid) {
          return false;
        } else {
          let apis = "login";
          let params = {};
          var tempAccount = "";
          if (this.info.account.indexOf("@") < 0) {
            // 邮箱登录
            tempAccount = `${this.area}-${this.info.account}`;
          } else {
            // 手机号登录
            tempAccount = this.info.account;
          }
          
          params = {
            login_type: 1,
            account: tempAccount,
            password: this.encryptPassword(this.$md5(this.info.password)),
            code: this.info.accountCode,
            code_key: this.verifyData.code_key,
            enterprise_id: this.info.enterprise_id || 0
          };
          this.is_logining = true;
          Api[apis](params).then(res => {
            if (res.code === 200) {
              this.info.mobile = res.data.token;
              this.isFirstLogin = res.data.is_first_login;
              this.login_show_module = 2;
              this.labelName = "短信验证登录";
              this.is_logining = false;
              this.sendCode();
            } else {
              this.is_logining = false;
              this.$message.error(res.message);
              this.getValidate();
            }
          });
        }
      });
    },
    // 短信验证
    testCode() {
      if (this.info.numberCode.length != 6) {
        this.$message.error("验证码格式不正确");
        return false;
      }
      let params = {
        login_type: 2,
        account: `${this.info.mobile}`,
        password: this.info.numberCode,
        enterprise_id: this.info.enterprise_id
      };
      this.is_logining = true;
      Api["login"](params).then(res => {
        if (res.code === 200) {
          this.is_logining = false;
          let tempData = res.data;
          this.hash = tempData.hash;
          this.token = tempData.token;
          localStorage.setItem("TOKEN", JSON.stringify(this.token));
          localStorage.setItem("LOGIN_HASH", JSON.stringify(this.hash));
          localStorage.setItem(
            "ENTERPRISE_NUM",
            JSON.stringify(res.data.enterprise_list.length)
          );
          // 存储当前登录企业标识
          /*  赫莲娜	=> HR	 羽西	=> YS	 理肤泉 =>	LRP	
              巴欧 =>	LRL	 美即 =>	MG
           */
          localStorage.setItem("ENTERPRISE_MARK", res.data.brand_code);
          // 判断是否为首次登录
          if (this.isFirstLogin) {
            this.login_show_module = 3;
            this.labelName = "首次登录，请修改密码";
            this.tempEnterpriseList = res.data.enterprise_list;
          } else {
            this.getEnterpriseList(res.data.enterprise_list);
          }
        } else {
          this.is_logining = false;
          this.$message.error(res.message);
        }
      });
    },
    // 修改密码并登录
    editLogin() {
      this.$refs["passInfo"].validate(valid => {
        if (!valid) {
          return false;
        } else {
          const params = {
            account: this.info.account,
            old_ps: this.$md5(this.passInfo.originalPass + ""),
            new_ps: this.$md5(this.passInfo.newPass + ""),
            new_ps_two: this.$md5(this.passInfo.confirmPass + ""),
            is_require_login: false,
            enterprise_id: this.info.enterprise_id || 0
          };
          this.is_logining = true;
          Api.editPassWord(params)
            .then(res => {
              if (res.code === 200) {
                this.is_logining = false;
                this.getEnterpriseList(this.tempEnterpriseList);
              }
            })
            .catch(err => {
              this.is_logining = false;
            });
        }
      });
    },
    // 获取企业列表
    getEnterpriseList(enterprise_list) {
      // 单企业登陆
      if (enterprise_list.length < 2) {
        this.getEnterpriseToken(enterprise_list[0]);
        // 多企业登陆
      } else {
        this.enterpriseList = enterprise_list;
        this.is_show_select_ent = true;
      }
    },
    //手动选择企业回调
    selectEntListenner(e) {
      this.getEnterpriseToken(e);
    },
    //通过token登陆所选企业
    getEnterpriseToken(target_enterprise) {
      localStorage.setItem("REAL_HASH", JSON.stringify(target_enterprise.hash)); //存储真实hash(非跳转hash)
      Api.getEnterpriseToken({
        admin_id: target_enterprise.admin_id,
        enterprise_id: target_enterprise.enterprise_id || this.info.enterprise_id
      }).then(
        res => {
          if (res.data.version == 1) {
            this.authorization(res.data);
          } else {
            this.authorizationNew(res.data);
          }
          // 登录成功，且用户勾选了7天自动登录，此时存储该账户
          if (this.sevenChecked) {
            setCookie(
              "autoAdmin",
              JSON.stringify({
                user: this.info.account,
                pass: this.info.password,
                userClickAuto: true
              }),
              7
            );
          }
          res.data;
        }
      );
    },
    // 权限解析处理于跳转V1
    authorization(res) {
      // console.log("V2=========" + JSON.stringify(res));
      if (res.data.serve_type === 3) {
        this.$message({
          showClose: false,
          message: "您的系统使用权限已到期，请联系管理员重新开通",
          type: "error"
        });
        return;
      }
      // 将用户信息存储至localstorage
      localStorage.setItem("LOGIN_HASH", JSON.stringify(res.hash));
      localStorage.setItem("TOKEN", JSON.stringify(res.token));
      localStorage.setItem("USERINFO", JSON.stringify(res.data));
      // boss key & role keys to crs keys
      CreateKeys(res.data).then(pass => {
        localStorage.setItem(
          "TEMP",
          Ciphertext("encode")("baseA", pass.join(","))
        );
        this.goHashPage(res.hash);
      });
    },
    // 权限解析处理于跳转V2
    authorizationNew(res) {
      // 试用过期提示
      let currentDate = parseInt(new Date().getTime()) / 1000;
      if (res.enterprise_admin_info.account_expire_at < currentDate) {
        this.overdue_show = true;
        return;
      }
      if (res.enterprise_admin_info.serve_type === 3) {
        this.$message({
          showClose: false,
          message: "您的系统使用权限已到期，请联系管理员重新开通",
          type: "error"
        });
        return;
      }
      // 将用户信息存储至localstorage
      const {
        hash,
        enterprise_admin_info,
        token,
        permissions_tree, //权限树
        role_list,
        permissions_list //权限列表
      } = res;
      localStorage.setItem("LOGIN_HASH", JSON.stringify(hash));
      localStorage.setItem(
        "ENTERPRISE_ID",
        enterprise_admin_info.enterprise_id
      );
      localStorage.setItem("TOKEN", JSON.stringify(token));
      localStorage.setItem("USERINFO", JSON.stringify(enterprise_admin_info));
      localStorage.setItem("KEYS", JSON.stringify(permissions_list));
      localStorage.setItem("PERMISSION", JSON.stringify(permissions_tree));
      localStorage.setItem("ROLES", JSON.stringify(role_list));

      // 暂时兼容新老权限
      if (permissions_tree.length === 0) {
        CreateKeysNew(res.enterprise_admin_info).then(pass => {
          localStorage.setItem(
            "TEMP",
            Ciphertext("encode")("baseA", pass.join(","))
          );
          this.goHashPage(res.hash);
        });
      } else {
        this.goHashPage(res.hash);
      }
    },
    // 根据hash跳转对应业务（容器）
    goHashPage(hash) {
      this.$message({
        type: "success",
        message: "登录成功!",
        duration: 1000
      });

      // 项目配置
      const projectConfig = [
        // 测试
        "b69395f18f85176a06336a960ef51541", // 德芙内购
        "edb03ea7a69a34c8f5b109ef17f2539c", // 德芙梦工厂
        // 线上
        "210004970da23b3b41aca96eba296800", // 德芙内购
        "fe1aa281fbe98b40254fb37d7d626ce1", // 德芙梦工厂
        "cd5fe71c9f76de2451cc3516f234ba42" // 欧舒丹
      ];
      setTimeout(() => {
        const redirect = this.$route.query.redirect || ''
        let rootPath = projectConfig.indexOf(hash) === -1 ? "web-admin" : hash;
        window.location.href = `${window.location.origin}/${rootPath}/${redirect}`;
      }, 1000);
    },
    //选择企 业
    selectEnterprise() {
      this.is_show_select_ent = true;
    },
    // 登录类型tab切换时，清空输入值
    loginChange() {
      this.checkCodeError = false;
      this.info.account = ""; //账号
      this.info.password = ""; //密码
      this.info.code = ""; //验证码
      this.info.mobile = ""; //手机号
      this.is_code = false; // 是否发送验证码flag置为初始值
    },
    // 点击眼睛icon,控制密码显隐
    seePass() {
      this.is_show = !this.is_show;
    },
    // 新的根据用户获取相关企业列表
    getNewEnterpriseList() {
      if(!this.info.account) return
      const params = {
        account: this.info.account
      };
      Api.getEnterpriseList(params).then(res => {
        this.newEnterpriseList = res.data.enterprise_admin_list;
        if (this.newEnterpriseList && this.newEnterpriseList.length === 1) {
          this.info.enterprise_id =
            this.newEnterpriseList[0].enterprise_id || 0;
        } else if (!this.newEnterpriseList.length) {
          this.info.enterprise_id = 0;
        }
      });
    },
    // 发送验证码
    sendCode() {
      let self = this;
      self.is_logining = false;
      let params = {
        mobile: `${this.info.mobile}`,
        type: 1,
        enterprise_id: this.info.enterprise_id
      };
      Api.sendCode(params).then(
        res => {
          if (res.code === 200) {
            self.is_code = true;
            self.valid_show = false;
            self.$message({
              type: "success",
              message: "验证码发送成功",
              duration: 1000
            });
            let interval = setInterval(() => {
              self.time--;
              if (self.time == 0) {
                // 停止定时器
                clearInterval(interval);
                self.is_code = false;
                self.time = 60;
                return;
              }
            }, 1000);
          } else {
            self.is_code = false;
            this.$message.error(res.message);
          }
        },
        () => {
          this.is_code = false;
        }
      );
    },
    // 滑块验证成功
    validSuccess() {
      let params = {
        mobile: `${this.info.mobile}`,
        type: 1
      };
      Api.sendCode(params).then(
        res => {
          if (res.code === 200) {
            let self = this;
            self.is_code = true;
            self.valid_show = false;
            self.$message({
              type: "success",
              message: "验证码发送成功",
              duration: 1000
            });
            let interval = setInterval(() => {
              self.time--;
              if (self.time == 0) {
                // 停止定时器
                clearInterval(interval);
                self.is_code = false;
                self.time = 60;
                return;
              }
            }, 1000);
          } else {
            self.is_code = false;
            this.$message.error(res.message);
          }
        },
        () => {
          this.is_code = false;
          this.valid_show = false;
        }
      );
    },
    validFail() {},
    closeValid() {
      this.valid_show = false;
    },
    // 点击企业微信登录按钮
    qrCodeLogin() {
      const origin = window.location.origin;
      const back_url =
        origin === "https://tapestry-dev.vchangyi.com"
          ? `${origin}/login`
          : origin;
      window.location.href = `${process.env.VUE_APP_API_ROOT}${process.env.VUE_APP_API_VERSION}gw-scrm/qy-code?ref=${back_url}`;
    },
    // 试用过期弹窗关闭按钮回调
    closeOverdue() {
      this.overdue_show = false;
      this.normal_login = true;
    },
    // 绑定已有账号弹窗确认按钮回调
    successBind(res) {
      this.bind_show = false;
      this.getEnterpriseToken(res.admin);
    },
    // 绑定已有账号弹窗取消按钮回调
    closeBind() {
      this.bind_show = false;
      this.normal_login = true;
    },
    // 绑定已有账号弹窗错误信息回调
    errorBind(err) {
      this.work_wechat_login_token = err.work_wechat_login_token;
    },
    // 申请试用填写信息组件成功回调
    successRegister(res) {
      this.getEnterpriseToken(res.admin);
    },
    // 申请试用填写信息组件关闭按钮回调
    closeRegister() {
      this.show_module = 1;
      this.normal_login = true; //走正常逻辑
      this.register_show = false;
    },
    // 企业微信登录接口
    wechatLogin() {
      const params = {
        type: this.wxtype,
        work_wechat_login_token: this.work_wechat_login_token
      };
      if (this.wxtype === "scan_code") {
        params.auth_code = this.auth_code;
      } else if (this.wxtype === "work_wechat") {
        params.corpid = this.corpid;
        params.token = this.wx_token;
      }
      Api.workWechatLogin(params)
        .then(res => {
          if (res.code === 200) {
            const {
              code,
              work_wechat_login_token,
              admin,
              tmp_token,
              auth_enterprise
            } = res.data;
            const map = {
              success: () => this.getEnterpriseToken(admin),
              enterprise_not_found: () => {
                this.auth_enterprise = auth_enterprise;
                this.register_show = true;
              },
              enterprise_expired: () => (this.overdue_show = true),
              admin_not_found: () => {
                this.epId = res.data.epId;
                this.bind_show = true;
              },
              auth_invalid: () => this.$message.error("授权无效")
            };

            localStorage.setItem("TOKEN", JSON.stringify(tmp_token));
            this.work_wechat_login_token = work_wechat_login_token;
            map[code]();
          } else {
            this.$message.error(res.message);
            this.normal_login = true;
          }
        })
        .catch(() => {
          this.normal_login = true;
        });
    },
    // 获取账号密码的验证码
    getValidate() {
      Api.validate()
        .then(res => {
          if (res.code === 200) {
            const { code_key, base64_img } = res.data;
            this.verifyData.code_key = code_key;
            this.verifyData.base64_img = base64_img;
          }
        })
        .catch(() => {
        });
    },
    /** 忘记密码 */
    handleForget() {
      // window.location.href = `${window.location.origin}/forget`;
      this.$router.push({
        path: '/login/forget'
      })
    }
  },
  mounted() {
    const { auth_code, corpid, token, accountId } = this.$route.query;
    if (
      this.$route.query &&
      this.$route.query.show_module &&
      parseInt(this.$route.query.show_module) === 2
    ) {
      this.show_module = 2;
      this.getWxInstallUrl();
    }
    // 扫码注册智慧导购回来直接去打开申请使用表单弹框
    if (this.$route.query && this.$route.query.installGuide) {
      Object.assign(this, {
        corpid,
        corp_name: this.$route.query.corp_name
      });
      this.register_show = true;
    } else if (auth_code) {
      // 扫码成功重定向而来
      Object.assign(this, {
        auth_code,
        corpid,
        wx_token: token,
        wxtype: "scan_code"
      });
      // 本地缓存记录登录方式
      localStorage.setItem("LOGIN_TYPE", "scan");
      this.wechatLogin();
    } else if (corpid && token) {
      // 企业微信后台跳转
      Object.assign(this, {
        auth_code,
        corpid,
        wx_token: token,
        wxtype: "work_wechat"
      });
      // 本地缓存记录登录方式
      localStorage.setItem("LOGIN_TYPE", "jump");
      this.wechatLogin();
    } else if (accountId) {
      // 通过邮件途径进来的
      this.info.account = accountId;
      this.normal_login = true;
    } else {
      // 本地缓存记录登录方式
      this.normal_login = true;
      localStorage.setItem("LOGIN_TYPE", "normal");
      // 用户勾选了七天自动登录
      if (this.sevenChecked && this.info.account && this.info.password) {
        this.goLogin("password");
      }
    }
    this.getValidate();
  },
  created() {
    delCookie("autoAdmin");
    let autoAdmin =
      (getCookie("autoAdmin") && JSON.parse(getCookie("autoAdmin"))) || null;
    if (autoAdmin && autoAdmin.user && autoAdmin.pass) {
      this.info = {
        account: autoAdmin.user, //账号
        password: autoAdmin.pass //密码
      };
      this.sevenChecked = true;
    }
  }
};
</script>

<style lang="sass" scoped>
.login
  width: 100vw
  height: 100vh
  min-width: 1200px
  min-height: 680px
  overflow: auto
  background: url('../../../assets/login-bac.png')
  background-repeat: no-repeat
  background-size: cover
  position: relative
  .links
    position: absolute
    bottom: 20px
    display: flex
    width: 1000px
    margin: 0 auto
    padding: 20px 0
    bottom: 0
    left: 50%
    transform: translateX(-50%)
    justify-content: space-between
    z-index: 99
    .a_box
      margin: 0 auto
      // display: flex
      text-align: center
      a
        display: flex
        align-items: center
        color: #909399
        text-decoration: none
    img
      width: 16px
      height: 16px
      padding-right: 5px
  /deep/.el-form-item__content
  .logo-click
    width: 200px
    height: 100px
    position: absolute
    top: 57px
    left: 60px
    cursor: pointer
    background: url('../../../assets/logo-black.png') no-repeat
    background-size: 108px 40px
  .copy-right
    color: #909399
    font-size: 16px
    position: absolute
    bottom: 20px
    left: 0
    right: 0
    margin: auto
    text-align: center
  .bg-mask
    width: 1333px
    height: 880px
    position: fixed
    left: 0
    top: 0
    bottom: 0
    right: 0
    margin: auto
    display: flex
.login-black-box
  // width: 1080px
  // height: 605px
  // position: fixed
  // left: 0
  // top: 0
  // bottom: 0
  // right: 0
  // margin: auto
  // display: flex
  display: flex
  -ms-flex-pack: justify
  justify-content: space-between
  position: relative
  width: 1080px
  height: 605px
  top: 50%
  left: 50%
  transform: translate(-50%,-50%)
  // 登录模块
  .login-box
    width: 520px
    height: 605px
    padding: 60px 60px
    border-radius: 10px
    background: #fff
    z-index: 50
    box-sizing: border-box
    margin: 0 auto
    // 顶部导航栏UI调整
    /deep/.el-tabs
      // padding-bottom: 30px
    /deep/.el-tabs__nav-wrap
      &:after
        height: 1px !important
        background-color: transparent
    /deep/.el-tabs__header
      margin: 0
      margin-bottom: 30px
    /deep/.el-tabs__item
      font-size: 20px
      color: #303133
      font-family: PingFangSC-Medium, PingFang SC
      font-weight: 500
    /deep/.el-tabs__item.is-active
      color: #000000
      font-size: 26px
    /deep/.el-tabs__active-bar
      height: 3px
      background-color: #000000
    /deep/.tab1-height
      height: 290px!important
    /deep/.tab2-height
      height: 290px!important
    /deep/.el-form
      margin-padding: 30px!important
    /deep/.mobile
      margin-bottom: 0 !important
    /deep/.el-form-item
      margin-bottom: 0px
    /deep/.el-form-item__error
      width: 350px
      color: #F87C7B
      line-height: 18px
    // 多选框UI调整
    /deep/.el-checkbox__input.is-checked .el-checkbox__inner
      background-color: #000000
      border-color: #000000
    /deep/.el-checkbox__inner:hover
      border-color: #000000
    /deep/.el-checkbox__input.is-checked+.el-checkbox__label
      color: #303133
    /deep/.el-checkbox__label
      padding-left: 3px
    // 单行文本输入框UI调整
    input
      text-indent: 0
      padding-left: 0 !important
      border-radius: 0px
      border-left: none
      border-right: none !important
      border-top: none
    /deep/.input-with-select .el-input-group__prepend
      background-color: #fff
    .el-input
      margin-top: 30px
      font-size: 16px
    /deep/.el-input__inner
      border: none
      border-bottom: solid 1px #E4E7EC
      border-radius: 0px
      padding-left: 0px
      color: #303133
    /deep/.el-input__inner:focus
      border-color: #000000
    /deep/.el-form-item.is-error .el-input__inner
      border-color: #F56C6C
    /deep/.el-form-item.is-success .el-input__inner:focus
      border-color: #000000
    /deep/.el-input__inner:-webkit-autofill
        -webkit-box-shadow:0 0 0 1000px white inset !important
    // 显示隐藏密码icon
    .see
      width: 19px
      height: 10px
      position: absolute
      bottom: 40px
      right: 0px
      cursor: pointer
      background-image: url('../../../assets/eye-on.png')
      background-size: 100% 100%
    .hide
      width: 19px
      height: 10px
      position: absolute
      bottom: 40px
      right: 0px
      cursor: pointer
      background-image: url('../../../assets/eye-off.png')
      background-size: 100% 100%
    // 验证码区域
    .code-row-box
      border-bottom: solid 1px #E4E7EC
      /deep/.el-input__inner
        border-bottom: none
        position: relative
      .code-box
        margin-left: 0px
        padding-left: 30px
        padding-top: 30px
        .code-btn
          border: none
          background: #fff
          color: #c0c4cc
          font-size: 16px
        .code-tip
          padding-left: 10px
        .hover-color:hover
          color: #000000
    .red-line
        border-bottom: solid 1px #F56C6C
    // 七天可登录
    .save-seven-box
      padding: 45px 0 20px 0
    // 登录提示
    .tip-box
      display: flex
      margin-bottom: 40px
      margin-top: 15px
      .tips
        color: #303133
        font-size: 14px
        font-weight: 400
        text-align: left
        line-height: 17px
      .cut-line
        width: 1px
        height: 14px
        background-color: #E4E7EC
        margin: 0 9px 0 9px
      .register-small-btn
        color: #000000
        font-size: 14px
        font-weight: 400
        text-align: left
        line-height: 17px
        cursor: pointer
    // 登录按钮
    .btn
      width: 100%
      color: #fff
      height: 46px
      font-size: 18px
      background-color: #000000
      border-color: #000000
      box-shadow: 0 11px 21px 0 rgba(0, 0, 0, 0.2)
    .btn:hover, .btn:focus
      color: #fff
      background-color: #000000
      border-color: #000000
    // 企业微信登录按钮
    .enterprise-wechat-btn
      height: 40px
      padding: 0
      font-size: 0
      border: 1px solid #1989FA
      background-color: #FFFFFF
      margin-bottom: 20px
      cursor: pointer
      border-radius: 4px
      .btn-wechat
        float: left
        width: 17px
        height: 16px
        margin: 12px 0 0 131px
        background-image: url('~@/assets/qywx.png')
        background-size: cover
      .btn-text
        float: left
        height: 20px
        line-height: 20px
        font-size: 18px
        color: #1989FA
        margin: 10px 0 0 10px
  // 注册模块
  .register-box
    width: 520px
    height: 605px
    padding: 60px 60px
    border-radius: 10px
    background: #fff
    z-index: 50
    box-sizing: border-box
    .register-top-box
      display: flex
      align-items: center
      justify-content: space-between
      .title-left
        display: flex
        align-items: center
        .left-first
          color: #303133
          font-size: 26px
          margin-right: 5px
          font-weight: 600
        .left-line
          color: #C0C4CC
          margin-right: 5px
        .left-second
          color: #C0C4CC
          font-size: 20px
      .title-right
        display: flex
        align-items: center
        .right-first
          color: #606266
          font-size: 20px
          margin-right: 3px
        .right-second
          color: #000000
          font-size: 20px
          cursor: pointer
    .register-center-box
      margin-top: 100px
      .center-first
        font-size: 22px
        color: #303133
        line-height: 30px
        text-align: center
      .center-second
        font-size: 18px
        color: #C0C4CC
        line-height: 25px
        text-align: center
        margin-top: 25px
    // 去企业微信扫码安装按钮
    .wechat-btn
      margin-top: 100px
      height: 46px
      padding: 0
      font-size: 0
      background-color: #0082EF
      margin-bottom: 45px
      cursor: pointer
      border-radius: 4px
      box-shadow: 0 11px 21px 0 rgba(98, 143, 255, 0.32)
      position: relative
      .btn-wechat
        float: left
        width: 22px
        height: 18px
        margin: 12px 0 0 103px
        background-image: url('~@/assets/qywx_white.png')
        background-size: cover
      .btn-text
        float: left
        height: 25px
        line-height: 25px
        font-size: 18px
        color: #FFFFFF
        font-weight: 500
        margin: 10px 0 0 10px
      .go-accredit
        position: absolute
        width: 100%
        height: 100%
        opacity: 0
        top: 0
        left: 0
    // 注册企业微信按钮
    .wechat-btn-register
      height: 46px
      padding: 0
      font-size: 0
      border: 1px solid #1989FA
      background-color: #FFFFFF
      cursor: pointer
      border-radius: 4px
      position: relative
      .btn-wechat
        float: left
        width: 22px
        height: 18px
        margin: 12px 0 0 130px
        background-image: url('~@/assets/qywx.png')
        background-size: cover
      .btn-text
        float: left
        height: 25px
        line-height: 25px
        font-size: 18px
        color: #0082EF
        font-weight: 500
        margin: 10px 0 0 10px
      .go-accredit
        position: absolute
        width: 100%
        height: 100%
        opacity: 0
        top: 0
        left: 0
  // 合作企业
  .cooperative-enterprise-box
    width: 549px
    height: 605px
    padding: 75px 50px 0 60px
    border-radius: 10px
    background: #fff
    z-index: 50
    box-sizing: border-box
    margin-left: 10px
    .cooperative-title
      font-size: 26px
      color: #303133
      text-align: center
      line-height: 37px
      font-weight: 600
    .cooperative-tip
      font-size: 22px
      color: #909399
      text-align: center
      line-height: 32px
      margin-top: 20px
    .cooperative-img-box
      margin-top: 25px
      .item-img
        width: 79px
        height: 38px
        display: inline-block
        margin-top: 45px
        margin-right: 38px
        background-position: center center
        background-repeat: no-repeat
        // filter: grayscale(1) opacity(0.5)
        background-size: 79px 38px
      .item-img:nth-child(4n+0)
        margin-right: 0
      .item-img:hover
        filter: none
      .opacity-35
        filter: grayscale(1) opacity(0.35)
      .opacity-4
        filter: grayscale(1) opacity(0.4)
      .opacity-45
        filter: grayscale(1) opacity(0.45)
      .opacity-5
        filter: grayscale(1) opacity(0.5)
      .opacity-55
        filter: grayscale(1) opacity(0.55)
      .opacity-6
        filter: grayscale(1) opacity(0.6)
      .opacity-65
        filter: grayscale(1) opacity(0.65)
  .code-image
    position: absolute
    right: 0
    bottom: 1px
    height: 45px
    width: 115px
    .image-validate
      width: 100%
      height: 100%
// .links
//   display: flex
//   width: 1000px
//   margin: 0 auto
//   padding: 20px 0
//   position: fixed
//   bottom: 0
//   left: 50%
//   transform: translateX(-50%)
//   justify-content: space-between
//   .a_box
//       display: flex
//       text-align: center
//     a
//       display: flex
//       align-items: center
//       color: white
//       text-decoration: none
//     img
//       width: 16px
//       height: 16px
//       padding-right: 5px
// @media screen and (max-width: 1680px)
//   .login-black-box, .logo-img, .bg-mask, .copy-right
//     transform: scale(0.94)
// @media screen and (max-width: 1440px)
//   .login-black-box, .logo-img, .bg-mask, .copy-right
//     transform: scale(0.8)
// @media screen and (max-width: 1280px)
//   .login-black-box, .logo-img, .bg-mask, .copy-right
//     transform: scale(0.7)
// @media screen and (max-width: 1165px)
//   .login-black-box, .logo-img, .bg-mask, .copy-right
//     transform: scale(0.6)
// @media screen and (max-width: 1024px)
//   .login-black-box, .logo-img, .bg-mask, .copy-right
//     transform: scale(0.5)
.send-vcode-btn
  position: absolute;
  top: 28px;
  right: 10px;
  color: #c0c4cc;
  font-weight: 600;
  cursor: pointer;
  z-index: 99;
.select-enterprise
  width: 100%;
  padding-top: 22px;
  /deep/ input
    font-size: 16px !important;
</style>
