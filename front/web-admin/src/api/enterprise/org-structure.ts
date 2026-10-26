import http from "@/utils/axios";

export default {
  // 获取数据架构tree
  getDashboardTree(data: any) {
    return http("/dashboard/area/tree", {
      method: "POST",
      data: data
      // data: { ...data, ...{ apiVersion: "testddd" } }
    });
  },
  // 创建数据架构节点
  createDashboardTree(data: any) {
    return http("/dashboard/area/save", {
      method: "POST",
      data: data
      // data: { ...data, ...{ apiVersion: "testddd" } }
    });
  },
  // 修改数据架构节点
  updateDashboardTree(data: any) {
    return http("/dashboard/area/update", {
      method: "POST",
      data: data
      // data: { ...data, ...{ apiVersion: "testddd" } }
    });
  },
  // 移除数据架构节点
  delDashboardTree(data: any) {
    return http("/dashboard/area/remove", {
      method: "POST",
      data: data
      // data: { ...data, ...{ apiVersion: "testddd" } }
    });
  },
  // 待分配门店列表
  getDashboardStoreList(data: any) {
    return http("/dashboard/area/not-relation-store-list", {
      method: "POST",
      data: data
      // data: { ...data, ...{ apiVersion: "testddd" } }
    });
  }
};
