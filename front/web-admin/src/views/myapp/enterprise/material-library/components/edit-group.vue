<template>
  <div class="edit-group-model">
    <el-dialog
      title="素材分组"
      :visible.sync="show_dialog"
      width="1000px"
      @close="closeDialog"
    >
      <div class="button-box">
        <el-button
          :disabled="!chooseCateList.length"
          class="button"
          @click="editStore"
          >批量设置分组权限</el-button
        >
      </div>
      <el-table
        ref="multipleTable"
        :data="groupList"
        :height="400"
        style="width: 100%;margin-bottom: 20px;"
        row-key="id"
        stripe
        default-expand-all
        :select-on-indeterminate="true"
        @select-all="selectAll"
        @select="select"
        @selection-change="selectionChange"
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
      >
        <el-table-column type="selection" width="55"></el-table-column>
        <el-table-column prop="group_name" label="分类名称" min-width="200">
        </el-table-column>
        <el-table-column v-if="ep.isSW" prop="group_name" label="主页面可见" min-width="200">
          <template slot-scope="scope">
            <span @click="handleClick(scope.row)">
              <el-switch :value="scope.row.ignore_index===0"></el-switch>
            </span>
          </template>
        </el-table-column>
        <el-table-column label="归属门店" width="180">
          <template slot-scope="scope">
            <el-button type="text" @click="showStoreList(scope.row)"
              >查看</el-button
            >
          </template>
        </el-table-column>
        <el-table-column label="操作" width="330">
          <template slot-scope="scope">
            <el-button type="text" @click="addnewgroup(scope.row, '编辑分组')"
              >重命名</el-button
            >
            <el-button type="text" @click="delGroup(scope.row)">删除</el-button>
            <el-button type="text" @click="showStore(scope.row)"
              >归属门店</el-button
            >
            <el-button
              v-if="scope.row.level < 4"
              type="text"
              @click="addnewgroup(scope.row, '添加分组')"
              >添加分组</el-button
            >
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
    <add-edit-group
      v-if="pushOrEditGroup"
      :title="addtit"
      :filterObj="filterObj"
      :rowData="rowData"
      @closeDialog="closeaddgroup"
    ></add-edit-group>
    <select-store
      v-if="selectStoreModel"
      :groupIds="groupIds"
      @close="closeSelectStore"
      @submit="submitStore"
      row-key="dp_third_id"
    ></select-store>
    <store-list
      v-if="showStoreModel"
      :group="groupData"
      @close="closeStoreModel"
    ></store-list>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import Api from "@/api/enterprise/material-library";
import addEditGroup from "./add-edit-group.vue";
import SelectStore from "@/components/business/select-store/index.vue";
import StoreList from "./store-list.vue";
import { areAllGroupsSelected } from "../group-selection-state";
@Component({
  components: {
    addEditGroup,
    SelectStore,
    StoreList
  }
})
export default class extends Vue {
  @Prop() list?: Array<any>;
  //------------------------------------------------------data-----------------------------------------------------------/
  private show_dialog = true;
  private parentId = 0; // 弹框最大层级的ID
  private selectStoreModel = false;
  private groupList = [];
  private multipleSelection = [];
  private isAllSelect = false;
  private chooseCateList = [];
  // 是否打开添加修改分组弹框
  private pushOrEditGroup = false;
  private addtit = "添加一级分组";
  private rowData = {}; //分组每一项
  private showStoreModel = false;
  private groupData = {};
  private groupIds = [];
  private filterObj: any = {
    keyword: "",
    create_end: 0,
    create_start: 0,
    children: {},
    just_me: 0 //查询条件
  };

  get ep() {
    return this.$store.getters['enterprise']
  }
  //----------------------------------------------------methods-----------------------------------------------------------/

