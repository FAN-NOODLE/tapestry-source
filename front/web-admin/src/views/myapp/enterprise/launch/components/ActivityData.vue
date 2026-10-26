<template>
  <layout-list>
    <!-- <template slot="buttons">
      <el-button type="primary" @click="handler(1)">导出数据</el-button>
    </template> -->
    <template slot="search">
      <search-item text="时间：">
        <el-radio-group v-model="fastTimeType" @change="handleFastTimeChange">
          <el-radio-button label="1">今天</el-radio-button>
          <el-radio-button label="2">昨天</el-radio-button>
          <el-radio-button label="3">近7天</el-radio-button>
          <el-radio-button label="4">近30天</el-radio-button>
        </el-radio-group>
        <el-date-picker
          class="cy-el-date-picker m-l-10"
          v-model="updateTime"
          type="datetimerange"
          size="medium"
          store_type-separator="至"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
          value-format="timestamp"
          :default-time="['00:00:00', '23:59:59']"
          @change="handleChangeTime"
        >
        </el-date-picker>
      </search-item>
      <search-item text="活动名称：">
        <cy-input
          size="medium"
          placeholder="请输入活动名称"
          v-model="searchParams.name_str"
          clearable
        ></cy-input>
      </search-item>
      <search-item text="活动点位名称：">
        <cy-input
          size="medium"
          placeholder="请输入活动点位名称"
          v-model="searchParams.point_name_str"
          clearable
        ></cy-input>
      </search-item>
      <search-item text="投放渠道：">
        <el-select v-model="searchParams.category_id" style="width: 280px;">
          <el-option label="全部" :value="0"></el-option>
          <el-option v-for="item,index in categoryList" :key="index" :label="item.name" :value="item.id"></el-option>
        </el-select>
      </search-item>
      <el-button size="medium" type="primary" @click="handleSearch">查询</el-button>
      <el-button size="medium" plain @click="handleReset">重置</el-button>
    </template>
    <template slot="table">
      <el-row class="list m-b-20" :gutter="20">
        <el-col :span="4">
          <div class="item">
            <div class="item-label">新增好友数
              <el-tooltip content="筛选时间内，新增的外部联系好友（不去重）">
                <i class="el-icon-question"></i>
              </el-tooltip>
            </div>
            <div class="item-value">{{statistics.add_relation_total}}</div>
          </div>
        </el-col>
        <el-col :span="4">
          <div class="item">
            <div class="item-label">新增顾客数
              <el-tooltip content="筛选时间内，新增的顾客数（去重）">
                <i class="el-icon-question"></i>
              </el-tooltip>
            </div>
            <div class="item-value">{{statistics.add_customer_total}}</div>
          </div>
        </el-col>
        <el-col :span="4">
          <div class="item">
            <div class="item-label">删除好友数
              <el-tooltip content="筛选时间内，员工删除外部联系好友或外部联系好友删除员工，或者互删（不去重）">
                <i class="el-icon-question"></i>
              </el-tooltip>
            </div>
            <div class="item-value">{{statistics.del_relation_total}}</div>
          </div>
        </el-col>
        <el-col :span="4">
          <div class="item">
            <div class="item-label">注册会员数
              <el-tooltip content="筛选时间内，添加的外部好友中是会员的人数（去重）">
                <i class="el-icon-question"></i>
              </el-tooltip>
            </div>
            <div class="item-value">{{statistics.add_vip_total}}</div>
          </div>
        </el-col>
      </el-row>
      <el-table :data="list" stripe>
        <el-table-column prop="activity_name" label="活动名称" min-width="100"></el-table-column>
        <el-table-column prop="category_name" label="投放渠道" min-width="100">
          <template slot-scope="scope">
            <el-tag v-if="scope.row.category_name && scope.row.category_name !=='-'">{{scope.row.category_name}}</el-tag>
            <span v-else>--</span>
          </template>
        </el-table-column>
        <el-table-column prop="point_name" label="活动点位" min-width="100">
          <template slot-scope="scope">
            <span v-if="scope.row.point_name">{{scope.row.point_name}}</span>
            <span v-else>--</span>
          </template>
        </el-table-column>
        <el-table-column prop="total_statistics.add_relation_total" label="新增好友数" min-width="100"></el-table-column>
        <el-table-column prop="total_statistics.add_customer_total" label="新增顾客数" min-width="100"></el-table-column>
        <el-table-column prop="total_statistics.del_relation_total" label="删除好友数" min-width="100"></el-table-column>
        <el-table-column prop="total_statistics.add_vip_total" label="注册会员数" min-width="100"></el-table-column>
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
import { getCategoryListApi, getStatisticsApi } from "@/api/enterprise/activity";
import { Component, Vue } from "vue-property-decorator";
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
  id?: number;
  activity_name: string;
  category_name: string;
  point_name: string;
  total_statistics: {
    add_customer_total: number;
    add_relation_total: number;
    add_vip_total: number;
    del_relation_total: number;
    state: string | null;
  };
}

