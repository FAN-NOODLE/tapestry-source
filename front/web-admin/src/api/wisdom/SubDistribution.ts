import http from "@/utils/axios";

export default {
  // 二级分销设置
  subDistributionSettingApi(data: any) {
    return http("/distribution-config/second-save", {
      method: "POST",
      data
    });
  },
  // 二级分销查看
  subDistributionCheckApi() {
    return http("/distribution-config/second-detail", {
      method: "POST"
    });
  }
};
