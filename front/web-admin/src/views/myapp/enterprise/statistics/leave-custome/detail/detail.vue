<template>
  <div>
    <!-- 客户详情 -->
    <customer-details></customer-details>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import CustomerDetails from "@/components/customer-detail/customer-detail.vue";
// vuex

@Component({
  components: {
    CustomerDetails
  }
})
export default class extends Vue {
  // data 变量
  private username = "demo";

  // 生命周期
  // private created() {
  //   console.log("created");
  // }
  // private mounted() {
  //   console.log("mounted");
  // }
  // private destroyed() {
  //   console.log("destroyed");
  // }
}
</script>
<style lang="scss" scoped>
// 公共颜色
</style>
