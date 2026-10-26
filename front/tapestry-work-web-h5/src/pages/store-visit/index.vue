<template>
  <div class="container">
    <div class="header">
      <div class="flex">
        <div class="user">
          <van-image class="avatar" :src="userInfo.avatar"></van-image>
          <div class="name">{{userInfo.name}}</div>
          <div class="dept">{{ userInfo.dept_name }}</div>
        </div>
        <div class="lang" :class="{zh: lang === 'zh-CN'}" @click="handleClickLang(lang)"></div>
      </div>
    </div>
    <div class="content">
      <div class="section">
        <calendar-card v-model="current" :schedules="schedules" @change-month="handleMonthChange">
          <!-- <div v-if="userInfo.is_store_visit_manager || userInfo.is_store_manager" class="btn" @click="handleReservation()">
            <span>{{ $t('storeVisit.btn.1') }}</span>
            <div class="icon-1"></div>
          </div> -->
        </calendar-card>
        <!-- 巡店记录 -->
        <div class="title">
          <div class="text">{{ $t('storeVisit.title.1') }}</div>
          <div class="record" @click="handleViewRecord">
            <span>{{ $t('storeVisit.btn.2') }}</span>
            <div class="icon-right"></div>
          </div>
        </div>
        <!-- 巡店记录列表 -->
        <div v-if="currentList.length" class="card" @click="handleShowRecord">
          <div class="card-inner">
            <div class="card-title">{{currentList[0].store_name}} · {{ $t('storeVisit.title.2') }}</div>
            <div class="flex">
              <div class="card-time">{{currentList[0].time_span}}</div>
              <div class="icon-date"></div>
            </div>
          </div>
        </div>
        <div v-else class="empty">{{ $t('storeVisit.empty.1') }}</div>
      </div>
      <!-- 巡店管理员或者店长可以发起巡店 -->
      <div class="btn-group">
        <div v-if="userInfo.is_store_visit_manager || userInfo.is_store_manager" class="btn" @click="handleReservation()">
          <div class="icon-1"></div>
          <span>{{ $t('storeVisit.btn.1') }}</span>
        </div>
        <div v-if="userInfo.is_store_visit_manager||userInfo.is_store_manager" class="start" @click="handleStartVisit()">
          <div class="icon-2"></div>
          <div class="start-text">{{ $t('storeVisit.btn.3') }}</div>
          <!-- <div class="icon-right-white"></div> -->
        </div>
      </div>
      <div class="tips">
        <div>{{ $t('storeVisit.tips.40') }}</div>
        <div>{{ $t('storeVisit.tips.41') }}</div>
      </div>
    </div>
    <div v-if="userInfo.is_store_visit_manager" class="footer">
      <div class="tabbar">
        <div class="tabbar-item active">{{ $t('storeVisit.tabs.1') }}</div>
        <div class="tabbar-item" @click="handleNav">{{ $t('storeVisit.tabs.2') }}</div>
      </div>
    </div>
    <select-record-popup v-model="showRecordPopup" :title="current.replaceAll('-', '.')" :list="currentList" @confirm="handleRecordConfirm"></select-record-popup>
    <select-store-popup v-model="showSelectStorePopup" :storeTre="storeTree" :multiple="true" @select="handleSelectChange" @confirm="handleSelectStoreConfirm"></select-store-popup>

  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { Button, Dialog, Image as VanImage } from "vant";
import SelectClerkPopup from "./components/select-clerk-popup.vue";
import SelectRecordPopup from "./components/select-record-popup.vue";
import SelectFeedbackPopup from "./components/select-feedback-popup.vue";
import SelectVisitorPopup from "./components/select-visitor-popup.vue";
import SelectStorePopup from "./components/SelectStorePopup/index.vue";
import CalendarCard from "./components/calendar-card.vue";
import ComConfirmPopup from "./components/ComConfirmPopup";
import SelectStore from "./components/SelectStorePopup";
import SelectRecord from "./components/SelectRecordPopup";
import ComCreatePopup from "./components/ComCreatePopup";
import { getConfigApi, getStoreVisitAuthApi, getStoreVisitListApi, saveConfigApi, startVisitApi } from "@/api/store-visit";
import dayjs from "dayjs";
import wx from 'weixin-js-sdk'
import { getStoreTreeApi } from "@/api/display-task";

import { Locale } from 'vant';
// 引入英文语言包
import enUS from 'vant/es/locale/lang/en-US';
import zhCN from 'vant/es/locale/lang/zh-CN';
import { i18n } from "@/locales";

