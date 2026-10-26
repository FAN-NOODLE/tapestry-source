<template>
  <div id="crs-scrm">
    <router-view></router-view>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import layout from "./layout/layout.vue";
import settingApis from "@/api/setting/setting";
import { ApiResult } from "./types";
import PerModule from "@/store/modules/permission";
import LoginAboat from "@/utils/LoginManager";
const isProduction = process.env.NODE_ENV === "production";
@Component({
  name: "App",
  components: {
    layout
  }
})
export default class extends Vue {
  private mounted() {
    // console.log('页面路由', this.$route);
    
    // if(this.$route.path === '/web-admin/callback/sso') return
    // this.getEpConfig();
  }

  /** 获取企业身份，区分ks、coach、sw */
  private async getEpConfig() {
    if(!LoginAboat.getToken()) return
    const res: ApiResult = await settingApis.getEpConfigApi();
    if (res.code !== 200) return 

    // 企业名称：coach;ks;sw
    const { name } = res.data;
    // 缓存企业名称，给详情页面用
    PerModule.setEqName(name);
  }
}
</script>
<style lang="scss">
// 自定义进度条颜色
#nprogress .bar {
  background: #000 !important; //自定义颜色
}
</style>
