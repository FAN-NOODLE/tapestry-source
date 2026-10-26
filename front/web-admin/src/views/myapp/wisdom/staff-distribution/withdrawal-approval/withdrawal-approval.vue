<template>
  <div class="withdrawal-approval">
    <layout-list>
      <!-- 顶部tab -->
      <template slot="top-tabs">
        <el-tabs v-model="active_name" @tab-click="handleTabClick">
          <el-tab-pane
            v-for="tab in tabList"
            :key="tab.name"
            :label="tab.label"
            :name="tab.name"
          ></el-tab-pane>
        </el-tabs>
      </template>

      <!-- 搜索部分 -->
      <template slot="search">
        <search-item text="筛选：">
          <el-input
            :placeholder="`请输入${search_types[search_type - 1].remind}`"
            class="search-input"
            size="medium"
            clearable
            v-model="search_name"
            @keyup.enter.native="search"
          >
            <el-select
              slot="prepend"
              v-model="search_type"
              class="input-with-select"
            >
              <el-option
                v-for="item in search_types"
                :value="item.key"
                :label="item.name"
                :key="item.key"
              ></el-option>
            </el-select>
          </el-input>
        </search-item>
        <search-item text="时间：">
          <el-date-picker
            v-model="date_value"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            value-format="timestamp"
            :default-time="['00:00:00', '23:59:59']"
            size="medium"
          >
          </el-date-picker>
        </search-item>
        <search-item text="导购：">
          <search-select-member
            :selectedList="hign_choose"
            @visible-change="highSearch"
            @remove-member="removeGuide"
          >
          </search-select-member>
        </search-item>
        <div class="search-button">
          <el-button type="primary" size="small" @click="search">
            查询
          </el-button>
          <el-button class="m-l-20" size="small" @click="reset">重置</el-button>
        </div>
      </template>

      <!-- 操作部分 -->
      <template slot="operate">
        <el-button
          v-if="active_name == 1"
          size="small"
          @click="approvalCase(1)"
        >
          批量批准
        </el-button>
        <el-button v-if="active_name == 1" size="small" @click="rejectCase(1)">
          批量驳回
        </el-button>
        <el-button size="small" @click="exportList">批量导出</el-button>
        <el-button v-if="active_name == 1" size="small" @click="onSetting">
          自动批准设置
        </el-button>
      </template>

      <!-- 表格部分 -->
      <template slot="table">
        <el-table
          class="list-last-line-none min-height-100"
          v-if="page_init"
          :data="list"
          size="mini"
          stripe
          @selection-change="handleSelectionChange"
          :row-key="getRowKeys"
          :header-cell-style="{
            fontSize: '14px',
            color: '#606266',
            fontWeight: '400',
            padding: '11px 0 10px 0'
          }"
        >
          <el-table-column
            type="selection"
            width="55"
            :reserve-selection="true"
          ></el-table-column>
          <el-table-column
            prop="id"
            label="提现申请单号"
            min-width="140"
          ></el-table-column>
          <el-table-column
            prop="application_amount"
            label="申请金额(元)"
            min-width="140"
            align="right"
          >
            <div slot-scope="{ row }">{{ row.application_amount / 100 }}</div>
          </el-table-column>
          <el-table-column
            prop="worker_id"
            label="导购ID"
            min-width="150"
          ></el-table-column>
          <el-table-column prop="is_sale" label="申请导购" min-width="140">
            <template slot-scope="scope">
              <div class="guide-name">
                <img
                  class="table-text-img"
                  :src="
                    scope.row.avatar
                      ? scope.row.avatar
                      : 'https://oss-crs.vchangyi.com/wechat/default-avatar-icon.png'
                  "
                  alt=""
                  width="30"
                  height="30"
                />
                <span v-text="scope.row.name"></span>
              </div>
            </template>
          </el-table-column>
          <el-table-column
            prop="mobile"
            label="导购手机号"
            min-width="130"
          ></el-table-column>
          <el-table-column
            prop="dp_store"
            :label="USER.version !== 3 ? '导购所在门店/部门' : '导购所在部门'"
            min-width="160"
            show-overflow-tooltip="show-overflow-tooltip"
          >
            <template slot-scope="scope">
              <span>{{ scope.row.dp_store ? scope.row.dp_store : "--" }}</span>
            </template>
          </el-table-column>
          <el-table-column label="申请发起时间" min-width="170">
            <template slot-scope="scope">{{
              scope.row.request_initiated_at | formatDate("yyyy-mm-dd hh:mm:ss")
            }}</template>
          </el-table-column>
          <el-table-column
            v-if="active_name != 1 && active_name != 5"
            label="申请处理时间"
            min-width="170"
          >
            <template slot-scope="scope">{{
              scope.row.request_handle_at | formatDate("yyyy-mm-dd hh:mm:ss")
            }}</template>
          </el-table-column>
          <el-table-column
            v-if="active_name == 4"
            prop="reason"
            label="驳回原因"
            min-width="260"
          >
            <template slot-scope="scope">
              <el-popover
                v-if="scope.row.reason.length > 20"
                placement="top"
                trigger="hover"
                width="250"
              >
                <div slot="reference" class="skuname character-hiding-o">
                  {{ scope.row.reason }}
                </div>
                <span>{{ scope.row.reason }}</span>
              </el-popover>
              <span v-else>{{ scope.row.reason }}</span>
            </template>
          </el-table-column>
          <el-table-column
            v-if="active_name == 3"
            prop="reason"
            label="提现失败原因"
            min-width="290"
          >
            <template slot-scope="scope">
              <el-popover
                placement="top"
                trigger="hover"
                width="250"
                v-if="scope.row.reason.length > 20"
              >
                <div slot="reference" class="skuname character-hiding-o">
                  {{ scope.row.reason }}
                </div>
                <span>{{ scope.row.reason }}</span>
              </el-popover>
              <span v-else>{{ scope.row.reason }}</span>
            </template>
          </el-table-column>
          <el-table-column
            v-if="active_name == 1"
            fixed="right"
            label="操作"
            min-width="120"
          >
            <template slot-scope="scope">
              <el-button
                @click="approvalCase(2, scope.row)"
                type="text"
                size="mini"
                class="operation-success"
              >
                批准
              </el-button>
              <el-button
                @click="rejectCase(2, scope.row)"
                type="text"
                size="mini"
                class="operation-success"
              >
                驳回
              </el-button>
            </template>
          </el-table-column>
          <div slot="empty">
            <layout-empty
              :message="'暂无待审批申请'"
              v-if="active_name == 1 && !searchLoading"
            ></layout-empty>
            <layout-empty
              :message="'暂无转账中记录'"
              v-if="active_name == 5 && !searchLoading"
            ></layout-empty>
            <layout-empty
              :message="'暂无提现成功记录'"
              v-if="active_name == 2 && !searchLoading"
            ></layout-empty>
            <layout-empty
              :message="'暂无提现失败记录'"
              v-if="active_name == 3 && !searchLoading"
            ></layout-empty>
            <layout-empty
              :message="'暂无已驳回记录'"
              v-if="active_name == 4 && !searchLoading"
            ></layout-empty>
            <layout-empty
              v-if="searchLoading"
              empty_type="tag"
              :message="'没有找到相关数据'"
            ></layout-empty>
          </div>
        </el-table>

        <!-- 自动审批设置 -->
        <auto-setting
          v-if="isSetting"
          :visible="isSetting"
          @done="successSetting"
          @close="clearAuto"
        ></auto-setting>

        <!-- 驳回弹框 -->
        <reject-case
          :visible="visible"
          :case_id="case_id"
          @success="fetchList"
          @close="clearCase"
        ></reject-case>

        <!-- 选择导购 -->
        <select-member
          v-if="member_show"
          :visible="member_show"
          dialog_title="选择导购"
          choice_name="已选导购"
          :selectArr="hign_choose"
          :tabName="USER.version == 3 ? 3 : 1"
          :is_show_guide="true"
          @close="member_show = false"
          @submit="submitDialog"
        ></select-member>
      </template>

      <template slot="footer">
        <el-pagination
          class="page-right font-weight-400"
          v-if="list && list.length !== 0 && page_init"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page="meta.page"
          :page-sizes="[10, 20, 30, 40]"
          :page-size="meta.size"
          layout="total, sizes, prev, pager, next, jumper"
          :total="meta.total"
        ></el-pagination>
      </template>
    </layout-list>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import UserModule from "@/store/modules/user";
