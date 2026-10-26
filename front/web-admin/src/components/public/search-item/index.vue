<template>
  <div class="search-item">
    <span v-text="text"></span>
    <slot></slot>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
@Component
export default class LayoutList extends Vue {
  @Prop() text!: string;
}
</script>
<style lang="scss" scoped>
.search-item {
  display: inline-block;
  color: $color-text-normal;
  margin-right: 20px;
  margin-bottom: 15px;
}
</style>
