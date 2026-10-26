<template>
  <van-popup v-model="show" class="popup" get-container="body">
    <div class="popup-title">{{title}}</div>
    <div class="popup-close" @click="onHide"></div>
    <div class="popup-content">
      <slot></slot>
    </div>
    <div class="popup-footer">
      <div class="btn plain" @click="onHide">取消</div>
      <div class="btn" @click="onConfirm">确定</div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
import { Popup } from 'vant';
@Component({
  components: {
    [Popup.name]: Popup
  }
})
export default class ConfirmPopup extends Vue {
  @Prop({ type: String }) title: string
  show = false;
  checked = false;

  showPopup() {
    this.show = true
  }
  onHide() {
    this.show = false
    this.$emit("cancel")
  }
  onConfirm() {
    this.$emit("confirm")
    this.onHide()
  }
}
</script>
<style lang="scss" scoped>
.popup {
  width: 325px;
  background: #FFFFFF;
  padding: 20px 15px;
  box-sizing: border-box;

  &-title {
    text-align: center;
    font-size: 16px;
    font-weight: 500;
    color: #000000;
    line-height: 22px;
  }

  &-close {
    position: absolute;
    width: 18px;
    height: 18px;
    right: 15px;
    top: 15px;
    background: url("~@/assets/display-task/icon_popup_close.png") center/100% no-repeat;
  }

  &-content {
    margin-top: 20px;
    text-align: center;
    font-size: 14px;
  }

  &-footer {
    display: flex;
    justify-content: center;
    margin-top: 35px;
    
   .btn {
      flex: 1;
      height: 36px;
      background: #000000;
      color: #ffffff;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 14px;
      &.plain {
        color: #000;
        background: #fff;
        border: 1px solid #000;
        margin-right: 15px;
      }
    }
  } 
}
</style>