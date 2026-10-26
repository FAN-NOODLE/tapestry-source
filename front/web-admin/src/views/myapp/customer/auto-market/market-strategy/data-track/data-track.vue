<template>
  <div class="pageDataTrack">
    <div class="container">
      <!-- <el-button plain @click="openSettingFn">业绩数据设置</el-button> -->
      <div class="search_container">
        <p class="time-show">
          数据更新时间：{{ updateTime | formatDate("yyyy-mm-dd hh:mm:ss") }}
        </p>
        <span>触达时间：</span>
        <el-date-picker
          class="el-date-sty"
          v-model="runTimeRange"
          size="middle"
          type="datetimerange"
          range-separator="至"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
          :default-time="['00:00:00', '23:59:59']"
        >
        </el-date-picker>
        <el-button
          type="primary"
          size="small"
          class="m-l-20"
          @click="searchFn"
        >查询</el-button>
        <el-button
          size="small"
          @click="resetFn"
        >重置</el-button>
      </div>
      <div class="container_btns">
        <el-button @click="downloadExeclFn">下载Execl</el-button>
      </div>
      <div v-if="isCoupon">
        <el-table
          :data="listcoupon"
          show-summary
          :span-method="objectSpanMethod"
          :summary-method="getSummaries"
          stripe
        >
          <el-table-column
            prop="type"
            label="渠道类型"
            min-width="150"
            :formatter="formatType"
            :show-overflow-tooltip="true"
            fixed="left"
          >
          </el-table-column>
          <el-table-column
            prop="send_num"
            :render-header="renderHeader"
            label="发送人数"
            min-width="120"
          >
          </el-table-column>
          <el-table-column
            prop="coupon_name"
            label="微信优惠券"
            :show-overflow-tooltip="true"
            min-width="120"
          >
          </el-table-column>
          <el-table-column
            prop="stock_id"
            label="批次号"
            min-width="120"
          >
            <template slot-scope="{ row }">
              {{ row.stock_id }}
            </template>
          </el-table-column>
          <el-table-column
            prop="touch_num"
            label="触达人数"
            sortable
            min-width="140"
            :show-overflow-tooltip="true"
            :render-header="renderHeader"
          >
          </el-table-column>
          <el-table-column
            prop="is_send_num"
            label="发放张数"
            sortable
            min-width="140"
            :render-header="renderHeader"
          >
          </el-table-column>
          <el-table-column
            prop="use_num"
            label="核销张数"
            sortable
            min-width="140"
            :render-header="renderHeader"
          >
          </el-table-column>
          <el-table-column
            prop="use_send"
            :render-header="renderHeader"
            label="核销率"
            sortable
            min-width="140"
          >
            <template slot-scope="{ row }">
              {{ row.use_send + "%" }}
            </template>
          </el-table-column>
          <el-table-column
            prop="pay_num"
            :render-header="renderHeader"
            label="支付买家数"
            sortable
            min-width="150"
          >
          </el-table-column>
          <el-table-column
            prop="pay_touch"
            :render-header="renderHeader"
            label="触达支付转化率"
            sortable
            min-width="180"
          >
            <template slot-scope="{ row }">
              {{ row.pay_touch + "%" }}
            </template>
          </el-table-column>
          <el-table-column
            prop="cash_fee"
            :render-header="renderHeader"
            label="实付金额"
            sortable
            min-width="150"
          >
          </el-table-column>
          <el-table-column
            prop="refund_customer_num"
            :render-header="renderHeader"
            label="退款买家数"
            sortable
            min-width="160"
          >
          </el-table-column>
          <el-table-column
            prop="refund_fee"
            :render-header="renderHeader"
            label="退款金额"
            sortable
            min-width="150"
          >
          </el-table-column>
          <div
            slot="empty"
            class="table-empty-label"
          >
            <img
              :src="searchLoading
                  ? 'https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/empty-search.png'
                  : require('@/assets/images/empty-box.png')
                "
              alt=""
            />
            <div>{{ searchLoading ? "没有找到业绩数据" : "暂无业绩数据" }}</div>
          </div>
        </el-table>
      </div>
      <div v-show="!isCoupon">
        <el-table
          :data="list"
          show-summary
          :summary-method="getSummaries"
          stripe
        >
          <el-table-column
            prop="type"
            label="渠道类型"
            min-width="150"
            :formatter="formatType"
            :show-overflow-tooltip="true"
            fixed="left"
          >
          </el-table-column>
          <el-table-column
            prop="send_member_number"
            label="发送人数"
            min-width="120"
          >
          </el-table-column>
          <el-table-column
            prop="reach_member_number"
            label="触达人数"
            :show-overflow-tooltip="true"
            min-width="120"
          >
          </el-table-column>
          <!-- <el-table-column prop="fans_number" label="加粉人数" min-width="120">
            <template slot-scope="scope">
              {{ scope.row.type === 2 ? "--" : scope.row.fans_number }}
            </template>
          </el-table-column>
          <el-table-column
            prop="visitor_number"
            label="访客数"
            min-width="120"
            :show-overflow-tooltip="true"
          >
          </el-table-column>
          <el-table-column
            prop="add_order_number"
            label="下单买家数"
            min-width="120"
          >
          </el-table-column>
          <el-table-column prop="pay_order" label="支付买家数" min-width="120">
          </el-table-column>
          <el-table-column prop="final_amount" label="实付金额" min-width="120">
          </el-table-column>
          <el-table-column
            prop="refund_amount_number"
            label="退款买家数"
            min-width="120"
          >
          </el-table-column>
          <el-table-column
            prop="refund_amount"
            label="退款金额"
            min-width="120"
          >
          </el-table-column> -->
          <div
            slot="empty"
            class="table-empty-label"
          >
            <img
              :src="searchLoading
                  ? 'https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/empty-search.png'
                  : require('@/assets/images/empty-box.png')
                "
              alt=""
            />
            <div>{{ searchLoading ? "没有找到业绩数据" : "暂无业绩数据" }}</div>
          </div>
        </el-table>
      </div>
    </div>
    <layout-form>
      <template slot="footer">
        <el-button
          @click="goBack()"
          plain
        >返回</el-button>
      </template>
    </layout-form>
    <el-dialog
      title="业绩数据设置"
      :visible.sync="dialogVisible"
      width="242px"
      :close-on-click-modal="false"
    >
      <div class="settingDay">
        <span class="title">业绩有效期：</span>
        <el-select
          size="small"
          placeholder="30"
          v-model="strategyDay"
        >
          <el-option
            v-for="(item, index) in strategyDays"
            :key="index"
            :label="item.label"
            :value="item.value"
          ></el-option>
        </el-select>
        <span class="day"> 天</span>
      </div>
      <span slot="footer">
        <div class="dialog-footer">
          <el-button @click="dialogVisible = false">取 消</el-button>
          <el-button
            type="primary"
            @click="saveUserBeFn"
          >确 定</el-button>
        </div>
      </span>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LayoutList from "@/components/public/layout-list/index.vue";
