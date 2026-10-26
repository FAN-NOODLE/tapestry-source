<template>
  <div class="item">
    <div class="item-status" :class="`status-${item.status}`">{{statusText}}</div>
    <div class="item-header">
      <div class="item-code">销售单号 {{ item.order_sn }}</div>
      <div class="item-num"><span v-if="item.send_quantity">发{{item.send_quantity}}件/</span>共{{item.total_quantity}}件</div>
    </div>
    <div class="item-content">
      <div class="item-row">
        <span class="item-row-label">下单客户</span>
        <span class="item-row-value">{{item.user_name}}</span>
      </div>
      <div class="item-row">
        <span class="item-row-label">创建时间</span>
        <span class="item-row-value">{{ $helper.dateFormat(item.created_at, "yyyy-MM-dd hh:mm:ss") }}</span>
      </div>
      <van-button color="#000000" plain block @click="handleViewDetail">查看详情</van-button>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';

export interface IOrderItem {
  [key: string]: any
}

@Component({

})
export default class OrderItem extends Vue {
  @Prop({ type: Object, default: () => ({}) }) item: IOrderItem;
  
  get statusText() {
    return ['', '待填写', '待发货', '已发货', '已拒单', '已作废'][this.item.status]
  }


  handleViewDetail() {
    this.$router.push({
      path: 'online-order-detail',
      query: {
        order_id: this.item.id
      }
    })
    // 神策数据上报 -- 订单列表点击
    this.$track.trackOrderListClick({
      sales_order_id: this.item.order_sn,
      order_status: this.statusText,
      order_creator: this.item.member_name,
      order_belonged_store: this.item.store_name,
      customer_id: this.item.external_user_id,
    })
  }
}
</script>
<style lang="scss" scoped>
.item {
  position: relative;
  background-color: #fff;
  margin-bottom: 15px;

  .item-header {
    padding: 12px 65px 12px 15px;
    position: relative;
    display: flex;
    align-items: center;

    &::after {
      content: '';
      position: absolute;
      left: 0;
      right: 0;
      bottom: 0;
      height: 1px;
      background-color: #F0F0F0;
    }

    .item-code {
      font-size: 14px;
      color: #000;
      font-weight: bold;
      flex: 1;
    }
    .item-num {
      font-weight: 400;
      font-size: 14px;
      color: #A6A6A6;
    }
  }

  .item-content {
    padding: 12px 15px;

    .item-row {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: 12px;
      line-height: 20px;
      
      .item-row-label {
        font-weight: 400;
        font-size: 14px;
        color: #A6A6A6;
      }

      .item-row-value {
        font-weight: 400;
        font-size: 14px;
        color: #000000;
      }
    }

    .van-button {
      height: 36px;
      margin-top: 15px;
    }
  }

  .item-status {
    width: 66px;
    height: 32px;
    background: url('~@/assets/online-order/icon_s_1.png') center/100% no-repeat;
    position: absolute;
    top: 8px;
    right: -6px;
    font-weight: 400;
    font-size: 14px;
    color: #000000;
    padding-bottom: 4px;
    display: flex;
    align-items: center;
    justify-content: center;

    &.status-1 {
      background-image: url('~@/assets/online-order/icon_s_0.png');
      color: #fff;
    }
    &.status-2 {
      background-image: url('~@/assets/online-order/icon_s_1.png');
    }
    &.status-3 {
      background-image: url('~@/assets/online-order/icon_s_2.png');
    }
    &.status-4 {
      background-image: url('~@/assets/online-order/icon_s_3.png');
    }
    &.status-5 {
      background-image: url('~@/assets/online-order/icon_s_4.png');
    }
  }
}
</style>