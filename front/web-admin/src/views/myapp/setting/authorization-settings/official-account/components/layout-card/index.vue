<template>
  <div class="layout-card" ref="LayoutCard" :style="`height: ${height}`">
    <div class="title" :class="{ 'has-title': hasTitle }">
      <slot name="title"></slot>
    </div>
    <div class="content">
      <slot name="content"></slot>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
@Component
export default class LayoutList extends Vue {
  @Prop({ default: "auto" }) height!: string;
  private hasTitle = !!this.$slots["title"];
}
</script>
<style lang="scss" scoped>
.layout-card {
  background-color: $color-white;
  padding: 30px;
  border-radius: 10px;
  .title {
    font-size: 18px;
    font-weight: bolder;
  }
  .has-title {
    margin-bottom: 20px;
  }
}
</style>