import API from "@/api/customer";
import { dateFormat, startAndEndOfday } from "../../util/index";
@Component({
  components: {
    LayoutList
  }
})
export default class extends Vue {
  private strategyDay = 30;
  private searchLoading = false;
  private dialogVisible = false;
  private updateTime = 0;
  private strategyName = "";
  private runTimeRange: any = ["", ""];
  private config = {
    send_num: "符合自动发券策略的可发送人数",
    touch_num: "策略中的该优惠券发放成功的人数",
    is_send_num: "自动发放该微信优惠券的张数",
    use_num: "自动发放该微信优惠券的核销张数",
    use_send: "核销张数/发放张数*100%",
    pay_num: "所有已支付完成的客户数",
    pay_touch: "支付买家数/触达人数*100%",
    cash_fee: "所有已支付订单的实付金额",
    refund_customer_num: "所有退款完成的客户数，不包含在途退款订单；",
    refund_fee: "所有退款完成的订单金额，不包含在途退款订单；"
  };
  private strategyDays = [
    {
      label: 1,
      value: 1
    },
    {
      label: 7,
      value: 7
    },
    {
      label: 15,
      value: 15
    },
    {
      label: 30,
      value: 30
    }
  ];
  private list = [];
  private listcoupon = [];
  private total: any = undefined;
  private isCoupon = false;
  formatType(row: any) {
    let label: any = "";
    switch (row.type) {
      case 1:
        label = "短信";
        break;
      case 2:
        label = "企微群发";
        break;
      case 3:
        label = "自动发券";
        break;
      case 4:
        label = "公众号";
        break;
      default:
        break;
    }
    return label;
  }
  searchFn() {
    this.searchLoading = true;
    this.getTrackListFn();
  }
  resetFn() {
    this.runTimeRange = ["", ""];
    this.searchLoading = false;
    this.getTrackListFn();
  }
  async initSettingDay() {
    const data: object = {
      id: Number(this.$route.query.id)
    };
    const res: any = await API.marketStrategyDetailApi(data);
    if (res.code == 200) {
      this.strategyDay = Number(res.data.setting_day) || 30;
      this.strategyName = res.data.strategy_name;
      const dayObj: any = startAndEndOfday(res.data.start_time);
      if (res.data.strategy_method == 1) {
        //长期执行
        this.runTimeRange = [
          new Date(res.data.start_time * 1000),
          new Date(res.data.strategy_end_time * 1000)
        ];
      } else {
        //立即执行
        this.runTimeRange = [dayObj.start, dayObj.end];
      }
    } else {
      this.$message.error(res.message);
    }
    this.getTrackListFn();
  }
  async saveUserBeFn() {
    try {
      const data: object = {
        id: Number(this.$route.query.id),
        setting_day: this.strategyDay
      };
      const res: any = await API.marketingUpdateSettingApi(data);
      if (res.code == 200) {
        this.dialogVisible = false;
        this.getTrackListFn();
      } else {
        this.$message.error(res.message);
      }
    } catch (error) {
      console.warn("业绩设置失败！");
    }
  }
  goBack() {
    this.$router.push({
      name: "myapp_customer_auto-market_market-strategy_market-list",
      params: JSON.parse(JSON.stringify(this.$route.query))
    });
  }
  renderHeader(h: any, b: any) {
    const { column } = b;
    return h(
      "el-tooltip",
      {
        props: {
          content: `${this.config[column.property]}`,
          placement: "top",
          effect: "light"
        },
        domProps: {
          innerHTML: column.label + ' <span class="el-icon-question"></span>'
        }
      },
      [h("span")]
    );
  }
  openSettingFn() {
    this.dialogVisible = true;
  }
  packageParamFn() {
    const param: any = {
      day: this.strategyDay,
      id: Number(this.$route.query.id),
      begin_time: this.runTimeRange[0]
        ? new Date(this.runTimeRange[0]).getTime() / 1000
        : 0,
      end_time: this.runTimeRange[1]
        ? new Date(this.runTimeRange[1]).getTime() / 1000
        : 0
    };
    return param;
  }
  async downloadExeclFn() {
    try {
      const data: any = this.packageParamFn();
      let res: any;
      if (this.$route.query.way == "自动发券") {
        res = await API.exportCouponApi(data);
      } else {
        res = await API.marketStrategyExportApi(data);
      }
      const blob: Blob = new Blob([res], {
        type: "application/x-xls;charset=utf-8"
      });
      const a: HTMLElement = document.createElement("a");
      const href = window.URL.createObjectURL(blob);
      const name = `${this.strategyName}${dateFormat(
        Date.now() / 1000,
        "yyyy_M_d"
      )}.xlsx`;
      a.setAttribute("download", name);
      a.setAttribute("href", href);
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      window.URL.revokeObjectURL(href);
    } catch (error) {
      console.warn("导出失败！");
    }
  }
  async getTrackListFn() {
    let list: any = [];
    const data: any = this.packageParamFn();
    try {
      let res: any;
      if (this.$route.query.way == "自动发券") {
        this.isCoupon = true;
        res = await API.couponDataApi(data);
        list = res.data.data.list.coupon_list;
        list.forEach((ele: any) => {
          ele.type = 3;
          ele.send_num = res.data.data.list.send_num;
          ele.use_send =
            ele.use_num && ele.is_send_num
              ? Math.round((ele.use_num / ele.is_send_num) * 10000) / 100
              : 0;
          ele.pay_touch =
            ele.pay_num && ele.touch_num
              ? Math.round((ele.pay_num / ele.touch_num) * 10000) / 100
              : 0;
        });
        const total = res.data.data.total.coupon_list;
        total.send_num = res.data.data.list.send_num;
        total.use_send =
          (total.use_num && total.is_send_num
            ? Math.round((total.use_num / total.is_send_num) * 10000) / 100
            : 0) + "%";
        total.pay_touch =
          (total.pay_num && total.touch_num
            ? Math.round((total.pay_num / total.touch_num) * 10000) / 100
            : 0) + "%";
        this.total = total;
        this.listcoupon = list;
      } else {
        this.isCoupon = false;
        res = await API.marketStrategyDataListApi(data);
        list = res.data.data;
        this.total = res.data.total;
        this.list = list;
      }
    } catch (error) {
      console.log(error);
      console.warn("获取表格信息失败！");
    }
  }
  private objectSpanMethod(item: any) {
    const { rowIndex, columnIndex } = item;
    if (columnIndex === 0 || columnIndex === 1) {
      if (rowIndex % this.listcoupon.length === 0) {
        return {
          rowspan: this.listcoupon.length,
          colspan: 1
        };
      } else {
        return {
          rowspan: 0,
          colspan: 0
        };
      }
    }
  }
  private getSummaries(params: any) {
    const { columns, data } = params;
    const isAllWx = data.every((item: any) => item.type == 2);
    const sums: Array<number | string> = [];
    columns.forEach((column: any, index: number) => {
      if (index === 0) {
        sums[index] = "合计";
        return;
      }
      if (index === 3 && isAllWx) {
        sums[index] = "--";
        return;
      }
      if (this.total) {
        sums[index] = this.total[column.property];
      } else {
        sums[index] = 0;
      }
    });
    return sums;
  }
  mounted() {
    this.initSettingDay();
    this.updateTime = new Date().getTime() / 1000;
  }
}
</script>
<style lang="scss" scoped>
.pageDataTrack {
  box-shadow: 0px 0px 20px 1px rgba(0, 0, 0, 0.08);
  background-color: #ffffff;
  padding: 0 0 30px 0;
  border-radius: 10px;

  .marginL20 {
    margin-left: 20px;
  }

  .table-empty-label {
    display: flex;
    height: 340px;
    flex-direction: column;

    img {
      width: 124px;
      margin: 80px auto 0;
    }

    div {
      line-height: 20px;
      margin-top: 14px;
    }
  }

  .search_container {
    padding: 20px;
    background: $color-background;
    margin-bottom: 20px;

    .time-show {
      color: $color-text-normal;
      font-size: 14px;
      line-height: 20px;
      margin-bottom: 20px;

      span {
        font-size: 14px;
        color: $color-main;
        line-height: 20px;
        margin-left: 10px;
        cursor: pointer;
      }
    }
  }

  .settingDay {
    .title {
      padding-right: 5px;
    }

    .el-select.el-select--small {
      width: 62px !important;
    }

    .day {
      padding-left: 10px;
    }
  }

  .container {
    padding: 30px;
    background: #ffffff;

    .container_btns {
      margin-bottom: 20px;
    }
  }

  .dialog-footer {
    text-align: left;
  }
}

/deep/ .el-date-sty {
  height: 36px;
  vertical-align: middle;
}</style>
