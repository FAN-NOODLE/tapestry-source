<template>
  <div class="layout-list">
    <!-- 微信授权支付配置未授权弹窗 payment_config为2表示支付配置未授权 -->
    <div
      v-if="shop_auth_detail.installed && this.basicConfig.payment_config === 2"
      class="no-auth-tip"
    >
      <i
        class="el-icon-warning"
        aria-hidden="true"
      ></i>
      小程序已授权成功,请完善支付配置
      <el-button
        type="text"
        @click="goPayDeatil"
      >去完善</el-button>
    </div>
    <layout-card>
      <!-- 授权成功  -->
      <template v-if="shop_auth_detail.installed">
        <template slot="content">
          <!-- 标题 -->
          <div class="top-title m-b-20">完成微信授权，开通更多业务</div>
          <!-- 小程序暂未获取授权-->
          <div class="flex-between-c-b ">
            <!-- <template v-if="!shop_auth_detail.name"> -->
            <img
              :src="shop_auth_detail.logo ? shop_auth_detail.logo : wx_app_logo"
              alt=""
              class="m-20"
              width="50px"
            />
            <div class="center">
              <p class="desc1 m-b-5">
                <span
                  v-if="shop_auth_detail.name"
                  class="m-r-5"
                >{{
                  shop_auth_detail.name
                }}</span>
                <span
                  class="m-r-5"
                  v-else
                >小程序名称暂未获取</span>
                <el-button
                  type="success"
                  size="mini"
                  plain
                >已授权</el-button>
              </p>
              <p
                class="desc2"
                v-if="!shop_auth_detail.name && !shop_auth_detail.logo"
              >
                小程序信息尚未完善
                <el-button
                  type="text"
                  class="m-l-5"
                  @click="dialogVisible = true"
                >去完善</el-button>
              </p>
            </div>
            <!-- </template> -->
            <div>
              <el-button
                type="primary"
                size="small"
                class="font-14 m-27"
                @click="goDetail"
              >查看详情</el-button>
              <el-button
                size="small"
                class="font-14 m-r-20"
                @click="update"
              >更新信息</el-button>
              <span
                class="m-r-20"
                v-if="shop_auth_detail.app_secret"
              >
                <a
                  :href="shop_auth_detail.auth_url"
                  v-if="shop_auth_detail.app_secret"
                >
                  <el-button
                    v-if="shop_auth_detail.app_secret"
                    size="small"
                    class="size"
                  >
                    更新权限集
                  </el-button>
                </a>
              </span>
              <span class="m-r-20">
                <el-button @click="updateToken">刷新Token</el-button>
              </span>
            </div>
          </div>
        </template>
      </template>
      <!-- 第一次绑定 -->
      <template v-else>
        <template slot="content">
          <div class="top-title m-b-20">将微信小程序绑定店铺，提升交易量</div>
          <div class="flex-between-c">
            <div class="card-content flex-c-s m-r-15">
              <img
                src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/mall-set/bind_app.png"
                alt=""
                class="m-t-34"
              />
              <p class="desc">
                为保证所有功能的正常使用，
              </p>
              <p class="desc">授权时请把所有权限统一授权给畅鲟</p>
              <div class="posr">
                <a
                  :href="shop_auth_detail.auth_url"
                  class="go-accredit"
                ></a>
                <el-button
                  type="primary"
                  size="small"
                  class="font-14 m-b-75 m-t-30"
                  @click="gotoAuthDetail"
                >已有小程序，立即绑定</el-button>
              </div>
            </div>
            <div class="card-content flex-c-s m-l-15">
              <img
                src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/mall-set/register_app.png"
                alt=""
                class="m-t-34"
              />
              <p class="desc">
                请在微信公众平台申请新的小程序，
              </p>
              <p class="desc">成功后再到畅鲟进行绑定</p>
              <el-button
                type="primary"
                size="small"
                class="font-14 m-b-75 m-t-30"
                @click="registerApp"
              >没有小程序，立即申请</el-button>
            </div>
          </div>
        </template>
      </template>
    </layout-card>
    <!-- 弹窗内容 -->
    <el-dialog
      title="完善信息引导"
      :visible.sync="dialogVisible"
      width="100%"
    >
      <div class="flex-between-c">
        <div class="card-content flex-c-s m-r-20">
          <img
            src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/mall-set/official-user.png"
            alt=""
            class="m-t-30"
          />
          <p class="desc-top m-10 ">
            官网/公众号注册小程序用户
          </p>
          <p class="desc-first">请通过邮箱账号密码登录小程序</p>
          <p class="desc-first">后台完善信息</p>
          <el-button
            type="primary"
            size="small"
            class="font-14 m-t-40 m-b-30"
            @click="goAdmin()"
          >前往小程序后台</el-button>
        </div>
        <div class="card-content flex-c-s ">
          <img
            src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/mall-set/registration.png"
            alt=""
            class="m-t-30"
          />
          <p class="desc-top m-10 ">
            代注册小程序用户
          </p>
          <p class="desc-last">1.打开小程序助手小程序</p>
          <p class="desc-last">2.选择对应的小程序</p>
          <p class="desc-last">3.点击登录邮箱设置项</p>
          <el-button
            type="primary"
            size="small"
            class="font-14 m-b-30 m-t-20 last"
            @click="replaceRegister"
          >完善登录信息</el-button>
        </div>
      </div>
      <span
        slot="footer"
        class="dialog-footer"
      >
        <el-button @click="goHelp">前往帮助中心</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LayoutCard from "./components/layout-card/index.vue";
