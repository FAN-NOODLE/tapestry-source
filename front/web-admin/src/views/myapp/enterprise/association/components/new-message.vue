<template>
  <div class="welcome">
    <layout-form>
      <template slot="content">
        <div class="content-box">
          <p class="time">
            {{ updateTime }}
          </p>
          <div class="content-left">
            <div class="content-msg">
              <p class="line" v-if="msgTextStr.length > 0">
                <span class="item-avatar-left"></span>
                <span class="item-content-right str" v-html="msgTextStr"></span>
              </p>
              <template
                v-if="info.content.fileData && info.content.fileData.length"
              >
                <p
                  class="line"
                  v-for="item in info.content.fileData"
                  :key="item.id"
                >
                  <span class="item-avatar-left"></span>
                  <span class="item-content-right str">
                    <span>
                      <div class="select-picture-box" v-if="item.fileType == 1">
                        <div class="item-picture">
                          <img class="item-img" :src="item.url" alt="" />
                        </div>
                      </div>
                      <div class="select-link-box" v-if="item.fileType == 2">
                        <div class="select-link-title">
                          {{ item.title }}
                        </div>
                        <div class="select-link-bottom">
                          <!-- 为了解决flex布局下超出省略失效问题，所以嵌套-->
                          <div>
                            <span class="select-link-description">{{
                              item.url ? item.url : item.link_url
                            }}</span>
                          </div>
                          <img
                            class="select-link-cover"
                            :src="
                              item.cover_url ? item.cover_url : item.pic_url
                            "
                            alt=""
                          />
                        </div>
                      </div>
                      <div class="select-applet-box" v-if="item.fileType == 3">
                        <div class="select-applet-name">
                          <img
                            class="minilogo"
                            src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/mini-applet-icon.png"
                            alt=""
                          />

                          <span class="smallName">{{ item.app_name }}</span>
                        </div>
                        <div class="select-applet-title">
                          {{ item.title }}
                        </div>
                        <img
                          class="select-applet-img"
                          :src="item.cover_url"
                          alt=""
                        />
                      </div>
                      <div class="select-picture-box" v-if="item.fileType == 4">
                        <div class="item-picture of-h">
                          <img
                            class="video-icon"
                            src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/play_icon.png"
                            alt=""
                          />
                          <video class="item-img video-box" :src="item.url" />
                        </div>
                      </div>
                      <div class="select-file-box" v-if="item.fileType == 5">
                        <div
                          class="item-file"
                          :class="`file-bg${item.file_type}`"
                        >
                          {{
                            item.file_type != 5
                              ? customIcon[item.file_type - 1]
                              : ""
                          }}
                        </div>
                      </div>
                    </span>
                  </span>
                </p>
              </template>
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
              <cy-form-item label="群发主题" prop="name" required>
                <cy-input
                  v-model="info.name"
                  clearable
                  maxlength="30"
                  placeholder="请输入群发主题"
                  show-word-limit
                ></cy-input>
              </cy-form-item>
              <cy-form-item label="群发类型" prop="send_type" required>
                <el-radio-group
                  v-model="info.send_type"
                  @change="massTypeChange"
                >
                  <el-radio
                    v-for="(item, index) in massSelectArray"
                    :key="index"
                    :label="item.type"
                    >{{ item.content }}</el-radio
                  >
                </el-radio-group>
              </cy-form-item>
              <cy-form-item
                label="执行者"
                prop="member_ids"
                type="array"
                required
                trigger="change"
              >
                <el-button
                  v-if="this.staffList[0].staff.length == 0"
                  size="small"
                  class="font-14"
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
                  <el-table-column prop="staff" label="群主">
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
                  <el-table-column fixed="right" label="操作" width="72">
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
                prop="type"
                required
                trigger="change"
              >
                <div :class="{ 'date-select-box': info.type == 2 }">
                  <el-radio-group v-model="info.type">
                    <el-radio
                      v-for="(item, index) in dateSelectArray"
                      :key="index"
                      :label="item.type"
                      >{{ item.content }}</el-radio
                    >
                  </el-radio-group>
                  <el-date-picker
                    v-if="info.type == 2"
                    type="datetime"
                    class="time-picker"
                    placeholder="请选择日期时间"
                    v-model="publish_time"
                    @change="selectTime"
                    :picker-options="pickerOptions"
                  ></el-date-picker>
                </div>
              </cy-form-item>
              <cy-form-item label="消息内容" required prop="content.text">
                <div
                  class="publish-content-box"
                  :class="{ 'border-color': showError }"
                >
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
                      maxlength="500"
                      show-word-limit
                    >
                    </el-input>
                  </div>
                </div>
                <div class="error-msg" v-if="showError">请输入消息内容</div>
              </cy-form-item>
              <cy-form-item
                label="消息素材"
                required
                class="m-r-30"
                prop="content.fileData"
                type="array"
                trigger="change"
                style="max-width: none;"
              >
                <!-- 图片、视频、链接、小程序选择区域-->
                <div class="material-select-box">
                  <el-popover placement="top-start" width="300" trigger="hover">
                    <div class="select-material-list">
                      <div
                        class="item-select"
                        v-for="(item, index) in materialSelectList"
                        @mouseover="selectStyle(item)"
                        @mouseout="outStyle(item)"
                        @click="openDialog(item.dialog, item.type)"
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
                      size="small"
                      slot="reference"
                      class="font-14"
                      type="primary"
                      >设置消息素材</el-button
                    >
                  </el-popover>
                </div>
                <div>
                  <div class="text-info font-12">
                    支持最多9个不同类型的素材。
                  </div>
                  <div class="m-l-20" v-if="info.content.fileData.length">
                    <div>
                      已选<span class="text-num"
                        >{{ info.content.fileData.length }}/9</span
                      >个素材
                      <span class="text-all">
                        <el-popconfirm
                          icon="el-icon-info"
                          confirm-button-text="确认"
                          cancel-button-type="Primary"
                          @onConfirm="deleteAll"
                          @confirm="deleteAll"
                          title="你确认移除全部素材吗？"
                        >
                          <el-button
                            type="text"
                            size="small"
                            class="text-all"
                            slot="reference"
                            >全部移除</el-button
                          >
                        </el-popconfirm>
                      </span>
                    </div>
                    <el-table
                      :data="info.content.fileData"
                      stripe
                      row-key="sid"
                    >
                      <el-table-column label="发送顺序" min-width="100">
                        <template slot-scope="{ $index }">
                          <span>{{ $index + 1 }} </span>
                        </template>
                      </el-table-column>
                      <el-table-column label="素材名称" min-width="220">
                        <template slot-scope="{ row }">
                          <span
                            v-if="
                              (row.fileType == 1 ||
                                row.fileType == 4 ||
                                row.fileType == 5) &&
                                !messageId
                            "
                            >{{ row.name ? row.name : "--" }}</span
                          >
                          <span v-else
                            >{{ row.title ? row.title : "--" }}
                          </span>
                        </template>
                      </el-table-column>
                      <el-table-column
                        label="素材类型"
                        prop="fileType"
                        min-width="120"
                      >
                        <template slot-scope="{ row }">
                          <span>{{
                            materialSelectList[row.fileType - 1].tip
                          }}</span>
                        </template>
                      </el-table-column>
                      <el-table-column
                        fixed="right"
                        label="操作"
                        min-width="100"
                      >
                        <template slot-scope="scope">
                          <el-popconfirm
                            icon="el-icon-info"
                            confirm-button-text="确认"
                            cancel-button-type="Primary"
                            @onConfirm="deleteRow(scope)"
                            @confirm="deleteRow(scope)"
                            title="你确认删除该素材吗？"
                          >
                            <el-button
                              size="small"
                              slot="reference"
                              type="text"
                              class="text-all"
                              >移除</el-button
                            >
                          </el-popconfirm>
                        </template>
                      </el-table-column>
                    </el-table>
                  </div>
                </div>
              </cy-form-item>
            </cy-form>
          </div>
        </div>
      </template>
      <template slot="footer" class="dialog-footer">
        <el-button
          @click="goBack()"
          class="font-14"
          v-if="this.messageId"
          size="small"
          >返 回</el-button
        >
        <el-button type="primary" @click="save()" class="font-14" size="small"
          >保 存</el-button
        >
        <el-button
          @click="cancel()"
          class="font-14"
          v-if="!this.messageId"
          size="small"
          >取 消</el-button
        >
      </template>
    </layout-form>
    <!-- 选择群主弹窗-->
    <selectMember
      v-if="showSelectStaff"
      :hasSelectArr="this.staffList[0].staff"
      @submit="getMember"
      dialogTitle="选择员工"
      @close="showSelectStaff = false"
    />
    <!-- 编辑离开页面提示 -->
    <edit-back-pop
      v-if="isEditPopDialog"
      @editGoBack="editGoBackFn"
    ></edit-back-pop>
    <!-- 新增离开页面提示 -->
    <back-pop v-if="isPopDialog" @createGoBack="createGoBackFn"></back-pop>
    <!-- 选择文本素材弹窗-->
    <select-text-material
      v-if="textDialog"
      @textData="getTextData"
    ></select-text-material>
    <!-- 开发中需求变动 -->
    <!-- 选择素材弹窗 -->
    <select-material
      v-if="dialog"
      :chooseInfo="info.content.fileData"
      :file="true"
      :type="infoType"
      @submit="getFileData"
    ></select-material>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Watch } from "vue-property-decorator";
