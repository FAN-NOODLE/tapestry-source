<template>
  <div
    :class="[info.chat_info ? 'portrait-association' : '', 'group-portrait']"
  >
    <div v-if="info.chat_info">
      <!-- 群信息 -->
      <div class="box-item">
        <div class="group-name">{{ info.chat_info.name || "-" }}</div>
        <div class="group-admin">
          <div class="tag">
            <span>群主</span>
          </div>
          {{ info.chat_info.owner || "-" }}
        </div>
        <div class="create-time">
          建群时间 {{ info.chat_info.create_time || "-" }}
        </div>
        <!-- <div class="base-info" :class="{ 'base-hide': isInfoHide }">
          <div class="info-item">
            <div class="label">群运营者</div>
            <div class="val">{{ info.chat_info.operator || "-" }}</div>
          </div>
          <div class="info-item">
            <div class="label">所属渠道</div>
            <div class="val">{{ info.chat_info.channel_name || "-" }}</div>
          </div>
          <div class="info-item">
            <div class="label">群类型</div>
            <div class="val">{{ info.chat_info.chat_type_name || "-" }}</div>
          </div>
          <div class="info-item">
            <div class="label">群目标</div>
            <div class="val">
              {{ info.chat_info.chat_target || "-" }}
            </div>
          </div>
        </div> -->
        <!-- <div class="hide-shadow" v-if="isInfoHide"></div>
        <div class="open-btn" @click="infoTrigger">
          <span>{{ isInfoHide ? "展开" : "收起" }}</span>
          <van-icon :name="isInfoHide ? 'arrow-down' : 'arrow-up'" />
        </div> -->
      </div>
      <!-- 时间组件 -->
      <div v-if="false" class="portrait-card portrait-middle portrait-date">
        <div class="date-tit">时间:</div>
        <div style="display: flex">
          <van-cell
            :class="{ 'change-cal': show_et }"
            :value="date_st"
            @click="openCal(0)"
          />
          <div class="portrait-cell">至</div>
          <van-cell
            :class="{ 'change-cal': show_et2 }"
            :value="date_et"
            @click="openCal(1)"
          />
          <van-calendar
            v-model="show_st"
            :color="brand === 'HR' ? '#D5BA8B' : '#5B89FE'"
            type="range"
            :show-confirm="false"
            :min-date="minDate"
            :formatter="formatter"
            :default-date="default_date"
            :allow-same-day="true"
            @confirm="onConfirm"
            @close="closeCal"
          />
        </div>
      </div>
      <!-- 数据统计 -->
      <!-- 客户概述 -->
      <div class="box-item">
        <div class="title">客户概述</div>
        <ul class="summary-ul">
          <li
            class="summary-li"
            v-for="(item, index) in c_summary_list"
            :key="index"
          >
            <p class="num">{{ info.target[item.prop] || 0 }}</p>
            <p class="num-tit">
              <span>{{ item.name }}</span>
              <van-popover
                trigger="click"
                v-model="item.showPop"
                :actions="item.actions"
              >
                <template #reference>
                  <van-icon name="warning-o" />
                </template>
              </van-popover>
            </p>
          </li>
        </ul>
      </div>
      <!-- 群会员统计 -->
      <div class="box-item">
        <div class="title">群会员统计</div>
        <ul class="g_summary-ul">
          <li
            class="g_summary-li"
            v-for="(item, index) in g_summary_list"
            :key="index"
          >
            <!-- <div v-if="index < 2" class="top-right">
              <van-icon name="arrow" color="#ffffff" />
            </div> -->
            <p class="num">{{ info.target[item.prop] || 0 }}</p>
            <p class="num-tit">
              <span>{{ item.name }}</span>
              <van-popover
                trigger="click"
                v-model="item.showPop"
                :actions="item.actions"
              >
                <template #reference>
                  <van-icon name="warning-o" />
                </template>
              </van-popover>
            </p>
          </li>
        </ul>
        <div ref="chart" :style="{ height: '382px', width: '100%' }"></div>
      </div>

      <!-- 会员等级 -->
      <div v-if="false" class="portrait-card portrait-middle pie">
        <div class="header">
          <div class="title"><i class="icon" /><i class="icon" />会员等级</div>
          <van-dropdown-menu
            :overlay="false"
            class="options"
            active-color="#B5B5B5"
          >
            <van-dropdown-item
              v-model="selectPieOpt"
              title-class="option"
              get-container=".options"
              :options="[
                { text: '新入群用户', value: 'new' },
                { text: '留存用户', value: 'valid' },
              ]"
              @open="setBorderRadius"
              @close="setBorderRadius('remove')"
              @change="changeOption"
            />
          </van-dropdown-menu>
        </div>
        <div ref="chart" :style="{ height: '382px', width: '100%' }"></div>
      </div>
      <!-- 群成员明细\活动参与排行榜 -->
      <div class="box-item rank-list">
        <van-tabs
          title-active-color="#000000"
          title-inactive-color="#000000"
          line-width="30px"
          v-model="currentTab"
          @change="changeTab"
        >
          <van-tab title="群成员列表"></van-tab>
          <van-tab title="活动参与排行榜"></van-tab>
        </van-tabs>

        <div class="list-box" v-if="list[tabsList[currentTab].prop].length > 0">
          <van-list>
            <userItem
              :ownerId="info.chat_info.user_id"
              :list="list[tabsList[currentTab].prop]"
              :type="!currentTab ? 0 : 2"
            />
            <div class="open-btn" @click="getPage(!currentTab ? 0 : 5)">
              <span>查看更多</span>
              <van-icon name="arrow" />
            </div>
          </van-list>
        </div>
        <div v-else style="border-radius: 0 0 12px 12px; background: #ffffff">
          <empty-page
            picture="emptyList"
            paddingTop="70"
            paddingBottom="60"
            :text="!currentTab ? '暂无成员' : '暂无用户参与活动'"
          ></empty-page>
        </div>
      </div>
    </div>
    <div v-else>
      <empty-page
        picture="noInfo"
        paddingTop="70"
        paddingBottom="60"
        text=""
      ></empty-page>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop, Provide } from "vue-property-decorator";
