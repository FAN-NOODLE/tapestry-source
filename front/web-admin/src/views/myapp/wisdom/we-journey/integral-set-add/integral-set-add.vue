<template>
  <integral :data="detail" :type="type" :pageType="1" />
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Integral from "../components/integral/index.vue";

@Component({
  components: { Integral }
})
export default class IntegralAdd extends Vue {
  /** 详细信息 */
  private detail = {};

  /** 类型 1.邀约维度 2.任务维度 3.客户互动维度 4.培训维度 5.线上能力 */
  private type = 1;

  /** 生命周期 */
  private created() {
    this.type = Number(this.$route.query.type);
  }
}
</script>
<style lang="scss" scoped></style>
