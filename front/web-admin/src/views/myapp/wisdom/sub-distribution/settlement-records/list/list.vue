<template>
  <div class="manage-list">
    <layout-list>
      <!-- 搜索 -->
      <template slot="search">
        <search-item text="编号：">
          <el-select
            v-model="search_type"
            class="search-select"
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
            class="name-input search-input"
            size="small"
          >
          </el-input>
        </search-item>
        <search-item text="状态：">
          <el-select
            class="search-select"
            placeholder="结算状态"
            size="medium"
            v-model="state_type"
            @change="changeType"
          >
            <el-option
              v-for="item in state_types"
              :key="item.key"
              :label="item.name"
              :value="item.key"
            ></el-option>
          </el-select>
          <!-- 结算状态值 -->
          <el-select
            v-if="this.state_type == 1"
            placeholder="请选择结算状态"
            clearable
            size="medium"
            :value="1"
            v-model="settle_type"
          >
            <el-option
              v-for="item in settle_types"
              :key="item.key"
              :label="item.name"
              :value="item.key"
            ></el-option>
          </el-select>
          <el-select
            v-if="this.state_type == 2"
            placeholder="请选择订单状态"
            clearable
            :value="2"
            size="medium"
            v-model="active_type"
          >
            <el-option
              v-for="item in active_types"
              :key="item.key"
              :label="item.name"
              :value="item.key"
            ></el-option>
          </el-select>
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
          <div class="data-button">
            <span>数据更新时间：</span>
            <span>{{ updateTime | formatDate("yyyy-mm-dd hh:mm:ss") }}</span>
            <el-button
              type="text"
              @click="reload"
              size="mini"
              class="mar-l-10 el-icon-refresh"
              >更新</el-button
            >
          </div>
        </div>

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

      <!-- 操作 -->
      <template slot="operate">
        <el-button
          size="small"
          @click="exportFile"
          :disabled="list.length === 0"
        >
          导出Excel
        </el-button>
      </template>

      <!-- 表格 -->
      <template slot="table">
        <el-table
          v-if="isLoaded"
          :data="list"
          ref="table"
          stripe
          @sort-change="sortChange"
          :header-cell-style="{
            fontSize: '14px',
            color: '#606266',
            fontWeight: '400',
            padding: '11px 0 10px 0'
          }"
        >
          <el-table-column
            label="结算单号"
            prop="id"
            min-width="170"
            fixed="left"
          ></el-table-column>
          <el-table-column
            label="结算状态"
            prop="settlement_status"
            min-width="120"
          >
            <template slot-scope="scope">
              <el-tag
                v-if="scope.row.settlement_status === 0"
                type="warning"
                size="small"
              >
                待结算
              </el-tag>
              <el-tag
                v-if="scope.row.settlement_status === 1"
                type="success"
                size="small"
              >
                已结算
              </el-tag>
              <el-tag
                v-if="scope.row.settlement_status === 2"
                type="info"
                size="small"
              >
                已取消
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column
            prop="order_no"
            label="订单编号"
            min-width="170"
            size="small"
          >
          </el-table-column>
          <el-table-column label="订单状态" prop="order_status" min-width="120">
            <template slot-scope="scope">
              <el-tag v-if="scope.row.order_status === 20" size="small">
                未支付
              </el-tag>
              <el-tag
                v-if="scope.row.order_status === 21"
                type="info"
                size="small"
              >
                已取消
              </el-tag>
              <el-tag v-if="scope.row.order_status === 23" size="small">
                待支付
              </el-tag>
              <el-tag v-if="scope.row.order_status === 40" size="small">
                待发货
              </el-tag>
              <el-tag
                v-if="scope.row.order_status === 60"
                type="info"
                size="small"
              >
                待收货
              </el-tag>
              <el-tag
                v-if="scope.row.order_status === 80"
                type="info"
                size="small"
              >
                已签收
              </el-tag>
              <el-tag
                v-if="[41, 61, 81].includes(scope.row.order_status)"
                type="info"
                size="small"
                >售后中</el-tag
              >
              <el-tag
                v-if="scope.row.order_status === 200"
                type="info"
                size="small"
              >
                退款中
              </el-tag>
              <el-tag
                v-if="scope.row.order_status === 205"
                type="info"
                size="small"
              >
                已退款
              </el-tag>
              <el-tag
                v-if="scope.row.order_status === 100"
                type="info"
                size="small"
              >
                已完成
              </el-tag>
              <el-tag
                v-if="scope.row.order_status === 201"
                type="info"
                size="small"
              >
                退款失败
              </el-tag>
              <el-tag
                v-if="scope.row.order_status === 83"
                type="info"
                size="small"
              >
                部分退款
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="订单支付时间" min-width="210">
            <template slot-scope="scope">
              {{ scope.row.order_pay_time | formatDate("yyyy-mm-dd hh:mm:ss") }}
            </template>
          </el-table-column>
          <el-table-column label="分销合伙人" width="230">
            <template slot-scope="scope">
              <div class="nameBox" v-if="scope.row.partner_name">
                <img
                  :src="
                    scope.row.avatar
                      ? scope.row.avatar
                      : 'https://oss-crs.vchangyi.com/wechat/default-avatar-icon.png'
                  "
                  alt=""
                />
                <span
                  class="m-l-5"
                  :class="
                    scope.row.partner_name.length > 10 ? 'partner-name' : ''
                  "
                  >{{ scope.row.partner_name }}</span
                >
              </div>
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column
            label="分销合伙人佣金（元）"
            min-width="200"
            prop="partner_commission"
            sortable="custom"
            align="right"
          >
            <div slot-scope="{ row }">
              {{ row.partner_commission / 100 }}
            </div>
          </el-table-column>
          <el-table-column label="所属导购" min-width="180">
            <template slot-scope="scope">
              <div class="nameBox" v-if="scope.row.worker_name">
                <img
                  :src="
                    scope.row.worker_avatar
                      ? scope.row.worker_avatar
                      : 'https://oss-crs.vchangyi.com/wechat/default-avatar-icon.png'
                  "
                  alt=""
                />
                <span class="m-l-5">{{ scope.row.worker_name }}</span>
              </div>
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column
            label="导购佣金（元）"
            min-width="210"
            prop="promotion_commission"
            sortable="custom"
            align="right"
          >
            <div slot-scope="{ row }">
              {{ row.promotion_commission / 100 }}
            </div>
          </el-table-column>
          <el-table-column label="操作" min-width="100" fixed="right">
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
              v-if="!condition"
              :message="'暂无结算记录'"
            ></layout-empty>
            <layout-empty
              v-if="condition"
              empty_type="tag"
              :message="'没有找到相关数据'"
            ></layout-empty>
          </div>
        </el-table>
      </template>

      <!-- 分页 -->
      <template slot="footer">
        <!-- v-if="isLoaded && list && list.length > 0" -->
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :page-sizes="[10, 20, 30, 40]"
          :current-page="meta.page"
          :page-size="meta.size"
          :total="meta.total"
          @current-change="handleCurrentPageChange"
          @size-change="handleSizeChange"
        ></el-pagination>
      </template>
    </layout-list>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import SelectMember from "@/components/sy-components/select-member.vue";
