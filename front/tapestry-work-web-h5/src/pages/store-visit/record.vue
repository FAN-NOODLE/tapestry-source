<template>
  <div :key="updateKey+''" class="container">
    <div class="header">
      <div class="flex">
        <div class="tab">
          <div class="tab-item" :class="{active: searchForm.status === item.value}" v-for="(item, index) in tabs" @click="handleTabChange(item.value)">{{ $t(item.label) }}</div>
        </div>
        <div class="flex">
          <div class="switch">
            <div class="label">{{ $t('storeVisit.btn.14') }}</div>
            <van-switch v-model="searchForm.data_type" :disabled="!userInfo.is_store_visit_manager" :active-value="1" :inactive-value="0" @change="handleDataTypeChange(1)"></van-switch>
          </div>
          <div v-if="userInfo.is_store_visit_manager" class="switch">
            <div class="label">{{ $t('storeVisit.btn.15') }}</div>
            <van-switch v-model="searchForm.data_type_2" :active-value="1" :inactive-value="0" @change="handleDataTypeChange(2)"></van-switch>
          </div>
        </div>
      </div>
      <div class="flex">
        <div class="select-item" @click="handleSelectDate">
          <div v-if="dateTypeIndex > -1" class="value">{{$t(dateName)}}</div>
          <div v-else class="placeholder">{{ $t('storeVisit.placeholder.3') }}</div>
          <div class="icon-down"></div>
          <div v-if="showSelectDatePopup" class="popup">
            <div class="item" v-for="(item, index) in dateType" :key="index" @click.stop="handleSelectDateConfirm(index)">{{ $t(item.label) }}</div>
          </div>
          <div v-if="showSelectDatePopup" class="popup-mask" @click.stop="handleSelectDateCancel"></div>
        </div>
        <div class="select-item" @click="handleSelectStore">
          <div v-if="storeName" class="value">{{storeName}}</div>
          <div v-else class="placeholder">{{ $t('storeVisit.placeholder.4') }}</div>
          <div class="icon-down" :class="{active: showSelectStorePopup}"></div>
        </div>
        <div class="reset" :class="{disabled: disabled}" @click="reset">{{ $t('storeVisit.btn.13') }}</div>
      </div>
    </div>
    <div class="content">
      <van-list class="list" 
        v-model="loading"
        :finished="finished"
        :finished-text="finishedText"
        :immediate-check="false"
        @load="getData">
        <visit-record-item v-for="(item, index) in list" :key="index" :item="item" @delete="handleDelete(index)"></visit-record-item>
        <!-- has -->
      </van-list>
      <van-empty
        v-if="list.length == 0 && !loading"
        image-size="66"
        :image="emptyIcon"
        :description="emptyText"
      >
      </van-empty>
    </div>
    <select-store-popup v-model="showSelectStorePopup" :storeTree="storeTree" :default-value="searchForm.third_id ? [searchForm.third_id]:[]" :multiple="false" immediate @select="handleSelectChange" @confirm="handleSelectStoreConfirm"></select-store-popup>
    <calendar-popup v-model="showDateRangePopup" :default-value="defaultDate" @cancel="handleSelectDateRangeCancel" @confirm="handleSelectDateRangeConfirm"></calendar-popup>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Watch } from "vue-property-decorator";
import { Button, List } from "vant";
import { getStoreVisitListApi } from "@/api/store-visit";
import SelectStorePopup from "./components/SelectStorePopup/index.vue";
import CalendarPopup from "./components/calendar-popup.vue";
import VisitRecordItem from "./components/visit-record-item.vue";
import dayjs from "dayjs";
import { getMonthListApi } from "@/api/online-order";
import { getStoreTreeApi } from "@/api/display-task";
@Component({
  components: {
    [Button.name]: Button,
    SelectStorePopup,
    CalendarPopup,
    VisitRecordItem
  },
})
export default class StoreVisitRecord extends Vue {
  tabs = [
    {
      label: 'storeVisit.item.status.3',
      value: 2
    },
    {
      label: 'storeVisit.item.status.2',
      value: 1
    },
  ]
  /** 搜索表单 */
  searchForm = {
    status: 1,
    third_id: 0,
    start_time: 0,
    end_time: 0,
    data_type: 1,
    data_type_2: 1
  }
  /** 分页 */
  meta: IMeta = {
    page: 1,
    size: 10,
  }
  /** 是否正在加载数据  */
  loading = false;
  /** 是否没有更多数据 */
  finished = false;
  /** 列表数据 */
  list:any[] = []
  /** 选择门店 */
  showSelectStorePopup = false;
  showSelectDatePopup = false;
  showDateRangePopup = false;
  defaultDate = []

