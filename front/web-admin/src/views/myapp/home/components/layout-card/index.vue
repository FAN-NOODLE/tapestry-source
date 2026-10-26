<template>
  <div
    class="layout-card"
    ref="LayoutCard"
  >
    <div
      class="title"
      :class="{ 'has-title': hasTitle }"
    >
      <slot name="title"></slot>
    </div>
    <div class="content">
      <slot name="content"></slot>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
@Component
export default class LayoutList extends Vue {
  private hasTitle = !!this.$slots["title"];
}
</script>
<style lang="scss" scoped>
.layout-card {
  background-color: $color-white;
  padding: 30px;
  box-shadow: 0px 0px 20px 1px rgba(0, 0, 0, 0.08);
  .title {
    font-size: 18px;
    font-weight: bolder;
  }
  .has-title {
    margin-bottom: 20px;
  }
}
</style>
