<template>
  <router-view></router-view>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
// vuex

@Component
export default class extends Vue {
  // 计算属性
  // get number(): number {
  //   return 1;
  // }
}
</script>
<style lang="scss" scoped>
// 公共颜色
</style>
