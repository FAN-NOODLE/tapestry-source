<template>
  <div>
    <!-- 客户详情 -->
    <customers-details></customers-details>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Watch } from "vue-property-decorator";
import CustomersDetails from "@/components/customer-detail/customer-detail.vue";
// vuex

@Component({
  components: {
    CustomersDetails
  }
})
export default class extends Vue {
  /********************************************data开始**********************************/

  /********************************************data结束**********************************/
  /***********************************************methods开始****************************/

  /**************************************************methods结束****************************/
  // 计算属性
  get number(): number {
    return 1;
  }
  // data 变量
  private username = "demo";
  // 监听
  @Watch("username")
  onChangeUsername(ewVal: string, oldVal: string) {
    console.log(ewVal, oldVal);
  }

  // 生命周期
  private created() {
    console.log("created", this.$router);
  }
  private mounted() {
    console.log("mounted");
  }
  private destroyed() {
    console.log("destroyed");
  }
}
</script>
<style lang="scss" scoped>
// 公共颜色
</style>
