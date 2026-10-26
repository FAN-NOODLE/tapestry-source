<template>
  <div class="container" :key="key">
    <el-tabs v-model="active" type="card" @tab-click="handleTabClick">
      <el-tab-pane label="播种模版" :name="0">
      </el-tab-pane>
      <el-tab-pane label="播种数据" :name="1">
        <!-- <data-list :title="title" :template_id="template_id"></data-list> -->
      </el-tab-pane>
    </el-tabs>
    <temp-list @view-data="handleViewData"></temp-list>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Watch } from "vue-property-decorator";
import TempList from '../components/temp-list.vue'
import DataList from '../components/data-list.vue'

@Component({
  name: 'SowList',
  components: {
    DataList,
    TempList
  },
})
export default class SowList extends Vue {
  active = 0
  handleViewData(row:any) {
    this.$router.push({
      name: "myapp_customer_sow_data-list",
      query: {
        id: row.id,
        name: row.name
      }
    })
  }

  handleTabClick(val) {
    if(val.index === '1') {
      this.$router.push({
        name: "myapp_customer_sow_data-list",
      })
    }
  }
}
</script>

<style lang="scss" scoped>
.container {
  background-color: #fff;
  padding: 20px;
}
::v-deep .layout-list {
  padding: 0
}
</style>