<template>
  <van-popup v-model="visible" @click-overlay="handleContinue">
    <div class="popup">
      <div class="popup-close" @click="handleContinue"></div>
      <div class="popup-content">
        <van-image v-if="showIcon" width="50" height="50" :src="require('@/assets/online-order/icon_success.png')"></van-image>
        <div class="title">{{title}}</div>
        <div class="text" v-for="(item, index) in content" :key="index">{{item}}</div>
      </div>
      <div v-if="showFooter" class="popup-footer">
        <van-button color="#000000" plain @click="handleContinue">{{ cancelText }}</van-button>
        <van-button color="#000000" @click="handleConfirm">发送至客户</van-button>
      </div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Model, Prop, Vue, Watch } from 'vue-property-decorator';

@Component({
  model: {
    prop: 'show',
    event: 'change'
  }
})
export default class OrderSuccessPopup extends Vue {
  @Prop({ type: Boolean, default: false }) show: boolean;
  @Prop({ type: Boolean, default: true }) showFooter: boolean;
  @Prop({ type: Boolean, default: true }) showIcon: boolean;
  @Prop({ type: String, default: '提交成功' }) title: string;
  @Prop({ type: String, default: '继续创建' }) cancelText: string;
  @Prop({ type: Array, default: () => ['订单已生成，请发给客户填写收货信息'] }) content: string[];

  visible: boolean = false;

  @Watch('show')
  onShowChange(val: boolean) {
    this.visible = val;
  }

  hidePopup() {
    this.$emit('change', false);
  }

  handleContinue() {
    this.hidePopup()
    this.$emit('close');
  }

  handleConfirm() {
    this.hidePopup()
    this.$emit('confirm');
  }
}
</script>
<style lang="scss" scoped>
.popup {
  width: 325px;
  background: #FFFFFF;

  position: relative;
  .popup-close {
    width: 48px;
    height: 48px;
    position: absolute;
    right: 0;
    top: 0;
    background: url(~@/assets/online-order/icon_popup_close.png) center/18px no-repeat;
  }
  
  .popup-content {
    text-align: center;
    padding: 40px 0;
    overflow: hidden;
    .van-image {
      display: block;
      margin: auto;
    }
    .title {
      width: 160px;
      margin: 15px auto;
      font-weight: 500;
      font-size: 16px;
      color: #000000;
      line-height: 22px;
    }
    .text {
      font-weight: 400;
      font-size: 14px;
      color: #000000;
      line-height: 20px;
    }
  }

  .popup-footer {
    padding: 20px 15px;
    display: flex;
    align-items: center;

    .van-button {
      flex: 1;
      height: 36px;
      &:nth-child(2) {
        margin-left: 10px;
      }
    }
  }
}
</style>