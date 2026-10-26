<template>
  <div>
    <layout-list>
      <template slot="search">
        <search-item text="数据时间：">
          <cy-select
            v-model="search.dataTime"
            placeholder="请选择数据时间"
            @change="timeChange"
            :optionsList="search.dataTimeList"
          >
          </cy-select>
        </search-item>
        <search-item text="门店选择：">
          <el-select
            style="width:280px"
            v-model="search.department.name"
            ref="selectBlur"
          >
            <el-option
              :value="search.department.value"
              :label="search.department.name"
              style="height: auto;padding: 0px;"
            >
              <el-tree
                ref="tree"
                class="treeStyle"
                v-if="search.department.isShowTree"
                :props="search.department.data"
                :load="loadNode"
                node-key="id"
                :expand-on-click-node="true"
                :check-strictly="true"
                :check-on-click-node="true"
                @node-click="handleNodeClick"
                lazy
              ></el-tree>
            </el-option>
          </el-select>
        </search-item>
        <search-item text="导购选择：">
          <cy-select
            v-model="search.staff.value"
            search-input
            :optionsList="searchStaffList"
            :options="{ label: 'mem_name', value: 'id' }"
          >
          </cy-select>
        </search-item>
        <search-item text="活动时间：">
          <el-date-picker
            style="height:36px"
            v-model="search.activityTime"
            @change="timeChooseChange"
            type="daterange"
            format="yyyy-MM-dd"
            value-format="timestamp"
            range-separator="-"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            :default-time="['00:00:00', '23:59:59']"
            :picker-options="pickerOptions"
          ></el-date-picker>
        </search-item>
        <el-button
          type="primary"
          class="font-14"
          @click="clickSearch()"
          size="small"
        >查询</el-button>
        <el-button @click="reset">重置</el-button>
      </template>
    </layout-list>
    <layout-card class="margin-top-30 card-box-shadow">
      <template slot="title">
        <div>社群数据</div>
      </template>
      <template slot="content">
        <div class="privateGuide">
          <div class="card">
            <div class="title">
              群总数
              <el-popover
                trigger="hover"
                :content="'当前查询条件下，有效的群（不包括删除的群）累计总数'"
              >
                <svg-icon
                  style="margin-left: 5px;"
                  icon-class="help-funnel"
                  slot="reference"
                ></svg-icon>
              </el-popover>
            </div>
            <div class="price">{{ communityData.group_count || 0 }}</div>
            <svg-icon
              class="icon-tip"
              icon-class="stay-fan"
            ></svg-icon>
          </div>
          <div class="card">
            <div class="title">
              累计群成员数
              <el-popover
                trigger="hover"
                :content="'当前查询条件下，有效的群（不包括删除的群）累计群成员总数*群成员不去重'
                  "
              >
                <svg-icon
                  style="margin-left: 5px;"
                  icon-class="help-funnel"
                  slot="reference"
                ></svg-icon>
              </el-popover>
            </div>
            <div class="price">{{ communityData.group_member_count || 0 }}</div>
            <svg-icon
              class="icon-tip"
              icon-class="stay-fan"
            ></svg-icon>
          </div>
          <div class="card">
            <div class="title">
              新入群人数
              <el-popover
                trigger="hover"
                :content="'当前查询条件下，在查询时间内入群（有效群）的用户的总数'
                  "
              >
                <svg-icon
                  style="margin-left: 5px;"
                  icon-class="help-funnel"
                  slot="reference"
                ></svg-icon>
              </el-popover>
            </div>
            <div class="price">
              {{ communityData.new_join_group_count || 0 }}
            </div>
            <svg-icon
              class="icon-tip"
              icon-class="new-people"
            ></svg-icon>
          </div>
          <div class="card">
            <div class="title">
              退群人数
              <el-popover
                trigger="hover"
                :content="'当前查询条件下，在查询时间累计退群（有效群）的用户的总数'
                  "
              >
                <svg-icon
                  style="margin-left: 5px;"
                  icon-class="help-funnel"
                  slot="reference"
                ></svg-icon>
              </el-popover>
            </div>
            <div class="price">{{ communityData.leave_group_count || 0 }}</div>
            <svg-icon
              class="icon-tip"
              icon-class="number-of-withdrawals"
            ></svg-icon>
          </div>
          <div class="card">
            <div class="title">
              客户入群率
              <el-popover
                trigger="hover"
                :content="'截止当前查询条件下，入群客户与全部客户占比，即客户入群率=新入群人数/累计群成员数*100%'
                  "
              >
                <svg-icon
                  style="margin-left: 5px;"
                  icon-class="help-funnel"
                  slot="reference"
                ></svg-icon>
              </el-popover>
            </div>
            <div class="price">{{ communityData.join_group_rate || "0%" }}</div>
            <svg-icon
              class="icon-tip"
              icon-class="percentage"
            ></svg-icon>
          </div>
          <div class="card">
            <div class="title">
              客户流失率
              <el-popover
                trigger="hover"
                :content="'截止当前查询条件下，退群客户数与累计群成员数的占比，即累计退群客户数/累计群成员数*100%'
                  "
              >
                <svg-icon
                  style="margin-left: 5px;"
                  icon-class="help-funnel"
                  slot="reference"
                ></svg-icon>
              </el-popover>
            </div>
            <div class="price">
              {{ communityData.leave_group_rate || "0%" }}
            </div>
            <svg-icon
              class="icon-tip"
              icon-class="leave_group_rate"
            ></svg-icon>
          </div>
          <div class="card">
            <div class="title">
              加好友人数
              <el-popover
                trigger="hover"
                :content="'截止当前查询条件下，已添加导购好友且已入群的累计客户数总和'
                  "
              >
                <svg-icon
                  style="margin-left: 5px;"
                  icon-class="help-funnel"
                  slot="reference"
                ></svg-icon>
              </el-popover>
            </div>
            <div class="price">
              {{ communityData.join_friend_count || 0 }}
            </div>
            <svg-icon
              class="icon-tip"
              icon-class="add-friend"
            ></svg-icon>
          </div>
          <div class="card">
            <div class="title">
              加好友率
              <el-popover
                trigger="hover"
                :content="'截止当前查询条件下，添加导购好友且已入群的客户数与累计群成员数的占比，即加好友率=加好友人数/去重后累计群成员数*100%'
                  "
              >
                <svg-icon
                  style="margin-left: 5px;"
                  icon-class="help-funnel"
                  slot="reference"
                ></svg-icon>
              </el-popover>
            </div>
            <div class="price">
              {{ communityData.join_friend_rate || "0%" }}
            </div>
            <svg-icon
              class="icon-tip"
              icon-class="percentage"
            ></svg-icon>
          </div>
        </div>
      </template>
    </layout-card>
    <layout-card class="margin-top-30 card-box-shadow">
      <template slot="title">
        <div>社群活动数据</div>
      </template>
      <template slot="content">
        <div class="privateGuide">
          <div class="card">
            <div class="title">
              活动群数
              <el-popover
                trigger="hover"
                :content="'截止当前查询条件下，与裂变活动相关联的群总数'"
              >
                <svg-icon
                  style="margin-left: 5px;"
                  icon-class="help-funnel"
                  slot="reference"
                ></svg-icon>
              </el-popover>
            </div>
            <div class="price">
              {{ communityData.group_activity_count || 0 }}
            </div>
            <svg-icon
              class="icon-tip"
              icon-class="group_activity_count"
            ></svg-icon>
          </div>
          <div class="card">
            <div class="title">
              群活动参与人数
              <el-popover
                trigger="hover"
                :content="'截止当前查询条件下，参与群裂变活动的累计人数总和'"
              >
                <svg-icon
                  style="margin-left: 5px;"
                  icon-class="help-funnel"
                  slot="reference"
                ></svg-icon>
              </el-popover>
            </div>
            <div class="price">
              {{ communityData.join_group_activity_count || 0 }}
            </div>
            <svg-icon
              class="icon-tip"
              icon-class="join_group_activity_count"
            ></svg-icon>
          </div>
          <div class="card">
            <div class="title">
              未参加群活动人数
              <el-popover
                trigger="hover"
                :content="'截止当前查询条件下，已入群但未参与过任何群裂变活动的人数总和'
                  "
              >
                <svg-icon
                  style="margin-left: 5px;"
                  icon-class="help-funnel"
                  slot="reference"
                ></svg-icon>
              </el-popover>
            </div>
            <div class="price">
              {{ communityData.not_join_group_activity_count || 0 }}
            </div>
            <svg-icon
              class="icon-tip"
              icon-class="stay-fan"
            ></svg-icon>
          </div>
          <div class="card">
            <div class="title">
              未参加群活动率
              <el-popover
                trigger="hover"
                :content="'截止当前查询条件下，未参加群活动人数与累计群成员数占比，即未参加群活动人数/去重后的累计群成员数*100%'
                  "
              >
                <svg-icon
                  style="margin-left: 5px;"
                  icon-class="help-funnel"
                  slot="reference"
                ></svg-icon>
              </el-popover>
            </div>
            <div class="price">
              {{ communityData.not_join_group_activity_rate || "0%" }}
            </div>
            <svg-icon
              class="icon-tip"
              icon-class="not_join_group_activity_rate"
            ></svg-icon>
          </div>
        </div>
      </template>
    </layout-card>
    <layout-card class="margin-top-30 card-box-shadow">
      <template slot="title">
        <div>社区数据明细</div>
      </template>
      <template slot="tabs">
        <el-tabs
          v-model="activTab"
          @tab-click="handleTabClick(activTab)"
        >
          <el-tab-pane
            label="按社群"
            name="1"
          > </el-tab-pane>
          <el-tab-pane
            label="按导购"
            name="2"
          > </el-tab-pane>
        </el-tabs>
      </template>
      <template slot="content">
        <el-radio-group
          v-model="shoppingGuideDate"
          @change="refreshTable"
        >
          <el-radio-button :label="private_guide_date.today">今日</el-radio-button>
          <el-radio-button :label="private_guide_date.yesterday">昨日</el-radio-button>
          <el-radio-button :label="private_guide_date.week">近7天</el-radio-button>
          <el-radio-button :label="private_guide_date.month">近30天</el-radio-button>
        </el-radio-group>
        <el-date-picker
          class="m-l-20 choice-picker"
          v-model="timeChoose"
          type="datetimerange"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
          size="medium"
          prefix-icon="el-icon-date"
          :clearable="false"
          :default-time="['00:00:00', '23:59:59']"
          @change="timeListChange"
          value-format="timestamp"
        >
        </el-date-picker>
        <el-button
          size="small"
          class="excel-button"
          :disabled="list.length <= 0"
          @click="statisticsCountExport"
        >导出Excel</el-button>
        <el-table
          v-if="handleTab === 1"
          ref="shopperGuideTable"
          style="width: 100%;margin: 20px 0 0 0;"
          class="border_btm_none"
          :data="list"
          height="500"
          stripe
        >
          <el-table-column
            label="社群名称"
            prop="name"
            min-width="150"
          ></el-table-column>
          <el-table-column
            label="累计群成员数"
            prop="group_member_count"
            min-width="150"
          ></el-table-column>
          <el-table-column
            label="新入群人数"
            prop="new_join_group_count"
            min-width="150"
          >
          </el-table-column>
          <el-table-column
            label="退群人数"
            prop="leave_group_count"
            min-width="100"
          ></el-table-column>
          <el-table-column
            label="客户入群率"
            prop="join_group_rate"
            min-width="110"
          ></el-table-column>
          <el-table-column
            label="客户流失率"
            prop="leave_group_rate"
            min-width="110"
          >
          </el-table-column>
          <el-table-column
            label="群活动数"
            prop="group_activity_count"
            min-width="100"
          >
          </el-table-column>
          <el-table-column
            label="加好友数"
            prop="join_friend_count"
            min-width="100"
          >
          </el-table-column>
          <el-table-column
            label="加好友率"
            prop="join_friend_rate"
            min-width="100"
          >
          </el-table-column>
          <el-table-column
            label="群活动参与人数"
            prop="join_group_activity_count"
            min-width="150"
          >
          </el-table-column>
          <el-table-column
            label="未参加群活动人数"
            prop="not_join_group_activity_count"
            min-width="160"
          >
          </el-table-column>
          <el-table-column
            label="未参加群活动率"
            prop="not_join_group_activity_rate"
            min-width="160"
          >
          </el-table-column>
          <div
            slot="empty"
            class="table-empty-label"
          >
            <img
              class="empty-img"
              src="@/assets/images/empty-box.png"
              alt=""
            />
            <div class="desc">没有找到相关数据</div>
          </div>
        </el-table>
        <el-table
          v-if="handleTab === 2"
          ref="shopperGuideTable"
          style="width: 100%;margin: 20px 0 0 0;"
          class="border_btm_none"
          :data="list"
          height="500"
          stripe
        >
          <el-table-column
            label="导购ID"
            prop="worker_no"
            min-width="100"
          ></el-table-column>
          <el-table-column
            label="导购"
            min-width="160"
          >
            <div
              v-if="handleTab === 2"
              slot-scope="scope"
            >
              <img
                :src="scope.row.avatar
                  ? scope.row.avatar
                  : 'https://oss-crs.vchangyi.com/wechat/default-avatar-icon.png'
                  "
                alt=""
                width="30"
                height="30"
                class="table-text-img"
              /><span v-text="scope.row.nickname"></span>
            </div>
          </el-table-column>
          <el-table-column
            label="导购所在门店"
            prop="department_name"
            min-width="150"
          ></el-table-column>
          <el-table-column
            label="累计群成员数"
            prop="group_member_count"
            min-width="150"
          ></el-table-column>
          <el-table-column
            label="新入群人数"
            prop="new_join_group_count"
            min-width="150"
          ></el-table-column>
          <el-table-column
            label="退群人数"
            prop="leave_group_count"
            min-width="100"
          ></el-table-column>
          <el-table-column
            label="客户入群率"
            prop="join_group_rate"
            min-width="150"
          ></el-table-column>
          <el-table-column
            label="客户流失率"
            prop="leave_group_rate"
            min-width="150"
          ></el-table-column>
          <el-table-column
            label="加好友数"
            prop="join_friend_count"
            min-width="100"
          ></el-table-column>
          <el-table-column
            label="加好友率"
            prop="join_friend_rate"
            min-width="100"
          ></el-table-column>
          <el-table-column
            label="群活动参与人数"
            prop="group_activity_count"
            min-width="150"
          ></el-table-column>
          <el-table-column
            label="未参加群活动人数"
            prop="not_join_group_activity_count"
            min-width="160"
          ></el-table-column>
          <el-table-column
            label="未参加群活动率"
            prop="not_join_group_activity_rate"
            min-width="150"
          ></el-table-column>
          <div
            slot="empty"
            class="table-empty-label"
          >
            <img
              class="empty-img"
              src="@/assets/images/empty-box.png"
              alt=""
            />
            <div class="desc">没有找到相关数据</div>
          </div>
        </el-table>
      </template>
    </layout-card>
    <layout-list class="footer-meta">
      <template slot="footer">
        <el-pagination
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page="meta.page"
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
import LayoutCard from "@/views/myapp/home/components/layout-card/index.vue";
import API from "@/api/enterprise/data-compass";
import storeAPI from "@/api/enterprise/enterprise";
import Layout from "@/layout/layout.vue";
@Component({
  components: {
    LayoutCard,
    Layout
  }
})
export default class extends Vue {
  // data
  private search: any = {
    dataTime: "1",
    dataTimeList: [
      { label: "今天", value: "1" },
      { label: "昨天", value: "2" },
      { label: "近7天", value: "3" },
      { label: "近30天", value: "4" }
    ],
    // 部门
    department: {
      value: 0, // 部门id
      name: "",
      isShowTree: false,
      data: {
        label: "name",
        children: [],
        isLeaf: "leaf"
      }
    },
    // 员工
    staff: {
      value: "全部" // 员工id
    },
    activityTime: 0
  };
  // 部门列表
  private departmentList: Array<any> = [];
  // 员工列表
  private searchStaffList: Array<any> = [];
  private pickerMinDate = 0; //第一次选中的时间
  private pickerOptions = {
    onPick: (obj: any) => {
      this.pickerMinDate = new Date(obj.minDate).getTime();
    },
    disabledDate: (time: any) => {
      if (this.pickerMinDate) {
        const day1 = 30 * 24 * 3600 * 1000;
        const maxTime = this.pickerMinDate + day1;
        const minTime = this.pickerMinDate - day1;
        return time.getTime() >= maxTime || time.getTime() <= minTime;
      }
    }
  };
  private communityData = {};
  // 是否正在搜索
  private searchLoading = false;
  private shoppingGuideDate = 0;
  private activTab = "1";
  private handleTab = 1;
  private timeChoose: Array<number> | number = 0; // 表格切换的时间
  private private_guide_date = {
    today: 0,
    yesterday: 1,
    week: 2,
    month: 3
  };
  // 表格数据
  private list: Array<any> = [];
  // 分页配置
  private meta: any = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  // 生命周期
  private mounted() {
    this.search.activityTime = this.timeChoose = this.getDateRange(
      new Date(),
      0
    );
    this.search.department.name = "全部";
    // 获取门店
    this.fetchStoreList();
    // 门店下的导购
    this.searchMemberList();
    this.loadData();
    this.getList(this.handleTab);
  }

