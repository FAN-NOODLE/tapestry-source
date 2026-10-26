<template>
  <div class="container display-task">
    <div class="header">
      <!-- 搜索 -->
      <div class="search">
        <div class="search-wrap">
          <van-search
            v-model="searchParams.name"
            clearable
            clear-trigger="always"
            placeholder="请输入任务名称"
            :left-icon="searchIcon"
            @input="handleSearch"
          />
          <div v-if="searchParams.name" class="search-clear" @click="handleClear"></div>
        </div>
        <select-status v-model="searchParams.status" :list="statusOptions" @change="reload"></select-status>
      </div>
      <!-- tab -->
      <van-tabs 
        v-model="searchParams.finish_status"
        color="#000"
        line-width="28px"
        line-height="3px"
        title-active-color="#000"
        title-inactive-color="#777777"
        :ellipsis="false"
        @change="reload">
        <van-tab 
          v-for="(item, index) in auditStatusOptions"
          :key="index"
          :name="item.value">
          <template #title>
            <span class="tab-title" :class="{dot: countInfo[item.key] > 0 && item.value !== -1 && item.value !== 3}">{{ item.label }}</span>
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
        <vm-display-task-item 
          v-for="(item, index) in list" 
          :key="index" 
          :item="item" 
          :status="searchParams.finish_status"
          @view="handleView(item)">
        </vm-display-task-item>
      </van-list>
      <empty-page
        v-if="list.length == 0 && !loading"
        :picture="emptyIcon"
        paddingTop="60"
        paddingBottom="0"
        :text="emptyText"
      ></empty-page>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import EmptyPage from "@/components/base/empty-page.vue";
import VmDisplayTaskItem from "./vm-display-task-item.vue";
import SelectStatus from "./select-status.vue";
import { auditStatusOptions, statusOptions } from "../const";
import { getCheckListApi, getListTipsApi } from '@/api/display-task'

@Component({
  components: {
    VmDisplayTaskItem,
    EmptyPage,
    SelectStatus
  },
})
export default class VmDisplayTaskList extends Vue {
  /** 搜索图标 */
  searchIcon = require('@/assets/display-task/icon_search.png')
  /** 是否正在加载数据  */
  loading = true;
  /** 是否没有更多数据 */
  finished = false;
  /** 搜索字段 */
  searchParams = {
    name: '',
    /** 任务状态 */
    status: -1,
    /** 审核状态 */
    finish_status: -1
  }
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
  }
  /** 列表数据 */
  list: IDisplayTaskItem[] = []
  /** 审核状态选项 */
  get auditStatusOptions() {
    return auditStatusOptions.filter(item => item.value !== 0)
  }
  /** 任务状态选项 */
  get statusOptions() {
    return statusOptions
  }
  /** 用户信息 */
  get userInfo():IUserInfo {
    return this.$store.getters["user/getUserInfo"]
  }
  /** 用户信息 */
  get updateInfo() {
    return this.$store.getters["displayTask/getUpdateList"]
  }
  /** 空页面提示 */
  get emptyText() {
    if(this.searchParams.name) {
      return '未找到相关任务'
    }
    if(this.searchParams.finish_status === -1) {
      return '暂无任务'
    }

    const label = this.auditStatusOptions.find(item => item.value === this.searchParams.finish_status)?.label
    return `暂无${label}任务`
  }
  /** 空页面图标 */
  get emptyIcon() {
    if(this.searchParams.name) {
      return 'noSearchData'
    }
    return 'noDisplayTask'
  }

  mounted() {
    this.getListTips()
    this.fetchData()
  }
  
  async activated() {
    if(this.updateInfo.update) {
      await this.getListTips()
      this.reload()
      this.$store.commit('displayTask/setUpdateList', {
        ...this.updateInfo,
        update: false,
      })
    }
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
   * @description 清空搜索
   */
  handleClear() {
    this.searchParams.name = ''
    this.reload()
  }
  handleSearch = this.$helper.throttle(() => {
    this.reload()
  })

  /**
   * @description 刷新
   */
  reload() {
    this.meta.page = 1;
    this.finished = false
    this.list = [];
    this.fetchData()
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
        ...this.searchParams
      }
      const {code, message, data}:IApiResult = await getCheckListApi(params)
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

  /** 
   * @description 获取角标数据
   */
  async getListTips() {
    try {
      const params = {
        is_vm: 1
      }
      const {code, message, data}:IApiResult = await getListTipsApi(params)
      if (code!== 200) {
        this.$toast.fail(message);
        return;
      }
      this.countInfo = data
    } catch (error) {
      console.log(error);
    }
  }


}
</script>
<style scoped lang="scss">
.tab-title {
  position: relative;
  &.dot::after {
    content: "";
    position: absolute;
    top: -2px;
    right: -4px;
    width: 6px;
    height: 6px;
    border: 1px solid #fff;
    background-color: #E02020;
    border-radius: 50%;
  }
}
.header {
  position: sticky;
  top: 0;
  z-index: 1;
  background-color: #fff;

  .search {
    padding: 10px 15px;
    display: flex;
    align-items: center;
  }
  .search-wrap {
    flex: 1;
    height: 38px;
    background: #F6F6F6;
    position: relative;
    .search-clear {
      position: absolute;
      right: 7px;
      top: 50%;
      transform: translateY(-50%);
      width: 24px;
      height: 24px;
      background: url("~@/assets/display-task/icon_search_clear.png") center/100% no-repeat;
    }
    ::v-deep .van-search {
      padding: 0;
      height: 100%;
      .van-search__content {
        height: 100%;
        padding-left: 10px;
      }
      .van-cell {
        height: 100%;
        padding: 0 15px 0 0 ;
        display: flex;
        align-items: center;
        font-size: 12px;
      }
      .van-cell__value {
        height: 100%;
      }
      .van-field__body {
        height: 100%;
      
      }
      .van-field__left-icon {
        display: flex;
        align-items: center;
        justify-items: center;
        margin-right: 7px;
      }
    }
  }

  
}
.list {
  padding: 15px;
}
</style>