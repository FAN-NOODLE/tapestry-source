<template>
  <el-dialog 
    class="dialog"
    title="关联部门传讯"
    :visible.sync="visible"
    width="1200px"
    :top="'10vh !important'"
    :append-to-body="false"
    :modal-append-to-body="false"
    @opened="handleOpened"
    @close="handleClose"
    >
    <div class="search">
      <search-item text="标题名称：">
        <cy-input
          size="medium"
          placeholder="请输入标题名称"
          v-model="searchParams.title"
        ></cy-input>
      </search-item>
      <search-item text="更新时间：">
        <el-date-picker
          size="medium"
          v-model="updateTime"
          type="datetimerange"
          store_type-separator="至"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
          value-format="timestamp"
          :default-time="['00:00:00', '23:59:59']"
        >
        </el-date-picker>
      </search-item>
      <!-- <search-item text="所属部门：">
        <el-select
          v-model="searchParams.department_name"
          size="medium"
          placeholder="请选择所属部门"
          filterable
          style="width: 280px;"
        >
          <el-option v-for="(item, index) in departmentOptions" :key="index" :label="item.label" :value="item.value"></el-option>
        </el-select>
      </search-item> -->
      <search-item text="传讯类型：">
        <cy-select
          v-model="searchParams.type"
          size="medium"
          placeholder="请选择传讯类型"
          :optionsList="typeOptions"
        ></cy-select>
      </search-item>
      <el-button size="medium" type="primary" @click="handleSearch">查询</el-button>
      <el-button size="medium" plain @click="handleReset">重置</el-button>
    </div>
    <el-table ref="table" :data="list" stripe row-key="id" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" reserve-selection :selectable="selectable"></el-table-column>
      <el-table-column prop="id" label="ID" width="180"></el-table-column>
      <el-table-column prop="title" label="标题名称" min-width="180" show-overflow-tooltip></el-table-column>
      <el-table-column prop="type" label="传讯类型" min-width="100">
        <template slot-scope="scope">
          {{ typeName(scope.row.type) }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="传讯状态" width="100">
        <el-tag size="small" type="success">已通过</el-tag>
      </el-table-column>
      <el-table-column prop="rel_status" label="发布状态" width="100">
        <template slot-scope="scope">
          <el-tag v-if="scope.row.rel_status === 0 || !selectable(scope.row)" size="small">已关联</el-tag>
          <el-tag v-else size="small" type="success">未发布</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="department_name" label="所属部门" min-width="120">
        <template slot-scope="scope">
          <span>{{scope.row.department_name || '--'}}</span>
        </template>
      </el-table-column>
      <el-table-column prop="update_at" label="最近更新时间" width="200">
        <template slot-scope="scope">
          {{ scope.row.update_at | formatDate('yyyy/mm/dd hh:mm:ss')}}
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
import { typeOptions } from '../../const';

const defaultSearchParams = {
  title: '',
  updated_start_time: 0,
  updated_end_time: 0,
  status: 4,
  department_name: '',
  type: -1,
  is_all: 1,
  no_release: 1 // 过滤已关联并发布的传序
}

@Component({
  model: {
    prop: 'value',
    event: 'change'
  }
})
export default class SelectContent extends Vue {
  @Prop({ type: Array, default: () => [] }) value: []
  @Prop({ type: Array, default: () => [] }) disabledList: number[]
  /** 是否显示 */
  visible = false
  department = ''
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
  /** 部门列表 */
  departmentOptions: any[] = []
  /** 类型列表 */
  typeOptions = [
    { label: '全部', value: -1 },
    ...typeOptions
  ]
  /** 选中的数据 */
  selection: any[] = []

    /** 创建时间 读取 */
  get updateTime() {
    const { updated_start_time, updated_end_time } = this.searchParams
    return [updated_start_time*1000, updated_end_time*1000].filter(Boolean)
  }

  /** 创建时间 写入 */
  set updateTime(value: number[]) {
    const [start_time, end_time] = value || []
    this.searchParams.updated_start_time = parseInt(start_time/1000+'');
    this.searchParams.updated_end_time = parseInt(end_time/1000+'');
  }

  get typeName() {
    return (value) => {
      return typeOptions.find(item => item.value === value).label
    }
  }

  @Watch('value')
  handleValueChange(val: any[]) {
    this.selection = JSON.parse(JSON.stringify(val));
  }

  selectable(row:any) {
    return row.rel_status === -1 && !this.disabledList.includes(row.id)
  }
  /**
   * 获取列表
   */
  async fetchData() {
    const params = {
      page: this.meta.page,
      size: this.meta.size,
      ...this.searchParams,
      department_name: this.department
    }
    const res: any = await Api.getContentListApi(params)
    if (res.code == 200) {
      // 设置列表选中项
      this.list = res.data.list;
      this.meta.total = res.data.total;

    } else {
      this.$message.error(res.message);
    }
  }

  /** 
   * 所属部门
   */
  async getDepartmentList() {
    const res: any = await Api.getDepartmentListApi({})
    if(res.code === 200) {
      this.departmentOptions = [
        { label: '全部', value: '' },
        ...res.data.map((item: any) => {
          return {
            label: item.name,
            value: item.name
          }
        })
      ]
    } else {
      this.$message.error(res.message)
    }
  }

  handleShow(department: string) {
    this.department = department
    this.searchParams = { ...defaultSearchParams }
    this.visible = true;
    this.getDepartmentList()
  }

  handleOpened() {
    this.value.forEach(item => {
      (this.$refs.table as any).toggleRowSelection(item, true)
    })
    this.handleSearch()
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