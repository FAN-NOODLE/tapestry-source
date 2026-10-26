<template>
  <empty-page
    picture="noData"
    paddingTop="60"
    v-if="!couponList.length"
    paddingBottom="60"
    text="暂无卡券"
  ></empty-page>
  <div class="container" v-else>
    <div class="all-coupon-list">
      <div
        class="coupon-content"
        v-for="(item, index) in couponList"
        :key="index"
      >
        <div class="coupon-content-top">
          <span>{{ item.type_str }}</span>
          <span>{{ item.status_str }}</span>
        </div>
        <div class="circle circle-left"></div>
        <div class="circle circle-right"></div>
        <div class="coupon-line"></div>
        <div class="coupon-content-bottom">
          <p>{{ item.name }}</p>
          <p>{{ item.startDate }} - {{ item.endDate }}</p>
          <p @click="showRule(index)">
            更多使用规则<van-icon
              :name="!item.isShow ? 'arrow-down' : 'arrow-up'"
              color="#4a4a4a"
            />
          </p>
          <p v-if="item.isShow" v-html="item.description">
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import { couponListApi } from "@/api/customer-portrait";
import emptyPage from "@/components/base/empty-page.vue";
import { Icon } from "vant";
@Component({
  components: {
    emptyPage,
    [Icon.name]: Icon,
  },
})
export default class extends Vue {
  @Prop({ default: "", type: String }) union_id!: string;
  private couponList = [];
  private mounted() {
    this.getCouponList();
  }
  private showRule(index: any) {
    const target = JSON.parse(JSON.stringify(this.couponList));
    target[index].isShow = !target[index].isShow;
    this.couponList = target;
  }
  async getCouponList(): Promise<void> {
    this.$toast.loading();
    const { code, data, message }: IApiResult = await couponListApi({
      union_id: this.union_id,
      status: 1,
      channel_type: "",
      source: 1,
      page: 1,
      size: 99999,
    });
    this.$toast.clear();
    if (code == 200) {
      this.couponList = data.records || [];
      this.couponList.forEach((item: any) => {
        item.isShow = false;
        // let timeData =
        //   item.Properties.binding_rule_properties[3].REDEEM_PERIOD_WITH_OFFSET;
        // item.endTime = timeData.endTime.split(" ")[0].replace(/-/g, "/");
        // item.startTime = timeData.startTime.split(" ")[0].replace(/-/g, "/");
      });
    } else {
      if (message) {
        this.$toast.fail(message);
      }
    }
  }
}
</script>
<style lang="scss" scoped>
div.container {
  position: fixed;
  height: 100%;
  width: 100%;
  background-color: #ffffff;
  .all-coupon-list {
    padding: 5px 15px;
    .coupon-content {
      width: 100%;
      margin-bottom: 15px;
      background-color: #f4f5f8;
      color: #4a4a4a;
      border-radius: 8px;
      position: relative;
      .circle {
        width: 14px;
        height: 14px;
        border-radius: 100%;
        background-color: #ffffff;
        position: absolute;
        top: 30px;
        z-index: 1;
        &.circle-left {
          left: -8px;
        }
        &.circle-right {
          right: -8px;
        }
      }

      .coupon-content-top {
        font-size: 12px;
        margin-top: 4px;
        padding: 10px 20px;
        font-weight: 400;
        display: flex;
        justify-content: space-between;
      }
      .coupon-line {
        width: 332px;
        margin: 0 auto;
        border-top: 1px dashed #cccccc;
      }
      .coupon-content-bottom {
        padding: 10px 20px;
        font-size: 10px;
        p:nth-child(1) {
          font-size: 14px;
          line-height: 20px;
          padding-bottom: 3px;
          font-weight: 500;
        }
        p:nth-child(2) {
          font-size: 12px;
          line-height: 17px;
          font-weight: 400;
        }
        p:nth-child(3) {
          display: flex;
          align-items: center;
          font-size: 12px;
          line-height: 17px;
          font-weight: 300;
          .van-icon {
            margin-left: 4px;
          }
        }
        p:nth-child(4) {
          margin-top: 6px;
          font-weight: 300;
          line-height: 14px;
          font-size: 10px;
        }
      }
    }
  }
}
</style>
