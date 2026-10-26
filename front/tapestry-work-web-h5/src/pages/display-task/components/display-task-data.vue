<template>
  <div class="container display-task-data">
    <div class="header">
      <select-task-popup 
        ref="taskRef"
        :list="task.list"
        :selected="task.selected"
        :defaultSelected="task.default"
        @confirm="handleTaskPopupConfirm"
      ></select-task-popup>
      <select-store-popup 
        ref="storeRef"
        :list="store.list"
        :selected="store.selected"
        :defaultSelected="store.default"
        @confirm="handleStorePopupConfirm"
      ></select-store-popup>
    </div>
    <div class="panel">
      <div class="panel-item" v-for="(item, index) in auditStatusOptions" :key="index" @click="handleViewDetail(item.value, info[item.key]||0)">
        <div class="panel-item__value">{{info[item.key]||0}}</div>
        <div class="panel-item__name">
          <div class="ellipsis">{{item.label}}门店</div>
          <div class="tips" @click.stop="handleTips(item.value)"></div>
        </div>
        <div class="panel-item__rate">占比：{{info[item.rate_key]||0}}</div>
        <div v-if="info[item.key]" class="next"></div>
      </div>
    </div>
    <div>
      <div class="title">管理范围下【未通过】门店</div>
      <div class="search">
        <van-search v-model="keyword" placeholder="请输入门店名称搜索" :left-icon="searchIcon" :clearable="false" />
        <div v-if="keyword" class="search-clear" @click="handleClear"></div>
      </div>
      <div class="task-title">任务名称：{{taskName || '--'}}</div>
      <div class="table-wrap">
        <com-table 
          :columns="columns" 
          :data="renderStoreList"
          @row-click="handleRowClick"
          >
          <template v-slot:review_time="{row}">
            <span>{{ $helper.dateFormat(row.review_time, "yyyy-MM-dd") }}</span>
          </template>
          <template v-slot:operation>
            <!-- <div class="icon-right"></div> -->
            <span>查看</span>
          </template>
          <template v-slot:empty>
            <empty-page
              picture="noSearchData"
              paddingTop="60"
              paddingBottom="100"
              text="暂无未通过门店哦"
            ></empty-page>
          </template>
        </com-table>
      </div>
      <div v-if="renderStoreList.length>0" class="load-more">～已经到底了～</div>
    </div>
    
    <!-- 提示信息 -->
    <van-popup v-model="tips.show">
      <div class="popup-tips">
        <div class="popup-tips__close" @click="tips.show = false"></div>
        <div class="popup-tips__title">{{tips.title}}</div>
        <div class="popup-tips__content">
          <div>{{tips.message}}</div>
          <div>（不同任务同一门店合并去重展示）</div>
        </div>
      </div>
    </van-popup>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop } from 'vue-property-decorator'
import SelectStorePopup, { IStoreItem } from './select-store-popup.vue';
import SelectTaskPopup, { ITaskItem } from './select-task-popup.vue';
import EmptyPage from "@/components/base/empty-page.vue";
import ComTable from './com-table.vue'
import { 
  getAllTaskListApi, 
  getDashboardInfoApi, 
  getFailListApi, 
  getMemberAuthApi, 
  getStoreTreeApi 
} from '@/api/display-task';
import { auditStatusOptions } from '../const';
import { Popup } from 'vant';
interface ITask {
  list: ITaskItem[];
  selected: number[];
  default: number[];
}
interface IStore {
  list: IStoreItem[];
  selected: number[];
  default: number[];
}

interface IFailStoreItem {
  allot_id: number;
  /** 门店名称 */
  store_name: string;
  /** 门店ID */
  store_id: number;
  /** 审核时间 */
  review_time: string;
}

@Component({
  components: {
    [Popup.name]: Popup,
    SelectStorePopup,
    SelectTaskPopup,
    ComTable,
    EmptyPage
  }
})
export default class DisplayTaskData extends Vue {
  info: Record<string, number> = {};

  task: ITask = {
    /** 任务列表 */
    list: [],
    /** 已选择的任务 */
    selected: [],
    default: [],
  }

  store: IStore = {
    /** 任务列表 */
    list: [],
    /** 已选择的任务 */
    selected: [],
    default: [],
  }

  /** 搜索关键字 */
  keyword = '';
  /** 搜索图标 */
  searchIcon = require('@/assets/display-task/icon_search.png')

