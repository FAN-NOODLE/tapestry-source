<template>
  <van-popup v-model="visible" position="bottom" round @click-overlay="handleClose">
    <div class="popup">
      <div class="popup-header">
        <div class="popup-title">{{ $t(title) }}</div>
        <div class="popup-close" @click.stop="handleClose"></div>
      </div>
      <div class="popup-content">
        <div class="tips">{{ $t(tips) }}：</div>
        <div class="list">
          <div class="item" v-for="item,index in list" :key="index">
            <div v-if="item.name && item.user_id" class="item-info" @click="handleOpenChat(item)">
              <div class="item-name">{{ item.name }}</div>
              <div class="item-dept">{{ item.dept_name }}</div>
              <div class="item-btn">{{ $t('storeVisit.popup.tips.8') }}</div>
              <div class="icon-right"></div>
            </div>
            <div v-if="item.email" class="item-email">
              <div class="icon-email"></div>
              <div class="item-email-text">{{ item.email }}</div>
            </div>
          </div>
        </div>
      </div>
      <div class="popup-footer">
        <van-button color="#000000" plain @click.stop="handleCancel">{{$t(cancelText)}}</van-button>
        <van-button v-if="showConfirmButton" color="#000000" @click.stop="handleConfirm">{{ $t('storeVisit.popup.btn.8') }}</van-button>
      </div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Model, Prop, Vue, Watch } from 'vue-property-decorator';
import wxUtil from  '@/utils/wx';


@Component({

})
export default class ContactPopup extends Vue {
  @Model('update:show', { type: Boolean, default: false }) show: boolean;
  @Prop({ type: String, default: 'storeVisit.popup.title.4' }) title: string;
  @Prop({ type: String, default: 'storeVisit.popup.tips.9' }) tips: string;
  @Prop({ type: String, default: 'storeVisit.popup.btn.9' }) cancelText: string;
  @Prop({ type: Boolean, default: false }) showConfirmButton: boolean;
  @Prop({ type: Array, default: () =>[] }) list: any[];
  
  visible: boolean = false;
  
  @Watch('show')
  onShowChange(val: boolean) {
    this.visible = val;
  }

  /**
   * 取消
   */
  handleClose() {
    this.$emit('update:show', false);
    this.$emit('close');
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
  }, 0, true)
  
  /**
   * 确认事件
   */
  emitConfirm() {
    this.$emit('update:show', false);
    this.$emit('confirm');
  }

  handleOpenChat(item) {
    console.log(item);

    this.$emit('contact', {
      track: {
        staff_name: item.name,
        staff_id: item.user_id,
        staff_department: item.dept_name,
        role_type: this.tips.replace('：', '')
      }
    });
    
    wxUtil.openEnterpriseChat({
      userIds: item.user_id
    })
    
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

  .popup-tips {
    padding: 12px 24px;
    font-size: 12px;
    line-height: 130%;
  }

  .popup-content {
    max-height: 450px;
    padding: 24px 0;
    overflow-y: auto;
    .tips {
      color: rgba(153, 153, 153, 1);
      font-size: 12px;
      font-weight: 600;
      margin-bottom: 12px;
      padding: 0 24px;
    }

    .list {
      padding: 0 24px;

      .item {
        // height: 100px;
        box-sizing: border-box;
        padding: 16px;
        transition: all 0.3s;
        background-color: rgba(248, 249, 250, 1);
        color: rgba(0, 0, 0, 1);
        border-radius: 12px;
        position: relative;

        &:not(:last-child) {
          margin-bottom: 8px;
        }

        .item-info {
          display: flex;
          align-items: center;
        }

        .item-name { 
          font-size: 14px;
          font-weight: 700;
        }

        .item-dept {
          font-size: 12px;
          font-weight: 500;
          color: rgba(153, 153, 153, 1);
          background-color: rgba(255, 255, 255, 1);
          padding: 0 4px;
          margin-left: 8px;
        }

        .item-btn {
          flex: 1;
          text-align: right;
          font-size: 10px;
          font-weight: 600;
        }

        .icon-right {
          width: 16px;
          height: 16px;
          background: url(~@/assets/store-visit/icon_right_black.png) center/16px no-repeat;
          margin-left: 8px;
        }

        .item-email {
          display: flex;
          align-items: center;
          background-color: rgba(255, 255, 255, 0.8);
          padding: 4px;
          margin-top: 12px;

          .icon-email {
            width: 16px;
            height: 16px;
            background: url(~@/assets/store-visit/icon_email.png) center/16px no-repeat;
          }

          .item-email-text {
            font-size: 12px;
            font-weight: 400;
            color: rgba(0, 0, 0, 0.8);
            margin-left: 8px;

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
    border-top: 1px solid rgba(233, 236, 239, 1);

    .van-button {
      height: 40px;
      border-radius: 4px;
      flex: 1;

      &.van-button--disabled {
        opacity: 1;
        background: rgba(222, 226, 230, 1) !important;
        border-color: rgba(222, 226, 230, 1) !important;
      }

      &:nth-child(2) {
        margin-left: 8px;
      }
    }
  }
}
</style>