  private mounted() {
    this.parentId = (this.list as any).id || 0;
    if (this.parentId === 0) {
      this.getGroup();
      return;
    }
    this.groupList = [this.list];
  }
  // 关闭弹框
  private closeDialog() {
    this.$emit("closeDialog");
  }
  private showStoreList(item) {
    this.groupData = item;
    this.showStoreModel = true;
  }
  private closeStoreModel() {
    this.showStoreModel = false;
  }
  private showStore(item) {
    this.groupIds = [item.id];
    this.selectStoreModel = true;
  }
  private closeSelectStore() {
    this.chooseCateList = [];
    (this.$refs.multipleTable as any).clearSelection();
    this.selectStoreModel = false;
  }
  private submitStore(arr: any) {
    const params = {
      source_group_ids: this.groupIds,
      third_ids: arr.map(item => {
        return item.dp_third_id;
      })
    };
    Api.saveGroupStore(params).then((res: any) => {
      if (res.code == 200) {
        this.$message.success("更新成功");
        this.closeSelectStore();
      } else {
        this.$message.error(res.message);
      }
    });
  }
  private closeaddgroup() {
    this.pushOrEditGroup = false;
    this.getGroup();
  }
  private getGroup() {
    try {
      Api.getSourceGroupListApi({ id: this.parentId - 0 }).then(res => {
        this.groupList = res.data;
      });
    } catch (error) {
      console.log("重置数据22333", error);
    }
  }
  private delGroup(data: any) {
    this.filterObj.group_id = data.id;
    if (data.children !== null) {
      this.$message.error("该分组下有子级分组,不能删除！");
    } else {
      this.$confirm("你确认要删除该分组吗？", "删除分组", {
        confirmButtonText: "确定删除",
        cancelButtonText: "我再想想",
        type: "warning"
      }).then(() => {
        Api.sourceGroupDeleteApi({ id: data.id }).then((res: any) => {
          if (res.code == 200) {
            this.filterObj.group_id = -1;
            this.getGroup();
            this.$message({
              type: "success",
              message: "删除成功!"
            });
          } else {
            this.$message.error(res.message);
          }
        });
      });
    }
  }
  private editCommand(command: any) {
    if (command.id) {
      this.addtit = "修改分组";
      this.rowData = command;
      this.pushOrEditGroup = true;
      return;
    }
  }
  private addnewgroup(command: any, title) {
    if (command.id) {
      this.addtit = title;
      this.rowData = command;
      this.filterObj.group_id = command.id;
      this.pushOrEditGroup = true;
      return;
    }
  }
  private setChildren(children, type) {
    // 编辑多个子层级
    children.map(j => {
      this.toggleSelection(j, type);
      if (j.children) {
        this.setChildren(j.children, type);
      }
    });
  }
  private editStore() {
    this.groupIds = this.chooseCateList.map(item => {
      return item.id;
    });
    this.selectStoreModel = true;
  }
  private select(selection, row) {
    if (
      selection.some(el => {
        return row.id === el.id;
      })
    ) {
      if (row.children) {
        // row.children.map(j => {
        //     this.toggleSelection(j, true)
        // })
        // 解决子组件没有被勾选到
        this.setChildren(row.children, true);
      }
    } else {
      if (row.children) {
        // row.children.map(j => {
        //     this.toggleSelection(j, false)
        // })
        this.setChildren(row.children, false);
      }
    }
  }
  private toggleSelection(row, select) {
    if (row) {
      this.$nextTick(() => {
        (this.$refs.multipleTable as any).toggleRowSelection(row, select);
      });
    }
  }
  private selectionChange(rows) {
    this.chooseCateList = rows;
    this.$nextTick(() => {
      const table = this.$refs.multipleTable as any;
      if (table) {
        table.store.states.isAllSelected = areAllGroupsSelected(
          this.groupList,
          this.chooseCateList
        );
      }
    });
  }
  private selectAll(selection) {
    // tabledata第一层只要有在selection里面就是全选
    const isSelect = selection.some(el => {
      const tableDataIds = this.groupList.map(j => j.id);
      return tableDataIds.includes(el.id);
    });
    // tableDate第一层只要有不在selection里面就是全不选
    const isCancel = !this.groupList.every(el => {
      const selectIds = selection.map(j => j.id);
      return selectIds.includes(el.id);
    });
    console.log(isSelect, "isSelect");
    if (isSelect) {
      selection.map(el => {
        if (el.children) {
          // el.children.map(j => {
          //     this.toggleSelection(j, true)
          // })
          // 解决子组件没有被勾选到
          this.setChildren(el.children, true);
        }
      });
    }
    if (isCancel) {
      this.groupList.map(el => {
        if (el.children) {
          // el.children.map(j => {
          //     this.toggleSelection(j, false)
          // })
          // 解决子组件没有被勾选到
          this.setChildren(el.children, false);
        }
      });
    }
    this.$emit("handleSelect", this.groupList);
  }


  handleClick(row) {
    const params = {
      id: row.id,
      ignore_index: row.ignore_index ? 0 : 1
    };
    Api.updateIgnoreIndexApi(params).then((res: any) => {
      if (res.code == 200) {
        this.$message.success("更新成功");
        this.getGroup();
      } else {
        this.$message.error(res.message);
      }
    });
  }
}
</script>
<style lang="scss" scoped>
/deep/.el-dialog {
  max-height: 700px;
}
.button-box {
  height: 50px;
}
/deep/.el-dialog__body {
  max-height: 500px;
}
</style>
