<template>
  <div>
    <div class="publish-content-box">
      <!-- 表情、文本选择区域-->
      <div class="select-box-insert">
        <add-emoji
          @select="addEmoji"
          class="add-emoji"
        >
          <template slot="reference">
            <el-button
              class="insert-btn"
              size="small"
            ><img
                class="btn-img"
                src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/face_icon.png"
                alt=""
              /><span>插入表情</span></el-button>
          </template>
        </add-emoji>
        <el-button
          class="insert-btn"
          size="small"
          @click="openSelectText"
        ><img
            class="btn-img"
            src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/text_icon.png"
            alt=""
          /><span>添加文本</span></el-button>
      </div>
      <!--多行文本区域-->
      <div>
        <el-input
          ref="textarea"
          class="textarea-box"
          placeholder="请输入营销话术内容"
          v-model="value.text.content"
          :rows="4"
          type="textarea"
          maxlength="500"
          show-word-limit
        >
        </el-input>
      </div>
      <!-- 图片、视频、链接、小程序选择区域-->
      <div class="material-select-box">
        <!-- 右侧素材综合选择展示区域-->
        <el-popover
          placement="top-start"
          :width="materialSelectList.length*60"
          trigger="hover"
          v-if="!currentMaterialType"
        >
          <div class="select-material-list">
            <div
              class="item-select"
              v-for="(item, index) in materialSelectList"
              @mouseover="selectStyle(item)"
              @mouseout="outStyle(item)"
              @click="openSelectMaterial(item.type)"
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
          <el-button
            slot="reference"
            v-if="true"
            class="select-material-btn"
            size="small"
          ><img
              class="material-btn-img"
              src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/add_icon.png"
              alt=""
            /><span>图片/链接/小程序/视频</span></el-button>
        </el-popover>
        <!-- 右侧图片选择展示区域-->
        <div
          class="select-picture-box"
          v-if="currentMaterialType == 'picture'"
        >
          <draggable
            v-model="value.picture"
            handle=".move"
            class="draggable-box"
          >
            <div
              class="item-picture move"
              v-for="(item, index) in value.picture"
              :key="index"
            >
              <img
                class="item-img"
                :src="item.url"
                alt=""
              />
              <div
                class="delete-img-box"
                @click="deleteMaterial('picture', index)"
              >
                <img
                  src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-delete-icon.png"
                  alt=""
                />
              </div>
            </div>
    
            <div
              @click="openSelectMaterial('picture')"
              class="upload-card"
              v-if="!ep.isCoach && value.picture.length < 9
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
              v-if="ep.isCoach && value.picture.length < 9"
            >
              <div class="select-material-list select-material-list-2">
                <div
                  class="item-select"
                  v-for="(item, index) in materialSelectList2"
                  @mouseover="selectStyle(item)"
                  @mouseout="outStyle(item)"
                  @click="openSelectMaterial(item.type)"
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
          v-if="currentMaterialType == 'video' &&
            value.video &&
            value.video.url
            "
        >
          <img
            @click="openVideo(value.video.url)"
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
            :src="value.video.url"
          ></video>
          <div
            class="delete-video-box"
            @click="deleteMaterial('video')"
          >
            <img
              src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-delete-icon.png"
              alt=""
            />
          </div>
        </div>
        <!-- 右侧链接选择展示区域-->
        <div
          class="select-link-box"
          v-if="currentMaterialType == 'link' &&
            value.link &&
            value.link.link_url
            "
        >
          <div class="select-link-title">
            {{ value.link.title }}
          </div>
          <div class="select-link-bottom">
            <!-- 为了解决flex布局下超出省略失效问题，所以嵌套-->
            <div>
              <span class="select-link-description">{{
                value.link.description
              }}</span>
            </div>
            <img
              class="select-link-cover"
              :src="value.link.cover_url"
              alt=""
            />
          </div>
          <div
            class="delete-link-box"
            @click="deleteMaterial('link')"
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
          v-if="currentMaterialType == 'applet' &&
            value.applet &&
            value.applet.app_name
            "
        >
          <div class="select-applet-name">
            <img
              class="mini-img"
              src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/mini-applet-icon.png"
              alt=""
            /><span>{{ value.applet.app_name }}</span>
          </div>
          <div class="select-applet-title">
            {{ value.applet.title }}
          </div>
          <img
            class="select-applet-img"
            :src="value.applet.cover_url"
            alt=""
          />
          <div
            class="delete-applet-box"
            @click="deleteMaterial('applet')"
          >
            <img
              src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-delete-icon.png"
              alt=""
            />
          </div>
        </div>
      </div>
      <p v-if="ep.isCoach && task_category === 0" class="tips">
        【本地上传】可支持上传图片及视频素材，最大不超过10M
      </p>
      <p v-if="ep.isCoach && task_category === 1" class="tips">
        <ul>
          <li>最多支持9个图片 类型，或者1个视频，或者1个链接。类型只能三选一</li>
          <li>
            选择【图片】素材：像素不超过1555200。图片大小不超过10M。最多支持传入9个；（普通图片尺寸建议不超过1440*1080px，长图片，长边建议不超过10800px）
          </li>
          <li>选择【视频】素材：最长不超过30S，最大不超过10MB。</li>
          <li>【本地上传】可支持上传图片及视频素材，最大不超过10M</li>
        </ul>
      </p>
    </div>

    <!-- 选择文本素材弹窗-->
    <select-text-material
      v-if="showSelectText"
      :isFriend="task_category === 1 ? true : false"
      :groupSendType="1"
      @textData="getTextData"
    ></select-text-material>
    <!-- 选择图片素材弹窗-->
    <select-picture-material
      v-if="showMaterialModal == 'picture'"
      :isFriend="task_category === 1 ? true : false"
      :groupSendType="1"
      @submit="getSelectMaterial"
    ></select-picture-material>
    <!-- 选择小程序素材弹窗-->
    <select-applet-material
      v-if="showMaterialModal == 'applet'"
      :isFriend="task_category === 1 ? true : false"
      :groupSendType="1"
      @submit="getSelectMaterial"
    ></select-applet-material>
    <!-- 选择链接素材弹窗-->
    <select-link-material
      v-if="showMaterialModal == 'link'"
      :isFriend="task_category === 1 ? true : false"
      :groupSendType="1"
      @submit="getSelectMaterial"
    ></select-link-material>
    <!-- 选择视频素材弹窗-->
    <select-video-material
      v-if="showMaterialModal == 'video'"
      :isFriend="task_category === 1 ? true : false"
      :groupSendType="1"
      @submit="getSelectMaterial"
    ></select-video-material>
    <!-- 预览视频弹框 -->
    <open-video
      v-if="isOpenVideo"
      :video_url="videoUrl"
      @closeModel="cancelVideo"
    ></open-video>
  </div>
