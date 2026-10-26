<template>
  <empty-page
    picture="noData"
    paddingTop="60"
    v-if="integralList.length == 0 && showEmpty"
    paddingBottom="60"
    text="暂无积分明细"
  ></empty-page>
  <div class="container" v-else>
    <div class="integral-all">
      当前积分：<span>{{ evalid_points }}</span>
    </div>
    <div
      class="integral-detail"
      v-for="(item, index) in integralList"
      :key="index"
    >
      <div class="detail-title">
        {{ integralDictionaries[item.pointType].name }}
      </div>
      <div class="cut-off-rule"></div>
      <div class="detail-data">
        <div class="detail-data-left">
          <p v-if="item.tradeNo && item.tradeNo !== ''">
            订单金额：{{ item.tradePrice }}
          </p>
          <p>{{ item.createDateTime }}</p>
        </div>
        <div class="detail-data-right">积分：{{ item.points }}</div>
      </div>
      <div class="detail-good-list" v-for="(itm, idx) in item.goods" :key="idx">
        <span>{{ itm.productName }}</span>
        <span>x {{ Math.abs(itm.productCount) }}</span>
      </div>
    </div>
    <div
      class="bottom-button"
      @click="toLoadData"
      v-if="integralNewList.length > 0 && integralList.length > 10"
    >
      ----------点击加载更多----------
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import { integralListApi, getOrderDetail } from "@/api/customer-portrait";
import emptyPage from "@/components/base/empty-page.vue";
@Component({
  components: {
    emptyPage,
  },
})
export default class extends Vue {
  @Prop({ default: "", type: String }) evalid_points!: string;
  @Prop({ default: "", type: String }) union_id!: string;
  private integralList = [];
  private integralNewList = [];
  private page = 1;
  private showEmpty = false;
  private goodList = [];
  private integralDictionaries = {
    CBP: {
      name: "消费积分",
      type: 1,
    },
    TBP: {
      name: "等级奖励积分",
      type: 1,
    },
    HBP: {
      name: "行为奖励积分",
      type: 1,
    },
    FGP: {
      name: "兑换失败积分",
      type: 1,
    },
    CGP: {
      name: "取消兑礼积分",
      type: 1,
    },
    OBP: {
      name: "生日奖励",
      type: 1,
    },
    ABP: {
      name: "手工加积分",
      type: 1,
    },
    BBP: {
      name: "退单",
      type: 2,
    },
    PBP: {
      name: "积分过期",
      type: 1,
    },
    PCMC: {
      name: "取消积分消费",
      type: 1,
    },
    FCP: {
      name: "取消兑礼失败积分",
      type: 2,
    },
    RDP: {
      name: "兑换积分",
      type: 2,
    },
    MBP: {
      name: "手工减积分",
      type: 2,
    },
    PCM: {
      name: "积分消费",
      type: 2,
    },
    // 线下
    DBP: {
      name: "消费积分",
      type: 1,
    },
    FBP: {
      name: "行为积分",
      type: 1,
    },
    HDP: {
      name: "退单",
      type: 2,
    },
    GDP: {
      name: "手工加积分",
      type: 1,
    },
    OFGP: {
      name: "兑换失败积分",
      type: 1,
    },
    OCGP: {
      name: "取消兑礼积分",
      type: 1,
    },
    EBP: {
      name: "积分过期",
      type: 2,
    },
    FDP: {
      name: "手工减积分",
      type: 2,
    },
    ODP: {
      name: "兑换积分",
      type: 2,
    },
    OFCP: {
      name: "取消兑礼失败积分",
      type: 2,
    },
  };
  private mounted() {
    this.getintegralList();
  }
  async getintegralList(): Promise<void> {
    const { code, data, message }: IApiResult = await integralListApi({
      union_id: this.union_id,
      page_no: this.page.toString(),
      limit: "10",
    });
    if (code == 200) {
      this.integralNewList = data || [];
      this.showEmpty = true;
      this.integralNewList.forEach((item, index) => {
        item.createDateTime = this.changeData(item.createDateTime);
        if (item.tradeNo && item.tradeNo !== "") {
          this.getOrderDetail(item, index);
        }
      });
      this.integralList = [...this.integralList, ...this.integralNewList];
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
      item.tradePrice = data.tradePrice || "";
      this.$set(this.integralNewList, index, data);
    } else {
      this.$toast.fail(message);
    }
  }
  // 改变时间
  private changeData(dataTime) {
    if (dataTime) {
      return dataTime.split(" ")[0];
    }
    return "-";
  }
  private toLoadData() {
    if (this.integralNewList.length > 0) {
      this.page++;
      this.getintegralList();
    }
  }
}
</script>
<style lang="scss" scoped>
div.container {
  padding: 15px;
  .integral-all {
    display: flex;
    justify-content: center;
    align-items: center;
    font-size: 16px;
    font-weight: 500;
    color: #96979e;
    padding: 30px 0;
    span {
      font-size: 24px;
      font-weight: bold;
      color: #33363c;
    }
  }
  .integral-detail {
    width: 100%;
    padding: 15px;
    box-sizing: border-box;
    background: #ffffff;
    border-radius: 4px;
    margin-bottom: 15px;
    .detail-title {
      font-size: 14px;
      font-weight: bold;
      color: #33363c;
    }
    .cut-off-rule {
      width: 100%;
      height: 1px;
      background: #f0f0f0;
      margin: 14px 0px;
    }
    .detail-data {
      display: flex;
      justify-content: space-between;
      align-items: center;
      font-size: 14px;
      line-height: 20px;
      margin-bottom: 20px;
      .detail-data-left {
        color: #96979e;
      }
      .detail-data-right {
        font-weight: 600;
        color: #000000;
      }
    }
    .detail-good-list {
      display: flex;
      justify-content: space-between;
      align-items: center;
      font-size: 14px;
      color: #33363c;
      line-height: 20px;
    }
  }
  .bottom-button {
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
    color: #96979e;
  }
}
</style>
