<template>
  <div class="code-fun-flex-col page">
    <div
      class="code-fun-flex-col code-fun-justify-start code-fun-flex-auto group_4"
    >
      <div
        v-if="resultList.length > 0"
        class="code-fun-flex-col list section_3 space-y-15"
      >
        <div
          class="code-fun-flex-row space-x-16"
          :key="i"
          v-for="(item, i) in resultList"
        >
          <img
            v-if="type == 5"
            class="image_7"
            :src="item.url ? item.url : require('@/assets/video-bg.png')"
          />

          <img v-else class="image_7" :src="item.url" />
          <div
            class="code-fun-flex-col code-fun-flex-auto code-fun-self-start group_5 space-y-16 code-fun-overflow-hidden"
          >
            <div
              class="code-fun-flex-row code-fun-justify-between code-fun-items-center"
            >
              <span class="font_1 code-fun-overflow-hidden code-fun-shrink">{{
                item.name
              }}</span>
              <span
                class="font_2 code-fun-shrink-0"
                :class="[[0, 2].includes(item.is_deal) ? 'text_6' : 'text_3']"
                >{{ item.is_deal | statusText }}</span
              >
            </div>
            <span class="code-fun-self-start font_2 text_4">{{
              item.reason
            }}</span>
          </div>
        </div>
      </div>
    </div>
    <!-- 视频预览 -->
    <div class="video-content" v-if="showModel">
      <div class="mask" @click="closeModel"></div>
      <div class="content">
        <video
          :src="videoOptions.src"
          :controls="videoOptions.controls"
          class="video-js vjs-big-play-centered vjs-fluid"
          playsinline="true"
          x-webkit-airplay="allow"
          x5-playsinline
          webkit-playsinline
          object-fit="cover"
          style="width: 100%; max-height: 50vh"
          ref="videoRef"
        ></video>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

import { getResult } from "@/api/marketing-promotion/index";
import wx from "weixin-js-sdk";
@Component({
  filters: {
    statusText(val) {
      const is_deal = {
        "0": "上传成功", //不压缩
        "1": "压缩中",
        "2": "上传成功",
        "3": "上传失败",
      };
      return is_deal[val];
    },
  },
})
export default class MaterialResult extends Vue {
  /** 素材类型 2图片 5视频 */
  private type: string | number = "";
  /**压缩结果 */
  private resultList = [];
  private showModel = false;
  private videoOptions = {
    controls: true,
    src: "", // url地址
  };
  private player = null;
  // 关闭视频
  private closeModel() {
    this.showModel = false;
    this.player.currentTime = 0;
    this.player && this.player.pause();
    this.player = null;
  }
  // 打开视频
  private openModel(item) {
    if (![0, 2].includes(item.is_deal)) {
      return;
    }
    localStorage.setItem("cacheTime", "0");
    this.videoOptions.src = item.url;
    this.showModel = true;
    setTimeout(() => {
      this.player = this.$refs.videoRef;
      this.player.play();
    }, 500);
  }
  /** 查询压缩结果 */
  private getResult(compress_id, source_type) {
    const params = {
      compress_id,
      source_type,
    };
    getResult(params)
      .then((res: IApiResult) => {
        const { code, data } = res;
        if (code != 200) return;
        this.resultList = data || [];
      })
      .catch((e) => {
        console.error(e);
      });
  }
  private popstate() {
    console.log("关闭微信页面");
    wx.closeWindow(); // 微信网页退出
  }
  private async created() {
    // 当前需要直接退出的页面
    window.history.pushState("_android_back", "", location.href);
    window.addEventListener("popstate", this.popstate, false);
  }
  /** 生命周期 */
  private mounted() {
    let id = this.$route.query.id ? Number(this.$route.query.id) : "",
      type = this.$route.query.type ? Number(this.$route.query.type) : "";
    this.type = type;
    this.getResult(id, type);
  }
  private async destroyed() {
    // 记得清除，不然单页面，其它页面也可以用
    window.removeEventListener("popstate", this.popstate, false);
  }
}
</script>
<style lang="scss" scoped>
/************************************************************
** 请将全局样式拷贝到项目的全局 CSS 文件或者当前页面的顶部 **
** 否则页面将无法正常显示                                  **
************************************************************/

