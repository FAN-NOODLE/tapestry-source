<template>
  <div class="supernatant">
    <layout-form>
      <template slot="content">
        <cy-form
          :model="detail"
          labelPosition="left"
          labelWidth="190px"
          class="margin-left--26"
        >
          <cy-form-item label="授权小程序" prop="name">
            <div>{{ detail.name }}</div>
          </cy-form-item>
          <cy-form-item label="授权状态">
            <div v-if="detail.installed">已授权</div>
            <div v-else>未授权</div>
          </cy-form-item>
          <cy-form-item label="授权时间" prop="auth_at">
            <div>{{ detail.auth_at }}</div>
          </cy-form-item>
          <cy-form-item label="小程序商城APP Secret" prop="app_secret">
            <div>{{ detail.app_secret }}</div>
          </cy-form-item>
          <cy-form-item label="企业微信Corpld" prop="corp_id">
            <div>{{ detail.corp_id }}</div>
          </cy-form-item>
          <cy-form-item label="商城应用Secret" prop="mini_program_secret"
            ><div>{{ detail.mini_program_secret }}</div></cy-form-item
          >
          <cy-form-item label="商城应用Agentld" prop="agent_id"
            ><div>{{ detail.agent_id }}</div></cy-form-item
          >
        </cy-form>
      </template>
      <template slot="footer">
        <el-button type="primary" size="small" @click="editData()">
          <span class="font-14">编辑</span>
        </el-button>
        <el-button size="small" @click="goBack()">
          <span class="font-14">返回</span>
        </el-button>
      </template>
    </layout-form>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LayoutForm from "@/components/public/layout-form/index.vue";

import Api from "@/api/setting/setting";

@Component({
  components: {
    LayoutForm
  }
})
export default class extends Vue {
  // 商城小程序授权详情
  private detail = {};
  // 跳转编辑页面
  private editData() {
    this.$router.push({
      name: "myapp_setting_authorization-settings_store-set_register-edit"
    });
  }
  // 返回
  private goBack() {
    this.$router.push({
      name: "myapp_setting_authorization-settings_store-set_store-sets"
    });
  }
  // 获取商城小程序授权详情接口
  private getshopAuthDetail() {
    Api.shopAuthDetail({}).then(res => {
      this.detail = res.data;
    });
  }

  mounted() {
    // 获取商城小程序授权详情接口
    this.getshopAuthDetail();
  }
}
</script>

<style lang="scss" scoped>
.margin-left--26 {
  margin-left: -26px;
}
/deep/.label {
  max-width: 143px !important;
}
/deep/.el-form-item__label {
  margin-bottom: 30px;
}
</style>
