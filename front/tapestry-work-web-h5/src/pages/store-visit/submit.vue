<template>
  <div v-if="detail">
    <!-- 巡店报告 -->
    <div class="container">
      <div class="header">
        <div class="store-name">{{detail.store_name}}</div>
        <div class="store-report-name">{{detail.creator.dept_name}} Store Visit Report</div>
        <div class="store-info">
          <div class="flex">
            <div class="label">{{ $t('storeVisit.tips.9') }}:</div>
            <div class="value">{{detail.date.replaceAll('-', '.')}} {{detail.time_span}}</div>
          </div>
        </div>
      </div>
      <div class="content">
        <div class="section">
          <!-- 概览 -->
          <div class="section-part">
            <div class="section-head">
              <div class="section-title">{{ $t('storeVisit.tips.17') }}</div>
              <div v-if="isCreator(detail) && detail.overview" class="icon-sync" @click="handleSyncOverview()"></div>
              <div v-if="isCreator(detail) && !detail.overview" class="icon-edit" @click="handleEditOverview"></div>
            </div>
            <div class="section-content" @click="handleEditOverview">
              <div class="edit-content">{{overview()}}<span v-if="!detail.overview" class="tips">{{ $t('storeVisit.tips.43') }}</span></div>
            </div>
          </div>

          <!-- 巡店反馈记录入口 -->
          <div class="section-part">
            <div class="section-head">
              <div class="section-title">{{ $t('storeVisit.tips.20') }}</div>
            </div>
            <div class="card">
              <img class="card-img" src="@/assets/store-visit/icon_y.png" alt="">
              <div class="flex-1">
                <div class="card-row">
                  <span class="label">{{ $t('storeVisit.tips.21') }}:</span>
                  <span class="name">@ {{ detail.creator.name }}</span>
                  <span class="dept">{{detail.creator.dept_name}}</span>
                </div>
                <div class="card-row">
                  <div class="label">{{ $t('storeVisit.tips.22') }}:</div>
                  <span class="value">{{$helper.dateFormat(detail.updated_at, "yyyy.MM.dd hh:mm:ss")}}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="section">
          <div class="section-head">    
            <div class="section-title">{{ $t('storeVisit.tips.11') }}</div>
            <div class="value" @click="handleAddVisitor">
              <div>{{ $t('storeVisit.btn.28') }}</div>
              <div class="icon-right"></div>
            </div>
          </div>
          <div class="section-content"> 
            <div class="user-list">
              <!-- 门店店长 -->
              <div class="user-item" v-for="(item, index) in managerList" :key="index">
                <div class="flex-1">
                  <div class="user-name">@ {{item.name}}</div>
                  <div class="user-dept">{{ item.position_abbr }}</div>
                  <!-- <div class="user-dept">{{ item.dept_name }}</div> -->
                </div>
                <div class="icon-close-circle" @click="handleDelete(index, 'manager')"></div>
              </div>
              <!-- 店长同步人员 -->
              <div class="user-item" v-for="(item, index) in copyMemberList" :key="index">
                <div class="flex-1">
                  <div class="user-name">@ {{item.name||item.email}}</div>
                  <div class="user-dept">{{ item.dept_name }}</div>
                </div>
                <div class="icon-close-circle" @click="handleDelete(index, 'copyMember')"></div>
              </div>
              <!-- 巡店人员 -->
              <div class="user-item" v-for="(item, index) in visitorMemberList" :key="index">
                <div class="flex-1">
                  <div class="user-name">@ {{item.name||item.email}}</div>
                  <div v-if="item.is_mm" class="user-dept">MM</div>
                  <div v-else class="user-dept" :class="{ext: item.dept_name === 'Email'}">{{item.dept_name }}</div>
                </div>
                <div class="icon-close-circle" @click="handleDelete(index, 'visitor')"></div>
              </div>
            </div>
          </div>
        </div>
  
      </div>
      
      <!-- footer -->

      <div class="footer">
        <div class="back" @click="handleBack">
          <div class="icon-back"></div>
          <div>{{ $t('storeVisit.btn.25') }}</div>
        </div>
        <van-button color="#000000" plain @click="handleSave">{{ $t('storeVisit.btn.22') }}</van-button>
        <van-button color="#000000" @click="handleComplete">{{ $t('storeVisit.btn.24') }}</van-button>
      </div>
      <!-- 选择相关巡店人员 -->
      <edit-visitor-popup v-model="showEditVisitorPopup" :entrance="entrance" :detail="detail" @cancel="handleEditVisitorCancel" @confirm="handleEditVisitorOk"></edit-visitor-popup>
      <!-- 编辑富文本内容 -->
      <com-editor-popup v-model="showEditContentPopup" :content="editContent" @confirm="handleEditContentConfirm"></com-editor-popup>
      <!-- 发送评论 -->
      <send-message-popup v-model="showEditOverviewPopup" :defaultValue="editOverview" only-text :confirm-text="$t('storeVisit.btn.26')" :placeholder="$t('storeVisit.btn.27')" @confirm="handleEditOverviewConfirm"></send-message-popup>
      <!-- 编辑巡店人员 -->
      <select-visitor-popup v-model="showSelectVisitorPopup" :default-selection="selection" @select="handleVisitorSelect" @cancel="handleSelectVisitorCancel" @confirm="handleSelectVisitorConfirm"></select-visitor-popup>
  
    </div>
    <!-- 导出报告 -->
    <com-export-popup v-model="showExportPopup" @confirm="handleExportConfirm"></com-export-popup>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
