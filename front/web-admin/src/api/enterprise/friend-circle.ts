import http from "@/utils/axios";

export default {
  // 编辑朋友圈形象
  momentAppearanceEditApi(data: any) {
    return http("moment/appearance-edit", {
      method: "POST",
      data: data
    });
  },
  // 获取朋友圈形象
  momentAppearanceGetApi(data: any) {
    return http("moment/appearance-get-admin", {
      method: "POST",
      data: data
    });
  },
  // 朋友圈新增 - 手动发布
  momentTopicCreateApi(data: any) {
    return http("moment/topic-create", {
      method: "POST",
      data: data
    });
  },
  // 朋友圈编辑
  momentTopicEditApi(data: any) {
    return http("moment/topic-edit", {
      method: "POST",
      data: data
    });
  },
  // 朋友圈展示
  momentTopicDisplayApi(data: any) {
    return http("moment/topic-display", {
      method: "POST",
      data: data
    });
  },
  // 朋友圈删除
  momentTopicDeleteApi(data: any) {
    return http("moment/topic-delete", {
      method: "POST",
      data: data
    });
  },
  //朋友圈按内容列表展示
  momentContentListApi(data: any) {
    return http("moment/content-list", {
      method: "POST",
      data: data
    });
  },
  //朋友圈按员工列表展示
  momentMemberListApi(data: any) {
    return http("moment/member-list", {
      method: "POST",
      data: data
    });
  },
  //上传图片
  momentUploadApi(data: any) {
    return http("file/upload", {
      method: "POST",
      data: data
    });
  },
  //朋友圈一个员工的朋友圈展示
  momentMemberTopicApi(data: any) {
    return http("moment/member-topic", {
      method: "POST",
      data: data
    });
  },
  // 删除一个员工的单条朋友圈
  memberMomentDeleteApi(data: any) {
    return http("moment/topic-member-delete", {
      method: "POST",
      data: data
    });
  },
  // 编辑一个员工的单条朋友圈
  memberMomentEditApi(data: any) {
    return http("moment/topic-member-edit", {
      method: "POST",
      data: data
    });
  }
};
