<template>
  <layout-list>
    <!-- 搜索区域 -->
    <template slot="search">
      <search-item text="导购ID：">
        <cy-input
          v-model="search.member_user_id"
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
      <search-item text="关联奖章规则名称：">
        <cy-input
          v-model="search.medal_rule_name"
          placeholder="请输入关联奖章规则名称"
          clearable
        >
        </cy-input>
      </search-item>
      <search-item text="奖章来源：">
        <cy-select
          v-model="search.source"
          size="medium"
          :optionsList="downOptions.sourceOption"
        >
        </cy-select>
      </search-item>
      <search-item text="所在部门：">
        <select-tree-department
          ref="selectTreeDepartment"
          @change="departmentChange"
        ></select-tree-department>
      </search-item>
      <search-item text="奖章规则维度：">
        <cy-select
          v-model="search.rule_dimension_id"
          size="medium"
          :optionsList="downOptions.ruleOption"
        >
        </cy-select>
      </search-item>
      <search-item text="奖章名称：">
        <cy-input
          v-model="search.medal_name"
          placeholder="请输入奖章名称"
          clearable
        >
        </cy-input>
      </search-item>

      <search-item text="获取时间：">
        <el-date-picker
          v-model="search.time"
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

    <!-- 按钮区域 -->
    <template slot="operate">
      <el-button type="primary" @click="openReissue">手动补发奖章</el-button>
    </template>

    <!-- 列表区域 -->
    <template slot="table">
      <el-table stripe :data="list">
        <el-table-column
          prop="member_user_id"
          width="140"
          label="导购ID"
        ></el-table-column>
        <el-table-column
          prop="member_name"
          width="130"
          label="导购姓名"
        ></el-table-column>
        <el-table-column
          prop="department_name"
          width="130"
          label="所在部门"
        ></el-table-column>
        <el-table-column width="130" label="关联规则名称" show-overflow-tooltip
          ><template slot-scope="scope">
            {{ scope.row.rule_name || "--" }}
          </template>
        </el-table-column>
        <el-table-column
          prop="medal_name"
          width="130"
          show-overflow-tooltip
          label="奖章名称"
        ></el-table-column>
        <el-table-column width="100" label="奖章来源">
          <template slot-scope="scope">
            {{ scope.row.source === 1 ? "手动上传" : "规则统计" }}
          </template>
        </el-table-column>
        <el-table-column width="130" label="奖章规则维度">
          <template slot-scope="scope">
            {{ scope.row.rule_dimension_name || "--" }}
          </template>
        </el-table-column>
        <el-table-column width="140" label="奖章规则类目">
          <template slot-scope="scope">
            {{ scope.row.rule_category_name || "--" }}
          </template>
        </el-table-column>
        <el-table-column
          prop="integral"
          width="70"
          label="积分"
        ></el-table-column>
        <el-table-column
          prop="coins"
          width="70"
          label="金币"
          v-if="epName === 'coach'"
        ></el-table-column>
        <el-table-column width="130" label="实物奖品" show-overflow-tooltip>
          <template slot-scope="scope">
            {{ scope.row.item_name || "--" }}
          </template>
        </el-table-column>
        <el-table-column min-width="180" label="获取时间">
          <template slot-scope="scope">
            {{ scope.row.created_at | formatDate("yyyy-mm-dd hh:mm:ss") }}
          </template>
        </el-table-column>
        <el-table-column prop="zip" label="操作" width="100" fixed="right">
          <template slot-scope="scope">
            <span class="operate-blue" @click="handleRecovery(scope.row.id)"
              >收回奖章
            </span>
          </template>
        </el-table-column>
        <div slot="empty" class="table-empty-label" v-if="!list.length">
          <img src="@/assets/images/empty-box.png" alt="" />
          <p>
            暂无领取数据~
          </p>
        </div>
      </el-table>
    </template>

    <!-- 底部分页 -->
    <template slot="footer">
      <!-- 手动上传弹窗 -->
      <upload-reissue
        :visible.sync="reissueOption.visible"
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
import { OptionType, sourceOptions } from "../config";
import SelectTreeDepartment from "@/components/business/select-tree-department/index.vue";
import UploadReissue from "../components/upload-reissue/index.vue";
import { ApiResult } from "@/types";
import {
  getRuleOptionsApi,
  getMedalRecordListApi,
  downloadMedalApi,
  recordDeleteApi
} from "@/api/we-journey/medal";
import DownLoad from "@/mixins/DownLoad";
import PerModule from "@/store/modules/permission";


