<template>
  <div class="select">
    <div class="select" @click="handleClick">
      <slot>{{ innerValue }}</slot>
    </div>
    <van-popup v-model="visible" position="bottom" :get-container="'body'" round @click-overlay="handleCancel">
      <div class="popup">
        <van-picker
          v-model="innerValue"
          :title="String($t(title))"
          show-toolbar
          :columns="columns()"
          @cancel="handleCancel"
          @confirm="handleConfirm"
        />
      </div>
    </van-popup>
  </div>
</template>
<script lang="ts">
import { Component, Model, Prop, Vue, Watch } from 'vue-property-decorator';
import { Picker } from "vant";

@Component({
  components: {
    [Picker.name]: Picker,
  }
})
export default class SelectTimePopup extends Vue {
  @Model('update:value', { type: [Number,String], default: false }) value: number|string;
  @Prop({ type: String, default: 'storeVisit.popup.title.12' }) title: string;
  @Prop({ type: Boolean, default: false }) type: boolean; // true:具体时间 false:上下午
  
  visible: boolean = false;
  innerValue: string = '';

  @Watch('value', { immediate: true })
  onValueChange(val) {
    if (!val) {
      this.innerValue = '';
      return;
    }
    this.innerValue = val;
  }

  get columns() {
    return () => {
      if(this.type) {
        // i.勾选具体时间时，10-18点每两小时一个时间区间；
        return [
          '10:00 - 12:00',
          '12:00 - 14:00',
          '14:00 - 16:00',
          '16:00 - 18:00',
        ]
      }
      // ii.未勾选具体时间时，上午/下午；
      console.log();
      
      return [
        this.$t('storeVisit.popup.tips.29'),
        this.$t('storeVisit.popup.tips.30'),
      ]
    }
  }

  handleClick() {
    this.visible = true;
  }
  /**
   * 取消
   */
  handleCancel() {
    this.visible = false;
    this.$emit('cancel');
  }

  /**
   * 确定
   */
  handleConfirm(value) {
    this.visible = false;
    this.$emit('update:value', value);
    this.$emit('confirm', value);
  }

}
</script>
<style lang="scss" scoped>
.select {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.popup {
  position: relative;
  width: 100vw;
  background: #FFFFFF;
  .popup-header {
    padding: 20px 24px;
    position: relative;
    
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
    background: rgba(248, 249, 250, 1);
    max-height: 450px;
    overflow-y: auto;

    .list {
      padding: 0 24px;

      .item {
        height: 64px;
        display: flex;
        align-items: center;

        &:not(:last-child) {
          border-bottom: 0.5PX solid rgba(233, 236, 239, 1);
        }

        .item-name {
          font-weight: 400;
          font-size: 14px;
          color: #000000;
          flex: 1;
          transition: all 0.3s;
        }

        .item-checkbox {
          width: 24px;
          height: 24px;
          background: url(~@/assets/store-visit/icon_check.png) center/24px no-repeat;
          opacity: 0;
          transition: all 0.3s;
        }

        &.active {
          .item-name {
            font-weight: 700;
          }
          .item-checkbox {
            opacity: 1;
          }
        }
      }
    }
  }

  .popup-footer {
    box-sizing: border-box;
    z-index: 2;
    display: flex;
    align-items: center;
    background-color: #fff;
    padding: 12px 24px;
    padding-bottom: calc(env(safe-area-inset-bottom) + 12px);
    padding-bottom: calc(constant(safe-area-inset-bottom) + 12px);

    .van-button {
      flex: 1;
      height: 40px;
      border-radius: 4px;

      &.van-button--normal {
        font-size: 12px;
        padding: 0 16px;
      }

      &:nth-child(2) {
        margin-left: 8px;
      }
    }
  }
}
</style>