import EditVisitorPopup from './components/edit-visitor-popup.vue'
import SelectVisitorPopup from "./components/select-visitor-popup.vue";
import SendMessagePopup from "./components/send-message-popup.vue";
import ComEditorPopup from "./components/com-editor-popup.vue";
import ComExportPopup from "./components/com-export-popup.vue";
import MessageList from "./message-list.vue";
import StoreVisitDetail from "./detail.vue";
import ComPopup from "./components/ComConfirmPopup";
import { sendPdfToEmailApi, getStoreVisitInfoApi, getVisitRecordListApi, saveVisitInfoApi, saveVisitReportAuthApi, getStoreVisitAuthApi, sendVisitReportApi, getConfigApi, getStoreManagerApi, completeVisitApi, getStoreMMApi } from "@/api/store-visit";
import wx from 'weixin-js-sdk';
import dayjs from "dayjs";
import { i18n } from "@/locales";
import { Locale } from "vant";
import enUS from "@/locales/lang/en-US";
import zhCN from "@/locales/lang/zh-CN";
@Component({
  components: {
    EditVisitorPopup,
    SelectVisitorPopup,
    MessageList,
    ComEditorPopup,
    SendMessagePopup,
    StoreVisitDetail,
    ComExportPopup,
  },
})
export default class StoreVisitReport extends Vue {
  @Prop({ type: Number, default: 0 }) visitId!: number;
  /** 确认弹窗 */
  showConfirmPopup = false;
  /** 选择访客弹窗 */
  showEditVisitorPopup = false;
  /** 选择人员弹窗 */
  showSelectVisitorPopup = false;
  /** 编辑内容 */
  showEditContentPopup = false;
  /** 编辑概览 */
  showEditOverviewPopup = false;
  /** 编辑概览 */
  showSyncPopup = false;
  /** 导出报告 */
  showExportPopup = false;
  /** 巡店详情 */
  detail:StoreVisit.IDetail = null;
  /** 巡店反馈列表 */
  recordList: StoreVisit.IFeedbackRecordItem[] = [];
  /** 巡店人员 */
  visitorMemberList: any[] = []
  /** 门店店长 */
  managerList: any[] = []
  /** 店长同步人员 */
  copyMemberList: any[] = []
  /** 当前门店的MM人员 */
  mmList: any[] = []

  editContent = ''
  editOverview = ''
  entrance = ''
  isEdit = false

  /** 需要推送的成员列表 */
  get sendMemberList() {
    if (!this.detail) {
      return []
    }
    const { visitor_is_show, store_is_show, visitor_member, store_manager, copy_member } = this.detail
    const list = []
    
    if(visitor_is_show) {
      list.push(...visitor_member)
    }

    if(store_is_show) {
      list.push(...store_manager, ...copy_member)
    }

    return list.reduce((pre, cur) => {
      const isExist = pre.some((item) => item.member_id === cur.member_id && item.name === cur.name && item.email === cur.email)
      // 去重
      if (!isExist) {
        pre.push(cur)
      }
      return pre
    }, [])
  }

  /** 概览 */
  get overview() {
    return () => {
      const { overview } = this.detail
      const typeNames = this.recordList.map(item => item.type_name||this.$t('storeVisit.feedbacks.10')).join('、')
      const text = typeNames ? `${this.$t('storeVisit.tips.34')}${typeNames}${this.$t('storeVisit.tips.35')}` : ``
      return overview || text
    }
  }

  get content() {
    if (!this.detail) {
      return ''
    }
    return this.detail.content
  }

  /** 用户信息 */
  get userInfo():IUserInfo {
    return this.$store.getters["user/getUserInfo"]
  }


  /**
   * 门店店长
   */
  get isStoreManager() {
    return (item) => item.store_manager.findIndex(i => i.member_id === this.userInfo.member_id) > -1
  }
  /**
   * 创建人
   */
  get isCreator() {
    return (item) => item.creator.member_id === this.userInfo.member_id
  }

  /**
   * 巡店人员
   */
  get isStoreVisitor() {
    return (item) => item.visitor_member.findIndex(i => i.member_id === this.userInfo.member_id) > -1
  }

