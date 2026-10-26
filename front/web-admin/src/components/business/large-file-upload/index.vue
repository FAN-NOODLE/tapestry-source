<template>
  <el-upload
    :multiple="false"
    :show-file-list="false"
    :accept="accept"
    action=""
    :before-upload="handleBeforeUpload"
  >
    <slot v-if="!isUploading"></slot>
    <div
      v-if="isUploading"
      class="progress"
      @click.prevent.stop=""
    >
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
          <!-- <span v-if="!isError"
            >已上传：{{ progressData.loaded | formatFileSize }}/共{{
              progressData.total | formatFileSize
            }}</span
          > -->
          <span
            v-if="!isError"
            class="m-l-5"
          >{{ uploadStatus }}</span>
        </div>
      </div>
      <div
        class="remove"
        @click.prevent.stop="remove"
      >
        <i
          class="el-icon el-icon-close"
          aria-hidden="true"
        ></i>
      </div>
    </div>
  </el-upload>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
import COSSDK from "cos-js-sdk-v5";
import { ApiResult } from "@/types";
import { getFileCredentialApi } from "@/api/enterprise/sensitive-word";
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
  private containerURL = "";
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
    // this.initCos();
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
  // 初始化腾讯对象存储Cos
  private async initCos() {
    const {
      code,
      message = "",
      data
    }: ApiResult<Credential> = await getFileCredentialApi();
    if (code !== 200) {
      this.$message(message);
      return false;
    }
    this.credential = data;
    this.cos = new COSSDK({
      getAuthorization: (options: any, callback: Function) => {
        console.log(options);
        callback({
          TmpSecretId: data?.secret_id,
          TmpSecretKey: data?.secret_key,
          XCosSecurityToken: data?.session_token,
          StartTime: data?.start_time,
          ExpiredTime: data?.expired_time
        });
      }
    });
    return true;
  }
  // 上传前的回调
  private handleBeforeUpload(file: File) {
    if (this.beforeUpload) {
      const flag = this.beforeUpload(file, "fileUplod");
      if (!flag) {
        return;
      }
    }
    return false;
  }
  //上传文件
  private async upload(file: File) {
    //判断是否有零时证书，如果没有或证书过期，则重新初始化
    if (!this.credential || this.credential.expired_time < Date.now() / 1000) {
      const flag = await this.initCos();
      // 初始化失败，停止上传
      if (!flag) {
        return;
      }
    }
    this.isUploading = true;
    this.fileName = file.name;
    this.file = file;
    this.cos.sliceUploadFile(
      {
        Bucket: this.credential?.bucket /* 必须 */,
        Region: this.credential?.region /* 存储桶所在地域，必须字段 */,
        Key: this.calcKey(file) /* 必须 */,
        Body: file /* 必须 */,
        onTaskReady: (taskId: string) => {
          this.taskId = taskId;
        },
        onProgress: (progressData: any) => {
          this.progressData = { ...progressData };
        }
      },
      (err: any, data: any) => {
        console.log(err, data);
        this.taskId = "";
        if (err) {
          this.isError = true;
          if (this.onError) {
            this.onError();
          }
          return;
        }
        const url = `https://img-crs.vchangyi.com/${data.Key}`;
        console.log(1111, url);
        
        this.$emit("change", url);
        if (this.onSuccess) {
          this.onSuccess({
            data: {
              url
            }
          });
        }
      }
    );
  }

  // 计算key 规则为upload/YYYY/MM/dd/32位随机字符串.文件扩展名
  private calcKey(file: File) {
    const chars = "ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz2345678";
    let fileName = "";
    for (let i = 0; i < 32; i++) {
      fileName += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    const arr = file.name.split(".");
    const date = new Date();
    const month = date.getMonth() + 1;
    const day = date.getDate();
    return `uploads/${date.getFullYear()}/${month > 10 ? month : "0" + month}/${day > 10 ? day : "0" + day
      }/${fileName}.${arr[arr.length - 1]}`;
  }

  // 移除已上传或上传中的文件
  private remove() {
    //存在taskId则是上传中，取消上传
    if (this.taskId) {
      this.cos.cancelTask(this.taskId);
    }
    this.reset();
    this.$emit("change", "");
    if (this.onRemove) {
      this.onRemove();
    }
  }
}
</script>

<style scoped lang="scss">
::v-deep .el-upload {
  width: 100%;
  text-align: left;
}

.progress {
  width: 80%;
  height: 36px;
  background: $color-second-line;
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
    padding: 0 15px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    box-sizing: border-box;

    .name {
      flex: 1;
      color: #ffffff;
      @include textOverflow;
    }

    .status {
      min-width: 200px;
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
