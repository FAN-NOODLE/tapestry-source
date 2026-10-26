<template>
  <van-popup v-model="visible" position="bottom" round :close-on-click-overlay="false">
    <div class="popup">
      <div class="popup-header">
        <div class="icon-left" @click="handleBack"></div>
        <div class="popup-title" @click="handleBack">{{$t(isEdit?'storeVisit.popup.title.10':'storeVisit.popup.title.9')}}</div>
        <div class="popup-close" @click.stop="handleCancel"></div>
      </div>
      <div class="popup-content">
        <input v-model="innerValue" ref="input" class="input" type="text" autofocus :placeholder="String($t('storeVisit.placeholder.6'))" @click="handleInputClick" />
      </div>
      <div class="popup-footer">
        <div v-if="!isEdit" class="checkbox" @click="checked=!checked">
          <div class="icon-check" :class="{active: checked}"></div>
          <div class="label">{{ $t('storeVisit.popup.tips.13') }}</div>
        </div>
        <div v-else class="checkbox"></div>
        <div class="add" @click="handleConfirm">{{$t(isEdit?'storeVisit.popup.btn.15':'storeVisit.popup.btn.16')}}</div>
      </div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Model, Prop, Vue, Watch } from 'vue-property-decorator';

@Component({

})
export default class AddEmailPopup extends Vue {
  @Model('update:show', { type: Boolean, default: false }) show: boolean;
  @Prop({ type: String, default: '' }) email: string;
  @Prop({ type: Boolean, default: false }) isEdit: boolean;
  visible: boolean = false;
  checked: boolean = false;

  innerValue: string = '';
  @Watch('email', { immediate: true })
  onValueChange(val: string) {
    this.innerValue = val || '';
  }
  
  @Watch('show')
  onShowChange(val: boolean) {
    this.visible = val;
    if(val) {
      this.$nextTick(() => {
        (this.$refs.input as any).focus();
      });
      this.checked = this.isEdit;
    }
  }
  get isEmail() {
    return (value) => /^[a-zA-Z0-9_.-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$/.test(value);
  }
  /**
   * 点击获取焦点
   */
  handleInputClick() {
    (this.$refs.input as any).focus()
  }

  /**
   * 取消
   */
  handleCancel() {
    this.$emit('update:show', false);
    this.$emit('cancel');
  }

  /**
   * 返回
   */
  handleBack() {
    this.$emit('update:show', false);
    this.$emit('back');
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
    // 正则校验邮箱格式
    if (!this.innerValue) {
      this.$toast(this.$t('common.toast.18'));
      return;
    }
    if (!this.isEmail(this.innerValue)) {
      this.$toast(this.$t('common.toast.19'));
      return;
    }
    this.$emit('update:show', false);
    this.$emit('confirm', {
      email: this.innerValue,
      checked: this.checked,
      isEdit: this.isEdit
    });
  }
}
</script>
<style lang="scss" scoped>


.popup {
  position: relative;
  width: 100vw;
  background: #FFFFFF;
  .popup-header {
    padding: 20px 24px;
    position: relative;
    display: flex;
    align-items: center;

    .icon-left {
      width: 24px;
      height: 24px;
      background: url(~@/assets/store-visit/icon_left.png) center/24px no-repeat;
    }
    
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
    padding: 8px 16px;

    .input {
      width: 100%;
      height: 40px;
      background-color: rgba(233, 236, 239, 1);
      border-radius: 12px;
      padding: 0 16px;
      box-sizing: border-box;
      font-size: 12px;

      &::placeholder {
        color: rgba(0, 0, 0, 0.2);
      }
    }
  }

  .popup-footer {
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 20px;

    .checkbox {
      height: 40px;
      display: flex;
      align-items: center;
      .icon-check {
        width: 16px;
        height: 16px;
        background: url(~@/assets/store-visit/icon_check_0.png) center/16px no-repeat;

        &.active {
          background: url(~@/assets/store-visit/icon_check_1.png) center/16px no-repeat;
        }
      }

      .label {
        font-size: 12px;
        color: rgba(0, 0, 0, 1);
        margin-left: 8px;
        margin-left: 4px;
      }
    }

    .add {
      height: 40px;
      padding: 0;
      color: rgba(13, 110, 253, 1);
      font-size: 12px;
      line-height: 40px;
      font-weight: 600;
    }
  }
}
</style>