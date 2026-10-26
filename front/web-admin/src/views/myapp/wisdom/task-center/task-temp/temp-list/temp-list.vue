<template>
  <div>
    <layout-list>
      <template slot="buttons">
        <el-button class="add-btn" @click="addTaskTemp">新增任务模板</el-button>
      </template>
      <template slot="search">
        <div>
          <search-item text="任务模板搜索：">
            <el-input
              :placeholder="placeholderName"
              class="search-input"
              size="medium"
              v-model.trim="search_key"
              maxlength="20"
              @input="changeNumberInput($event)"
              clearable
            >
              <el-select
                slot="prepend"
                v-model="search_type"
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
          <search-item text="创建人：">
            <cy-input
              size="small"
              width-size="small"
              placeholder="请输入创建人名称"
              v-model.trim="admin_name"
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
              :picker-options="pickerOptions"
            >
            </el-date-picker>
          </search-item>
          <search-item text="任务类型：">
            <el-select
              v-model="task_category"
              placeholder="请选择任务类型"
            >
              <el-option :value="-1" label="全部"></el-option>
              <el-option :value="0" label="沟通任务"></el-option>
              <el-option :value="1" label="朋友圈任务"></el-option>
            </el-select>
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
          <el-table-column label="任务模板ID" width="180">
            <template slot-scope="scope">
              <span class="c-blue plan-id-box" @click="viewEvn(scope.row)">
                {{ scope.row.id }}</span
              >
            </template>
          </el-table-column>
          <el-table-column label="任务模板名称" min-width="220">
            <template slot-scope="scope">
              <span class="c-blue plan-id-box" @click="viewEvn(scope.row)">
                {{ scope.row.template_name }}</span
              >
            </template>
          </el-table-column>
          <el-table-column min-width="130" label="任务类型">
            <template slot-scope="scope">
              <span v-if="scope.row.task_category === 0">沟通任务</span>
              <span v-else-if="scope.row.task_category === 1">朋友圈任务</span>
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="模板状态" min-width="110">
            <template slot-scope="scope">
              <span>{{ scope.row.status === 0 ? "禁用" : "启用" }}</span>
            </template>
          </el-table-column>
          <el-table-column label="创建/编辑时间" width="200">
            <template slot-scope="scope">
              <div>
                {{ scope.row.updated_at | formatDate("yyyy-mm-dd hh:mm: ss") }}
              </div>
            </template>
          </el-table-column>
          <el-table-column
            prop="admin_name"
            label="创建人"
            min-width="220"
          ></el-table-column>

          <el-table-column prop="zip" label="操作" width="310" fixed="right">
            <template slot-scope="scope">
              <span class="c-blue" @click="viewEvn(scope.row)">查看模板</span>
              <span class="c-blue" @click="editEvn(scope.row)">
                编辑
              </span>
              <span class="c-blue" @click="viewTask(scope.row)"
                >查看关联任务</span
              >
              <span class="c-blue" @click="updataStatus(scope.row)">
                {{ scope.row.status === 0 ? "启用模板" : "禁用模板" }}
              </span>
            </template>
          </el-table-column>
          <div slot="empty" class="table-empty-label">
            <img
              v-if="searchLoading"
              src="@/assets/images/empty-search.png"
              alt=""
            />
            <img v-else src="@/assets/images/empty-box.png" alt="" />
            <p>{{ searchLoading ? "没有找到相关任务模板" : "暂无任务模板" }}</p>
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
import Api from "@/api/wisdom/temp";
@Component({
  components: {}
})
export default class extends Vue {
  [x: string]: any;
  type_list: Array<any> = [
    { name: "模板ID", key: 2 },
    { name: "模板名称", key: 1 }
  ];
  searchLoading = false; // 是否正在搜索
  search_type = 1; // 任务搜索  任务名称/任务ID
  placeholderName = "请输入模板名称"; // 搜索框placeholder
  task_status = ""; // 任务状态
  search_key = ""; // 搜索val
  createTime: any = ""; // 创建时间
  admin_name = ""; // 创建人名称
  list: Array<any> = [];
  page_init = false;
  meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  dialogType = ""; //弹窗类型
  showCreateOrEditDialog = false;
  selectedTask = null; // 选中的任务
  task_category = -1; // 任务类型
  pickerOptions = {
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
  showuploadTagDialog = false;
  // 生命周期
  mounted() {
    this.getTaskList();
  }
  changeSelect() {
    this.search_key = "";
    this.placeholderName =
      this.search_type == 1 ? "请输入模板名称" : "请输入模板ID";
  }
  // 新增任务模板
  addTaskTemp() {
    this.$router.push({
      name: "myapp_wisdom_task-center_task-temp_create-temp"
    });
  }
  // 查询
  searchEvn() {
    this.meta.page = 1;
    // 是否正在搜索
    this.searchLoading = true;
    this.getTaskList();
  }
  changeNumberInput(val: any) {
    if (this.search_type == 1) {
      this.search_key = val;
      return;
    }
    const inputValue = String(val);
    const str: any = inputValue.replace(/\D/g, "").replace(/^0{1,}/g, "");
    this.search_key = str;
  }
  // 重置
  reset() {
    // 初始化数据
    this.search_type = 1;
    this.search_key = "";
    this.createTime = "";
    this.admin_name = "";
    this.meta.page = 1;
    this.meta.size = 10;
    this.placeholderName =
      this.search_type == 1 ? "请输入模板名称" : "请输入模板ID";
    this.task_category = -1; // 任务类型
    // 是否正在搜索
    this.searchLoading = false;
    this.getTaskList();
  }
  handleSizeChange(val: number) {
    this.meta.size = val;
    this.getTaskList();
  }
  handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getTaskList();
  }
  getTaskList() {
    const params: any = {
      page: this.meta.page,
      size: this.meta.size,
      search_type: this.search_type,
      start_time: (this.createTime && this.createTime[0] / 1000) || 0,
      end_time: (this.createTime && this.createTime[1] / 1000) || 0,
      admin_name: this.admin_name,
      task_category: this.task_category
    };
    if (this.search_type === 1) {
      params.search_name = this.search_key;
    }
    if (this.search_type === 2 && this.search_key !== "") {
      params.search_id = Number(this.search_key);
    }
    Api.getTaskTempList(params).then((res: any) => {
      if (res.code === 200) {
        this.page_init = true;
        this.list = res.data.list;
        this.meta.total = res.data.total;
      }
    });
  }
  // 更新任务状态
  updataStatus(row: any) {
    const params = {
      id: row.id,
      status: row.status === 0 ? 1 : 0
    };
    Api.changeTaskTempStatus(params).then((res: any) => {
      if (res.code === 200) {
        this.$message.success("更新成功");
        this.getTaskList();
      } else {
        this.$message.error(res.message);
      }
    });
  }
  /**
   * @description: 编辑任务模板
   * @param {*} row
   */

  editEvn(row: any) {
    this.$router.push({
      name: "myapp_wisdom_task-center_task-temp_create-temp",
      query: {
        id: row.id,
        is_edit: "1"
      }
    });
  }
  /**
   * @description: 模板详情
   * @param {*} val
   */

  viewEvn(row: any) {
    this.$router.push({
      name: "myapp_wisdom_task-center_task-temp_temp-info",
      query: {
        id: row.id
      }
    });
  }
  /**
   * @description: 任务计划列表
   * @param {*} row
   */

  viewTask(row: any) {
    this.$router.push({
      name: "myapp_wisdom_task-center_plan-task_plan-list",
      query: {
        plan_template_id: row.id
      }
    });
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
