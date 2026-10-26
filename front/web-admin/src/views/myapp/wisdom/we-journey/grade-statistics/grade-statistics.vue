<template>
  <layout-list>
    <!-- 搜索区域 -->
    <template slot="search">
      <search-item text="导购ID：">
        <cy-input
          v-model="search.member_userid"
          placeholder="请输入导购ID"
          clearable
        >
        </cy-input>
      </search-item>
      <search-item text="导购姓名：">
        <cy-input
          v-model="search.member_name"
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
      <search-item text="最近积分获取时间：">
        <el-date-picker
          v-model="search.pointTime"
          type="datetimerange"
          range-separator="至"
          :default-time="['00:00:00', '23:59:59']"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
        >
        </el-date-picker>
      </search-item>
      <search-item text="等级：">
        <cy-select
          v-model="search.level"
          size="medium"
          :optionsList="levelOptions"
          :options="{ label: 'level_name', value: 'level_name' }"
        >
        </cy-select>
      </search-item>
      <search-item text="最近等级升级时间：">
        <el-date-picker
          v-model="search.levelTime"
          type="datetimerange"
          range-separator="至"
          :default-time="['00:00:00', '23:59:59']"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
        >
        </el-date-picker>
      </search-item>
      <search-item text="最近奖章获取时间：">
        <el-date-picker
          v-model="search.medalTime"
          type="datetimerange"
          range-separator="至"
          :default-time="['00:00:00', '23:59:59']"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
        >
        </el-date-picker>
      </search-item>
      <el-button size="small" type="primary" @click="handleSearch">
        查询
      </el-button>
      <el-button class="mar-l-20" size="small" @click="handleReset">
        重置
      </el-button>
      <el-button class="mar-l-20" size="small" @click="handleDownload">
        导出
      </el-button>
    </template>

    <!-- 列表区域 -->
    <template slot="table">
      <el-table stripe :data="list" @sort-change="sortChange">
        <el-table-column
          prop="member_userid"
          width="130"
          label="导购ID"
        ></el-table-column>
        <el-table-column
          prop="member_name"
          width="130"
          label="导购姓名"
        ></el-table-column>
        <el-table-column
          prop="third_name"
          width="130"
          label="所在部门"
        ></el-table-column>
        <el-table-column
          width="100"
          label="积分"
          prop="point_value"
          sortable="custom"
          ><template slot-scope="scope">
            <span
              class="jump-column"
              @click="jumpIntegral(scope.row.member_userid)"
            >
              {{ scope.row.point_value }}</span
            >
          </template>
        </el-table-column>
        <el-table-column
          width="130"
          label="当前金币"
          prop="coins"
          sortable="custom"
          v-if="epName === 'coach'"
          ><template slot-scope="scope">
            <span
              class="jump-column"
              @click="jumpCoins(scope.row.member_userid)"
            >
              {{ scope.row.coins }}</span
            >
          </template>
        </el-table-column>
        <el-table-column
          prop="total_coins"
          width="130"
          label="累计总金币"
          sortable="custom"
          v-if="epName === 'coach'"
        ></el-table-column>
        <el-table-column
          prop="spend_coins"
          width="140"
          label="累计消耗金币"
          sortable="custom"
          v-if="epName === 'coach'"
        ></el-table-column>
        <el-table-column width="180" label="最近积分获取时间">
          <template slot-scope="scope">
            {{ scope.row.point_updated_at | formatDate("yyyy-mm-dd hh:mm:ss") }}
          </template>
        </el-table-column>
        <el-table-column width="100" label="等级">
          <template slot-scope="scope">
            {{ scope.row.level || "--" }}
          </template>
        </el-table-column>
        <el-table-column width="180" label="最近等级升级时间">
          <template slot-scope="scope">
            {{ scope.row.level_updated_at | formatDate("yyyy-mm-dd hh:mm:ss") }}
          </template>
        </el-table-column>
        <el-table-column
          width="100"
          label="奖章数"
          prop="medal_num"
          sortable="custom"
          ><template slot-scope="scope">
            <span
              class="jump-column"
              @click="jumpMedal(scope.row.member_userid)"
            >
              {{ scope.row.medal_num }}</span
            >
          </template>
        </el-table-column>
        <el-table-column min-width="180" label="最近奖章获取时间">
          <template slot-scope="scope">
            {{ scope.row.medal_updated_at | formatDate("yyyy-mm-dd hh:mm:ss") }}
          </template>
        </el-table-column>
        <!-- todo:历史最高等级 -->
        <el-table-column width="150" label="历史最高等级" fixed="right">
          <template slot-scope="scope">
            {{ scope.row.max_level || "--" }}
          </template>
        </el-table-column>
        <!-- todo：最高等级获取时间 -->
        <el-table-column width="180" label="最高等级获取时间" fixed="right">
          <template slot-scope="scope">
            {{ scope.row.max_level_time | formatDate("yyyy-mm-dd hh:mm:ss") }}
          </template>
        </el-table-column>
        <div slot="empty" class="table-empty-label" v-if="!list.length">
          <img src="@/assets/images/empty-box.png" alt="" />
          <p>
            暂无等级数据~
          </p>
        </div>
      </el-table>
    </template>

    <!-- 底部分页 -->
    <template slot="footer">
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
import {
  getMemberListApi,
  getLevelListApi,
  downloadTemplateApi
} from "@/api/we-journey/grade";
import { ApiResult } from "@/types";
import DownLoad from "@/mixins/DownLoad";
import PerModule from "@/store/modules/permission";


