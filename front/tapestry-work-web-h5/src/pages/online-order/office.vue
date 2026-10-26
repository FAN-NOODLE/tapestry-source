<template>
  <div class="container">
    <van-tabs 
     v-model="active"
      title-inactive-color="#A6A6A6"
      color="#000000"
      line-width="40px"
      line-height="2px"
      :ellipsis="false"
      sticky
      @change="handleTabChange"
    >
      <van-tab title="订单数据">
        <order-data></order-data>
      </van-tab>
      <van-tab title="销售排行榜">
        <sale-rank></sale-rank>
      </van-tab>
      <van-tab title="商品排行榜">
        <goods-rank></goods-rank>
      </van-tab>
    </van-tabs>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import OrderData from './components/order-data.vue';
import SaleRank from './components/sale-rank.vue';
import GoodsRank from './components/goods-rank.vue';

@Component({
  components: {
    OrderData,
    SaleRank,
    GoodsRank
  }
})
export default class OnlineOrderOffice extends Vue {
  active: number = 0;
  handleTabChange(index: number) {
    
    // 神策数据上报 -- tab切换
    this.$track.trackTabSwitch({
      tab_name: ['订单数据', '销售排行榜'][index]
    })
  }
}
</script>
<style lang="scss" scoped>
::v-deep {
  .van-tabs__line {
    bottom: 22px
  }

  .van-tab {
    position: relative;

    &:nth-child(1)::after,
    &:nth-child(2)::after {
      content: '';
      position: absolute;
      right: 0;
      top: 50%;
      transform: translateY(-50%);
      width: 1px;
      height: 12px;
      background-color: #D1BFA7;
    }
  }
}
.container {
  width: 100vw;
  min-height: 100vh;
  background: #F7F7F7;
}

</style>