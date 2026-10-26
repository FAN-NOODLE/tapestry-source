<template>
  <div class="guide-data">
    <layout-list>
      <template slot="search">
        <search-item text="日期：">
          <el-date-picker
            v-model="times"
            type="month"
            placeholder="选择月"
            format="yyyy年MM月"
            :clearable="false"
            value-format="yyyy-MM"
            @change="getTimes()"
          >
          </el-date-picker>
        </search-item>
        <search-item text="大区主管：">
          <cy-select
            v-if="userInfo.is_super === 1"
            v-model="dept_leader_mobile"
            @change="baSourceDepte()"
            :optionsList="supervisorList"
          >
          </cy-select>
          <span v-else>{{ counter_name }}</span>
        </search-item>
        <search-item text="分区主管：">
          <cy-select
            v-model="staff.value"
            :optionsList="searchStaffList"
            :options="{ label: 'ba_name', value: 'ba_mobile' }"
          >
          </cy-select>
        </search-item>
        <div class="operate-btn">
          <el-button type="primary" size="small" @click="search"
            >查询</el-button
          >
          <el-button type="primary" size="small" class="reset" @click="reset"
            >重置</el-button
          >
        </div>
      </template>
      <template slot="tabs">
        <h4 class="title">销售数据总览</h4>
        <div class="chart-data-wrpaer">
          <div class="left-box">
            <div v-for="(item, index) in data1" :key="index" class="info-item">
              <span>{{ item.name }}</span>
              <el-popover
                :ref="`popoverA${index}`"
                placement="top"
                width="270"
                trigger="hover"
                visible-arrow="true"
                :content="item.desc"
              >
                <svg-icon
                  slot="reference"
                  class="help-icon"
                  icon-class="group-help"
                ></svg-icon>
              </el-popover>
              <div class="sum">{{ item.value }}</div>
            </div>
          </div>
          <div class="right-box">
            <!-- 同比增长柱状图 -->
            <div id="histogram-1" class="histogram"></div>
            <div class="histogram-text">
              <span>同比增长</span>
              <el-popover
                placement="top"
                width="270"
                trigger="hover"
                visible-arrow="true"
                :content="'对比去年今月数据（MTD）'"
              >
                <svg-icon
                  slot="reference"
                  class="help-icon"
                  icon-class="group-help"
                ></svg-icon>
              </el-popover>
              <div class="sum">{{ growth_data.growth_rate }}</div>
              <p>
                {{ this.lastTimes }}销售额¥{{ growth_data.current_sale_amount }}
              </p>
              <p>{{ this.times }}销售额¥{{ growth_data.last_sale_amount }}</p>
            </div>
          </div>
        </div>
        <div class="chart-data—1-wrpaer">
          <div class="left-box">
            <div v-for="(item, index) in data2" :key="index" class="info-item">
              <span>{{ item.name }}</span>
              <el-popover
                v-if="item.name"
                :ref="`popoverA${index}`"
                placement="top"
                width="270"
                trigger="hover"
                visible-arrow="true"
                :content="item.desc"
              >
                <svg-icon
                  slot="reference"
                  class="help-icon"
                  icon-class="group-help"
                ></svg-icon>
              </el-popover>
              <div class="sum">{{ item.value }}</div>
            </div>
          </div>
          <div class="right-box">
            <div id="histogram-2" class="histogram-2"></div>
          </div>
        </div>
        <h4 class="title">运营数据总览</h4>
        <div class="chart-data-item-wraper chart-data—1-wrpaer">
          <div class="left-box">
            <div v-for="(item, index) in data3" :key="index" class="info-item">
              <span>{{ item.name }}</span>
              <el-popover
                :ref="`popoverA${index}`"
                placement="top"
                width="270"
                trigger="hover"
                visible-arrow="true"
                :content="item.desc"
              >
                <svg-icon
                  slot="reference"
                  class="help-icon"
                  icon-class="group-help"
                ></svg-icon>
              </el-popover>
              <div class="sum">{{ item.value }}</div>
            </div>
          </div>
          <div class="right-item-box">
            <div id="pie-chart-1" class="pie-chart"></div>
          </div>
        </div>
        <div class="chart-data-item-wraper chart-data—1-wrpaer">
          <div class="left-box">
            <div v-for="(item, index) in data4" :key="index" class="info-item">
              <span>{{ item.name }}</span>
              <el-popover
                :ref="`popoverA${index}`"
                placement="top"
                width="270"
                trigger="hover"
                visible-arrow="true"
                :content="item.desc"
              >
                <svg-icon
                  slot="reference"
                  class="help-icon"
                  icon-class="group-help"
                ></svg-icon>
              </el-popover>
              <div class="sum">{{ item.value }}</div>
            </div>
          </div>
          <div class="right-item-box">
            <div id="pie-chart-2" class="pie-chart"></div>
          </div>
        </div>
        <div class="chart-data-item-wraper chart-data—1-wrpaer">
          <div class="left-box">
            <div v-for="(item, index) in data5" :key="index" class="info-item">
              <span>{{ item.name }}</span>
              <el-popover
                v-if="item.name"
                :ref="`popoverA${index}`"
                placement="top"
                width="270"
                trigger="hover"
                visible-arrow="true"
                :content="item.desc"
              >
                <svg-icon
                  slot="reference"
                  class="help-icon"
                  icon-class="group-help"
                ></svg-icon>
              </el-popover>
              <div class="sum">{{ item.value }}</div>
            </div>
          </div>
          <div class="right-item-box">
            <div id="pie-chart-3" class="pie-chart" style="display:none"></div>
          </div>
        </div>
        <h4 class="title">销售数据明细</h4>
        <div class="listTitle">
          <!-- <el-tabs
            v-model="sale_list_data_type"
            type="card"
            @tab-click="changeSaleListType"
          >
            <el-tab-pane label="默认展示" name="1"></el-tab-pane>
            <el-tab-pane label="排行展示" name="2"></el-tab-pane>
          </el-tabs> -->
          <el-button
            :disabled="saleList.length == 0"
            size="small"
            @click="exportList1"
            class="export"
            >导出Excel
          </el-button>
          <el-table
            ref="dataList"
            :data="saleList"
            size="mini"
            stripe
            class="mar-t-15 list-last-line-none min-height-100"
            :header-cell-style="{
              fontSize: '14px',
              color: '#606266',
              fontWeight: '400',
            }"
          >
            <el-table-column label="导购" prop="ba_name"></el-table-column>
            <el-table-column label="销售额" min-width="120">
              <template slot-scope="scope">
                {{ scope.row.sales.toLocaleString() }}
              </template>
            </el-table-column>
            <el-table-column label="目标金额" min-width="120">
              <template slot-scope="scope">
                {{ scope.row.monthly_sales.toLocaleString() }}
              </template>
            </el-table-column>
            <el-table-column label="完成率">
              <template slot-scope="scope">
                {{ scope.row.finish_rate + "%" }}
              </template>
            </el-table-column>
            <el-table-column label="差额" min-width="120">
              <template slot-scope="scope">
                {{ scope.row.difference_sales.toLocaleString() }}
              </template>
            </el-table-column>
            <el-table-column label="订单数" prop="order_num"></el-table-column>
            <el-table-column label="客单价" min-width="120">
              <template slot-scope="scope">
                {{ scope.row.order_num_price.toLocaleString() }}
              </template>
            </el-table-column>
            <el-table-column
              label="客人数"
              prop="customer_num"
              min-width="100"
            ></el-table-column>
            <el-table-column
              label="新客人数"
              prop="new_customer_num"
            ></el-table-column>
            <el-table-column label="新客销售额" min-width="120">
              <template slot-scope="scope">
                {{ scope.row.new_customer_sales.toLocaleString() }}
              </template>
            </el-table-column>
            <el-table-column label="新客销售占比">
              <template slot-scope="scope">
                {{ scope.row.new_customer_rate + "%" }}
              </template>
            </el-table-column>
            <el-table-column
              label="新二购人数"
              prop="new_second_customer_num"
            ></el-table-column>
            <el-table-column label="新二购销售额" min-width="140">
              <template slot-scope="scope">
                {{ scope.row.new_second_customer_sales.toLocaleString() }}
              </template>
            </el-table-column>
            <el-table-column label="新二购销售占比">
              <template slot-scope="scope">
                {{ scope.row.new_second_customer_sales_rate + "%" }}
              </template>
            </el-table-column>
            <el-table-column
              label="老客人数"
              prop="old_customer_num"
            ></el-table-column>
            <el-table-column label="老客销售额" min-width="120">
              <template slot-scope="scope">
                {{ scope.row.old_customer_sales.toLocaleString() }}
              </template>
            </el-table-column>
            <el-table-column label="老客销售占比">
              <template slot-scope="scope">
                {{ scope.row.old_customer_sales_rate + "%" }}
              </template>
            </el-table-column>
            <div slot="empty">
              <layout-empty :message="'暂无数据'"></layout-empty>
            </div>
          </el-table>
          <el-pagination
            v-if="saleList && saleList.length > 0"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
            :current-page="sale_meta.page"
            :page-sizes="[10, 20, 30, 40]"
            :page-size="sale_meta.size"
            class="page-right font-weight-400"
            layout="total, sizes, prev, pager, next, jumper"
            :total="sale_meta.total"
          >
          </el-pagination>
        </div>
        <h4 class="title mt40">运营数据明细</h4>
        <div class="listTitle">
          <!-- <el-tabs
            v-model="operation_list_data_type"
            type="card"
            @tab-click="changeOperationListType"
          >
            <el-tab-pane label="默认展示" name="1"></el-tab-pane>
            <el-tab-pane label="排行展示" name="2"></el-tab-pane>
          </el-tabs> -->
          <el-button
            :disabled="operationList.length == 0"
            size="small"
            @click="exportList2"
            class="export"
            >导出Excel
          </el-button>
          <el-table
            ref="dataList"
            :data="operationList"
            size="mini"
            stripe
            class="mar-t-15 list-last-line-none min-height-100"
            :header-cell-style="{
              fontSize: '14px',
              color: '#606266',
              fontWeight: '400',
            }"
          >
            <el-table-column label="导购" prop="ba_name"></el-table-column>
            <el-table-column
              label="新二购且一码三关注指标人数"
              prop="monthly_second_buy_and_three_scan"
            ></el-table-column>
            <el-table-column
              label="新二购且一码三关注完成人数"
              prop="new_second_and_three_scan"
            ></el-table-column>
            <el-table-column label="新二购一码三关注率">
              <template slot-scope="scope">
                {{ scope.row.new_second_and_three_scan_rate + "%" }}
              </template>
            </el-table-column>
            <el-table-column label="新二购工资完成进度">
              <template slot-scope="scope">
                {{ scope.row.new_second_and_three_scan_wage_progress + "%" }}
              </template>
            </el-table-column>
            <el-table-column
              label="老客且一码三关注指标人数"
              prop="monthly_old_customer_and_three_scan"
            ></el-table-column>
            <el-table-column
              label="老客且一码三关注完成人数"
              prop="old_and_three_scan"
            ></el-table-column>
            <el-table-column
              label="老客人数"
              prop="old_x_customer_num"
            ></el-table-column>
            <el-table-column label="老客一码三关注率">
              <template slot-scope="scope">
                {{ scope.row.monthly_old_customer_and_three_scan_rate + "%" }}
              </template>
            </el-table-column>
            <el-table-column label="老客工资完成进度">
              <template slot-scope="scope">
                {{
                  scope.row.monthly_old_customer_and_three_wage_progress + "%"
                }}
              </template>
            </el-table-column>
            <el-table-column label="新客且一码三关注指标人数">
              <template>0</template>
            </el-table-column>
            <el-table-column
              label="新客且一码三关注完成人数"
              prop="new_customer_and_three_scan"
            ></el-table-column>
            <el-table-column
              label="新客人数"
              prop="monthly_new_customer_and_three_scan"
            ></el-table-column>
            <el-table-column label="新客人一码三关注率">
              <template slot-scope="scope">
                {{ scope.row.new_customer_and_three_scan_rate + "%" }}
              </template>
            </el-table-column>
            <div slot="empty">
              <layout-empty :message="'暂无数据'"></layout-empty>
            </div>
          </el-table>
          <el-pagination
            v-if="operationList && operationList.length > 0"
            @size-change="handleSizeChange1"
            @current-change="handleCurrentChange1"
            :current-page="operation_meta.page"
            :page-sizes="[10, 20, 30, 40]"
            :page-size="operation_meta.size"
            class="page-right font-weight-400"
            layout="total, sizes, prev, pager, next, jumper"
            :total="operation_meta.total"
          >
          </el-pagination>
        </div>
      </template>
    </layout-list>
  </div>
