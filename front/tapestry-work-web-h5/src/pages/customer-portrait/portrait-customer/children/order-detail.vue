<template>
  <div class="order-detail-nav">
    <van-nav-bar title="订单详情" left-arrow @click-left="backOrderList" />
    <div class="order-status">{{ statusObj[order.status] }}</div>
    <div class="order-personal-info">
      <div class="container">
        <p>
          <span class="name">{{ order.delivery_info.name }}</span>
          <span class="phone">{{
            order.delivery_info.mobile | phoneFormat
          }}</span>
        </p>
        <p class="address">
          {{
            order.delivery_info.province +
              order.delivery_info.city +
              order.delivery_info.city +
              order.delivery_info.area +
              order.delivery_info.address || "-"
          }}
        </p>
      </div>
    </div>
    <div class="good-list">
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
            <span v-if="good.goods_type == 1">{{
              good.sell_price == "-" ? "-" : "¥" + good.sell_price / 100
            }}</span>
            <van-button size="mini" v-else>{{
              good.goods_type == 2 ? "赠品" : "奖品"
            }}</van-button>
            <span>x{{ good.quantity }}</span>
          </p>
        </div>
      </div>
    </div>
    <div class="order-info">
      <p class="line">
        <span class="key">订单编号</span>
        <span class="value"
          ><span>{{ order.order_no }}</span>
          <van-button size="mini" @click="copyOrderNo(order.order_no)"
            >复制</van-button
          >
        </span>
      </p>
      <p class="line">
        <span class="key">支付方式 </span>
        <span class="value"
          >{{ payChannel[String(order.payment_channel)] || "-" }}
        </span>
      </p>
      <p class="line">
        <span class="key">支付时间 </span>
        <span class="value">{{ formatTime(order.pay_at) }} </span>
      </p>
      <p class="line">
        <span class="key">交易流水号 </span>
        <span class="value">{{ order.third_payment_no || "-" }} </span>
      </p>
    </div>
    <div class="good-info">
      <p class="line">
        <span class="key">商品总额 </span>
        <span class="value"
          >{{
            order.goods_pay_price == "-"
              ? "-"
              : "¥" + order.goods_pay_price / 100
          }}
        </span>
      </p>
      <p class="line">
        <span class="key">运费 </span>
        <span class="value"
          >{{
            order.logistics_price == "-"
              ? "-"
              : "¥" + order.logistics_price / 100
          }}
        </span>
      </p>
      <p class="line">
        <span class="key">商场优惠 </span>
        <span class="value blod-500"
          >{{ order.other_fee == "-" ? "-" : "¥" + order.other_fee / 100 }}
        </span>
      </p>
      <p class="line">
        <span class="key">微信优惠券 </span>
        <span class="value blod-500"
          >{{ order.coupon_fee == "-" ? "-" : "¥" + order.coupon_fee / 100 }}
        </span>
      </p>
      <p class="line">
        <span class="key key-color-33363C">应付款 </span>
        <span class="value blod-600"
          >{{ order.cash_fee == "-" ? "-" : "¥" + order.cash_fee / 100 }}
        </span>
      </p>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { NavBar, Image, Loading } from "vant";