// 调取接口
import Api from "@/api/setting/setting";
import ConfigApi from "@/api/home/home";
@Component({
  components: {
    LayoutCard,
  },
})
export default class extends Vue {
  // 获取基础配置情况  payment_config代表支付配置 1已配置 2未配置
  private basicConfig = {};
  // 授权状态下小程序的数据
  private shop_auth_detail = {};
  //微信授权回调地址
  private back_url = "";
  private auth_url = "";
  private wx_app_logo = "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/mall-set/bind_app.png";
  // 弹窗默认状态
  private dialogVisible = false;
  // 前往帮助中心
  private goHelp() {
    const protocol = "http";
    window.open(
      `${protocol}://xueyuan.vchangyi.com/bin/view/11.%E8%AE%BE%E7%BD%AE/%E5%95%86%E5%9F%8E%E8%AE%BE%E7%BD%AE/`,
      "_blank"
    );
  }
  // 前往管理后台
  private goAdmin() {
    window.open("https://mp.weixin.qq.com", "_blank");
  }
  // 完善信息
  private goInfo() {
    window.open(
      "https://developers.weixin.qq.com/community/minihome/doc/000c4c6ffb43f809859a0da4456001?blockType=99",
      "_blank"
    );
  }
  // 没有小程序去往register-app
  private registerApp() {
    this.$router.push({
      name: "myapp_setting_authorization-settings_store-set_register-app",
    });
  }
  // 查看详情
  private goDetail() {
    this.$router.push({
      name: "myapp_setting_authorization-settings_store-set_register-detail",
    });
  }
  // 微信配置未完善 -------去完善
  private goPayDeatil() {
    this.$router.push({
      name: "myapp_setting_authorization-settings_wx-pay_pay-detail",
    });
  }
  // 获取商城小程序授权详情
  private getshopAuthDetail() {
    Api.shopAuthDetail({
      back_url: this.back_url,
    }).then((res) => {
      this.shop_auth_detail = res.data;
    });
  }
  //  调取更新信息接口
  private update() {
    Api.registerUpdataInfo({}).then(() => {
      this.$message({
        type: "success",
        message: "更新成功",
      });
    });
  }
  // 弹窗-------完善信息
  private replaceRegister() {
    window.open(
      "https://developers.weixin.qq.com/community/minihome/doc/000c4c6ffb43f809859a0da4456001?blockType=99",
      "_blank"
    );
  }

  // 调取支付配置检查接口
  private getNormalCheck() {
    ConfigApi.getNormalCheck({}).then((res) => {
      this.basicConfig = res.data;
    });
  }
  // 更新token
  private updateToken() {
    Api.refreshToken().then(() => {
      this.$message({
        message: "刷新成功",
        type: "success",
      });
    });
  }

