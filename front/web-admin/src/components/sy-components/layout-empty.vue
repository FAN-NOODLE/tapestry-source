<template>
  <div class="empty-page">
    <div class="default">
      <img :src="image_obj[empty_type]" alt="" class="empty-list-img" />
      <div v-if="!$slots.default" class="empty-data">{{ message }}</div>
      <slot v-else></slot>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
@Component({
  name: "LayoutEmpty"
})
export default class extends Vue {
  @Prop({ default: "暂无数据" }) message?: string; //默认文案

  /**
   * 缺省类别
   * list: 列表缺省
   * tag: 用户标签缺省
   * draw: 用户画像缺省
   * nowifi: 网络不佳
   * 404: 404
   */
  @Prop({ default: "list" }) empty_type?: any; //默认状态
  private image_obj = {
    list: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/empty-box.png",
    tag: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/empty-search.png",
    draw: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/default/default_noClient.png",
    collect: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/default/default_noClient.png",
    nowifi: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/default/default_netWorkError.png",
    noOpen: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/default/default_noOpen.gif",
    noRedEnvelope: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/default/default_noRedEnvelope.png",
    404: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/default/default_404.png"
  };
}
</script>
<style lang="scss" scoped>
.empty-page {
  width: 100%;
  height: 100%;
  text-align: center;
  padding-top: 100px;
  .default {
    .empty-data {
      // color: #292a32;
      color: #909399;
      font-size: 14px;
    }
    .empty-list-img {
      // width: 229px;
      // height: 197px;
      width: 124px;
      height: 124px;
    }
  }
}
</style>
