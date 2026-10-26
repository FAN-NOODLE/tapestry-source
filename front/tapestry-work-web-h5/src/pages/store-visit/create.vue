<template>
  <div class="container">
    <div class="header">
      <div class="flex">
        <div class="title">{{userInfo.dept_name}}</div>
        <!-- <div class="dept">部门：销售反馈</div> -->
      </div>
    </div>
    <div class="content">
      <div class="section">
        <div class="section-head">
          <div class="section-title">{{ $t('storeVisit.tips.7') }}</div>
          <div class="use-time">
            <span>{{ $t('storeVisit.tips.8') }}</span>
            <van-switch v-model="useTime" @change="handleUseTimeChange" />
          </div>
        </div>
        <div class="section-content"> 
          <div class="store-item" v-for="(item, index) in formData.store_list" :key="index">
            <div class="store-item-head">
              <div class="store-item-title">{{ item.store_name }}</div>
              <div v-if="!isStoreManager" class="icon-close-circle" @click="handleDeleteStore(index)"></div>
            </div>
            <div class="store-item-body">
              <div class="row">
                <div class="label">{{ $t('storeVisit.tips.9') }}</div>
                <div class="value">
                  <select-date-popup class="month" v-model="item.date" @confirm="handleDateChange(index)"></select-date-popup>
                  <select-time-popup ref="timePopup" class="time" :type="useTime" v-model="item.time_span"></select-time-popup>
                </div>
              </div>
            </div>
          </div>
          <div v-if="!isStoreManager" class="store-item" @click="handleSelectStore">
            <div class="store-item-head">
              <div class="store-item-title add">{{ $t('storeVisit.tips.10') }}</div>
              <div class="icon-plus-blue"></div>
            </div>
          </div>
        </div>
      </div>
      <div class="section">
        <div class="section-head">
          <div class="section-title">{{ $t('storeVisit.tips.46') }}</div>
          <div class="value" @click="handleSelectVisitor">
            <div>{{ $t('storeVisit.btn.28') }}</div>
            <div class="icon-right"></div>
          </div>
        </div>
        <div class="section-content"> 
          <div class="user-list">
            <!-- 自己 -->
            <div class="user-item">
              <div class="flex-1">
                <div class="user-name">@ {{userInfo.name||userInfo.email}}</div>
                <div class="user-dept" :class="{ext: userInfo.dept_name === 'Email'}">{{ userInfo.dept_name }}</div>
              </div>
            </div>
            <!-- 门店店长 -->
            <div class="user-item" v-for="(item, index) in managerList" :key="index">
              <div class="flex-1">
                <div class="user-name">@ {{item.member_name || item.name}}</div>
                <div class="user-dept">{{ item.store_code }}</div>
                <div class="user-dept">{{ item.position_abbr }}</div>
              </div>
              <div class="icon-close-circle" @click="handleDeleteVisitor(index, 'manager')"></div>
            </div>
            <!-- 巡店人员 -->
            <div class="user-item" v-for="(item, index) in formData.visitor_list" :key="index">
              <div class="flex-1">
                <div class="user-name">@ {{item.name||item.email}}</div>
                <div class="user-dept" :class="{ext: item.dept_name === 'Email'}">{{ item.dept_name }}</div>
              </div>
              <div class="icon-close-circle" @click="handleDeleteVisitor(index, 'visitor')"></div>
            </div>
          </div>
        </div>
      </div>
      <div class="section">
        <div class="section-head">
          <div class="section-title">{{ $t('storeVisit.tips.6') }}</div>
        </div>
        <div class="section-content"> 
          <textarea v-model="formData.remark" class="textarea" :placeholder="String($t('storeVisit.tips.6'))"></textarea>
        </div>
      </div>
    </div>
    <div class="footer">
      <van-button color="#000000" @click="handleBeforeCreate">{{ $t('storeVisit.btn.19') }}</van-button>
    </div>
    <select-store-popup v-model="showSelectStorePopup" :storeTree="storeTree" :defaultValue="defaultValue" @select="handleStoreSelect" @confirm="handleSelectStoreConfirm"></select-store-popup>
    <select-visitor-popup v-model="showSelectVisitorPopup" :default-selection="selection" @select="handleVisitorSelect" @confirm="handleSelectVisitorConfirm"></select-visitor-popup>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Watch } from "vue-property-decorator";