  //绑定小程序
  private gotoAuthDetail() {
    // window.location.href = this.auth_url;
  }
  // dialog显示
  private isShowModal() {
    // 从审核完成路由传过来的参数
    const isShowModal = this.$route.query.isShowModal;
    // 传回来的参数为1--------弹窗展示
    if (isShowModal === "1") {
      this.dialogVisible = true;
    }
  }

  mounted() {
    const back_url = this.$router.resolve({
      path: "/myapp/setting/authorization-settings/store-set/empty",
    }).href;
    this.back_url = window.location.origin + back_url;
    this.isShowModal();
    // 获取支付配置检查接口
    this.getNormalCheck();
    // 获取商城小程序授权详情接口
    this.getshopAuthDetail();
  }
}
</script>
<style lang="scss" scoped>
div {
  user-select: none;
}

// 授权成功样式
.top-title {
  color: #303133;
  font-size: 16px;
  font-weight: 500;
}

// 第一次绑定样式
.flex-between-c {
  display: flex;
  // 主轴对齐方向：两端对齐
  justify-content: space-between;
  // 垂直方向居中对齐
  align-items: center;
}

// 框里面的内容样式
.flex-c-s {
  display: flex;
  // 水平方向居中对齐
  justify-content: center;
  // 正向列
  flex-direction: column;
  // 垂直居中对齐
  align-items: center;
}

.card-content {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  flex: 1;

  img {
    width: 120px;
    height: 120px;
  }

  .desc {
    color: #606266;
    font-size: 14px;
    text-align: center;
  }

  .m-b-75 {
    margin-bottom: 75px;
  }

  .el-button {
    width: 180px !important;
    height: 40px !important;
  }
}

.flex-between-c-b {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #f5f6fa;
  height: 90px;
  border-radius: 4px;

  //  中间文字按钮样式
  .center {
    position: absolute;
    left: 140px;

    // 第一行文字按钮样式
    .desc1 {
      color: #303133;
      font-size: 14px;
      text-align: left;
      font-weight: 500;

      .el-button {
        width: 46px;
        height: 22px !important;
        line-height: 0px !important;
        padding: 0px !important;
      }

      .el-button:hover,
      .el-button:focus {
        color: #6fc688;
        background: #f1f9f3;
        border-color: #c5e8cf;
      }
    }

    // 第二行文字按钮样式
    .desc2 {
      color: #909399;
      font-size: 12px;
      text-align: left;
      font-weight: 400;

      .el-button {
        font-size: 12px;
      }
    }
  }
}

// 弹窗样式
.el-dialog {
  .desc-top {
    font-weight: 400;
    color: #303133;
  }

  .desc-first {
    width: 196px;
    font-size: 12px;
    font-weight: 400;
    color: #909399;
    text-align: center;
  }

  .desc-last {
    width: 137px;
    font-size: 12px;
    font-weight: 400;
    color: #909399;
    text-align: left;
    margin-left: 10px;
  }

  .el-button {
    width: 180px !important;
    height: 40px !important;
  }
}

// 提示框样式
.no-auth-tip {
  margin-bottom: 30px;
  color: #606266;
  width: 100%;
  height: 50px;
  text-align: center;
  background: #fef7ec;
  border-radius: 4px;
  line-height: 50px;

  .el-icon-warning {
    color: #f6ab3c;
    width: 16px;
    height: 16px;
  }
}

// 间距样式
.m-b-20 {
  margin-bottom: 20px;
}

.m-t-34 {
  margin-top: 34px;
}

.m-20 {
  margin: 20px 0px 20px 20px;
}

.m-27 {
  margin: 27px 0px;
}

.m-b-15 {
  margin-bottom: 15px;
}

.m-10 {
  margin: 10px 0px;
}

.m-t-30 {
  margin-top: 30px;
}

.m-t-40 {
  margin-top: 40px;
}

.m-b-30 {
  margin-bottom: 30px;
}

.m-b-5 {
  margin-bottom: 5px;
}

.m-l-5 {
  margin-left: 5px;
}

.m-r-5 {
  margin-right: 5px;
}

.posr {
  position: relative;

  .go-accredit {
    position: absolute;
    top: 30px;
    left: 0;
    right: 0;
    bottom: 0;
    height: 40px;
  }
}

.size {
  font-size: 14px;
}
</style>
