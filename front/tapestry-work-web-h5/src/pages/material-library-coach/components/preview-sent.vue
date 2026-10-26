<template>
  <div class="confirm-material-box">
    <ul class="msg-ui">
      <li class="msg-li" v-for="(item, index) in chooseArr" :key="index">
        <!-- 删除按钮 -->
        <img
          :src="require('@/assets/msg-del.png')"
          class="delete-icon"
          @click="deleteFile(index)"
        />
        <!-- 发送主体 -->
        <div class="msg-cont">
          <!-- 文本 -->
          <div class="text" v-if="item.type === 1">
            <van-field
              v-model="item.content"
              rows="2"
              class="input-field"
              autosize
              type="textarea"
              readonly
            />
            <div class="point-icon">
              <van-icon name="play" />
            </div>
          </div>
          <!-- 图片 -->
          <div class="image" v-else-if="item.type === 2">
            <img :src="item.url" @click="prewImg(item.url)" />
          </div>
          <!-- 视频 -->
          <div class="video" v-else-if="item.type === 5">
            <img :src="require('@/assets/video-bg.png')" />
          </div>
          <!-- 链接 -->
          <div class="link" v-else-if="item.type === 3">
            <div class="link-name">{{ item.title }}</div>
            <div class="desc-box">
              <div class="desc-text">{{ item.description }}</div>
              <div class="link-pic">
                <img :src="item.cover_url" />
              </div>
            </div>
          </div>
          <!-- 小程序 -->
          <div class="mini-program" v-else-if="item.type === 4">
            <div class="title">
              <img :src="require('@/assets/mini.png')" />
              <div class="text-div">
                {{ item.app_name }}
              </div>
            </div>
            <p class="name">{{ item.title }}</p>
            <div class="code-pic">
              <img :src="item.cover_url" />
            </div>
            <div class="bottom">小程序</div>
          </div>
          <!-- 文件 -->
          <div class="files" v-else>
            <div class="left-box">
              <p class="file-name">
                {{ item.name }}
              </p>
              <p class="file-size">
                {{ (item.size / 1024 / 1024).toFixed(2) }}M
              </p>
            </div>
            <div class="right-box">
              <img :src="item.file_type | imgUrl" />
            </div>
          </div>
        </div>
        <!-- 头像 -->
        <img :src="require('@/assets/default/brand-logo.png')" class="avatar" />
      </li>
    </ul>
    <div class="btn-box">
      <div class="back-btn" @click="backPage">返回</div>
      <div class="submit-btn" @click="combinationSend">发送</div>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import { ImagePreview, Icon } from "vant";