</template>

<script lang="ts">
interface Params {
  month: string | number;
  default_mobile: string;
  dept_leader_mobile: string;
  ba_mobile: string | number;
  enterprise_id: string | number;
  page?: number;
  page_size?: number;
}
import { Component, Vue } from "vue-property-decorator";
import DataCompassApi from "@/api/wisdom/dataCompass";
import LayoutEmpty from "@/components/sy-components/layout-empty.vue";
@Component({
  components: {
    LayoutEmpty,
  },
})
export default class extends Vue {
  // data
  private userInfo: any = {};
  private times: any = ""; //用户信息
  private lastTimes: any = "";
  private counter_name = "";
  private dept_leader_mobile = "";
  private supervisorList: Array<any> = [
    {
      label: "陈静",
      value: "15692163738",
    },
    {
      label: "杨林川",
      value: "13980001181",
    },
    {
      label: "王彪",
      value: "18643147744",
    },
  ];
  private staff = {
    // 员工
    value: "全部",
  };
  // 员工列表
  private searchStaffList: Array<any> = [];
  // 销售数据
  private data1 = [
    {
      name: "销售目标 (元）",
      desc: "当月目标设置",
      key: "monthly_sales",
      value: 0,
    },
    {
      name: "销售额 (元）",
      desc: "当月销售业绩",
      key: "sales",
      value: 0,
    },
    {
      name: "销售完成率",
      desc: "销售额 / 完成率 *100%",
      key: "finish_rate",
      value: "0%",
    },
    {
      name: "销售差额 (元）",
      desc: "目标金额 - 销售额",
      key: "difference_sales",
      value: 0,
    },
    {
      name: "客单价 (元）",
      desc: " 当月销售额/当月全部订单",
      key: "order_num_price",
      value: 0,
    },
    {
      name: "订单数 (笔）",
      desc: "当月全部订单",
      key: "order_num",
      value: 0,
    },
  ];
  //
  private data2: any = [
    {
      name: "客人数 (人）",
      desc: "当月发生购买的客人数量",
      value: 0,
    },
    {
      name: "",
      desc: "",
      value: null,
    },
    {
      name: "",
      desc: "",
      value: null,
    },
    {
      name: "新客数 (人）",
      desc: "当月在品牌内发生首单购买的客人数量",
      value: 0,
    },
    {
      name: "新客销售额 (元）",
      desc: "新客人数对应的销售额",
      value: 0,
    },
    {
      name: "新客销售占比",
      desc: "新客销售额 / 总销售额 * 100%",
      value: 0,
    },
    {
      name: "老客数 (人）",
      desc:
        "当月的购买为全渠道的第三单或以上的购买，且品牌内首单购买距离当月<=12个自然月的客人数量 或，在品牌内首购时间距离当月>12个自然月的客人数量",
      value: 0,
    },
    {
      name: "老客销售额 (元）",
      desc: "老客人数对应的销售额",
      value: 0,
    },
    {
      name: "老客销售占比",
      desc: "老客销售额 / 总销售额 * 100%",
      value: 0,
    },
    {
      name: "新二购人数 (人）",
      desc:
        "当月的购买为全渠道的第二单购买，且品牌内首单购买距离当月<=12个自然月的客人数量",
      value: 0,
    },
    {
      name: "新二购销售额 (元）",
      desc: "新二购人数对应的销售额",
      value: 0,
    },
    {
      name: "新二购销售占比",
      desc: "新二购销售额 / 总销售额 * 100%",
      value: 0,
    },
    {
      name: "游客人数 (人）",
      desc: "参与购买但未注册会员的人员，1单算一个游客",
      value: 0,
    },
    {
      name: "游客销售额 (元）",
      desc: "游客人数对应的销售额",
      value: 0,
    },
    {
      name: "游客销售占比",
      desc: "游客销售额 / 总销售额 * 100%",
      value: 0,
    },
  ];
  private data3: any = [
    {
      name: "工资完成进度",
      desc:
        "新二购且一码三关注的完成人数 / 新二购且一码三关注的指标人数 * 100%",
      value: 0,
    },
    {
      name: "新二购且一码三关注指标 (人）",
      desc: "新二购且一码三关注的指标 人数",
      value: 0,
    },
    {
      name: "新二购且一码三关注完成 (人）",
      desc: "新二购且一码三关注的完成人数",
      value: 0,
    },
    {
      name: "新二购人数 (人）",
      desc:
        "当月的购买为全渠道的第二单购买，且品牌内首单购买距离当月<=12个自然月的客人数量导购工资计算逻辑：- 若第二单购买渠道为线下，则第二单的接待导购为获取新二购奖金的导购- 若第二单购买渠道为线上，则首单的接待导购为获取新二购奖金的导购",
      value: 0,
    },
    {
      name: "一码三关注率",
      desc: " 老客且一码三关注的完成人数 / 老客人数 * 100%",
      value: "0%",
    },
  ];
  private data4: any = [
    {
      name: "工资完成进度",
      desc: "老客且一码三关注的完成人数 / 老客且一码三关注的指标人数 * 100%",
      value: 0,
    },
    {
      name: "老客且一码三关注指标 (人）",
      desc: "老客且一码三关注的指标人数",
      value: 0,
    },
    {
      name: "老客且一码三关注完成 (人）",
      desc: "老客且一码三关注的完成人数",
      value: 0,
    },
    {
      name: "老客人数 (人）",
      desc:
        "当月的购买为全渠道的第三单或以上的购买，且品牌内首单购买距离当月<=12个自然月的客人数量 或，在品牌内首购时间距离当月>12个自然月的客人数量 导购工资计算逻辑：- 若最近一单购买渠道为线下，则这一单的接待导购为获取老客奖金的导购- 若最近一单购买渠道为线上，则该客人的管理导购为获取老客奖金的导购",
      value: 0,
    },
    {
      name: "一码三关注率",
      desc: "老客且一码三关注的完成人数 / 老客人数 * 100%",
      value: "0%",
    },
  ];
  private data5: any = [
    {
      name: "新客人数 (人）",
      desc: "当月在品牌内发生首单购买的客人数量",
      value: 0,
    },
    {
      name: "新客且一码三关注完成 (人）",
      desc: "新客且一码三关注的完成人数",
      value: 0,
    },
    {
      name: "一码三关注率",
      desc: "新客且一码三关注的完成人数 / 新客人数 * 100%",
      value: 0,
    },
    {
      name: "",
      desc: "",
      key: "",
      value: "",
    },
  ];
  private growth_data = {
    // 同比增长
    current_sale_amount: 0,
    last_sale_amount: 0,
    growth_rate: "0%",
  };
  private h_1_data: {
    color: Array<string>;
    title: any;
    grid: any;
    xAxis: any;
    yAxis: any;
    series: any;
  } = {
    color: ["#FAC858"],
    title: { text: "" },
    grid: {
      left: "2%",
      right: "2%",
      bottom: "0%",
      top: "15%",
      containLabel: true,
    },
    xAxis: {
      type: "category",
      data: [],
    },
    yAxis: [
      {
        type: "value",
      },
    ],
    series: [
      {
        type: "bar",
        data: [],
      },
    ],
  };
  private h_2_data: {
    legend: any;
    color: Array<string>;
    title: any;
    grid: any;
    dataset: any;
    xAxis: any;
    yAxis: any;
    series: any;
  } = {
    legend: {},
    color: ["#91CD75", "#FAC858", "#EF6767"],
    title: { text: "" },
    grid: {
      left: "2%",
      right: "2%",
      bottom: "0%",
      top: "15%",
      containLabel: true,
    },
    dataset: {
      source: [],
    },
    xAxis: {
      type: "category",
    },
    yAxis: {},
    series: [
      {
        type: "bar",
      },
      {
        type: "bar",
      },
      {
        type: "bar",
      },
    ],
  };
  private p_1_data: {
    color: Array<string>;
    title: any;
    graphic: any;
    series: any;
  } = {
    color: ["#EE6767", "#72BFDE"],
    title: {
      text: "",
      left: "center",
      top: "55%",
      textStyle: {
        fontSize: 24,
        align: "center",
      },
    },
    graphic: {
      type: "text",
      left: "center",
      top: "38%",
      style: {
        text: "新二购指\n标完成率",
        textAlign: "center",
        fill: "#333",
        fontSize: 18,
      },
    },
    series: [
      {
        name: "",
        type: "pie",
        radius: ["85%", "70%"],
        avoidLabelOverlap: false,
        hoverAnimation: false,
        label: {
          normal: {
            show: false,
            position: "center",
          },
        },
        data: [],
      },
    ],
  };
  private p_2_data: {
    color: Array<string>;
    title: any;
    graphic: any;
    series: any;
  } = {
    color: ["#EE6767", "#72BFDE"],
    title: {
      text: "",
      left: "center",
      top: "55%",
      textStyle: {
        fontSize: 24,
        align: "center",
      },
    },
    graphic: {
      type: "text",
      left: "center",
      top: "38%",
      style: {
        text: "老客指标\n完成率",
        textAlign: "center",
        fill: "#333",
        fontSize: 18,
      },
    },
    series: [
      {
        name: "",
        type: "pie",
        radius: ["85%", "70%"],
        avoidLabelOverlap: false,
        hoverAnimation: false,
        label: {
          normal: {
            show: false,
            position: "center",
          },
        },
        data: [],
      },
    ],
  };
  private p_3_data: {
    color: Array<string>;
    title: any;
    graphic: any;
    series: any;
  } = {
    color: ["#EE6767", "#72BFDE"],
    title: {
      text: "",
      left: "center",
      top: "55%",
      textStyle: {
        fontSize: 24,
        align: "center",
      },
    },
    graphic: {
      type: "text",
      left: "center",
      top: "38%",
      style: {
        text: "新客指标\n完成率",
        textAlign: "center",
        fill: "#333",
        fontSize: 18,
      },
    },
    series: [
      {
        name: "",
        type: "pie",
        radius: ["85%", "70%"],
        avoidLabelOverlap: false,
        hoverAnimation: false,
        label: {
          normal: {
            show: false,
            position: "center",
          },
        },
        data: [],
      },
    ],
  };
  private params: Params = {
    month: "",
    ba_mobile: "",
    dept_leader_mobile: "",
    default_mobile: "",
    enterprise_id: "",
  };
  private sale_meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10, // 每页数量
  };
  private sale_list_data_type = "1";
  private saleList = [];
  private operation_meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10, // 每页数量
  };
  private operation_list_data_type = "1";
  private operationList = [];
  // 生命周期
  private mounted() {
    this.userInfo = JSON.parse(localStorage.getItem("USERINFO") || "{}");
    const date: any = new Date();
    this.times = date.getFullYear() + "-" + (date.getMonth() + 1);
    this.lastTimes = date.getFullYear() - 1 + "-" + (date.getMonth() + 1);
    if (this.userInfo.is_super === 1) {
      this.dept_leader_mobile = this.supervisorList[0].value;
    } else {
      this.dept_leader_mobile = "杨林川";
    }
    this.baSourceDepte();
  }
  private baSourceDepte() {
    this.searchStaffList = [];
    this.staff.value = "全部";
    const params: any = {
      enterprise_id: this.userInfo.enterprise_id,
    };
    if (this.userInfo.is_super === 1) {
      params.mobile = this.dept_leader_mobile;
    } else {
      params.mobile = "13980001181";
    }
    DataCompassApi.baSourceDepte(params).then((res: any) => {
      if (res.code === 200) {
        this.searchStaffList = res.data.leader_ba;
        this.counter_name = res.data.counter_name;
        if (
          this.userInfo.is_super !== 1 &&
          this.searchStaffList &&
          this.searchStaffList.length > 0
        ) {
          this.dept_leader_mobile = this.searchStaffList[0].regional_leader_mobie;
        }
        this.search();
      }
    });
  }
  private getTimes() {
    const [year, month] = this.times.split("-");
    this.lastTimes = [year - 1, month].join("-");
    this.search();
  }
  private search() {
    this.params = {
      month: new Date(this.times).getTime() / 1000,
      default_mobile: this.dept_leader_mobile || "13980001181", //默认查询当前人的手机号
      dept_leader_mobile: this.dept_leader_mobile,
      ba_mobile: this.staff.value === "全部" ? "" : this.staff.value,
      enterprise_id: this.userInfo.enterprise_id,
    };
    Promise.all([
      DataCompassApi.baKanBan(this.params),
      DataCompassApi.baKanBanGrowthRate(this.params),
    ]).then((res: any) => {
      const data = res[0].data;
      const growth_data = res[1].data;
      // 销售数据
      this.data1[0].value = data.monthly_sales.toLocaleString();
      this.data1[1].value = data.sales.toLocaleString();
      this.data1[2].value = data.finish_rate + "%";
      this.data1[3].value = data.difference_sales.toLocaleString();
      this.data1[4].value = data.order_num_price.toLocaleString();
      this.data1[5].value = data.order_num;
      // 同比增长数据
      this.growth_data.current_sale_amount = growth_data.beforeyearMonth.sales;
      this.growth_data.last_sale_amount = growth_data.nowMonth.sales;
      this.growth_data.growth_rate = growth_data.nowMonth.growth_rate + "%";
      this.h_1_data.xAxis.data = [
        new Date(this.lastTimes).getFullYear(),
        new Date(this.times).getFullYear(),
      ];
      this.h_1_data.series[0].data = [
        growth_data.beforeyearMonth.sales,
        growth_data.nowMonth.sales,
      ];
      // 客群数据
      this.data2[0].value = data.customer_num;
      this.data2[3].value = data.new_customer_num;
      this.data2[4].value = data.new_customer_sales.toLocaleString();
      this.data2[5].value = data.new_customer_rate + "%";
      this.data2[6].value = data.old_customer_num;
      this.data2[7].value = data.old_customer_sales.toLocaleString();
      this.data2[8].value = data.old_customer_sales_rate + "%";
      this.data2[9].value = data.new_second_customer_num;
      this.data2[10].value = data.new_second_customer_sales.toLocaleString();
      this.data2[11].value = data.new_second_customer_sales_rate + "%";
      this.data2[12].value = data.visit_customer_num;
      this.data2[13].value = data.visit_customer_sales.toLocaleString();
      this.data2[14].value = data.visit_customer_sales_rate + "%";
      this.h_2_data.dataset.source = [
        ["product", "销售额", "人数", "占比"],
        [
          "新客",
          data.new_customer_sales,
          this.data2[3].value,
          data.new_customer_rate,
        ],
        [
          "老客",
          data.old_customer_sales,
          this.data2[6].value,
          data.old_customer_sales_rate,
        ],
        [
          "新二购",
          data.new_second_customer_sales,
          this.data2[9].value,
          data.new_second_customer_sales_rate,
        ],
        ["游客", data.visit_customer_sales, this.data2[12].value, 0],
      ];
      // // 运营数据总览
      // 新二购
      this.data3[0].value = data.new_second_and_three_scan_wage_progress + "%";
      this.data3[1].value = data.monthly_second_buy_and_three_scan;
      this.data3[2].value = data.new_second_and_three_scan;
      this.data3[3].value = data.new_second_customer_y_num;
      this.data3[4].value = data.new_second_and_three_scan_rate + "%";
      this.p_1_data.title.text = this.data3[0].value;
      if (data.new_second_and_three_scan_wage_progress >= 100) {
        this.p_1_data.series[0].data = [
          { value: 1, name: "已完成" },
          { value: 0, name: "未完成" },
        ];
      } else {
        this.p_1_data.series[0].data = [
          {
            value: data.new_second_and_three_scan_wage_progress,
            name: "已完成",
          },
          {
            value: 100 - data.new_second_and_three_scan_wage_progress,
            name: "未完成",
          },
        ];
      }
      // 老客
      this.data4[0].value =
        data.monthly_old_customer_and_three_wage_progress + "%";
      this.data4[1].value = data.monthly_old_customer_and_three_scan;
      this.data4[2].value = data.old_and_three_scan;
      this.data4[3].value = data.old_x_customer_num;
      this.data4[4].value = data.monthly_old_customer_and_three_scan_rate + "%";
      this.p_2_data.title.text = this.data4[0].value;
      if (data.monthly_old_customer_and_three_wage_progress >= 100) {
        this.p_2_data.series[0].data = [
          { value: 1, name: "已完成" },
          { value: 0, name: "未完成" },
        ];
      } else {
        this.p_2_data.series[0].data = [
          {
            value: data.monthly_old_customer_and_three_wage_progress,
            name: "已完成",
          },
          {
            value: 100 - data.monthly_old_customer_and_three_wage_progress,
            name: "未完成",
          },
        ];
      }
      // 新客
      this.data5[0].value = data.monthly_new_customer_and_three_scan;
      this.data5[1].value = data.new_customer_and_three_scan;
      this.data5[2].value = data.new_customer_and_three_scan_rate + "%";
      this.p_3_data.title.text = this.data5[2].value + "%";
      this.p_3_data.series[0].data = [{ value: 1, name: "已完成" }];
      this.draw();
    });
    this.getSalesList();
    this.getOperationList();
  }

  private reset() {
    const date: any = new Date();
    this.times = date.getFullYear() + "-" + (date.getMonth() + 1);
    if (this.userInfo.is_super === 1) {
      this.dept_leader_mobile = this.supervisorList[0].value;
    }
    this.staff.value = "全部";
    this.search();
  }

  private draw() {
    const category1 = this.$echarts.init(
      document.getElementById("histogram-1")
    );
    const category2 = this.$echarts.init(
      document.getElementById("histogram-2")
    );
    const pieChart1 = this.$echarts.init(
      document.getElementById("pie-chart-1")
    );
    const pieChart2 = this.$echarts.init(
      document.getElementById("pie-chart-2")
    );
    const pieChart3 = this.$echarts.init(
      document.getElementById("pie-chart-3")
    );
    category1.setOption(this.h_1_data, true);
    category2.setOption(this.h_2_data, true);
    pieChart1.setOption(this.p_1_data, true);
    pieChart2.setOption(this.p_2_data, true);
    pieChart3.setOption(this.p_3_data, true);
    window.onresize = function() {
      category1.resize();
      category2.resize();
      pieChart1.resize();
      pieChart2.resize();
      pieChart3.resize();
    };
  }

  // 切换销售数据
  private changeSaleListType() {
    this.sale_meta.page = 1;
  }

  // 切换运营数据
  private changeOperationListType() {
    this.operation_meta.page = 1;
  }

  // 获取销售数据列表
  private getSalesList() {
    const params = {
      page: this.sale_meta.page,
      page_size: this.sale_meta.size,
      ...this.params,
    };
    DataCompassApi.baKanBanSalesOrOperation(params).then((res: any) => {
      if (res.code == 200) {
        this.saleList = res.data.list || [];
        this.sale_meta.total = res.data.total;
      } else {
        this.$message.error(res.message);
      }
    });
  }

  // 获取运营数据列表
  private getOperationList() {
    const params = {
      page: this.operation_meta.page,
      page_size: this.operation_meta.size,
      ...this.params,
    };
    DataCompassApi.baKanBanSalesOrOperation(params).then((res: any) => {
      if (res.code == 200) {
        this.operationList = res.data.list || [];
        this.operation_meta.total = res.data.total;
      } else {
        this.$message.error(res.message);
      }
    });
  }

  // 销售数据每页
  private handleSizeChange(val: number) {
    this.sale_meta.page = 1;
    this.sale_meta.size = val;
    this.getSalesList();
  }
  // 销售数据当前页
  private handleCurrentChange(val: number) {
    this.sale_meta.page = val;
    this.getSalesList();
  }

  // 运营数据每页
  private handleSizeChange1(val: number) {
    this.operation_meta.page = 1;
    this.operation_meta.size = val;
    this.getOperationList();
  }
  // 运营数据当前页
  private handleCurrentChange1(val: number) {
    this.operation_meta.page = val;
    this.getOperationList();
  }

  // 导出销售数据
  private exportList1() {
    DataCompassApi.baKanBanSalesExport(this.params).then((res: any) => {
      const blob = new Blob([res], {
        type: "application/x-xls;charset=utf-8",
      });
      let fileName = "";
      const nowTime = new Date();
      fileName = "销售数据明细_" + this.formateDate(nowTime) + ".xlsx";

      //下载的文件名称及其后缀，后缀要和后台保持的一致

      if ("download" in document.createElement("a")) {
        // 非IE下载
        const elink = document.createElement("a");
        elink.download = fileName;
        elink.style.display = "none";
        elink.href = URL.createObjectURL(blob);
        document.body.appendChild(elink);
        elink.click();
        URL.revokeObjectURL(elink.href); // 释放URL 对象
        document.body.removeChild(elink);
      } else {
        // IE10+下载
        (navigator as any).msSaveBlob(blob, fileName);
      }
    });
  }

  // 导出运营数据
  private exportList2() {
    DataCompassApi.baKanBanOperationExport(this.params).then((res: any) => {
      const blob = new Blob([res], {
        type: "application/x-xls;charset=utf-8",
      });
      let fileName = "";
      const nowTime = new Date();
      fileName = "运营数据明细_" + this.formateDate(nowTime) + ".xlsx";

      //下载的文件名称及其后缀，后缀要和后台保持的一致

      if ("download" in document.createElement("a")) {
        // 非IE下载
        const elink = document.createElement("a");
        elink.download = fileName;
        elink.style.display = "none";
        elink.href = URL.createObjectURL(blob);
        document.body.appendChild(elink);
        elink.click();
        URL.revokeObjectURL(elink.href); // 释放URL 对象
        document.body.removeChild(elink);
      } else {
        // IE10+下载
        (navigator as any).msSaveBlob(blob, fileName);
      }
    });
  }

  // 时间格式转换
  private formateDate(time: any) {
    const year = time.getFullYear();
    let month = time.getMonth() + 1;
    let day = time.getDate();

    month = month < 10 ? "0" + month : month;
    day = day < 10 ? "0" + day : day;

    return year + "-" + month + "-" + day;
  }
}
</script>
<style lang="scss" scoped>
/deep/ .layout-list {
  padding: 10px 0 0 0 !important;
}
.operate-btn {
  display: inline-block;
  .reset {
    color: #303133;
    background-color: #fff;
    border-color: #c7cad0;
    margin-left: 15px;
  }
}
.chart-data-wrpaer {
  display: flex;
  .left-box {
    flex: 0.65;
    display: flex;
    flex-wrap: wrap;
    background: #f0f2f6;
    border-radius: 2px;
    padding: 20px;
  }
  .right-box {
    background: #f0f2f6;
    border-radius: 2px;
    margin-left: 20px;
    flex: 0.35;
    display: flex;
    padding: 10px 20px;
    .histogram {
      flex: 0.6;
      background: #fff;
      margin-right: 10px;
    }
    .histogram-text {
      flex: 0.4;
      margin: auto;
      .sum {
        font-size: 36px;
        display: block;
        font-weight: 600;
        margin-top: 20px;
        margin-bottom: 10px;
      }
    }
  }
}
.chart-data—1-wrpaer {
  display: flex;
  background: #f0f2f6;
  border-radius: 2px;
  margin-top: 10px;
  padding: 20px;
  .left-box {
    flex: 1;
    display: flex;
    flex-wrap: wrap;
  }
  .right-box {
    flex: 1;
    padding: 40px 20px;
    .histogram-2 {
      height: 100%;
      background: #fff;
    }
  }
}
.chart-data-item-wraper {
  min-height: 220px;
  .left-box {
    flex: 1;
  }
  .right-item-box {
    width: 220px;
    padding: 0;
    padding-right: 40px;
    .pie-chart {
      height: 100%;
      background: #fff;
    }
  }
}

.title {
  margin: 20px 0;
}
.info-item {
  width: 30%;
  margin-left: 0;
  text-align: center;
  padding: 20px 0;
  .info-tilte {
    font-size: 14px;
    padding: 20px 25px;
    border-radius: 5px;
    background: #f0f2f6;
    span {
      display: inline-block;
      margin-bottom: 15px;
    }
  }
  .sum {
    font-size: 26px;
    display: block;
    font-weight: 600;
    margin-top: 20px;
  }
}
/deep/ .el-popover__reference-wrapper {
  margin-left: 2px !important;
}
/deep/ .el-pagination {
  margin: 20px 0;
}
.mt40 {
  margin-top: 40px;
}
.export {
  margin-bottom: 20px;
}
</style>
