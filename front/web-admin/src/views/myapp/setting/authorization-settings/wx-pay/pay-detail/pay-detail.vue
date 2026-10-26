<template>
  <div class="supernatant">
    <layout-form>
      <template slot="content">
        <cy-form
          :model="info"
          labelPosition="left"
          labelWidth="150px"
          :rules="rules"
          class="margin-left--26"
        >
          <cy-form-item label="商户ID" prop="id">
            <div>{{ info.id }}</div>
          </cy-form-item>
          <cy-form-item label="微信支付API key" prop="key">
            <div>{{ info.key }}</div>
          </cy-form-item>
          <cy-form-item label="API_V3_KEY" prop="api_v3_key">
            <div>{{ info.api_v3_key }}</div>
          </cy-form-item>
          <cy-form-item label="支付证书编号" prop="pay">
            <div>{{ info.pay }}</div>
          </cy-form-item>
          <cy-form-item label="微信支付证书私钥">
            <div v-if="!is_first">已上传证书私钥</div>
            <div v-else>未上传</div>
          </cy-form-item>
          <cy-form-item label="微信支付证书">
            <div v-if="!is_first">已上传支付证书</div>
            <div v-else>未上传</div>
          </cy-form-item>
        </cy-form>
      </template>
      <template slot="footer">
        <el-button type="primary" size="small" @click="editData">
          <span class="font-14">重新配置</span>
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
  private rules = {};
  private is_first = true;
  private title_name = "支付配置"; //主标题
  private edit = true;
  private info = {
    id: "***",
    key: "***", // 微信支付API key
    api_v3_key: "***", // API_V3_KEY
    pay: "***", // 支付证书编号
    private_file: "***", //微信支付证书私钥
    public_file: "***" //微信支付证书
  };
  private editData() {
    this.$router.push({
      name: "myapp_setting_authorization-settings_wx-pay_pay-edit"
    });
  }

  // 获取详情
  private getDetail() {
    Api.getPayApi().then(res => {
      const data = res.data;
      this.is_first = Object.keys(data).length === 0;
      if (data.merchant_id) {
        this.info = {
          id: data.merchant_id,
          key: data.merchant_key,
          pay: data.serial_no,
          api_v3_key: data.api_v3_key,
          private_file: "",
          public_file: ""
        };
      }
    });
  }
  mounted() {
    this.getDetail();
  }
}
</script>

<style lang="scss" scoped>
.margin-left--26 {
  margin-left: -26px;
}
/deep/.label {
  max-width: 112px !important;
}
/deep/.el-form-item__label {
  margin-bottom: 30px;
}
/deep/.el-form-item {
  margin-left: 10px;
}
</style>
