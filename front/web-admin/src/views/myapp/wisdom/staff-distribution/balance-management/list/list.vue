<template>
  <layout-list title_name="余额管理" class="banlance-list">
    <template slot="search">
      <search-item text="筛选：">
        <el-select
          v-model="search_type"
          class="name-select"
          size="small"
          @change="selectChange"
        >
          <el-option
            v-for="item in search_types"
            :key="item.key"
            :label="item.name"
            :value="item.key"
          >
          </el-option>
        </el-select>
        <el-input
          :placeholder="`请输入${search_types[search_type - 1].name}`"
          v-model="search_name"
          class="name-input"
          size="small"
        >
        </el-input>
      </search-item>
      <search-item text="导购：">
        <search-select-member
          :selectedList="hign_choose"
          @visible-change="highSearch"
          @remove-member="removeGuide"
        >
        </search-select-member>
      </search-item>
      <div class="operate-btn">
        <el-button type="primary" size="small" @click="search">查询</el-button>
        <el-button type="primary" size="small" class="reset" @click="reset">
          重置
        </el-button>
      </div>
    </template>
    <template slot="tabs">
      <div class="clearbox">
        <div class="moneyStatus">
          <span class="spa">待结算金额（元）</span>
          <span>{{ no_settlement_price }}</span>
        </div>
        <div class="moneyStatus">
          <span class="spa">可提现金额（元）</span>
          <span>{{ available_cash_amount }}</span>
        </div>
        <div class="moneyStatus">
          <span class="spa">已提现金额（元）</span>
          <span>{{ cash_withdrawal }}</span>
        </div>
      </div>
      <el-button
        size="small"
        @click="exportList"
        :class="{ exportBut: exportBut }"
        class="export-btn"
      >
        导出Excel
      </el-button>
    </template>
    <!-- 表格部分 -->
    <template slot="table">
      <el-table
        :data="list"
        stripe
        class="mar-t-15 list-last-line-none min-height-100"
        v-if="page_init"
        :header-cell-style="{
          fontSize: '14px',
          color: '#606266',
          fontWeight: '400',
          padding: '11px 0 10px 0'
        }"
      >
        <el-table-column label="导购" min-width="180">
          <template slot-scope="scope">
            <img
              style="width:30px;height:30px;"
              :src="
                scope.row.avatar
                  ? scope.row.avatar
                  : 'https://oss-crs.vchangyi.com/wechat/default-avatar-icon.png'
              "
              alt=""
              class="table-text-img"
            />
            <span v-text="scope.row.name"></span>
          </template>
        </el-table-column>
        <el-table-column label="手机号" width="150" show-overflow-tooltip>
          <template slot-scope="scope">
            <div>
              {{ scope.row.mobile ? scope.row.mobile : "--" }}
            </div>
          </template>
        </el-table-column>
        <el-table-column
          label="导购所在门店/部门"
          width="200"
          show-overflow-tooltip
        >
          <template slot-scope="scope">
            <span>
              <el-tooltip
                v-if="tooltip"
                effect="dark"
                :content="scope.row.store_name"
                placement="top"
              ></el-tooltip>
              <span>{{ scope.row.store_name_format || "--" }}</span>
            </span>
            <span v-for="(item, index) in scope.row.dpName" :key="index"
              >{{ item }}
              <span
                v-if="index > 1 && item !== scope.row.dpName.slice(-1)"
              ></span>
            </span>
          </template>
        </el-table-column>
        <el-table-column
          prop="expect_cash_amount"
          label="待结算金额(元)"
          width="160"
          align="right"
        >
          <template slot-scope="scope">
            <span>{{ scope.row.expect_cash_amount / 100 }}</span>
          </template>
        </el-table-column>
        <el-table-column
          prop="available_cash_amount"
          label="可提现金额(元)"
          width="160"
          align="right"
        >
          <template slot-scope="scope">
            <span>{{ scope.row.available_cash_amount / 100 }}</span>
          </template>
        </el-table-column>
        <el-table-column
          prop="cash_withdrawal"
          label="已提现金额(元)"
          width="160"
          align="right"
        >
          <template slot-scope="scope">
            <span>{{ scope.row.cash_withdrawal / 100 }}</span>
          </template>
        </el-table-column>
        <el-table-column
          prop="cash_withdrawal"
          label="操作"
          width="100"
          fixed="right"
        >
          <template slot-scope="scope">
            <el-button
              @click.native.prevent="goDetail(scope.row)"
              type="text"
              class="operation-success"
            >
              查看明细
            </el-button>
          </template>
        </el-table-column>
        <div slot="empty">
          <layout-empty
            v-if="!search_name"
            :message="'暂无余额记录'"
          ></layout-empty>
          <layout-empty
            v-if="search_name"
            empty_type="tag"
            :message="'没有找到相关导购'"
          ></layout-empty>
        </div>
      </el-table>

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
        v-if="page_init"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="meta.page"
        :page-sizes="[10, 20, 30, 40]"
        :page-size="meta.size"
        class="page-right font-weight-400"
        layout="total, sizes, prev, pager, next, jumper"
        :total="meta.total"
      >
      </el-pagination>
    </template>
    <!-- 高级搜索组件 -->
  </layout-list>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Api from "@/api/wisdom/balanceManage";
