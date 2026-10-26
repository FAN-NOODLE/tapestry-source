<template>
  <div class="order-container">
    <van-row class="search-container">
      <van-col span="12" class="search-content">
        <van-search
          v-model.trim="orderNo"
          placeholder="请输入订单编号"
          @search="searchOrderList"
          @clear="resetList"
        />
      </van-col>
      <van-col span="12" class="search-content" ref="dropdown-box">
        <van-dropdown-menu active-color="#0095F6">
          <van-dropdown-item
            v-model="orderStatus"
            :options="orderStatusList"
            @change="changeOrderStatus"
            @open="openDropdown"
          />
        </van-dropdown-menu>
      </van-col>
    </van-row>
    <empty-page
      picture="noData"
      paddingTop="60"
      v-if="orderList.length == 0"
      paddingBottom="60"
      :text="orderNo ? '没有找到订单' : '暂无订单'"
    ></empty-page>
    <div class="order-list" v-else>
      <van-list
        v-model="loading"
        :finished="finished"
        finished-text="没有更多了"
        offset="100"
        @load="getMoreOrder"
      >
        <div
          class="order-item"
          v-for="(order, index) in orderList"
          :key="index"
        >
          <div class="desc">
            <div class="order-no">
              <p>
                <span class="key">订单编号: </span>
                <span class="value">{{ order.order_no }}</span>
                <img
                  src="@/assets/customer-portrait/copy_icon.png"
                  class="copy"
                  @click="copyOrderNo(order.order_no)"
                  alt=""
                />
              </p>
              <span class="order_status">
                {{ statusObj[String(order.status)] }}
              </span>
            </div>
            <div class="pay-time">
              <span class="key">支付时间: </span>
              <span class="value">{{ formatTime(order.pay_at) }}</span>
            </div>
          </div>
          <div
            class="order-container"
            v-for="(good, ind) in order.goods_list || []"
            :key="ind"
          >
            <div class="left">
              <van-image :src="good.sku_picture" width="100%" height="100%" fit>
                <template v-slot:loading>
                  <van-loading type="spinner" size="20" />
                </template>
              </van-image>
            </div>
            <div class="right-content">
              <div>
                <p class="good-name">
                  {{ good.spu_name }}
                </p>
                <p class="good-standards">
                  {{ formateAttr(good.sku_attribute_list) }}
                </p>
              </div>

              <p class="good-money">
                <span v-if="good.goods_type == 1">
                  <span v-if="order.final_amount != '-'">¥</span>
                  {{ formateMoney(good.sell_price) }}</span
                >
                <vant-button size="mini" v-else>{{
                  good.goods_type == 2 ? "赠品" : "奖品"
                }}</vant-button>
                <span>x{{ good.quantity }}</span>
              </p>
            </div>
          </div>
          <div class="order-summary">
            <div class="order-summary-money">
              <span></span>
              <p>
                <span>实付金额: </span
                ><span>
                  <span v-if="order.cash_fee != '-'">¥</span>
                  {{ formateMoney(order.cash_fee) }}</span
                >
              </p>
            </div>
            <div class="order-summary-detail">
              <span></span>
              <van-button
                class="order-btn"
                size="small"
                @click="lookOrderDetail(order.order_no)"
              >
                查看详情</van-button
              >
            </div>
          </div>
        </div>
      </van-list>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Watch } from "vue-property-decorator";
import emptyPage from "@/components/base/empty-page.vue";
import { Row, Col, Image, Loading, Toast } from "vant";
import { customerOrderListApi } from "@/api/customer-portrait/index";
interface sku_attribute_list {
  attribute_name: string;
  attribute_item_value: string;
}
@Component({
  components: {
    emptyPage,
    [Col.name]: Col,
    [Row.name]: Row,
    [Image.name]: Image,
    [Loading.name]: Loading,
  },
})
export default class extends Vue {
  private orderNo = "";
  private orderStatus = 0;
  private statusObj = {
    "0": "全部订单状态",
    "20": "待支付",
    "25": "待确认",
    "40": "待发货",
    "60": "待收货",
    "80": "已签收",
    "81": "售后中",
    "41": "售后中",
    "200": "退款中",
    "201": "退款失败",
    "83": "部分退款",
    "205": "已退款",
    "100": "已完成",
    "21": "已取消",
  };
  private orderStatusList = [
    { text: "全部订单状态", value: 0 },
    { text: "待支付", value: 20 },
    { text: "待确认", value: 25 },
    { text: "待发货", value: 40 },
    { text: "待收货", value: 60 },
    { text: "已签收", value: 80 },
    { text: " 售后中", value: 41 },
    { text: " 退款中", value: 200 },
    { text: "退款失败", value: 201 },
    { text: "部分退款", value: 83 },
    { text: "已退款", value: 205 },
    { text: "已完成", value: 100 },
    { text: "已取消", value: 21 },
  ];
  private loading = false;
  private finished = false;
  private page = 1;
  private size = 10;
  private orderList = [];
  public get externalId(): string {
    return this.$store.getters["user/getCustomerId"];
  }

