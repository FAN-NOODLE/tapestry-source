<template>
  <div class="com-task-center" v-loading.fullscreen.lock="loading">
    <layout-form>
      <template slot="content">
        <cy-form
          ref="cform"
          label-width="90px"
          :rules="taskRules"
          :model="taskInfo"
        >
          <!-- 任务名称 -->
          <cy-form-item
            label="任务名称"
            prop="task_title"
            ref="taskNameInfo"
            :required="status !== 3"
          >
            <cy-input
              v-if="status !== 3"
              size="small"
              width-size="medium"
              placeholder="请输入任务名称"
              v-model="taskInfo.task_title"
              class="input-with-select"
              maxlength="20"
              show-word-limit
              @keyup.native="trimLR()"
            ></cy-input>
            <div class="font-re show-content-text" v-if="status == 3">
              {{ taskInfo.task_title }}
            </div>
          </cy-form-item>
          <!-- 可见范围 -->
          <cy-form-item
            :class="
              taskInfo.range && taskInfo.range.length > 0
                ? ''
                : 'plan-form-item'
            "
            label="可见范围"
            :required="status !== 3"
          >
            <el-button
              v-if="status !== 3"
              class="range-choose-btn"
              size="small"
              @click="showChoiceStaff = true"
              >选择可见范围</el-button
            >
            <!-- 标题 -->
            <div
              v-if="
                (taskInfo.range && taskInfo.range.length > 0) || status == 3
              "
              class="title-item flex"
            >
              <div class="flex">
                <div class="font-me">
                  {{
                    status !== 3 ? "可见范围" : `按${taskInfo.rangeType}已选`
                  }}
                </div>
                <div class="font-re" v-if="status == 3">
                  共计{{ taskInfo.rangeSumNum }}人
                </div>
              </div>
              <div v-if="status !== 3">操作</div>
            </div>
            <!-- 可见范围显示区 -->
            <div
              v-if="taskInfo.range && taskInfo.range.length > 0"
              class="range-content flex flex-align-center"
            >
              <div class="left">
                <div
                  v-for="(rangeItem, rangeIndex) in taskInfo.range"
                  :key="rangeIndex"
                  class="flex item flex-align-center"
                >
                  <img src="@/assets/images/member_normal.png" alt="" />
                  <div class="item-name font-re">
                    {{ rangeItem.name }}
                  </div>
                </div>
              </div>
              <!-- 修改按钮 -->
              <div
                v-if="status != 3"
                class="right font-re pointer-cursor"
                @click="showChoiceStaff = true"
              >
                修改
              </div>
            </div>
          </cy-form-item>
          <!-- 任务时间 -->
          <cy-form-item
            label="任务时间"
            prop="taskTime"
            trigger="change"
            :required="status !== 3"
          >
            <el-date-picker
              v-if="status !== 3"
              style="height: 36px; width:400px !important;"
              class="activity-time"
              size="small"
              v-model="taskInfo.taskTime"
              type="datetimerange"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              value-format="timestamp"
              format="yyyy-MM-dd HH:mm"
              :picker-options="pickerOption"
              :default-time="['00:00:00', '23:59:59']"
            >
            </el-date-picker>
            <div class="font-re show-content-text" v-if="status == 3">
              {{
                (taskInfo.taskTime[0] / 1e3) | formatDate("yyyy-mm-dd hh:mm:ws")
              }}
              至
              {{
                (taskInfo.taskTime[1] / 1e3) | formatDate("yyyy-mm-dd hh:mm:ws")
              }}
            </div>
          </cy-form-item>
          <!-- 任务类型 -->
          <cy-form-item
            label="任务类型"
            prop="type"
            trigger="change"
            :required="status !== 3"
          >
            <cy-select
              v-if="status !== 3"
              v-model="taskInfo.type"
              size="medium"
              :optionsList="taskTypeList"
              :options="{ label: 'label', value: 'value' }"
              @change="typeChange"
            >
            </cy-select>
            <div class="font-re show-content-text" v-if="status == 3">
              {{ taskInfo.taskType }}
            </div>
          </cy-form-item>
          <!-- 客户标签 -->
          <cy-form-item
            :class="
              taskInfo.customerLabel && taskInfo.customerLabel.length > 0
                ? ''
                : 'plan-form-item'
            "
            v-if="taskInfo.type == 3"
            label="客户标签"
            :required="status !== 3"
          >
            <el-button
              v-if="status !== 3"
              class="range-choose-btn"
              size="small"
              @click="tagDialog = true"
              >选择标签</el-button
            >
            <!-- 标题 -->
            <div
              v-if="(showTagsList && showTagsList.length > 0) || status == 3"
              class="title-item flex"
            >
              <div class="flex">
                <div class="font-me">
                  已选标签
                </div>
              </div>
              <div v-if="status !== 3">操作</div>
            </div>
            <!-- 可见范围显示区 -->
            <div
              v-if="showTagsList && showTagsList.length > 0"
              class="range-content flex flex-align-center"
            >
              <div class="left">
                <div
                  v-for="(item, index) in showTagsList"
                  :key="index"
                  class="flex item flex-align-center"
                >
                  <div class="item-name font-re">
                    {{ item }}
                  </div>
                </div>
              </div>
              <!-- 修改按钮 -->
              <div
                v-if="status != 3"
                class="right font-re pointer-cursor"
                @click="tagDialog = true"
              >
                修改
              </div>
            </div>
          </cy-form-item>
          <!-- 发送提醒 -->
          <cy-form-item label="发送提醒" prop="sendReminder">
            <div class="flex sendReminderBox" v-if="status !== 3">
              <el-checkbox v-model="taskInfo.sendReminder.selectStart">
              </el-checkbox>
              开始时间前
              <el-input
                v-model="taskInfo.sendReminder.startTime"
                type="text"
                @input="changeInputStartTime()"
              />
              分钟
            </div>
            <div
              class="font-re show-content-text"
              v-if="status == 3 && taskInfo.sendReminder.startTime"
            >
              开始前 {{ taskInfo.sendReminder.startTime }} 分钟
            </div>
            <div class="flex sendReminderBox" v-if="status !== 3">
              <el-checkbox v-model="taskInfo.sendReminder.selectEnd">
              </el-checkbox>
              结束时间前
              <el-input
                v-model="taskInfo.sendReminder.endTime"
                type="text"
                @input="changeInputEndTime()"
              />
              分钟
            </div>
            <div
              class="font-re show-content-text"
              v-if="status == 3 && taskInfo.sendReminder.endTime"
            >
              结束前 {{ taskInfo.sendReminder.endTime }} 分钟
            </div>
            <div class="flex" v-if="status !== 3">
              <el-checkbox v-model="taskInfo.sendReminder.selectfixed">
              </el-checkbox>
              <p class="mlr-12">固定时间提醒</p>
              <el-date-picker
                v-model="taskInfo.sendReminder.fixedTime"
                size="small"
                type="datetime"
                placeholder="选择日期时间"
                format="yyyy-MM-dd HH:mm"
                :picker-options="pickerOption"
              >
              </el-date-picker>
            </div>
            <div
              class="font-re show-content-text"
              v-if="status == 3 && taskInfo.sendReminder.fixedTime"
            >
              固定时间
              {{
                taskInfo.sendReminder.fixedTime
                  | formatDate("yyyy-mm-dd hh:mm:ws")
              }}
            </div>
            <div
              class="font-re show-content-text"
              v-if="
                status == 3 &&
                  taskInfo.sendReminder.startTime == 0 &&
                  taskInfo.sendReminder.endTime == 0 &&
                  taskInfo.sendReminder.fixedTime == 0
              "
            >
              暂无发送提醒
            </div>
          </cy-form-item>
          <!-- 任务内容 -->
          <cy-form-item
            class="risk-form-item"
            label="任务内容"
            trigger="blur"
            prop="text"
            :required="status !== 3"
          >
            <!-- 新增素材展示区 -->
            <div v-if="status !== 3" class="risk-content-box">
              <!-- 表情选择区域-->
              <div class="select-box-insert">
                <add-emoji @select="addEmoji" class="add-emoji">
                  <template slot="reference">
                    <el-button class="insert-btn" size="small"
                      ><img
                        class="btn-img"
                        src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/face_icon.png"
                        alt=""
                      /><span id="insertEmoji">插入表情</span></el-button
                    >
                  </template>
                </add-emoji>
              </div>
              <!--多行文本区域-->
              <div>
                <el-input
                  ref="textarea"
                  class="textarea-box"
                  placeholder="请输入任务内容文字描述"
                  v-model="taskInfo.text"
                  :rows="4"
                  type="textarea"
                  maxlength="500"
                  show-word-limit
                ></el-input>
              </div>
              <!-- 图片、视频、小程序选择区域-->
              <div class="material-select-box">
                <!-- 右侧素材综合选择展示区域-->
                <el-popover
                  placement="top-start"
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
                    /><span v-show="materialSelectList.length === 3"
                      >图片/视频/小程序</span
                    ><span v-show="materialSelectList.length === 1"
                      >小程序</span
                    ></el-button
                  >
                </el-popover>
                <!-- 右侧图片选择展示区域-->
                <div
                  class="select-picture-box"
                  v-if="currentMaterialType == 'picture'"
                >
                  <div
                    class="item-picture"
                    v-for="(pItem, pIndex) in taskInfo.content.picture"
                    :key="pIndex"
                  >
                    <img class="item-img" :src="pItem.url" alt="" />
                    <div
                      class="delete-img-box"
                      @click="deleteMaterial('picture', pIndex)"
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
                    v-if="taskInfo.content.picture.length < 9"
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
                      taskInfo.content.video &&
                      taskInfo.content.video.url
                  "
                >
                  <img
                    @click="openVideo(taskInfo.content.video.url)"
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
                    :src="taskInfo.content.video.url"
                  ></video>
                  <div
                    class="delete-video-box"
                    @click="deleteMaterial('video', 1)"
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
                      taskInfo.content.applet &&
                      taskInfo.content.applet.app_name
                  "
                >
                  <div class="select-applet-name">
                    <img
                      class="mini-img"
                      src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/mini-applet-icon.png"
                      alt=""
                    /><span>{{ taskInfo.content.applet.app_name }}</span>
                  </div>
                  <div class="select-applet-title">
                    {{ taskInfo.content.applet.title }}
                  </div>
                  <img
                    class="select-applet-img"
                    :src="taskInfo.content.applet.cover_url"
                    alt=""
                  />
                  <div
                    class="delete-applet-box"
                    @click="deleteMaterial('applet', 1)"
                  >
                    <img
                      src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-delete-icon.png"
                      alt=""
                    />
                  </div>
                </div>
              </div>
            </div>
            <!-- 查看素材展示 -->
            <div class="font-re check-risk-content" v-if="status == 3">
              <!-- 任务内容 -->
              <div class="showRiskContent">
                {{ taskInfo.text }}
              </div>
              <!-- 图片 -->
              <div class="flex" v-if="currentMaterialType == 'picture'">
                <div
                  class="material-img-show"
                  v-for="picTtem in taskInfo.content.picture"
                  :key="picTtem.index"
                >
                  <img
                    @click.prevent="checkNowImage(picTtem.url)"
                    :src="picTtem.url"
                    alt=""
                  />
                </div>
              </div>
              <!-- 视频 -->
              <div class="flex m-t-15" v-if="currentMaterialType == 'video'">
                <div class="select-video-box">
                  <img
                    @click="openVideo(taskInfo.content.video.url)"
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
                    :src="taskInfo.content.video.url"
                  ></video>
                </div>
              </div>
              <!-- 小程序 -->
              <div class="flex m-t-15" v-if="currentMaterialType == 'applet'">
                <div class="select-applet-box">
                  <div class="select-applet-name">
                    <img
                      class="mini-img"
                      src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/mini-applet-icon.png"
                      alt=""
                    /><span>{{ taskInfo.content.applet.app_name }}</span>
                  </div>
                  <div class="select-applet-title">
                    {{ taskInfo.content.applet.title }}
                  </div>
                  <img
                    class="select-applet-img"
                    :src="taskInfo.content.applet.cover_url"
                    alt=""
                  />
                </div>
              </div>
            </div>
          </cy-form-item>
        </cy-form>
        <!-- 任务完成情况 -->
        <div v-if="status == 3" class="risk-title flex">
          <div class="flex flex-align-center">
            任务完成情况
            <el-popover
              placement="top-start"
              width="212"
              trigger="hover"
              content="在智慧导购小程序的社群任务点击任意一个按钮即视为任务完成"
            >
              <img
                slot="reference"
                class="risk-status-icon"
                src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/question-icon.png"
                alt=""
              />
            </el-popover>
          </div>
          <div class="m-r-20">
            ：{{ taskInfo.complete_amount }}/{{ taskInfo.all_amount }}人
          </div>
          <div class="copy-text" @click.stop="downloadExeclFn()">
            下载未完成人员名单
          </div>
        </div>
      </template>
      <template slot="footer">
        <div v-if="status != 3">
          <el-button class="range-choose-btn" @click="save">保存</el-button>
          <el-button @click="cancel">取消</el-button>
        </div>
        <div v-if="status == 3">
          <el-button class="range-choose-btn" @click="cancel">返回</el-button>
          <el-button @click="edited" v-if="taskStatus == 1">编辑</el-button>
        </div>
      </template>
    </layout-form>
    <!-- 可见范围弹窗  -->
    <select-member
      v-if="showChoiceStaff"
      :visible="showChoiceStaff"
      :reback_type="true"
      :tabName="rangeTabName == 0 ? (USER.version == 3 ? 3 : 1) : rangeTabName"
      :selectArr="selectArr"
      :is_show_guide="true"
      :worker_status="1"
      choice_name="已选导购"
      dialog_title="选择可见范围"
      @close="showChoiceStaff = false"
      @submit="submitDialog"
    ></select-member>
    <!-- 客户标签弹窗  -->
    <select-tag-dialog
      v-if="tagDialog"
      :updateSelectedTags="selectedTags"
      @selectedTags="selectedTagsFn"
    ></select-tag-dialog>
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
    <!-- 选择视频素材弹窗-->
    <select-video-material
      v-if="showMaterialModal == 'video'"
      @submit="getSelectMaterial"
    ></select-video-material>
    <!-- 图片视频弹框 -->
    <preview
      v-if="preview_show"
      :urls="checkImgArr"
      @close="preview_show = false"
    ></preview>
    <!-- 预览视频弹框 -->
    <open-video
      v-if="isOpenVideo"
      :video_url="videoUrl"
      @closeModel="cancelVideo"
    ></open-video>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { dateFormat } from "../../../../customer/auto-market/util/index";
