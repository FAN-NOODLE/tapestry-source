<template>
  <div class="upload-image-batch clearfix">
    <preview
      v-if="preview_show"
      :urls="banner_list.map(item => item.url)"
      :index="preview_index"
      @close="preview_show = false"
    ></preview>
    <!-- 上传器-->
    <el-upload
      :show-file-list="false"
      :file-list="banner_list"
      class="upload-machine"
      :class="disabled ? 'no-upload' : ''"
      :action="action"
      :accept="accept"
      :disabled="disabled"
      :http-request="handleRequest"
      :before-upload="beforeUpload"
      :on-remove="onRemove"
      :limit="limit"
      :multiple="multiple"
      :on-exceed="handleExceed"
      enctype="multipart/form-data"
    >
      <slot>
        <el-button
          size="small"
          type="primary"
          :disabled="disabled"
        >本地上传</el-button>
      </slot>
    </el-upload>
    <div
      class="image-list"
      :style="{ '--width': `${width}px`, '--height': `${height}px`, 'marginTop': `${disabled ? '0px' : '50px'}` }"
    >
      <!-- 开启拖拽 -->
      <div class="list-wrap">
        <draggable
          v-model="banner_list"
          :options="{
            animation: 150,
            ghostClass: 'sortable-ghost',
            chosenClass: 'chosenClass',
            scroll: true,
            scrollSensitivity: 200,
            disabled: file_list.some(item => item.progress_percent < 100)
          }"
          @end="endListenter"
        >
          <!-- 上传成功的图片 -->
          <div 
            v-for="($item, $index) in banner_list"
            :key="$item.url + $index"
            @mouseover.prevent="onMouseOver($index)"
            @mouseleave.prevent="$item.is_hover = false"
            >

            <div
              class="image-item"
              :style="{ backgroundImage: `url(${$item.url})` }"
            >
              <!-- 小图布局 -->
              <template v-if="size==='small'">
                <div class="preview" @click.prevent="toPreview($index)"></div>
                <el-button 
                  v-if="!disabled" 
                  type="text" 
                  circle 
                  class="close" 
                  icon="el-icon-close"
                  @click="deleteImage($index)"
                ></el-button>
              </template>
              <!-- 大图布局 -->
              <template v-else>
                <div
                  class="no-allowed"
                  v-if="disabled"
                ></div>
                <div
                  class="mask"
                  v-else
                  v-show="$item.is_hover"
                >
                  <i
                    v-if="!allowPreview"
                    class="el-icon-delete bin"
                    @click="deleteImage($index)"
                    aria-hidden="true"
                  ></i>
                  <div
                    class="double-btn"
                    v-if="allowPreview"
                  >
                    <i
                      class="el-icon-zoom-in"
                      @click.prevent="toPreview($index)"
                      aria-hidden="true"
                    ></i>
                    <i
                      class="el-icon-delete has-preview"
                      @click.prevent="deleteImage($index)"
                      aria-hidden="true"
                    ></i>
                  </div>
                </div>
              </template>
              111
            </div>
          </div>
          <!-- 正在上传的图片,显示滚动条 -->
          <div
            class="progess-list"
            v-for="(uploading, index) in file_list"
            v-show="uploading.progress_percent < 100"
            :key="index"
          >
            <el-progress
              class="progess"
              :stroke-width="10"
              type="circle"
              color="#6F3112"
              :percentage="uploading.progress_percent"
            ></el-progress>
          </div>
        </draggable>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