import { customerOrderDetailApi } from "@/api/customer-portrait";
interface sku_attribute_list {
  attribute_name: string;
  attribute_item_value: string;
}
@Component({
  components: {
    [NavBar.name]: NavBar,
    [Image.name]: Image,
    [Loading.name]: Loading,
  },
  filters: {
    phoneFormat(phone: number | string) {
      let _phone = String(phone);
      if (_phone && _phone != "-") {
        _phone =
          _phone.slice(0, 3) + " " + _phone.slice(3, 7) + " " + _phone.slice(7);
      }
      return _phone;
    },
  },
})
export default class extends Vue {
  private order = {
    id: 1,

    order_no: "-", //订单编号
    payment_channel: "-",
    status: "-",
    third_payment_no: "-", //交易流水号
    logistics_price: "-",
    coupon_fee: "-",
    other_fee: "-",
    total_amount: "-",
    cash_fee: "-",
    pay_at: "-",
    goods_pay_price: "-",
    delivery_info: {
      name: "-",
      mobile: "-",
      address: "",
      province: "",
      city: "",
      area: "",
    },
    goods_list: [
      {
        sku_picture: "",
        spu_name: "-",
        sku_attribute_list: [],
        sell_price: "-",
        quantity: "-",
      },
    ],
  };
  private statusObj = {
    "0": "全部订单状态",
    "20": "待支付",
    "25": "待确认",
    "40": "待发货",
    "60": "待收货",
    "80": "已签收",
    "41": "售后中",
    "81": "售后中",
    "200": "退款中",
    "201": "退款失败",
    "83": "部分退款",
    "205": "已退款",
    "100": "已完成",
    "21": "已取消",
  };
  private payChannel = {
    "1": "微信小程序支付",
    "2": " 美团支付",
  };
  private created() {
    this.initPage();
  }
  private async initPage() {
    const { code, data, message }: IApiResult = await customerOrderDetailApi(
      Number(this.$route.query.orderNo)
    );
    if (code == 200) {
      this.order = data;
    } else {
      this.$toast.fail(message);
    }
  }
  private backOrderList() {
    this.$router.go(-1);
  }
  private copyOrderNo(msg: string) {
    this.$helper.copyMessage(msg);
  }
  private formatTime(time: number) {
    return this.$helper.dateFormat(time, "yyyy.MM.dd hh:mm:ss");
  }
  private formateAttr(sku_attribute_list: Array<sku_attribute_list>) {
    let stanards = "";
    (sku_attribute_list || []).forEach((element: sku_attribute_list) => {
      stanards =
        stanards + `${element.attribute_name}:${element.attribute_item_value};`;
    });
    return stanards.slice(0, -1);
  }
}
</script>
<style lang="scss" scoped>
::v-deep {
  .van-nav-bar .van-icon {
    color: #999;
  }
  .van-button--default {
    background: #f5f5f5;
  }
  .van-nav-bar__content {
    background: #f4f5f8;
  }
}
.order-detail-nav {
  background: #f4f5f8;
  .order-status {
    background: linear-gradient(295deg, #d7d7dd 0%, #a6a7b6 100%);
    font-size: 16px;
    color: #ffffff;
    height: 58px;
    line-height: 58px;
    padding-left: 15px;
  }
  .order-personal-info {
    padding: 15px;
    padding-bottom: 0;
    .container {
      border-bottom: 1px solid #d5d6d8;
      padding-bottom: 20px;
      font-size: 15px;
      color: #33363c;
      line-height: 21px;
      font-weight: 600;
      .name {
        padding-right: 10px;
      }
      .address {
        margin-top: 10px;
        font-size: 14px;
        color: #33363c;
        line-height: 20px;
        font-weight: 400;
        word-break: break-all;
      }
    }
  }
  .good-list {
    padding: 15px 0 0 0;
    margin: 0 15px;
    border-bottom: 1px solid #d5d6d8;
    .order-container {
      font-size: 14px;
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
  }
  .order-info {
    font-size: 14px;
    padding: 20px 0 0 0;
    margin: 0 15px;
    border-bottom: 1px solid #d5d6d8;
    .line {
      margin-bottom: 20px;
      display: flex;
      justify-content: space-between;
      overflow: hidden;
      .key {
        color: #96979e;
        flex: 0 0 80px;
      }
      .value {
        overflow: hidden;
        text-overflow: ellipsis;
        color: #272727;
        display: flex;
        justify-content: center;
        align-items: center;
        button {
          margin-left: 10px;
        }
      }
    }
  }
  .good-info {
    font-size: 14px;
    padding: 20px 0 0 0;
    margin: 0 15px;
    overflow: hidden;
    .line {
      margin-bottom: 20px;
      display: flex;
      justify-content: space-between;
      overflow: hidden;
      .blod-500 {
        font-weight: 500;
      }
      .blod-600 {
        font-weight: 600;
      }

      .key {
        color: #96979e;
        flex: 0 0 80px;
      }
      .key-color-33363C {
        color: #33363c;
      }
      .value {
        overflow: hidden;
        text-overflow: ellipsis;
        color: #33363c;
      }
    }
  }
}
</style>
