<template>
  <div>
    <el-dialog
      :title="dialogTitle[dialogType]"
      :visible="true"
      @close="closeDialogFn"
      :close-on-click-modal="false"
    >
      <div class="section">
        <el-form
          ref="formObject"
          :rules="rules"
          :model="formObject"
          label-width="120px"
          label-position="right"
          size="medium"
        >
          <el-form-item label="任务计划名称" prop="plan_name">
            <el-input
              v-model="formObject.plan_name"
              maxlength="20"
              minlength="1"
              placeholder="请输入任务计划名称"
              :show-word-limit="true"
              size="medium"
              clearable
            ></el-input>
          </el-form-item>
          <el-form-item label="任务计划时间" prop="taskPlanTime">
            <el-date-picker
              v-model="formObject.taskPlanTime"
              type="datetimerange"
              format="yyyy-MM-dd HH:mm:ss"
              value-format="timestamp"
              range-separator="-"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              :default-time="['00:00:00', '23:59:59']"
              :picker-options="pickerOptions"
            ></el-date-picker>
          </el-form-item>
          <el-form-item label="任务计划描述" prop="description" required>
            <el-input
              ref="activityInfo"
              type="textarea"
              placeholder="请输入任务计划描述"
              v-model="formObject.description"
              :maxlength="200"
              resize="none"
              :showWordLimit="true"
            >
            </el-input>
          </el-form-item>
          <el-form-item label="是否开启" prop="whetherToStart">
            <el-switch v-model="formObject.whether_open"> </el-switch>
            <span class="desc">
              开启，任务计划在前端显示；关闭，任务计划在前端不显示
            </span>
          </el-form-item>
        </el-form>
      </div>
      <span slot="footer" class="dialog-footer">
        <el-button size="small" @click="closeDialogFn">取 消</el-button>
        <el-button size="small" type="primary" @click="submitFn"
          >保 存</el-button
        >
      </span>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop,Watch} from "vue-property-decorator";
