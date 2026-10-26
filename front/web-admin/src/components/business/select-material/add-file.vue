<template>
  <div>
    <div>
      <cy-form ref="cform" :rules="rules" :model="formData" label-width="90px">
        <cy-form-item label="素材分组" prop="libraryGroup">
          <cy-select
            v-model="formData.libraryGroup"
            size="medium"
            placeholder="请选择素材分组"
            search-input
            :popperAppendToBody="true"
          >
            <cy-option
              :label="item.group_name"
              :value="item.id"
              :key="index"
              v-for="(item, index) in newGroupData"
            ></cy-option>
          </cy-select>
        </cy-form-item>
        <cy-form-item label="本地文件" prop="cover_url" required>
          <large-file-upload
            v-model="formData.cover_url"
            :beforeUpload="beforeUploadFn"
            :onSuccess="onSuccessFn"
            :onRemove="removeFn"
            :onError="errorFn"
          >
            <el-button type="primary" size="small">上传文件</el-button>
          </large-file-upload>
        </cy-form-item>
      </cy-form>
    </div>
    <div class="dialog-footer">
      <el-button @click="goBack" size="small">返 回</el-button>
      <el-button size="small" type="primary" @click="addConfirm"
        >确 认</el-button
      >
    </div>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import Api from "@/api/enterprise/material-library";
import LargeFileUpload from "@/components/business/large-file-upload/index.vue";

@Component({
  components: {
    LargeFileUpload
  }
})
export default class extends Vue {
  //------------------------prop-------------------------------
  @Prop() newGroupData?: any;
  private fileObj: any = {}; //上传文件信息
  private fileSuccessObj: any = {}; //上传成功信息
  private isShowBtn = true; //是否显示按钮
  private isClosePop = true; //是否关闭弹框
  private isFail = false; //是否上传失败
  private formData = {
    libraryGroup: "",
    // group_id: 0,
    cover_url: "", //视频封面
    name: "", //名称
    size: "", //大小
    url: "", //	视频
    ext: "" //后缀
  };
  private rules = {};
  //--------------------data-----------------------
  //------------------method-----------------------
  private goBack() {
    this.$emit("submit");
  }
  private addConfirm() {
    if (this.isShowBtn) {
      this.$message.warning("请上传文件");
      return false;
    }
    if (!this.isClosePop) {
      this.$message.warning("文件正在上传中");
      return false;
    }
    if (this.isFail) {
      this.$message.error("文件上传失败");
      return false;
    }
    (this.$refs["cform"] as any).validate((valid: any) => {
      if (valid) {
        const suffix = this.fileObj.name.split(".")[
          this.fileObj.name.split(".").length - 1
        ];
        Api.fileAddApi({
          group_id: parseInt(this.formData.libraryGroup) || -1,
          name: this.fileObj.name,
          size: Number(this.fileObj.size),
          url: this.fileSuccessObj.url,
          ext: suffix
        }).then((res: any) => {
          if (res.code == 200) {
            this.$message.success("新增成功");
            this.$emit("submit");
          } else {
            this.$message.error(res.message);
          }
        });
      } else {
        return;
      }
    });
  }
  private beforeUploadFn(file: any) {
    if (
      (file.type == "image/png" || file.type == "image/jpeg") &&
      file.size > 2 * 1024 * 1024
    ) {
      this.$message.error("请上传2MB大小以内的图片");
      return false;
    }

    if (file.type == "video/mp4" && file.size > 1024 * 1024 * 10) {
      this.$message.error("请上传10M大小以内的视频");
      return false;
    }
    if (
      file.type !== "video/mp4" &&
      file.type !== "image/png" &&
      file.type !== "image/jpeg" &&
      file.size > 1024 * 1024 * 20
    ) {
      this.$message.error("请上传20M大小以内的文件");
      return false;
    }

    this.isShowBtn = false;
    this.isClosePop = false;
    this.isFail = false; //上传失败
    this.fileObj = file;
    return true;
  }
  // 上传错误
  private errorFn() {
    this.isClosePop = true;
    this.isFail = true; //上传失败
  }
  //上传成功
  private onSuccessFn(res: any) {
    this.isClosePop = true;
    this.fileSuccessObj = res.data;
    this.isFail = false;
    const { url } = res.data;
    this.$set(this.formData, "cover_url", url);
  }
  //移除上传的文件
  private removeFn() {
    this.isShowBtn = true;
  }
}
</script>
<style lang="scss">
.dialog-footer {
  text-align: right;
  margin-top: 20px;
}
</style>
