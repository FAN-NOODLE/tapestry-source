<template>
  <medal-add :info="detail" :pageType="3"></medal-add>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { medalDetailApi } from "@/api/we-journey/medal";
import MedalAdd from "../components/medal/index.vue";
import { ApiResult } from "@/types";

@Component({
  components: { MedalAdd }
})
export default class MedalLook extends Vue {
  /** 详细信息 */
  private detail = {};

  /** 获取详情 */
  private async getDetail() {
    const res: ApiResult = await medalDetailApi({
      id: Number(this.$route.query.id)
    });
    if (res.code !== 200) return;
    this.detail = res.data;
    console.log("编辑页面查看详情接口", this.detail);
  }

  /** 生命周期 */
  private created() {
    this.getDetail();
  }
}
</script>
<style lang="scss" scoped></style>
