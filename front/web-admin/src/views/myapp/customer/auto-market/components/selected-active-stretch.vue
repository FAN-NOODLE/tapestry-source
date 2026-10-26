<template>
  <el-table :data="list" height="427">
    <el-table-column
      label="活动名称"
      min-width="300"
      prop="name"
      stripe
      :show-overflow-tooltip="true"
    >
    </el-table-column>
    <el-table-column label="操作" width="100">
      <template slot-scope="scope">
        <el-button @click="delFn(scope.$index)" type="text" size="small"
          >删除</el-button
        >
      </template>
    </el-table-column>
  </el-table>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
@Component({})
export default class extends Vue {
  /************data**************/
  @Prop() list!: Array<any>;
  delFn(ind: number) {
    this.list.splice(ind, 1);
    this.$emit("getActiveChange", this.list);
  }
}
</script>
<style lang="scss" scoped>
.el-range-editor {
  width: 280px !important;
}
.input-with-select {
  .el-select {
    width: 100px !important;
  }
}
.el-button--text {
  color: #f67170;
}
.el-button--text:hover {
  color: #f56c6c;
}
.marginL10 {
  margin-left: 10px;
}

.tableCell {
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}
table ::-webkit-scrollbar {
  width: 6px;
}
</style>
