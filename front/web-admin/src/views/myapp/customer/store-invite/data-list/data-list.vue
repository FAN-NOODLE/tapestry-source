<template>
  <div class="data-list">
    <layout-list>
      <template slot="search">
        <search-item text="创建时间： ">
          <el-date-picker
            v-model="search.time"
            type="datetimerange"
            prefix-icon="el-icon-date"
            format="yyyy-MM-dd HH:mm:ss"
            value-format="timestamp"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            :default-time="['00:00:00', '23:59:59']"
            style="height:36px"
          >
          </el-date-picker>
        </search-item>

        <search-item text="订单时间： ">
          <el-date-picker
            v-model="search.reserve_order_time"
            type="datetimerange"
            prefix-icon="el-icon-date"
            format="yyyy-MM-dd HH:mm:ss"
            value-format="timestamp"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            :default-time="['00:00:00', '23:59:59']"
            style="height:36px"
          >
          </el-date-picker>
        </search-item>
        <search-item v-if="epNameDetail != 'coach'" text="评价时间： ">
          <el-date-picker
            v-model="search.comment_time"
            type="datetimerange"
            prefix-icon="el-icon-date"
            format="yyyy-MM-dd HH:mm:ss"
            value-format="timestamp"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            :default-time="['00:00:00', '23:59:59']"
            style="height:36px"
          >
          </el-date-picker>
        </search-item>

        <search-item text="邀约门店： ">
          <SelectStore ref="selectStore"></SelectStore>
        </search-item>

        <search-item text="邀约导购： ">
          <SelectStaff ref="selectStaff"></SelectStaff>
        </search-item>

        <search-item text="邀约状态：">
          <el-select
            style="width:250px"
            class="select"
            v-model="search.status"
            size="medium"
          >
            <el-option
              v-for="item in statusList"
              :key="item.key"
              :label="item.label"
              :value="item.key"
            ></el-option>
          </el-select>
        </search-item>
        <search-item text="是否成单：">
          <el-select
            style="width:250px"
            class="select"
            v-model="search.is_pay"
            size="medium"
          >
            <el-option
              v-for="item in ifOrNotLists"
              :key="item.key"
              :label="item.label"
              :value="item.key"
            ></el-option>
          </el-select>
        </search-item>
        <search-item text="是否到店购买：">
          <el-select
            style="width:250px"
            class="select"
            v-model="search.is_store"
            size="medium"
          >
            <el-option
              v-for="item in isStoreOptions"
              :key="item.key"
              :label="item.label"
              :value="item.key"
            ></el-option>
          </el-select>
        </search-item>
        <search-item v-if="epNameDetail != 'coach'" text="是否评价：">
          <el-select
            style="width:250px"
            class="select"
            v-model="search.is_comment"
            size="medium"
          >
            <el-option
              v-for="item in ifOrNotLists"
              :key="item.key"
              :label="item.label"
              :value="item.key"
            ></el-option>
          </el-select>
        </search-item>
        <search-item text="邀约单号：">
          <cy-input
            size="small"
            width-size="small"
            placeholder="请输入邀约单号"
            v-model.trim="search.reserve_number"
            prefix-icon="el-icon-search"
            clearable
          ></cy-input>
        </search-item>
        <search-item text="邀约来源：">
          <el-select
            style="width:250px"
            class="select"
            v-model="search.order_origin"
            size="medium"
          >
            <el-option
              v-for="item in orderOriginList"
              :key="item.key"
              :label="item.label"
              :value="item.key"
            ></el-option>
          </el-select>
        </search-item>

        <search-item>
          <el-button type="primary" @click="goSearch" size="small"
            >查询</el-button
          >
          <el-button plain @click="reset" size="small">重置</el-button>
          <el-button v-if="epNameDetail == 'coach'" plain @click="exportList" size="small">导出列表</el-button>
          <el-button v-else plain @click="exportList" size="small">导出</el-button>
          <el-button v-if="epNameDetail == 'coach'" plain @click="exportOrderList" size="small">导出订单明细</el-button>
        </search-item>
      </template>
      <template slot="tabs" v-if="nums">
        <div class="data-content">
          <el-card class="box-card">
            <div class="title">
              <el-tooltip
                effect="light"
                :content="'导购仅创建邀约次数及创建并发送的邀约次数总和'"
                placement="top-end"
              >
                <span
                  >邀约次数
                  <svg-icon
                    slot="reference"
                    class="help-icon help-icon-m"
                    icon-class="group-help"
                  ></svg-icon>
                </span>
              </el-tooltip>
            </div>
            <div class="content">{{ nums.invite_num }}</div>
          </el-card>

          <el-card class="box-card">
            <div class="title">
              <el-tooltip
                effect="light"
                :content="'邀约状态为【邀约成功】的单数；'"
                placement="top-end"
              >
                <span
                  >邀约成功数
                  <svg-icon
                    slot="reference"
                    class="help-icon help-icon-m"
                    icon-class="group-help"
                  ></svg-icon>
                </span>
              </el-tooltip>
            </div>
            <div class="content">{{ nums.invite_success_num }}</div>
            <div class="rag">
              <el-tooltip
                effect="light"
                placement="top-end"
                v-if="ep.isCoach"
              >
                <span
                  >到店购买
                  <svg-icon
                    slot="reference"
                    class="help-icon help-icon-m"
                    icon-class="group-help"
                  ></svg-icon>
                </span>
                <template #content>
                  <div>
                    <div>邀约成功的订单中，到店购买订单数的统计包含以下两种状态</div>
                    <div>状态0：导购未选择（默认状态）</div>
                    <div>状态1：导购选择了“到店购买”</div>
                  </div>
                </template>
              </el-tooltip>
              <el-tooltip
                effect="light"
                :content="'邀约成功的订单中，到店购买的订单数'"
                placement="top-end"
                v-else
              >
                <span
                  >到店购买
                  <svg-icon
                    slot="reference"
                    class="help-icon help-icon-m"
                    icon-class="group-help"
                  ></svg-icon>
                </span>

              </el-tooltip>
              {{ nums.pay_success_num }}
            </div>
          </el-card>
          <el-card class="box-card">
            <div class="title">
              <el-tooltip
                effect="light"
                :content="'筛选条件下全部成单金额汇总'"
                placement="top-end"
              >
                <span
                  >导购成单金额
                  <svg-icon
                    slot="reference"
                    class="help-icon help-icon-m"
                    icon-class="group-help"
                  ></svg-icon>
                </span>
              </el-tooltip>
            </div>
            <div class="content">{{ nums.success_price / 100 }}</div>
          </el-card>

          <el-card class="box-card">
            <div class="title">
              <el-tooltip
                effect="light"
                :content="'客户打开导购分享邀约单的人数；'"
                placement="top-end"
              >
                <span
                  >小程序打开数
                  <svg-icon
                    slot="reference"
                    class="help-icon help-icon-m"
                    icon-class="group-help"
                  ></svg-icon>
                </span>
              </el-tooltip>
            </div>
            <div class="content">{{ nums.open_num }}</div>
            <div class="rag">
              <el-tooltip
                effect="light"
                :content="'筛选条件下，客户打开小程序的人数在邀约用户的占比'"
                placement="top-end"
              >
                <span
                  >占比
                  <svg-icon
                    slot="reference"
                    class="help-icon help-icon-m"
                    icon-class="group-help"
                  ></svg-icon>
                </span>
              </el-tooltip>
              {{ nums.open_per }}
            </div>
          </el-card>
          <el-card class="box-card">
            <div class="title">
              <el-tooltip
                effect="light"
                :content="'客户确认邀约单的人数；'"
                placement="top-end"
              >
                <span
                  >顾客确认数
                  <svg-icon
                    slot="reference"
                    class="help-icon help-icon-m"
                    icon-class="group-help"
                  ></svg-icon>
                </span>
              </el-tooltip>
            </div>
            <div class="content">{{ nums.commit_num }}</div>
            <div class="rag">
              <el-tooltip
                effect="light"
                :content="'筛选条件下，客户确认邀约单的用户在邀约用户的占比'"
                placement="top-end"
              >
                <span
                  >占比
                  <svg-icon
                    slot="reference"
                    class="help-icon help-icon-m"
                    icon-class="group-help"
                  ></svg-icon>
                </span>
              </el-tooltip>
              {{ nums.commit_per }}
            </div>
          </el-card>
        </div>
      </template>
      <template slot="table">
        <el-table
          tooltip-effect="light"
          :data="list"
          ref="multipleTable"
          @selection-change="selectChange"
        >
          <!-- <el-table-column type="selection" width="55"></el-table-column> -->
          <el-table-column label="编号" prop="content" min-width="180">
            <template slot-scope="scope">
              {{ scope.row.reserve_number }}
            </template>
          </el-table-column>
          <el-table-column label="顾客昵称" prop="content" min-width="150">
            <template slot-scope="scope">
              {{ scope.row.user_name }}
            </template>
          </el-table-column>
          <el-table-column label="邀约导购" prop="content" min-width="150">
            <template slot-scope="scope">
              {{ scope.row.member_name }}
            </template>
          </el-table-column>
          <el-table-column label="邀约门店" prop="content" min-width="140">
            <template slot-scope="scope"> {{ scope.row.store_name }} </template>
          </el-table-column>

          <el-table-column label="邀请来源" key="order_origin" min-width="140">
            <template slot-scope="scope">
              {{ orderOriginList[scope.row.order_origin].label || "--" }}
            </template>
          </el-table-column>

          <el-table-column label="邀约次数" prop="content" min-width="140">
            <template slot="header">
              <span class="tag-header">邀约次数</span>
              <el-popover placement="top-start" trigger="hover">
                <div>
                  0代表导购仅创建了邀约，未发送邀约；1代表导购创建并发送了邀约单
                </div>
                <svg-icon
                  slot="reference"
                  class="help-icon help-icon-m"
                  icon-class="group-help"
                ></svg-icon>
              </el-popover>
            </template>
            <template slot-scope="scope"> {{ scope.row.invite_num }} </template>
          </el-table-column>

          <el-table-column align="center" min-width="140px" v-if="epNameDetail">
            <template slot="header">
              <span class="tag-header">邀约状态</span>
              <el-popover placement="top-start" trigger="hover">
                <div v-if="epNameDetail === 'sw'">
                  <b>邀约生效中：</b>用户未打开邀约链接或用户已打开链接，但距邀约单上的到店时间未超过30天<br />
                  <b>邀约匹配失败：</b>到店时间超过30天，且未匹配到订单，状态变更为已过期<br />
                  <b>已取消：</b>SA/SM手动取消邀约单<br />
                  <b>邀约完成：</b>用户到店，且导购手动确认用户到店，该阶段可邀请客户评价<br />
                  <b>邀约成功：</b>用户到店且有关联订单，若用户在邀约完成阶段未评价，SA在该阶段仍可发送邀请评价链接<br />
                </div>
                <div v-if="epNameDetail==='coach'">
                  <b>进行中：</b>用户未打开邀约链接或用户已打开链接，但距邀约单发出时间（不含当天）未超过30天<br />
                  <b>已过期：</b>邀约发出时间超过30天，且未匹配到订单<br />
                  <b>已取消：</b>SA/SM手动取消邀约单<br />
                  <b>邀约完成：</b>用户到店，且导购手动确认用户到店<br />
                  <b>邀约成功：</b>邀约单发出（不含当天）30天内有关联成功订单<br />
                </div>
                <div v-if="epNameDetail==='ks'">
                  <b>进行中：</b>用户未打开邀约链接或用户已打开链接，但距邀约单上的到店时间未超过30天<br />
                  <b>已过期：</b>到店时间超过30天，且未匹配到订单，状态变更为已过期<br />
                  <b>已取消：</b>SA/SM手动取消邀约单<br />
                  <b>邀约完成：</b>用户到店，且导购手动确认用户到店，该阶段可邀请客户评价<br />
                  <b>邀约成功：</b>用户到店且有关联订单，若用户在邀约完成阶段未评价，SA在该阶段仍可发送邀请评价链接<br />
                </div>
                <svg-icon
                  slot="reference"
                  class="help-icon help-icon-m"
                  icon-class="group-help"
                ></svg-icon>
              </el-popover>
            </template>
            <template slot-scope="scope">
              <span>{{ statusMap[scope.row.status] || "-" }}</span>
            </template>
          </el-table-column>

          <el-table-column label="签到状态" prop="content" min-width="140" v-if="epNameDetail != 'coach'">
            <template slot-scope="scope">
              <span v-if="scope.row.order_origin == 2">--</span>
              <span v-else>{{
                scope.row.is_sign === 1 ? "已签到" : "未签到"
              }}</span>
            </template>
          </el-table-column>

          <el-table-column label="是否打开" prop="content" min-width="140">
            <template slot-scope="scope">
              {{ scope.row.is_open > 0 ? "是" : "否" }}
            </template>
          </el-table-column>
          <el-table-column label="是否确认" prop="content" min-width="140">
            <template slot-scope="scope">
              {{ scope.row.is_commit > 0 ? "是" : "否" }}
            </template>
          </el-table-column>

          <el-table-column
            label="顾客类型"
            key="customer_status"
            min-width="100"
          >
            <template slot-scope="scope">
              {{ customerStatusList[scope.row.customer_status].label || "--" }}
            </template>
          </el-table-column>
          <el-table-column label="是否成单" min-width="100" key="is_pay">
            <template slot-scope="scope">
              {{ ifOrNotLists[scope.row.is_pay].label || "--" }}
            </template>
          </el-table-column>
          <el-table-column
            label="是否到店购买"
            key="is_store"
            min-width="170"
            prop="is_store_str"
          >
            <template slot-scope="scope">
              {{ ["--", "是", "否"][scope.row.is_store] }}
            </template>
          </el-table-column>
          <el-table-column
            v-if="epNameDetail != 'coach'"
            key="reserve_order_price"
            label="导购成单金额"
            min-width="140"
          >
            <template slot-scope="scope">
              ￥{{ scope.row.reserve_order_price / 100 }}
            </template>
          </el-table-column>
          <el-table-column
            v-if="epNameDetail == 'coach'"
            label="订单总金额"
            min-width="170"
            key="reserve_order_price"
            prop="reserve_order_price"
          >
            <template slot-scope="scope">
              ￥{{ scope.row.reserve_order_price / 100 }}
            </template>
          </el-table-column>
          <el-table-column
            v-if="epNameDetail == 'coach'"
            label="关联订单数"
            min-width="170"
            key="relation_order_num"
            prop="relation_order_num"
          >
            <template slot-scope="scope">
              {{ scope.row.relation_order_num||'--' }}
            </template>
          </el-table-column>
          <el-table-column
            label="最近订单时间"
            key="reserve_order_time_str"
            min-width="180"
          >
            <template slot-scope="scope">
              {{ scope.row.reserve_order_time_str }}</template
            >
          </el-table-column>
          <el-table-column label="是否评价" key="is_comment" min-width="170"  v-if="epNameDetail != 'coach'">
            <template slot-scope="scope">
              {{ ifOrNotLists[scope.row.is_comment].label || "--" }}
            </template>
          </el-table-column>
          <el-table-column
            label="评价时间"
            key="comment_time_str"
            min-width="190"
            v-if="epNameDetail != 'coach'"
          >
            <template slot-scope="scope">
              {{ scope.row.comment_time_str }}
            </template>
          </el-table-column>

          <el-table-column label="创建时间" prop="content" min-width="180">
            <template slot-scope="scope">
              {{
                scope.row.created_at | formatDate("yyyy-mm-dd hh:mm:ss")
              }}</template
            >
          </el-table-column>

          <div slot="empty" v-if="!emptyFlag" style="margin-top: 50px;">
            <img class="empty-img" src="@/assets/images/empty-box.png" alt="" />
            <div class="empty-box">
              <div class="desc">暂无数据</div>
            </div>
          </div>
          <div slot="empty" v-if="emptyFlag" style="margin-top: 50px;">
            <img
              class="empty-img"
              src="@/assets/images/empty-search.png"
              alt=""
            />
            <div class="empty-box">
              <div class="desc">没有找到相关数据</div>
            </div>
          </div>
        </el-table>
      </template>
      <template slot="footer" v-if="list.length > 0">
        <el-pagination
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page.sync="meta.page"
          :page-sizes="[10, 20, 30, 40]"
          :page-size="meta.size"
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
import SelectStore from "../component/select-store.vue";
import SelectStaff from "../component/select-staff.vue";
import Api from "@/api/store-invite/index";
import PerModule from "@/store/modules/permission";
import {
  customerStatusList,
  orderOriginList,
  ifOrNotLists,
  statusList,
  statusMap,
  swStatusList,
  swStatusMap,
  isStoreOptions
} from "../contants";
@Component({
  components: {
    SelectStore,
    SelectStaff
  }
})
export default class extends Vue {
  private search = {
    time: [],
    reserve_order_time: [],
    comment_time: [],
    status: -1,
    is_pay: 0,
    is_store: -1,
    is_comment: 0,
    order_origin: 0,
    reserve_number: ''
  };
  private meta = {
    page: 1,
    size: 10,
    total: 0
  };
  private list = [];
  private epNameDetail = "";
  private nums = {
    commit_num: 0,
    open_num: 0,
    invite_num: 0
  };
  private multipleSelection = [];
  private customerStatusList = customerStatusList;
  private orderOriginList = orderOriginList;
  private ifOrNotLists = ifOrNotLists;
  private isStoreOptions = isStoreOptions;
  private statusList = [];
  private statusMap = {};
  private emptyFlag = false;
  private ids = [];
  get ep() {
    return PerModule.enterprise
  }
  // 生命周期
  private created() {
    console.log("created");
    Api.getEpConfig().then(res => {
      // // 企业名称：coach;ks;sw
      const { name } = res.data;
      this.epNameDetail = name;

      // todo:邀约状态判断
      if (name === "sw") {
        this.statusList = swStatusList;
        this.statusMap = swStatusMap;
      } else {
        this.statusList = statusList;
        this.statusMap = statusMap;
      }
    });
  }
  private mounted() {
    this.ids = this.$route.params.ids ? [this.$route.params.ids] : []; // 预约列表传进来的数据
    this.getList();
  }
  private async exportList() {
    let params = {};
    if (this.multipleSelection.length === 0) {
      const { page, size, ...obj } = this.getParams();
      params = obj;
    } else {
      params = { ids: this.multipleSelection };
    }
   // todo: 跳转下载中心
   try {
      const downRes = await Api.dataExportApi(params)
      if(downRes.code ===200){
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
    } catch (error) {
      console.log(error);
    }
  }
  private async exportOrderList() {
    let params = {};
    if (this.multipleSelection.length === 0) {
      const { page, size, ...obj } = this.getParams();
      params = obj;
    } else {
      params = { ids: this.multipleSelection };
    }
   // todo: 跳转下载中心
   try {
      const downRes = await Api.dataOrderExportApi(params)
      if(downRes.code ===200){
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
    } catch (error) {
      console.log(error);
    }
  }
  private getPer(val, num) {
    if (val === 0 || num === 0) {
      return "0%";
    }
    return `${((val / num) * 100).toFixed(2)}%`;
  }
  private goSearch() {
    this.meta.page = 1;
    this.ids = [];
    this.getList();
  }
  // 分页
  private handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.getList();
  }
  private handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getList();
  }
  private getParams() {
    const selectStore: any = this.$refs.selectStore;
    const selectStaff: any = this.$refs.selectStaff;
    const params = {
      store_ids: selectStore.list.map(item => {
        return item.dp_third_id;
      }),
      member_ids: selectStaff.list.map(item => {
        return item.member_id;
      }),
      page: this.meta.page,
      size: this.meta.size,
      start_time: this.search.time ? (this.search.time[0] || 0) / 1000 : 0,
      end_time: this.search.time ? (this.search.time[1] || 0) / 1000 : 0,

      //TODO
      reserve_order_start_time: this.search.reserve_order_time
        ? (this.search.reserve_order_time[0] || 0) / 1000
        : 0,
      reserve_order_end_time: this.search.reserve_order_time
        ? (this.search.reserve_order_time[1] || 0) / 1000
        : 0,

      comment_start_time: this.search.comment_time
        ? (this.search.comment_time[0] || 0) / 1000
        : 0,
      comment_end_time: this.search.comment_time
        ? (this.search.comment_time[1] || 0) / 1000
        : 0,
      status: this.search.status,
      is_store: this.search.is_store,
      is_comment: this.search.is_comment,
      is_pay: this.search.is_pay,
      order_origin: this.search.order_origin,
      ids: this.ids,
      reserve_number: this.search.reserve_number
    };
    return params;
  }
  // 当选择项发生变化时会触发该事件，数组
  selectChange(val: any) {
    this.multipleSelection = [];
    val.map((item: any) => {
      this.multipleSelection.push(item.id);
    });
  }
  private reset() {
    console.log("重置");
    this.search = {
      time: [],
      reserve_order_time: [],
      comment_time: [],
      status: -1,
      is_pay: 0,
      is_store: -1,
      is_comment: 0,
      order_origin: 0,
      reserve_number: ''
    };
    this.ids = [];
    (this.$refs.selectStore as any).clear();
    (this.$refs.selectStaff as any).clear();
    this.meta.page = 1;
    this.getList();
  }
  private getList() {
    const params = this.getParams();
    Api.staticDataApi(params).then(res => {
      const { total, list, nums } = res.data;
      this.list = list;
      this.nums = nums;
      this.meta.total = total;
    });
  }
  private destroyed() {
    console.log("destroyed");
  }
}
</script>
<style lang="scss" scoped>
.data-content {
  padding: 20px;
  margin: 20px 0;
  background-color: $color-background;
  border-radius: 4px;
  display: flex;
  justify-content: flex-start;
  .box-card {
    margin-right: 20px;
    width: 200px;
    position: relative;
    /deep/.el-card__body {
      padding: 10px 20px 20px;
    }
    .title {
      line-height: 30px;
      font-size: 16px;
      text-align: center;
    }
    .content {
      height: 60px;
      line-height: 60px;
      text-align: center;
      font-size: 20px;
      font-weight: bold;
    }
    .rag {
      position: absolute;
      right: 10px;
      bottom: 10px;
    }
  }
}
.help-icon-m {
  margin: -3px 4px 0 1px;
}
</style>
