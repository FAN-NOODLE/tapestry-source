<template>
  <layout-list>
    <!-- 搜索区域 -->
    <template slot="search">
      <search-item text="导购ID：">
        <cy-input v-model="search.user_id" placeholder="请输入导购ID" clearable>
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
          @change="departmentChange"
          ref="selectTreeDepartment"
        ></select-tree-department>
      </search-item>
      <search-item text="积分维度：">
        <cy-select
          v-model="search.incentive_type"
          size="medium"
          :optionsList="dimensionOption"
        >
        </cy-select>
      </search-item>
      <search-item text="积分来源：">
        <cy-select
          v-model="search.source_type"
          size="medium"
          :optionsList="integralSourceOption"
        >
        </cy-select>
      </search-item>
      <search-item text="关联规则名称：">
        <cy-input
          v-model="search.integral_rule_name"
          placeholder="请输入规则名称"
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
      <el-button class="mar-l-20" size="small" @click="downLoad">
        导出
      </el-button>
    </template>

    <!-- 按钮区域 -->
    <template slot="operate">
      <el-button type="primary" @click="openReissue">手动补发积分</el-button>
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
        <el-table-column width="130" label="关联规则名称"
          ><template slot-scope="scope">
            {{ scope.row.integral_rule_name || "--" }}
          </template>
        </el-table-column>
        <el-table-column width="100" label="积分维度"
          ><template slot-scope="scope">
            {{ getIncentiveName(scope.row.incentive_type) }}
          </template>
        </el-table-column>
        <el-table-column width="130" label="规则类目"
          ><template slot-scope="scope">
            {{ scope.row.rule_category_name || "--" }}
          </template>
        </el-table-column>
        <el-table-column
          prop="points_remark"
          width="130"
          label="积分获取条件"
        ></el-table-column>
        <el-table-column
          prop="points"
          width="130"
          label="积分"
        ></el-table-column>
        <el-table-column
          prop="coins"
          width="130"
          label="金币"
          v-if="epName === 'coach'"
        ></el-table-column>
        <el-table-column width="130" label="积分来源"
          ><template slot-scope="scope">
            {{ getRuleSourceName(scope.row.source_type) }}
          </template>
        </el-table-column>
        <el-table-column min-width="180" label="获取时间">
          <template slot-scope="scope">
            {{ scope.row.created_at | formatDate("yyyy-mm-dd hh:mm:ss") }}
          </template>
        </el-table-column>
        <div slot="empty" class="table-empty-label" v-if="!list.length">
          <img src="@/assets/images/empty-box.png" alt="" />
          <p>
            暂无积分兑领数据~
          </p>
        </div>
      </el-table>
    </template>

    <!-- 底部分页 -->
    <template slot="footer">
      <!-- 手动上传弹窗 -->
      <upload-reissue
        :visible.sync="reissueOption.visible"
        :type="1"
        @success="handleSearch"
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
import { Component, Mixins, Watch } from "vue-property-decorator";
import { interfaceDimension, integralSource } from "../config";
import SelectTreeDepartment from "@/components/business/select-tree-department/index.vue";
import UploadReissue from "../components/upload-reissue/index.vue";
import { getPointsListApi, downloadPointApi } from "@/api/we-journey/Integral";
import { ApiResult } from "@/types";
import DownLoad from "@/mixins/DownLoad";
import PerModule from "@/store/modules/permission";

/** 搜索条件type */
type searchType = {
  /** 导购姓名 */
  user_name: string;
  /** 导购id */
  user_id: string;
  /** 企微部门id */
  third_id: number;
  /** 关联规则名称 */
  integral_rule_name: string;
  /** 积分维度 */
  incentive_type: string | number;
  /** 积分来源 */
  source_type: string | number;
  /** 获取时间 */
  time?: Array<Date>;
};

/** 生成搜索条件 */
function generateSearch(): searchType {
  return {
    /** 导购姓名 */
    user_name: "",
    /** 导购id */
    user_id: "",
    /** 企微部门id */
    third_id: -1,
    /** 关联规则名称 */
    integral_rule_name: "",
    /** 积分维度 */
    incentive_type: 0,
    /** 积分来源 */
    source_type: 0,
    /** 获取时间 */
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
export default class IntegralList extends Mixins(DownLoad) {
  /** 搜索规则维度option */
  private dimensionOption = interfaceDimension;

  /** 积分来源options */
  private integralSourceOption = integralSource;

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

  // todo: 当前品牌
  private epName = PerModule.getEpName;

  /** 监听企业信息变化 */
  @Watch("$store.state.permission.EQ_NAME", { immediate: true })
  change(value: string) {
    if (!value) return;

    // Coach展示五种，KS&SW本期暂不展示培训维度及线上能力
    if (value !== "coach") {
      this.dimensionOption = interfaceDimension.filter(item =>
        [0, 1, 2, 3].includes(item.value)
      );
    }
  }

  /** 点击导出 */
  private async downLoad() {
    try {
      const params = this.mergeSearch();
      const res: any = await downloadPointApi(params);
      if(res.code ===200){
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

  /** 打开手动补发弹窗 */
  private openReissue() {
    this.reissueOption.visible = true;
  }

  /** 获取积分维度名称 */
  private getIncentiveName(id: number) {
    if (id === 0) return "--";
    const checkedItem = interfaceDimension.filter(item => item.value === id);
    return checkedItem.length ? checkedItem[0].label : "--";
  }

  /** 积分来源获取名称 */
  private getRuleSourceName(id: number) {
    if (!id) return "--";
    const checkedItem = this.integralSourceOption.filter(
      item => item.value === id
    );
    return checkedItem.length ? checkedItem[0].label : "--";
  }

  /** 选择部门变化 */
  private departmentChange(id: number) {
    this.search.third_id = id;
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
    // 别的组件重置
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

  /** 整合搜索参数 */
  private mergeSearch() {
    const time = this.search.time || [];
    const params = {
      ...this.search,
      third_id: Number(this.search.third_id),
      page: this.meta.page,
      size: this.meta.size,
      start_time: time.length ? ~~((time[0] as Date).getTime() / 1000) : 0,
      end_time: time.length ? ~~((time[1] as Date).getTime() / 1000) : 0
    };

    delete params.time;
    return params;
  }

  /** 获取列表数据 */
  private async getList() {
    const params = this.mergeSearch();
    const res: ApiResult = await getPointsListApi(params);

    if (res.code !== 200) return this.$message.error(res.message as string);
    this.list = res.data.list || [];
    this.meta.total = res.data.total;
  }

  /** 生命周期初始化 */
  private mounted() {
    const { type, userId } = this.$route.query;

    this.search.incentive_type = Number(type) || 0;
    this.search.user_id = (userId as string) || "";

    this.getList();
  }
}
</script>
<style lang="scss" scoped>
.table-empty-label {
  margin-top: 50px;
}
</style>
