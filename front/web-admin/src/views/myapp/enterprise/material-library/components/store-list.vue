<template>
  <div class="edit-group-model">
    <el-dialog
      :title="'分组名称：' + group.group_name"
      :visible.sync="show_dialog"
      width="400px"
      @close="closeDialog"
    >
      <el-table
        ref="multipleTable"
        :data="list"
        :height="400"
        style="width: 100%;margin-bottom: 20px;"
        row-key="id"
        stripe
      >
        <el-table-column
          prop="department_name"
          label="门店名称"
          min-width="200"
        >
        </el-table-column>
      </el-table>
      <div class="footer-wrapper" v-if="list.length > 0 && false">
        <el-pagination
          @size-change="sizeChange"
          @current-change="currentChange"
          :current-page.sync="page"
          :page-sizes="[10, 20, 30, 40]"
          :page-size="10"
          layout="prev, pager, next"
          :total="total"
        >
        </el-pagination>
      </div>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import Api from "@/api/enterprise/material-library";
@Component
export default class extends Vue {
  @Prop() group?: any;
  //------------------------------------------------------data-----------------------------------------------------------/
  private show_dialog = true;
  private list = [];
  private page = 1;
  private size = 10;
  private total = 0;
  //----------------------------------------------------methods-----------------------------------------------------------/
  // 关闭弹框
  private closeDialog() {
    this.$emit("close");
  }
  private getList() {
    Api.getGroupStoreData({
      source_group_id: this.group.id,
      page: this.page,
      size: this.size
    }).then(res => {
      this.list = res.data.list;
      // this.total = res.data.total;
    });
  }
  private sizeChange(val: number) {
    this.size = val;
    this.page = 1;
    this.getList();
  }
  private currentChange(val: number) {
    this.page = val;
    this.getList();
  }
  private mounted() {
    this.getList();
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
.footer-wrapper {
  text-align: right;
}
</style>
