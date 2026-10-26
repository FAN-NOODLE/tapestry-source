import http from "@/utils/axios";
// 【自设任务】- 列表
export const getSowTaskListApi = (data: any) =>
  http("customer/sow-task/list", {
    method: "POST",
    data,
  });


// 【自设任务】- 创建
export const createSowTaskApi = (data: any) =>
  http("customer/sow-task/create", {
    method: "POST",
    data,
  });


// 【自设任务】- 删除
export const deleteSowTaskApi = (data: any) =>
  http("customer/sow-task/delete", {
    method: "POST",
    data,
  });


// 【自设任务】- 编辑
export const updateSowTaskApi = (data: any) =>
  http("customer/sow-task/update", {
    method: "POST",
    data,
  });


// 【自设任务】- 完成
export const finishedSowTaskApi = (data: any) =>
  http("customer/sow-task/finished", {
    method: "POST",
    data,
  });