</template>

<script lang="ts">
import { Component, Model, Prop, Vue, Watch } from "vue-property-decorator";
import UploadImage from "@/components/business/upload-image/index.vue";
import Api from "@/api/wisdom/temp";
// 选择表情弹窗
import AddEmoji from "@/views/myapp/enterprise/components/add-emoji.vue";
// 选择员工弹窗
import SelectMember from "@/components/business/select-member/index.vue";
// 选择文本素材弹窗
import selectTextMaterial from "@/views/myapp/enterprise/material-library/components/text-material-dialog.vue";
// 选择图片素材弹窗
import selectPictureMaterial from "@/views/myapp/enterprise/material-library/components/pic-mterial-dialog.vue";
// 选择小程序素材弹窗
import selectAppletMaterial from "@/views/myapp/enterprise/material-library/components/small-mterial-dialog.vue";
// 选择链接素材弹窗
import selectLinkMaterial from "@/views/myapp/enterprise/material-library/components/link-material-dialog.vue";
// 选择视频素材弹窗
import selectVideoMaterial from "@/views/myapp/enterprise/material-library/components/video-material-dialog.vue";
import openVideo from "@/views/myapp/enterprise/material-library/components/com-open-video.vue";
import PerModule from "@/store/modules/permission";
import LocalUpload from '@/mixins/LocalUpload';
import { TASK_MATERIAL_SELECT_LIST_1, TASK_MATERIAL_SELECT_LIST_2 } from "@/contanst";
import draggable from "vuedraggable";