import LayoutForm from "@/components/public/layout-form/index.vue";
import UserModule from "@/store/modules/user";
// 可见范围
import SelectMember from "@/components/sy-components/select-member.vue";
// 选中标签
import selectTagDialog from "../../../../enterprise/components/select-tag-dialog.vue";
// 选择表情弹窗 ../components/add-emoji.vue
import AddEmoji from "@/views/myapp/wisdom/marketing-calendar/components/add-emoji.vue";
// 视频预览
import OpenVideo from "@/views/myapp/wisdom/marketing-calendar/components/com-open-video.vue";
// 图片预览
import preview from "@/components/business/preview/index.vue";
// 选择图片素材弹窗
import selectPictureMaterial from "../../../../enterprise/material-library/components/pic-mterial-dialog.vue";
// 选择视频素材弹窗
import selectVideoMaterial from "../../../../enterprise/material-library/components/video-material-dialog.vue";
// 选择小程序素材弹窗
import selectAppletMaterial from "../../../../enterprise/material-library/components/small-mterial-dialog.vue";
import {
  TaskCenterCreateParams,
  taskCenterCreate,
  taskCenterUpdate,
  taskCenterDetail,
  marketingPlanExport
} from "@/api/wisdom/taskCenter";
@Component({
  components: {
    LayoutForm,
    SelectMember,
    selectTagDialog,
    AddEmoji,
    OpenVideo,
    preview,
    selectPictureMaterial,
    selectVideoMaterial,
    selectAppletMaterial
  }
})
export default class extends Vue {
  public get tagsToTableData(): any {
    return [
      {
        list: this.showTagsList
      }
    ];
  }
  private taskTypeList: Array<any> = [
    {
      label: "培训",
      value: 1
    },
    {
      label: "邀约",
      value: 2
    },
    {
      label: "群发",
      value: 3
    }
  ];
  private status = 0; // 页面状态 0-新增 1-编辑 2-复制 3-查看
  private taskStatus = 0; //任务状态1-未开始 2-进行中 3-已结束
  private loading = true;
  private taskInfo: any = {
    task_title: "", // 任务名称
    range: [], // 可见范围
    rangeType: "", // 可见范围类型
    taskTime: [], // 任务时间
    type: "", // 任务类型
    taskType: "", // 任务类型文本
    customerLabel: [], // 客户标签
    complete_amount: 0, // 完成人数
    all_amount: 0, // 总人数
    sendReminder: {
      // 发送提醒
      selectStart: false,
      startTime: "",
      selectEnd: false,
      endTime: "",
      selectfixed: false,
      fixedTime: ""
    },
    text: "", // 任务内容
    content: {
      text: "",
      video: {},
      applet: {},
      picture: []
    } //素材内容
  };
  // 时间限制
  private pickerOption = {
    disabledDate(time: any) {
      // TODO减少科学计算
      return time.getTime() < Date.now() - 8.64e7;
    }
  };
  private validateTaskName = (rule: any, value: any, callback: Function) => {
    if (value == "") {
      callback(new Error("请输入任务名称"));
    }
    if (value.length > 20) {
      callback(new Error("任务名称不能超过20个字"));
    } else {
      callback();
    }
  };
  private taskRules = {
    task_title: [
      { validator: this.validateTaskName, trigger: ["blur", "change"] }
    ],
    text: [{ message: "请输入营销话术内容", trigger: ["blur", "change"] }]
  };
  private showChoiceStaff = false; // 选择可见范围弹窗展示/隐藏
  private rangeTabName = 0; // 0 选中默认
  private selectArr: Array<any> = []; // 可见范围选中的数据

