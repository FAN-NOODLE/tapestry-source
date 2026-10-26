<template>
  <layout-list>
    <template slot="table">
      <el-tabs v-model="active" @tab-click="handleTabClick">
        <el-tab-pane label="VM陈列任务" name="1">
          <display-task-list @view="handleChange"></display-task-list>
        </el-tab-pane>
        <el-tab-pane label="陈列数据" name="2">
          <display-data-list :title="current" :current-id="current_id"></display-data-list>
        </el-tab-pane>
      </el-tabs>
    </template>
  </layout-list>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import DisplayDataList from "../components/display-data-list.vue";
import DisplayTaskList, {IItem} from "../components/display-task-list.vue";
@Component({
  name: 'DisplayTask',
  components: {
    DisplayDataList,
    DisplayTaskList
  }
})
export default class DisplayTask extends Vue {
  active = '1'
  current = ''
  current_id = ''

  handleTabClick() {
    this.current = ''
    this.current_id = ''
  }
  handleChange(row: IItem) {
    this.active = '2'
    this.current = row.display_name
    this.current_id = String(row.id)
  }

}
</script>