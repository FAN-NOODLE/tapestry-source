<template>
  <layout-list>
    <template slot="buttons">
      <el-button type="primary" @click="handler(1)">新建活动</el-button>
    </template>
    <template slot="search">
      <search-item text="活动名称：">
        <cy-input
          size="medium"
          placeholder="请输入活动名称"
          v-model="searchParams.name"
          clearable
        ></cy-input>
      </search-item>
      <search-item text="投放渠道：">
        <el-select v-model="searchParams.category_id" style="width: 280px;">
          <el-option label="全部" :value="0"></el-option>
          <el-option v-for="item,index in categoryList" :key="index" :label="item.name" :value="item.id"></el-option>
        </el-select>
      </search-item>
      <search-item text="活动有效期：">
        <el-date-picker
          class="cy-el-date-picker"
          v-model="effectTime"
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
      <el-button size="medium" type="primary" @click="handleSearch">查询</el-button>
      <el-button size="medium" plain @click="handleReset">重置</el-button>
    </template>
    <template slot="table">
      <el-table :data="list" stripe>
        <el-table-column prop="id" label="活动ID" min-width="180"></el-table-column>
        <el-table-column prop="name" label="活动名称" min-width="140" show-overflow-tooltip></el-table-column>
        <el-table-column prop="alias_name" label="对外活动名称" min-width="140" show-overflow-tooltip>
          <template slot-scope="scope">
            <span>{{ scope.row.alias_name || '--' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="category_names" label="投放渠道" min-width="180">
          <template slot-scope="scope">
            <!-- <el-tag v-for="(item, index) in scope.row.category_names" :key="index">{{item}}</el-tag> -->
            <el-tooltip placement="top" effect="light" v-if="scope.row.category_names.length > 3">
              <div class="tags" style="cursor: pointer;">
                <!-- 显示前2个 -->
                <el-tag class="tag" v-for="item,index in scope.row.category_names.slice(0,3)" :key="index">{{item}}</el-tag>
                <span>...</span>
              </div>
              <div class="tags" slot="content">
                <el-tag v-for="item,index in scope.row.category_names" :key="index">{{item}}</el-tag>
              </div>
            </el-tooltip>
            <div v-else-if="scope.row.category_names.length <= 3 && scope.row.category_names.length > 0" class="tags">
              <el-tag v-for="name,index in scope.row.category_names" :key="index">{{name}}</el-tag>
            </div>
            <span v-else>--</span>
          </template>
        </el-table-column>
        <el-table-column prop="tag_names" label="自定义标签" min-width="180">
          <template slot-scope="scope">
            <el-tooltip placement="top" effect="light" v-if="scope.row.tag_names.length > 3">
              <div class="tags" style="cursor: pointer;">
                <!-- 显示前2个 -->
                <el-tag class="tag" v-for="item,index in scope.row.tag_names.slice(0,3)" :key="index">{{item}}</el-tag>
                <span>...</span>
              </div>
              <div class="tags" slot="content">
                <el-tag v-for="item,index in scope.row.tag_names" :key="index">{{item}}</el-tag>
              </div>
            </el-tooltip>
            <div v-else-if="scope.row.tag_names.length <= 3 && scope.row.tag_names.length > 0" class="tags">
              <el-tag v-for="name,index in scope.row.tag_names" :key="index">{{name}}</el-tag>
            </div>
            <span v-else>--</span>
          </template>
        </el-table-column>
        <el-table-column prop="has_welcome" label="是否配置欢迎语" min-width="160">
          <template slot-scope="scope">
            <el-tag v-if="scope.row.has_welcome" type="success">已配置</el-tag>
            <el-tag v-else type="info">未配置</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="time_type" label="活动有效期" width="185">
          <template slot-scope="scope">
            <div v-if="scope.row.time_type === 1">
              <div>{{ scope.row.start_time | formatDate('yyyy/mm/dd hh:mm:ss') }}</div> 
              <div>{{ scope.row.end_time | formatDate('yyyy/mm/dd hh:mm:ss') }}</div>
            </div>
            <span v-else>长期有效</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="185">
          <template slot-scope="scope">
            {{ scope.row.created_at | formatDate('yyyy/mm/dd hh:mm:ss') }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="400" fixed="right">
          <template slot-scope="scope">
            <el-button type="text" @click="handler(6, scope.row)">下载活动码</el-button>
            <el-button type="text" @click="handler(5, scope.row)">复制链接</el-button>
            <!-- <el-button type="text" @click="handler(7, scope.row)">复制活动页面参数</el-button> -->

            <el-button type="text" @click="handler(8, scope.row)">查看数据</el-button>
            
            <el-dropdown
              class="del-dropdown"
              @command="handleCommand($event, scope.row)"
            >
              <el-button class="m-l-10" type="text">更多</el-button>
              <el-dropdown-menu slot="dropdown" style="width: 80px;">
                <el-dropdown-item command="2">编辑</el-dropdown-item>
                <el-dropdown-item command="3">删除</el-dropdown-item>
              </el-dropdown-menu>
            </el-dropdown>
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
      <download-popup ref="downloadPopupRef"></download-popup>
    </template>
  </layout-list>
</template>
<script lang="ts">
import { deleteActivityApi, getActivityListApi, getCategoryListApi } from "@/api/enterprise/activity";
import { Component, Vue } from "vue-property-decorator";
import DownloadPopup from "./download-popup.vue";
/**
 * 1.新建;  
 * 2.编辑;  
 * 3.删除;  
 * 4.查看;  
 * 5.复制链接;
 * 6.复制太阳码;
 * 7.复制活动页面参数;
 * 8.查看数据;
 */
type IOperationType = 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8;
/**
 * 列表项
 */
interface IItem {
  id: number;
  name: string;
  alias_name: string;
  banner_url: string;
  time_type: 1 | 2;
  start_time: number;
  end_time: number;
  skip_verify: 0 | 1;
  add_name_tag: 0 | 1;
  add_category_tag: 0 | 1;
  created_at: number;
  tag_names: string[];
  category_names: string[];
}

interface ISearchParams {
  name?: string;              // 可选：名称模糊匹配
  category_id?: number;       // 可选：分类 ID
  effect_time: {
    start_time: number;       // Unix 时间戳（秒）
    end_time: number;
  };
  create_time: {
    start_time: number;
    end_time: number;
  }; 
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
  effect_time: {
    start_time: 0,
    end_time: 0
  },
  create_time: {
    start_time: 0,
    end_time: 0
  },
  name: '',
  category_id: 0,
}

@Component({
  components: {
    DownloadPopup
  }

})
export default class ActivityList extends Vue {
  /** 搜索参数 */
  searchParams: ISearchParams = JSON.parse(JSON.stringify(defaultSearchParams))
  /** 分页 */
  meta: IMate = {
    page: 1,
    size: 10,
    total: 0
  }
  /** 列表数据 */
  list: IItem[] = []
  categoryList: any[] = []

  /** 有效期 读取 */
  get effectTime() {
    const { start_time, end_time } = this.searchParams.effect_time
    return [start_time*1000, end_time*1000].filter(Boolean)
  }

  /** 有效期 写入 */
  set effectTime(value: number[]) {
    const [start_time, end_time] = value || []
    this.searchParams.effect_time.start_time = parseInt(start_time/1000+'');
    this.searchParams.effect_time.end_time = parseInt(end_time/1000+'');
  }

  /** 创建时间 读取 */
  get createTime() {
    const { start_time, end_time } = this.searchParams.create_time
    return [start_time*1000, end_time*1000].filter(Boolean)
  }

  /** 创建时间 写入 */
  set createTime(value: number[]) {
    const [start_time, end_time] = value || []
    this.searchParams.create_time.start_time = parseInt(start_time/1000+'');
    this.searchParams.create_time.end_time = parseInt(end_time/1000+'');
  }

  mounted() {
    this.fetchData();
    this.getCategory()
  }
  // activated() {
  //   this.fetchData();
  // }

  /**
   * 获取列表
   */
  fetchData() {
    const params = {
      page: this.meta.page,
      size: this.meta.size,
      ...this.searchParams,
    }
    getActivityListApi(params).then((res: any) => {
      if (res.code == 200) {
        this.list = res.data.list;
        this.meta.total = res.data.total;
      } else {
        this.$message.error(res.message);
      }
    });
  }

  /**
   * 获取渠道列表
   */
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
          name: "myapp_enterprise_launch_create",
        });
        break;
      case 2:
        this.$router.push({
          name: "myapp_enterprise_launch_edit",
          query: {
            id: String(row.id)
          }
        });
        break;
      case 3:
        this.handleRemove(row)
        break;
      case 5:
        (this.$refs.downloadPopupRef as any).show(row.id, 2)
        break;
      case 6:
        (this.$refs.downloadPopupRef as any).show(row.id, 1)
        break;
      case 7:
        (this.$refs.downloadPopupRef as any).show(row.id, 3)
        break;
      case 8:
        window.open(`/web-admin/myapp/enterprise/launch/list?name=${row.name}&active=2`)
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
    this.searchParams = JSON.parse(JSON.stringify(defaultSearchParams))
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
      '确定删除当前活动内容？', 
      '删除提醒', 
      {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning"
      })
      .then(() => {
        deleteActivityApi({
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

  handleCommand(command, row) {
    this.handler(Number(command) as IOperationType, row)
  }
}
</script>

<style lang="scss" scoped>
::v-deep.layout-list {
  padding: 0;
}
::v-deep.layout-list .search-content {
  padding-bottom: 0 !important;
}
.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
}

.tag {
  max-width: 80px;
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

</style>