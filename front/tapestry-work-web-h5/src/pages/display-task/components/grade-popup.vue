<template>
  <van-popup v-model="show" class="popup" position="bottom" get-container="body">
    <div class="popup-close" @click="handleCancel"></div>
    <div class="popup-header">请对该门店陈列评分</div>
    <div class="popup-content">
      <div class="input">
        <input id="input" ref="input" v-if="show" class="input__inner" v-model="content" placeholder="请输入0~100之间的整数"/>
      </div>
      <div class="btn" @click="handleConfirm">确定</div>
      <div v-if="showTips" class="error-tips">请输入0~100之间的整数</div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Vue, Prop, Watch } from 'vue-property-decorator'
import { Popup } from 'vant';
@Component({
  components: {
    [Popup.name]: Popup
  }
})
export default class CommentPopup extends Vue {
  show: boolean = false
  content: string = ''

  get showTips() {
    return this.content && !/^(0|[1-9][0-9]{0,1}|100)$/.test(this.content)
  }

  @Watch('show')
  handleShowChange(val) {
    if (val) {
      this.$nextTick(() => {
        setTimeout(() => {
          (this.$refs.input as any).focus();
          (this.$refs.input as any).scrollIntoView({ behavior: 'smooth', block: 'center' });
        }, 100);
      });
    }
  }

  showPopup(v) {
    this.show = true
    this.content = v || ''
  }

  hidePopup() {
    this.show = false
    this.content = ''
  }
  // handleChange() {

  // }

  handleConfirm() {
    if(!/^(0|[1-9][0-9]{0,1}|100)$/.test(this.content)) {
      this.$toast('请输入0~100之间的整数')
      return
    }
    this.$emit('confirm', Number(this.content))
    this.hidePopup()
  }

  handleCancel() {
    this.$emit('cancel')
    this.hidePopup()
  }
}
</script>
<style lang="scss" scoped>
.popup {
  padding: 20px 15px;
  // padding-bottom: calc(20px + constant(safe-area-inset-bottom)); /** iOS 11.0+ */
  // padding-bottom: calc(20px + env(safe-area-inset-bottom)); /** iOS 11.0+ */

  &-close {
    position: absolute;
    width: 18px;
    height: 18px;
    right: 15px;
    top: 15px;
    background: url("~@/assets/display-task/icon_popup_close.png") center/100% no-repeat;
  }

  &-header {
    text-align: center;
    font-size: 16px;
    font-weight: 500;
    color: #000000;
    line-height: 22px;
  }

  &-content {
    margin-top: 20px;
    display: flex;
    align-items: center;
    position: relative;

    .input {
      flex: 1;
      height: 47px;
      background: #F6F5F5;
      display: flex;
      align-items: center;
      justify-content: center;

      &__inner {
        font-size: 14px;
        padding-left: 10px;
        background-color: transparent;
        width: 100%;
        height: 100%;
      }
    }

    .btn {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 103px;
      height: 47px;
      background: #000000;
      font-size: 14px;
      font-weight: 400;
      color: #FFFFFF;
      line-height: 20px;
    }

    .error-tips {
      position: absolute;
      bottom: -16px;
      left: 0;
      font-size: 10px;
      color: #E02020;
    }
  }
}
</style>