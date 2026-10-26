<template>
  <layout-list class="integral-set-list">
    <!-- 按钮区域 -->
    <template slot="buttons">
      <el-button type="primary" @click="openAddRule" :disabled="disabledAddBtn"
        >新增积分规则</el-button
      >
      <div class="tips">
        注：同一【规则维度】的规则仅可创建一条，当前列表最多可新建{{
          tipNum
        }}条规则；创建后可在【编辑积分规则】中配置不同的规则类目条件
      </div>
    </template>

    <!-- 搜索区域 -->
    <template slot="search">
      <search-item text="规则名称：">
        <cy-input
          v-model="search.rule_name"
          placeholder="请输入规则名称"
          clearable
        >
        </cy-input>
      </search-item>
      <search-item text="规则维度：">
        <cy-select
          v-model="search.type"
          size="medium"
          :optionsList="dimensionOption"
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
        <el-table-column prop="id" width="80" label="ID"></el-table-column>
        <el-table-column
          prop="rule_name"
          width="130"
          label="规则名称"
        ></el-table-column>
        <el-table-column
          prop="type_name"
          width="130"
          label="规则维度"
        ></el-table-column>
        <el-table-column
          prop="rule_num"
          width="110"
          label="规则类目数"
        ></el-table-column>
        <el-table-column prop="status" width="130" label="关联能力图谱">
          <template slot-scope="scope">
            {{ scope.row.status === 1 ? "已开启" : "未开启" }}
          </template>
        </el-table-column>
        <el-table-column width="180" label="创建时间">
          <template slot-scope="scope">
            {{ scope.row.created_at | formatDate("yyyy-mm-dd hh:mm:ss") }}
          </template>
        </el-table-column>
        <el-table-column min-width="180" label="最近更新时间">
          <template slot-scope="scope">
            {{ scope.row.updated_at | formatDate("yyyy-mm-dd hh:mm:ss") }}
          </template>
        </el-table-column>
        <el-table-column prop="zip" label="操作" width="330" fixed="right">
          <template slot-scope="scope">
            <span class="c-blue" @click="goDetail(scope.row.id)">查看 </span>
            <span class="c-blue" @click="goDetail(scope.row.id, 2)"
              >编辑积分规则
            </span>
            <span
              class="c-blue"
              @click="
                relationAtlas(
                  scope.row.type,
                  scope.row.status,
                  scope.row.relation_id
                )
              "
              >{{ scope.row.status === 1 ? "编辑" : "关联" }}能力图谱
            </span>
            <el-dropdown>
              <span class="el-dropdown-link c-blue">
                更多<i class="el-icon-arrow-down el-icon--right"></i>
              </span>
              <el-dropdown-menu slot="dropdown">
                <el-dropdown-item
                  class="operate-blue"
                  @click.native="handleGoExplain(scope.row.type)"
                >
                  查看数据</el-dropdown-item
                >
                <el-dropdown-item
                  class="operation-error"
                  @click.native="handleDel(scope.row.id)"
                  >删除</el-dropdown-item
                >
              </el-dropdown-menu>
            </el-dropdown>
          </template>
        </el-table-column>
        <div slot="empty" class="table-empty-label" v-if="!list.length">
          <img src="@/assets/images/empty-box.png" alt="" />
          <p>
            您还没有创建规则，快去创建吧~
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

      <!-- 弹窗区域 -->
      <add-dialog
        :visible.sync="addDialogOptions.visible"
        :list="addDialogOptions.list"
        @relationAtlas="relationAtlas"
      ></add-dialog>

      <!-- 能力图谱弹窗 -->
      <atlas-dialog
        :type="atlasDialogOptions.modelId"
        :info="atlasDialogOptions.info"
        :visible.sync="atlasDialogOptions.visible"
        @success="getList"
      ></atlas-dialog>
    </template>
  </layout-list>
</template>
<script lang="ts">
import { Component, Vue, Watch } from "vue-property-decorator";
import AddDialog from "./components/add-dialog.vue";
import AtlasDialog from "./components/atlas-dialog.vue";
import { integralConfig, integralType, interfaceDimension } from "../config";
import {
  getIntegralListApi,
  getIntegralDeleteApi,
  getRelationDetailApi,
  getIncentiveConfigApi
} from "@/api/we-journey/Integral";
import { ApiResult } from "@/types";

/** 积分维度变化 */
let integralDefaultConfig = integralConfig;

/** 列表类型 */
type IntegralItemType = {
  id: number;
  /** 规则名称 */
  rule_name: string;
  /** 维度id */
  type: number;
  /** 维度名称 */
  type_name: string;
  /** 开启状态 1开启 2关闭 */
  status: number;
};

/** 搜索条件类型 */
type searchType = {
  /** 规则名称 */
  rule_name: string;
  /** 规则维度 */
  type: number;
  /** 时间 */
  time?: Array<Date>;
};

/** 生成搜索条件 */
function generateSearch(): searchType {
  return {
    rule_name: "",
    type: 0,
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
  components: { AddDialog, AtlasDialog }
})
export default class IntegralList extends Vue {
  /** 搜索规则维度option */
  private dimensionOption = interfaceDimension;

  /** 是否打开新增积分规则弹窗 */
  private addDialogOptions = {
    visible: false,
    list: [] as Array<integralType>
  };

