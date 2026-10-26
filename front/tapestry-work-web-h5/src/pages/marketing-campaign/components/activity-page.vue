<template>
  <div id="activity">
    <van-nav-bar
      :title="activityDetail.type_text"
      left-arrow
      @click-left="goBack()"
    />
    <div class="activity-detail">
      <div class="hint-info">
        <img src="@/assets/marketing-campaign/hint.png" alt="" />
        当前为预览页面，回答不计入结果
      </div>
      <div class="detail-top">
        <!-- 整图 -->
        <div
          class="imgTop-detail"
          v-if="activityDetail.page_config.attribute_type === 2"
        >
          <img class="bgImg" :src="attribute.img.url" alt="" />
          <div class="detail-title">
            <img src="@/assets/marketing-campaign/title-bg.png" alt="" />
            <span
              :style="{
                background: attribute.button_color.solid_button,
                color: attribute.button_color.solid_word,
              }"
              >{{
                activityDetail.type === 1
                  ? "报名信息"
                  : activityDetail.type === 2
                  ? "签到信息"
                  : "问卷信息"
              }}</span
            >
          </div>
        </div>
        <!-- 配置 -->
        <div
          v-else
          class="top-detail"
          :style="{
            background: attribute.word_color.backgroup,
          }"
        >
          <div
            class="detail-top-bg"
            :style="{
              paddingTop: activityDetail.title.length > 14 ? '18px' : '35px',
            }"
          >
            <h3
              :style="{
                color: attribute.word_color.title,
                marginBottom:
                  activityDetail.title.length > 14 ? '15px' : '35px',
              }"
            >
              {{ activityDetail.title }}
            </h3>
            <div class="detail-content">
              <div
                class="detail-date"
                :style="{
                  color: attribute.word_color.period,
                }"
              >
                活动时间：{{
                  $helper.dateFormat(
                    activityDetail.start_time,
                    "yyyy/MM/dd hh:mm"
                  )
                }}-{{
                  $helper.dateFormat(
                    activityDetail.end_time,
                    "yyyy/MM/dd hh:mm"
                  )
                }}
              </div>
              <img
                class="detail-content-img"
                :src="activityDetail.cover_url"
                alt=""
              />
              <div
                class="detail-explain"
                :style="{
                  color: attribute.word_color.rule,
                }"
                @click="toShowPopup"
              >
                <div class="detail-explain-text" v-html="activity_rule">
                  <!-- {{ activityDetail.activity_rule }} -->
                </div>
                <div>
                  <van-icon name="arrow" color="#000" />
                </div>
              </div>
            </div>
          </div>
          <div
            class="detail-title-bg"
            :style="{
              backgroundImage: 'url(' + titletBg + ')',
              backgroundSize: '100% 100%',
            }"
          >
            <div class="detail-title">
              <img src="@/assets/marketing-campaign/title-bg.png" alt="" />
              <span
                :style="{
                  background: attribute.button_color.solid_button,
                  color: attribute.button_color.solid_word,
                }"
                >{{
                  activityDetail.type === 1
                    ? "报名信息"
                    : activityDetail.type === 2
                    ? "签到信息"
                    : "问卷信息"
                }}</span
              >
            </div>
          </div>
        </div>
      </div>
      <!-- 签到组件 -->
      <sign-list
        v-if="activityDetail.type === 2"
        :signInfoList="signInfoList"
      ></sign-list>
      <sign-page v-else :moduleContent="listContent.templateList"></sign-page>
    </div>
    <div class="signPage-bottom">
      <div
        class="bottom-button"
        @click="toSend"
        :style="{
          width: '100%',
          background:
            activityDetail.page_config.attribute_type === 1
              ? attribute.button_color.solid_button
              : '#AF8B4E',
          color:
            activityDetail.page_config.attribute_type === 1
              ? attribute.button_color.solid_word
              : '#EFEFEF',
        }"
      >
        立即发送
      </div>
    </div>
    <van-action-sheet v-model="popupShow" title="活动规则">
      <div class="content" v-html="activity_rule"></div>
      <div class="signPage-bottom">
        <div class="new-bottom-button" @click="popupShow = false">确认</div>
      </div>
    </van-action-sheet>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import signPage from "./sign-page.vue";
