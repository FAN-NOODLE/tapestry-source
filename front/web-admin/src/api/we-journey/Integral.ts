import http from "@/utils/axios";

/** 积分规则列表 */
export const getIntegralListApi = (data: any) => {
  return http("excitation/integral-rule/list", {
    method: "POST",
    data: { ...data },
  });
};

/** 积分规则详情 */
export const getIntegralDetailApi = (data: any) => {
  return http("excitation/integral-rule/detail", {
    method: "POST",

    data: {
      ...data,
    },
  });
};

/** 积分规则详情 */
export const getIntegralDeleteApi = (data: any) => {
  return http("excitation/integral-rule/destroy", {
    method: "POST",
    data: {
      ...data,
    },
  });
};

/** 积分规则新增 */
export const getIntegralAddApi = (data: any) => {
  return http("excitation/integral-rule/create", {
    method: "POST",
    data: {
      ...data,
    },
  });
};

/** 积分规则编辑 */
export const getIntegralEditApi = (data: any) => {
  return http("excitation/integral-rule/update", {
    method: "POST",
    data: {
      ...data,
    },
  });
};

/** 导购历史积分列表 */
export const getPointsListApi = (data: any) => {
  return http("excitation/points/record-list", {
    method: "POST",
    data: { ...data },
  });
};

/** 关联图谱规则详情 */
export const getRelationDetailApi = (data: any) => {
  return http("excitation/incentive/relation/detail", {
    method: "POST",
    data: { ...data },
  });
};

/** 关系能力图谱-新增或修改指标配置 */
export const getRelationSaveApi = (data: any) => {
  return http("excitation/incentive/relation/save", {
    method: "POST",
    data: { ...data },
  });
};

/** 积分规则维度配置 */
export const getIncentiveConfigApi = () => {
  return http("excitation/incentive/config", {
    method: "POST",
  });
};

/** 积分兑领数据导出 */
export const downloadPointApi = (data: any) => {
  return http("excitation/points/download-record", {
    method: "POST",
    data: { ...data },
  });
};

/** 下载手动添加积分模板 */
export const downloadTemplateApi = () => {
  return http("excitation/download-manual-points-template", {
    method: "POST",
    responseType: "blob",
    data:{ apiVersion: "old" },
  });
};

/** 导入手动积分模板 */
export const importPointApi = (data: any) => {
  return http("excitation/import-manual-points", {
    method: "POST",
    data: { ...data, ...{ apiVersion: "old" } },
  });
};

/** 获取部门列表 */
export const getThirdListApi = (data: any) => {
  return http("department/third-list", {
    method: "POST",
    data: { ...data, ...{ apiVersion: "old" } },
  });
};

/** 富文本自定义上传图片 */
export const uploadImageApi = (data: any) => {
  return http("common/upload-img", {
    method: "POST",
    data,
  });
};

/** 积分规则说明-详情 */
export const ruleDetailAPi = () => {
  return http("excitation/integral-rule/description-detail", {
    method: "POST",
  });
};

/** 积分规则说明-创建或修改 */
export const ruleSaveAPi = (data: any) => {
  return http("excitation/integral-rule/description-save", {
    method: "POST",
    data: { ...data },
  });
};
