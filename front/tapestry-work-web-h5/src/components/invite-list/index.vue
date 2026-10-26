<template>
  <div class="page-record-list" @click="goDetail($event)">
    <div class="status-box-left" @click="onTips($event)" :style="{backgroundImage:'url('+statusList[item.status].detail_status_bg+')'}">
      <!-- <img class="bg-icon" :src="statusList[item.status].detail_status_bg" /> -->
      <div class="txt status-txt">
        {{
          isSW
            ? statusList[item.status].sw_status_name
            : statusList[item.status].status_name
        }}
      </div>
      <img
        class="info-img"
        :style="{
          right: item.status === 4 ? '12rpx' : '18rpx',
        }"
        src="https://cdn-cos.coach.com.cn/uploads/2023/v3.4/info.png"
      />
    </div>

    <!-- 注意：代码不要删除，只是暂时隐藏邀请评价，因为第三方有这个功能 -->
    <div
      class="status-box-right"
      :style="{
        opacity: item.is_comment === 2 ? 1 : 0.4,
      }"
      v-if="isComment && singleChat"
    >
      <img
        class="edit-icon"
        src="https://cdn-cos.coach.com.cn/uploads/2022/v3.0.1/weclient/edit_icon.png"
      />
      <div class="txt">邀请评价</div>
    </div>
    <div
      class="status-box-right"
      :style="{
        opacity: item.is_comment === 2 ? 1 : 0.4,
      }"
      v-if="isComment && !singleChat"
    >
      <img
        class="edit-icon"
        src="https://cdn-cos.coach.com.cn/uploads/2022/v3.0.1/weclient/edit_icon.png"
      />
      <div @click="inviteComments($event, item)" class="txt">邀请评价</div>
    </div>

    <div class="order-box item-bottom m-t-45" v-if="item.status == 1">
      <!-- ks-顾客类型 -->
      <div class="order-item item-bottom" v-if="isKS">
        <div class="item-label">顾客类型</div>
        <div class="item-box" @click="radiochange($event)">
          <van-radio-group
            v-model="customer_status"
            icon-size="16"
            direction="horizontal"
            :disabled="statusDisabled"
            @change="onChangeStatus"
          >
            <van-radio :name="1" checked-color="#000">新客</van-radio>
            <van-radio :name="2" checked-color="#000">老客</van-radio>
          </van-radio-group>
        </div>
      </div>

      <!-- sw-顾客类型 -->
      <div class="order-item item-bottom" v-if="isSW">
        <div class="item-label">
          顾客类型
          <img
            class="tips-icon"
            src="https://cdn-cos.coach.com.cn/uploads/2022/v3.0.1/weclient/tips_icon.png"
            @click="onTips($event, 'customer_type')"
          />
        </div>
        <div class="item-value">
          {{ item.customer_status_str || "--" }}
        </div>
      </div>

      <!-- 是否到店购买 -->
      <div class="order-item item-bottom">
        <div class="item-label">是否到店购买</div>
        <div class="item-box" @click="radiochange($event)">
          <van-radio-group
            v-model="is_store"
            icon-size="16"
            direction="horizontal"
            :disabled="storeDisabled"
            @change="onChangeStore"
          >
            <van-radio :name="1" checked-color="#000">是</van-radio>
            <van-radio :name="2" checked-color="#000">否</van-radio>
          </van-radio-group>
        </div>
      </div>

      <!-- 订单日期 -->
      <div class="order-item item-bottom" v-if="isSW">
        <div class="item-label">订单日期</div>
        <div class="item-value">
          {{ filterTime(item.reserve_order_time) }}
        </div>
      </div>

      <div :class="['order-item', isSW ? 'item-bottom' : '']">
        <div class="item-label">
          订单金额
          <img
            class="tips-icon"
            src="https://cdn-cos.coach.com.cn/uploads/2022/v3.0.1/weclient/tips_icon.png"
            @click="onTips($event, 'reserve_order_price')"
          />
        </div>
        <div class="item-value">
          {{ formatMoney(item.reserve_order_price) }}
        </div>
      </div>

      <!-- 订单门店 -->
      <div class="order-item" v-if="isSW">
        <div class="item-label">订单门店</div>
        <div class="item-value">
          {{ item.store_name }}
        </div>
      </div>
    </div>

    <div
      class="m-t-15"
      :class="{
        'form-item': true,
      }"
    >
      <div class="label">邀约时间</div>
      <div class="value">{{ item.create_time_str }}</div>
    </div>
    <div class="form-item m-t-15">
      <div class="label">到店时间</div>
      <div class="value">{{ item.reserve_time_str }}</div>
    </div>
    <div class="form-item m-t-15">
      <div class="label">顾客姓名</div>
      <div class="value">{{ item.user_name }}</div>
    </div>
    <div class="form-item m-t-15">
      <div class="label">服务名称</div>
      <div class="value">{{ item.reserve_name }}</div>
    </div>
    <div class="form-item m-t-15">
      <div class="label">邀约导购</div>
      <div class="value">{{ item.member_name }}</div>
    </div>
    <div class="form-item m-t-15">
      <div class="label">邀约门店</div>
      <div class="value">{{ item.store_name }}</div>
    </div>
    <div class="form-item m-t-15">
      <div class="label">客户留言</div>
      <div class="value">{{ item.user_remark || "无" }}</div>
    </div>
    <!-- 提示文案弹窗 -->
    <van-popup v-model="toastOption.show" closeable class="popup-content">
      <div class="popup-content-title">{{ toastOption.title }}</div>
      <div class="popup-content-content">{{ toastOption.content }}</div>
    </van-popup>
  </div>
