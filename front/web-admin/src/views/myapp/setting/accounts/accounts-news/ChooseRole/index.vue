<template>
  <el-dialog 
    class="dialog"
    title="选择角色"
    :visible.sync="visible"
    width="1000px"
    :top="'10vh !important'"
    :append-to-body="false"
    :modal-append-to-body="false"
    @opened="handleOpened"
    @close="handleClose"
    >
    <div class="search">
      <search-item text="角色名称：">
        <cy-input
          size="small"
          placeholder="请输入角色名称"
          v-model="searchParams.role_name"
        ></cy-input>
      </search-item>
      <el-button size="small" type="primary" @click="handleSearch">查询</el-button>
      <el-button size="small" plain @click="handleReset">重置</el-button>
    </div>
    <el-table ref="table" :data="list" stripe row-key="id" @selection-change="handleSelectionChange">
      <el-table-column :selectable="checkSelectable" type="selection" width="55" reserve-selection></el-table-column>
      <el-table-column label="角色名称" prop="name"></el-table-column>
      <el-table-column label="适用后台" prop="app_code_name"></el-table-column>
    </el-table>
    <el-pagination
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      :current-page="meta.page"
      :page-sizes="[10, 20, 30, 40]"
      :page-size="meta.size"
      class="page-right font-weight-400"
      layout="total, prev, pager, next, jumper"
      :total="meta.total"
    ></el-pagination>

    <div slot="footer">
      <el-button @click="handleClose">取 消</el-button>
      <el-button type="primary" @click="handleConfirm">确 定</el-button>
    </div>
  </el-dialog>
</template>

<script lang="ts">
import { Component, Prop, Vue, Watch } from 'vue-property-decorator';
import BaseData from "../base-data";
import Api from "@/api/setting/Account";
import { roleAppName } from "../../roles/app-options";

const defaultSearchParams = {
  role_name: '',
}

@Component({
  mixins: [BaseData],
  model: {
    prop: 'value',
    event: 'change'
  }
})
export default class SelectContent extends Vue {
  @Prop({ type: Array, default: () => [] }) value: []
  /** 是否显示 */
  visible = false
  /** 搜索条件 */
  searchParams: {role_name: string} = {...defaultSearchParams}
  /** 分页 */
  meta = {
    page: 1,
    size: 10,
    total: 0
  }
  /** 列表数据 */
  list = []

  /** 选中的数据 */
  selection: any[] = []


  @Watch('value')
  handleValueChange(val: any[]) {
    this.selection = JSON.parse(JSON.stringify(val));
  }

  /**
   * 获取列表
   */
  async fetchData() {
    const params = {
      page: this.meta.page,
      size: this.meta.size,
      ...this.searchParams,
    }
    const res: any = await Api.roleList(params)
    if (res.code == 200) {
      // 设置列表选中项
      this.list = res.data.list.map((item: any) => ({
        ...item,
        app_code_name: roleAppName(item.app_code),
      }));
      this.meta.total = res.data.total;

    } else {
      this.$message.error(res.message);
    }
  }

  handleShow() {
    this.visible = true;
  }

  handleOpened() {
    this.value.forEach(item => {
      (this.$refs.table as any).toggleRowSelection(item, true)
    })
    this.handleSearch()
  }

  checkSelectable(row: any) {
    return !this["inner_no_choice"].includes(row.mark);
  }

  handleSelectionChange(selection) {
    this.selection = selection
  }

  /**
   * @description 查询
   */
   handleSearch() {
    this.meta.page = 1
    this.fetchData()
  }
  
  /**
   * @description 重置
   */
  handleReset() {
    this.searchParams = { ...defaultSearchParams }
    this.handleSearch()
  }

  /**
   * 每页数量变化
   * @param val 
   */
  handleSizeChange(val: number) {
    this.meta.size = val;
    this.meta.page = 1;
    this.fetchData();
  }

  /**
   * 当前页
   * @param val 
   */
  handleCurrentChange(val: number) {
    this.meta.page = val;
    this.fetchData();
  }

  /** 
   * 取消
   */
  handleClose() {
    this.visible = false;
    (this.$refs.table as any).clearSelection();
  }

  /**
   * 确定
   */
  handleConfirm() {
    this.$emit('change', this.selection);
    this.visible = false;
  }
}
</script>

<style lang="scss" scoped>
#crs-scrm .dialog ::v-deep .el-dialog {
  max-height: none !important;
  margin-top:  10vh !important;
  transform: none;
  top: 0;
}


.el-pagination {
  margin-top: 10px;
  text-align: right;
}
</style>
