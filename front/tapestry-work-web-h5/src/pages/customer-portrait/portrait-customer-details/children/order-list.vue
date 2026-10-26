<template>
  <div class="order-list">
    <empty-page
      picture="noData"
      paddingTop="60"
      v-if="!orderListCurrent.length"
      paddingBottom="60"
      text="暂无订单"
    ></empty-page>
    <div class="service-info" v-else>
      <template v-for="(item, index) in orderListCurrent">
        <div class="list-item" :key="index">
          <div class="list-title">
            <span>下单时间：</span>
            <span>{{
              $helper.dateFormat(
                item.purchase_date / 1000,
                "yyyy-MM-dd hh:mm:ss"
              )
            }}</span>
          </div>
          <div class="list-line"></div>
          <div
            class="list-goods"
            v-for="(itm, idx) in item.purchase_skus"
            :key="idx"
          >
            <span class="sku-name">{{ itm.sku_name }}</span>
            <span>x {{ Math.abs(itm.quantity) }}</span>
          </div>
          <div class="list-info">
            <p>订单编号：{{ item.order_no }}</p>
            <p>订单金额：￥{{ item.tradePrice }}</p>
            <p>来源：{{ item.store_name }}</p>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import { NavBar, Image, Loading } from "vant";
import emptyPage from "@/components/base/empty-page.vue";
import { getOrderApi, getOrderDetail } from "@/api/customer-portrait/index";
@Component({
  components: {
    [NavBar.name]: NavBar,
    [Image.name]: Image,
    [Loading.name]: Loading,
    emptyPage,
  },
})
export default class extends Vue {
  @Prop({ default: "", type: String }) union_id!: string;
  private orderListCurrent = [];
  private created() {
    if (this.union_id != "union_id") {
      this.getOrderList();
    }
  }
  // 获取订单详情列表
  async getOrderList(): Promise<void> {
    this.$toast.loading();
    const { code, data, message }: IApiResult = await getOrderApi({
      union_id: this.union_id,
      // page: 1,
      // size: 99999,
    });
    this.$toast.clear();
    if (code == 200) {
      this.orderListCurrent = data.records || [];
      this.orderListCurrent.forEach((item, index) => {
        const amountArr = item.purchase_skus.map((skus) => Number(skus.amount));
        item.tradePrice = this.tradePriceReduce(amountArr);
        // this.getOrderDetail(item, index);
      });
    } else {
      if (message) {
        this.$toast.fail(message);
      }
    }
  }
  // 计算订单总额
  private tradePriceReduce(amountList: any) {
    return amountList.reduce((res, cur) => {
      return res + cur;
    });
  }
  // 获取订单详情
  async getOrderDetail(item: GetOrderDetail, index: number): Promise<void> {
    this.$toast.loading();
    const params: GetOrderDetailParams = {
      union_id: this.union_id,
      trade_no: item.tradeNo,
    };
    const { code, data, message }: IApiResult = await getOrderDetail(params);
    this.$toast.clear();
    if (code == 200) {
      item.goods = data.list || [];
      item.list = data.list || [];
      item.tradePrice = data.tradePrice;
      this.$set(this.orderListCurrent, index, item);
    } else {
      this.$toast.fail(message);
    }
  }
}
</script>
<style lang="scss" scoped>
.service-info {
  padding: 15px;
  .list-item {
    width: 100%;
    padding: 15px;
    background: #fff;
    border-radius: 4px;
    margin-bottom: 15px;
    .list-title {
      display: flex;
      justify-content: space-between;
      align-items: center;
      font-size: 14px;
      font-weight: 400;
      color: #33363c;
      line-height: 24px;
      span:last-child {
        font-weight: bold;
      }
    }
    .list-line {
      width: 100%;
      height: 1px;
      background: #f0f0f0;
      margin: 10px 0px;
    }
    .list-goods {
      display: flex;
      justify-content: space-between;
      align-items: center;
      font-size: 14px;
      font-weight: bold;
      color: #33363c;
      line-height: 28px;
      .sku-name {
        width: 260px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
      span:last-child {
        font-weight: 400;
      }
    }
    .list-info {
      margin-top: 10px;
      font-size: 14px;
      color: #33363c;
      line-height: 28px;
      font-weight: 400;
    }
  }
}
</style>
