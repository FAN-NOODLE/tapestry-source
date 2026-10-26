<template>
  <layout-list>
    <template slot="buttons">
      <el-button type="primary" @click="handler(1)">新建企业传讯</el-button>
      <remind-config></remind-config>
    </template>
    <template slot="search">
      <search-item text="标题：">
        <cy-input
          size="medium"
          placeholder="请输入标题名称"
          v-model="searchParams.title"
          clearable
        ></cy-input>
      </search-item>
      <search-item text="提交者：">
        <cy-input
          size="medium"
          placeholder="请输入提交者名称"
          v-model="searchParams.member_name"
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
      <el-button size="medium" type="primary" @click="handleSearch">查询</el-button>
      <el-button size="medium" plain @click="handleReset">重置</el-button>
    </template>
    <template slot="table">
      <el-table :data="list" stripe>
        <el-table-column prop="id" label="ID" min-width="180"></el-table-column>
        <el-table-column prop="title" label="标题名称" min-width="200"></el-table-column>
        <el-table-column prop="cover" label="封面图" min-width="100">
          <template slot-scope="scope">
            <el-image :src="scope.row.cover||defaultCover" fit="contain" lazy :preview-src-list="[scope.row.cover||defaultCover]"></el-image>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="传讯类型" min-width="140">
          <template slot-scope="scope">
            <span v-if="scope.row.type === 1">周传讯</span>
            <span v-if="scope.row.type === 2">转货传讯</span>
            <span v-if="scope.row.type === 3">紧急传讯</span>
          </template>
        </el-table-column>
        <el-table-column prop="send_type" label="发布时间" min-width="200">
          <template slot-scope="scope">
            <div v-if="scope.row.send_type == 0 && scope.row.status == 0">定时发布</div>
            <div v-if="scope.row.send_type == 2 && !scope.row.send_time">暂未发布</div>
            <div v-if="scope.row.send_time">{{scope.row.send_time | formatDate('yyyy/mm/dd hh:mm:ss')}}</div>
          </template>
        </el-table-column>
        <el-table-column prop="member_name" label="提交者" min-width="100">
          <template slot-scope="scope">
            {{ scope.row.member_name || '--' }}
          </template>
        </el-table-column>
        <el-table-column prop="updated_at" label="最近更新时间" width="200">
          <template slot-scope="scope">
            {{ scope.row.updated_at | formatDate('yyyy/mm/dd hh:mm:ss')}}
          </template>
        </el-table-column>
        <el-table-column prop="send_num" label="下发人数" width="120">
          <template slot-scope="scope">
            {{ scope.row.send_num || '--' }}
          </template>
        </el-table-column>
        <el-table-column prop="view_people_num" label="查看人数" width="120">
          <template slot-scope="scope">
            {{ scope.row.view_people_num || '--' }}
          </template>
        </el-table-column>
        <el-table-column prop="view_num" label="查看次数" width="120">
          <template slot-scope="scope">
            {{ scope.row.view_num || '--' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="240" fixed="right">
          <template slot-scope="scope">
            <!-- RO及管理员可以对传讯内容进行【查看】【编辑】【发布】 -->
            <!-- Function Team及SA仅可【查看】 -->
            <el-button type="text" @click="handler(4, scope.row)">查看</el-button>
            <el-button type="text" :disabled="scope.row.status === 1" @click="handler(2, scope.row)">编辑</el-button>
            <el-button v-if="scope.row.status === 0 && scope.row.send_type === 2" type="text" @click="handler(3, scope.row)">发布</el-button>
            <el-button v-if="scope.row.status === 1 && scope.row.type === 1" type="text" @click="handler(5, scope.row)">复制邮件文本</el-button>
          </template>
        </el-table-column>
      </el-table>
      <copy-text v-if="detail" ref="copyText" :value="detail"></copy-text>
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
import CopyText from '../components/copy-text.vue'
import RemindConfig from '../components/remind-config.vue'
import { defaultCover } from "../../const";
/**
 * 1.新建;  
 * 2.编辑;  
 * 3.立即发布;  
 * 4.查看; 
 * 5.复制邮件文本; 
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
  /** 提交者 */
  member_name: string;
  /** 更新时间 开始 */
  start_time: number;
  /** 更新时间 结束 */
  end_time: number;
}

/** 默认筛选条件 */
const defaultSearchParams: ISearchParams = {
  title: '',
  member_name: '',
  start_time: 0,
  end_time: 0,
}

@Component({
  name: 'IntegrationList',
  components: {
    CopyText,
    RemindConfig
  }
})
export default class IntegrationList extends Vue {
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

  detail: any = null

  defaultCover = defaultCover

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

  created() {
    this.fetchData();
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
    const res: any = await Api.getMergeListApi(params)
    if (res.code == 200) {
      this.list = res.data.list;
      this.meta.total = res.data.total;
    } else {
      this.$message.error(res.message);
    }
  }

  /**
   * 获取详情
   */
  async getDetail(row: IItem) {
    const res: any = await Api.getMergeDetailApi({
      id: Number(row.id)
    })
    if (res.code == 200) {
      res.data.content = res.data.content_list
      this.detail = res.data
    } else {
      this.$message.error(res.message);
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
        this.$router.push({
          name: "myapp_wisdom_communication_integration_add",
          query: {
            t: Date.now().toString()
          }
        });
        break;
      case 2:
        this.$router.push({
          name: "myapp_wisdom_communication_integration_edit",
          query: {
            id: String(row.id)
          }
        });
        break;
      case 3:
        this.handleRelease(row)
        break;
      case 4:
        this.$router.push({
          name: "myapp_wisdom_communication_integration_detail",
          query: {
            id: String(row.id)
          }
        });
        break;
      case 5:
        // 复制邮件文本
        this.handleCopyText(row)
        break;
    
      default:
        break;
    }
  }

  /**
   * 复制邮件文本
   * @param row 
   */
  async handleCopyText(row: IItem) {
    await this.getDetail(row);
    (this.$refs.copyText as any).copy()
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
   * 立即发布
   */
  handleRelease(row: IItem) {
    this.$confirm('确定立即发布？发布后该传讯内容将通知到接收SA。', '发布提醒', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(() => {
      Api.releaseMergeApi({
        id: row.id
      }).then((res: any) => {
        if (res.code == 200) {
          this.$message.success('发布成功!')
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

.el-image {
  width: 60px;
  height: 60px;
}
  
</style>