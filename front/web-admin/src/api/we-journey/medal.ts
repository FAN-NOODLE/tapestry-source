import http from "@/utils/axios";

/** 获取奖章规则下拉列表 */
export const getRuleOptionsApi = (data: any) => {
  return http("excitation/medal/rule-list", {
    method: "POST",
    data: { ...data },
  });
};

/** 奖章规则列表 */
export const getMedalListApi = (data: any) => {
  return http("excitation/medal/list", {
    method: "POST",
    data: { ...data },
  });
};

/** 奖章规则删除 */
export const medalDeleteApi = (data: any) => {
  return http("excitation/medal/remove", {
    method: "POST",
    data: { ...data },
  });
};

/** 奖章规则禁用启用 */
export const medalStatusApi = (data: any) => {
  return http("excitation/medal/status", {
    method: "POST",
    data: { ...data },
  });
};

/** 获取奖章规则详情 */
export const medalDetailApi = (data: any) => {
  return http("excitation/medal/detail", {
    method: "POST",
    data: {
      ...data,
    },
  });
};

/** 奖章规则创建 */
export const medalAddApi = (data: any) => {
  return http("excitation/medal/create", {
    method: "POST",
    data: {
      ...data,
    },
  });
};

/** 奖章规则编辑 */
export const medalEditApi = (data: any) => {
  return http("excitation/medal/update", {
    method: "POST",
    data: {
      ...data,
    },
  });
};

/** 奖章领取列表 */
export const getMedalRecordListApi = (data: any) => {
  return http("excitation/medal/record-list", {
    method: "POST",
    data: { ...data },
  });
};

/** 导出奖章领取记录 */
export const downloadMedalApi = (data: any) => {
  return http("excitation/download-medal-record", {
    method: "POST",
    responseType: "blob",
    data: { apiVersion: "old", ...data },
  });
};

/** 下载手动补发奖章模板 */
export const downloadReissueApi = () => {
  return http("excitation/download-reissue", {
    method: "POST",
    data: { apiVersion: "old" },
    responseType: "blob",
  });
};

/** 导入手动积分奖章 */
export const importMedalApi = (data: any) => {
  return http("excitation/read-reissue-data", {
    method: "POST",
    data: { apiVersion: "old", ...data },
  });
};

/** 下载人群包模板 */
export const downloadCrowPackApi = () => {
  return http("excitation/download-crowd-pack", {
    method: "POST",
    data: { apiVersion: "old" },
    responseType: "blob",
  });
};

/** 读取人群包 */
export const readCrowPackApi = (data: any) => {
  return http("excitation/read-crowd-pack-data", {
    method: "POST",
    data: { apiVersion: "old", ...data },
  });
};

/** 回收奖章 */
export const recordDeleteApi = (data: any) => {
  return http("excitation/medal/record-delete", {
    method: "POST",
    data: { ...data },
  });
};

/** 奖章领取补发 */
export const recordReissueApi = (data: any) => {
  return http("excitation/medal/record-reissue", {
    method: "POST",
    data: { ...data },
  });
};