import Api from "@/api/wisdom/Withdrawal";
import SelectMember from "@/components/sy-components/select-member.vue";
import SearchSelectMember from "@/components/sy-components/search-select-member.vue";
import LayoutEmpty from "@/components/sy-components/layout-empty.vue";
import RejectCase from "./model/reject-case.vue";
import AutoSetting from "./model/auto-setting.vue";
import ShareModule from "@/shared/index";

@Component({
  components: {
    SelectMember,
    SearchSelectMember,
    LayoutEmpty,
    RejectCase,
    AutoSetting
  }
})
export default class extends Vue {
  private tabList: Array<any> = [
    { label: "待审批", name: "1" },
    { label: "转账中", name: "5" },
    { label: "提现成功", name: "2" },
    { label: "提现失败", name: "3" },
    { label: "已驳回", name: "4" }
  ];
  private page_init = false; // 页面初始状态
  private meta = {
    total: 40, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  private active_name = "1";
  private searchLoading = false; // 是否正在搜索
  private list: Array<any> = []; // 列表
  private search_name = ""; // 搜索活动名称
  private search_type = 1; // 搜索类型 1单号 2手机号
  private date_value: Array<any> = []; // 日期数据
  private search_types: Array<any> = [
    { name: "申请单号", key: 1, label: "id", remind: "申请单号" },
    { name: "申请人", key: 2, label: "mem_user_name", remind: "申请姓名" },
    {
      name: "导购手机号",
      key: 3,
      label: "mem_user_phone",
      remind: "导购手机号"
    }
  ];
  private visible = false;
  private case_id: Array<any> = [];
  private member_show = false;
  private hign_choose: Array<any> = [];
  private guide_workers: Array<any> = [];
  private omit = "...";
  private tab_name = 1;
  private isSetting = false;

  get USER() {
    return UserModule.getUserInfo;
  }

  private successSetting(data: any) {
    this.isSetting = false;
    if (data) {
      Api.saveAutoSetting(data).then((res: any) => {
        if (res.code == 200) {
          this.$message.success("设置成功");
        }
      });
    }
  }

  private clearAuto() {
    this.isSetting = false;
  }

  private onSetting() {
    this.isSetting = true;
  }

  private clearCase() {
    this.case_id = [];
  }

  // 高级搜索
  private highSearch() {
    this.member_show = true;
  }

  // 高级搜索选择完成
  private submitDialog(list: any) {
    this.hign_choose = list;
    this.guide_workers = list.map((v: any) => v.id);
    this.fetchList();
  }

  // 多选选择
  private handleSelectionChange(val: Array<any>) {
    this.case_id = [...val];
  }

  // 返回Id
  private getRowKeys(row: any) {
    return row.id;
  }

  //  搜索
  private search() {
    this.meta.page = 1;
    this.searchLoading = true;
    this.fetchList();
  }

  // 重置
  private reset() {
    this.search_type = 1;
    this.search_name = "";
    this.date_value = [];
    this.hign_choose = [];
    this.searchLoading = false;
    this.guide_workers = [];
    this.tab_name = 1;
    this.meta.page = 1;
    this.fetchList();
  }

  // 每页
  private handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.fetchList();
  }

