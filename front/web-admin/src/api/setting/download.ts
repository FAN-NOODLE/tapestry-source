import http from "@/utils/axios";
export default {
  // 下载中心列表
  DownloadList(data: any) {
    return http("download-center/export-list", {
      method: "POST",
      data,
    });
  },
  // 下载文件
  DownloadFile(data: any) {
    return http("download-center/download", {
      method: "GET",
      params: data,
    });
  },
  // 删除文件
  downloadDel(data: any) {
    return http("download-center/delete", {
      method: "POST",
      data: data,
    });
  },
  // 重新生成 /gw-console/v1/download-center/retry
  reDownload(data: any) {
    return http("download-center/retry", {
      method: "POST",
      data: data,
    });
  },
  // 不走下载中心-下载
  downLoadApi(data: any) {
    return http("download-center/detail", {
      method: "GET",
      params: data,
    });
  },
};