  //动态加载节点
  private loadNode(node: any, resolve: any) {
    // 一级节点处理
    if (node.level === 0) {
      this.setChecked(this.departmentList);
      return resolve(this.departmentList);
    }
    // 其余节点处理
    if (node.level >= 1) {
      this.getIndex(node, resolve);
    }
  }
  //判断部门是否有子节点
  private setChecked(arr: any) {
    arr.forEach((item: any) => {
      if (!item.is_child_department) item.leaf = true;
    });
  }
  // 加载叶子节点数据函数
  private getIndex(node: any, resolve: any) {
    console.log("node", node);
    storeAPI
      .storeDataApi({
        dp_id: node.data.id
      })
      .then((res: any) => {
        if (res.code == 200) {
          const data = res.data.list;
          data.forEach((item: any) => {
            item.id = item.dp_id;
            item.name = item.dp_name;
          });
          if (data.length <= 0) {
            return;
          }
          this.setChecked(data);
          resolve(data);
        } else {
          this.$message.error(res.message);
        }
      });
  }
  // 选择部门
  private handleNodeClick(node: any) {
    this.search.department.name = node.name;
    this.search.department.value = node.id;
    (this.$refs.selectBlur as any).blur();
    this.search.staff.value = "全部";
    this.searchStaffList = [];
    this.searchMemberList();
  }
  // 获取部门
  fetchStoreList() {
    storeAPI.storeDataApi({ pid: 0 }).then((res: any) => {
      if (res.code == 200) {
        const obj = { id: 0, name: "全部" };
        res.data.list.forEach((item: any) => {
          item.id = item.dp_id;
          item.name = item.dp_name;
        });
        this.departmentList = res.data.list;
        this.departmentList.unshift(obj);
        this.search.department.isShowTree = true;
      } else {
        this.$message.error(res.message);
        return;
      }
    });
  }
  // 搜索员工列表
  private searchMemberList() {
    const params = {
      page: 1,
      size: 999,
      search: {
        dp_id: this.search.department.value
      }
    };
    storeAPI.shoppingGuideListApi(params).then((res: any) => {
      if (res.code == 200) {
        const obj: any = { id: 0, mem_name: "全部" };
        this.searchStaffList = res.data.list || [];
        this.searchStaffList.unshift(obj);
      } else {
        this.$message.error(res.message);
      }
    });
  }