import LayoutEmpty from "@/components/sy-components/layout-empty.vue";
import LayoutForm from "@/components/public/layout-form/index.vue";
import LayoutList from "@/components/public/layout-list/index.vue";
import SelectMember from "@/components/sy-components/select-member.vue";
import SearchSelectMember from "@/components/sy-components/search-select-member.vue";
import UserModule from "@/store/modules/user";
import ShareModule from "@/shared/index";

@Component({
  components: {
    LayoutEmpty,
    LayoutForm,
    LayoutList,
    SelectMember,
    SearchSelectMember
  }
})
export default class extends Vue {
  private dialog_title = "余额管理";
  private times: any = "";
  private omit = "...";
  private page_init = false; // 页面初始状态
  private search_type = 1; // 搜索类型 1姓名 2手机号
  private search_types = [
    { name: "导购姓名", key: 1 },
    { name: "导购手机号", key: 2 }
  ];
  private search_name = ""; // 搜索活动名称
  private meta = {
    total: 40, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };

  // 接口数据
  private list = [];
  private no_settlement_price = 0;
  private available_cash_amount = 0;
  private cash_withdrawal = 0;
  private member_show = false;
  private tooltip = false;
  private choose_list = [];
  private tab_name = 1;
  private dataArr = [];
  private dataArg = 1;
  // 导出按钮是否置灰
  private exportBut = false;
  private hign_choose: Array<any> = [];
  private guide_workers: Array<any> = [];

  // 计算属性
  get USER() {
    return UserModule.getUserInfo;
  }

  // methods
  // 下拉框改变的时候
  private selectChange() {
    this.search_name = "";
  }
  // private changedTypeList(val: any) {
  //   this.search_type = val;
  //   this.search_name = "";
  // }
  // 员工输入框即高级搜索
  private highSearch() {
    this.member_show = true;
  }
  // 高级搜索完成
  private submitDialog(list: any) {
    this.hign_choose = list;
    this.guide_workers = list.map((v: any) => v.id);
    this.fetchList();
  }

  // 删除已选导购
  private removeGuide(val: any) {
    const idx = this.hign_choose.findIndex(item => item.id === val);
    idx > -1 && this.hign_choose.splice(idx, 1);
    this.guide_workers = this.hign_choose.map((v: any) => v.id);
    this.fetchList();
  }

