<template>
  <div class="record-list">
    <layout-list>
      <!-- tab切换 -->
      <template slot="buttons">
        <div class="settlement">
          <div
            @click="settlement(0)"
            :class="{ 'active-tab': flag === 0 }"
            class="info-tilte"
          >
            <span>待结算（元）</span>
            <el-popover
              ref="popoverAN"
              placement="top-start"
              width="270"
              trigger="hover"
              visible-arrow="true"
              content="对应的分销订单还未达到可以结算的状态"
            >
              <svg-icon
                slot="reference"
                class="help-icon"
                icon-class="group-help"
              ></svg-icon>
            </el-popover>
            <div class="sum">{{ wait_settle }}</div>
          </div>
          <div
            @click="settlement(1)"
            :class="{ 'active-tab': flag === 1 }"
            class="info-tilte"
          >
            <span>已结算（元）</span>
            <el-popover
              ref="popoverAM"
              placement="top-start"
              width="200"
              trigger="hover"
              visible-arrow="true"
              content="对应的分销订单达到可结算状态"
            >
              <svg-icon
                slot="reference"
                class="help-icon"
                icon-class="group-help"
              ></svg-icon>
            </el-popover>
            <div class="sum">{{ adopt_settle }}</div>
          </div>
          <div
            @click="settlement(2)"
            :class="{ 'active-tab': flag === 2 }"
            class="info-tilte"
          >
            <span>已取消（元）</span>
            <el-popover
              ref="popoverAP"
              placement="top-start"
              width="240"
              trigger="hover"
              visible-arrow="true"
              content="对应的分销订单取消或被完全退款"
            >
              <svg-icon
                slot="reference"
                class="help-icon"
                icon-class="group-help"
              ></svg-icon>
            </el-popover>
            <div class="sum">{{ cancel_settle }}</div>
          </div>
        </div>
      </template>

      <!-- 搜索部分 -->
      <template slot="search">
        <search-item text="编号：">
          <el-input
            :placeholder="`请输入${search_types[search_type - 1].name}`"
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
        <search-item text="订单状态：" class="order_type">
          <el-select
            v-if="flag === 0"
            class="search-select"
            size="medium"
            v-model="active_type"
            clearable
          >
            <el-option
              v-for="item in active_types0"
              :key="item.key"
              :label="item.name"
              :value="item.key"
            ></el-option>
          </el-select>
          <el-select
            v-if="flag === 1"
            class="search-select"
            size="medium"
            v-model="active_type"
            clearable
          >
            <el-option
              v-for="item in active_types1"
              :key="item.key"
              :label="item.name"
              :value="item.key"
            ></el-option>
          </el-select>
          <el-select
            v-if="flag === 2"
            class="search-select"
            size="medium"
            v-model="active_type"
            clearable
          >
            <el-option
              v-for="item in active_types2"
              :key="item.key"
              :label="item.name"
              :value="item.key"
            ></el-option>
          </el-select>
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
            ref="searchSelectMember"
            :hasPrepend="true"
            :memberTypeList="guide_types"
            :selectedList="hign_choose"
            @member-type-change="onGuideTypeChange"
            @visible-change="highSearch"
            @remove-member="removeGuide"
          ></search-select-member>
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
          size="small"
          @click="exportFile"
          :disabled="list.length === 0"
        >
          导出Excel
        </el-button>
      </template>

      <!-- 表格部分 -->
      <template slot="table">
        <el-table
          v-if="isLoaded"
          :data="list"
          stripe
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
          <el-table-column label="结算状态" min-width="100">
            <template slot-scope="scope">
              <el-tag
                v-if="scope.row.settlement_status === 0"
                type="success"
                size="small"
              >
                待结算
              </el-tag>
              <el-tag
                v-if="scope.row.settlement_status === 1"
                type="info"
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
          <el-table-column label="订单状态" min-width="100">
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
          <el-table-column label="订单支付时间" min-width="180">
            <template slot-scope="scope">
              <span>{{
                scope.row.order_pay_time | formatDate("yyyy-mm-dd hh:mm:ss")
              }}</span>
            </template>
          </el-table-column>
          <el-table-column label="推广导购所在门店/部门" min-width="200">
            <template slot-scope="scope">
              <span>{{ scope.row.promotion_dept || "--" }}</span>
            </template>
          </el-table-column>
          <el-table-column label="推广导购" min-width="120">
            <template slot-scope="scope">
              <div class="nameBox" v-if="scope.row.promotion_name">
                <img
                  v-if="scope.row.promotion_avatar"
                  :src="scope.row.promotion_avatar"
                  alt=""
                />
                <span class="m-l-5">{{ scope.row.promotion_name }}</span>
              </div>
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column label="推广佣金(元)" min-width="120" align="right">
            <template slot-scope="scope">
              <span
                v-if="
                  (scope.row.promotion_commission ||
                    scope.row.promotion_commission == 0) &&
                    scope.row.promotion_name
                "
              >
                {{ scope.row.promotion_commission / 100 }}
              </span>
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column
            v-if="USER.version !== 3"
            label="服务导购"
            min-width="120"
          >
            <template slot-scope="scope">
              <div class="nameBox" v-if="scope.row.service_name">
                <img
                  v-if="scope.row.service_avatar"
                  :src="scope.row.service_avatar"
                  alt=""
                />
                <span class="m-l-5">{{ scope.row.service_name }}</span>
              </div>
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column
            v-if="USER.version !== 3"
            label="服务佣金(元)"
            min-width="120"
            align="right"
          >
            <template slot-scope="scope">
              <span
                v-if="
                  (scope.row.service_commission ||
                    scope.row.service_commission == 0) &&
                    scope.row.service_name
                "
              >
                {{ scope.row.service_commission / 100 }}
              </span>
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column
            v-if="USER.version !== 3"
            prop="service_dept"
            label="服务导购所在门店"
            min-width="160"
          >
            <template slot-scope="scope">
              <span>{{
                scope.row.service_dept ? scope.row.service_dept : "--"
              }}</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" min-width="100" fixed="right">
            <template slot-scope="scope">
              <el-button
                @click.native.prevent="goDetail(scope.row)"
                type="text"
                class="operation-success"
              >
                结算明细
              </el-button>
            </template>
          </el-table-column>
          <div slot="empty">
            <layout-empty
              v-if="flag === 0 && !searchLoading"
              :message="'暂无待结算记录'"
            ></layout-empty>
            <layout-empty
              v-if="flag === 1 && !searchLoading"
              :message="'暂无已结算记录'"
            ></layout-empty>
            <layout-empty
              v-if="flag === 2 && !searchLoading"
              :message="'暂无已取消记录'"
            ></layout-empty>
            <layout-empty
              v-if="searchLoading"
              empty_type="tag"
              :message="'没有找到相关数据'"
            ></layout-empty>
          </div>
        </el-table>

        <!-- 选择导购 -->
        <select-member
          v-if="member_show"
          dialog_title="选择导购"
          choice_name="已选导购"
          :visible="member_show"
          :selectArr="hign_choose"
          :tabName="USER.version == 3 ? 3 : 1"
          :is_show_guide="true"
          @close="member_show = false"
          @submit="submitDialog"
        ></select-member>
      </template>

      <template slot="footer">
        <el-pagination
          v-if="isLoaded && list && list.length > 0"
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
import UserModule from "@/store/modules/user";
import SettlementApi from "@/api/wisdom/Settlement";
import SelectMember from "@/components/sy-components/select-member.vue";
import LayoutEmpty from "@/components/sy-components/layout-empty.vue";
import SearchSelectMember from "@/components/sy-components/search-select-member.vue";
import ShareModule from "@/shared/index";
// import DownLoad from "@/mixins/DownLoad";