interface ISearchParams {
  /** 名称 */
  name_str: string
  point_name_str: string
  category_id: number;
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
  name_str: "",
  point_name_str:"",
  category_id: 0,
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
  statistics =  {
    state: '',
    add_relation_total: 0, //新增好友数
    add_customer_total: 0, //新增顾客数
    del_relation_total: 0, //删除好友数
    add_vip_total: 0 //注册会员数
  }

  categoryList:any[] = []

  fastTimeType = '3'

  pickerOptions = {
    shortcuts: [
      {
        text: '今天',
        onClick(picker) {
          const end = new Date();
          const start = new Date();
          // 今天 00:00:00 - 23:59:59
          start.setHours(0, 0, 0, 0);
          end.setHours(23, 59, 59, 999);
        
          picker.$emit('pick', [start, end]);
        }
      }, 
      {
        text: '昨天',
        onClick(picker) {
          const end = new Date();
          const start = new Date();
          // 00:00:00 - 23:59:59
          start.setHours(0, 0, 0, 0);
          end.setHours(23, 59, 59, 999);
          start.setTime(start.getTime() - 3600 * 1000 * 24);
          end.setTime(end.getTime() - 3600 * 1000 * 24);

          picker.$emit('pick', [start, end]);
        }
      }, 
      {
        text: '近7天',
        onClick(picker) {
          const end = new Date();
          const start = new Date();
          start.setTime(start.getTime() - 3600 * 1000 * 24 * 7);
          picker.$emit('pick', [start, end]);
        }
      }, 
      {
        text: '近30天',
        onClick(picker) {
          const end = new Date();
          const start = new Date();
          start.setTime(start.getTime() - 3600 * 1000 * 24 * 30);
          picker.$emit('pick', [start, end]);
        }
      }
    ]
  }

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
    this.searchParams.name_str = String(this.$route.query.name ||'')
    this.handleFastTimeChange(this.fastTimeType)
    this.fetchData();
    this.getCategory()
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
    getStatisticsApi(params).then((res: any) => {
      if (res.code == 200) {
        this.list = res.data.list;
        this.meta.total = res.data.total;
        this.statistics = res.data.total_statistics;
      } else {
        this.$message.error(res.message);
      }
    });
  }

  handleFastTimeChange(e) {

    const end = new Date();
    const start = new Date();
    // 今天 00:00:00 - 23:59:59
    start.setHours(0, 0, 0, 0);
    end.setHours(23, 59, 59, 999);

    switch (e) {
      case '1':
        // 今天
        start.setTime(start.getTime() - 3600 * 1000 * 24 * 0)
        end.setTime(end.getTime() - 3600 * 1000 * 24 * 0)
        break;
      case '2':
        // 昨天
        start.setTime(start.getTime() - 3600 * 1000 * 24 * 1)
        end.setTime(end.getTime() - 3600 * 1000 * 24 * 1)
        break;
      case '3':
        // 默认近7天
        start.setTime(start.getTime() - 3600 * 1000 * 24 * 6);
        end.setTime(end.getTime() - 3600 * 1000 * 24 * 0);
        break;
      case '4':
        // 默认近30天
        start.setTime(start.getTime() - 3600 * 1000 * 24 * 29);
        end.setTime(end.getTime() - 3600 * 1000 * 24 * 0);
        break;
    }

    this.searchParams.start_time = parseInt(start.getTime()/1000+'');
    this.searchParams.end_time = parseInt(end.getTime()/1000+'');
  }

  handleChangeTime() {
    this.fastTimeType = ''
  }

  getCategory() {
    getCategoryListApi({}).then((res:any) => {
      if (res.code == 200) {
        this.categoryList = res.data
      } else {
        this.$message.error(res.message);
      }
    })
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
    this.fastTimeType = '3'
    this.searchParams = { ...defaultSearchParams }
    this.handleFastTimeChange(this.fastTimeType)
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
        // Api.removeWhitelistApi({
        //   id: row.id
        // }).then((res: any) => {
        //   if (res.code == 200) {
        //     this.$message.success('删除成功!')
        //     this.refresh()
        //   } else {
        //     this.$message.error(res.message);
        //   }
        // });
      })
  }
}
</script>

<style lang="scss" scoped>
::v-deep.layout-list {
  padding: 0;
}
::v-deep.layout-list .search-content {
  padding-bottom: 0 !important;
  margin-top: 0;
}
.list {
  padding: 15px;
  border-radius: 4px;
  background-color: #f7f7f7;
}

.item {
  border-radius: 4px;
  background-color: #fff;
  padding: 15px;
  .item-label {
    line-height: 30px;
    font-size: 16px;
    color: #666;
    text-align: center;
  }
  .item-value {
    line-height: 50px;
    font-size: 30px;
    text-align: center;
    font-weight: bold;
  }
}
  
</style>