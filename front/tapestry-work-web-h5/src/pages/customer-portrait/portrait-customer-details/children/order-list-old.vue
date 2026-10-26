<template>
  <div class="order-list">
    <empty-page
      picture="noData"
      paddingTop="60"
      v-if="orderListCurrent.length == 0"
      paddingBottom="60"
      text="暂无订单"
    ></empty-page>
    <div class="service-info" v-else>
      <template v-for="(item, index) in orderListCurrent">
        <div class="list-item" :key="index">
          <template v-for="(i, j) in orderDefault">
            <template v-if="i.showDetile">
              <div
                class="goods-wrap"
                v-for="(x, y) in item.goods"
                :key="j + '' + y"
              >
                <p class="line" v-for="(x1, y1) in i.list" :key="y + '' + y1">
                  <span class="key" v-if="x1.title">{{ x1.title }}:</span>
                  <span class="value" :class="x1.class">{{ x[x1.key] }}</span>
                </p>
              </div>
            </template>
            <p class="line" :key="index + '' + j" v-else>
              <span class="key">{{ i.title }}:</span>
              <span class="value overflow" :class="item.class">{{
                i.enum ? i.enum[item[i.key]] : item[i.key]
              }}</span>
              <span v-if="i.call" class="img pointer">
                <img
                  src="@/assets/customer-portrait/copy-code.png"
                  alt="copy"
                  @click="copyItem(item[i.key])"
                />
              </span>
            </p>
          </template>
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
  private orderDefault = [
    {
      title: "订单编号",
      key: "tradeNo",
      call: true,
    },
    {
      title: "订单类型",
      key: "orderType",
      enum: { 1: "正单", 2: "积分兑换" },
    },
    {
      title: "支付时间",
      key: "tradeTime",
    },
    {
      title: "购买柜台",
      key: "store_name",
    },
    {
      key: "goods",
      showDetile: true,
      list: [
        {
          key: "productName",
        },
        {
          key: "skuCode",
          title: "SKU编号",
          class: "sku",
        },
      ],
    },
  ];
  private copyItem(data: string) {
    this.$helper.copyMessage(data);
  }
  private created() {
    if (this.union_id != "union_id") {
      this.getOrderList();
    }
  }
  // 获取订单详情列表
  async getOrderList(): Promise<void> {
    const { code, data, message }: IApiResult = await getOrderApi({
      union_id: this.union_id,
    });
    if (code == 200) {
      this.orderListCurrent = data || [];
      this.orderListCurrent.forEach((item, index) => {
        this.getOrderDetail(item, index);
      });
    } else {
      this.$toast.fail(message);
    }
  }
  // 获取订单详情
  async getOrderDetail(item: GetOrderDetail, index: number): Promise<void> {
    const params: GetOrderDetailParams = {
      union_id: this.union_id,
      trade_no: item.tradeNo,
    };
    const { code, data, message }: IApiResult = await getOrderDetail(params);
    if (code == 200) {
      item.goods = data.list || [];
      item.list = data.list || [];
      this.$set(this.orderListCurrent, index, item);
    } else {
      this.$toast.fail(message);
    }
  }
}
</script>
<style lang="scss" scoped>
.order-list {
  position: fixed;
  height: 100%;
  width: 100%;
  background-color: #ffffff;
}
.service-info {
  font-size: 14px;
  .list-item {
    padding: 20px 15px;
    position: relative;
    & + .list-item {
      &::after {
        content: "";
        position: absolute;
        height: 1px;
        top: 0;
        right: 0;
        left: 15px;
        background-color: #ebebeb;
      }
    }
    .line {
      &:nth-of-type(1) {
        margin-top: 0;
      }
    }
  }
  .goods-wrap {
    background: #f4f5f8;
    border-radius: 2px;
    padding: 15px;
    font-size: 12px;
    &.goods-wrap {
      margin-top: 15px;
    }
    .line {
      margin: 0;
      & + .line {
        margin-top: 10px;
      }
      .key {
        flex: 0 0 65px;
        padding-right: 5px;
      }
      .value {
        flex: 1;
        line-height: 18px;
        color: #333;
        &.sku {
          color: #888;
        }
      }
    }
  }
  .line {
    display: flex;
    align-items: center;
    margin: 15px 0;
    .key {
      flex: 0 0 70px;
      min-width: 80px;
      color: #808080;
      padding-right: 10px;
    }
    .value {
      color: #272727;
      white-space: normal;
      word-break: break-all;
      &.active-color {
        color: #0095f6;
      }
    }
    .overflow {
      height: 20px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
    .img {
      display: flex;
      align-items: center;
      margin-left: 10px;
      img {
        touch-action: none;
        width: 16px;
        height: 16px;
      }
    }
  }
}
</style>