import Api from "@/api/wisdom/planTask";
import PerModule from "@/store/modules/permission";
@Component
export default class extends Vue {
  @Prop({ default: () => null }) selectedTask?: any;
  @Prop({ default: () => "1" }) dialogType?: string;
  @Watch("dialogType", { immediate: true, deep: true })
  onDialogTypeChange(val: any) {
    if (val === "2" || val === "3") {
      this.getPlanTaskDetail();
    }
  }
  /************data**************/
  get ep() {
    return PerModule.enterprise
  }
  private dialogTitle = {
    "1": "新增任务计划",
    "2": "编辑任务计划",
    "3": "复制任务计划",
  };
  private formObject: any = {
    plan_name: "",
    taskPlanTime: "",
    description: "",
    whether_open: false,
  };
  private validateMap = {
    plan_name(rule: any, value: any, callback: Function) {
      if (!value) {
        callback(new Error("请输入任务计划名称"));
      } else if (value.length > 20) {
        callback(new Error("任务计划名称最多20个字"));
      } else {
        callback();
      }
    },
    taskPlanTime(rule: any, value: any, callback: Function) {
      if (!value) {
        callback(new Error("请选择任务计划时间"));
      } else {
        callback();
      }
    },
    description(rule: any, value: any, callback: Function) {
      if (!value) {
        callback(new Error("请输入任务计划描述"));
      } else {
        callback();
      }
    },
  };
  private rules = {
    plan_name: {
      validator: this.validateMap.plan_name,
      trigger: ["blur"],
    },
    taskPlanTime: {
      validator: this.validateMap.taskPlanTime,
      trigger: [ "change"],
    },
    description: {
      validator: this.validateMap.description,
      trigger: ["blur"],
    },
  };
  private pickerMinDate = 0; //第一次选中的时间
  private pickerOptions = {
    onPick: (obj: any) => {
      this.pickerMinDate = new Date(obj.minDate).getTime();
    },
    disabledDate: (time: any) => {
      if (this.pickerMinDate) {
        const day1 = 90 * 24 * 3600 * 1000;
        const maxTime = this.pickerMinDate + day1;
        const minTime = this.pickerMinDate - day1;
        return (
          time.getTime() >= maxTime ||
          time.getTime() <= minTime ||
          time.getTime() < Date.now() - 8.64e7
        );
      } else {
        return time.getTime() < Date.now() - 8.64e7;
      }
    },
  };
  // 关闭弹框
  private closeDialogFn() {
    this.$emit("channelDialog", 0);
  }
  // 获取任务详情
  getPlanTaskDetail() {
    this.formObject = {
      plan_name: this.selectedTask.name,
      taskPlanTime: 0,
      description: this.selectedTask.remark,
      whether_open: this.selectedTask.is_open,
    };
    if (this.dialogType === "2") {
      this.formObject.taskPlanTime = [
        this.selectedTask.start_time * 1000,
        this.selectedTask.end_time * 1000,
      ];
    } else if (this.dialogType === "3") {
      const planName = this.selectedTask.name + "-副本";
      this.formObject.plan_name = planName;
      this.$forceUpdate();
    }
  }
  // 保存
  private submitFn() {
    (this.$refs["formObject"] as any).validate(async (valid: boolean) => {
      if (valid) {
        if (this.formObject.plan_name.length > 20) {
          this.$message.error("请注意任务计划名称长度");
          return;
        }
        // 任务开始时间校验
        if (this.formObject.taskPlanTime[0] < new Date().getTime()) {
          this.$message.error("任务计划开始时间应大于当前时间");
          return;
        }

        // 校验拦截
        // if(this.ep.isCoach && this.dialogType === '3') {
        //   console.log(this.selectedTask);
        //   // 判断是否需要提示
        //   const res = await Api.canCopyApi({plan_id: this.selectedTask.id})
        //   console.log(res);
        //   if(res.data.show) {
        //     await this.$confirm('您当前复制的任务中，存在历史不可复制任务类型的子任务。为不影响新复制任务正常下发，我们将为您自动过滤掉不可复制的子任务，确认继续？', '温馨提示', {
        //       confirmButtonText: "确定",
        //       cancelButtonText: "取消",
        //       type: "warning"
        //     })
        //   }
        // }
        
        const params: any = {
          plan_name: this.formObject.plan_name,
          start_time:Math.floor(this.formObject.taskPlanTime[0] / 1000),
          end_time:Math.floor(this.formObject.taskPlanTime[1] / 1000),
          description: this.formObject.description,
          whether_open: this.formObject.whether_open,
        };
        if (this.dialogType != "1") {
          params.category = this.selectedTask.category;
        }
        const userInfo = JSON.parse(localStorage.getItem("USERINFO") || "{}"); //用户信息
        if (this.dialogType === "1") {
          params.operation_user_id = userInfo.admin_id;
          this.planTaskAdd(params);
        } else if (this.dialogType === "2") {
          params.id = this.selectedTask.id;
          params.created_at = this.selectedTask.created_at;
          params.operation_user_id = this.selectedTask.operation_user_id;
          params.status = this.selectedTask.status;
          params.task_num = this.selectedTask.task_num;
          this.planTaskEdit(params);
        } else if (this.dialogType === "3") {
          params.id = this.selectedTask.id;
          params.task_num = this.selectedTask.task_num;
          params.operation_user_id = userInfo.admin_id;
          this.planTaskCopy(params);
        }
      }
    });
  }
  // 新增
  private planTaskAdd(params: any) {
    Api.planTaskAdd(params).then((res: any) => {
      if (res.code === 200) {
        this.$message.success("新增成功");
        this.$emit("channelDialog", 1);
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 编辑保存
  private planTaskEdit(params: any) {
    Api.planTaskEdit(params).then((res: any) => {
      if (res.code === 200) {
        this.$message.success("更新成功");
        this.$emit("channelDialog", 1);
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 复制
  private planTaskCopy(params: any) {
    Api.planTaskCopy(params).then((res: any) => {
      if (res.code === 200) {
        this.$message.success("复制成功");
        this.$emit("channelDialog", 1);
      } else {
        this.$message.error(res.message);
      }
    });
  }
}
</script>
<style lang="scss" scoped>
.section {
  /deep/ .el-input--medium {
    width: 400px;
    .el-textarea__inner {
      height: 100px;
    }
  }
  .desc {
    margin-left: 14px;
    color: #999999;
  }
}
</style>
