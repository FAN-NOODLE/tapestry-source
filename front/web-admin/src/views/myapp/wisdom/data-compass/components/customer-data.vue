<template>
  <div class="customer-data">
    <layout-list>
      <template slot="search">
        <search-item text="日期：">
          <el-radio-group
            class="search-info m-t-2 "
            v-model="time_type"
            size="small"
            @change="changeTime"
          >
            <el-radio-button
              class="search-data"
              v-for="(item, index) in time_list"
              :key="index"
              :label="item.label"
            >
              {{ item.name }}</el-radio-button>
          </el-radio-group>
          <el-date-picker
            class="m-l-20 choice-picker"
            v-model="times"
            type="datetimerange"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            size="medium"
            :default-time="['00:00:00', '23:59:59']"
            clearable
            @change="timeChange"
            value-format="timestamp"
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
        <div class="operate-btn">
          <el-button
            type="primary"
            size="small"
            @click="search"
          >查询</el-button>
          <el-button
            type="primary"
            size="small"
            class="reset"
            @click="reset"
          >重置</el-button>
        </div>
      </template>
      <template slot="tabs">
        <div>
          <div class="mar-t-20 flex-s-c">
            <div
              v-for="(item, index) in summary_data_customer"
              :key="index"
              class="info-tilte info-card mar-b-20"
            >
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
          <div class="mar-t-20 flex-s-c">
            <div
              v-for="(item, index) in summary_data_customer_add"
              :key="index"
              class="info-tilte info-card mar-b-20"
            >
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
          <div class="mar-t-20 flex-s-c customize">
            <div
              v-for="(item, index) in summary_data_customer_third"
              :key="index"
              class="info-tilte info-card mar-b-20"
              :class="{ 'no-data': !item.name }"
            >
              <div v-if="item.name">
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
              <div
                class
                v-else
              ></div>
            </div>
          </div>
          <!-- <div v-if="active_name != 2" class="mar-t-20 flex-s-c">
        <div class="flex-s-c-auto m-b-20">
          <div class="info-title">
            <span v-if="active_name == 1">推广订单数</span>
            <span v-if="active_name == 3">分销结算笔数</span>
            <el-popover
              ref="popoverA"
              placement="top"
              width="270"
              trigger="hover"
              visible-arrow="true"
              :content="
                active_name == 1
                  ? '导购使用在本店/本企业的导购身份，累计的推广订单总数'
                  : '属于导购名下的已结算状态的分销结算单笔数'
              "
            >
              <svg-icon
                slot="reference"
                class="help-icon"
                icon-class="group-help"
              ></svg-icon>
            </el-popover>
            <div class="sum">{{ summary_data.data1 || "0" }}</div>
          </div>
          <div class="info-title" v-if="active_name == 1">
            <span>导购推广订单数</span>
            <div class="sum">{{ summary_data.data5 || "0" }}</div>
          </div>
          <div class="info-title" v-if="active_name == 1">
            <span>合伙人推广订单数</span>
            <div class="sum">{{ summary_data.data6 || "0" }}</div>
          </div>
        </div>
        <div class="m-l-20 m-b-20">
          <div class="info-tilte">
            <span v-if="active_name == 1">推广订单实收金额(元)</span>
            <span v-if="active_name == 3">分销结算金额（元）</span>
            <el-popover
              ref="popoverB"
              placement="top"
              width="210"
              trigger="hover"
              visible-arrow="true"
              :content="
                active_name == 1
                  ? '导购名下推广订单的实收金额合计'
                  : '属于导购名下的已结算状态的分销结算单中的该导购推广佣金+服务佣金合计金额'
              "
            >
              <svg-icon
                slot="reference"
                class="help-icon"
                icon-class="group-help"
              ></svg-icon>
            </el-popover>
            <div class="sum">
              {{ summary_data.data2 / 100 || "0" }}
            </div>
          </div>
        </div>
        <div
          class="m-l-20 m-r-20 m-b-20"
          v-if="active_name == 1 && USER.version !== 3"
        >
          <div class="info-tilte">
            <span v-if="active_name == 1">服务订单数</span>
            <el-popover
              ref="popoverC"
              placement="top"
              width="240"
              trigger="hover"
              visible-arrow="true"
              :content="
                active_name == 1
                  ? '店员/店长在本店累计接单/核销的服务订单总数'
                  : ''
              "
            >
              <svg-icon
                slot="reference"
                class="help-icon"
                icon-class="group-help"
              ></svg-icon>
            </el-popover>
            <div class="sum">{{ summary_data.data3 || "0" }}</div>
          </div>
        </div>
        <div v-if="active_name == 1 && USER.version !== 3" class="m-b-20">
          <div class="info-tilte">
            <span v-if="active_name == 1">服务订单实收金额(元)</span>
            <el-popover
              ref="popoverD"
              placement="top"
              width="240"
              trigger="hover"
              visible-arrow="true"
              :content="
                active_name == 1 ? '导购名下的服务订单的实收金额合计' : ''
              "
            >
              <svg-icon
                slot="reference"
                class="help-icon"
                icon-class="group-help"
              ></svg-icon>
            </el-popover>
            <div class="sum">{{ summary_data.data4 / 100 || "0" }}</div>
          </div>
        </div>
      </div> -->
          <!-- =========折线图========= -->
          <div
            class="m-t-30"
            style="width: 100%"
            v-if="active_name === '2'"
          >
            <div
              id="diagram-box"
              v-show="e_data.xAxis.data.length > 0"
              class="echartBox"
            ></div>
            <div
              v-show="e_data.xAxis.data.length < 1"
              class="noData"
            >
              <span>所选日期暂无数据</span>
            </div>
          </div>
          <!-- =========饼图========= -->
          <div
            class="pte-container"
            v-if="active_name === '2'"
          >
            <div class="left-pie">
              <p class="pie-title">会员等级分布（新增会员）</p>
              <div
                id="pieChart1-box"
                class="pieChartBox"
                v-show="pice_data.series[0].data.length > 0"
              ></div>
              <div
                v-show="pice_data.series[0].data.length < 1"
                class="noData"
              >
                <span>所选日期暂无数据</span>
              </div>
            </div>
            <div class="right-pie">
              <p class="pie-title">会员等级分布（会员总数）</p>
              <div
                id="pieChart-box"
                class="pieChartBox"
                v-show="p_data.series[0].data.length > 0"
              ></div>
              <div
                v-show="p_data.series[0].data.length < 1"
                class="noData"
              >
                <span>所选日期暂无数据</span>
              </div>
            </div>
          </div>
          <div
            class="listTitle"
            v-if="active_name === '2'"
          >
            <el-tabs
              v-model="list_data_type"
              type="card"
              @tab-click="changeListType"
            >
              <el-tab-pane
                label="按日期"
                name="1"
              ></el-tab-pane>
              <el-tab-pane
                label="按导购"
                name="2"
              ></el-tab-pane>
            </el-tabs>
            <!-- <el-button size="mini" @click="exportExcel" class="export"
          >导出列表数据</el-button
        > -->
            <!-- <el-button
          size="mini"
          @click="exportOrderDetailExcel"
          class="export"
          v-if="active_name === '1' || active_name === '2'"
        >
          <span v-if="active_name === '1'">导出订单明细</span>
          <span v-if="active_name === '2'">导出绑客明细</span>
        </el-button> -->
          </div>
          <el-table
            ref="dataList"
            :data="table_list"
            @sort-change="sortChange"
            size="mini"
            stripe
            class="mar-t-15 list-last-line-none min-height-100"
            :header-cell-style="{
              fontSize: '14px',
              color: '#606266',
              fontWeight: '400'
            }"
          >
            <!-- 公共部分 -->
            <el-table-column
              v-if="list_data_type == 1"
              :fixed="list_data_type == 1"
              label="日期"
              min-width="120"
              prop="time"
            ></el-table-column>
            <el-table-column
              v-if="list_data_type == 2"
              :fixed="list_data_type == 2"
              label="导购ID"
              min-width="140"
              prop="worker_no"
            ></el-table-column>
            <el-table-column
              v-if="list_data_type == 2"
              label="导购"
              :fixed="list_data_type == 2"
              min-width="150"
              prop="totally"
            >
              <div
                slot-scope="scope"
                v-if="list_data_type == 2"
              >
                <img
                  :src="
                    scope.row.avatar
                      ? scope.row.avatar
                      : 'https://oss-crs.vchangyi.com/wechat/default-avatar-icon.png'
                  "
                  alt=""
                  width="30"
                  height="30"
                  class="table-text-img"
                />
                <span v-text="scope.row.worker_name"></span>
              </div>
            </el-table-column>
            <el-table-column
              v-if="list_data_type == 2 && USER.version !== 3"
              label="导购所在门店"
              min-width="140"
              prop="store_name"
            >
              <div slot-scope="scope">
                <span>{{
                  scope.row.store_name ? scope.row.store_name : "--"
                }}</span>
              </div>
            </el-table-column>
            <el-table-column
              v-if="list_data_type == 2"
              label="导购所在部门"
              min-width="140"
            >
              <template slot-scope="scope">
                <div>
                  {{ scope.row.dp_name ? scope.row.dp_name : "--" }}
                </div>
              </template>
            </el-table-column>
            <!-- 订单数据 -->
            <el-table-column
              v-if="active_name == 1"
              label="推广订单数"
              min-width="140"
              prop="extension_num"
              sortable="custom"
              key="1"
            ></el-table-column>
            <el-table-column
              v-if="active_name == 1"
              label="推广订单实收金额（元）"
              min-width="210"
              prop="extension_amount"
              sortable="custom"
              key="2"
            >
              <div
                slot-scope="{ row }"
                v-if="active_name == 1"
              >
                {{ row.extension_amount / 100 }}
              </div>
            </el-table-column>
            <el-table-column
              label="服务订单数"
              min-width="150"
              prop="service_num"
              sortable="custom"
              key="3"
              v-if="active_name == 1 && this.USER.version !== 3"
            ></el-table-column>
            <el-table-column
              label="服务订单实收金额（元）"
              min-width="210"
              prop="service_amount"
              sortable="custom"
              key="4"
              v-if="active_name == 1 && this.USER.version !== 3"
            >
              <div
                slot-scope="{ row }"
                v-if="active_name == 1"
              >
                {{ row.service_amount / 100 }}
              </div>
            </el-table-column>
            <!-- 顾客数据 -->
            <el-table-column
              v-if="active_name == 2"
              label="顾客总数"
              min-width="140"
              prop="customer_total"
              sortable="custom"
            >
            </el-table-column>
            <el-table-column
              v-if="active_name == 2"
              label="会员总数"
              min-width="120"
              prop="member_total"
              sortable="custom"
            >
            </el-table-column>
            <el-table-column
              v-if="active_name == 2"
              label="非会员总数"
              min-width="120"
              prop="non_member_total"
              sortable="custom"
            >
            </el-table-column>
            <el-table-column
              v-if="active_name == 2"
              label="导购名片添加总数"
              min-width="150"
              prop="guide_customer_total"
              sortable="custom"
              width="180"
            ></el-table-column>
            <el-table-column
              v-if="active_name == 2"
              label="其他方式添加总数"
              min-width="150"
              prop="other_customer_total"
              sortable="custom"
              width="180"
            ></el-table-column>
            <el-table-column
              v-if="active_name == 2"
              label="新增顾客数"
              min-width="150"
              prop="customer_new"
              sortable="custom"
            >
            </el-table-column>
            <el-table-column
              v-if="active_name == 2"
              label="新增会员数"
              min-width="150"
              prop="member_new"
              sortable="custom"
            >
            </el-table-column>
            <el-table-column
              v-if="active_name == 2"
              label="新增非会员数"
              min-width="150"
              prop="non_member_new"
              sortable="custom"
            >
            </el-table-column>
            <el-table-column
              v-if="active_name == 2"
              label="导购名片新加数"
              min-width="150"
              prop="guide_customer_new"
              sortable="custom"
              width="180"
            ></el-table-column>
            <el-table-column
              v-if="active_name == 2"
              label="其他方式新加数"
              min-width="150"
              prop="other_customer_new"
              sortable="custom"
              width="180"
            ></el-table-column>
            <el-table-column
              v-if="active_name == 2"
              label="发起会话数"
              min-width="140"
              prop="chat_total"
              sortable="custom"
              width="180"
            ></el-table-column>
            <el-table-column
              v-if="active_name == 2"
              label="咨询应答率"
              min-width="150"
              prop="reply_percentage"
              sortable="custom"
            >
              <template slot-scope="scope">
                {{ (scope.row.reply_percentage * 100).toFixed(2) + "%" }}
              </template>
            </el-table-column>
            <!-- <el-table-column
          v-if="active_name == 2"
          label="绑定顾客总数"
          min-width="150"
          prop="binding_customer_total"
          sortable="custom"
        >
        </el-table-column>
        <el-table-column
          width="125"
          v-if="active_name == 2"
          label="新增绑客数"
          min-width="150"
          prop="binding_customer_new"
          sortable="custom"
        >
        </el-table-column> -->
            <!-- 分销数据 -->
            <el-table-column
              v-if="active_name == 3"
              label="分销结算笔数"
              min-width="140"
              prop="settlement_num"
              sortable="custom"
              key="11"
            ></el-table-column>
            <el-table-column
              v-if="active_name == 3"
              label="分销结算金额（元）"
              min-width="190"
              sortable="custom"
              key="12"
            >
              <div
                slot-scope="{ row }"
                v-if="active_name == 3"
              >
                {{ row.settlement_money / 100 }}
              </div>
            </el-table-column>
            <div slot="empty">
              <layout-empty
                v-if="searchLoading"
                empty_type="tag"
                :message="'没有找到相关数据'"
              ></layout-empty>
              <layout-empty
                v-else
                :message="'暂无数据'"
              ></layout-empty>
            </div>
          </el-table>
          <!-- 选择导购 -->
          <select-member
            v-if="member_show"
            :show.sync="member_show"
            :value="member_show"
            @cancel="member_show = false"
            @confirm="submitDialog"
          ></select-member>
        </div>
      </template>
      <template slot="footer">
        <el-pagination
          v-if="table_list && table_list.length !== 0"
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
    </layout-list>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Api from "@/api/wisdom/dataCompass";