  dateType = [
    {
      label: 'storeVisit.popup.tips.25',
      value: 1
    },
    {
      label: 'storeVisit.popup.tips.26',
      value: 2
    },
    {
      label: 'storeVisit.popup.tips.27',
      value: 3
    },
    {
      label: 'storeVisit.popup.tips.28',
      value: 4
    },
  ]
  dateTypeIndex = -1

  monthList = []

  storeName = ''

  updateKey = 0

  storeTree = []

  get disabled() {
    return !this.searchForm.third_id && !this.searchForm.start_time
  }

  get dateName() {
    if(this.dateTypeIndex === 3) {
      return dayjs(this.searchForm.start_time*1000).format('YYYY.MM.DD') + ' - ' + dayjs(this.searchForm.end_time*1000).format('YYYY.MM.DD')
    }
    return this.dateType[this.dateTypeIndex].label
  }

  get finishedText() {
    return this.list.length > 0 ? this.$t('storeVisit.empty.5'):''
  }
  /** 空页面提示 */
  get emptyText() {
    if(this.searchForm.third_id || this.searchForm.start_time) {
      return this.$t('storeVisit.empty.4')
    }
    return this.$t('storeVisit.empty.3')
  }
  /** 空页面图标 */
  get emptyIcon() {
    return require('@/assets/store-visit/empty.png')
  }

  /** 用户信息 */
  get userInfo():IUserInfo {
    return this.$store.getters["user/getUserInfo"]
  }

  /** 是否更新列表 */
  get isUpdate() {
    return this.$store.getters["storeVisit/getUpdate"];
  }

  @Watch('dateTypeIndex')
  async handleDateTypeChange(index) {
    console.log(index);
    if(index === -1) return
    if(index === 0) {
      // 最近一周
      this.searchForm.start_time = dayjs().set('h',0).set('m',0).set('s',0).subtract(7, 'day').unix()
      this.searchForm.end_time = dayjs().set('h',23).set('m',59).set('s',59).subtract(1, 'day').unix()
      this.reload()
    } else if(index === 1) {
      // 最近一个财月
      const index = this.monthList.findIndex(item => item.now == 1)
      this.searchForm.start_time = this.monthList[index].start_time
      this.searchForm.end_time = this.monthList[index].end_time
      this.reload()
    } else if(index === 2) {
      // 最近三个财月
      const index = this.monthList.findIndex(item => item.now == 1)
      this.searchForm.start_time = this.monthList[index - 2].start_time
      this.searchForm.end_time = this.monthList[index].end_time
      this.reload()
    } else {
      // 自定义时间
    }

    // 神策数据上报 -- 巡店筛选
    this.$track.trackSVFilter({
      current_page: '巡店记录',
      filter_type_sv: '日期查询',
      filter_option: this.dateName,
      filter_method: '切换维度',
    })
    
  }
  

  mounted() {
    if(this.$route.query.store_name) {
      this.storeName = String(this.$route.query.store_name)
      this.searchForm.data_type = 0
      this.searchForm.data_type_2 = Number(this.$route.query.data_type) || 1
      this.searchForm.start_time = Number(this.$route.query.start_time)
      this.searchForm.end_time = Number(this.$route.query.end_time)
      this.searchForm.third_id = Number(this.$route.query.third_id)
      this.dateTypeIndex = 3
    }
    this.$store.commit("storeVisit/setUpdate", true)
    this.getMonthList()
    this.getStoreTree()
  }
  activated() {
    
    if(this.isUpdate) {
      this.updateKey = Date.now()
      this.showSelectStorePopup = false;
      this.showSelectDatePopup = false;
      this.showDateRangePopup = false;
      if(this.$route.query.store_name) {
        this.storeName = String(this.$route.query.store_name)
        this.searchForm.data_type = 0
        this.searchForm.data_type_2 = Number(this.$route.query.data_type) || 1
        this.searchForm.start_time = Number(this.$route.query.start_time)
        this.searchForm.end_time = Number(this.$route.query.end_time)
        this.searchForm.third_id = Number(this.$route.query.third_id)
        this.dateTypeIndex = 3
        this.reload()
      } else {
        this.reset()
      }
      
      this.$store.commit("storeVisit/setUpdate", false)
    }

    // 
    const id = this.$store.getters["storeVisit/getUpdateId"];
    if(id) {
      const index = this.list.findIndex(item => item.id === id);
      // this.list[index].status = 3;
      this.list.splice(index, 1);
      this.$store.commit("storeVisit/getUpdateId", 0);
    }
    this.getStoreTree()
  }