@Component({
  components: {
    [Button.name]: Button,
    [VanImage.name]: VanImage,
    SelectClerkPopup,
    SelectRecordPopup,
    SelectFeedbackPopup,
    SelectVisitorPopup,
    SelectStorePopup,
    CalendarCard
  },
})
export default class StoreVisitIndex extends Vue {
  /** 当前时间 */
  current = dayjs().format('YYYY-MM-DD')
  /** 搜索条件 */
  searchForm = {
    status: 1,
    start_time: 0,
    end_time: 0
  }
  /** 巡店记录列表 */
  list: StoreVisit.IDetail[] = []
  /** 显示巡店记录 */
  showRecordPopup = false
  /** 显示选择门店弹窗 */
  showSelectStorePopup = false;
  /** 显示确认弹窗 */
  showConfirmPopup = false;
  /** 当前选中的门店ID */
  currentThirdId = 0;
  /** 当前选中的门店 */
  currentStore = null
  /** 是否 多选门店 */
  multiple = false

  /** 上报场景 */
  trackScene = ''
  /** 上报入口 */
  trackEntrance = ''
  storeTree = []

  // 当前语言
  get lang () {
    return this.$store.getters["lang/getLang"]
  }

  /** 巡店记录日期列表 */
  get schedules() {
    return this.list.map(item => item.date)
  }

  /** 今天的巡店记录 */
  get currentList() {
    return this.list.filter(item => item.date === this.current)
  }

  /** 用户信息 */
  get userInfo():IUserInfo {
    return this.$store.getters["user/getUserInfo"]
  }

  /**
   * 创建人
   */
  get isCreator() {
    return (item) => item.creator.member_id === this.userInfo.member_id
  }


  get currentItem () {
    return (third_id) => {
      const current = this.list.find(item => item.third_id === third_id &&  this.isCreator(item))
      return current
    }
  }

  get allList() {
    return this.list.filter(item => this.isCreator(item))
  }

  get recordByStore() {
    return (third_id) => {
      return this.allList.filter(item => item.third_id === third_id)
    }
  }

  /**生命周期 */
  async created() {
    await this.getConfig()
    // 当前需要直接退出的页面
    window.history.pushState("_android_back", "", location.href);
    window.addEventListener("popstate", this.popstate, false);
    
    const [year, month] = this.current.split('-').map(Number);
    // 开始时间为每月的第一天00:00:00
    this.searchForm.start_time = dayjs(new Date(year, month - 1, 1, 0, 0, 0)).unix()
    // 结束时间为每月的最后一天23:59:59
    this.searchForm.end_time = dayjs(new Date(year, month, 0, 23, 59, 59)).unix()
    console.log(this.userInfo);

  }
  destroyed() {
    // 记得清除!!!!不然单页面，其它页面也可以用
    window.removeEventListener("popstate", this.popstate, false);
    
  }

  activated() {
    // 当前需要直接退出的页面
    window.addEventListener("popstate", this.popstate, false);
    this.getData();
    const toast = this.$store.getters["storeVisit/getToast"]
    if (toast) {
      this.$toast(toast);
      this.$store.commit("storeVisit/setToast", '');
    }
  }
  deactivated() {
    // 记得清除!!!!不然单页面，其它页面也可以用
    window.removeEventListener("popstate", this.popstate, false);
  }

