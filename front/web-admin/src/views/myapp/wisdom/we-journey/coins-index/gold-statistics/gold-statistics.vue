<template>
  <layout-list>
    <!-- 搜索区域 -->
    <template slot="search">
      <search-item text="导购ID：">
        <cy-input
          v-model="search.user_id"
          placeholder="请输入导购ID"
          clearable
        >
        </cy-input>
      </search-item>
      <search-item text="导购姓名：">
        <cy-input
          v-model="search.user_name"
          placeholder="请输入导购姓名"
          clearable
        >
        </cy-input>
      </search-item>
      <search-item text="所在部门：">
        <select-tree-department
          ref="selectTreeDepartment"
          @change="departmentChange"
        ></select-tree-department>
      </search-item>
      <search-item text="类型：">
        <cy-select
          v-model="search.change_type"
          size="medium"
          :optionsList="levelOptions"
          :options="{ label: 'name', value: 'id' }"
        >
        </cy-select>
      </search-item>
      <search-item text="变更时间：">
        <el-date-picker
          v-model="search.goldTime"
          type="datetimerange"
          range-separator="至"
          :default-time="['00:00:00', '23:59:59']"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
        >
        </el-date-picker>
      </search-item>
      <el-button
        size="small"
        type="primary"
        @click="handleSearch"
      >
        查询
      </el-button>
      <el-button
        class="mar-l-20"
        size="small"
        @click="handleReset"
      >
        重置
      </el-button>
      <el-button
        class="mar-l-20"
        size="small"
        @click="handleDownload"
      >
        导出
      </el-button>
    </template>
    <!-- 按钮区域 -->
    <template slot="operate">
      <el-button
        type="primary"
        @click="openGold"
      >手动补发金币</el-button>
    </template>
    <!-- 列表区域 -->
    <template slot="table">
      <el-table
        stripe
        :data="list"
        @sort-change="sortChange"
      >
        <el-table-column
          prop="member_user_id"
          width="150"
          label="导购ID"
        ></el-table-column>
        <el-table-column
          prop="member_name"
          width="130"
          label="导购姓名"
        ></el-table-column>
        <el-table-column
          prop="department_name"
          width="120"
          label="所在部门"
        ></el-table-column>
        <el-table-column
          width="100"
          label="类型"
          prop="change_type"
        ><template slot-scope="scope">
            <span>
              {{ scope.row.change_type === 1 ? '增加' : '减少' }}</span>
          </template>
        </el-table-column>
        <el-table-column
          width="150"
          label="金币数"
          prop="coins"
        ><template slot-scope="scope">
            <span>
              <span>{{ scope.row.coins > 0 ? '+' : '' }}</span>
              <span>{{ scope.row.coins }}</span>
            </span>
          </template>
        </el-table-column>
        <el-table-column
          width="120"
          label="来源"
          prop="source_type_str"
        ><template slot-scope="scope">
            <span>
              {{ scope.row.source_type_str }}</span>
          </template>
        </el-table-column>
        <el-table-column
          prop="coins_remark"
          width="250"
          label="变动原因"
          show-overflow-tooltip
        ></el-table-column>
        <el-table-column
          width="200"
          label="兑礼单号"
          fixed="right"
        >
          <template slot-scope="scope">
            {{ scope.row.order_num || "--" }}
          </template>
        </el-table-column>
        <el-table-column
          label="变更时间"
          width="200"
          fixed="right"
        >
          <template slot-scope="scope">
            {{ scope.row.created_time_str }}
          </template>
        </el-table-column>
        <div
          slot="empty"
          class="table-empty-label"
          v-if="!list.length"
        >
          <img
            src="@/assets/images/empty-box.png"
            alt=""
          />
          <p>
            暂无金币数据~
          </p>
        </div>
      </el-table>
    </template>

    <!-- 底部分页 -->
    <template slot="footer">
      <!-- 手动上传弹窗 -->
      <upload-gold
        :visible.sync="uploadGoldvisible"
        @success="uploadSuccess"
      />
      <el-pagination
        v-if="list && list.length > 0"
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
import { Component, Mixins } from "vue-property-decorator";
import SelectTreeDepartment from "@/components/business/select-tree-department/index.vue";
import { getGoldRecordList, exportGoldRecordList } from "@/api/we-journey/gold";
import { ApiResult } from "@/types";
import DownLoad from "@/mixins/DownLoad";
// 手动补发金币
import uploadGold from "../components/upload-gold/index.vue";
/** 生成搜索条件 */
function generateSearch() {
  return {
    /** 导购姓名 */
    user_name: "",
    /** 导购id */
    user_id: "",
    /** 企微部门id */
    third_id: 0,
    /** 类型 */
    change_type: 0,
    goldTime: [] as string[],
    /** 变更时间 */
    start_time: "" as string | number,
    end_time: "" as string | number
  };
}

