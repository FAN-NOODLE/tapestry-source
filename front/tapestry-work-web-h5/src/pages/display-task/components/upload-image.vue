<template>
  <van-uploader
    v-model="list"
    :preview-image="false"
    :max-count="maxLength"
    :multiple="multiple"
    :disabled="disabled"
    :accept="accept"
    :max-size="maxSize"
    :before-read="beforeRead"
    :after-read="handleAfterRead"
  >
    <slot>
      <div class="upload">
        <!-- <div class="upload-image"></div>
        <div class="upload-text">上传图片</div> -->
      </div>
    </slot>
  </van-uploader>
</template>

<script lang="ts">
import { Component, Model, Prop, Vue } from "vue-property-decorator";
import { Uploader } from 'vant';
import { cosUploader } from '@/utils/upload'
const isArray = (value: any) => Object.prototype.toString.call(value).includes("Array")
@Component({
  components: {
    [Uploader.name]: Uploader
  }
})
export default class UploadImage extends Vue {
  @Model('change', { type: Array }) value: [];
  @Prop({ type: Number, default: 10 }) maxLength: number;
  @Prop({ type: Number, default: 10 * 1024 * 1024 }) maxSize: number;
  @Prop({ type: Boolean, default: true }) multiple: boolean;
  @Prop({ type: Boolean, default: false }) disabled: boolean;
  @Prop({ type: String, default: 'image/jpg,image/jpeg,image/png' }) accept: string;

  get list(): any[] {
    return this.value || [];
  }
  set list(value) {
    value.forEach(item => {
      item._content = item.content || '';
      item._upload_status = typeof item.status === 'string' ? item.status : item._upload_status || '';
      item._message = item.message || '';
      item.comment = item.comment || [];
      item.num = item.num || 0;
      item.vm_num = item.vm_num || 0;
      delete item.content;
      delete item.message;
      // if(typeof item.status !== 'string') {
      //   delete item.status;
      // }

    })
    this.$emit('change', value);
  }
 // 返回布尔值
  beforeRead(file: any) {
    if(isArray(file)) {
      for(let i = 0; i < file.length; i++) {
        if(file[i].size > this.maxSize) {
          this.$toast.fail('图片大小不能超过' + 10 + 'MB');
          return false;
        }
      }
    } else {
      if(file.size > this.maxSize) {
        this.$toast.fail('图片大小不能超过' + 10 + 'MB');
        return false;
      }
    }
    return true;
  }
  /**
   * 获取文件后，执行上传文件
   * @param file 
   */
  async handleAfterRead(file: any) {
    if(isArray(file)) {
      await Promise.all(file.map(item => this.uploadImage(item)));
    } else {
      await this.uploadImage(file);
    }
  }

  /**
   * 上传图片
   */
  uploadImage(file: any) {
    file._upload_status = 'uploading';
    file._message = '上传中...';
    return new Promise((resolve, reject) => {
      cosUploader.uploadFile(file.file).then((res: any) => {
        file.url = res.url
        file._upload_status = 'done';
        file._message = '上传成功';
        // file._content = ''
        resolve(file);
      }).catch(err => {
        file._upload_status = 'fail';
        file._message = '上传失败';
        console.log(err);
      })
    });
  }
}
</script>
<style lang="scss" scoped>
.upload {
  width: 100px;
  height: 100px;
  background: #F6F6F6;
  // border: 1px dashed #A6A6A6;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  margin-bottom: 15px;
  box-sizing: border-box;
  position: relative;

  background: url("~@/assets/display-task/icon_upload_image.png") center/100% no-repeat;
  // &::after {
  //   content: '';
  //   position: absolute;
  //   top: 0;
  //   left: 0;
  //   width: 100%;
  //   height: 100%;
  //   border: 1px dashed #A6A6A6;
  //   box-sizing: border-box;
  // }
  // &-image {
  //   width: 21px;
  //   height: 18px;
  //   background: url("~@/assets/display-task/icon_photo.png") center/100% no-repeat;
  // }

  // &-text {
  //   font-size: 12px;
  //   font-weight: 400;
  //   color: #9E9E9E;
  //   line-height: 17px;
  //   margin-top: 10px;
  // }
}
  
</style>