export default {
  data() {
    return {
      no_choice: [
        "guide_store_member",
        "console_community_guide",
        "guide_community_guide",
        "console_enterprise_guide",
        // "guide_enterprise_guide",
        "guide_store_manager",
        "console_store_manager",
        "guide_enterprise_admin",
        "console_enterprise_admin"
      ], //不可选择的角色
      inner_no_choice: ["guide_enterprise_guide", "guide_enterprise_admin"],
      inner_no_choice_admin: [
        "guide_enterprise_guide",
        "guide_enterprise_admin"
      ]
    };
  }
};
