<template>
  <div v-if="detail.task_info" class="container">
    <!-- 任务详情 -->
    <div class="detail">
      <div class="flex-center">
        <div class="title">{{detail.task_info.display_name}}</div>
      </div>
      <div class="tips">
        <span>任务接收时间：{{ dateFormat(detail.task_info.created_at, 'yyyy年MM月dd日 hh:mm:ss') }}</span>
      </div>
      <div class="flex-center line">
        <span class="status" :data-status="detail.task_info.status">{{statusLabel}}</span>
        <span class="time">执行时间：{{ dateFormat(detail.task_info.start_time) }}-{{ dateFormat(detail.task_info.end_time) }} </span>
      </div>
      <div class="desc" :class="{
        show: expend,
      }">
        <span class="ellipse">{{ detail.task_info.content.length > 150 && !expend ? detail.task_info.content.slice(0, 150)+'...':detail.task_info.content }}</span>
        <span v-if="detail.task_info.content.length>150" class="more" @click="toggleExpend">{{expend?'收起':'全部'}}</span>
      </div>
      <div class="link">
        <div v-for="(item, index) in detail.task_info.class_url" :key="index" class="link-item flex-center">
          <div class="link-text">{{item.class_name}}</div>
          <div class="link-btn flex-center" @click="handleToLink(item)">去学习</div>
        </div>
      </div>
      <div class="countdown flex-center" :class="{ start: show(2) }">
        <van-count-down v-if="show(2)" :time="beforeStartTime" format="距开始DD天HH时mm分" />
        <van-count-down v-if="show(3)" :time="beforeEndTime" format="距结束DD天HH时mm分" />
      </div>
    </div>
    <!-- 筛选 -->
    <div class="header">
      <div class="tab">
        <div 
          class="tab-item" 
          :class="{ active: searchParams.finish_status === item.value }" 
          @click="handleTabChange(item.value)" 
          v-for="(item, index) in auditStatusOptions" 
          :key="index">
          <div class="tab-item__value">{{ countInfo[item.key]||0 }}</div>
          <div class="tab-item__label">{{ item.label }}</div>
        </div>
      </div>
      <div class="search">
        <div class="search-wrap">
          <van-search
            v-model="searchParams.name"
            clearable
            clear-trigger="always"
            placeholder="请输入门店名称"
            :left-icon="searchIcon"
            @input="handleSearch"
          />
          <div v-if="searchParams.name" class="search-clear" @click="handleClear"></div>
        </div>
        <div v-if="countInfo.fail_num !== 0" class="btn" @click="handleRemind()">一键提醒</div>
      </div>
    </div>
    <!-- 内容 -->
    <div class="content">
      <div class="table">
        <div class="table__header">
          <div class="table__row">
            <div class="table__col" v-for="(column, index) in columns" :key="index">{{ column.label }}</div>
          </div>
        </div>
        <div class="table__content">
          <div class="table__row" v-for="(item, index) in storeList" :key="index" @click="handleViewDetail(item)">
            <div class="table__col" v-for="(column, index) in columns" :key="index">
              <div v-if="column.prop === 'store_name'" class="title">{{ item[column.prop] }}</div>
              <div v-else-if="column.prop === 'finish_status'">
                <div class="status-label" :data-status="item[column.prop]">{{ auditStatusLabel(item[column.prop]) }}</div>
                <div v-if="item[column.prop] === 2" class="status-btn" @click.stop="handleRemind(item)">去提醒 ></div>
              </div>
              <div v-else-if="column.prop === 'grade'">
                {{ item.finish_status === 3||item.finish_status === 4 ? item[column.prop]:'--' }}
              </div>
              <span v-else>{{ item[column.prop]||'--' }}</span>
            </div>
          </div>
          <empty-page
            v-if="storeList.length == 0"
            picture="noDisplayTask"
            paddingTop="60"
            paddingBottom="0"
            text="暂无门店信息"
          ></empty-page>
        </div>
      </div>
    </div>
    <div class="footer">
      <div class="btn" @click="handleView">预览他店陈列图</div>
    </div>
    <!-- 确认提交 -->
    <confirm-popup ref="remindPopup" title="确认提醒" @confirm="handleRemindConfirm">
      <div class="text-align-left">点击确定后，会向【未通过】门店导购发送任务通知</div>
    </confirm-popup>
    <select-store ref="selectStore" @confirm="handleViewConfirm"></select-store>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Watch } from "vue-property-decorator";