  private tagDialog = false; // 选择客户标签
  private selectedTags: Array<any> = [];
  private showTagsList: Array<string> = [];

  private currentMaterialType = ""; // 当前选择素材的类型
  private showMaterialModal = ""; // 控制当前选择素材的弹框显隐
  private isOpenVideo = false; //是否打开视频
  private videoUrl = ""; // 当前播放视频的url
  private materialSelectList: Array<any> = [];
  private materialSelectListOne: Array<any> = [
    {
      type: "picture",
      hoverIcon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-img-icon2.png",
      icon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-img-hover-icon2.png",
      tip: "图片"
    },
    {
      type: "video",
      hoverIcon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-video-icon2.png",
      icon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-video-hover-icon2.png",
      tip: "视频"
    },
    {
      type: "applet",
      hoverIcon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-mini-icon2.png",
      icon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-mini-hover-icon2.png",
      tip: "小程序"
    }
  ];
  private materialSelectListTwo: Array<any> = [
    {
      type: "applet",
      hoverIcon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-mini-icon2.png",
      icon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-mini-hover-icon2.png",
      tip: "小程序"
    }
  ];
  private preview_show = false;
  private checkImgArr: Array<any> = [];
  // 计算属性
  get USER() {
    return UserModule.getUserInfo;
  }
  /* 生命周期 */
  private mounted() {
    const status = this.$route.query.status;
    const taskStatus = this.$route.query.taskStatus;
    this.status = Number(status);
    this.taskStatus = Number(taskStatus);
    this.materialSelectList = this.materialSelectListOne;
    this.init();
  }