import UserModule from "@/store/modules/user";
import HelpStore from "@/store/modules/help";
import LayoutEmpty from "@/components/sy-components/layout-empty.vue";
// import SelectMember from "@/components/sy-components/select-member.vue";
import SelectMember from "@/views/myapp/wisdom/communication/components/select-member.vue";
import SearchSelectMember from "@/components/sy-components/search-select-member.vue";
import ShareModule from "@/shared/index";
@Component({
  components: {
    SelectMember,
    SearchSelectMember,
    LayoutEmpty,
  },
})
export default class extends Vue {
  // data
  get HELPSHOW() {
    return HelpStore.getHelp;
  }
  private times: any = "";
  private active_name: any = "2";
  private time_type: any = 1;
  private hign_choose: Array<any> = [];
  private guide_workers: Array<any> = [];
  private searchLoading = false; // 是否正在搜索
  private time_list = [
    { name: "今日", label: 0 },
    { name: "昨日", label: 3 },
    { name: "近7天", label: 1 },
    { name: "近30天", label: 2 },
  ];
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10, // 每页数量
  };
  private summary_data_customer = [
    {
      name: "顾客总数",
      desc: "导购名下顾客中，累计由该导购发展的顾客总数累计到当前时间的顾客总数",
      key: "data1",
      value: 0,
    },
    {
      name: "会员总数",
      desc: "会员是指已在DL系统注册的用户",
      key: "data13",
      value: 0,
    },
    {
      name: "非会员总数",
      desc: "非会员是指未在DL系统注册、已在我们系统注册的用户",
      key: "data15",
      value: 0,
    },
    {
      name: "导购名片添加总数",
      desc: "累计通过导购名片添加的外部联系人数量",
      key: "data5",
      value: 0,
    },
    {
      name: "其他方式添加总数",
      desc: "累计通过非导购名片的其他方式（例如原生二维码或其他联系我二维码）添加的外部联系人数量",
      key: "data6",
      value: 0,
    },
  ];
  private summary_data_customer_add = [
    {
      name: "新增顾客数",
      desc: "筛选时间范围内由该导购新发展的顾客数量",
      key: "data2",
      value: 0,
    },
    {
      name: "新增会员数",
      desc: "对应统计时段和导购/门店下的新增会员数",
      key: "data14",
      value: 0,
    },
    {
      name: "新增非会员数",
      desc: "对应统计时段和导购/门店下的新增非会员数",
      key: "data16",
      value: 0,
    },
    {
      name: "导购名片新加数",
      desc: "通过导购名片添加的外部联系人新增数量",
      key: "data7",
      value: 0,
    },
    {
      name: "其他方式新加数",
      desc: "通过非导购名片的其他方式（例如原生二维码或其他联系我二维码）添加的外部联系人新增数量",
      key: "data8",
      value: 0,
    },
  ];
  private summary_data_customer_third = [
    // {
    //   name: "绑定顾客总数",
    //   desc: "统计时段内导购名下的绑定顾客总数，该数据不受时间筛选影响",
    //   key: "data11",
    //   value: 0
    // },
    // {
    //   name: "新增绑客数",
    //   desc: "统计时段内导购名下新增的绑定业绩关系的顾客数量",
    //   key: "data12",
    //   value: 0
    // },
    // {
    //   name: "",
    //   desc: "",
    //   key: "",
    //   value: 0
    // },
    {
      name: "发起会话数",
      desc: "统计时效T+1，统计时段内导购主动发起的会话总数",
      key: "data9",
      value: 0,
    },
    {
      name: "咨询应答率",
      desc: "统计时效T+1，统计时段内顾客主动发起聊天后，指定导购在一个自然日内，回复聊天数/客户发起聊天数的比率平均值",
      key: "data10",
      value: 0,
    },
    {
      name: "",
      desc: "",
      key: "",
      value: 0,
    },
    {
      name: "",
      desc: "",
      key: "",
      value: 0,
    },
    {
      name: "",
      desc: "",
      key: "",
      value: 0,
    },
  ];
  private list_data_type = "1";
  private order_type = "";
  private order_desc = "";
  private summary_data = {};
  private table_list = []; // 列表数据
  private member_show = false;

  // 订单图表数据
  private order_data = {
    item_name: [
      "推广订单数",
      "推广订单实收金额（元）",
      "服务订单数",
      "服务订单实收金额（元）",
    ],
    yAxis: [{ type: "value" }, { type: "value", splitLine: { show: false } }],
    series: [
      { name: "推广订单数", type: "line", smooth: true, data: [] },
      {
        name: "推广订单实收金额（元）",
        type: "line",
        smooth: true,
        data: [],
        yAxisIndex: 1,
      },
      { name: "服务订单数", type: "line", smooth: true, data: [] },
      {
        name: "服务订单实收金额（元）",
        type: "line",
        smooth: true,
        data: [],
        yAxisIndex: 1,
      },
    ],
  };

  // 顾客图表数据
  private customer_data = {
    item_name: [
      "顾客总数",
      "新增顾客数",
      "导购名片添加总数",
      "其他方式添加总数",
      "导购名片新加数",
      "其他方式新加数",
      "发起会话数",
      "咨询应答率（%）",
      // "绑定顾客总数",
      // "新增绑客数",
      "会员总数",
      "非会员总数",
      "新增会员数",
      "新增非会员数",
    ],
    yAxis: [
      { type: "value" },
      {
        type: "value",
        min: 0,
        max: 100,
      },
    ],
    series: [
      {
        name: "顾客总数",
        type: "line",
        smooth: true,
        data: [],
        yAxisIndex: 0,
      },
      {
        name: "新增顾客数",
        type: "line",
        smooth: true,
        data: [],
        yAxisIndex: 0,
      },
      {
        name: "导购名片添加总数",
        type: "line",
        smooth: true,
        data: [],
        yAxisIndex: 0,
      },
      {
        name: "其他方式添加总数",
        type: "line",
        smooth: true,
        data: [],
        yAxisIndex: 0,
      },
      {
        name: "导购名片新加数",
        type: "line",
        smooth: true,
        data: [],
        yAxisIndex: 0,
      },
      {
        name: "其他方式新加数",
        type: "line",
        smooth: true,
        data: [],
        yAxisIndex: 0,
      },
      {
        name: "发起会话数",
        type: "line",
        smooth: true,
        data: [],
        yAxisIndex: 0,
      },
      {
        name: "咨询应答率（%）",
        type: "line",
        smooth: true,
        data: [],
        yAxisIndex: 1,
      },
      // {
      //   name: "绑定顾客总数",
      //   type: "line",
      //   smooth: true,
      //   data: [],
      //   yAxisIndex: 0
      // },
      // {
      //   name: "新增绑客数",
      //   type: "line",
      //   smooth: true,
      //   data: [],
      //   yAxisIndex: 0
      // },
      {
        name: "会员总数",
        type: "line",
        smooth: true,
        data: [],
        yAxisIndex: 0,
      },
      {
        name: "非会员总数",
        type: "line",
        smooth: true,
        data: [],
        yAxisIndex: 0,
      },
      {
        name: "新增会员数",
        type: "line",
        smooth: true,
        data: [],
        yAxisIndex: 0,
      },
      {
        name: "新增非会员数",
        type: "line",
        smooth: true,
        data: [],
        yAxisIndex: 0,
      },
    ],
  };

  // 分销图表数据
  private task_data = {
    item_name: ["分销结算笔数", "分销结算金额（元）"],
    yAxis: [{ type: "value" }, { type: "value", splitLine: { show: false } }],
    series: [
      { name: "分销结算笔数", type: "line", smooth: true, data: [] },
      {
        name: "分销结算金额（元）",
        type: "line",
        smooth: true,
        data: [],
        yAxisIndex: 1,
      },
    ],
  };

  // 折现图表数据
  private e_data: {
    color: Array<string>;
    title: any;
    tooltip: any;
    legend: any;
    grid: any;
    xAxis: any;
    yAxis: any;
    series: any;
  } = {
    color: [
      "#5B89FE",
      "#E16565",
      "#7EC5DC",
      "#A8BAEA",
      "#EC913E",
      "#E8C10A",
      "#0D9884",
      "#BDAF62",
      // "#7E62C2",
      // "#A2C780",
      "#101010",
      "#6CB13A",
      "#3A54B1",
      "#B13A80",
    ],
    title: { text: "" },
    tooltip: { trigger: "axis" },
    legend: {
      left: "1.5%",
      icon: "circle",
      data: [],
    },
    grid: {
      left: "2%",
      right: "2%",
      bottom: "0%",
      top: "15%",
      containLabel: true,
    },
    xAxis: {
      type: "category",
      boundaryGap: false,
      data: [],
      axisLabel: {
        // interval: 0,
        // fontSize: 10
        // inside: true,
        // formatter: function (value) {
        //     return value.splice(0,4)
        // }
        // rotate: 30
      },
    },
    yAxis: [],
    series: [],
  };

  // 左侧饼图图表数据
  private p_data: {
    legend: any;
    color: Array<string>;
    series: any;
  } = {
    legend: {
      top: "bottom",
    },
    color: ["#71D4D4", "#A3D2f6", "#CEDD79", "#FACFA9", "#C3B6E6"],
    series: [
      {
        name: "面积模式",
        type: "pie",
        radius: [0, 120],
        center: ["50%", "50%"],
        roseType: "area",
        itemStyle: {
          borderColor: "#fff",
          borderWidth: 1,
        },
        label: {
          formatter: "{name|{b}:}\n {number|{d}%,} {number|共{c}人}",
          minMargin: 5,
          edgeDistance: 2,
          lineHeight: 15,
          rich: {
            number: {
              fontSize: 10,
              color: "#333",
              fontWeight: 600,
            },
          },
        },
        data: [],
      },
    ],
  };

  // 右侧饼图图表数据
  private pice_data: {
    legend: any;
    color: Array<string>;
    series: any;
  } = {
    legend: {
      top: "bottom",
    },
    color: ["#71D4D4", "#A3D2f6", "#CEDD79", "#FACFA9", "#C3B6E6"],
    series: [
      {
        name: "面积模式",
        type: "pie",
        radius: [0, 120],
        center: ["50%", "50%"],
        roseType: "area",
        itemStyle: {
          borderColor: "#fff",
          borderWidth: 1,
        },
        label: {
          formatter: "{name|{b}:}\n {number|{d}%,} {number|共{c}人}",
          minMargin: 5,
          edgeDistance: 2,
          lineHeight: 15,
          rich: {
            number: {
              fontSize: 10,
              color: "#333",
              fontWeight: 600,
            },
          },
        },
        data: [],
      },
    ],
  };

  // 计算属性
  get USER() {
    return UserModule.getUserInfo;
  }
  // 切换日期/导购tab栏
  private changeListType() {
    this.meta.page = 1;
    this.order_type = "";
    this.order_desc = "";
    this.table_list = [];
    (this.$refs["dataList"] as any).clearSort();
    this.getList();
  }

  // 隐藏选择员工下拉选项
  private hideDropDownMenu(is_show: boolean) {
    if (is_show) {
      const dom = document.getElementsByClassName("is-multiple")[0];
      (dom as any).style.opacity = 0;
      this.member_show = true;
    }
  }

  // 员工输入框即高级搜索
  private highSearch() {
    this.member_show = true;
  }
  // 高级搜索完成
  private submitDialog(list: any) {
    this.hign_choose = list.map(item => {
      return {
        id: item.value,
        name: item.label,
      }
    });
    this.guide_workers = list.map((v: any) => v.value);
    this.fetchList();
    this.getList();
  }

  // 删除已选导购
  private removeGuide(val: any) {
    const idx = this.hign_choose.findIndex((item) => item.id === val);
    idx > -1 && this.hign_choose.splice(idx, 1);
    this.guide_workers = this.hign_choose.map((v: any) => v.id);
    this.fetchList();
    this.getList();
  }

  // 搜索
  private search() {
    if (!this.times) {
      this.$message.info("请先选择时间");
      return;
    }
    this.meta.page = 1;
    this.fetchList();
    this.getList();
    (this.$refs["dataList"] as any).clearSort();
    this.searchLoading = true;
  }
  // 重置
  private reset() {
    this.time_type = 1;
    this.searchLoading = false;
    (this.$refs["dataList"] as any).clearSort();
    this.times = this.getDates(this.time_type);
    this.hign_choose = [];
    this.guide_workers = [];
    this.fetchList();
    this.getList();
  }
  // 选择自定义时间，清空时间类型
  private timeChange() {
    this.time_type = "";
  }
  // 切换时间类型，清空自定义时间
  private changeTime(data: any) {
    this.times = this.getDates(data);
  }
  //获取时间
  private getDates(type: any) {
    const getDate = (n: number, m: number) => {
      const date = new Date();
      const year = date.getFullYear();
      const month = date.getMonth() + 1;
      const day = date.getDate();
      const nowDay = new Date(`${year}/${month}/${day}`).getTime();
      return [nowDay - 86400000 * n, nowDay - 1000 - 86400000 * m];
    };

    const arr = [
      () => {
        return getDate(0, -1);
      },
      () => {
        return getDate(6, -1);
      },
      () => {
        return getDate(29, -1);
      },
      () => {
        return getDate(1, 0);
      },
    ];
    return arr[type]();
  }
  // 排序
  private sortChange(column: any) {
    this.order_type = column.prop;
    // ascending 升序 descending 降序
    this.order_desc =
      column.order === "ascending"
        ? "asc"
        : column.order === "descending"
        ? "desc"
        : "";
    this.meta.page = 1;
    this.getList();
  }
  // 每页
  private handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.getList();
  }
  // 当前页
  private handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getList();
  }
  // 导出表格
  private exportExcel() {
    const obj = {
      "1": "orderStatisticalExport",
      "2": "guideCustomerExport",
      "3": "taskStatisticalEXPORT",
    };
    const params = {
      start_time: this.times[0] / 1000,
      end_time: this.times[1] / 1000,
      worker_ids: this.guide_workers,
      export_type: Number(this.list_data_type),
      sort_type: this.order_type,
      statistical_sort: this.order_desc,
      version: this.USER.version,
    };
    Api[obj[this.active_name]](params).then(() => {
      this.$confirm(
        "文件已进入下载列表，请前往【下载中心】查看及下载文件。",
        "文件下载提示",
        {
          distinguishCancelAndClose: true,
          confirmButtonText: "前往下载中心",
          cancelButtonText: "取消",
          type: "success",
        }
      ).then(() => {
        ShareModule.getShared().pushState("myapp_download-center");
      });
    });
  }
  //导出订单|绑客明细
  private exportOrderDetailExcel() {
    const obj = {
      "1": "orderDetailExport",
      "2": "customerDetailExport",
    };
    let params;
    if (this.active_name == 1) {
      params = {
        worker_ids: this.guide_workers,
        start_time: this.times[0] / 1000,
        end_time: this.times[1] / 1000,
        version: this.USER.version,
      };
    } else {
      params = {
        worker_ids: this.guide_workers,
      };
    }
    Api[obj[this.active_name]](params).then(() => {
      this.$confirm(
        "文件已进入下载列表，请前往【下载中心】查看及下载文件。",
        "文件下载提示",
        {
          distinguishCancelAndClose: true,
          confirmButtonText: "前往下载中心",
          cancelButtonText: "取消",
          type: "success",
        }
      ).then(() => {
        ShareModule.getShared().pushState("myapp_download-center");
      });
    });
  }

  // 初始图表数据
  private initChart() {
    //如果是顾客数据，头部距离增加
    if (this.active_name == 2) {
      this.e_data.grid.top = "25%";
    } else {
      this.e_data.grid.top = "15%";
    }
    const obj = {
      "1": "order_data",
      "2": "customer_data",
      "3": "task_data",
    };
    const arr = (name: any) => {
      this.e_data.series = this.$data[name].series;
      this.e_data.legend.data = this.$data[name].item_name;
      this.e_data.yAxis = this.$data[name].yAxis;
    };
    arr(obj[this.active_name]);
  }

  // 获取图表数据
  private async fetchList() {
    const obj = {
      "1": "orderStatistical",
      "2": "statisticGuideCustomer",
      "3": "taskStatisticalChart",
    };
    const params = {
      worker_ids: this.guide_workers,
      version: this.USER.version,
      start_time: this.times[0] / 1000,
      end_time: this.times[1] / 1000,
    };

    await Api[obj[this.active_name]](params).then((res: any) => {
      const { detail_data } = res.data;
      const y_data = detail_data.y;
      // 饼图数据
      this.p_data.series[0].data = res.data.member_level_statistics;
      this.pice_data.series[0].data =
        res.data.member_level_statistics_by_time_range;
      // v3.2.21兼容顾客数据
      if (this.active_name == 1) {
        this.summary_data = {
          data1: res.data.data1,
          data2: res.data.data2,
          data3: res.data.data3,
          data4: res.data.data4,
          data5: res.data.data5,
          data6: res.data.data6,
        };
      } else if (this.active_name == 3) {
        this.summary_data = {
          data1: res.data.data1,
          data2: res.data.data2,
        };
      } else {
        this.summary_data_customer.forEach((ele: any) => {
          ele.value = res.data[ele.key] || 0;
        });
        this.summary_data_customer_add.forEach((ele: any) => {
          ele.value = res.data[ele.key] || 0;
        });
        this.summary_data_customer_third.forEach((ele: any) => {
          ele.value = res.data[ele.key] || 0;
          if (ele.key === "data10") {
            ele.value = (res.data[ele.key] * 100).toFixed(2) + "%" || "0.00%";
          }
        });
      }

      this.e_data.xAxis.data = detail_data.x;
      if (this.active_name == 1) {
        y_data.data2 = y_data.data2.map((v: number) => Number(v) / 100);
        y_data.data4 = y_data.data4.map((v: number) => Number(v) / 100);
      }
      if (this.active_name == 3) {
        y_data.data2 = y_data.data2.map((v: number) => Number(v) / 100);
      }
      const arr: any = Object.values(y_data);
      //订单数据与分销数据按照原先逻辑走
      if (this.active_name == 1 || this.active_name == 3) {
        arr.forEach((item: any, index: any) => {
          if (this.e_data.series[index]) {
            this.e_data.series[index].data = item;
          }
        });
      } else {
        this.e_data.series[0].data = arr[0];
        this.e_data.series[1].data = arr[1];
        this.e_data.series[2].data = arr[4];
        this.e_data.series[3].data = arr[5];
        this.e_data.series[4].data = arr[6];
        this.e_data.series[5].data = arr[7];
        this.e_data.series[6].data = arr[8];
        arr[9].forEach((item: any, index: any) => {
          arr[9][index] = (item * 100).toFixed(2);
        });
        this.e_data.series[7].data = arr[9];
        // this.e_data.series[8].data = arr[10];
        // this.e_data.series[9].data = arr[11];
        this.e_data.series[8].data = arr[12];
        this.e_data.series[9].data = arr[13];
        this.e_data.series[10].data = arr[14];
        this.e_data.series[11].data = arr[15];
      }
    });
    this.draw();
  }

  // 获取列表数据
  private getList() {
    let apis = "";
    const type: any = this.list_data_type;
    if (this.active_name == 1) {
      apis = type == 1 ? "orderStatisticalDate" : "orderStatisticalGuide";
    }
    if (this.active_name == 2) {
      apis = type == 1 ? "guideCustomerDate" : "guideCustomerList";
    }
    if (this.active_name == 3) {
      apis = type == 1 ? "taskStatisticalDate" : "taskStatisticalWorker";
    }
    const params = {
      worker_ids: this.guide_workers,
      version: this.USER.version,
      start_time: this.times[0] / 1000,
      end_time: this.times[1] / 1000,
      page: this.meta.page,
      size: this.meta.size,
      sort_type: this.order_type,
      statistical_sort: this.order_desc,
    };
    Api[apis](params).then((res: any) => {
      const { list, meta } = res.data;
      if (meta) {
        this.meta = meta;
      } else {
        const { page, size, total } = res.data;
        this.meta = { page, size, total };
      }
      this.table_list = list;
    });
  }

  // 绘制图表
  private draw() {
    const is_exist = document.getElementById("diagram-box");
    if (is_exist == null) return;
    // 折现图
    const my_chart = this.$echarts.init(document.getElementById("diagram-box"));
    // 左侧饼图
    const pie_chart = this.$echarts.init(
      document.getElementById("pieChart-box")
    );
    // 右侧饼图
    const pie_chart1 = this.$echarts.init(
      document.getElementById("pieChart1-box")
    );
    // 中心商城无门店企业不可见-服务订单和服务订单支付金额
    if (
      this.active_name == 1 &&
      this.USER.version == 3 &&
      this.e_data.series.length > 2
    ) {
      this.e_data.series = this.e_data.series.slice(0, 2);
    }
    my_chart.setOption(this.e_data, true);
    pie_chart.setOption(this.p_data, true);
    pie_chart1.setOption(this.pice_data, true);
    window.onresize = function () {
      my_chart.resize();
      pie_chart.resize();
      pie_chart1.resize();
    };
  }

  // 生命周期
  private mounted() {
    this.times = this.getDates(this.time_type);
    this.initChart();
    this.fetchList();
    this.getList();
  }
}
</script>
<style lang="scss" scoped>
.customer-data {
  color: #303133;
}
/deep/.search-item {
  display: inline-flex;
  align-items: center;
  margin-top: 0;
  margin-bottom: 20px;
}
/deep/.el-tabs__active-bar {
  width: 60px !important;
}
.search-select-member {
  width: 280px;
  display: flex;
  display: inline-block;
  vertical-align: middle;
  align-items: center;
  overflow: hidden;
  border: 1px solid $color-second-line;
  background-color: $color-white;
  border-radius: 0 !important;
  /deep/.el-input__inner {
    border: none;
    border-radius: 0;
    color: $color-text-second;
  }

  .select-collapse {
    /deep/.el-input__suffix {
      display: none;
    }
  }
}

