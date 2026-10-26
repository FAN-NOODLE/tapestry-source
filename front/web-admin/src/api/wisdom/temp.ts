import http from "@/utils/axios";

export default {
  // 任务模板列表
  getTaskTempList: (data: any) => {
    return http({
      url: "plan/template/lists",
      method: "post",
      data: { ...data },
    });
  },
  // 创建任务模板
  createTaskTemp: (data: any) => {
    return http({
      url: "plan/template/create",
      method: "post",
      data: { ...data },
    });
  },
  // 编辑任务模板
  editTaskTemp: (data: any) => {
    return http({
      url: "plan/template/update",
      method: "post",
      data: { ...data },
    });
  },
  // 任务模板详情
  getTaskTempDetail: (data: any) => {
    return http({
      url: "plan/template/details",
      method: "post",
      data: { ...data },
    });
  },
  // 任务模板启用禁用
  changeTaskTempStatus: (data: any) => {
    return http({
      url: "plan/template/status",
      method: "post",
      data: { ...data },
    });
  },
};