import { CountDown } from 'vant';
import ConfirmPopup from "./components/confirm-popup.vue";
import SelectStore from "./components/select-store.vue";
import EmptyPage from "@/components/base/empty-page.vue";
import { auditStatusOptions, statusOptions } from "./const";
import { getTaskInfoApi, getListTipsApi, getCheckStoreListApi, sendNoticeApi, getAuthApi } from '@/api/display-task'
import { getRedirectParams } from "@/utils/welearning";
import wx from 'weixin-js-sdk';
interface ICheckStoreItem {
  /** 分配id */
  allot_id: number;
  /** 审核状态，1为待审核，2为未通过，3为已通过 */
  finish_status: number;
  /** 评分 */
  grade: number;
  /** 任务id */
  id: number;
  /** 门店id */
  store_id: number;
  /** 门店code */
  store_code: string;
  /** 门店名称 */
  store_name: string;
}
@Component({
  components: {
    [CountDown.name]: CountDown,
    ConfirmPopup,
    SelectStore,
    EmptyPage
  },
  beforeRouteLeave(to, from, next) {
    this.scrollTop = document.documentElement.scrollTop;
    next()
  }
})
export default class DisplayTaskDetail extends Vue {
  /** 搜索图标 */
  searchIcon = require('@/assets/display-task/icon_search.png')
  /** 门店列表 */
  storeList: ICheckStoreItem[] = []
  /** 搜索参数 */
  searchParams = {
    name: '',
    finish_status: -1
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
  /** 表格列 */
  columns = [
    {
      label: '门店编码',
      prop: 'store_code',
    },
    {
      label: '门店名称',
      prop: 'store_name',
    },
    {
      label: '状态',
      prop: 'finish_status',
    },
    {
      label: '评分',
      prop: 'grade',
    }
  ]
  /** 当前去提醒的门店 */
  currentStore: ICheckStoreItem = null
  scrollTop = 0
  expend = false

  /** 任务详情 */
  get detail(): IDisplayTaskDetail {
    return this.$store.getters["displayTask/getDetail"]
  }

  /** 任务状态 */
  get statusLabel() {
    return statusOptions.find((item) => item.value === this.detail.task_info.status).label
  }

  /** 剩余多少时间开始 */
  get beforeStartTime() {
    return this.detail.task_info.start_time * 1000 - Date.now()
  }

  /** 剩余多少时间结束 */
  get beforeEndTime() {
    return this.detail.task_info.end_time * 1000 - Date.now()
  }

  /** 审核状态选项 */
  get auditStatusOptions() {
    return auditStatusOptions.filter(item => item.value !== 0)
  }

  /** 审核状态 */
  get auditStatusLabel() {
    return (value: number) => {
      return auditStatusOptions.find((item) => item.value === value).label
    }
  }

  /** 控制页面区域展示 */
  get show() {
    const { status } = this.detail.task_info
    return (key: number) => {
      switch (key) {
        case 2:
          // 任务未开始
          return status === 0
        case 3:
          // 任务进行中
          return status === 1
        default:
          return true
      }
    }
  }


  /** 用户信息 */
  get updateInfo() {
    return this.$store.getters["displayTask/getUpdateList"]
  }

  /** 用户信息 */
  get userInfo(): IUserInfo {
    return this.$store.getters["user/getUserInfo"] || {}
  }

  @Watch("detail", { deep: true })
  onDetailChange(value: IDisplayTaskDetail) {
    this.$store.commit("displayTask/setDetail", value)
  }

  async mounted() {
    if(this.userInfo.app_source !== 1) {
      await this.getAuth()
    }
    this.searchParams.finish_status = Number(this.$route.query.finish_status) || -1
    this.getDetail()
    this.getListTips()
    this.getCheckStoreList()
    if(this.$route.query.from !== '1') {
      // 神策数据上报 -- 进入VM
      this.$track.trackEnterDisplayTask({
        vm_entrance: 2
      })
    }
  }
  async activated() {
    if(!this.detail.task_info) {
      this.searchParams.finish_status = Number(this.$route.query.finish_status) || -1
      this.scrollTop = 0
      this.getDetail()
      await this.getListTips()
      this.getCheckStoreList()
    }
    if(this.updateInfo.updateDetail) {
      this.$toast.success(this.updateInfo.message || '审核完成')
      await this.getListTips()
      this.getCheckStoreList()
      this.$store.commit('displayTask/setUpdateList', {
        ...this.updateInfo,
        updateDetail: false,
      })
    }

    this.$nextTick(() => {
      if (this.scrollTop > 0) {
        console.log('activated', this.scrollTop);
  
        document.documentElement.scrollTo({
          top: this.scrollTop,
          behavior: "auto"
        });
      }
    })
    if(this.$route.query.from !== '1') {
      window.addEventListener("popstate", this.popstate, false);
    }
  }
  deactivated() {
    if(this.$route.query.from !== '1') {
      window.removeEventListener("popstate", this.popstate, false);
    }
  }
  /**生命周期 */
  private async created() {
    // 当前需要直接退出的页面
    if(this.$route.query.from !== '1') {
      window.history.pushState("_android_back", "", location.href);
      window.addEventListener("popstate", this.popstate, false);
    }
  }
  private async destroyed() {
    // 记得清除!!!!不然单页面，其它页面也可以用
    if(this.$route.query.from !== '1') {
      window.removeEventListener("popstate", this.popstate, false);
    }
  }

  popstate() {
    wx.closeWindow(); // 微信网页退出
  }
  toggleExpend() {
    this.expend = !this.expend;
  }
  /**
   * @description 格式化时间
   * @param time 
   * @param format 
   */
  dateFormat(time: number, format = "yyyy/MM/dd") {
    return this.$helper.dateFormat(time, format);
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
    } catch (error) {
      console.log(error);
    }
  }

