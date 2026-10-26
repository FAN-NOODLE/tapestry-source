<template>
  <!--分组弹窗 -->
  <van-popup
    v-model="visible"
    :close-on-click-outside="false"
    :close-on-click-action="false"
    position="center"
  >
    <div class="popup">
      <div class="popup-close" @click="handleClose"></div>
      <div class="popup-header">
        <div class="title">编辑</div>
      </div>
      <div class="popup-content">
        <textarea
          v-model="content"
          class="input"
          ref="fileNameInput"
          rows="5"
        ></textarea>
      </div>
      <div class="popup-footer">
        <div class="button" @click="handleSend">发送</div>
      </div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Vue,PropSync, Prop,Watch } from "vue-property-decorator";
import { Popup, Field } from "vant";
@Component({
  components: {
    [Popup.name]: Popup,
    [Field.name]: Field,
  },
})
export default class PreviewVideoPopup extends Vue {
  @PropSync("show", { type: Boolean }) visible;
  @Prop({ type: Object, default: () => ({}) }) item: any;

  content = ''

  @Watch("item")
  onItemChange(val: any) {
    this.content = val.content;
  }

  @Watch("visible")
  onShowChange(val:Boolean){
    if(val){
      this.$nextTick(() => {
        (this.$refs.fileNameInput as HTMLElement)?.focus();
      });
    }
  }
  handleClose() {
    this.visible = false
    this.$emit("close");
  }
  /**
   * @description: 确定
   * @return {*}
   */
  handleSend() {
    
    this.$emit("send", {
      ...this.item,
      content: this.content
    });
  }
}
</script>
<style lang="scss" scoped>
::v-deep .van-field__control {
  font-size: 14px;
}
.popup {
  background-color: #fff;
  position: relative;
  width: 280px;
  .popup-close {
    width: 40px;
    height: 40px;
    background: url("https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/icon_close.png") no-repeat center;
    background-size: 20px 20px;
    position: absolute;
    top: 0;
    right: 0;
  }

  .popup-header {
    font-size: 16px;
    font-weight: bold;
    text-align: center;
    height: 40px;
    line-height: 40px;
  }

  .popup-content {
    padding: 0 20px;
    font-size: 14px;
    height: 100px;
    .input {
      width: 100%;
      height: 100%;
    }
  }

  .popup-footer {
    display: flex;
    justify-content: center;
    padding: 15px;
    .button {
      width: 100%;
      height: 34px;
      font-size: 14px;
      background-color: #4a4a4a;
      color: #fff;
      display: flex;
      align-items: center;
      justify-content: center;
    }
  }
}
</style>
