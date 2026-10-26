<template>
  <div class="supernatant">
    <layout-form>
      <template slot="content">
        <cy-form
          :model="info"
          labelPosition="left"
          labelWidth="150px"
          :rules="info_rules"
        >
          <cy-form-item label="商户ID" prop="id">
            <cy-input
              widthSize="medium"
              :value="info.id"
              placeholder="请输入商户ID"
              maxlength="32"
            ></cy-input>
          </cy-form-item>
          <cy-form-item label="微信支付API key" prop="key">
            <cy-input
              widthSize="medium"
              :value="info.key"
              placeholder="微信支付API key"
              maxlength="32"
            >
            </cy-input>
          </cy-form-item>
          <cy-form-item label="支付证书编号" prop="key">
            <cy-input
              widthSize="medium"
              :value="info.pay"
              placeholder="请输入支付证书编号"
            >
            </cy-input>
          </cy-form-item>
          <cy-form-item label="微信支付证书私钥" prop="private_file">
            <div class="pos-box">
              <el-button class="btns">
                <input
                  class="class-model"
                  type="file"
                  name="file"
                  @change="getPrivate($event)"
                />
                <span>点击上传</span>
                <span class="po-100" v-if="info.private_file"
                  >已上传证书私钥</span
                >
                <span class="po-100" v-else>请选择apiclient_key.pem 文件</span>
              </el-button>
            </div>
          </cy-form-item>
          <cy-form-item label="微信支付证书" prop="public_file">
            <div class="pos-box">
              <el-button class="btns">
                <input
                  class="class-model"
                  type="file"
                  name="file"
                  @change="getPublic($event)"
                />
                <span>点击上传</span>
                <span class="po-100" v-if="info.private_file"
                  >已上传证书私钥</span
                >
                <span class="po-100" v-else>请选择apiclient_key.pem 文件</span>
              </el-button>
            </div>
          </cy-form-item>
          <cy-form-item label="">
            <el-button
              class="default-button"
              type="primary"
              size="small"
              @click="saveData"
            >
              保存
            </el-button>
          </cy-form-item>
        </cy-form>
      </template>
    </layout-form>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Api from "@/api/setting/setting";

// key 校验
const check_range = (rule: any, value: any, callback: any) => {
  if (value.length !== 32) {
    callback(new Error("微信支付API key为32位"));
  } else {
    callback();
  }
};
export default class extends Vue {
  private info = {
    id: "", // 商户id
    key: "", // 微信支付API key
    pay: "", // 支付证书编号
    private_file: null,
    public_file: null
  };
  private info_rules = {
    id: [{ required: true, message: "商户ID必填" }],
    pay: [{ required: true, message: "请输入支付证书编号" }],
    key: [
      { required: true, message: "微信支付API key必填" },
      { validator: check_range, trigger: ["blur", "change"] }
    ],
    private_file: [{ required: true, message: "必须上传微信支付证书私钥" }],
    public_file: [{ required: true, message: "必须上传微信支付证书" }]
  };
  // 保存前校验
  private saveData() {
    (this.$refs["info"] as any).validate((valid: any) => {
      if (valid) {
        this.saveApi();
      } else {
        return false;
      }
    });
  }

  // 保存接口
  private saveApi() {
    const params = {
      private_key: this.info.private_file,
      public_key: this.info.public_file,
      merchant_id: this.info.id,
      merchant_key: this.info.key,
      serial_no: this.info.pay
    };
    Api.setPayApi(params).then(() => {
      this.$message({
        type: "success",
        message: "保存成功"
      });
      this.getDetail();
    });
  }

  // 公共 文件
  private getPublic(event: any) {
    const publicFile = event.target.files[0];
    event.target.value = "";
    this.info.public_file = publicFile;
  }

  // 私密文件
  private getPrivate(event: any) {
    const privateFile = event.target.files[0];
    event.target.value = "";
    this.info.private_file = privateFile;
  }

  // 获取详情
  private getDetail() {
    Api.getPayApi().then(res => {
      if (res.data.merchant_id) {
        this.$router.push({
          name: "myapp_setting_basic-setup_wx-pay_pay-detail"
        });
      }
    });
  }
  mounted() {
    this.getDetail();
  }
}
</script>

<style lang="scss" scoped>
.pos-box {
  position: relative;
  .tip {
    position: absolute;
    z-index: 999;
    top: 12px;
    left: 123px;
  }
  .btns {
    position: relative;
    .class-model {
      position: absolute;
      top: 0;
      left: 0;
      z-index: 999;
      opacity: 0;
    }
    .po-100 {
      position: absolute;
      line-height: 36px;
      top: 0;
      left: 120px;
    }
  }
}
</style>
