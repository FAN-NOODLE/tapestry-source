<template>
  <layout-list>
    <!-- 按钮区域 -->
    <template slot="buttons">
      <el-button type="primary" @click="goDetail('', 1)"
        >新增奖章规则</el-button
      >
    </template>

    <!-- 搜索区域 -->
    <template slot="search">
      <search-item text="奖章规则名称：">
        <cy-input v-model="search.name" placeholder="请输入规则名称" clearable>
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
      <search-item text="奖章规则维度：">
        <cy-select
          v-model="search.rule_dimension_id"
          size="medium"
          :optionsList="downOptions.ruleOption"
        >
        </cy-select>
      </search-item>
      <search-item text="其它奖励：">
        <cy-select
          v-model="search.other_rewards"
          size="medium"
          :optionsList="downOptions.otherOption"
        >
        </cy-select>
      </search-item>
      <search-item text="规则状态：">
        <cy-select
          v-model="search.is_effective"
          size="medium"
          :optionsList="downOptions.validOption"
        >
        </cy-select>
      </search-item>
      <search-item text="是否启用：">
        <cy-select
          v-model="search.status"
          size="medium"
          :optionsList="downOptions.enableOption"
        >
        </cy-select>
      </search-item>
      <search-item text="创建时间：">
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
      <search-item text="规则有效期：">
        <el-date-picker
          v-model="search.ruleTime"
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
    </template>

    <!-- 列表区域 -->
    <template slot="table">
      <el-table stripe :data="list">
        <el-table-column prop="id" width="180" label="ID"></el-table-column>
        <el-table-column
          prop="name"
          width="130"
          label="奖章规则名称"
        ></el-table-column>
        <el-table-column width="130" label="奖章来源">
          <template slot-scope="scope">
            {{ getName("sourceOption", scope.row.source) }}
          </template>
        </el-table-column>
        <el-table-column width="130" label="奖章规则维度">
          <template slot-scope="scope">
            {{ scope.row.medal_dimension_name || "--" }}
          </template>
        </el-table-column>
        <el-table-column width="130" label="奖章规则类目">
          <template slot-scope="scope">
            {{ scope.row.rule_category_name || "--" }}
          </template>
        </el-table-column>
        <el-table-column width="130" label="其它奖励">
          <template slot-scope="scope">
            {{ getName("otherOption", scope.row.other_rewards) }}
          </template>
        </el-table-column>
        <el-table-column width="130" label="奖章类型">
          <template slot-scope="scope">
            {{ getName("medalTypeOption", scope.row.meda_type) }}
          </template>
        </el-table-column>
        <el-table-column width="130" label="规则状态">
          <template slot-scope="scope">
            {{
              scope.row.source === 1
                ? "--"
                : getName("validOption", scope.row.is_effective)
            }}
          </template>
        </el-table-column>
        <el-table-column width="200" label="规则有效期">
          <template slot-scope="scope">
            <span v-if="scope.row.rule_type === 1">长期有效</span>
            <span v-else-if="scope.row.rule_type === 2">
              {{
                scope.row.rule_start_time | formatDate("yyyy-mm-dd hh:mm:ss")
              }}
              ~
              {{ scope.row.rule_end_time | formatDate("yyyy-mm-dd hh:mm:ss") }}
            </span>
            <span v-else>仅生效一次</span>
          </template>
        </el-table-column>
        <el-table-column width="180" label="创建时间">
          <template slot-scope="scope">
            {{ scope.row.create_time | formatDate("yyyy-mm-dd hh:mm:ss") }}
          </template>
        </el-table-column>
        <el-table-column width="120" label="是否启用">
          <template slot-scope="scope">
            <el-switch
              @change="changeStatus(scope.row)"
              v-model="scope.row.status"
              active-color="#000000"
              inactive-color="#909399"
              :active-value="1"
              :inactive-value="2"
            >
            </el-switch>
          </template>
        </el-table-column>
        <el-table-column prop="zip" label="操作" width="250" fixed="right">
          <template slot-scope="scope">
            <el-button
              type="text"
              class="operate-blue"
              @click="goDetail(scope.row.id, 0)"
              >查看
            </el-button>
            <!-- 仅一次生效及已失效的奖章规则【编辑】置灰不可点 -->
            <el-button
              :disabled="scope.row.rule_type === 3"
              type="text"
              class="operate-blue"
              @click="goDetail(scope.row.id, 2)"
              >编辑
            </el-button>
            <el-button
              type="text"
              class="operate-blue"
              @click="goDetail(scope.row.id, 1)"
              >复制
            </el-button>
            <span
              type="text"
              class="operate-blue"
              @click="goStatistics(scope.row.id)"
              >查看数据
            </span>
            <el-button
              type="text"
              v-if="scope.row.status === 2"
              class="operate-error"
              @click="handleDel(scope.row.id)"
              >删除
            </el-button>
          </template>
        </el-table-column>
        <div slot="empty" class="table-empty-label" v-if="!list.length">
          <img src="@/assets/images/empty-box.png" alt="" />
          <p>
            暂无内容
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
import { ApiResult } from "@/types";
import { Component, Vue } from "vue-property-decorator";
import {
  sourceOptions,
  OptionType,
  otherOptions,
  validOptions,
  enableOptions,
  medalTypeOptions
} from "../config";
import {
  getRuleOptionsApi,
  getMedalListApi,
  medalDeleteApi,
  medalStatusApi
} from "@/api/we-journey/medal";

/** 给下拉框补充元素 */
function replenishOption(list: Array<OptionType>) {
  return [{ label: "全部", value: 0 }].concat(list);
}

