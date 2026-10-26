<template>
  <div class="star-box">
    <div v-for="(item, index) in list" :key="index" class="star-item">
      <img
        v-if="item.checked"
        class="star-icon"
        src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/customer/checked_star_icon.png"
      />
      <img
        v-else
        class="star-icon"
        src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/customer/unchecked_star_icon.png"
      />
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
type StarItem = {
  checked: boolean;
};
@Component
export default class extends Vue {
  @Prop({ default: {} }) list: Array<StarItem>; // 右侧选择标题
}
</script>

<style lang="scss" scoped>
.star-box {
  display: flex;
  width: 100%;
  .star-item {
    display: flex;
    margin-right: 4px;
    .star-icon {
      width: 24px;
      height: 22px;
      border-radius: 0 !important;
    }
  }
}
</style>
