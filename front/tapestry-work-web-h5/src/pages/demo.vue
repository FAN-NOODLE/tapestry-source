<template>
  <div class="container">
    <!-- <div class="title m-b-15">链接直接点击</div>
    <a class="btn m-b-15" :href="url['excel'].url" :download="url['excel'].name">
      <van-button type="primary" block>下载excel</van-button>
    </a>
    <a class="btn m-b-15" :href="url['word'].url" :download="url['word'].name">
      <van-button type="primary" block>下载word</van-button>
    </a>
    
    <a class="btn m-b-15" :href="url['ppt'].url" :download="url['ppt'].name">
      <van-button type="primary" block>下载ppt</van-button>
    </a>
    
    <a class="btn m-b-15" :href="url['pdf'].url" :download="url['pdf'].name">
      <van-button type="primary" block>下载pdf</van-button>
    </a>
    
    <a class="btn m-b-15" :href="url['img'].url" :download="url['img'].name">
      <van-button type="primary" block>下载图片</van-button>
    </a>
    
    <div class="title m-b-15">预览下载</div>
    <van-button class="m-b-15" type="primary" block @click="preview('word')">预览下载word</van-button>
    <van-button class="m-b-15" type="primary" block @click="preview('excel')">预览下载excel</van-button>
    <van-button class="m-b-15" type="primary" block @click="preview('ppt')">预览下载ppt</van-button>
    <van-button class="m-b-15" type="primary" block @click="preview('pdf')">预览下载pdf</van-button>
    <van-button class="m-b-15" type="primary" block @click="preview('img')">预览下载图片</van-button> -->

    <div class="title m-b-15">选择图片/视频</div>
    <van-button class="btn m-b-15" type="primary" block>
      选择图片或拍照
      <input class="input" type="file" accept="image/*" />
    </van-button>
    <van-button class="btn m-b-15" type="primary" block>
      拍照
      <input class="input" type="file" accept="image/*" capture="user" />
    </van-button>
    <van-button class="btn m-b-15" type="primary" block>
      选择视频或拍摄
      <input class="input" type="file" accept="video/*" />
    </van-button>
    <van-button class="btn m-b-15" type="primary" block>
      拍摄
      <input class="input" type="file" accept="video/*" capture="user" />
    </van-button>

    <van-button class="btn m-b-15" type="primary" block>
      选择文件
      <input class="input" type="file" accept="*" />
    </van-button>
    <van-button class="btn m-b-15" type="primary" block>
      拍照
      <input class="input" type="file" accept="*" capture="user" />
    </van-button>
    <div class="title m-b-15">企微API</div>
    <van-button class="btn m-b-15" type="primary" block @click="handleChooseImage">
      选择图片
    </van-button>
    <van-button class="btn m-b-15" type="primary" block @click="handleChooseVideo">
      选择视频
    </van-button>
    <!-- <img style="width: 200px; height: 200px;" src="https://wework.qpic.cn/wwpic3az/67042_UpxRwYHcSIOruDL_1751421317/0" alt=""> -->
    <van-button class="btn m-b-15" type="primary" block @click="handleGetLocation()">
      获取位置信息
    </van-button>
    <van-button class="btn m-b-15" type="primary" block @click="handleGetLocation()">
      获取位置信息不获取config
    </van-button>
    <van-button class="btn m-b-15" type="primary" block @click="show = true">
      显示popup
    </van-button>
    <van-button class="btn m-b-15" type="primary" block @click="handlShareToExternalMoments">
      发送朋友圈
    </van-button>




    <!-- 签名组件 -->
    <van-button class="btn m-b-15" type="primary" block @click="show = true">
      打开签名弹框
    </van-button>
    <van-popup v-model="show" :closeable="true" class="popup-container">
      <div class="signature-wrapper">
        <h3 class="signature-title">请在下方区域签名</h3>
        <div class="canvas-container">
          <canvas 
            id="canvas" 
            ref="canvas" 
            width="300" 
            height="200"
            @touchstart="startDrawing"
            @touchmove="draw"
            @touchend="stopDrawing"
          ></canvas>
        </div>
        <div class="signature-buttons">
          <van-button type="default" block @click="clearCanvas">清空</van-button>
          <van-button type="primary" block @click="saveSignature">保存</van-button>
        </div>
        <div v-if="signatureImage" class="signature-result">
          <h4>签名预览：</h4>
          <img :src="signatureImage" alt="签名结果" />
        </div>
      </div>
    </van-popup>


    <van-button class="btn m-b-15" type="primary" block @click="handleShowInput">
      显示输入框
    </van-button>

    <div class="input-container" :class="{ 'input-container--show': showInput }">
      <input ref="input" type="text" placeholder="请输入内容">
    </div>

    <van-button class="btn m-b-15" type="primary" block @click="handleChooseMaterial">
      打开素材库
    </van-button>
  </div>

