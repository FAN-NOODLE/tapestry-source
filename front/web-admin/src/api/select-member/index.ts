import http from "@/utils/axios";

export default {
  // 获取树形数据
  getDepartmentTreeApi(data) {
    return http("/department/tree-list", {
      method: "POST",
      data: { ...data, apiVersion: "old" },
    });
  },
  treePeople(data: any) {
    // 获取同步状态
    return http("members/momentList", {
      method: "POST",
      data: { ...data,},
    });
  },
  // 部门列表
  depListApi(data: any) {
    return http("departments/moment-list", {
      method: "POST",
      data: { ...data},
    });
  },
  treePeopleNew(data: any) {
    // 新的获取人员接口
    return http("members/ba/list", {
      method: "POST",
      data: { ...data},
    });
  },
  dpMembers(data: any) {
    // 新的获取人员接口
    return http("members/dp-members", {
      method: "POST",
      data: { ...data},
    });
  },
  // 奥莱/零售列表
  getStoreDataList() {
    return http("departments/sa-store/lists", {
      method: "POST",
    });
  },
  // 根据部门id 查询人员
  getMemberFromId(data: { dp_ids: Array<number> }) {
    return http("department/member-lists-by-dp-ids", {
      method: "POST",
      data: { apiVersion: "old", ...data },
    });
  },

  // 获取数据架构tree
  getDashboardTreeApi(data: any) {
    return http("/dashboard/area/tree", {
      method: "POST",
      data: data
      // data: { ...data, ...{ apiVersion: "testddd" } }
    });
  },
};
