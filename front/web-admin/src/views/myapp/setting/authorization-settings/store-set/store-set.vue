<template>
  <router-view></router-view>
</template>
<script lang="ts">
import { Component, Vue, Watch } from "vue-property-decorator";
// vuex

@Component
export default class extends Vue {
  // 计算属性
  get number(): number {
    return 1;
  }
  // data 变量
  private username = "demo";
  // 监听
  @Watch("username")
  onChangeUsername(ewVal: string, oldVal: string) {
    console.log(ewVal, oldVal);
  }

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
<style lang="scss" scoped>
// 公共颜色
</style>
