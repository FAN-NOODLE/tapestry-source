<template>
  <div>
    <cy-form
      class="main-box"
      ref="cform"
      :model="ruleForm"
      :rules="rules"
      label-width="90px"
    >
      <cy-form-item label="素材分组" prop="libraryGroup">
        <cy-select
          v-model="ruleForm.sucai_id"
          size="medium"
          placeholder="请选择素材分组"
          search-input
          :optionsList="newGroupData"
          :options="{ label: 'group_name', value: 'id' }"
        >
        </cy-select>
      </cy-form-item>
      <cy-form-item label="本地图片" required prop="urls" trigger="change">
        <upload-image
          class="upload-template"
          :urls="ruleForm.urls"
          accept="image/jpeg,image/png"
          img-type="batch"
          :before-upload="beforeUpload"
          @update="update"
          :limit="9"
          multiple
        >
        </upload-image>
        <span>支持JPG、PNG格式，图片小于2MB</span>
      </cy-form-item>
    </cy-form>

    <div class="fotter">
      <el-button @click="goBack" size="small">返 回</el-button>
      <el-button type="primary" @click="childSure" size="small"
        >确 认</el-button
      >
    </div>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import Api from "@/api/enterprise/material-library";
import UploadImage from "@/components/business/upload-image/index.vue";
@Component({
  components: {
    UploadImage
  }
})
export default class extends Vue {
  //------------------------prop-------------------------------
  @Prop() newGroupData?: any;
  //------------------------data-------------------------------
  private ruleForm: any = {
    sucai_id: "",
    urls: []
  };
  private rules = {};
  private isUploading = false;
  //----------------------------------------------------methods-----------------------------------------------------------/
  private update(res: Array<string>) {
    this.ruleForm.urls = res;
    this.isUploading = false;
  }
  // 上传文件限制规范
  private beforeUpload(file: any) {
    const name = file.name
      .substring(file.name.lastIndexOf(".") + 1, file.name.length)
      .toLocaleLowerCase();
    if (name !== "jpg" && name !== "png") {
      this.$message.error("请上传JPG、PNG格式的图片");
      return false;
    }
    if (file.size > 2 * 1024 * 1024) {
      this.$message.error("请上传2MB大小以内的图片");
      return false;
    }
    this.isUploading = true;
  }
  // 内弹框点击确定
  private childSure() {
    if (this.isUploading) {
      this.$message.warning("图片正在上传中");
      return;
    }
    (this.$refs["cform"] as any).validate((valid: any) => {
      if (valid) {
        this.imageAdd();
      }
    });
  }
  private goBack() {
    this.$emit("submit");
  }
  // 图片上传接口
  imageAdd() {
    const params = {
      group_id: this.ruleForm.sucai_id || -1,
      images: this.ruleForm.urls
    };
    Api.imageAddApi(params).then((res: any) => {
      if (res.code == 200) {
        console.log("成功！");
        this.$emit("submit");
      } else {
        this.$message.error(res.message);
      }
    });
  }
}
</script>
<style lang="scss">
.fotter {
  text-align: right;
  margin-top: 20px;
}
</style>
