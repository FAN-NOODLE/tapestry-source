<template>
  <div class="supernatant">
    <layout-form>
      <template slot="content">
        <cy-form
          ref="cform"
          :model="detail"
          labelPosition="left"
          labelWidth="160px"
          class="margin-left--26"
          :rules="rules"
        >
          <cy-form-item label="授权小程序" prop="name">
            <div>{{ detail.name }}</div>
          </cy-form-item>
          <cy-form-item label="授权状态">
            <div v-if="detail.installed">已授权</div>
            <div v-else>未授权</div>
          </cy-form-item>
          <cy-form-item label="授权时间">
            <div>{{ detail.auth_at }}</div>
          </cy-form-item>
          <cy-form-item
            label="小程序商城APP Secret"
            desc="注：小程序商城secret"
            prop="app_secret"
            required
          >
            <cy-input v-model="detail.app_secret"></cy-input>
          </cy-form-item>
          <cy-form-item label="企业微信Corpld" prop="corp_id">
            <div>{{ detail.corp_id }}</div>
          </cy-form-item>
          <cy-form-item
            label="商城应用Secret"
            prop="mini_program_secret"
            desc="注：企业微信secret"
          >
            <cy-input :value="detail.mini_program_secret" clearable></cy-input>
          </cy-form-item>
          <cy-form-item label="商城应用Agentld" prop="agent_id">
            <cy-input :value="detail.agent_id" clearable></cy-input>
          </cy-form-item>
        </cy-form>
      </template>
      <template slot="footer">
        <el-button type="primary" size="small" class="font-14" @click="saveData"
          >保存</el-button
        >
        <el-button size="small" @click="dialogVisible = true" class="font-14"
          >取消</el-button
        >
      </template>
    </layout-form>
    <!-- 弹窗  -->
    <el-dialog title="提示" :visible.sync="dialogVisible" width="30%">
      <span>你确认要放弃编辑吗？放弃后编辑的信息不回被保存</span>
      <span slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="goBack">确 定</el-button>
      </span>
    </el-dialog>
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
  // 弹窗
  private dialogVisible = false;
  // 校验规则
  private rules = {};
  // 商城小程序授权详情
  private detail = {
    app_secret: "", //  商城应用APP Secret
    mini_program_secret: "", //商城应用Secret
    agent_id: "" //商城应用Agentld
  };
  // 取消
  private goBack() {
    this.$router.go(-1);
  }
  // 获取商城小程序授权详情接口
  private getshopAuthDetail() {
    Api.shopAuthDetail({}).then(res => {
      this.detail = res.data;
    });
  }

  // 保存接口
  private saveData() {
    // 保存之前进行校验
    (this.$refs["cform"] as any).validate((valid: boolean) => {
      if (valid) {
        const params = {
          app_secret: this.detail.app_secret,
          mini_program_secret: this.detail.mini_program_secret,
          agent_id: this.detail.agent_id
        };
        // 保存用户输入信息接口
        Api.EditShopApi(params).then(() => {
          this.$message({
            type: "success",
            message: "保存成功"
          });
          this.$router.push({
            name:
              "myapp_setting_authorization-settings_store-set_register-detail"
          });
        });
      }
    });
  }

  mounted() {
    // 获取商城小程序授权详情接口
    this.getshopAuthDetail();
  }
}
</script>

<style lang="scss" scoped>
/deep/.label {
  max-width: 143px !important;
}
.cy-input {
  width: 400px !important;
}
/deep/.el-form-item__label {
  margin-bottom: 30px;
  width: 173px !important;
}
/deep/.el-form-item__content {
  margin-left: 173px !important;
}
/deep/.desc {
  margin-left: 173px !important;
}
</style>