html {
  font-size: 16px;
}

body {
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", "Oxygen",
    "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue",
    "Microsoft Yahei", sans-serif;
}

body * {
  box-sizing: border-box;
  flex-shrink: 0;
}

#app {
  width: 100vw;
  height: 100vh;
}
.video-content {
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  right: 0;
  z-index: 99;
  .mask {
    position: absolute;
    width: 100%;
    height: 100%;
    background: rgba($color: #000000, $alpha: 0.6);
    z-index: 2;
  }
  .content {
    position: absolute;
    width: 100%;
    top: 0;
    max-height: 50vh;
    z-index: 8;
    background: white;
    font-size: 0;
  }
}
.code-fun-overflow-hidden {
  overflow: hidden;
}

.code-fun-flex-row {
  display: flex;
  flex-direction: row;
}

.code-fun-flex-col {
  display: flex;
  flex-direction: column;
}

.code-fun-justify-start {
  justify-content: flex-start;
}

.code-fun-justify-end {
  justify-content: flex-end;
}

.code-fun-justify-center {
  justify-content: center;
}

.code-fun-justify-between {
  justify-content: space-between;
}

.code-fun-justify-around {
  justify-content: space-around;
}

.code-fun-justify-evenly {
  justify-content: space-evenly;
}

.code-fun-items-start {
  align-items: flex-start;
}

.code-fun-items-end {
  align-items: flex-end;
}

.code-fun-items-center {
  align-items: center;
}

.code-fun-items-baseline {
  align-items: baseline;
}

.code-fun-items-stretch {
  align-items: stretch;
}

.code-fun-self-start {
  align-self: flex-start;
}

.code-fun-self-end {
  align-self: flex-end;
}

.code-fun-self-center {
  align-self: center;
}

.code-fun-self-baseline {
  align-self: baseline;
}

.code-fun-self-stretch {
  align-self: stretch;
}

.code-fun-flex-1 {
  flex: 1 1 0%;
}

.code-fun-flex-auto {
  flex: 1 1 auto;
}

.code-fun-grow {
  flex-grow: 1;
}

.code-fun-grow-0 {
  flex-grow: 0;
}

.code-fun-shrink {
  flex-shrink: 1;
}

.code-fun-shrink-0 {
  flex-shrink: 0;
}

.code-fun-relative {
  position: relative;
}

.page {
  background-color: #f4f7fb;
  width: 100%;
  min-height: 100vh;
  .group_4 {
    padding: 16px 0;
    .list {
      margin: 0 15px;
      .space-x-16 {
        & > *:not(:first-child) {
          margin-left: 16px;
        }
        .image_7 {
          width: 88px;
          height: 88px;
        }
        .group_5 {
          margin-top: 2px;
          .font_1 {
            font-size: 14px;
            font-family: PingFangSC;
            color: #000000;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
          }
          .font_2 {
            font-size: 12px;
            font-family: PingFangSC;
            line-height: 11px;
          }
          .text_3 {
            color: #d84242;
          }
          .text_4 {
            color: #777777;
            line-height: 11.5px;
          }
          .text_6 {
            display: inline-block;
            height: 100%;
            color: #3bc255;
          }
        }
        .space-y-16 {
          & > *:not(:first-child) {
            margin-top: 16px;
          }
        }
      }
    }
    .section_3 {
      padding: 16px 14px;
      background-color: #ffffff;
    }
    .space-y-15 {
      & > *:not(:first-child) {
        margin-top: 15px;
      }
    }
  }
}
</style>
