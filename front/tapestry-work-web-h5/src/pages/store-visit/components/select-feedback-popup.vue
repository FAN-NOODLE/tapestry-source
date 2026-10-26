<template>
  <van-popup v-model="visible" position="bottom" round @click-overlay="handleCancel">
    <div class="popup">
      <div class="popup-header">
        <div class="popup-title">{{ $t('storeVisit.popup.title.3') }}</div>
        <div class="popup-close" @click.stop="handleCancel"></div>
      </div>
      <div class="popup-content">
        <div class="list">
          <div class="item" v-for="item,index in list" :key="index" @click="handleClick(item)">
            <div class="item-name">{{ item.label }}</div>
            <div class="item-checkbox"></div>
          </div>
        </div>
      </div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Model, Prop, Vue, Watch } from 'vue-property-decorator';
export const feedbacks = [
  'storeVisit.feedbacks.1',
  'storeVisit.feedbacks.2',
  'storeVisit.feedbacks.3',
  'storeVisit.feedbacks.4',
  'storeVisit.feedbacks.5',
  'storeVisit.feedbacks.6',
  'storeVisit.feedbacks.7',
  'storeVisit.feedbacks.8',
  'storeVisit.feedbacks.9',
]

@Component({

})
export default class SelectFeedbackPopup extends Vue {
  @Model('update:show', { type: Boolean, default: false }) show: boolean;
  
  visible: boolean = false;
  
  list: any[] = [
    { value: '店铺陈列', label: '店铺陈列' },
    { value: '店铺升级', label: '店铺升级' },
    { value: '市场活动', label: '市场活动' },
    { value: '市场拓展', label: '市场拓展' },
    { value: '产品/库存', label: '产品/库存' },
    { value: '销售反馈', label: '销售反馈' },
    { value: '服务技巧', label: '服务技巧' },
    { value: '店铺运营', label: '店铺运营' },
    { value: '系统反馈', label: '系统反馈' },
    { value: '自定义分类', label: '自定义分类' },
  ]
  
  @Watch('show', {  immediate: true })
  onShowChange(val: boolean) {
    this.visible = val;
    if(val) {
      this.list = feedbacks.map(item => {
        return {
          value: this.$t(item),
          label: this.$t(item),
        }
      })
    }
  }

  /**
   * 取消
   */
  handleCancel() {
    this.$emit('update:show', false);
    this.$emit('cancel');
  }

  handleClick(item) {
    this.$emit('update:show', false);
    this.$emit('confirm', item);
  }
}
</script>
<style lang="scss" scoped>


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
    padding-bottom: calc(env(safe-area-inset-bottom));
    padding-bottom: calc(constant(safe-area-inset-bottom));

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
}
</style>