<template>
  <layout-list>
    <template slot="buttons">
      <el-button type="primary" @click="handler(1)">新建投放渠道</el-button>
    </template>
    <template slot="search">
      <search-item text="投放渠道名称：">
        <cy-input
          size="medium"
          placeholder="请输入投放渠道名称"
          v-model="searchParams.name"
          clearable
        ></cy-input>
      </search-item>
      <search-item text="创建时间：">
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
      <el-button size="medium" type="primary" @click="handleSearch">查询</el-button>
      <el-button size="medium" plain @click="handleReset">重置</el-button>
    </template>
    <template slot="table">
      <el-table :data="list" stripe>
        <el-table-column prop="id" label="ID" min-width="100"></el-table-column>
        <el-table-column prop="name" label="投放渠道名称" min-width="100"></el-table-column>
        <el-table-column prop="created_at" label="创建时间" min-width="100">
          <template slot-scope="scope">
            {{ scope.row.created_at | formatDate('yyyy/mm/dd hh:mm:ss') }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180">
          <template slot-scope="scope">
            <el-button type="text" @click="handler(2, scope.row)">编辑</el-button>
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
import { Component, Vue } from "vue-property-decorator";
import Api from '@/api/wisdom/communication';
/**
 * 1.新建;  
 * 2.编辑;  
 * 3.删除;  
 * 4.查看;  
 */
type IOperationType = 1 | 2 | 3 | 4 ;
/**
 * 列表项
 */
interface IItem {
  id: number;
  name: string;
  num: number;
  create_by: number;
  update_at: number;
}

interface ISearchParams {
  /** 名称 */
  name: string;
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
  start_time: 0,
  end_time: 0
}

@Component({
  name: 'Whitelist',
})
export default class Whitelist extends Vue {
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
    Api.getWhitelistApi(params).then((res: any) => {
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
          name: "myapp_wisdom_communication_whitelist_add",
        });
        break;
      case 2:
        this.$router.push({
          name: "myapp_wisdom_communication_whitelist_edit",
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
          name: "myapp_wisdom_communication_whitelist_detail",
          query: {
            id: String(row.id)
          }
        });
        break;
    
      default:
        break;
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
  handleRemove(row: IItem) {
    this.$confirm(
      '确定删除当前白名单内容？', 
      '删除提醒', 
      {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning"
      })
      .then(() => {
        Api.removeWhitelistApi({
          id: row.id
        }).then((res: any) => {
          if (res.code == 200) {
            this.$message.success('删除成功!')
            this.refresh()
          } else {
            this.$message.error(res.message);
          }
        });
      })
  }
}
</script>

<style lang="scss" scoped>
::v-deep.layout-list .search-content {
  padding-bottom: 0 !important;
}
  
</style>