import http from "@/utils/axios";
export default {
  //贺卡列表
  templateCard(data: any) {
    return http("marketing-strategies/template/list", {
      method: "POST",
      data: data
    });
  },
  //添加/编辑贺卡模板
  templateInfo(data: any) {
    return http("marketing-strategies/template/info", {
      method: "POST",
      data: data
    });
  },
  //检查贺卡模板是否在使用
  templateCheck(data: any) {
    return http("marketing-strategies/template/check", {
      method: "POST",
      data: data
    });
  },
  //删除贺卡模板
  templateDelete(data: any) {
    return http("marketing-strategies/template/delete", {
      method: "POST",
      data: data
    });
  },
  //获取贺卡详情接口
  templateDetail(data: any) {
    return http("marketing-strategies/template/detail", {
      method: "POST",
      data: data
    });
  }
};
