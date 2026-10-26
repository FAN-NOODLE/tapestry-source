<template>
  <!-- 图片、视频、小程序选择区域-->
  <div class="material-select-box">
    <!-- 素材综合选择展示区域-->
    <el-popover
      placement="top-start"
      trigger="hover"
      v-if="!item.hasMaterial"
    >
      <div class="select-material-list">
        <div
          class="item-select"
          v-for="(materialItem, materialIndex) in materialSelectList"
          @mouseover="selectStyle(materialItem)"
          @mouseout="outStyle(materialItem)"
          @click="openSelectMaterial(materialItem.type, index)"
          :key="materialIndex"
        >
          <div>
            <img
              v-show="materialItem.active"
              class="item-icon"
              :src="materialItem.icon"
              alt=""
            />
            <img
              v-show="!materialItem.active"
              class="item-icon"
              :src="materialItem.hoverIcon"
              alt=""
            />
            <div class="item-tip">
              {{ materialItem.tip }}
            </div>
          </div>
        </div>
      </div>
      <!-- 添加小程序/视频/图片按钮 -->
      <el-button
        slot="reference"
        v-if="true"
        class="select-material-btn"
        size="small"
      >
        <img
          class="material-btn-img"
          src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/add_icon.png"
          alt=""
        />
        <span>{{ planInfo.task_category === 0 ? '图片/视频/小程序/链接' : '图片/视频/链接' }}</span>
      </el-button>
    </el-popover>
    <!-- 右侧图片选择展示区域-->
    <div
      class="select-picture-box"
      v-if="item.currentMaterialType == 'picture'"
    >
      <draggable
        v-model="item.content.picture"
        handle=".move"
        class="draggable-box"
      >
        <div
          class="item-picture move"
          v-for="(pItem, pindex) in item.content.picture"
          :key="pindex"
        >
          <img
            class="item-img"
            :src="pItem.url"
            alt=""
          />
          <div
            class="delete-img-box"
            @click="deleteMaterial('picture', index, pindex)"
          >
            <img
              src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-delete-icon.png"
              alt=""
            />
          </div>
        </div>
        <div
          @click="openSelectMaterial('picture', index)"
          class="upload-card"
          v-if="!ep.isCoach && item.content.picture && item.content.picture.length < 9
            "
        >
          <div class="center">
            <i
              class="el-icon-plus avatar-uploader-icon"
              aria-hidden="true"
            ></i>
            <div
              class="tips"
              v-text="'上传图片'"
            ></div>
          </div>
        </div>

        <el-popover
          placement="top"
          width="140"
          trigger="hover"
          v-if="ep.isCoach && item.content.picture && item.content.picture.length < 9"
        >
          <div class="select-material-list select-material-list-2">
            <div
              class="item-select"
              v-for="(item, index) in materialSelectList2"
              @mouseover="selectStyle(item)"
              @mouseout="outStyle(item)"
              @click="openSelectMaterial(item.type, index)"
              :key="index"
            >
              <img
                v-show="item.active"
                class="item-icon"
                :src="item.icon"
                alt=""
              />
              <img
                v-show="!item.active"
                class="item-icon"
                :src="item.hoverIcon"
                alt=""
              />
              <div class="item-tip">{{ item.tip }}</div>
            </div>
          </div>
          <div
            class="upload-card"
            slot="reference"
          >
            <div class="center">
              <i
                class="el-icon-plus avatar-uploader-icon"
                aria-hidden="true"
              ></i>
              <div
                class="tips"
                v-text="'上传图片'"
              ></div>
            </div>
          </div>
        </el-popover>
      </draggable>
    </div>
    <!-- 右侧视频选择展示区域-->
    <div
      class="select-video-box"
      v-if="item.hasMaterial &&
        planInfo.termArr[index].content.video &&
        planInfo.termArr[index].content.video.url
        "
    >
      <img
        @click="
          openVideo(planInfo.termArr[index].content.video.url)
          "
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
        :poster="planInfo.termArr[index].content.video.video_image_url || planInfo.termArr[index].content.video.cover_url"
        :src="planInfo.termArr[index].content.video.url"
      ></video>
      <div
        class="delete-video-box"
        @click="deleteMaterial('video', index)"
      >
        <img
          src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-delete-icon.png"
          alt=""
        />
      </div>
    </div>
    <!-- 右侧小程序卡片选择展示区域-->
    <div
      class="select-applet-box"
      v-if="item.currentMaterialType == 'applet'"
    >
      <div class="select-applet-name">
        <img
          class="mini-img"
          src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/mini-applet-icon.png"
          alt=""
        /><span>{{ item.content.applet.app_name }}</span>
      </div>
      <div class="select-applet-title">
        {{ item.content.applet.title }}
      </div>
      <img
        class="select-applet-img"
        :src="item.content.applet.cover_url"
        alt=""
      />
      <div
        class="delete-applet-box"
        @click="deleteMaterial('applet', index)"
      >
        <img
          src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-delete-icon.png"
          alt=""
        />
      </div>
    </div>
    <!-- 右侧链接展示区域 -->
    <div
      class="select-link-box"
      v-if="item.currentMaterialType == 'link'"
    >
      <div class="select-link-title">
        {{ item.content.link.title }}
      </div>
      <div class="select-link-bottom">
        <span class="select-link-description">
          {{ item.content.link.description }}
        </span>
        <img
          class="select-applet-img"
          :src="item.content.link.cover_url"
          alt=""
        />
      </div>
      <div
        class="delete-link-box"
        @click="deleteMaterial('link', index)"
      >
        <img
          src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-delete-icon.png"
          alt=""
        />
      </div>
    </div>
    <!-- 贺卡 -->
    <div
      class="select-applet-box"
      v-if="item.hasMaterial && item.currentMaterialType === 'card'
        "
    >
      <div class="select-applet-name">
        <img
          class="mini-img"
          src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/mini-applet-icon.png"
          alt=""
        /><span>贺卡</span>
      </div>
      <div class="select-applet-title">
        {{ item.content.card.template_name }}
      </div>
      <img
        v-if="item.content.card.template_picture_url"
        class="select-applet-img"
        :src="item.content.card.template_picture_url[0]"
        alt=""
      />
      <div
        class="delete-applet-box"
        @click="deleteMaterial('card', index)"
      >
        <img
          src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-delete-icon.png"
          alt=""
        />
      </div>
    </div>
    <!-- 预览视频弹框 -->
    <open-video
      v-if="isOpenVideo"
      :video_url="videoUrl"
      @closeModel="cancelVideo"
    ></open-video>
    <!-- 文本选择框 -->
    <select-text
      v-if="textDialog"
      @textData="textData"
    />
    <!-- 选择贺卡弹框 -->
    <select-card
      v-if="showMaterialModal == 'card'"
      @submit="getSelectCard"
    ></select-card>
    <!-- 选择图片素材弹窗-->
    <select-picture-material
      v-if="showMaterialModal == 'picture'"
      :isFriend="planInfo.task_category === 1 ? true : false"
      :groupSendType="planInfo.group_send_type"
      @submit="getSelectMaterial"
    ></select-picture-material>
    <!-- 选择视频素材弹窗-->
    <select-video-material
      v-if="showMaterialModal == 'video'"
      :isFriend="planInfo.task_category === 1 ? true : false"
      :groupSendType="planInfo.group_send_type"
      @submit="getSelectMaterial"
    ></select-video-material>
    <!-- 选择小程序素材弹窗-->
    <select-wxapp-material
      v-if="showMaterialModal == 'applet'"
      :isFriend="planInfo.task_category === 1 ? true : false"
      :groupSendType="planInfo.group_send_type"
      @submit="getSelectMaterial"
    ></select-wxapp-material>
    <!-- 选择链接素材弹窗-->
    <select-link-material
      v-if="showMaterialModal == 'link'"
      :isFriend="planInfo.task_category === 1 ? true : false"
      :groupSendType="planInfo.group_send_type"
      @submit="getSelectMaterial"
    ></select-link-material>
  </div>
