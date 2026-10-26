<template>
  <van-popup v-model="visible" position="bottom" @click-overlay="handleCancel">
    <div class="popup">
      <div class="popup-close" @click.stop="handleCancel"></div>
      <div class="popup-title">共{{list.length}}件商品</div>
      <div class="popup-content">
        <div class="sticky">
          <div class="tips" v-if="showLoadingTips">网仓库存仍在查询中，若无库存则存在拒单风险，请确认是否继续下单</div>
          <div class="tips" v-else-if="showTips">监测到网仓中部分商品无库存，存在拒单风险，请确认是否继续下单</div>
          <div class="tips" v-else-if="showStockTips">监测到网仓中部分商品库存不足，存在拒单风险，请确认是否继续下单</div>
          <div class="title">请核对商品信息，提交后将无法修改！</div>
        </div>
        <div class="list">
          <div class="item" v-for="item, index in list" :key="index">
            <van-image 
              width="80px"
              height="80px"
              fit="contain"
              lazy-load
              :src="item.image||emptyImage"
            ></van-image>
            <div class="item-info">
              <div class="item-code">
                <span>{{item.product_code}}</span>
                <span>x{{item.quantity}}</span>
              </div>
              <div class="item-name">{{item.product_name}}</div>
              <!-- 加载中 -->
              <div v-if="item.loading" class="item-stock">
                <span>网仓库存</span>
                <van-loading type="spinner" size="12" />
              </div>
              <div v-else-if="item.dcl_on_hand==-1" class="item-stock error">商品库存未知</div>
              <div v-else-if="item.dcl_on_hand>0" class="item-stock">网仓库存 {{item.dcl_on_hand}}</div>
              <div v-else class="item-stock error">网仓无库存</div>
            </div>
          </div>
        </div>
      </div>
      <div class="popup-footer">
        <van-button color="#000000" plain @click.stop="handleCancel">返回修改</van-button>
        <van-button color="#000000" @click.stop="handleConfirm">继续下单</van-button>
      </div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Model, Prop, Vue, Watch } from 'vue-property-decorator';

@Component({

})
export default class OrderConfirmPopup extends Vue {
  @Model('update:show', { type: Boolean, default: false }) show: boolean;
  // 商品列表
  @Prop({ type: Array, default: () => [] }) list: any[];
  // 销售单号
  @Prop({ type: String, default: '' }) orderSn: string;
  emptyImage = require('@/assets/material-library/v4.2.0/icon_empty.png')

  visible: boolean = false;

  @Watch('show')
  onShowChange(val: boolean) {
    this.visible = val;
  }

  get showTips() {
    return this.list.some(item => item.dcl_on_hand == 0);
  }
  get showStockTips() {
    return this.list.some(item => item.dcl_on_hand < item.quantity);
  }
  get showLoadingTips() {
    return this.list.some(item => item.loading)
  }

  get riskType() {
    if(this.showLoadingTips) return '库存仍在查询中'
    if(this.showTips) return '部分商品无库存'
    if(this.showStockTips) return '部分商品库存不足'
    return ''
  }

  handleCancel() {
    // 神策数据上报 -- 提交浮窗点击
    this.$track.trackSubmitPopupClick({
      sales_order_id: this.orderSn,
      has_reject_risk: this.showTips || this.showLoadingTips || this.showStockTips,
      risk_type: this.riskType,
      button_name: '返回修改',
      customer_id: this.$store.getters["user/getExternalId"],
    })
    this.$emit('cancel');
  }

  handleConfirm() {
    // 神策数据上报 -- 提交浮窗点击
    this.$track.trackSubmitPopupClick({
      sales_order_id: this.orderSn,
      has_reject_risk: this.showTips || this.showLoadingTips || this.showStockTips,
      risk_type: this.riskType,
      button_name: '继续下单',
      customer_id: this.$store.getters["user/getExternalId"],
    })
    this.$emit('confirm', {
      risk_type: this.riskType,
      has_reject_risk: this.showTips || this.showLoadingTips || this.showStockTips
    });
  }
}
</script>
<style lang="scss" scoped>


.popup {
  position: relative;
  .popup-close {
    width: 48px;
    height: 48px;
    position: absolute;
    right: 0;
    top: 0;
    background: url(~@/assets/online-order/icon_popup_close.png) center/18px no-repeat;
  }
  .popup-title {
    height: 52px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 500;
    font-size: 16px;
    color: #000000;
  }

  .popup-content {
    max-height: 60vh;
    overflow-y: scroll;

    .sticky {
      padding: 5px 30px 20px;
    }

    .tips {
      background: #f7f7f7;
      font-weight: 400;
      font-size: 12px;
      color: #97989F;
      line-height: 17px;
      padding: 8px;
    }

    .title {
      font-weight: 400;
      font-size: 15px;
      color: #000000;
      line-height: 21px;
      text-align: center;
      margin-top: 15px;
    }

    .list {
      padding: 0 30px 20px;
    }
    .item {
      display: flex;
      position: relative;

      &:not(:last-child) {
        margin-bottom: 15px;
      }

      .van-image {
        background: #EFEFEF;
      }

      .item-info {
        margin-left: 10px;
        flex: 1;
        overflow: hidden;

        .item-code {
          font-weight: 400;
          font-size: 14px;
          color: #000000;
          line-height: 20px;
          display: flex;
          align-items: center;
          justify-content: space-between;
        }

        .item-name {
          font-weight: 400;
          font-size: 12px;
          color: #96979E;
          line-height: 17px;
          margin-top: 5px;
          height: 17px;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }

        .item-stock {
          font-weight: 400;
          font-size: 12px;
          color: #000000;
          line-height: 17px;
          margin-top: 14px;
          display: flex;
          align-items: center;

          &.error {
            color: #F64040;
          }
        }
      }
    }
  }

  .popup-footer {
    box-shadow: 0px 1px 6px 0px rgba(0,0,0,0.15);
    padding: 10px 15px;
    z-index: 2;
    display: flex;
    align-items: center;
    background-color: #fff;
    padding-bottom: calc(10px + contain(safe-area-inset-bottom));
    padding-bottom: calc(10px + env(safe-area-inset-bottom));

    .van-button {
      flex: 1;

      &:nth-child(2) {
        margin-left: 15px;
      }
    }
  }
}
</style>