</template>

<script lang="ts">
// eslint-disable-next-line @typescript-eslint/no-unused-vars
import { Component, Vue, Prop, Watch } from "vue-property-decorator";
import { Toast, RadioGroup, Radio, Popup } from "vant";
import { StatusList } from "./constants";
import {
  changeStatus,
  commentShareApi,
  verifyShareApi,
} from "@/api/guide-invite/index";
import Helper from "@/utils/helper";
import store from "@/store";
@Component({
  components: {
    // Toast,
    [RadioGroup.name]: RadioGroup,
    [Radio.name]: Radio,
    [Popup.name]: Popup,
  },
})
export default class extends Vue {
  @Prop(Number) role!: number;
  // eslint-disable-next-line @typescript-eslint/ban-types
  @Prop(Object) item!: any;
  @Prop(Number) index!: number;

  /** 是否是sw */
  private isSW = Helper.isSW();
  /** 是否是ks */
  private isKS = Helper.isKS();
  customer_status = 0;
  is_store = 0;
  statusDisabled = false;
  storeDisabled = false;
  statusList = StatusList;
  isOpenTipsModal = false;
  isShowPriceTips = false;
  isComment = false;
  appid = "";
  /** 是否从侧边栏直接打开 */
  singleChat = false;
  /** 提示弹窗信息 */
  toastOption = {
    show: false,
    content: "",
    title: "",
  };
  @Watch("item", { deep: true, immediate: true })
  changeItem(value) {
    if (value) {
      this.isComment = this.item.status === 1 || this.item.status === 2;
      this.customer_status = this.item.customer_status;
      this.is_store = this.item.is_store;
      this.statusDisabled =
        this.item.customer_status === 1 || this.item.customer_status === 2;
      this.storeDisabled =
        this.item.customer_status === 1 || this.item.customer_status === 2;
    }
  }

  created() {
    // 获取小程序的appid
    this.appid = Helper.getAppid();
    console.log("appid", this.appid);
  }

  /** 获取用户是否从侧边栏打开 */
  judgeEntrance() {
    window.wx.invoke("getContext", {}, function (res: any) {
      console.log("getContext返回res,=======>", res);
      if (res.err_msg == "getContext:ok") {
        var entry = res.entry; // 返回进入小程序的入口类型
        if (entry === "single_chat_tools") {
          this.singleChat = true;
        }
      }
    });
  }
  moneyString(str) {
    str = str.toString();
    if (!str) return; //判空

    //切割数组
    //以.分隔为数组
    var arr = str.split(".");

    var left = arr[0];
    var right = arr.length > 1 ? "." + arr[1] : "";
    //取第一项用正则判断，每隔3位加，
    left = left.replace(new RegExp("\B(?=(\d{3})+$)", "g"), ",");
    return left + right;
  }
  formatMoney(money) {
    if (money) {
      return "¥" + this.moneyString(money / 100);
    }
    return "--";
  }

