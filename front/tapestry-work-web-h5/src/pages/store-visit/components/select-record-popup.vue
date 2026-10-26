<template>
  <van-popup v-model="visible" position="bottom" round @click-overlay="handleCancel">
    <div class="popup">
      <div class="popup-header">
        <div class="popup-title">{{title}}{{ $t('storeVisit.title.2') }}</div>
        <div class="popup-close" @click.stop="handleCancel"></div>
      </div>
      <!-- <div class="popup-tips">您有{{list.length}}场未开始的巡店计划</div> -->
      <div class="popup-content">
        <div class="list">
          <div class="item" v-for="item,index in list" :key="index" @click="handleClick(index)">
            <div class="flex-1">
              <div class="item-name">{{ item.store_name }}</div>
              <div class="item-row">
                <div class="label">{{ $t('storeVisit.item.label.1') }}:</div>
                <div class="value">{{item.date.replaceAll('-', '.')}} {{item.time_span}}</div>
              </div>
              
              <div class="item-row">
                <div class="label">{{ $t('storeVisit.item.label.2') }}:</div>
                <div class="value">
                  <span>{{item.creator.name}}</span>
                  <span class="dept" v-if="item.creator.dept_name">{{item.creator.dept_name}}</span>
                </div>
              </div>
              <div class="item-row">
                <div class="label">{{ $t('storeVisit.item.label.3') }}:</div>
                <div v-if="item.store_manager.length" class="value flex">
                  <div class="name" v-for="(m, index) in item.store_manager" :key="index">{{ m.name }}</div>
                </div>
                <div v-else class="value">--</div>
              </div>
            </div>
            <div class="item-checkbox"></div>
          </div>
        </div>
      </div>
      <!-- <div v-if="userInfo.is_store_manager || userInfo.is_store_visit_manager" class="popup-footer">
        <van-button color="#000000" plain @click.stop="handleRestart">发起新的巡店</van-button>
        <van-button color="#000000" :disabled="disabled(current)" @click.stop="handleConfirm">使用该巡店计划继续</van-button>
      </div> -->
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Model, Prop, Vue, Watch } from 'vue-property-decorator';

@Component({

})
export default class SelectRecordPopup extends Vue {
  @Model('update:show', { type: Boolean, default: false }) show: boolean;
  @Prop({ type: Array, default: () =>[] }) list: any[];
  @Prop({ type: String, default: '2025.03.05' }) title: string;
  
  visible: boolean = false;
  
  current: number = -1;

  /** 用户信息 */
  get userInfo():IUserInfo {
    return this.$store.getters["user/getUserInfo"]
  }
  /**
   * 创建人
   */
  get isCreator() {
    return (item) => item.creator.member_id === this.userInfo.member_id
  }

  get disabled() {
    return (index) => {
      if(index===-1) return true;
      const item = this.list[index];
      if(!this.isCreator(item)) return true;
      return false
    }
  }
  
  @Watch('show')
  onShowChange(val: boolean) {
    this.visible = val;
    if(!val) {
      this.current = -1
    }
  }

  /**
   * 取消
   */
  handleCancel() {
    this.$emit('update:show', false);
    this.$emit('cancel');
  }

  /**
   * 取消
   */
  handleRestart() {
    this.$emit('update:show', false);
    this.$emit('restart');
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
    this.$emit('confirm', this.list[this.current]);
  }

  handleClick(index) {
    this.current = index;
    this.handleConfirm()
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
    overflow-y: auto;
    padding-bottom: calc(env(safe-area-inset-bottom) + 24px);
    padding-bottom: calc(constant(safe-area-inset-bottom) + 24px);

    .list {
      padding: 8px 24px;

      .item {
        // height: 100px;
        display: flex;
        align-items: center;
        box-sizing: border-box;
        padding: 18px 16px;
        transition: all 0.3s;
        background-color: rgba(248, 249, 250, 1);
        color: rgba(0, 0, 0, 1);
        border-radius: 12px;
        position: relative;

        &:not(:last-child) {
          margin-bottom: 8px;
        }

        .flex-1 {
          flex: 1;
          overflow: hidden;
        }

        .item-name {
          font-weight: 600;
          font-size: 14px;
          line-height: 140%;
          margin-bottom: 12px;
        }

        .item-row {
          display: flex;
          font-size: 12px;
          line-height: 130%;

          .label {
            color: rgba(153, 153, 153, 1);
          }

          .value {
            margin-left: 8px;
            flex: 1;
            overflow: hidden;

            .dept {
              padding: 2px 4px;
              border-radius: 2px;
              font-size: 10px;
              line-height: 12px;
              font-weight: 700;
              margin-left: 4px;
              background-color: rgba(206, 212, 218, 1);
            }

            &.flex {
              display: flex;
              flex-wrap: wrap;
            }

            

            .name {
              margin-right: 8px;
              &:not(:last-child)::after {
                content: ',';
              }
            }
          }

          &:not(:last-child) {
            margin-bottom: 4px;
          }
        }

        .item-checkbox {
          width: 24px;
          height: 24px;
          background: url(~@/assets/store-visit/icon_check_circle.png) center/24px no-repeat;
          opacity: 0;
        }

        &.active {
          background-color: rgba(13, 110, 253, 1);
          color: #fff;
          box-shadow: 0 4px 12px 0 rgba(33, 37, 41, 0.2);
          
          .item-checkbox {
            opacity: 1;
          }

          .label {
            color: rgba(222, 226, 230, 1);
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
      height: 40px;
      border-radius: 4px;

      &:nth-child(1) {
        width: 120px;
      }
      &:nth-child(2) {
        flex: 1;
        margin-left: 8px;
      }

      &.van-button--disabled {
        opacity: 1;
        background: rgba(222, 226, 230, 1) !important;
        border-color: rgba(222, 226, 230, 1) !important;
      }
    }
  }
}
</style>