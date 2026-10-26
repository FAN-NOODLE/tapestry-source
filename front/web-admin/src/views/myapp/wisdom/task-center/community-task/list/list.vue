<template>
  <div class="list">
    <layout-list>
      <template slot="buttons">
        <el-button
          class="add-btn"
          @click="addTask"
        >新增任务</el-button>
      </template>
      <template slot="search">
        <div>
          <search-item text="任务搜索：">
            <el-input
              :placeholder="`请输入${type_list[gudie_type].name}`"
              class="search-input"
              size="medium"
              v-model.trim="search_key"
              clearable
            >
              <el-select
                slot="prepend"
                v-model="gudie_type"
                class="input-with-select"
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
          <search-item text="任务类型：">
            <cy-select
              v-model="task_type"
              size="medium"
              :optionsList="taskTypeList"
              :options="{ label: 'label', value: 'value' }"
            >
            </cy-select>
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
          <search-item text="创建时间：">
            <el-date-picker
              v-model="createTime"
              type="daterange"
              range-separator="-"
              start-placeholder="开始时间"
              end-placeholder="结束时间"
            >
            </el-date-picker>
          </search-item>
          <search-item text="编辑时间：">
            <el-date-picker
              v-model="editTime"
              type="daterange"
              range-separator="-"
              start-placeholder="开始时间"
              end-placeholder="结束时间"
            >
            </el-date-picker>
          </search-item>
          <search-item text="创建人：">
            <cy-input
              size="small"
              width-size="small"
              placeholder="请输入创建人名称"
              v-model.trim="creatName"
              prefix-icon="el-icon-search"
              clearable
            ></cy-input>
          </search-item>
          <div>
            <el-button
              class=""
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
          </div>
        </div>
      </template>
      <template slot="table">
        <el-table
          stripe
          :data="list"
          class="mar-t-20"
          v-if="page_init"
        >
          <el-table-column
            label="任务ID"
            width="150"
            fixed="left"
          >
            <template slot-scope="scope">
              <span
                class="c-blue plan-id-box"
                @click="checkTaskDetail(scope.row.task_id, scope.row.status)"
              >{{ scope.row.task_sn }}</span>
            </template>
          </el-table-column>
          <el-table-column
            label="任务名称"
            min-width="220"
          >
            <template slot-scope="scope">
              <el-popover
                placement="top-start"
                width="220"
                trigger="hover"
                :disabled="scope.row.task_title.length <= 20"
                :content="scope.row.task_title"
              >
                <div
                  :class="
                    scope.row.task_title.length > 20 ? 'work-name-wrap' : ''
                  "
                  slot="reference"
                >
                  {{ scope.row.task_title }}
                </div>
              </el-popover>
            </template>
          </el-table-column>
          <!-- 可见范围 -->
          <el-table-column
            label="可见范围"
            width="250"
          >
            <template slot-scope="scope">
              <span v-if="scope.row.object_array.length <= 1">
                {{
                  scope.row.object_array.length > 0
                    ? scope.row.object_array[0].name
                    : ""
                }}
              </span>
              <el-popover
                v-if="scope.row.object_array.length > 1"
                popper-class="marketing-calendar-popover-wrap"
                placement="top-start"
                width="220"
                trigger="hover"
              >
                <div class="object-array-content">
                  {{ scope.row.object_array_item }}
                </div>
                <div slot="reference">
                  <span>
                    {{
                      scope.row.object_array.length > 0
                        ? scope.row.object_array[0].name
                        : ""
                    }}等{{ scope.row.object_array.length }}个导购
                  </span>
                </div>
              </el-popover>
            </template>
          </el-table-column>
          <el-table-column
            prop="task_type"
            min-width="110"
            label="任务类型"
          ></el-table-column>
          <el-table-column
            prop="plan_resource"
            label="包含素材"
            min-width="200"
          >
          </el-table-column>
          <el-table-column
            label="创建时间/编辑时间"
            width="220"
          >
            <template slot-scope="scope"><span>{{ scope.row.created_at | formatDate("yyyy-mm-dd") }}/{{
                  scope.row.updated_at | formatDate("yyyy-mm-dd")
                }}</span></template>
          </el-table-column>
          <el-table-column
            prop="operator_name"
            label="创建人"
            min-width="250"
          ></el-table-column>
          <el-table-column
            prop="task_status"
            label="任务状态"
            min-width="110"
          ></el-table-column>
          <el-table-column
            prop="zip"
            label="操作"
            width="200"
            fixed="right"
          >
            <template slot-scope="scope">
              <span
                class="c-blue operate-btn"
                v-if="scope.row.status == 1"
                @click="editEvn(scope.row.task_id)"
              >
                编辑
              </span>
              <span
                class="c-blue operate-btn"
                @click="copyEvn(scope.row.task_id)"
              >复制
              </span>
              <span
                class="c-red operate-btn"
                v-if="scope.row.status == 1"
                @click="deleteEvn(scope.row.task_id)"
              >删除
              </span>
            </template>
          </el-table-column>

          <!-- <div slot="empty">
            <layout-empty
              empty_type="tag"
              v-if="search_key || createTime || editTime || creatName"
              :message="'没有找到相关任务'"
            ></layout-empty>
            <layout-empty
              v-if="!search_key && !createTime && !editTime && !creatName"
              :message="'暂无任务'"
            ></layout-empty>
          </div> -->
          <div
            slot="empty"
            class="table-empty-label"
          >
            <img
              v-if="searchLoading"
              src="@/assets/images/empty-search.png"
              alt=""
            />
            <img
              v-else
              src="@/assets/images/empty-box.png"
              alt=""
            />
            <div>{{ searchLoading ? "没有找到相关任务" : "暂无任务" }}</div>
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
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { taskCenterList, taskCenterDelete } from "@/api/wisdom/taskCenter";
@Component
export default class extends Vue {
  private type_list: Array<any> = [
    { name: "任务ID", key: 0, keyword: "task_no" },
    { name: "任务名称", key: 1, keyword: "task_name" },
  ];
  private taskTypeList: Array<any> = [
    {
      label: "培训",
      value: 1,
    },
    {
      label: "邀约",
      value: 2,
    },
    {
      label: "群发",
      value: 3,
    },
  ];
  private taskStatusList: Array<any> = [
    {
      label: "未开始",
      value: 1,
    },
    {
      label: "进行中",
      value: 2,
    },
    {
      label: "已结束",
      value: 3,
    },
  ];
  private searchLoading = false; // 是否正在搜索
  private gudie_type = 1; // 任务搜索  任务名称/任务ID
  private task_type = ""; // 任务类型
  private task_status = ""; // 任务状态
  private search_key = ""; // 搜索val
  private list: Array<any> = [
    // {
    //   plan_sn: "JH20201030001", // 任务ID
    //   plan_title: "这里显示计划名称最多显示20个字", // 任务名称
    //   task_type: "培训", // 任务类型
    //   plan_resource: "图片", // 计划包含素材
    //   created_at: "1621331021", // 创建编辑时间
    //   updated_at: "1621331021", // 创建编辑时间
    //   operator_name: "创建人", // 创建人
    //   task_status: "未开始", // 任务状态
    //   object_sum: 1, // 可见范围数量
    //   object_type: 1, // 可见范围类型 1门店 2部门
    //   object_array: [{ worker_id: "11111", name: "部门或门店名称" }], // 部门或门店ID 部门或门店名称
    //   object_array_item: "部门或门店名称"
    // }
  ]; // 表格数据
  private editTime = ""; // 编辑时间
  private createTime = ""; // 创建时间
  private creatName = ""; // 创建人名称

