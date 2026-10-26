<template>
  <el-dialog 
    class="dialog"
    title="选择白名单"
    :visible.sync="visible"
    width="1000px"
    :append-to-body="false"
    :modal-append-to-body="false"
    @opened="handleOpened"
    @close="handleClose"
    >
    <div class="search">
      <search-item text="白名单名称：">
        <cy-input
          size="medium"
          placeholder="请输入白名单名称"
          v-model="searchParams.name"
        ></cy-input>
      </search-item>
      <el-button size="medium" type="primary" @click="handleSearch">查询</el-button>
      <el-button size="medium" plain @click="handleReset">重置</el-button>
      <el-button size="medium" type="primary" @click="handleAdd">新增白名单</el-button>
    </div>
    <el-table ref="table" :data="list" stripe row-key="id" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" reserve-selection></el-table-column>
      <el-table-column prop="id" label="ID" width="180"></el-table-column>
      <el-table-column prop="name" label="白名单名称" min-width="100"></el-table-column>
      <el-table-column prop="member_num" label="白名单人数" min-width="100"></el-table-column>
      <el-table-column prop="admin_name" label="创建人" min-width="100"></el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="200">
        <template slot-scope="scope">
          {{ scope.row.created_at | formatDate('yyyy/mm/dd hh:mm:ss') }}
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      :current-page="meta.page"
      :page-sizes="[10, 20, 30, 40]"
      :page-size="meta.size"
      class="page-right font-weight-400"
      layout="total, sizes, prev, pager, next, jumper"
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
import Api from '@/api/wisdom/communication';

const defaultSearchParams = {
  name: '',
}

@Component({
  model: {
    prop: 'value',
    event: 'change'
  }
})
export default class SelectWhitelist extends Vue {
  @Prop({ type: Array, default: () => [] }) value: any[]
  /** 是否显示 */
  visible = false
  /** 搜索条件 */
  searchParams = {...defaultSearchParams}
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

  _visible = false
  _selection: any[] = []



  @Watch('value', {deep: true})
  handleValueChange(val: any[]) {
    this.selection = val
  }
  
  activated() {
    this.visible = this._visible
    this._visible = false
  }
  deactivated() {
    this._visible = this.visible
    this.visible = false
    this._selection = JSON.parse(JSON.stringify(this.selection))
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
    const res: any = await Api.getWhitelistApi(params)
    if (res.code == 200) {
      this.list = res.data.list;
      this.meta.total = res.data.total;
    } else {
      this.$message.error(res.message);
    }
  }

  handleShow() {
    this.visible = true
  }

  handleOpened() {
    this.value.forEach(item => {
      (this.$refs.table as any).toggleRowSelection(item, true)
    })
    
    // if(this._selection) {
    //   this._selection.forEach(item => {
    //     (this.$refs.table as any).toggleRowSelection(item, true)
    //   })
    // }
    // this._selection = []
    this.handleSearch();
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
    this.value = this.selection
    this.$emit('change', this.selection);
    this.visible = false;
  }

  handleAdd() {
    this.$router.push({
      name: "myapp_wisdom_communication_whitelist_add",
    });
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