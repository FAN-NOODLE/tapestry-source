<template>
  <div class="supernatant">
    <layout-form>
      <template slot="content">
        <cy-form
          :model="info"
          label-position="left"
          ref="info"
          label-width="140px"
          :rules="rules"
        >
          <cy-form-item label="商户ID" prop="id" required>
            <cy-input
              widthSize="medium"
              v-model="info.id"
              placeholder="请输入商户ID"
              maxlength="32"
            >
            </cy-input>
          </cy-form-item>
          <cy-form-item label="微信支付API key" prop="key" required>
            <cy-input
              widthSize="medium"
              v-model="info.key"
              placeholder="微信支付API key"
              maxlength="32"
            >
            </cy-input>
          </cy-form-item>
          <cy-form-item label="API_V3_KEY" prop="api_v3_key">
            <cy-input
              widthSize="medium"
              v-model="info.api_v3_key"
              placeholder="请输入API_V3_KEY"
              maxlength="32"
            >
            </cy-input>
          </cy-form-item>
          <cy-form-item label="支付证书编号" prop="key" required>
            <cy-input
              widthSize="medium"
              v-model="info.pay"
              placeholder="请输入支付证书编号"
            >
            </cy-input>
          </cy-form-item>
          <cy-form-item label="微信支付证书私钥" prop="private_file">
            <cy-upload
              ref="upload-private"
              class="upload-template"
              type="file"
              :file-list="filelist"
              :http-request="handleRequest2"
              :on-success="onSuccessFun"
            >
              <el-button class="" type="primary" v-if="!info.private_file"
                >点击上传</el-button
              >
              <el-button class=" " @click="clearPrivate" v-else
                >重新上传</el-button
              >
            </cy-upload>
            <div class="po-100 tips">
              上传文件后，文件不可见，仅用于存储设置
            </div>
            <div class="po-100 tips m-t-10">
              {{ this.info.private_file.name }}
            </div>
          </cy-form-item>
          <!-- <div class="pos-box">
              <el-button class="btns">
                <input
                  class="class-model"
                  type="file"
                  name="file"
                  @change="getPrivate($event)"
                />
                <span>点击上传</span>
                <span class="po-100" v-if="info.private_file || !is_first"
                  >已上传证书私钥</span
                >
                <span class="po-100" v-if="!info.private_file && is_first"
                  >请选择apiclient_key.pem 文件</span
                >
              </el-button>
            </div> -->
          <cy-form-item label="微信支付证书" prop="public_file">
            <cy-upload
              ref="upload-public"
              class="upload-template"
              type="file"
              :file-list="filelist"
              :on-success="onSuccess"
              :http-request="handleRequest1"
            >
              <el-button class="" type="primary" v-if="!info.public_file"
                >点击上传</el-button
              >
              <el-button class="" @click="clearPublic" v-else
                >重新上传</el-button
              >
            </cy-upload>
            <div class="po-100 tips ">
              上传文件后，文件不可见，仅用于存储设置
            </div>
            <div class="po-100 tips m-t-10">
              {{ this.info.public_file.name }}
            </div>
            <!-- <div class="pos-box">
              <el-button class="btns">
                <input
                  class="class-model"
                  type="file"
                  name="file"
                  @change="getPublic($event)"
                />
                <span>点击上传</span>
                <span class="po-100" v-if="info.public_file || !is_first"
                  >已上传证书私钥</span
                >
                <span class="po-100" v-if="!info.public_file && is_first"
                  >请选择apiclient_key.pem 文件</span
                >
              </el-button>
            </div> -->
          </cy-form-item>
        </cy-form>
      </template>
      <template slot="footer">
        <el-button type="primary" size="small" @click="saveData" class="font-14"
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
import CyUpload from "@/components/public/cy-upload/index.vue";
import Api from "@/api/setting/setting";
import LoginAboat from "@/utils/LoginManager";
@Component({
  components: {
    LayoutForm,
    CyUpload
  }
})
export default class extends Vue {
  private filelist = [];
  // 把上传的值传到详情里面去
  // 密钥上传
  private onSuccessFun(res: any) {
    this.info.private_file = res.data.url;
    console.log(res.data.url);
  }
  // 证书上传
  private onSuccess(res: any) {
    this.info.public_file = res.data.url;
  }
  // 拿取上传微信支付证书文件
  private handleRequest1(data: any) {
    this.info.public_file = data.file;
    console.log(this.info.public_file);
    this.$message({
      type: "success",
      message: "保存成功"
    });
  }
  // 上传私钥
  private handleRequest2(data: any) {
    this.info.private_file = data.file;
    console.log(this.info.private_file);
    this.$message({
      type: "success",
      message: "保存成功"
    });
  }

  // 弹窗
  private dialogVisible = false;
  // key 校验
  private check_file = (rule: any, value: any, callback: any) => {
    const message = {
      private_file: "必须上传微信支付证书私钥",
      public_file: "必须上传微信支付证书"
    };
    this.is_first || callback();
    if (!value) {
      callback(new Error(message[rule]));
    } else {
      callback();
    }
  };
  private is_first = true;
  private title_name = "支付配置"; //主标题
  private info = {
    id: "",
    key: "", // 微信支付API key
    api_v3_key: "", // api_v3_key
    pay: "", // 支付证书编号
    private_file: "", //微信支付证书私钥
    public_file: "" //微信支付证书
  };
  private rules = {};
  private goBack() {
    this.$router.go(-1);
  }
  // 删除dom上传的文件
  private clearPrivate() {
    // 获取dom,并调取组件中的方法（清除文件）
    (this.$refs["upload-private"] as HTMLFormElement).clearFiles();
  }
  private clearPublic() {
    (this.$refs["upload-public"] as HTMLFormElement).clearFiles();
  }

  // 保存前校验
  private saveData() {
    // console.log("chengg");
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
      api_v3_key: this.info.api_v3_key,
      serial_no: this.info.pay
    };
    const url = this.$VUE_APP_API_ROOT + "gw-console/v1/payment/save";
    const form_data = new FormData();
    form_data.append("private_key", params.private_key);
    form_data.append("public_key", params.public_key);
    form_data.append("merchant_id", params.merchant_id);
    form_data.append("merchant_key", params.merchant_key);
    form_data.append("serial_no", params.serial_no);
    form_data.append("api_v3_key", params.api_v3_key);
    const config = {
      headers: {
        "Enterprise-Hash": LoginAboat.getRealHash(),
        "Auth-Token": LoginAboat.getToken()
      }
    };
    this.$axios
      .post(url, form_data, config)
      .then(() => {
        this.$message({
          type: "success",
          message: "保存成功"
        });
        this.$router.push({
          name: "myapp_setting_authorization-settings_wx-pay_pay-detail"
        });
      })
      .catch((err: any) => {
        const {
          data: { message, code }
        } = err.response;
        this.$message.error(message);
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
    Api.getPayApi().then((res) => {
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
      height: 36px;
    }
    .po-100 {
      position: absolute;
      line-height: 36px;
      top: 0;
      left: 120px;
    }
  }
}
/deep/.tips {
  color: #909399 !important;
}
/deep/.label {
  max-width: 112px !important;
  padding-left: 10px !important;
}
/deep/.in-item {
  margin-left: -7px !important;
}
/deep/.el-input__inner {
  width: 400px !important;
}
/deep/.el-form-item__label {
  width: 142px !important;
}
/deep/.el-form-item {
  margin-left: -10px;
}
/deep/.el-form-item__content {
  margin-left: 142px !important;
}
.m-l-170 {
  margin-left: -170px;
}
.m-t-10 {
  margin-top: -10px;
}
</style>
