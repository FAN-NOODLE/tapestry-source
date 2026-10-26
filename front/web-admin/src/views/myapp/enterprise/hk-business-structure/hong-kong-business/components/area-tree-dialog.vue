<template>
  <el-dialog
    title="选择管理范围"
    :visible.sync="visible"
    width="760px"
    :close-on-click-modal="false"
    @opened="handleOpened"
    @close="handleCancel"
  >
    <div class="area-tree-dialog">
      <div class="area-tree-dialog__tree">
        <el-tree
          ref="tree"
          :data="treeList"
          :props="defaultProps"
          :default-expand-all="true"
          :show-checkbox="true"
          node-key="id"
          :check-strictly="false"
          :expand-on-click-node="true"
          @check-change="handleCheckChange"
        >
          <span slot-scope="{ data }" class="custom-tree-node">
            <svg-icon
              v-if="data.type === 4"
              class-name="m-r-5"
              icon-class="store"
            ></svg-icon>
            <svg-icon
              v-else
              class-name="m-r-5"
              icon-class="organization"
            ></svg-icon>
            <span>{{ getNodeName(data) }}</span>
          </span>
        </el-tree>
      </div>
      <div class="area-tree-dialog__selected">
        <div class="selected-title">已选管理范围</div>
        <div
          v-for="item in normalizedCheckedList"
          :key="item.id"
          class="selected-item"
        >
          <span class="selected-item__name">{{ getNodeName(item) }}</span>
          <el-button type="text" size="mini" @click="handleDel(item)">
            <i class="el-icon-delete"></i>
          </el-button>
        </div>
      </div>
    </div>
    <template slot="footer">
      <el-button @click="handleCancel">取消</el-button>
      <el-button type="primary" @click="handleConfirm">确定</el-button>
    </template>
  </el-dialog>
</template>

<script lang="ts">
import { Component, Prop, Vue, Watch } from "vue-property-decorator";

@Component
export default class AreaTreeDialog extends Vue {
  @Prop({ default: false }) value!: boolean;
  @Prop({ default: () => [] }) selectedIds!: number[];
  @Prop({ default: () => [] }) treeList!: any[];

  visible = false;
  checkedList: any[] = [];
  defaultProps = {
    children: "children",
    label: "name"
  };

  get normalizedCheckedList() {
    return this.normalizeNodes(this.checkedList);
  }

  @Watch("value")
  handleVisibleChange(value: boolean) {
    this.visible = value;
  }

  handleOpened() {
    (this.$refs.tree as any).setCheckedKeys(this.selectedIds || []);
    this.handleCheckChange();
  }

  handleCheckChange() {
    this.checkedList = (this.$refs.tree as any).getCheckedNodes();
  }

  handleDel(item: any) {
    (this.$refs.tree as any).setChecked(item, false, true);
    this.handleCheckChange();
  }

  handleConfirm() {
    const ids = this.normalizedCheckedList.map((item: any) => item.id);
    this.$emit("confirm", ids);
    this.$emit("input", false);
  }

  handleCancel() {
    this.$emit("input", false);
  }

  getNodeName(node: any) {
    return node.store_name || node.name || "--";
  }

  normalizeNodes(nodes: any[]) {
    const selectedIds = new Set(nodes.map(item => item.id));
    return nodes.filter(item => {
      let parentId = item.parent_id;
      while (parentId) {
        if (selectedIds.has(parentId)) return false;
        const parent = this.findNode(parentId);
        parentId = parent ? parent.parent_id : 0;
      }
      return true;
    });
  }

  findNode(id: number) {
    const stack = [...this.treeList];
    while (stack.length) {
      const node = stack.shift();
      if (node.id === id) return node;
      stack.push(...(node.children || []));
    }
    return null;
  }
}
</script>

<style lang="scss" scoped>
.area-tree-dialog {
  display: flex;
  height: 460px;
}

.area-tree-dialog__tree {
  flex: 1;
  overflow-y: auto;
  padding-right: 16px;
}

.area-tree-dialog__selected {
  width: 280px;
  overflow-y: auto;
  padding-left: 16px;
  border-left: 1px solid #e6e6e6;
}

.selected-title {
  margin-bottom: 10px;
  color: $color-text-normal;
}

.selected-item {
  display: flex;
  align-items: center;
  padding: 4px 0;

  &:hover {
    background-color: #f5f7fa;
  }
}

.selected-item__name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.custom-tree-node {
  display: inline-flex;
  align-items: center;
}
</style>