import signList from "./sign-list.vue";
import Store from "@/store/index";
const store: any = Store;
import { Tab, Tabs, NavBar, Icon, ActionSheet } from "vant";
import { activityDetailApi } from "@/api/marketing-campaign/index";
@Component({
  components: {
    signPage,
    signList,
    [Tab.name]: Tab,
    [Tabs.name]: Tabs,
    [Icon.name]: Icon,
    [NavBar.name]: NavBar,
    [ActionSheet.name]: ActionSheet,
  },
})
export default class extends Vue {
  @Prop({ default: 0, type: Number }) id!: number;
  @Prop({ default: 0, type: Number }) defaultTab!: number;
  private signInfoList = [];
  private titletBg = require("@/assets/marketing-campaign/gradual-bg.png");
  private activityDetail = {
    type: 0,
    cover_url: "",
    title: "",
    page_config: {
      share_pic: "",
      share_title: "",
    },
    activity_rule: "",
  };
  private activity_rule: "";
  private attribute = {
    button_color: {
      hollow_button: "",
      hollow_word: "",
      solid_button: "",
      solid_word: "",
    },
    word_color: {
      backgroup: "",
      period: "",
      rule: "",
      title: "",
    },
  };
  private listContent = {
    templateList: [],
  };
  private popupShow = false;
  private created() {
    window.scrollTo(0, 0);
  }
  private mounted() {
    this.getActivityDetail();
  }
  async getActivityDetail(): Promise<void> {
    const that = this;
    await activityDetailApi({
      id: this.id,
    }).then((res: any) => {
      if ((res.httpCode as number) === 200) {
        that.attribute = JSON.parse(res.data.page_config.attribute);
        that.activityDetail = res.data;
        that.activity_rule = res.data.activity_rule.replace(/\n/g, "<br>");
        if (res.data.type === 2) {
          that.listContent = JSON.parse(res.data.sign_in.content);
          that.signInfoList = that.changeTime(
            res.data.sign_in.sign_start,
            res.data.sign_in.sign_end
          );
        } else if (res.data.type === 1) {
          that.listContent = JSON.parse(res.data.sign_up.content);
        } else if (res.data.type === 3) {
          that.listContent = JSON.parse(res.data.question.content);
        }
        this.listContent.templateList.forEach((item, index) => {
          let sequence = index + 1 < 10 ? "0" + (index + 1) : index + 1;
          return (item.sequence = sequence);
        });
        // console.log(this.attribute, "头部样式");
        // console.log(this.listContent.templateList, "列表样式");
      } else {
        that.$toast.fail(res.status);
      }
    });
  }
  // 格式化时间
  private changeTime(start: any, end: any) {
    const startTime =
      this.$helper.dateFormat(start, "yyyy/MM/dd") + " 00:00:00";
    const endTime = this.$helper.dateFormat(end, "yyyy/MM/dd") + " 00:00:00";
    const timeLead =
      Number(Date.parse(endTime)) - Number(Date.parse(startTime));
    const dateLead = timeLead / 86400 / 1000;
    let arr = [];
    for (var i = 0; i <= dateLead; i++) {
      let time = Date.parse(startTime) / 1000 + i * 86400;
      let newTime = this.$helper.dateFormat(time, "MM-dd");
      arr = [...arr, newTime];
    }
    return arr;
  }
  // 转发小程序
  private toSend() {
    let pageUrl = "";
    // const jointUrl = `.html?activityId=${this.$route.query.id}&userId=${this.$store.state.user.USERINFO.user_id}`;
    const jointUrl = `.html?activityId=${this.id}&userId=${this.$store.state.user.USERINFO.user_id}&hash=${store.state.common.ENTERPRISE_HASH}`;
    if (this.activityDetail.type === 1) {
      pageUrl = `pages/sign-up/sign-up${jointUrl}`;
    } else if (this.activityDetail.type === 2) {
      pageUrl = `pages/sign-in/sign-in${jointUrl}`;
    } else if (this.activityDetail.type === 3) {
      pageUrl = `pages/questionnaire/questionnaire${jointUrl}`;
    }
    this.$toast.loading({
      message: "发送中...",
      forbidClick: true,
    });
    let appid = "";
    if (window.location.hostname === "tapestry.vchangyi.com") {
      appid = "wx387d289f4bf306ed";
    } else {
      appid = "wx5e1eb0d63e504f98";
    }
    let that = this;
    const message = {
      msgtype: "miniprogram",
      miniprogram: {
        appid: appid,
        title: that.activityDetail.page_config.share_title,
        imgUrl: that.activityDetail.page_config.share_pic,
        page: pageUrl,
      },
    };
    window.wx.invoke("sendChatMessage", message, (res) => {
      if (res.err_msg !== "sendChatMessage:ok") {
        if (res.err_msg == "sendChatMessage:cancel") {
          return;
        }
        //发送成功
        that.$toast.fail("企微官方后台未关联小程序");
      }
    });
  }
  // 返回按钮
  private goBack() {
    const item = {
      type: 1,
      defaultTab: this.defaultTab,
    };
    this.$emit("changeType", item);
  }
  private toShowPopup() {
    this.popupShow = true;
  }
}
</script>
<style lang="scss" scoped>
#activity {
  background: #f6f5f3;
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
    .van-action-sheet {
      padding-top: 10px;
    }
    .van-action-sheet__close {
      padding-top: 10px;
    }
  }
  .activity-detail {
    padding-bottom: 80px;
    margin-top: 45px;
    .hint-info {
      img {
        width: 18px;
        height: 18px;
        margin-right: 6px;
      }
      font-size: 14px;
      color: #fa4046;
      line-height: 20px;
      background: #fec6c8;
      padding: 10px 15px;
      display: flex;
      align-items: center;
      // justify-content: center;
    }
    .detail-top {
      width: 100%;
      // position: relative;
      .imgTop-detail {
        height: 480px;
        position: relative;
        // position: absolute;
        .bgImg {
          width: 100%;
          height: 480px;
        }
        .detail-title {
          position: absolute;
          bottom: 50px;
          left: 15px;
          img {
            width: 100%;
            height: 1px;
          }
          span {
            width: 167px;
            height: 32px;
            display: flex;
            justify-content: center;
            align-items: center;
            background: #caa361;
            // box-shadow: 0px 0px 9px 0px #caa362, 0px 0px 6px 0px #b68c44;
            border-radius: 16px;
            border: 2px solid #fdffff;
            position: absolute;
            top: 20px;
            left: 25%;
            font-size: 18px;
            font-weight: bold;
            color: #ffffff;
            line-height: 25px;
          }
        }
      }
      .top-detail {
        height: 480px;
        .detail-top-bg {
          padding: 30px 15px 0px;
          h3 {
            font-size: 24px;
            font-weight: bold;
            color: #af8b4e;
            line-height: 40px;
            margin-bottom: 30px;
            text-align: center;
          }
          .detail-content {
            width: 100%;
            height: 273px;
            padding: 0px 15px 15px 15px;
            background: linear-gradient(
              180deg,
              rgba(255, 255, 255, 0.3) 0%,
              rgba(255, 255, 255, 0.6) 100%
            );
            box-shadow: 0px 0px 5px 0px rgba(202, 163, 97, 0.3);
            border-radius: 12px;
            .detail-date {
              width: 316px;
              height: 30px;
              background: #ffffff;
              box-shadow: 0px 0px 5px 0px rgba(202, 163, 97, 0.3),
                0px 2px 8px 0px rgba(0, 0, 0, 0.09);
              border-radius: 0px 0px 12px 12px;
              margin: 0 auto;
              text-align: center;
              font-size: 12px;
              color: #808080;
              line-height: 30px;
              text-shadow: 0px 0px 5px rgba(202, 163, 97, 0.3);
            }
            .detail-content-img {
              display: block;
              margin: 15px 0px;
              width: 100%;
              height: 158px;
              border-radius: 12px;
            }
            .detail-explain {
              font-size: 14px;
              color: #33363c;
              line-height: 20px;
              text-shadow: 0px 0px 5px rgba(202, 163, 97, 0.3);
              display: flex;
              align-items: flex-end;
              .detail-explain-text {
                width: 98%;
                overflow: hidden;
                text-overflow: ellipsis;
                display: -webkit-box;
                -webkit-line-clamp: 2;
                line-clamp: 2;
                -webkit-box-orient: vertical;
              }
            }
          }
        }
        .detail-title-bg {
          width: 100%;
          height: 133px;
          margin-top: -30px;
          padding: 0px 15px 0px;
          .detail-title {
            position: relative;
            img {
              width: 100%;
              height: 1px;
              position: absolute;
              top: 75px;
              box-shadow: 0px 0px 4px 0px #af8b4e;
            }
            span {
              width: 167px;
              height: 32px;
              display: flex;
              justify-content: center;
              align-items: center;
              background: #caa361;
              border-radius: 16px;
              border: 2px solid #fdffff;
              position: absolute;
              top: 60px;
              left: 25%;
              font-size: 18px;
              font-weight: bold;
              color: #ffffff;
              line-height: 25px;
            }
          }
        }
      }
    }
  }
  .signPage-bottom {
    width: 100%;
    position: fixed;
    bottom: 0px;
    background: #fff;
    padding: 5px 15px 15px;
    font-size: 14px;
    font-weight: 400;
    line-height: 20px;
    display: flex;
    justify-content: space-between;
    z-index: 99;
    .bottom-button {
      display: flex;
      justify-content: center;
      align-items: center;
      height: 42px;
      border-radius: 12px;
      color: #ffffff;
      margin: 10px 0px;
    }
    .new-bottom-button {
      width: 100%;
      background: #caa361;
      display: flex;
      justify-content: center;
      align-items: center;
      height: 42px;
      border-radius: 12px;
      color: #ffffff;
      margin: 10px 0px;
    }
  }
  .content {
    padding: 0px 30px 110px;
    font-size: 14px;
    color: #808080;
    line-height: 24px;
  }
}
</style>