  columns = [
    {
      label: '门店名称',
      prop: 'store_name',
    },
    {
      label: '审批时间',
      prop: 'review_time',
      width: 140
    },
    {
      label: '操作',
      prop: 'operation',
      width: 80
    },
  ]

  failInfo: {
    /** 任务ID */
    id: number;
    /** 任务名称 */
    display_name: string;
    /** 门店列表 */
    list: IFailStoreItem[];
  } = {
    id: 0,
    display_name: '',
    list: []
  }

  tips = {
    show: false,
    title: '',
    message: ''
  }
  
  get renderTaskValue() {
    const { list, selected } = this.task
    const value = list.filter(item => selected.includes(item.id))
    return {
      name: value.map(item => item.display_name).join('、'),
      count: value.length
    }
  }

  get auditStatusOptions() {
    return auditStatusOptions.filter(item => item.value >= 0)
  }

  get renderStoreList() {
    const { list } = this.failInfo
    return list.filter(item => item.store_name.indexOf(this.keyword) >= 0)
  }

  get taskName () {
    const { list } = this.task
    return this.failInfo.display_name || list[0]?.display_name || '--'
  }


  async mounted() {
    await this.getAllTaskList()
    await this.getMemberAuth()
    this.getStoreTree()
    this.getDashboardInfo()
    this.getFailList()
  }

  async getAllTaskList() {
    const res:any = await getAllTaskListApi()
    if (res.code !== 200) {
      this.$toast.fail(res.message);
      return
    }
    this.task.list = res.data
    // 如果有任务做处理
    if(res.data.length) {
      this.task.selected = [res.data[0].id]
      this.task.default = [res.data[0].id]
    }
  }

  async getDashboardInfo() {
    const params = {
      ids: this.task.selected,
      area_ids: this.store.selected
    }
    const res:any = await getDashboardInfoApi(params)
    if (res.code !== 200) {
      this.$toast.fail(res.message);
      return
    }
    this.info = res.data
  }

  async getFailList() {
    const res: any = await getFailListApi({
      id: this.task.list[0].id
    })
    if (res.code !== 200) {
      this.$toast.fail(res.message);
      return
    }
    
    this.failInfo = res.data
  }

  async getMemberAuth() {
    const res: any = await getMemberAuthApi()
    if (res.code !== 200) {
      this.$toast.fail(res.message);
      return
    }
    
    this.store.selected = [...res.data]
    this.store.default = [...res.data]
  }

  async getStoreTree() {
    const res: any = await getStoreTreeApi({
      has_children: true
    })
    if (res.code !== 200) {
      this.$toast.fail(res.message);
      return
    }
    
    this.store.list = res.data
  }

  handleTaskPopupConfirm(value: ITask['selected']) {
    this.task.selected = value
    this.getDashboardInfo()
  }


  handleStorePopupConfirm(value: IStore['selected']) {
    this.store.selected = value
    this.getDashboardInfo()
  }

  /**
   * 清空搜索
   */
   handleClear() {
    this.keyword = ''
  }

  /**
   * handle row click
   * @param row 
   */
  handleRowClick(item: any) {
    this.$store.commit('displayTask/setDetail', {
      task_info: null,
      store_info: null
    });
    this.$store.commit('displayTask/setUpdate', true);
    this.$router.push({
      path: 'display-task-detail',
      query: {
        id: this.failInfo.id + '',
        store_id: item.store_id + '',
        from: '1',
        manager: '1' // 从数据管理进入
      }
    });
  }
  /**
   * 下钻
   */
  handleViewDetail(value: number, num: number) {
    if(num === 0) return
    const taskInfo = (this.$refs.taskRef as any).getInfo()
    const storeInfo = (this.$refs.storeRef as any).getInfo()
    
    this.$store.commit('displayTask/setSearchParams', {
      ...taskInfo,
      ...storeInfo
    })
    this.$router.push({
      name: 'display-store-list',
      query: {
        finish_status: value+''
      }
    })
  }

  /**
   * 提示
   */
  handleTips(value: number) {
    const message = [
      '所选范围内有未提交任务的门店数',
      '所选范围内有待审核任务的门店数',
      '所选范围内有未通过任务的门店数',
      '所选范围内有已通过任务的门店数',
      '所选范围内有待复审任务的门店数',
    ][value]
    const title = [
      '未提交门店',
      '待审核门店',
      '未通过门店',
      '已通过门店',
      '待复审门店',
    ][value]
    this.tips = {
      show: true,
      title: title,
      message: message
    }
  }
}
</script>

