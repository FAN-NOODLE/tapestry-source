<template>
  <van-popup v-model="visible" overlay-class="overlay" position="center" :close-on-click-overlay="false">
    <div class="popup">
      <div class="title">{{ $t('storeVisit.title.4') }}</div>
      <div class="add">
        <div class="icon-plus"></div>
        <span>{{ $t('storeVisit.btn.4') }}</span>
      </div>
      <div class="tips">
        <div class="icon-feedback-white"></div>
        <div class="text">{{ $t('storeVisit.tips.5') }}</div>
      </div>
      <div class="btn" @click="handleConfirm">{{ $t('storeVisit.btn.8') }}</div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Model, Prop, Vue, Watch } from 'vue-property-decorator';

@Component({

})
export default class GuidePopup extends Vue {
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
.icon-feedback {
  width: 20px;
  height: 20px;
  background: url("~@/assets/store-visit/icon_feedback.png") center/100% no-repeat;
}
.icon-feedback-white {
  width: 32px;
  height: 32px;
  background: url("~@/assets/store-visit/icon_feedback_white.png") center/100% no-repeat;
}
.icon-plus {
  width: 20px;
  height: 20px;
  background: url("~@/assets/store-visit/icon_plus.png") center/20px no-repeat;
}
.van-popup {
  background-color: transparent;
}

.popup {
  width: 100vw;
  height: 100vh;
  .title {
    padding-top: 112px;
    text-align: center;
    font-size: 16px;
    line-height: 24px;
    color: #fff;
    font-weight: 600;
  }

  .add {
    width: 343px;
    height: 52px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-size: 14px;
    font-weight: 600;
    box-sizing: border-box;
    margin: 48px auto 0;
    padding: 0 16px;
    background-color: rgba(255, 255, 255, 1);
    border-radius: 12px;
  }

  .tips {
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
    color: #fff;
    margin-top: 44px;
    .icon-feedback-white {
      position: relative;

      &::after {
        content: '';
        width: 15px;
        height: 20px;
        position: absolute;
        top: -20px;
        left: -15px;
        background: url("~@/assets/store-visit/icon_x.png") center/100% auto no-repeat;
      }
    }
  }

  .btn {
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 12px;
    width: 104px;
    height: 48px;
    margin: 125px auto 0;
    color: #fff;
    font-size: 12px;
    font-weight: 600;
    background-color: rgba(13, 110, 253, 1);
  }
}
</style>