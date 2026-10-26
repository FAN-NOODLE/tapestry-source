<template>
  <el-dialog
    width="544px"
    :visible.sync="visible"
    :title="config.title || '上传文件'"
    :before-close="() => closeModel()"
    :close-on-click-modal="false"
  >
    <div class="upload">
      <!-- 模版下载部分 -->
      <div class="clearfix">
        <div class="pull-left" v-text="config.downloadText || '模版下载'"></div>
        <div class="pull-left m-l-20">
          <slot name="downloadTemplate"></slot>
        </div>
      </div>
      <!-- 拖拽上传部分 -->
      <div class="clearfix p-t-20">
        <div class="pull-left" v-text="config.uploadText || '模版上传'"></div>
        <div class="pull-left m-l-20">
          <cy-upload
            ref="uploads"
            drag
            show-file-list
            :file-list="fileLists"
            :on-change="handleChange"
            :before-upload="beforeUploadFile"
          >
            <i class="el-icon-upload" aria-hidden="true"></i>
            <div class="el-upload__text">
              点击或将文件<em>拖拽到这里上传</em>
            </div>
            <div class="el-upload__tip">
              {{
                config.extType ? `只能上传${config.extType.join("、")}文件` : ""
              }}
            </div>
          </cy-upload>
        </div>
      </div>
      <!-- 上传结果部分 -->
      <div class="clearfix p-t-20">
        <div class="pull-left" v-text="config.resultText || '上传结果'"></div>
      </div>
    </div>
    <div slot="footer">
      <el-button size="small" @click="closeModel">取消</el-button>
      <el-button class="confirm" size="small" type="primary"
        >确认更新</el-button
      >
    </div>
  </el-dialog>
</template>
<script lang="ts">
interface IUploadConfig {
  title?: string;
  downloadText?: string;
  uploadText?: string;
  resultText?: string;
  extType: Array<string>;
  uploadUrl?: string;
}
import { Component, Vue, Prop } from "vue-property-decorator";
import CyUpload from "@/components/public/cy-upload/index.vue";
@Component({
  components: {
    CyUpload
  }
})
export default class UploadDrag extends Vue {
  /**
   * @param title
   * @description 上传弹框的标题
   * */
  @Prop() config!: IUploadConfig;
  @Prop({ default: false }) visible!: boolean;
  private defaultUrl = `${process.env.VUE_APP_API_ROOT}${process.env.VUE_APP_API_CONSOLE}upload`; //gw-console
  private fileLists: Array<any> = [];
  private handleChange(files: any, fileList: Array<any>) {
    if (fileList.length === 0) {
      return false;
    } else {
      const file = fileList[fileList.length - 1];
      //覆盖上一次的文件
      this.fileLists = [file];
    }
  }
  // 上传文件前
  private beforeUploadFile(file: any) {
    // 验证传入文件的类型
    if (this.config.extType) {
      const file_msg = file.name.substring(file.name.lastIndexOf(".") + 1);
      const could_type = this.config.extType.join("、");
      if (this.config.extType.indexOf(file_msg) === -1) {
        this.$message({
          message: `文件格式有误，请上传${could_type}格式的文件`,
          type: "warning",
          duration: 2000
        });
        this.fileLists = [];
        return false;
      }
    }
  }
  private closeModel() {
    this.$emit("update:visible", false);
    (this.$refs["uploads"] as any).clearFiles();
  }
}
</script>
<style lang="scss" scoped>
/deep/ {
  .el-upload {
    width: 300px;
    height: 180px;
  }
  .el-upload-dragger {
    @extend .el-upload;
    border: 1px solid $color-second-line !important;
    border-radius: 4px !important;
  }
  .el-upload__tip {
    color: $color-text-second !important;
  }
  .el-upload-dragger {
    .el-icon-upload {
      font-size: 80px !important;
      margin-top: 30px !important;
      margin-bottom: 20px !important;
    }
  }
}
.confirm {
  margin-left: 20px !important;
}
</style>
