<template>
  <div class="portrait-association">
    <div class="portrait-top">
      <img
        class="portrait-img"
        :src="
          info.group_image_url
            ? info.group_image_url
            : require('@/assets/default/default-person.png')
        "
        alt=""
      />
      <div class="portrait-name">{{ info.name }}</div>
    </div>
    <div class="border-line"></div>
    <div class="portrait-middle">
      <div class="portrait-item">
        <div class="item-label">建群时间</div>
        <div class="item-content">{{ info.create_time }}</div>
      </div>
      <div class="portrait-item">
        <div class="item-label">成员性别</div>
        <div class="item-content sex-box">
          <div class="sex-item">
            <img
              class="sex-icon"
              src="@/assets/customer-portrait/man-sex.png"
              alt=""
            />
            <div class="sex-num">{{ info.male_total || 0 }}人</div>
          </div>
          <div class="sex-item">
            <img
              class="sex-icon"
              src="@/assets/customer-portrait/woman-sex.png"
              alt=""
            />
            <div class="sex-num">{{ info.gender_female_total || 0 }}人</div>
          </div>
          <div class="sex-item">
            <img
              class="sex-icon"
              src="@/assets/customer-portrait/unknown-sex.png"
              alt=""
            />
            <div class="sex-num">{{ info.unknown_gender_total || 0 }}人</div>
          </div>
        </div>
      </div>
      <div class="portrait-item">
        <div class="item-label">会员分布</div>
        <div class="item-content sex-box">
          <div class="sex-item">
            <img
              class="sex-icon"
              src="@/assets/customer-portrait/vip.png"
              alt=""
            />
            <div class="sex-num">{{ info.group_vip_member_num || 0 }}人</div>
          </div>
          <div class="sex-item">
            <img
              class="sex-icon"
              src="@/assets/customer-portrait/un-vip.png"
              alt=""
            />
            <div class="sex-num">
              {{ info.group_member_num - info.group_vip_member_num || 0 }}人
            </div>
          </div>
        </div>
      </div>
      <div
        class="statistics-box m-b-15"
        v-if="JSON.stringify(info.group_vip_group) !== '{}'"
      >
        <div class="statistics-label">等级分布</div>
        <div class="statistics-card vip">
          <div
            class="item-statistics"
            v-for="(item, key, index) in info.group_vip_group"
            :key="index"
          >
            <div class="item-title">{{ key }}</div>
            <div class="item-num">{{ item || 0 }}</div>
          </div>
        </div>
      </div>
      <div class="portrait-item">
        <div class="item-label">群主</div>
        <div class="item-content">{{ info.group_leader_name }}</div>
      </div>
      <div class="statistics-box m-b-6">
        <div class="statistics-label">统计</div>
        <div class="statistics-card">
          <div class="item-statistics">
            <div class="item-num">{{ info.today_active_user || 0 }}</div>
            <div class="item-title">今日活跃</div>
          </div>
          <div class="item-statistics">
            <div class="item-num">{{ info.today_join_user || 0 }}</div>
            <div class="item-title">今日入群</div>
          </div>
          <div class="item-statistics">
            <div class="item-num">{{ info.today_leave_user || 0 }}</div>
            <div class="item-title">今日退群</div>
          </div>
        </div>
      </div>
    </div>

    <!--  -->
    <van-tabs
      v-model="currentTab"
      title-active-color="#0095F6"
      border
      ref="tabs"
      @click="changeTab"
    >
      <van-tab title="成员列表" :name="1">
        <div class="list-box">
          <van-list
            v-if="currentTab == 1 && list.length > 0"
            v-model="loading"
            :finished="noMore"
            finished-text="没有更多了~"
            @load="getMoreOrder"
          >
            <div class="item-list" v-for="(item, index) in list" :key="index">
              <div class="item-img">
                <img
                  :src="
                    item.avatar
                      ? item.avatar
                      : require('@/assets/default/default-person.png')
                  "
                  alt=""
                />
              </div>
              <div class="item-right">
                <div class="item-name-box">
                  <div class="item-name">{{ item.nick_name || "-" }}</div>
                  <div
                    class="item-btn"
                    @click="showPopup(item)"
                    v-if="item.external_userid"
                  >
                    查看添加员工
                  </div>
                  <div class="item-btn-else" v-else>暂无添加员工</div>
                </div>
                <div class="item-form">
                  <div class="item-form-label">入群时间:</div>
                  <div class="item-form-content">{{ item.join_time }}</div>
                  <div class="item-type">
                    {{
                      item.join_scene == 1
                        ? "直接邀请入群"
                        : item.join_scene == 2
                        ? "邀请链接入群"
                        : "扫码入群"
                    }}
                  </div>
                </div>
              </div>
            </div>
          </van-list>
        </div>
        <empty-page
          v-if="list.length == 0 && !loading"
          picture="noData"
          paddingTop="60"
          paddingBottom="60"
          :text="currentTab == 1 ? '暂无成员' : '暂无活跃成员'"
        ></empty-page>
      </van-tab>
      <!-- <van-tab title="微信优惠券" name="2"><coupon-goods-list /></van-tab> -->
      <van-tab title="活跃排行榜" :name="2">
        <div class="list-box">
          <van-list
            v-if="currentTab == 2 && list.length > 0"
            v-model="loading"
            :finished="noMore"
            finished-text="没有更多了~"
            @load="getMoreOrder"
          >
            <div class="item-list" v-for="(item, index) in list" :key="index">
              <div class="item-img">
                <img
                  :src="
                    item.avatar
                      ? item.avatar
                      : require('@/assets/default/default-person.png')
                  "
                  alt=""
                />
              </div>
              <div class="item-right">
                <div class="item-name-box">
                  <div class="item-name">{{ item.nick_name || "-" }}</div>
                  <div
                    class="item-btn"
                    @click="showPopup(item)"
                    v-if="item.external_userid"
                  >
                    查看添加员工
                  </div>
                  <div class="item-btn-else" v-else>暂无添加员工</div>
                </div>
                <div class="item-form">
                  <div class="item-form-label">最新聊天时间:</div>
                  <div class="item-form-content">
                    {{ item.last_chat_at || "-" }}
                  </div>
                  <div class="item-type">
                    {{
                      item.join_scene == 1
                        ? "直接邀请入群"
                        : item.join_scene == 2
                        ? "邀请链接入群"
                        : "扫码入群"
                    }}
                  </div>
                </div>
                <div class="item-form">
                  <div class="item-form-label">累计聊天数:</div>
                  <div class="item-form-content">{{ item.message_total }}</div>
                </div>
                <div class="item-form">
                  <div class="item-form-label">今日聊天数:</div>
                  <div class="item-form-content">
                    {{ item.today_message_total }}
                  </div>
                </div>
              </div>
            </div>
          </van-list>
        </div>
        <empty-page
          v-if="list.length == 0 && !loading"
          picture="noData"
          paddingTop="60"
          paddingBottom="60"
          :text="currentTab == 1 ? '暂无成员' : '暂无活跃成员'"
        ></empty-page>
      </van-tab>
      <!-- 活动参与排行榜 -->
      <van-tab title="活动参与排行榜" :name="4">
        <div class="list-box">
          <van-list
            v-if="currentTab == 4 && list.length > 0"
            v-model="loading"
            :finished="noMore"
            finished-text="没有更多了~"
            @load="getMoreOrder"
          >
            <div class="item-list" v-for="(item, index) in list" :key="index">
              <div class="item-index">
                {{ index &lt; 9 ? "0" + (index + 1) : index + 1 }}
              </div>
              <div class="item-img">
                <img
                  :src="
                    item.avatar
                      ? item.avatar
                      : require('@/assets/default/default-person.png')
                  "
                  alt=""
                />
              </div>
              <div class="item-right">
                <div class="item-name-box">
                  <div class="item-name">{{ item.nick_name || "-" }}</div>
                  <div class="item-btn" @click="showActivity(item)" v-if="true">
                    查看参加活动
                  </div>
                  <div class="item-btn-else" v-else>暂无参加活动</div>
                </div>
                <div class="item-form">
                  <div class="item-form-label">参加活动数:</div>
                  <div class="item-form-content">{{ item.join_num }}</div>
                  <div class="item-type activity">
                    <div
                      class="item-btn"
                      @click="showPopup(item)"
                      v-if="item.unionid"
                    >
                      查看添加员工
                    </div>
                    <div class="item-btn-else" v-else>暂无添加员工</div>
                  </div>
                </div>
              </div>
            </div>
          </van-list>
        </div>
        <empty-page
          v-if="list.length == 0 && !loading"
          picture="noData"
          paddingTop="60"
          paddingBottom="60"
          text="暂无活动"
        ></empty-page>
      </van-tab>
    </van-tabs>

    <!-- 添加的员工列表弹层-->
    <staff-popup
      v-if="showStaffPopup"
      :external="currentExternal"
      @close="showStaffPopup = false"
    ></staff-popup>

    <!-- 参与活动列表列表弹层-->
    <activity-popup
      v-if="showActivityPopup"
      :unionId="union_id"
      @close="showActivityPopup = false"
    ></activity-popup>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import staffPopup from "@/components/customer-portrait/staff-popup.vue";