  // 初始化页面
  private init() {
    // 0-新增 1-编辑 2-复制 3-查看
    if (this.status == 0) {
      this.loading = false;
      return;
    }
    // 查看任务详情
    this.getTaskCenterDetail();
  }
  // 输入框禁止输入空格
  private trimLR() {
    this.taskInfo.task_title = this.taskInfo.task_title.replace(/[ ]/g, "");
  }

  // 获取任务详情
  getTaskCenterDetail() {
    const params = {
      task_id: String(this.$route.query.val)
    };
    taskCenterDetail(params).then((res: any) => {
      if (res.code === 200) {
        const resData = res.data;
        this.taskInfo.task_title =
          this.status && this.status == 2
            ? resData.task_title + "-副本"
            : resData.task_title; // 计划名称

        this.taskInfo.complete_amount = resData.complete_amount;
        this.taskInfo.all_amount = resData.all_amount;

        // 可见范围类型
        this.taskInfo.rangeType = resData.object_type === 2 ? "部门" : "门店"; // 接口)门店是1 部门是2
        this.taskInfo.range = resData.object_array; // 可见范围
        this.taskInfo.rangeSumNum = resData.object_array.length;
        const temSelectArr = resData.object_array.map((item: any) => {
          return {
            name: item.name,
            worker_no: item.worker_no,
            id: item.worker_id,
            mem_uid: item.mem_uid
          };
        });
        this.$set(this.$data, "selectArr", temSelectArr);
        this.rangeTabName = resData.object_type === 2 ? 3 : 1; // 3 部门
        this.taskInfo.taskTime = [
          resData.start_time * 1000,
          resData.end_time * 1000
        ];
        this.taskInfo.type = resData.type;
        if (Number(this.taskInfo.type) === 3) {
          this.materialSelectList = this.materialSelectListOne;
        } else {
          this.materialSelectList = this.materialSelectListTwo;
        }
        this.taskInfo.taskType =
          resData.type == 1 ? "培训" : resData.type == 2 ? "邀约" : "群发";

        if (resData.before_start_time != 0) {
          this.taskInfo.sendReminder.startTime = resData.before_start_time;
          this.taskInfo.sendReminder.selectStart = true;
        }
        if (resData.before_end_time != 0) {
          this.taskInfo.sendReminder.endTime = resData.before_end_time;
          this.taskInfo.sendReminder.selectEnd = true;
        }
        if (resData.fixed_time != 0) {
          if (this.status == 3) {
            this.taskInfo.sendReminder.fixedTime = resData.fixed_time;
          } else {
            this.taskInfo.sendReminder.fixedTime = resData.fixed_time * 1000;
          }
          this.taskInfo.sendReminder.selectfixed = true;
        }
        // 客户标签
        if (resData.tags && resData.tags.length > 0) {
          resData.tags.forEach((e: any) => {
            this.showTagsList.push(e.name);
            this.selectedTags.push(`${e.group_id},${e.id}`);
          });
        }

        //任务内容
        this.taskInfo.text = resData.content;
        if (resData.resource && resData.resource.length > 0) {
          if (resData.resource[0].type == "1") {
            // 图片
            this.currentMaterialType = "picture";
            this.taskInfo.content.picture = resData.resource;
            this.taskInfo.content.picture.forEach((item: any) => {
              item.url = item.resource_url;
            });
          } else if (resData.resource[0].type == "2") {
            // 小程序海报
            this.currentMaterialType = "applet";
            this.taskInfo.content.applet = resData.resource[0];
            this.taskInfo.content.applet.cover_url =
              resData.resource[0].resource_url;
          } else if (resData.resource[0].type == "3") {
            this.currentMaterialType = "video";
            this.taskInfo.content.video = resData.resource[0];
            this.taskInfo.content.video.url = resData.resource[0].resource_url;
          }
        } else {
          this.currentMaterialType = "";
        }
        this.loading = false;
      }
    });
  }

