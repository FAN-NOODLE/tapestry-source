<template>
  <div class="display-task">
    <div class="header">
      <van-tabs 
        v-model="current"
        color="#000"
        line-width="28px"
        line-height="3px"
        title-active-color="#000"
        title-inactive-color="#777777"
        @change="reload">
        <van-tab 
          v-for="(item, index) in statusOptions"
          :key="index"
          :name="item.value">
          <template #title>
            <span>{{ item.label + (countInfo[item.key]||'') }}</span>
          </template>
        </van-tab>
      </van-tabs>
    </div>
    <div class="content">
      <van-list
        class="list"
        v-model="loading"
        :finished="finished"
        :finished-text="list.length > 0 ? '没有更多了~':''"
        @load="fetchData"
      >
        <display-task-item 
          v-for="(item, index) in list" 
          :key="index" 
          :item="item" 
          @view="handleView(item)"></display-task-item>
      </van-list>
      <empty-page
        v-if="list.length == 0 && !loading"
        picture="noDisplayTask"
        paddingTop="60"
        paddingBottom="0"
        text="暂无任务"
      ></empty-page>
    </div>
    <select-store ref="selectStore" @confirm="handleViewConfirm"></select-store>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import EmptyPage from "@/components/base/empty-page.vue";
import DisplayTaskItem from "./display-task-item.vue";
import SelectStore from "./select-store.vue";
import { auditStatusOptions } from "../const";
import { getTaskListApi, getListTipsApi } from '@/api/display-task'

@Component({
  components: {
    DisplayTaskItem,
    EmptyPage,
    SelectStore
  },
})
export default class DisplayTaskList extends Vue {
  /** 当前选中的tab */
  current = -1;
  /** 是否正在加载数据  */
  loading = true;
  /** 是否没有更多数据 */
  finished = false;
  /** 分页 */
  meta: IMeta = {
    page: 1,
    size: 10,
    total: 0,
  }
  /** 统计信息 */
  countInfo = {
    /** 待审核 */
    check_num: 0,
    /** 不通过 */
    fail_num: 0,
    /** 已通过 */
    success_num: 0,
    /** 待提交 */
    wait_num: 0,
    /** 全部 */
    all_num: 0,
    /** 待复审 */
    double_check_num: 0,
  }
  /** 列表数据 */
  list: IDisplayTaskItem[] = []
  /** 审核状态选项 */
  get statusOptions() {
    return auditStatusOptions.filter(item => item.value !== 4)
  }
  /** 用户信息 */
  get userInfo():IUserInfo {
    return this.$store.getters["user/getUserInfo"]
  }
  /** 用户信息 */
  get updateInfo() {
    return this.$store.getters["displayTask/getUpdateList"]
  }

  mounted() {
    this.getListTips()
    this.fetchData()
  }

  async activated() {
    // 更新列表
    if(this.updateInfo.update) {
      this.$toast.success('已提交，等待VM审核')
      this.current = this.updateInfo.current
      await this.getListTips()
      this.reload()
      this.$store.commit('displayTask/setUpdateList', {
        update: false,
        current: -1
      })
    }
  }

  /**
   * @description tab切换
   */
  reload() {
    this.meta.page = 1;
    this.finished = false
    this.list = [];
    this.fetchData()
  }

  /**
   * @description 查看店铺
   * @param item 
   */
  handleView(item: IDisplayTaskItem) {
    (this.$refs.selectStore as any).showPopup({
      id: item.id,
      store_id: item.store_id
    })
  }

  /**
   * @description 他店陈列
   */
  handleViewConfirm(value:any) {
    this.$store.commit('displayTask/setDetail', {
      task_info: null,
      store_info: null
    });
    this.$router.push({
      name: 'display-task-preview',
      query: {
        ...value,
        other: '1'
      }
    })
  }

  /** 
   * @description 获取列表数据
   */
  async fetchData() {
    this.loading = true;
    try {
      const params = {
        page: this.meta.page,
        size: this.meta.size,
        finish_status: this.current,
      }
      const {code, message, data}:IApiResult = await getTaskListApi(params)
      if (code !== 200) {
        this.finished = true;
        this.loading = false;
        this.$toast.fail(message);
        return;
      }
      this.loading = false;
      this.finished = data.list.length === 0;
      this.list = this.list.concat(data.list)
      this.meta.page = this.meta.page + 1

    } catch (error) {
      console.log(error);
      this.finished = true;
      this.loading = false;
    }
  }

  async getListTips() {
    try {
      const params = {
        is_vm: 0
      }
      const {code, message, data}:IApiResult = await getListTipsApi(params)
      if (code!== 200) {
        this.$toast.fail(message);
        return;
      }
      this.countInfo = {
        check_num: data.check_num,
        fail_num: data.fail_num,
        success_num: 0,
        wait_num: data.wait_num,
        all_num: 0,
        double_check_num: data.double_check_num,
      }
    } catch (error) {
      console.log(error);
    }
  }

}
</script>
<style scoped lang="scss">
  .header {
    position: sticky;
    top: 0;
    z-index: 1;
    background-color: #fff;
    height: 44px;
  }
  .list {
    padding: 15px;
  }

</style>