  /**
   * 获取门店树
   */
  async getStoreTree() {
    const res: any = await getStoreTreeApi({
      has_children: true,
      ignore_close_store: true
    })
    if (res.code !== 200) {
      this.$toast.fail(res.message);
      return
    }
    
    this.storeTree = res.data.filter((item, index) => index<2)

  }

  getMonthList() {
    getMonthListApi().then(res => {
      this.monthList = res.data
    })
  }

  async getData() {

    const params = {
      ...this.meta,
      ...this.searchForm,
    }
    if(params.data_type == 1) {
      params.data_type = 0
    } else if(params.data_type == 0 && params.data_type_2 == 1) {
      params.data_type = 1
    } else {
      params.data_type = 2
    }
    delete params.data_type_2
    try {
      const {code, message, data}: IApiResult = await getStoreVisitListApi(params)
      if (code !== 200) {
        this.finished = true;
        this.loading = false;
        this.$toast.fail(message);
        return;
      }
      this.loading = false;
      this.list = this.meta.page == 1 ? data.list : this.list.concat(data.list)
      this.finished = this.list.length === data.total || data.list.length == 0;
      this.meta.page = this.meta.page + 1

    } catch (error) {
      console.log(error);
      this.finished = true;
      this.loading = false;
      this.list = []
    }
  }


  /**
   * @description 刷新
   */
  reload() {
    if(this.loading) return
    this.loading = true;
    this.meta.page = 1;
    this.finished = false
    this.list = [];
    this.getData()
  }

  reset() {
    this.storeName = ''
    this.dateTypeIndex = -1
    this.searchForm.third_id = 0
    this.searchForm.start_time = 0
    this.searchForm.end_time = 0
    this.searchForm.status = 1
    this.defaultDate = []
    this.reload()
  }

  /**
   * 切换tab
   */
  handleTabChange(value) {
    this.searchForm.status = value
    this.reload()

    // 神策数据上报 -- 巡店筛选
    this.$track.trackSVFilter({
      current_page: '巡店记录',
      filter_type_sv: '记录状态',
      filter_option: {1:'进行中', 2:'已完成'}[value],
      filter_method: '切换维度',
    })
  }

  /**
   * 选择日期
   */
  handleSelectDate() {
    this.showSelectDatePopup = true
  }

  /**
   * 选择日期确定
   */
  handleSelectDateConfirm(index) {
    this.showSelectDatePopup = false
    if(index === 3) {
      this.handleSelectDateRange()
    } else {
      this.dateTypeIndex = index
    }
  }

  /**
   * 选择日期取消
   */
  handleSelectDateCancel() {
    this.showSelectDatePopup = false
  }

  /**
   * 选择门店
   */
   handleSelectStore() {
    this.showSelectStorePopup = true
  }

  /**
   * 选择门店取消
   */
  handleSelectStoreCancel() {
    this.showSelectStorePopup = false
  }

  /**
   * 选择门店确定
   * @param value 
   */
  handleSelectStoreConfirm(value) {
    this.showSelectStorePopup = false
    const item = value[0]
    this.searchForm.third_id = item.third_id
    this.storeName = item.store_name
    this.reload()

    // 神策数据上报 -- 巡店筛选
    this.$track.trackSVFilter({
      current_page: '巡店记录',
      filter_type_sv: '门店查询',
      filter_option: item.store_name,
      filter_method: '切换维度',
    })
  }


  handleSelectChange(e) {
    const { track, item } = e
    // 神策数据上报 -- 选择门店
    this.$track.trackChooseStore({
      filter_option: '',
      search_word: '',
      choose_scene: '巡店记录筛选',
      choose_entrance: '门店筛选',
      store_name: item.store_name,
      store_code: item.store_code,
      ...track,
    })
  }

  handleSelectDateRange() {
    this.showDateRangePopup = true

  }
  handleSelectDateRangeCancel() {
    this.showDateRangePopup = false

  }

  handleSelectDateRangeConfirm(data:any[]) {
    this.defaultDate = data
    this.showDateRangePopup = false
    console.log(data[1]);
    this.searchForm.start_time = dayjs(`${data[0]} 00:00:00`).unix()
    this.searchForm.end_time =  dayjs(`${data[1]} 23:59:59`).unix()
    this.dateTypeIndex = 3
    this.reload()
  }

