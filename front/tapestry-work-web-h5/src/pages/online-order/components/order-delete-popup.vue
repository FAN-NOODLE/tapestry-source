<template>
  <van-popup v-model="visible" @click-overlay="hidePopup">
    <div class="popup">
      <div class="popup-close" @click="hidePopup"></div>
      <div class="popup-title">{{title}}</div>
      <div v-if="item" class="popup-content">
        <div class="code">{{item.product_code}}</div>
        <div class="name">{{ item.product_name }}</div>
      </div>
      <div class="popup-footer">
        <van-button color="#F64040" @click="handleConfirm">{{okText}}</van-button>
        <van-button color="#000000" plain @click="hidePopup">取消</van-button>
      </div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Model, Prop, Vue, Watch } from 'vue-property-decorator';

@Component({

})
export default class OrderDeletePopup extends Vue {
  @Model('change', { type: Boolean, default: false }) show: boolean;
  // 商品信息
  @Prop({ type: Object, default: () => (null) }) item: any;
  @Prop({ type: String, default: '确定删除该商品吗？' }) title: string;
  @Prop({ type: String, default: '确认删除' }) okText: string;

  visible: boolean = false;

  @Watch('show')
  onShowChange(val: boolean) {
    this.visible = val;
  }

  hidePopup() {
    this.$emit('change', false);
  }

  handleConfirm = this.$helper.debounce(() => {
    this.hidePopup()
    this.$emit('confirm');
  }, 500, true)
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
  .popup-title {
    height: 52px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 500;
    font-size: 16px;
    color: #000000;
  }

  .popup-content {
    text-align: center;
    .code {
      font-weight: 400;
      font-size: 14px;
      color: #000000;
      line-height: 20px;
      margin-top: 5px;
    }

    .name {
      font-weight: 400;
      font-size: 12px;
      color: #96979E;
      line-height: 17px;
      margin-top: 5px;
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