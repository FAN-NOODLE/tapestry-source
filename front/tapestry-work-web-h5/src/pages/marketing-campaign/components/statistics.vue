<template>
  <div id="statistics">
    <div class="statistics-detail">
      <van-tabs
        class="tab-box"
        v-model="currentTab"
        animated
        sticky
        @click="changeTabClick"
      >
        <van-tab title="全部" :name="2">
          <div class="statistics-content">
            <van-list
              v-model="loading"
              :finished="finished"
              :finished-text="activeList.length ? '没有更多了' : ''"
              @load="onLoad"
              :immediate-check="false"
              offset="100"
            >
              <div
                class="single-statistics"
                v-for="(item, index) in activeList"
                :key="index"
                @click="toDetail(item)"
              >
                <div class="top-title">
                  <img :src="item.cover_url" alt="" />
                  <div class="title-text">
                    <h3>{{ item.title }}</h3>
                    <p>
                      累计报名：<span>{{ item.activity_user_count }}</span
                      >人
                    </p>
                  </div>
                  <span class="headline">
                    {{
                      item.type == 1 ? "报名" : item.type == 2 ? "签到" : "问卷"
                    }}
                  </span>
                </div>
                <p class="activity-time">
                  <span class="left-time">
                    活动有效期:{{
                      $helper.dateFormat(item.start_time, "yyyy/MM/dd hh:mm")
                    }}-{{ $helper.dateFormat(item.end_time, "MM/dd hh:mm") }}
                  </span>
                  <span class="right-status">{{ item.status }}</span>
                </p>
              </div>
            </van-list>
            <!-- 列表缺省-->
            <empty-page
              v-if="activeList.length == 0 && !loading"
              picture="noData"
              paddingTop="100"
              paddingBottom="0"
              text="暂无活动"
            ></empty-page>
          </div>
        </van-tab>
        <van-tab title="进行中" :name="1">
          <div class="statistics-content">
            <van-list
              v-model="loading"
              :finished="finished"
              :finished-text="activeList.length ? '没有更多了' : ''"
              @load="onLoad"
              :immediate-check="false"
              offset="100"
            >
              <div
                class="single-statistics"
                v-for="(item, index) in activeList"
                :key="index"
                @click="toDetail(item)"
              >
                <div class="top-title">
                  <img :src="item.cover_url" alt="" />
                  <div class="title-text">
                    <h3>{{ item.title }}</h3>
                    <p>
                      累计报名：<span>{{ item.activity_user_count }}</span
                      >人
                    </p>
                  </div>
                  <span class="headline">
                    {{
                      item.type == 1 ? "报名" : item.type == 2 ? "签到" : "问卷"
                    }}
                  </span>
                </div>
                <p class="activity-time">
                  <span class="left-time">
                    活动有效期:{{
                      $helper.dateFormat(item.start_time, "yyyy/MM/dd hh:mm")
                    }}-{{ $helper.dateFormat(item.end_time, "MM/dd hh:mm") }}
                  </span>
                  <span class="right-status">{{ item.status }}</span>
                </p>
              </div>
            </van-list>
            <!-- 列表缺省-->
            <empty-page
              v-if="activeList.length == 0 && !loading"
              picture="noData"
              paddingTop="100"
              paddingBottom="0"
              text="暂无活动"
            ></empty-page>
          </div>
        </van-tab>
        <van-tab title="已结束" :name="3">
          <div class="statistics-content">
            <van-list
              v-model="loading"
              :finished="finished"
              :finished-text="activeList.length ? '没有更多了' : ''"
              @load="onLoad"
              :immediate-check="false"
              offset="100"
            >
              <div
                class="single-statistics"
                v-for="(item, index) in activeList"
                :key="index"
                @click="toDetail(item)"
              >
                <div class="top-title">
                  <img :src="item.cover_url" alt="" />
                  <div class="title-text">
                    <h3>{{ item.title }}</h3>
                    <p>
                      累计报名：<span>{{ item.activity_user_count }}</span
                      >人
                    </p>
                  </div>
                  <span class="headline">
                    {{
                      item.type == 1 ? "报名" : item.type == 2 ? "签到" : "问卷"
                    }}
                  </span>
                </div>
                <p class="activity-time">
                  <span class="left-time">
                    活动有效期:{{
                      $helper.dateFormat(item.start_time, "yyyy/MM/dd hh:mm")
                    }}-{{ $helper.dateFormat(item.end_time, "MM/dd hh:mm") }}
                  </span>
                  <span class="right-status">{{ item.status }}</span>
                </p>
              </div>
            </van-list>
            <!-- 列表缺省-->
            <empty-page
              v-if="activeList.length == 0 && !loading"
              picture="noData"
              paddingTop="100"
              paddingBottom="0"
              text="暂无活动"
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
import { activityStatisticsApi } from "@/api/marketing-campaign/index";
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
  private currentTab = 2;
  private activeList = [];
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
    this.getStatisticsList();
  }
  // 触底加载
  private onLoad() {
    if (this.isRequest) {
      if (this.total > this.activeList.length) {
        this.page = this.page + 1;
        this.getStatisticsList();
      } else {
        this.finished = true;
      }
    }
  }
  // tab切换
  private changeTabClick() {
    this.activeList = [];
    this.page = 1;
    this.finished = false;
    this.loading = false;
    this.getStatisticsList();
  }
  async getStatisticsList(): Promise<void> {
    let params: any = {
      page: this.page,
      size: 10,
    };
    if (this.currentTab !== 2) {
      params.status = this.currentTab;
    }
    await activityStatisticsApi(params).then((res: any) => {
      if ((res.httpCode as number) == 200) {
        this.activeList = [...this.activeList, ...res.data.list];
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
  private toDetail(item) {
    const obj = {
      id: item.id,
      type: 4,
    };
    this.$emit("changeId", obj);
    this.$emit("changeStatus", item.status);
    this.$emit("saveTab", this.currentTab);
  }
}
</script>
<style lang="scss" scoped>
#statistics {
  background: #ffffff;
  min-height: 100vh;
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
  .statistics-content {
    padding: 15px 15px 80px 15px;
    box-sizing: border-box;
    .single-statistics {
      padding: 15px;
      background: #ffffff;
      box-shadow: 0px 9px 16px 0px rgba(0, 0, 0, 0.05);
      border-radius: 5px;
      margin-bottom: 15px;
      .top-title {
        display: flex;
        margin-bottom: 10px;
        position: relative;
        img {
          width: 84px;
          height: 84px;
          border-radius: 6px;
          margin-right: 10px;
        }
        .title-text {
          h3 {
            font-size: 16px;
            font-weight: bold;
            color: #33363c;
            line-height: 26px;
            text-shadow: 0px 9px 16px rgba(0, 0, 0, 0.05);
          }
          p {
            font-size: 12px;
            color: #808080;
            line-height: 20px;
            text-shadow: 0px 9px 16px rgba(0, 0, 0, 0.05);
            span {
              color: #33363c;
              font-weight: bold;
              margin-right: 1px;
            }
          }
        }
        .headline {
          position: absolute;
          top: 0px;
          left: 0px;
          width: 44px;
          height: 24px;
          display: flex;
          justify-content: center;
          align-items: center;
          background: #caa361;
          box-shadow: 0px 9px 16px 0px rgba(0, 0, 0, 0.05),
            0px 2px 8px 0px rgba(0, 0, 0, 0.09);
          border-radius: 0px 0px 12px 12px;
          font-size: 12px;
          color: #ffffff;
          line-height: 16px;
          text-shadow: 0px 9px 16px rgba(0, 0, 0, 0.05);
        }
      }
      .activity-time {
        padding: 5px 7px;
        background: rgba(202, 163, 97, 0.1);
        box-shadow: 0px 9px 16px 0px rgba(0, 0, 0, 0.05);
        border-radius: 5px;
        font-size: 12px;
        display: flex;
        align-items: center;
        justify-content: space-between;
        .left-time {
          color: #caa361;
          text-shadow: 0px 9px 16px rgba(0, 0, 0, 0.05);
        }
        .right-status {
          font-weight: bold;
          color: #af8b4e;
          margin-left: 5px;
        }
      }
    }
  }
}
</style>
