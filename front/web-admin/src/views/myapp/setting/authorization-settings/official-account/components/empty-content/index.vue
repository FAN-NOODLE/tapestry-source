<template>
  <div class="empty-content">
    <img src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/default/default_noList.png" alt="" />
    <p>{{ description }}</p>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
@Component
export default class EmptyContent extends Vue {
  @Prop({ default: "暂无数据" }) description!: string;
}
</script>
<style lang="scss" scoped>
.empty-content {
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
  img {
    width: 124px;
    height: 124px;
  }
  p {
    margin-top: 15px;
    color: #909399;
    font-size: 14px;
  }
}
</style>
