export default {
  data() {
    return {
      //内置的角色,不可删除编辑
      default_roles: [
        "super_admin",
        "area",
        "store_manager",
        "store_member",
        "guide_store_member",
        "guide_store_manager",
        "guide_enterprise_guide",
        "guide_community_guide",
        "guide_enterprise_admin",
        "console_admin_role" //内置角色管理员
      ],
      // 管理后台角色,可编辑不可删除
      console_default_roles: [
        "console_store_manager", //管理后台店长
        "console_enterprise_guide", //管理后台企业导购
        "console_community_guide", //管理后台社群导购
        "console_enterprise_admin"
      ]
    };
  }
};