<style lang="scss" scoped>
.container.display-task-data {
  min-height: 100vh;
  background-color: #F7F7F7;
  .header {
    padding: 10px 15px;
    display: flex;
    align-items: center;
    background-color: #fff;
    font-size: 12px;
    box-shadow: 0px 6px 12px 0px rgba(222,224,227,0.5);
    
  }

  .panel {
    padding: 15px;
    display: flex;
    flex-wrap: wrap;
    
    .panel-item {
      width: 165px;
      // height: 80px;
      background: #FFFFFF;
      border-radius: 5px;
      box-sizing: border-box;
      position: relative;
      border: 1px solid #000000;

      &:nth-child(1),
      &:nth-child(2),
      &:nth-child(3) {
        margin-bottom: 15px;
        width: 105px;
      } 

      &:nth-child(1),
      &:nth-child(2),
      &:nth-child(4) {
        margin-right: 15px;
      } 

      .next {
        position: absolute;
        right: -1px;
        bottom: -1px;
        width: 24px;
        height: 18px;
        background: url('~@/assets/display-task/icon_11.png') center/contain no-repeat;
      }

      &__value {
        font-size: 16px;
        font-weight: 500;
        line-height: 22px;
        padding-top: 10px;
        text-align: center;
      }

      &__name {
        margin-top: 6px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 13px;
        // line-height: 18px;
        color: rgba(0,0,0,0.85);

        .tips {
          margin-left: 4px;
          width: 14px;
          height: 14px;
          background: url('~@/assets/display-task/icon_tips_black.png') center/20px no-repeat;
        }
      }

      &__rate {
        font-size: 10px;
        line-height: 14px;
        padding: 5px 0;
        text-align: center;
        color: #777
      }
    }
  }

  .title {
    font-size: 14px;
    font-weight: 600;
    color: #000000;
    line-height: 20px;
    text-align: center;
  }

  .search {
    padding: 15px;
    margin-top: 3px;
    position: relative;
    .van-search {
      height: 38px;
      padding: 0;
      border: 1px solid #A6A6A6;
      background-color: transparent;
    }

    .van-search__content {
      height: 100%;
      background-color: transparent;
    }
    .van-cell {
      height: 100%;
      padding: 0 15px 0 0 ;
      display: flex;
      align-items: center;
      font-size: 12px;
    }
    ::v-deep {
      .van-cell__value {
        height: 100%;
      }
      .van-field__body {
        height: 100%;
      }
      input.van-field__control {
        height: 100%;
      }
    }
    .van-field__left-icon {
      display: flex;
      align-items: center;
      justify-items: center;
      margin-right: 7px;
    }
    .search-clear {
      position: absolute;
      right: 20px;
      top: 50%;
      transform: translateY(-50%);
      width: 20px;
      height: 20px;
      background: url("~@/assets/display-task/icon_search_clear_2.png") center/100% no-repeat;
    }
  }

  .task-title {
    padding: 0 20px 15px;
    font-size: 12px;
    font-weight: 500;
    color: #000000;
    line-height: 17px;
  }

  .table-wrap {
    padding: 0 15px;
    ::v-deep {
      .table {
        background-color: #fff;
      }
      .table__header {
        background-color: rgba($color: #000000, $alpha: 0.1);

        .table__col {
          padding: 15px 0;
          color: #777;
          font-weight: normal;
        }
      }

      .empty {
        background-color: #f7f7f7;
      }
      .default {
        height: auto;
      }
    } 
  }
  .icon-right {
    width: 11px;
    height: 17px;
    background: url('~@/assets/display-task/icon_right_black.png') center/100% no-repeat;
  }
  .load-more {
    padding: 10px;
    text-align: center;
    font-size: 12px;
    color: #A6A6A6;
  }

  .popup-tips {
    width: 249px;
    height: 124px;
    background: #FFFFFF;
    padding: 20px 0;
    position: relative;
    &__close {
      position: absolute;
      right: 10px;
      top: 10px;
      width: 18px;
      height: 18px;
      background: url('~@/assets/display-task/icon_popup_close.png') center/100% no-repeat;
    }
    &__title {
      font-size: 14px;
      font-weight: 500;
      color: #000000;
      line-height: 20px;
      text-align: center;
    }
    &__content {
      font-size: 12px;
      font-weight: 400;
      color: #676767;
      line-height: 17px;
      padding: 15px;
      text-align: center;
    }
  }
}
</style>