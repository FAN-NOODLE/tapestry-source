<template>
  <van-popup v-model="visible" position="center" @click-overlay="handleCancel">
    <div class="popup">
      <div class="popup-close" @click.stop="handleCancel"></div>
      <div class="popup-title">提示</div>
      <div class="popup-content">
        <div class="text">请输入完整的“商品款_色_码”或“商品款_色”才可以手动添加成功</div>
        <input ref="input" type="text" v-model="inputValue" placeholder="请输入商品款_色_码或商品款_色" class="input" @click="handleClick('input')" :class="{error: isError()}" />
        <div class="error-tip" v-if="isError()">请按要求填写</div>
      </div>
      <div class="popup-footer">
        <van-button color="#000000" plain @click.stop="handleCancel">取消</van-button>
        <van-button color="#000000" @click.stop="handleConfirm">手动添加</van-button>
      </div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Model, Prop, Vue, Watch } from 'vue-property-decorator';

@Component({

})
export default class AddGoodsPopup extends Vue {
  @Model('update:show', { type: Boolean, default: false }) show: boolean;
  // sku
  @Prop({ type: String, default: '' }) sku: string;
  

  visible: boolean = false;
  inputValue: string = '';

  isError() {
    return !/^[0-9A-Za-z]+_[0-9A-Za-z/]+_[0-9A-Za-z/.]+$|^[0-9A-Za-z]+_[0-9A-Za-z/]+$/.test(this.inputValue.trim());
  }

  @Watch('show')
  onShowChange(val: boolean) {
    this.visible = val;
    if(val) {
      this.inputValue = this.sku;
    }
  }

  handleClick(key) {
    (this.$refs[key] as any).focus()
  }


  handleCancel() {
    this.$emit('cancel');
  }

  handleConfirm = this.$helper.debounce(() => {
    
    if(this.isError()) {
      return;
    }
    this.emitConfirm();
  }, 500, true)
  
  emitConfirm() {
    this.$emit('confirm', this.inputValue.trim());
  }
}
</script>
<style lang="scss" scoped>


.popup {
  position: relative;
  width: 325px;
  background: #FFFFFF;
  .popup-close {
    width: 48px;
    height: 48px;
    position: absolute;
    right: 0;
    top: 0;
    background: url(~@/assets/online-order/icon_popup_close.png) center/18px no-repeat;
  }
  .popup-title {
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 500;
    font-size: 16px;
    line-height: 22px;
    padding: 20px;
    color: #000000;
  }

  .popup-content {
    padding: 0 15px;
    .text {
      font-weight: 400;
      font-size: 14px;
      color: #000000;
      line-height: 20px;
      text-align: center;
      padding: 0 10px;
    }

    .input {
      height: 40px;
      box-sizing: border-box;
      border: 1px solid #A6A6A6;
      font-size: 14px;
      width: 100%;
      padding: 0 10px;
      margin-top: 10px;

      &.error {
        border-color: #F64040;
      }
    }

    .error-tip {
      font-size: 10px;
      color: #F64040;
      margin-top: 5px;
    }
  }

  .popup-footer {
    box-sizing: border-box;
    z-index: 2;
    display: flex;
    align-items: center;
    background-color: #fff;
    padding: 20px 15px;

    .van-button {
      flex: 1;
      height: 36px;

      &:nth-child(2) {
        margin-left: 15px;
      }
    }
  }
}
</style>