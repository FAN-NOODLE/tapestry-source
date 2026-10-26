<template>
  <div class="lager">
    <!-- 视频预览 -->
    <div class="video-content" v-if="showVideo">
      <div class="mask" @click="closeVideo"></div>
      <div class="content">
        <video
          :src="videoOptions.src"
          :controls="videoOptions.controls"
          class="video-js vjs-big-play-centered vjs-fluid"
          playsinline="true"
          x-webkit-airplay="allow"
          x5-playsinline
          webkit-playsinline
          object-fit="cover"
          style="width: 100%; max-height: 50vh"
          ref="videoRef"
        ></video>
      </div>
    </div>
    <!-- 视频上传  添加此属性后capture="camcorder"，安卓只能录制视频，无法选择相册视频-->
    <van-uploader
      v-show="!fileList.length"
      :before-read="beforeRead"
      :after-read="afterRead"
      :max-count="1"
      accept="video/*"
    >
      <div class="upload-icon">
        <div></div>
        <p>上传视频</p>
      </div>
    </van-uploader>
    <template v-if="fileList.length">
      <div v-for="(file, index) in fileList" :key="index" class="upload">
        <div class="preview">
          <div class="delete-icon" @click.stop.prevent="deleteFile(0)"></div>
          <div v-show="[2, 3, 4].includes(file.status)" class="mask"></div>
          <div
            v-if="file.status == 3"
            class="retry"
            @click.stop.prevent="retryHandle(0)"
          >
            <img class="retry-icon" :src="retryIcon" alt="" />
            <p>点击重试</p>
          </div>
          <img
            v-if="file.status == 2"
            class="fail-icon"
            :src="failIcon"
            alt=""
          />
          <!-- 视频播放 -->
          <img
            v-if="file.status == 4"
            class="play-icon"
            :src="playIcon"
            @click.stop.prevent="playHandle(file)"
            alt=""
          />
          <video :poster="file.poster" :id="file.id">
            <source :src="file.url" alt="视频" />
          </video>
        </div>
        <div class="detail">
          <div class="file-name-wrap">
            <div class="file-name">
              <p class="ellipsis">{{ file.name }}</p>
              <p class="f-s-0">.{{ file.type }}</p>
            </div>
            <!-- 10MB内可选是否原视频上传 -->
            <div
              v-show="file.status === 4 && file.size <= 10 * 1024 * 1024"
              class="original"
            >
              <p>原视频</p>
              <van-checkbox v-model="file.original" @change="changeOriginal">
                <template #icon>
                  <img
                    class="check-icon"
                    :src="file.original ? activeIcon : inactiveIcon"
                  />
                </template>
              </van-checkbox>
            </div>
          </div>
          <div class="file-status" v-show="file.status !== 0">
            <div
              v-if="[2, 3, 4].includes(file.status)"
              class="file-message"
              :class="{
                success: file.status === 4,
                fail: [2, 3].includes(file.status),
              }"
            >
              {{ file.message }}
            </div>
            <div class="w60">
              <van-progress
                v-if="[1].includes(file.status)"
                track-color="#eee"
                color="#000"
                :show-pivot="false"
                :percentage="file.progress"
              />
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { Loading, Uploader, Checkbox, Progress } from "vant";
import { getFileCredentialApi } from "@/api/marketing-promotion/index";
import COSSDK from "cos-js-sdk-v5";
import { API } from "@/utils/env-config";
import wx from "weixin-js-sdk";
interface ICredential {
  bucket: string;
  region: string;
  secret_id: string;
  secret_key: string;
  session_token: string;
  start_time: number;
  expired_time: number;
}
interface IInfo {
  TmpSecretId: string;
  TmpSecretKey: string;
  XCosSecurityToken: string;
  StartTime: number;
  ExpiredTime: number;
}

