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
    </template>
    <template v-if="ep.isCoach" slot="search">
      <div>
        <search-item text="子任务ID：">
          <cy-input
            size="small"
            width-size="small"
            placeholder="请输入子任务ID"
            v-model.trim="searchForm.task_id"
            prefix-icon="el-icon-search"
            clearable
          ></cy-input>
        </search-item>
        <search-item text="任务模板ID：">
          <cy-input
            size="small"
            width-size="small"
            placeholder="请输入任务模板ID"
            v-model.trim="searchForm.template_id"
            prefix-icon="el-icon-search"
            clearable
          ></cy-input>
        </search-item>
        <search-item text="子任务名称：">
          <cy-select
            v-model="searchForm.task_name"
            size="medium"
            search-input
            :optionsList="taskNameList"
            :options="{ label: 'label', value: 'value' }"
          >
          </cy-select>
        </search-item>
        <search-item text="任务类型：">
          <cy-select
            v-model="searchForm.category"
            size="medium"
            :optionsList="categoryList"
            :options="{ label: 'label', value: 'value' }"
          >
          </cy-select>
        </search-item>
        <search-item text="任务群发类型：">
          <cy-select
            v-model="searchForm.group_type"
            size="medium"
            :optionsList="groupTypeList"
            :options="{ label: 'label', value: 'value' }"
          >
          </cy-select>
        </search-item>
        <search-item v-if="ep.isCoach" text="子任务来源：">
          <cy-select
            v-model="searchForm.origin_type"
            size="medium"
            :optionsList="originTypeList"
            :options="{ label: 'label', value: 'value' }"
          >
          </cy-select>
        </search-item>
        <search-item text="子任务创建时间：">
          <el-date-picker
            v-model="createTime"
            type="datetimerange"
            range-separator="-"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            :default-time="['00:00:00', '23:59:59']"
            :picker-options="pickerOptions"
          >
          </el-date-picker>
        </search-item>
        <search-item text="子任务开始：">
          <el-date-picker
            v-model="startTime"
            type="datetimerange"
            range-separator="-"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            :default-time="['00:00:00', '23:59:59']"
            :picker-options="pickerOptions"
          >
          </el-date-picker>
        </search-item>
        <el-button size="small" type="primary" @click="handleSearch">
          查询
        </el-button>
        <el-button class="mar-l-20" size="small" @click="handleReset">
          重置
        </el-button>
      </div>
    </template>
    <template slot="table">
      <el-table stripe :data="list" class="mar-t-20">
        <el-table-column label="子任务ID" width="180">
          <template slot-scope="scope">
            <span
              class="c-blue plan-id-box"
              @click="viewChildEvn(scope.row.id)"
              >{{ scope.row.id }}</span
            >
          </template>
        </el-table-column>
        <el-table-column label="子任务名称" width="150">
          <template slot-scope="scope">
            <span
              class="c-blue plan-id-box"
              @click="viewChildEvn(scope.row.id)"
              >{{ scope.row.name }}</span
            >
          </template>
        </el-table-column>
        <el-table-column min-width="130" label="任务类型">
          <template slot-scope="scope">
            <span v-if="scope.row.task_category === 0">沟通任务</span>
            <span v-if="scope.row.task_category === 1">朋友圈任务</span>
          </template>
        </el-table-column>
        <el-table-column min-width="130" label="任务群发类型">
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
        <el-table-column min-width="130" label="子任务维度">
          <template slot-scope="scope">
            <span v-if="scope.row.dimension === 1">社群数据</span>
            <span v-if="scope.row.dimension === 2">用户画像</span>
            <span v-if="scope.row.dimension === 3">导入顾客名单</span>
            <span v-if="scope.row.dimension === 4">导入顾客&导购匹配关系</span>
          </template>
        </el-table-column>
        <el-table-column min-width="130" label="子任务条件">
          <template slot-scope="scope">
            <span v-if="scope.row.dimension <= 2">{{
              scope.row.condition_name
            }}</span>
            <span v-else>--</span>
          </template>
        </el-table-column>
        <el-table-column v-if="ep.isCoach" min-width="130" label="任务完成标准">
          <template slot-scope="scope">
            <span v-if="scope.row.finish_rate && !(scope.row.task_category === 0 &&scope.row.group_send_type === 1)">{{scope.row.finish_rate}}%</span>
            <span v-else>--</span>
          </template>
        </el-table-column>
        <el-table-column v-if="ep.isCoach" label="子任务有效期" width="230">
          <template slot-scope="scope">
            <div v-if="scope.row.task_start_time">
              {{ scope.row.task_start_time | formatDate("yyyy-mm-dd hh:mm:ss") }} ~
            </div>
            <div v-if="scope.row.task_start_time">
              {{ scope.row.task_end_time | formatDate("yyyy-mm-dd hh:mm:ss") }}
            </div>
            <span v-if="!scope.row.task_start_time">--</span>
          </template>
        </el-table-column>
        <el-table-column v-if="ep.isCoach" min-width="130" label="子任务创建时间">
          <template slot-scope="scope">
            <span>{{scope.row.created_at | formatDate("yyyy-mm-dd hh:mm:ss")}}</span>
          </template>
        </el-table-column>
        <el-table-column v-if="ep.isCoach" min-width="130" label="任务模板ID">
          <template slot-scope="scope">
            <span>{{scope.row.plan_template_id || '--'}}</span>
          </template>
        </el-table-column>
        <el-table-column v-if="!ep.isCoach" label="最近更新时间" width="230">
          <template slot-scope="scope">
            <div>
              {{ scope.row.created_at | formatDate("yyyy-mm-dd hh:mm:ss") }} ~
            </div>
            <div>
              {{ scope.row.updated_at | formatDate("yyyy-mm-dd hh:mm:ss") }}
            </div>
          </template>
        </el-table-column>
        <el-table-column v-if="ep.isCoach" label="子任务来源" width="230">
          <template slot-scope="scope">
            <div>{{ originTypeName(scope.row.origin_type) }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="zip" label="操作" width="200" fixed="right">
          <template slot-scope="scope">
            <span class="c-blue" @click="viewChildEvn(scope.row.id)"
              >查看
            </span>
            <span
              class="c-blue"
              v-if="scope.row.status === 0"
              @click="copyChildEvn(scope.row.id)"
              >复制
            </span>
            <!-- 任务未开始 & (企微沟通任务 | 朋友圈任务 & 用户画像) 可以编辑 -->
            <span
              class="c-blue"
              v-if="
                scope.row.status === 0 && (scope.row.task_category === 1 && (scope.row.dimension === 1 || scope.row.dimension === 2) || scope.row.task_category === 0)
              "
              @click="editChildEvn(scope.row.id)"
            >
              编辑
            </span>
            <span
              class="c-blue"
              v-if="scope.row.status === 0"
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
import PerModule from "@/store/modules/permission";
@Component
export default class extends Vue {
  get ep() {
    return PerModule.enterprise
  }
  searchForm = {
    task_id: null,
    template_id: null,
    category: -1,
    task_name: '',
    group_type: -1,
    created_start_time: 0,
    created_end_time: 0,
    task_star_time: 0,
    task_end_time: 0,
    origin_type: -1, // '0:默认手动创建，1:convertlab创建，2:startDt创建
  };

  taskNameList = []

  categoryList = [
    { label: "全部", value: -1 },
    { label: "沟通任务", value: 0 },
    { label: "朋友圈任务", value: 1 },
  ]
  originTypeList = [
    { label: "全部", value: -1 },
    { label: "默认手动创建", value: 0 },
    { label: "convertlab创建", value: 1 },
    { label: "startDt创建", value: 2 },
    { label: "CGE创建", value: 3 },
  ]
  groupTypeList = [
    { label: "全部", value: -1 },
    { label: "标签任务群发", value: 0 },
    { label: "企微一键群发", value: 1 },
    { label: "灵活沟通转发", value: 2 },
  ]

  get originTypeName() {
    return (value) => {
      return this.originTypeList.find((item) => item.value === value)?.label || '--';
    }
  }

  /** 创建时间 读取 */
  get createTime() {
    const { created_start_time, created_end_time } = this.searchForm
    return [created_start_time*1000, created_end_time*1000].filter(Boolean)
  }

  /** 创建时间 写入 */
  set createTime(value: number[]) {
    const [start_time, end_time] = value || []
    this.searchForm.created_start_time = parseInt(start_time/1000+'');
    this.searchForm.created_end_time = parseInt(end_time/1000+'');
  }

  /** 开始时间 读取 */
  get startTime() {
    const { task_star_time, task_end_time } = this.searchForm
    return [task_star_time*1000, task_end_time*1000].filter(Boolean)
  }

  /** 开始时间 写入 */
  set startTime(value: number[]) {
    const [start_time, end_time] = value || []
    this.searchForm.task_star_time = parseInt(start_time/1000+'');
    this.searchForm.task_end_time = parseInt(end_time/1000+'');
  }

  private taskObj: any = {};
  private page_init = false;
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  private list: Array<any> = [];
  private mounted() {
    this.taskObj = this.$route.query.selectedTask;
    this.taskObj = JSON.parse(this.taskObj);
    this.getChildTaskList();
    this.getPlanName()
  }
  handleSearch() {
    this.meta.page = 1;
    this.getChildTaskList();
  }
  handleReset() {
    this.searchForm = {
      task_id: null,
      template_id: null,
      category: -1,
      task_name: '',
      group_type: -1,
      created_start_time: 0,
      created_end_time: 0,
      task_star_time: 0,
      task_end_time: 0,
      origin_type: -1,
    }
    this.handleSearch()
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
      ...this.searchForm
    };
    params.task_id = params.task_id ? Number(params.task_id): 0
    params.template_id = params.template_id ? Number(params.template_id): 0
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
            case 8:
              item.condition_name = "顾客生日";
              break;
            case 9:
              item.condition_name = "全部客户";
              break;
            default:
              item.condition_name = "";
              break;
          }
        });
      }
    });
  }
  private getPlanName() {
    const params = {
      plan_id: this.taskObj.id,
    };
    Api.getPlanNameApi(params).then((res: any) => {
      if (res.code === 200) {
        this.taskNameList = res.data.map(item => {
          return {
            label: item,
            value: item
          }
        })
        
      }
    });
  }
  // 新增子任务
  private addChildTask() {
    this.$router.push({
      name: "myapp_wisdom_task-center_plan-task_child-task-add",
      query: {
        status: "0",
        plan_id: this.taskObj.id
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
  //复制子任务
  private copyChildEvn(task_id: any) {
    this.$router.push({
      name: "myapp_wisdom_task-center_plan-task_child-task-copy",
      query: {
        status: "2",
        plan_id: this.taskObj.id,
        st: this.taskObj.start_time,
        et: this.taskObj.end_time,
        task_id,
      }
    });
  }
  //编辑子任务
  private editChildEvn(task_id: any) {
    this.$router.push({
      name: "myapp_wisdom_task-center_plan-task_child-task-edit",
      query: {
        status: "1",
        plan_id: this.taskObj.id,
        st: this.taskObj.start_time,
        et: this.taskObj.end_time,
        task_id
      }
    });
  }
  //删除子任务
  private deleteChildEvn(task_id: string | number) {
    this.$confirm("你确认删除该任务吗？", "删除确认", {
      confirmButtonText: "确认",
      cancelButtonText: "取消"
    })
      .then(() => {
        const params = {
          task_id: task_id
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
</style>
