<template>
  <router-view></router-view>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
@Component
export default class extends Vue {
  // 生命周期
  private created() {
    console.log("created");
  }
  private mounted() {
    console.log("mounted");
  }
  private destroyed() {
    console.log("destroyed");
  }
}
</script>
<style lang="scss" scoped></style>
