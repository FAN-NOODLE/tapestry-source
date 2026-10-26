<template>
  <van-popup 
    v-model="visible"
    @click-overlay="handleCancel">
    <div class="com-message">
      <div class="com-message__close" @click="handleCancel"></div>
      <div class="com-message__header">{{ title }}</div>
      <div class="com-message__content">
        {{ content }}
      </div>
      <div class="com-message__footer">
        <div class="btn plain" @click="handleCancel">取消</div>
        <div class="btn" @click="handleConfirm">确定</div>
      </div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Vue, Prop } from 'vue-property-decorator';
import { Popup } from "vant";
@Component({
  components: {
    [Popup.name]: Popup,
  },
  model: {
    prop: 'visible',
    event: 'change'
  }
})
export default class ComMessage extends Vue {
  @Prop({ type: Boolean, default: false }) visible: boolean;
  @Prop({ type: String, default: '' }) title: string;
  @Prop({ type: String, default: '' }) content: string;

  handleCancel() {
    this.$emit('change', false);
    this.$emit('cancel');
  }

  handleConfirm() {
    this.$emit('confirm');
  }
}
</script>
<style lang="scss" scoped>
.com-message {
  width: 325px;
  padding: 20px;
  box-sizing: border-box;
  background: #FFFFFF;
  position: relative;

  &__close {
    position: absolute;
    top: 0;
    right: 0;
    width: 48px;
    height: 48px;
    background: url('~@/assets/communication/icon_popup_close.png') center/18px no-repeat;
  }

  &__header {
    height: 22px;
    font-weight: 500;
    font-size: 16px;
    color: #000000;
    line-height: 22px;
    text-align: center;
  }

  &__content {
    width: 200px;
    margin: 30px auto;
    text-align: center;
    font-weight: 400;
    font-size: 14px;
    color: #000000;
    line-height: 20px;
  }

  &__footer {
    display: flex;
    align-items: center;

    .btn {
      flex: 1;
      height: 36px;
      background: #000000;
      color: #fff;
      font-weight: 400;
      font-size: 14px;
      display: flex;
      align-items: center;
      justify-content: center;
      box-sizing: border-box;

      &.plain {
        background: #fff;
        color: #000;
        border: 1px solid #000;
      }

      &:nth-child(2) {
        margin-left: 10px;
      }
    }
  }
}
</style>