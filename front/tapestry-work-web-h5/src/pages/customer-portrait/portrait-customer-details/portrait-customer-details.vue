<template>
  <div>
    <van-nav-bar :title="$route.query.title" left-arrow @click-left="goBack" />
    <integral-detail
      v-if="$route.query.title == '积分明细'"
      :evalid_points="$route.query.evalid_points"
      :union_id="$route.query.union_id"
    />
    <coupon-list
      v-if="isKS && $route.query.title == '卡券详情'"
      :union_id="$route.query.union_id"
    />
    <sw-coupon-list
      v-if="isSW && $route.query.title == '卡券详情'"
      :union_id="$route.query.union_id"
    />
    <order-list
      v-if="$route.query.title == '消费记录'"
      :union_id="$route.query.union_id"
    />
    <footprint-list v-if="$route.query.title == '客户足迹'" />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import couponList from "./children/coupon-list.vue";
import swCouponList from "./children/sw-coupon-list.vue";
import footprintList from "./children/foot-print-list.vue";
import orderList from "./children/order-list.vue";
import integralDetail from "./children/integral-detail.vue";
import { NavBar } from "vant";
import Helper from "@/utils/helper";
@Component({
  components: {
    [NavBar.name]: NavBar,
    couponList,
    footprintList,
    orderList,
    integralDetail,
    swCouponList
  },
})
export default class extends Vue {
  // 是否是sw
  private isSW = Helper.isSW();
  // 是否是ks
  private isKS = Helper.isKS();
  private goBack() {
    this.$router.go(-1);
  }
}
</script>
<style lang="scss" scoped>
::v-deep {
  .van-nav-bar .van-icon {
    color: #000000;
  }
  .van-nav-bar__title {
    font-weight: bold;
    color: #000000;
    font-size: 17px;
  }
}
</style>