</template>
<script lang="ts">
import { Component, Vue, Watch } from "vue-property-decorator";
import { downloadFile } from '@/utils/download';
import wxUtil from '@/utils/wx';

@Component
export default class Demo extends Vue {
  show = false;
  signatureImage = '';
  isDrawing = false;
  lastX = 0;
  lastY = 0;
  ctx: CanvasRenderingContext2D | null = null;
  
  url = {
    word: {
      url: 'https://tapestry-uat-1308608494.cos.ap-shanghai.myqcloud.com/uploads/2023/06/05/7XjRFBnRAyt5RSAWnDsset7cjHFd2ies.docx',
      name: '蔻驰测试文档.docx',
      size: 43469
    },
    ppt: {
      url: 'https://tapestry-uat-1308608494.cos.ap-shanghai.myqcloud.com/uploads/2023/010/19/zp6Z3kiphM6GHsxQy7Sw7TFi876am5nM.pptx',
      name: "清新植物叶子PPT模板.pptx",
      size: 3670296
    },
    excel: {
      url: 'https://tapestry-uat-1308608494.cos.ap-shanghai.myqcloud.com/uploads/2023/010/19/AsJ6GNyT8j2FpR6FTZNwMPjDDwkpwQSN.xlsx',
      name: "企微素材库文件操作.xlsx",
      size: 9142
    },
    pdf: {
      url: 'https://tapestry-uat-1308608494.cos.ap-shanghai.myqcloud.com/uploads/2023/05/30/2pKkMGhJbt3NcZKjcicrMHnQZCSket88.pdf',
      name: "正大方略工程咨询有限公司.pdf",
      size: 2496302
    },
    img: {
      url: 'https://tapestry-uat-1308608494.cos.ap-shanghai.myqcloud.com/uploads/2023/10/23/a3jMEb5ECpz3dYf42TkPMZNJGiYwCt4T.jpg',
      name: 'test.jpg',
      size: 281616
    }
  };

  showInput = false;

  @Watch('show')
  onShowChange(newVal: boolean, oldVal: boolean) {
    if (newVal) {
      this.$nextTick(() => {
        this.initCanvas();
      });
    }
  }

  handleShowInput() {
    this.showInput = true;
    setTimeout(() => {
      const input = this.$refs.input as any;
      input.focus();
    }, 1000);
  }
  
  download(type: string) {
    try {
      downloadFile(this.url[type].url, {
        filename: 'test',
      });
    } catch (error) {
      console.log(error);
    }
  }
  
  preview(type: string) {
    const current = this.url[type];
    window.wx.invoke("previewFile", {
      ...current
    }, (err: any) => {
      console.log(err);
    });
  }
  
  handleChooseImage() {
    wxUtil.chooseImage().then(res => {
      console.log(res);
    }).catch(err => {
      console.error(err);
    });
  }
  
  handleChooseVideo() {
    wxUtil.chooseVideo().then(res => {
      console.log(res);
    }).catch(err => {
      console.error(err);
    });
  }

  handleGetLocation() {
    wxUtil.getLocation().then(res => {
      console.log(res);
    }).catch(err => {
      console.error(err);
    });
  }
  
  handlShareToExternalMoments() {
    wxUtil.shareToExternalMoments().then(res => {
      console.log(res);
    }).catch(err => {
      console.error(err);
    });
  }
  
  // 签名组件方法
  mounted() {
    
  }
  
  initCanvas() {
    const canvas = this.$refs.canvas as HTMLCanvasElement;
    this.ctx = canvas.getContext('2d');
    if (this.ctx) {
      this.ctx.lineWidth = 2;
      this.ctx.lineCap = 'round';
      this.ctx.strokeStyle = '#000';
      // 设置画布背景为白色
      this.ctx.fillStyle = '#fff';
      this.ctx.fillRect(0, 0, canvas.width, canvas.height);
      // 添加边框
      this.ctx.strokeStyle = '#ccc';
      this.ctx.strokeRect(0, 0, canvas.width, canvas.height);
      this.ctx.strokeStyle = '#000';
    }
  }
  