  // 当前页
  private handleCurrentChange(val: number) {
    this.meta.page = val;
    this.fetchList();
  }

  // 标签页切换
  private handleTabClick() {
    this.page_init = false;
    this.meta.page = 1;
    this.list = [];
    this.hign_choose = [];
    this.guide_workers = [];
    this.tab_name = 1;
    this.reset();
  }

  // 获取列表
  private fetchList() {
    const params = {
      search_type: this.search_type,
      search_value: this.search_name,
      page: this.meta.page,
      size: this.meta.size,
      status: Number(this.active_name),
      start_time: this.date_value ? this.date_value[0] / 1000 : "",
      end_time: this.date_value ? this.date_value[1] / 1000 : "",
      worker_ids: this.guide_workers
    };
    Api.approvalList(params).then(res => {
      this.list = res.data.list;
      const { total, page, size } = res.data;
      this.meta = { total, page, size };
      this.page_init = true;
    });
  }

  // 导出列表
  private exportList() {
    if (!this.case_id.length) {
      this.$message.warning("请先选择要批量导出的数据");
      return;
    }
    const params = {
      search_type: this.search_type,
      search_value: this.search_name,
      status: Number(this.active_name),
      start_time: this.date_value ? this.date_value[0] / 1000 : "",
      end_time: this.date_value ? this.date_value[1] / 1000 : "",
      worder_ids: this.guide_workers,
      ids: this.case_id.map((v: any) => v.id)
    };
    Api.exportApprovalList(params).then(() => {
      this.$confirm(
        "文件已进入下载列表，请前往【下载中心】查看及下载文件。",
        "文件下载提示",
        {
          distinguishCancelAndClose: true,
          confirmButtonText: "前往下载中心",
          cancelButtonText: "取消",
          type: "success"
        }
      ).then(() => {
        ShareModule.getShared().pushState("myapp_download-center");
      });
    });
  }