  /** 提示文案显示数字 */
  private tipNum = 5;

  /** 监听企业信息变化 */
  @Watch("$store.state.permission.EQ_NAME", { immediate: true })
  change(value: string) {
    if (!value) return;

    // Coach展示五种，KS&SW本期暂不展示培训维度及线上能力
    if (value !== "coach") {
      integralDefaultConfig = integralDefaultConfig.filter(item =>
        [1, 2, 3].includes(item.id)
      );

      this.dimensionOption = interfaceDimension.filter(item =>
        [0, 1, 2, 3].includes(item.value)
      );

      this.tipNum = 3;
    }
  }

  /** 分页数据 */
  private meta = generateMeta();

  /** 能力图谱options */
  private atlasDialogOptions = {
    visible: false,
    info: {},
    modelId: 0 // 模块（维度）id - 区分维度
  };

  /** 搜素内容 */
  private search = generateSearch();

  /** 列表数据 */
  private list = [] as Array<IntegralItemType>;

  /** 是否置灰新增按钮 */
  private disabledAddBtn = false;

  /** 点击查询 */
  private handleSearch() {
    console.log("搜索参数为", this.search);
    this.meta = generateMeta();
    this.getList();
  }

  /** 点击重置 */
  private handleReset() {
    this.meta = generateMeta();
    this.search = generateSearch();
    this.getList();
  }

  /** 点击新增积分规则按钮 */
  private async openAddRule() {
    const res: ApiResult = await getIncentiveConfigApi();

    if (res.code !== 200) return this.$message(res.message as string);

    // 处理子组件所需要的数据，从列表中筛选
    this.addDialogOptions.list = integralDefaultConfig.map(item => {
      const checkedItem = res.data.filter(
        (listItem: { type: number }) => listItem.type === item.id
      );

      return {
        ...item,
        listId: checkedItem.length ? checkedItem[0]["integral_rule"].id : 0,
        relationId: checkedItem.length
          ? checkedItem[0]["incentive_relation"].id
          : item.relationId,
        relationStatus: checkedItem.length
          ? checkedItem[0]["incentive_relation"].status
          : 2
      };
    });
    // 打开弹窗、
    this.addDialogOptions.visible = true;
  }

  /** 点击查看 - table 操作 */
  private goDetail(id: number, type = 1) {
    // type === 1 表示查看 其他表示编辑
    const name =
      type === 1
        ? "myapp_wisdom_we-journey_integral-set-look"
        : "myapp_wisdom_we-journey_integral-set-edit";
    this.$router.push({
      name,
      query: {
        id: id + ""
      }
    });
  }

  /** 关联能力图谱 - 通过emit 触发/按钮触发 */
  private async relationAtlas(modelId: number, status?: number, id?: number) {
    // 判断是否已经存在能力图谱，获取旧数据
    let info = {};
    if (id) {
      const res: ApiResult = await getRelationDetailApi({ id });
      if (res.code !== 200) return this.$message.error(res.message as string);

      info = {
        ...res.data.relation_detail,
        relation_rule: res.data.relation_rule_details
      };
    }
    this.atlasDialogOptions.modelId = modelId;
    this.atlasDialogOptions.visible = true;
    this.atlasDialogOptions.info = info;
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

  /** 获取数据列表 */
  private async getList() {
    this.getTotal();
    const time = this.search.time || [];
    const params = {
      ...this.search,
      page: this.meta.page,
      size: this.meta.size,
      start_time: time.length ? ~~((time[0] as Date).getTime() / 1000) : 0,
      end_time: time.length ? ~~((time[1] as Date).getTime() / 1000) : 0
    };

    delete params.time;
    const res: ApiResult = await getIntegralListApi(params);
    if (res.code !== 200) return this.$message.error(res.message as string);
    this.list = res.data.list || [];
    this.meta.total = res.data.total;
  }

  /** 判断是否可以点击新增按钮 */
  private async getTotal() {
    const res: ApiResult = await getIntegralListApi({});
    if (res.code !== 200) return this.$message.error(res.message as string);
    // eslint-disable-next-line no-constant-condition
    this.disabledAddBtn =
      this.$store.state.permission.EQ_NAME === "coach"
        ? res.data.total === 5
        : res.data.total === 3;
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
      const res: ApiResult = await getIntegralDeleteApi({
        id
      });
      if (res.code !== 200) return this.$message.error(res.message as string);

      this.$message.success("删除成功");
      // 刷新列表
      this.getList();
    } catch (e) {
      console.error(e);
    }
  }

  /** 跳转到数据明细 */
  private handleGoExplain(type: number) {
    this.$router.push({
      name: "myapp_wisdom_we-journey_integral-exchange",
      query: {
        type: type + ""
      }
    });
  }

  /** 生命周期 */
  private created() {
    this.getList();
  }
}
</script>
<style lang="scss" scoped>
.table-empty-label {
  margin-top: 50px;
}

.operation-error {
  color: #f56c6c !important;
}

.operate-blue {
  color: #6f3112;
}

.integral-set-list {
  .c-blue {
    cursor: pointer;
    color: #6f3112;
    margin-right: 8px;
  }
}

.tips {
  margin-top: 10px;
  color: #909399;
  font-size: 12px;
}
</style>