// 表情
import AddEmoji from "../../components/add-emoji.vue";
// 时间戳转换
import { dateFormat } from "@/views/myapp/customer/auto-market/util/index";
// 新增离开页面提示
import backPop from "../../components/back-pop.vue";
// 编辑离开页面提示
import editBackPop from "../../components/edit-back-pop.vue";
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
// 选择文件素材弹窗
import selectFileMterial from "../../material-library/components/file-mterial-dialog.vue";
import selectMaterial from "@/components/business/select-material/select-material.vue";
import {
  groupSendAdd,
  taskInfoApi,
  groupSendEdit
} from "@/api/enterprise/association";
import { ApiResult } from "@/types";
@Component({
  components: {
    AddEmoji,
    SelectMember,
    editBackPop,
    backPop,
    selectTextMaterial,
    selectPictureMaterial,
    selectAppletMaterial,
    selectLinkMaterial,
    selectVideoMaterial,
    selectFileMterial,
    selectMaterial
  }
})
export default class extends Vue {
  // 所选群主数据
  private staffList: Array<any> = [
    {
      staff: []
    }
  ];
  private customIcon = ["W", "X", "P", "PDF"]; //自定义icon
  // 消息ID
  private messageId: any = 0;
  // 弹窗tab类型
  private infoType: any = 1;
  // 弹窗开关
  private dialog: any = false;
  // 消息内容素材相关数据
  private materialSelectList: any = [
    {
      type: 1,
      hoverIcon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-img-icon.png",
      icon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-img-hover-icon.png",
      tip: "图片",
      dialog: "pictureDialog"
    },
    {
      type: 2,
      hoverIcon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-link-icon.png",
      icon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-link-hover-icon.png",
      tip: "链接",
      dialog: "urlDialog"
    },
    {
      type: 3,
      hoverIcon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-mini-icon.png",
      icon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-mini-hover-icon.png",
      tip: "小程序",
      dialog: "appletDialog"
    },
    {
      type: 4,
      hoverIcon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-video-icon.png",
      icon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-video-hover-icon.png",
      tip: "视频",
      dialog: "videoDialog"
    },
    {
      type: 5,
      hoverIcon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-file-icon.png",
      icon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-file-hover-icon.png",
      tip: "文件夹",
      dialog: "fileDialog"
    }
  ];
  // 发布时间
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
  // 群发类型
  private massSelectArray: Array<any> = [
    {
      content: "群发客户",
      type: 1
    },
    {
      content: "群发客户群",
      type: 2
    }
  ];
  private rules = {};
  private showError = false; // 控制消息内容错误提示是否展示
  private publish_time = "";
  private info: any = {
    name: "",
    member_ids: [],
    type: 1, //1-立即 2-定时
    send_type: 1, //2-群发客户群 1-群发客户
    send_time: 0,
    content: {
      text: "",
      fileData: [],
      fileType: 0
    }
  };
  private timer: any = undefined;
  // 选择群主弹窗显隐
  private showSelectStaff = false;
  // 文本素材弹窗显隐
  private textDialog = false;
  // 图片素材弹窗显隐
  private pictureDialog = false;
  // 选择链接素材弹窗显隐
  private urlDialog = false;
  // 选择小程序素材弹窗显隐
  private appletDialog = false;
  // 选择视频弹窗显隐
  private videoDialog = false;
  // 选择文件弹窗显隐
  private fileDialog = false;
  // 新增离开页面提示显隐
  private isPopDialog = false;
  // 编辑离开页面提示显隐
  private isEditPopDialog = false;
  @Watch("info.content", { deep: true })
  handleWatch(content: any) {
    if (content.text || content.fileData) {
      this.showError = false;
    } else {
      this.showError = true;
    }
  }
  // --------------------------生命周期---------------------

