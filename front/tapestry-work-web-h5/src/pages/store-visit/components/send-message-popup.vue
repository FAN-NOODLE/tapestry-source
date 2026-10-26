<template>
  <van-popup v-model="visible" position="bottom" round @click-overlay="handleCancel">
    <div class="popup">
      <div class="popup-content">
        <div class="wrap">
          <textarea v-model="content" ref="input" class="input" type="text" :placeholder="String($t(placeholder))" @click="handleInputClick"> </textarea>
          <div v-if="!onlyText" class="upload">
            <com-upload v-model="fileList" accept="image/jpg,image/jpeg,image/png,video/*" @upload-success="handleUploadSuccess">
              <div class="icon-img"></div>
            </com-upload>
          </div>
        </div>
      </div>
      <div class="popup-footer">
        <div class="list">
          <div class="item" v-for="item, index in fileList" :key="index">
            <van-image v-if="item.media_type === 1" :src="item.content"></van-image>
            <video v-if="item.media_type === 2" :src="item.content" :poster="item.video_image_url"></video>
            <div class="icon-close-circle" @click="handleDeleteFile(index)"></div>
          </div>
        </div>
        <div class="add" :class="{disabled: disabled}" @click="handleConfirm">{{ String($t(confirmText)) }}</div>
      </div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Model, Prop, Vue, Watch } from 'vue-property-decorator';
import ComUpload from "./com-upload.vue";
import { getVideoCover, getVideoDuration } from '@/utils/upload/utils';
import { base64ToImage } from '@/utils/tools';
import { cosUploader } from '@/utils/upload';
import { getFileData } from '@/utils/download';

@Component({
  components: {
    ComUpload,
  },
})
export default class SendMessagePopup extends Vue {
  @Model('update:show', { type: Boolean, default: false }) show: boolean;
  @Prop({ type: String, default: 'storeVisit.tips.24' }) placeholder: string;
  @Prop({ type: String, default: 'storeVisit.popup.btn.19' }) confirmText: string;
  @Prop({ type: Boolean, default: false }) onlyText: boolean;
  @Prop({ type: String, default: '' }) defaultValue: string;
  visible: boolean = false;
  
  content = ''
  fileList = []

  get disabled () {
    const hasUploading = this.fileList.some(item => item.status === 'uploading');
    return !this.content.trim() && !this.fileList.length || hasUploading;
  }

  get p() {
    return this.placeholder
  }

  @Watch('show', { immediate: true })
  onShowChange(val: boolean) {
    this.visible = val;
    if(val) {
      this.$nextTick(() => {
        (this.$refs.input as any).focus();
        (this.$refs.input as any).scrollIntoView({ behavior: 'smooth', block: 'center' });
      });
      this.content = this.defaultValue;
    }
  }

  @Watch('defaultValue', { immediate: true })
  onDefaultValueChange(val: string) {
    this.content = val;
  }
  /**
   * 点击获取焦点
   */
  handleInputClick() {
    (this.$refs.input as any).focus()
  }

  /**
   * 取消
   */
  handleCancel() {
    this.content = '';
    this.fileList = [];
    this.$emit('update:show', false);
    this.$emit('cancel');
  }


  /**
   * 确认
   */
  handleConfirm = this.$helper.debounce(() => {
    // 触发确认事件
    this.emitConfirm();
  }, 500, true)
  
  /**
   * 确认事件
   */
  emitConfirm() {
    if(this.disabled) return
    // this.$emit('update:show', false);
    this.$emit('confirm', {
      content: this.content,
      media_list: this.fileList.map(item => {
        return {
          url: item.url,
          media_type: item.media_type,
          video_second: item.video_second,
          video_image_url: item.video_image_url
        }
      })
    });
    this.content = '';
    this.fileList = [];
  }


  handleDeleteFile(index) {
    this.fileList.splice(index, 1);
  }

  /**
   * 上传成功后
   * @param value 
   */
  async handleUploadSuccess(e) {
    console.log(e);
  }
}
</script>
<style lang="scss" scoped>
.icon-img {
  width: 16px;
  height: 16px;
  background: url('~@/assets/store-visit/icon_img.png') center/100% no-repeat;
}
.icon-close-circle {
  width: 16px;
  height: 16px;
  background: url('~@/assets/store-visit/icon_close_circle.png') center/100% no-repeat;
}

.popup {
  position: relative;
  width: 100vw;
  background: #FFFFFF;
  .popup-content {
    padding: 16px 16px 8px;
    .wrap {
      height: 84px;
      box-sizing: border-box;
      padding: 12px 16px;
      border-radius: 12px;
      background-color: rgba(233, 236, 239, 1);
      display: flex;
      flex-direction: column;
    }

    .input {
      width: 100%;
      flex: 1;
      line-height: 18px;
      font-size: 12px;
      padding: 0;
      background-color: transparent;

      &::placeholder {
        color: rgba(153, 153, 153, 0.5);
      }
      &:focus {
        outline: none;
      }
    }

    .upload {
      width: 16px;
      height: 16px;
      display: flex;
      align-items: center;
      justify-content: center;
    }
  }

  .popup-footer {
    display: flex;
    justify-content: space-between;

    .list {
      flex: 1;
      padding: 12px 16px;

      .item {
        width: 40px;
        height: 40px;
        border-radius: 8px;
        overflow: hidden;
        position: relative;

        .icon-close-circle {
          position: absolute;
          top: 0;
          right: 0;
        }

        .img {
          width: 100%;
          height: 100%;
        }
      }
    }

    .add {
      height: 40px;
      padding: 0;
      color: rgba(13, 110, 253, 1);
      font-size: 12px;
      line-height: 40px;
      font-weight: 600;
      padding-right: 20px;
      &.disabled {
        color: rgba(153, 153, 153, 1);
      }
    }
  }
}
</style>
