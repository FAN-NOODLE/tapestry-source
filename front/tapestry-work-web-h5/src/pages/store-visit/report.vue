<template>
  <div v-if="detail">
    <!-- 巡店报告 -->
    <div v-if="detail.status == 3" class="container">
      <div class="header">
        <div class="store-name">{{detail.store_name}}</div>
        <div class="store-report-name">{{detail.creator.dept_name}} Store Visit Report</div>
        <div class="store-info">
          <div class="flex">
            <div class="label">{{ $t('storeVisit.tips.9') }}:</div>
            <div class="value">{{detail.date.replaceAll('-', '.')}} {{detail.time_span}}</div>
          </div>
          <div class="flex">
            <div class="label">{{ $t('storeVisit.tips.14') }}:</div>
            <div class="value">{{$helper.dateFormat(detail.updated_at, "yyyy.MM.dd hh:mm:ss")}}</div>
          </div>
        </div>
      </div>
      <!-- <div v-if="isCreator(detail) && !detail.is_report" class="remind" @click="handleConfirm()">
        <div class="remind-text">{{ $t('storeVisit.tips.15') }}</div>
        <div class="remind-btn">{{ $t('storeVisit.tips.16') }}</div>
      </div> -->
      <div class="content">
        <div v-if="showOverview(detail)||isCreator(detail)" class="section">
          <!-- 概览 -->
          <div v-if="showOverview(detail)" class="section-part">
            <div class="section-head">
              <div class="section-title">{{ $t('storeVisit.tips.17') }}</div>
              <div v-if="isCreator(detail) && detail.overview" class="icon-sync" @click="handleSyncOverview()"></div>
              <div v-if="isCreator(detail) && !detail.overview" class="icon-edit" @click="handleEditOverview"></div>
            </div>
            <div class="section-content" @click="handleEditOverview">
              <div class="edit-content">{{overview()}} <span v-if="!detail.overview" class="tips">{{ $t('storeVisit.tips.43') }}</span></div>
            </div>
          </div>
          <!-- 巡店反馈记录入口 -->
          <div v-if="showFeedback(detail)" class="section-part" @click="handleEditRecord">
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
          <!-- 创建人可见 巡店相关人员可见权限 -->
          <div v-if="isCreator(detail)" class="section-footer" @click="handleEditVisitor(true)">
            <div class="text">{{ $t('storeVisit.tips.11') }}</div>
            <div class="icon-right-black"></div>
          </div>
        </div>
  
      </div>
      <!-- 留言板 -->
      <message-list ref="message" :detail="detail"></message-list>
      <!-- footer -->
      <div class="footer">
        <div v-if="showOverview(detail)" class="download" @click="handleExport">
          <div class="icon-download"></div>
          <div class="text">{{ $t('storeVisit.tips.23') }}</div>
        </div>
        <van-button color="rgba(248, 249, 250, 1)" block @click="handleSendMessage">{{ $t('storeVisit.tips.24') }}</van-button>
      </div>
      <!-- 消息推送  弹窗 -->
      <!-- <com-confirm-popup v-model="showConfirmPopup" cancelText="更改名单" confirmText="确认并推送巡店完成通知" @cancel="handleConfirmCancel" @confirm="handleConfirmOk">
        <div class="tips">请确认以下通知人名单</div>
        <div class="list">
          <div class="item" v-for="item, index in sendMemberList" :key="index">
            <div class="item-name">@ {{item.name||item.email}} </div>
            <div class="item-dept" :class="{ext: item.dept_name === 'Email'}">{{item.dept_name}}</div>
          </div>
        </div>
      </com-confirm-popup> -->
      <!-- 选择相关巡店人员 -->
      <edit-visitor-popup v-model="showEditVisitorPopup" :entrance="entrance" :detail="detail" @cancel="handleEditVisitorCancel" @confirm="handleEditVisitorOk"></edit-visitor-popup>
      <!-- 编辑富文本内容 -->
      <com-editor-popup v-model="showEditContentPopup" :content="editContent" @confirm="handleEditContentConfirm"></com-editor-popup>
      <!-- 发送评论 -->
      <send-message-popup v-model="showEditOverviewPopup" :defaultValue="editOverview" only-text :confirm-text="$t('storeVisit.btn.26')" :placeholder="$t('storeVisit.btn.27')" @confirm="handleEditOverviewConfirm"></send-message-popup>
    </div>
    <!-- 巡店详情 -->
    <store-visit-detail v-else :detail="detail" @update:detail="getDetail"></store-visit-detail>
    <!-- 导出报告 -->
    <com-export-popup v-model="showExportPopup" @confirm="handleExportConfirm"></com-export-popup>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import EditVisitorPopup from './components/edit-visitor-popup.vue'
