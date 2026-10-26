<template>
  <div class="data-compass">
    <div v-if="ENTERPRISE_MARK === 'HR'">
      <layout-list>
        <template slot="top-tabs">
          <el-tabs v-model="active_name">
            <!-- <el-tab-pane label="订单数据" name="1"></el-tab-pane> -->
            <el-tab-pane label="顾客数据" name="2"></el-tab-pane>
            <el-tab-pane label="导购数据" name="4"></el-tab-pane>
            <!-- <el-tab-pane label="分销数据" name="3"></el-tab-pane> -->
          </el-tabs>
        </template>
        <template slot="footer">
          <!-- 顾客数据 -->
          <customer-data v-if="active_name === '2'"></customer-data>
          <!-- 导购数据 -->
          <guide-data v-if="active_name === '4'"></guide-data>
        </template>
      </layout-list>
    </div>
    <div v-else>
      <!-- 除了赫莲娜外，通用的数据罗盘 -->
      <compass-ys />
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import customerData from "./components/customer-data.vue";
import guideData from "./components/guide-data.vue";
import compassYs from "./components/compass-ys.vue";

@Component({
  components: {
    customerData,
    guideData,
    compassYs
  }
})
export default class extends Vue {
  // data
  private active_name = "2"; // 默认展示
  private ENTERPRISE_MARK = localStorage.getItem("ENTERPRISE_MARK");
}
</script>
<style lang="scss" scoped>
/deep/ .table-content {
  margin-bottom: 0 !important;
}
</style>
