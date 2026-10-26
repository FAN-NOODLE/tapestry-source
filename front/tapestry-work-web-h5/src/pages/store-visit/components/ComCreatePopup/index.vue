<template>
  <div>
    <van-popup v-model="visible" position="bottom" round @click-overlay="handleClose">
      <div class="popup">
        <div class="popup-content">
          <div class="section">
            <div class="section-head">
              <div class="section-title">{{ $t(title) }}</div>
              <div class="use-time">
                <span>{{ $t('storeVisit.tips.8') }}</span>
                <van-switch v-model="useTime" @change="handleUseTimeChange" />
              </div>
            </div>
            <div class="section-content"> 
              <div class="store-item" v-for="(item, index) in formData.store_list" :key="index">
                <div class="store-item-body">
                  <div class="row" @click="handleSelectStore()">
                    <div class="label">{{ $t('storeVisit.tips.38') }}</div>
                    <div class="value">
                      <span>{{ item.store_name }}</span>
                      <span class="icon icon-right"></span>
                    </div>
                  </div>
                  <div class="row">
                    <div class="label">{{ $t('storeVisit.tips.9') }}</div>
                    <div class="value">
                      <select-date-popup class="month" v-model="item.date" @confirm="handleDateChange(index)"></select-date-popup>
                      <select-time-popup ref="timePopup" class="time" :type="useTime" v-model="item.time_span"></select-time-popup>
                    </div>
                  </div>
                  <div class="row">
                    <div class="label">{{ $t('storeVisit.tips.39') }}</div>
                    <div class="value">
                      <div class="manager">
                        <span v-if="managerList.length === 0">--</span>
                        <span v-for="manager in managerList" class="manager-item" :key="manager.member_id">{{ manager.member_name }}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div v-if="footerStyle=='2'" class="popup-footer-2">
          <van-button color="#000000" block @click.stop="handleConfirm">{{ $t(confirmButtonText) }}</van-button>
          <van-button color="#000000" block plain @click.stop="handleCancel">{{ $t(cancelButtonText) }}</van-button>
        </div>
        <div v-else class="popup-footer">
          <van-button color="#000000" plain @click.stop="handleCancel">{{ $t(cancelButtonText) }}</van-button>
          <van-button color="#000000" @click.stop="handleConfirm">{{ $t(confirmButtonText) }}</van-button>
        </div>
      </div>
    </van-popup>
  </div>
</template>
<script lang="ts">
import { Component, Model, Prop, Vue, Watch } from 'vue-property-decorator';
import SelectDatePopup from "../select-date-popup.vue";
import SelectTimePopup from "../select-time-popup.vue";
import SelectStore from "../SelectStorePopup";
import { Button, Switch } from "vant";
import { getStoreManagerApi } from '@/api/store-visit';
import dayjs from 'dayjs';

@Component({
  components: {
    [Button.name]: Button,
    [Switch.name]: Switch,
    SelectDatePopup,
    SelectTimePopup
  }
})
export default class ComCreatePopup extends Vue {
  @Model('update:show', { type: Boolean, default: false }) show: boolean;
  @Prop({ type: String, default: 'storeVisit.popup.title.24' }) title: string;
  @Prop({ type: String, default: 'storeVisit.popup.btn.3' }) cancelButtonText: string;
  @Prop({ type: String, default: 'storeVisit.popup.btn.23' }) confirmButtonText: string;
  @Prop({ type: Array, default: [] }) store_list: any[];
  @Prop({ type: Array, default: [] }) storeTree: any[];
  @Prop({ type: String, default: '' }) footerStyle: string;
  @Prop({ type: Boolean, default: false }) is_store_manager: boolean;
  
  visible: boolean = false;
  useTime: boolean = false;

  formData = {
    store_list: [],
    visitor_list: [],
    send_msg_to_store: 1,
    remark: "",
  }
  managerList:any[] = [];

  current: number = -1;
  
  @Watch('show')
  onShowChange(val: boolean) {
    this.visible = val;
    if(val) {
      this.formData.store_list = this.store_list.map((item: any) => {
        return {
          ...item,
          date: dayjs().format('YYYY-MM-DD'),
          time_span: dayjs().isAfter(dayjs().hour(12).minute(59).second(59)) ? '下午' : '上午',
          unix_time: dayjs().unix(),
        }
      })
      this.getStoreManager()
    }
  }