import LayoutEmpty from "@/components/sy-components/layout-empty.vue";
import SearchSelectMember from "@/components/sy-components/search-select-member.vue";
import settlementRecordsApi from "@/api/wisdom/settlementRecords";
import UserModule from "@/store/modules/user";
import ShareModule from "@/shared/index";

@Component({
  components: {
    SelectMember,
    LayoutEmpty,
    SearchSelectMember
  }
})
export default class extends Vue {
  private search_type = 1; // 搜索类型 1订单编号 2结算单号
  private search_types = [
    { name: "订单编号", key: 1 },
    { name: "结算单号", key: 2 }
  ];
  private search_name = ""; // 搜索编号活动名称
  private search_status = ""; // 搜索状态活动名称
  private meta = {
    total: 40, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  private updateTime = 0;
  private timerId = 0;
  private statusFlag = true;
  private state_type = 1; //状态
  private state_types = [
    { name: "结算状态", key: 1 },
    { name: "订单状态", key: 2 }
  ];
  // 当state_type为1结算状态时的对应下拉菜单
  private settle_type = "";
  private settle_types = [
    { name: "待结算", key: 0 },
    { name: "已结算", key: 1 },
    { name: "已取消", key: 2 }
  ];
  // 当state_type为2订单状态时的对应下拉菜单
  private active_type = ""; // 订单状态
  private active_types = [
    { name: "待收货", key: 60 },
    { name: "待发货", key: 40 },
    { name: "已签收", key: 80 },
    { name: "售后中", key: 41 },
    { name: "退款中", key: 200 },
    { name: "已完成", key: 100 },
    { name: "退款失败", key: 201 },
    { name: "部分退款", key: 83 },
    { name: "已退款", key: 205 },
    { name: "已取消", key: 21 }
  ];
  private hign_choose: Array<any> = [];
  private isLoaded = false; // 数据是否加载完成
  private list: Array<any> = []; // 列表数据
  private guide_workers: Array<any> = [];
  private order_desc = ""; //排序方式 asc 正序desc 倒序
  private order_type = 0; //排序类型1合伙人佣金排序2导购佣金排序
  private condition = false; //是否查询

  // 接口数据
  private member_show = false;

  // methods
  // 下拉框改变的时候
  private selectChange() {
    this.search_name = "";
    this.search_status = "";
  }

  //刷新
  private reload() {
    this.list = [];
    this.reset();
    this.getList();
    this.updata();
  }
  private updata() {
    const date = new Date().getTime();
    //返回格式化后的日期
    this.updateTime = date / 1000;
    console.log(this.updateTime);
  }

  // 获取列表数据
  private getList() {
    const params = {
      search_type: this.search_type,
      search_value: this.search_name,
      status_type: this.state_type, //是否为订单状态或结算状态
      status_value: this.state_type == 1 ? this.settle_type : this.active_type, // 状态值
      worker_ids: this.guide_workers,
      page: this.meta.page, //  页数
      size: this.meta.size, // 请求数量
      is_page: true, // 接口合并
      order_by_type: this.order_type, //排序类型1合伙人佣金排序2导购佣金排序
      orderby: this.order_desc //asc 正序desc 倒序
    };
    settlementRecordsApi.getSettlementListApi(params).then(res => {
      const data = { ...res.data };
      this.list = data.list;
      const { total, page, size } = data;
      this.meta = { total, page, size };
      const time = Date.now() / 1000;
      this.updateTime = time;
      this.isLoaded = true;
    });
  }

  // 状态改变
  private changeType() {
    // this.statusFlag = !this.statusFlag;
    this.state_type = 2;
    this.active_type = "";
    this.settle_type = "";
  }

  // 每页显示条数
  private handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.getList();
  }

