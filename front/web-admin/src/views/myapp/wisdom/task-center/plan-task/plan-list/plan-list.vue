<template>
  <div>
    <layout-list>
      <template slot="buttons">
        <el-button class="add-btn" @click="addTaskPlan">新增任务计划</el-button>
        <el-button @click="addUpload" v-if="false">上传标签</el-button>
      </template>
      <template slot="search">
        <div>
          <search-item text="任务计划搜索：">
            <el-input
              :placeholder="`请输入${type_list[gudie_type].name}`"
              class="search-input"
              size="medium"
              v-model.trim="search_key"
              maxlength="20"
              @input="changeNumberInput($event)"
              clearable
            >
              <el-select
                slot="prepend"
                v-model="gudie_type"
                class="input-with-select"
                @change="changeSelect"
              >
                <el-option
                  v-for="item in type_list"
                  :key="item.key"
                  :value="item.key"
                  :label="item.name"
                ></el-option>
              </el-select>
            </el-input>
          </search-item>
          <search-item text="任务状态：">
            <cy-select
              v-model="task_status"
              size="medium"
              :optionsList="taskStatusList"
              :options="{ label: 'label', value: 'value' }"
            >
            </cy-select>
          </search-item>
          <search-item text="任务标签：">
            <cy-select
              v-model="category"
              size="medium"
              :optionsList="taskTagList"
              :options="{ label: 'label', value: 'value' }"
            >
            </cy-select>
          </search-item>
          <search-item text="创建人：">
            <cy-input
              size="small"
              width-size="small"
              placeholder="请输入创建人名称"
              v-model.trim="creatName"
              prefix-icon="el-icon-search"
              maxlength="20"
              clearable
            ></cy-input>
          </search-item>
          <search-item text="创建时间：">
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
          <el-button size="small" type="primary" @click="searchEvn">
            查询
          </el-button>
          <el-button class="mar-l-20" size="small" @click="reset">
            重置
          </el-button>
        </div>
      </template>
      <template slot="table">
        <el-table stripe :data="list" class="mar-t-20" v-if="page_init">
          <el-table-column label="任务计划ID" width="180">
            <template slot-scope="scope">
              <span class="c-blue plan-id-box" @click="viewEvn(scope.row)">
                {{ scope.row.id }}</span
              >
            </template>
          </el-table-column>
          <el-table-column label="任务计划名称" min-width="220">
            <template slot-scope="scope">
              <span class="c-blue plan-id-box" @click="viewEvn(scope.row)">
                {{ scope.row.name }}</span
              >
            </template>
          </el-table-column>
          <el-table-column
            prop="taskTag"
            label="任务标签"
            min-width="200"
            text-align
          >
            <template slot-scope="scope">
              <div v-if="scope.row.category === 3">
                MA任务创建
                <!-- <div>({{ scope.row.template_id }})</div> -->
              </div>
              <div v-if="scope.row.category === 2">
                MA任务模板创建
                <!-- <div>({{ scope.row.template_id }})</div> -->
              </div>
              <span v-if="scope.row.category === 1">
                手动创建
              </span>
            </template>
          </el-table-column>
          <el-table-column
            prop="task_num"
            min-width="130"
            label="包含子任务数"
          ></el-table-column>
          <el-table-column label="创建时间/编辑时间" width="200">
            <template slot-scope="scope">
              <div>
                {{ scope.row.created_at | formatDate("yyyy-mm-dd hh:mm:ss") }}~
              </div>
              <div>
                {{ scope.row.updated_at | formatDate("yyyy-mm-dd hh:mm:ss") }}
              </div>
            </template>
          </el-table-column>
          <el-table-column label="任务计划有效期" width="200">
            <template slot-scope="scope">
              <div>
                {{ scope.row.start_time | formatDate("yyyy-mm-dd hh:mm:ss") }}~
              </div>
              <div>
                {{ scope.row.end_time | formatDate("yyyy-mm-dd hh:mm:ss") }}
              </div>
            </template>
          </el-table-column>
          <el-table-column
            prop="create_name"
            label="创建人"
            min-width="220"
          ></el-table-column>
          <el-table-column label="是否开启" width="100">
            <template slot-scope="scope">
              <el-switch
                v-model="scope.row.is_open"
                @change="updataStatu(scope.row.id, scope.row.is_open)"
              >
              </el-switch>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="任务状态" min-width="110">
            <template slot-scope="scope">
              <span>{{
                scope.row.status === 0
                  ? "未开始"
                  : scope.row.status === 1
                  ? "进行中"
                  : "已结束"
              }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="zip" label="操作" width="310" fixed="right">
            <template slot-scope="scope">
              <span
                class="c-blue"
                v-if="scope.row.category !== 2"
                @click="copyEvn(scope.row)"
                >复制
              </span>
              <span
                v-if="scope.row.status === 0"
                class="c-blue"
                @click="editEvn(scope.row)"
              >
                编辑
              </span>
              <span
                v-if="scope.row.status === 0"
                class="c-blue"
                @click="deleteEvn(scope.row.id)"
                >删除
              </span>
              <span
                class="c-blue"
                v-if="
                  scope.row.status !== 0 ||
                    (scope.row.status === 0 && scope.row.task_num > 0)
                "
                @click="viewEvn(scope.row)"
                >查看任务</span
              >
              <span
                v-if="scope.row.status === 0 && scope.row.category !== 2"
                class="c-blue"
                @click="addChildEvn(scope.row)"
                >新增子任务
              </span>
              <span class="c-blue" @click="viewEvnData(scope.row)"
                >任务数据</span
              >
            </template>
          </el-table-column>
          <div slot="empty" class="table-empty-label">
            <img
              v-if="searchLoading"
              src="@/assets/images/empty-search.png"
              alt=""
            />
            <img v-else src="@/assets/images/empty-box.png" alt="" />
            <p>{{ searchLoading ? "没有找到相关任务计划" : "暂无任务计划" }}</p>
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
    <createPlan
      v-if="showCreateOrEditDialog"
      :selectedTask="selectedTask"
      :dialogType="dialogType"
      @channelDialog="channelDialogFn"
    >
    </createPlan>
    <upload-tag
      v-if="showuploadTagDialog"
      @closeUploadDialog="() => (showuploadTagDialog = false)"
    >
    </upload-tag>
    <child-task-type
      v-if="showChildTaskTypeDialog"
      @createTaskType="ChildTaskTypeSubmit"
    >
    </child-task-type>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import createPlan from "../components/create-plan.vue";
import uploadTag from "../components/uploadTag.vue";
import childTaskType from "../components/child-task-type.vue";
import Api from "@/api/wisdom/planTask";
@Component({
  components: {
    createPlan,
    uploadTag,
    childTaskType
  }
})
export default class extends Vue {
  private type_list: Array<any> = [
    { name: "任务ID", key: 0, keyword: "task_no" },
    { name: "任务名称", key: 1, keyword: "task_name" }
  ];
  private taskStatusList: Array<any> = [
    {
      label: "未开始",
      value: 0
    },
    {
      label: "进行中",
      value: 1
    },
    {
      label: "已结束",
      value: 2
    }
  ];
  private taskTagList = [
    {
      label: "手动创建",
      value: 1
    },
    // {
    //   label: "任务模板创建",
    //   value: 2
    // },
    // {
    //   label: "MA任务模板创建",
    //   value: 2
    // },
    {
      label: "MA任务创建",
      value: 3
    },
  ];
  private category = "";
  private searchLoading = false; // 是否正在搜索
  private gudie_type = 1; // 任务搜索  任务名称/任务ID
  private task_type = ""; // 任务类型
  private task_status = ""; // 任务状态
  private search_key = ""; // 搜索val
  private createTime: any = ""; // 创建时间
  private creatName = ""; // 创建人名称
  private list: Array<any> = [];
  private page_init = false;
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  private dialogType = ""; //弹窗类型
  private showCreateOrEditDialog = false;
  private selectedTask = null;
  private pickerOptions = {
    onPick: (obj: any) => {
      this.createTime = new Date(obj.minDate).getTime();
    },
    disabledDate: (time: any) => {
      if (this.createTime) {
        const day1 = 90 * 24 * 3600 * 1000;
        const maxTime = this.createTime + day1;
        const minTime = this.createTime - day1;
        return time.getTime() >= maxTime || time.getTime() <= minTime;
      }
    }
  };
  private showuploadTagDialog = false;
  private plan_template_id: any = "";
  // todo: 新增子任务类型
  private showChildTaskTypeDialog = false;
  private selectedTaskId: any = null;
  // 生命周期
  private mounted() {
    const plan_template_id = this.$route.query.plan_template_id;
    if (plan_template_id) {
      this.plan_template_id = plan_template_id;
    }
    this.getTaskList();
  }
  private changeSelect() {
    this.search_key = "";
  }
  // 新增任务计划
  private addTaskPlan() {
    this.dialogType = "1";
    this.showCreateOrEditDialog = true;
  }
  // 查询
  private searchEvn() {
    this.meta.page = 1;
    // 是否正在搜索
    this.searchLoading = true;
    this.getTaskList();
  }
  private changeNumberInput(val: any) {
    if (this.gudie_type == 1) {
      this.search_key = val;
      return;
    }
    const inputValue = String(val);
    const str: any = inputValue.replace(/\D/g, "").replace(/^0{1,}/g, "");
    this.search_key = str;
  }
  // 重置
  private reset() {
    // 初始化数据
    this.gudie_type = 1;
    this.task_type = "";
    this.task_status = "";
    this.search_key = "";
    this.createTime = "";
    this.creatName = "";
    this.category = "";
    this.meta.page = 1;
    this.meta.size = 10;
    // 是否正在搜索
    this.searchLoading = false;
    this.getTaskList();
  }
  private handleSizeChange(val: number) {
    this.meta.size = val;
    this.getTaskList();
  }
  private handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getTaskList();
  }
  private getTaskList() {
    const params: any = {
      page: this.meta.page,
      size: this.meta.size,
      search_type: this.gudie_type,
      status: this.task_status === "" ? -1 : this.task_status,
      start_time: (this.createTime && this.createTime[0] / 1000) || 0,
      end_time: (this.createTime && this.createTime[1] / 1000) || 0,
      create_name: this.creatName,
      task_label: this.category ? Number(this.category) : 0,
      plan_template_id: this.plan_template_id
        ? Number(this.plan_template_id)
        : 0
    };
    if (this.gudie_type === 1) {
      params.name = this.search_key;
    }
    if (this.gudie_type === 0 && this.search_key !== "") {
      params.id = Number(this.search_key);
    }
    Api.planTaskList(params).then((res: any) => {
      if (res.code === 200) {
        this.page_init = true;
        this.list = res.data.list;
        this.meta.total = res.data.total;
      }
    });
  }
  // 更新任务状态
  private updataStatu(id: number, isOpen: boolean) {
    const params = {
      id: id,
      is_open: isOpen
    };
    Api.planTaskStatusUpdata(params).then((res: any) => {
      if (res.code === 200) {
        this.$message.success("更新成功");
        this.getTaskList();
      }
    });
  }
  private copyEvn(row: any) {
    this.dialogType = "3";
    this.showCreateOrEditDialog = true;
    this.selectedTask = row;
  }
  private editEvn(row: any) {
    this.dialogType = "2";
    this.showCreateOrEditDialog = true;
    this.selectedTask = row;
  }
  private deleteEvn(val: number | string) {
    this.$confirm("你确认删除该任务吗？", "删除确认", {
      confirmButtonText: "确认",
      cancelButtonText: "取消"
    })
      .then(() => {
        Api.planTaskDelete({ id: val })
          .then((res: any) => {
            if (res.code == 200) {
              this.$message.success("删除成功");
              this.getTaskList();
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
  private viewEvn(val: any) {
    this.$router.push({
      name: "myapp_wisdom_task-center_plan-task_child-list",
      query: {
        selectedTask: JSON.stringify(val),
        type: "details"
      }
    });
  }
  private viewEvnData(val: any) {
    this.$router.push({
      name: "myapp_wisdom_task-center_plan-task_child-task-data",
      query: {
        selectedTask: JSON.stringify(val),
        type: "data"
      }
    });
  }
  private addChildEvn(row: any) {
    // todo: 新增子任务类型
    this.showChildTaskTypeDialog = true;
    this.selectedTaskId = row;
  }
  // 新增子任务类型提交
  ChildTaskTypeSubmit(typeObj: any) {
    console.log(typeObj);
    this.showChildTaskTypeDialog = false;
    if (typeObj.type || typeObj.typeName) {
      this.$router.push({
        name: "myapp_wisdom_task-center_plan-task_child-task-add",
        query: {
          status: "0",
          plan_id: this.selectedTaskId.id,
          st: this.selectedTaskId.start_time,
          et: this.selectedTaskId.end_time,
          type: typeObj.type
        }
      });
    }
  }
  private channelDialogFn(status: number) {
    this.showCreateOrEditDialog = false;
    this.selectedTask = null;
    if (status === 1) {
      this.getTaskList();
    }
  }
  private addUpload() {
    this.showuploadTagDialog = true;
  }
}
</script>
<style lang="scss" scoped>
.search-input {
  width: 280px !important;
}

.input-with-select {
  width: 110px;
}

.add-btn {
  color: #fff;
  background: $color-main;
}

.c-blue {
  cursor: pointer;
  color: #6f3112;
  margin-right: 8px;
}

.table-empty-label {
  text-align: center;
  height: 340px;
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;

  img {
    width: 124px;
  }

  p {
    line-height: 20px;
    margin-top: 14px;
  }
}
</style>
