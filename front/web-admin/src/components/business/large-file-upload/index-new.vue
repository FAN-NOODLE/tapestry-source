<template>
  <el-upload
    ref="upload"
    :multiple="false"
    :show-file-list="false"
    :accept="accept"
    action=""
    :before-upload="handleBeforeUpload"
    :limit="1"
    :disabled="disabled"
    :http-request="handleRequest"
  >
    <slot></slot>
    <div v-if="isUploading" class="progress" @click.prevent.stop="">
      <!-- 文件进度条样式 -->
      <div
        class="progress-inner"
        :class="{ error: isError }"
        :style="innerStyle"
      ></div>
      <!-- 文件描述 -->
      <div class="progress-text">
        <!-- 文件名称 -->
        <div class="name">{{ fileName }}</div>
        <!-- 文件上传状态 -->
        <div class="status">
          <span v-if="isError">上传失败</span>
          <span v-if="!isError"
            >已上传：{{ progressData.loaded | formatFileSize }}/共{{
              progressData.total | formatFileSize
            }}</span
          >
          <span v-if="!isError" class="m-l-5">{{ uploadStatus }}</span>
        </div>
      </div>
      <div class="remove" @click.prevent.stop="remove">
        <i class="el-icon el-icon-close" aria-hidden="true"></i>
      </div>
    </div>
  </el-upload>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
import COSSDK from "cos-js-sdk-v5";
import { ApiResult } from "@/types";
import { getFileCredentialApi } from "@/api/enterprise/sensitive-word";
import { API } from "@/utils/env-config";
import { cosUploader } from "@/utils/upload";
import { getCosCoverBlob } from "@/utils/upload/utils";
interface Credential {
  bucket: string;
  region: string;
  secret_id: string;
  secret_key: string;
  session_token: string;
  start_time: number;
  expired_time: number;
}

@Component({
  model: {
    prop: "value",
    event: "change"
  }
})
export default class LargeFileUpload extends Vue {
  @Prop({ default: "" }) value!: string;
  @Prop(Function) onSuccess!: Function;
  @Prop(Function) beforeUpload?: Function;
  @Prop(Function) onRemove?: Function;
  @Prop(Function) onError?: Function;
  @Prop({ default: "" }) accept?: string;
  @Prop({ default: false }) disabled!: boolean;
  private isUploading = false;
  private isError = false;
  private cos: any = null;
  private progressData = {
    loaded: 0,
    total: 0,
    speed: 0,
    percent: 0
  };
  private file?: File;
  private fileName = "";
  private credential?: Credential;
  private taskId = "";
  get innerStyle() {
    return {
      width: `${this.progressData.percent * 100}%`
    };
  }

  get uploadStatus() {
    if (this.progressData.percent === 1) {
      return "上传完成";
    }
    return `${(this.progressData.percent * 100).toFixed(2)}%`;
  }

  private mounted() {
    this.reset();
  }
  private selectedFile($event: any) {
    console.log($event.target.files, "file");
  }
  /**
   * 数据重置
   */
  private reset() {
    this.isError = false;
    this.isUploading = false;
    this.progressData = {
      loaded: 0,
      total: 0,
      speed: 0,
      percent: 0
    };
    this.fileName = "";
    this.file = undefined;
  }
  
  // 上传前的回调
  private async handleBeforeUpload(file: File) {
    if (this.beforeUpload) {
      const flag = await this.beforeUpload(file);
      console.log(flag);
      
      if (!flag) {
        throw new Error("上传前校验失败");
      }
    }
    return true;
  }

  async handleRequest(e: any) {
    try {
      
      console.log(e);
      const { file } = e;

      this.isUploading = true;
      this.fileName = file.name;
      this.file = file;

      
      const videoRes = await cosUploader.sliceUploadFile(file, {
        onProgress: (data:any) => {
          this.progressData = { ...data };
        }
      })
      console.log(videoRes);
      // 获取视频封面图
      const poster = await getCosCoverBlob(videoRes.url);
      console.log(poster);
      

      const posterRes = await cosUploader.uploadFile(poster, {
        filename: `${file.name.split(".")[0]}_poster.jpg`
      });

      this.isUploading = false;
      
      this.$emit("change",videoRes.url);
      if (this.onSuccess) {
        this.onSuccess({
          data: {
            ...videoRes,
            fileName: file.name,
            poster: posterRes.url
          }
        });
      }
    } catch (error) {
      console.log(error);
      this.isError = true;
      if (this.onError) {
        this.onError();
      }
      
    }

  }

  // 移除已上传或上传中的文件
  private remove() {
    this.$emit("change", "");
    if (this.onRemove) {
      this.onRemove();
    }
    (this.$refs.upload as any).clearFiles();
  }
}
</script>

<style scoped lang="scss">
::v-deep .el-upload {
  width: 100%;
  text-align: left;
}
.progress {
  margin-top: 50px;
  width: 600px;
  height: 36px;
  background: $color-second-line;
  border-radius: 4px;
  position: relative;
  color: $color-text-main;
  cursor: default;
  .progress-inner {
    height: 36px;
    border-radius: 4px;
    background: $color-main;
    &.error {
      width: 100% !important;
      background: $color-danger;
    }
  }
  .progress-text {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 36px;
    border-radius: 4px;
    padding: 0 15px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    box-sizing: border-box;
    .name {
      color: #ffffff;
      flex: 1;
      color: #ffffff;
      @include textOverflow;
    }
    .status {
      color: #ffffff;
      min-width: 200px;
      color: #ffffff;
      text-align: right;
      color: #ffffff;
    }
  }
  .remove {
    width: 14px;
    height: 14px;
    border-radius: 7px;
    background: $color-placeholder;
    color: $color-white;
    font-size: 12px;
    position: absolute;
    top: -7px;
    right: -7px;
    cursor: pointer;
    line-height: 14px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
}
</style>