  /** 格式化时间 */
  private filterTime(time) {
    return this.$helper.dateFormat(time, "yyyy/MM/dd hh:mm");
  }

  // 状态信息弹窗
  onTips(e, type) {
    e.stopPropagation();

    let title = this.isSW
      ? this.statusList[this.item.status].sw_status_name
      : this.statusList[this.item.status].status_name;
    let message = this.isKS
      ? this.statusList[this.item.status].ks_desc
      : this.statusList[this.item.status].desc;
    // sw 顾客类型提示文案
    if (type === "customer_type") {
      title = "";
      message = `*客人会是VIP也会是新客，或者是VIP也是老客，沉睡客人中没有VIP 新客：该财年的第一次购买的客人 老客：该财年的客人去掉该财年的新客 VIP：在有效期内的VIP 沉睡客：非活跃客人（即排除了新客，老客和VIP 后其他的客人都属于沉睡客）`;
    }

    // sw/ks 订单金额提示文案
    if (type === "reserve_order_price") {
      title = "订单金额";
      message = "此处的订单金额为含税金额";
    }
    this.toastOption.content = message || "";
    this.toastOption.title = title || "";
    this.toastOption.show = true;
  }
  goDetail(e: any) {
    console.log("goDetail", e);
    if (
      e.target.className.indexOf("van-popup") > -1 ||
      e.target.className.indexOf("van-overlay") > -1
    ) {
      return;
    }
    e.stopPropagation();
    const _this = this;
    window.wx.invoke("launchMiniprogram", {
      appid: this.appid, // 需跳转的小程序appid
      path: `/pages/store-invite/invite-detail/invite-detail.html?id=${_this.item.id}&reserveNumber=${_this.item.reserve_number}&role=${this.role}&index=${this.index}`,
      function(res) {
        if (res.err_msg == "launchMiniprogram:ok") {
          // 正常
          console.log("跳转成功", res);
          console.log("acceptDataFromOpenedPage", res);
          // 状态与当前不一致
          if (_this.item.status !== res.status) {
            this.$emit("statusChange");
          }
        } else {
          console.log("跳转出错了", res);
        }
      },
    });
  }
  // radiochange父元素阻止冒泡
  radiochange(e) {
    e.stopPropagation();
  }
  // 操作 - 是否到店购买
  onChangeStore(value) {
    this.is_store = value;
    const params = {
      id: this.item.id,
      is_store: value,
      reserve_number: this.item.reserve_number,
    };
    changeStatus(params)
      .then((res: any) => {
        if (res.code === 200) {
          this.storeDisabled = true;
        }else{
          this.$toast(res.message);
        }
      })
      .catch((err) => {
        console.log(err);
      });
  }
  // 操作 - 顾客类型
  onChangeStatus(value) {
    this.customer_status = value;
    const params = {
      id: this.item.id,
      customer_status: value,
      reserve_number: this.item.reserve_number,
    };
    changeStatus(params)
      .then((res: any) => {
        if (res.code === 200) {
          this.statusDisabled = true;
        } else {
          this.$toast(res.message);
        }
      })
      .catch((err) => {
        console.log(err);
      });
  }

  /** 侧边栏打开返送邀请评价 */
  sidebarSend(id, reserve_number, title, url) {
    // 兼容不同品牌imgurl
    let imgUrl = "";
    if (this.isKS) {
      imgUrl =
        "https://cdn-cos.coach.com.cn/uploads/2022/5/27/ks-sharePic.png";
    } else if (this.isSW) {
      imgUrl =
        "https://cdn-cos.coach.com.cn/uploads/2022/5/27/sw-sharePic.png?";
    } else {
      imgUrl =
        "https://cdn-cos.coach.com.cn/uploads/2022/5/23/60f86724d69e71c10979b4e5bb17fa23.png";
    }
    console.log("imgUrl", imgUrl);
    const message = {
      msgtype: "miniprogram",
      enterChat: true,
      miniprogram: {
        appid: this.appid,
        title: title, // 小程序消息的title
        imgUrl: imgUrl, // 小程序消息的封面图。
        page: url, // 小程序消息打开后的路径，注意要以.html作为后缀，否则在微信端打开会提示找不到页面
      },
    };
    window.wx.invoke("sendChatMessage", message, function (res) {
      if (res.err_msg == "sendChatMessage:ok") {
        console.log("inviteComments", res, id, reserve_number);
        // 记录邀请评价日志
        const params = {
          id: id, // 主键id
          reserve_number: reserve_number, // 预约单号
        };
        commentShareApi(params);
      }
    });
  }