  /**
   * @description 获取详情
   */
  async getDetail() {
    try {
      const params = {
        id: Number(this.$route.query.id),
      }
      const { code, message, data }: IApiResult<ITaskInfo> = await getTaskInfoApi(params)
      
      if(code!== 200) {
        this.$toast.fail(message)
        return
      }
      // 这个是必须的
      data.template_info = data.template_info.map(item => {
        return {
          ...item,
          image_list: []
        }
      })
      this.$store.commit("displayTask/setDetail", {
        ...this.detail,
        task_info: data
      })
    } catch (error) {
      console.log(error);
    }
  }

  /** 
   * @description 获取角标数据
   */
  async getListTips() {
    try {
      const params = {
        id: Number(this.$route.query.id),
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
  handleSearch = this.$helper.throttle(() => {
    this.getCheckStoreList()
  })
  /** 
   * @description 获取角标数据
   */
  async getCheckStoreList() {
    try {
      const params = {
        id: Number(this.$route.query.id),
        ...this.searchParams
      }
      const {code, message, data}:IApiResult<ICheckStoreItem[]> = await getCheckStoreListApi(params)
      if (code!== 200) {
        this.$toast.fail(message);
        return;
      }
      this.storeList = data
    } catch (error) {
      console.log(error);
    }
  }

  /**
   * @description 清空搜索
   */
  handleClear() {
    this.searchParams.name = ''
    this.getCheckStoreList()
  }

  /**
   * 状态切换
   * @param value 
   */
  handleTabChange(value: number) {
    this.searchParams.finish_status = value;
    this.getCheckStoreList()
  }

  /**
   * @description 去学习
   * @param item 
   */
  async handleToLink(item: IClassUrlItem) {
    const redirect = getRedirectParams(item.url)
    const uri = await redirect.uri()
    // 神策数据上报 -- 课程学习
    this.$track.trackCourseLearing({
      course_title: item.class_name
    })
    location.href = uri;
  }

  /**
   * @description 一键提醒
   */
  handleRemind(item?: ICheckStoreItem) {
    this.currentStore = item || null;
    (this.$refs.remindPopup as any).showPopup();
  }

  /**
   * @description 确认提醒
   */
  async handleRemindConfirm() {
    try {
      const params = {
        id: this.detail.task_info.id,
        allot_id: 0
      }
      if(this.currentStore) {
        params.allot_id = this.currentStore.allot_id
      }
      const { code, message }:IApiResult = await sendNoticeApi(params)
      if (code!== 200) {
        this.$toast.fail(message);
        return;
      }
      this.$toast('已提醒');
      // 神策数据上报 -- 陈列提醒（vm视角）
      this.$track.trackDisplayReminder({
        reminder_type: this.currentStore ? 2 : 1
      })

    } catch (error) {
      console.log(error);
      
    }
  }


  /**
   * @description 查看店铺
   * @param item 
   */
  handleView() {
    (this.$refs.selectStore as any).showPopup({
      id: this.detail.task_info.id,
    })
  }

  /**
   * @description 他店陈列
   */
  handleViewConfirm(value: any) {
    this.$store.commit('displayTask/setDetail', {
      ...this.detail,
      store_info: null
    });
    this.$router.push({
      name: 'display-task-preview',
      query: {
        ...value,
        id: this.detail.task_info.id,
        is_vm: '1',
        other: '1'
      }
    })
  }

  /** 
   * @description 查看详情
   */
  handleViewDetail(item: ICheckStoreItem) {
    this.$store.commit('displayTask/setDetail', {
      ...this.detail,
      store_info: null
    });
    if(item.finish_status === 1 || item.finish_status === 4) {
      this.$router.push({
        name: 'display-task-audit',
        query: {
          id: item.id + '',
          store_id: item.store_id + '',
          from: '1'
        }
      })
    } else {
      this.$router.push({
        name: 'display-task-preview',
        query: {
          id: item.id + '',
          store_id: item.store_id + '',
          is_vm: '1'
        }
      })
    }
  }
  
  /**
   * @description 返回
   */
  handleBack() {
    this.$router.go(-1)
  }
}
</script>
<style lang="scss" scoped>
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
.container {
  background-color: #F7F7F7;
  min-height: 100vh;
  box-sizing: border-box;
  padding-bottom: calc(56px + constant(safe-area-inset-bottom)); /** iOS 11.0+ */
  padding-bottom: calc(56px + env(safe-area-inset-bottom)); /** iOS 11.0+ */
  .detail {
    padding: 15px;
    background-color: #fff;
    margin-bottom: 15px;
    .title {
      font-size: 16px;
      font-weight: 600;
      color: #000000;
      line-height: 22px;
      display: flex;
      flex: 1;
    }

    .tips {
      font-size: 12px;
      font-weight: 400;
      color: #A6A6A6;
      line-height: 17px;
      margin-top: 10px;
      position: relative;
    }

   .line {
      width: 345px;
      height: 37px;
      background: rgba(119, 119, 119, 0.1);
      justify-content: flex-start;
      padding: 0 15px;
      margin-top: 10px;

      .status {
        font-size: 12px;
        font-weight: bold;
        margin-right: 10px;
        &[data-status="0"] {
          color: #C6AD7F;
        }
        &[data-status="1"] {
          color: #70A560;
        }
        &[data-status="2"] {
          color: #A6A6A6;
        }
      }
      .time {
        font-size: 12px;
        font-weight: 400;
        color: #000000;
        line-height: 17px;
      }
    }

    .desc {
      font-size: 12px;
      font-weight: 400;
      color: #777777;
      line-height: 17px;
      margin-top: 10px;
      text-align: justify;
      word-break: break-all;
      .more {
        color: #A1A1A1;
        font-size: 12px;
        font-weight: 400;
        line-height: 17px;
        text-align: right;
        cursor: pointer;
        transition: all .3s;
        padding-left: 10px;
        &::after {
          content: "";
          margin-left: 5px;
          display: inline-block;
          width: 14px;
          height: 9px;
          background: url('~@/assets/display-task/icon_down.png') right center/14px 9px no-repeat;
        }
      }

      &.show .more::after {
        transform: rotate(180deg);
      }
    }

    .link {
      margin-top: 15px;

      &-item {
        border: 1px solid #A6A6A6;
        box-sizing: border-box;
        padding: 10px 15px;

        &:not(:last-child) {
          margin-bottom: 15px;
        }
      }

      &-text {
        font-size: 12px;
        font-weight: 400;
        color: #000000;
        line-height: 17px;
        flex: 1;
        padding-right: 15px;
        word-break: break-all;
      }

      &-btn {
        width: 62px;
        height: 24px;
        background: #000000;
        font-size: 11px;
        font-weight: 400;
        color: #FFFFFF;
        line-height: 17px;
      }
    }

    .countdown {
      margin-top: 15px;
      .van-count-down {
        line-height: 18px;
        font-size: 12px;
        color: #70A560;
      }

      &.start .van-count-down {
        color: #C6AD7F;
      }
    }
  }
  
  .header {
    padding: 15px;
    background: #FFFFFF;
    box-shadow: 0px 3px 17px 0px rgba(54,68,58,0.06);
    margin-bottom: 15px;
    // position: sticky;
    // top: 0;
    z-index: 1;
    .tab {
      display: flex;
      justify-content: space-between;
      align-items: center;

     .tab-item {
        width: 78px;
        height: 60px;
        background: #fff;
        border: 1px solid #000000;
        display: flex;
        align-items: center;
        justify-content: center;
        flex-direction: column;
        transition: all 0.2s;
        margin-left: -1px;

        &__value {
          font-size: 16px;
          font-weight: 500;
          color: #000000;
          line-height: 22px;
        }

        &__label {
          font-size: 12px;
          font-weight: 400;
          color: #777777;
          line-height: 17px;
          margin-top: 2px;
        }

        &.active {
          background: #000000;
          .tab-item__value,
          .tab-item__label {
            color: #fff;
          }
        }
      }
    }

    .search {
      margin-top: 15px;
      display: flex;
      align-items: center;
      
      .btn {
        width: 88px;
        height: 38px;
        background: #000000;
        border-radius: 19px 0px 0px 19px;
        font-size: 14px;
        font-weight: 400;
        color: #FFFFFF;
        display: flex;
        align-items: center;
        justify-content: center;
        margin-left: 15px;
        margin-right: -15px;
      }
    }
  }

  .content {
    background-color: #fff;

    .table {
      font-size: 14px;

      &__row {
        padding: 15px 0;
        display: flex;
        align-items: center;
      }

      &__col {
        box-sizing: border-box;
        padding-right: 20px;
        word-break: break-all;
        line-height: 20px;
        &:nth-child(1) {
          width: 80px;
        }
        &:nth-child(2) {
          flex: 1;
        }
        &:nth-child(3) {
          width: 80px;
        }
        &:nth-child(4) {
          width: 70px;
        }
      }

      &__header {
        padding: 0 15px;
        height: 50px;
        background: #777777;
        color: #fff;
        box-shadow: 0px 3px 17px 0px rgba(54,68,58,0.06);
      }

      &__content {
        color: #000;
        padding: 0 15px;

        .table__row {
          position: relative;

          &::before {
            content: "";
            position: absolute;
            left: 0;
            bottom: 0;
            width: 345px;
            height: 1px;
            border: 1px solid #F2F2F2;
          }
          &::after {
            content: "";
            position: absolute;
            right: 0;
            top: 50%;
            transform: translateY(-50%);
            width: 9px;
            height: 14px;
            background: url("~@/assets/display-task/icon_right.png") no-repeat center/100%;
          }

        }
      }


    }

    .status-label {
      &[data-status="1"],
      &[data-status="4"] {
        color: #C6AD7F;
      }
      &[data-status="2"] {
        color: #E02020;
      }
      &[data-status="3"] {
        color: #70A560;
      }
    }

    .status-btn {
      font-size: 11px;
      font-weight: 400;
      color: #000000;
      line-height: 16px;
      margin-top: 3px;
    }

    .title {
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 3;
      line-clamp: 3;
      -webkit-box-orient: vertical;
    }
  }

  .footer {
    position: fixed;
    bottom: 0;
    width: 100vw;
    background-color: #fff;
    padding: 10px 30px;
    padding-bottom: calc(10px + constant(safe-area-inset-bottom)); /** iOS 11.0+ */
    padding-bottom: calc(10px + env(safe-area-inset-bottom)); /** iOS 11.0+ */
    left: 0;
    z-index: 1;

    .btn {
      height: 36px;
      background: #FFFFFF;
      border: 1px solid #000000;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 14px;
      font-weight: 400;
      color: #000000;
    }
  }
}
</style>