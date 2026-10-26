<template>
  <div class="display-task-list">
    <!-- 按钮 -->
    <el-row class="buttons">
      <el-button type="primary" @click="handleAdd">新建任务</el-button>
    </el-row>
    <!-- 筛选 -->
    <el-row class="search">
      <search-item text="标题：">
        <cy-input
          size="medium"
          placeholder="请输入标题名称"
          v-model="searchParams.display_name"
          clearable
        ></cy-input>
      </search-item>
      <search-item text="执行时间：">
        <el-date-picker
          class="cy-el-date-picker"
          v-model="executeTime"
          type="datetimerange"
          store_type-separator="至"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
          value-format="timestamp"
          :default-time="['00:00:00', '23:59:59']"
        >
        </el-date-picker>
      </search-item>
      <search-item text="创建时间：">
        <el-date-picker
          class="cy-el-date-picker"
          v-model="createTime"
          type="datetimerange"
          store_type-separator="至"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
          value-format="timestamp"
          :default-time="['00:00:00', '23:59:59']"
        >
        </el-date-picker>
      </search-item>
      <search-item text="状态：">
        <cy-select
          class="search-select"
          size="medium"
          v-model="searchParams.status"
        >
          <cy-option
            v-for="item in statusOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          ></cy-option>
        </cy-select>
      </search-item>
      <search-item text="是否上架：">
        <cy-select
          class="search-select"
          size="medium"
          v-model="searchParams.show_type"
        >
          <cy-option
            v-for="item in stateOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          ></cy-option>
        </cy-select>
      </search-item>
      <el-button class="m-b-20" size="medium" type="primary" @click="handleSearch">查询</el-button>
      <el-button class="m-b-20" size="medium" plain @click="handleReset">重置</el-button>
    </el-row>
    <!-- 列表 -->
    <el-row class="table">
      <el-table :data="list" stripe>
        <el-table-column prop="id" label="ID" width="200"></el-table-column>
        <el-table-column prop="display_name" label="标题名称" min-width="200"></el-table-column>
        <el-table-column prop="store_type" label="通知范围" min-width="100">
          <template slot-scope="{row}">
            <span>{{ ['全部门店','指定门店'][row.store_type] }}</span>
          </template>
        </el-table-column>
        <el-table-column label="执行时间" width="250">
          <template slot-scope="{row}">
            <div v-if="row.start_time">
              <span>{{ row.start_time | formatDate('yyyy/mm/dd') }}</span>
              <span> ~ </span>
              <span>{{ row.end_time | formatDate('yyyy/mm/dd') }}</span>
            </div>
            <span v-else>--</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" min-width="100">
          <template slot-scope="{row}">
            <span>{{ status(row.status) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="show_type" label="是否上架" min-width="100">
          <template slot-scope="{row}">
            <span @click="handleChangeStatus(row)">
              <el-switch 
                :value="row.show_type" 
                :active-value="1"
                :inactive-value="2"
                :disabled="!row.start_time"
              ></el-switch>
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="200">
          <template slot-scope="{row}">
            <span>{{ row.created_at | formatDate('yyyy/mm/dd hh:mm:ss') }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="260" fixed="right">
          <template slot-scope="{row}">
            <el-button type="text" size="small" @click="handleViewDetail(row)">查看</el-button>
            <el-button type="text" size="small" @click="handleEdit(row)" :disabled="row.status !== 0">编辑</el-button>
            <el-button type="text" size="small" @click="handleViewData(row)">查看数据</el-button>
            <el-button type="text" size="small" @click="handleCopy(row)">复制</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-row>
    <!-- 分页 -->
    <el-row class="footer">
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
    </el-row>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Api from '@/api/wisdom/displayTask';
/** 查询条件 */
interface ISearchParams {
  /** 任务标题 */
  display_name: string;
  /** 执行时间 开始 */
  start_time: number;
  /** 执行时间 结束 */
  end_time: number;
  /** 创建时间 开始 */
  created_start_at: number;
  /** 创建时间 结束 */
  created_end_at: number;
  /** 状态，-1全部，0未开始，1进行中，2已结束 */
  status: number;
  /** 是否上架，-1全部，0下架，1上架 */
  show_type: number;
}
/** 下拉选项 */
interface ISelectItem {
  label: string;
  value: number;
}
/** 分页 */
interface IMate {
  /** 页数 */
  page: number;
  /** 每页数量 */
  size: number;
  /** 总数 */
  total: number;
}
/** 数据 */
export interface IItem {
  /** ID */
  id: number;
  /** 标题名称 */
  display_name: string;
  /** 通知范围 通知门店类型，0全部门店,1选择门店 */
  store_type: number;
  /** 执行时间 开始时间 */
  start_time: number;
  /** 执行时间 结束时间 */
  end_time: number;
  /** 任务状态：0-未开始，1-进行中，2-已结束 */
  status: number;
  /** 上下架状态，0-保存，1-保存并上架 */
  show_type: number;
  /** 创建时间 */
  created_at: number;
}
const defaultSearchParams: ISearchParams = {
  display_name: '',
  start_time: 0,
  end_time: 0,
  created_start_at: 0,
  created_end_at: 0,
  show_type: -1,
  status: -1
}
@Component
export default class DisplayTaskList extends Vue {
  /** 状态选项 */
  private statusOptions: ISelectItem[] = [
    { label: '全部', value: -1 },
    { label: '进行中', value: 1 },
    { label: '已结束', value: 2 },
    { label: '未开始', value: 0 },
  ]
  /** 上下架选项 */
  private stateOptions: ISelectItem[] = [
    { label: '全部', value: -1 },
    { label: '是', value: 1 },
    { label: '否', value: 0 },
  ]
  /** 筛选数据 */ 
  private searchParams: ISearchParams = {...defaultSearchParams}
  /** 分页 */
  private meta: IMate = {
    page: 1,
    size: 10,
    total: 0
  }
  /** 列表数据 */
  private list: IItem[] = []
  /** 状态文案 */
  get status() {
    return (value: number) => {
      return this.statusOptions.find(item => item.value === value)?.label
    }
  }

  /** 执行时间 读取 */
  get executeTime() {
    const { start_time, end_time } = this.searchParams
    return [start_time*1000, end_time*1000].filter(Boolean)
  }
  /** 执行时间 写入 */
  set executeTime(value: number[]) {
    const [start_time, end_time] = value || []
    this.searchParams.start_time = parseInt(start_time/1000+'');
    this.searchParams.end_time = parseInt(end_time/1000+'');
  }

  /** 创建时间 读取 */
  get createTime() {
    const { created_start_at, created_end_at } = this.searchParams
    return [created_start_at*1000, created_end_at*1000].filter(Boolean)
  }
  /** 创建时间 写入 */
  set createTime(value: number[]) {
    const [start_time, end_time] = value || []
    this.searchParams.created_start_at = parseInt(start_time/1000+'');
    this.searchParams.created_end_at = parseInt(end_time/1000+'');
  }

  activated() {
    this.fetchData()
  }
  
  /**
   * @description 查询
   */
  private handleSearch() {
    this.meta.page = 1
    this.fetchData()
  }
  
  /**
   * @description 重置
   */
  private handleReset() {
    this.searchParams = { ...defaultSearchParams }
    this.handleSearch()
  }

  /**
   * 每页数量变化
   * @param val 
   */
  private handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.fetchData();
  }

  /**
   * 当前页
   * @param val 
   */
  private handleCurrentChange(val: number) {
    this.meta.page = val;
    this.fetchData();
  }
  /** 
   * @description 修改状态
   */
  private handleChangeStatus(row: IItem) {
    const isCopy = !row.start_time
    if(isCopy) return
    const show_type = row.show_type === 1 ? 0 : 1
    const message = {
      0: '关闭后该任务将不在SA手机端展示，确认关闭？',
      1: '开启后该任务将在SA手机端展示，确认开启？',
    }[show_type]
    const title = ['下架提醒', '上架提醒'][show_type]
    this.$confirm(message, title, {
      confirmButtonText: "确认",
      cancelButtonText: "取消"
    })
      .then(() => {
        Api.updateDisplayTaskStatusApi({id: row.id})
        .then((res: any) => {
          if (res.code == 200) {
            row.show_type = show_type
          } else {
            this.$message.error(res.message);
          }
        })
      })
      .catch((err: any) => {
        console.log(err);
      });
  }

  /** 
   * @description 查看
   */
  private handleViewDetail(row: IItem) {
    console.log(row);
    this.$router.push({
      name: "myapp_wisdom_display-task_detail",
      query: { id: String(row.id) }
    });
  }

  /** 
   * @description 编辑
   */
  private handleEdit(row: IItem) {
    this.$router.push({
      name: "myapp_wisdom_display-task_edit",
      query: { id: String(row.id) }
    });
  }

  /** 
   * @description 查看数据
   */
  private handleViewData(row: IItem) {
    this.$emit("view", row);
  }

  /** 
   * @description 复制
   */
  private handleCopy(row: IItem) {
    this.$confirm('确认复制当前任务？复制后请重新选择任务时间。', '复制提醒', {
      confirmButtonText: "确认",
      cancelButtonText: "取消"
    })
      .then(() => {
        Api.copyDisplayTaskApi({id: row.id})
          .then((res: any) => {
            if (res.code == 200) {
              this.fetchData()
              this.$message.success('复制成功')
            } else {
              this.$message.error(res.message);
            }
          });
      })
      .catch((err: any) => {
        console.log(err);
      });
    
  }

  /** 
   * @description 添加
   */
  private handleAdd() {
    this.$router.push({
      name: "myapp_wisdom_display-task_add"
    });
  }
  
  /**
   * @description 获取数据
   */
  private fetchData() {
    const params = {
      page: this.meta.page,
      size: this.meta.size,
      ...this.searchParams,
    }
    Api.getDisplayTaskListApi(params).then((res: any) => {
      if (res.code == 200) {
        this.list = res.data.list;
        this.meta.total = res.data.total;
      } else {
        this.$message.error(res.message);
      }
    });
  }

  mounted() {
    this.fetchData()
  }
}
</script>

<style lang="scss" scoped> 
/deep/.cy-select.search-select .el-select .el-input--medium, 
/deep/.cy-select.search-select .el-select .el-input--medium .el-input__inner, 
/deep/.cy-select.search-select .el-select .el-input--small, 
/deep/.cy-select.search-select .el-select .el-input--small .el-input__inner {
  height: 40px !important;
  line-height: 40px !important;
}
.display-task-list {
  .buttons {
    margin-bottom: 20px;
  }
  .search {
    padding: 20px 20px 0px 20px;
    margin-bottom: 20px;
    background-color: #f7f7f7;
    border-radius: 4px;
  }

  .footer {
    text-align: right;
    margin-top: 20px;
  }
}
</style>