

<template>
  <van-popup v-model="visible" position="bottom" round @click-overlay="handleCancel">
    <div class="popup">
      <div class="popup-content">
        <div class="item" @click="handleClick(1)">{{ $t('storeVisit.popup.btn.17') }}</div>
        <div class="item" @click="handleClick(2)">{{ $t('storeVisit.popup.btn.18') }}</div>
      </div>
      <div class="popup-footer">
        <div v-if="!disabledDelete" class="item" @click="handleClick(3)">{{ $t('storeVisit.popup.btn.13') }}</div>
      </div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Model, Prop, Vue, Watch } from 'vue-property-decorator';
import ComUpload from "./com-upload.vue";

@Component({
  components: {
    ComUpload,
  },
})
export default class SheetmenuPopup extends Vue {
  @Model('update:show', { type: Boolean, default: false }) show: boolean;
  @Prop({ type: Boolean, default: false }) disabledDelete: boolean;

  visible: boolean = false;

  @Watch('show', { immediate: true })
  onShowChange(val: boolean) {
    this.visible = val;
  }

  /**
   * 取消
   */
  handleCancel() {
    this.$emit('update:show', false);
    this.$emit('cancel');
  }

  /**
   * 确认事件
   */
  handleClick(value) {
    this.$emit('update:show', false);
    this.$emit('confirm', value);
  }
}
</script>
<style lang="scss" scoped>
.popup {
  width: 100vw;
  position: relative;
  background: #FFFFFF;
  
  .item {
    height: 68px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
    font-weight: 600;
    color: rgba(0, 0, 0, 1);
    &:not(:last-child) {
      border-bottom: 1px solid rgba(248, 249, 250, 1);
    }
  }

  .popup-content {
    padding: 0 24px;
  }

  .popup-footer {
    border-top: 4px solid rgba(248, 249, 250, 1);
    padding-bottom: calc(env(safe-area-inset-bottom) + 12px);
    padding-bottom: calc(constant(safe-area-inset-bottom) + 12px);


    .item {
      font-weight: normal;
    }
  }
}
</style>