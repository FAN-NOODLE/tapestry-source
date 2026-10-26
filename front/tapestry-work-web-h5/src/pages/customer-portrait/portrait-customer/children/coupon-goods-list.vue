<template>
  <div>
    <empty-page
      picture="noData"
      paddingTop="60"
      v-if="couponList.length == 0"
      paddingBottom="60"
      text="暂无优惠券"
    ></empty-page>
    <van-list
      v-model="loading"
      :finished="finished"
      finished-text="没有更多了"
      offset="100"
      @load="getMorecoupon"
      v-else
    >
      <div
        class="outline-box"
        v-for="(coupon, index) in couponList"
        :key="index"
      >
        <div
          class="coupon-box"
          :class="{
            'color-blue': coupon.status == 2,
            'color-blue': coupon.status == 3,
            'color-gray': coupon.status == 4,
            'color-green': coupon.status == 1,
          }"
        >
          <div class="coupon-box-left">
            <p>
              <span v-if="coupon.coupon_value != '-'" class="mUnit">¥</span>
              <span>{{ formateMoney(coupon.coupon_value) }}</span>
            </p>
            <p>{{ coupon.coupon_name }}</p>
          </div>
          <div class="coupon-box-right">
            <div class="coupon-box-right-top">
              <p class="title">
                <span>批次号: </span><span> {{ coupon.coupon_stock_id }}</span>
              </p>
              <p class="mark">
                <span>有效期: </span
                ><span>
                  {{
                    dateformat(coupon.start_time) +
                    "-" +
                    dateformat(coupon.end_time)
                  }}</span
                >
              </p>
              <p class="mark" v-if="coupon.status == 2 || coupon.status == 3">
                <span>使用时间: </span
                ><span>
                  {{ dateformat(coupon.use_time, "yyyy-MM-dd hh:mm") }}</span
                >
              </p>
            </div>
            <div class="coupon-box-right-bottom">
              <span>{{
                coupon.is_all == 1 ? "全部商品可用" : "部分商品可用"
              }}</span>
              <van-icon
                name="arrow"
                @click="showGoodList(coupon.coupon_id)"
                v-if="coupon.is_all == 2 && coupon.status != 3"
              />
            </div>
          </div>
        </div>
        <!-- 这次不用 -->
        <p v-if="false" class="coupon-copy-btn">
          <van-button
            color="linear-gradient(130deg, #0285EC 0%, #0DADF6 100%)"
            size="small"
            type="primary"
            @click="copyOrderNo(coupon.order_no)"
          >
            复制订单编号
          </van-button>
        </p>
      </div>
    </van-list>

    <available-Product
      v-if="showavailableProduct"
      @close="showavailableProduct = false"
      :coupon-id="_coupon_id"
    ></available-Product>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import availableProduct from "@/components/marketing/available-product.vue";
import { Icon } from "vant";
import { customerCouponListApi } from "@/api/customer-portrait";
import emptyPage from "@/components/base/empty-page.vue";

@Component({
  components: {
    emptyPage,
    [Icon.name]: Icon,
    availableProduct,
  },
})
export default class extends Vue {
  private showavailableProduct = false;
  private _coupon_id = 0;
  private couponList = [
    // {
    //   coupon_id: 1,
    //   user_coupon_id: 1,
    //   coupon_name: "xxx",
    //   coupon_value: "xxx",
    //   effective: "xxx",
    //   status: 2, //使用状态（1:待使用，2:使用中，3:已过期）
    //   start_time: Date.now() / 1000,
    //   end_time: "",
    //   is_all: 2, //可用范围(1=全场，2=部分)
    //   range: "xxx",
    //   range_goods: [],
    //   order_no: "-",
    // },
    // {
    //   coupon_id: 1,
    //   user_coupon_id: 1,
    //   coupon_name: "xxx",
    //   coupon_value: "xxx",
    //   effective: "xxx",
    //   status: 3, //使用状态（1:待使用，2:使用中，3:已过期）
    //   start_time: Date.now() / 1000,
    //   end_time: "",
    //   is_all: 2, //可用范围(1=全场，2=部分)
    //   range: "xxx",
    //   range_goods: [],
    // },
    // {
    //   coupon_id: 1,
    //   user_coupon_id: 1,
    //   coupon_name: "xxx",
    //   coupon_value: "xxx",
    //   effective: "xxx",
    //   status: 1, //使用状态（1:待使用，2:使用中，3:已过期）
    //   start_time: Date.now() / 1000,
    //   end_time: "",
    //   is_all: 1, //可用范围(1=全场，2=部分)
    //   range: "xxx",
    //   range_goods: [],
    // },
  ];
  private finished = false;
  private loading = false;
  private page = 1;
  private size = 10;
  private mounted() {
    if (this.$store.getters["user/getCustomerId"]) {
      this.getMorecoupon();
    }
  }
  private async getMorecoupon() {
    this.loading = true;
    const { code, data, message }: IApiResult = await customerCouponListApi({
      user_id: this.$store.getters["user/getCustomerId"],
      page: this.page,
      size: this.size,
    });
    if (code == 200) {
      this.couponList = this.couponList.concat(data.list || []);
      if (
        data.total == this.couponList.length ||
        !data.list ||
        (data.list && data.list.length == 0)
      ) {
        this.finished = true;
      }
      this.loading = false;
      this.page++;
    } else {
      this.$toast.fail(message);
    }
  }
  private dateformat(time: number, format = "yyyy-MM-dd") {
    return this.$helper.dateFormat(time, format);
  }
  private formateMoney(money: number) {
    return money / 100;
  }
  private copyOrderNo(msg: string) {
    this.$helper.copyMessage(msg);
  }
  private showGoodList(id: number) {
    this._coupon_id = id;
    this.showavailableProduct = true;
  }
}
</script>
<style lang="scss" scoped>
.outline-box:not(:last-child) {
  border-bottom: 1px solid #ebebeb;
  padding-bottom: 20px;
}
.outline-box {
  margin: 20px 15px;
  background: white;
  .coupon-copy-btn {
    text-align: right;
    color: #fff;
    margin-top: 20px;
  }
  .color-blue {
    color: #0387ed;
    background: url(~@/assets/customer-portrait/coupon_2.png) no-repeat;
  }
  .color-gray {
    color: #868686;
    background: url(~@/assets/customer-portrait/coupon_3.png) no-repeat;
  }
  .color-green {
    color: #67c35b;
    background: url(~@/assets/customer-portrait/coupon_1.png) no-repeat;
  }
  .coupon-box {
    width: 100%;
    height: 115px;
    background-size: 100% 100%;
    display: flex;
    overflow: hidden;
    justify-content: center;
    align-items: center;
    .coupon-box-left {
      flex: 0 0 110px;
      width: 110px;
      overflow: hidden;
      .mUnit {
        font-size: 12px;
      }
      p {
        text-align: center;
        line-height: 20px;
      }
      p:first-child {
        font-weight: 600;
        line-height: 30px;
      }
    }
    .coupon-box-right {
      flex: 1;
      display: flex;
      overflow: hidden;
      justify-content: center;
      flex-direction: column;
      height: 100%;
      padding-left: 12px;
      .coupon-box-right-top {
        font-size: 12px;
        overflow: hidden;
        line-height: 24px;
        padding: 10px 2px;
        padding-top: 5px;
        height: 82px;
        .title {
          font-weight: 600;
        }
      }
      .coupon-box-right-bottom {
        display: flex;
        justify-content: space-between;
        justify-items: center;
        height: 32px;
        line-height: 32px;
        padding-right: 12px;
        i {
          padding-top: 10px;
        }
      }
    }
  }
}
</style>
