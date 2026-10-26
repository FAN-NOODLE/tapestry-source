<template>
  <layout-list>
    <template slot="table"> 
      <el-tabs 
        v-model="active"
        @tab-click="handleTabClick"
        type="card"
      >
        <el-tab-pane label="活动列表" name="1">
          <ActivityList></ActivityList>
        </el-tab-pane>
        <el-tab-pane label="投放数据" name="2">
          <ActivityData></ActivityData>
        </el-tab-pane>
      </el-tabs>
    </template>
  </layout-list>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import ActivityList from '../components/ActivityList.vue'
import ActivityData from '../components/ActivityData.vue'

@Component({
  components: {
    ActivityList,
    ActivityData
  }
})
export default class LaunchList extends Vue {
  active = '1'
  handleTabClick(e) {
    this.active = e.name;
  }

  mounted() {
    this.active = String(this.$route.query.active || '1');
  }
}
</script>

<style lang="scss" scoped>
::v-deep.layout-list .search-content {
  padding-bottom: 0 !important;
}
/deep/ .el-tabs__nav {
  .is-active {
    color: #fff;
    background-color: #000000;
  }
}
</style>