</template>

<script lang="ts">
import { Vue, Component, Prop, Watch } from "vue-property-decorator";

// 视频预览
import OpenVideo from "@/views/myapp/wisdom/marketing-calendar/components/com-open-video.vue";
// 添加文本
import selectText from "../../../../enterprise/material-library/components/text-material-dialog.vue";
// 选择图片素材弹窗
import selectPictureMaterial from "../../../../enterprise/material-library/components/pic-mterial-dialog.vue";
// 选择视频素材弹窗
import selectVideoMaterial from "../../../../enterprise/material-library/components/video-material-dialog.vue";
// 选择小程序
import selectWxappMaterial from "../../../../enterprise/material-library/components/small-mterial-dialog.vue";
// 选择链接
import selectLinkMaterial from "../../../../enterprise/material-library/components/link-material-dialog.vue";
// 选择贺卡
import selectCard from "./select-card.vue";
import { TASK_MATERIAL_SELECT_LIST_2, TASK_MATERIAL_SELECT_LIST_1 } from "@/contanst";
import draggable from "vuedraggable";
// 上传素材按钮信息字段类型
interface IMaterialItem {
  type: string;
  hoverIcon: string | any;
  icon: string | any;
  tip: string;
  active?: boolean;
}
@Component({
  components: {
    OpenVideo,
    selectText,
    selectPictureMaterial,
    selectVideoMaterial,
    selectWxappMaterial,
    selectLinkMaterial,
    selectCard,
    draggable,
  },
})
export default class extends Vue {
  // 任务信息
  @Prop({ default: () => ({}) }) planInfo!: any;
  @Watch("planInfo", { immediate: true, deep: true })
  onPlanInfoChange(val: any) {
    const selectList = this.ep.isCoach ? TASK_MATERIAL_SELECT_LIST_1 : TASK_MATERIAL_SELECT_LIST_1.slice(0, 4)
    if (val.task_category === 1) {
      this.materialSelectList = selectList.filter((item: any) => item.type !== 'applet')
    } else {
      this.materialSelectList = selectList
    }
  }
  // 当前选中的任务
  @Prop({ default: () => ({}) }) item!: any;
  // 当前选中任务下标
  @Prop({ default: () => 0 }) index!: number;
  // 文本弹窗
  @Prop({ default: false }) textDialog!: boolean;
  // 控制当前选择素材的弹框显隐
  @Prop({ default: "" }) showMaterialModal!: string;
  // 上传素材按钮信息
  materialSelectList: any = []
  //是否打开视频
  isOpenVideo = false;
  // 当前播放视频的url
  videoUrl = "";
  materialSelectList2: Array<any> = TASK_MATERIAL_SELECT_LIST_2;