  handleDataTypeChange(type) {
    if(type == 1 && this.searchForm.data_type == 1) {
      this.searchForm.data_type_2 = 1
    }
    this.reload()
    if(type == 2) {
      // 神策数据上报 -- 巡店筛选
      this.$track.trackSVFilter({
        current_page: '巡店记录',
        filter_type_sv: '是否我发起的巡店',
        filter_option: this.searchForm.data_type_2 == 1 ? '是' : '否',
      filter_method: '切换维度',
      })
    } else {
      // 神策数据上报 -- 巡店筛选
      this.$track.trackSVFilter({
        current_page: '巡店记录',
        filter_type_sv: '是否我参与的巡店',
        filter_option: this.searchForm.data_type == 1 ? '是' : '否',
        filter_method: '切换维度',
      })
    }
  }

  handleDelete(index) {
    this.list.splice(index, 1)
  }
}
</script>

<style lang="scss" scoped>

.van-switch {
  width: 26px;
  height: 16px;
  border: none;
  background-color: rgba(173, 181, 189, 1);
  overflow: hidden;
}

::v-deep .van-switch__node {
  width: 12px;
  height: 12px;
  top: 50%;
  transform: translate(2px, -50%);
  box-shadow: none;
}

::v-deep .van-switch--on {
  background-color: rgba(13, 110, 253, 1);
  .van-switch__node {
    transform: translate(12px, -50%);
  }
}

.icon-down {
  width: 16px;
  height: 16px;
  background: url('~@/assets/store-visit/icon_caret_down.png') center/100% no-repeat;
  transition: all 0.3s ease;
  &.active {
    transform: rotate(180deg);
  }
}
.container { 
  background-color: rgba(248, 249, 250, 1);
  min-height: 100vh;

  .header {
    padding: 16px 12px;
    background-color: #fff;
    position: sticky;
    top: 0;
    z-index: 1;

    .flex {
      display: flex;
      align-items: center;
      justify-content: space-between;

      &:nth-child(2) {
        margin-top: 16px;
      }

      .flex {
        margin-top: 0;
      }
    }

    
    .switch {
      display: flex;
      align-items: center;
      margin-left: 16px;
      .label {
        font-size: 10px;
        font-weight: 700;
        margin-right: 4px;
      }
    }

    .tab {
      display: flex;
      align-items: center;
      height: 24px;
      border-radius: 4px;
      padding: 2px;
      background-color: rgba(248, 249, 250, 1);

      .tab-item {
        height: 20px;
        // width: 52px;
        padding: 0 8px;
        border-radius: 3px;
        color: rgba(0, 0, 0, 0.3);
        font-size: 12px;
        text-align: center;
        line-height: 20px;
        transition: all 0.3s;

        &.active {
          font-weight: 600;
          color: rgba(0, 0, 0, 1);
          background: rgba(255, 255, 255, 1);
        }

        &:not(:last-child) {
          margin-right: 2px;
        }
      }
    }

    .select-item {
      display: flex;
      align-items: center;
      font-size: 12px;
      color: rgba(0, 0, 0, 1);
      height: 32px;
      border-radius: 8px;
      background-color: rgba(248, 249, 250, 1);
      box-sizing: border-box;
      padding: 0 12px;
      margin-right: 8px;
      position: relative;

      .value {
        flex: 1;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
      .placeholder {
        flex: 1;
        color: rgba(153, 153, 153, 1);
      }

      &:nth-child(1) {
        width: 178px;
      }

      &:nth-child(2) {
        width: 125px;
      }
    }

    .reset {
      width: 32px;
      height: 24px;
      text-align: center;
      font-size: 12px;
      line-height: 24px;
      color: rgba(0, 0, 0, 1);

      &.disabled {
        color: rgba(153, 153, 153, 1);
      }
    }
  }

  .list {
    padding: 16px 12px;

    
  }
}

.popup {
  width: 210px;
  height: 240px;
  background: var(--Grey-100, rgba(248, 249, 250, 1));
  border: 1px solid var(--Grey-200, rgba(233, 236, 239, 1));
  box-shadow: 4px 4px 16px 0px rgba(0, 0, 0, 0.05);

  border-radius: 8px;
  position: absolute;
  top: 52px;
  left: 0;
  z-index: 3;
  box-sizing: border-box;
  padding: 8px 12px;

  &-mask {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0);
    z-index: 2;
  }

  .item {
    height: 56px;
    font-size: 12px;
    line-height: 56px;
    color: rgba(0, 0, 0,1);
    &:not(:last-child) {
      border-bottom: 1px solid var(--Grey-200, rgba(233, 236, 239, 1));
    }
  }
}
</style>