  // 页面初始化设置
  private page_init = false;
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10, // 每页数量
  };

  // 查询
  private searchEvn() {
    this.meta.page = 1;
    // 是否正在搜索
    this.searchLoading = true;
    this.getTaskList();
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
    this.editTime = "";
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

  // 获取任务列表数据
  getTaskList() {
    const param = {
      page: this.meta.page,
      size: this.meta.size, //条数：默认10
      task_title: this.gudie_type == 1 ? this.search_key : "", // 任务名称
      task_sn: this.gudie_type == 0 ? this.search_key : "", //任务ID
      create_time: {
        create_start: this.createTime
          ? Date.parse(this.createTime[0]) / 1000
          : "",
        create_end: this.createTime ? Date.parse(this.createTime[1]) / 1e3 : "",
      },
      //  TODO 时间选择组件可以自定义返回时间格式
      edit_time: {
        edit_start: this.editTime ? Date.parse(this.editTime[0]) / 1000 : "",
        edit_end: this.editTime ? Date.parse(this.editTime[1]) / 1e3 : "",
      },
      type: this.task_type,
      status: this.task_status,
      operator_name: this.creatName, // 创建人
    };
    taskCenterList(param)
      .then((res: any) => {
        if (res.code == 200) {
          this.page_init = true;
          res.data.list.forEach((item: any) => {
            // 处理导购名称
            if (item.object_array && item.object_array.length > 0) {
              item.object_array_item = item.object_array
                .map((item: any) => {
                  return item.name;
                })
                .join("、");
            }
            // 素材
            let resource_text = "";
            if (item.plan_resource && item.plan_resource.length <= 0) {
              resource_text = "无素材";
            } else {
              resource_text =
                item.plan_resource[0] == 1
                  ? "图片"
                  : item.plan_resource[0] == 2
                  ? "小程序"
                  : "视频";
            }
            item.plan_resource = resource_text;
            // 任务类型
            const task_type =
              item.type == 1 ? "培训" : item.type == 2 ? "邀约" : "群发";
            item.task_type = task_type;
            // 任务状态
            const task_status =
              item.status == 1
                ? "未开始"
                : item.status == 2
                ? "进行中"
                : "已结束";
            item.task_status = task_status;
          });
          this.list = res.data.list;
          this.meta.total = res.data.total;
        }
      })
      .catch((err: any) => {
        console.log("获取任务列表失败", err);
      });
  }

  // 查看任务详情
  checkTaskDetail(val: any, taskStatus: any) {
    this.$router.push({
      name: "myapp_wisdom_task-center_community-task_task-detail",
      query: {
        status: "3",
        taskStatus,
        val,
      },
    });
  }

  // 新增任务
  private addTask() {
    this.$router.push({
      name: "myapp_wisdom_task-center_community-task_task-add",
      query: {
        status: "0",
      },
    });
  }

  // 编辑任务
  editEvn(val: any) {
    this.$router.push({
      name: "myapp_wisdom_task-center_community-task_task-edit",
      query: {
        status: "1",
        val,
      },
    });
  }

  // 复制任务
  copyEvn(val: any) {
    this.$router.push({
      name: "myapp_wisdom_task-center_community-task_task-copy",
      query: {
        status: "2",
        val,
      },
    });
  }

  // 删除任务
  deleteEvn(val: number | string) {
    this.$confirm("你确认删除该任务吗？", "删除确认", {
      confirmButtonText: "确认",
      cancelButtonText: "取消",
    })
      .then(() => {
        const params = {
          task_id: val,
        };
        taskCenterDelete(params)
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

  // 生命周期
  private mounted() {
    this.getTaskList();
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
.operate-btn {
  cursor: pointer;
  margin-right: 20px;
}
.plan-id-box {
  cursor: pointer;
}

.c-blue {
  color: #6f3112;
}
.c-red {
  color: #ff0051;
}
.m-t-20 {
  margin-top: 20px;
}
.m-b-30 {
  margin-bottom: 30px;
}
.table-empty-label {
  text-align: center;
  height: 340px;
  display: flex;
  justify-content: center;
  flex-direction: column;
  align-items: center;
  img {
    width: 124px;
  }
  div {
    line-height: 20px;
    margin-top: 14px;
  }
}
</style>
