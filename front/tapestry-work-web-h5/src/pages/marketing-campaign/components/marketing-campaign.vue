<template>
  <div id="activity">
    <div class="activity-detail">
      <van-tabs
        class="tab-box"
        v-model="currentTab"
        animated
        sticky
        @click="changeTabClick"
      >
        <van-tab title="报名" :name="1">
          <div class="activity-content">
            <van-list
              v-model="loading"
              :finished="finished"
              :finished-text="activityList.length ? '没有更多了' : ''"
              :immediate-check="false"
              @load="onLoad"
              offset="100"
            >
              <div
                class="single-activity"
                v-for="(item, index) in activityList"
                :key="index"
                @click="toDetail(item)"
              >
                <div class="activity-img">
                  <img :src="item.cover_url" alt="" />
                  <span>进行中</span>
                </div>
                <p>{{ item.title }}</p>
              </div>
            </van-list>
            <!-- 列表缺省-->
            <empty-page
              v-if="activityList.length == 0 && !loading"
              picture="noData"
              paddingTop="100"
              paddingBottom="0"
              text="暂无报名活动"
            ></empty-page>
          </div>
        </van-tab>
        <van-tab title="问卷" :name="3">
          <div class="activity-content">
            <van-list
              v-model="loading"
              :finished="finished"
              :finished-text="activityList.length ? '没有更多了' : ''"
              :immediate-check="false"
              @load="onLoad"
              offset="100"
            >
              <div
                class="single-activity"
                v-for="(item, index) in activityList"
                :key="index"
                @click="toDetail(item)"
              >
                <div class="activity-img">
                  <img :src="item.cover_url" alt="" />
                  <span>进行中</span>
                </div>
                <p>{{ item.title }}</p>
              </div>
            </van-list>
            <!-- 列表缺省-->
            <empty-page
              v-if="activityList.length == 0 && !loading"
              picture="noData"
              paddingTop="100"
              paddingBottom="0"
              text="暂无问卷活动"
            ></empty-page>
          </div>
        </van-tab>
        <van-tab title="签到" :name="2">
          <div class="activity-content">
            <van-list
              v-model="loading"
              :finished="finished"
              :finished-text="activityList.length ? '没有更多了' : ''"
              :immediate-check="false"
              @load="onLoad"
              offset="100"
            >
              <div
                class="single-activity"
                v-for="(item, index) in activityList"
                :key="index"
                @click="toDetail(item)"
              >
                <div class="activity-img">
                  <img :src="item.cover_url" alt="" />
                  <span>进行中</span>
                </div>
                <p>{{ item.title }}</p>
              </div>
            </van-list>
            <!-- 列表缺省-->
            <empty-page
              v-if="activityList.length == 0 && !loading"
              picture="noData"
              paddingTop="100"
              paddingBottom="0"
              text="暂无签到活动"
            ></empty-page>
          </div>
        </van-tab>
      </van-tabs>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import { Tab, Tabs, List } from "vant";
import emptyPage from "@/components/base/empty-page.vue";
import { activityListApi } from "@/api/marketing-campaign/index";
@Component({
  components: {
    emptyPage,
    [Tab.name]: Tab,
    [Tabs.name]: Tabs,
    [List.name]: List,
  },
})
export default class extends Vue {
  @Prop({ default: 0, type: Number }) defaultTab!: number;
  private currentTab = 1;
  private activityList = [];
  private page = 1;
  private loading = false;
  private finished = false;
  private total = 0;
  private isRequest = false;

  private created() {
    window.scrollTo(0, 0);
  }
  private mounted() {
    if (this.defaultTab !== 0) {
      this.currentTab = this.defaultTab;
    }
    this.getActivityList();
  }
  // 触底加载
  private onLoad() {
    if (this.isRequest) {
      if (this.total > this.activityList.length) {
        this.page = this.page + 1;
        this.getActivityList();
      } else {
        this.finished = true;
      }
    }
  }
  // tab切换
  private changeTabClick() {
    this.activityList = [];
    this.page = 1;
    this.finished = false;
    this.loading = false;
    this.getActivityList();
  }
  private getActivityList() {
    activityListApi({
      type: this.currentTab,
      page: this.page,
      size: 10,
    }).then((res: any) => {
      if ((res.httpCode as number) == 200) {
        this.activityList = [...this.activityList, ...res.data.list];
        this.total = res.data.meta.total;
        // 加载状态结束
        this.loading = false;
        this.isRequest = true;
      } else {
        this.$toast.fail(res.status);
      }
    });
  }
  // 跳转详情
  private toDetail(item: any) {
    const obj = {
      id: item.id,
      type: 3,
    };
    this.$emit("changeId", obj);
    this.$emit("saveTab", this.currentTab);
  }
}
</script>
<style lang="scss" scoped>
#activity {
  min-height: 100vh;
  background: #ffffff;
  ::v-deep {
    .van-tabs__wrap {
      height: 51px;
    }
    .van-tab {
      color: #33363c;
      font-size: 16px;
    }
    .van-tabs__line {
      width: 16px;
      height: 2px;
      border-radius: 2px;
      background-color: #af8b4e;
    }
    .van-tab--active {
      color: #af8b4e;
    }
    .van-tabs__wrap {
      // position: sticky;
      // top: -2px;
      width: 100%;
      background: #ffffff;
      z-index: 99;
      border-bottom: 1px solid #f0f0f0;
    }
  }
  .activity-detail {
    .activity-content {
      padding: 15px 15px 80px 15px;
      box-sizing: border-box;
      .single-activity {
        padding: 15px 15px 0px 15px;
        .activity-img {
          position: relative;
          img {
            height: 158px;
            border-radius: 12px;
          }
          span {
            position: absolute;
            top: 15px;
            right: 15px;
            background: rgba(0, 0, 0, 0.4);
            box-shadow: 0px 2px 8px 0px rgba(0, 0, 0, 0.09);
            border-radius: 15px;
            padding: 7px 10px;
            font-size: 12px;
            color: #ffffff;
            line-height: 16px;
          }
        }
        p {
          text-align: center;
          font-size: 14px;
          color: #33363c;
          line-height: 20px;
          text-shadow: 0px 9px 16px rgba(0, 0, 0, 0.05);
        }
        border-radius: 12px;
        margin-bottom: 15px;
        padding-bottom: 10px;
        box-shadow: 0px 9px 16px 0px rgba(0, 0, 0, 0.05);
      }
    }
  }
}
</style>
