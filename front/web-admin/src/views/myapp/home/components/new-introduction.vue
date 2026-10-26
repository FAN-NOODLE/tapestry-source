<template>
  <div class="new-introduce">
    <div class="tabs">
      <div
        v-for="(item, index) in editableTabs"
        :key="item.name"
        class="normal-tab"
        :class="tabIndex === index ? 'active-bg' : ''"
        @click="tabClick(index)"
      >
        {{ item.title }}
      </div>
    </div>
    <div class="line"></div>
    <!-- 基础配置 -->
    <div class="content" v-show="tabIndex === editableContent.normal">
      <div class="m-b-20">
        <div class="title">
          <span>小程序授权</span>
          <span v-if="baseConfig.shop_auth === disposeState.on">
            <i
              aria-hidden="true"
              class="el-icon-success success-color m-l-10"
            ></i>
            <span class="m-l-2">已完成</span>
          </span>
          <span v-else>
            <i
              aria-hidden="true"
              class="el-icon-warning-outline success-color m-l-10"
            ></i>
            <span class="m-l-2">未完成</span>
          </span>
        </div>
        <div class="m-b-10">
          如果你已注册小程序使用微信小程序管理员帐号扫码进行授权，授权过程中请勾选所有权限以确保小程序功能完整性
        </div>
        <div>
          如果你未注册小程序，可以直接前往<span class="theme-color m-l-10"
            ><a
              href="https://mp.weixin.qq.com/"
              target="_blank"
              rel="noopener noreferrer"
              >微信公众平台</a
            ></span
          >
        </div>
      </div>
      <div class="m-b-20">
        <div class="title">
          <span>支付配置</span>
          <span v-if="baseConfig.payment_config === disposeState.on">
            <i
              aria-hidden="true"
              class="el-icon-success m-l-10 success-color"
            ></i>
            <span class="m-l-2">已完成</span>
          </span>
          <span v-else>
            <i
              aria-hidden="true"
              class="el-icon-warning-outline m-l-10 undo-color"
            ></i>
            <span class="m-l-2">未完成</span>
          </span>
        </div>
        <div
          class="m-b-10"
          v-if="baseConfig.payment_config === disposeState.off"
        >
          如果你已开通微信支付特约商户并且完成小程序授权后，可直接进行支付配置。<span
            class="theme-color m-l-10"
            @click="goConfig"
            >去配置</span
          >
        </div>
        <div>
          如果你未开通特约客户，请开通微信支付特约客户。开通完成可联系服务商相关人员进行配置后以完成后续支付配置。
        </div>
      </div>
      <div class="m-b-20">
        <div class="title">
          <span>企微配置</span>
          <span v-if="baseConfig.qw_config === disposeState.on">
            <i
              aria-hidden="true"
              class="el-icon-success m-l-10 success-color"
            ></i>
            <span class="m-l-2">已完成</span>
          </span>
          <span v-else>
            <i
              aria-hidden="true"
              class="el-icon-warning-outline m-l-10 undo-color"
            ></i>
            <span class="m-l-2">未完成</span>
          </span>
        </div>
        <div v-if="baseConfig.qw_config === disposeState.off">
          使用“企业微信”功能，需前往企微配置相关参数。<span
            class="theme-color m-l-10"
            @click="goAuth"
            >去配置</span
          >
        </div>
      </div>
    </div>
    <!-- 店铺装修 -->
    <div class="content" v-show="tabIndex === editableContent.shopdesc">
      <div class="m-b-20">
        <div class="title">
          <span>店铺风格</span>
        </div>
        <div class="m-b-10">
          店铺风格指商城小程序端的风格色系，小程序首次发布之前需要先设置店铺风格。
        </div>
      </div>
      <div class="m-b-20">
        <div class="title">
          <span>店铺导航</span>
        </div>
        <div class="m-b-10">
          店铺导航是指商城小程序底部的导航图标，小程序首次发布之前你需要设置它。
        </div>
      </div>
      <div class="m-b-20">
        <div class="title">
          <span>店铺模版</span>
        </div>
        <div class="m-b-10">
          店铺模版是指商城小程序端的首页和活动页的相关内容配置，品牌日常运营的重点模块。
        </div>
      </div>
    </div>
    <!-- 商品上新 -->
    <div class="content" v-show="tabIndex === editableContent.goodnew">
      <div class="m-b-20">
        <div class="title">
          <span>库存管理</span>
        </div>
        <div class="m-b-20">
          库存管理主要管理物料信息，即货的信息，对接第三方系统，分配店铺物料，管理企业综合库存等。
        </div>
      </div>
      <div class="m-b-20">
        <div class="title">
          <span>邮费模版</span>
        </div>
        <div class="m-b-10">
          邮费模板用于商品邮费设置。可设置各商品的配送范围以及邮费条件。
        </div>
      </div>
      <div class="m-b-20">
        <div class="title">
          <span>规格模版</span>
        </div>
        <div class="m-b-10">
          规格模板可用于同类商品，系统管理商品规格。为顾客选择商品提供商品说明。
        </div>
      </div>
      <div class="m-b-20">
        <div class="title">
          <span>商品管理</span>
        </div>
        <div class="m-b-10">
          商品管理主要管理商品展示、售卖价格、售卖库存、商品邮费等信息。
        </div>
      </div>
    </div>
    <!-- 版本发布 -->
    <div class="content" v-show="tabIndex === editableContent.edition">
      <div class="m-b-20">
        <div class="title">
          <span>提审管理</span>
        </div>
        <div class="m-b-10">
          选择需要发布的页面模版 ，统一提交给微信官方进行审核。
        </div>
      </div>
      <div class="m-b-20">
        <div class="title">
          <span>发布管理</span>
        </div>
        <div class="m-b-10">
          审核通过后，在版本管理进行发布，发布后商城小程序即可展示最新版本。
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import ShareModule from "@/shared/index";
@Component
export default class NewIntroduction extends Vue {
  @Prop() baseConfig!: object;
  private tabPosition = "left";
  private editableTabs = [
    {
      title: "1.基础配置",
      name: 1
    },
    {
      title: "2.店铺装修",
      name: 2
    },
    {
      title: "3.商品上新",
      name: 3
    },
    {
      title: "4.版本发布",
      name: 4
    }
  ];
  private tabIndex = 0;
  private editableContent = {
    normal: 0, //基础配置
    shopdesc: 1, //店铺装修
    goodnew: 2, //商品上新
    edition: 3 //版本发布
  };
  private tabClick(index: number): void {
    this.tabIndex = index;
  }
  private goAuth() {
    this.$router.push({
      name: "myapp_setting_authorization-settings_enterprise-config"
    });
  }
  private goConfig() {
    this.$router.push({
      name: "myapp_setting_authorization-settings_wx-pay_pay-detail"
    });
  }
  //配置状态
  private disposeState = {
    on: 1, //已配置
    off: 2 //未配置
  };
}
</script>
<style lang="scss" scoped>
.new-introduce {
  color: #606266;
  display: flex;
  justify-content: flex-start;
  height: 100%;
  .tabs {
    display: flex;
    flex-direction: column;
    justify-content: flex-start;
    cursor: pointer;
    .normal-tab {
      padding: 20px;
      box-sizing: border-box;
      color: #303133;
      font-size: 14px;
    }
    .active-bg {
      background: $color-background;
      border-radius: 4px;
    }
  }
  .line {
    width: 1px;
    background: $color-third-line;
    margin: 0 20px;
  }
  .content {
    .title {
      margin-bottom: 10px;
      span {
        font-size: 14px;
        color: $color-text-second;
      }
      span:first-child {
        color: $color-text-main;
        font-weight: bolder;
      }
    }
  }
}

.theme-color {
  color: $color-main;
  a {
    color: $color-main;
    text-decoration: none;
  }
}
.success-color {
  color: #6fc688;
}
.undo-color {
  color: #f5f6fa;
}
.m-l-2 {
  margin-left: 2px;
}
</style>