  // 时间改变
  private timeChange(value: string) {
    const date = new Date();
    switch (value) {
      case "1":
        this.search.activityTime = this.getDateRange(date, 0);
        break;
      case "2":
        this.search.activityTime = this.getDateRange(date, 1);
        break;
      case "3":
        this.search.activityTime = this.getDateRange(date, 6);
        break;
      case "4":
        this.search.activityTime = this.getDateRange(date, 29);
    }
  }

  // 计算时间差
  private getDateRange(dateNow: any, intervalDays: number) {
    let computeTime: Array<number> = [];
    const oneDayTime = 24 * 3600 * 1000; // 毫秒数
    // 当前时间戳减去一天的毫秒数(6近一周的时间，30近一个月的时间)
    // 过去时间
    const lastDay =
      new Date(new Date().toLocaleDateString()).getTime() -
      intervalDays * oneDayTime;
    // 当前时间
    const nowDays =
      new Date(dateNow.toLocaleDateString()).getTime() + 86400000 * 1 - 1000;
    // const lastDay = new Date(dateNow.getTime() - intervalDays * oneDayTime);
    if (intervalDays === 1) {
      const lastDays =
        new Date(dateNow.toLocaleDateString()).getTime() + 86400000 * 0 - 1000;
      computeTime = [lastDay, lastDays];
    } else {
      computeTime = [lastDay, nowDays];
    }
    return computeTime;
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

  //日期时间选择器
  private timeChooseChange() {
    this.search.dataTime = "";
  }

  // 查询
  private clickSearch() {
    this.searchLoading = true;
    this.loadData();
  }

  // 重置
  private reset() {
    this.search.dataTime = "1";
    this.search.department.value = 0;
    this.search.department.name = "全部";
    this.search.staff.value = "全部";
    (this.$refs.tree as any).setCurrentKey(null);
    // 获取部门
    this.fetchStoreList();
    // 部门下的员工
    this.searchMemberList();
    this.timeChange("1");
    this.clickSearch();
  }

  // 社群数据
  private loadData() {
    const params = {
      start_time:
        (this.search.activityTime && this.search.activityTime[0] / 1000) ||
        undefined,
      end_time:
        (this.search.activityTime && this.search.activityTime[1] / 1000) ||
        undefined,
      department_ids: this.search.department.value
        ? [this.search.department.value]
        : [],
      worker_ids:
        this.search.staff.value !== "全部" ? [this.search.staff.value] : []
    };
    API.getCommunityIndexApi(params).then((res: any) => {
      if (res.code === 200) {
        this.communityData = res.data;
      } else {
        this.$message.error(res.message);
      }
    });
  }

  // 切换tabs
  handleTabClick(type: any) {
    this.meta.page = 1;
    this.handleTab = Number(type);
    this.shoppingGuideDate = 0;
    this.refreshTable(0);
    this.getList(this.handleTab);
  }

  // 选择自定义时间，清空时间类型
  private timeListChange() {
    this.shoppingGuideDate = -1;
    this.getList(this.handleTab);
  }

  // 表格数据
  getList(type: number) {
    const params = {
      start_time: (this.timeChoose && this.timeChoose[0] / 1000) || undefined,
      end_time: (this.timeChoose && this.timeChoose[1] / 1000) || undefined,
      page: this.meta.page,
      size: this.meta.size
    };
    let api = "";
    if (type === 1) {
      api = "getCommunityListApi";
    } else {
      api = "getCommunityListByMemberApi";
    }
    API[api](params).then((res: any) => {
      if (res.code === 200) {
        this.list = res.data.list || [];
        this.meta.total = res.data.total;
        this.meta.page = res.data.page;
      } else {
        this.$message.error(res.message);
      }
    });
  }

  // 分页
  handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.getList(this.handleTab);
  }
  handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getList(this.handleTab);
  }

  // 切换表格数据
  refreshTable(value: number) {
    const date = new Date();
    switch (value) {
      case 0:
        this.timeChoose = this.getDateRange(date, 0);
        break;
      case 1:
        this.timeChoose = this.getDateRange(date, 1);
        break;
      case 2:
        this.timeChoose = this.getDateRange(date, 6);
        break;
      case 3:
        this.timeChoose = this.getDateRange(date, 29);
    }
    this.getList(this.handleTab);
  }

  // 导出表格
  statisticsCountExport() {
    const params = {
      type: this.handleTab,
      start_time: (this.timeChoose && this.timeChoose[0] / 1000) || undefined,
      end_time: (this.timeChoose && this.timeChoose[1] / 1000) || undefined
    };
    API.exportCommunityApi(params).then((res: any) => {
      const blob = new Blob([res], {
        type: "application/x-xls;charset=utf-8"
      });
      let fileName = "";
      const nowTime = new Date();
      fileName = "社区数据明细_" + this.formateDate(nowTime) + ".xlsx";

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
}
</script>
<style lang="scss" scoped>
/deep/ .layout-list {
  padding-bottom: 1px;

  .search-content {
    background-color: #fff;
    padding: 0;
    margin-bottom: 0;
  }
}

.margin-top-30 {
  margin-top: 30px;
}

.card-box-shadow {
  box-shadow: 0px 0px 20px 1px rgba(0, 0, 0, 0.08);
}

.privateGuide {
  display: flex;
  justify-content: flex-start;
  align-items: center;
  margin-top: 20px;
  flex-wrap: wrap;

  .card {
    border: 1px solid #000000;
    width: 30%;
    min-width: 30%;
    max-width: 30%;
    margin-right: 40px;
    margin-top: 30px;
    position: relative;

    div {
      line-height: 1.5;
    }

    padding: 24px 30px;
    box-sizing: border-box;

    .title {
      font-size: 16px;
      color: #606266;
    }

    .price {
      color: #606266;
      font-weight: bolder;
      font-size: 26px;
      position: relative;
      top: 12px;
    }

    .desc {
      margin-top: 10px;
      color: #606266;
      font-size: 14px;
    }

    .icon-tip {
      width: 92px;
      height: 92px;
      position: absolute;
      bottom: 0;
      right: 0;
    }
  }
}

.excel-button {
  margin-left: 30px;
}

.footer-meta {
  margin-top: -60px;
  padding-bottom: 20px;
}

.table-text-img {
  border-radius: 100%;
  vertical-align: middle;
  margin-right: 4px;
}

/deep/ .el-table__body-wrapper {
  height: 460px !important;
}
</style>