/** 生成分页数据 */
function generateMeta() {
  return {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
}

@Component({
  components: { SelectTreeDepartment, uploadGold }
})
export default class GradeStatistics extends Mixins(DownLoad) {
  [x: string]: any;
  /** 积分来源options */
  levelOptions = [{
    id: 0,
    name: "全部"
  },
  {
    id: 1,
    name: "增加"
  }, {
    id: 2,
    name: "减少"
  }];

  /** 分页数据 */
  meta = generateMeta();

  /** 搜素内容 */
  search = generateSearch();

  /** 列表数据 */
  list = [];

  /** 列表排序 */
  sort = {
    prop: "",
    order: ""
  };
  /** 手动补发金币弹窗  */
  uploadGoldvisible = false;
  /** 列表排序 */
  sortChange({ prop, order }: { prop: string; order: string }) {
    // todo: 排序
    this.sort = { prop, order: order === "ascending" ? "asc" : "desc" };
    this.getList({ [prop]: order === "ascending" ? "asc" : "desc" });
  }

  /** 导出 */
  async handleDownload() {
    try {
      const params = this.mergeSearch();
      const res: any = await exportGoldRecordList(params);
      if (res.code === 200) {
        this.$confirm("文件生成中，稍后将自动下载，详情记录请前往【下载中心】查看。", "文件下载提示", {
          confirmButtonText: "去下载中心",
          cancelButtonText: "留在当前页",
          type: "success"
        })
          .then(() => {
            this.$router.push({
              path: "/myapp/setting/download-center/list"
            });
          })
          .catch(() => {
            this.$message({
              type: "info",
              message: "已取消前往下载中心"
            });
          });
      }
    } catch (e) {
      console.error(e);
    }
  }

  /** 选择部门变化 */
  departmentChange(id: number) {
    this.search.third_id = Number(id);
  }

  /** 点击查询 */
  handleSearch() {
    this.meta = generateMeta();
    this.getList();
    console.log("搜索参数为", this.search);
  }

  /** 点击重置 */
  handleReset() {
    // 部门组件重置
    (this.$refs as HTMLFormElement).selectTreeDepartment.reset();
    this.meta = generateMeta();
    this.search = generateSearch();
    console.log("点击了重置", this.search);
    this.getList();
  }

  /** 分页变化 - 每页 */
  handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.getList();
  }

  /** 分页变化 - 当前页 */
  handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getList();
  }
  /** 整合搜索数据 */
  mergeSearch() {
    let { goldTime } = this.search;

    goldTime = goldTime || [];

    const params: Partial<ReturnType<typeof generateSearch> &
      ReturnType<typeof generateMeta>> = {
      ...this.search,
      change_type: this.search.change_type,
      page: this.meta.page,
      size: this.meta.size,
      start_time: goldTime.length
        ? ~~((goldTime[0] as any).getTime() / 1000)
        : 0,
      end_time: goldTime.length
        ? ~~((goldTime[1] as any).getTime() / 1000)
        : 0,
    };

    delete params.goldTime;
    return params;
  }

  /** 获取列表数据 */
  async getList(other = {}) {
    const params = this.mergeSearch();
    console.log("获取列表", params);
    console.log("列表排序", this.sort);
    const sort = {};
    if (this.sort.prop) {
      sort[this.sort.prop] = this.sort.order;
    }

    const res: ApiResult = await getGoldRecordList({
      ...params,
      ...other,
      ...sort
    });

    if (res.code !== 200) return this.$message.error(res.message as string);
    this.list = res.data.list || [];
    this.meta.total = res.data.total;
  }

  /** 打开手动补发金币弹窗 */
  openGold() {
    this.uploadGoldvisible = true;
  }

  /** 生命周期初始化开始 */
  created() {
    // 获取路由传过来的导购id
    const { userId } = this.$route.query;
    console.log("路由参数", this.$route.query);
    if (userId) {
      this.search.user_id = userId as string;
    }
    this.getList();
  }
}
</script>
<style lang="scss" scoped>
.table-empty-label {
  margin-top: 50px;
}

.jump-column {
  color: #409eff;
  cursor: pointer;
}</style>
