import http from "@/utils/axios";
interface TaskList {
  page: number;
  page_size: number;
  status: number;
  start_time?: number;
  end_time?: number;
}
interface TaskAdd {
  name: string;
  transfer_type: number;
  is_run?: number;
}
export default {
  // 根据ID，查询继承任务
  transferCustomerTask(data: any){
    return http("transfer_customer/task", {
      method: "POST",
      data
    });
  },
  // 获取好友继承管理列表
  getBuddyInheritanceListApi(data: TaskList) {
    return http("transfer_customer/tasks_list", {
      method: "POST",
      data
    });
  },
  // 新增继承任务
  addTask(data: TaskAdd) {
    return http("transfer_customer/tasks_add", {
      method: "POST",
      data
    });
  },
  // 更新继承任务
  tasksUpdate(data: any){
    return http("transfer_customer/tasks_update", {
      method: "POST",
      data
    });
  },
  // 删除继承任务
  tasksDelete(data: any){
    return http("transfer_customer/tasks_delete", {
      method: "POST",
      data
    });
  },
  // 继承任务报表
  tasksReport(data: any){
    return http("transfer_customer/tasks_report", {
      method: "POST",
      data
    });
  }
}