  /**
   * 同步人员
   */
  get isCopyMember() {
    return (item) => item.copy_member.findIndex(i => i.member_id === this.userInfo.member_id) > -1
  }

  /** 显示反馈记录 */
  get showFeedback() {
    return (item:StoreVisit.IDetail) => {


      const hasFeedback = this.recordList.length > 0
      if(this.isCreator(item)) {
        // 有反馈记录，如果是创建人，显示
        return true
      } else if(this.isCopyMember(item) || this.isStoreManager(item)) {
        // 有反馈记录，如果是门店相关人员，通过权限控制
        return item.store_is_show && hasFeedback
      } else if(this.isStoreVisitor(item)) {
        // 有反馈记录，如果是巡店人员，通过权限控制
        return item.visitor_is_show && hasFeedback
      } else {
        return false
      }
    }
  }

  /** 
   * 概览
   */
  get showOverview() {
    return (item:StoreVisit.IDetail) => {
      const hasFeedback = this.recordList.length > 0
      if(!hasFeedback) {
        // 如果没有反馈记录 不显示
        console.log(555);
        return false
      } else if(this.isCreator(item)) {
        console.log(4444);
        // 有反馈记录，如果是创建人，显示
        return true
      } else if(this.isCopyMember(item) || this.isStoreManager(item)) {
        console.log(333);
        // 有反馈记录，如果是门店相关人员，通过权限控制
        return item.store_is_show
      } else if(this.isStoreVisitor(item)) {
        console.log(222);
        // 有反馈记录，如果是巡店人员，通过权限控制
        return item.visitor_is_show
      } else {
        console.log(111);
        
        return false
      }
    }
  }
  get show() {
    return (item:StoreVisit.IDetail) => {
      if(this.recordList.length === 0) return false
      return item.store_is_show &&(this.isCopyMember(item) || this.isStoreManager(item)) || item.visitor_is_show && (this.isStoreVisitor(item) || this.isCreator(item)) || this.isCreator(item)
    }
  }



  get trackInfo () {
    return {
      SV_id: this.detail.id,
      store_name: this.detail.store_name,
      store_code: this.detail.store_code,
      visit_time: `${this.detail.date} ${this.detail.time_span}`,
      is_synced: !!this.detail.is_report,
    }
  }


  // 当前语言
  get lang () {
    return this.$store.getters["lang/getLang"]
  }


  get selection() {
    return this.visitorMemberList.concat([...this.managerList, ...this.copyMemberList])
  }


  destroyed() {
    if(!this.$route.query.from) {
      // 记得清除!!!!不然单页面，其它页面也可以用
      window.removeEventListener("popstate", this.popstate, false);
    }
  }
   
  created() {
    // 当前需要直接退出的页面
    if(!this.$route.query.from) {
      this.getConfig()
      window.history.pushState("_android_back", "", location.href);
      window.addEventListener("popstate", this.popstate, false);
    }
  }

  mounted() {
    this.init()
  }

  popstate() {
    wx.closeWindow(); // 微信网页退出
  }

  async init() {
    await this.getAuth()
    await this.getDetail()
    
    this.getVisitRecordList()

    // 神策数据上报 -- 巡店计划浏览
    this.$track.trackStoreVisitAppointmentView({
      ...this.trackInfo,
      SV_initiator: this.detail.creator.name,
      SV_initiator_id: this.detail.creator.user_id,
      SV_initiator_department: this.detail.creator.dept_name,
    })


  }

  async getConfig() {
    const res = await getConfigApi({
      key: 'language'
    });
    console.log(JSON.parse(res.data.value));
    const { en, zh } = JSON.parse(res.data.value)
    
    this.changeLang(res.data.member_config === 'en' ? 'zh-CN' : 'en-US')
    // i18n.setLocaleMessage('zh-CN', zh)
    // i18n.setLocaleMessage('en-US', en)
  }