.search-info /deep/ .el-radio-button .el-radio-button__inner {
  height: 36px;
  line-height: 15px;
  display: inline-block;
  vertical-align: middle;
}

.search-info {
  font-size: 14px;
  color: #606266;
  &.m-t-2 {
    margin-top: -2px;
  }
}

.search-data:last-child {
  border-radius: 0 4px 4px 0;
}

/deep/ .el-row {
  margin: 20px 0;
}

.search-info /deep/ .el-range-editor--small.el-input__inner {
  height: 36px;
}

/deep/ .hover-span:hover {
  color: #404040;
  border-color: #404040;
}

.help-icon {
  margin-left: 5px;
  cursor: pointer;
}

/deep/ .el-table .caret-wrapper {
  position: absolute;
  top: -6px;
}

.customer-data /deep/ .search-content {
  padding: 20px;
  margin-bottom: 20px;
}

.customer-data /deep/ .buttons-content {
  display: none;
}

.customer-data /deep/ .operate-content {
  display: none;
}

.customer-data /deep/ .slide-box {
  display: none;
}

.customer-data /deep/ .table-content {
  display: none;
}

.page-right {
  margin-top: 20px;
}

.table-text-img {
  border-radius: 100%;
  vertical-align: middle;
  margin-right: 4px;
}

