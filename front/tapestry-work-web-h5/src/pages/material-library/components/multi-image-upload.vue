<template>
  <div class="lager">
    <!-- 原生图片上传 -->
    <van-uploader
      v-show="fileList.length < maxLength"
      v-model="vanFileList"
      :preview-image="false"
      :max-count="maxLength"
      :before-read="beforeRead"
      :after-read="afterRead"
      multiple
      accept="image/jpeg,image/png"
    >
      <div class="wrap">
        <div class="upload-icon">
          <div></div>
          <p>上传图片</p>
        </div>
      </div>
    </van-uploader>
    <!-- 全部原图 -->
    <div class="choose-all" v-if="fileList.length > 0">
      <van-checkbox
        v-model="chooseAll"
        @click="toggleOriginAll"
        label-position="right"
      >
        <template #icon>
          <img
            class="check-icon"
            :src="chooseAll ? activeIcon : inactiveIcon"
          />
        </template>
        <p>全部原图</p>
      </van-checkbox>
      <!-- todo:已上传几张 -->
      <div
        class="upload-length"
        v-if="isCoach"
      >
        已上传:{{ fileList.length }}张
      </div>
      <!-- todo:长按图片拖拽 -->
      <div
        class="drag-image"
        v-if="isCoach"
      >
        长按图片可拖拽排序
      </div>
    </div>

    <template v-if="fileList.length">
      <draggable
        v-model="fileList"
        :options="{
          animation: 150,
          ghostClass: 'sortable-ghost',
          chosenClass: 'chosenClass',
          scroll: true,
          scrollSensitivity: 200,
          disabled: !isCoach
        }"
        handle=".preview"
        filter=".delete-icon"
        @end="endListenter"
      >
        <div
          v-for="(file, index) in fileList"
          :key="index"
          class="upload"
        >
          <div class="preview">
            <div
              class="delete-icon"
              @click.stop.prevent="deleteFile(index)"
            ></div>
            <div
              v-show="[2, 3].includes(file.status)"
              class="mask"
            ></div>
            <div
              v-if="file.status == 3"
              class="retry"
              @click.stop.prevent="retryHandle(0)"
            >
              <img
                class="retry-icon"
                :src="retryIcon"
                alt=""
              />
              <p>点击重试</p>
            </div>
            <img
              v-if="file.status == 2"
              class="fail-icon"
              :src="failIcon"
              alt=""
            />
            <!-- 图片预览 -->
            <img
              class="drag-image"
              @click="previewImg(file.url)"
              :src="file.local"
              alt="图片"
              style="object-fit: contain; pointer-events: none;"
            />
          </div>
          <div class="detail">
            <div class="file-name-wrap">
              <div class="file-name">
                <template>
                  <div class="file-name-space">
                    <p class="ellipsis">{{ file.name }}</p>
                    <p class="f-s-0">.{{ file.type }}</p>
                  </div>
                  <i
                    class="edit-icon"
                    @click="editStart(file,index)"
                  ></i>
                </template>
              </div>
            </div>
            <!-- 10MB内可选是否原图上传 -->
            <div
              v-if="file.status === 4 && file.size <= mB10"
              class="original"
            >
              <p>原图</p>
              <van-checkbox
                v-model="file.original"
                @change="changeOriginal"
                label-position="right"
              >
                <template #icon>
                  <img
                    class="check-icon"
                    :src="file.original ? activeIcon : inactiveIcon"
                  />
                </template>
              </van-checkbox>
            </div>
            <div
              class="file-status"
              v-show="file.status !== 0"
            >
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
      </draggable>
    </template>
    <!-- 企微API -->
    <div
      class="wrap"
      v-show="false && fileList.length < maxLength"
      @click="apiChoose"
    >
      <div class="upload-icon">
        <div></div>
        <p>上传图片</p>
      </div>
    </div>
    <!-- 上传视频限制 desc -->
    <ul class="desc">
      <li>1. 图片格式仅支持*.jpg；*.png</li>
      <li>2. 单张展示每次最多支持上传20张图片，合并展示每次最多支持9张，图片不可超过50M</li>
      <li>3. 图片上传时会进行压缩，10M以内图片可以选择原图上传</li>
      <li>4. 苹果手机jpeg格式图片原图上传后，图片大小会有变动</li>
    </ul>
    <!-- todo:编辑图片弹窗 -->
    <edit-name-popup
      :showPopup.sync="showEditPopup"
      :file="currentFile"
      :index="currentFileIndex"
      @enterFileName="enterFileName"
    />
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import { Loading, Uploader, Checkbox, Progress } from "vant";
import { getFileCredentialApi } from "@/api/marketing-promotion/index";
import COSSDK from "cos-js-sdk-v5";
import { API } from "@/utils/env-config";
import { ImagePreview, Field } from "vant";
// todo:拖拽排序
import draggable from 'vuedraggable'
// todo:编辑图片名称
import editNamePopup from "./edit-name-popup.vue";
import Helper from "@/utils/helper";
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
    [Field.name]: Field,
    draggable,
    editNamePopup
  },
})
export default class MultiImageUpload extends Vue {
  // todo:最大数量
  @Prop({ default: 20 }) readonly maxLength!: number;
  private showEditPopup = false;
  private isCoach = Helper.isCoach();
  private currentFile:any = {};
  private currentFileIndex = 0;
  /**选中，未选中图片 */
  private activeIcon = require("@/assets/material-library/active_icon.png");
  private inactiveIcon = require("@/assets/material-library/inactive_icon.png");
  private retryIcon = require("@/assets/material-library/retry_icon.png");
  private failIcon = require("@/assets/material-library/fail_icon.png");
  private playIcon = require("@/assets/material-library/play_icon.png");
  /**是否全部原图 */
  private chooseAll = false;
  /** 正在上传的文件 */
  private fileList = [];
  /**组件双向绑定的list  解决多选bug*/
  private vanFileList = [];
  /**上传签名 */
  private credential?: ICredential;
  /** 存储实例 */
  private cosInstance: any = null;
  /** 10MB */
  private mB10 = 10 * 1024 * 1024;
  //todo:移动监听改变数组
  endListenter() {
    const sortFileList = this.fileList
    if (sortFileList && sortFileList.length > 0) {
      for (let i = 0; i < sortFileList.length; i++) {
        sortFileList[i].sort = i
      }
      this.$emit("change", this.findSuccess(sortFileList))
    }
  }
  /**是否展示全部原图 */
  private get showAll() {
    let arr = this.fileList.filter((file) => {
      return file.status == 4 && file.size <= this.mB10;
    });

    return arr.length > 0;
  }
  /** 读取前格式校验 */
  private beforeRead(file) {
    console.log("beforeRead", file);
    const currentFileLength = this.fileList.length;

    //批量选择
    if (Object.prototype.toString.call(file).includes("Array")) {
      if (currentFileLength + file.length > this.maxLength) {
        this.$toast("选择图片大于20张，已为您保留所选图片的前20张");
      }
      return true;
    } else {
      //单选文件,限制上传个数
      if (currentFileLength + 1 > this.maxLength) {
        return false;
      }
      //校验文件格式等等
      return this.dealRead(file);
    }
  }
  //单文件校验
  private dealRead(file) {
    const nameType = file.name
      .substring(file.name.lastIndexOf(".") + 1, file.name.length)
      .toLocaleLowerCase();
    console.log("上传文件名", nameType);

    // if (nameType === "jpeg") {
    //   this.$toast("企业微信不支持JPEG的图片");
    //   return false;
    // }

    if (!["image/jpeg", "image/png"].includes(file.type)) {
      this.$toast("图片格式仅支持*.jpg、*.png");
      return false;
    }
    return true;
  }
  /** 上传之前 */
  private afterRead(upFile) {
    console.log("afterRead", upFile);
    //一次选择多张
    if (Object.prototype.toString.call(upFile).includes("Array")) {
      //校验文件是否合法
      for (let i = 0; i < upFile.length; i++) {
        if (this.dealRead(upFile[i].file)) {
          //符合上传要求
          this.updateList(upFile[i]);
        }
      }
      return;
    }
    this.updateList(upFile);
  }
  /**渲染到页面 */
  private updateList(upFile) {
    //一次选择一张
    const fileArray = upFile.file.name.split(".");
    let original = false;
    //文件小于等于10M时，根据全选按钮状态初始化original
    if (upFile.file.size <= this.mB10 && this.chooseAll) {
      original = true;
    }
    let uploadFile = {
      local: upFile.content, //本地文件流
      name: fileArray[fileArray.length - 2], //文件名称
      type: fileArray[fileArray.length - 1], //文件类型
      size: upFile.file.size, //文件大小
      file: upFile.file, //文件对象
      original: original, //是否原图
      id: "", // task-id
      url: "", //上传成功的服务器地址
      status: 0, // 0待上传 1上传中 2校验失败 3上传失败 4上传成功
      progress: 0, //上传进度
      message: "", //错误信息
      edit: false, //编辑名称开关
    };
    /**本地读取成功，但是size 大于50MB */
    if (upFile.file.size > 50 * 1024 * 1024) {
      uploadFile.status = 2;
      uploadFile.message = "图片过大！";
      this.fileList.push(uploadFile);
      //更新finished状态 控制上传按钮
      this.$emit("change", this.findSuccess(this.fileList));
      return;
    }
    this.fileList.push(uploadFile);
    //更新finished状态 控制上传按钮
    this.$emit("change", this.findSuccess(this.fileList));
    //调用COS
    this.singleUpload(uploadFile.file, this.fileList[this.fileList.length - 1]);
  }
  /** wxApi选择文件 */
  private apiChoose() {
    const self = this;
    const count = self.maxLength - self.fileList.length;
    window.wx.chooseImage({
      count: count, // 默认9
      sizeType: ["original", "compressed"], // 可以指定是原图还是压缩图，默认二者都有
      sourceType: ["album", "camera"], //'camera'可以指定来源是相册还是相机，默认二者都有
      success: function (res) {
        console.log("apiChoose", res);
        // andriod中localId可以作为img标签的src属性显示图片；
        // iOS应当使用 getLocalImgData 获取图片base64数据，从而用于img标签的显示（在img标签内使用 wx.chooseImage 的 localid 显示可能会不成功）
        // if () {
        //   window.wx.getLocalImgData({
        //     localId: "", // 图片的localID
        //     success: function (res) {
        //       var localData = res.localData; // localData是图片的base64数据，可以用img标签显示
        //     },
        //   });
        // }
      },
      fail: function (err) {
        console.log("chooseImage", err);
      },
    });
  }
  /** 上传到服务端 */
  private async singleUpload(file: File, uploadFile) {
    //判断是否有零时证书，如果没有或证书过期，则重新初始化
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
    this.cosInstance.sliceUploadFile(
      {
        Bucket: this.credential?.bucket /* 必须 */,
        Region: this.credential?.region /* 存储桶所在地域，必须字段 */,
        Key: this.calcKey(file) /* 必须 */,
        Body: file /* 必须 */,
        SliceSize:
          1024 *
          1024 *
          10 /* 触发分块上传的阈值，超过10MB使用分块上传，非必须 */,
        onTaskReady: (taskId: string) => {
          uploadFile.id = taskId; // 任务id
        },
        onProgress: (progressData: { percent: number }) => {
          uploadFile.progress = progressData.percent * 100; //上传进度
        },
      },
      (err: any, data: any) => {
        if (err) {
          uploadFile.message = "网络波动,加载失败！"; //错误信息
          uploadFile.status = 3; // 0待上传 1上传中 2校验失败 3上传失败 4上传成功
          return;
        }
        let api = "";
        if (API.name.indexOf("PROD") > -1) {
          api = "https://cdn-cos.coach.com.cn";
        } else if (API.name.indexOf("UAT") > -1) {
          api = "https://tapestry-uat-1308608494.cos.ap-shanghai.myqcloud.com";
        } else {
          api = "https://tapestry-uat-1308608494.cos.ap-shanghai.myqcloud.com";
        }

        const url = `${api}/${data.Key}`;
        uploadFile.status = 4;
        uploadFile.url = url; //上传成功的服务器地址
        uploadFile.message = "加载成功"; //成功信息
        this.$emit("change", this.findSuccess(this.fileList));
      }
    );
  }
  /** 计算key 规则为upload/YYYY/MM/dd/32位随机字符串.文件扩展名 */
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
  /**初始化腾讯对象存储Cos*/
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
    if (this.fileList[index] && this.fileList[index].id) {
      this.cosInstance.cancelTask(this.fileList[index].id);
    }
    this.fileList.splice(index, 1);
    this.vanFileList.splice(index, 1);
    this.$emit("change", this.findSuccess(this.fileList));
  }
  /**重新上传文件 */
  private retryHandle(index) {
    console.log("retryHandle");
    this.cosInstance.restartTask(this.fileList[index].id);
  }
  /** 预览图片 */
  private previewImg(url) {
    let arr = [url];
    ImagePreview({
      images: arr,
      startPosition: 0,
    });
  }
  /**返回文件列表中成功的文件 */
  private findSuccess(fileList) {
    let finished = false;
    if (fileList.length > 0) {
      finished = fileList.every((file) => {
        return file.status === 4 && !file.edit;
      });
      // todo:增加sort字段
      for (let i = 0; i < fileList.length; i++) {
        fileList[i].sort = i
      }
    }
    return {
      fileList: fileList.filter((file) => file.status === 4),
      finished,
    };
  }
  /**
   * 反向处理全部勾选
   * 视频文件勾选是否原视频上传
   */
  private changeOriginal() {
    let supportOrigin = this.fileList.filter((file) => {
      return file.status === 4 && file.size <= this.mB10;
    });
    if (supportOrigin.length) {
      //some 不对空数组执行 有任一个false则灭掉全选
      let someNoChoose = supportOrigin.some((file) => {
        return !file.original;
      });
      // 任意一个false 存在，则chooseAll false;
      // 任意一个false 都不存在，则chooseAll true;
      this.chooseAll = !someNoChoose;
    }

    this.$emit("change", this.findSuccess(this.fileList));
  }
  /** 选择源文件上传 */
  private toggleOriginAll() {
    console.log("toggleOriginAll", this.chooseAll);
    //小于10M的资源 跟随开关是否受控
    this.fileList.forEach((file) => {
      if (file.size <= this.mB10) {
        file.original = this.chooseAll;
      }
    });
  }
  /** 确认-编辑名称 */
  private editConfirm(file) {
    if (file.name.trim()) {
      // const editName = file.name.trim() + "." + file.type;
      if (file.name.trim().length > 50) {
        this.$toast("图片名称不能超过50个字");
        return;
      }
      file.edit = false;
      this.$emit("change", this.findSuccess(this.fileList));
    } else {
      this.$toast("图片名称不能为空");
    }
  }
  /** 变为编辑状态 */
  private editStart(file:any,index:any) {
    this.currentFile = file;
    this.currentFileIndex = index;
    this.showEditPopup = true;
  }
  /** 编辑确认 */
  enterFileName(file:any){
    const {index,name} = file
    this.fileList[index].name = name
    this.$emit("change", this.findSuccess(this.fileList));
    this.showEditPopup = false
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
  flex-shrink: 0;

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

.wrap {
  display: flex;
  align-items: center;
  padding-bottom: 20px;
}

.desc {
  padding-top: 20px;
  font-size: 9px;
  font-family: PingFangSC-Regular, PingFang SC;
  font-weight: 400;
  color: #9e9e9e;
  border-top: 1px solid #e5e5e5;

  li {
    line-height: 20px;
  }
}

.choose-all {
  display: flex;
  align-items: center;
  font-size: 12px;
  font-weight: 400;
  color: #000000;
  margin-top: 8px;
  margin-bottom: 25px;

  .check-icon {
    margin-left: 0;
    display: block;
    width: 18px;
    height: 18px;
  }

  .upload-length {
    margin-left: 8px;
    font-size: 12px;
    font-weight: 400;
    color: #999999;
    line-height: 17px;
  }

  .drag-image {
    margin-left: 16px;
    width: 120px;
    height: 24px;
    background: #CCCCCC;
    font-size: 12px;
    font-weight: 400;
    color: #000000;
    line-height: 24px;
    padding-left: 4px;
    pointer-events:none;
    -webkit-user-select:none;
    -moz-user-select:none;user-select:none
  }
}

.upload {
  padding-top: 14px;
  display: flex;
  margin-bottom: 8px;

  &:last-of-type {
    margin-bottom: 28px;
  }

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
      align-items: center;
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
        align-items: center;
        height: 26px;
        line-height: 26px;

        .file-name-space {
          display: flex;
          overflow: hidden;
          height: 26px;
          line-height: 26px;
        }

        .ellipsis {
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }

        .f-s-0 {
          flex-shrink: 0;
        }
        .confirm-icon,
        .edit-icon {
          flex-shrink: 0;
          width: 18px;
          height: 18px;
          background: url(../../../assets/material-library/edit_icon.png) no-repeat;
          background-size: 18px 18px;
          margin-left: 6px;
        }

        .confirm-icon {
          background: url(../../../assets/material-library/confirm_icon.png) no-repeat;
          background-size: 18px 18px;
        }
      }
    }

    .original {
      display: flex;
      align-items: center;
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
}</style>
