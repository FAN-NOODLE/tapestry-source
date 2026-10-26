<template>
  <div class="child-task-data">
    <!-- 头部 -->
    <layout-list>
      <template slot="search">
        <div class="content-item m-b-20">
          <span>任务名称：{{ mainTaskData.plan_name || "-" }}</span>
          <span>子任务数：{{ taskObj.task_num || "0" }}</span>
          <span>任务有效期：{{
            mainTaskData.plan_start_time | formatDate("yyyy-mm-dd hh:mm:ss")
          }}～{{
  mainTaskData.plan_end_time | formatDate("yyyy-mm-dd hh:mm:ss")
}}</span>
        </div>
        <div class="content-item m-b-20 flex-start">
          时间：
          <el-date-picker
            v-model="timeRange"
            type="daterange"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            value-format="timestamp"
            :picker-options="pickerOptions"
            :clearable="false"
            :default-time="['00:00:00', '23:59:59']"
            @change="changeSelector()"
          >
          </el-date-picker>
        </div>
        <div class="top-box width-60">
          <div
            v-for="(item, index) in topOneList"
            :key="index"
            :class="['content-box']"
          >
            <p>
              <span class="title-name">{{ item.name }}</span>
              <el-popover
                :ref="`popoverA${index}`"
                placement="top"
                width="100"
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
            </p>
            <p class="sum">
              <template v-if="item.key !== 'subtasks_finish_percentage'">
                <div>{{ (mainTaskData[item.key] || 0) | formatNumber }}</div>
              </template>
              <template v-else>
                <div>
                  {{ mainTaskData[item.key] || 0 }}
                </div>
              </template>
            </p>
          </div>
        </div>

        <div class="top-box width-60">
          <div
            v-for="(item, index) in topTwoList"
            :key="index"
            :class="['content-box']"
          >
            <p>
              <span class="title-name">{{ item.name }}</span>
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
            </p>
            <p class="sum">
              <template v-if="item.key !== 'task_finish_percentage'">
                <div>{{ (mainTaskData[item.key] || 0) | formatNumber }}</div>
              </template>
              <template v-else>
                <div>
                  {{ mainTaskData[item.key] || 0 }}
                </div>
              </template>
            </p>
          </div>
        </div>
        <div class="top-box">
          <div
            v-for="(item, index) in topThreeList"
            :key="index"
            :class="[
              'content-box',
              index !== 0 && index % 4 == 0 ? 'mar-r-0' : ''
            ]"
            style="width:20%;"
          >
            <p>
              <span class="title-name">{{ item.name }}</span>
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
            </p>
            <p class="sum">
              {{ (mainTaskData[item.key] || 0) | formatNumber }}
            </p>
          </div>
        </div>
      </template>
    </layout-list>
    <layout-list class="margin-top-30">
      <template slot="table">
        <div class="table-title">数据明细</div>
        <div style="padding:20px 0px;">
          <search-item text="子任务名称：">
            <cy-input
              size="small"
              width-size="small"
              placeholder="请输入子任务名称查询"
              v-model.trim="taskName"
              prefix-icon="el-icon-search"
              maxlength="20"
              clearable
            ></cy-input>
          </search-item>
          <!-- todo：子任务类型 -->
          <search-item text="子任务类型：">
            <cy-select
              v-model="task_category"
              size="medium"
              :optionsList="taskTypeList"
              :options="{ label: 'name', value: 'value' }"
              @change="changeTaskCategory"
            >
            </cy-select>
          </search-item>
          <!-- todo：任务群发类型 -->
          <search-item text="任务群发类型：">
            <cy-select
              v-model="group_send_type"
              size="medium"
              :optionsList="taskGroupTypeList"
              :options="{ label: 'name', value: 'value' }"
              @change="changeGroupTask"
            >
            </cy-select>
          </search-item>
          <!-- todo：条件维度 -->
          <search-item text="条件维度：">
            <cy-select
              v-model="taskType"
              size="medium"
              :optionsList="taskStatusList"
              :options="{ label: 'name', value: 'value' }"
            >
            </cy-select>
          </search-item>
          <el-button
            size="small"
            type="primary"
            @click="searchEvn"
          >
            查询
          </el-button>
          <el-button
            class="mar-l-20"
            size="small"
            @click="reset"
          >
            重置
          </el-button>
          <el-button
            class="mar-l-20"
            size="small"
            @click="exportFile"
          >
            导出
          </el-button>
        </div>
        <el-table
          stripe
          :data="list"
          class="mar-t-20 task-data"
          ref="taskDataTable"
        >
          <el-table-column
            width="200"
            label="子任务名称"
          >
            <template slot-scope="scope">
              <span
                class="c-blue plan-id-box"
                @click="viewChildEvn(scope.row.id)"
              >{{ scope.row.task_name }}</span>
            </template>
          </el-table-column>
          <!-- todo:子任务类型 -->
          <el-table-column
            min-width="140"
            label="子任务类型"
          >
            <template slot-scope="scope">
              <span v-if="scope.row.task_category === 0">沟通任务</span>
              <span v-if="scope.row.task_category === 1">朋友圈任务</span>
            </template>
          </el-table-column>
          <!-- todo:任务群发类型 -->
          <el-table-column
            min-width="140"
            label="任务群发类型"
          >
            <template slot-scope="scope">
              <span v-if="scope.row.task_category === 1">
                <span v-if="scope.row.group_send_type === 0">标签任务群发</span>
                <span v-if="scope.row.group_send_type === 1">企微一键群发</span>
              </span>
              <span v-if="scope.row.task_category === 0">
                <span v-if="scope.row.group_send_type === 1">企微一键群发</span>
                <span v-if="scope.row.group_send_type === 2">灵活沟通转发</span>
              </span>
            </template>
          </el-table-column>
          <el-table-column
            min-width="125"
            label="条件维度"
          >
            <template slot-scope="scope">
              {{ getTaskStatusName(scope.row) }}
            </template>
          </el-table-column>
          <!-- 客户画像 -->
          <el-table-column
            v-if="!taskOrPortrait"
            min-width="140"
            label="触达导购总人数"
          >
            <template slot-scope="scope">
              <span>
                {{ scope.row.reach_guide_number | formatNumber }}
              </span>
            </template>
          </el-table-column>
          <el-table-column
            v-if="!taskOrPortrait"
            min-width="140"
            label="匹配客户总人数"
          >
            <template slot-scope="scope">
              <span>
                {{ scope.row.match_customers_number | formatNumber }}
              </span>
            </template>
          </el-table-column>
          <el-table-column
            v-if="!taskOrPortrait"
            min-width="200"
            label="匹配客户总人数（去重）"
          >
            <template slot-scope="scope">
              <span>
                {{ scope.row.match_customers_number2 | formatNumber }}
              </span>
            </template>
          </el-table-column>
          <!-- 社群数据 -->
          <el-table-column
            v-if="taskOrPortrait"
            min-width="140"
            label="触达群主总人数"
          >
            <template slot-scope="scope">
              <span>
                {{ scope.row.reach_lord_number | formatNumber }}
              </span>
            </template>
          </el-table-column>
          <el-table-column
            v-if="taskOrPortrait"
            min-width="140"
            label="匹配社群总数"
          >
            <template slot-scope="scope">
              <span>
                {{ scope.row.match_group_number | formatNumber }}
              </span>
            </template>
          </el-table-column>
          <el-table-column
            min-width="140"
            label="任务总数（件）"
          >
            <template slot-scope="scope">
              <span>
                {{ scope.row.task_number | formatNumber }}
              </span>
            </template>
          </el-table-column>
          <el-table-column
            min-width="170"
            label="已完成任务数（件）"
          >
            <template slot-scope="scope">
              <span>
                {{
                  scope.row.task_finish_number == 0
                  ? "-"
                  : scope.row.task_finish_number | formatNumber
                }}
              </span>
            </template>
          </el-table-column>
          <el-table-column
            min-width="160"
            label="任务完成率（%）"
          >
            <template slot-scope="scope">
              <span>
                {{ scope.row.task_finish_percentage }}
              </span>
            </template>
          </el-table-column>
          <el-table-column
            prop="zip"
            label="操作"
            width="150"
            fixed="right"
          >
            <template slot-scope="scope">
              <span
                class="c-blue"
                @click="viewChildEvn(scope.row.id)"
              >子任务完成情况
              </span>
            </template>
          </el-table-column>
          <div
            slot="empty"
            class="table-empty-label"
          >
            <img
              src="@/assets/images/empty-box.png"
              alt=""
            />
            <p>暂无子任务计划</p>
          </div>
        </el-table>
      </template>
      <!-- 分页 -->
      <template slot="footer">
        <el-pagination
          v-if="list && list.length > 0 && page_init"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page="meta.page"
          :page-sizes="[10, 20, 30, 40]"
          :page-size="meta.size"
          class="page-right font-weight-400"
          layout="total, sizes, prev, pager, next, jumper"
          :total="meta.total"
        ></el-pagination>
      </template>
    </layout-list>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Api from "@/api/wisdom/planTask";