/** 生成搜索条件 */
function generateSearch() {
  return {
    /** 导购姓名 */
    member_name: "",
    /** 导购id */
    member_userid: "",
    /** 企微部门id */
    member_third_id: 0,
    /** 等级 */
    level: "全部",
    /** 获取时间 */
    pointTime: [] as Array<Date>,
    levelTime: [] as Array<Date>,
    medalTime: [] as Array<Date>,
    /** 积分筛选时间 */
    point_updated_start_time: "" as string | number,
    point_updated_end_time: "" as string | number,
    /** 等级筛选时间 */
    level_updated_start_time: "" as string | number,
    level_updated_end_time: "" as string | number,
    /** 奖章筛选时间 */
    medal_updated_start_time: "" as string | number,
    medal_updated_end_time: "" as string | number
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
  components: { SelectTreeDepartment }
})
export default class GradeStatistics extends Mixins(DownLoad) {
  /** 积分来源options */
  private levelOptions = [{ level_name: "全部" }];

  /** 分页数据 */
  private meta = generateMeta();

  /** 搜素内容 */
  private search = generateSearch();

  /** 列表数据 */
  private list = [];

  /** 列表排序 */
  private sort = {
    prop: "",
    order: ""
  };

  // 当前品牌
  private epName = PerModule.getEpName;

  /** 列表排序 */
  private sortChange({ prop, order }: { prop: string; order: string }) {
    // todo: 排序
    this.sort = { prop, order: order === "ascending" ? "asc" : "desc" };
    this.getList({ [prop]: order === "ascending" ? "asc" : "desc" });
  }

  /** 导出 */
  private async handleDownload() {
    try {
      const params = this.mergeSearch();
      const res = await downloadTemplateApi(params);
      const nowTime = new Date().getTime();
      const fileName = "导购等级数据_" + nowTime + ".xlsx";
      this.downLoadFile(res, fileName);
    } catch (e) {
      console.error(e);
    }
  }

  /** 选择部门变化 */
  private departmentChange(id: number) {
    this.search.member_third_id = Number(id);
  }

  /** 点击查询 */
  private handleSearch() {
    this.meta = generateMeta();
    this.getList();
    console.log("搜索参数为", this.search);
  }

  /** 点击重置 */
  private handleReset() {
    // 部门组件重置
    (this.$refs as HTMLFormElement).selectTreeDepartment.reset();
    this.meta = generateMeta();
    this.search = generateSearch();
    console.log("点击了重置", this.search);
    this.getList();
  }

  /** 分页变化 - 每页 */
  private handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.getList();
  }

  /** 分页变化 - 当前页 */
  private handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getList();
  }

  /** 跳转到积分明细 */
  private jumpIntegral(userId: string) {
    this.$router.push({
      name: "myapp_wisdom_we-journey_integral-exchange",
      query: {
        userId
      }
    });
  }

  /** 跳转到奖章明细 */
  private jumpMedal(userId: string) {
    this.$router.push({
      name: "myapp_wisdom_we-journey_medal-statistics",
      query: {
        userId
      }
    });
  }

  /** 整合搜索数据 */
  private mergeSearch() {
    let { pointTime, levelTime, medalTime } = this.search;

    pointTime = pointTime || [];
    medalTime = medalTime || [];
    levelTime = levelTime || [];

    const params: Partial<ReturnType<typeof generateSearch> &
      ReturnType<typeof generateMeta>> = {
      ...this.search,
      level: this.search.level === "全部" ? "" : this.search.level,
      page: this.meta.page,
      size: this.meta.size,
      point_updated_start_time: pointTime.length
        ? ~~((pointTime[0] as Date).getTime() / 1000)
        : 0,
      point_updated_end_time: pointTime.length
        ? ~~((pointTime[1] as Date).getTime() / 1000)
        : 0,
      level_updated_start_time: levelTime.length
        ? ~~((levelTime[0] as Date).getTime() / 1000)
        : 0,
      level_updated_end_time: levelTime.length
        ? ~~((levelTime[1] as Date).getTime() / 1000)
        : 0,
      medal_updated_start_time: medalTime.length
        ? ~~((medalTime[0] as Date).getTime() / 1000)
        : 0,
      medal_updated_end_time: medalTime.length
        ? ~~((medalTime[1] as Date).getTime() / 1000)
        : 0
    };

    delete params.pointTime;
    delete params.levelTime;
    delete params.medalTime;
    return params;
  }

  /** 获取列表数据 */
  private async getList(other = {}) {
    const params = this.mergeSearch();
    console.log("获取列表", params);
    console.log("列表排序", this.sort);
    const sort = {};
    if (this.sort.prop) {
      sort[this.sort.prop] = this.sort.order;
    }

    const res: ApiResult = await getMemberListApi({
      ...params,
      ...other,
      ...sort
    });

    if (res.code !== 200) return this.$message.error(res.message as string);
    this.list = res.data.list || [];
    this.meta.total = res.data.total;
  }

  /** 获取等级列表 */
  private async getLevelList() {
    const res: ApiResult = await getLevelListApi();
    if (res.code !== 200) return;

    const defaultOptions = [{ level_name: "全部" }];
    this.levelOptions = defaultOptions.concat(res.data) || defaultOptions;
  }
  /** 跳转到积分明细 */
  private jumpCoins(userId: string) {
    this.$router.push({
      name: "myapp_wisdom_we-journey_coins-index_gold-statistics",
      query: {
        userId
      }
    });
  }
  /** 生命周期初始化开始 */
  private created() {
    this.getLevelList();
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
}
</style>
