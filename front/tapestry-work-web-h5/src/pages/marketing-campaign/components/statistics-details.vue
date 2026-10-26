<template>
  <div id="activity">
    <van-nav-bar title="统计" left-arrow @click-left="goBack" />
    <div class="activity-detail">
      <div
        class="detail-top"
        :style="{
          background: attribute.word_color.backgroup,
          paddingTop: statisticsDetail.title.length > 14 ? '13px' : '26px',
        }"
        v-if="statisticsDetail.attribute_type === 1"
      >
        <div class="top-all-detail">
          <div class="top-detail">
            <h3
              :style="{
                color: attribute.word_color.title,
                marginBottom:
                  statisticsDetail.title.length > 14 ? '10px' : '31px',
              }"
            >
              {{ statisticsDetail.title }}
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
                    statisticsDetail.start_time,
                    "yyyy/MM/dd hh:mm"
                  )
                }}-{{
                  $helper.dateFormat(
                    statisticsDetail.end_time,
                    "yyyy/MM/dd hh:mm"
                  )
                }}
              </div>
              <div class="detail-topImg">
                <img :src="statisticsDetail.cover_url" alt="" />
                <span>{{ status }}</span>
              </div>
              <div class="detail-explain">
                <div class="explain-left">
                  活动类型：<span>{{
                    statisticsDetail.type === 1
                      ? "报名活动"
                      : statisticsDetail.type === 2
                      ? "签到活动"
                      : "问卷活动"
                  }}</span>
                </div>
                <div class="explain-right">
                  报名人数：<span>{{ statisticsDetail.total }}</span
                  >人
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
          ></div>
        </div>
      </div>
      <div
        class="join-list"
        v-if="clientList.lenght !== 0"
        :style="{
          marginTop: statisticsDetail.attribute_type === 1 ? '-125px' : '0px',
        }"
      >
        <van-list
          v-model="loading"
          :finished="finished"
          :finished-text="clientList.length ? '没有更多了' : ''"
          :immediate-check="false"
          @load="onLoad"
        >
          <div
            class="list-content"
            v-for="(item, index) in clientList"
            :key="index"
          >
            <div class="list-content-top">
              <div class="content-top-left">
                <img class="content-top-img" :src="item.img_url" alt="" />
                <div class="content-top-info">
                  <p>
                    {{
                      item.nickname.length > 10
                        ? item.nickname.substring(0, 10) + "..."
                        : item.nickname
                    }}
                  </p>
                  <p class="top-info-type">
                    来源类型：{{
                      item.source_type == 1 ? "企业社群" : "个人聊天"
                    }}
                  </p>
                </div>
              </div>
              <div v-if="item.user_id !== '' && item.user_id !== null">
                <div
                  v-if="item.relation === 0"
                  class="content-top-button"
                  @click="toAddFriend(item)"
                >
                  <!-- style="background: rgba(202, 163, 97, 0.3)" -->
                  加好友
                </div>
                <div
                  v-else
                  class="content-top-button"
                  @click="toSendMessage(item)"
                >
                  打招呼
                </div>
              </div>
            </div>
            <div class="list-content-bottom">
              <p v-if="statisticsDetail.type !== 2">
                提交时间：{{
                  $helper.dateFormat(item.time, "yyyy/MM/dd hh:mm:ss")
                }}
              </p>
              <p v-if="statisticsDetail.type === 2">
                最新签到时间：{{
                  $helper.dateFormat(item.time, "yyyy/MM/dd hh:mm:ss")
                }}
              </p>
              <p v-if="statisticsDetail.type === 1 && item.finish_status !== 0">
                审核状态：{{
                  item.finish_status == 1
                    ? "审核中"
                    : item.finish_status == 2
                    ? "审核通过"
                    : "审核失败"
                }}
              </p>
              <p v-if="statisticsDetail.type === 2">
                签到天数：{{ item.total }}
              </p>
            </div>
          </div>
        </van-list>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import { Tab, Tabs, NavBar } from "vant";
import { statisticsDetailApi } from "@/api/marketing-campaign/index";
@Component({
  components: {
    [Tab.name]: Tab,
    [Tabs.name]: Tabs,
    [NavBar.name]: NavBar,
  },
})
export default class extends Vue {
  @Prop({ default: 0, type: Number }) id!: number;
  @Prop({ default: "", type: String }) status!: string;
  @Prop({ default: 0, type: Number }) defaultTab!: number;
  private statisticsDetail = {};
  private clientList = [];
  private titletBg = require("@/assets/marketing-campaign/gradual-bg.png");

  private total = 0;
  private isRequest = false;
  private page = 1;
  private loading = false;
  private finished = false;
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