  /** 判断是否导购自己的邀请单 */
  verifyShareSend(id, reserve_number) {
    return new Promise((resolve) => {
      const params = {
        id,
        reserve_number,
      };
      verifyShareApi(params)
        .then((res: any) => {
          if (res.code !== 200) {
            resolve(false);
            this.$toast(res.message);
          }
          resolve(true);
        })
        .catch((e) => {
          resolve(false);
          console.error(e, "verifyShareSend 出错");
        });
    });
  }

  // 邀请顾客评价
  async inviteComments(e, item) {
    e.stopPropagation();
    const { id, reserve_number, is_comment } = item;
    const isOwn = await this.verifyShareSend(id, reserve_number);
    console.log("isOwn=================", isOwn);
    // 判断是否导购自己的单子
    if (!isOwn) return;
    // 1 未评价 2已评价
    console.log("is_comment=================", is_comment);
    if (is_comment === 1) return;
    let title = "邀约评价";
    let url = `pagesOrder/pages/store-invite/comment/comment.html?id=${id}&reserve_number=${reserve_number}`;
    this.sidebarSend(id, reserve_number, title, url);
  }
}
</script>
<style lang="sass" scoped>
.popup-content
      width: 240px
      padding: 30px 20px
      font-size: 14px
      text-align: center
      &-title
        font-size: 14px
        margin-bottom: 12px
      &-content
        font-size: 12px
        color: #96979E
      ::v-deep .van-popup__close-icon--top-right
          position: absolute
          top: 10px
          right: 10px
.page-record-list
    .m-t-15
        margin-top: 15px
    .m-t-45
        margin-top: 45px
    .form-item
        display: flex
        justify-content: space-between
        align-items: center
        font-size: 14px
        line-height: 20px
        .label
            color: #96979E
            flex-shrink: 0
            margin-right: 30px

    .status-box-left
        position: relative
        top: -8px
        left: -24px
        display: flex
        align-items: center
        width: 112px
        height: 32px
        padding-bottom: 4px
        box-sizing: border-box
        display: flex
        justify-content: flex-start
        align-items: center
        background-size: cover
        background-position: center
        background-repeat: no-repeat
        // padding-left: 24px
        .txt
            color: #FFFFFF
            text-align: center
            font-size: 14px
            margin-right: 5px
            white-space: nowrap
        .status-txt
            color: #000
        // .bg-icon
        //     width: 100%
        //     height: 32px
        .info-img
           width: 20px
           height: 20px
    .status-box-right
        cursor: pointer
        position: absolute
        top: 7px
        right: -6px
        width: 92px
        height: 32px
        display: flex
        align-items: center
        .txt
            position: absolute
            top: 0
            left: 0
            right: 0
            color: #000
            line-height: 33px
            text-align: center
            font-size: 14px
            background: transparent
        .edit-icon
            width: 16px
            height: 16px
    .order-box
        font-size: 14px
        padding: 15px
        background-color: #F3F4F6
        .order-item
            display: flex
            justify-content: space-between
            flex-direction: row
            .item-label
                flex-shrink: 0
                color: #96979E
                display: flex
                .tips-icon
                    width: 20px
                    height: 20px
                    margin-left: 5px
            .item-value
                color: #000
                word-break: break-all

            .item-box
                display: flex
                .radio-icon
                    width: 54px

                .radio-mr
                    margin-right: 24px
        .item-bottom
            margin-bottom: 15px
</style>
