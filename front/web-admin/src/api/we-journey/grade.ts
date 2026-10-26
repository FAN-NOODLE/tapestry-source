import http from "@/utils/axios";

/** 等级规则详情 */
export const getLevelDetailApi = (data: any) => {
  return http("excitation/level-rule/detail-list", {
    method: "POST",
    data: { ...data },
  });
};

/** 新增或修改等级规则 */
export const getLevelSaveApi = (data: any) => {
  return http("excitation/level-rule/save", {
    method: "POST",
    data: { ...data },
  });
};

/** 导购等级数据 */
export const getMemberListApi = (data: any) => {
  return http("excitation/member-data/list", {
    method: "POST",
    data: { ...data },
  });
};

/** 等级列表-options */
export const getLevelListApi = () => {
  return http("excitation/level-rule/list", {
    method: "POST",
  });
};

/** 导购等级数据列表导出 */
export const downloadTemplateApi = (data: any) => {
  return http("excitation/member-data/list-export", {
    method: "POST",
    responseType: "blob",
    data: { ...data },
  });
};

/** 获取附加条件下拉列表 */
export const getOptionApi = () => {
  return http("excitation/incentive/relation/detail-list", {
    method: "POST",
  });
};