import emptyPage from "@/components/base/empty-page.vue";
import userItem from "./components/userItem.vue";
import {
  chatGroupMembers,
  chatIndexList,
  getRefreshTime,
} from "@/api/customer-portrait";
import { Icon, Tab, Tabs } from "vant";
const echarts = require("./echarts.min.js");
@Component({
  components: {
    emptyPage,
    userItem,
    [Icon.name]: Icon,
    [Tab.name]: Tab,
    [Tabs.name]: Tabs,
  },
})
export default class extends Vue {
  @Prop({ default: 0, type: String }) groupChatId!: string;

  /*************************** data start ***************************/
  private refreshTime = ""; // 数据更新时间
  private selectPieOpt = "valid";
  private currentTab = 0; // 默认选中的tab
  private loading = false;
  private noMore = false;
  private isInfoHide = false;
  private meta = {
    page: 1,
    size: 20,
    total: 0,
  };
  private info: any = {}; // 群信息

  private target = {}; // 指标
  private list = {
    user: [],
    acts: [],
  };
  private date_st = "";
  private date_et = "";
  private show_st = false;
  private show_et = false;
  private show_et2 = false;
  private minDate = new Date(2011, 0, 1);
  private maxDate = new Date(new Date().getTime() - 24 * 60 * 60 * 1000);
  private default_date = [
    new Date(new Date().getTime() - 24 * 60 * 60 * 1000),
    new Date(new Date().getTime() - 24 * 60 * 60 * 1000),
  ];
  private chart = null;
  private showPopover = false;
  private c_summary_list = [
    {
      name: "群成员总数",
      showPop: false,
      actions: [{ text: "当前查询日期，该社群存在的群成员总数" }],
      prop: "valid_uv",
    },
    {
      name: "当月新入群人数",
      showPop: false,
      actions: [{ text: "当前查询日期内，该社群新入群人员数（人数去重）" }],
      prop: "add_uv",
    },
    {
      name: "累计退群人数",
      showPop: false,
      actions: [{ text: "群创建日起至查询日期该社群每日退群人数之和" }],
      prop: "all_loss_uv",
    },
    {
      name: "退群率",
      showPop: false,
      actions: [
        {
          text: "当前查询日期，新增退群人数与退群前累计群留存成员成员数的占比：新增退群人数/（新增退群人数+群留存成员数）*100%",
        },
      ],
      prop: "loss_uv_fee",
    },
  ];
  private g_summary_list = [
    {
      name: "群会员数",
      showPop: false,
      actions: [{ text: "当前查询日期，群内留存用户中属于会员的总人数" }],
      prop: "member_uv",
    },
    {
      name: "非会员数",
      showPop: false,
      actions: [{ text: "当前查询日期，群内留存用户中属于非会员的总人数" }],
      prop: "general_uv",
    },
    {
      name: "会员率",
      showPop: false,
      actions: [
        {
          text: "当前查询日期，群内会员人数在群留存成员中的占比：（会员总人数/群留存成员数（去除员工））*100%",
        },
      ],
      prop: "member_card_uv",
    },
    // {
    //   name: "白金卡及以下卡等用户入群率",
    //   showPop: false,
    //   actions: [
    //     {
    //       text: "当前查询日期，群内白金卡及以下用户数（白金卡、金卡、名媛、GC）在群留存成员中的占比：（群内白金卡及以下用户数/群留存成员数（去除员工）)*100%",
    //     },
    //   ],
    //   prop: "bp_card_fee",
    // },
  ];
  private tabsList = [
    { tit: "群成员列表", prop: "user" },
    { tit: "活动参与排行榜", prop: "acts" },
  ];
  public get brand(): string {
    return this.$store.getters["user/getEnterpriseMark"];
  }
  /*************************** 生命周期 start ***************************/
  private created() {
    if (this.brand !== "HR") document.body.classList.add("other");
    document.title = "社群画像";

    this.date_st = this.formatDate(
      new Date(new Date().getTime() - 24 * 60 * 60 * 1000)
    );
    this.date_et = this.formatDate(
      new Date(new Date().getTime() - 24 * 60 * 60 * 1000)
    );
    this.getMemberList();
    this.getRefreshTime();
  }
  private mounted() {
    this.getChatInfo();
  }
  private setBorderRadius(type = "add") {
    if (type === "add") {
      window.document.querySelector(".options").classList.add("borderR");
      window.document
        .querySelector(".options .van-dropdown-menu__bar")
        .classList.add("borderR");
    } else {
      window.document.querySelector(".options").classList.remove("borderR");
      window.document
        .querySelector(".options .van-dropdown-menu__bar")
        .classList.remove("borderR");
    }
  }
  // 获取数据更新时间
  private getRefreshTime() {
    getRefreshTime({}).then((res: any) => {
      if (res.code === 200) this.refreshTime = res.data;
    });
  }
  /*************************** method start ***************************/
  private initChart(data) {
    this.$nextTick(() => {
      data = data || [];
      this.chart = echarts.init(this.$refs.chart);
      const options = {
        silent: true,
        legend: [
          {
            orient: "horizontal",
            bottom: 24,
            icon: "rect",
            itemGap: 15,
            x: "center",
            itemWidth: 10,
            itemHeight: 10,
            itemStyle: {
              opacity: 1,
              fontSize: 12,
            },
            data: ["黑钻卡", "钻卡", "白金卡", "金卡"],
          },
          {
            orient: "horizontal",
            bottom: 0,
            icon: "rect",
            itemGap: 15,
            itemWidth: 10,
            itemHeight: 10,
            itemStyle: {
              opacity: 1,
              fontSize: 12,
            },
            data: ["名媛", "GC", "会员", "非会员"],
          },
        ],
        color: [
          "#EBD95D",
          "#E19E57",
          "#D67B5E",
          "#CDB7CB",
          "#D7E9E9",
          "#6A99A7",
          "#116F6A",
        ],
        series: [
          {
            type: "pie",
            radius: "45%",
            center: ["50%", "40%"],
            data,
            itemStyle: {
              opacity: 0.2,
            },
            labelLine: {
              smooth: 0.2,
              length: document.body.clientWidth / 20,
              length2: document.body.clientWidth / 25,
              lineStyle: {
                type: "dashed",
              },
            },
            label: {
              opacity: 1,
              show: true,
              align: "left",
              fontSize: 14,
              color: "#33363C",
              lineHeight: 19,
              formatter: function (params) {
                return `${params.name}\n{x|${params.value}}人\n`;
              },
              rich: {
                x: {
                  color: "#000000",
                  fontSize: 16,
                },
              },
            },
          },
          {
            type: "pie",
            radius: "40%",
            center: ["50%", "40%"],
            data,
            labelLine: false,
            label: {
              opacity: 1,
              show: true,
              position: "inner",
              align: "left",
              fontSize: 10,
              color: "#ffffff",
              lineHeight: 21,
              formatter: function (params) {
                const num = params.percent;
                return `${Math.round(num) === num ? num : num.toFixed(1)}%`;
              },
            },
          },
        ],
      };
      this.chart.setOption(options);
    });
  }
  // 社群画像首页接口
  private async getChatInfo() {
    this.loading = true;
    this.$toast.loading({ message: "加载中...", forbidClick: true });
    let st = new Date(this.date_st);
    let et = new Date(this.date_et);
    let date_st = st.getTime() / 1000; // 精确到毫秒开始时间
    let date_et = (et.getTime() + 24 * 60 * 60 * 1000 - 1) / 1000; // 精确到毫秒结束时间
    const params = {
      chat_id: this.groupChatId,
      sts: date_st,
      ets: parseInt(String(date_et)),
    };
    const {
      code,
      message = "",
      data,
    }: IApiResult = await chatIndexList(params);
    this.loading = false;
    this.$toast.clear();
    if (code === 200) {
      data.chat_info.create_time = this.$helper.dateFormat(
        data.chat_info.create_time,
        "yyyy/MM/dd hh:mm"
      );
      this.info = data;
      this.initChart(this.info[this.selectPieOpt] || []);
    } else {
      this.$toast.fail(message);
      this.initChart(this.info[this.selectPieOpt] || []);
    }
  }
  // 跳转到公共页面
  private getPage(type) {
    let et = new Date(this.date_et);
    let date_et = (et.getTime() + 24 * 60 * 60 * 1000 - 1) / 1000; // 精确到毫秒结束时间
    this.$router.push(
      `/public-page?wnerId=${this.info.chat_info.user_id}&type=${
        type || "4"
      }&ets=${String(parseInt(String(date_et)))}`
    );
  }
  // 选择时间组件
  private formatDate(date) {
    return `${date.getFullYear()}/${date.getMonth() + 1}/${date.getDate()}`;
  }
  private formatter(day) {
    // 当前日期及往后不可选
    let dis_day =
      day.date > new Date(new Date().getTime() - 24 * 60 * 60 * 1000);
    if (dis_day) {
      day.type = "disabled";
    }
    return day;
  }
  private onConfirm(date) {
    this.show_st = false;
    this.show_et = false;
    this.show_et2 = false;
    this.date_st = `${this.formatDate(date[0])}`;
    this.date_et = `${this.formatDate(date[1])}`;
    this.getChatInfo();
  }
  // 点击tab切换
  private changeTab(index) {
    // this.currentTab = index;
    if (!this.list.user.length || !this.list.acts.length) {
      this.$toast.loading({ message: "加载中...", forbidClick: true });
      this.getMemberList();
    }
  }
  // 打开日历组件
  private openCal(type) {
    console.log("打开日历");
    this.show_st = true;
    if (type == 1) {
      this.show_et2 = true;
    } else this.show_et = true;
  }
  // 关闭日历组件
  private closeCal() {
    this.show_et = false;
    this.show_et2 = false;
  }
  // 切换下拉菜单
  private changeOption() {
    this.initChart(this.info[this.selectPieOpt]);
  }
  // 群成员列表
  private async getMemberList() {
    this.loading = true;
    const params = {
      chat_id: this.groupChatId,
      page: 1,
      size: 5,
      friend: 0,
      list_type: this.currentTab,
    };
    const {
      code,
      message = "",
      data,
    }: IApiResult = await chatGroupMembers(params);
    this.$toast.clear();
    this.loading = false;
    if (code === 200) {
      if (data.items && data.items.length) {
        // data.items.forEach(ite => {
        //   ite.join_time = this.$helper.dateFormat(
        //     ite.join_time,
        //     "yyyy/MM/dd hh:mm"
        //   );
        // });
        if (!this.currentTab) this.list.user = data.items;
        else this.list.acts = data.items;
      }
      console.log(this.list.user);
    } else {
      this.$toast.fail(message);
    }
  }
  // 群消息收起或者展开
  private infoTrigger() {
    this.isInfoHide = !this.isInfoHide;
  }
}
</script>
<style lang="scss" scoped>
.group-portrait {
  background: #f4f7fb;
  font-size: 14px;
  color: #000000;
  padding-top: 20px;
  .box-item {
    width: 343px;
    background: #ffffff;
    box-shadow: 0px 6px 12px 0px rgba(222, 224, 227, 0.5);
    border-radius: 20px 20px 19px 19px;
    margin: 0 auto 20px;
    padding: 20px 16px 24px;
    .title {
      font-size: 15px;
      font-weight: 500;
    }
    .group-name {
      line-height: 20px;
      font-weight: 500;
      color: #33363c;
    }
    .group-admin {
      font-size: 12px;
      font-weight: 400;
      color: #000000;
      line-height: 17px;
      padding: 8px 0;
      .tag {
        display: inline-block;
        width: 34px;
        height: 16px;
        background: #222222;
        box-shadow: 0px 10px 23px 0px rgba(0, 0, 0, 0.14);
        border-radius: 8px;
        font-size: 10px;
        font-weight: 500;
        color: #dfbf7e;
        line-height: 16px;
        text-align: center;
        margin-right: 6px;
        span {
          display: inline-block;
          transform: scale(0.83);
        }
      }
    }
    .create-time {
      font-size: 12px;
      font-weight: 400;
      color: #a1a1a1;
      line-height: 17px;
    }
    .base-info {
      margin-top: 15px;
      padding-top: 15px;
      border-top: 1px dashed #f0f0f0;
      .info-item {
        font-size: 12px;
        padding-bottom: 12px;
        display: flex;
        .label {
          color: #a1a1a1;
          width: 63px;
        }
        .val {
          width: 250px;
        }
      }
      &.base-hide {
        height: 60px;
        overflow: hidden;
        margin-bottom: 16px;
      }
    }
    position: relative;
    .hide-shadow {
      position: absolute;
      z-index: 1;
      bottom: 45px;
      left: 16px;
      width: 312px;
      height: 28px;
      background: linear-gradient(
        180deg,
        rgba(255, 255, 255, 0) 0%,
        #ffffff 100%
      );
      transform: translateY(-30px);
    }
    .open-btn {
      width: 311px;
      height: 34px;
      background: #f4f7fb;
      border-radius: 8px;
      color: #a1a1a1;
      line-height: 34px;
      text-align: center;
      font-size: 11px;
      span {
        margin-right: 5px;
      }
    }
    .invite-btn {
      width: 311px;
      height: 33px;
      background: #dfbf7e;
      border-radius: 8px;
      line-height: 33px;
      text-align: center;
      font-size: 11px;
      color: #ffffff;
      font-weight: 400;
    }
    .summary-ul {
      display: flex;
      flex-wrap: wrap;
      padding-top: 16px;
      .summary-li {
        width: 145px;
        height: 106px;
        background: #f4f7fb;
        border-radius: 20px;
        text-align: center;
        padding-top: 22px;
        .num {
          font-size: 26px;
          line-height: 37px;
          font-weight: 500;
        }
        .num-tit {
          color: #777777;
          line-height: 18px;
          padding-top: 5px;
          font-size: 13px;
          span {
            margin-right: 5px;
          }
        }
        &:nth-child(even) {
          margin-left: 20px;
        }
        &:first-child,
        &:nth-child(2) {
          margin-bottom: 20px;
        }
      }
      .line {
        width: 1px;
        height: 45px;
        background: #d9d9d9;
      }
    }
    .g_summary-ul {
      display: flex;
      padding-top: 16px;
      .g_summary-li {
        width: 90px;
        height: 79px;
        background: #f4f7fb;
        border-radius: 16px;
        margin-right: 20px;
        text-align: center;
        padding-top: 15px;
        position: relative;
        overflow: hidden;
        .top-right {
          position: absolute;
          width: 22px;
          height: 20px;
          text-align: center;
          line-height: 20px;
          background: #d9e0ea;
          border-radius: 0px 0px 0px 12px;
          top: -1px;
          right: 0;
          font-size: 12px;
          .van-icon {
            transform: scale(0.8);
          }
        }
        .num {
          font-size: 20px;
          line-height: 28px;
          font-weight: 400;
        }
        .num-tit {
          color: #777777;
          line-height: 18px;
          padding-top: 2px;
          font-size: 13px;
          span {
            margin-right: 4px;
          }
        }
        &:last-child {
          margin-right: 0;
        }
      }
    }
  }
  .rank-list {
    .c-ul {
      margin-bottom: 20px;
      .c-li {
        display: flex;
        align-items: center;
        margin-top: 20px;
        .avatar {
          width: 44px;
          height: 44px;
          border-radius: 100%;
          display: block;
        }
        .info-div {
          width: 203px;
          padding-left: 8px;
          .c-name {
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
            width: 168px;
            font-weight: 500;
            line-height: 20px;
          }
          .other-p {
            font-size: 11px;
            color: #a1a1a1;
            font-weight: 400;
            line-height: 16px;
            padding-top: 4px;
            .line {
              margin: 0 6px;
            }
            span {
              display: inline-block;
              transform: scale(0.95);
            }
          }
        }
        .btn {
          width: 56px;
          height: 20px;
          border-radius: 4px;
          border: 1px solid #5d5d5d;
          line-height: 20px;
          text-align: center;
          font-size: 12px;
          font-weight: 400;
        }
      }
    }
  }
}
</style>
<style lang="scss">
.box-item {
  .van-tabs__line {
    bottom: 18px;
    background-color: #000000;
  }
}
.group-portrait {
  .van-dropdown-item__content {
    width: 120px;
    left: initial;
    right: 30px;
    background: #ffffff;
    box-shadow: 0px 4px 13px 0px rgba(0, 0, 0, 0.06);
    border-radius: 0px 0px 8px 8px;
  }
  .van-cell__title {
    white-space: nowrap;
  }
  .van-dropdown-menu__bar {
    background: #f4f4f4;
    border-radius: 8px;
  }
  .van-dropdown-menu__title {
    font-size: 14px;
    font-weight: 400;
    color: #808080 !important;
    &::after {
      border-width: 4px;
      margin-top: -7px;
      border-color: transparent transparent #808080 #808080;
    }
  }
  .van-dropdown-menu__title--down::after {
    margin-top: -1px;
  }
  .van-ellipsis {
    padding-right: 4px;
  }
  .van-cell {
    color: #b5b5b5;
  }
  .van-icon-success:before {
    color: #caa361;
  }
}
.van-popover__action {
  width: 110px;
  padding: 7px;
  height: auto;
  font-size: 12px;
  color: #0b0a09;
  line-height: 17px;
}

.other .group-portrait {
  .portrait-middle > div .portrait-list {
    background: rgba(91, 137, 254, 0.1) !important;
    &::before {
      content: "";
      border-bottom-color: rgba(91, 137, 254, 0.1) !important;
    }
  }
  .item-statistics {
    background: linear-gradient(180deg, #fafcff 0%, #eef3fe 100%) !important;
  }
  .pie .title .icon {
    background: #cddbfe !important;
  }
  .pie .title .icon:nth-child(2) {
    background: #5b89fe !important;
  }
  .van-icon-success:before {
    color: #5b89fe !important;
  }
  .van-calendar__day.van-calendar__day--start,
  .van-calendar__day.van-calendar__day--end,
  .van-calendar__day.van-calendar__day--start-end {
    background-color: #5b89fe !important;
  }
  .item-list .item-right .item-btns .item-type.underLine {
    color: #5b89fe !important;
  }
}
</style>
