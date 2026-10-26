<template>
  <div class="detail">
    <div class="header">
      <div class="store-name">{{detail.store_name}}</div>
      <div class="store-report-name">{{detail.creator.dept_name}} Store Visit</div>
      <div class="store-info">
        <div class="flex">
          <div class="label">{{ $t('storeVisit.item.label.2') }}:</div>
          <div class="value">
            <div class="name">{{ detail.creator.name }}</div>
            <div class="dept">{{ detail.creator.dept_name }}</div>
          </div>
        </div>
        <div class="flex">
          <div class="label">{{ $t('storeVisit.item.label.1') }}:</div>
          <div class="value">{{detail.date.replaceAll('-', '.')}} {{detail.time_span}}</div>
        </div>
      </div>
    </div>
    <div v-if="detail.remark" class="content">
      <div class="remark">
        <div class="remark-title">{{ $t('storeVisit.tips.6') }}</div>
        <div class="remark-content">{{ detail.remark }}</div>
      </div>
    </div>
    <div class="footer">
      <van-button color="#000000" plain @click="handleContact(1)">{{ $t('storeVisit.btn.16') }}</van-button>
      <van-button v-if="isStoreManager && copyMember.length>0" color="#000000" plain @click="handleContact(2)">{{ $t('storeVisit.btn.17') }}</van-button>
      <van-button v-if="isStoreManager && copyMember.length==0" color="#000000" plain @click="handleSelectVisitor">{{ $t('storeVisit.btn.18') }}</van-button>
    </div>
    <contact-popup v-model="showContactPopup" v-bind="contactPopup" @contact="handleContactClick" @cancel="handleContactCancel"></contact-popup>
    <select-visitor-popup v-model="showSelectVisitorPopup" :default-selection="copyMember" confirmText="storeVisit.popup.btn.11" @select="handleVisitorSelect" @confirm="handleConfirm"></select-visitor-popup>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";

import ContactPopup from './components/contact-popup.vue'
import SelectVisitorPopup from './components/select-visitor-popup.vue'
import ComConfirmPopup from "./components/ComConfirmPopup";
import { addCopyMemberApi } from "@/api/store-visit";
@Component({
  components: {
    ContactPopup,
    SelectVisitorPopup,
  },
})
export default class StoreVisitDetail extends Vue {
  @Prop({ type: Object, default: () => ({}) }) detail!: StoreVisit.IDetail;
  showContactPopup = false
  showSelectVisitorPopup = false
  // showConfirmPopup = false

  contactPopup:any = {
    title: 'storeVisit.popup.title.5',
    tips: 'storeVisit.popup.tips.10',
    list: []
  }

  choose_scene = '同步指定人员'
  choose_entrance = '新增'

  get vistor() {
    const list = [...this.detail.visitor_member]
    if(!list.find(item => item.member_id === this.detail.creator.member_id)) {
      list.unshift(this.detail.creator)
    }
    return list
  }
  get copyMember() {
    return this.detail.copy_member
  }

  copyMemberList:any[] = []

  /** 用户信息 */
  get userInfo():IUserInfo {
    return this.$store.getters["user/getUserInfo"]
  }

  get isStoreManager() {
    if (!this.detail) return false
    
    return this.detail.store_manager.findIndex(i => i.member_id === this.userInfo.member_id) > -1 && this.userInfo.is_store_manager
  }

  get trackInfo() {
    return {
      SV_id: this.detail.id,
      store_name: this.detail.store_name,
      store_code: this.detail.store_code,
      SV_initiator: this.detail.creator.name,
      SV_initiator_id: this.detail.creator.user_id,
      SV_initiator_department: this.detail.creator.dept_name,
      visit_time: this.detail.date + ' ' + this.detail.time_span,
    }
  }

  created() {
    document.title = String(this.$t('common.pageTitle.6'))
  }

  handleContact(type) {
    this.showContactPopup = true
    if(type == 1) {
      this.contactPopup = {
        type,
        title: 'storeVisit.popup.title.5',
        tips: 'storeVisit.popup.tips.10',
        list: this.vistor,
        showConfirmButton: false,
      }

      // 神策数据上报 -- 巡店计划点击
      this.$track.trackStoreVisitAppointmentClick({
        ...this.trackInfo,
        button_name: '联系巡店人员',
      })
    } else {
      this.contactPopup = {
        type,
        title: 'storeVisit.popup.title.6',
        tips: 'storeVisit.popup.tips.11',
        list: this.copyMember,
        showConfirmButton: true,
        cancelText: 'storeVisit.popup.btn.10'
      }

      // 神策数据上报 -- 巡店计划点击
      this.$track.trackStoreVisitAppointmentClick({
        ...this.trackInfo,
        button_name: '查看同步的人员名单',
      })
    }
  }
  handleContactCancel() {
    if(this.contactPopup.type == 2) {
      this.showSelectVisitorPopup = true

      // 神策数据上报 -- 巡店计划点击
      this.$track.trackStoreVisitAppointmentClick({
        ...this.trackInfo,
        button_name: '新增同步人员',
      })

      this.choose_scene = '同步指定人员'
      this.choose_entrance = '新增'
    }
  }
  handleContactClick(e) {
    const { track } = e
    this.$track.trackFindContact({
      ...this.trackInfo,
      ...track,
    })
  }

