<template>
  <div>
    <div class="page-invite-config flex-center" v-if="!showLoading && !isError">
      <div>
        <div class="page-invite-config-header">
          <div class="page-invite-config-header-name flex">
            <img
              class="letter-icon"
              src="https://cdn-cos.coach.com.cn/uploads/2022/v3.2/weclient/letter-icon.png"
            />
            <div class="font-14 label">邀请函名称：</div>
            <div v-if="id">{{ detail.title_name }}</div>
            <div class="page-invite-config-header-name-choose" v-else>
              <div
                class="page-invite-config-header-name-choose-content flex"
                @click="handleOpenOptions"
              >
                <div class="page-invite-config-header-name-choose-content-name">
                  {{ optionsChoose.title_name }}
                </div>
                <img
                  class="page-invite-config-header-name-choose-content-icon"
                  src="https://cdn-cos.coach.com.cn/uploads/2023/v3.2.0.1/updown.png"
                />
              </div>
              <div
                class="page-invite-config-header-name-choose-content-options"
                v-if="optionsShow"
              >
                <div class="scroll-view">
                  <div v-for="(item, index) in options" :key="index">
                    <div
                      :class="{ active: optionsChoose.id === item.id }"
                      class="page-invite-config-header-name-choose-content-options-item"
                      @click="handleChangeOptions(item)"
                    >
                      {{ item.title_name }}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="service-relation flex">
            <div class="font-14">关联预约服务：</div>
            <div class="font-14">
              {{
                detail.reserve_service_name ||
                optionsChoose.reserve_service_name
              }}
            </div>
          </div>
        </div>
        <div class="page-invite-config-tabs flex-center">
          <div
            :class="{ active: styleIndex === 1 }"
            @click="handleChangeStyle(1)"
          >
            样式一
          </div>
          <div
            :class="{ active: styleIndex === 2 }"
            class="second"
            @click="handleChangeStyle(2)"
          >
            样式二
          </div>
        </div>
        <div class="page-invite-config-content">
          <div v-if="styleIndex === 1">
            <div class="page-invite-config-content-title">
              尊敬的顾客：{{ userName || "--" }}
            </div>
            <!-- sa编辑内容 -->
            <div class="edit-box">
              <van-field
                type="textarea"
                class="edit"
                maxlength="150"
                autosize
                show-word-limit
                v-model="editValue"
              />
            </div>
            <div class="page-invite-config-content-sign">
              您的专属顾问：{{ member_name }}
            </div>
          </div>
          <div class="page-invite-config-content-backgroundImg">
            <img
              class="page-invite-config-content-backgroundImg-image"
              :src="
                id
                  ? detail.background_url[currentIndex]
                  : optionsChoose.background_url[currentIndex]
              "
            />
          </div>
          <div v-if="styleIndex === 2" class="m-t-20">
            <div
              class="page-invite-config-content-title"
              style="margin-top: 16px"
            >
              尊敬的顾客：{{ userName || "--" }}
            </div>
            <!-- sa编辑内容 -->
            <div class="edit-box">
              <van-field
                v-model="editValue"
                type="textarea"
                class="edit"
                maxlength="150"
                autosize
                show-word-limit
              />
            </div>
            <div class="page-invite-config-content-sign m-b-80">
              您的专属顾问：{{ member_name }}
            </div>
          </div>
          <!-- 自定义指示点 -->
          <div class="dot-container flex-center">
            <img
              :class="{ disabled: currentIndex === 0 }"
              class="dot-container-icon-left icon"
              src="https://cdn-cos.coach.com.cn/uploads/2023/v3.2.0.1/icon_right_light.png"
              @click="handleChangePic(1)"
            />
            <div class="dot-container-tag">背景{{ chinese[currentIndex] }}</div>
            <img
              class="dot-container-icon-right icon"
              src="https://cdn-cos.coach.com.cn/uploads/2023/v3.2.0.1/icon_right_light.png"
              @click="handleChangePic(2)"
            />
          </div>
        </div>
        <div class="footer" v-if="!id && !hiddenButton">
          <button
            class="footer-sendBtn font-12"
            @click="sendingInvitation"
            v-throttle
          >
            发送邀请函
          </button>
        </div>
      </div>
    </div>
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
import { Image, Field, Picker, Popup } from "vant";
import emptyPage from "@/components/base/empty-page.vue";
import Helper from "@/utils/helper";
import {
  getInviteConfigList,
  sendInvite,
  getInvitationInfo,
  recordInfo,
  userInfoApi,
  deleteInvited,
} from "@/api/invitation-letter/index";
import { storeInfoApi } from "@/api/guide-invite/index";
Vue.use(Field);
@Component({
  components: {
    emptyPage,
    [Popup.name]: Popup,
    [Image.name]: Image,
    [Picker.name]: Picker,
  },
})
export default class extends Vue {
  //样式
  styleIndex = 1;
  userName = "";
  // 创建按钮是否可以点击，防止多次点击
  canClick = false;
  //背景下标
  currentIndex = 0;
  chinese = ["一", "二", "三", "四", "五", "六", "七", "八", "九", "十"];
  //邀请函名称选择
  optionsChoose: any = {};
  //选择邀请函弹窗
  optionsShow = false;
  //是否出错
  isError = false;
  //是否加载
  showLoading = true;
  id = "";
  //邀请函模板详情
  detail: any = {};
  //是否显示邀请函模板详情
  invite_record_id = "";
  //专属顾问
  member_name = "--";
  options: any = [];
  appid = "";
  wxwork = true;
  errorTxt = "获取邀请函配置信息失败";
  //邀请函内容 */
  editValue = "";
  userInfo: any = {};
  //是否隐藏button 如果是从客户端直接重定向过来的，需要隐藏
  hiddenButton: any = false;
  created() {
    // 获取小程序appid
    this.appid = Helper.getAppid();
    console.log("this.appid", this.appid);
    this.id = this.getQueryString("id");
    this.invite_record_id = this.getQueryString("invite_record_id");
    const isWorkWx = this.envjudge();
    console.log("isWorkWx=====" + isWorkWx);
    this.wxwork = isWorkWx == "com-wx-mobile" || isWorkWx == "com-wx-pc";
    this.hiddenButton = !!(this.wxwork && this.invite_record_id);
    console.log("hiddenButton=====" + this.hiddenButton);
  }
  mounted() {
    // 获取导购名片信息
    this.getGuideInfo();
    if (this.id) {
      this.getInvitationInfo();
    } else {
      // 存在ID的时候是后台预览，不需要请求顾客信息
      if (!this.hiddenButton) {
        this.getUserId();
        // this.getUserInfo("wm9hx_EQAAGBQYNvdhy_btZBC-m65d7g");
        this.getInvitation();
      } else {
        // 隐藏返回按钮
        // wx.hideHomeButton();
        this.getRecordInfo();
      }
    }
  }
  /**
   * @description: 判断h5环境
   */

