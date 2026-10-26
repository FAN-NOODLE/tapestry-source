<template>
  <div id="public-page">
    <van-nav-bar
      :title="type == 6 ? '未入群用户' : tit"
      left-arrow
      @click-left="$router.back()"
    />
    <div style="padding-top: 45px">
      <div class="loading" v-show="loadingShow">
        <van-loading type="spinner" color="#ffffff" vertical
          >加载中...</van-loading
        >
      </div>
      <!-- 群成员明细 -->
      <van-tabs
        v-if="type"
        title-active-color="#000000"
        title-inactive-color="#000000"
        v-model="currentTab"
        border
        sticky
        ref="tabs"
        @click="changeTab"
      >
        <!-- 全部\好友\非好友 -->
        <van-tab
          v-for="(item, index) in tabsList"
          :key="index"
          :title="item.tit"
          :name="index + 1"
        >
          <div class="list-box">
            <div class="activites-tit">
              {{ tit
              }}{{
                type !== 6 && currentTab - 1
                  ? "-" + tabsList[currentTab - 1].tit
                  : ""
              }}：<span class="activites-num">共{{ meta.total }}人</span>
            </div>
            <van-list
              v-if="list.length > 0"
              v-model="loading"
              offset="30"
              :finished="noMore"
              finished-text="没有更多了~"
              @load="getMoreOrder"
            >
              <userItem
                :ownerId="$route.query.ownerId"
                :list="list"
                :total="meta.total"
                :type="type"
              />
            </van-list>
          </div>
          <empty-page
            v-if="!list.length && !loadingShow"
            picture="emptyList"
            paddingTop="87"
            paddingBottom="60"
            :text="
              type == 7 || type == 8
                ? '暂无' + tit
                : type == 6
                ? '暂无未入群用户'
                : type == 5
                ? '暂无参与用户'
                : '暂无群成员'
            "
          ></empty-page>
        </van-tab>
      </van-tabs>
      <userItem
        v-else
        :ownerId="$route.query.ownerId"
        :list="list"
        :type="type"
      />
    </div>
  </div>
