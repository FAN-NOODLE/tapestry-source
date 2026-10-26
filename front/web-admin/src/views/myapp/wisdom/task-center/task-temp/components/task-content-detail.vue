<template>
  <div class="material-container" v-if="content">
    <!-- 图片 -->
    <div class="select-picture-box" v-if="content.picture">
      <div
        class="item-picture"
        v-for="(pItem, pindex) in content.picture"
        :key="pindex"
      >
        <el-image class="item-img" :src="pItem.url" :preview-src-list="[pItem.url]" alt="" />
      </div>
    </div>
    <!-- 视频 -->
    <div
      class="select-video-box"
      v-if="content.video && content.video.url"
    >
      <img
        @click="openVideo(content.video.url)"
        class="video-play-icon"
        src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/play_icon.png"
        alt=""
      />
      <video
        width="80"
        height="80"
        x5-video-player-fullscreen="true"
        x5-playsinline
        playsinline
        webkit-playsinline
        muted
        preload="auto"
        :src="content.video.url"
      ></video>
    </div>
    <!-- 小程序 -->
    <div
      class="select-applet-box"
      v-if="
        content.applet &&
          Object.keys(content.applet).length > 0
      "
    >
      <div class="select-applet-name">
        <img
          class="mini-img"
          src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/mini-applet-icon.png"
          alt=""
        /><span>{{ content.applet.app_name }}</span>
      </div>
      <div class="select-applet-title">
        {{ content.applet.title }}
      </div>
      <img
        class="select-applet-img"
        :src="content.applet.cover_url"
        alt=""
      />
    </div>
    <div
      class="select-link-box"
      v-if="
        content.link &&
          Object.keys(content.link).length > 0
      "
    >
      <div class="select-link-title">
        {{ content.link.title }}
      </div>
      <div class="select-link-bottom">
        <span class="select-link-description">
          {{ content.link.description }}
        </span>
        <img
          class="select-applet-img"
          :src="content.link.cover_url"
          alt=""
        />
      </div>
    </div>
  </div>
</template>
<script lang="ts">
  import {Vue, Component, Prop} from "vue-property-decorator";
  @Component({})
  export default class TaskContentDetail extends Vue {
    @Prop({default: () => {}}) content:any;
    openVideo(url: string) {
      this.$emit("openVideo", url);
    }
  }
</script>
<style lang="scss" scoped>
div {
  box-sizing: border-box;
}
.content {
  .font-re {
    display: flex;
    line-height: 38px;
    .title {
      min-width: 92px;
    }
    .content {
      background: none;
      margin-left: 10px;
      box-shadow: none;
      padding: 0;
    }
  }
  .material-container {
    margin: 20px 0 0 0;
    div {
      margin-bottom: 10px;
    }
    .select-picture-box {
      display: flex;
      flex-wrap: wrap;
      margin-top: 12px;
      .item-picture {
        position: relative;
        width: 80px;
        height: 80px;
        margin-right: 10px;
        margin-bottom: 10px;
        .item-img {
          width: 100%;
          height: 100%;
          border-radius: 4px;
        }
        .delete-img-box {
          width: 20px;
          height: 20px;
          display: flex;
          align-items: center;
          justify-content: center;
          position: absolute;
          cursor: pointer;
          right: -10px;
          top: -10px;
          z-index: 2;
          img {
            width: 14px;
            height: 14px;
            display: block;
          }
        }
      }
      .upload-card {
        background-color: $color-white;
        border: 1px solid $color-second-line;
        border-radius: 4px;
        position: relative;
        width: 80px;
        height: 80px;
        &:hover {
          border: 1px solid $color-main;
        }
        .center {
          position: absolute;
          bottom: 0;
          top: 0;
          left: 0;
          right: 0;
          margin: auto;
          text-align: center;
          height: 33px;
          .avatar-uploader-icon {
            font-size: 14px;
            display: block;
            margin-bottom: 5px;
          }
          .tips {
            font-size: 12px;
            color: $color-text-second;
            height: 12px;
            line-height: 12px;
            overflow: hidden;
          }
        }
      }
    }
    .select-video-box {
      position: relative;
      width: 80px;
      height: 80px;
      border-radius: 4px;
      margin-top: 12px;
      video {
        border-radius: 4px;
      }
      .video-play-icon {
        position: absolute;
        left: 0;
        right: 0;
        bottom: 0;
        top: 0;
        margin: auto;
        width: 25px;
        height: 25px;
        z-index: 99;
        cursor: pointer;
      }
      .delete-video-box {
        width: 20px;
        height: 20px;
        display: flex;
        align-items: center;
        justify-content: center;
        position: absolute;
        cursor: pointer;
        right: -10px;
        top: -10px;
        z-index: 2;
        img {
          width: 14px;
          height: 14px;
          display: block;
        }
      }
    }
    .select-applet-box {
      width: 165px;
      height: 211px;
      border: solid 1px #dcdfe6;
      border-radius: 4px;
      position: relative;
      padding: 10px;
      .select-applet-name {
        display: flex;
        font-size: 12px;
        color: #707277;
        line-height: 20px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        margin-bottom: 4px;
        .mini-img {
          width: 20px;
          height: 20px;
          margin-right: 5px;
        }
      }
      .select-applet-title {
        font-size: 14px;
        color: #303133;
        line-height: 20px;
        overflow: hidden;
        text-overflow: ellipsis;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
      }
      .select-applet-img {
        width: 145px;
        height: 116px;
        margin-top: 10px;
        border-radius: 4px;
      }
      .delete-applet-box {
        width: 20px;
        height: 20px;
        display: flex;
        align-items: center;
        justify-content: center;
        position: absolute;
        cursor: pointer;
        right: -10px;
        top: -10px;
        z-index: 2;
        img {
          width: 14px;
          height: 14px;
          display: block;
        }
      }
    }
    .select-link-box {
      width: 225px;
      // height: 83px;
      border: 1px solid #e4e7ed;
      border-radius: 4px;
      position: relative;
      padding: 10px;
      line-height: normal;
      .select-link-title {
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
      .select-link-bottom {
        display: flex;
        margin-top: 3px;
        .select-link-description {
          width: 155px;
          font-size: 12px;
          color: #909399;
          line-height: 17px;
          overflow: hidden;
          text-overflow: ellipsis;
          display: -webkit-box;
          -webkit-line-clamp: 2;
          -webkit-box-orient: vertical;
        }
        img {
          width: 40px;
          height: 40px;
          display: block;
          border-radius: 4px;
          margin-left: 10px;
        }
      }
      .delete-link-box {
        width: 20px;
        height: 20px;
        display: flex;
        align-items: center;
        justify-content: center;
        position: absolute;
        cursor: pointer;
        right: -10px;
        top: -10px;
        z-index: 2;
        img {
          width: 14px;
          height: 14px;
          display: block;
        }
      }
    }
  }
}
</style>