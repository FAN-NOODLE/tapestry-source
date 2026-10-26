<template>
  <van-popup v-model="visible" position="bottom" round @click-overlay="handleCancel">
    <div class="popup">
      <div class="popup-header">
        <div class="popup-title">同步给指定店员</div>
        <div class="popup-close" @click.stop="handleCancel"></div>
      </div>
      <div class="popup-content">
        <div class="list">
          <div class="item" v-for="item,index in list" :key="index" :class="{active: item.checked}" @click="handleClick(item)">
            <div class="item-name">{{ item.name }}</div>
            <div class="item-checkbox"></div>
          </div>
        </div>
      </div>
      <div class="popup-footer">
        <van-button color="#000000" plain @click.stop="handleCancel">取消</van-button>
        <van-button color="#000000" @click.stop="handleConfirm">同步</van-button>
      </div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Model, Prop, Vue, Watch } from 'vue-property-decorator';

@Component({

})
export default class SelectClerkPopup extends Vue {
  @Model('update:show', { type: Boolean, default: false }) show: boolean;
  
  visible: boolean = false;
  
  list: any[] = [
    { id: '1', name: '张三', checked: false },
    { id: '2', name: '张三', checked: false },
    { id: '3', name: '张三', checked: false },
    { id: '4', name: '张三', checked: false },
    { id: '5', name: '张三', checked: false },
  ]
  
  @Watch('show')
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
   * 确认
   */
  handleConfirm = this.$helper.debounce(() => {
    // 触发确认事件

    this.emitConfirm();
  }, 500, true)
  
  /**
   * 确认事件
   */
  emitConfirm() {
    this.$emit('update:show', false);
    this.$emit('confirm');
  }

  handleClick(item) {
    item.checked = !item.checked;
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