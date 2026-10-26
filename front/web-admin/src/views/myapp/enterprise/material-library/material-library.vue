<template>
  <keep-alive include="MaterialList">
    <router-view></router-view>
  </keep-alive>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
// vuex

@Component
export default class extends Vue {}
</script>
<style lang="scss" scoped>
// 公共颜色
</style>
