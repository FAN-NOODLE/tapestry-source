<template>
  <layout-list>
    <!-- 头部 -->
    <template slot="search">
      <div class="content-item content-topBox">
        <span>任务名称：{{ mainTaskData.task_name || "-" }}</span>
        <span>子任务数：{{ mainTaskData.subtasks_number || "0" }}</span>
        <span>任务有效期：{{
            mainTaskData.plan_start_time | formatDate("yyyy-mm-dd hh:mm:ss")
          }}～{{
            mainTaskData.plan_end_time | formatDate("yyyy-mm-dd hh:mm:ss")
          }}</span>
      </div>
      <div>
        <div>
          时间：
          <el-date-picker
            v-model="searchTime"
            type="date"
            placeholder="选择日期"
            value-format="yyyy-MM-dd"
            :picker-options="pickerOptions"
            @blur="changeSelector()"
          >
          </el-date-picker>
          <!-- <el-date-picker
            v-model="searchTime"
            type="daterange"
            range-separator="-"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            value-format="yyyy-MM-dd"
            :picker-options="pickerOptions"
          >
          </el-date-picker> -->
        </div>
        <div
          class="top-box"
          style="width:60%;"
        >
          <div
            v-for="(item, index) in topOneList"
            :key="index"
            :class="[
              'content-box',
              index == topOneList.length - 1 ? '' : 'last-content-box',
            ]"
            style="width:33.33%;"
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
            <p class="sum">{{ mainTaskData[item.key] }}</p>
          </div>
        </div>
        <div
          class="top-box"
          style="width:60%;"
        >
          <div
            v-for="(item, index) in topTwoList"
            :key="index"
            :class="[
              'content-box',
              index == topTwoList.length - 1 ? '' : 'last-content-box',
            ]"
            style="width:33.33%;"
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
            <p class="sum">{{ mainTaskData[item.key] }}</p>
          </div>
        </div>
        <div
          class="top-box"
          style="width:100%;"
        >
          <div
            v-for="(item, index) in topThreeList"
            :key="index"
            :class="[
              'content-box',
              index == topThreeList.length - 1 ? '' : 'last-content-box',
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
            <p class="sum">{{ mainTaskData[item.key] }}</p>
          </div>
        </div>
      </div>
    </template>
    <!-- 列表 -->
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
        <search-item text="任务类型：">
          <cy-select
            v-model="taskType"
            size="medium"
            :optionsList="taskStatusList"
            :options="{ label: 'label', value: 'value' }"
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
        class="mar-t-20"
      >
        <el-table-column
          width="125"
          label="子任务名称"
        >
          <template slot-scope="scope">
            <span
              class="c-blue plan-id-box"
              @click="viewChildEvn(scope.row.id)"
            >{{ scope.row.subtasks_number }}</span>
          </template>
        </el-table-column>
        <el-table-column
          min-width="125"
          label="任务类型"
        >
          <template slot-scope="scope">
            <span>
              {{ scope.row.task_type === 1 ? "社群数据" : "用户画像" }}
            </span>
          </template>
        </el-table-column>
        <!-- 客户画像 -->
        <el-table-column
          v-if="!taskOrPortrait"
          min-width="140"
          label="累计触达导购人数"
        >
          <template slot-scope="scope">
            <span>
              {{ scope.row.reach_guide_number }}
            </span>
          </template>
        </el-table-column>
        <el-table-column
          v-if="!taskOrPortrait"
          min-width="130"
          label="累计匹配客户数"
        >
          <template slot-scope="scope">
            <span>
              {{ scope.row.match_customers_number }}
            </span>
          </template>
        </el-table-column>
        <el-table-column
          v-if="!taskOrPortrait"
          min-width="130"
          label="累计匹配客户数（导购去重）"
        >
          <template slot-scope="scope">
            <span>
              {{ scope.row.match_customers_number2 }}
            </span>
          </template>
        </el-table-column>
        <!-- 社群数据 -->
        <el-table-column
          v-if="taskOrPortrait"
          min-width="140"
          label="触达群主"
        >
          <template slot-scope="scope">
            <span>
              {{ scope.row.reach_lord_number }}
            </span>
          </template>
        </el-table-column>
        <el-table-column
          v-if="taskOrPortrait"
          min-width="140"
          label="匹配社群数"
        >
          <template slot-scope="scope">
            <span>
              {{ scope.row.match_group_number }}
            </span>
          </template>
        </el-table-column>
        <!-- ------- -->
        <el-table-column
          min-width="130"
          label="任务总数（件）"
        >
          <template slot-scope="scope">
            <span>
              {{ scope.row.task_number }}
            </span>
          </template>
        </el-table-column>
        <el-table-column
          min-width="130"
          label="已完成任务数（件）"
        >
          <template slot-scope="scope">
            <span>
              {{
                scope.row.task_finish_number == 0
                  ? "-"
                  : scope.row.task_finish_number
              }}
            </span>
          </template>
        </el-table-column>
        <el-table-column
          min-width="130"
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
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Api from "@/api/wisdom/planTask";
// import { Col } from "node_modules/_element-ui@2.15.5@element-ui/types";
@Component
export default class extends Vue {
  private taskObj: any = {};
  private mainTaskData: any = {
    task_name: "",
    subtasks_number: "",
    plan_start_time: "",
    plan_end_time: "",
    period_subtasks_number: 0, //时间段内的子任务数
    subtasks_finish_number: 0, //子任务完成数
    subtasks_finish_percentage: 0, //子任务完成率
    task_number: 0, //任务总件数
    task_finish_number: 0, //已完成任务数
    task_finish_percentage: 0, //任务完成率
    reach_guide_number: 0, //触达导购人数
    match_customers_number: 0, //匹配客户数
    match_customers_number2: 0, //匹配客户数（去重）
    reach_lord_number: 0, //触达群主数量
    match_group_number: 0, //匹配社群数
  };
  private searchTime: any = "";
  private endTime: any = "";
  private taskName = "";
  private taskType = 1;
  private taskOrPortrait = true;
  private topOneList = [
    {
      name: "子任务数",
      desc: "任务计划所包含的所有子任务个数",
      key: "period_subtasks_number",
    },
    {
      name: "子任务完成数",
      desc: "任务计划所包含的所有子任务已完成个数",
      key: "subtasks_finish_number",
    },
    {
      name: "子任务完成率",
      desc: "任务计划的子任务完成率",
      key: "subtasks_finish_percentage",
    },
  ];
  private topTwoList = [
    {
      name: "任务总数 (件）",
      desc: "当前所有子任务所包含的任务件数总和（包含新增的社群任务件数+用户画像任务件数）",
      key: "task_number",
    },
    {
      name: "已完成任务数 (件）",
      desc: "当前所有子任务内状态为已完成的任务件数总和",
      key: "task_finish_number",
    },
    {
      name: "任务完成率",
      desc: "当前所有子任务的任务件数完成率",
      key: "task_finish_percentage",
    },
  ];
  private topThreeList = [
    {
      name: "触达导购总人数",
      desc: "当前所有用户画像类子任务所能触达的导购人数总和",
      key: "reach_guide_number",
    },
    {
      name: "累计匹配客户总人数",
      desc: "当前所有用户画像类子任务所成功匹配的客户人数总和",
      key: "match_customers_number",
    },
    {
      name: "累计匹配客户总人数（去重）",
      desc: "当前所有用户画像类子任务所成功匹配的客户人数总和",
      key: "match_customers_number2",
    },
    {
      name: "触达群主总人数",
      desc: "当前所有社群类子任务成功触达的社群群主人数之和",
      key: "reach_lord_number",
    },
    {
      name: "匹配社群总数",
      desc: "当前所有社群类子任务成功触达的社群个数之和",
      key: "match_group_number",
    },
  ];
  private taskStatusList: Array<any> = [
    {
      label: "用户画像",
      value: 2,
    },
    {
      label: "社群数据",
      value: 1,
    },
  ];
  private page_init = false;
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10, // 每页数量
  };
  private list: Array<any> = [];
  private mounted() {
    this.searchTime = new Date();
    this.endTime = Date.parse(this.searchTime) / 1000;
    console.log(this.endTime);
    this.taskObj = this.$route.query.selectedTask;
    this.taskObj = JSON.parse(this.taskObj);
    this.getTaskData();
    this.getChildTaskList();
    // console.log("this.taskObj", this.taskObj);
    this.pickerOptions = {
      disabledDate: (time: any) => {
        return (
          time.getTime() / 1000 < this.taskObj.start_time - 24 * 3600 ||
          time.getTime() / 1000 > this.taskObj.end_time
        );
      },
    };
  }
  // 选择时间
  private pickerOptions = {};
  // 根据时间查询
  private changeSelector() {
    this.searchTime = new Date(this.searchTime);
    this.getTaskData();
    this.getChildTaskList();
  }
  // 查询
  private searchEvn() {
    this.meta.page = 1;
    this.list = [];
    this.meta.total = 0;
    if (this.taskType === 1) {
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
    // 初始化数据
    this.taskType = 1;
    this.taskName = "";
    this.meta.page = 1;
    this.meta.size = 10;
    this.taskOrPortrait = true;
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
      task_type: this.taskType,
    };
    Api.mainTaskDataExport(params).then((res: any) => {
      const blob = new Blob([res], {
        type: "application/zip;charset=utf-8",
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
      start_time: 0,
      end_time: Date.parse(this.getSearchTime()) / 1000,
    };
    Api.mainTask(params).then((res: any) => {
      if (res.code === 200) {
        this.mainTaskData = res.data;
      }
    });
  }
  // 数据明细
  private getChildTaskList() {
    const params = {
      start_time: 0,
      end_time: Date.parse(this.getSearchTime()) / 1000,
      plan_id: this.taskObj.id,
      subtasks_number: this.taskName,
      task_type: this.taskType,
      page: this.meta.page,
      size: this.meta.size,
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
        taskObj: JSON.stringify(this.taskObj),
      },
    });
  }
}
</script>
<style lang="scss" scoped>
.content-item {
  font-size: 16px;
  display: flex;
  justify-content: space-between;
  line-height: 30px;
  margin-bottom: 20px;
}
.content-topBox {
  border-bottom: 1px dashed #009999;
  padding-bottom: 15px;
}
.top-box {
  display: flex;
  justify-content: space-around;
  align-items: center;
  padding: 15px 0px;
  .content-box {
    p {
      text-align: center;
      line-height: 26px;
      font-size: 17px;
    }
  }
  .last-content-box {
    border-right: 1px solid #009999;
  }
}
.table-title {
  width: 99%;
  background: #f2f2f2;
  padding: 8px;
  margin: 10px 0px;
}
.c-blue {
  cursor: pointer;
  color: $color-main;
  margin-right: 8px;
}
/deep/ .el-table__body-wrapper {
  height: 400px !important;
  overflow: auto;
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
</style>
