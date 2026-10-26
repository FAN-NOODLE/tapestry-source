<template>
  <div class="communication">
    <!-- 顶部tab -->
    <van-tabs v-if="tabs.length>1" v-model="active" :line-width="100/tabs.length + 'vw'" line-height="0.05rem" color="#000" >
      <van-tab 
        v-for="(item, index) in tabs" 
        :key="index" 
        :title="item.title" 
        :badge="detail[item.prop]||''"
      >
        <component :is="item.component"></component>
      </van-tab>
    </van-tabs>
    <component v-else :is="tabs[active].component"></component>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import CommunicationCorporateList from './corporate-list.vue'
import CommunicationMyList from './my-list.vue';
import CommunicationAuditList from './audit-list.vue';
import { getManagerInfoApi, getReadNumApi } from '@/api/communication';
import wx from 'weixin-js-sdk'
@Component({
  components: {
    CommunicationCorporateList,
    CommunicationMyList,
    CommunicationAuditList
  },
  beforeRouteLeave(to, from, next) {
    this.scrollTop = document.documentElement.scrollTop;
    next()
  }
})
export default class Communication extends Vue {
  detail = {}
  scrollTop = 0
  /** current */
  active = 0
  /** tab */
  get tabs() {
    const list: any[] = [
      {
        title: '企业传讯',
        prop: 'merge_info_num',
        component: CommunicationCorporateList
      }
    ]
    // ro
    if (this.userInfo.is_ro) {
      list.unshift({
        title: '传讯审核',
        prop: 'info_checking_num',
        component: CommunicationAuditList
      })
    }
    // ft&ro
    if(this.userInfo.is_ro || this.userInfo.is_ft) {
      list.unshift({
        title: '我的传讯',
        prop: 'info_fail_num',
        component: CommunicationMyList
      })
    }
    
    return list
  }
  /** 用户信息 */
  get userInfo(): IUserInfo {
    return this.$store.getters["user/getUserInfo"]
  }

  /**生命周期 */
  private async mounted() {
    await this.getAuth()
    this.getReadNum()
    // 当前需要直接退出的页面
    window.history.pushState("_android_back", "", location.href);
    window.addEventListener("popstate", this.popstate, false);
    // 神策数据上报 -- 进入传讯
    this.$track.trackEnterTransmittal({
      transmittal_entrance: 1
    })
  }
  private destroyed() {
    // 记得清除!!!!不然单页面，其它页面也可以用
    window.removeEventListener("popstate", this.popstate, false);
  }

  activated() {
    this.getReadNum()
    // 当前需要直接退出的页面
    window.addEventListener("popstate", this.popstate, false);

    this.$nextTick(() => {
      document.documentElement.scrollTo({
        top: this.scrollTop,
        behavior: "auto"
      });
    })
  }
  deactivated() {
    // 记得清除!!!!不然单页面，其它页面也可以用
    window.removeEventListener("popstate", this.popstate, false);
  }

  popstate() {
    wx.closeWindow(); // 微信网页退出
  }

  /**
   * @description 获取用户权限
   */
  async getAuth() {
    try {
      const {code, message, data}:IApiResult = await getManagerInfoApi({})
      if (code!== 200) {
        this.$toast.fail(message);
        return;
      }
      this.$store.commit("user/setUserInfo", {
        app_source: 2,
        ...this.userInfo, 
        ...data,
      })
    } catch (error) {
      console.log(error);
    }
  }

  async getReadNum() {
    try {
      const {code, message, data}:IApiResult = await getReadNumApi({})
      if (code!== 200) {
        this.$toast.fail(message);
        return;
      }
      this.detail = data
    } catch (error) {
      console.log(error);
    }
  }
}
</script>
<style lang="scss" scoped>
::v-deep {
  .van-info {
    background-color: #000;
    font-size: 7px;
    width: 13px;
    height: 13px;
    display: flex;
    align-items: center;
    justify-content: center;
    min-width: 13px;
    right: -3px
  }

  // .van-tabs {
  //   height: 100%;
  //   display: flex;
  //   flex-direction: column;
  // }

  // .van-tabs__content {
  //   flex: 1;
  //   overflow: hidden;
  // }

  // .van-tab__pane {
  //   height: 100%;
  // }
}

.communication {
  width: 100vw;
  min-height: 100vh;
  background: #F7F7F7;
}
</style>