import { Button, Switch } from "vant";
import SelectStorePopup from "./components/SelectStorePopup/index.vue";
import SelectVisitorPopup from "./components/select-visitor-popup.vue";
import SelectDatePopup from "./components/select-date-popup.vue";
import SelectTimePopup from "./components/select-time-popup.vue";
import ComPopup from "./components/ComConfirmPopup";
import { saveVisitPlanApi, getStoreManagerApi } from "@/api/store-visit";
import dayjs from "dayjs";
import { getStoreTreeApi } from "@/api/display-task";
@Component({
  components: {
    [Button.name]: Button,
    [Switch.name]: Switch,
    SelectStorePopup,
    SelectVisitorPopup,
    SelectDatePopup,
    SelectTimePopup
  },
})
export default class StoreVisitCreate extends Vue {
  useTime = false;
  showSelectStorePopup = false;
  showSelectVisitorPopup = false;

  formData = {
    store_list: [],
    visitor_list: [],
    send_msg_to_store: 1,
    remark: "",
  }
  storeTree = []
  managerList = []


  get defaultValue() {
    return this.formData.store_list.map(item => item.third_id)
  }

  /** 用户信息 */
  get userInfo():IUserInfo {
    return this.$store.getters["user/getUserInfo"]
  }

  get isStoreManager() {
    return this.userInfo.is_store_manager && !this.userInfo.is_store_visit_manager
  }

  get selection() {
    return this.formData.visitor_list.concat(this.managerList)
  }

  @Watch('formData.store_list.length')
  onStoreListChange() {
    this.getStoreManager()
  }

