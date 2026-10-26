<template>
  <div id="activity">
    <marketing-campaign
      v-if="pageType === 1"
      @changeId="changeId"
      @saveTab="saveTab"
      :defaultTab="defaultTab"
    ></marketing-campaign>
    <statistics
      v-if="pageType === 2"
      @changeId="changeId"
      @saveTab="saveTab"
      :defaultTab="defaultTab"
      @changeStatus="changeStatus"
    ></statistics>
    <activity-page
      v-if="pageType === 3"
      @changeType="changeType"
      :id="id"
      :defaultTab="defaultTab"
    ></activity-page>
    <statistics-details
      v-if="pageType === 4"
      @changeType="changeType"
      :id="id"
      :defaultTab="defaultTab"
      :status="status"
    ></statistics-details>

    <div class="signPage-bottom" v-if="pageType === 1 || pageType === 2">
      <div
        class="bottom-tab"
        :style="{ color: pageType === 2 ? '#a4a5aa' : '#7c889f' }"
        @click="toSkip(1)"
      >
        <img
          v-if="pageType === 1"
          src="@/assets/marketing-campaign/home-active.png"
          alt=""
        />
        <img
          v-if="pageType === 2"
          src="@/assets/marketing-campaign/home.png"
          alt=""
        />
        <span>首页</span>
      </div>
      <div
        class="bottom-tab"
        :style="{ color: pageType === 1 ? '#a4a5aa' : '#7c889f' }"
        @click="toSkip(2)"
      >
        <img
          v-if="pageType === 1"
          src="@/assets/marketing-campaign/statistics.png"
          alt=""
        />
        <img
          v-if="pageType === 2"
          src="@/assets/marketing-campaign/statistics-active.png"
          alt=""
        />
        <span>统计</span>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import marketingCampaign from "./components/marketing-campaign.vue";
import activityPage from "./components/activity-page.vue";
import statisticsDetails from "./components/statistics-details.vue";
import statistics from "./components/statistics.vue";
import wx from "weixin-js-sdk";
@Component({
  components: {
    marketingCampaign,
    activityPage,
    statisticsDetails,
    statistics,
  },
})
export default class extends Vue {
  private pageType = 1; // 1首页 2统计 3活动详情 4统计详情
  private defaultTab = 0;
  private id = 0;
  private title = "";
  private status = "";
  // 统计页面
  private toSkip(type: number) {
    if (type === 1) {
      this.defaultTab = 1;
    } else {
      this.defaultTab = 2;
    }
    this.pageType = type;
  }
  // 传递id
  private changeId(item: any) {
    this.id = item.id;
    this.pageType = item.type;
  }
  // 返回保持tab选项
  private saveTab(val: number) {
    this.defaultTab = val;
  }
  // 返回的状态
  private changeStatus(item: any) {
    this.status = item;
  }
  // 返回按钮
  private changeType(item: any) {
    this.pageType = item.type;
    if (item.type === 1 || item.type === 2) {
      this.defaultTab = item.defaultTab;
    }
  }
  private popstate() {
    wx.closeWindow(); // 微信网页退出
  }
  private async created() {
    // 当前需要直接退出的页面
    window.history.pushState("_android_back", "", location.href);
    window.addEventListener("popstate", this.popstate, false);
  }
  private async destroyed() {
    // 记得清除，不然单页面，其它页面也可以用
    window.removeEventListener("popstate", this.popstate, false);
  }
}
</script>
<style lang="scss" scoped>
.signPage-bottom {
  width: 100%;
  position: fixed;
  bottom: -2px;
  background: #fff;
  padding: 5px 0px 15px;
  font-size: 14px;
  font-weight: 400;
  line-height: 20px;
  display: flex;
  justify-content: space-between;
  box-shadow: 0px 2px 9px 0px rgba(0, 0, 0, 0.31);
  .bottom-tab {
    display: flex;
    flex-direction: column;
    width: 50%;
    justify-content: center;
    align-items: center;
    font-size: 10px;
    font-weight: 500;
    color: #7c889f;
    line-height: 14px;
    margin: 10px 0px;
    img {
      width: 23px;
      height: 23px;
      margin-bottom: 3px;
    }
  }
}
</style>
