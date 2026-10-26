<template>
  <layout-list>
    <template slot="buttons">
      <el-button type="primary" @click="handler(1)">新建播种模版</el-button>
    </template>
    <template slot="search">
      <search-item text="播种模版名称：">
        <cy-input
          size="medium"
          placeholder="请输入播种模版名称"
          v-model="searchParams.name"
          clearable
        ></cy-input>
      </search-item>
      <search-item text="是否上架：">
        <cy-select
          size="medium"
          placeholder="请选择是否上架"
          v-model="searchParams.status"
          clearable
        >
        <cy-option label="全部" :value="0" ></cy-option>
        <cy-option label="已上架" :value="1" ></cy-option>
        <cy-option label="未上架" :value="2" ></cy-option>
      </cy-select>
      </search-item>
      <search-item text="创建时间时间：">
        <el-date-picker
          class="cy-el-date-picker"
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
      <el-button class="m-b-20" size="small" type="primary" @click="handleSearch">查询</el-button>
      <el-button class="m-b-20" size="small" plain @click="handleReset">重置</el-button>
      <el-button class="m-b-20" size="small" plain @click="handleExport">导出数据</el-button>
    </template>
    <template slot="table">
      <el-table :data="list" stripe>
        <el-table-column prop="id" label="模版ID" width="200"></el-table-column>
        <el-table-column prop="name" label="播种模版名称" min-width="200" show-overflow-tooltip></el-table-column>
        <el-table-column prop="task_num" label="播种次数" min-width="100"></el-table-column>
        <el-table-column prop="status" label="是否上架" min-width="100">
          <template slot-scope="scope">
            <span @click="handler(5, scope.row)">
              <el-switch :value="scope.row.status" :active-value="1" :inactive-value="2"></el-switch>
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="use_num" label="使用模版次数" min-width="100"></el-table-column>
        <el-table-column prop="use_sa_num" label="引用导购人数" min-width="120"></el-table-column>
        <el-table-column prop="admin_name" label="创建者" min-width="100"></el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="200">
          <template slot-scope="scope">
            {{ scope.row.created_at | formatDate('yyyy/mm/dd hh:mm:ss') }}
          </template>
        </el-table-column>
        <el-table-column prop="updated_at" label="最近更新时间" width="200">
          <template slot-scope="scope">
            {{ scope.row.updated_at | formatDate('yyyy/mm/dd hh:mm:ss') }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="260" fixed="right">
          <template slot-scope="scope">
            <el-button type="text" @click="handler(4, scope.row)">查看</el-button>
            <el-button type="text" @click="handler(2, scope.row)">编辑</el-button>
            <el-button type="text" :disabled="scope.row.status == 1" @click="handler(3, scope.row)">删除</el-button>
            <el-button type="text" @click="handler(6, scope.row)">播种数据</el-button>
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
import { Component, Vue, Watch } from "vue-property-decorator";
import Api from '@/api/customer/sow';
import ShareModule from "@/shared/index";
import { download } from "@/utils/download";
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
  name: string;
  /** 模版状态 */
  status: number;
  /** 模版任务数 */
  task_num: number;
  /** 模版引用数 */
  use_num: number;
  /** 模版引用导购数 */
  use_sa_num: number;
  /** 创建时间 */
  create_at: number;
  /** 更新时间 */
  update_at: number;
  /** 创建人 */
  admin_name: string;
}

interface ISearchParams {
  /** 名称 */
  name: string;
  /** 是否上架 0 全部 1 是 2 否 */
  status: 0|1|2;
  /** 执行时间 开始 */
  start_time: number;
  /** 执行时间 结束 */
  end_time: number;
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
  name: '',
  status: 0,
  start_time: 0,
  end_time: 0
}

@Component({
  name: 'SowTempList',
})
export default class SowTempList extends Vue {
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
  get updateTime() {
    const { start_time, end_time } = this.searchParams
    return [start_time*1000, end_time*1000].filter(Boolean)
  }

  /** 创建时间 写入 */
  set updateTime(value: number[]) {
    const [start_time, end_time] = value || []
    this.searchParams.start_time = parseInt(start_time/1000+'');
    this.searchParams.end_time = parseInt(end_time/1000+'');
  }
  mounted() {
    this.fetchData();
  }
  activated() {
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
    }
    Api.getSowListApi(params).then((res: any) => {
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
        this.$emit('view-data', row)
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
    Api.exportSowListApi(this.searchParams).then((res:any) => {
      download(res, 'application/x-xls;charset=utf-8', '播种模版导出.xlsx')
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