  mounted() {
    const { store_list } = this.$route.query
    if (store_list) {
      this.formData.store_list = JSON.parse(store_list as string).map((item: any) => {
        return {
          ...item,
          date: '',
          time_span: '',
          unix_time: '',
        }
      })
      
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


  handleUseTimeChange() {
    this.formData.store_list.forEach(item => {
      item.time_span = '';
    });
  }

  handleSelectStore() {
    this.showSelectStorePopup = true;
  }
  handleSelectStoreConfirm(data) {
    this.formData.store_list = data.map(item => {
      const val = this.formData.store_list.find(store => store.third_id === item.third_id) || {};
      return {
        third_id: item.third_id,
        store_name: item.store_name,
        store_code: item.store_code,
        date: '',
        time_span: '',
        unix_time: '',
        ...val, // 保留已有的值
      }
    });
    this.showSelectStorePopup = false;
  }
  handleSelectStoreCancel() {
    this.showSelectStorePopup = false;
  }

  handleDeleteStore(index) {
    this.formData.store_list.splice(index, 1);
  }

  handleDateChange(index) {
    this.$refs.timePopup[index].visible = true;
  }


  handleStoreSelect(e) {
    const { track, item } = e
    // 神策数据上报 -- 选择门店
    this.$track.trackChooseStore({
      filter_option: '',
      search_word: '',
      choose_scene: '预约',
      choose_entrance: '新增巡访门店',
      store_name: item.store_name,
      store_code: item.store_code,
      ...track,
    })
  }

  handleSelectVisitor() {
    this.showSelectVisitorPopup = true;
  }
  handleSelectVisitorConfirm(data) {
    const storeList = this.formData.store_list.map(item => item.third_id)
    
    const memberIds = this.managerList.map(item => item.member_id)
    const managers = data.filter(item => item.position_abbr && storeList.includes(item.third_id)).map(item => {
      const currentStore = this.formData.store_list.find(store => store.third_id === item.third_id)
      return {
        ...item,
        member_name: item.member_name || item.name,
        store_code: currentStore?.store_code || '',
      }
    })
    const visitors = data.filter(item => !(item.position_abbr && storeList.includes(item.third_id)))
    // 过滤
    this.formData.visitor_list = visitors;
    this.managerList = managers
    this.showSelectVisitorPopup = false;
  }
  handleVisitorSelect(e) {
    const { track } = e
    // 神策数据上报 -- 选择门店
    this.$track.trackChooseStaff({
      SV_id: '',
      store_name: '',
      store_code: '',
      SV_initiator: this.userInfo.name,
      SV_initiator_id: this.userInfo.user_id,
      SV_initiator_department: this.userInfo.dept_name,
      choose_scene: '添加巡店人员',
      choose_entrance: '添加',
      ...track,
    })
  }
  handleSelectVisitorCancel() {
    this.showSelectVisitorPopup = false;
  }

  async handleDeleteVisitor(index, type) {
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
    } else {
      this.formData.visitor_list.splice(index, 1);
    }
  }

  async getStoreManager() {
    const res:any = await getStoreManagerApi({
      third_ids: this.formData.store_list.map(item => item.third_id)
    });
    if (res.code !== 200) {
      this.$toast.fail(res.message);
      return
    }
    this.managerList = (res.data || []).map(item => {
      return {
        ...item,
        name: item.member_name,
        dept_name: item.store_code,
        is_manager: 1,
      }
    })
  }

  handleBeforeCreate = this.$helper.debounce(async () => {
    this.handleCreate()
  }, 500, true)

  async handleCreate() {
    const params = JSON.parse(JSON.stringify(this.formData));

    if (!params.store_list.length) {
      this.$toast(this.$t('common.toast.3'));
      return;
    }

    const noTime = params.store_list.some(item => !item.time_span || !item.date)
    if (noTime) {
      this.$toast(this.$t('common.toast.4'));
      return;
    }
    

    this.$toast.loading(this.$t('common.toast.5'));
    console.log(params);
    params.store_list.forEach(item => {
      
      if(item.time_span === this.$t('storeVisit.popup.tips.29')) {
        item.unix_time = dayjs(item.date).hour(10).minute(0).second(0).unix();
      } else if(item.time_span === this.$t('storeVisit.popup.tips.30')) {
        item.unix_time = dayjs(item.date).hour(14).minute(0).second(0).unix();
      } else {
        const [hour] = item.time_span.split(':');
        item.unix_time = dayjs(item.date).hour(Number(hour)).minute(0).second(0).unix();
      }

      item.store_manager = this.managerList.filter(v => v.third_id === item.third_id).map(v => v.member_id || 0)
    })
    params.visitor_list = params.visitor_list.map(item => {
      return {
        email: item.email,
        member_id: item.member_id || 0,

      }
    });
    try {
      const res:any = await saveVisitPlanApi(params);
      this.$toast.clear()
      if (res.code === 200) {
        this.track(res.data)
        // this.$toast('创建成功');
        this.$store.commit('storeVisit/setToast', this.$t('common.toast.2'))
        this.$router.back()
      } else {
        this.$toast(res.message);
      }
      
    } catch (error) {
      this.$toast.clear()
    }
  }

  track(res) {
    const params = JSON.parse(JSON.stringify(this.formData))
    const sync_staff_list = params.visitor_list.map(item => {
      return `${item.name || item.email}:${item.dept_name || 'Email'}`;
    });
    const list = params.store_list.map(item => {
      return {
        third_id: item.third_id,
        store_name: item.store_name,
        store_code: item.store_code,
        visit_time: `${item.store_name}:${item.date} ${item.time_span}`,
      }
    })
    const data = {
      is_specific_time: this.useTime,
      is_SM_notified: params.send_msg_to_store === 1,
      sync_staff_list
    }
    // 神策数据上报 -- 发布巡店计划
    this.$track.trackSubmitAppointment({
      store_name_list: list.map(item => item.store_name),
      store_code_list: list.map(item => item.store_code),
      visit_time_list: list.map(item => item.visit_time),
      ...data
    })
    // 神策数据上报 -- 发布巡店计划详情
    list.forEach(item => {
      const val = res.find(i => i.third_id === item.third_id) || {};
      this.$track.trackSubmitAppointmentDetail({
        SV_id: val.id,
        store_name: item.store_name,
        store_code: item.store_code,
        visit_time: item.visit_time,
        ...data
      })
    })
  }
}
</script>

<style lang="scss" scoped>
.icon-plus-blue {
  width: 16px;
  height: 16px;
  background: url(~@/assets/store-visit/icon_plus_blue.png) center/100% no-repeat;
}
.icon-close-circle {
  width: 16px;
  height: 16px;
  background: url(~@/assets/store-visit/icon_close_circle_2.png) center/16px no-repeat;
}
.icon-right {
  width: 16px;
  height: 16px;
  background: url(~@/assets/store-visit/icon_right.png) center/16px no-repeat;
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
.container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: rgba(248, 249, 250, 1);
  .header {
    flex: 0 0 56px;
    background: linear-gradient(180deg, rgba(57, 117, 198, 1) 0%, rgba(24, 112, 232, 1) 100%);
    
    .flex {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 0 24px;
      height: 100%;

      .title {
        font-size: 16px;
        font-weight: 700;
        color: rgba(255, 255, 255, 1);
      }

      .dept {
        font-size: 12px;
        font-weight: 600;
        color: rgba(255, 255, 255, 0.8);
      }
    }
  }
  .content {
    flex: 1;
    overflow-y: auto;
    padding: 8px;
    padding-bottom: calc(env(safe-area-inset-bottom) + 76px);
    padding-bottom: calc(constant(safe-area-inset-bottom) + 76px);

    .section {
      background-color: rgba(255, 255, 255, 1);
      border-radius: 16px;
      padding: 16px;
      margin-bottom: 8px;

      .section-head {
        display: flex;
        align-items: center;
        justify-content: space-between;
        height: 40px;

        .section-title {
          font-size: 16px;
          font-weight: 600;
          color: rgba(0, 0, 0, 1);
        }

        .use-time {
          font-size: 10px;
          font-weight: 600;
          display: flex;
          align-items: center;
          &>span {
            margin-right: 4px;
          }
        }

        .value {
          font-size: 14px;
          color: rgba(13, 110, 253, 1);
          display: flex;
          align-items: center;
        }
      }
      .section-content {
        margin-top: 16px;

        .textarea {
          width: 100%;
          height: 112px;
          padding: 0;
          display: block;
          font-size: 12px;
          line-height: 20px;

          &::placeholder {
            color: rgba(153, 153, 153, 0.5);
          }
        }

        .row {
          display: flex;
          align-items: center;
          height: 28px;

          &:not(:last-child) {
            margin-bottom: 12px;
          }

          .label {
            flex: 1;
            font-size: 12px;
            color: rgba(0, 0, 0, 0.8);
          }

          .value {
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

        .store-item {
          &:not(:last-child) {
            padding-bottom: 12px;
            margin-bottom: 12px;
            border-bottom: 1px solid var(--Grey-100, rgba(248, 249, 250, 1))
          }

          .store-item-head {
            height: 28px;
            display: flex;
            align-items: center;

            .store-item-title {
              flex: 1;
              font-size: 14px;

              &.add {
                color: rgba(13, 110, 253, 1);
              }
            }
          }

          .store-item-body {
            margin-top: 8px;
            
          }
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
      }
    }
  }
  .footer {
    position: fixed;
    bottom: 0;
    left: 0;
    width: 100vw;
    box-sizing: border-box;
    padding: 12px 24px;
    display: flex;
    padding-bottom: calc(env(safe-area-inset-bottom) + 12px);
    padding-bottom: calc(constant(safe-area-inset-bottom) + 12px);

    .van-button {
      flex: 1;
      border-radius: 8px;
      font-size: 14px;
      font-weight: 600;
      height: 44px;
    }
  }
}

</style>