  async handleSelectStore() {
    if(this.is_store_manager) {
      return
    }
    const current = this.formData.store_list[0]
    const res = await SelectStore.open({
      multiple: false,
      storeTree: this.storeTree,
      defaultValue: current.third_id,
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
    if(res.action === 'confirm') {
      this.formData.store_list = [{
        ...res.data[0],
        date: current.date,
        time_span: current.time_span,
        unix_time: current.unix_time,
      }]
      this.getStoreManager()
    }

  }

  /**
   * 切换使用时间
   */
  handleUseTimeChange() {
    this.formData.store_list.forEach(item => {
      item.time_span = '';
    });

    console.log(this.formData);
    
  }

  handleDateChange(index) {
    this.$refs.timePopup[index].visible = true;
  }

  async getStoreManager() {
    const res:any = await getStoreManagerApi({
      third_ids: this.formData.store_list.map(item => item.third_id)
    });
    if (res.code !== 200) {
      this.$toast.fail(res.message);
      return
    }
    this.managerList = res.data || []
  }

  /**
   * 取消
   */
  handleClose() {
    this.$emit('update:show', false);
    this.$emit('close');
  }
  /**
   * 取消
   */
  handleCancel() {
    this.$emit('update:show', false);
    this.$emit('confirm', {
      action: 'cancel'
    });
  }

  /**
   * 确认
   */
  handleConfirm = this.$helper.debounce(() => {
    // 触发确认事件
    this.emitConfirm();
  }, 500, true)
  
  /**
   * 确认事件
   */
  emitConfirm() {
    const storeInfo = this.formData.store_list[0]
        
    const params:any = {
      third_id: storeInfo.third_id,
      time_span: storeInfo.time_span,
      date: storeInfo.date
    }
    if(storeInfo.time_span === this.$t('storeVisit.popup.tips.29')) {
      params.unix_time = dayjs(storeInfo.date).hour(10).minute(0).second(0).unix();
    } else if(storeInfo.time_span === this.$t('storeVisit.popup.tips.30')) {
      params.unix_time = dayjs(storeInfo.date).hour(14).minute(0).second(0).unix();
    } else {
      const [hour] = storeInfo.time_span.split(':');
      params.unix_time = dayjs(storeInfo.date).hour(Number(hour)).minute(0).second(0).unix();
    }
    if(params.time_span === '') {
      this.$toast(this.$t('common.toast.4'));
      return
    }


    this.$emit('update:show', false);
    this.$emit('confirm', {
      action: 'confirm',
      data: params
    });
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

.icon-right {
  width: 16px;
  height: 16px;
  background: url(~@/assets/store-visit/icon_right_black.png) center/100% no-repeat;
}

.popup {
  position: relative;
  width: 100vw;
  background: #FFFFFF;
  .popup-header {
    padding: 20px 24px;
    position: relative;
    
    .popup-title {
      font-weight: 600;
      font-size: 16px;
      line-height: 24px;
      color: #000000;
    }

    .popup-close {
      width: 72px;
      height: 64px;
      position: absolute;
      right: 0;
      top: 0;
      background: url(~@/assets/store-visit/icon_popup_close.png) center/24px no-repeat;
    }
  }
  .popup-content {
    padding: 16px 24px;
    font-size: 12px;
    line-height: 130%;
  }

  .popup-footer {
    z-index: 2;
    display: flex;
    align-items: center;
    background-color: #fff;
    padding: 24px 24px;
    box-sizing: border-box;

    .van-button {
      height: 40px;
      border-radius: 4px;

      &:nth-child(1) {
        width: fit-content;
        min-width: 88px;
      }
      &:nth-child(2) {
        flex: 1;
        margin-left: 8px;
      }
      &.van-button--normal {
        font-size: 12px;
        padding: 0 16px;
      }

      &.van-button--disabled {
        opacity: 1;
        background: rgba(222, 226, 230, 1) !important;
        border-color: rgba(222, 226, 230, 1) !important;
      }
    }
  }

  .popup-footer-2 {
    padding: 24px 24px;
    .van-button {
      height: 40px;
      border-radius: 4px;

      &:nth-child(2) {
        margin-top: 8px;
      }
    }
  }
}


.manager {
  // display: flex;
  // align-items: center;
  // flex-wrap: wrap;
  // justify-content: flex-end;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;

  .manager-item {
    &:not(:last-child):after {
      content: '，'
    }
  }
}

.section {
  background-color: rgba(255, 255, 255, 1);
  border-radius: 16px;
  // padding: 16px;
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
        width: 60px;
        font-size: 12px;
        color: rgba(0, 0, 0, 0.8);
        flex-shrink: 0;
      }

      .value {
        flex: 1;
        display: flex;
        align-items: center;
        justify-content: flex-end;
        overflow: hidden;

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
</style>