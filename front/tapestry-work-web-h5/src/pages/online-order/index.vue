<template>
  <div>
    <!-- 订单列表 -->
    <OnlineOrderList v-if="role === 3"></OnlineOrderList>
    <!-- 数据统计 -->
    <OnlineOrderOffice v-if="role === 5"></OnlineOrderOffice>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import OnlineOrderList from './list.vue';
import OnlineOrderOffice from './office.vue';
import { getUserAuthApi } from '@/api/customer-portrait';
import wx from 'weixin-js-sdk'
@Component({
  components: {
    OnlineOrderList,
    OnlineOrderOffice
  }
})
export default class OnlineOrderIndex extends Vue {
  role = 0
  
  private async mounted() {
    await this.getUserAuth()
    // 神策数据上报 -- 进入网仓
    this.$track.trackEnterOnlineWarehouse({
      OnlineWarehouse_entrance: '企微工作台',
      message_type: '企微工作台',
    })
  }
  /**生命周期 */
  private async created() {
    // 当前需要直接退出的页面
    window.history.pushState("_android_back", "", location.href);
    window.addEventListener("popstate", this.popstate, false);
    
  }
  private destroyed() {
    // 记得清除!!!!不然单页面，其它页面也可以用
    window.removeEventListener("popstate", this.popstate, false);
  }

  private activated() {
    // 当前需要直接退出的页面
    window.addEventListener("popstate", this.popstate, false);

  }
  private deactivated() {
    // 记得清除!!!!不然单页面，其它页面也可以用
    window.removeEventListener("popstate", this.popstate, false);
  }


  popstate() {
    wx.closeWindow(); // 微信网页退出
  }
  async getUserAuth() {
    const res = await getUserAuthApi({})
    this.role = res.data.role
    this.$store.commit('user/setRole', res.data.role)
  }

}
</script>
<style lang="scss" scoped>
</style>