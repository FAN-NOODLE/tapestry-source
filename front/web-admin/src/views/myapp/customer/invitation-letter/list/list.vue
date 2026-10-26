<template>
  <div class="welcome_container">
    <layout-list>
      <!-- tab切换 -->
      <template slot="top-tabs">
        <el-tabs v-model="activeName" type="card">
          <el-tab-pane label="邀请函列表" name="first"> </el-tab-pane>
          <el-tab-pane label="邀请数据" name="second"> </el-tab-pane>
        </el-tabs>
      </template>
      <!-- 列表 -->
      <template slot="table">
        <invitation-list
          v-if="activeName === 'first'"
          @handleToViewData="handleToViewData"
        ></invitation-list>
        <invitation-data
          v-if="activeName === 'second'"
          :emitData="emitData"
        ></invitation-data>
      </template>
    </layout-list>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import invitationList from "../components/invite-list.vue";
import invitationData from "../components/invite-data.vue";
@Component({
  components: {
    invitationList,
    invitationData
  }
})
export default class extends Vue {
  /** data START */
  private activeName = "first";

  private emitData = {};

  /** data END */

  /** methods */
  private handleToViewData(data: any) {
    this.activeName = "second";
    this.emitData = data;
  }
}
</script>
<style lang="scss" scoped>
.welcome_container {
  background-color: #fff;
  // padding: 20px;
  border-radius: 10px;
}

/deep/ .el-tabs__nav {
  .is-active {
    color: #fff;
    background-color: #000000;
  }
}
::v-deep {
  .layout-list .table-content {
    margin-bottom: 0;
  }
}
</style>
