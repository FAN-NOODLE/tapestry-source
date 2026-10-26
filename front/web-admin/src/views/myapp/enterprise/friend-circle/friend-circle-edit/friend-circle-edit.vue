<template>
  <div class="friend-circle-edit">
    <layout-form>
      <template slot="content">
        <div class="content-box">
          <div class="content-left">
            <div class="phone-box">
              <div
                class="friend-top"
                :style="
                  'background-image: url(' +
                    (appearanceInfo.background
                      ? appearanceInfo.background
                      : 'https://wesociastg.blob.core.chinacloudapi.cn/wesocial-uat/uploads/2021/3/12/f0ea071d94dd85390d45c7d249f72951.jpeg') +
                    ')'
                "
              >
                <img
                  class="mobile-phone-top"
                  src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/mobile-header.png"
                  alt=""
                />
                <img
                  class="head-portrait-box"
                  :src="
                    appearanceInfo.avatar
                      ? appearanceInfo.avatar
                      : 'https://wesociastg.blob.core.chinacloudapi.cn/wesocial-uat/uploads/2021/3/12/5f136e2f41e2a7e63540a03b1f4d6f56.png'
                  "
                  alt=""
                />
                <div class="staff-name">
                  {{ appearanceInfo.name || "员工姓名" }}
                </div>
                <div class="signature">
                  {{ appearanceInfo.sign || "朋友圈签名" }}
                </div>
              </div>
              <div class="friend-middle" v-if="!contentEmpty">
                <div class="middle-content">
                  <div class="item-list">
                    <img
                      class="item-avatar"
                      :src="
                        appearanceInfo.avatar
                          ? appearanceInfo.avatar
                          : 'https://wesociastg.blob.core.chinacloudapi.cn/wesocial-uat/uploads/2021/3/12/5f136e2f41e2a7e63540a03b1f4d6f56.png'
                      "
                      alt=""
                    />
                    <div class="item-right">
                      <div class="item-name">
                        {{ appearanceInfo.name || "员工姓名" }}
                      </div>
                      <div class="item-content" v-html="textStr"></div>
                      <!-- 图片-->
                      <div
                        class="item-img-box"
                        :class="{
                          'one-img-box': info.content.picture.length == 1,
                          'four-img-box': info.content.picture.length == 4
                        }"
                        v-if="
                          info.content.picture &&
                            info.content.picture.length > 0
                        "
                      >
                        <img
                          class="item-img"
                          v-for="(itemImg, index) in info.content.picture"
                          :key="index"
                          :src="itemImg.url"
                          alt=""
                        />
                      </div>
                      <!-- 视频-->
                      <div
                        class="item-video"
                        v-if="info.content.video && info.content.video.url"
                      >
                        <img
                          @click="openVideo(info.content.video.url)"
                          class="play-icon"
                          src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/play_icon.png"
                          alt=""
                        />
                        <video
                          width="145"
                          height="180"
                          x5-video-player-fullscreen="true"
                          x5-playsinline
                          playsinline
                          webkit-playsinline
                          muted
                          preload="auto"
                          :src="info.content.video.url"
                        ></video>
                      </div>
                      <!-- 链接-->
                      <div
                        class="item-link"
                        v-if="info.content.link && info.content.link.link_url"
                      >
                        <img
                          class="link-img"
                          :src="info.content.link.cover_url"
                          alt=""
                        />
                        <div class="link-text">
                          {{ info.content.link.title }}
                        </div>
                      </div>
                      <!-- 小程序卡片-->
                      <div
                        id="launch-btn-box"
                        v-if="
                          info.content.applet && info.content.applet.app_name
                        "
                      >
                        <div class="item-card">
                          <div class="card-top">
                            <img
                              class="mini-img"
                              src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/mini-applet-icon.png"
                              alt=""
                            />
                            <div class="mini-name">
                              {{ info.content.applet.app_name }}
                            </div>
                          </div>
                          <div class="mini-content">
                            {{ info.content.applet.title }}
                          </div>
                          <img
                            class="mini-bg"
                            :src="info.content.applet.cover_url"
                            alt=""
                          />
                          <div class="card-bottom">
                            <img
                              class="mini-icon"
                              src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/mini-icon.png"
                              alt=""
                            />
                            <span class="mini-text">小程序</span>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="content-right">
            <cy-form
              class="main-box"
              ref="cform"
              :rules="rules"
              :model="info"
              label-width="90px"
            >
              <cy-form-item label="主题名称" prop="name" required>
                <cy-input
                  v-model="info.name"
                  clearable
                  maxlength="30"
                  placeholder="请输入主题名称"
                  show-word-limit
                ></cy-input
              ></cy-form-item>
              <cy-form-item
                label="员工"
                prop="member_ids"
                type="array"
                required
                trigger="change"
              >
                <el-button
                  v-if="this.staffList[0].staff.length == 0"
                  size="small"
                  type="primary"
                  @click="openSelectStaff"
                  >选择员工</el-button
                >
                <el-table
                  v-if="this.staffList[0].staff.length > 0"
                  :data="staffList"
                  style="width: 100%"
                  max-height="500"
                >
                  <el-table-column prop="staff" label="员工">
                    <template slot-scope="scope">
                      <div
                        class="item-staff"
                        v-for="(item, index) in scope.row.staff"
                        :key="index"
                      >
                        <img
                          class="staff-icon"
                          src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/staff_icon.png"
                          alt=""
                        />
                        <span>{{ item.name }}</span>
                      </div>
                    </template>
                  </el-table-column>
                  <el-table-column
                    fixed="right"
                    label="操作"
                    width="72"
                    v-if="!edit_member_id"
                  >
                    <template>
                      <el-button
                        @click.native.prevent="openSelectStaff"
                        type="text"
                      >
                        修改
                      </el-button>
                    </template>
                  </el-table-column>
                </el-table>
              </cy-form-item>
              <cy-form-item
                label="发布时间"
                prop="publish_type"
                required
                trigger="change"
              >
                <div :class="{ 'date-select-box': info.publish_type == 2 }">
                  <el-radio-group v-model="info.publish_type">
                    <el-radio
                      v-for="(item, index) in dateSelectArray"
                      :key="index"
                      :label="item.type"
                      >{{ item.content }}</el-radio
                    >
                  </el-radio-group>
                  <el-date-picker
                    v-if="info.publish_type == 2"
                    type="datetime"
                    class="time-picker"
                    placeholder="请选择日期时间"
                    v-model="info.publish_date"
                    @change="selectTime"
                    :picker-options="pickerOptions"
                  ></el-date-picker>
                </div>
              </cy-form-item>
              <cy-form-item label="发布内容" class="content-form">
                <div class="publish-content-box">
                  <!-- 表情、文本选择区域-->
                  <div class="select-box-insert">
                    <add-emoji @select="addEmoji" class="add-emoji">
                      <template slot="reference">
                        <el-button class="insert-btn" size="small"
                          ><img
                            class="btn-img"
                            src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/face_icon.png"
                            alt=""
                          /><span>插入表情</span></el-button
                        >
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
                      /><span>添加文本</span></el-button
                    >
                  </div>
                  <!--多行文本区域-->
                  <div>
                    <el-input
                      ref="textarea"
                      class="textarea-box"
                      placeholder="请输入消息内容"
                      v-model="info.content.text"
                      :rows="4"
                      type="textarea"
                      maxlength="1000"
                      show-word-limit
                    >
                    </el-input>
                  </div>
                  <!-- 图片、视频、链接、小程序选择区域-->
                  <div class="material-select-box">
                    <!-- 右侧素材综合选择展示区域-->
                    <el-popover
                      placement="top-start"
                      width="224"
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
                        /><span>图片/链接/小程序/视频</span></el-button
                      >
                    </el-popover>
                    <!-- 右侧图片选择展示区域-->
                    <div
                      class="select-picture-box"
                      v-if="currentMaterialType == 'picture'"
                    >
                      <div
                        class="item-picture"
                        v-for="(item, index) in info.content.picture"
                        :key="index"
                      >
                        <img class="item-img" :src="item.url" alt="" />
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
                        v-if="info.content.picture.length < 9"
                      >
                        <div class="center">
                          <i
                            class="el-icon-plus avatar-uploader-icon"
                            aria-hidden="true"
                          ></i>
                          <div class="tips" v-text="'上传图片'"></div>
                        </div>
                      </div>
                    </div>
                    <!-- 右侧视频选择展示区域-->
                    <div
                      class="select-video-box"
                      v-if="
                        currentMaterialType == 'video' &&
                          info.content.video &&
                          info.content.video.url
                      "
                    >
                      <img
                        @click="openVideo(info.content.video.url)"
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
                        :src="info.content.video.url"
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
                      v-if="
                        currentMaterialType == 'link' &&
                          info.content.link &&
                          (info.content.link.link_url ||
                            info.content.link.author)
                      "
                    >
                      <div class="select-link-title">
                        {{ info.content.link.title }}
                      </div>
                      <div class="select-link-bottom">
                        <!-- 为了解决flex布局下超出省略失效问题，所以嵌套-->
                        <div>
                          <span class="select-link-description">{{
                            info.content.link.description
                          }}</span>
                        </div>
                        <img
                          class="select-link-cover"
                          :src="info.content.link.cover_url"
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
                      v-if="
                        currentMaterialType == 'applet' &&
                          info.content.applet &&
                          info.content.applet.app_name
                      "
                    >
                      <div class="select-applet-name">
                        <img
                          class="mini-img"
                          src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/mini-applet-icon.png"
                          alt=""
                        /><span>{{ info.content.applet.app_name }}</span>
                      </div>
                      <div class="select-applet-title">
                        {{ info.content.applet.title }}
                      </div>
                      <img
                        class="select-applet-img"
                        :src="info.content.applet.cover_url"
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
                </div>
                <div class="error-msg" v-if="showError">请填写发布内容</div>
              </cy-form-item>
            </cy-form>
          </div>
        </div>
      </template>
      <template slot="footer" class="dialog-footer">
        <el-button type="primary" @click="save()" size="small">发 布</el-button>
        <el-button @click="cancel()" size="small">取 消</el-button>
      </template>
    </layout-form>
    <!-- 选择员工弹窗-->
    <selectMember
      v-if="showSelectStaff"
      :hasSelectArr="this.staffList[0].staff"
      @submit="getMember"
      @close="showSelectStaff = false"
    />
    <!-- 选择文本素材弹窗-->
    <select-text-material
      v-if="showSelectText"
      @textData="getTextData"
    ></select-text-material>
    <!-- 选择图片素材弹窗-->
    <select-picture-material
      v-if="showMaterialModal == 'picture'"
      @submit="getSelectMaterial"
    ></select-picture-material>
    <!-- 选择小程序素材弹窗-->
    <select-applet-material
      v-if="showMaterialModal == 'applet'"
      @submit="getSelectMaterial"
    ></select-applet-material>
    <!-- 选择链接素材弹窗-->
    <select-link-material
      v-if="showMaterialModal == 'link'"
      @submit="getSelectMaterial"
    ></select-link-material>
    <!-- 选择视频素材弹窗-->
    <select-video-material
      v-if="showMaterialModal == 'video'"
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
import { Component, Vue, Watch } from "vue-property-decorator";
import UploadImage from "@/components/business/upload-image/index.vue";
import API from "@/api/enterprise/friend-circle";
// 选择表情弹窗
import AddEmoji from "../../components/add-emoji.vue";
// 选择员工弹窗
import SelectMember from "@/components/business/select-member/index.vue";
// 选择文本素材弹窗
import selectTextMaterial from "../../material-library/components/text-material-dialog.vue";
// 选择图片素材弹窗
import selectPictureMaterial from "../../material-library/components/pic-mterial-dialog.vue";
// 选择小程序素材弹窗
import selectAppletMaterial from "../../material-library/components/small-mterial-dialog.vue";
// 选择链接素材弹窗
import selectLinkMaterial from "../../material-library/components/link-material-dialog.vue";
// 选择视频素材弹窗
import selectVideoMaterial from "../../material-library/components/video-material-dialog.vue";
import openVideo from "../../material-library/components/com-open-video.vue";
@Component({
  components: {
    UploadImage,
    AddEmoji,
    SelectMember,
    selectTextMaterial,
    selectPictureMaterial,
    selectAppletMaterial,
    selectLinkMaterial,
    selectVideoMaterial,
    openVideo
  }
})
export default class extends Vue {
  /**************************************************data开始**************************************************/
  private is_edit = 0; // 是否是编辑
  private id = 0; // 朋友圈id
  private edit_member_id = 0; // 当前员工的member_id,编辑单人的单条朋友圈时，必传
  private moment_topic_member_id = 0; // 编辑单人的单条朋友圈时，必传
  private mark: any = ""; // 链接加密mark
  private materialSelectList: Array<any> = [
    {
      type: "picture",
      hoverIcon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-img-icon.png",
      icon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-img-hover-icon.png",
      tip: "图片"
    },
    {
      type: "link",
      hoverIcon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-link-icon.png",
      icon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-link-hover-icon.png",
      tip: "链接"
    },
    {
      type: "applet",
      hoverIcon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-mini-icon.png",
      icon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-mini-hover-icon.png",
      tip: "小程序"
    },
    {
      type: "video",
      hoverIcon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-video-icon.png",
      icon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-video-hover-icon.png",
      tip: "视频"
    }
  ];
  private contentEmpty = false; // 判断发布内容是否为空
  private showError = false; // 控制发布内容错误提示是否展示
  private currentMaterialType = ""; // 当前选择素材的类型
  private showMaterialModal = ""; // 控制当前选择素材的弹框显隐
  private dateSelectArray: Array<any> = [
    {
      content: "立即发送",
      type: 1
    },
    {
      content: "定时发送",
      type: 2
    }
  ];
  private isOpenVideo = false; //是否打开视频
  private videoUrl = ""; // 当前播放视频的url
  // 日期时间选择器今天之前置灰
  private pickerOptions = {
    disabledDate(time: Date) {
      return time.getTime() < Date.now() - 8.64e7;
    }
  };
  private staffList: Array<any> = [
    {
      staff: []
    }
  ];
  private checkName = (rule: any, value: any, callback: any) => {
    if (!value) {
      return callback(new Error("请输入主题名称"));
    }
    callback();
  };
  private rules = {};
  private appearanceInfo = {
    avatar: "", // 头像
    background: "", // 背景图
    name: "", // 名称
    sign: "" // 签名
  };
  private info: any = {
    name: "",
    member_ids: [],
    publish_type: 1,
    publish_date: "",
    publish_time: 0,
    content: {
      text: "",
      video: {},
      applet: {},
      link: {},
      picture: []
    }
  };
  private showSelectStaff = false; // 控制选择员工弹框显隐
  private showSelectText = false; // 控制选择文本弹框显隐
  /**************************************************data结束**************************************************/
  /***********************************************methods开始**************************************************/
  // 选人组件回调
  private getMember(memberList: any) {
    this.showSelectStaff = false;
    console.log("组件传递回来的数据-> ", memberList);
    if (memberList.length > 0) {
      const list: Array<any> = [];
      memberList.map((item: any) => {
        const value = {
          name: item.name,
          user_id: item.user_id,
          member_id: item.member_id,
          is_dp: false
        };
        list.push(value);
      });
      this.staffList[0].staff = list;
      this.info.member_ids = list.map((item: any) => item.member_id);
    } else {
      this.staffList[0].staff = [];
      this.info.member_ids = [];
    }
    (this.$refs["cform"] as any).validateField("member_ids");
  }
  // 选择素材弹窗后回调
  private getSelectMaterial(material?: any) {
    console.log("material", material);
    // 选择图片点击取消按钮回调处理
    if (
      this.currentMaterialType == "picture" &&
      (!material || material.length == 0) &&
      this.info.content.picture.length > 0
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
        if (material.length + this.info.content.picture.length > 9) {
          this.$message.warning("图片素材最多上传9张");
          return;
        }
        this.info.content.picture.push(...material);
        this.showMaterialModal = "";
        break;
      case "video":
      case "link":
      case "applet":
        this.info.content[this.showMaterialModal] = material;
        this.showMaterialModal = "";
        break;
    }
  }
  // 预览视频
  private openVideo(url: string) {
    this.videoUrl = url;
    this.isOpenVideo = true;
  }
  // 关闭预览视频
  private cancelVideo() {
    this.isOpenVideo = false;
  }
  // 选择时间修改回调
  private selectTime() {
    this.info.publish_time = new Date(this.info.publish_date).getTime() / 1000;
  }
  // 打开选择员工弹窗
  private openSelectStaff() {
    this.showSelectStaff = true;
  }
  // 打开选择文本弹窗
  private openSelectText() {
    this.showSelectText = true;
  }
  // 鼠标划入事件
  private selectStyle(item: any) {
    const _this = this;
    this.$nextTick(function() {
      this.materialSelectList.forEach(function(item: any) {
        _this.$set(item, "active", false);
      });
      this.$set(item, "active", true);
    });
  }
  // 鼠标划出事件
  private outStyle(item: any) {
    this.$set(item, "active", false);
  }
  // 打开对应素材弹窗
  private openSelectMaterial(type: string) {
    const _this = this;
    this.materialSelectList.forEach(function(item: any) {
      _this.$set(item, "active", false);
    });
    this.currentMaterialType = type;
    this.showMaterialModal = type;
  }
  // 点击素材的叉叉按钮
  private deleteMaterial(type: string, index?: number) {
    switch (type) {
      case "picture":
        this.info.content.picture.splice(index, 1);
        if (this.info.content.picture.length <= 0)
          this.currentMaterialType = "";
        break;
      case "video":
      case "link":
      case "applet":
        this.info.content[type] = {};
        this.currentMaterialType = "";
        break;
    }
  }
  // 添加表情/添加文本
  private async addEmoji(emj: any) {
    const myField = (this.$refs.textarea as any).$refs.textarea;
    let strTem: string = emj || "";
    const surplus: number = 1000 - this.info.content.text.length;
    if (surplus - strTem.length < 0) {
      strTem = strTem.substring(0, surplus);
      this.$message.warning("发布内容长度不能超过1000字符，超出部分则被截断。");
    }
    if (myField.selectionStart || myField.selectionStart === 0) {
      const startPos = myField.selectionStart;
      const endPos = myField.selectionEnd;
      this.info.content.text =
        myField.value.substring(0, startPos) +
        strTem +
        myField.value.substring(endPos, myField.value.length);
      await this.$nextTick();
      myField.focus();
      myField.setSelectionRange(endPos + strTem.length, endPos + strTem.length);
    } else {
      this.info.content.text += strTem;
    }
  }
  // 添加文本回调
  private getTextData(data: any) {
    this.showSelectText = false;
    if (data) {
      this.addEmoji(data.content);
    }
  }
  // 处理文本内容
  public get textStr(): string {
    return this.info.content.text
      .replace(/\n/g, "<br/>")
      .replace(/\s/g, "&nbsp;");
  }
  // 新增朋友圈接口
  private createFriend() {
    const params = {
      ...this.info
    };
    API.momentTopicCreateApi(params).then((res: any) => {
      if (res.code === 200) {
        this.$message.success("保存成功");
        this.$router.go(-1);
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 编辑多人朋友圈接口
  private editFriend() {
    const params = {
      ...this.info,
      id: this.id
    };
    API.momentTopicEditApi(params).then((res: any) => {
      if (res.code === 200) {
        this.$message.success("保存成功");
        this.$router.go(-1);
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 编辑单人朋友圈接口
  private editSingleFriend() {
    const params = {
      ...this.info,
      moment_topic_id: this.id,
      moment_topic_member_id: this.moment_topic_member_id || undefined
    };
    API.memberMomentEditApi(params).then((res: any) => {
      if (res.code === 200) {
        this.$message.success("保存成功");
        this.$router.go(-1);
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 获取朋友圈详情
  private getDeatail() {
    const params = {
      edit_member_id: this.edit_member_id || undefined,
      id: this.id,
      page: 1,
      size: 9999
    };
    API.momentTopicDisplayApi(params).then((res: any) => {
      if (res.code === 200) {
        if (res.data.publish_time) {
          res.data.publish_date = new Date(res.data.publish_time * 1000);
        }
        if (res.data.list && res.data.list.length > 0 && !this.edit_member_id) {
          const memberList: Array<any> = [];
          res.data.list.map((item: any) => {
            const value = {
              name: item.user_name,
              user_id: item.userid,
              member_id: item.member_id,
              is_dp: false
            };
            memberList.push(value);
          });
          this.staffList[0].staff = memberList;
          this.info.member_ids = memberList.map((item: any) => item.member_id);
        }
        const materialType = ["picture", "link", "applet", "video"];
        materialType.forEach((item: string) => {
          if (
            res.data.content[item] instanceof Array &&
            res.data.content[item].length > 0
          ) {
            this.currentMaterialType = "picture";
          } else if (
            res.data.content[item] &&
            Object.keys(res.data.content[item]).length > 0
          ) {
            this.currentMaterialType = item;
          } else {
            res.data.content[item] = item == "picture" ? [] : {};
          }
          // console.log("currentMaterialType", res.data.content[item] instanceof Object, this.currentMaterialType);
        });
        this.info = { ...this.info, ...res.data };
      } else {
        this.$message.error(res.message);
      }
    });
  }
  //获取朋友圈形象(左侧)
  private friendDetailRight() {
    const params = {
      mark: this.mark
    };
    API.momentAppearanceGetApi(params)
      .then((res: any) => {
        if (res.code === 200) {
          this.appearanceInfo = res.data;
        }
      })
      .catch((error: any) => {
        console.log(error);
      });
  }
  // 点击取消按钮
  private cancel() {
    this.$router.go(-1);
  }
  // 点击发布按钮
  private save() {
    (this.$refs["cform"] as any).validate((valid: any) => {
      if (valid) {
        if (this.info.member_ids.length == 0) {
          this.$message({
            message: "请选择员工",
            type: "warning",
            duration: 2000
          });
          return;
        }
        if (this.info.publish_type == 2 && !this.info.publish_date) {
          this.$message({
            message: "请选择定时发送时间",
            type: "warning",
            duration: 2000
          });
          return;
        }
        if (
          this.info.publish_type == 2 &&
          this.info.publish_time &&
          this.info.publish_time <= new Date().getTime() / 1000
        ) {
          this.$message({
            message: "定时发送时间必须大于当前时间",
            type: "warning",
            duration: 2000
          });
          return;
        }
        if (this.contentEmpty) {
          this.showError = true;
          return;
        }
        if (!this.is_edit) {
          // 新增朋友圈
          this.createFriend();
        } else if (this.edit_member_id && this.is_edit) {
          // 编辑单人单个朋友圈
          this.editSingleFriend();
        } else {
          // 编辑多人的单个朋友圈
          this.editFriend();
        }
      } else {
        return;
      }
    });
  }
  @Watch("info.content", { deep: true })
  handleWatch(content: any) {
    if (
      content.text ||
      (content.picture && content.picture.length > 0) ||
      (content.link && content.link.link_url) ||
      (content.link && content.link.author) ||
      (content.video && content.video.url) ||
      (content.applet && content.applet.app_name)
    ) {
      this.contentEmpty = false;
      this.showError = false;
    } else {
      this.contentEmpty = true;
      this.showError = true;
    }
  }
  /**************************************************methods结束***********************************************/
  /**************************************************生命周期结束***********************************************/
  private created() {
    const {
      is_edit,
      member_id,
      moment_topic_id,
      moment_topic_member_id,
      user_name,
      mark
    } = this.$route.query;
    this.is_edit = Number(is_edit) || 0;
    this.id = Number(moment_topic_id) || 0;
    this.edit_member_id = Number(member_id) || 0;
    this.moment_topic_member_id = Number(moment_topic_member_id) || 0;
    // 选择员工模块展示该员工
    if (member_id && user_name) {
      const list: Array<any> = [
        {
          name: user_name,
          member_id: member_id,
          is_dp: false
        }
      ];
      this.staffList[0].staff = list;
      this.info.member_ids = list.map((item: any) => parseInt(item.member_id));
    }
    if (!this.is_edit) {
      this.contentEmpty = true;
      this.showError = false;
    }
    if (this.id) {
      this.getDeatail();
    }
    if (mark) {
      this.mark = mark || "";
      this.friendDetailRight();
    }
  }
}
</script>
<style lang="scss" scoped>
div {
  box-sizing: border-box;
}
.select-material-list {
  width: 224px;
  height: 70px;
  display: flex;
  align-items: center;
  justify-content: center;
  .item-select {
    text-align: center;
    margin-right: 20px;
    cursor: pointer;
    .item-icon {
      width: 38px;
      height: 38px;
      display: block;
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
.friend-circle-edit {
  .content-box {
    display: flex;
    .content-left {
      flex: 0 0 405px;
      width: 405px;
      .phone-box {
        width: 375px;
        height: 667px;
        border: solid 1px #ebeef5;
        overflow: auto;
        box-sizing: content-box;
        -ms-overflow-style: none;
        scrollbar-width: none;
      }
      .friend-top {
        background-color: #dcdfe6;
        width: 375px;
        height: 320px;
        position: relative;
        background-image: url("https://wesociastg.blob.core.chinacloudapi.cn/wesocial-uat/uploads/2021/3/12/f0ea071d94dd85390d45c7d249f72951.jpeg");
        background-size: cover;
        .mobile-phone-top {
          width: 375px;
          height: 22px;
        }
        .head-portrait-box {
          width: 71px;
          height: 71px;
          position: absolute;
          right: 15px;
          bottom: -35px;
          // background-color: #f5f6fa;
        }
        .staff-name {
          position: absolute;
          font-size: 16px;
          color: #ffffff;
          right: 106px;
          bottom: 8px;
        }
        .signature {
          position: absolute;
          font-size: 14px;
          color: #606266;
          right: 15px;
          left: 15px;
          top: 364px;
          text-align: right;
        }
      }
      .friend-middle {
        display: flex;
        position: relative;
        padding-top: 106px;
        overflow: hidden;
        .item-list {
          padding: 15px 0 15px 15px;
          display: flex;
          .item-avatar {
            flex: 0 0 0 50px;
            width: 50px;
            height: 50px;
            background: #f0f0f0;
          }
          .item-right {
            flex: 1;
            margin-left: 10px;
            padding-right: 15px;
            .item-name {
              font-size: 14px;
              color: #333333;
              line-height: 20px;
              font-weight: 600;
            }
            .item-content {
              font-size: 14px;
              color: #333333;
              line-height: 20px;
              margin-top: 2px;
              word-break: break-all;
            }
            .item-img-box {
              margin-top: 10px;
              font-size: 0;
              .item-img {
                width: 80px;
                height: 80px;
                background-size: cover;
                background-color: #f0f0f0;
                margin-right: 10px;
                margin-bottom: 10px;
                display: inline-flex;
              }
            }
            .one-img-box {
              .item-img {
                width: 120px;
                height: 120px;
              }
            }
            .four-img-box {
              width: 200px;
            }
            .item-video {
              margin-top: 10px;
              position: relative;
              width: fit-content;
              .play-icon {
                cursor: pointer;
                position: absolute;
                left: 0;
                right: 0;
                bottom: 0;
                top: 0;
                margin: auto;
                width: 25px;
                height: 25px;
                z-index: 99;
              }
              video {
                display: inline-block;
                vertical-align: baseline;
                object-fit: fill;
              }
            }
            .item-link {
              width: 230px;
              height: 60px;
              background-color: #f5f6fa;
              display: flex;
              align-items: center;
              padding: 8px;
              margin-top: 23px;
              position: relative;
              .link-img {
                width: 44px;
                height: 44px;
                background-color: #d8d8d8;
              }
              .link-text {
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
                font-size: 14px;
                margin-left: 10px;
                color: #333333;
                font-weight: 600;
              }
            }
            .item-card {
              width: 150px;
              max-height: 234px;
              background-color: #f5f6fa;
              padding: 8px;
              margin-top: 15px;
              .card-top {
                display: flex;
                align-items: center;
                margin-bottom: 4px;
              }
              .mini-img {
                width: 22px;
                height: 22px;
                background-color: #d8d8d8;
                border-radius: 50%;
              }
              .mini-name {
                color: #909399;
                font-size: 14px;
                transform: scale(0.8);
                margin-left: 5px;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
              }
              .mini-content {
                overflow: hidden;
                text-overflow: ellipsis;
                display: -webkit-box;
                -webkit-line-clamp: 2;
                -webkit-box-orient: vertical;
                width: 132px;
                color: #606266;
                font-size: 12px;
                line-height: 15px;
              }
              .mini-bg {
                width: 135px;
                height: 135px;
                margin-top: 5px;
                background-color: #dcdfe6;
                display: block;
              }
              .card-bottom {
                display: flex;
                align-items: center;
                margin-top: 5px;
              }
              .mini-icon {
                width: 12px;
                height: 11px;
              }
              .mini-text {
                color: #606266;
                font-size: 14px;
                transform: scale(0.8);
                margin-left: 5px;
                position: relative;
                top: -1px;
                left: -3px;
              }
            }
          }
        }
      }
    }
    .content-right {
      flex: 1;
      border-left: solid 1px #dcdfe6;
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
      .content-form {
        /deep/.el-form-item__label .label:before {
          content: "*";
          color: #f67170;
          position: absolute;
          left: 2px;
          top: 50%;
          -webkit-transform: translateY(-45%);
          transform: translateY(-45%);
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
            border-radius: 14px !important;
            display: inline-flex;
            align-items: center;
            justify-content: center;
            padding: 0;
            font-size: 12px;
            width: 87px;
            height: 28px !important;
            background-color: $color-background;
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
            -ms-overflow-style: none;
            scrollbar-width: none;
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
        text-overflow: -o-ellipsis-lastline;
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
}
</style>
