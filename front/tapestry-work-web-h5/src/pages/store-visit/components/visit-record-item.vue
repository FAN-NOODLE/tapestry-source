<template>
  <div v-if="item.id" class="item">
    <div ref="draggable" class="item-inner" @touchstart="handleTouchStart" @touchmove="handleTouchMove">
      <div class="item-title">{{ item.store_name }}</div>
      <div class="item-row">
        <div class="item-label">{{ $t('storeVisit.item.label.1') }}:</div>
        <div class="item-value">{{item.date.replaceAll('-', '.')}} {{item.time_span}}</div>
      </div>
      <div class="item-row">
        <div class="item-label">{{ $t('storeVisit.item.label.2') }}:</div>
        <div class="item-value">
          <span>{{item.creator.name}}</span>
          <span v-if="item.creator.dept_name" class="dept">{{item.creator.dept_name}}</span>
        </div>
      </div>
      <div class="item-row">
        <div class="item-label">{{ $t('storeVisit.item.label.3') }}:</div>
        <div v-if="item.store_manager.length" class="item-value">
          <div class="name" v-for="(m, index) in item.store_manager" :key="index">{{ m.name }}</div>
        </div>
        <div v-else class="item-value">--</div>
      </div>
      <div class="item-row" v-if="item.copy_member.length">
        <div class="item-label">{{ $t('storeVisit.item.label.4') }}:</div>
        <div class="item-value">
          <div class="name" v-for="(m, index) in item.copy_member" :key="index">{{ m.name||m.email }}</div>
        </div>
      </div>
      <div class="item-footer" :class="{en: lang === 'en-US'}">
        <van-button color="#000000" plain @click="handleContactStoreManager">{{ $t('storeVisit.btn.9') }}</van-button>
        <!-- 创建人并且是预约巡店后 -->
        <van-button v-if="isCreator(item) && item.status === 1" color="#000000" @click="handleStart('开始巡店')">{{ $t('storeVisit.btn.11') }}</van-button>
        <van-button v-else-if="isCreator(item) && item.status === 2" color="#000000" @click="handleStart('继续巡店')">{{ $t('storeVisit.btn.12') }}</van-button>
        <van-button v-else color="#000000" plain @click="handleViewDetail">{{ $t('storeVisit.btn.10') }}</van-button>
      </div> 
      <!-- 右侧删除 -->
      <div ref="right" class="item-right">
        <div class="btn" @click="handleDelete">
          <div class="icon-delete"></div>
          <div>{{ $t('storeVisit.btn.7') }}</div>
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Model, Prop, Vue, Watch } from 'vue-property-decorator';
import ContactPopup from './ContactPopup';
import dayjs from 'dayjs';
import ComConfirmPopup from './ComConfirmPopup';

@Component({
  components: {
    // ContactPopup
  },
})
export default class VisitRecordItem extends Vue {
  @Prop({ type: Object, default: () => {} }) item: StoreVisit.IDetail

  // showContactPopup = false
  startX = 0
 
  /** 用户信息 */
  get userInfo():IUserInfo {
    return this.$store.getters["user/getUserInfo"]
  }

