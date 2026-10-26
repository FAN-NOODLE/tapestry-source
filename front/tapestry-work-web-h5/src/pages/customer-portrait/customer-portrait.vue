<template>
  <div>
    <div v-if="!isError">
      <new-portrait-customer v-if="type == 1" />
      <group-portrait v-if="type == 2 && isCoach" :groupChatId="groupChatId" />
    </div>
    <!-- 微信打开无权限 -->
    <empty-page
      v-if="isError"
      :picture="'emptyList'"
      :text="errorTxt"
    ></empty-page>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import groupPortrait from "./group-portrait/group-portrait.vue";
import newPortraitCustomer from "./new-portrait-customer/index.vue";
import emptyPage from "@/components/base/empty-page.vue";
import wxUtils from "@/utils/wx";
import wx from "weixin-js-sdk";
import Helper from "@/utils/helper";
@Component({
  components: {
    groupPortrait,
    newPortraitCustomer,
    emptyPage,
  },
})
export default class extends Vue {
  // 是否是sw
  isSW = Helper.isSW();
  // 是否是ks
  isKS = Helper.isKS();
  //  是否是coach
  isCoach = Helper.isCoach();
  isError = false;
  errorTxt = "";
  public get type(): number {
    return this.$store.getters["common/getPageType"];
  }
  public get groupChatId(): string {
    return this.$store.getters["common/getGroupChatId"];
  }
  private getWxContent() {
    const self = this;
    console.log("getContext=============");
    window.wx.invoke("getContext", {}, function (res) {
      if (res.err_msg == "getContext:ok") {
        console.log("getContext", res.entry);
        const entry = res.entry;
        // 从单聊会话的工具栏进入-打开客户画像-客户页面
        if (entry == "single_chat_tools") {
          window.wx.invoke("getCurExternalContact", {}, (res) => {
            if (res.err_msg == "getCurExternalContact:ok") {
              console.log("customer:getCurExternalContact", res.userId);
              self.$store.commit("user/setExternalId", res.userId);
              self.$store.commit("common/setPageType", 1);
              // 神策数据上报 -- 点击侧边栏
              self.$track.trackClickItem({
                feature: self.$route.path,
                customer_id: res.userId
              })
            } else {
              //错误处理
              console.log("getCurExternalContacterr:" + JSON.stringify(res));
            }
          });
        } else if (entry == "group_chat_tools" && self.isCoach) {
          // 从群聊会话的工具栏进入-打开客户画像-社群页面
          window.wx.invoke("getCurExternalChat", {}, function (res) {
            if (res.err_msg == "getCurExternalChat:ok") {
              console.log("res.chatId", res.chatId);
              self.$store.commit("common/setPageType", 2);
              self.$store.commit("common/setGroupChatId", res.chatId);
            } else {
              //错误处理
              console.log("getCurExternalChaterr:" + JSON.stringify(res));
            }
          });
        } else {
          self.isError = true;
          self.errorTxt = "通过私聊窗口发送";
        }
      } else {
        //错误处理
        console.log("getContexterr:" + JSON.stringify(res));
      }
    });
  }
  private popstate() {
    wx.closeWindow(); // 微信网页退出
  }
  private async created() {
    // 当前需要直接退出的页面
    window.history.pushState("_android_back", "", location.href);
    window.addEventListener("popstate", this.popstate, false);
  }
  private async destroyed() {
    // 记得清除，不然单页面，其它页面也可以用
    window.removeEventListener("popstate", this.popstate, false);
  }
  private async mounted() {
    console.log("render this");
    if (this.type) return;
    try {
      this.getWxContent();
    } catch (e) {
      this.isError = true;
      this.errorTxt = "获取用户ID失败";
    }
  }
}
</script>