  envjudge() {
    var isMobile = window.navigator.userAgent.match(
      /(phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone)/i
    ); // 是否手机端
    var isWx = /micromessenger/i.test(navigator.userAgent); // 是否微信
    var isComWx = /wxwork/i.test(navigator.userAgent); // 是否企业微信

    if (isComWx && isMobile) {
      //手机端企业微信
      return "com-wx-mobile";
    } else if (isComWx && !isMobile) {
      //PC端企业微信
      return "com-wx-pc";
    } else if (isWx && isMobile) {
      // 手机端微信
      return "wx-mobile";
    } else if (isWx && !isMobile) {
      // PC端微信
      return "wx-pc";
    } else {
      return "other";
    }
  }
  /**
   * @description: 获取地址栏参数
   */
  getQueryString(name: string) {
    const reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)", "i");
    const r = window.location.search.substr(1).match(reg);
    if (r != null) {
      return unescape(r[2]);
    }
    return null;
  }
  /**
   * @description: 获取导购信息
   */
  async getGuideInfo() {
    const res: any = await storeInfoApi();
    if (res.code === 200) {
      this.member_name = res.data.members_name;
    }
  }
  /**
   * @description:  获取用户ID
   */
  getUserId() {
    const _this = this;
    try {
      window.wx.invoke("getContext", {}, function (res: any) {
        console.log("getContext返回res,=======>", res);
        if (res.err_msg == "getContext:ok") {
          // 返回进入小程序的入口类型
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
          console.error("wx.qy.getContext 异常");
          console.log(_this, "xxxxxx");
          // _this.getUserInfo('wmvjp0DgAA-vjnccwO0FzDLkE2WZZvdA')
          _this.isError = true;
          _this.showLoading = false;
          _this.errorTxt = "获取用户ID失败";
        }
      });
    } catch (error) {
      this.isError = true;
      this.showLoading = false;
      this.errorTxt = "获取用户ID失败";
    }
  }
  /**
   * @description:  获取用户信息
   */
  async getUserInfo(external_userid) {
    try {
      const { code, data, message }: any = await userInfoApi({
        external_userid,
      });
      if (code === 200) {
        this.userInfo = data;
        this.userName = data.name || data.nick;
      } else {
        this.$toast(message);
      }
    } catch (error) {
      console.error(error, "获取用户userInfo错误");
    }
  }
  /**
   * @description:  获取邀请函模板详情
   */
  async getInvitationInfo() {
    this.$toast.loading({
      message: "加载中...",
      duration: 0,
    });
    try {
      let params = {
        id: Number(this.id),
      };
      let res: any = await getInvitationInfo(params);
      if (res.code === 200) {
        let { data } = res;
        try {
          data.share_url = JSON.parse(data.share_url);
          data.background_url = JSON.parse(data.background_url);
        } catch (error) {}

        this.detail = data;
        this.showLoading = false;
        this.editValue = data.content;
      } else {
        this.$toast(res.message);
      }
    } catch (error) {
      console.log(error);
    } finally {
      this.$toast.clear();
    }
  }
  /**
   * @description:  获取邀请函内容
   */
  async getInvitation() {
    this.$toast.loading({
      message: "加载中...",
      duration: 0,
    });
    try {
      const res: any = await getInviteConfigList();
      if (res.code === 200) {
        let optionsChoose = res.data[0];
        optionsChoose.share_url = JSON.parse(optionsChoose.share_url);
        optionsChoose.background_url = JSON.parse(optionsChoose.background_url);
        this.options = res.data;
        this.optionsChoose = optionsChoose;
        this.editValue = optionsChoose.content;
      } else {
        this.$toast(res.message);
      }
      this.$toast.clear();
      this.showLoading = false;
    } catch (error) {
      this.isError = true;
      this.showLoading = false;
      this.$toast.clear();
    }
  }
  /**
   * @description:客户端跳转到导购端重定向 获取邀请函信息
   */
  getRecordInfo() {
    recordInfo({ invite_record_id: Number(this.invite_record_id) })
      .then((res: any) => {
        if (res.code !== 200) return this.$toast(res.message);
        let optionsChoose = {
          ...res.data,
          background_url: [res.data.background_url],
          content: res.data.send_content,
          share_url: JSON.parse(res.data.share_url),
        };

        this.isError = false;
        this.showLoading = false;
        this.currentIndex = 0;
        this.styleIndex = optionsChoose.style_type;
        this.options = [optionsChoose];
        this.optionsChoose = optionsChoose;
        this.editValue = optionsChoose.content;
      })
      .catch((e) => console.error("获取接口失败"));
  }
  /**
   * @description:切换邀请函
   */
  handleChangeOptions(item: any) {
    try {
      item.share_url = JSON.parse(item.share_url);
      item.background_url = JSON.parse(item.background_url);
    } catch (error) {}

    this.optionsChoose = item;
    this.editValue = item.content;
    this.optionsShow = false;
  }
  /**
   * @description: 打开切换邀请函下拉
   */
  handleOpenOptions() {
    this.optionsShow = !this.optionsShow;
  }
  /**
   * @description:切换样式
   */
  handleChangeStyle(index: number) {
    this.styleIndex = index;
  }
  /**
   * @description:切换背景图
   */
  handleChangePic(index: number) {
    console.log(index, "index");
    if (index == 1) {
      if (this.currentIndex <= 0) return;
      --this.currentIndex;
    } else {
      if (this.currentIndex >= this.optionsChoose.background_url.length - 1) {
        return;
      }
      ++this.currentIndex;
    }
    this.currentIndex = this.currentIndex;
  }
  /**
   * @description:发送邀请函
   */
  sendingInvitation() {
    const { optionsChoose, userInfo, styleIndex, currentIndex, editValue } =
      this;
    const params = {
      invite_template_id: optionsChoose.id,
      user_openid: userInfo.external_userid,
      style_type: styleIndex,
      background_url: optionsChoose.background_url[currentIndex],
      send_content: editValue,
    };
    if (this.canClick) return;
    this.canClick = true;

    sendInvite(params)
      .then((res: any) => {
        if (res.code === 200) {
          this.sendMessage(res.data.invite_record_id);
        } else {
          this.$toast(res.message);
          this.canClick = false;
        }
      })
      .catch((error) => {
        this.canClick = false;
        console.log(error);
      });
  }
  /**
   * @description:发送
   */
  sendMessage(invite_record_id: String) {
    const { appid, optionsChoose, userInfo } = this;
    const self = this;
    const message = {
      msgtype: "miniprogram", // 消息类型，必填
      enterChat: true,
      miniprogram: {
        // 小程序的appid
        appid,
        // 小程序消息的title
        title: optionsChoose.share_title,
        imgUrl: `${optionsChoose.share_url[0]}?imageMogr2/thumbnail/!50p`,
        // 小程序消息打开后的路径，注意要以.html作为后缀，否则在微信端打开会提示找不到页面
        page: `pagesOrder/pages/invite-check/invite-check.html?id=${invite_record_id}`,
      },
    };
    window.wx.invoke("sendChatMessage", message, function (res) {
      console.log(res, "res");
      if (res.err_msg == "sendChatMessage:ok") {
        console.log(res, "++++");
        self.canClick = false;
        // 神策数据上报 -- 发送邀请函
        self.$track.trackInvitationSend({
          invite_record_id: invite_record_id,
          invite_template_id: optionsChoose.id,
          user_external_id: userInfo.external_userid,
          entrance: 1
        })
      } else {
        console.log(res, "-----error");
        self.deleteInvitedRecord(invite_record_id);
      }
    });
  }
  /**
   * @description: 删除邀请记录
   */
  deleteInvitedRecord(invite_record_id: String) {
    const params = {
      invite_record_id: invite_record_id,
    };
    deleteInvited(params)
      .then((res: any) => {
        if (res.code === 200) {
          console.log("删除邀请记录成功");
          this.canClick = false;
        }
      })
      .catch((err) => {
        console.error(err, "删除邀请记录错误");
      });
  }
}
</script>
<style lang="sass" scoped>
.page-invite-config
    width: 100%
    // height: 100vh
    flex-direction: column
    line-height: 20px
    .font-12
        font-size: 12px
    .font-14
        font-size: 14px
    &-header
        width: 100%
        height: 91px
        padding: 15px
        background-color: #FFFFFF

        &-name
            justify-content: space-between
            align-items: center
            color: #000000
            font-size: 14px

            .letter-icon
                width: 15px
                height: 15px
                margin-right: 5px

            .label
                flex: 1

            &-choose
                width: 173px
                height: 28px
                border: 1px solid #A6A6A6
                padding: 4px 10px
                position: relative

                &-content
                    width: 100%
                    height: 100%
                    align-items: center

                    &-name
                        flex: 1
                        overflow: hidden
                        white-space: nowrap
                        text-overflow: ellipsis
                        display: flex
                        justify-content: center

                    &-icon
                        width: 7px
                        height: 4px
                        margin-left: 20px

                    &-options
                        z-index: 1
                        width: 100%
                        height: 145px
                        position: absolute
                        background-color: #FFFFFF
                        top: 27px
                        left: 0
                        // overflow: hidden
                        color: #A6A6A6
                        font-size: 14px
                        box-shadow: 0 2px 4px 0 rgba(0,0,0,0.2)
                        border-bottom-left-radius: 2px
                        border-bottom-right-radius: 2px

                        .scroll-view
                            height: 100%
                            overflow-y: auto

                        &-item
                            flex: 1
                            padding: 6px 21px
                            overflow: hidden
                            white-space: nowrap
                            text-overflow: ellipsis
                            font-family: PingFangSC-Regular, PingFang SC
                            font-weight: 400
                            display: flex
                            justify-content: center

                        .active
                            color: #000000
                            font-weight: 700
                            background-color: #f0f0f0

        .service-relation
            justify-content: space-between
            color: #A6A6A6
            margin-top: 15px

    &-tabs
        width: 100%
        font-size: 14px
        color: #A6A6A6
        margin-top: 16px

        .second
            margin-left: 100px

        .active
            color: #000000
            // border-bottom: 2px solid #000000
            font-weight: 700
            position: relative

            &::after
                content: ''
                width: 100%
                height: 2px
                position: absolute
                left: 0
                bottom: -2px
                border-radius: 1px
                background-color: #000000

    &-content
        width: 345px
        padding: 20px
        background-color: #FFFFFF
        margin-top: 15px
        font-size: 14px
        color: #000000
        font-weight: 700

        &-title
            margin-bottom: 16px

        &-write
            font-size: 12px
            font-weight: normal
            text-indent: 28px
            line-height: 22px

        &-sign
            display: flex
            justify-content: flex-end
            margin-top: 41px
            margin-bottom: 25px
            font-size: 14px
            font-family: PingFangSC-Regular, PingFang SC
            font-weight: 400
            color: #000000

            &.m-b-80
                margin-bottom: 80px

        &-backgroundImg

            &-image
                width: 305px
                height: 263px

        .dot-container
            margin-top: 16px

            &-tag
                margin: 0 15px
                font-size: 14px
                font-family: PingFangSC-Regular, PingFang SC
                font-weight: 400
                color: #000000

            .icon
                width: 8px
                height: 13px

            &-icon-left
                transform: rotate(180deg)

            .disabled
                opacity: .3

    .footer
        width: 100%
        margin-bottom: 30px
        display: flex
        align-items: center
        justify-content: center

        &-sendBtn
            width: 260px
            height: 40px
            background-color: #000000
            color: #FFFFFF
            line-height: 20px
            margin-top: 15px
            display: flex
            align-items: center
            justify-content: center
            font-family: PingFangSC-Regular, PingFang SC
            font-weight: 400
            font-size: 14px

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

.edit-box
    margin-top: 15px
    margin-bottom: 25px
    background-color: #F6F6F9
    position: relative
    .van-cell
      background: none
    .edit
        padding: 10px
        font-size: 14px
        font-family: PingFangSC-Regular, PingFang SC
        font-weight: 400
        color: #000000
    .word-total
        font-size: 12px
        font-family: PingFangSC-Regular, PingFang SC
        font-weight: 400
        color: #BDBEC8
        display: flex
        justify-content: flex-end
        position: absolute
        bottom: 10px
        right: 10px
</style>