@Component({
  components: {
    SelectMember,
    LayoutEmpty,
    SearchSelectMember
  }
})
export default class extends Vue {
  /*** data ***/
  private member_show = false;
  private hign_choose: Array<any> = [];
  private searchLoading = false; // 是否正在搜索
  private isLoaded = false; // 数据是否加载完成
  private search_name = ""; // 输入的具体单号或姓名
  private search_type = 1; // 订单编号或者导购姓名
  private search_types = [
    { name: "订单编号", key: 1, label: "order_no" },
    { name: "导购姓名", key: 2, label: "guide_user" },
    { name: "结算单号", key: 3, label: "settlement_no" }
  ];
  private active_type = ""; // 订单状态
  private active_types0 = [
    { name: "待收货", key: 60 },
    { name: "待发货", key: 40 },
    { name: "已签收", key: 80 },
    { name: "售后中", key: 61 },
    { name: "退款中", key: 200 }
  ];
  private active_types1 = [
    { name: "已完成", key: 100 },
    { name: "退款失败", key: 201 },
    { name: "部分退款", key: 83 },
    { name: "已退款", key: 205 }
  ];
  private active_types2 = [
    { name: "已退款", key: 205 },
    { name: "已取消", key: 21 }
  ];
  private guide_types = [
    { name: "不限", key: 1 },
    { name: "推广导购", key: 2 },
    { name: "服务导购", key: 3 }
  ];
  private guide_type: string | number = 1; // 导购类型
  private guide_workers: Array<any> = [];
  private date_value: any = [];
  private flag = 0; // 切换结算状态
  private wait_settle = 0; // 待结算金额
  private adopt_settle = 0; // 已结算
  private cancel_settle = 0; // 已取消
  private list: Array<any> = []; // 列表数据
  private meta = {
    total: 40, // 总数
    page: 1, // 当前页
    size: 10 // 每页显示数量
  };
  private order_status_list: Array<any> = [];

  /*** 计算属性 ***/
  get USER() {
    return UserModule.getUserInfo;
  }