.employee {
  font-size: 14px;
  height: 32px;
  border: 1px solid #dcdfe6;
}

.operate-btn {
  .reset {
    color: #303133;
    background-color: #fff;
    border-color: #c7cad0;
    margin-left: 15px;
  }
}

.guide-info {
  font-size: 12px;
  margin-top: 10px;
}

.export {
  margin-bottom: 20px;
}

.pickTab {
  border: 1px solid #c7cad0;
}

.order-title {
  font-size: 14px;
  font-weight: 600;
  margin-top: 20px;
}

.info-tilte {
  font-size: 14px;
  padding: 20px 25px;
  border-radius: 5px;
  background: #f0f2f6;
}

.info-tilte:nth-child(2),
.info-tilte:nth-child(3) {
  padding: 20px 20px;
}

.info-tilte span {
  display: inline-block;
  margin-bottom: 15px;
}

.flex-s-c {
  width: 100%;
  margin-left: 0;
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-start;
  align-items: center;
}

.flex-s-c-auto {
  display: flex;
  justify-content: flex-start;
  align-items: center;
  .info-title {
    font-size: 14px;
    padding: 20px 25px;
    background: #f0f2f6;
    span {
      display: inline-block;
      margin-bottom: 15px;
    }
  }
  .info-title:nth-child(1) {
    border-top-left-radius: 5px;
    border-bottom-left-radius: 5px;
  }
  .info-title:nth-child(3) {
    border-top-right-radius: 5px;
    border-bottom-right-radius: 5px;
  }
}

