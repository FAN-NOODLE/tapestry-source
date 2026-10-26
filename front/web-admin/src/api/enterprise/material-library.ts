import http from "@/utils/axios";

export default {
  // 查询分组信息
  getGroupStoreData(data: any) {
    return http("/source/group-relation/list", {
      method: "POST",
      data: data,
      // data: { ...data, ...{ apiVersion: "testddd" } }
    });
  },

  saveGroupStore(data: any) {
    return http("/source/group-relation/save", {
      method: "POST",
      data: data,
    });
  },
  //
  getTreeList() {
    return http("department/tree-list", {
      method: "POST",
      data: { apiVersion: "old" },
    });
  },
  // 查询组名
  inquireGroupName(data: any) {
    return http("source/get/group_name", {
      method: "POST",
      data: data,
      // data: { ...data, ...{ apiVersion: "testddd" } }
    });
  },
  // 移动分组
  moveGroupListApi(data: any) {
    return http("source-group/move", {
      method: "POST",
      data: data,
      // data: { ...data, ...{ apiVersion: "testddd" } }
    });
  },
  // 获取分组列表
  getSourceGroupListApi(data: any) {
    return http("source-group/list", {
      method: "POST",
      data: data,
      // data: { ...data, ...{ apiVersion: "testddd" } }
    });
  },
  // 分组列表
  sourceGroupListApi(data: any) {
    return http("source-group/list", {
      method: "POST",
      data: data,
    });
  },
  // 移动分组列表
  moveGroupApi(data: any) {
    return http("source/move-group", {
      method: "POST",
      data: data,
      // data: { ...data, ...{ apiVersion: "testddd" } }
    });
  },
  // 新增分组列表
  sourceGroupAddApi(data: any) {
    return http("source-group/add", {
      method: "POST",
      data: data,
      // data: { ...data, ...{ apiVersion: "testddd" } }
    });
  },
  // 编辑分组列表
  sourceGroupEditApi(data: any) {
    return http("source-group/edit", {
      method: "POST",
      data: data,
    });
  },
  // 删除分组列表
  sourceGroupDeleteApi(data: any) {
    return http("source-group/delete", {
      method: "POST",
      data: data,
      // data: { ...data, ...{ apiVersion: "testddd" } }
    });
  },
  // 各类素材详情
  // type 类型(1 文本 2 图片 3 链接 4 小程序 5 视频 6 文件)
  sourceDetail(data: any) {
    return http("source/source-detail", {
      method: "POST",
      data: data,
    });
  },
  // 各类素材编辑保存
  // source_type 类型(1 文本 2 图片 3 链接 4 小程序 5 视频 6 文件)
  sourceUpdate(data: any) {
    return http("source/source-update", {
      method: "POST",
      data: data,
    });
  },
  // 素材列表
  consoleList(data: any) {
    return http("source/console-list", {
      method: "POST",
      data: data,
    });
  },
  // 文本素材弹框列表
  textListApi(data: any) {
    return http("source/text-list", {
      method: "POST",
      data: data,
      // data: { ...data, ...{ apiVersion: "testddd" } }
    });
  },
  // 文本素材新增编辑
  textSaveApi(data: any) {
    return http("source/text-save", {
      method: "POST",
      data: data,
      // data: { ...data, ...{ apiVersion: "testddd" } }
    });
  },
  // 图片素材弹框列表
  imageListApi(data: any) {
    return http("source/image-list", {
      method: "POST",
      data: data,
    });
  },
  // 图片素材弹框列表
  imageAddApi(data: any) {
    return http("source/image-add", {
      method: "POST",
      data: data,
    });
  },
  // 小程序素材弹框列表
  miniprogramListApi(data: any) {
    return http("source/mini-program-list", {
      method: "POST",
      data: data,
      // data: { ...data, ...{ apiVersion: "testddd" } }
    });
  },
  // 小程序素材新增编辑
  miniProgramSaveApi(data: any) {
    return http("source/mini-program-save", {
      method: "POST",
      data: data,
    });
  },
  // 素材列表-批量启用/禁用素材
  updateStatus(data: any) {
    return http("source/update-status", {
      method: "POST",
      data: data,
    });
  },
  // 小程序素材新增编辑
  miniProgramEditApi(data: any) {
    return http("mini-program/edit", {
      method: "POST",
      data: data,
    });
  },
  // 小程序素材弹框列表
  miniListApi(data: any) {
    return http("mini-program/list", {
      method: "POST",
      data,
    });
  },
  // 小程序素材新增编辑
  miniProgramAddApi(data: any) {
    return http("/mini-program/add", {
      method: "POST",
      data: data,
    });
  },
  // 链接素材列表
  linkListApi(data: any) {
    return http("source/link-list", {
      method: "POST",
      data: data,
      // data: { ...data, ...{ apiVersion: "testddd" } }
    });
  },
  // 链接素材列表
  linkDetailApi(data: any) {
    return http("source/link-detail", {
      method: "POST",
      data: data,
    });
  },
  // 链接素材新增 source/link-save
  linkSaveApi(data: any) {
    return http("source/link-save", {
      method: "POST",
      data: data,
    });
  },
  // 文件上传
  fileUploadApi(data: any) {
    const form_data = new FormData();
    form_data.append("file", data.file);
    return http("file/upload", {
      method: "POST",
      data: form_data,
    });
  },
  // 视频素材列表
  videoListApi(data: any) {
    return http("source/video-list", {
      method: "POST",
      data: data,
      // data: { ...data, ...{ apiVersion: "testddd" } }
    });
  },
  // 视频素材列表
  videoSaveApi(data: any) {
    return http("source/video-save", {
      method: "POST",
      data: data,
    });
  },
  // 删除小程序
  miniProgramDelApi(data: any) {
    return http("mini-program/del", {
      method: "POST",
      data: data,
    });
  },
  //文件素材新增
  fileAddApi(data: any) {
    return http("source/file-add", {
      method: "POST",
      data: data,
    });
  },
  //文件素材新增
  fileListApi(data: any) {
    return http("source/file-list", {
      method: "POST",
      data: data,
    });
  },
  //素材删除
  materialDeleteApi(data: any) {
    return http("source/delete", {
      method: "POST",
      data: data,
    });
  },
  //小程序素材详情
  miniProgramDetailApi(data: any) {
    return http("source/mini-program-detail", {
      method: "POST",
      data: data,
    });
  },
  //素材标签-列表
  categoryList(data: any) {
    return http("source/category/list", {
      method: "POST",
      data: data,
    });
  },
  //素材标签-创建
  categoryCreate(data: any) {
    return http("source/category/create", {
      method: "POST",
      data: data,
    });
  },
  //素材标签-删除
  categoryDelete(data: any) {
    return http("source/category/delete", {
      method: "POST",
      data: data,
    });
  },
  //素材标签-编辑
  categoryUpdate(data: any) {
    return http("source/category/update", {
      method: "POST",
      data: data,
    });
  },
  //素材标签-获取管理素材ID集合
  categorySourceRelation(data: any) {
    return http("source/category/source_relation", {
      method: "POST",
      data: data,
    });
  },
  //各类素材-移除素材标签
  removeCategory(data: any) {
    return http("source/remove-category", {
      method: "POST",
      data: data,
    });
  },
  // 朋友圈素材新增
  friendCircleCreate(data: any) {
    return http("source/friend/create", {
      method: "POST",
      data: data,
    });
  },
  // 朋友圈编辑
  friendCircleUpdate(data: any) {
    return http("source/friend/update", {
      method: "POST",
      data: data,
    });
  },
  // 朋友圈文件校验
  getMediaIdApi(data: any) {
    return http("source/get-media-info", {
      method: "POST",
      data: data,
      hideLoading: true,
    } as any);
  },
  // 朋友圈图片压缩
  compressImage(data: any) {
    return http("source/compress-file", {
      method: "POST",
      data: data,
      hideLoading: true,
    } as any);
  },
  setTopApi(data: any) {
    return http("source/set-top", {
      method: "POST",
      data: data,
    });
  },
  updateIgnoreIndexApi(data: any) {
    return http("source-group/update-ignore-index", {
      method: "POST",
      data: data,
    });
  },
  cardListApi(data: any) {
    return http("marketing-strategies/template/list", {
      method: "POST",
      data: { ...data, ...{ apiVersion: "old" } },
    });
  },
  searchGoodsListApi(data: any) {
    return http("source/brand/search-list", {
      method: "POST",
      data: data,
    });
  },
  getTalkieImageSearchApi(data: any) {
    return http("source/brand/talkie-image-search", {
      method: "POST",
      data: data,
      hideLoading: true,
    });
  },
};