  /*** methods ***/
  // 获取列表数据
  private getList() {
    const params = {
      version: this.USER.version,
      guide_workers: this.guide_workers,
      search_type: this.search_type,
      search_value: this.search_name,
      order_status: this.active_type, // 订单状态
      start_time: this.date_value ? this.date_value[0] / 1000 : "", //  开始时间
      end_time: this.date_value ? this.date_value[1] / 1000 : "", // 结束时间
      settlement_status: this.flag, // 结算状态
      page: this.meta.page, //  页数
      size: this.meta.size, // 请求数量
      is_page: true, // 接口合并
      guide_workers_scope_type: this.guide_type
    };
    SettlementApi.getSettlementListApi(params).then(res => {
      const data = { ...res.data };
      this.list = data.list;
      const { total, page, size } = data;
      this.meta = { total, page, size };
      this.wait_settle = data.statistics.wait_settle / 100;
      this.adopt_settle = data.statistics.adopt_settle / 100;
      this.cancel_settle = data.statistics.cancel_settle / 100;
      this.isLoaded = true;
    });
  }

  // 每页显示条数
  private handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.getList();
  }

  // 当前页
  private handleCurrentPageChange(val: number) {
    this.meta.page = val;
    this.getList();
  }

  private onGuideTypeChange(val: string | number) {
    this.guide_type = val;
  }

  // 高级搜索
  private highSearch() {
    this.member_show = true;
  }

  // 高级搜索选择完成
  private submitDialog(list: any) {
    this.hign_choose = list;
    this.guide_workers = list.map((v: any) => v.id);
    this.getList();
  }

  // 删除已选导购
  private removeGuide(val: any) {
    const idx = this.hign_choose.findIndex(item => item.id === val);
    idx > -1 && this.hign_choose.splice(idx, 1);
    this.guide_workers = this.hign_choose.map((v: any) => v.id);
    this.getList();
  }

  // 查询结算记录
  private search() {
    this.meta.page = 1;
    this.searchLoading = true;
    this.getList();
  }

  // 重置
  private reset() {
    this.search_type = 1;
    this.search_name = "";
    this.active_type = "";
    this.guide_type = 1;
    this.searchLoading = false;
    (this.$refs.searchSelectMember as any).initMemberType();
    this.guide_workers = [];
    this.hign_choose = [];
    this.date_value = "";
    this.meta.page = 1;
    this.getList();
  }
  private getOrderStatusList() {
    switch (this.flag) {
      case 0:
        return [
          { name: "待收货", skey: 60 },
          { name: "待发货", key: 40 },
          { name: "已签收", key: 80 },
          { name: "售后中", key: 61 },
          { name: "退款中", key: 200 }
        ];
      case 1:
        return [
          { name: "已完成", key: 100 },
          { name: "退款失败", key: 201 },
          { name: "部分退款", key: 83 },
          { name: "已退款", key: 205 }
        ];
      case 2:
        return [
          { name: "已退款", key: 205 },
          { name: "已取消", key: 21 }
        ];
      default:
        return [];
    }
  }

  // 获取当前时间
  private getformatDate() {
    const newDate = new Date();
    const year = newDate.getFullYear();
    let month: string | number = newDate.getMonth() + 1;
    let day: string | number = newDate.getDate();
    month = month < 10 ? "0" + month : month;
    day = day < 10 ? "0" + day : day;
    return year.toString() + month.toString() + day.toString();
  }

  // 导出excel
  private exportFile() {
    const params = {
      version: this.USER.version,
      search_type: this.search_type,
      search_value: this.search_name,
      guide_workers: this.guide_workers,
      order_status: this.active_type, // 订单状态
      start_time: this.date_value[0] ? this.date_value[0] / 1000 : "", //  开始时间
      end_time: this.date_value[1] ? this.date_value[1] / 1000 : "", // 结束时间
      settlement_status: this.flag, // 结算状态
      guide_workers_scope_type: this.guide_type
    };
    SettlementApi.exportSettlementListApi(params).then(() => {
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

  // 切换结算状态
  private settlement(type: any) {
    this.flag = type;
    console.log("settlement: ", type);
    this.reset();
  }

  // 去往结算明细
  private goDetail(params: any) {
    this.$router.push({
      name: "myapp_wisdom_staff-distribution_settlement-record_detail",
      query: {
        id: params.id
      }
    });
  }

  private mounted() {
    this.getList();
  }
}
</script>

<style lang="scss" scoped>
.record-list {
  color: #303133;
}
/deep/.layout-list {
  .tabs-content,
  .slide-box {
    display: none;
  }
}

/deep/.search-item {
  display: inline-flex;
  align-items: center;
}

.el-select-dropdown__item.selected {
  font-weight: normal;
  height: auto !important;
}

.input-with-select {
  width: 110px;
}

.search-input {
  width: 280px !important;
  margin-bottom: 0px !important;
}

.settlement {
  overflow: hidden;
}

.info-tilte {
  float: left;
  font-size: 14px;
  width: 203px;
  height: 73px;
  margin-right: 20px;
  padding: 20px 0 0 15px;
  background: $color-background;
  border-radius: 5px;
  cursor: pointer;

  span {
    display: inline-block;
    margin-bottom: 15px;
  }
}

.active-tab {
  background-color: $color-main-01;
  border: 1px solid $color-main;
}

.sum {
  font-weight: 700;
  font-size: 20px;
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
    width: 85px;
    flex-shrink: 0;
  }
}

.search-select {
  width: 280px;
}
</style>
