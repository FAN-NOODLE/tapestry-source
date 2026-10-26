<template>
  <!--分组弹窗 -->
  <van-popup
    v-model="visible"
    :close-on-click-outside="false"
    :close-on-click-action="false"
    position="bottom"
  >
    <div class="popup">
      <div class="popup-header">
        <div class="btn cancel" @click="handleClose">取消</div>
        <div class="btn" @click="handleConfirm">确定</div>
      </div>
      <div class="popup-content">
        <div :class="['item', {active: current === index}]" v-for="(item,index) in list" :key="index" @click="handleChange(index)">{{item.title}}</div>
      </div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Vue, PropSync, Prop,Watch } from "vue-property-decorator";
import { Popup, Field } from "vant";
@Component({
  components: {
    [Popup.name]: Popup,
    [Field.name]: Field,
  },
})
export default class SelectSearchPopup extends Vue {
  @PropSync("show", { type: Boolean }) visible;
  @Prop({ type: Array, default: () => [] }) list: any[];
  @Prop({ type: Number, default: 0}) active: number;

  current = 0

  @Watch('active')
  handleActiveChange(val) {
    this.current = val;
  }

  handleChange(val) {
    this.current = val
  }

  handleClose() {
    this.visible = false
    this.$emit("close");
  }
  /**
   * @description: 确定
   * @return {*}
   */
  handleConfirm() {
    this.visible = false
    this.$emit("confirm", this.current);
  }
}
</script>
<style lang="scss" scoped>
.popup {
  background-color: #fff;
  position: relative;
  padding-bottom: 20px;
  
  .popup-header {
    font-size: 14px;
    font-weight: bold;
    text-align: center;
    height: 50px;
    display: flex;
    padding: 0 15px;
    align-items: center;
    justify-content: space-between;
    position: relative;

    &::after {
      content: "";
      position: absolute;
      bottom: 0;
      left: 50%;
      transform: translateX(-50%);
      height: 1px;
      width: 345px;
      background-color: #F0F0F0;
    }

    .cancel {
      color: #A6A6A6;
    }

  }

  .popup-content {
    padding: 20px 15px;

    .item {
      width: 345px;
      height: 42px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 14px;
      margin-bottom: 10px;

      &.active {
        font-weight: 500;
        background: #F6F5F5;
      }
    }
  }
}
</style>
