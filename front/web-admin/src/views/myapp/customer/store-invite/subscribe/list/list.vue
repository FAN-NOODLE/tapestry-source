<template>
  <div class="subscribe-list">
    <layout-list>
      <template slot="search">
        <search-item text="创建时间： ">
          <el-date-picker
            v-model="search.created_time"
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
        <search-item text="到店时间： ">
          <el-date-picker
            v-model="search.invite_time"
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
        <search-item text="邀约来源：" key="order_origin">
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
        <search-item :text="`邀约门店： `">
          <SelectStore ref="selectStore"></SelectStore>
        </search-item>
        <search-item :text="`邀约导购： `">
          <SelectStaff ref="selectStaff"></SelectStaff>
        </search-item>
        <search-item text="服务名称： ">
          <el-select
            style="width:250px"
            class="select"
            v-model="search.service_id"
            size="medium"
          >
            <el-option
              v-for="item in serverList"
              :key="item.id"
              :label="item.reserve_name"
              :value="item.id"
            ></el-option>
          </el-select>
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

        <search-item text="是否成单：" key="is_pay">
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
        <search-item text="是否到店购买：" key="is_store">
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
        <search-item text="是否评价：" key="is_comment">
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

        <search-item>
          <el-button type="primary" @click="goSearch" size="small"
            >查询</el-button
          >
          <el-button plain @click="reset" size="small">重置</el-button>
          <el-button plain @click="exportList" size="small">导出</el-button>
          <el-button v-if="epNameDetail == 'coach'" plain @click="exportOrderList" size="small">导出订单明细</el-button>
        </search-item>
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
          <el-table-column :label="`邀约服务`" prop="content" min-width="170">
            <template slot-scope="scope">
              {{ scope.row.reserve_name || "--" }}
            </template>
          </el-table-column>
          <el-table-column label="邀约来源" min-width="100" key="order_origin">
            <template slot-scope="scope">
              {{ orderOriginList[scope.row.order_origin].label || "--" }}
            </template>
          </el-table-column>
          <el-table-column label="创建时间" prop="content" min-width="190">
            <template slot-scope="scope">
              {{ scope.row.create_time_str }}
            </template>
          </el-table-column>
          <el-table-column label="到店时间" prop="content" min-width="190">
            <template slot-scope="scope">
              {{ scope.row.reserve_time_str }}
            </template>
          </el-table-column>
          <el-table-column :label="`邀约门店`" prop="content" min-width="150">
            <template slot-scope="scope">
              {{ scope.row.store_name }}
            </template>
          </el-table-column>
          <el-table-column :label="`邀约导购`" prop="content" min-width="100">
            <template slot-scope="scope">
              {{ scope.row.member_name }}
            </template>
          </el-table-column>
          <el-table-column align="center" min-width="140px">
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

          <el-table-column
            label="签到状态"
            min-width="100"
            v-if="epNameDetail !== 'coach'"
          >
            <template slot-scope="scope">
              <span v-if="scope.row.order_origin == 2">--</span>
              <span v-else>{{
                scope.row.is_sign === 1 ? "已签到" : "未签到"
              }}</span>
            </template>
          </el-table-column>

          <el-table-column label="是否成单" min-width="100" key="is_pay">
            <template slot-scope="scope">
              {{
                scope.row.is_pay ? ifOrNotLists[scope.row.is_pay].label : "--"
              }}
            </template>
          </el-table-column>
          <el-table-column
            label="顾客类型"
            min-width="100"
            key="customer_status"
          >
            <template slot-scope="scope">
              {{ customerStatusList[scope.row.customer_status].label || "--" }}
            </template>
          </el-table-column>
          <el-table-column
            label="是否到店购买"
            min-width="170"
            prop="is_store"
          >
            <template slot-scope="scope">
              {{ ["--", "是", "否"][scope.row.is_store] }}
            </template>
          </el-table-column>
          <el-table-column
            v-if="epNameDetail == 'coach'"
            label="订单总金额"
            min-width="170"
            prop="reserve_order_price"
          >
            <template slot-scope="scope">
              ￥{{ Number(scope.row.reserve_order_price / 100).toLocaleString() }}
            </template>
          </el-table-column>
          <el-table-column
            v-if="epNameDetail == 'coach'"
            label="关联订单数"
            min-width="170"
            prop="relation_order_num"
          >
            <template slot-scope="scope">
              {{ scope.row.relation_order_num||'--' }}
            </template>
          </el-table-column>
          <el-table-column
            v-if="epNameDetail == 'coach'"
            label="最近订单时间"
            key="reserve_order_time_str"
            min-width="180"
          >
            <template slot-scope="scope">
              {{ scope.row.reserve_order_time | formatDate("yyyy-mm-dd hh:mm:ss") }}</template
            >
          </el-table-column>
          <el-table-column 
            v-if="epNameDetail != 'coach'"
            label="是否评价" 
            min-width="100" >
            <template slot-scope="scope">
              {{
                scope.row.is_comment
                  ? ifOrNotLists[scope.row.is_comment].label
                  : "--"
              }}
            </template>
          </el-table-column>
          <el-table-column
            v-if="epNameDetail != 'coach'"
            label="评价时间"
            min-width="190"
          >
            <template slot-scope="scope">
              {{ scope.row.comment_time | formatDate("yyyy-mm-dd hh:mm:ss") }}
            </template>
          </el-table-column>

          <el-table-column
            label="备注"
            min-width="150"
            v-if="epNameDetail === 'sw'"
          >
            <template slot-scope="scope">
              {{ scope.row.reserve_remark }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200" fixed="right">
            <template slot-scope="scope">
              <el-button type="text" size="small" @click="goDetail(scope.row)"
                >查看详情</el-button
              >
              <el-button type="text" size="small" @click="goDataList(scope.row)"
                >查看数据</el-button
              >
            </template>
          </el-table-column>
          <div slot="empty" v-if="!emptyFlag" style="margin-top: 50px;">
            <img class="empty-img" src="@/assets/images/empty-box.png" alt="" />
            <div class="empty-box">
              <div class="desc">暂无邀约单</div>
            </div>
          </div>
          <div slot="empty" v-if="emptyFlag" style="margin-top: 50px;">
            <img
              class="empty-img"
              src="@/assets/images/empty-search.png"
              alt=""
            />
            <div class="empty-box">
              <div class="desc">没有找到相关邀约单</div>
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
    <!-- <subscribe-detail
      v-if="showDetailModel"
      @close="showDetailModel = false"
    ></subscribe-detail> -->
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import SelectStore from "../../component/select-store.vue";
import SelectStaff from "../../component/select-staff.vue";
import Api from "@/api/store-invite/index";
import {
  customerStatusList,
  orderOriginList,
  ifOrNotLists,
  isStoreOptions,
  statusList,
  statusMap,
  swStatusMap,
  swStatusList
} from "../../contants";
@Component({
  components: {
    SelectStore,
    SelectStaff
  }
})
export default class extends Vue {
  private search = {
    created_time: [],
    invite_time: [],
    // todo:成单时间
    reserve_order_time: [],
    service_id: "",
    status: -1,
    is_store: -1,
    is_comment: 0,
    is_pay: 0,
    order_origin: 0,
    reserve_number:''
  };
  private meta = {
    page: 1,
    size: 10,
    total: 5
  };
  private list = [];
  private serverList = [];
  private customerStatusList = customerStatusList;
  private orderOriginList = orderOriginList;
  private ifOrNotLists = ifOrNotLists;
  private isStoreOptions = isStoreOptions;
  private statusList = [];
  private statusMap = {};
  private multipleSelection = [];
  // private showDetailModel = false;
  private emptyFlag = false;
  private epNameDetail = "";

  // 生命周期
  private created() {
    console.log("created");
    this.getEpConfig();
  }
  private mounted() {
    console.log("mounted");
    this.getServerList();
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
      const downRes = await Api.serverExportApi(params)
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
  private getServerList() {
    const params = {
      page: 1,
      size: 999,
      reserve_name: ""
    };
    Api.serviceListApi(params).then(res => {
      this.serverList = res.data.list;
    });
  }
  private goSearch() {
    this.meta.page = 1;
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
  private goDetail(row) {
    this.$router.push({
      name: "myapp_customer_store-invite_subscribe_detail",
      query: {
        id: row.id
      }
    });
  }
  private goDataList(row) {
    this.$router.push({
      name: "myapp_customer_store-invite_data-list",
      params: {
        ids: row.id
      }
    });
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
      start_time: this.search.created_time
        ? (this.search.created_time[0] || 0) / 1000
        : 0,
      end_time: this.search.created_time
        ? (this.search.created_time[1] || 0) / 1000
        : 0,
      reserve_start_time: this.search.invite_time
        ? (this.search.invite_time[0] || 0) / 1000
        : 0,
      reserve_end_time: this.search.invite_time
        ? (this.search.invite_time[1] || 0) / 1000
        : 0,
      // todo:成单时间
      reserve_order_start_time: this.search.reserve_order_time
        ? (this.search.reserve_order_time[0] || 0) / 1000
        : 0,
      reserve_order_end_time: this.search.reserve_order_time
        ? (this.search.reserve_order_time[1] || 0) / 1000
        : 0,
      page: this.meta.page,
      size: this.meta.size,
      service_id: this.search.service_id || 0,
      status: this.search.status,
      is_store: this.search.is_store,
      is_comment: this.search.is_comment,
      is_pay: this.search.is_pay,
      order_origin: this.search.order_origin,
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
    this.search = {
      created_time: [],
      invite_time: [],
      reserve_order_time: [],
      service_id: "",
      status: -1,
      is_store: -1,
      is_comment: 0,
      is_pay: 0,
      order_origin: 0,
      reserve_number:''
    };
    (this.$refs.selectStore as any).clear();
    (this.$refs.selectStaff as any).clear();
    this.meta.page = 1;
    this.getList();
  }
  private getList() {
    const params = this.getParams();
    Api.orderListApi(params).then(res => {
      const { total, list } = res.data;
      this.list = list || [];
      this.meta.total = total;
    });
  }
  // 获取企业
  getEpConfig() {
    Api.getEpConfig().then(res => {
      // // 企业名称：coach;ks;sw
      const { name } = res.data;
      this.epNameDetail = name;
      // 缓存企业名称，给详情页面用
      localStorage.setItem("epNameDetail", name);
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
  private destroyed() {
    console.log("destroyed");
  }
}
</script>
<style lang="scss" scoped>
.help-icon-m {
  margin: -3px 4px 0 1px;
}
</style>