import activityPopup from "@/components/customer-portrait/activity-popup.vue";
import emptyPage from "@/components/base/empty-page.vue";
import { groupChatInfoApi, groupChatMemberApi } from "@/api/customer-portrait";
@Component({
  components: {
    emptyPage,
    staffPopup,
    activityPopup,
  },
})
export default class extends Vue {
  @Prop({ default: 0, type: String }) groupChatId!: string;
  /*************************** data start ***************************/
  private currentTab = 1; // 默认选中的tab
  private showStaffPopup = false; // 控制添加员工弹框显隐
  private loading = false;
  private noMore = false;
  private meta = {
    page: 1,
    size: 20,
    total: 0,
  };
  private currentExternal = ""; // 当前点击用户的外部联系人id
  private info = {};
  private list = [];
  private showActivityPopup = false; // 控制参与活动弹窗
  private enterpriseMark = ""; // 企业标识
  private union_id = "";
  /*************************** 生命周期 start ***************************/
  private created() {
    this.enterpriseMark = this.$store.getters["user/getEnterpriseMark"];
    this.groupChatDetail();
    this.groupChatMember();
  }
  /*************************** method start ***************************/
  // 社群详情接口
  private async groupChatDetail() {
    const {
      code,
      message = "",
      data,
    }: IApiResult = await groupChatInfoApi(this.groupChatId);
    if (code === 200) {
      data.create_time = this.$helper.dateFormat(
        data.create_time,
        "yyyy-MM-dd hh:mm"
      );
      this.info = data;
    } else {
      this.$toast.fail(message);
    }
  }
  // 用户列表接口
  private async groupChatMember() {
    this.loading = true;
    const params = {
      group_chat_id: this.groupChatId,
      type: this.currentTab,
      page: this.meta.page,
      size: this.meta.size,
    };
    // 活动排行榜仅展示Top10
    if (this.currentTab == 4) params.size = 10;
    const {
      code,
      message = "",
      data,
    }: IApiResult = await groupChatMemberApi(params);
    if (code === 200) {
      data.list.forEach((item: any) => {
        item.join_time =
          item.join_time &&
          this.$helper.dateFormat(item.join_time, "yyyy-MM-dd hh:mm");
        item.last_chat_at =
          item.last_chat_at &&
          this.$helper.dateFormat(item.last_chat_at, "yyyy-MM-dd hh:mm");
        item.pay_at =
          item.pay_at &&
          this.$helper.dateFormat(item.pay_at, "yyyy-MM-dd hh:mm");
      });
      this.list = this.list.concat(data.list);
      this.meta.total = data.total;
      this.loading = false;
      this.meta.page++;
      this.noMore = this.list.length >= data.total;
      if (
        data.total == this.list.length ||
        !data.list ||
        (data.list && data.list.length == 0)
      ) {
        this.noMore = true;
      }
    } else {
      this.loading = false;
      this.noMore = true;
      this.$toast.fail(message);
    }
  }
  // 页面滑到底部触发加载更多
  private getMoreOrder() {
    // 活动排行榜仅展示Top10,不做滚动加载
    if (this.currentTab != 4) {
      this.groupChatMember();
    }
  }