  @Watch("externalId")
  initOrderListByuserId(newVal: string): void {
    if (newVal) {
      this.getMoreOrder();
    }
  }
  mounted(): void {
    if (this.$store.getters["user/getCustomerId"]) {
      this.getMoreOrder();
    }
  }
  private formatTime(time: number) {
    return this.$helper.dateFormat(time, "yyyy-MM-dd hh:mm:ss");
  }
  private formateMoney(money: number) {
    return money / 100;
  }
  private formateAttr(sku_attribute_list: Array<sku_attribute_list>) {
    let stanards = "";
    (sku_attribute_list || []).forEach((element: sku_attribute_list) => {
      stanards =
        stanards + `${element.attribute_name}:${element.attribute_item_value};`;
    });
    return stanards.slice(0, -1);
  }
  private lookOrderDetail(orderNo: number) {
    this.$router.push({
      path: "/order-detail",
      query: Object.assign({ orderNo: orderNo }, this.$route.query),
    });
  }
  private async getMoreOrder() {
    this.loading = true;
    const no: any = Number(this.orderNo);
    if (String(no) && (typeof no != "number" || isNaN(no))) {
      Toast({
        duration: 3000,
        message: "请输入有效的订单编号!",
      });
      this.loading = false;
      return;
    }
    const { code, data, message }: IApiResult = await customerOrderListApi({
      user_id: this.$store.getters["user/getCustomerId"],
      order_no: no || undefined,
      status: this.orderStatus,
      page: 1,
      size: 20,
    });

    if (code == 200) {
      this.loading = false;
      this.orderList = this.orderList.concat(data.list || []);
      if (
        data.total == this.orderList.length ||
        !data.list ||
        (data.list && data.list.length == 0)
      ) {
        this.finished = true;
      }
    } else {
      this.$toast.fail(message);
    }

    this.page++;
  }
  private copyOrderNo(orderNo: string) {
    this.$helper.copyMessage(orderNo);
  }
  searchUpdateList(): void {
    this.page = 1;
    this.finished = false;
    this.orderList = [];
    this.getMoreOrder();
  }
  private openDropdown() {
    const box: any = this.$refs["dropdown-box"]; // 获取元素
    this.$nextTick(() => {
      const top = box.getBoundingClientRect().top;
      const windowHeight: any = document.body.clientHeight;
      if (windowHeight - top < 300) {
        window.scrollBy(0, 300);
      }
    });
  }
  private changeOrderStatus() {
    this.searchUpdateList();
  }
  private resetList() {
    this.searchUpdateList();
  }
  private searchOrderList() {
    this.searchUpdateList();
  }
}
</script>
<style lang="scss" scoped>
.search-container {
  margin: 0 15px;
}
::v-deep .van-dropdown-menu__title {
  font-size: 14px;
}
.van-search {
  padding: 0;
}
.van-search .van-cell {
  height: 32px;
  line-height: 24px;
}
.search-content:nth-child(1) {
  padding: 15px 7.5px 15px 0;
}
.search-content:nth-child(2) {
  padding: 15px 0 15px 7.5px;
}
.order-item:not(:nth-child(1)) {
  border-top: 1px solid #ebebeb;
  margin-top: 15px;
  padding-top: 15px;
}
.order-item .order-container:not(:nth-child(1)) {
  margin-top: 5px;
}

.order-item {
  margin: 0 15px;
  .desc {
    .key {
      color: #808080;
    }
    .order-no {
      display: flex;
      justify-content: space-between;
      margin-bottom: 15px;
      .copy {
        margin-left: 10px;
        width: 13px;
        height: 15px;
        position: relative;
        top: 2px;
      }
      .order_status {
        color: #272727;
        font-weight: bold;
        font-size: 12px;
        display: inline-block;
        height: 20px;
        line-height: 20px;
      }
    }
    .pay-time {
      margin-bottom: 15px;
    }
  }
  .order-container {
    background: #f4f5f8;
    border-radius: 2px;
    border: 1px solid #f4f5f8;
    height: 100px;
    padding: 10px 9px;
    display: flex;
    .left {
      flex: 0 0 75px;
      max-width: 75px;
      height: 80px;
    }
    .right-content {
      flex: 1;
      padding-left: 8px;
      overflow: hidden;
      display: flex;
      flex-direction: column;
      justify-content: space-between;
      .good-standards {
        font-size: 12px;
        color: #96979e;
        line-height: 12px;
        margin-top: 5px;
      }
      .good-name {
        word-wrap: break-word;
        overflow: hidden;
        text-overflow: ellipsis;
        display: -webkit-box;
        -webkit-box-orient: vertical;
        -webkit-line-clamp: 2;
        color: #313131;
        line-height: 18px;
        font-size: 12px;
        min-height: 36px;
      }

      .good-money {
        display: flex;
        justify-content: space-between;
        font-size: 12px;
        color: #313131;
      }
    }
  }
  .order-summary {
    margin-top: 15px;
    .order-summary-money {
      display: flex;
      justify-content: space-between;
    }
    .order-summary-detail {
      margin-top: 15px;
      display: flex;
      justify-content: space-between;
      .order-btn {
        height: 26px;
        border-radius: 4px;
        padding: 0 16px;
        border: solid 1px #b9bec4;
        font-size: 12px;
        color: #272727;
        line-height: 24px;
      }
    }
  }
}
</style>