  popstate() {
    wx.closeWindow(); // 微信网页退出
  }
  async mounted() {
    // this.$toast({
    //   duration: 0,
    //   message: this.$t('common.toast.13') as string,
    //   forbidClick: true,
    // })
    await this.getAuth()
    this.getData();
    this.getStoreTree()
  }
  async getConfig() {
    const res = await getConfigApi({
      key: 'language'
    });
    console.log(res.data);
    console.log(JSON.parse(res.data.value));
    this.changeLang(res.data.member_config === 'en' ? 'zh-CN' : 'en-US')

    // const { en, zh } = JSON.parse(res.data.value)
    // i18n.setLocaleMessage('zh-CN', zh)
    // i18n.setLocaleMessage('en-US', en)
  }
  /**
   * @description 获取用户权限
   */
  async getAuth() {
    try {
      const {code, message, data}:IApiResult = await getStoreVisitAuthApi({})
      if (code!== 200) {
        this.$toast.fail(message);
        return;
      }
      this.$store.commit("user/setUserInfo", {
        ...this.userInfo, 
        ...data,
        // is_store_visit_manager: true
      })
    } catch (error) {
      console.log(error);
    }
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

  /** 获取巡店记录 */
  async getData() {
    const res:any = await getStoreVisitListApi({
      // ...this.searchForm,
      status: 1,
      start_time: 0,
      end_time: 0,
      page: 1,
      size: 999
    });
    if(res.code === 200) {
      this.list = res.data.list;
    }
  }
  /**
   * 开始巡店
   */
  async startVisit(params) {

    const res:IApiResult = await startVisitApi(params)

    if(res.code === 200) {
      return res.data
    }
  }

  /** 切换月份 */
  handleMonthChange(e) {
    const { year, month } = e
    // 开始时间为每月的第一天00:00:00
    this.searchForm.start_time = dayjs(new Date(year, month - 1, 1, 0, 0, 0)).unix()
    // 结束时间为每月的最后一天23:59:59
    this.searchForm.end_time = dayjs(new Date(year, month, 0, 23, 59, 59)).unix()
    this.getData()
  }

  /** 
   * 查看记录
   */
  handleShowRecord() {
    this.showRecordPopup = true


    // 神策数据上报 -- 日历页点击
    this.$track.trackCalendarPageClick({
      button_type: '巡店列表',
      button_name: '巡店列表',
    })
  }
  /**
   * 使用该巡店计划开始巡店
   * @param data 
   */
  handleRecordConfirm(data) {
    this.showRecordPopup = false
    if(this.isCreator(data)) {
      this.navigateToStart({
        id: data.id,
        entry: '巡店计划',
        btn: '巡店计划',
      })
    } else {
      this.$nextTick(() => {
        this.$router.push({
          name: 'store-visit-report',
          query: {
            id: String(data.id),
            from: '1',
          }
        })
      })
    }
  }

  /**
   * 巡店预约
   */
  async handleReservation() {

    // 神策数据上报 -- 日历页点击
    this.$track.trackCalendarPageClick({
      button_type: '预约巡店',
      button_name: '预约巡店',
    })

    if(this.userInfo.is_store_visit_manager) {
      // 巡店管理员，选择门店 - 单选
      // this.showSelectStorePopup = true
      // this.multiple = true
      // this.trackScene = '预约'
      // this.trackEntrance = '预约巡店'

      const res = await SelectStore.open({
        multiple: true,
        storeTree: this.storeTree,
        defaultValue: [],
        onSelect: (e) => {
        const { track, item } = e
          // 神策数据上报 -- 选择门店
          this.$track.trackChooseStore({
            filter_option: '',
            search_word: '',
            choose_scene: '预约',
            choose_entrance: '预约巡店',
            store_name: item.store_name,
            store_code: item.store_code,
            ...track,
          })
        }
      })

      if(res.action === 'confirm') {
        const store = res.data.map(item => {
          return {
            third_id: item.third_id,
            store_name: item.store_name,
            store_code: item.store_code,
          }
        })
        this.$nextTick(() => {
          this.$router.push({
            path: '/store-visit-create',
            query: {
              store_list: JSON.stringify(store),
            }
          })
        })
      }


    } else {
      // 店长只能巡店自己的门店 - 预约自己的门店
      this.handleClick([
        {
          third_id: this.userInfo.third_id,
          store_name: this.userInfo.store_name,
          store_code: this.userInfo.store_code,
        }
      ])
    }
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

    if(this.multiple) {
      // 多选 的时候 去预约
      this.handleClick(value)
      return
    }
    console.log(1111111, value[0]);
    
    // 单选 
    const { third_id } = value[0]
    this.currentStore = {...value[0]}
    if(this.currentItem(third_id)) {
      // 是否已存在巡店预约记录，已存在，弹框确认是否选择当前记录开始巡店 
      this.currentThirdId = third_id
      this.handleConfirm()
    } else {
      // 否则 直接开始巡店
      this.navigateToStart({
        third_id: third_id,
        entry: '发起巡店',
        btn: '完成',
      })
    }
  }

  handleSelectChange(e) {
    const { track, item } = e
    // 神策数据上报 -- 选择门店
    this.$track.trackChooseStore({
      filter_option: '',
      search_word: '',
      choose_scene: this.trackScene,
      choose_entrance: this.trackEntrance,
      store_name: item.store_name,
      store_code: item.store_code,
      ...track,
    })
  }

  handleConfirm() {
    ComConfirmPopup.confirm({
      title: this.$t('storeVisit.confirm.title.1'),
      content: this.$t('storeVisit.confirm.content.3'),
      confirmButtonText: this.$t('storeVisit.confirm.confirmBtn.3'),
      cancelButtonText: this.$t('storeVisit.confirm.cancelBtn.3'),
    }).then((res) => {
      if(res.action === 'confirm') {
        // 使用该巡店计划开始巡店
        this.navigateToStart({
          id: this.currentItem(this.currentThirdId).id,
          entry: '发起巡店',
          btn: '使用该巡店计划继续',
        })

      } else {
        // 发起新的巡店记录
        this.navigateToStart({
          third_id: this.currentThirdId,
          entry: '发起巡店',
          btn: '发起新的巡店',
        })
      }
    }).catch(() => {
      console.log('关闭');
    })
  }

  /**
   * 查看记录
   */
  handleViewRecord() {
    // 神策数据上报 -- 日历页点击
    this.$track.trackCalendarPageClick({
      button_type: '巡店记录',
      button_name: '巡店记录',
    })
    this.$store.commit("storeVisit/setUpdate", true)
    this.$router.push({
      path: '/store-visit-record',
    })

  }

  /**
   * 跳转巡店
   * @param data 
   */
  async navigateToStart(data) {
    
    if(!data.id) {
      const res:any = await ComCreatePopup.open({
        store_list: [this.currentStore],
        is_store_manager: data.is_store_manager,
        storeTree: this.storeTree,
      })
      if(res.action == 'confirm') {
        data.id = await this.startVisit(res.data)
      } else {
        return
      }
      
    }
    this.$nextTick(() => {
      this.$router.push({
        path: '/store-visit-start',
        query: data
      })
    })
  }

  /**
   * 预约巡店
   */
  handleClick(value) {
    const store = value.map(item => {
      return {
        third_id: item.third_id,
        store_name: item.store_name,
        store_code: item.store_code,
      }
    })
    this.$nextTick(() => {
      this.$router.push({
        path: '/store-visit-create',
        query: {
          store_list: JSON.stringify(store),
        }
      })
    })
  }

  handleNav() {

    // 神策数据上报 -- tab切换
    this.$track.trackTabSwitch({
      tab_name: '看板'
    })

    this.$router.replace({
      path: '/store-visit-dashboard',
    })
  }

  /**
   * 发起巡店 - 巡店发起人
   * 1. 判断是否有未开始的巡店记录
   *    1.1 有未开始的巡店记录 - 显示提示弹框展示未开始的巡店记录
   *        1.1.1 选择未开始的巡店记录 - 使用该巡店计划继续
   *        1.1.2 发起新的巡店 - 关闭提示弹框，并选择门店
   *    1.2 无未开始的巡店记录 - 选择门店
   * 2. 选择门店 - 判断当前所选门店是否有未开始的巡店记录
   *     2.1 有未开始的巡店记录 - 显示提示弹框展示未开始的巡店记录
   *        2.1.1 选择未开始的巡店记录 - 使用该巡店计划继续
   *        2.1.2 发起新的巡店 - 关闭提示弹框，并选择门店
   *     2.2 无未开始的巡店记录 - 开始巡店
   */
  async handleStartVisit() {
    // 神策数据上报 -- 日历页点击
    this.$track.trackCalendarPageClick({
      button_type: '发起巡店',
      button_name: '发起巡店',
    })

    if(this.userInfo.is_store_visit_manager) {
      
      const res1 = this.allList.length ? await SelectRecord.open({
        list: this.allList,
        tips: `${this.$t('storeVisit.popup.tips.1')} ${this.allList.length} ${this.$t('storeVisit.popup.tips.2')}`
      }) : {
        action: 'cancel',
        data: null
      }
      
      if(res1.action === 'cancel') {
        // 1.2
        const res2 = await SelectStore.open({
          multiple: false,
          storeTree: this.storeTree,
          defaultValue: [],
          onSelect: (e) => {
          const { track, item } = e
            // 神策数据上报 -- 选择门店
            this.$track.trackChooseStore({
              filter_option: '',
              search_word: '',
              choose_scene: '直接发起',
              choose_entrance: '发起巡店',
              store_name: item.store_name,
              store_code: item.store_code,
              ...track,
            })
          }
        })
        console.log(res2);
        if(res2.action === 'confirm') {
          // 2
          const record = this.recordByStore(res2.data[0].third_id)
          if(record.length) {
            // 2.1
            const res3 = await SelectRecord.open({
              list: record,
              tips: `${this.$t('storeVisit.popup.tips.3')} ${record.length} ${this.$t('storeVisit.popup.tips.4')}`
            })
            if(res3.action === 'confirm') {
              // 2.1.1
              this.navigateToStart({
                id: res3.data.id,
                entry: '发起巡店',
                btn: '二节点使用该巡店计划继续',
              })
            } else if(res3.action === 'cancel') {
              // 2.1.2
              this.currentStore = {...res2.data[0]}
              this.navigateToStart({
                third_id: res2.data[0].third_id,
                entry: '发起巡店',
                btn: '发起新的巡店',
              })

            }

          } else {
            // 2.2
            // 店长只能巡店自己的门店
            this.currentStore = {...res2.data[0]}
            this.navigateToStart({
              third_id: res2.data[0].third_id,
              entry: '发起巡店',
              btn: '完成',
            })
          }
        }

      } else if(res1.action === 'confirm') {
        // 1.1
        this.navigateToStart({
          id: res1.data.id,
          entry: '发起巡店',
          btn: '一节点使用该巡店计划继续',
        })
        
      }
      
      

    } else {
      // 店长只能巡店自己的门店
      this.currentStore = {...this.userInfo}
      this.navigateToStart({
        third_id: this.userInfo.third_id,
        is_store_manager: true,
        entry: '发起巡店',
        btn: '开始巡店',
      })
    }
  }

  handleClickLang(lang) {
    this.changeLang(lang)
    setTimeout(() => {
      saveConfigApi({})
    }, 100)
  }

  changeLang(lang) {
    console.log(lang);
    
    if(lang === 'zh-CN') {
      this.$i18n.locale = 'en-US'
      this.$store.dispatch('lang/setLang', 'en-US')
      Locale.use('en-US', enUS);

    } else {
      this.$i18n.locale = 'zh-CN'
      this.$store.dispatch('lang/setLang', 'zh-CN')
      Locale.use('zh-CN', zhCN);
    }
    document.title = String(this.$t('common.pageTitle.1'));
    // 刷新页面
    this.getStoreTree()
    this.getData();
    this.getAuth()
  }
}
</script>

<style lang="scss" scoped>
.container {
  .header {
    height: 176px;
    background: linear-gradient(180deg, rgba(57, 117, 198, 1) 0%, rgba(24, 112, 232, 1) 100%);
    .flex {
      align-items: center;
    }

    .user {
      display: flex;
      align-items: center;
      padding: 16px 24px;
      flex: 1;
      box-sizing: border-box;

      .avatar {
        width: 40px;
        height: 40px;
        border-radius: 50%;
        margin-right: 12px;
        overflow: hidden;
      }

      .name {
        font-weight: 700;
        font-size: 14px;
        line-height: 20px;
        color: #FFFFFF;
      }

      .dept {
        font-weight: 700;
        font-size: 10px;
        line-height: 20px;
        color: #FFFFFF;
        background-color: rgba(33, 37, 41, 1);
        padding: 0 4px;
        border-radius: 2px;
        margin-left: 4px;
      }
    }

    .lang {
      width: 68px;
      height: 40px;
      background: url('~@/assets/store-visit/language-en.png') center/20px no-repeat;

      &.zh {
        background: url('~@/assets/store-visit/language-cn.png') center/20px no-repeat;
      }
    }

  }

  .content {
    margin-top: -112px;
    padding-bottom: calc(env(safe-area-inset-bottom) + 70px);
    padding-bottom: calc(constant(safe-area-inset-bottom) + 70px);

    .section {
      width: 359px;
      // height: 492px;
      margin: auto;
      border-radius: 20px;
      background-color: #fff;
      padding-bottom: 20px;

      .btn {
        display: flex;
        height: 40px;
        justify-content: center;
        align-items: center;
        font-size: 12px;
        font-weight: 600;
        
        .icon-1 {
          width: 16px;
          height: 16px;
          background: url('~@/assets/store-visit/icon_1.png') no-repeat center center/100%;
          margin-left: 8px;
        }
      }

      .title {
        margin-top: 16px;
        display: flex;
        justify-content: space-between;
        font-size: 12px;
        font-weight: 600;
        padding: 20px 24px 0;
        position: relative;
        &::before {
          content: '';
          position: absolute;
          height: 1px;
          background: rgba(248, 249, 250, 1);
          top: 0;
          left: 24px;
          right: 24px;

        }

        .record {
          font-size: 12px;
          font-weight: 400;
          display: flex;
          align-items: center;
          .icon-right {
            width: 16px;
            height: 16px;
            background: url('~@/assets/store-visit/icon_right_black.png') no-repeat center center/100%;
            
          }
        }

      }
      .card {
        width: 327px;
        height: 48px;
        margin: 16px auto 0;
        position: relative;
        z-index: 1;

        .card-inner {
          position: relative;
          z-index: 1;
          border-radius: 12px;
          padding: 16px;
          background: var(--Bootstrap-blue, rgba(13, 110, 253, 1));
          box-shadow: 0px 6px 12px 0px rgba(33, 37, 41, 0.2);
          box-sizing: border-box;
          display: flex;
          align-items: center;
          justify-content: space-between;

          .card-title {
            font-size: 12px;
            font-weight: 600;
            color: #fff;
            margin-right: 8px;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
            flex: 1;
          }
          .flex {
            display: flex;
            align-items: center;
          }

          .card-time {
            font-size: 12px;
            font-weight: 700;
            font-style: italic;
            color: #fff;
          }

          .icon-date {
            margin-left: 8px;
            width: 16px;
            height: 16px;
            background: url('~@/assets/store-visit/icon_date.png') no-repeat center center/100%;
          }
        }

        &::after {
          content: "";
          width: 307px;
          height: 48px;
          background-color: rgba(148, 193, 255, 1);
          position: absolute;
          top: 8px;
          left: 50%;
          transform: translateX(-50%);
          box-shadow: 0px 4px 8px 0px rgba(0, 0, 0, 0.05);
          border-radius: 8px;
          z-index: -1;
        }
      }
      .empty {
        height: 84;
        font-weight: 300;
        font-size: 12px;
        text-align: center;
        padding: 20px;
        box-sizing: border-box;
      }
    }


    .btn-group {
      display: flex;
      padding: 16px;

      .btn { 
        flex: 1;
        height: 44px;
        // margin: 16px auto;
        background-color: rgba(255, 255, 255, 1);
        border-radius: 8px;
        display: flex;
        align-items: center;
        justify-content: center;
        box-sizing: border-box;
        padding: 0;
        color: rgba(0, 0, 0, 1);
        border: 1px solid rgba(0, 0, 0, 1);
        font-size: 12px;
        font-weight: 600;
        margin-right: 16px;
        
        .icon-1 {
          width: 24px;
          height: 24px;
          background: url('~@/assets/store-visit/icon_1.png') no-repeat center center/100%;
          margin-right: 12px;
        }
      }
      .start {
        // width: 343px;
        flex: 1;
        height: 44px;
        // margin: 16px auto;
        background-color: rgba(0, 0, 0, 1);
        border-radius: 8px;
        display: flex;
        align-items: center;
        justify-content: center;
        box-sizing: border-box;
        padding: 0;
  
        .icon-2 {
          width: 24px;
          height: 24px;
          background: url('~@/assets/store-visit/icon_2.png') no-repeat center center/100%;
          margin-right: 12px;
        }
        .icon-right-white {
          width: 24px;
          height: 24px;
          background: url('~@/assets/store-visit/icon_right_white.png') no-repeat center center/100%;
          margin-right: 8px;
        }
  
        .start-text {
          // flex: 1;
          color: rgba(255, 255, 255, 1);
          font-size: 12px;
          font-weight: 600;
        }
  
      }
    }

    .tips {
      padding: 0 16px;
      font-size: 12px;
      font-weight: 400;
      color: rgba(0, 0, 0, 0.6);
    }
  }

  .footer {
    background-color: #fff;
    padding-bottom: calc(env(safe-area-inset-bottom));
    padding-bottom: calc(constant(safe-area-inset-bottom));
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    z-index: 10;
  }

  .tabbar {
    height: 60px;
    display: flex;
    justify-content: center;

    .tabbar-item {
      width: 40%;
      text-align: center;
      font-size: 16px;
      line-height: 60px;
      font-weight: 400;
      height: 60px;

      &.active {
        color: rgba(0, 0, 0, 1);
        font-weight: 600;
      }
      
      &:not(.active) {
        color: rgba(0, 0, 0, 0.2);
      }
    }
  }
}
</style>

<style>
</style>