  private mounted() {
    this.timer = setInterval(() => {
      this.updateTime = dateFormat(
        Number.parseInt(String(Date.now() / 1000)),
        "hh:mm"
      );
    }, 1000);
    const groupLocalData = window.sessionStorage.getItem("groupLocalData");
    const groupData = window.sessionStorage.getItem("groupstaffData");
    this.messageId = this.$route.query.id ? Number(this.$route.query.id) : null;
    if (this.messageId) {
      this.initData();
    } else if (!groupLocalData) {
      window.sessionStorage.setItem(
        "groupLocalData",
        JSON.stringify(this.packageData())
      );
    } else {
      this.analyzingData(JSON.parse(groupLocalData));
      if (groupData) {
        this.staffList[0].staff = JSON.parse(groupData);
      }
    }
  }
  private destroyed() {
    clearInterval(this.timer);
  }
  // --------------------------methods---------------------
  // 选择群发类型
  private massTypeChange() {
    this.staffList[0].staff = [];
  }
  // 移除全部内容
  private deleteAll() {
    this.info.content.fileData = [];
  }
  // 移除单行内容
  private deleteRow(e) {
    this.info.content.fileData.splice(e.$index, 1);
  }
  // 选择时间修改回调
  private selectTime() {
    this.info.send_time = new Date(this.publish_time).getTime() / 1000;
  }
  // 日期时间选择器今天之前置灰
  private pickerOptions = {
    disabledDate(time: Date) {
      return time.getTime() < Date.now() - 8.64e7;
    }
  };
  // 同步手机消息内容
  public get msgTextStr(): string {
    return this.info.content.text
      .replace(/\n/g, "<br/>")
      .replace(/\s/g, "&nbsp;");
  }
  // 当前时间 hh:mm
  private updateTime = dateFormat(
    Number.parseInt(String(Date.now() / 1000)),
    "hh:mm"
  );
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
  // 选人组件回调
  private getMember(memberList: any) {
    this.showSelectStaff = false;
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
    // 对部分表单字段进行校验
    (this.$refs["cform"] as any).validateField("member_ids");
  }

