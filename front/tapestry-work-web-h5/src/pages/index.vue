<template>
  <div>
    <van-button> 点击我</van-button>
    <cy-select
      placeholder="全部订单状态"
      :value="selectVal"
      :options="list"
      @select="selectFn"
    />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { Button } from "vant";
import cySelect from "@/components/base/cy-select.vue";
@Component({
  components: {
    [Button.name]: Button,
    cySelect,
  },
})
export default class extends Vue {
  private selectVal = "";
  private selectFn(val) {
    console.log(val);
  }
  private list = [
    { label: "全部订单状态", value: 0 },
    { label: "待支付", value: 20 },
    { label: "待确认", value: 25 },
    { label: "待发货", value: 40 },
    { label: "已签收", value: 80 },
    { label: " 售后中", value: 41 },
    { label: " 退款中", value: 200 },
    { label: "退款失败", value: 201 },
    { label: "部分退款", value: 83 },
    { label: "已退款", value: 205 },
    { label: "已完成", value: 100 },
    { label: "已取消", value: 21 },
  ];
}
</script>
