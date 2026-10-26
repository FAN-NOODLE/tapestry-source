<template>
  <div>
    <div v-if="useSheet" @click="visible = true">
      <slot></slot>
    </div>
    <van-uploader
      v-else
      ref="uploader"
      v-model="list"
      :preview-image="previewImage"
      :max-count="maxLength"
      :multiple="multiple"
      :disabled="disabled"
      :accept="accept"
      :max-size="maxSize * 1024 * 1024"
      :before-read="beforeRead"
      :after-read="handleAfterRead"
    >
      <slot>
        <div class="upload">
          <div class="upload-image"></div>
          <div class="upload-text">上传图片</div>
        </div>
      </slot>
    </van-uploader>

    <van-popup v-model="visible" position="bottom" round @click-overlay="handleCancel">
      <div class="popup">
        <div class="popup-content">
          <div class="item" @click="visible = false">
            <van-uploader
              ref="uploader"
              v-model="list"
              :preview-image="previewImage"
              :max-count="maxLength"
              :multiple="multiple"
              :disabled="disabled"
              :accept="'image/jpg,image/jpeg,image/png,video/*'"
              :max-size="maxSize * 1024 * 1024"
              :before-read="beforeRead"
              :after-read="handleAfterRead"
            >
              <div class="item-text">{{ $t('storeVisit.popup.btn.21') }}</div>
            </van-uploader>
          </div>
          <div class="item" @click="visible = false">
            <van-uploader
              ref="uploader"
              v-model="list"
              :preview-image="previewImage"
              :max-count="maxLength"
              :multiple="multiple"
              :disabled="disabled"
              :accept="'video/*'"
              capture="camera"
              :max-size="maxSize * 1024 * 1024"
              :before-read="beforeRead"
              :after-read="handleAfterRead"
            >
              <div class="item-text">{{ $t('storeVisit.popup.btn.22') }}</div>
            </van-uploader>
          </div>
        </div>
        <div class="popup-footer">
          <div class="item" @click="handleCancel">{{ $t('storeVisit.popup.btn.3') }}</div>
        </div>
      </div>
    </van-popup>
  </div>
</template>

<script lang="ts">
import { Component, Model, Prop, Vue, Watch } from "vue-property-decorator";
import { Uploader } from 'vant';
import { cosUploader } from '@/utils/upload'
import { getVideoDuration } from "@/utils/upload/utils";
import { getFileData } from "@/utils/download";
const isArray = (value: any) => Object.prototype.toString.call(value).includes("Array")
@Component({
  components: {
    [Uploader.name]: Uploader
  }
})
export default class ComUpload extends Vue {
  @Model('change', { type: Array }) value: [];
  @Prop({ type: Number, default: 10 }) maxLength: number;
  @Prop({ type: Number, default: 100 }) maxSize: number;
  @Prop({ type: Boolean, default: true }) multiple: boolean;
  @Prop({ type: Boolean, default: false }) disabled: boolean;
  @Prop({ type: Boolean, default: false }) previewImage: boolean;
  @Prop({ type: String, default: 'image/*' }) accept: string;
  @Prop({ type: Boolean, default: false }) useSheet: boolean;

  visible = false

  get list(): any[] {
    return this.value || [];
  }
  set list(value) {
    this.$emit('change', value);
  }
 // 返回布尔值
  beforeRead(file: any) {
    // console.log(file.size, this.maxSize);
    
    if(isArray(file)) {
      for(let i = 0; i < file.length; i++) {
        if(file[i].size > this.maxSize * 1024 * 1024) {
          this.$toast(String(this.$t('common.toast.20')) + this.maxSize + 'MB');
        }
      }
    } else {
      if(file.size > this.maxSize * 1024 * 1024) {
        this.$toast(String(this.$t('common.toast.20')) + this.maxSize + 'MB');
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
      await Promise.all(file.map(item => this.uploadFile(item)));
    } else {
      await this.uploadFile(file);
    }
  }

  /**
   * 上传图片
   */
  async uploadFile(file: any) {
    console.log(file);
    
    this.$set(file, 'progress', 0)
    this.$set(file, 'status', 'uploading')
    this.$set(file, 'message', '上传中...')
    this.$set(file, 'media_type', file.file.type.includes('image')?1:2)
    const d = file.media_type == 2 ? 0.9 : 1
    try  {
      const res:any = await cosUploader.uploadFile(file.file, {
        onProgress: (progress: any) => {
          file.progress = progress.percent * d;
        }
      })

      if(file.media_type === 1) {
        // 图片
        file.url = res.url
        file.status = 'done';
        file.message = '上传成功';
      } else {
        // 视频
        file.url = res.url
        // 获取视频时长
        const duration = await getVideoDuration(file.file)
        // 视频封面
        const buffer:any = await getFileData(`${res.url}?ci-process=snapshot&time=1&format=jpg`)
        // 上传视频封面
        const res2:any = await cosUploader.uploadFile(buffer, { filename: res.url.split('/').pop().split('.')[0] + '_poster.jpg'})
        
        const blob = new Blob([buffer], { type: 'image/png' })
        this.$set(file, 'poster', URL.createObjectURL(blob))
        this.$set(file, 'video_second', duration)
        this.$set(file, 'video_image_url', res2.url)
        
        file.progress = 1
        file.status = 'done';
        file.message = '上传成功';
      }


      this.visible = false
      this.$emit('upload-success', file)


      return file
      
    } catch (error) {
      console.log(error);
      file.status = 'failed';
      file.message = '上传失败';
    }
  }


  handleCancel() {
    this.visible = false
    this.$emit('cancel')
  }
}
</script>
<style lang="scss" scoped>
.upload {
  width: 100px;
  height: 100px;
  background: #F6F6F6;
  border: 1px dashed #A6A6A6;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  margin-bottom: 15px;
  &-image {
    width: 21px;
    height: 18px;
    background: url("~@/assets/display-task/icon_photo.png") center/100% no-repeat;
  }

  &-text {
    font-size: 12px;
    font-weight: 400;
    color: #9E9E9E;
    line-height: 17px;
    margin-top: 10px;
  }
}
  
.popup {
  width: 100vw;
  position: relative;
  background: #FFFFFF;
  
  .item {
    height: 68px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
    font-weight: 600;
    color: rgba(0, 0, 0, 1);
    &:not(:last-child) {
      border-bottom: 1px solid rgba(248, 249, 250, 1);
    }
    .item-text {
      height: 68px;
      width: 100vw;
      display: flex;
      align-items: center;
      justify-content: center;
    }
  }

  .popup-content {
    padding: 0 24px;
  }

  .popup-footer {
    border-top: 4px solid rgba(248, 249, 250, 1);
    padding-bottom: calc(env(safe-area-inset-bottom) + 12px);
    padding-bottom: calc(constant(safe-area-inset-bottom) + 12px);


    .item {
      font-weight: normal;
    }
  }
}
</style>