.info-card {
  margin-bottom: 20px;
  border-radius: 0;
  width: 20%;
  flex: 1;
}

.info-card:nth-child(4),
.info-card:nth-child(5),
.info-card:nth-child(8) {
  margin-left: 20px;
  border-radius: 5px;
}

.info-card:nth-child(1),
.info-card:nth-child(5) {
  border-top-left-radius: 5px !important;
  border-bottom-left-radius: 5px !important;
}

.info-card:nth-child(3),
.info-card:nth-child(7) {
  border-top-right-radius: 5px !important;
  border-bottom-right-radius: 5px !important;
}

.sum {
  font-size: 36px;
  font-weight: 600;
}

.serach-tab {
  background: #ffffff;
  padding: 10px 20px 20px 20px;
  margin-top: 20px;
}

.echartBox {
  width: 100%;
  height: 240px;
  margin-bottom: 20px;
}
.pieChartBox {
  flex: 1;
  height: 440px;
  margin-bottom: 20px;
  #pieChart-box {
    width: 100%;
    height: 100%;
  }
}

.noData {
  min-height: 330px;
  margin-bottom: 20px;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #333;
}

.listTitle {
  position: relative;
  margin-top: 25px;
}

.hover-span {
  margin-left: 3px;
  cursor: pointer;
}
.no-data {
  background: transparent;
}

.pie-title {
  width: 90%;
  color: $color-main;
  border-bottom: 1px solid;
  margin-top: 30px;
  padding-bottom: 10px;
}
.customize :nth-of-type(2) {
  border-top-right-radius: 5px !important;
  border-bottom-right-radius: 5px !important;
}
.pte-container {
  display: flex;
  .left-pie {
    flex: 1;
  }
  .right-pie {
    flex: 1;
  }
}
/deep/ .layout-list {
  padding: 10px 0 0 0 !important;
}
</style>