</template>
<script lang="ts">
import { NavBar } from "vant";
import { Component, Vue, Prop } from "vue-property-decorator";
import userItem from "../components/userItem.vue";
import emptyPage from "@/components/base/empty-page.vue";
import { chatUserList, chatGroupMembers } from "@/api/customer-portrait";
import wxUtils from "@/utils/wx";
@Component({
  components: {
    [NavBar.name]: NavBar,
    userItem,
    emptyPage,
  },
})
export default class extends Vue {
  private list = [];
  private loading = false;
  private loadingShow = true;
  private currentTab = 1; // 默认选中的tab
  private noMore = false;
  private meta = {
    page: 1,
    size: 8,
    total: 0,
  };
  private type = 0;
  private tit = "";
  private tabsList = [{ tit: "全部" }, { tit: "好友" }, { tit: "非好友" }];
  public get brand(): string {
    return this.$store.getters["user/getEnterpriseMark"];
  }
  private async created() {
    if (
      process.env.NODE_ENV !== "test" &&
      process.env.NODE_ENV !== "development"
    ) {
      await wxUtils.wxAgentConfig();
    }
    if (this.brand !== "HR") document.body.classList.add("other");

    this.type = Number(this.$route.query.type);
    if (this.type == 4) document.title = "群成员列表";
    else if (this.type == 5) document.title = "活动参与排行榜";
    else if (this.type == 6) document.title = "未入群用户";
    else if (this.type == 7) document.title = "会员用户";
    else if (this.type == 8) document.title = "非会员用户";
    this.tit = document.title;
    if (this.type == 6) this.tit = "白金卡及以下卡等未入群的管理用户";
    if (this.type == 4 || this.type == 5) this.getMemberList();
    else this.getUserList();
  }
  // 点击tab切换
  private changeTab(name) {
    this.loadingShow = true;
    this.list = [];
    this.meta.page = 1;
    this.meta.total = 0;
    this.currentTab = name;
    if (this.type == 4 || this.type == 5) this.getMemberList();
    else this.getUserList();
  }
  // 用户列表接口
  private async getUserList() {
    this.loading = true;
    const params = {
      chat_id: this.$store.getters["common/getGroupChatId"],
      page: this.meta.page,
      size: this.meta.size,
      list_type: this.type == 6 ? 0 : this.type == 7 ? 1 : 2,
      friend: this.currentTab - 1,
      // ets: Number(this.ets)
      ets: Number(this.$route.query.ets),
    };
    const { code, message = "", data }: IApiResult = await chatUserList(params);
    if (code === 200) {
      this.loadingShow = false;
      if (data.items && data.items.length) {
        data.items.forEach((item: any) => {
          item.join_time =
            item.join_time &&
            this.$helper.dateFormat(item.join_time, "yyyy-MM-dd hh:mm");
        });
        this.list = this.list.concat(data.items);
        this.meta.total = data.total_count;
        this.loading = false;
        this.meta.page++;
        this.noMore = this.list.length >= data.total_count;
        if (
          data.total_count == this.list.length ||
          !data.items ||
          (data.items && data.items.length == 0)
        ) {
          this.noMore = true;
        }
      }
    } else {
      this.loading = false;
      this.noMore = true;
      this.loadingShow = false;
      this.$toast.fail(message);
    }
  }
  // 群成员列表
  private async getMemberList() {
    this.loading = true;
    const params = {
      chat_id: this.$store.getters["common/getGroupChatId"],
      page: this.meta.page,
      size: this.meta.size,
      friend: this.currentTab - 1,
      list_type: this.type == 4 ? 0 : 1,
    };
    const {
      code,
      message = "",
      data,
    }: IApiResult = await chatGroupMembers(params);
    if (code === 200) {
      this.loading = false;
      this.loadingShow = false;
      if (data.items && data.items.length) {
        this.list = this.list.concat(data.items);
        this.meta.total = data.total_count;
        this.loading = false;
        this.meta.page++;
        this.noMore = this.list.length >= data.total_count;
        if (
          data.total_count == this.list.length ||
          !data.items ||
          (data.items && data.items.length == 0)
        ) {
          this.noMore = true;
        }
      }
    } else {
      this.loading = false;
      this.noMore = true;
      this.loadingShow = false;
      this.$toast.fail(message);
    }
  }
  // 页面滑到底部触发加载更多
  private getMoreOrder() {
    // 活动排行榜仅展示Top10,不做滚动加载
    if (this.type == 4 || this.type == 5) this.getMemberList();
    else this.getUserList();
  }
}
</script>
<style lang="scss" scoped>
#public-page {
  background: #f4f7fb;
  ::v-deep {
    .van-nav-bar .van-icon {
      color: #000000;
    }
    .van-nav-bar__title {
      font-weight: bold;
      color: #000000;
      font-size: 17px;
    }
    .van-nav-bar__content {
      position: fixed;
      top: 0px;
      width: 100%;
      background: #ffffff;
    }
    .van-sticky--fixed {
      top: auto;
    }
  }
  .activites-tit {
    height: 22px;
    font-size: 14px;
    font-weight: 400;
    color: #808080;
    line-height: 22px;
    margin: 15px 15px 0;
    position: relative;
    .activites-num {
      font-weight: 600;
      color: #000000;
      position: absolute;
      right: 0;
    }
  }
  .loading {
    position: absolute;
    top: 0;
    bottom: 0;
    right: 0;
    left: 0;
    margin: auto;
    z-index: 1;
    .van-loading {
      position: absolute;
      top: 0;
      left: 0;
      bottom: 0;
      right: 0;
      margin: auto;
      background: rgba(0, 0, 0, 0.5);
      width: 140px;
      height: 100px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
      .van-loading__text {
        color: #ffffff;
      }
    }
  }
}
</style>

<style lang="scss">
.van-tabs__line {
  background-color: #000000;
}
</style>