@Component({
  components: {
    [Checkbox.name]: Checkbox,
    [Loading.name]: Loading,
    [Uploader.name]: Uploader,
    [Progress.name]: Progress,
  },
})
export default class extends Vue {
  /**选中，未选中图片 */
  private activeIcon = require("@/assets/material-library/active_icon.png");
  private inactiveIcon = require("@/assets/material-library/inactive_icon.png");
  private retryIcon = require("@/assets/material-library/retry_icon.png");
  private failIcon = require("@/assets/material-library/fail_icon.png");
  private playIcon = require("@/assets/material-library/play_icon.png");
  /** 视频预览 */
  private videoOptions = {
    controls: true,
    src: "", // url地址
  };

  //视频预览控制
  private showVideo = false;
  private player = null;
  /** 正在上传的文件 */
  private fileList = [];

  /**上传签名 */
  private credential?: ICredential;
  /** 存储实例 */
  private cosInstance: any = null;

  /** 上传前格式校验 */
  private beforeRead(file) {
    console.log("beforeRead", file, file.type, file.name);
    //TODO: ios 相册不是 mp4
    if (!file.type.includes("video")) {
      this.$toast("格式仅支持视频");
      return false;
    }
    return true;
  }
  /** 上传之前 */
  private afterRead(upFile) {
    console.log("afterRead", upFile);
    const fileArray = upFile.file.name.split(".");
    let uploadFile = {
      local: upFile.content, //本地文件流
      name: fileArray[fileArray.length - 2], //文件名称
      type: fileArray[fileArray.length - 1], //文件类型
      size: upFile.file.size, //文件大小
      file: upFile.file, //文件对象
      original: false, //是否原图
      id: "", // task-id
      url: "", //上传成功的服务器地址
      status: 0, // 0待上传 1上传中 2校验失败 3上传失败 4上传成功
      progress: 0, //上传进度
      message: "", //错误信息
    };

    /**本地读取成功，但是size 大于10MB */
    if (upFile.file.size > 10 * 1024 * 1024) {
      uploadFile.status = 2;
      uploadFile.message = "视频大小不能超过10M！";
      this.fileList.push(uploadFile);
      return;
    }
    this.fileList.push(uploadFile);
    //调用COS
    this.singleUpload(uploadFile.file, this.fileList[this.fileList.length - 1]);
  }
  /** 上传到服务端 */
  private async singleUpload(file: File, uploadFile) {
    //判断是否有证书，如果没有或证书过期，则重新初始化
    if (!this.credential || this.credential.expired_time < Date.now() / 1000) {
      const flag = await this.initCos();
      // 初始化失败，停止上传
      if (!flag) {
        uploadFile.url = ""; //上传成功的服务器地址
        uploadFile.message = "加载失败"; //错误信息
        uploadFile.status = 2; // 0待上传 1上传中 2校验失败 3上传失败 4上传成功
        return;
      }
    }

    uploadFile.url = ""; //上传成功的服务器地址
    uploadFile.message = ""; //错误信息
    uploadFile.status = 1; // 0待上传 1上传中 2校验失败 3上传失败 4上传成功
    console.log("file", file, this.cosInstance);
    //  分片上传视频
    this.cosSliceUpload(file, uploadFile);
  }
  /**
   * @description: cos分片上传
   */
  private cosSliceUpload(file: File, uploadFile) {
    this.cosInstance.sliceUploadFile(
      {
        Bucket: this.credential?.bucket /* 必须 */,
        Region: this.credential?.region /* 存储桶所在地域，必须字段 */,
        Key: this.calcKey(file) /* 必须 */,
        Body: file /* 必须 */,
        SliceSize:
          1024 *
          1024 *
          10 /* 触发分块上传的阈值，超过5MB使用分块上传，非必须 */,
        onTaskReady: (taskId: string) => {
          uploadFile.id = taskId; // 任务id
        },
        onProgress: (progressData: { percent: number }) => {
          console.log("progressData", progressData);
          uploadFile.progress = progressData.percent * 100; //上传进度
        },
      },
      async (err: any, data: any) => {
        if (err) {
          uploadFile.message = "网络波动,加载失败！"; //错误信息
          uploadFile.status = 3; // 0待上传 1上传中 2校验失败 3上传失败 4上传成功
          return;
        }
        const url = this.getCosUrl(data.Key);
        // 获取视频文件第一帧图片
        const firstImageURl = await this.findVideoPoster(url);
        console.log(firstImageURl);
        if (!firstImageURl) {
          return;
        }
        uploadFile.status = 4;
        uploadFile.url = url; //上传成功的服务器地址
        uploadFile.message = "加载成功"; //成功信息
        // 上传成功，更新视频第一帧图片
        uploadFile.poster = firstImageURl;
        this.$emit("change", this.findSuccess(this.fileList));
      }
    );
  }
  /**
   * @description: 获取cos URl
   * @param {string} key
   */
  private getCosUrl(Key) {
    let api = "";
    if (API.name.indexOf("PROD") > -1) {
      api = "https://cdn-cos.coach.com.cn";
    } else if (API.name.indexOf("UAT") > -1) {
      api = "https://tapestry-uat-1308608494.cos.ap-shanghai.myqcloud.com";
    } else {
      api = "https://tapestry-uat-1308608494.cos.ap-shanghai.myqcloud.com";
    }
    const url = `${api}/${Key}`;
    return url;
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
  // 初始化腾讯对象存储Cos
  private async initCos() {
    try {
      const { code, data }: IApiResult<ICredential> =
        await getFileCredentialApi();
      if (code === 200) {
        this.credential = data;
        this.cosInstance = new COSSDK({
          getAuthorization: (options: any, callback: (IInfo) => void) => {
            console.log(options);
            callback({
              TmpSecretId: data?.secret_id,
              TmpSecretKey: data?.secret_key,
              XCosSecurityToken: data?.session_token,
              StartTime: data?.start_time,
              ExpiredTime: data?.expired_time,
            });
          },
        });
        return true;
      }
      return false;
    } catch (e) {
      return false;
    }
  }
  /** 删除文件 */
  private deleteFile(index) {
    //停止上传任务
    if (this.fileList[index].id) {
      this.cosInstance.cancelTask(this.fileList[index].id);
    }
    this.fileList.splice(index, 1);
    this.$emit("change", this.findSuccess(this.fileList));
  }
  /**重新上传文件 */
  private retryHandle(index) {
    console.log("retryHandle");
    this.cosInstance.restartTask(this.fileList[index].id);
  }
  /**关闭视频预览 */
  private closeVideo() {
    this.showVideo = false;
    this.player.currentTime = 0;
    this.player && this.player.pause();
    this.player = null;
  }
  /**播放当前视频预览 */
  private playHandle(file) {
    this.videoOptions.src = file.local;
    this.showVideo = true;
    setTimeout(() => {
      this.player = this.$refs.videoRef;
      this.player.play();
    }, 500);
  }
  /**返回文件列表中成功的文件 */
  private findSuccess(fileList) {
    console.log(fileList);
    return fileList.filter((file) => file.status === 4);
  }
  /** 视频文件勾选是否原视频上传 */
  private changeOriginal() {
    this.$emit("change", this.findSuccess(this.fileList));
  }
  /**
   * @description: 获取上传视频文件的第一帧
   */
  async findVideoPoster(videoUrl) {
    let self = this;
    console.log("video=====================" + videoUrl);
    return new Promise(function (resolve) {
      let base64URL = "";
      let video = document.createElement("video");
      video.setAttribute("crossOrigin", "anonymous"); //处理跨域
      video.setAttribute("src", videoUrl);
      // x5-video-player-type="h5"
      video.setAttribute("x5-video-player-type", "h5");
      video.currentTime = 1;
      video.load();
      video.addEventListener(
        "loadeddata",
        async function () {
          let canvas = document.createElement("canvas");
          //使用视频的宽高作为canvas、预览图的宽高
          let width = video.videoWidth;
          let height = video.videoHeight;
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
          self.cosInstance.putObject(
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
<style lang="scss" scoped>
.upload-icon {
  width: 88px;
  height: 88px;
  background: #ffffff;
  border: 1px solid #000000;
  position: relative;
  div {
    width: 30px;
    height: 1px;
    background: #000;
    position: absolute;
    top: 36px;
    left: 50%;
    transform: translateX(-50%);
    &::before {
      content: "";
      width: 1px;
      height: 30px;
      background: #000;
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
    }
  }
  p {
    position: absolute;
    bottom: 10px;
    left: 0;
    width: 100%;
    text-align: center;
    color: #9e9e9e;
    font-size: 12px;
    font-family: PingFangSC-Regular, PingFang SC;
    font-weight: 400;
    color: #9e9e9e;
  }
}
.upload {
  display: flex;
  .preview {
    flex-shrink: 0;
    position: relative;
    width: 88px;
    height: 88px;
    background-color: rgba(240, 240, 240, 1);
    margin-right: 15px;
    .delete-icon {
      width: 20px;
      height: 20px;
      background: url(../../../assets/material-library/close_icon.png) no-repeat;
      background-size: cover;
      position: absolute;
      top: -10px;
      right: -10px;
      z-index: 3;
    }
    .mask {
      background: rgba(0, 0, 0, 0.7);
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
    }
    .play-icon,
    .fail-icon {
      width: 30px;
      height: 30px;
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      z-index: 4;
    }
    .retry {
      width: 100%;
      height: 100%;
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      z-index: 1;
      p {
        width: 100%;
        text-align: center;
        position: absolute;
        bottom: 17px;
        left: 50%;
        font-size: 12px;
        font-family: PingFangSC-Regular, PingFang SC;
        font-weight: 400;
        color: #ffffff;
        transform: translateX(-50%);
      }
    }
    .retry-icon {
      width: 30px;
      height: 30px;
      position: absolute;
      top: 19px;
      left: 50%;
      transform: translateX(-50%);
    }
    video,
    img {
      display: block;
      width: 100%;
      height: 100%;
    }
  }
  .detail {
    overflow: hidden;
    flex-grow: 1;
    .file-name-wrap {
      display: flex;
      justify-content: space-between;
      .file-name {
        overflow: hidden;
        flex-grow: 1;
        font-size: 14px;
        font-family: PingFangSC-Regular, PingFang SC;
        font-weight: 400;
        color: #000000;
        display: flex;
        .ellipsis {
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
        .f-s-0 {
          flex-shrink: 0;
        }
      }
      .original {
        margin-left: 16px;
        display: flex;
        flex-shrink: 0;
        font-size: 12px;
        font-family: PingFangSC-Regular, PingFang SC;
        font-weight: 400;
        color: #000000;

        .check-icon {
          margin-top: 1px;
          margin-left: 7px;
          display: block;
          width: 18px;
          height: 18px;
        }
      }
    }
    .file-status {
      margin-top: 10px;
      .file-message {
        font-size: 12px;
        font-family: PingFangSC-Regular, PingFang SC;
        font-weight: 400;
        &.success {
          color: #3bc255;
        }
        &.fail {
          color: #d84242;
        }
      }
      .w60 {
        width: 60px;
      }
    }
  }
}
.video-content {
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  right: 0;
  z-index: 99;
  .mask {
    position: absolute;
    width: 100%;
    height: 100%;
    background: rgba($color: #000000, $alpha: 0.6);
    z-index: 2;
  }
  .content {
    position: absolute;
    width: 100%;
    top: 0;
    max-height: 50vh;
    z-index: 8;
    background: white;
    font-size: 0;
  }
}
</style>
