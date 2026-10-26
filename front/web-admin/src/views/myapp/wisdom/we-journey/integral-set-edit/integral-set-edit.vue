<template>
  <integral :info="detail" :type="type" :pageType="2" />
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Integral from "../components/integral/index.vue";
import { getIntegralDetailApi } from "@/api/we-journey/Integral";
import { ApiResult } from "@/types";

@Component({
  components: { Integral }
})
export default class IntegralEdit extends Vue {
  /** 详细信息 */
  private detail = {};

  /** 类型 1.邀约维度 2.任务维度 3.客户互动维度 4.培训维度 5.线上能力 */
  private type = 1;

  /** 获取详情 */
  private async getDetail() {
    const res: ApiResult = await getIntegralDetailApi({
      id: Number(this.$route.query.id)
    });
    if (res.code !== 200) return this.$message.error(res.message as string);
    this.detail = res.data;
    this.type = res.data.type;
    console.log("编辑页面查看详情接口", this.detail);
  }

  /** 生命周期 */
  private created() {
    this.getDetail();
  }
}
</script>
<style lang="scss" scoped></style>