import SelectVisitorPopup from "./components/select-visitor-popup.vue";
import SendMessagePopup from "./components/send-message-popup.vue";
import ComEditorPopup from "./components/com-editor-popup.vue";
import ComExportPopup from "./components/com-export-popup.vue";
import MessageList from "./message-list.vue";
import StoreVisitDetail from "./detail.vue";
import ComPopup from "./components/ComConfirmPopup";
import { sendPdfToEmailApi, getStoreVisitInfoApi, getVisitRecordListApi, saveVisitInfoApi, saveVisitReportAuthApi, getStoreVisitAuthApi, sendVisitReportApi, getConfigApi } from "@/api/store-visit";
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
    ComExportPopup
  },
})
export default class StoreVisitReport extends Vue {
  /** 确认弹窗 */
  showConfirmPopup = false;
  /** 选择访客弹窗 */
  showEditVisitorPopup = false;
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
      const text = typeNames ? `${this.$t('storeVisit.tips.34')}${typeNames}${this.$t('storeVisit.tips.35')}` : `${this.$t('storeVisit.tips.42')}`
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
  get isStoreVistor() {
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
      } else if(this.isStoreVistor(item)) {
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
      } else if(this.isStoreVistor(item)) {
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
      return item.store_is_show &&(this.isCopyMember(item) || this.isStoreManager(item)) || item.visitor_is_show && (this.isStoreVistor(item) || this.isCreator(item)) || this.isCreator(item)
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
      id: Number(this.$route.query.id)
    })
    if(res.code === 200) {

      res.data.store_manager.forEach((item:any) => {
        item.third_id = res.data.third_id || 0
      })
      this.detail = res.data
    }
  }


  /**
   * 获取巡店反馈列表
   */
  async getVisitRecordList() {
    const res:IApiResult<StoreVisit.IFeedbackRecordItem[]> = await getVisitRecordListApi({
      id: Number(this.$route.query.id),
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
    const list = this.sendMemberList.filter(item => {
      if(onlyNew) {
        return item.is_new
      }
      return true
    })

    const {action} = await ComPopup.confirm({
      title: this.$t('storeVisit.confirm.title.1'),
      content: this.generateContent(list, onlyNew ? this.$t('storeVisit.popup.tips.18') : this.$t('storeVisit.popup.tips.14')),
      confirmButtonText: this.$t('storeVisit.confirm.confirmBtn.8'),
      cancelButtonText: this.$t('storeVisit.confirm.cancelBtn.4'),
      footerStyle: '2'
    })
    if(action === 'confirm') {
      const res:IApiResult = await sendVisitReportApi({
        id: this.detail.id,
      })
      if(res.code === 200) {
        this.$toast(this.$t('common.toast.13'))
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
    } else {
      this.handleEditVisitor()
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
    const all = [...value.add_visitor_list, ...value.add_manager_list]
    if(this.detail.is_report === 1 && all.length) {
      const { action } = await ComPopup.confirm({
        title: this.$t('storeVisit.confirm.title.1'),
        content: this.generateContent(all),
        confirmButtonText: this.$t('storeVisit.confirm.confirmBtn.7'),
        cancelButtonText: this.$t('storeVisit.confirm.cancelBtn.2'),
        footerStyle: '2',
        tips: this.$t('storeVisit.tips.45'),
      })
      if(action === 'cancel') {
        return
      }
    }
    
    this.$toast.loading()
    try {
      const res:IApiResult<StoreVisit.IFeedbackRecordItem[]> = await saveVisitReportAuthApi({
        id: this.detail.id,
        ...value,
        add_store_manager: value.add_manager_list.map(item => item.member_id),
      })
      this.$toast.clear()
      if(res.code === 200) {
        await this.getDetail()
        this.showEditVisitorPopup = false;
        if(this.detail.is_report === 1) {
          this.$toast(this.$t('common.toast.13'))
          // 神策数据上报 -- 推送报告
          this.$track.trackSyncReport({
            ...this.trackInfo,
            is_new_added: true,
            is_sa_visible: !!value.store_is_show,
            is_relevant_visible: !!value.visitor_is_show,
            sync_staff_list: all.map(item => `${item.name || item.email}:${item.dept_name||'Email'}`)
          })
        } else {
          this.handleConfirm()
        }
      }
    } catch (error) {
      this.$toast.clear()
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
    if(!this.isCreator(this.detail)) return
    if(!this.detail.overview) {
      await this.beforeEdit('概览')
    }
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
      title: this.$t('storeVisit.confirm.title.1'),
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
}
</script>

<style lang="scss" scoped>
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
          color: rgba(0, 0, 0, 1);
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
    padding: 8px 24px;
    display: flex;
    align-items: center;
    padding-bottom: calc(env(safe-area-inset-bottom) + 8px);
    padding-bottom: calc(constant(safe-area-inset-bottom) + 8px);
    background-color: #fff;

    .download {
      width: 40px;
      height: 40px;
      display: flex;
      align-items: center;
      justify-content: center;
      flex-direction: column;
      margin-right: 24px;

      .icon-download  {
        width: 24px;
        height: 24px;
        background: url("~@/assets/store-visit/icon_download.png") center/100% no-repeat;
        margin-bottom: 8px;
      }

      .text { 
        font-size: 10px;
        font-weight: 600;
        line-height: 12px;
        white-space: nowrap;
      }
    }

    .van-button {
      flex: 1;
      height: 40px;
      font-size: 12px;
      border-radius: 40px;
      color: rgba(153, 153, 153, 1) !important;
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