<template>
  <div class="upload">
    <el-upload
      ref="upload"
      :type="type"
      :file-list="file_list"
      :multiple="multiple"
      :show-file-list="showFile"
      :action="action"
      :drag="drag"
      :disabled="disabled"
      :list-type="listType"
      :on-success="file => uploadSuccess(file)"
      :on-error="onError"
      :on-remove="() => removeFile()"
      :on-change="onChange"
      :headers="headers"
      :before-upload="file => uploadTodo(file)"
      :http-request="httpRequest"
    >
      <i
        v-if="!$slots.default && this.listType === 'picture-card'"
        class="el-icon-plus avatar-uploader-icon"
        aria-hidden="true"
      ></i>
      <div
        class="default"
        v-if="!$slots.default && this.listType != 'picture-card'"
      >
        <transition-group name="img-list">
          <div
            class="image-success m-r-10"
            v-for="(item, $index) in file_list"
            :key="item.url"
            @mouseover.prevent="item.is_hover = true"
            @mouseleave.prevent="item.is_hover = false"
          >
            <img
              alt=""
              :src="item.url"
              :style="`backgroundImage: url(${item.url})`"
            />
            <div class="mask" @click.stop="() => {}">
              <i
                v-show="item.is_hover"
                class="el-icon-delete delete-icon"
                aria-hidden="true"
                @click="remove($index)"
              ></i>
            </div>
          </div>
        </transition-group>
        <el-button
          class="upload-btn m-r-10"
          type="primary"
          plain
          v-if="type === 'image'"
          v-show="file_list.length < limit"
        >
          上传图片
        </el-button>
        <el-button
          class="upload-btn m-r-10"
          type="primary"
          plain
          v-if="type === 'file'"
          v-show="!showFile"
        >
          上传文件
        </el-button>
      </div>
      <!-- 自定义插槽部分 -->
      <slot></slot>
    </el-upload>
  </div>
</template>
<script lang="ts">
const fn = () => {
  return;
};
import { Component, Prop, Mixins, Watch } from "vue-property-decorator";
import UploadValidate from "@/mixins/UploadValidate";
import LoginAboat from "@/utils/LoginManager";
import { Loading } from "element-ui";
@Component({})
export default class CyUpload extends Mixins(UploadValidate) {
  @Prop({ default: "file" }) type?: string; // image file
  @Prop({ default: false, type: Boolean }) multiple!: boolean; //是否批量上传
  @Prop({ default: false, type: Boolean }) showFileList!: boolean;
  @Prop({
    default: `${process.env.VUE_APP_API_ROOT}${process.env.VUE_APP_API_CONSOLE}upload` //gw-console
  })
  action?: string;
  @Prop({ default: false, type: Boolean }) drag!: boolean; //是否开启拖拽上传
  @Prop({ default: fn }) onSuccess?: Function;
  @Prop({ default: fn }) onError?: Function;
  @Prop({ default: fn }) beforeUpload?: Function;
  @Prop({ default: fn }) onRemove?: Function;
  @Prop({ default: fn }) onChange?: Function;
  @Prop({ default: false, type: Boolean }) disabled!: boolean;
  @Prop() httpRequest?: Function;
  @Prop({
    default: () => {
      return {
        "Enterprise-Hash": LoginAboat.getAppHash(),
        "Auth-Token": LoginAboat.getToken()
      };
    }
  })
  headers?: object;
  @Prop({ default: "text" }) listType?: string;
  @Prop({ default: 1 }) limit?: number;
  @Prop({ default: () => [] }) fileList?: Array<any>;
  @Watch("fileList")
  onchangeList(val: Array<any>) {
    this.file_list = [...val];
  }
  mounted() {
    this.file_list = this.fileList || [];
    if (this.listType === "picture-card" && !this.$slots.default) {
      this.showFile = true;
    }
  }
  /**
   * data
   * */

  private file_list: Array<any> = [];
  private showFile = this.showFileList;
  private defaultFile = false;
  /**
   * methods
   * */
  public clearFiles() {
    return (this.$refs["upload"] as any).clearFiles();
  }
  // 上传前处理
  private uploadTodo(file: any) {
    if (this.beforeUpload) {
      return this.beforeUpload(file);
    } else {
      // 默认执行一些
      // if (this.type === "image") {
      //   return this._uploadValidate(file, [
      //     { type: "png", limit: 500 },
      //     { type: "jpg", limit: 500 },
      //     { type: "gif", limit: 1536 }
      //   ]);
      // }
    }
  }
  // 上传成功处理
  private uploadSuccess(file: any) {
    const { code, data } = file;
    if (code === 200) {
      this.$message({ message: "上传成功", type: "success", duration: 2000 });
      Loading.service({ fullscreen: true }).close();
    }
    // 没有自定义插槽并且为上传图片类型(default)
    if (!this.$slots.default && this.type === "image") {
      const img: any = {
        url: data.data,
        is_hover: false
      };
      this.file_list.push(img);
    }
    // 没有自定义插槽并且为上传文件类型
    if (!this.$slots.default && this.type === "file") {
      this.defaultFile = true;
      this.showFile = true;
    }
    this.onSuccess && this.onSuccess(file);
  }
  // 移除文件的处理
  private removeFile(i?: number) {
    this.onRemove && this.onRemove(i);
    // 没有自定义插槽并且为上传文件类型
    if (!this.$slots.default && this.type === "file") {
      this.showFile = false;
    }
  }
  // 删除
  private remove(i: number): void {
    this.file_list.splice(i, 1);
    this.removeFile(i);
  }
}
</script>
<style lang="scss" scoped>
.upload {
  display: inline-block;
  .default {
    vertical-align: top;
    display: inline-block;
    .class-btn {
      display: inline-block;
      vertical-align: top;
    }
    .image-success {
      width: 100px;
      height: 100px;
      border-radius: 8px;
      overflow: hidden;
      position: relative;
      display: inline-block;
      vertical-align: top;
      img {
        width: 100%;
        height: 100%;
        background-size: cover;
        background-position: center;
      }
      .mask {
        width: 100%;
        height: 100%;
        line-height: 100px;
        background-color: rgba(0, 0, 0, 0);
        position: absolute;
        top: 0;
        left: 0;
        &:hover {
          display: block;
          background: rgba(0, 0, 0, 0.5);
          transition: all 0.5s;
        }
        .delete-icon {
          color: #ffffff;
        }
      }
    }
  }
}
.img-list-item {
  display: inline-block;
  margin-right: 10px;
}
.img-list-enter-active,
.img-list-leave-active {
  transition: all 1s;
}
.img-list-enter, .img-list-leave-to
/* .list-leave-active for below version 2.1.8 */ {
  opacity: 0;
  transform: translateY(30px);
}
</style>
