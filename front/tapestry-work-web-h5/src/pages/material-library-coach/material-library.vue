<template>
  <div class="container">
    <div :class="['header', { 'sticky': active===1}]">
      <van-tabs 
        v-if="isCoach"
        class="tabs-wrap"
        v-model="active">
        <van-tab title="KOS素材">
        </van-tab>
        <van-tab title="官网商品素材">
        </van-tab>
      </van-tabs>
    </div>
    <normal-material v-if="active===0"></normal-material>
    <goods-material v-if="active===1"></goods-material>

    <!-- 悬浮球 -->
    <draggable-ball :width="95" :height="34" @click="goCollectPage">
      <template v-slot:icon>
        <van-icon :name="require('@/assets/material-library/v5.1.5/m0529-collect01.png')"></van-icon>
      </template>
      <template v-slot:text>
        <span>我的收藏</span>
      </template>
    </draggable-ball>
  
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import wx from "weixin-js-sdk";
import wxUtils from "@/utils/wx";
import Helper from "@/utils/helper";
import NormalMaterial from "./normal-material.vue";
import GoodsMaterial from "./goods-material.vue";
import DraggableBall from '@/components/material-library/draggable-ball.vue';
@Component({
  components: {
    NormalMaterial,
    GoodsMaterial,
    DraggableBall
  },
})
export default class MaterialLibrary extends Vue {
  active = 0

  /** 判断企业 */
  get isCoach() {
    return Helper.isCoach()
  }
  /**生命周期 */
  activated() {
    // 当前需要直接退出的页面
    window.history.pushState("_android_back", "", location.href);
    window.addEventListener("popstate", this.popstate, false);
  }
  deactivated() {
    // 记得清除!!!!不然单页面，其它页面也可以用
    window.removeEventListener("popstate", this.popstate, false);
  }

  /** 去我的收藏 */
  goCollectPage() {
    // 神策数据上报 -- 点击收藏入口
    this.$track.trackMyFavoriteEntranceClick({
      operation_type: '我的收藏',
      entry_point: '1V1聊天窗口',
      library_type: '素材库'
    })

    this.$router.push('/my-collect')

  }

  async mounted() {
    try {
      await wxUtils.wxAgentConfig();
      this.getWxContent();
    } catch (error) {
      console.log(error);
    }
    // 神策数据上报 -- 进入素材库
    this.$track.trackAssetLibrary({
      entry_point: 2,
      library_type: 2
    })
  }

  popstate() {
    wx.closeWindow(); // 微信网页退出
  }
  getWxContent() {
    const self = this;
    if (!window.wx.invoke) return;
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
        } else if (entry == "group_chat_tools") {
          // 从群聊会话的工具栏进入-打开客户画像-社群页面
          window.wx.invoke("getCurExternalChat", {}, function (res) {
            if (res.err_msg == "getCurExternalChat:ok") {
              console.log("res.chatId", res.chatId);
              self.$store.commit("common/setGroupChatId", res.chatId);
            } else {
              //错误处理
              console.log("getCurExternalChaterr:" + JSON.stringify(res));
            }
          });
        }
      } else {
        //错误处理
        console.log("getContexterr:" + JSON.stringify(res));
      }
    });
  }

}
</script>
<style lang="scss" scoped>
::v-deep .van-tabs__wrap {
  height: 50px;
}

.container {
  min-height: 100vh;
  background: #f4f7fb;
  box-sizing: border-box;
  padding-bottom: calc(50px + constant(safe-area-inset-bottom)); /** iOS 11.0+ */
  padding-bottom: calc(50px + env(safe-area-inset-bottom)); /** iOS 11.0+ */

  .header {
    &.sticky {
      position: sticky;
      top: 0;
      z-index: 10;
    }
  }
}
</style>