/** 给下拉框补充元素 */
function replenishOption(list: Array<OptionType>) {
  return [{ label: "全部", value: 0 }].concat(list);
}

/** 搜索条件type */
type searchType = {
  /** 导购姓名 */
  member_name: string;
  /** 导购id */
  member_user_id?: string;
  /** 企微部门id */
  department_id: number;
  /** 关联规则名称 */
  medal_rule_name: string;
  /** 奖章维度 */
  rule_dimension_id: number;
  /** 奖章 */
  source: number;
  /** 奖章名称 */
  medal_name: string;
  /** 获取时间 */
  time?: Array<Date>;
  /** 奖章列表id */
  medal_id?: number;
};

/** 生成搜索条件 */
function generateSearch(): searchType {
  return {
    member_name: "",
    member_user_id: "",
    department_id: 0,
    rule_dimension_id: 0,
    medal_rule_name: "",
    source: 0,
    medal_name: "",
    time: []
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
  components: { SelectTreeDepartment, UploadReissue }
})
export default class MedalStatistics extends Mixins(DownLoad) {
  /** 下拉框配置 */
  private downOptions = {
    /** 奖章来源 */
    sourceOption: replenishOption(sourceOptions),
    /** 奖章规则维度 */
    ruleOption: [] as Array<OptionType>
  };

  /** 分页数据 */
  private meta = generateMeta();

  /** 搜素内容 */
  private search = generateSearch();

  /** 手动上传参数 */
  private reissueOption = {
    visible: false
  };

  /** 列表数据 */
  private list = [];

  // 当前品牌
  private epName = PerModule.getEpName;

  /** 导出 */
  private async handleDownload() {
    try {
      const params = this.mergeSearch();
      const res = await downloadMedalApi(params);
      const nowTime = new Date().getTime();
      const fileName = "奖章领取数据_" + nowTime + ".xlsx";
      this.downLoadFile(res, fileName);
    } catch (e) {
      console.error(e);
    }
  }

  /** 打开手动补发弹窗 */
  private openReissue() {
    this.reissueOption.visible = true;
  }

  /** 选择部门变化 */
  private departmentChange(id: string) {
    this.search.department_id = Number(id);
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

  /** 获取业绩排行下拉 */
  private async getRuleOptions() {
    const res: ApiResult = await getRuleOptionsApi({ pid: 0 });
    if (res.code !== 200) return;

    const resetOption = res.data.list.map(
      (item: { id: number; rule_name: string }) => ({
        value: item.id,
        label: item.rule_name
      })
    );

    this.downOptions.ruleOption = replenishOption(resetOption);
  }

  /** 回收奖章 */
  private async handleRecovery(id: number) {
    try {
      const confirmRes = await this.$confirm(
        "确认收回该奖章，收回后当前导购已获得的奖章及其它奖励将被收回。",
        "收回奖章",
        {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning"
        }
      );
      if (confirmRes !== "confirm") return;

      // 请求回收奖章接口
      const res: ApiResult = await recordDeleteApi({ id });
      if (res.code !== 200) return this.$message.error(res.message as string);

      this.$message.success("操作成功");
      // 刷新列表
      this.getList();
    } catch (e) {
      console.error(e);
    }
  }

  /** 整合搜索数据 */
  private mergeSearch() {
    const time = this.search.time || [];
    const params = {
      ...this.search,
      page: this.meta.page,
      size: this.meta.size,
      start_time: time.length ? ~~((time[0] as Date).getTime() / 1000) : 0,
      end_time: time.length ? ~~((time[1] as Date).getTime() / 1000) : 0
    };

    // 判断是否从奖章列表跳转过来
    const listID = Number(this.$route.query.id);
    if (listID) {
      params.medal_id = listID;
    }
    delete params.time;
    return params;
  }

  /** 子页面手动补发成功 */
  private uploadSuccess() {
    this.getList();
  }

  /** 获取列表数据 */
  private async getList() {
    const params = this.mergeSearch();
    const res: ApiResult = await getMedalRecordListApi(params);

    if (res.code !== 200) return this.$message.error(res.message as string);
    this.list = res.data.list || [];
    this.meta.total = res.data.total;
  }

  /** 生命周期初始化开始 */
  private created() {
    this.getRuleOptions();
  }

  private mounted() {
    const { userId } = this.$route.query;

    this.search.member_user_id = (userId as string) || "";

    this.getList();
  }
}
</script>
<style lang="scss" scoped>
.table-empty-label {
  margin-top: 50px;
}
.operate-blue {
  cursor: pointer;
  color: #6f3112;
}
</style>
