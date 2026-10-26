<template>
  <div>
    <div class="page-create" v-if="!showLoading && !isError">
      <div class="content">
        <van-image
          class="top-img"
          v-if="serverIndex > -1"
          :src="serverList[serverIndex].reserve_img"
        />
        <div class="form-item flex-start-center has-line">
          <div class="label flex-start-center">
            <van-image
              class="icon"
              src="https://cdn-cos.coach.com.cn/uploads/2022/v3.0.1/weclient/icon-01.png"
            />
            <div>顾客姓名</div>
          </div>
          <div class="value font-500">
            <van-field
              class="input-box"
              v-model="userName"
              placeholder="请输入用户名称"
            />
          </div>
        </div>
        <div class="form-item flex-start-center has-line">
          <div class="label flex-start-center">
            <van-image
              class="icon"
              src="https://cdn-cos.coach.com.cn/uploads/2022/v3.0.1/weclient/icon-02.png"
            />
            <div>服务类型</div>
          </div>
          <div class="value" @click="showServicePicker = true">
            <div class="font-500" v-if="serverIndex > -1">
              {{ serverList[serverIndex].reserve_name }}
            </div>
            <div v-else class="placeholder">选择服务类型</div>
            <van-image
              class="right-icon"
              src="https://cdn-cos.coach.com.cn/uploads/2022/v3.0.1/weclient/right-icon.png"
            />
          </div>
        </div>
        <div
          class="form-item flex-start-center has-line"
          @click="handleSelectTime"
        >
          <div class="label flex-start-center">
            <van-image
              class="icon"
              src="https://cdn-cos.coach.com.cn/uploads/2022/v3.0.1/weclient/icon-03.png"
            />
            <div>到店时间</div>
          </div>
          <div class="value">
            <div class="font-500" v-if="selectTime.select_name">
              {{ selectTime.select_name }} {{ selectTime.name }}
            </div>
            <div class="placeholder" v-else>选择到店时间</div>
            <van-image
              class="right-icon"
              src="https://cdn-cos.coach.com.cn/uploads/2022/v3.0.1/weclient/right-icon.png"
            />
          </div>
        </div>
        <div class="form-item flex-start-center has-line">
          <div class="label flex-start-center">
            <van-image
              class="icon"
              src="https://cdn-cos.coach.com.cn/uploads/2022/v3.0.1/weclient/icon-09.png"
            />
            <div>邀约门店</div>
          </div>
          <div class="value font-500">
            <van-field
              class="input-box"
              v-model="storeName"
              placeholder="请输入门店名称"
            />
          </div>
        </div>
        <div class="form-item flex-start-center has-line">
          <div class="label flex-start-center">
            <van-image
              class="icon"
              src="https://cdn-cos.coach.com.cn/uploads/2022/v3.0.1/weclient/icon-06.png"
            />
            <div>门店地址</div>
          </div>
          <div class="value font-500">
            <van-field
              class="textarea"
              rows="1"
              type="textarea"
              autosize
              maxlength="45"
              v-model="storeAddress"
              placeholder="请输入门店地址"
            />
          </div>
        </div>
        <div class="form-item flex-start-center">
          <div class="label flex-start-center">
            <van-image
              class="icon"
              src="https://cdn-cos.coach.com.cn/uploads/2022/v3.0.1/weclient/icon-10.png"
            />
            <div>导购备注</div>
          </div>
        </div>
        <div class="textarea-box" @click.prevent="inputFocus('remarkTxt')">
          <van-field
            ref="remarkTxt"
            class="textarea"
            type="textarea"
            rows="5"
            autosize
            maxlength="50"
            v-model="reserveRemark"
            placeholder="请输入备注信息"
          />
        </div>
      </div>
      <div class="footer-btn">
        <div @click="goRecordList" class="btn btn-default">查看邀约记录</div>
        <div @click="create" class="btn" v-throttle>创建并发送</div>
      </div>
    </div>
    <!-- 顾客类型picker -->
    <van-popup v-model="showServicePicker" round position="bottom">
      <van-picker
        show-toolbar
        @confirm="changeServer"
        @cancel="showServicePicker = false"
        value-key="reserve_name"
        :default-index="serverIndex"
        :columns="serverList"
      >
      </van-picker>
    </van-popup>
    <!-- 到店时间picker -->
    <invited-Date
      :allDay="allDay"
      :timeObj="timeObj"
      :timeSlot="timeSlot"
      v-if="showInvitedDate"
      @close="closeInvitedDate"
      @sumbit="sumbitInvitedDate"
    ></invited-Date>
    <!-- 微信打开无权限 -->
    <empty-page
      v-if="!showLoading && isError"
      :picture="'emptyList'"
      :text="errorTxt"
    ></empty-page>
  </div>