let index = -1;
import { Component, Vue, Prop, Watch } from "vue-property-decorator";
import { Loading } from "element-ui";
import { IFileParam } from "./type";
import draggable from "vuedraggable";
import CyUpload from "@/components/public/cy-upload/index.vue";
import preview from "@/components/business/preview/index.vue";
import LoginAboat from "@/utils/LoginManager";
import { VUE_APP_API_ROOT } from "@/utils/env-config";
// todo:增加获取图片信息
import { getImageWh } from "@/utils/get-image-info";
import Api from "@/api/enterprise/material-library";
@Component({
  components: { draggable, preview, CyUpload }
})
export default class UploadImageBatch extends Vue {
  @Prop(Function) onSuccess?: Function;
  @Prop(Function) onError?: Function;
  @Prop(Function) beforeUpload?: Function;
  @Prop(Function) onRemove?: Function;
  @Prop({ default: 0 }) limit?: number;
  @Prop() urls!: Array<any>;
  @Prop({ default: false }) multiple?: boolean;
  // @Prop({ default: 100 }) width?: number;
  // @Prop({ default: 100 }) height?: number;
  @Prop({ default: "medium" }) size?: string;
  // @Prop({ default: true }) preview?: boolean;
  @Prop({ default: "上传图片" }) tips?: string;
  @Prop({ default: "" }) accept?: string;
  @Prop({ default: false, type: Boolean }) disabled!: boolean;
  @Prop({
    default: `${VUE_APP_API_ROOT}file/upload`
  })
  action?: string;
  @Watch("urls", { immediate: true, deep: true })
  onchangeUrls() {
    this.banner_list = this.urls.map(url => {
      let obj: IFileParam = {
        is_hover: false
      };
      if (Object.prototype.toString.call(url) === "[object Object]") {
        obj = { ...url };
      } else {
        obj.url = url;
      }
      return obj;
    });
    if (this.size === "small") {
      this.width = 50;
      this.height = 50;
      this.allowPreview = false;
    } else {
      this.width = 100;
      this.height = 100;
    }
  }
  destroyed() {
    index = -1;
  }
  width = 100;
  height = 100;
  allowPreview = true;
  banner_list: Array<IFileParam> = [];
  file_list: Array<any> = [];
  preview_show = false;
  preview_index = 0;
  // 鼠标放到图片上
  onMouseOver(index) {
    this.$set(this.banner_list[index], "is_hover", true);
  }
  // 删除图片
  deleteImage(i: number) {
    this.onRemove && this.onRemove(i);
    this.banner_list.splice(i, 1);
    this.$emit(
      "change",
      this.banner_list.map(item => item.url)
    );
    this.$emit("update", this.banner_list);
  }
  handleExceed() {
    this.$message.warning(`限制上传${this.limit}张图片,请重新选择`);
  }
  toPreview(index: number): void {
    this.preview_show = true;
    this.preview_index = index;
  }
  //移动监听改变数组
  endListenter() {
    this.$emit(
      "change",
      this.banner_list.map(item => item.url)
    );
    this.$emit("update", this.banner_list);
  }
  // 成功回调
  uploadSuccess(file: any) {
    const obj = {
      url: file.data.url,
      is_hover: false
    };
    this.banner_list.push(obj);
    this.onSuccess && this.onSuccess();
  }
  handleRequest(data: any) {
    index++;
    const { file } = data;
    this.file_list.push({ ...file });
    (_index => {
      const formdata = new FormData();
      formdata.append("file", data.file);
      const { name, size } = data.file;
      const fileParam: IFileParam = { name, ext: data.file.type, size };
      const onUploadProgress = (progressEvent: any) => {
        const { loaded, total } = progressEvent;
        const progress_percent = ~~((loaded / total) * 100);
        this.$set(this.file_list[_index], "progress_percent", progress_percent);
      };
      const config = {
        onUploadProgress,
        headers: {
          "Enterprise-Hash": LoginAboat.getRealHash(),
          "Auth-Token": LoginAboat.getToken()
        }
      };
      this.$axios
        .post(this.action, formdata, config)
        .then(async (res: any) => {
          const { data, code, message } = res.data;
          if (code === 200) {
            Loading.service({ fullscreen: true }).close();
            const { width, height } = await getImageWh(data.url);
            // todo:如果width 大于 1440 或者 height 大于 1080
            if (width > 1440 || height > 1080) {
              // todo:压缩图片
              const compressData: any = await Api.compressImage({
                file_type: 'image',
                file_path: data.url,
              });
              if (compressData.code === 200) {
                const newData = {
                  url: compressData.data,
                  is_hover: false
                }
                this.uploadCompressSuccess(fileParam, newData);
              } else {
                this.$message({
                  message: compressData.message,
                  type: "error",
                  duration: 2000
                });
                this.$emit("fail");
              }
            } else {
              this.uploadCompressSuccess(fileParam, data);
            }

          } else {
            this.$message({
              message: message,
              type: "error",
              duration: 2000
            });
            this.$emit("fail");
          }
        })
        .catch((err: any) => {
          const { message } = err;
          this.$message({
            message: message,
            type: "error",
            duration: 2000
          });
          this.$emit("fail");
        });
    })(index);
  }
  // todo:上传且压缩成功
  uploadCompressSuccess(fileParam: any, data: any) {
    this.$message({
      message: "上传成功",
      type: "success",
      duration: 2000
    });
    fileParam.url = data.url;
    fileParam.is_hover = false;
    this.banner_list.push(fileParam);

    this.onSuccess && this.onSuccess();
    this.$emit(
      "change",
      this.banner_list.map(item => item.url)
    );
    this.$emit("update", this.banner_list);
  }
}
</script>
<style lang="scss" scoped>
.dynamic-size {
  width: var(--width);
  height: var(--height);
}