// todo: 企业任务类型
import {
  tag_task,
  group_task,
  task_type,
  task_group_type
} from "../components/dimension_data";
import { topOneList, topTwoList, topThreeList } from "./list-data";
import PerModule from "@/store/modules/permission";
// import { Col } from "node_modules/_element-ui@2.15.5@element-ui/types";
@Component
export default class extends Vue {
  get ep() {
    return PerModule.enterprise
  }
  private taskObj: any = {};
  private mainTaskData: any = {
    plan_name: "",
    period_subtasks_number: "",
    plan_start_time: "",
    plan_end_time: ""
  };
  private searchTime: any = "";
  private startTime: any = "";
  private endTime: any = "";
  private taskName = "";
  private taskType: any = "";
  // todo: 查询条件
  private task_category: any = "";
  private group_send_type: any = "";
  private taskOrPortrait = false;
  private topOneList = topOneList;
  private topTwoList = topTwoList;
  private topThreeList = topThreeList;
  // todo: 数据抽离
  private taskStatusList: Array<any> = group_task;
  //  todo: 任务类型
  private taskTypeList: Array<any> = task_type;
  
    
  private page_init = false;
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  private list: Array<any> = [];

  get taskGroupTypeList() {
    return this.ep.isCoach ? task_group_type : JSON.parse(JSON.stringify(task_group_type)).slice(1);
  }