  startDrawing(e: MouseEvent | TouchEvent) {
    this.isDrawing = true;
    
    let clientX: number;
    let clientY: number;
    
    if (e instanceof TouchEvent) {
      e.preventDefault();
      clientX = e.touches[0].clientX;
      clientY = e.touches[0].clientY;
    } else {
      clientX = e.clientX;
      clientY = e.clientY;
    }
    
    const canvas = this.$refs.canvas as HTMLCanvasElement;
    const rect = canvas.getBoundingClientRect();
    this.lastX = clientX - rect.left;
    this.lastY = clientY - rect.top;
  }
  
  draw(e: MouseEvent | TouchEvent) {
    if (!this.isDrawing || !this.ctx) return;
    
    let clientX: number;
    let clientY: number;
    
    if (e instanceof TouchEvent) {
      e.preventDefault();
      clientX = e.touches[0].clientX;
      clientY = e.touches[0].clientY;
    } else {
      clientX = e.clientX;
      clientY = e.clientY;
    }
    
    const canvas = this.$refs.canvas as HTMLCanvasElement;
    const rect = canvas.getBoundingClientRect();
    const currentX = clientX - rect.left;
    const currentY = clientY - rect.top;
    
    this.ctx.beginPath();
    this.ctx.moveTo(this.lastX, this.lastY);
    this.ctx.lineTo(currentX, currentY);
    this.ctx.stroke();
    
    this.lastX = currentX;
    this.lastY = currentY;
  }
  
  stopDrawing() {
    this.isDrawing = false;
  }
  
  clearCanvas() {
    const canvas = this.$refs.canvas as HTMLCanvasElement;
    if (this.ctx) {
      this.ctx.clearRect(0, 0, canvas.width, canvas.height);
      // 重新设置背景和边框
      this.ctx.fillStyle = '#fff';
      this.ctx.fillRect(0, 0, canvas.width, canvas.height);
      this.ctx.strokeStyle = '#ccc';
      this.ctx.strokeRect(0, 0, canvas.width, canvas.height);
      this.ctx.strokeStyle = '#000';
      this.signatureImage = '';
    }
  }
  
  saveSignature() {
    const canvas = this.$refs.canvas as HTMLCanvasElement;
    this.signatureImage = canvas.toDataURL('image/png');
    console.log('签名结果:', this.signatureImage);
    // 这里可以添加保存到服务器或其他处理逻辑
    this.$toast('签名已保存');
  }

  handleChooseMaterial() {
    window.open('https://we.coach.com.cn/work/f7c9d26cf6441e9dbad633bab7be9af6/material-detail?debug=1&upc_code=198685101430');
  }
}
</script>
<style lang="scss" scoped>
.m-b-15 {
  margin-bottom: 15px;
}
.container {
  background-color: #fff;
  min-height: 100vh;
  padding: 20px;
  font-size: 14px;

  .title {
    font-size: 14px;
  }

  .btn {
    display: block;
    position: relative;
    .input {
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      opacity: 0;
      z-index: 1;
    }
  }
}
.popup-container {
  padding: 20px;
  width: 90%;
  max-width: 400px;
}

.signature-wrapper {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.signature-title {
  font-size: 18px;
  font-weight: bold;
  text-align: center;
  margin: 0;
}

.canvas-container {
  border: 1px solid #ddd;
  border-radius: 4px;
  overflow: hidden;
  background-color: #f9f9f9;
}

#canvas {
  display: block;
  width: 100%;
  height: auto;
  cursor: crosshair;
  background-color: white;
}

.input-container {
  position: fixed;
  bottom: 0;
  left: 0;
  width: 100%;
  padding: 10px;
  box-sizing: border-box;
  visibility: hidden;
  pointer-events: none;

  &--show {
    visibility: visible;
    pointer-events: auto;
  }

  input {
    width: 100%;
    height: 40px;
    padding: 10px;
    box-sizing: border-box;
    border: 1px solid #ccc;
    border-radius: 4px;
    font-size: 14px;
  }
}

.signature-buttons {
  display: flex;
  gap: 10px;
}

.signature-result {
  margin-top: 10px;
  padding: 10px;
  border: 1px solid #eee;
  border-radius: 4px;
  background-color: #f9f9f9;
}

.signature-result img {
  max-width: 100%;
  height: auto;
  display: block;
  margin-top: 10px;
}
</style>