  // 可见范围
  private submitDialog(val: any) {
    if (!val.list || !val.type) return;
    this.$set(this.$data, "selectArr", val.list);

    this.rangeTabName = val.type;
    // 获取可见范围
    this.taskInfo.range = val.list.map((item: any) => {
      return {
        // object_id: item.id, // 门店/部门id
        worker_id: item.id, // 门店/部门id,
        mem_uid: item.mem_uid,
        name: item.name // 门店/部门名称
      };
    });
    this.taskInfo.rangeType = val.type === 3 ? "部门" : "门店";
  }

  changeInputStartTime() {
    const pattern = /^[1-9][0-9]*$/; // 正整数的正则表达式
    // 不符合正整数时
    if (!pattern.test(this.taskInfo.sendReminder.startTime)) {
      // input 框绑定的内容为空
      this.taskInfo.sendReminder.startTime = "";
    }
  }

  changeInputEndTime() {
    const pattern = /^[1-9][0-9]*$/; // 正整数的正则表达式
    // 不符合正整数时
    if (!pattern.test(this.taskInfo.sendReminder.endTime)) {
      // input 框绑定的内容为空
      this.taskInfo.sendReminder.endTime = "";
    }
  }

  selectedTagsFn(list: any) {
    this.tagDialog = false;
    if (list) {
      this.selectedTags = list.ids;
      this.showTagsList = list.showList;
    }
  }