  get updateTime() {
    return (Date.parse(this.searchTime) - 60 * 60 * 24 * 1000) / 1000;
  }

  get timeRange() {
    console.log([this.startTime * 1000, this.endTime* 1000]);
    if(this.startTime) {
      return [this.startTime * 1000, this.endTime* 1000]
    }
    return []
  }
  set timeRange(value) {
    this.startTime = parseInt(value[0]/1000 + '')
    this.endTime = parseInt(value[1]/1000 + '')
  }

  private mounted() {
    this.searchTime = new Date();
    this.taskObj = this.$route.query.selectedTask;
    this.taskObj = JSON.parse(this.taskObj);
    console.log(11111, this.taskObj);
    this.startTime = this.taskObj.created_at
    this.endTime = parseInt(Date.now() / 1000 + '');
    
    this.getTaskData();
    this.getChildTaskList();
    // console.log("this.taskObj", this.taskObj);
    this.pickerOptions = {
      disabledDate: (time: any) => {

        return (
          time.getTime() / 1000 < this.taskObj.created_at - 24 * 3600 ||
          time.getTime() / 1000 > new Date().setHours(23,59,59) / 1000
        );
      }
    };
    // todo: el-table不对齐
    this.$nextTick(() => {
      const table = this.$refs.taskDataTable as any;
      table.doLayout();
    });
  }
  // todo:子任务类型change
  changeTaskCategory(val: number) {
    // 清空dimension
    this.taskStatusList = [];
    // 清空任务类型
    this.taskType = "";
    // 清空任务群发类型
    this.group_send_type = "";
    if (val === 1) {
      // 如果是朋友圈任务删除dimension_list社群数据维度
      this.taskStatusList = tag_task.filter((item: any) => {
        return item.value !== 1;
      });
    } else {
      this.taskStatusList = tag_task;
    }
  }
  // todo:任务群发类型change
  changeGroupTask(val: number) {
    // 清空dimension
    this.taskStatusList = [];
    // 清空任务类型
    this.taskType = "";
    if (val === 1) {
      this.taskStatusList = group_task;
    } else {
      // 如果是朋友圈任务删除dimension_list社群数据维度
      if (this.task_category === 1) {
        // 如果是朋友圈任务删除dimension_list社群数据维度
        this.taskStatusList = tag_task.filter((item: any) => {
          return item.value !== 1;
        });
      } else {
        this.taskStatusList = tag_task;
      }
    }
  }
  // todo: 条件维度获取名称
  getTaskStatusName(row: any) {
    let name: any = "";
    if (row.group_send_type === 0) {
      name = tag_task.find((item: any) => item.value === row.dimension)?.name;
    } else {
      name = group_task.find((item: any) => item.value === row.dimension)?.name;
    }
    return name;
  }
  // 选择时间
  private pickerOptions = {};
  // 根据时间查询
  private changeSelector() {
    // this.searchTime = new Date(this.searchTime);
    this.getTaskData();
    this.getChildTaskList();
  }
  // 查询
  private searchEvn() {
    this.meta.page = 1;
    this.list = [];
    this.meta.total = 0;
    if (this.group_send_type === 0 && this.taskType === 1) {
      console.log("社群数据");
      this.taskOrPortrait = true;
    } else {
      this.taskOrPortrait = false;
    }
    // 是否正在搜索
    // this.searchLoading = true;
    this.getChildTaskList();
  }
  // 重置
  private reset() {
    // todo：初始化数据
    this.task_category = "";
    this.group_send_type = "";
    this.taskType = "";
    this.taskName = "";
    this.meta.page = 1;
    this.meta.size = 10;
    if (this.group_send_type === 0 && this.taskType === 1) {
      console.log("社群数据");
      this.taskOrPortrait = true;
    } else {
      this.taskOrPortrait = false;
    }
    // 是否正在搜索
    // this.searchLoading = false;
    this.getChildTaskList();
  }
  // 数据明细导出
  private exportFile() {
    // console.log(this.searchTime);
    const params = {
      start_time: 0,
      end_time: Date.parse(this.getSearchTime()) / 1000,
      plan_id: this.taskObj.id,
      task_name: this.taskName,
      dimension: this.taskType !== '' ? this.taskType : -1,
      // todo:查询条件
      task_category: this.task_category !== '' ? this.task_category : -1,
      group_send_type: this.group_send_type !== '' ? this.group_send_type : -1,
    };
    Api.mainTaskDataExport(params).then((res: any) => {
      const blob = new Blob([res], {
        type: "application/zip;charset=utf-8"
      });
      let fileName = "";
      if (this.taskType === 1) {
        fileName =
          "社群数据数据明细" +
          new Date().toLocaleString().substr(0, 10) +
          ".xlsx";
      } else {
        fileName =
          "数据明细" + new Date().toLocaleString().substr(0, 10) + ".xlsx";
      }
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
  private handleSizeChange(val: number) {
    this.meta.size = val;
    this.getChildTaskList();
  }
  private handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getChildTaskList();
  }
  private getSearchTime() {
    const year = this.searchTime.getFullYear();
    let month = this.searchTime.getMonth() + 1;
    let day = this.searchTime.getDate();
    if (month < 10) {
      month = "0" + month;
    }
    if (day < 10) {
      day = "0" + day;
    }
    const nowDate = `${year}-${month}-${day} 23:59:59`;
    return nowDate;
  }
  // 顶部统计数据
  private getTaskData() {
    // const startTime = this.getSearchTime();
    const params = {
      plan_id: this.taskObj.id,
      start_time: this.startTime,
      end_time: this.endTime
    };
    Api.mainTask(params).then((res: any) => {
      if (res.code === 200) {
        const nowData = res.data;
        // 交易转化率：(下单客户数/触达客户数)*100%
        nowData.trade_transform =
          nowData.total_order_customer_number > 0 &&
            nowData.total_reach_customer_number > 0
            ? `${Math.floor(
              (nowData.total_order_customer_number /
                nowData.total_reach_customer_number) *
              1000
            ) / 10}%`
            : "0%";
        nowData.total_order_fee = nowData.total_order_fee / 100;
        this.mainTaskData = nowData;
      }
    });
  }
  // 数据明细
  private getChildTaskList() {
    const params = {
      start_time: this.startTime,
      end_time: this.endTime,
      plan_id: this.taskObj.id,
      task_name: this.taskName,
      dimension: this.taskType !== '' ? this.taskType : -1,
      // todo:查询条件
      task_category: this.task_category !== '' ? this.task_category : -1,
      group_send_type: this.group_send_type !== '' ? this.group_send_type : -1,
      page: this.meta.page,
      size: this.meta.size
    };
    Api.mainTaskData(params).then((res: any) => {
      if (res.code === 200) {
        console.log(res.data, 48968597459);
        this.page_init = true;
        this.list = res.data.data;
        this.meta.total = res.data.count;
      }
    });
  }
  //查看子任务
  private viewChildEvn(task_id: any) {
    this.$router.push({
      name: "myapp_wisdom_task-center_plan-task_child-task-detail",
      query: {
        task_id,
        taskObj: JSON.stringify(this.taskObj)
      }
    });
  }
}
</script>
<style lang="scss" scoped>
.child-task-data {
  .margin-top-30 {
    margin-top: 30px;
  }

  /deep/.search-content {
    background: white;
    padding: 0;
    margin-bottom: 0;
  }
}

.content-item {
  font-size: 14px;
  display: flex;
  justify-content: space-between;
  line-height: 30px;
  color: #303133;
  background: #f7f7f7;
  padding: 20px 30px;
  font-family: PingFangSC-Medium, PingFang SC;

  &.m-b-20 {
    margin-bottom: 20px;
  }

  &.flex-start {
    justify-content: flex-start;
    align-items: center;
  }
}

.top-box {
  display: flex;
  justify-content: flex-start;
  align-items: center;

  &.width-60 {
    width: calc((100% - 120px) / 5 * 3 + 90px);

    .content-box {
      width: calc(100% / 3);
      margin-right: 30px;
    }
  }

  .content-box {
    height: 137px;
    background: #f7f7f7;
    display: flex;
    justify-content: center;
    align-items: center;
    flex-direction: column;
    margin-bottom: 30px;
    margin-right: 30px;

    &.mar-r-0 {
      margin-right: 0;
    }

    &.width-30 {
      width: calc(20% - 30px);
    }

    &.m-b-0 {
      margin-bottom: 0;
    }

    p {
      text-align: center;
      line-height: 26px;
      font-size: 17px;
    }
  }
}

.table-title {
  color: #303133;
  font-family: PingFangSC-Medium, PingFang SC;
  font-weight: 500;
  font-size: 22px;
  margin-bottom: 20px;
}

.c-blue {
  cursor: pointer;
  color: $color-main;
  margin-right: 8px;
}
.help-icon {
  margin-left: 5px;
  cursor: pointer;
  width: 13px;
}

.sum {
  font-size: 24px !important;
  font-weight: 600;
  margin-top: 10px;
}

.title-name {
  font-size: 14px !important;
}

.data-info {
  .time {
    color: $color-text-normal;
    font-size: 14px;
    line-height: 20px;
  }
}
</style>