.upload-image-batch {
  .image-list {
    float: left;

    &:after {
      display: block;
      clear: both;
      content: "";
    }

    .list-wrap {
      float: left;
    }

    .progess-list {
      @extend .dynamic-size;
      margin-right: 10px;
      margin-bottom: 10px;
      float: left;
      position: relative;

      .progess {
        @extend .dynamic-size;
        border: 1px solid $color-second-line;
        border-radius: 4px;
        position: absolute;
        top: 0;
        left: 0;
        margin: auto;

        ::v-deep.el-progress-circle {
          width: calc(var(--width) * 0.6) !important;
          height: calc(var(--height) * 0.6) !important;
          position: absolute;
          top: 0;
          left: 0;
          bottom: 0;
          right: 0;
          margin: auto;
        }
      }
    }

    .image-item {
      @extend .dynamic-size;
      position: relative;
      margin-right: 10px;
      margin-bottom: 10px;
      background-size: 100% 100%;
      border-radius: 6px;
      float: left;
      // overflow: hidden;
      cursor: pointer;

      .no-allowed {
        cursor: not-allowed;
        border-radius: 6px;
        background-color: #dfe1e4;
        opacity: 0.3;
        position: relative;
        width: 100%;
        height: 100%;
      }

      .preview {
        width: 100%;
        height: 100%;
        position: absolute;
        top: 0;
        left: 0;
      }
      .close {
        position: absolute;
        top: 0;
        right: 0;
        transform: translate(50%, -50%);
        background-color: #dfe1e4;
        color: #fff !important;
        border-radius: 50% !important;
        font-size: 12px
      }

      .mask {
        @extend .dynamic-size;
        border-radius: 6px;
        background-color: rgba(0, 0, 0, 0.5);
        position: relative;


        .bin {
          color: $color-white;
          font-size: 20px;
          position: absolute;
          left: 32%;
          top: 36%;
        }

        .double-btn {
          width: 60px;
          position: absolute;
          top: 0;
          left: 0;
          bottom: 0;
          right: 0;
          margin: auto;

          .el-icon-zoom-in {
            cursor: pointer;
            color: white;
            font-size: 20px;
            position: absolute;
            top: 1px;
            left: 0;
            line-height: var(--height) !important;
          }

          .has-preview {
            cursor: pointer;
            color: white;
            font-size: 20px;
            position: absolute;
            top: 0;
            right: 0px;
            line-height: var(--height) !important;
          }
        }
      }
    }
  }

  .upload-machine {
    float: left;
  }
}

::v-deep.no-upload {
  .el-upload {
    cursor: not-allowed;
    background: $color-background;

    &:hover {
      border: 1px solid $color-second-line;
    }
  }
}

::v-deep {
  .el-upload {
    @extend .dynamic-size;
    line-height: var(--height);
  }

  .el-upload--picture-card {
    background-color: $color-white;
    border: 1px solid $color-second-line;
    border-radius: 4px;
    position: relative;

    &:hover {
      border: 1px solid $color-main;
    }
  }
}
</style>
<style lang="scss" scoped>
.img-list-item {
  display: inline-block;
  margin-right: 10px;
}

.img-list-enter-active,
.img-list-leave-active {
  transition: all 1s;
}

.img-list-enter,
.img-list-leave-to {
  opacity: 0;
  transform: translateY(30px);
}
</style>