  get ep() {
    return this.$store.getters['enterprise']
  }
  selectStyle(item: any) {
    const _this = this;
    this.$nextTick(function () {
      this.materialSelectList.forEach(function (item: any) {
        _this.$set(item, "active", false);
      });
      this.$set(item, "active", true);
    });
  }
  // 鼠标划出事件
  outStyle(item: any) {
    this.$set(item, "active", false);
  }
  // 添加对应素材
  openSelectMaterial(type: string, index: number) {
    const _this = this;
    this.materialSelectList.forEach(function (item: any) {
      _this.$set(item, "active", false);
    });
    this.$emit("openSelectMaterial", type, index);
  }
  // 点击素材的叉叉按钮
  deleteMaterial(type: string, index: number, fIndex?: number) {
    this.$emit("deleteMaterial", type, index, fIndex);
  }
  // 选择素材弹窗后回调
  getSelectMaterial(material?: any) {
    this.$emit("getSelectMaterial", material);
  }
  // 文本弹窗回调
  textData(data: any) {
    this.$emit("textData", data);
  }
  // 选择贺卡回调
  getSelectCard(item: any) {
    this.$emit("getSelectCard", item);
  }
  // 预览视频
  openVideo(url: string) {
    this.videoUrl = url;
    this.isOpenVideo = true;
  }
  // 关闭预览视频
  cancelVideo() {
    this.isOpenVideo = false;
  }
}
</script>
<style lang="scss" scoped>
.material-select-box {
  padding: 15px;
  margin-top: 12px;

  /deep/.el-input {
    .el-input__inner {
      border: 1px solid #DCDFE6;
    }
  }

  .select-material-btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 0;
    height: 28px;

    /deep/span {
      display: inline-flex;
      align-items: center;
    }

    .material-btn-img {
      width: 14px;
      height: 14px;
      margin-right: 5px;
    }
  }

  .select-applet-poster-box {
    cursor: pointer;
    width: 100px;
    height: 100px;
    // border: solid 1px #dcdfe6;

    position: relative;

    .select-applet-poster {
      width: 100%;
      height: 100%;
      border-radius: 4px;
    }
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
    z-index: 1;

    img {
      width: 14px;
      height: 14px;
      display: block;
    }
  }
}

.tips {
  color: #999;
  font-size: 12px;
}


.select-material-list {
  // width: 140px;
  height: 70px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;

  &.select-material-list-2 {
    width: 140px; 
  }

  .item-select {
    text-align: center;
    margin-right: 20px;
    cursor: pointer;

    .item-icon {
      width: 38px;
      height: 38px;
      display: block;
      margin: 0 auto;
    }

    .item-tip {
      font-size: 12px;
      color: #606266;
      line-height: 17px;
      margin-top: 7px;
    }
  }

  .item-select:last-child {
    margin-right: 0;
  }

  .item-select:hover {
    .item-tip {
      color: $color-main;
    }
  }
}

.select-picture-box {
  display: flex;
  flex-wrap: wrap;
  margin-top: 12px;

  .draggable-box {
    display: flex;
    flex-wrap: wrap;
  }

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
  height: 83px;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  position: relative;
  padding: 10px;

  .select-link-title {
    width: 220px;
    text-overflow: ellipsis;
    overflow: hidden;
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
</style>