function  getMaterialType(template_content: any) {
  if (template_content.picture && template_content.picture.length > 0) {
    return "picture";
  } else if (template_content.link && template_content.link.link_url) {
    return "link";
  } else if (template_content.video && template_content.video.url) {
    return "video";
  } else if (template_content.applet && template_content.applet.app_name) {
    return "applet";
  } else {
    return "text";
  }
}

@Component(
  {
    components: {
      UploadImage,
      AddEmoji,
      SelectMember,
      selectTextMaterial,
      selectPictureMaterial,
      selectAppletMaterial,
      selectLinkMaterial,
      selectVideoMaterial,
      openVideo,
      draggable
    },
    mixins: [LocalUpload],
  }
)
export default class TaskContent extends Vue { 
  @Model('change', { default: ()=>({
    text: {
      content: ""
    },
    picture: [],
    video: {},
    link: {},
    applet: {}
  }), type: Object }) value: any;
  @Prop({ default: 0 }) task_category: number
  materialSelectList: Array<any> = TASK_MATERIAL_SELECT_LIST_1;
  materialSelectList2: Array<any> = TASK_MATERIAL_SELECT_LIST_2;
  contentEmpty = false; // 判断发布内容是否为空
  showError = false; // 控制发布内容错误提示是否展示
  currentMaterialType = ""; // 当前选择素材的类型
  showMaterialModal = ""; // 控制当前选择素材的弹框显隐
  isOpenVideo = false; //是否打开视频
  videoUrl = ""; // 当前播放视频的url
  showSelectText = false;

  get ep() {
    return this.$store.getters['enterprise']
  }
  @Watch("value", {deep: true, immediate: true})
  onValueChange(val: any) {
    this.$emit('change', val);
    this.currentMaterialType = getMaterialType(val);
    if (this.currentMaterialType == "text") {
      this.currentMaterialType = "";
    }
  }

