<template>
  <div
    class="com-upload-image"
    :style="{
      '--width': `${width}px`,
      '--height': `${height}px`,
    }"
  >
    <preview
      v-if="preview_show"
      :urls="[imgUrl]"
      @close="preview_show = false"
    ></preview>
    <el-upload
      class="avatar-uploader"
      :class="disabled ? 'no-upload' : ''"
      v-if="!imgUrl"
      ref="upsingle"
      :action="action"
      :accept="accept"
      :show-file-list="false"
      :before-upload="beforeUpload"
      :http-request="handleRequest"
      :disabled="disabled"
      v-show="!progess_show"
    >
      <div class="no-img">
        <i
          class="el-icon-plus avatar-uploader-icon"
          :class="{ 'disabled-tips': disabled }"
          aria-hidden="true"
        ></i>
        <div
          class="tips"
          :class="{ 'disabled-tips': disabled }"
          v-text="tips"
        ></div>
      </div>
    </el-upload>
    <el-progress
      class="progess"
      :stroke-width="10"
      v-if="progess_show"
      type="circle"
      color="#6F3112"
      :width="width"
      :percentage="progress_percent"
    ></el-progress>
    <div
      class="img-info no-allowed"
      v-if="disabled && imgUrl"
    >
      <img
        :src="imgUrl"
        alt=""
      />
    </div>
    <transition
      name="img"
      v-else-if="!disabled && imgUrl"
    >
      <div
        class="img-info"
        @mouseover.prevent="is_hover = true"
        @mouseleave.prevent="is_hover = false"
      >
        <img
          :src="imgUrl"
          alt=""
        />
        <div class="delete-info">
          <div
            class="double-btn"
            v-if="allowPreview"
            v-show="is_hover"
          >
            <i
              class="el-icon-zoom-in"
              @click.prevent="toPreview"
              aria-hidden="true"
            ></i>
            <i
              class="el-icon-delete has-preview"
              @click.prevent="remove"
              aria-hidden="true"
            ></i>
          </div>
          <i
            class="el-icon-delete only-delete"
            v-if="!allowPreview"
            v-show="is_hover"
            @click.prevent="remove"
            aria-hidden="true"
          ></i>
        </div>
      </div>
    </transition>
  </div>
</template>
<script lang="ts">
/**
 * @author LeeYunxiang
 * @desc 不用card-list的组件，可以不用再绑数组，只适用于传一张图片的情景
 * @param {Function} onSuccess 成功回调
 * @param {Function} onError 失败回调
 * @param {Function} beforeUpload 上传前回调
 * @param {Function} onRemove 移除回调
 * @param {String} urls 图片url
 */
import { Component, Vue, Prop, Watch } from "vue-property-decorator";
import { Loading } from "element-ui";
import preview from "@/components/business/preview/index.vue";
import CyUpload from "@/components/public/cy-upload/index.vue";
import LoginAboat from "@/utils/LoginManager";
import { VUE_APP_API_ROOT } from "@/utils/env-config";
@Component({
  components: { preview, CyUpload },
})
export default class ComUploadImage extends Vue {
  @Prop(Function) onError?: Function;
  @Prop(Function) beforeUpload?: Function;
  @Prop(Function) onSuccess?: Function;
  @Prop(Function) onRemove?: Function;
  @Prop(String) urls!: string;
  @Prop({ default: "上传图片" }) tips?: string;
  @Prop({ default: "" }) accept?: string;
  @Prop({ default: false, type: Boolean }) disabled!: boolean;
  @Prop({ default: true }) preview?: boolean;
  @Prop({ default: "medium" }) size?: string;
  @Prop({
    default: `${VUE_APP_API_ROOT}file/upload`, //gw-console
  })
  action?: string;
  @Watch("urls")
  onChangeUrls() {
    this.imgUrl = this.urls;
  }
  mounted() {
    this.imgUrl = this.urls;
    // 预设两种size 100px和50px 不允许自定义宽高了,当宽高为50px时 没有预览按钮了
    if (this.size === "small") {
      this.width = 50;
      this.height = 50;
      this.allowPreview = false;
    } else {
      this.width = 100;
      this.height = 100;
    }
  }
  private width = 100;
  private height = 100;
  private allowPreview = true;
  private preview_show = false;
  private is_hover = false;
  private progress_percent = 0;
  private progess_show = false;
  private imgUrl = "";
  private toPreview(): void {
    this.preview_show = true;
  }

