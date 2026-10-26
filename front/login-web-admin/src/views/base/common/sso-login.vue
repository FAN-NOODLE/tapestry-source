<template>
  <div class="sso-login">
    <div class="login-box">
      <div class="login-head">Weclient管理后台</div>
      <div class="login-content">
        <!-- <div class="title">Tapestry</div> -->
        <div class="btn" @click="handleLogin">使用Tapestry SSO 登录</div>
        <div class="tips">- Tapestry -</div>
      </div>
    </div>
  </div>
</template>

<script>
const enterpriseMap = {
  // coach
  "tapestry-web-uat.coach.com.cn": 3363603667487424,
  "we.coach.com.cn": 3360610327843968,
  // ks
  "tapestry-web-uat.katespade.cn": 3376048056335040,
  "we.katespade.cn": 3383774945595776
  // sw
  // "tapestry-web-uat.stuartweitzman.cn": 3373738709328576,
  // "we.stuartweitzman.cn": 3374939108972352
};
const appIdMap = {
  // weclient Coach QA
  "we.coach.com.cn": "e8e595a5-8f5d-45f6-be56-3547ed10754f",
  // weclient Coach Prod
  "tapestry-web-uat.coach.com.cn": "030b24ab-8fbb-4a01-8f32-3f791807ba56",
  // weclient Kate Spade QA
  "we.katespade.cn": "d5c3898b-33e9-40d2-a111-2186c3078894",
  // weclient Kate Spade Prod
  "tapestry-web-uat.katespade.cn": "ea970c83-9e69-4933-b78e-e0bff6fbeb53"
};

export default {
  name: "SsoLogin",
  data() {
    return {
      enterpriseId: enterpriseMap[window.location.host],
      appId:
        appIdMap[window.location.host] || "e8e595a5-8f5d-45f6-be56-3547ed10754f"
    };
  },
  mounted() {},
  methods: {
    async handleLogin() {
      // 未登录
      window.location.href = `https://login.microsoftonline.com/74da624f-e27a-4b03-8702-f5e0a389e009/oauth2/v2.0/authorize?client_id=${
        this.appId
      }&response_type=code&response_mode=query&scope=User.Read&state=${Date.now()}`;
    }
  }
};
</script>
<style lang="scss" scoped>
.sso-login {
  width: 100vw;
  height: 100vh;
  background: url("~@/assets/bg.jpg") center no-repeat;
  background-size: cover;
  position: relative;
  &::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    z-index: 0;
  }

  .login-box {
    width: 450px;
    height: 300px;
    border-radius: 10px;
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    overflow: hidden;
    background-color: #fff;
    z-index: 1;

    .login-head {
      height: 60px;
      background-color: #fed64b;
      font-weight: bold;
      font-size: 20px;
      line-height: 60px;
      text-align: center;
    }

    .login-content {
      padding: 40px 40px;
      box-sizing: border-box;
    }

    .title {
      font-size: 16px;
      color: #666;
    }

    .btn {
      height: 50px;
      line-height: 50px;
      text-align: center;
      margin-top: 40px;
      background-color: #f6f7f8;

      cursor: pointer;
      border-radius: 4px;
    }

    .tips {
      position: absolute;
      bottom: 20px;
      left: 0;
      right: 0;
      text-align: center;
      font-size: 14px;
      color: #666;
    }
  }
}
</style>
