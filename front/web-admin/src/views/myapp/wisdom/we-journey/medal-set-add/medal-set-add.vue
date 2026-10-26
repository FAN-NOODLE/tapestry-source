<template>
  <medal :info="detail" :pageType="1"></medal>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { medalDetailApi } from "@/api/we-journey/medal";
import Medal from "../components/medal/index.vue";
import { ApiResult } from "@/types";

@Component({
  components: { Medal }
})
export default class MedalAdd extends Vue {
  /** 详细信息 */
  private detail = {};

  /** 获取详情 */
  private async getDetail() {
    const id = this.$route.query.id;

    // 有id 表示当前是复制页面
    if (!id) return;
    const res: ApiResult = await medalDetailApi({
      id: Number(id)
    });
    if (res.code !== 200) return;
    delete res.data.id;
    this.detail = {
      ...res.data,
      rule_end_time: 0,
      rule_start_time: 0,
      name: res.data.name + "-副本"
    };
    console.log("编辑页面查看详情接口", this.detail);
  }

  /** 生命周期 */
  private created() {
    this.getDetail();
  }
}
</script>
<style lang="scss" scoped></style>
