<template>
  <div class="select">
    <div class="select-content" @click="handleClick">
      <slot>{{ innerValue }}</slot>
    </div>
    <van-popup v-model="visible" position="bottom" round @click-overlay="handleCancel">
      <div class="popup">
        <div class="popup-header">
          <div class="popup-cancel" @click="handleCancel">{{ $t('storeVisit.popup.btn.3') }}</div>
          <div class="popup-title">{{ $t(title) }}</div>
          <div class="popup-confirm" @click="handleConfirm">{{ $t('storeVisit.popup.btn.20') }}</div>
        </div>
        <div class="popup-tabs">
          <div class="tabs">
            <div class="tab" :class="{active: current === item.value}" v-for="item, index in tabs" :key="index" @click="handleTabClick(item.value)">{{ $t(item.label) }}</div>
          </div>
        </div>
        <div class="popup-content">
          <div v-if="current === 0" class="item">
            <select-year v-model="innerValue[0]" :defaultColumns="columns"></select-year>
          </div>
          <div v-if="current === 1" class="item">
            <select-quarter v-model="innerValue[1]" :defaultColumns="columns"></select-quarter>
          </div>
          <div v-if="current === 2" class="item">
            <select-month v-model="innerValue[2]" :defaultColumns="columns"></select-month>
          </div>
          <div v-if="current === 3" class="item">
            <select-calendar v-model="innerValue[3]"></select-calendar>
          </div>
        </div>
      </div>
    </van-popup>
  </div>
</template>
<script lang="ts">
import { Component, Model, Prop, Vue, Watch } from 'vue-property-decorator';
import SelectYear from './select-year.vue'
import SelectQuarter from './select-quarter.vue'
import SelectMonth from './select-month.vue'
import SelectCalendar from './select-calendar.vue'
import dayjs from 'dayjs';

@Component({
  components: {
    SelectYear,
    SelectQuarter,
    SelectMonth,
    SelectCalendar
  }
})
export default class SelectDaterangePopup extends Vue {
  @Prop({ type: String, default: 'storeVisit.popup.title.16' }) title: string;
  @Prop({ type: Array, default: () => [] }) columns: any[];
  @Prop({ type: Number, default: 0 }) defaultCurrent: number;
  @Prop({ type: Object, default: () => ({
    0: new Date().getFullYear(),
    1: [new Date().getFullYear(), 1],
    2: [new Date().getFullYear(), new Date().getMonth() + 1],
    3: []
  }) }) defaultValue: any;
  
  tabs = [
    { label: 'storeVisit.popup.title.17', value: 0 },
    { label: 'storeVisit.popup.title.18', value: 1 },
    { label: 'storeVisit.popup.title.19', value: 2 },
    { label: 'storeVisit.popup.title.20', value: 3 },
  ]

  current = 0
  
  visible: boolean = false;

  innerValue = {
    0: new Date().getFullYear(),
    1: [new Date().getFullYear(), 1],
    2: [new Date().getFullYear(), new Date().getMonth() + 1],
    3: []
  }

  @Watch('defaultValue', { immediate: true })
  onValueChange(val) {
    this.innerValue = JSON.parse(JSON.stringify(val));
  }
  
  handleClick() {
    this.visible = true;
    this.innerValue = JSON.parse(JSON.stringify(this.defaultValue));
    this.current = this.defaultCurrent
  }
  /**
   * 取消
   */
  handleCancel() {
    this.visible = false;
    this.$emit('cancel');
  }

  /**
   * 确定
   */
  handleConfirm() {
    const value = this.innerValue[this.current]
    // console.log('value', value)
    switch (this.current) {
      case 0:
        // 财年
        this.$emit('confirm', {
          value,
          current: this.current,
          data: this.columns[0].children.find(item => item.value === value).data
        });
        break;
      case 1:
        // 财季
        this.$emit('confirm', {
          value,
          current: this.current,
          data: this.columns[1].children.find(item => item.value === value[1] && item.year === value[0]).data
        });
        break;
      case 2:
        this.$emit('confirm', {
          value,
          current: this.current,
          data: this.columns[2].children.find(item => item.value === value[1] && item.year === value[0]).data
        });
        break;
      case 3:
        if(value.length === 0) return
        this.$emit('confirm', {
          value,
          current: this.current,
          data: {
            start_time: dayjs(`${value[0]} 00:00:00`).unix(),
            end_time: dayjs(`${value[1]} 23:59:59`).unix()
          }
        });
        break;
    }


    this.visible = false;
    
  }

  handleTabClick(value) {
    this.current = value;
  }

}
</script>
<style lang="scss" scoped>
.tabs {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: rgba(248, 249, 250, 1);
  padding: 2px;
  border-radius: 4px;

  .tab {
    flex: 1;
    height: 20px;
    line-height: 20px;
    text-align: center;
    font-size: 12px;
    color: rgba(0, 0, 0, 0.3);
    border-radius: 3px;
    transition: all 0.2s;

    &.active {
      font-weight: bold;
      color: rgba(0, 0, 0, 1);
      background-color: rgba(255, 255, 255, 1);
    }

    &:nth-child(1) {
      margin-right: 2px;
    }
  }
}

.popup {
  position: relative;
  width: 100vw;
  background: #FFFFFF;
  .popup-header {
    padding: 20px 24px;
    position: relative;
    display: flex;
    justify-content: space-between;
    align-items: center;
    
    .popup-title {
      font-weight: 600;
      font-size: 16px;
      line-height: 24px;
      color: #000000;
    }

    .popup-cancel {
      width: 32px;
      height: 24px;
      font-size: 12px;
      font-weight: 600;
      color: rgba(108, 117, 125, 1);
      display: flex;
      align-items: center;
      justify-content: center;
    }
    .popup-confirm {
      width: 32px;
      height: 24px;
      font-size: 12px;
      font-weight: 600;
      color: rgba(13, 110, 253, 1);
      display: flex;
      align-items: center;
      justify-content: center;
    }

  }
  .popup-tabs {
    padding: 8px 16px;


  }

  .popup-content {
    background-color: rgba(248, 249, 250, 1);
    min-height: 304px;
    padding-bottom: env(safe-area-inset-bottom);
    padding-bottom: constant(safe-area-inset-bottom);

  }
}
</style>