  // 批准
  private approvalCase(type: string | number, row: any) {
    if (type === 1 && !this.case_id.length) {
      this.$message.warning("请先选择要批量批准的数据");
      return;
    }
    this.$confirm(`你确认要审批通过吗？`, "批准确认", {
      confirmButtonText: "确定",
      cancelButtonText: "取消"
    }).then(() => {
      const params = {
        action: "pass",
        ids: type === 1 ? this.case_id.map((v: any) => v.id) : [row.id]
      };
      Api.operationStatus(params).then(() => {
        this.$message.success("批准成功!");
        this.fetchList();
      });
    });
  }

  // 驳回
  private rejectCase(type: number, row: any) {
    if (type === 1 && !this.case_id.length) {
      this.$message.warning("请先选择要批量驳回的数据");
      return;
    }
    this.case_id = type === 1 ? this.case_id.map((v: any) => v.id) : [row.id];
    this.visible = !this.visible;
  }

  // 删除已选导购
  removeGuide(val: any) {
    const idx = this.hign_choose.findIndex(item => item.id === val);
    idx > -1 && this.hign_choose.splice(idx, 1);
    this.guide_workers = this.hign_choose.map((v: any) => v.id);
    this.fetchList();
  }

  private mounted() {
    this.fetchList();
  }
}
</script>

<style lang="scss" scoped>
/deep/.layout-list {
  .slide-box {
    display: none;
  }
}

/deep/.search-item {
  display: inline-flex;
  align-items: center;
}

.input-with-select {
  width: 116px;
}

.search-input {
  width: 280px !important;
  margin-bottom: 0px !important;
}

.skuname {
  width: 260px;
}

.el-scrollbar {
  overflow-x: hidden !important;
  .el-scrollbar__view {
    .el-select-dropdown__item {
      height: auto;
    }
  }
}

.el-select-dropdown__item.selected {
  font-weight: normal;
  height: auto !important;
}

.el-icon-refresh-right {
  color: #409eff;
  cursor: pointer;
}

/deep/.el-tree--highlight-current
  .el-tree-node.is-current
  > .el-tree-node__content {
  color: #409eff !important;
}

.guide-name {
  display: flex;
  align-items: center;
}

.table-text-img {
  border-radius: 50%;
  margin-right: 10px;
}
</style>