</template>

<script lang="ts">
// eslint-disable-next-line @typescript-eslint/no-unused-vars
import { Component, Vue } from "vue-property-decorator";
import invitedDate from "@/components/invited-date/index.vue";
import { Image, Field, Picker, Popup } from "vant";
import emptyPage from "@/components/base/empty-page.vue";
import Helper from "@/utils/helper";
import {
  guideTimeApi,
  serviceListApi,
  userInfoApi,
  createOrderApi,
  storeInfoApi,
  orderShareNumApi,
  deleteOrderApi,
  verifyTimeApi,
  showOrderApi,
} from "@/api/guide-invite/index";

/** 获取配置图片 */
function getSharePic() {
  const time = Date.now();
  if (Helper.isCoach()) {
    return `https://cdn-cos.coach.com.cn/uploads/2022/5/23/60f86724d69e71c10979b4e5bb17fa23.png?time=${time}`;
  }

  if (Helper.isKS()) {
    return `https://cdn-cos.coach.com.cn/uploads/2022/5/27/ks-sharePic.png?time=${time}`;
  }

  return `https://cdn-cos.coach.com.cn/uploads/2022/5/27/sw-sharePic.png?time=${time}`;
}

@Component({
  components: {
    invitedDate,
    [Field.name]: Field,
    [Popup.name]: Popup,
    [Image.name]: Image,
    [Picker.name]: Picker,
    emptyPage,
  },
})
export default class extends Vue {
  showInvitedDate = false;
  showServicePicker = false;
  dateModel = false;
  // 服务选择的下标
  serverIndex = -1;
  timeObj = {};
  showTime = false;
  // 创建按钮是否可以点击，防止多次点击
  canClick = false;
  // 用户的信息
  userInfo: any = {};
  // 选择时间
  selectTime: any = {};
  serverList = [];
  storeName = "";
  storeAddress = "";
  reserveRemark = "";
  userName = "";
  storeInfo: any = {};
  showModel = false;
  isSingle = false;
  modelTitle = "提醒";
  description = "该时段已有顾客预约服务，";
  descriptionTwo = "建议重新选择时段";
  showLoading = true;
  orderInfo = {};
  isError = false;
  errorTxt = "";
  allDay = 30;
  appid = "";
  timeSlot = []
  getEnterpriseHash() {
    return this.$store.getters["common/getEnterpriseHash"];
  }
  created() {
    // 获取小程序的appid
    this.appid = Helper.getAppid();
    console.log("this.appid", this.appid);
    // 初始化数据
    this.init();
  }
  handleSelectTime() {
    console.log("点击了打开弹窗");
    // todo: 判断是否选择了服务类型
    if (this.serverIndex === -1) {
      this.$toast("请选择服务类型");
      return;
    }
    this.showInvitedDate = true;
  }
  /**
   * @description: input框聚焦
   */
  inputFocus(refName) {
    (this.$refs[refName] as HTMLElement)?.focus();
  }
  /**
   * @description:初始化数据
   */
  init() {
    this.$toast.loading({ forbidClick: true, message: "加载中..." });
    let p = [
      // 获取服务列表
      this.getServerList(),
      // 获取门店
      this.getStoreInfo(),
      // 获取一个月之内时间段预约的人数
      this.getTime(),
    ];
    Promise.all(p)
      .then(() => {
        console.log("Promise.all成功");
        this.showLoading = false;
        // 获取用户userId 然后获得当前用户的基本信息
        this.getUserId();
        // this.getUserInfo("wm9hx_EQAAGBQYNvdhy_btZBC-m65d7g");
      })
      .catch((err) => {
        this.showLoading = false;
        console.log(err, "出错");
      }).finally(() => {
        this.$toast.clear();
      });
  }
  /**
   * @description:服务类型切换
   */
  changeServer(item, index) {
    console.log(item, index);
    let allDay = 30;
    let serverList = this.serverList;
    if (index !== -1) {
      const item = serverList[index];
      if (item.service_type === 0) {
        const nowDate = this.getTodayTime(null);
        const endTime = this.getTodayTime(item.end_time * 1000);
        allDay = Math.min(
          Math.floor((endTime - nowDate) / (24 * 60 * 60 * 1000)),
          30
        );
      }
    }
    console.log(allDay, "allDay");
    this.serverIndex = Number(index);
    this.allDay = allDay;
    this.selectTime = {};
    this.showServicePicker = false;
    this.timeSlot = item.time_slot
  }
  /**
   * @description: 选择到店时间
   */
  sumbitInvitedDate(e) {
    console.log(e, "选择到店时间");
    // 如果不是coach需要验证到店时间
    if (Helper.isCoach() === false) {
      this.verifyTime(e);
    } else {
      this.selectTime = e;
      this.showInvitedDate = false;
    }
  }
  /**
   * @description: 验证到店时间
   */
  verifyTime(timeStamp) {
    const params = {
      reserve_time: timeStamp.startTime, // 到店时间戳
      store_id: this.storeInfo.department_third_id, // 门店id
    };
    verifyTimeApi(params)
      .then((res) => {
        console.log(res, "验证到店时间");
        const { check, msg } = res.data;
        if (check) {
          this.selectTime = timeStamp;
        } else {
          this.$toast(msg);
        }
        this.showInvitedDate = false;
      })
      .catch((err) => {
        console.log(err, "验证到店时间出错");
        this.$toast(err.msg);
      });
  }
  /**
   * @description: 关闭到店时间弹窗
   */
  closeInvitedDate() {
    console.log(222222);
    this.showInvitedDate = false;
  }
  /**
   * @description: 跳转到邀约记录小程序
   */
  goRecordList() {
    const { appid } = this;
    const self = this
    console.log("跳转到邀约记录小程序");
    
    window.wx.invoke("launchMiniprogram", {
      appid: appid, // 需跳转的小程序appid
      path: `pages/store-invite/record-list/record-list.html?union_id=${this.userInfo.union_id}`, // 所需跳转的小程序内页面路径及参数。非必填 },
      function(res) {
        console.log(res.err_msg);
        
        if (res.err_msg == "launchMiniprogram:ok") {
          // 正常
          console.log("跳转成功", res);
        } else {
          console.log("跳转出错了", res);
        }
      },
    });
    // 神策数据上报 -- 查看邀约记录
    self.$track.trackViewAppointmentHistory({
      entrance: 3
    })
  }
  /**
   * @description:获取到店时间，预约人数
   */
  getTime() {
    return new Promise((resolve, reject) => {
      guideTimeApi()
        .then((res: any) => {
          if (res.code === 200) {
            this.timeObj = res.data;
            this.showTime = true;
            resolve(res);
          }
          reject(res.message);
        })
        .catch((err) => {
          reject(err);
        });
    });
  }
  /**
   * @description:获取门店信息
   */
  getStoreInfo() {
    return new Promise((resolve, reject) => {
      storeInfoApi()
        .then((res: any) => {
          if (res.code === 200) {
            this.storeInfo = res.data;
            this.storeName = res.data.department_name;
            this.storeAddress = res.data.address;

            resolve(res);
          }
          reject(res.message);
        })
        .catch((err) => {
          reject(err);
        });
    });
  }
  /**
   * @description:获取用户ID
   */
  getUserId() {
    const _this = this;
    window.wx.invoke("getContext", {}, function (res: any) {
      console.log("getContext返回res,=======>", res);
      if (res.err_msg == "getContext:ok") {
        const entry = res.entry;
        // 从单聊会话的工具栏进入-打开客户画像-客户页面
        if (entry == "single_chat_tools") {
          window.wx.invoke("getCurExternalContact", {}, (res) => {
            if (res.err_msg == "getCurExternalContact:ok") {
              console.log("外部联系人返回res,=======>", res);
              _this.getUserInfo(res.userId); // 返回当前外部联系人userId
              // 神策数据上报 -- 点击侧边栏
              _this.$track.trackClickItem({
                feature: _this.$route.path,
                customer_id: res.userId
              })
            } else {
              console.error(res, "getCurExternalContact===");
            }
          });
        } else {
          _this.isError = true;
          _this.showLoading = false;
          _this.errorTxt = "通过私聊窗口发送";
          console.log("wx.qy.getCurExternalContact==", res);
        }
      } else {
        console.error("wx.qy.getContext 异常", res);
        console.log(_this, "xxxxxx");
        // _this.getUserInfo('wmvjp0DgAA-vjnccwO0FzDLkE2WZZvdA')
        _this.isError = true;
        _this.showLoading = false;
        _this.errorTxt = "获取用户ID失败";
      }
    });
  }
  /**
   * @description:获取用户信息
   */
  async getUserInfo(external_userid) {
    try {
      const { code, data, message }: any = await userInfoApi({
        external_userid: external_userid,
      });
      if (code === 200) {
        this.userInfo = data;
        this.userName = data.nick;
      } else {
        this.$toast(message);
      }
    } catch (error) {
      console.error(error, "获取用户userInfo错误");
    }
  }
  /**
   * @description:获取服务列表
   */
  getServerList() {
    const params = {
      page: 1,
      size: 9999,
    };
    return new Promise((resolve, reject) => {
      serviceListApi(params)
        .then((res: any) => {
          if (res.code === 200) {
            const serverList = res.data.list.filter((i) => {
              if (i.service_type === 0) {
                const nowDate = this.getTodayTime("");
                const startTime = this.getTodayTime(i.start_time * 1000);

                if (startTime > nowDate) {
                  return false;
                }
              }
              return true;
            });
            this.serverList = serverList;
            resolve(res.message);
          }
          reject(res.message);
        })
        .catch((err) => {
          reject(err);
        });
    });
  }
  /**
   * @description:发送邀约单
   */
  sendMessage(param) {
    const { orderInfo, serverList, serverIndex, storeInfo }: any = this;
    const self = this;
    console.log(this.appid, "appid");
    const message = {
      // 向当前会话发送消息
      msgtype: "miniprogram", // 消息类型，必填
      enterChat: true,
      miniprogram: {
        appid: this.appid,
        title: "到店邀约", // 小程序消息的title
        imgUrl: getSharePic(),
        page: `pagesOrder/pages/store-invite/share-detail/share-detail.html?id=${orderInfo.id}&reserve_number=${orderInfo.reserve_number}`, // 小程序消息打开后的路径，注意要以.html作为后缀，否则在微信端打开会提示找不到页面
      },
    };
    window.wx.invoke("sendChatMessage", message, function (res) {
      if (res.err_msg == "sendChatMessage:ok") {
        console.log(res, "++++");
        self.canClick = false;
        // 发送成功后清空信息
        self.init();
        self.getUserId();
        setTimeout(() => {
          self.serverIndex = -1;
          self.selectTime = {};
          self.reserveRemark = ""; // 备注
          self.storeAddress = "";
        }, 100);
        const params = {
          id: orderInfo.id,
          reserve_number: orderInfo.reserve_number,
        };
        orderShareNumApi(params);
        showOrderApi(params);
        // 神策数据上报 -- 发送到店邀约
        self.$track.trackAppointmentSend({
          reserve_number: orderInfo.reserve_number,
          user_external_id: param.user_external_id,
          reserve_time: param.reserve_time,
          time_period: param.reserve_time,
          reserve_name: serverList[serverIndex].reserve_name,
          store_name: param.store_name,
          store_code: storeInfo.store_code || param.store_name,
        })
      } else {
        console.log(res, "-----error");
        self.canClick = false;
        // self.deleteOrder();
      }
    });
  }
  /**
   * @description: 删除邀约记录
   */
  deleteOrder() {
    const { orderInfo }: any = this;
    const params = {
      id: orderInfo.id,
      reserve_number: orderInfo.reserve_number,
    };
    deleteOrderApi(params)
      .then((res: any) => {
        if (res.code === 200) {
          console.log("删除邀约记录成功");
          this.canClick = false;
        }
      })
      .catch((err) => {
        console.error(err, "删除邀约记录错误");
      });
  }
  /**
   * @description: 创建预约单
   */
  create() {
    const params = this.getParams();
    if (typeof params === "string") {
      this.$toast(params);
      return;
    }
    if (this.canClick) return;

    this.canClick = true;

    this.$toast.loading({
      message: "加载中...",
      duration: 0,
    });
    createOrderApi(params)
      .then((res: any) => {
        this.$toast.clear();
        if (res.code === 200) {
          this.orderInfo = res.data;

          this.sendMessage(params);
          return;
        }

        this.canClick = false;
        this.$toast(res.message);
      })
      .catch((err) => {
        this.canClick = false;
        console.error(err, "=+++");
        this.$toast.clear();
      });
  }
  /**
   * @description: 获取当日00:00:00时间戳
   */
  getTodayTime(target) {
    let time = 0;

    if (target) {
      time = new Date(
        new Date(
          new Date(new Date(target).setHours(0)).setMinutes(0)
        ).setSeconds(0)
      ).setMilliseconds(0);
    } else {
      time = new Date(
        new Date(new Date(new Date().setHours(0)).setMinutes(0)).setSeconds(0)
      ).setMilliseconds(0);
    }

    return time;
  }
  /**
   * @description: 获取参数
   */
  getParams() {
    const {
      serverIndex,
      userInfo, // 用户的信息
      selectTime, // 选择时间
      serverList,
      storeName,
      storeAddress,
      reserveRemark,
      storeInfo,
      userName,
    } = this;
    if (!userName) {
      return "请填写用户名称";
    }
    if (serverIndex === -1) {
      return "请选择服务类型";
    }
    if (!selectTime.startTime) {
      return "请选择到店时间";
    }
    if (!storeName) {
      return "请填写门店名称";
    }
    // if (!storeAddress) {
    //   return '请填写门店地址'
    // }
    const params = {
      reserve_id: serverList[serverIndex].id,
      reserve_time: selectTime.startTime,
      reserve_end_time: selectTime.endTime,
      store_id: storeInfo.department_third_id,
      store_name: storeName,
      store_address: storeAddress.slice(0, 45),
      user_name: userName,
      user_unionid: userInfo.union_id,
      user_external_id: userInfo.external_userid,
      reserve_remark: reserveRemark,
    };
    return params;
  }
}
</script>
<style lang="sass" scoped>
.page-create
  padding: 15px 15px 46px
  .top-img
    width: 315px
    height: 107px
    background: #F7F7F7
  .flex-start-center
    display: flex
    justify-content: flex-start
    align-items: center
  .font-500
    font-family: PingFangSC-Medium, PingFang SC
    font-weight: 500
  .content
    width: 345px
    padding: 10px 15px
    background: white
    margin-bottom: 138px
    .van-cell
      padding: 0
      background: none
    .form-item
      padding: 15px 0
      line-height: 20px
      font-size: 14px
      position: relative
      &.has-line
        &::after
          content: ""
          clear: both
          width: 314px
          height: 1px
          background: #F0F0F0
          position: absolute
          bottom: 0
          left: 0
      .picker-box
        position: absolute
        top: 0
        left: 0
        right: 0
        bottom: 0
        .picker
          height: 48px
          opacity: 0
      .label
        flex-shrink: 0
        color: #96979E
        margin-right: 32px
        .icon
          width: 20px
          height: 20px
          margin-right: 6px
      .value
        position: relative
        color: #4A4A4A
        width: 100%
        .placeholder
          font-weight: 400
          color: #A6A6A6
        .right-icon
          position: absolute
          top: 50%
          right: 5px
          transform: translateY(-50%)
          width: 9px
          height: 14px
        .textarea
          width:201px !important
          font-size: 14px
          line-height: 14px !important
    .textarea-box
      width: 315px
      height: 119px
      background: #F6F5F5
      padding: 10px
      .textarea
        font-size: 14px
        line-height: 20px
        width: 100%
        height: 119px
      .placeholder-textarea
        color: #BDBEC8
  .footer-btn
    background: #FFFFFF
    box-shadow: 0px -1px 4px 0px rgba(0,0,0,0.1)
    padding: 10px  30px 30px
    // background: #F7F7F7
    position: fixed
    bottom: 0
    left: 0
    display: flex
    justify-content: space-between
    width: 100%
    .btn
      width: 150px
      line-height: 44px
      background: #4A4A4A
      text-align: center
      font-size: 14px
      color: #FFFFFF
    .btn-default
      color: #96979E
      background: transparent
      border: 1px solid #96979E
  .empty
    display: flex
    flex-direction: column
    align-items: center
    justify-content: center
    height: 200px
    &-icon
      margin-top: 10px
      width: 66px
      height: 66px
</style>