  changeLang(lang) {
    if(lang === 'zh-CN') {
      this.$i18n.locale = 'en-US'
      this.$store.dispatch('lang/setLang', 'en-US')
      Locale.use('en-US', enUS);

    } else {
      this.$i18n.locale = 'zh-CN'
      this.$store.dispatch('lang/setLang', 'zh-CN')
      Locale.use('zh-CN', zhCN);
    }
    document.title = String(this.$t('common.pageTitle.7'));
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
      })
    } catch (error) {
      console.log(error);
    }
  }
  /**
   * 获取详情
   */
  async getDetail() {
    const res:IApiResult<StoreVisit.IDetail> = await getStoreVisitInfoApi({
      id: this.visitId
    })
    if(res.code === 200) {
      res.data.store_manager.forEach((item:any) => {
        item.third_id = res.data.third_id || 0
      })
      this.detail = res.data
      this.visitorMemberList = JSON.parse(JSON.stringify(res.data.visitor_member))
      this.managerList = JSON.parse(JSON.stringify(res.data.store_manager))
      this.copyMemberList = JSON.parse(JSON.stringify(res.data.copy_member))
      this.getStoreMM()
    }
  }
  /**
   * 获取详情
   */
  async getStoreMM() {
    const res:any = await getStoreMMApi({
      this_ids: [this.detail.third_id]
    })
    if(res.code === 200) {
      this.mmList = res.data
    }
  }


  /**
   * 获取巡店反馈列表
   */
  async getVisitRecordList() {
    const res:IApiResult<StoreVisit.IFeedbackRecordItem[]> = await getVisitRecordListApi({
      id: this.visitId,
    })

    if(res.code === 200) {
      this.recordList = res.data
    }
  }

  generateContent(list, message = this.$t('storeVisit.popup.tips.18')) {
    let content = `<div class="u-tips">${message}:</div>
      <div class="u-list">`;

    list.forEach((item, index) => {
      const deptClass = item.dept_name === 'Email' ? 'ext' : '';
      content += `
        <div class="item" data-index="${index}">
          <div class="item-name">@ ${item.name || item.email}</div>
          <div class="item-dept" class="${deptClass}">${item.is_mm ? 'MM' : item.dept_name || 'Email'}</div>
        </div>`;
    });

    content += `</div>`;
    return content;
  }

  /**
   * 推送消息
   */
  async handleConfirm(onlyNew = false) {
    const list = this.visitorMemberList.filter(item => item)
    const res:IApiResult = await sendVisitReportApi({
      id: this.detail.id,
    })
    if(res.code === 200) {
      // this.$toast(this.$t('common.toast.13'))
      this.detail.is_report = 1

      // 神策数据上报 -- 推送报告
      this.$track.trackSyncReport({
        ...this.trackInfo,
        is_new_added: onlyNew,
        is_sa_visible: !!this.detail.store_is_show,
        is_relevant_visible: !!this.detail.visitor_is_show,
        sync_staff_list: list.map(item => `${item.name || item.email}:${item.dept_name||'Email'}`)
      })
    }
  }

  /**
   * 编辑名单
   */
  handleEditVisitor(isEdit = false) {
    this.showEditVisitorPopup = true;
    this.entrance = '新增'
    this.isEdit = isEdit
  }

  /**
   * 编辑名单 取消
   */
  handleEditVisitorCancel() {
    this.showEditVisitorPopup = false;
  }

  /**
   * 编辑名单 确认
   */
  async handleEditVisitorOk(value) {
    // console.log(value, this.detail.visitor_member);
    // return
    
    // if(this.detail.is_report === 1 && value.add_visitor_list.length) {
    //   const { action } = await ComPopup.confirm({
    //     title: this.$t('storeVisit.confirm.title.1'),
    //     content: this.generateContent(value.add_visitor_list),
    //     confirmButtonText: this.$t('storeVisit.confirm.confirmBtn.7'),
    //     cancelButtonText: this.$t('storeVisit.confirm.cancelBtn.2'),
    //     footerStyle: '2'
    //   })
    //   if(action === 'cancel') {
    //     return
    //   }
    // }
    
    // this.$toast.loading()
    try {
      const res:IApiResult<StoreVisit.IFeedbackRecordItem[]> = await saveVisitReportAuthApi({
        id: this.detail.id,
        ...value
      })
    } catch (error) {
      // this.$toast.clear()
    }
  }

  /**
   * 发送评论
   */
  handleSendMessage() {
    (this.$refs.message as any).handleSendMessage();
  }

  /**
   * 编辑之前提示
   */
  async beforeEdit(type) {
    const res = await ComPopup.confirm({
      title: this.$t('storeVisit.confirm.title.1'),
      content: this.$t('storeVisit.confirm.content.6'),
      confirmButtonText: this.$t('storeVisit.confirm.confirmBtn.5'),
      cancelButtonText: this.$t('storeVisit.confirm.cancelBtn.2'),
      footerStyle: '2'
    })
    if(res.action === 'confirm') {

      // 神策数据上报 -- 编辑报告概览
      this.$track.trackEditReportSummary({
        ...this.trackInfo,
        summary_type: type,
        edit_type: '编辑',
      })
      return true
    }  else {
      throw  new Error('取消')
    }
  }

  async handleEditContent() {
    if(!this.isCreator(this.detail)) return
    if(!this.detail.content) {
      await this.beforeEdit('整体表现')
    }
    this.showEditContentPopup = true
    this.editContent = this.detail.content || (this.$refs.content as any)?.innerHTML
  }

  async handleEditContentConfirm(e) {
    this.showEditContentPopup = false
    const res:IApiResult = await saveVisitInfoApi({
      id: this.detail.id,
      content: e || ' '
    })
    if(res.code === 200) {
      this.$toast(this.$t('storeVisit.popup.tips.20'))
      this.getDetail()
    }
  }

  /**
   * 重置内容
   */
  async handleResetContent() {
    
    const res:IApiResult = await saveVisitInfoApi({
      id: this.detail.id,
      reset_content: 1
    })
    if(res.code === 200) {
      this.$toast(this.$t('storeVisit.popup.tips.20'))
      this.getDetail()
    }
  }

  /**
   * 编辑概览
   */
  async handleEditOverview() {
    // if(!this.isCreator(this.detail)) return
    // if(!this.detail.overview) {
    //   await this.beforeEdit('概览')
    // }
    this.showEditOverviewPopup = true
    this.editOverview = this.overview()
  }

  /**
   * 确认概览
   */
  async handleEditOverviewConfirm(e) {
    this.showEditOverviewPopup = false
    const res:IApiResult = await saveVisitInfoApi({
      id: this.detail.id,
      overview: e.content || ' '
    })
    if(res.code === 200) {
      this.$toast(this.$t('storeVisit.popup.tips.20'))
      this.getDetail()
    }

  }


  /**
   * 同步概览
   */
   async handleResetOverview() {
    
    const res:IApiResult = await saveVisitInfoApi({
      id: this.detail.id,
      reset_overview: 1
    })
    if(res.code === 200) {
      this.$toast(this.$t('storeVisit.popup.tips.20'))
      this.getDetail()
    }
  }

  /**
   * 编辑巡店记录
   */
  handleEditRecord() {
    this.$router.push({
      name: 'store-visit-edit',
      query: {
        id: String(this.detail.id)
      }
    })
  }

  /**
   * 同步整体内容
   * @param key 
   */
  async handleSyncContent() {
    const res = await ComPopup.confirm({
      title: this.$t('storeVisit.confirm.title.1'),
      content: this.$t('storeVisit.confirm.content.7'),
      confirmButtonText: this.$t('storeVisit.confirm.confirmBtn.6'),
      cancelButtonText: this.$t('storeVisit.confirm.cancelBtn.2'),
      footerStyle: '2'
    })
    if(res.action === 'confirm') {
      this.handleResetContent()


      // 神策数据上报 -- 编辑报告概览
      this.$track.trackEditReportSummary({
        ...this.trackInfo,
        summary_type: '整体表现',
        edit_type: '同步',
      })
    }
  }

  /**
   * 同步概览内容
   * @param key 
   */
  async handleSyncOverview() {
    const res = await ComPopup.confirm({
      title: this.$t('storeVisit.confirm.title.3'),
      content: this.$t('storeVisit.confirm.content.7'),
      confirmButtonText: this.$t('storeVisit.confirm.confirmBtn.6'),
      cancelButtonText: this.$t('storeVisit.confirm.cancelBtn.2'),
      footerStyle: '2'
    })
    if(res.action === 'confirm') {
      this.handleResetOverview()
      

      // 神策数据上报 -- 编辑报告概览
      this.$track.trackEditReportSummary({
        ...this.trackInfo,
        summary_type: '概览',
        edit_type: '同步',
      })
    }
  }
  async handleSync(key) {
    this.showSyncPopup = true
    this.$set(this.detail, 'sync', key)
  }

  /**
   * 导出
   */
  handleExport() {
    this.showExportPopup = true;
  }

  /**
   * 导出确认
   * @param value 
   */
  async handleExportConfirm(value) {
    this.showExportPopup = false;
    const res:IApiResult = await sendPdfToEmailApi({
      id: this.detail.id,
      email: value
    })
    if(res.code === 200) {
      this.$toast(this.$t('common.toast.14'))
      // 神策数据上报 -- 下载报告
      this.$track.trackDownloadReport({
        ...this.trackInfo,
        SV_initiator: this.detail.creator.name,
        SV_initiator_id: this.detail.creator.user_id,
        SV_initiator_department: this.detail.creator.dept_name,
        report_time: dayjs(this.detail.updated_at*1000).format('YYYY-MM-DD HH:mm:ss'),
      })
    } else {
      this.$toast(res.message)
    }
  }

  /**
   * 添加巡店人员
   */
  handleAddVisitor() {
    this.showSelectVisitorPopup = true;
  }

  handleSelectVisitorCancel() {
    this.showSelectVisitorPopup = false;
  }

  handleSelectVisitorConfirm(data) {
    const copyIds = this.copyMemberList.map(item => item.member_id)
    const managers = data.filter(item => item.position_abbr && this.detail.third_id ===item.third_id)
    const visitors = data.filter(item => !(item.position_abbr && this.detail.third_id ===item.third_id))
    const copyMember = visitors.filter(item => copyIds.includes(item.member_id))
    
    this.showSelectVisitorPopup = false;
    this.$emit('update:list', data);
    this.visitorMemberList = visitors.filter(item => !copyIds.includes(item.member_id));
    this.managerList = managers;
    this.copyMemberList = copyMember;
  }

  async handleDelete(index, type) {
    const { action } = await ComPopup.confirm({
      title: this.$t('storeVisit.confirm.title.4'),
      content: this.$t('storeVisit.confirm.content.9'),
      confirmButtonText: this.$t('storeVisit.confirm.confirmBtn.9'),
      cancelButtonText: this.$t('storeVisit.confirm.cancelBtn.2'),
      footerStyle: '1'
    })
    if(action === 'cancel') {
      return
    }
    if(type === 'manager') {
      this.managerList.splice(index, 1);
    } else if(type === 'copyMember') {
      this.copyMemberList.splice(index, 1);
    } else {
      this.visitorMemberList.splice(index, 1);
     }
  }

  handleVisitorSelect(e) {
    const { track } = e
    // 神策数据上报 -- 选择门店
    this.$track.trackChooseStaff({
      SV_id: this.detail.id,
      store_name: this.detail.store_name,
      store_code: this.detail.store_code,
      SV_initiator: this.detail.creator.name,
      SV_initiator_id: this.detail.creator.user_id,
      SV_initiator_department: this.detail.creator.dept_name,
      visit_time: this.detail.date + ' ' + this.detail.time_span,
      choose_scene: '报告查看权限',
      choose_entrance: this.detail.is_report ? '新增':'添加',
      ...track,
    })
  }

  async handleSave() {
    await this.handleSaveVisitor()
    this.$toast(this.$t('common.toast.24'))
  }


  async handleComplete() {
    await this.completeVisit()
  }

  async handleSaveVisitor() {
    // 删除的巡店人员
    const deleteVisitor = this.detail.visitor_member.filter(item => !this.visitorMemberList.some(sel => 
      sel.member_id === item.member_id
      && sel.email === item.email
    )).map(item=> {
      return {
        member_id: item.member_id,
        email: item.email,
        dept_name: item.dept_name,
        name: item.name
      }
    })
    // 新增的巡店人员
    const addVisitor = this.visitorMemberList.filter(item => !this.detail.visitor_member.some(sel => 
      sel.member_id === item.member_id
      && sel.email === item.email
    )).map(item=> {
      return {
        member_id: item.member_id,
        email: item.email,
        dept_name: item.dept_name,
        name: item.name
      }
    })

    // 删除的门店店长
    const deleteManager = this.detail.store_manager.filter(item => !this.managerList.some(sel => 
      sel.member_id === item.member_id
    )).map(item=> {
      return {
        member_id: item.member_id,
        email: item.email,
        dept_name: item.dept_name,
        name: item.name
      }
    })
    // 新增的门店店长
    const addManager = this.managerList.filter(item => !this.detail.store_manager.some(sel => 
      sel.member_id === item.member_id
    )).map(item=> {
      return item.member_id
    })
    // 删除的店长关联人员
    


    const value = {
      store_is_show: this.detail.store_is_show,
      visitor_is_show: this.detail.visitor_is_show,
      del_visitor_list: [...deleteVisitor, ...deleteManager],
      add_visitor_list: addVisitor,
      add_store_manager: addManager,
    }
    await this.handleEditVisitorOk(value)
  }
    /**
   * 完成巡店
   */
  async completeVisit() {
    
    const { action } = await ComPopup.confirm({
      title: this.$t('storeVisit.confirm.title.2'),
      content: this.$t('storeVisit.confirm.content.8'),
      confirmButtonText: this.$t('storeVisit.confirm.confirmBtn.2'),
      cancelButtonText: this.$t('storeVisit.confirm.cancelBtn.2'),
      footerStyle: this.lang ==='zh-CN' ? '' : '2',
    })
    if(action === 'cancel') {
      return
    }
    this.$toast.loading()

    const res:IApiResult = await completeVisitApi({
      id: this.detail.id
    })

    if(res.code === 200) {
      this.track()

      this.$store.commit("storeVisit/setUpdateId", this.detail.id)
      this.$router.replace({
        name: 'store-visit-report',
        query: {
          id: String(this.detail.id),
          from: '2',
        },
      })
    }

    await this.handleSaveVisitor()
    await this.handleConfirm()

    this.$toast.clear();
    this.$toast(this.$t('common.toast.23'))
  }
  track() {
    const list = this.recordList.map(item => {
      return {
        feedback_id: String(item.id),
        feedback_category: item.type_name,
        is_customized: item.is_custom,
        is_picture: item.media_list.length > 0,
        is_text: !!(item.content && item.content.length > 0),
      }
    })
    // 神策数据上报 -- 提交巡店记录
    this.$track.trackSubmitReport({
      ...this.trackInfo,
      feedback_quantity: list.length ,
      feedback_category_list: list.map(item => item.feedback_category) ,
      is_synced: false,
      button_name: '提交巡店记录',
    })

    list.forEach(item => {
      // 神策数据上报 -- 提交巡店记录详情
      this.$track.trackSubmitReportDetail({
        ...this.trackInfo,
        ...item,
        is_synced: false,
        button_name: '提交巡店记录',
      })
    })
  }
  handleBack() {
    this.$emit('back')
  }
}
</script>