  // 排序
  private sortChange(column: any) {
    //排序类型1合伙人佣金排序2导购佣金排序
    // 排序方式 asc 正序desc 倒序
    if (column.order == null) {
      this.order_type = 0;
      this.order_desc = "";
    } else {
      this.order_type = column.prop === "partner_commission" ? 1 : 2;
      // ascending 升序 descending 降序
      this.order_desc = column.order === "ascending" ? "asc" : "desc";
    }
    this.getList();
  }

  // 当前页
  private handleCurrentPageChange(val: number) {
    this.meta.page = val;
    this.getList();
  }

  // 去往结算明细
  private goDetail(params: any) {
    this.$router.push({
      name: "myapp_wisdom_sub-distribution_settlement-records_detail",
      query: {
        id: params.id
      }
    });
  }

  // 计算属性
  get USER() {
    return UserModule.getUserInfo;
  }

  // 员工输入框即高级搜索
  private highSearch() {
    this.member_show = true;
  }
  // 高级搜索完成
  private submitDialog(list: any) {
    this.hign_choose = list;
    this.guide_workers = list.map((v: any) => v.id);
    // 调用数据
    this.getList();
  }

  // 删除已选导购
  private removeGuide(val: any) {
    const idx = this.hign_choose.findIndex(item => item.id === val);
    idx > -1 && this.hign_choose.splice(idx, 1);
    this.guide_workers = this.hign_choose.map((v: any) => v.id);
    // 调用数据
    this.getList();
  }

  // 查询结算记录
  private search() {
    this.condition = true;
    this.meta.page = 1;
    (this.$refs["table"] as any).clearSort();
    // 调用查询数据
    this.getList();
  }
  // 重置
  private reset() {
    this.search_type = 1;
    this.search_name = "";
    this.state_type = 1;
    this.active_type = "";
    this.settle_type = "";
    (this.$refs["table"] as any).clearSort();
    this.hign_choose = [];
    this.guide_workers = [];
    this.meta.page = 1;
    this.order_desc = ""; //排序方式
    this.order_type = 0;
    // 调用查询数据
    this.getList();
  }
  // 导出excel
  private exportFile() {
    const params = {
      search_type: this.search_type,
      search_value: this.search_name,
      status_type: this.state_type, //是否为订单状态或结算状态
      status_value: this.state_type == 1 ? this.settle_type : this.active_type, // 状态值
      worker_ids: this.guide_workers,
      orderby: this.order_desc,
      page: this.meta.page, //  页数
      size: this.meta.size // 请求数量
    };
    settlementRecordsApi.exportSettlementListApi(params).then(() => {
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
  // 生命周期
  private mounted() {
    this.getList();
  }
}
</script>
<style lang="scss" scoped>
/deep/.search-item {
  display: inline-flex;
  align-items: center;
}

.search-item /deep/ .el-input--small .el-input__inner {
  height: 36px;
  line-height: 36px;
}

.search-item .search-select {
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
  text-align: center;
}

.search-input /deep/ .el-input__inner {
  text-align: left;
}

.data-button {
  display: inline-block;
  margin-left: 20px;
  color: #666666;
  font-size: 12px;
}

.search-select /deep/ .el-input__inner {
  border-right: none;
  border-top-right-radius: 0;
  border-bottom-right-radius: 0;
}

/deep/.el-input--medium .el-input__inner {
  border-top-left-radius: 0;
  border-bottom-left-radius: 0;
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

.com-section {
  border: 1px solid #dcdfe6;
  border-top-right-radius: 4px;
  border-bottom-right-radius: 4px;
  background-color: #fff;
}

.com-section /deep/ .el-input {
  width: 77px;
  display: inline-flex;
}

.com-section /deep/ .el-input .el-input__inner {
  border: none;
  height: 34px;
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

.nameBox {
  display: flex;
  align-items: center;

  img {
    width: 30px;
    height: 30px;
    border-radius: 50%;
  }

  span {
    width: 140px;
    flex-shrink: 0;
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
.partner-name {
  display: block;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}
</style>
