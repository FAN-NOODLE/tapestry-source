<template>
  <el-upload
    :multiple="false"
    :show-file-list="false"
    :accept="accept"
    action=""
    :before-upload="handleBeforeUpload"
    :limit="1"
    :disabled="disabled"
  >
    <slot></slot>
    <div v-if="isUploading && uploadStatus !== '上传完成'" class="progress" @click.prevent.stop="">
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
    this.initCos();
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
      const flag = this.beforeUpload(file);
      if (!flag) {
        return;
      }
    }
    this.upload(file);
    return false;
  }

  //上传文件
  private async upload(file: File) {
    console.log(file)
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
      async (err: any, data: any) => {
        console.log(err, data);
        this.taskId = "";
        if (err) {
          this.isError = true;
          if (this.onError) {
            this.onError();
          }
          return;
        }
        let api = "";

        if (API) {
          if (API.name.indexOf("PROD") > -1) {
            api =
              "https://cdn-cos.coach.com.cn";
          } else{
            api =
              "https://cdn-cos-uat.coach.com.cn";
          }
        } else {
          api = "https://cdn-cos-uat.coach.com.cn";
        }

        const url = `${api}/${data.Key}`;
        console.log("url::", url);
        console.log("file::", file);
        // 判断是不是视频文件
        const isVideo = file.type.indexOf("video") > -1;
        let poster: any = "";
        // 如果是视频文件，获取视频第一帧
        if (isVideo) {
          // todo:需要创建一个video标签，然后获取视频的宽高
          const video = document.createElement("video");
            video.preload = 'metadata'
            video.src = url;
            // 定义一个promise
            const promise = new Promise((resolve) => {
              video.onloadedmetadata = () => {
                const { duration } = video;
                console.log(duration)
                // 最长不超过30S
                if (duration > 30) {
                  this.$message.error("视频时长不能超过30S,请重新选择");
                  resolve(false);
                }
                resolve(true);
              };
              video.onerror = () => {
                this.$message.error("视频加载失败,请重新选择");
                resolve(false);
              };
            });
          // 等待promise执行完毕
          const result = await promise;
          video.remove();
          if(!result){
            return
          }
          poster = await this.findVideoPoster(url);
          console.log("poster::", poster);
        }
        this.$emit("change",url);
        if (this.onSuccess) {
          this.onSuccess({
            data: {
              url,
              fileName: file.name,
              poster
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
    return `uploads/${date.getFullYear()}/${month > 10 ? month : "0" + month}/${
      day > 10 ? day : "0" + day
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
  /**
   * @description: 获取上传视频文件的第一帧
   */
   async findVideoPoster(videoUrl: string) {
    const self = this;
    console.log("video=====================" + videoUrl);
    return new Promise(function (resolve) {
      let base64URL = "";
      const video = document.createElement("video");
      video.setAttribute("crossOrigin", "anonymous"); //处理跨域
      video.setAttribute("src", videoUrl);
      // x5-video-player-type="h5"
      video.setAttribute("x5-video-player-type", "h5");
      video.currentTime = 1;
      video.load();
      video.addEventListener(
        "loadeddata",
        async function () {
          const canvas = document.createElement("canvas");
          //使用视频的宽高作为canvas、预览图的宽高
          const width = video.videoWidth;
          const height = video.videoHeight;
          canvas.width = width;
          canvas.height = height;
          canvas.getContext("2d").drawImage(video, 0, 0, width, height); //绘制canvas
          base64URL = canvas.toDataURL(); //转换为base64，图片格式默认为png，这里修改为jpeg
          const imgfile = self.data64toFile(base64URL);
          console.log(imgfile);
          const flag = await self.initCos();
          // 初始化失败，停止上传
          if (!flag) {
            return;
          }
          // 上传图片到cos
          self.cos.putObject(
            {
              Bucket: self.credential.bucket,
              Region: self.credential.region,
              // 去除后缀
              Key: self.randomFileName() + ".jpg",
              Body: imgfile,
            },
            function (err, data) {
              if (err && err.error) {
                resolve("");
              } else {
                const url = `https://${data.Location}`;
                resolve(url);
              }
            }
          );
        },
        false
      );
    });
  }
  /**
   * @description: base64转blob
   * @param {String} base64URL
   */
  data64toFile(base64URL) {
    const arr = base64URL.split(",");
    const mime = arr[0].match(/:(.*?);/)[1];
    const bstr = atob(arr[1]);
    let n = bstr.length;
    const u8arr = new Uint8Array(n);
    while (n--) {
      u8arr[n] = bstr.charCodeAt(n);
    }
    return new Blob([u8arr], { type: mime });
  }
  /**
   * @description: 随机生成文件名
   */
  randomFileName() {
    const date = new Date();
    const year = date.getFullYear();
    const month = date.getMonth() + 1;
    const day = date.getDate();
    const hour = date.getHours();
    const minute = date.getMinutes();
    const second = date.getSeconds();
    const millisecond = date.getMilliseconds();
    return `${year}${month}${day}${hour}${minute}${second}${millisecond}`;
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
  width: 100%;
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