  // 添加表情
  async addEmoji(emj: any) {
    const myField = (this.$refs.textarea as any).$refs.textarea;
    let strTem: string = emj || "";
    const surplus: number = 500 - this.taskInfo.text.length;
    if (surplus - strTem.length < 0) {
      strTem = strTem.substring(0, surplus);
      this.$message.warning("消息内容长度不能超过50字符，超出部分则被截断。");
      return;
    }
    if (myField.selectionStart || myField.selectionStart === 0) {
      const startPos = myField.selectionStart;
      const endPos = myField.selectionEnd;
      this.taskInfo.text =
        myField.value.substring(0, startPos) +
        strTem +
        myField.value.substring(endPos, myField.value.length);
      await this.$nextTick();
      myField.focus();
      myField.setSelectionRange(endPos + strTem.length, endPos + strTem.length);
    } else {
      this.taskInfo.text += strTem;
    }
    //获取元素id触发click事件
    (document.getElementById("insertEmoji") as any).click();
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
  private deleteMaterial(type: string, index: number) {
    switch (type) {
      case "picture":
        this.taskInfo.content.picture.splice(index, 1);
        if (this.taskInfo.content.picture.length <= 0)
          this.currentMaterialType = "";
        break;
      case "video":
      case "link":
      case "applet":
        this.taskInfo.content[type] = {};
        this.currentMaterialType = "";
        break;
    }
  }

  // 选择素材弹窗后回调
  private getSelectMaterial(material: any) {
    // 选择图片点击取消按钮回调处理
    if (
      this.currentMaterialType == "picture" &&
      (!material || material.length == 0) &&
      this.taskInfo.content.picture.length > 0
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
        if (material.length + this.taskInfo.content.picture.length > 9) {
          this.$message.warning("图片素材最多上传9张");
          return;
        }
        this.taskInfo.content.picture.push(...material);
        this.showMaterialModal = "";
        break;
      case "video":
      case "link":
      case "applet":
        this.taskInfo.content[this.showMaterialModal] = material;
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

  // 监听任务类型改变
  typeChange() {
    const typeId = this.taskInfo.type;
    this.currentMaterialType = "";
    this.taskInfo.content = {
      text: "",
      video: {},
      applet: {},
      picture: []
    };
    if (Number(typeId) === 3) {
      this.materialSelectList = this.materialSelectListOne;
    } else {
      this.materialSelectList = this.materialSelectListTwo;
    }
  }

  // 查看图片
  private checkNowImage(item: string) {
    this.checkImgArr = [item];
    this.preview_show = true;
  }

  // 保存
  save() {
    if (this.taskInfo.task_title && this.taskInfo.range.length <= 0) {
      this.$message.error("请先选择可见范围");
      // 定位到可选范围
      (this.$refs["taskNameInfo"] as any).$el.scrollIntoViewIfNeeded();
      return;
    }
    if (this.taskInfo.type == "3" && this.selectedTags.length == 0) {
      this.$message.error("请先选择客户标签");
      // 定位到可选范围
      (this.$refs["taskNameInfo"] as any).$el.scrollIntoViewIfNeeded();
      return;
    }
    // 计划校验
    (this.$refs.cform as any).validate((valid: boolean) => {
      if (valid) {
        // 任务开始时间校验
        if (this.taskInfo.taskTime[0] < new Date().getTime()) {
          this.$message.error("任务开始时间应大于当前时间");
          return;
        }
        // 发送提醒存在的话，对其验证
        if (
          this.taskInfo.sendReminder.selectStart &&
          this.taskInfo.sendReminder.startTime
        ) {
          const timeDifference =
            (this.taskInfo.taskTime[0] - new Date().getTime()) / (1e3 * 60);
          if (this.taskInfo.sendReminder.startTime > timeDifference) {
            this.$message.error("开始前的提醒时间应在当前时间之后");
            return;
          }
        }
        if (
          this.taskInfo.sendReminder.selectEnd &&
          this.taskInfo.sendReminder.endTime
        ) {
          const timeDifference =
            (this.taskInfo.taskTime[1] - this.taskInfo.taskTime[0]) /
            (1e3 * 60);
          if (
            this.taskInfo.sendReminder.endTime >= Math.floor(timeDifference)
          ) {
            this.$message.error("结束前的提醒时间应处于任务期内");
            return;
          }
        }
        if (
          this.taskInfo.sendReminder.selectfixed &&
          this.taskInfo.sendReminder.fixedTime
        ) {
          const fixedTime =
            typeof this.taskInfo.sendReminder.fixedTime == "number"
              ? this.taskInfo.sendReminder.fixedTime
              : this.taskInfo.sendReminder.fixedTime.getTime();
          if (fixedTime < new Date().getTime()) {
            this.$message.error("固定提醒时间应晚于当前时间");
            return;
          }
          if (fixedTime > this.taskInfo.taskTime[1]) {
            this.$message.error("固定提醒时间应早于任务结束时间");
            return;
          }
        }

        // 处理任务内容入参
        const temTaskArray: Array<any> = [];
        const taskInfo = this.taskInfo;
        temTaskArray.push({
          task_title: taskInfo.task_title,
          start_time:
            taskInfo.taskTime.length > 0
              ? typeof taskInfo.taskTime[0] == "number"
                ? taskInfo.taskTime[0] / 1e3
                : Date.parse(taskInfo.taskTime[0]) / 1e3
              : "",
          end_time:
            taskInfo.taskTime.length > 0
              ? typeof taskInfo.taskTime[1] == "number"
                ? taskInfo.taskTime[1] / 1e3
                : Date.parse(taskInfo.taskTime[1]) / 1e3
              : "",
          type: taskInfo.type,
          content: taskInfo.text
        });

        // 任务类型为群发
        if (taskInfo.type == "3") {
          temTaskArray[0].tag_ids = this.selectedTags.map((v: any) => {
            const labelArr: Array<any> = v.split(",");
            const req: any = Number(labelArr[1]);
            return req;
          });
        }

        // 是否勾选发送提醒的某一项
        if (taskInfo.sendReminder.selectStart) {
          temTaskArray[0].before_start_time =
            taskInfo.sendReminder.startTime || 0;
        } else {
          temTaskArray[0].before_start_time = 0;
        }
        if (taskInfo.sendReminder.selectEnd) {
          temTaskArray[0].before_end_time = taskInfo.sendReminder.endTime || 0;
        } else {
          temTaskArray[0].before_end_time = 0;
        }
        if (taskInfo.sendReminder.selectfixed) {
          temTaskArray[0].fixed_time =
            taskInfo.sendReminder.fixedTime > 0
              ? typeof taskInfo.sendReminder.fixedTime == "number"
                ? taskInfo.sendReminder.fixedTime / 1e3
                : Date.parse(taskInfo.sendReminder.fixedTime) / 1e3
              : 0;
        } else {
          temTaskArray[0].fixed_time = 0;
        }

        // 有任务资源
        if (this.currentMaterialType) {
          const temResource: Array<any> = [];
          if (this.currentMaterialType == "picture") {
            taskInfo.content.picture.forEach((picItem: any) => {
              picItem.resource_url = picItem.url;
              picItem.type = 1;
              temResource.push(picItem);
            });
          } else if (this.currentMaterialType == "video") {
            temResource.push(taskInfo.content.video);
            temResource[0].type = 3;
            temResource[0].resource_url = temResource[0].url;
          } else if (this.currentMaterialType == "applet") {
            temResource.push(taskInfo.content.applet);
            temResource[0].type = 2;
            temResource[0].resource_url = temResource[0].cover_url;
          }
          temTaskArray[0].resource = temResource;
        } else {
          temTaskArray[0].resource = [];
        }
        const params: TaskCenterCreateParams = {
          object_type: this.taskInfo.rangeType == "门店" ? 1 : 2, // 可见范围类型 (接口)门店是1 部门是2
          object_array: this.taskInfo.range, // 可见范围数据
          task_array: temTaskArray // 任务内容
        };
        // 编辑状态
        if (this.status == 1) {
          this.loading = true;
          params.task_id = String(this.$route.query.val);
          this.updateSave(params);
        } else {
          this.loading = true;
          this.createSave(params);
        }
      }
    });
  }

  // 新增任务保存
  private createSave(params: TaskCenterCreateParams) {
    taskCenterCreate(params)
      .then((res: any) => {
        this.loading = false;
        if (res.code == 200) {
          this.$message.success("保存成功");
          this.$router.go(-1);
        }
      })
      .catch((err: any) => {
        this.loading = false;
        console.log("新增任务失败", err);
      });
  }

  // 下载未完成人员名称
  async downloadExeclFn() {
    try {
      const data: any = {
        task_id: String(this.$route.query.val)
      };
      const res: any = await marketingPlanExport(data);
      const blob: Blob = new Blob([res], {
        type: "application/x-xls;charset=utf-8"
      });
      const a: HTMLElement = document.createElement("a");
      const href = window.URL.createObjectURL(blob);
      const name = `未完成人员名单${dateFormat(
        Date.now() / 1000,
        "yyyy_M_d"
      )}.xlsx`;
      a.setAttribute("download", name);
      a.setAttribute("href", href);
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      window.URL.revokeObjectURL(href);
    } catch (error) {
      console.warn("导出失败！");
    }
  }

  // 编辑任务保存
  private updateSave(params: TaskCenterCreateParams) {
    taskCenterUpdate(params)
      .then((res: any) => {
        this.loading = false;
        if (res.code == 200) {
          this.$message.success("保存成功");
          this.$router.go(-1);
        }
      })
      .catch((err: any) => {
        this.loading = false;
        console.log("更新任务失败", err);
      });
  }

  // 编辑
  edited() {
    this.$router.push({
      name: "myapp_wisdom_task-center_community-task_task-edit",
      query: {
        status: "1",
        val: this.$route.query.val
      }
    });
  }

  // 取消
  cancel() {
    this.$confirm("你确认要离开吗？", "取消确认", {
      confirmButtonText: "确认",
      cancelButtonText: "取消"
    }).then(() => {
      this.$router.go(-1);
    });
  }
}
</script>
<style lang="scss" scoped>
.com-task-center {
  background: #fff;
  .in-item {
    width: 100%;
  }
  .range-choose-btn {
    color: #fff;
    background: $color-main;
  }
  .range-content {
    padding: 0 20px 5px;
    justify-content: space-between;
    margin-top: 15px;
    border-bottom: 1px solid #ddd;
    .left {
      display: flex;
      flex-wrap: wrap;
      max-height: 132px;
      overflow: auto;
      .item {
        background: #f0f3fe;
        border-radius: 3px;
        border: 1px solid #d6dffd;
        padding: 8px 10px;
        margin: 0 10px 10px 0;
        img {
          width: 14px;
          height: 14px;
          margin-right: 4px;
        }
      }
      .item-name {
        font-size: 12px;
        line-height: 12px;
      }
    }
    .right {
      color: #6f3112;
      font-size: 14px;
      min-width: 60px;
      text-align: right;
    }
  }
  /deep/ .el-collapse-item__header {
    background: #ebebf1;
    height: 36px;
  }
  /deep/ .el-collapse-item__arrow {
    margin-right: 20px;
  }
  /deep/ .el-collapse {
    border-top: none;
  }
  /deep/ .el-collapse-item__wrap {
    border-bottom: none;
  }
  /deep/ .el-collapse-item__content {
    padding-bottom: 0;
  }
}
.risk-content-box {
  width: 600px;
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
      height: 200px !important;
      padding: 15px 0;
      font-family: PingFangSC-Medium, PingFang SC;
      font-size: 14px;
      line-height: 20px;
      color: #606266;
      margin-bottom: 40px;
    }
    /deep/.el-input__count {
      right: 0;
    }
  }
  .material-select-box {
    padding: 15px;
    margin-top: 12px;
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
}
.select-material-list {
  width: auto;
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
.font-re {
  font-family: PingFangSC-Regular, PingFang SC;
  font-weight: 400;
}
.font-me {
  font-family: PingFangSC-Medium, PingFang SC;
  font-weight: 500;
  font-size: 14px;
}
.flex {
  display: flex;
}
.flex-align-center {
  align-items: center;
}
.flex-wrap {
  flex-wrap: wrap;
}
.title-item {
  margin-top: 10px;
  box-sizing: border-box;
  line-height: 20px;
  height: 36px;
  background: #ebebf1;
  padding: 8px 20px;
  justify-content: space-between;
}
.pointer-cursor {
  cursor: pointer;
}
.sendReminderBox {
  margin-bottom: 10px;
  .el-checkbox {
    margin-right: 12px;
  }
  .el-input {
    width: 60px;
    margin: 0 12px;
  }
  /deep/.el-input__inner {
    height: 28px;
  }
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
.check-risk-content {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  padding: 15px;
  .material-img-show {
    width: 80px;
    height: 80px;
    margin: 15px 15px 0 0;
    border-radius: 4px;
    overflow: hidden;
    img {
      width: 100%;
      height: 100%;
    }
  }
}
/deep/ .el-checkbox__inner {
  width: 16px;
  height: 16px;
  border-radius: 50%;
}
.mlr-12 {
  margin: 0 12px;
}
.risk-title {
  height: 40px;
  line-height: 40px;
  width: 340px;
  padding: 0 10px 0 20px;
  background: rgba(187, 187, 187, 0.2);
  .risk-status-icon {
    width: 14px;
    height: 14px;
    margin: 0 5px;
  }
  .copy-text {
    color: #708df5;
    cursor: pointer;
  }
}
::v-deep {
  .exchange-picker {
    input {
      height: 100%;
      width: 400px;
    }
  }
  .activity-time {
    .el-range__close-icon {
      position: absolute;
      right: 5px;
    }
  }
}
/deep/ .el-dialog {
  min-width: 700px;
}
</style>
