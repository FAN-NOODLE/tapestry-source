<template>
  <layout-list>
    <template slot="buttons">
      <div class="content-item">
        <span>{{ taskObj.name || "-" }}</span>
        <span
          >任务计划状态：{{
            taskObj.status === 0
              ? "未开始"
              : taskObj.status === 1
              ? "进行中"
              : "已结束" || "-"
          }}
        </span>
        <span>{{ taskObj.is_open ? "已开启" : "已关闭" }}</span>
      </div>
      <div class="content-item">
        任务计划有效期：{{
          taskObj.start_time | formatDate("yyyy-mm-dd hh:mm:ss")
        }}～{{ taskObj.end_time | formatDate("yyyy-mm-dd hh:mm:ss") }}
      </div>
      <div class="content-item">任务计划描述：{{ taskObj.remark || "-" }}</div>
      <el-button class="add-btn" type="primary" @click="addChildTask"
        >新增子任务</el-button
      >
    </template>
    <template slot="table">
      <el-table stripe :data="list" class="mar-t-20">
        <el-table-column label="子任务ID" width="150">
          <template slot-scope="scope">
            <span
              class="c-blue plan-id-box"
              @click="viewChildEvn(scope.row.id)"
              >{{ scope.row.id }}</span
            >
          </template>
        </el-table-column>
        <el-table-column label="子任务名称" width="300">
          <template slot-scope="scope">
            <span
              class="c-blue plan-id-box"
              @click="viewChildEvn(scope.row.id)"
              >{{ scope.row.name }}</span
            >
          </template>
        </el-table-column>
        <el-table-column min-width="130" label="子任务维度">
          <template slot-scope="scope">
            <span>
              {{ scope.row.dimension === 1 ? "社群数据" : "用户画像" }}
            </span>
          </template>
        </el-table-column>
        <el-table-column min-width="130" label="子任务条件">
          <template slot-scope="scope">
            <span>
              {{ scope.row.condition_name }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="创建时间/编辑时间" width="230">
          <template slot-scope="scope"
            ><span
              >{{ scope.row.created_at | formatDate("yyyy-mm-dd hh:mm:ss") }} ~
              {{
                scope.row.updated_at | formatDate("yyyy-mm-dd hh:mm:ss")
              }}</span
            ></template
          >
        </el-table-column>
        <el-table-column prop="zip" label="操作" width="200">
          <template slot-scope="scope">
            <span class="c-blue" @click="viewChildEvn(scope.row.id)"
              >查看
            </span>
            <span
              class="c-blue"
              v-if="taskObj.status === 0"
              @click="copyChildEvn(scope.row.id)"
              >复制
            </span>
            <span
              class="c-blue"
              v-if="taskObj.status === 0"
              @click="editChildEvn(scope.row.id)"
            >
              编辑
            </span>
            <span
              class="c-blue"
              v-if="taskObj.status === 0"
              @click="deleteChildEvn(scope.row.id)"
              >删除
            </span>
          </template>
        </el-table-column>
        <div slot="empty" class="table-empty-label">
          <img src="@/assets/images/empty-box.png" alt="" />
          <p>暂无子任务计划</p>
        </div>
      </el-table>
    </template>
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
@Component
export default class extends Vue {
  // props: {
  //   taskObj: any;
  // };
  private taskObj: any = {};
  private page_init = false;
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10, // 每页数量
  };
  private list: Array<any> = [];
  private mounted() {
    this.taskObj = this.$route.query.selectedTask;
    this.taskObj = JSON.parse(this.taskObj);
    this.getChildTaskList();
  }
  private handleSizeChange(val: number) {
    this.meta.size = val;
    this.getChildTaskList();
  }
  private handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getChildTaskList();
  }
  private getChildTaskList() {
    const params = {
      plan_id: this.taskObj.id,
      page: this.meta.page,
      size: this.meta.size,
    };
    Api.childTaskList(params).then((res: any) => {
      if (res.code === 200) {
        this.page_init = true;
        this.list = res.data.list;
        this.meta.total = res.data.total;
        this.list.forEach((item: any) => {
          switch (item.condition_type) {
            case 1:
              item.condition_name = "标签";
              break;
            case 2:
              item.condition_name = "群人数";
              break;
            case 3:
              item.condition_name = "群内新增人数";
              break;
            case 4:
              item.condition_name = "参与活动";
              break;
            case 5:
              item.condition_name = "购买行为";
              break;
            case 6:
              item.condition_name = "群活跃度";
              break;
            case 7:
              item.condition_name = "退群人群";
              break;
            default:
              item.condition_name = "";
              break;
          }
        });
      }
    });
  }
  // 新增子任务
  private addChildTask() {
    this.$router.push({
      name: "myapp_wisdom_task-center_plan-task_child-task-add",
      query: {
        status: "0",
        plan_id: this.taskObj.id,
      },
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
  //复制子任务
  private copyChildEvn(task_id: any) {
    this.$router.push({
      name: "myapp_wisdom_task-center_plan-task_child-task-copy",
      query: {
        status: "2",
        plan_id: this.taskObj.id,
        task_id,
      },
    });
  }
  //编辑子任务
  private editChildEvn(task_id: any) {
    this.$router.push({
      name: "myapp_wisdom_task-center_plan-task_child-task-edit",
      query: {
        status: "1",
        plan_id: this.taskObj.id,
        task_id,
      },
    });
  }
  //删除子任务
  private deleteChildEvn(task_id: string | number) {
    this.$confirm("你确认删除该任务吗？", "删除确认", {
      confirmButtonText: "确认",
      cancelButtonText: "取消",
    })
      .then(() => {
        const params = {
          task_id: task_id,
        };
        Api.childTaskDelete(params)
          .then((res: any) => {
            if (res.code == 200) {
              this.$message.success("删除成功");
              this.getChildTaskList();
            }
          })
          .catch((err: any) => {
            console.log("taskCenterDelete:", err);
          });
      })
      .catch((err: any) => {
        console.log(err);
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
.c-blue {
  cursor: pointer;
  color: $color-main;
  margin-right: 8px;
}
/deep/ .el-table__body-wrapper {
  height: 400px !important;
  overflow: auto;
}
</style>