<style lang="scss" scoped>
.icon-right {
  width: 16px;
  height: 16px;
  background: url(~@/assets/store-visit/icon_right.png) center/16px no-repeat;
}
.icon-delete {
  width: 16px;
  height: 16px;
  background: url(~@/assets/store-visit/icon_popup_close.png) center/100% no-repeat;
}

.van-switch {
  width: 26px;
  height: 16px;
  border: none;
  background-color: rgba(173, 181, 189, 1);
  overflow: hidden;
  // display: block;
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

.row {
  display: flex;
  align-items: center;
  height: 28px;

  &:not(:last-child) {
    margin-bottom: 12px;
  }

  &-label {
    flex: 1;
    font-size: 12px;
    color: rgba(0, 0, 0, 0.8);
  }

  &.add .row-label {
    color: #0D6EFD;
  }

  &-value {
    display: flex;
    align-items: center;

    .month,
    .time {
      width: 100px;
      height: 32px;
      border-radius: 4px;
      background-color: rgba(248, 249, 250, 1);
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 12px;
    }

    .month {
      margin-right: 4px;
    }
  }

}
.icon-close-circle {
  width: 16px;
  height: 16px;
  background: url(~@/assets/store-visit/icon_close_circle.png) center/16px no-repeat;
}
.user-list:not(:last-child) { 
  margin-bottom: 12px;
}
.user-item {
  height: 32px;
  background-color: rgba(248, 249, 250, 1);
  border-radius: 4px;
  display: flex;
  align-items: center;
  padding-left: 8px;
  .flex-1 {
    flex: 1;
    display: flex;
    align-items: center;
  }

  .user-name {
    font-size: 12px;
    color: rgba(0, 0, 0, 1);
  }

  .user-dept {
    font-weight: 700;
    font-size: 10px;
    color: rgba(0, 0, 0, 1);
    background-color: rgba(206, 212, 218, 1);
    line-height: 16px;
    padding: 0 4px;
    border-radius: 2px;
    margin-left: 4px;

    &.ext {
      color: rgba(153, 153, 153, 1);
      background-color: rgba(255, 255, 255, 1);
    }
  }

  .icon-close-circle {
    width: 32px;
    height: 32px;
  }

  &:not(:last-child) {
    margin-bottom: 4px;
  }
}

.icon-back {
  width: 20px;
  height: 20px;
  background: url("~@/assets/store-visit/icon_back.png") center/100% no-repeat;
}

.icon-edit {
  width: 24px;
  height: 24px;
  background: url("~@/assets/store-visit/icon_edit.png") center/100% no-repeat;
}
.icon-right-black {
  width: 16px;
  height: 16px;
  background: url("~@/assets/store-visit/icon_right_black.png") center/100% no-repeat;
}
.icon-sync {
  width: 24px;
  height: 24px;
  background: url("~@/assets/store-visit/icon_sync.png") center/100% no-repeat;
}
.container {
  width: 100%;
  min-height: 100vh;
  background: rgba(248, 249, 250, 1) linear-gradient(180deg, #3975C6 0%, #0D6EFD 100%) center top/100%  286px no-repeat;
  padding-bottom: calc(env(safe-area-inset-bottom) + 56px);
  padding-bottom: calc(constant(safe-area-inset-bottom) + 56px);

  .header {
    padding: 16px 24px;
    .store-name {
      font-size: 12px;
      font-weight: 700;
      color: rgba(248, 249, 250, 1);
      line-height: 16px;
    }
    .store-report-name {
      font-size: 20px;
      font-weight: 600;
      color: rgba(255, 255, 255, 1);
      line-height: 32px;
    }

    .store-info  { 
      margin-top: 12px;
      .flex {
        display: flex;
        line-height: 16px;
        .label {
          color: rgba(222, 226, 230, 1);
          font-size: 12px;
        }
        .value {
          color: rgba(255, 255, 255, 1);
          font-size: 12px;
          margin-left: 8px;
          
        }

        &:not(:last-child) {
          margin-bottom: 4px;
        }
      }
    }
  }

  .remind {
    width: 359px;
    margin: 8px auto 0;
    display: flex;
    align-items: center;
    justify-content: space-between;
    box-sizing: border-box;
    border-radius: 12px;
    padding: 16px;
    box-shadow: 4px 6px 11.6px 0px rgba(1, 27, 158, 0.2);
    background: linear-gradient(0deg, var(--Bootstrap-blue, #0D6EFD), var(--Bootstrap-blue, #0D6EFD)), linear-gradient(91.38deg, rgba(255, 255, 255, 0.1) 0%, rgba(255, 255, 255, 0.2) 100%);


    .remind-text {
      font-size: 12px;
      color: rgba(255, 255, 255, 1);
      font-weight: 600;
    }
    .remind-btn {
      font-size: 10px;
      color: rgba(255, 255, 255, 1);
      font-weight: 700;
      line-height: 12px;
      padding: 8px 12px;
      border-radius: 8px;
      background-color: #000;
      white-space: nowrap;
    }
  }

  .edit-content {
    padding: 12px;
    border-radius: 12px; 
    background-color: rgb(248, 249, 250);
    word-break: break-word;
    white-space: wrap;
    .tips {
      color: #ccc;
    }
  }

  .content {
    padding: 8px;

    .section {
      border-radius: 12px;
      background-color: #fff;
      padding: 24px 16px;
      margin-bottom: 8px;

      .section-part:not(:last-child) {
        margin-bottom: 24px;
      }

      .section-head {
        display: flex;
        align-items: center;
        justify-content: space-between;

        .section-title {
          font-size: 14px;
          font-weight: 600;
          color: rgba(51, 51, 51, 1);
          line-height: 20px;
        }
        .value {
          font-size: 14px;
          color: rgba(13, 110, 253, 1);
          display: flex;
          align-items: center;
        }
      }

      .section-content {
        margin-top: 8px;
        line-height: 20px;
        padding: 8px 4px;
        font-size: 12px;

        .label {
          color: #999999;
          font-weight: 600;
          opacity: 0.5;
        }
        .value {
          margin-top: 12px;
          
        }

        ::v-deep ul {
          list-style-type: disc;
          padding-left: 20px;
        }
        ::v-deep ol {
          list-style-type: decimal;
          padding-left: 20px;
        }
      }

      .section-footer {
        padding: 0;
        display: flex;
        align-items: center;
        justify-content: space-between;

        .text {
          font-size: 12px;
          font-weight: 600;
          color: rgba(153, 153, 153, 1);
          white-space: nowrap;
        }
      }

      .section-part ~ .section-footer {
        border-top: 1px solid rgba(248, 249, 250, 1);
        padding-top: 24px;
      }

      .card {
        border-radius: 8px;
        padding: 12px;
        margin-top: 16px;
        background: 
        linear-gradient(0deg, rgba(255, 255, 255, 0.95), rgba(255, 255, 255, 0.95)) center center/100% 100% no-repeat, 
        linear-gradient(0deg, var(--Bootstrap-blue, #0D6EFD), var(--Bootstrap-blue, #0D6EFD)) center center/100% 100% no-repeat;
        display: flex;
        align-items: center;

        .card-img {
          width: 48px;
          height: 48px;
          margin-right: 16px;
        }

        .card-row {
          display: flex;
          color: rgba(153, 153, 153, 1);
          font-size: 12px;
          line-height: 16px;

          &:nth-child(2) {
            margin-top: 4px;
          }

          .name {
            color: rgba(0, 0, 0, 1);
            margin-left: 8px;
          }

          .dept {
            margin-left: 4px;
            border-radius: 2px;
            background-color: rgba(206, 212, 218, 1);
            line-height: 16px;
            padding: 0 4px;
            color: rgba(0, 0, 0, 1);
            font-size: 10px;
            font-weight: 700;
          }

          .value {
            margin-left: 8px;
          }
        }
      }
    }
  }


  .footer {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    z-index: 10;
    padding: 12px 24px;
    display: flex;
    padding-bottom: calc(env(safe-area-inset-bottom) + 12px);
    padding-bottom: calc(constant(safe-area-inset-bottom) + 12px);
    background-color: rgba(248, 249, 250, 1);

    .back {
      width: auto;
      height: 44px;
      font-size: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
      flex-direction: column;
      box-sizing: border-box;
    }

    .van-button {
      height: 44px;
      border-radius: 8px;
      font-size: 14px;
      font-weight: 600;
      flex: 1;

      &:not(:first-child) {
        margin-left: 8px;
      }

      &.van-button--disabled {
        opacity: 1;
        background: rgba(222, 226, 230, 1) !important;
        border-color: rgba(222, 226, 230, 1) !important;
      }
    }
  }
  
}
</style>

<style lang="scss">
.popup-content { 
  .u-list {
    margin-top: 16px;
    padding: 12px 0;
  
    .item {
      height: 16px;
      display: flex;
      &:not(:last-child) {
        margin-bottom: 8px;
      }

      .item-dept {
        margin-left: 4px;
        height: 16px;
        padding: 0 4px;
        line-height: 16px;
        font-size: 10px;
        background-color: rgba(206, 212, 218, 1);
        color: rgba(0, 0, 0, 1);
        font-weight: 700;
        border-radius: 2px;
  
        &.ext {
          color: rgba(153, 153, 153, 1);
          background-color: rgba(248, 249, 250, 1);
        }
      }
    }
  }
}
</style>