  @Watch("task_category", {immediate: true})
  onTaskCategoryChange(val: number) {
    if(val === 1) {
      this.materialSelectList = TASK_MATERIAL_SELECT_LIST_1.filter((item: any) => item.type !== "applet");
    } else {
      this.materialSelectList = TASK_MATERIAL_SELECT_LIST_1;
    }
  }
  /**************************************************data结束**************************************************/
  /***********************************************methods开始**************************************************/
  // 选择素材弹窗后回调
  getSelectMaterial(material?: any) {
    console.log("material", material);
    // 选择图片点击取消按钮回调处理
    if (
      this.currentMaterialType == "picture" &&
      (!material || material.length == 0) &&
      this.value.picture.length > 0
    ) {
      this.showMaterialModal = "";
      return;
    }
    // 若素材返回为空，则直接关闭弹窗
    if (
      !material ||
      Object.keys(material).length == 0 ||
      material.length == 0
    ) {
      this.currentMaterialType = "";
      this.showMaterialModal = "";
      return;
    }
    switch (this.showMaterialModal) {
      case "picture":
        if (material.length + this.value.picture.length > 9) {
          this.$message.warning("图片素材最多上传9张");
          return;
        }
        this.value.picture.push(...material);
        this.showMaterialModal = "";
        break;
      case "video":
      case "link":
      case "applet":
        this.value[this.showMaterialModal] = material;
        this.showMaterialModal = "";
        break;
    }
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
  // 打开选择文本弹窗
  openSelectText() {
    this.showSelectText = true;
  }
  // 鼠标划入事件
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
  // 打开对应素材弹窗
  openSelectMaterial(type: string) {
    if(type === 'local-image') {
      this.handleUpload('image');
    } else if(type === 'local-video') {
      this.handleUpload('video');
    } else {
      const _this = this;
      this.materialSelectList.forEach(function (item: any) {
        _this.$set(item, "active", false);
      });
      this.currentMaterialType = type;
      this.showMaterialModal = type;
    }
  }
  // 点击素材的叉叉按钮
  deleteMaterial(type: string, index?: number) {
    switch (type) {
      case "picture":
        this.value.picture.splice(index, 1);
        if (this.value.picture.length <= 0)
          this.currentMaterialType = "";
        break;
      case "video":
      case "link":
      case "applet":
        this.value[type] = {};
        this.currentMaterialType = "";
        break;
    }
  }
  // 添加表情/添加文本
  async addEmoji(emj: any) {
    const myField = (this.$refs.textarea as any).$refs.textarea;
    if (myField.selectionStart || myField.selectionStart === 0) {
      const startPos = myField.selectionStart;
      const endPos = myField.selectionEnd;
      // 添加不能超过500个字符
      if (this.value.text.content.length + emj.length > 500) {
        this.$message.warning("文本内容不能超过500个字符");
        return;
      }
      this.value.text.content =
        myField.value.substring(0, startPos) +
        emj +
        myField.value.substring(endPos, myField.value.length);
      await this.$nextTick();
      myField.focus();
      myField.setSelectionRange(endPos + emj.length, endPos + emj.length);
    } else {
      this.value.text.content += emj;
    }
  }
  async handleUpload(type: string) {
    try {
      // 上传多个媒体文件
      const results = await (this as any).uploadMedias({
        multiple: type === 'image',
        maxCount: type === 'image' ? 9-this.value.picture.length : 1,
        maxSize: 10 * 1024 * 1024, // 10MB
        acceptTypes: type === 'image' ? ['image/jpeg','image/png'] : ['video/mp4', 'video/m4v'],
        validate: this.task_category !== 0,
        compress: this.task_category !== 0,
      });
      console.log('多文件上传结果:', results);
      if(results.length === 0) {
        return;
      }
      
      if(type === 'image') {
        this.currentMaterialType = "picture";
        this.$set(this.value, 'picture', [
          ...this.value.picture, 
          ...results.map(item => ({
            url: item.url,
            type: 2,
            id: 0
          }))
        ]);
      } else {
        this.currentMaterialType = "video";
        this.$set(this.value, 'video', {
          url:  results[0].url,
          type: 5,
          id: 0,
          duration: results[0].duration,
          video_image_url: results[0].cover_url,
        });
      }
      console.log(this.value);
      
    } catch (error) {
      console.error('上传失败:', error);
      // this.$message.error(error.message);
    }
  }
  
  // 添加文本回调
  getTextData(data: any) {
    this.showSelectText = false;
    if (data) {
      this.addEmoji(data.content);
    }
  }
  // 点击取消按钮
  cancel() {
    this.$router.go(-1);
  }
}
</script>

<style lang="scss" scoped> 
.publish-content-box {
  width: 600px !important;
  border: solid 1px #dcdfe6;
  border-radius: 4px;
  min-height: 425px;

  .select-box-insert {
    height: 60px;
    margin: 0 15px;
    display: flex;
    align-items: center;

    .add-emoji {
      display: inline-flex;
      justify-content: center;
      margin-right: 10px;

      /deep/span {
        display: inline-flex;
        align-items: center;
      }
    }

    .insert-btn {
      border: none;
      border-radius: 25px;
      display: inline-flex;
      align-items: center;
      justify-content: center;
      padding: 0;
      width: 87px;
      height: 28px;
      background-color: #f5f6fa;

      /deep/span {
        display: inline-flex;
        align-items: center;
      }

      .btn-img {
        width: 14px;
        height: 14px;
        margin-right: 5px;
      }
    }
  }

  .textarea-box {
    width: calc(100% - 30px) !important;
    margin-left: 15px;
    border-top: solid 1px #ebeef5;
    border-bottom: solid 1px #ebeef5;

    /deep/textarea {
      border: none;
      resize: none;
      height: 255px !important;
      padding: 15px 0;
      font-family: PingFangSC-Medium, PingFang SC;
    }

    /deep/.el-input__count {
      right: 0;
    }
  }

  .material-select-box {
    padding: 15px;

    .select-material-btn {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      padding: 0;
      width: 185px;
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

    .select-video-box {
      position: relative;
      width: 80px;
      height: 80px;
      border-radius: 4px;

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

    .select-picture-box {
      display: flex;
      flex-wrap: wrap;
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

    .select-link-box {
      width: 225px;
      height: 83px;
      border: solid 1px #e4e7ed;
      border-radius: 4px;
      position: relative;
      padding: 10px;

      .select-link-title {
        font-size: 14px;
        color: #303133;
        line-height: 20px;
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

        .select-link-cover {
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
  }
}

div {
  box-sizing: border-box;
}

.select-material-list {
  // width: 350px;
  height: 70px;
  display: flex;
  align-items: center;
  justify-content: center;

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
      white-space: nowrap;
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

.tips {
  color: #999;
  font-size: 12px;
}

.friend-circle-edit {
  .content-box {
    display: flex;
    flex: 1;
    padding-left: 30px;

    /deep/.in-item {
      min-width: 500px;
      max-width: calc(100% - 20px);
      width: auto;
    }

    .date-select-box {
      display: flex;
      align-items: center;

      .time-picker {
        margin-left: 20px;
        flex: 1;
      }
    }

    .item-staff {
      display: inline-flex;
      align-items: center;
      height: 28px;
      background-color: #f0f3fe;
      border: solid 1px #d6dffd;
      padding: 0 10px;
      border-radius: 3px;
      margin-right: 10px;
      margin-bottom: 10px;
      min-width: 87px;
      justify-content: center;

      .staff-icon {
        width: 14px;
        height: 14px;
        margin-right: 5px;
      }
    }


    .publish-content-box {
      width: 100%;
      border: solid 1px #dcdfe6;
      border-radius: 4px;
      min-height: 425px;

      .select-box-insert {
        height: 60px;
        margin: 0 15px;
        display: flex;
        align-items: center;

        .add-emoji {
          display: inline-flex;
          justify-content: center;
          margin-right: 10px;

          /deep/span {
            display: inline-flex;
            align-items: center;
          }
        }

        .insert-btn {
          border: none;
          border-radius: 25px;
          display: inline-flex;
          align-items: center;
          justify-content: center;
          padding: 0;
          width: 87px;
          height: 28px;
          background-color: #f5f6fa;

          /deep/span {
            display: inline-flex;
            align-items: center;
          }

          .btn-img {
            width: 14px;
            height: 14px;
            margin-right: 5px;
          }
        }
      }

      .textarea-box {
        width: calc(100% - 30px) !important;
        margin-left: 15px;
        border-top: solid 1px #ebeef5;
        border-bottom: solid 1px #ebeef5;

        /deep/textarea {
          border: none;
          resize: none;
          height: 255px !important;
          padding: 15px 0;
          font-family: PingFangSC-Medium, PingFang SC;
        }

        /deep/.el-input__count {
          right: 0;
        }
      }

      .material-select-box {
        padding: 15px;

        .select-material-btn {
          display: inline-flex;
          align-items: center;
          justify-content: center;
          padding: 0;
          width: 185px;
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

        .select-video-box {
          position: relative;
          width: 80px;
          height: 80px;
          border-radius: 4px;

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

        .select-picture-box {
          display: flex;
          flex-wrap: wrap;

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

        .select-link-box {
          width: 225px;
          height: 83px;
          border: solid 1px #e4e7ed;
          border-radius: 4px;
          position: relative;
          padding: 10px;

          .select-link-title {
            font-size: 14px;
            color: #303133;
            line-height: 20px;
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

            .select-link-cover {
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
      }
    }

    .error-msg {
      padding-top: 0;
      margin-top: 5px;
      line-height: 18px;
      max-height: 36px;
      color: #f67170;
      font-size: 12px;
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      line-clamp: 2;
      -webkit-box-orient: vertical;
      position: absolute;
      left: 0px;
    }
  }
}



.message-info {
  width: 608px;
  font-size: 12px;
  color: red;
  line-height: 20px;
}


</style>