  // 打开选择群主弹窗
  private openSelectStaff() {
    this.showSelectStaff = true;
  }
  // 打开选择文本弹窗
  private openSelectText() {
    this.textDialog = true;
  }
  // 打开素材弹框
  private openDialog(dialog: string, type: number) {
    this.dialog = true;
    this.infoType = type;
  }
  // 文本弹框回调
  private getTextData(data: any) {
    this.textDialog = false;
    if (data) {
      this.addEmoji(this.getDataFromArrayOrObject(data).content);
    }
  }
  // 判断类型
  private getDataFromArrayOrObject(data: any) {
    if (Object.prototype.toString.call(data) == "[object Array]") {
      return data[0];
    } else if (Object.prototype.toString.call(data) == "[object Object]") {
      return data;
    } else {
      return null;
    }
  }
  // 素材弹框回调
  private getFileData(data: any) {
    this.dialog = false;
    this.$nextTick(() => {
      data && (this.$refs["cform"] as any).clearValidate(["content.fileData"]);
    });
    console.log(data, "数据");
    data.forEach(item => {
      if (item.fileType == 1 || item.fileType == 4 || item.fileType == 5) {
        item.title = item.name ? item.name : item.title;
      }
    });
    this.info.content.fileData = JSON.parse(JSON.stringify(data));
  }
  // 素材点击删除
  private deleteMaterial() {
    this.info.content.fileType = 0;
    this.info.content.fileData = null;
  }
  // 添加表情
  async addEmoji(emj: any) {
    const myField = (this.$refs.textarea as any).$refs.textarea;
    let strTem: string = emj || "";
    const surplus: number = 500 - this.info.content.text.length;
    if (surplus - strTem.length < 0) {
      strTem = strTem.substring(0, surplus);
      this.$message.warning("消息内容长度不能超过500字符，超出部分则被截断。");
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
  // 新增挽留
  private createGoBackFn(type: string) {
    this.isPopDialog = false;
    if (type == "save") {
      //数据缓存到本地；
      window.sessionStorage.setItem(
        "groupLocalData",
        JSON.stringify(this.packageData())
      );
      window.sessionStorage.setItem(
        "groupstaffData",
        JSON.stringify(this.staffList[0].staff)
      );
      this.$router.push({
        name: "myapp_enterprise_association_group-sent-message_list"
      });
    } else if (type == "cancel") {
      //直接离开
      window.sessionStorage.removeItem("groupLocalData");
      this.$router.push({
        name: "myapp_enterprise_association_group-sent-message_list"
      });
    }
  }
  // 编辑挽留
  editGoBackFn(type: string) {
    this.isEditPopDialog = false;
    if (type != "close") {
      this.$router.push({
        name: "myapp_enterprise_association_group-sent-message_list"
      });
    }
  }
  // 新增、点击取消按钮
  private cancel() {
    const local: any = window.sessionStorage.getItem("groupLocalData");
    const compare: any = JSON.stringify(this.packageData());
    // console.log(local);
    // console.log(compare);
    if (local === compare) {
      this.$router.push({
        name: "myapp_enterprise_association_group-sent-message_list"
      });
    } else {
      this.isPopDialog = true;
    }
  }
  // 编辑、返回
  private goBack() {
    const compare = window.sessionStorage.getItem("editdefaultData");
    // console.log(compare);
    // console.log(JSON.stringify(this.info));
    if (JSON.stringify(this.info) == compare) {
      this.$router.push({
        name: "myapp_enterprise_association_group-sent-message_list"
      });
    } else {
      this.isEditPopDialog = true;
    }
  }
  // 页面数据整合
  private packageData() {
    const data: any = {};
    data.text = {
      content: this.info.content.text
    };
    console.log(this.info.content.fileData);
    const arr = this.info.content.fileData.map(item => {
      const obj: any = {};
      if (item.fileType == 1) {
        obj["image"] = {
          fileType: item.fileType,
          pic_url: this.messageId
            ? item.pic_url
              ? item.pic_url
              : item.url
            : item.url,
          url: item.url,
          title: this.messageId ? item.title : item.name
        };
      } else if (item.fileType == 2) {
        obj["link"] = {
          fileType: item.fileType,
          pic_url: this.messageId
            ? item.pic_url
              ? item.pic_url
              : item.cover_url
            : item.cover_url,
          title: item.title,
          desc: this.messageId
            ? item.desc
              ? item.desc
              : item.description
            : item.description,
          url: this.messageId
            ? item.url
              ? item.url
              : item.link_url
            : item.link_url
        };
      } else if (item.fileType == 3) {
        obj["miniprogram"] = {
          fileType: item.fileType,
          appid: item.app_id,
          page: this.messageId
            ? item.page
              ? item.page
              : item.path
            : item.path,
          title: item.title,
          app_name: item.app_name,
          cover_url: item.cover_url
        };
      } else if (item.fileType == 4) {
        obj["video"] = {
          fileType: item.fileType,
          pic_url: this.messageId
            ? item.pic_url
              ? item.pic_url
              : item.url
            : item.url,
          title: item.title,
          url: item.url
        };
      } else if (item.fileType == 5) {
        obj["files"] = {
          fileType: item.fileType,
          file_type: item.file_type,
          pic_url: this.messageId
            ? item.pic_url
              ? item.pic_url
              : item.url
            : item.url,
          title: item.title,
          url: item.ur
        };
      }
      return obj;
    });
    data.content = arr;
    const params: any = {
      name: this.info.name,
      type: this.info.type, //1-立即 2-定时
      send_type: this.info.send_type,
      send_time: this.info.send_time,
      member_ids: this.info.member_ids,
      content: data
    };
    return params;
  }
  //详情赋值
  private analyzingData(data: any) {
    this.info.name = data.name;
    this.info.type = data.type;
    this.publish_time = data.send_time
      ? new Date(data.send_time * 1000) + ""
      : "";
    this.info.send_time = data.send_time;
    this.info.content.text = data.content.text.content;
    this.info.member_ids = data.member_ids;
    if (data.content.content && Array.isArray(data.content.content)) {
      this.info.content.fileData = data.content.content.map(function(
        item: any
      ) {
        const keys = Object.keys(item);
        const key = keys[0];
        return item[key];
      });
    } else {
      if ("image" in data.content) {
        this.info.content.fileType = 1;
        this.info.content.fileData = [
          {
            url: data.content.image.pic_url,
            pic_url: data.content.image.pic_url,
            fileType: 1
          }
        ];
      } else if ("link" in data.content) {
        this.info.content.fileType = 2;
        this.info.content.fileData = [
          {
            pic_url: data.content.link.pic_url,
            title: data.content.link.title,
            desc: data.content.link.desc,
            url: data.content.link.url,
            fileType: 2
          }
        ];
      } else if ("miniprogram" in data.content) {
        this.info.content.fileType = 3;
        this.info.content.fileData = [
          {
            app_id: data.content.miniprogram.appid,
            page: data.content.miniprogram.page,
            title: data.content.miniprogram.title,
            app_name: data.content.miniprogram.app_name,
            cover_url: data.content.miniprogram.cover_url,
            fileType: 3
          }
        ];
      }
    }
  }
  // 详情接口
  private async initData() {
    const { code, message = "", data }: ApiResult = await taskInfoApi(
      this.messageId
    );
    if (code !== 200) {
      this.$message.error(message);
      return;
    }
    this.analyzingData(data);
    this.staffList[0].staff = data.members;
    this.info.member_ids = data.members.map((item: any) => item.member_id);
    window.sessionStorage.setItem("editdefaultData", JSON.stringify(data));
  }
  // 新增接口
  private async addGroupSend() {
    const params = this.packageData();
    console.log(params);
    const { code, message = "" }: ApiResult = await groupSendAdd(params);
    if (code !== 200) {
      this.$message.error(message);
      return;
    }
    this.$message.success("保存成功");
    window.sessionStorage.removeItem("groupLocalData");
    window.sessionStorage.removeItem("groupstaffData");
    window.sessionStorage.removeItem("editdefaultData");
    console.log("移除缓存");
    this.$router.push({
      name: "myapp_enterprise_association_group-sent-message_list"
    });
  }
  // 编辑接口
  private async editGroupSend() {
    const params = this.packageData();
    params.id = this.messageId;
    const { code, message = "" }: ApiResult = await groupSendEdit(params);
    if (code !== 200) {
      this.$message.error(message);
      return;
    }
    this.$message.success("保存成功");
    this.$router.push({
      name: "myapp_enterprise_association_group-sent-message_list"
    });
  }
  // 点击保存按钮
  private save() {
    console.log(this.info);

    (this.$refs["cform"] as any).validate(async (valid: any) => {
      if (valid) {
        if (this.info.member_ids.length == 0) {
          this.$message({
            message: "请选择群主",
            type: "warning",
            duration: 2000
          });
          return;
        }
        if (this.info.type == 2 && !this.publish_time) {
          this.$message({
            message: "请选择定时发送时间",
            type: "warning",
            duration: 2000
          });
          return;
        }
        if (
          this.info.type == 2 &&
          this.info.send_time &&
          this.info.send_time <= new Date().getTime() / 1000
        ) {
          this.$message({
            message: "定时发送时间必须大于当前时间",
            type: "warning",
            duration: 2000
          });
          return;
        }
        if (this.info.content.text || this.info.content.fileData) {
          this.showError = false;
        } else {
          this.showError = true;
          return;
        }
        if (this.messageId) {
          this.editGroupSend();
        } else {
          this.addGroupSend();
        }
      }
    });
  }
}
</script>
<style lang="scss" scoped>
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
.welcome {
  box-shadow: 0px 0px 20px 1px rgba(0, 0, 0, 0.08);
  .el-table__body-wrapper::-webkit-scrollbar {
    width: 6px !important;
  }
}
.main-box-label {
  .in-item:last-child {
    .label::before {
      content: "*";
      color: #f67170;
      position: absolute;
      left: 2px;
      top: 50%;
      transform: translateY(-45%);
    }
  }
}
div {
  box-sizing: border-box;
}
.select-material-list {
  width: 300px;
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
}
.select-link-box {
  width: 225px;
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
    display: -webkit-box;
    -webkit-line-clamp: 1;
    -webkit-box-orient: vertical;
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
    margin-bottom: 5px;
    .minilogo {
      display: inline-block;
      margin-right: 5px;
      width: 20px;
      height: 20px;
      border-radius: 50%;
    }
    .smallName {
      color: #707277;
      line-height: 20px;
      overflow: hidden;
      display: -webkit-box;
      -webkit-line-clamp: 1;
      -webkit-box-orient: vertical;
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
.welcome {
  .content-box {
    position: relative;
    display: flex;
    .time {
      position: absolute;
      top: 1px;
      font-weight: bold;
      left: 156px;
      font-size: 14px;
      background: rgb(237, 237, 237);
      width: 80px;
      color: $color-text-normal;
      text-align: center;
    }
    .content-left {
      // 1125*2436
      flex: 0 0 390px;
      width: 390px;
      height: 694px;
      padding: 72px 0;
      margin-right: 30px;
      padding-left: 20px;
      padding-right: 20px;
      background: url("https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/welcome-phone.png") no-repeat
        100% top;
      background-size: 390px 694px;
      .content-msg::-webkit-scrollbar {
        width: 6px;
      }
      .content-msg {
        padding: 10px 0px;
        font-size: 12px;
        line-height: 18px;
        word-break: break-all;
        max-height: 500px;
        overflow-y: auto;
        padding-right: 2px;
        overflow-x: hidden;
        text-overflow: ellipsis;
        .line {
          display: flex;
          position: relative;
          margin: 20px 0;
          .item-avatar-left {
            min-width: 50px;
            width: 50px;
            height: 50px;
            background: #fff;
            line-height: 50px;
            font-size: 18px;
            font-weight: bold;
            color: #708df5;
            text-align: center;
            border-radius: 4px;
          }
          .item-content-right {
            margin-left: 15px;
            padding: 15px;
            background: #fff;
            border-radius: 4px;
            color: $color-text-normal;
            max-width: 250px;
          }
          .item-content-right::before {
            content: "";
            position: absolute;
            width: 0px;
            height: 0px;
            font-size: 0;
            line-height: 0;
            overflow: hidden;
            border-width: 20px;
            border-style: dashed dashed solid dashed;
            border-color: transparent #fff transparent transparent;
            left: 35px;
            top: 5px;
          }
        }
      }
    }
    .content-right {
      flex: 1;
      width: calc(100% - 390px);
      border-left: solid 1px #dcdfe6;
      padding-left: 30px;
      /deep/.in-item {
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
        // min-height: 425px;
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
            width: 87px;
            font-size: 12px;
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
          }
          /deep/.el-input__count {
            right: 0;
          }
        }
      }
    }
  }
}
.text {
  line-height: 16px;
}
.text-info {
  color: #c0c4cc;
}
.text-num {
  color: rgba(255, 0, 0, 0.8);
}
.text-all {
  color: #6f3112;
  font-weight: 500;
  font-size: 14px;
}
// /deep/ .content-right .el-table tr td .cell {
//   height: 105px !important;
// }
.of-h {
  overflow: hidden;
  display: flex;
  justify-content: center;
  align-items: center;
  background: #000;
  border-radius: 5px;
  .video-icon {
    position: absolute;
    left: 30px;
    top: 30px;
    width: 20px;
    height: 20px;
  }
  .video-box {
    height: auto !important;
  }
}
.item-file {
  width: 80px;
  height: 80px;
  text-align: center;
  color: #fff;
  font-size: 40px;
  line-height: 80px;
  border-radius: 5px;
}
.file-bg1 {
  background: #25589a;
}
.file-bg2 {
  background: #02843f;
}
.file-bg3 {
  background: #d34729;
}
.file-bg4 {
  background: #0e88e7;
}
.file-bg5 {
  background: #c9c9c9;
}
/deep/ .content-right .el-table tr td:first-child .cell {
  overflow-y: scroll !important;
}
</style>