  private close(): void {
    this.preview_show = false;
  }
  private remove(): void {
    this.is_hover = false;
    this.imgUrl = "";
    this.onRemove && this.onRemove();
    this.$emit("change", this.imgUrl);
  }
  private handleRequest(data: any) {
    const formdata = new FormData();
    formdata.append("file", data.file);
    const config = {
      headers: {
        "Enterprise-Hash": LoginAboat.getAppHash(),
        "Auth-Token": LoginAboat.getToken(),
      },
      onUploadProgress: (progressEvent: any) => {
        this.progess_show = true;
        const { loaded, total } = progressEvent;
        this.progress_percent = ~~((loaded / total) * 100);
      },
    };
    this.$axios
      .post(this.action, formdata, config)
      .then((res: any) => {
        const { data, code, message } = res.data;
        if (code === 200) {
          Loading.service({ fullscreen: true }).close();
          this.progess_show = false;
          this.$message({
            message: "上传成功",
            type: "success",
            duration: 2000,
          });
          this.imgUrl = data.url;
          this.onSuccess && this.onSuccess();
          this.$emit("change", this.imgUrl);
        } else {
          this.progess_show = false;
          this.$message({
            message: message,
            type: "error",
            duration: 2000,
          });
          // 上传失败 清空已上传的文件列表
          (this.$refs["upsingle"] as any).clearFiles();
          this.$emit("fail");
        }
      })
      .catch((err: any) => {
        const { message } = err;
        this.progess_show = false;
        this.$message({
          message: message,
          type: "error",
          duration: 2000,
        });
        // 上传失败 清空已上传的文件列表
        (this.$refs["upsingle"] as any).clearFiles();
        this.$emit("fail");
      });
  }
}
</script>
<style lang="scss" scoped>
.dynamic-size {
  width: var(--width);
  height: var(--height);
}

.com-upload-image {
  @extend .dynamic-size;
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
      width: 60px !important;
      height: 60px !important;
      position: absolute;
      top: 0;
      left: 0;
      bottom: 0;
      right: 0;
      margin: auto;
    }
  }
}

::v-deep.avatar-uploader {
  width: 100%;
  height: 100%;
  display: inline-block;
  margin-right: 20px;
  vertical-align: top;
  margin-bottom: 10px;

  .el-upload {
    width: 100%;
    height: 100%;
    background-color: $color-white;
    border: 1px solid $color-second-line;
    border-radius: 4px;

    &:hover {
      border-color: $color-main;
      color: $color-main;
    }
  }
}

.no-img {
  @extend .dynamic-size;
  position: relative;
  line-height: calc(var(--height) * 0.74);

  .tips {
    position: absolute;
    font-size: 12px;
    width: 100%;
    height: 100%;
    text-align: center;
    line-height: calc(var(--height) + 20px);
    color: $color-text-second;
    top: 0;
  }

  i {
    font-size: 14px;
    color: $color-text-second;
  }
}

.img-info {
  width: 100%;
  height: 100%;
  display: inline-block;
  border-radius: 6px;
  position: relative;
  margin-right: 20px;
  margin-bottom: 10px;
  vertical-align: top;

  img {
    width: 100%;
    height: 100%;
    border-radius: 6px;
    object-fit: contain;
  }

  .delete-info {
    display: block;
    position: absolute;
    width: 100%;
    height: 100%;
    z-index: 2;
    background: rgba(0, 0, 0, 0);
    top: 0;
    left: 0;
    text-align: center;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.5s;
    overflow: hidden;

    &:hover {
      display: block;
      background: rgba(0, 0, 0, 0.5);
      transition: all 0.5s;
    }

    .only-delete {
      cursor: pointer;
      color: white;
      font-size: 20px;
      position: absolute;
      left: 30%;
      line-height: var(--height) !important;
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

.no-allowed {
  cursor: not-allowed;
  background: $color-black;
  opacity: 0.7;
}

.disabled-tips {
  color: $color-disable !important;
}

::v-deep.no-upload {
  .el-upload {
    cursor: not-allowed;
    background: $color-background;

    &:hover {
      border: 1px solid $color-second-line;
    }
  }
}</style>
<style lang="scss" scoped>.img-item {
  display: inline-block;
  margin-right: 10px;
}

.img-enter-active,
.img-leave-active {
  transition: all 1s;
}

.img-enter,
.img-leave-to

/* .list-leave-active for below version 2.1.8 */
  {
  opacity: 0;
  transform: translateY(30px);
}</style>
