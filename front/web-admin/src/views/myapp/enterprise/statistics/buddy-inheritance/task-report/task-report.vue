<template>
  <layout-list>
    <template slot="search">
      <div class="top-title">
        <p>
          任务名称：
          <span v-if="taskInfo.name">
            {{
              taskInfo.name.length > 10
                ? taskInfo.name.slice(0, 10) + "..."
                : taskInfo.name
            }}</span
          >
        </p>
        <p>客户范围：全部客户</p>
        <p>继承类型：{{ taskInfo.type === 1 ? "唯一继承" : "随机继承" }}</p>
        <p>
          任务状态：<span class="task-status">{{
            taskInfo.status === 1
              ? "未开始"
              : taskInfo.status === 2
              ? "进行中"
              : taskInfo.status === 3
              ? "已完成"
              : "任务取消"
          }}</span>
        </p>
      </div>
      <div class="top-data">
        <div
          class="top-data-content"
          v-for="(item, index) in topTabList"
          :key="index"
        >
          <p>{{ item.name }}</p>
          <p>
            {{ reportTitle[item.key] }}
          </p>
        </div>
      </div>
    </template>
    <template slot="table">
      <el-table
        tooltip-effect="light"
        :data="staffList"
        stripe
        :header-cell-style="{ textAlign: 'center' }"
        :cell-style="{ textAlign: 'center' }"
        style="width: 100%"
      >
        <el-table-column label="客户名称" width="200" show-overflow-tooltip>
          <template slot-scope="scope">
            <span>{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column
          prop="mobile"
          label="好友数"
          min-width="200"
          show-overflow-tooltip
        >
          <template slot-scope="scope">
            <span v-if="scope.row.friend_num">{{ scope.row.friend_num }}</span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column
          label="当前管理导购"
          min-width="200"
          show-overflow-tooltip
        >
          <template slot-scope="scope">
            <span v-if="scope.row.current_admin_ba_name">
              {{ scope.row.current_admin_ba_name }}
            </span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="所属部门" min-width="200" show-overflow-tooltip>
          <template slot-scope="scope">
            <span v-if="scope.row.ba_department_name">{{
              scope.row.ba_department_name
            }}</span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="继承状态" min-width="150">
          <template slot-scope="scope">
            <span v-if="scope.row.transfer_status_msg">{{
              scope.row.transfer_status_msg
            }}</span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="失败原因" min-width="150">
          <template slot-scope="scope">
            <span>{{ scope.row.fail_msg }}</span>
          </template>
        </el-table-column>
        <div slot="empty" class="table-empty-label">
          <img class="empty-img" src="@/assets/images/empty-box.png" alt="" />
          <div>暂无数据</div>
        </div>
      </el-table>
    </template>
    <template slot="footer">
      <el-button size="medium" @click="goBack">返回</el-button>
      <el-pagination
        v-if="meta.total > 0"
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
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Api from "@/api/enterprise/buddy-inheritance";
@Component
export default class extends Vue {
  private topTabList = [
    {
      name: "客户总数 (人）",
      key: "customer_total_num"
    },
    {
      name: "成功继承客户数(人)",
      key: "customer_success_num"
    },
    {
      name: "继承失败客户数（人）",
      key: "customer_fail_num"
    },
    {
      name: "完成继承导购数（人）",
      key: "finish_ba_num"
    },
    {
      name: "继承完成率",
      key: "finish_rate"
    }
  ];
  private taskInfo = {};
  private reportTitle = {
    customer_total_num: 0,
    customer_success_num: 0,
    customer_fail_num: 0,
    finish_ba_num: 0,
    finish_rate: "0"
  };
  private staffId = 0;
  // 表格数组
  private staffList: Array<any> = [];
  private onlyList = false;
  // 分页配置
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  // 生命周期
  private mounted() {
    this.staffId = Number(this.$route.query.id);
    this.onlyList = false;
    this.getList();
  }
  // 分页
  private handleSizeChange(val: string) {
    this.onlyList = true;
    this.$set(this.meta, "size", val);
    this.meta.page = 1;
    this.getList();
  }
  private handleCurrentChange(val: string) {
    this.onlyList = true;
    this.$set(this.meta, "page", val);
    this.getList();
  }
  private getList() {
    const params = {
      id: this.staffId,
      is_only_list: this.onlyList,
      page: this.meta.page,
      page_size: this.meta.size
    };
    Api.tasksReport(params).then((res: any) => {
      if (res.code == 200) {
        this.staffList = res.data.report_body.data_list
          ? res.data.report_body.data_list
          : [];
        if (!this.onlyList) {
          this.taskInfo = res.data.task_info;
          this.reportTitle = res.data.report_title;
          this.reportTitle.finish_rate = this.reportTitle.finish_rate + "%";
          this.meta.total = res.data.report_title.customer_total_num;
        }
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 返回
  private goBack() {
    this.$router.go(-1);
  }
}
</script>
<style lang="scss" scoped>
.top-title {
  display: flex;
  padding-bottom: 15px;
  border-bottom: 1px dashed #aaaaaa;
  color: #555555;
  p {
    width: 25%;
    text-align: center;
    .task-status {
      color: $color-text-button;
    }
  }
}
.top-data {
  display: flex;
  padding-top: 15px;
  .top-data-content {
    width: 20%;
    box-sizing: border-box;
    border-right: 1px solid #e6e7e8;
    text-align: center;
    color: #333333;
    line-height: 30px;
    p:last-child {
      font-size: 24px;
      font-weight: 600;
    }
  }
  .top-data-content:last-child {
    border-right: none;
  }
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
</style>
