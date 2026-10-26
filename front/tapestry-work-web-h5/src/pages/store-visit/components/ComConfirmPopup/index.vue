<template>
  <van-popup v-model="visible" position="center" round @click-overlay="handleClose">
    <div class="popup">
      <div class="popup-header">
        <div class="popup-title">
          {{ title }}
          <com-popover v-if="tips" placement="bottom-start" theme="light" :offsetX="-50" top="0">
            <template #reference>
              <div class="popup-tip"></div>
            </template>
            <div class="tip-content-2">{{ tips }}</div>
          </com-popover>

        </div>
        <div class="popup-close" @click.stop="handleClose"></div>
      </div>
      <div class="popup-content">
        <slot>
          <div class="popup-content-text" v-html="content"></div>
        </slot>
      </div>
      <div v-if="footerStyle=='2'" class="popup-footer-2">
        <van-button color="#000000" block @click.stop="handleConfirm">{{confirmButtonText}}</van-button>
        <van-button color="#000000" block plain @click.stop="handleCancel">{{ cancelButtonText }}</van-button>
      </div>
      <div v-else class="popup-footer">
        <van-button color="#000000" plain @click.stop="handleCancel">{{ cancelButtonText }}</van-button>
        <van-button color="#000000" @click.stop="handleConfirm">{{confirmButtonText}}</van-button>
      </div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Model, Prop, Vue, Watch } from 'vue-property-decorator';
import ComPopover from '../com-popover.vue';

@Component({
  components: {
    ComPopover,
  },
})
export default class ComConfirmPopup extends Vue {
  @Model('update:show', { type: Boolean, default: false }) show: boolean;
  @Prop({ type: String, default: '温馨提示' }) title: string;
  @Prop({ type: String, default: '发起新的巡店' }) cancelButtonText: string;
  @Prop({ type: String, default: '使用该巡店计划继续' }) confirmButtonText: string;
  @Prop({ type: String, default: '所选门店已有一场进行中的巡店计划！' }) content: string;
  @Prop({ type: String, default: '' }) footerStyle: string;
  @Prop({ type: String, default: '' }) tips: string;
  
  visible: boolean = false;
  

  current: number = -1;
  
  @Watch('show')
  onShowChange(val: boolean) {
    this.visible = val;
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
    this.$emit('update:show', false);
    this.$emit('confirm', {
      action: 'confirm'
    });
  }
}
</script>
<style lang="scss" scoped>


.popup {
  position: relative;
  width: 311px;
  background: #FFFFFF;
  .popup-header {
    padding: 20px 24px;
    position: relative;
    
    .popup-title {
      font-weight: 600;
      font-size: 16px;
      line-height: 24px;
      color: #000000;
      display: flex;
      align-items: center;
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
</style>
<style lang="scss"> 

.popup-tip {
  width: 24px;
  height: 24px;
  background: url(~@/assets/store-visit/icon_wen.png) center/16px no-repeat;
}

.tip-content-2 {
  font-size: 12px;
  font-weight: 400;
  width: 180px;
}
</style>