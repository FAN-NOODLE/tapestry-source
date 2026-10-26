<template>
  <el-table :data="list" :height="operationBtn ? 427 : 340" stripe>
    <el-table-column prop="id" label="活动ID" min-width="150">
    </el-table-column>
    <el-table-column
      label="活动名称"
      min-width="300"
      prop="template_name"
      :show-overflow-tooltip="true"
    >
    </el-table-column>
    <el-table-column label="操作" width="100" v-if="operationBtn">
      <template slot-scope="scope">
        <el-button @click="delPageFn(scope.$index)" type="text" size="small"
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
  @Prop({ default: true }) operationBtn!: boolean;
  delPageFn(ind: number) {
    this.list.splice(ind, 1);
    this.$emit("getPageChange", this.list);
  }
}
</script>
<style lang="scss">
.el-table--scrollable-y {
  .el-table__body-wrapper::-webkit-scrollbar {
    width: 6px !important;
  }
}
</style>
<style lang="scss" scoped>
.el-range-editor {
  width: 280px !important;
}
.el-button--text {
  color: #f67170;
}
.el-button--text:hover {
  color: #f56c6c;
}
.input-with-select {
  .el-select {
    width: 100px !important;
  }
}
.marginL10 {
  margin-left: 10px;
}
</style>