  get lang() {
    return this.$store.getters["lang/getLang"]
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

  get trackInfo() {
    const data: any = {
      SV_id: this.item.id,
      store_name: this.item.store_name,
      store_code: this.item.store_code,
      SV_initiator: this.item.creator.name,
      SV_initiator_id: this.item.creator.user_id,
      SV_initiator_department: this.item.creator.dept_name,
      visit_time: this.item.date + ' ' + this.item.time_span,
      store_manager: this.item.store_manager.map(i => i.name).join(','),
      plan_status: ['进行中', '进行中', '已完成'][this.item.status -1]
    }

    if(this.item.status === 3) {
      data.report_time = dayjs(this.item.updated_at*1000).format('YYYY-MM-DD HH:mm:ss')
    }

    return data
  }


  handleContactStoreManager() {
    // this.showContactPopup = true
    ContactPopup.open({
      title: this.$t('storeVisit.popup.title.4'),
      tips: this.$t('storeVisit.popup.tips.9'),
      list: this.item.store_manager,
      cancelBtn: true,
    })
    // 神策数据上报 -- 巡店记录点击
    this.$track.trackSVPlanClick({
      ...this.trackInfo,
      button_name: '联系店长',
    })
  }

  handleStart(name) {
    // 神策数据上报 -- 巡店记录点击
    this.$track.trackSVPlanClick({
      ...this.trackInfo,
      button_name: name,
    })
    this.$router.push({
      name: 'store-visit-start',
      query: {
        id: String(this.item.id),
        entry: '巡店记录',
        btn: name,
      }
    })
    this.item.status = 2
  }

  handleViewDetail() {
    // 神策数据上报 -- 巡店记录点击
    this.$track.trackSVPlanClick({
      ...this.trackInfo,
      button_name: '查看详情',
    })
    this.$router.push({
      name: 'store-visit-report',
      query: {
        id: String(this.item.id),
        from: '1',
      }
    })
  }
  
  handleTouchStart(e) {
    this.startX = e.changedTouches[0].pageX;
  }
  handleTouchMove(e) {
    const endX = e.changedTouches[0].pageX;
    const dx = endX - this.startX
    
    if (Math.abs(dx) > 50) {
      const dom = this.$refs.draggable as HTMLElement
      const rightDom = this.$refs.right as HTMLElement
      
      if(dx < 0) {
        // dom.style.transform = `translateX(${0-rightDom.clientWidth}px)`
        dom.style.transform = `translateX(0px)`
      } else {
        dom.style.transform = `translateX(0px)`
      }
    }
  }

  async handleDelete() {
    const { action } = await ComConfirmPopup.confirm({
      title: "温馨提示",
      content: "删除后无法恢复，请确认是否删除",
      confirmButtonText: "删除巡店记录",
      cancelButtonText: "取消",
    })
    if(action === 'cancel') {
      return
    }
    // TODO 删除
    (this.$refs.draggable as HTMLElement).style.transform = `translateX(0px)`
    this.$emit('delete', this.item)
    this.$toast(this.$t('common.toast.21'))
  }

}

</script>
<style lang="scss" scoped>
.icon-delete {
  width: 24px;
  height: 24px;
  background: url("~@/assets/store-visit/icon_delete_w.png") center/24px no-repeat;
}
.item {
  border-radius: 12px;
  background-color: #fff;
  margin-bottom: 12px;
  overflow: hidden;
  .item-inner {
    padding: 16px;
    position: relative;
    transition: all .3s;

    .item-right {
      position: absolute;
      right: 0;
      top: 0;
      bottom: 0;
      width: 64px;
      background: #DC3545;
      display: flex;
      align-items: center;
      justify-content: center;
      transform: translateX(100%);

      .btn {
        display: flex;
        align-items: center;
        justify-content: center;
        flex-direction: column;
        color: #fff;
        font-size: 12px;

        .icon-delete {
          margin-bottom: 4px;
        }
      }
    }
  }

  .item-title {
    font-size: 14px;
    font-weight: 600;
    line-height: 20px;
    margin-bottom: 12px;
  }

  .item-row {
    display: flex;
    &:nth-child(n+1) {
      margin-bottom: 4px;
    }
  }
  .item-label {
    font-size: 12px;
    line-height: 16px;
    color: rgba(153, 153, 153, 1);
    margin-right: 8px;

  }
  .item-value { 
    font-size: 12px;
    line-height: 16px;
    color: rgba(0, 0, 0, 1);
    flex: 1;
    overflow: hidden;
    display: flex;
    flex-wrap: wrap;

    .name:not(:last-child) {
      margin-right: 4px;
      &::after {
        content: ',';
      }
    }

    .dept {
      padding: 2px 4px;
      border-radius: 2px;
      font-size: 10px;
      line-height: 12px;
      font-weight: 700;
      margin-left: 4px;
      background-color: rgba(206, 212, 218, 1);
    }
  }

  .item-footer {
    margin-top: 16px;
    padding-top: 16px;
    border-top: 1px solid rgba(233, 236, 239, 1);
    display: flex;

    .van-button {
      height: 40px;
      border-radius: 4px;
      font-size: 12px;
      font-weight: 600;
      flex: 1;
      

      &:nth-child(2) {
        margin-left: 8px;
      }
    }

    &.en {
      .van-button:nth-child(1) {
        flex: none;
        width: fit-content;
      }
    }
  }
}
</style>