<template>
  <div>
    <layout-list>
      <template slot="buttons">
        <el-button class="add_btn" type="primary" @click="addInheritTask">
          新增继承任务
        </el-button>
      </template>
      <template slot="search">
        <search-item text="任务时间：">
          <el-date-picker
            v-model="createTime"
            :default-time="['00:00:00', '23:59:59']"
            value-format="timestamp"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            clearable
          >
          </el-date-picker>
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
        <el-button size="small" type="primary" @click="searchEvn">
          查询
        </el-button>
        <el-button class="mar-l-20" size="small" @click="reset">
          重置
        </el-button>
      </template>
      <template slot="table">
        <el-table
          :data="list"
          stripe
          style="width: 100%"
          tooltip-effect="light"
        >
          <el-table-column
            prop="transfer_id"
            label="任务id"
            min-width="240"
          ></el-table-column>
          <el-table-column
            prop="name"
            label="任务名称"
            min-width="240"
          ></el-table-column>
          <el-table-column label="客户数" width="120">
            <template slot-scope="scope">
              <span>
                {{
                  scope.row.status === 1 && scope.row.customer_num === 0
                    ? "计数中..."
                    : scope.row.customer_num
                }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="任务创建时间" width="230">
            <template slot-scope="scope">
              <span>
                {{ scope.row.created_at | formatDate("yyyy-mm-dd hh:mm:ss") }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="任务状态" width="120">
            <template slot-scope="scope">
              {{
                scope.row.status === 1
                  ? "未开始"
                  : scope.row.status === 2
                  ? "进行中"
                  : scope.row.status === 3
                  ? "已结束"
                  : "已取消"
              }}
            </template>
          </el-table-column>
          <el-table-column prop="zip" label="操作" width="310" fixed="right">
            <template slot-scope="scope">
              <span
                class="c-blue"
                v-if="scope.row.status === 1"
                @click="goRun(scope.row)"
                >立即执行</span
              >
              <span class="c-blue" @click="check(scope.row)">查看 </span>
              <span class="c-blue" @click="taskReport(scope.row)"
                >任务报表
              </span>
              <span
                class="c-blue"
                v-if="scope.row.status !== 4"
                @click="taskCancel(scope.row)"
              >
                任务取消
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
            <p>{{ searchLoading ? "没有找到相关任务" : "暂无任务" }}</p>
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
    <el-dialog
      :title="dialogTitle"
      :visible="isShowDialog"
      @close="closeDialogFn"
      :close-on-click-modal="false"
    >
      <div class="dialog-content">
        <div class="row-item">
          数据条件： 已有绑定管理导购，且管理导购未与客户建立好友关系的客户
        </div>
        <div class="row-item input-line">
          任务名称：
          <cy-input
            v-if="activeStatus === 1 || activeStatus === 0"
            size="small"
            width-size="medium"
            placeholder="请输入任务名称"
            v-model="taskName"
            @input="accountInput"
            class="input-with-select"
            maxlength="40"
            show-word-limit
          ></cy-input>
          <span v-else style="max-width: 80%">{{ taskName }}</span>
        </div>
        <div class="row-item">客户范围： 全部客户</div>
        <div class="row-item">
          继承类型：
          <el-radio-group
            v-model="inheritType"
            :disabled="activeStatus !== 1 && activeStatus !== 0"
          >
            <el-radio :label="1">唯一继承</el-radio>
            <div class="text-des">
              建立与管理导购的唯一好友关系，解除客户与其他非管理导购的好友关系；
            </div>
            <el-radio :label="2">随机继承</el-radio>
            <div class="text-des">
              从客户已有好友中随机选择一人解除好友关系，并与管理导购建立新好友关系；
            </div>
          </el-radio-group>
        </div>
      </div>
      <span slot="footer" class="dialog-footer">
        <el-button
          v-if="activeStatus === 1 || activeStatus === 0"
          size="small"
          type="primary"
          @click="saveAndExecute"
        >
          保存并执行
        </el-button>
        <el-button
          v-if="activeStatus === 1 || activeStatus === 0"
          size="small"
          @click="saveTask"
        >
          保存任务
        </el-button>
        <el-button size="small" @click="cancelTak">取消</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Api from "@/api/enterprise/buddy-inheritance";
@Component
export default class extends Vue {
  // data
  private createTime: Array<number> | number = 0; // 创建时间
  private task_status = 0; // 任务状态
  private searchLoading = false; // 是否在搜索
  private page_init = false;
  private list: Array<any> = [];
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  private taskStatusList: Array<any> = [
    {
      label: "全部",
      value: 0
    },
    {
      label: "未开始",
      value: 1
    },
    {
      label: "进行中",
      value: 2
    },
    {
      label: "已结束",
      value: 3
    },
    {
      label: "已取消",
      value: 4
    }
  ];
  private taskName = "";
  private isShowDialog = false;
  private dialogTitle = "";
  private activeStatus = 0;
  private activeId = 0;
  private inheritType = 1;
  // 生命周期
  private mounted() {
    this.getList();
  }

  // 新增继承任务
  private addInheritTask() {
    this.dialogTitle = "新增好友关系继承";
    this.activeStatus = 0;
    this.activeId = 0;
    this.isShowDialog = true;
  }

  // 关闭弹窗
  private closeDialogFn() {
    this.isShowDialog = false;
    this.taskName = "";
  }

  // 搜索
  private searchEvn() {
    this.meta.page = 1;
    this.searchLoading = true;
    this.getList();
  }

  // 重置
  private reset() {
    this.createTime = 0;
    this.task_status = 0;
    this.getList();
  }

  private getList() {
    const params = {
      page: this.meta.page,
      page_size: this.meta.size,
      start_time: (this.createTime && this.createTime[0] / 1000) || undefined,
      end_time: (this.createTime && this.createTime[1] / 1000) || undefined,
      status: this.task_status
    };
    Api.getBuddyInheritanceListApi(params).then((res: any) => {
      if (res.code == 200) {
        this.page_init = true;
        this.list = res.data.data;
        this.meta.total = res.data.total;
      } else {
        this.$message.error(res.message);
      }
    });
  }

  private accountInput(val: string) {
    const str = val.replace(/[^\a-\z\A-\Z0-9\u4E00-\u9FA5]/g, "");
    this.$nextTick(() => {
      this.taskName = str;
    });
  }

  private handleSizeChange(val: number) {
    this.meta.size = val;
    this.getList();
  }
  private handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getList();
  }

  // 立即执行
  private goRun(row: any) {
    const params = {
      id: row.id,
      status: 2
    };
    Api.tasksUpdate(params).then((res: any) => {
      if (res.code == 200) {
        this.$message.success("任务开始执行");
        this.getList();
      } else {
        this.$message.error(res.message);
      }
    });
  }

  // 查看
  private check(row: any) {
    this.dialogTitle = "查看好友关系继承";
    this.activeStatus = row.status;
    this.taskName = row.name;
    this.inheritType = row.type;
    this.activeId = row.id;
    this.isShowDialog = true;
  }

  // 任务报表
  private taskReport(row: any) {
    this.$router.push({
      name: "myapp_enterprise_statistics_buddy-inheritance_task-report",
      query: {
        id: row.id
      }
    });
  }

  // 任务取消
  private taskCancel(row: any) {
    const params = {
      id: row.id,
      status: 4
    };
    Api.tasksUpdate(params).then((res: any) => {
      if (res.code == 200) {
        this.$message.success("任务已取消");
        this.getList();
      } else {
        this.$message.error(res.message);
      }
    });
  }

  // 保存并执行
  private saveAndExecute() {
    if (this.taskName === "") {
      this.$message.warning("请先输入任务名称");
      return;
    }
    if (this.activeStatus === 0) {
      const params = {
        name: this.taskName,
        transfer_type: this.inheritType,
        is_run: 1
      };
      Api.addTask(params).then((res: any) => {
        if (res.code == 200) {
          this.$message.success("新增成功");
          this.closeDialogFn();
          this.getList();
        } else {
          this.$message.error(res.message);
        }
      });
    } else if (this.activeStatus === 1) {
      const params = {
        id: this.activeId,
        name: this.taskName,
        transfer_type: this.inheritType,
        status: 2
      };
      Api.tasksUpdate(params).then((res: any) => {
        if (res.code == 200) {
          this.$message.success("保存成功");
          this.closeDialogFn();
          this.getList();
        } else {
          this.$message.error(res.message);
        }
      });
    }
  }

  // 保存任务
  private saveTask() {
    if (this.taskName === "") {
      this.$message.warning("请先输入任务名称");
      return;
    }
    if (this.activeStatus === 0) {
      const params = {
        name: this.taskName,
        transfer_type: this.inheritType
      };
      Api.addTask(params).then((res: any) => {
        if (res.code == 200) {
          this.$message.success("新增成功");
          this.closeDialogFn();
          this.getList();
        } else {
          this.$message.error(res.message);
        }
      });
    } else if (this.activeStatus === 1) {
      const params = {
        id: this.activeId,
        name: this.taskName,
        transfer_type: this.inheritType,
        status: 1
      };
      Api.tasksUpdate(params).then((res: any) => {
        if (res.code == 200) {
          this.$message.success("保存成功");
          this.closeDialogFn();
          this.getList();
        } else {
          this.$message.error(res.message);
        }
      });
    }
  }

  // 取消
  private cancelTak() {
    this.closeDialogFn();
  }
}
</script>
<style lang="scss" scoped>
/deep/ .search-content {
  padding-bottom: 0 !important;
}
.c-blue {
  cursor: pointer;
  color: $color-text-button;
  margin-right: 8px;
}
.dialog-content {
  border-top: 1px solid #d7d7d7;
  padding: 30px 30px 0 30px;
  .row-item {
    display: flex;
    margin-bottom: 28px;
  }
  .input-line {
    line-height: 40px;
  }
}
.el-dialog__footer {
  text-align: left;
}
.text-des {
  font-size: 12px;
  color: #8e8e8e;
  margin: 20px 0;
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
