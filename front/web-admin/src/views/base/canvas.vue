<template>
  <div class="container">
    <!-- <textarea
      v-model="textareaValue"
      placeholder="请输入内容"
      rows="5"
      class="textarea"
    ></textarea> -->
    <!-- 截图区域 -->
    <div ref="captureArea">
      <!-- 一个图标 -->
      <!-- <div ref="captureArea" class="icon">
        <img src="@/assets/images/upload_local_image_2.png" alt="">
      </div> -->
      <div class="upload">
        <div class="upload-image"></div>
        <div class="upload-text">上传图片</div>
      </div>
    </div>

    <!-- 操作按钮 -->
    <el-button type="primary" @click="handleCapture">截图并下载</el-button>
  </div>
</template>

<script lang="ts">
import { download } from "@/utils/download";
import { Vue, Component } from "vue-property-decorator";
import html2canvas from '@whowasme/html2canvas';

@Component({})
export default class extends Vue {
  /** 文本域值 */
  textareaValue = ''
  // 截图并下载
  async handleCapture() {
    const el = this.$refs.captureArea as HTMLElement;
    if (!el) return;

    try {
      const canvas = await html2canvas(el, {
        useCORS: true, // 支持跨域图片
        scale: 2,      // 放大倍数，提高清晰度
        backgroundColor: 'transparent'
      });

      // 将 canvas 转为 blob 并下载
      canvas.toBlob((blob) => {
        if (blob) {
          download(blob, '', 'icon.png');
        }
      });
    } catch (error) {
      console.error('截图失败:', error);
      this.$message.error('截图失败，请重试');
    }
  }
}
</script>

<style lang="scss" scoped>
.container {
  width: 100vw;
  height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 24px;

  .textarea {
    width: 800px;
    height: 200px;
    padding: 12px;
    border-radius: 10px;
    border: 1px solid #000;
    font-size: 16px;
    display: block;
  }
}


.icon {
  width: 100px;
  height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  overflow: hidden;
  background-color: #000;
  img {
    width: 60px;
    height: 60px;
  }
}

.demo-box {
  display: flex;
  justify-content: center;
  gap: 16px;
  font-size: 18px;
}

.upload {
  width: 100px;
  height: 100px;
  background: #F6F6F6;
  // border: 1px dashed #A6A6A6;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  box-sizing: border-box;
  position: relative;
  &::after {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    border: 1px dashed #A6A6A6;
    box-sizing: border-box;
  }
  &-image {
    width: 21px;
    height: 18px;
    background: url("~@/assets/images/icon_photo.png") center/100% no-repeat;
  }

  &-text {
    font-size: 12px;
    font-weight: 400;
    color: #9E9E9E;
    line-height: 17px;
    margin-top: 10px;
  }
}
</style>
