<template>
  <div></div>
</template>
<script>
import Api from "@/api/public/public";
import UserModule from "@/store/modules/user";
import PerModule from "@/store/modules/permission";
import settingApis from "@/api/setting/setting";
const enterpriseMap = {
  // coach
  "tapestry-web-uat.coach.com.cn": 3363603667487424,
  "we.coach.com.cn": 3360610327843968,
  // ks
  "tapestry-web-uat.katespade.cn": 3376048056335040,
  "we.katespade.cn": 3383774945595776
  // sw
  // "tapestry-web-uat.stuartweitzman.cn": 3373738709328576,
  // "we.stuartweitzman.cn": 3374939108972352
};
export default {
  name: "Sso",
  data() {
    return {
      enterpriseId: enterpriseMap[window.location.host],
    };
  },
  mounted() {
    this.login();
  },
  methods: {
    /**
     * 登录
     */
    async login() {
      const code = this.$route.query.code;
      if (code) {
        // 登录
        const res = await Api.ssologin({
          enterprise_id: this.enterpriseId,
          code,
        })

        if(res.code !== 200) {
          this.$message.error(res.message || '登录失败')
          return
        }
        const { token, enterprise_admin_list } = res.data;
        localStorage.setItem("TOKEN", JSON.stringify(token));
        // localStorage.setItem("LOGIN_HASH", JSON.stringify(hash));
        this.getEnterpriseToken(enterprise_admin_list[0])
      }
    },
    /**
     * 获取企业token
     */
    async getEnterpriseToken(target_enterprise) {
      const res = await Api.getEnterpriseToken({
        admin_id: target_enterprise.admin_id,
        enterprise_id: target_enterprise.enterprise_id || this.enterpriseId
      });
      if(res.code === 200) {
        const {
          token,
          enterprise_admin_info,
          permissions_list,
          permissions_tree,
          role_list,
          hash,
        } = res.data;
        localStorage.setItem("LOGIN_HASH", JSON.stringify(hash));
        localStorage.setItem("ENTERPRISE_ID", enterprise_admin_info.enterprise_id);
        localStorage.setItem("TOKEN", JSON.stringify(token));
        localStorage.setItem("KEYS", JSON.stringify(permissions_list));
        PerModule.setPermission(permissions_tree); //存储权限树
        PerModule.setRoles(role_list); //存储角色list
        UserModule.setUserInfo(enterprise_admin_info);
        await this.getEpConfig();
        this.$router.push({
          name: "myapp_home"
        });
      } else {
        this.$message.error(res.msg || '登录失败')
      }
    },
    /** 获取企业身份，区分ks、coach、sw */
    async getEpConfig() {
      const res = await settingApis.getEpConfigApi();
      if (res.code !== 200) return 

      // 企业名称：coach;ks;sw
      const { name } = res.data;
      // 缓存企业名称，给详情页面用
      PerModule.setEqName(name);
    }
  }
}
</script>
<style lang="scss" scoped></style>
