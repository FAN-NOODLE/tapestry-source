<template>
  <div class="container">
    <el-tabs v-model="active" type="card" @tab-click="handleTabClick">
      <el-tab-pane label="播种模版" :name="0"></el-tab-pane>
      <el-tab-pane label="播种数据" :name="1"></el-tab-pane>
    </el-tabs>
    <data-list></data-list>
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
  active = 1

  handleTabClick(val) {
    if(val.index === '0') {
      this.$router.push({
        name: "myapp_customer_sow_list",
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