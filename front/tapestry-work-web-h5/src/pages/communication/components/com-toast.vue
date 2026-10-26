<template>
  <van-popup 
    v-model="visible"
    @click-overlay="handleCancel">
    <div class="com-toast">
      <div class="com-toast__close" @click="handleCancel"></div>
      <div class="com-toast__content">
        <div class="icon"></div>
        <div class="title">提交成功</div>
        <div class="text">下载内容稍后将发送到指定邮箱</div>
      </div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Vue, Prop, Watch } from 'vue-property-decorator'
import { Popup } from "vant";

@Component({
  components: {
    [Popup.name]: Popup
  },
  model: {
    prop: 'visible',
    event: 'change'
  }
})
export default class ComToast extends Vue {
  @Prop({ type: Boolean, default: false }) visible: boolean;
  
  @Watch('visible')
  handleVisibleChange(val: boolean) {
    if (val) {
      setTimeout(() => {
        this.handleCancel();
      }, 2000);
    }
  }

  handleCancel() {
    this.$emit('change', false);
  }
}
</script>

<style lang="scss" scoped> 
.com-toast {
  width: 325px;
  padding: 40px;
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

  &__content {
    text-align: center;
    color: #000000;

    .icon {
      margin: auto;
      width: 50px;
      height: 50px;
      background: url('~@/assets/communication/icon_success.png') center/50px no-repeat;
    }
    .title {
      font-weight: 500;
      font-size: 16px;
      line-height: 22px;
      margin-top: 15px
    }

    .text {
      margin-top: 15px;
      font-weight: 400;
      font-size: 14px;
      line-height: 20px;
    }
  }

}
</style>