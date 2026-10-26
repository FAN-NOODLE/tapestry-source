<template>
  <el-popover
    popper-class="logo-hover"
    trigger="hover"
    width="313"
    v-model="modelShow"
  >
    <div
      v-if="USER.logo"
      class="brand-logo pull-left m-r-20 m-t-10"
      :style="{ backgroundImage: `url(${USER.logo})` }"
      slot="reference"
    ></div>
    <div
      v-else
      class="default-logo pull-left m-r-20 m-t-10"
      slot="reference"
    ></div>
    <div class="message-wrap">
      <div class="message-header">
        <!-- 展示的logo(大) -->
        <div
          v-if="USER.qr_code"
          class="qrcode"
          :style="{ backgroundImage: `url(${USER.qr_code})` }"
        ></div>
        <div
          v-else
          class="header-default-logo"
        ></div>
        <div class="header-message">
          <div class="user-name">{{ USER.account_name }}</div>
          <div class="role-tags">
            <template v-if="USER.is_super != 1">
              <el-tag
                class="tags m-r-10"
                v-for="(tag, index) in ROLES"
                :key="index"
                size="mini"
              >{{ tag.name }}</el-tag>
            </template>
            <template v-else>
              <el-tag
                size="mini"
                type="warning"
              >超管</el-tag>
            </template>
          </div>
          <div class="enterprise-name m-t-10">{{ USER.shop_name }}</div>
        </div>
      </div>
      <div class="message-body">
        <!-- <div
          class="body-item"
          v-if="enterpriseNum > 0"
          @click="selectEnterprise"
        >
          <svg-icon
            class="m-l-20 icon"
            icon-class="enterprise-change"
          ></svg-icon>
          <span class="item m-l-10">切换企业</span>
        </div> -->
        <div
          class="body-item"
          @click="goChangePassword"
        >
          <svg-icon
            class="m-l-20 icon"
            icon-class="account"
          ></svg-icon>
          <span class="item m-l-10">修改密码</span>
        </div>
        <div
          class="body-item"
          @click="logout"
        >
          <svg-icon
            class="m-l-20 icon"
            icon-class="quit"
          ></svg-icon>
          <span class="item m-l-10">退出</span>
        </div>
      </div>
    </div>
  </el-popover>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import UserModule from "@/store/modules/user";
import PerModule from "@/store/modules/permission";
import LoginAboat from "@/utils/LoginManager";
@Component
export default class LogoHover extends Vue {
  get USER() {
    return UserModule.getUserInfo;
  }
  get ROLES() {
    return PerModule.getRoleList;
  }
  private enterpriseNum = 0;
  private modelShow = false; //弹窗的控制

  /**
   * @description 退出登录
   * */
  private async logout() {
    await this.$confirm("确定退出当前账户, 是否继续?", "温馨提示", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning"
    });
    this.$message({
      type: "success",
      message: "退出成功!"
    });
    localStorage.clear();
    sessionStorage.clear();
    this.$emit("logout");
    LoginAboat.backLogin();
  }
  private selectEnterprise() {
    this.$emit("selectEnterprise");
  }
  private goChangePassword() {
    this.$router.push({
      name: "myapp_setting_accounts_accounts-news_change-password"
    })
  }
  mounted() {
    // 本地登录不存在切换企业功能
    this.enterpriseNum = JSON.parse(
      localStorage.getItem("ENTERPRISE_NUM") || "0"
    );
  }
}
</script>
<style lang="scss" scoped>
.round {
  border-radius: 50%;
  background-size: cover;
  cursor: pointer;
}

.small-logo {
  width: 30px;
  height: 30px;
}

.default-logo {
  @extend .small-logo;
  @extend .round;
  background-color: $color-main;
  text-align: center;
  color: #ffffff;
  line-height: 30px;
  background-size: cover;
  background-image: url("https://crs-prod-1251064102.cos.ap-shanghai.myqcloud.com/tapestry-wxapp/default-avatar-gray.png");
}

.brand-logo {
  @extend .small-logo;
  @extend .round;
}

.message-wrap {
  .message-header {
    width: 100%;
    height: 132px;
    padding-top: 25px;
    border-bottom: 1px solid #e4e7ed;
    box-sizing: border-box;

    .big-logo {
      width: 90px;
      height: 90px;
      float: left;
      margin-left: 20px;
    }

    .qrcode {
      @extend .round;
      @extend .big-logo;
    }

    .header-default-logo {
      @extend .round;
      @extend .big-logo;
      background-color: $color-main;
      text-align: center;
      color: #ffffff;
      background-image: url("https://crs-prod-1251064102.cos.ap-shanghai.myqcloud.com/tapestry-wxapp/default-avatar-gray.png");
      line-height: 90px;
      font-size: 20px;
    }

    .header-message {
      width: 170px;
      float: left;
      margin-left: 13px;

      .user-name {
        color: $color-text-main;
        font-size: 16px;
        line-height: 24px;
        font-weight: 500;
        margin-bottom: 3px;
      }

      .role-tags {
        width: 200px;
        height: 20px;
        overflow: hidden;

        .tags {
          max-width: 185px;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
      }

      .enterprise-name {
        width: 170px;
        color: $color-text-second;
        font-size: 14px;
        line-height: 24px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }
  }

  .message-body {
    padding-top: 15px;

    .body-item {
      width: 100%;
      height: 54px;
      line-height: 54px;
      cursor: pointer;

      &:hover {
        background-color: $color-main;

        .item {
          color: #ffffff;
        }
      }

      .icon {
        vertical-align: middle;
      }

      .item {
        color: $color-text-normal;
        font-size: 16px;
        font-weight: 500;
        vertical-align: middle;
      }
    }
  }
}
</style>
