<template>
  <layout-list>
    <template slot="search">
      <search-item text="模版ID：">
        <cy-input
          size="medium"
          placeholder="请输入模版ID"
          v-model="searchParams.template_id"
          clearable
        ></cy-input>
      </search-item>
      <search-item text="播种模版名称：">
        <cy-input
          size="medium"
          placeholder="请输入播种模版名称"
          v-model="searchParams.template_name"
          clearable
        ></cy-input>
      </search-item>
      <search-item text="播种类型：">
        <cy-select
          size="medium"
          placeholder="请选择播种类型"
          v-model="searchParams.template_type"
          clearable
        >
          <cy-option label="全部" :value="0" ></cy-option>
          <cy-option label="默认播种" :value="1" ></cy-option>
          <cy-option label="播种模版" :value="2" ></cy-option>
        </cy-select>
      </search-item>
      <search-item text="导购姓名：">
        <cy-input
          size="medium"
          placeholder="请输入导购姓名"
          v-model="searchParams.sa_name"
          clearable
        ></cy-input>
      </search-item>
      <search-item text="顾客姓名：">
        <cy-input
          size="medium"
          placeholder="请输入顾客姓名"
          v-model="searchParams.username"
          clearable
        ></cy-input>
      </search-item>
      <search-item text="会员卡号：">
        <cy-input
          size="medium"
          placeholder="请输入会员卡号"
          v-model="searchParams.union_code"
          clearable
        ></cy-input>
      </search-item>
      <search-item text="播种状态：">
        <cy-select
          size="medium"
          placeholder="请选择播种状态"
          v-model="searchParams.sow_status"
          clearable
        >
          <cy-option label="全部" :value="0" ></cy-option>
          <cy-option label="待完成" :value="1" ></cy-option>
          <cy-option label="已完成" :value="2" ></cy-option>
        </cy-select>
      </search-item>
      <search-item text="创建时间：">
        <el-date-picker
          class="cy-el-date-picker"
          v-model="createdTime"
          type="datetimerange"
          store_type-separator="至"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
          value-format="timestamp"
          :default-time="['00:00:00', '23:59:59']"
        >
        </el-date-picker>
      </search-item>
      <search-item text="完成时间：">
        <el-date-picker
          class="cy-el-date-picker"
          v-model="noticeTime"
          type="datetimerange"
          store_type-separator="至"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
          value-format="timestamp"
          :default-time="['00:00:00', '23:59:59']"
        >
        </el-date-picker>
      </search-item>
      <el-button class="m-b-20" size="small" type="primary" @click="handleSearch">查询</el-button>
      <el-button class="m-b-20" size="small" plain @click="handleReset">重置</el-button>
      <el-button class="m-b-20" size="small" plain @click="handleExport">导出数据</el-button>
    </template>
    <template slot="table">
      <el-table :data="list" stripe>
        <el-table-column prop="template_type" label="播种类型" min-width="100">
          <template slot-scope="scope">
            <span>{{ scope.row.template_id?'播种模版':'默认播种' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="template_name" label="播种模版名称" min-width="200" show-overflow-tooltip>
          <template slot-scope="scope">
            <span>{{ scope.row.template_name||'--' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="template_id" label="模版ID" width="200">
          <template slot-scope="scope">
            <span>{{ scope.row.template_id||'--' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="sa_name" label="导购姓名" min-width="100">
          <template slot-scope="scope">
            <span>{{ scope.row.sa_name||'--' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="sa_dp_name" label="所在部门" min-width="100">
          <template slot-scope="scope">
            <span>{{ scope.row.sa_dp_name||'--' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="username" label="顾客姓名" min-width="120">
          <template slot-scope="scope">
            <span>{{ scope.row.username||'--' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="union_code" label="会员卡号" width="180">
          <template slot-scope="scope">
            <span>{{ scope.row.union_code||'--' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="播种状态" min-width="100">
          <template slot-scope="scope">
            <el-tag v-if="scope.row.status === 2" type="success">已完成</el-tag>
            <el-tag v-else type="info">待完成</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="200">
          <template slot-scope="scope">
            {{ scope.row.created_at | formatDate('yyyy/mm/dd hh:mm:ss') }}
          </template>
        </el-table-column>
        <el-table-column prop="remind_time" label="提醒时间" width="200">
          <template slot-scope="scope">
            {{ scope.row.remind_time | formatDate('yyyy/mm/dd hh:mm:ss') }}
          </template>
        </el-table-column>
        <el-table-column prop="finished_at" label="完成时间" width="200">
          <template slot-scope="scope">
            {{ scope.row.finished_at | formatDate('yyyy/mm/dd hh:mm:ss') }}
          </template>
        </el-table-column>
        <el-table-column prop="content" label="播种内容" width="200" show-overflow-tooltip>
          <template slot-scope="scope">
            {{ scope.row.content|| '--' }}
          </template>
        </el-table-column>
      </el-table>
    </template>
    <template slot="footer">
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
    </template>
  </layout-list>
</template>
<script lang="ts">
import { Component, Prop, Vue, Watch } from "vue-property-decorator";
import Api from '@/api/customer/sow';
/**
 * 1.新建;  
 * 2.编辑;  
 * 3.删除;  
 * 4.查看;  
 * 5.上下架
 * 6.播种数据
 */
type IOperationType = 1 | 2 | 3 | 4 | 5 | 6;
/**
 * 列表项
 */
interface IItem {
  /** 模版ID */
  id: number;
  /** 模版名称 */
  template_name: string;
  /** 是否使用模版  0 全部 1 默认播种 2 使用模版 */
  template_type: 0|1|2;
  /** 播种状态 */
  status: 1|2;
  /** 导购姓名 */
  sa_name: string;
  /** 所在部门 */
  sa_dp_name: string;
  /** 顾客姓名 */
  username: string;
  /** 会员卡号 */
  union_code: string;
  /** 创建时间 */
  created_at: number;
  /** 提醒时间 */
  remind_time: number;
  /** 播种内容 */
  content: string;
}

interface ISearchParams {
  /** 播种模版名称 */
  template_name: string;
  /** 播种模版id */
  template_id: number;
  /** 是否使用模版  0 全部 1 默认播种 2 使用模版 */
  template_type: 0|1|2;
  /** 导购姓名 */
  sa_name: string;
  /** 顾客姓名 */
  username: string;
  /** 会员编码 */
  union_code: string;
  /** 完成状态 1 未完成 2 已完成 */
  sow_status: 0|1|2;
  /** 创建时间开始 */
  start_created_at: number;
  /** 创建时间结束 */
  end_created_at: number;
  /** 完成时间开始 */
  start_finished_at: number;
  /** 完成时间结束 */
  end_finished_at: number;
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

/** 默认筛选条件 */
const defaultSearchParams: ISearchParams = {
  template_name: "",
  template_id: null,
  template_type: 0,
  sa_name: "",
  username: "",
  union_code: "",
  sow_status: 0,
  start_created_at: 0,
  end_created_at: 0,
  start_finished_at: 0,
  end_finished_at: 0
}

@Component({
  name: 'SowDataList',
})
export default class SowDataList extends Vue {
  @Prop({ type: String, default: '' }) title:string;
  @Prop({ type: Number, default: '' }) template_id:number;
  /** 搜索参数 */
  searchParams: ISearchParams = {...defaultSearchParams}
  /** 分页 */
  meta: IMate = {
    page: 1,
    size: 10,
    total: 0
  }
  /** 列表数据 */
  list: IItem[] = []

  /** 创建时间 读取 */
  get createdTime() {
    const { start_created_at, end_created_at } = this.searchParams
    return [start_created_at*1000, end_created_at*1000].filter(Boolean)
  }

  /** 创建时间 写入 */
  set createdTime(value: number[]) {
    const [start_time, end_time] = value || []
    this.searchParams.start_created_at = parseInt(start_time/1000+'');
    this.searchParams.end_created_at = parseInt(end_time/1000+'');
  }

  /** 提醒时间 读取 */
  get noticeTime() {
    const { start_finished_at, end_finished_at } = this.searchParams
    return [start_finished_at*1000, end_finished_at*1000].filter(Boolean)
  }

  /** 提醒时间 写入 */
  set noticeTime(value: number[]) {
    const [start_time, end_time] = value || []
    this.searchParams.start_finished_at = parseInt(start_time/1000+'');
    this.searchParams.end_finished_at = parseInt(end_time/1000+'');
  }

  mounted() {

    this.searchParams.template_name = this.$route.query.name as string;
    this.searchParams.template_id = Number(this.$route.query.id)||null;
    this.fetchData();
  }

  /**
   * 获取列表
   */
  fetchData() {
    const params = {
      page: this.meta.page,
      size: this.meta.size,
      ...this.searchParams,
      template_id: Number(this.searchParams.template_id)
    }
    Api.getSowRecordListApi(params).then((res: any) => {
      if (res.code == 200) {
        this.list = res.data.list;
        this.meta.total = res.data.total;
      } else {
        this.$message.error(res.message);
      }
    });
  }

  /**
   * 操作
   * @param key 
   * @param row 
   */
  handler(key: IOperationType, row?: IItem) {
    switch (key) {
      case 1:
        this.$router.push({
          name: "myapp_customer_sow_add",
        });
        break;
      case 2:
        this.$router.push({
          name: "myapp_customer_sow_edit",
          query: {
            id: String(row.id)
          }
        });
        break;
      case 3:
        this.handleRemove(row)
        break;
      case 4:
        this.$router.push({
          name: "myapp_customer_sow_info",
          query: {
            id: String(row.id)
          }
        });
        break;
      case 5:
        this.handleUpdate(row)
        break;
      case 6:
        console.log(row);
        
        break;
    
      default:
        break;
    }
  }

  /** 更新上下架 */
  async handleUpdate(row: IItem) {
    try {
      const [title, message, tips] = {
        1: ['下架提醒','确定下架该模版吗？', '下架成功'],
        2: ['上架提醒','确定上架该模版吗？', '上架成功'],
      }[row.status]
      await this.$confirm(message, title, {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning"
      })
      const res:any = await Api.updateSowStatusApi({
        id: row.id,
        status: row.status === 1 ? 2 : 1
      })
      if (res.code !== 200) return this.$message.error(res.message || '操作失败')
      this.$message.success(tips)
      row.status = row.status === 1 ? 2 : 1
      
    } catch (error) {
      console.log(error);
      
    }
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
   * @description 导出数据
   */
  handleExport() {
    Api.exportSowRecordListApi({
      ...this.searchParams,
      template_id: Number(this.searchParams.template_id)
    }).then(() => {
      this.$confirm(
        "文件已进入下载列表，请前往【下载中心】查看及下载文件。",
        "文件下载提示",
        {
          distinguishCancelAndClose: true,
          confirmButtonText: "前往下载中心",
          cancelButtonText: "取消",
          type: "success",
        }
      ).then(() => {
        this.$router.push({
          path: "/myapp/setting/download-center/list"
        });
      });
    })
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
   * 删除后刷新数据
   */
  refresh() {
    const { total, page, size } = this.meta
    if((total - 1) % size === 0 && (total - 1) / size < page) {
      this.meta.page = Math.max(page-1, 1)
    }
    this.fetchData()
  }

  /**
   * 删除
   */
  async handleRemove(row: IItem) {
    await this.$confirm(
      '确定删除当前模版内容？', 
      '删除提醒', 
      {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning"
      })
    const res:any = await Api.removeSowApi({
      id: row.id
    })
    if (res.code == 200) {
      this.$message.success('删除成功!')
      this.refresh()
    } else {
      this.$message.error(res.message);
    }
  }
}
</script>

<style lang="scss" scoped>
::v-deep.layout-list .search-content {
  padding-bottom: 0 !important;
}
  
</style>