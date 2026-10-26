<template>
  <el-select
    class="search-select"
    :value="select_dp.name"
    placeholder="请选择所属部门"
    ref="selectBlur"
    size="medium"
    @clear="handleClear"
  >
    <el-option :value="select_dp.name" :label="select_dp.name" class="options">
      <el-tree
        :props="dp_props"
        :load="loadNode"
        node-key="dp_third_id"
        :check-strictly="true"
        :check-on-click-node="true"
        @node-click="handleNodeClick"
        highlight-current
        :expand-on-click-node="false"
        lazy
        ref="tree"
      ></el-tree>
    </el-option>
  </el-select>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { getThirdListApi } from "@/api/we-journey/Integral";

@Component
export default class SelectTreeDepartment extends Vue {
  private select_dp = { name: "", value: "" };

  /**  data数据 */
  private dp_props = {
    label: "dp_name",
    children: [],
    isLeaf: "leaf"
  };

  /** 选择部门 */
  private handleNodeClick(node: any) {
    this.select_dp.name = node.dp_name;
    this.select_dp.value = node.dp_third_id;
    // 部门选中之后关闭弹窗
    (this.$refs.selectBlur as any).blur();
    this.$emit("change", node.dp_third_id);
  }

  /** 动态加载节点 */
  private async loadNode(node: any, resolve: any) {
    if (node.level === 0) {
      const {
        data: { list }
      } = await getThirdListApi({ third_parent_id: 0 });
      list.forEach((item: any) => {
        if (!item.is_child_department) item.leaf = true;
      });

      const newList = [
        { dp_name: "全部", dp_id: -1, is_child_department: 0, leaf: true }
      ];
      resolve([...newList, ...list]);
    }
    if (node.level >= 1) {
      this.getIndex(node, resolve);
    }
  }

  /** 加载叶子节点数据函数 */
  private getIndex(node: any, resolve: any) {
    const params = { third_parent_id: node.data.dp_third_id };
    getThirdListApi(params).then(res => {
      this.setChecked(res.data.list);
      resolve(res.data.list);
    });
  }

  /** 判断部门是否有子节点 */
  private setChecked(arr: any) {
    arr.forEach((item: any) => {
      if (!item.is_child_department) item.leaf = true;
    });
  }

  /** 清空选中内容 */
  private handleClear() {
    this.select_dp.value = "";
    this.select_dp.name = "";
    this.$emit("change", "-1");
  }

  /** 父组件调用重置 */
  private reset() {
    this.select_dp = { name: "全部", value: "-1" };
    this.$emit("change", this.select_dp.value);
  }

  /** 初始化赋值 */
  private created() {
    this.reset();
  }
}
</script>
<style lang="scss" scoped>
.el-scrollbar .el-scrollbar__view .el-select-dropdown__item {
  height: auto;
  max-height: 274px;
  min-height: 32px;
  padding: 0;
  overflow: hidden;
  overflow-y: auto;
  width: 280px;
}
.el-select-dropdown__item.selected {
  font-weight: normal;
}
ul li >>> .el-tree .el-tree-node__content {
  height: auto;
  padding: 0 20px;
}
.el-tree-node__label {
  font-weight: normal;
}
.el-tree >>> .is-current .el-tree-node__label {
  color: #409eff;
  font-weight: 700;
}
.el-tree >>> .is-current .el-tree-node__children .el-tree-node__label {
  color: #606266;
  font-weight: normal;
}
.search-select {
  width: 280px;
}
</style>
