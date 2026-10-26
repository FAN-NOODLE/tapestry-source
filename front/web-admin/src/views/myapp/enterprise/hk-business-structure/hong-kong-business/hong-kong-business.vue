<template>
  <div class="hong-kong-business">
    <div class="business-tabs">
      <el-tabs v-model="activeTab" @tab-click="handleTabClick">
        <el-tab-pane label="门店管理" name="store"></el-tab-pane>
        <el-tab-pane label="职位权限管理" name="permission"></el-tab-pane>
      </el-tabs>
    </div>
    <router-view></router-view>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Watch } from "vue-property-decorator";

@Component
export default class HongKongBusiness extends Vue {
  activeTab = "store";

  @Watch("$route", { immediate: true })
  handleRouteChange() {
    this.activeTab = this.$route.path.indexOf("permission-") > -1
      ? "permission"
      : "store";
  }

  handleTabClick() {
    const name = this.activeTab === "store"
      ? "myapp_enterprise_hk-business-structure_hong-kong-business_store-management"
      : "myapp_enterprise_hk-business-structure_hong-kong-business_permission-management";
    if (this.$route.name !== name) {
      this.$router.push({ name });
    }
  }
}
</script>

<style lang="scss" scoped>
.business-tabs {
  padding: 10px 30px 0;
  background: $color-white;

  /deep/ .el-tabs__header {
    margin-bottom: 0;
  }

  /deep/ .el-tabs__item {
    min-width: 136px;
    height: 42px;
    line-height: 42px;
    text-align: center;
    font-size: 16px;
  }
}
</style>
