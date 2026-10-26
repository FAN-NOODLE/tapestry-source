<template>
  <router-view></router-view>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
@Component
export default class extends Vue {
  private mounted() {
    // console.log("mounted");
  }
  private destroyed() {
    // console.log("destroyed");
  }
}
</script>
<style lang="scss" scoped>
// 公共颜色
</style>
