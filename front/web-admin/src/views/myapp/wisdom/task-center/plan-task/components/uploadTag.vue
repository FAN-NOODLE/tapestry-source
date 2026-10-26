<template>
  <div>
    <el-dialog
      title="上传标签"
      :visible="true"
      @close="closeUploadDialog"
      :close-on-click-modal="false"
    >
      <cy-form
        ref="cform"
        label-width="90px"
        :rules="planRules"
        :model="params"
      >
        <cy-form-item
          label="任务对象"
          required
        >
          <el-radio
            v-model="params.ba_type"
            label="1"
          >关联导购</el-radio>
          <el-radio
            v-model="params.ba_type"
            label="2"
          >管理导购</el-radio>
        </cy-form-item>
        <cy-form-item
          label="标签组"
          required
        >
          <cy-input
            size="small"
            width-size="medium"
            placeholder="请输入标签组名称"
            v-model="params.tag_group"
          />
        </cy-form-item>
        <cy-form-item
          label="标签名称"
          required
        >
          <cy-input
            size="small"
            width-size="medium"
            placeholder="请输入标签名称"
            v-model="params.tag_name"
          />
        </cy-form-item>
        <cy-form-item
          label="union_id是否加密"
          required
        >
          <el-popover
            placement="top-start"
            width="240"
            trigger="hover"
            content="unionid加密方式支持AES-256-CBC。"
          >
            <span slot="reference">
              <i class="el-icon-warning-outline"></i>
            </span>
          </el-popover>
          <el-radio
            v-model="params.union_id_encrypt"
            label="1"
          >是</el-radio>
          <el-radio
            v-model="params.union_id_encrypt"
            label="2"
          >否</el-radio>
          <span
            class="upload-template"
            @click="downTemplate"
          >下载模板</span>
        </cy-form-item>
      </cy-form>
      <div class="upload-wraper">
        <el-upload
          class="upload-demo"
          drag
          accept=".csv"
          name="file"
          :multiple="false"
          :limit="1"
          :data="params"
          :headers="headers"
          :action="action"
          :file-list="fileList"
          :before-upload="beforeUpload"
          :on-success="uploadSuccess"
          :on-error="uploadError"
        >
          <i class="el-icon-upload"></i>
          <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
        </el-upload>
        <p class="tip-text">{{ uploadMessage }}</p>
      </div>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import LoginAboat from "@/utils/LoginManager";
import { VUE_APP_API_ROOT, VUE_APP_API_VERSION } from "@/utils/env-config";
@Component
export default class extends Vue {
  /************data**************/
  private headers = {
    "Auth-Token": LoginAboat.getToken(),
  };
  private params = {
    tag_group: "",
    tag_name: "",
    ba_type: "1",
    union_id_encrypt: "1",
  };
  private planRules = {};
  private action = `${VUE_APP_API_VERSION}tag/csvFile/upload`;
  private fileList: Array<any> = [];
  private uploadMessage = "暂未检测到上传信息";
  private downloadUrl = `${VUE_APP_API_ROOT}ba-performance/download`;
  private beforeUpload(file: any) {
    if (this.params.tag_group === "" || this.params.tag_name === "") {
      this.$message.warning("请先输入标签组名称和标签名称");
      return false;
    }
    const fileSuffix = file.name.substring(file.name.lastIndexOf(".") + 1);
    const whiteList = ["csv"];
    if (whiteList.indexOf(fileSuffix) === -1) {
      this.$message.error("上传文件只能是csv");
      return false;
    }
    const size = file.size / 1024 / 1024;
    if (size > 20) {
      this.$message.warning("文件过大,不能超过20M");
      return false;
    }
  }
  private uploadSuccess(res: any) {
    if (res.code !== 200) {
      this.uploadMessage = res.message;
      this.fileList = [];
      return;
    }
    this.$message.success(
      "人群包上传成功！标签同步完成后，请在【客户标签】中查看"
    );
    setTimeout(() => {
      this.closeUploadDialog();
    }, 500);
  }

  private uploadError(res: any) {
    const error: any = JSON.parse(res.message).message;
    this.$message.error("上传失败,请重试");
    this.uploadMessage = error;
  }

  private downTemplate() {
    // console.log("downloadUrl", this.downloadUrl);
    window.location.href = this.downloadUrl;
  }

  private closeUploadDialog() {
    this.$emit("closeUploadDialog");
  }
}
</script>
<style lang="scss" scoped>
.upload-wraper {
  margin-left: 90px;
}
/deep/ .in-item {
  margin-bottom: 20px;
}
.tip-text {
  margin: 10px 0 20px 0;
}
/deep/ .el-icon-warning-outline {
  font-size: 18px;
  position: relative;
  left: -20px;
  top: 1px;
}
.upload-template {
  margin-left: 58px;
  color: #6f3112;
  cursor: pointer;
}
</style>
