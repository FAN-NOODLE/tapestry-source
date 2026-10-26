import http from "@/utils/axios";

export default {
  syncStores(data: any = {}) {
    return http("/dashboard/hkmo-store/sync", { method: "POST", data });
  },
  getLastSyncTime(data: any = {}) {
    return http("/dashboard/hkmo-store/last-sync-time", { method: "POST", data });
  },
  getStoreList(data: any) {
    return http("/dashboard/hkmo-store/list", { method: "POST", data });
  },
  getStoreDetail(data: any) {
    return http("/dashboard/hkmo-store/detail", { method: "POST", data });
  },
  saveStore(data: any) {
    return http("/dashboard/hkmo-store/save", { method: "POST", data });
  },
  closeStore(data: any) {
    return http("/dashboard/hkmo-store/close", { method: "POST", data });
  },
  getAreaTree(data: any) {
    return http("/dashboard/hkmo-area/list", { method: "POST", data });
  },
  getPermissionList(data: any) {
    return http("/dashboard/hkmo-manager/auth-list", { method: "POST", data });
  },
  getPermissionDetail(data: any) {
    return http("/dashboard/hkmo-manager/auth-detail", { method: "POST", data });
  },
  savePermission(data: any) {
    return http("/dashboard/hkmo-manager/auth-save", { method: "POST", data });
  },
  deletePermission(data: any) {
    return http("/dashboard/hkmo-manager/auth-delete", { method: "POST", data });
  },
  getConfiguredMemberIds(data: any = {}) {
    return http("/dashboard/hkmo-manager/all-member-ids", { method: "POST", data });
  }
};