  // 点击tab切换
  private changeTab(name) {
    this.list = [];
    this.meta.page = 1;
    this.currentTab = name;
    this.groupChatMember();
  }
  // 打开添加的员工弹框
  private showPopup(item) {
    this.showStaffPopup = true;
    this.currentExternal = item.external_userid;
  }
  // 打开参与活动弹窗
  private showActivity(item) {
    this.showActivityPopup = true;
    this.union_id = item.unionid;
  }
}
</script>
<style lang="scss" scoped>
div {
  box-sizing: border-box;
}
.m-t-10 {
  margin-top: 10px;
}
::v-deep {
  .van-list__finished-text {
    text-align: center;
    font-size: 12px;
    color: #c8c8c8;
    line-height: 17px;
    padding-bottom: 20px;
    padding-top: 10px;
  }
  .van-loading__text {
    font-size: 12px;
    color: #c8c8c8;
  }
}

.portrait-association {
  background-color: #f5f5f5;
  overflow-x: hidden;
  .portrait-top {
    height: 80px;
    width: 100%;
    display: flex;
    align-items: center;
    padding: 20px 15px;
    background-color: #ffffff;
    .portrait-img {
      width: 40px;
      height: 40px;
      border-radius: 4px;
    }
    .portrait-name {
      font-size: 16px;
      color: #272727;
      width: 270px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      margin-left: 10px;
      margin-right: 16px;
      font-weight: 600;
    }
    .arrow-right {
      width: 6px;
      height: 11px;
    }
  }
  .border-line {
    width: calc(100vw - 15px);
    height: 1px;
    background: #ebebeb;
    position: relative;
    left: 15px;
  }
  .portrait-middle {
    padding: 19px 15px 5px 15px;
    margin-bottom: 6px;
    background-color: #ffffff;
    .portrait-item {
      display: flex;
      align-items: center;
      font-size: 14px;
      line-height: 20px;
      margin-bottom: 15px;
      .item-label {
        color: #808080;
        margin-right: 10px;
      }
      .item-content {
        color: #313131;
        .sex-item {
          display: flex;
          align-items: center;
          margin-right: 15px;
          .sex-icon {
            width: 14px;
            height: 14px;
            margin-right: 5px;
            display: block;
          }
          .sex-num {
            color: #808080;
          }
        }
      }
      .sex-box {
        display: flex;
        align-items: center;
      }
    }
    .statistics-box {
      &.m-b-15 {
        margin-bottom: 15px;
      }
      &.m-b-6 {
        margin-bottom: 6px;
      }
      .statistics-label {
        font-size: 14px;
        line-height: 20px;
        color: #808080;
        margin-bottom: 5px;
      }
      .statistics-card {
        background: #f4f5f8;
        height: 73px;
        width: 100%;
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 0 15px;
        &.vip {
          height: auto;
          padding: 11px 0 0 0;
          flex-wrap: wrap;
          justify-content: flex-start;
          .item-statistics {
            width: 33.33%;
            display: flex;
            justify-content: center;
            align-items: center;
            margin-bottom: 11px;
            .item-num {
              margin-bottom: 0;
              padding-left: 5px;
            }
            .item-title {
              font-size: 12px;
            }
          }
        }
        .item-statistics {
          text-align: center;
          .item-num {
            color: #313131;
            font-size: 16px;
            font-weight: 600;
            margin-bottom: 4px;
            .currency {
              font-size: 12px;
            }
          }
          .item-title {
            font-size: 12px;
            line-height: 17px;
            color: #999999;
          }
        }
      }
      .fold-btn {
        font-size: 12px;
        color: #999999;
        display: flex;
        align-items: center;
        height: 37px;
        width: 100%;
        justify-content: center;
        cursor: pointer;
        img {
          height: 15px;
          width: 15px;
          margin-left: 4px;
        }
      }
    }
  }
  .tab-box {
    background: #ffffff;
    height: 40px;
    width: 100%;
    color: #828282;
    font-size: 14px;
    padding: 0 15px;
    display: flex;
    align-items: center;
    border-bottom: solid 1px #ebebeb;
    .item-tab {
      margin-right: 30px;
      cursor: pointer;
    }
    .active {
      color: #000000;
      font-weight: 600;
    }
  }
  .list-box {
    background-color: #ffffff;
    padding: 0 15px 0 15px;
    overflow: hidden;
    .item-list {
      display: flex;
      margin-bottom: 20px;
      background: #ffffff;
      flex-shrink: 0;
      .item-index {
        font-size: 12px;
        padding: 0 12px;
        margin-bottom: 10px;
        display: flex;
        justify-content: center;
        align-items: center;
      }
      .item-img {
        width: 40px;
        height: 40px;
        img {
          width: 100%;
          height: 100%;
        }
      }
      .item-right {
        flex: 1;
        margin-left: 10px;
        .item-name-box {
          display: flex;
          justify-content: space-between;
          .item-name {
            font-size: 14px;
            line-height: 20px;
            font-weight: 600;
            color: #272727;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
          }
          .item-btn {
            font-size: 14px;
            line-height: 20px;
            color: #008fee;
            cursor: pointer;
          }
          .item-btn-else {
            font-size: 14px;
            line-height: 20px;
            color: #808080;
            cursor: pointer;
          }
        }
        .item-form {
          display: flex;
          align-items: center;
          margin-top: 7px;
          position: relative;
          .item-form-label {
            font-size: 12px;
            line-height: 17px;
            color: #808080;
            margin-right: 5px;
          }
          .item-form-content {
            font-size: 12px;
            line-height: 17px;
            color: #808080;
          }
          .item-type {
            position: absolute;
            right: 0;
            font-size: 12px;
            line-height: 17px;
            color: #808080;
          }
          .activity {
            font-size: 14px;
          }
        }
      }
    }
    .item-list:first-child {
      margin-top: 15px;
    }
    .no-more {
      text-align: center;
      font-size: 12px;
      color: #c8c8c8;
      line-height: 17px;
      padding-bottom: 20px;
      padding-top: 10px;
    }
  }
}
</style>
