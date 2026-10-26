<template>
  <div v-if="current !== -1" class="container">
    <div v-if="tabs.length>1" class="tab" :data-length="tabs.length">
      <div 
        class="tab-item" 
        :class="{
          active: current === index
        }"
        v-for="(item, index) in tabs" 
        :key="index"
        @click="handleTabChange(index)"
      >
        {{ item.label }}
      </div>
    </div>
    <display-task-data v-if="tabs[current].value === 0"></display-task-data>
    <vm-display-task-list v-if="tabs[current].value === 1"></vm-display-task-list>
    <display-task-list v-if="tabs[current].value === 2"></display-task-list>
    
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import DisplayTaskList from "./components/display-task-list.vue";
import VmDisplayTaskList from "./components/vm-display-task-list.vue";
import DisplayTaskData from "./components/display-task-data.vue";
import { getAuthApi } from '@/api/display-task'
import wx from 'weixin-js-sdk'
@Component({
  components: {
    DisplayTaskList,
    VmDisplayTaskList,
    DisplayTaskData
  },
  beforeRouteLeave(to, from, next) {
    this.scrollTop = document.documentElement.scrollTop;
    console.log(this.scrollTop);
    
    next()
  }
})
export default class DisplayTask extends Vue {
  scrollTop = 0;
  current:number = -1;
  /** 用户信息 */
  get userInfo():IUserInfo {
    return this.$store.getters["user/getUserInfo"]
  }

  get tabs() {
    const list = [
      {
        label: '陈列任务',
        value: 2
      }
    ]

    if(this.userInfo.is_vm === 1) {
      list.unshift({
        label: 'VM审核',
        value: 1
      })
    }

    if(this.userInfo.is_manager === 1) {
      list.unshift({
        label: '陈列数据',
        value: 0
      })
    }

    return list
  }

  async mounted() {
    await this.getAuth()
    // 神策数据上报 -- 进入VM
    this.$track.trackEnterDisplayTask({
      vm_entrance: 1
    })
  }
  
  /**生命周期 */
  private async created() {
    // 当前需要直接退出的页面
    window.history.pushState("_android_back", "", location.href);
    window.addEventListener("popstate", this.popstate, false);
    
  }
  private async destroyed() {
    // 记得清除!!!!不然单页面，其它页面也可以用
    window.removeEventListener("popstate", this.popstate, false);
  }

  activated() {
    // 当前需要直接退出的页面
    window.addEventListener("popstate", this.popstate, false);

    this.$nextTick(() => {
      if (this.scrollTop > 0) {
        console.log('activated', this.scrollTop);
  
        document.documentElement.scrollTo({
          top: this.scrollTop,
          behavior: "auto"
        });
      }
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
      const {code, message, data}:IApiResult = await getAuthApi()
      if (code!== 200) {
        this.$toast.fail(message);
        return;
      }
      this.$store.commit("user/setUserInfo", {
        app_source: 1,
        ...this.userInfo, 
        ...data,
      })
      this.current = 0
    } catch (error) {
      console.log(error);
    }
  }

  /**
   * 切换
   * @param index 
   */
  handleTabChange(index:number) {
    this.current = index
  }

}
</script>
<style scoped lang="scss">
.container {
  min-height: 100vh;
  background: #F7F7F7;
}
.tab {
  padding: 12px 36px;
  background-color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  &[data-length="3"] {
    padding: 12px 15px;
  }
  &-item {
    height: 36px;
    font-size: 14px;
    font-weight: 400;
    color: #777777;
    line-height: 20px;
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    &:nth-child(2) {
      margin-left: 15px;
    }

    &.active {
      color: #fff;
      background-color: #000;
    }
  }
}
</style>