  /**
   * 同步指定人员
   */
  handleSelectVisitor() {
    this.showSelectVisitorPopup = true
    // 神策数据上报 -- 巡店计划点击
    this.$track.trackStoreVisitAppointmentClick({
      ...this.trackInfo,
      button_name: '同步给指定人员',
    })
    this.choose_scene = '同步指定人员'
    this.choose_entrance = '添加'
  }

  handleVisitorSelect(e) {
    const { track } = e
    // 神策数据上报 -- 选择门店
    this.$track.trackChooseStaff({
      ...track,
      ...this.trackInfo,
      choose_scene: this.choose_scene,
      choose_entrance: this.choose_entrance,
    })
  }

  /**
   * 同步其他人员确定
   */
  async handleConfirm(value) {
    const list = value.filter(item => item.is_new).map(item =>{
      return {
        name: item.name,
        member_id: item.member_id||0,
        email: item.email,
        dept_name: item.dept_name
      }
    })
    if(list.length == 0) return
    
    const { action } = await ComConfirmPopup.confirm({
      title: this.$t('storeVisit.confirm.title.1'),
      content: `${this.$t('storeVisit.confirm.content.4')}${list.length}${this.$t('storeVisit.confirm.content.5')}`,
      confirmButtonText: this.$t('storeVisit.confirm.confirmBtn.4'),
      cancelButtonText: this.$t('storeVisit.confirm.cancelBtn.2'),
      footerStyle: '2'
    })

    if(action === 'confirm') {
      const params = {
        id: this.detail.id,
        copy_member_list: list
      }
      
      const res:IApiResult = await addCopyMemberApi(params)
      if(res.code === 200) {
        this.showSelectVisitorPopup = false
        this.$emit('update:detail')

        // 神策数据上报 -- 同步巡店计划
        this.$track.trackSyncStoreVisit({
          ...this.trackInfo,
          is_new_added: this.detail.copy_member.length > 0,
          sync_staff_list: list.map(item => `${item.name||item.email}: ${item.dept_name||'Email'}`)
        })

        // 神策数据上报 -- 巡店计划浏览
        this.$track.trackStoreVisitAppointmentView({
          ...this.trackInfo,
          is_synced: true,
        })

        this.$toast(this.$t('common.toast.1'))
      } else {
        this.$toast(res.message)
      }
    }
  }

}
</script>

<style lang="scss" scoped>
.detail {
  width: 100%;
  min-height: 100vh;
  background: rgba(248, 249, 250, 1) linear-gradient(180deg, #3975C6 0%, #0D6EFD 100%) center top/100%  148px no-repeat;
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
          display: flex;
          align-items: center;
          .name {
            font-size: 12px;
            margin-right: 4px;
          }

          .dept {
            font-weight: 700;
            font-size: 10px;
            padding: 0 4px;
            border-radius: 2px;
            background-color: rgba(255, 255, 255, 1);
            color: rgba(153, 153, 153, 1);
          }
        }

        &:not(:last-child) {
          margin-bottom: 4px;
        }
      }
    }
  }

  .content {
    padding: 8px;

    .remark {
      border-radius: 12px;
      background-color: rgba(255, 255, 255, 1);
      padding: 24px 16px;

      .remark-title {
        font-size: 14px;
        font-weight: 600;
        color: rgba(0, 0, 0, 1);
        line-height: 20px;
      }
      .remark-content {
        font-size: 12px;
        font-weight: 400;
        color: rgba(0, 0, 0, 1);
        line-height: 20px;
        margin-top: 16px;
      }
    }
  }

  .footer {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    padding: 12px 24px;
    display: flex;
    padding-bottom: calc(env(safe-area-inset-bottom) + 12px);
    padding-bottom: calc(constant(safe-area-inset-bottom) + 12px);

    .van-button {
      height: 44px;
      border-radius: 8px;
      font-size: 14px;
      font-weight: 600;
      flex: 1;
      
      &:nth-child(2) {
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