  // 获取列表
  private fetchList() {
    const params = {
      search_text: this.search_name,
      search_type: this.search_type,
      worker_ids: this.guide_workers,
      page: this.meta.page,
      size: this.meta.size
    };
    Api.balanceListListApi(params).then(res => {
      this.no_settlement_price = res.data.expect_cash_amount / 100;
      this.available_cash_amount = res.data.available_cash_amount / 100;
      this.cash_withdrawal = res.data.cash_withdrawal / 100;
      res.data.list.forEach((item: any) => {
        if (item.store_name.length > 20) {
          this.tooltip = true;
          item.store_name_format = item.store_name.slice(0, 20) + "...";
        } else {
          item.store_name_format = item.store_name;
        }
      });
      this.list = res.data.list;
      const { total, page, size } = res.data;
      this.meta = { total, page, size };
      this.page_init = true;
      if (res.data.list.length == 0) {
        this.exportBut = true;
      } else {
        this.exportBut = false;
      }
    });
  }

  //  搜索
  private search() {
    this.meta.page = 1;
    this.fetchList();
  }

  // 导出列表
  private exportList() {
    if (this.exportBut == true) {
      return;
    }
    const params = {
      version: this.USER.version,
      search_text: this.search_name,
      search_type: this.search_type,
      worker_ids: this.guide_workers
    };
    Api.balanceListLisExporttApi(params).then(() => {
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

  // 重置
  private reset() {
    this.choose_list = [];
    this.tab_name = 1;
    this.search_type = 1;
    this.search_name = "";
    this.hign_choose = [];
    this.guide_workers = [];
    this.meta.page = 1;
    this.fetchList();
  }

  // 每页
  private handleSizeChange(val: number) {
    this.meta.size = val;
    this.fetchList();
  }
  // 当前页
  private handleCurrentChange(val: number) {
    this.meta.page = val;
    this.fetchList();
  }

  // 去往结算明细
  private goDetail(params: any) {
    this.$router.push({
      name: "myapp_wisdom_staff-distribution_balance-management_detail",
      query: {
        id: params.worker_id
      }
    });
  }

  // 生命周期
  private mounted() {
    this.fetchList();
  }
}
</script>
<style lang="scss" scoped>
.banlance-list {
  color: #303133;
}
/deep/.search-item {
  display: inline-flex;
  align-items: center;
  margin-bottom: 0;
}

.search-item > span {
  width: 70px;
}

.search-item /deep/ .el-input--small .el-input__inner {
  height: 36px;
  line-height: 36px;
}

.search-item .name-select {
  width: 115px;
}

.search-item .name-input {
  width: 165px;
}

.operate-btn {
  margin: 20px 0 0;

  .reset {
    color: #303133;
    background-color: #fff;
    border-color: #c7cad0;
    margin-left: 15px;
  }
}

.name-input /deep/ .el-input__inner {
  border-top-left-radius: 0;
  border-bottom-left-radius: 0;
}
.name-select /deep/ .el-input__inner {
  border-right: none;
  border-top-right-radius: 0;
  border-bottom-right-radius: 0;
}

.banlance-list /deep/ .buttons-content {
  display: none;
}

.banlance-list /deep/ .operate-content {
  display: none;
}

.banlance-list /deep/ .slide-box {
  display: none;
}

.banlance-list /deep/ .search-content {
  padding: 20px;
}

.export-btn {
  margin-bottom: 20px;
}

.time {
  display: inline-block;
  margin-left: 20px;
}

/deep/.el-input__inner {
  border-radius: 4px;
}

.select-input /deep/ .el-input__inner:first-child {
  width: 120px;
}

.table-text-img {
  border-radius: 100%;
  vertical-align: middle;
  margin-right: 4px;
}

.clearbox {
  display: flex;

  .moneyStatus {
    background-color: $color-background;
    font-size: 14px;
    width: 203px;
    height: 55px;
    margin: 10px 20px 20px 0;
    padding: 10px;
    border-radius: 5px;

    span {
      display: block;
      font-size: 20px;
      font-weight: 700;
    }

    .spa {
      margin-bottom: 10px;
      font-weight: 400;
      font-size: 14px;
    }
  }
}

.el-scrollbar {
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
  color: #708df5;
  cursor: pointer;
}

/deep/.el-tree--highlight-current
  .el-tree-node.is-current
  > .el-tree-node__content {
  color: #708df5 !important;
}
</style>