import { sourceMakeMediaApi,reportTimes } from "@/api/marketing-promotion/index";
@Component({
  components: {
    [Icon.name]: Icon,
  },
  filters: {
    imgUrl: (value) => {
      const wordBj = require("@/assets/material-library/Word.png");
      const excelBj = require("@/assets/material-library/Excel.png");
      const pptBj = require("@/assets/material-library/PPT.png");
      const pdfBj = require("@/assets/material-library/PDF.png");
      const otherBj = require("@/assets/material-library/Other.png");
      switch (value) {
        case 1:
          return wordBj;
        case 2:
          return excelBj;
        case 3:
          return pptBj;
        case 4:
          return pdfBj;
        default:
          return otherBj;
      }
    },
  },
})
export default class extends Vue {
  @Prop({ default: [] }) chooseArr: Array<any>;
  @Prop({ default: 1 }) is_official: number;
  // 记录本次已发送素材
  public get sentData(): any {
    return this.$store.getters["materialLibrary/getNowSentData"];
  }
  //判断被选中的数据的类型
  private fileType(type: any) {
    let text: any = "";
    switch (type) {
      case 1:
        text = "text";
        break;
      case 2:
        text = "image";
        break;
      case 3:
        text = "link";
        break;
      case 4:
        text = "miniprogram";
        break;
      case 5:
        text = "video";
        break;
      default:
        text = "file";
        break;
    }
    return text;
  }
  private backPage() {
    this.$emit("backPage");
    // this.$router.push("/material-library");
  }
  // 删除文件
  private deleteFile(index: number) {
    if (this.chooseArr.length <= 1) {
      this.$toast("至少保留一项");
      return;
    }
    this.chooseArr.splice(index, 1);
    sessionStorage.setItem("chooseMaterial", JSON.stringify(this.chooseArr));
  }
  // 预览图片
  private prewImg(url) {
    let arr = [url];
    ImagePreview({
      images: arr,
      startPosition: 0,
    });
  }
  // 组合发送
  private combinationSend() {
    this.$emit('send')
  }
}
</script>
<style lang="scss" scoped>
.confirm-material-box {
  background: #f6f7f9;
  height: 100vh;
  overflow: auto;
  .msg-ui {
    padding: 39px 10px;
    .msg-li {
      display: flex;
      // align-items: center;
      position: relative;
      margin-bottom: 20px;
      .msg-cont {
        width: 256px;
        margin-left: 18px;
        .text {
          border-radius: 8px;
          width: 256px;
          background: #cbe6fb;
          padding: 15px 10px;
          font-size: 14px;
          font-weight: 400;
          color: #4a4a4a;
          line-height: 20px;
          position: relative;
          .input-field {
            font-size: 14px;
            width: 100%;
            background: #cbe6fb;
            padding: 0;
            &::after {
              border: 0 none;
            }
          }
          .point-icon {
            position: absolute;
            right: -10px;
            top: 22px;
            transform: translateY(-50%);
            font-size: 18px;
            color: #cbe6fb;
          }
        }
        .video,
        .image {
          text-align: right;
          img {
            width: 114px;
            border-radius: 8px;
          }
        }
        // 文件
        .files {
          background: #ffffff;
          border-radius: 8px;
          width: 251px;
          padding: 12px 26px 15px 10px;
          display: flex;
          justify-content: space-between;
          .left-box {
            .file-name {
              font-size: 16px;
              font-weight: 500;
              color: #454545;
              line-height: 22px;
              margin-bottom: 8px;
              width: 150px;
              text-overflow: ellipsis;
              overflow: hidden;
              max-height: 44px;
              -webkit-line-clamp: 2;
              display: -webkit-box;
              -webkit-box-orient: vertical;
            }
            .file-size {
              font-size: 12px;
              font-weight: 400;
              color: #999999;
            }
          }
          .right-box {
            transform: translateY(3px);
            width: 46px;
            height: 46px;
            background: #95a2b4;
            display: flex;
            align-items: center;
            justify-content: center;
            img {
              width: 22px;
              height: 26px;
            }
          }
        }
        // 链接
        .link {
          background: #ffffff;
          border-radius: 8px;
          width: 251px;
          padding: 10px 12px;
          .link-name {
            font-size: 16px;
            font-weight: 500;
            color: #454545;
            line-height: 22px;
            text-overflow: ellipsis;
            overflow: hidden;
            max-height: 44px;
            -webkit-line-clamp: 2;
            display: -webkit-box;
            -webkit-box-orient: vertical;
          }
          .desc-box {
            display: flex;
            justify-content: space-between;
            // align-items: center;
            padding-top: 5px;
            .desc-text {
              font-size: 12px;
              font-weight: 400;
              color: #999999;
              line-height: 17px;
              width: 165px;
            }
            .link-pic {
              width: 46px;
              height: 46px;
              img {
                width: 100%;
                height: 100%;
                object-fit: cover;
              }
              // object-fit: cover;
              // margin-left: 20px;
            }
          }
        }
        // 小程序
        .mini-program {
          width: 251px;
          background-color: #ffffff;
          border-radius: 8px;
          .title {
            padding: 15px;
            font-size: 12px;
            font-weight: 400;
            color: #999999;
            line-height: 17px;
            display: flex;
            align-items: center;
            .text-div {
              text-overflow: ellipsis;
              white-space: nowrap;
              overflow: hidden;
            }
            img {
              margin-right: 5px;
              width: 15px;
              height: 15px;
            }
          }
          .name {
            height: 22px;
            font-size: 16px;
            font-weight: 500;
            color: #454545;
            line-height: 22px;
            padding: 0 15px;
            text-overflow: ellipsis;
            white-space: nowrap;
            overflow: hidden;
            width: 100%;
          }
          .code-pic {
            margin: 8px auto 23px;
            width: 157px;
            height: 157px;
            img {
              width: 100%;
              height: 100%;
            }
          }
          .bottom {
            font-size: 12px;
            font-weight: 400;
            color: #999999;
            line-height: 32px;
            padding: 0 10px;
            border-top: 1px solid #eeeeee;
          }
        }
      }
      .delete-icon {
        width: 28px;
        height: 28px;
      }
      .avatar {
        width: 45px;
        height: 45px;
        position: absolute;
        top: 0;
        right: 0;
        border-radius: 4px;
        object-fit: cover;
      }
    }
  }
  .btn-box {
    position: fixed;
    bottom: 0;
    left: 0;
    width: 100%;
    padding: 8px 15px;
    display: flex;
    background-color: #ffffff;
    box-shadow: 0px 1px 6px 0px rgba(0, 0, 0, 0.15);
    .back-btn {
      width: 90px;
      height: 34px;
      border: 1px solid #4a4a4a;
      text-align: center;
      font-size: 14px;
      font-weight: 500;
      color: #4a4a4a;
      line-height: 34px;
    }
    .submit-btn {
      width: 240px;
      height: 34px;
      background: #4a4a4a;
      font-size: 14px;
      font-weight: 500;
      color: #ffffff;
      line-height: 34px;
      text-align: center;
      margin-left: 15px;
    }
  }
}
</style>
