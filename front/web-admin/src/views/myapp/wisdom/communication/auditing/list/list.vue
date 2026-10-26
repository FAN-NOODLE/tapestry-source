<template>
  <layout-list>
    <template slot="buttons"></template>
    <template slot="search">
      <search-item text="标题名称：">
        <cy-input
          size="medium"
          placeholder="请输入标题名称"
          v-model="searchParams.title"
          clearable
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
      <search-item text="审核时间：">
        <el-date-picker
          v-model="reviewTime"
          size="medium"
          type="datetimerange"
          store_type-separator="至"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
          value-format="timestamp"
          :default-time="['00:00:00', '23:59:59']"
        >
        </el-date-picker>
      </search-item>
      <search-item text="传讯状态：">
        <cy-select
          size="medium"
          placeholder="请选择状态"
          v-model="searchParams.status"
          :optionsList="statusOptions"
        ></cy-select>
      </search-item>
      <search-item text="所属部门：">
        <cy-select
          size="medium"
          placeholder="请选择所属部门"
          v-model="searchParams.department_name"
          :optionsList="departmentOptions"
        ></cy-select>
      </search-item>
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
    </template>
    <template slot="table">
      <el-table :data="list" stripe>
        <el-table-column prop="id" label="ID" min-width="180"></el-table-column>
        <el-table-column prop="title" label="标题名称" min-width="200"></el-table-column>
        <el-table-column prop="type" label="传讯类型" min-width="100">
          <template slot-scope="scope">
            {{ typeName(scope.row.type) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="传讯状态" min-width="100">
          <template slot-scope="scope">
            <el-tag size="small" :type="{
              1: 'warning',
              3: 'danger',
              4: 'success'
            }[scope.row.status]||''">{{ statusName(scope.row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="rel_status" label="发布状态" width="100">
          <template slot-scope="scope">
            <span>{{
            {
              '-1': "未发布",
              0: "未发布",
              1: "已发布",
            }[scope.row.rel_status]
            }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="department_name" label="所属部门" min-width="120">
          <template slot-scope="scope">
            <span>{{scope.row.department_name || '--'}}</span>
          </template>
        </el-table-column>
        <el-table-column prop="submit_name" label="提交者" min-width="100"></el-table-column>
        <el-table-column prop="update_at" label="最近更新时间" width="200">
          <template slot-scope="scope">
            {{ scope.row.update_at | formatDate('yyyy/mm/dd hh:mm:ss')}}
          </template>
        </el-table-column>
        <el-table-column prop="review_name" label="最近审核人员" width="180">
          <template slot-scope="scope">
            {{ scope.row.review_name || '--' }}
          </template>
        </el-table-column>
        <el-table-column prop="review_time" label="最近审核时间" width="200">
          <template slot-scope="scope">
            {{ scope.row.review_time | formatDate('yyyy/mm/dd hh:mm:ss') }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template slot-scope="scope">
            <el-button type="text" :disabled="scope.row.status === 5" @click="handler(5, scope.row)">作废</el-button>
            <el-button type="text" @click="handler(4, scope.row)">查看</el-button>
            <el-button type="text" :disabled="scope.row.status !== 1" @click="handler(2, scope.row)">审核</el-button>
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
import { typeOptions, statusOptions } from '../../const'
/**
 * 1.新建;  
 * 2.审核;  
 * 3.撤回;  
 * 4.查看; 
 * 5.作废; 
 */
type IOperationType = 1 | 2 | 3 | 4 | 5 ;
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

/** 分页 */
interface IMate {
  /** 页数 */
  page: number;
  /** 每页数量 */
  size: number;
  /** 总数 */
  total: number;
}

interface ISearchParams {
  /** 名称 */
  title: string;
  /** 更新时间 开始 */
  updated_start_time: number;
  /** 更新时间 结束 */
  updated_end_time: number;
  /** 审核时间 开始 */
  review_start_time: number;
  /** 审核时间 结束 */
  review_end_time: number;
  /** 传讯状态 */
  status: number;
  /** 所属部门 */
  department_name: string;
  /** 传讯类型 */
  type: number;
  /** RO|管理员身份，值为1，FT类型值为0 */
  is_all: 1;
}

/** 默认筛选条件 */
const defaultSearchParams: ISearchParams = {
  title: '',
  updated_start_time: 0,
  updated_end_time: 0,
  status: -1,
  department_name: '',
  type: -1,
  review_start_time: 0,
  review_end_time: 0,
  is_all: 1
}

@Component({
  name: 'AuditingList'
})
export default class AuditingList extends Vue {
  /** 传讯类型选项 */
  typeOptions = [
    { label: '全部', value: -1 },
    ...typeOptions
  ]
  /** 所属门店选项 */
  departmentOptions = []
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

  /** 传讯状态选项 */
  get statusOptions() {
    return [
      { label: '全部', value: -1 },
    ].concat(statusOptions.filter(item => item.value !== 0 && item.value !== 2))
  }

  /** 更新时间 读取 */
  get updateTime() {
    const { updated_start_time, updated_end_time } = this.searchParams
    return [updated_start_time*1000, updated_end_time*1000].filter(Boolean)
  }

  /** 更新时间 写入 */
  set updateTime(value: number[]) {
    const [start_time, end_time] = value || []
    this.searchParams.updated_start_time = parseInt(start_time/1000+'');
    this.searchParams.updated_end_time = parseInt(end_time/1000+'');
  }
  /** 审核时间 读取 */
  get reviewTime() {
    const { review_start_time, review_end_time } = this.searchParams
    return [review_start_time*1000, review_end_time*1000].filter(Boolean)
  }

  /** 审核时间 写入 */
  set reviewTime(value: number[]) {
    const [start_time, end_time] = value || []
    this.searchParams.review_start_time = parseInt(start_time/1000+'');
    this.searchParams.review_end_time = parseInt(end_time/1000+'');
  }


  get typeName() {
    return (value) => {
      return typeOptions.find(item => item.value === value).label
    }
  }

  get statusName() {
    return (value) => {
      const current = statusOptions.find(item => item.value === value)
      if(!current) {
        return '--'
      }
      return current.label
    }
  }

  created() {
    this.fetchData();
    this.getDepartmentList()
  }

  activated() {
    this.fetchData();
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
    const res: any = await Api.getContentListApi(params)
    if (res.code == 200) {
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

  /**
   * 操作
   * @param key 
   * @param row 
   */
  handler(key: IOperationType, row?: IItem) {
    switch (key) {
      case 1:
        break;
      case 2:
        this.$router.push({
          name: "myapp_wisdom_communication_auditing_edit",
          query: {
            id: String(row.id)
          }
        });
        break;
      case 3:
        break;
      case 4:
        this.$router.push({
          name: "myapp_wisdom_communication_auditing_detail",
          query: {
            id: String(row.id)
          }
        });
        break;
      case 5:
        this.handleCancel(row);
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
   * 作废
   */
   handleCancel(row: IItem) {
    this.$confirm('确定作废当前内容？作废后该内容将不可被再次修改编辑。', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(() => {
      Api.cancelContentApi({
        id: row.id
      }).then((res: any) => {
        if (res.code == 200) {
          this.$message.success('作废成功!')
          this.fetchData()
        } else {
          this.$message.error(res.message);
        }
      });
    });
  }
}
</script>

<style lang="scss" scoped>
::v-deep.layout-list .search-content {
  padding-bottom: 0 !important;
}
  
</style>