/** 搜索条件类型 */
type searchType = {
  /** 奖章规则名称 */
  name: string;
  /** 奖章来源 */
  source: number;
  /** 奖章规则维度 */
  rule_dimension_id: number;
  /** 其它奖励 */
  other_rewards: number;
  /** 是否禁用 */
  status: number;
  /** 规则状态 */
  is_effective: number;
  /** 创建时间 */
  time?: Array<Date>;
  /** 规则有效期 */
  ruleTime?: Array<Date>;
};

/** 生成搜索条件 */
function generateSearch(): searchType {
  return {
    name: "",
    source: 0,
    rule_dimension_id: 0,
    other_rewards: 0,
    status: 0,
    is_effective: 0,
    time: [],
    ruleTime: []
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

@Component({})
export default class MedalList extends Vue {
  /** 下拉框配置 */
  private downOptions = {
    /** 奖章来源 */
    sourceOption: replenishOption(sourceOptions),
    /** 奖章规则维度 */
    ruleOption: [] as Array<OptionType>,
    /**  其它奖励 */
    otherOption: replenishOption(otherOptions),
    /** 是否有效 */
    validOption: replenishOption(validOptions),
    /** 启用禁用 */
    enableOption: replenishOption(enableOptions),
    /** 奖章类型 */
    medalTypeOption: medalTypeOptions
  };

  /** 分页数据 */
  private meta = generateMeta();

  /** 搜素内容 */
  private search = generateSearch();

  /** 列表数据 */
  private list = [] as Array<{ status: number }>;

  /** 表格根据value 查询label */
  private getName(optionName: string, value: number) {
    const checkedItem = this.downOptions[optionName].filter(
      (item: OptionType) => item.value === value
    );
    return checkedItem.length ? checkedItem[0].label : "--";
  }

  /** 跳转到对应页面/新增/编辑/详情 */
  private goDetail(id: number, index: number) {
    const names = ["look", "add", "edit"];
    this.$router.push({
      name: `myapp_wisdom_we-journey_medal-set-${names[index]}`,
      query: { id: id + "" }
    });
  }

  /** 查看数据 */
  private goStatistics(id: number) {
    this.$router.push({
      name: `myapp_wisdom_we-journey_medal-statistics`,
      query: { id: id + "" }
    });
  }

  /** 删除 */
  private async handleDel(id: number) {
    try {
      const confirmRes = await this.$confirm(
        "此操作将删除该规则, 是否继续?",
        "提示",
        {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning"
        }
      );
      if (confirmRes !== "confirm") return;

      // 请求删除接口
      const res: ApiResult = await medalDeleteApi({ id });
      if (res.code !== 200) return this.$message.error(res.message as string);

      this.$message.success("删除成功");
      // 刷新列表
      this.getList();
    } catch (e) {
      console.error(e);
    }
  }

  /** 改变启用禁用状态 */
  private async changeStatus(row: { status: number; id: number }) {
    const res: ApiResult = await medalStatusApi({
      id: row.id,
      status: row.status
    });

    if (res.code !== 200) {
      row.status = row.status === 1 ? 2 : 1;
      this.$message.error(res.message as string);
    }

    // 手动更改状态- 不管成功失败都需要更改状态， 因为dom 是v-model 绑定
    const status = row.status === 1 ? 1 : 2;
    row.status = status;
  }

  /** 点击查询 */
  private handleSearch() {
    this.meta = generateMeta();
    this.getList();
  }

  /** 点击重置 */
  private handleReset() {
    this.meta = generateMeta();
    this.search = generateSearch();
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
    if (res.code !== 200) return this.$message.error(res.message as string);

    const resetOption = res.data.list
      ? res.data.list.map((item: { id: number; rule_name: string }) => ({
          value: item.id,
          label: item.rule_name
        }))
      : [];

    this.downOptions.ruleOption = replenishOption(resetOption);
  }

  /** 获取列表详情 */
  private async getList() {
    // 整合查询参数
    const time = this.search.time || [];
    const ruleTime = this.search.ruleTime || [];
    const params = {
      ...this.search,
      page: this.meta.page,
      size: this.meta.size,
      start_time: time.length ? ~~(time[0].getTime() / 1000) : 0,
      end_time: time.length ? ~~(time[1].getTime() / 1000) : 0,
      rule_start_time: ruleTime.length ? ~~(ruleTime[0].getTime() / 1000) : 0,
      rule_end_time: ruleTime.length ? ~~(ruleTime[1].getTime() / 1000) : 0
    };

    delete params.time;
    delete params.ruleTime;
    console.log("搜索参数", params);
    const res: ApiResult = await getMedalListApi(params);

    // 判断请求是否成功
    if (res.code !== 200) return this.$message.error(res.message as string);
    this.list = res.data.list || [];
    this.meta.total = res.data.total;
  }

  /** 生命周期初始化开始 */
  private created() {
    this.getRuleOptions();
    this.getList();
  }
}
</script>
<style lang="scss" scoped>
.table-empty-label {
  margin-top: 50px;
}

.operate-error {
  color: #f56c6c;
  cursor: pointer;
  margin-right: 6px;
  margin-left: 0 !important;
  &[disabled] {
    color: #909399 !important;
    cursor: not-allowed;
    margin-right: 6px;
  }
}

.operate-blue {
  color: #6f3112;
  cursor: pointer;
  margin-right: 6px;
  margin-left: 0 !important;
  &[disabled] {
    color: #909399;
    cursor: not-allowed;
    margin-right: 6px;
  }
}
</style>
