import http from "@/utils/axios";
interface AddPlanTask {
  plan_name: string;
  plan_status: number | string;
  start_time: number | undefined;
  end_time: number | undefined;
  description: string;
  whether_open: boolean;
}
export default {
  // 任务计划-列表
  planTaskList(data: any) {
    return http("plan/list", {
      method: "POST",
      data: { ...data },
    });
  },
  // 任务计划状态更新
  planTaskStatusUpdata(data: any) {
    return http("plan/update/status", {
      method: "POST",
      data: { ...data },
    });
  },
  // 任务计划-新增
  planTaskAdd(data: AddPlanTask) {
    return http("plan/add", {
      method: "POST",
      data: { ...data },
    });
  },
  // 任务计划-编辑更新
  planTaskEdit(data: any) {
    return http("plan/edit", {
      method: "POST",
      data: { ...data },
    });
  },
  // 任务计划-删除
  planTaskDelete(data: any) {
    return http("plan/remove", {
      method: "POST",
      data: { ...data },
    });
  },
  // 任务计划-复制
  planTaskCopy(data: any) {
    return http("plan/copy", {
      method: "POST",
      data: { ...data },
    });
  },
  // 子任务列表
  childTaskList(data: any) {
    return http("task/list", {
      method: "POST",
      data: { ...data },
    });
  },
  // 子任务删除
  childTaskDelete(data: any) {
    return http("plan_tasks/delete", {
      method: "POST",
      data: { ...data },
    });
  },

  //子任务新增
  addChildTask(data: any) {
    return http("plan_tasks/add", {
      method: "POST",
      data: { ...data },
    });
  },
  // 子任务详情
  childTaskDetail(data: any) {
    return http("plan_tasks/info", {
      method: "POST",
      data: { ...data },
    });
  },
  // 子任务编辑修改
  editChildTask(data: any) {
    return http("plan_tasks/update", {
      method: "POST",
      data: { ...data },
    });
  },
  // 条件完成情况
  conditionComplete(data: any) {
    return http("plan_tasks/conditions/completion", {
      method: "POST",
      data: { ...data },
    });
  },
  // 导购完成情况
  userCompletion(data: any) {
    return http("plan_tasks/conditions/user_completion", {
      method: "POST",
      data: { ...data },
    });
  },
  // 主任务统计
  mainTask(data: any) {
    return http("plan/task/data/new-stat", {
      method: "POST",
      data: { ...data },
    });
  },
  // 主任务统计数据明细
  mainTaskData(data: any) {
    return http("plan/task/data/stat/new-detail", {
      method: "POST",
      data: { ...data },
    });
  },
  // 主任务统计数据明细导出
  mainTaskDataExport(data: any) {
    return http("plan/task/data/stat/new-detail/export", {
      method: "POST",
      data: { ...data },
      responseType: "blob",
    });
  },
  // 子任务社群任务顶部统计
  subTaskAssociation(data: any) {
    return http("plan/subtask/group_task", {
      method: "POST",
      data: { ...data },
    });
  },
  // 子任务社群任务列表导出
  subTaskAssociationExport(data: any) {
    return http("plan/subtask/group_task_export", {
      method: "POST",
      data: { ...data },
      responseType: "blob",
    });
  },
  // 子任务用户画像顶部统计
  subTaskPortrayal(data: any) {
    return http("plan/subtask/customer", {
      method: "POST",
      data: { ...data },
    });
  },
  // 子任务----明细数据导出
  subTaskPortrayalExport(data: any) {
    return http("plan/task-detail-export-new", {
      method: "POST",
      data: { ...data },
    });
  },
  // 子任务用户画像列表数据导出
  subTaskPortrayalListExport(data: any) {
    return http("plan/subtask/customer/list_export", {
      method: "POST",
      data: { ...data },
      responseType: "blob",
    });
  },
  // 子任务社群任务列表
  subTaskList(data: any) {
    return http("plan/subtask/group_task_list", {
      method: "POST",
      data: { ...data },
    });
  },
  // 子任务用户画像列表
  subPortrayalList(data: any) {
    return http("plan/subtask/customer_task_list", {
      method: "POST",
      data: { ...data },
    });
  },
  // 是否可以复制
  canCopyApi(data: any) {
    return http("plan/is-show", {
      method: "POST",
      data: { ...data },
    });
  },
  // 获取子任务名名称列表
  getPlanNameApi(data: any) {
    return http("plan/task/plan-name-list", {
      method: "POST",
      data: { ...data },
    });
  },
};