  private created() {
    window.scrollTo(0, 0);
  }
  private mounted() {
    this.getStatisticsDetail(1); //活动统计
    this.getStatisticsDetail(2); //列表
  }
  // 触底加载
  private onLoad() {
    if (this.isRequest) {
      if (this.total > this.clientList.length) {
        this.page = this.page + 1;
        this.getStatisticsDetail(2);
      } else {
        this.finished = true;
      }
    }
  }
  async getStatisticsDetail(type: number): Promise<void> {
    await statisticsDetailApi({
      id: this.id,
      scope: type,
      page: this.page,
      size: 10,
    }).then((res: any) => {
      if ((res.httpCode as number) == 200) {
        if (type == 1) {
          this.statisticsDetail = res.data;
          if (res.data.attribute_type === 1) {
            this.attribute = JSON.parse(res.data.attribute);
          }
        } else {
          // this.clientList = res.data.list;
          this.clientList = [...this.clientList, ...res.data.list];
          this.total = res.data.meta.total;
          // 加载状态结束
          this.loading = false;
          this.isRequest = true;
        }
      } else {
        this.$toast.fail(res.status);
      }
    });
  }
  private toSendMessage(item) {
    const message = {
      userIds: "", //参与会话的企业成员列表，格式为userid1;userid2;...，用分号隔开。
      externalUserIds: item.user_id, // 参与会话的外部联系人列表，格式为userId1;userId2;…，用分号隔开。
      groupName: "", // 会话名称。单聊时该参数传入空字符串""即可。
      chatId: "", // 若要打开已有会话，需指定此参数。如果是新建会话，chatId必须为空串
    };
    // 调用企微接口
    window.wx.invoke("openEnterpriseChat", message, (res) => {
      if (res.err_msg == "openEnterpriseChat:ok") {
        //发送成功
        console.log("成功回调");
      }
    });
  }
  private toAddFriend(item) {
    const message = {
      type: 2, //1表示该userid是企业成员，2表示该userid是外部联系人
      userid: item.user_id, //可以是企业成员，也可以是外部联系人
    };
    // 调用企微接口
    window.wx.invoke("openUserProfile", message, (res) => {
      if (res.err_msg == "openUserProfile:ok") {
        //发送成功
        console.log("跳转成功");
      }
    });
  }
  // 返回按钮
  private goBack() {
    const item = {
      type: 2,
      defaultTab: this.defaultTab,
    };
    this.$emit("changeType", item);
  }
}
</script>
<style lang="scss" scoped>
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
}
#activity {
  // background: #f6f5f3;
  // min-height: 100vh;
  .activity-detail {
    margin-top: 45px;
    .detail-top {
      height: 480px;
      padding: 25px 0px;
      .top-detail {
        padding: 0px 15px;
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
          height: 248px;
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
          .detail-topImg {
            position: relative;
            img {
              display: block;
              margin: 15px 0px 13px;
              width: 100%;
              height: 158px;
              border-radius: 12px;
            }
            span {
              position: absolute;
              top: 15px;
              right: 15px;
              padding: 7px 10px;
              background: rgba(0, 0, 0, 0.4);
              box-shadow: 0px 0px 5px 0px rgba(202, 163, 97, 0.3),
                0px 2px 8px 0px rgba(0, 0, 0, 0.09);
              border-radius: 15px;
              font-size: 12px;
              color: #ffffff;
              line-height: 16px;
              text-shadow: 0px 0px 5px rgba(202, 163, 97, 0.3);
            }
          }
          .detail-explain {
            display: flex;
            justify-content: space-between;
            padding-right: 30px;
            font-size: 14px;
            color: #33363c;
            line-height: 20px;
            text-shadow: 0px 0px 5px rgba(202, 163, 97, 0.3);
            span {
              color: #caa361;
              font-weight: bold;
            }
            .explain-right {
              span {
                font-size: 18px;
                margin-right: 2px;
              }
            }
          }
        }
      }
      .detail-title-bg {
        width: 100%;
        height: 150px;
        margin-top: -10px;
      }
    }
  }
  .join-list {
    width: 100%;
    padding: 7px 15px 25px 15px;
    .list-content {
      width: 100%;
      padding: 15px;
      background: #ffffff;
      box-shadow: 0px 9px 16px 0px rgba(0, 0, 0, 0.05);
      border-radius: 5px;
      margin-bottom: 15px;
      .list-content-top {
        display: flex;
        justify-content: space-between;
        align-items: center;
        .content-top-left {
          display: flex;
          align-items: center;
          .content-top-img {
            width: 50px;
            height: 50px;
            border-radius: 50%;
            margin-right: 8px;
          }
          .content-top-info {
            font-size: 16px;
            color: #33363c;
            line-height: 26px;
            .top-info-type {
              font-size: 12px;
              color: #808080;
              line-height: 22px;
            }
          }
        }
        .content-top-button {
          display: flex;
          align-items: center;
          justify-content: center;
          height: 28px;
          padding: 6px 10px;
          font-size: 12px;
          color: #ffffff;
          background: #caa361;
          border-radius: 2px;
        }
      }
      .list-content-bottom {
        background: #f6f6f9;
        border-radius: 4px;
        padding: 7px 15px;
        font-size: 12px;
        color: #808080;
        line-height: 17px;
        margin-top: 10px;
      }
    }
  }
}
</style>
