<template>
  <van-popup v-model="visible" position="center" round :close-on-click-overlay="false">
    <div class="popup">
      <div class="popup-content">
        <img class="img" src="@/assets/store-visit/clock_banner.png" alt="">
        <div class="tips">{{ $t('storeVisit.popup.tips.7') }}</div>
      </div>
      <div class="popup-footer">
        <van-button color="#0D6EFD" block @click="handleConfirm">{{ $t('storeVisit.popup.btn.6') }}</van-button>
        <div class="btn" @click="handleCancel">{{ $t('storeVisit.popup.btn.7') }}</div>
      </div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Model, Prop, Vue, Watch } from 'vue-property-decorator';

@Component({

})
export default class ClockPopup extends Vue {
  @Model('update:show', { type: Boolean, default: false }) show: boolean;
  
  visible: boolean = false;
  
  @Watch('show', { immediate: true })
  onShowChange(val: boolean) {
    this.visible = val;
  }

  /**
   * 取消
   */
  handleCancel() {
    this.$emit('update:show', false);
    this.$emit('cancel');
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
    this.$emit('update:show', false);
    this.$emit('confirm');
  }
}
</script>
<style lang="scss" scoped>


.popup {
  position: relative;
  width: 311px;
  background: #FFFFFF;
  .img  {
    width: 100%;
    display: block;
  }

  .tips {
    padding: 16px 24px;
    font-size: 12px;
    color: #000000;
    text-align: center;
  }

  .popup-footer {
    padding: 12px 24px;

    .van-button {
      height: 48px;
      border-radius: 12px;
      font-size: 12px;
      font-weight: 600;
    }
    .btn {
      // width: 80px;
      height: 32px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 12px;
      font-weight: 600;
      margin: 8px auto 0;
    }
  }
}
</style>