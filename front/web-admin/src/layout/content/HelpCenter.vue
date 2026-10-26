<template>
  <div class="help-center">
    <!-- 公告弹窗 -->
    <announcement-model :announcement="announcement"></announcement-model>
    <!-- 帮助中心内容 -->
    <div class="header-title clearfix">
      <div class="title-jump ">
        <el-link
          target="_blank"
          :underline="false"
          class="title pull-left"
        >
          <span>帮助中心</span>
          <i class="el-icon-arrow-right title-icon" aria-hidden="true"></i>
        </el-link>
      </div>
      <div class="hide pull-right" @click="hideHelp">
        <span>收起</span>
        <div class="hide-icon"></div>
        <!-- <svg-icon icon-class="help-hide" class="hide-icon"></svg-icon> -->
      </div>
    </div>
    <div class="function-introduce clearfix" v-if="HELPCONTENT.type === 'text'">
      <h4 class="title">功能介绍</h4>
      <p class="content">
        {{ HELPCONTENT.content }}
      </p>
    </div>
    <!-- 公告列表组件 -->
    <announcement-list
      v-if="HELPCONTENT.type === 'list'"
      :announceList="HELPCONTENT.content"
    ></announcement-list>
    <!-- 在线客服 -->
    <!-- <div
      class="customer-service"
      :class="{ 'm-b-20': HELPCONTENT.key === 'myapp_home' }"
    >
      <h4 class="title">在线客服</h4>
      <div class="code-wrap clearfix">
        <el-popover
          class="sun-code"
          placement="left-start"
          trigger="hover"
          width="127"
        >
          <div class="code pull-left" slot="reference"></div>
          <img
            src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/help-center/qrcode.png"
            alt=""
            class="big-code"
          />
        </el-popover>

        <div class="service-text pull-left">
          添加官方客服微信1对1专业免费咨询
        </div>
      </div>
    </div> -->
    <!-- 公众号 (只有在首页时才会出现) -->
    <div class="the-public" v-if="HELPCONTENT.key === 'myapp_home'">
      <h4 class="title">畅移公众号</h4>
      <div class="code-wrap clearfix">
        <el-popover
          class="sun-code"
          placement="left-start"
          trigger="hover"
          width="127"
        >
          <div class="code pull-left" slot="reference"></div>
          <img
            src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/help-center/changyi.png"
            alt=""
            class="big-code"
          />
        </el-popover>
        <div class="service-text pull-left">
          扫码关注畅移公众号了解更多行业最新资讯
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
import HelpStore from "@/store/modules/help";
import AnnouncementModel from "../components/announcement-model.vue";
import AnnouncementList from "../components/announcement-list.vue";
@Component({
  components: {
    AnnouncementModel, //查看公告弹窗组件
    AnnouncementList //公告列表组件
  }
})
export default class HelpCenter extends Vue {
  get HELPCONTENT() {
    return HelpStore.getHelpContent;
  }
  private modelShow = true; //公告弹窗是否展示
  private announcement = {}; //要查看的公告
  private wiki =
    "http://xueyuan.vchangyi.com/bin/view/CRS%E5%8A%9F%E8%83%BD%E8%AF%B4%E6%98%8E/";
  // 收起帮助中心
  private hideHelp() {
    HelpStore.setHelp("2");
  }
}
</script>
<style lang="scss" scoped>
.help-center {
  width: 100%;
  box-sizing: border-box;
  padding: 20px;
  /deep/.el-link--default {
    color: #6f3112;
  }
  .header-title {
    border-bottom: 1px solid $color-third-line;
    .title-jump {
      cursor: pointer;
      opacity: 0;
      &:hover {
        .title,
        .title-icon {
          color: $color-main;
        }
      }
      .title {
        font-size: 16px;
        line-height: 26px;
        color: $color-text-main;
        font-weight: 500;
        margin-bottom: 20px;
      }
      .title-icon {
        font-weight: 600;
        font-size: 14px;
        color: $color-text-main;
      }
    }
    .hide {
      color:#6f3112;
      font-size: 14px;
      line-height: 20px;
      margin-top: 2px;
      font-weight: 400;
      cursor: pointer;
      &:hover {
        color: $color-main;
        .hide-icon {
          background-image: url("https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/help-center/hide-arrow-dark.png");
        }
      }
      span {
        vertical-align: middle;
      }
      .hide-icon {
        width: 9px;
        height: 9px;
        margin-left: 5px;
        display: inline-block;
        background-image: url("https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/help-center/hide-arrow-dark.png");
        background-size: cover;
        vertical-align: middle;
      }
    }
  }
  .function-introduce,
  .customer-service,
  .the-public {
    .title {
      margin: 20px 0 5px 0;
      font-size: 16px;
      color: $color-text-main;
      font-weight: 500;
      line-height: 26px;
      text-align: left;
    }
  }
  .function-introduce {
    // border-bottom: 1px solid $color-third-line;
    .content {
      color: $color-text-normal;
      line-height: 20px;
      font-size: 14px;
      // margin-bottom: 25px;
      text-align: left;
    }
  }
  .customer-service {
    .code-wrap {
      .code {
        width: 50px;
        height: 50px;
        background-image: url("https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/help-center/qrcode.png");
        background-size: cover;
      }
      .service-text {
        width: 116px;
        color: $color-text-normal;
        font-size: 14px;
        line-height: 20px;
        margin: 6px 0 0 10px;
        text-align: left;
      }
    }
  }
  .the-public {
    border-top: 1px solid $color-third-line;
    .code-wrap {
      .code {
        width: 50px;
        height: 50px;
        background-image: url("https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/help-center/changyi.png");
        background-size: cover;
      }
      .service-text {
        width: 140px;
        color: $color-text-normal;
        font-size: 14px;
        line-height: 20px;
        margin: 6px 0 0 10px;
        text-align: left;
      }
    }
  }
}
.big-code {
  width: 150px;
  height: 150px;
}
</style>
