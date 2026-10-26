<template>
  <div class="select">
    <div class="select" @click="handleClick">
      <slot>{{ showTime }}</slot>
    </div>
    <van-popup v-model="visible" position="bottom" :get-container="'body'" round @click-overlay="handleCancel">
      <div class="popup">
        <van-datetime-picker
          v-model="currentDate"
          :type="type"
          :title="String($t(title))"
          @cancel="handleCancel"
          @confirm="handleConfirm"
        />
      </div>
    </van-popup>
  </div>
</template>
<script lang="ts">
import { Component, Model, Prop, Vue, Watch } from 'vue-property-decorator';
import {DatetimePicker } from "vant";
import dayjs from 'dayjs';

@Component({
  components: {
    [DatetimePicker.name]: DatetimePicker,
  }
})
export default class SelectDatePopup extends Vue {
  @Model('update:value', { type: [Number,String], default: false }) value: number|string;
  @Prop({ type: String, default: 'storeVisit.popup.title.11' }) title: string;
  @Prop({ type: String, default: 'date' }) type: string;
  // 最小日期 默认为T+1
  @Prop({ type: Date, default: () => dayjs().add(1, 'day').toDate() }) minDate: Date;

  
  visible: boolean = false;

  currentDate = new Date()

  @Watch('value', { immediate: true })
  onValueChange(val) {
    if (!val) {
      this.currentDate = this.minDate;
      return;
    }
    this.currentDate = dayjs(val).toDate();
  }

  get showTime() {
    if (!this.value) {
      return '';
    }
    return dayjs(this.value).format('YYYY.MM.DD');
  }

  formatter(type, val) {
    if (type === 'year') {
      return `${val}年`;
    } else if (type === 'month') {
      return `${val}月`;
    }
    return `${val}日`;
  }
  handleClick() {
    this.visible = true;
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
  handleConfirm(value) {
    this.visible = false;
    this.$emit('update:value', dayjs(this.currentDate).format('YYYY-MM-DD'));
    this.$emit('confirm', this.currentDate);
  }

}
</script>
<style lang="scss" scoped>
.select {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
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
    background: rgba(248, 249, 250, 1);
    max-height: 450px;
    overflow-y: auto;

    .list {
      padding: 0 24px;

      .item {
        height: 64px;
        display: flex;
        align-items: center;

        &:not(:last-child) {
          border-bottom: 0.5PX solid rgba(233, 236, 239, 1);
        }

        .item-name {
          font-weight: 400;
          font-size: 14px;
          color: #000000;
          flex: 1;
          transition: all 0.3s;
        }

        .item-checkbox {
          width: 24px;
          height: 24px;
          background: url(~@/assets/store-visit/icon_check.png) center/24px no-repeat;
          opacity: 0;
          transition: all 0.3s;
        }

        &.active {
          .item-name {
            font-weight: 700;
          }
          .item-checkbox {
            opacity: 1;
          }
        }
      }
    }
  }

  .popup-footer {
    box-sizing: border-box;
    z-index: 2;
    display: flex;
    align-items: center;
    background-color: #fff;
    padding: 12px 24px;
    padding-bottom: calc(env(safe-area-inset-bottom) + 12px);
    padding-bottom: calc(constant(safe-area-inset-bottom) + 12px);

    .van-button {
      flex: 1;
      height: 40px;
      border-radius: 4px;

      &.van-button--normal {
        font-size: 12px;
        padding: 0 16px;
      }

      &:nth-child(2) {
        margin-left: 8px;
      }
    }
  }
}
</style>