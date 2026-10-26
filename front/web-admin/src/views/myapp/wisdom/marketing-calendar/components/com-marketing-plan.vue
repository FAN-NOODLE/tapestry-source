<template>
  <div class="com-marketing-plan" v-loading.fullscreen.lock="loading">
    <layout-form>
      <template slot="content">
        <cy-form
          ref="cform"
          label-width="90px"
          :rules="planRules"
          :model="planInfo"
        >
          <!-- 计划名称 -->
          <cy-form-item
            label="计划名称"
            ref="planNameInput"
            prop="planName"
            :required="status !== 3"
          >
            <cy-input
              v-if="status !== 3"
              size="small"
              width-size="medium"
              placeholder="请输入计划名称"
              v-model="planInfo.planName"
              class="input-with-select"
              maxlength="20"
              show-word-limit
              @keyup.native="trimLR('planName')"
            ></cy-input>
            <div class="font-re show-content-text" v-if="status == 3">
              {{ planInfo.planName }}
            </div>
          </cy-form-item>
          <!-- 可见范围 -->
          <cy-form-item
            :class="
              planInfo.range && planInfo.range.length > 0
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
            <div v-if="status !== 3" class="des-wrap font-re">
              选择后仅可见范围内的导购，可以在小程序端看到这个营销计划，并接收到该营销计划的“营销日历提醒”
            </div>
            <!-- 标题 -->
            <div
              v-if="
                (planInfo.range && planInfo.range.length > 0) || status == 3
              "
              class="title-item flex"
            >
              <div class="flex">
                <div class="font-me">
                  {{
                    status !== 3 ? "可见范围" : `按${planInfo.rangeType}已选`
                  }}
                </div>
                <div class="font-re" v-if="status == 3">
                  共计{{ planInfo.rangeSumNum }}人
                </div>
              </div>
              <div v-if="status !== 3">操作</div>
            </div>
            <!-- 可见范围显示区 -->
            <div
              v-if="planInfo.range && planInfo.range.length > 0"
              class="range-content flex flex-align-center"
            >
              <div class="left">
                <div
                  v-for="(rangeItem, rangeIndex) in planInfo.range"
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
          <!--营销任务  -->
          <div class="position-relative">
            <div
              class="lable-title"
              :class="status != 3 ? 'lable-title-berfore' : ''"
            >
              营销任务
            </div>
            <div class="risk-content-wrap">
              <el-collapse v-model="activeName">
                <el-collapse-item
                  v-for="(item, index) in planInfo.taskArr"
                  :ref="'riskItem' + index"
                  :key="index"
                  :name="index + 1"
                  :class="{ 'm-b-30': index < planInfo.taskArr.length - 1 }"
                >
                  <template slot="title">
                    <div class="risk-title font-me flex">
                      <div class="flex">
                        <div class="m-r-20">任务{{ index + 1 }}</div>
                        <!-- 任务完成情况 -->
                        <div v-if="status == 3" class="flex">
                          <div class="flex flex-align-center">
                            任务完成情况
                            <el-popover
                              placement="top-start"
                              width="212"
                              trigger="hover"
                              content="在智慧导购小程序营销日历的任务中点击任意一个按钮即视为任务完成"
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
                            ：{{ item.complete_amount }}/{{ item.all_amount }}人
                          </div>
                          <!-- <div
                            class="down-loda-list"
                            @click.stop="downloadExeclFn(item)"
                          >
                            下载未完成人员名单
                          </div> -->
                        </div>
                      </div>
                      <div v-if="status != 3" class="operable-btn flex font-re">
                        <div
                          v-if="planInfo.taskArr && planInfo.taskArr.length > 1"
                          class="delete-text"
                          @click.stop="deleteRiskEvn(index)"
                        >
                          删除
                        </div>
                        <div class="copy-text" @click.stop="copyRiskEvn(item)">
                          复制
                        </div>
                      </div>
                    </div>
                  </template>
                  <!-- 营销任务详情 -->
                  <div>
                    <!-- 任务名称 -->
                    <cy-form-item
                      class="risk-form-item"
                      max="30"
                      :prop="`taskArr.${index}.taskInfo.name`"
                      label="任务名称"
                      :required="status !== 3"
                      trigger="blur"
                    >
                      <cy-input
                        v-if="status != 3"
                        size="small"
                        width-size="medium"
                        placeholder="任务名称"
                        v-model="item.taskInfo.name"
                        maxlength="30"
                        show-word-limit
                        @keyup.native="trimLR(index)"
                      ></cy-input>
                      <div class="font-re show-content-text" v-if="status == 3">
                        {{ item.taskInfo.name }}
                      </div>
                    </cy-form-item>
                    <!-- 任务时间 -->
                    <cy-form-item
                      class="risk-form-item"
                      label="任务时间"
                      trigger="change"
                      :prop="`taskArr.${index}.taskInfo.time`"
                      :required="status !== 3"
                    >
                      <el-date-picker
                        v-if="status !== 3"
                        v-model="item.taskInfo.time"
                        type="daterange"
                        range-separator="-"
                        start-placeholder="开始日期"
                        end-placeholder="结束日期"
                        :clearable="false"
                        :picker-options="pickerOption"
                      >
                      </el-date-picker>
                      <div class="font-re show-content-text" v-if="status == 3">
                        {{ item.taskInfo.time[0] | formatDate("yyyy-mm-dd") }}
                        至
                        {{ item.taskInfo.time[1] | formatDate("yyyy-mm-dd") }}
                      </div>
                    </cy-form-item>
                    <!-- 任务内容 -->
                    <cy-form-item
                      class="risk-form-item"
                      label="任务内容"
                      trigger="blur"
                      :prop="`taskArr.${index}.taskInfo.text`"
                      :required="status !== 3"
                    >
                      <!-- 新增素材展示区 -->
                      <div v-if="status !== 3" class="risk-content-box">
                        <!-- 表情选择区域-->
                        <div class="select-box-insert">
                          <add-emoji
                            @select="addEmoji($event, item, index)"
                            class="add-emoji"
                          >
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
                        </div>
                        <!--多行文本区域-->
                        <div>
                          <el-input
                            :ref="'textarea' + index"
                            class="textarea-box"
                            placeholder="请输入营销话术内容"
                            v-model="item.taskInfo.text"
                            :rows="4"
                            type="textarea"
                            maxlength="500"
                            show-word-limit
                          ></el-input>
                        </div>
                        <!-- 图片、视频、小程序选择区域-->
                        <div class="material-select-box">
                          <!-- 素材综合选择展示区域-->
                          <el-popover
                            placement="top-start"
                            width="118"
                            trigger="hover"
                            v-if="!item.hasMaterial"
                          >
                            <div class="select-material-list">
                              <div
                                class="item-select"
                                v-for="(materialItem,
                                materialIndex) in materialSelectList"
                                @mouseover="selectStyle(materialItem)"
                                @mouseout="outStyle(materialItem)"
                                @click="
                                  openSelectMaterial(materialItem.type, index)
                                "
                                :key="materialIndex"
                              >
                                <div v-if="materialItem.type == 'picture'">
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

                                <!-- <div v-if="materialItem.type == 'applet'">
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
                                </div> -->
                                <div v-if="materialItem.type == 'video'">
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
                              <span>图片/视频</span>
                            </el-button>
                          </el-popover>
                          <!-- 右侧图片选择展示区域-->
                          <div
                            class="select-picture-box"
                            v-if="item.currentMaterialType == 'picture'"
                          >
                            <div
                              class="item-picture"
                              v-for="(pItem, pindex) in item.taskInfo.content
                                .picture"
                              :key="pindex"
                            >
                              <img class="item-img" :src="pItem.url" alt="" />
                              <div
                                class="delete-img-box"
                                @click="
                                  deleteMaterial('picture', index, pindex)
                                "
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
                              v-if="item.taskInfo.content.picture.length < 9"
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
                              item.hasMaterial &&
                                planInfo.taskArr[index].taskInfo.content
                                  .video &&
                                planInfo.taskArr[index].taskInfo.content.video
                                  .url
                            "
                          >
                            <img
                              @click="
                                openVideo(
                                  planInfo.taskArr[index].taskInfo.content.video
                                    .url
                                )
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
                              :src="
                                planInfo.taskArr[index].taskInfo.content.video
                                  .url
                              "
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
                          <div
                            v-if="!item.hasMaterial"
                            class="des-wrap font-re m-t-0"
                          >
                            <div class="m-t-5">
                              图片素材最多上传9张，限制上传格式为png、jpg格式，JPEG格式图片不可上传，大小不超过1M
                            </div>
                            <div class="m-t-5">
                              视频素材最多上传1个，限制上传格式为MP4，大小不超过10M
                            </div>
                            <!-- <div class="m-t-5">
                              小程序海报最多上传1张；不同类型素材不可同时使用
                            </div> -->
                          </div>
                        </div>
                      </div>
                      <div class="error-msg" v-if="showError && status != 3">
                        请输入营销话术内容
                      </div>
                      <!-- 查看素材展示 -->
                      <div
                        class="font-re check-risk-content"
                        v-if="status == 3"
                      >
                        <!-- 任务内容 -->
                        <div
                          class="showRiskContent"
                          :class="
                            item.currentMaterialType
                              ? 'border-bottom p-b-15'
                              : ''
                          "
                        >
                          {{ item.taskInfo.text }}
                        </div>
                        <div>
                          <!-- 图片 -->
                          <div
                            class="flex"
                            v-if="item.currentMaterialType == 'picture'"
                          >
                            <div
                              class="material-img-show"
                              v-for="picTtem in item.taskInfo.content.picture"
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
                          <div
                            v-if="item.currentMaterialType == 'video'"
                            class="m-t-15"
                          >
                            <div class="select-video-box">
                              <img
                                @click="
                                  openVideo(item.taskInfo.content.video.url)
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
                                :src="item.taskInfo.content.video.url"
                              ></video>
                            </div>
                          </div>
                        </div>
                      </div>
                    </cy-form-item>
                  </div>
                </el-collapse-item>
              </el-collapse>
            </div>
            <el-button
              v-if="status != 3"
              class="add-risk-btn range-choose-btn"
              size="small"
              @click="addRiskEvn"
              >添加营销任务</el-button
            >
          </div>
        </cy-form>
      </template>
      <template slot="footer">
        <div v-if="status != 3">
          <el-button class="range-choose-btn" @click="save">保存</el-button>
          <el-button @click="cancel">取消</el-button>
        </div>
        <div v-if="status == 3">
          <el-button class="range-choose-btn" @click="cancel">返回</el-button>
          <el-button @click="edited">编辑</el-button>
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
    <!-- 预览视频弹框 -->
    <open-video
      v-if="isOpenVideo"
      :video_url="videoUrl"
      @closeModel="cancelVideo"
    ></open-video>
    <!-- 小程序带参海报弹窗 -->
    <with-ginseng-poster-pop
      v-if="showPoster"
      :visible="showPoster"
      :appletPosterMess="nowAppletMess"
      @close="closePoster"
      @save="savePoster"
    ></with-ginseng-poster-pop>
    <preview
      v-if="preview_show"
      :urls="checkImgArr"
      @close="preview_show = false"
    ></preview>
    <!-- 编辑时时间过期保存提示 -->
    <el-dialog title="保存提示" :visible.sync="dialogVisible" width="420px">
      <span>{{ `${this.taskName}，是否继续保存` }}</span>
      <span slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="submitEvn">确 定</el-button>
      </span>
    </el-dialog>
    <!-- 选择图片素材弹窗-->
    <select-picture-material
      v-if="showMaterialModal == 'picture'"
      @submit="getSelectMaterial"
    ></select-picture-material>
    <!-- 选择视频素材弹窗-->
    <select-video-material
      v-if="showMaterialModal == 'video'"
      @submit="getSelectMaterial"
    ></select-video-material>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import UserModule from "@/store/modules/user";
// 选择表情弹窗 ../components/add-emoji.vue
import AddEmoji from "@/views/myapp/wisdom/marketing-calendar/components/add-emoji.vue";
// 视频预览
import OpenVideo from "@/views/myapp/wisdom/marketing-calendar/components/com-open-video.vue";
// 小程序带参海报弹窗
import WithGinsengPosterPop from "@/views/myapp/wisdom/marketing-calendar/components/with-ginseng-poster-pop.vue";
// 可见范围
import SelectMember from "@/components/sy-components/select-member.vue";
// 上传图片
import UploadImage from "@/components/business/upload-image/index.vue";
import ShareModule from "@/shared/index";
import preview from "@/components/business/preview/index.vue";
// 选择图片素材弹窗
import selectPictureMaterial from "../../../enterprise/material-library/components/pic-mterial-dialog.vue";
// 选择视频素材弹窗
import selectVideoMaterial from "../../../enterprise/material-library/components/video-material-dialog.vue";

import {
  IMarketingPlanCreateParams,
  marketingPlanCreate,
  marketingPlanUpdate,
  IMarketingPlanExportParams,
  marketingPlanExport,
  IMarketingPlanDetailParams,
  marketingPlanDetail,
  IUploadImageSuccess
} from "@/api/wisdom/marketingPlan";
import { IUploadImageBeforeParam } from "@/api/public/public-interface";
import { VUE_APP_API_ROOT } from "@/utils/env-config";

interface ITaskInfo {
  name?: string;
  time?: string | number;
  text?: string;
  id?: string | number;
  content: {
    videoUrl?: string;
    videoId?: string;
    videoFlag?: boolean;
    applet?: any;
    appletId?: string;
    picture?: Array<any>;
    pictureId?: Array<any>;
    sun_code?: string;
  };
}
// 任务字段类型
interface ITaskArr {
  taskInfo?: ITaskInfo;
  hasMaterial?: boolean;
  currentMaterialType?: string;
  complete_amount?: string | number;
  all_amount?: string | number;
}
// 上传素材按钮信息字段类型
interface IMaterialItem {
  type: string;
  hoverIcon: string | any;
  icon: string | any;
  tip: string;
}
@Component({
  components: {
    AddEmoji,
    OpenVideo,
    WithGinsengPosterPop,
    SelectMember,
    UploadImage,
    preview,
    selectPictureMaterial,
    selectVideoMaterial
  }
})
export default class extends Vue {
  private action = "";
  private loading = true;
  private status = 0; // 页面状态 0-新增 1-编辑 2-复制 3-查看
  private activeName: Array<any> = [1]; // 营销任务初始状态只有一个
  private temVideoUrl = "";
  // 上传素材按钮信息
  private materialSelectList: Array<IMaterialItem> = [
    {
      type: "picture",
      hoverIcon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-img-icon2.png",
      icon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-img-hover-icon2.png",
      tip: "图片"
    },
    // {
    //   type: "applet",
    //   hoverIcon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-mini-icon.png",
    //   icon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-mini-hover-icon.png",
    //   tip: "小程序",
    // },
    {
      type: "video",
      hoverIcon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-video-icon2.png",
      icon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-video-hover-icon2.png",
      tip: "视频"
    }
  ];
  private showError = false;

  // 营销计划
  private planInfo: any = {
    planName: "", // 计划名称
    id: "0", // 计划id(更新有值 默认0)
    rangeType: "", // 可见范围类型
    range: [
      // { worker_id: "", mem_uid: "", name: "" }
    ], // 可见范围
    rangeSumNum: 0,
    // 营销任务
    taskArr: [
      {
        taskInfo: {
          name: "", // 任务名称
          time: [], // 任务时间
          text: "", // 任务内容
          id: "0", //任务id(更新有值 默认0)
          content: {
            videoId: "0",
            videoUrl: "",
            videoFlag: false,
            appletId: "0",
            applet: {
              /*  applet_card_cover_url: "" //小程序封面地址
              applet_card_title: "" //小程序卡片标题
              applet_path: "" //小程序码跳转路径
              resource_id: 19
              resource_url: "https://img-crs-test.vchangyi.com/goods16099878126220.jpg" //海报图
              type: 2 // 1图片 2 小程序海报 3.视频 */
            },
            pictureId: [],
            picture: [],
            video: {}
          }
        },
        hasMaterial: false, // 是否有已选素材
        currentMaterialType: "", // 当前选择素材的类型
        all_amount: "", // 导购总数
        complete_amount: "" // 完成数量
      }
    ]
  };
  // 时间限制
  private pickerOption = {
    disabledDate(time: any) {
      // TODO减少科学计算
      return time.getTime() < Date.now() - 8.64e7;
    }
  };
  private validatePlanName = (rule: any, value: any, callback: Function) => {
    if (value == "") {
      callback(new Error("请输入计划名称"));
    }
    if (value.length > 20) {
      callback(new Error("计划名称不能超过20个字"));
    } else {
      callback();
    }
  };
  private planRules = {
    planName: [
      { validator: this.validatePlanName, trigger: ["blur", "change"] }
    ],
    name: [{ message: "请输入任务名称", trigger: ["blur", "change"] }],
    time: [{ message: "请选择任务时间", trigger: "change" }],
    text: [{ message: "请输入营销话术内容", trigger: ["blur", "change"] }]
  };
  private isOpenVideo = false; //是否打开视频
  private videoUrl = ""; // 当前播放视频的url
  private showPoster = false; // 小程序带参海报弹窗展示flag
  private showChoiceStaff = false; // 选择可见范围弹窗展示/隐藏
  private showMaterialModal = ""; // 控制当前选择素材的弹框显隐
  private imageArr = [];
  private nowRisk = {}; // 当前编辑的任务
  private nowRiskIndex = 0; // 当前任务的下标

  private nowAppletMess = {};
  // 上传图片
  // 2021 2 1
  private newUploadImg = `${VUE_APP_API_ROOT}common/upload-img`;
  private selectArr: Array<any> = []; // 可见范围选中的数据
  private rangeTabName = 0; // 0 选中默认
  private preview_show = false;
  private checkImgArr: Array<any> = [];
  private resourceHasPoster = false; // 任务中包含小程序海报
  private dialogVisible = false;
  private taskName = ""; // 时间过期提示语
  private updataParamse: any; // 编辑更新保存入参
  // 计算属性
  get USER() {
    return UserModule.getUserInfo;
  }
  /* 生命周期 */
  private mounted() {
    this.action = `${this.$VUE_APP_API_VERSION}file/credential`;
    const status = this.$route.query.status;
    this.status = Number(status);
    this.init();
  }

  // 初始化页面
  private init() {
    // 0-新增 1-编辑 2-复制 3-查看
    if (this.status == 0) {
      this.loading = false;
      return;
    }
    this.getMarketingPlanDetail();
  }
  // 输入框禁止输入空格
  private trimLR(index: number | string) {
    if (index == "planName") {
      this.planInfo.planName = this.planInfo.planName.replace(/[ ]/g, "");
    } else {
      this.planInfo.taskArr[index].taskInfo.name = this.planInfo.taskArr[
        index
      ].taskInfo.name.replace(/[ ]/g, "");
    }
  }

  // 可见范围
  private submitDialog(val: any) {
    if (!val.list || !val.type) return;
    this.$set(this.$data, "selectArr", val.list);

    this.rangeTabName = val.type;
    // 获取可见范围
    this.planInfo.range = val.list.map((item: any) => {
      return {
        // object_id: item.id, // 门店/部门id
        worker_id: item.id, // 门店/部门id,
        mem_uid: item.mem_uid,
        name: item.name // 门店/部门名称
      };
    });
    this.planInfo.rangeType = val.type === 3 ? "部门" : "门店";
  }
  // 增加营销任务
  private addRiskEvn() {
    const temRiskItem = {
      taskInfo: {
        name: "", // 任务名称
        time: "", // 任务时间
        text: "",
        id: "",
        content: {
          videoId: "0",
          videoUrl: "",
          videoFlag: false,
          appletId: "0",
          applet: {},
          pictureId: [],
          picture: []
        }
      },
      hasMaterial: false,
      currentMaterialType: "",
      all_amount: "",
      complete_amount: ""
    };
    this.nowAppletMess = {};
    this.planInfo.taskArr.push(temRiskItem);
    this.activeName.push(this.planInfo.taskArr.length);
    // 定位到当前任务
    this.$nextTick(() => {
      (this.$refs as any)[
        "riskItem" + (this.planInfo.taskArr.length - 1)
      ][0].$el.scrollIntoViewIfNeeded();
    });
  }
  // 删除营销日历
  private deleteRiskEvn(index: number) {
    this.$confirm("你确认要删除吗？", "删除确认", {
      confirmButtonText: "确认",
      cancelButtonText: "取消"
    })
      .then(() => {
        this.planInfo.taskArr.splice(index, 1);
        this.$message({
          type: "success",
          message: "删除成功!"
        });
      })
      .catch((err: any) => {
        console.log("删除营销日历", err);
      });
  }
  // 复制营销任务
  private copyRiskEvn(item: object) {
    const newTask = JSON.parse(JSON.stringify(item));
    newTask.taskInfo.name += `${this.planInfo.taskArr.length + 1}`;
    newTask.taskInfo.time = "";
    newTask.taskInfo.id = "0";
    newTask.taskInfo.content.appletId = "0";
    newTask.taskInfo.content.pictureId = [];
    newTask.taskInfo.content.videoId = "0";
    console.log(newTask, "newTask==========");
    this.planInfo.taskArr.push(newTask);
    this.activeName.push(this.planInfo.taskArr.length);
    this.$message.success("复制成功");
    // 定位到当前任务
    this.$nextTick(() => {
      (this.$refs as any)[
        "riskItem" + (this.planInfo.taskArr.length - 1)
      ][0].$el.scrollIntoViewIfNeeded();
    });
  }

  // 鼠标划入事件
  private selectStyle(item: object) {
    const _this = this;
    this.$nextTick(function() {
      this.materialSelectList.forEach(function(item: any) {
        _this.$set(item, "active", false);
      });
      this.$set(item, "active", true);
    });
  }
  // 鼠标划出事件
  private outStyle(item: object) {
    this.$set(item, "active", false);
  }

  // 添加对应素材
  private openSelectMaterial(type: string, index: number) {
    const _this = this;
    this.materialSelectList.forEach(function(item: any) {
      _this.$set(item, "active", false);
    });
    this.nowRiskIndex = index;
    this.planInfo.taskArr[this.nowRiskIndex].currentMaterialType = type;
    this.showMaterialModal = type;
  }
  // 上传图片方法
  private imgBeforeUpload(file: IUploadImageBeforeParam) {
    const fileType = file
      ? file.name.substring(file.name.lastIndexOf(".") + 1)
      : "";
    if (
      file.type != "image/png" &&
      file.type != "image/jpg" &&
      fileType != "png" &&
      fileType != "jpg"
    ) {
      this.$message.error("请上传JPG、PNG格式的图片");
      return false;
    }
    if (file.size > 1 * 1024 * 1024) {
      this.$message.error("请上传1MB大小以内的图片");
      return false;
    }
  }
  // _--------------------------
  private beforeRemove() {
    return this.$confirm("你确认要删除吗？", "删除确认").catch(
      (err: any) => err
    );
  }
  private imgUpdate(res: Array<any>, item: ITaskArr) {
    (item.taskInfo as ITaskInfo).content.picture = JSON.parse(
      JSON.stringify(res)
    );
    // 更新图片素材
    if (res.length == 0) {
      item.currentMaterialType = "";
      item.hasMaterial = false;
    }
  }

  private imageOnSuccessEv(item: ITaskArr) {
    (item as ITaskArr).hasMaterial = true;
  }

  // 查看图片
  private checkNowImage(item: string) {
    this.checkImgArr = [item];
    this.preview_show = true;
  }
  // 点击素材的叉叉按钮
  private deleteMaterial(type: string, index: number, fIndex?: number) {
    switch (type) {
      case "picture":
        this.planInfo.taskArr[index].taskInfo.content.picture.splice(fIndex, 1);
        if (this.planInfo.taskArr[index].taskInfo.content.picture.length <= 0) {
          this.planInfo.taskArr[index].currentMaterialType = "";
          this.planInfo.taskArr[index].hasMaterial = false;
        }
        break;
      case "video":
        this.planInfo.taskArr[index].taskInfo.content.video = {};
        this.planInfo.taskArr[index].currentMaterialType = "";
        this.planInfo.taskArr[index].hasMaterial = false;
        break;
    }
  }

  // 视频

  // 上传文件前
  private videoBeforeUpload(file: any, item: any) {
    const fileMaxSize = 10 * 1024 * 1024; //视频最大上传10M
    const isMP4 = file.type === "video/mp4";
    if (!isMP4) {
      this.$message.error("请上传MP4格式的视频");
      return false;
    }
    if (file.size > fileMaxSize) {
      this.$message.warning("请上传不超过10M以内的视频");
      return false;
    }
    (item.taskInfo as ITaskInfo).content.videoFlag = true;
    item.hasMaterial = true;
  }
  //上传成功
  private videoOnSuccess(file: IUploadImageSuccess) {
    const { data } = file;
    if (this.nowRisk) {
      const temRisk = (this.nowRisk as ITaskArr).taskInfo;
      (temRisk as ITaskInfo).content.videoUrl = data.url;
      this.loading = false;
      (this.nowRisk as ITaskArr).hasMaterial = true;
    }
  }
  private videoChange(res: any, item: any) {
    (item.taskInfo as ITaskInfo).content.videoUrl = JSON.parse(
      JSON.stringify(res)
    );

    // debugger;
    // 更新图片素材
    if (!res) {
      item.currentMaterialType = "";
      item.hasMaterial = false;
      (item.taskInfo as ITaskInfo).content.videoFlag = false;
    }
  }
  // 选择素材弹窗后回调
  private getSelectMaterial(material?: any) {
    console.log("[ material ]", material);
    // 对图片类型处理
    if (material && material.length > 0) {
      material.forEach((item: any) => {
        if (item.type == 2) {
          item.type = 1;
        }
      });
    }
    // 选择图片点击取消按钮回调处理
    if (
      this.planInfo.taskArr[this.nowRiskIndex].currentMaterialType ==
        "picture" &&
      (!material || material.length == 0) &&
      this.planInfo.taskArr[this.nowRiskIndex].taskInfo.content.picture.length >
        0
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
      this.planInfo.taskArr[this.nowRiskIndex].currentMaterialType = "";
      this.showMaterialModal = "";
      return;
    }
    switch (this.showMaterialModal) {
      case "picture":
        if (
          material.length +
            this.planInfo.taskArr[this.nowRiskIndex].taskInfo.content.picture
              .length >
          9
        ) {
          this.$message.warning("图片素材最多上传9张");
          return;
        }
        this.planInfo.taskArr[this.nowRiskIndex].taskInfo.content.picture.push(
          ...material
        );
        this.planInfo.taskArr[this.nowRiskIndex].hasMaterial = true;
        this.showMaterialModal = "";
        break;
      case "video":
        this.planInfo.taskArr[
          this.nowRiskIndex
        ].taskInfo.content.video = material;
        this.planInfo.taskArr[this.nowRiskIndex].hasMaterial = true;
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

  /* 
    添加表情
    emj: 选中的表情
    item: 当前任务详情
    index: 当前第几个任务
   */
  private async addEmoji(emj: any, item: ITaskArr, index: number) {
    const myField = (this.$refs[`textarea${index}`] as any)[0].$refs.textarea;
    // 多行文本内容长度
    if ((this.$refs[`textarea${index}`] as any)[0].value) {
      const textLength = (this.$refs[`textarea${index}`] as any)[0].value
        .length;
      if (textLength >= 499) return;
    }
    if (myField.selectionStart || myField.selectionStart === 0) {
      const startPos = myField.selectionStart;
      const endPos = myField.selectionEnd;
      (item.taskInfo as ITaskInfo).text =
        myField.value.substring(0, startPos) +
        emj +
        myField.value.substring(endPos, myField.value.length);

      this.$set(
        this.planInfo,
        `taskArr[${index}].taskInfo.text`,
        (item.taskInfo as ITaskInfo).text
      );
      await this.$nextTick();
      myField.focus();
      myField.setSelectionRange(endPos + emj.length, endPos + emj.length);
    } else {
      (item.taskInfo as ITaskInfo).text += emj;
    }
  }
  private closePoster() {
    this.showPoster = false;
  }
  // 查看小程序带参海报
  private checkPoster(val: object, item: any) {
    this.showPoster = true;
    this.nowAppletMess = val;
    this.nowRisk = item;
  }
  // 小程序带参海报
  private savePoster(val: any) {
    if (!val.app_name) return ((this.nowRisk as ITaskArr).hasMaterial = false);
    if (this.nowRisk) {
      ((this.nowRisk as ITaskArr).taskInfo as ITaskInfo).content.applet = val;
      (this.nowRisk as ITaskArr).hasMaterial = true;
    }
  }
  private async downloadExeclFn(item: ITaskArr) {
    this.loading = true;
    const data: IMarketingPlanExportParams = {
      task_id: String((item.taskInfo as ITaskInfo).id),
      plan_id: this.planInfo.id
    };
    await marketingPlanExport(data);
    this.loading = false;
    this.$confirm(
      "文件已进入下载列表，请前往【下载中心】查看及下载文件。",
      "文件下载提示",
      {
        distinguishCancelAndClose: true,
        confirmButtonText: "前往下载中心",
        cancelButtonText: "取消",
        type: "success"
      }
    ).then(() => {
      ShareModule.getShared().pushState("myapp_download-center");
    });
  }
  // 计划详情
  private getMarketingPlanDetail() {
    const params: IMarketingPlanDetailParams = {
      plan_id: String(this.$route.query.val)
    };
    marketingPlanDetail(params).then((res: any) => {
      if (res.code === 200) {
        const resData = res.data;
        this.planInfo.planName =
          this.status && this.status == 2
            ? resData.plan_title + "-副本"
            : resData.plan_title; // 计划名称
        this.planInfo.id = resData.plan_id; // 计划id
        // 可见范围类型
        this.planInfo.rangeType = resData.object_type === 2 ? "部门" : "门店"; // 接口)门店是1 部门是2
        this.planInfo.range = resData.object_array; // 可见范围
        this.planInfo.rangeSumNum = resData.object_array.length;
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
        this.planInfo.taskArr = [];
        this.activeName = [];
        // 任务内容
        resData.task_array.forEach((item: any, index: number) => {
          const tempPic: Array<any> = [];
          const pictureId: Array<any> = [];
          let video = {};
          const videoFlag = true;
          let applet = {};
          let appletId = "0";
          let type = "";
          item.resource.forEach((resItem: any) => {
            // 图片
            if (resItem.type == "1") {
              type = "picture";
              resItem.url = resItem.resource_url;
              if (resItem.sun_code && resItem.sun_code.length == 0) {
                resItem.sun_code = "";
              }
              tempPic.push(resItem);
              pictureId.push(resItem.resource_id);
            }
            // 小程序海报
            if (resItem.type == "2") {
              type = "applet";
              applet = resItem;
              appletId = resItem.resource_id;
            }
            // 视频
            if (resItem.type == "3") {
              type = "video";
              resItem.url = resItem.resource_url;
              video = resItem;
            }
          });
          // status 0 新增3-查看
          let temTime = [item.start_time, item.end_time];
          if (this.status !== 0 && this.status !== 3) {
            item.start_time = item.start_time * 1000;
            item.end_time = item.end_time * 1000;
            temTime = [item.start_time, item.end_time];
          }
          // 复制清空时间
          if (this.status == 2) {
            temTime = [];
          }
          const temTask = {
            taskInfo: {
              name: item.task_title, // 任务名称
              time: temTime, // 任务时间
              text: item.content, // 任务内容
              id: item.task_id,
              content: {
                videoFlag,
                video,
                applet,
                appletId,
                picture: tempPic,
                pictureId
              }
            },
            hasMaterial: type ? true : false, // 是否有已选素材
            currentMaterialType: type ? type : "", // 当前选择素材的类型
            complete_amount: item.complete_amount,
            all_amount: item.all_amount
          };
          this.planInfo.taskArr.push(temTask);
          this.activeName.push(index + 1);
        });
        this.loading = false;
      }
    });
  }
  // 保存
  private async save() {
    if (this.planInfo.planName && this.planInfo.range.length <= 0) {
      this.$message.error("请先选择可见范围");
      // 定位到可选范围
      (this.$refs["planNameInput"] as any).$el.scrollIntoViewIfNeeded();
      return;
    }
    // 计划校验
    (this.$refs.cform as any).validate((valid: boolean) => {
      if (!valid) {
        if (!this.planInfo.planName) {
          this.$nextTick(() => {
            (this.$refs["planNameInput"] as any).$el.scrollIntoViewIfNeeded();
          });
          return;
        }
        const temArr = this.planInfo.taskArr;
        for (let i = 0; i < temArr.length; i++) {
          // 定位到当前任务
          if (
            !temArr[i].taskInfo.name ||
            temArr[i].taskInfo.time.length == 0 ||
            !temArr[i].taskInfo.text
          ) {
            (this.$refs as any)["riskItem" + i][0].$el.scrollIntoViewIfNeeded();

            break;
          }
        }
      } else {
        // 处理任务内容入参
        const temTaskArray: Array<any> = [];
        this.planInfo.taskArr.forEach((item: any) => {
          const temResource: Array<any> = [];
          const type =
            item.currentMaterialType == "picture"
              ? "1"
              : item.currentMaterialType == "video"
              ? "3"
              : "";
          // 任务资源
          if (item.hasMaterial) {
            const temContent = item.taskInfo.content;
            const resourceItem = {
              type: type, //1图片  3.视频
              // 资源id(更新有值 默认0)
              resource_id: type == "3" ? temContent.video.resource_id : "",
              resource_url: type == "3" ? temContent.video.url : "",
              applet_path: ""
            };

            // 所选材料为图片
            if (item.currentMaterialType == "picture") {
              item.taskInfo.content.picture.forEach((picItem: any) => {
                picItem.resource_url = picItem.url;
                temResource.push(picItem);
              });
            } else {
              temResource.push(resourceItem);
            }
          }
          const temTime = item.taskInfo.time;
          const riskItem = {
            task_title: item.taskInfo.name,
            task_id: item.taskInfo.id,
            start_time:
              temTime.length > 0
                ? typeof temTime[0] == "number"
                  ? temTime[0] / 1e3
                  : Date.parse(temTime[0]) / 1e3
                : "",
            end_time: temTime
              ? typeof temTime[1] == "number"
                ? temTime[1] / 1e3
                : Date.parse(temTime[1]) / 1e3 + 86399
              : "",
            content: item.taskInfo.text,
            resource: temResource //资源数据
          };
          temTaskArray.push(riskItem);
        });

        const params: IMarketingPlanCreateParams = {
          plan_title: this.planInfo.planName, // 计划名称
          plan_id: this.planInfo.id, // 计划id(更新有值 默认0)
          object_type: this.planInfo.rangeType == "门店" ? 1 : 2, // 可见范围类型 (接口)门店是1 部门是2
          object_array: this.planInfo.range, // 可见范围数据
          task_array: temTaskArray // 任务内容
        };
        // 编辑状态检查任务时间是否已过期，过期提示用户是否继续保存
        if (this.status == 1) {
          const startTimeTask: Array<any> = []; // 开始时间小于当前时间任务
          const expiredTask: Array<any> = []; // 过期任务
          this.planInfo.taskArr.forEach((item: any, index: number) => {
            const itemTime = item.taskInfo.time;
            const nowDate =
              new Date().getFullYear() +
              "-" +
              (new Date().getMonth() + 1) +
              "-" +
              new Date().getDate();
            const nowTime = new Date(nowDate).getTime();
            // 任务开始时间小于当前时间
            if (itemTime[0] < nowTime && itemTime[1] > nowTime) {
              startTimeTask.push(index + 1);
            }
            // 任务时间已过期
            if (itemTime[0] < nowTime && itemTime[1] < nowTime) {
              expiredTask.push(index + 1);
            }
          });

          this.taskName = "";
          let startTimeTaskName = "";
          let expiredTaskName = "";
          if (startTimeTask.length > 0 || expiredTask.length > 0) {
            // 仅任务开始时间早于当前时间
            if (startTimeTask.length > 0) {
              startTimeTask.forEach((item: number, index: number) => {
                startTimeTaskName += `任务${item}${
                  index == startTimeTask.length - 1 ? "" : "、"
                }`;
              });
              this.taskName = `${startTimeTaskName}开始时间早于当前时间`;
            }
            // 仅任务时间过期
            if (expiredTask.length > 0) {
              expiredTask.forEach((item: number, index: number) => {
                expiredTaskName += `任务${item}${
                  index == expiredTask.length - 1 ? "" : "、"
                }`;
              });
              this.taskName = `${expiredTaskName}时间已过期`;
            }

            // 任务开始时间早于当前时间 && 任务时间过期
            if (expiredTask.length > 0 && startTimeTask.length > 0) {
              this.taskName = `${startTimeTaskName}开始时间早于当前时间, ${expiredTaskName}时间已过期`;
            }
            // 任务中包含小程序海报
            if (this.resourceHasPoster) {
              this.dialogVisible = true;
              this.updataParamse = params;
            } else {
              this.$confirm(`${this.taskName}，是否继续保存`, "保存提示")
                .then(() => {
                  this.loading = true;
                  this.updateSave(params);
                })
                .catch((err: any) => {
                  console.log(err);
                });
            }
          }
          // 任务时间正常
          if (startTimeTask.length <= 0 && expiredTask.length <= 0) {
            if (this.resourceHasPoster) {
              this.$alert(
                "小程序素材生成会有延迟，预计数分钟后导购端可以获取素材",
                "保存提示",
                {
                  confirmButtonText: "我知道了",
                  callback: (action: any) => {
                    if (action == "confirm") {
                      this.loading = true;
                      this.updateSave(params);
                    }
                  }
                }
              );
            } else {
              this.loading = true;
              this.updateSave(params);
            }
          }
        } else {
          // 新增任务保存
          if (this.resourceHasPoster) {
            this.$alert(
              "小程序素材生成会有延迟，预计数分钟后导购端可以获取素材",
              "保存提示",
              {
                confirmButtonText: "我知道了",
                callback: (action: any) => {
                  if (action == "confirm") {
                    this.loading = true;
                    this.createSave(params);
                  }
                }
              }
            );
          } else {
            this.loading = true;
            this.createSave(params);
          }
        }
      }
    });
  }
  // 编辑任务时间过期弹窗确定按钮事件
  private submitEvn() {
    this.$alert(
      "小程序素材生成会有延迟，预计数分钟后导购端可以获取素材",
      "保存提示",
      {
        confirmButtonText: "我知道了",
        callback: (action: any) => {
          if (action == "confirm") {
            this.dialogVisible = false;
            this.loading = true;
            this.updateSave(this.updataParamse);
          }
        }
      }
    );
  }
  // 新增营销计划保存
  private createSave(params: IMarketingPlanCreateParams) {
    marketingPlanCreate(params)
      .then((res: any) => {
        this.loading = false;
        if (res.code == 200) {
          this.$message.success("保存成功");
          this.$router.go(-1);
        }
      })
      .catch((err: any) => {
        this.loading = false;
        console.log("新增营销计划 marketingPlanCreate:", err);
      });
  }
  // 编辑更新
  private updateSave(params: IMarketingPlanCreateParams) {
    marketingPlanUpdate(params)
      .then((res: any) => {
        this.loading = false;
        if (res.code == 200) {
          this.$message.success("保存成功");
          this.$router.go(-1);
        }
        if (res.code == 2200003) {
          this.$message.error("该营销计划已被删除");
          this.$router.go(-1);
        }
      })
      .catch((err: any) => {
        this.loading = false;
        console.log("更新营销计划 marketingPlanUpdate:", err);
      });
  }
  // 取消
  private cancel() {
    this.$confirm("你确认要离开吗？", "取消确认", {
      confirmButtonText: "确认",
      cancelButtonText: "取消"
    }).then(() => {
      this.$router.go(-1);
    });
  }
  // 编辑
  private edited() {
    this.$router.push({
      name: "myapp_wisdom_marketing-calendar_edit-marketing-plan",
      query: {
        status: "1",
        val: this.$route.query.val
      }
    });
  }
}
</script>
<style lang="scss" scoped>
.com-marketing-plan {
  background: #fff;
  .in-item {
    width: 100%;
  }
  .range-choose-btn {
    color: #fff;
    background: #708cf6;
  }
  .des-wrap {
    margin: 5px 0 15px;
    color: #909399;
    font-size: 12px;
    line-height: 20px;
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
        color: #708df5;
        font-size: 12px;
        line-height: 12px;
      }
    }
    .right {
      color: #6d86fd;
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
  .plan-form-item {
    margin-bottom: 15px;
  }
  .risk-title {
    width: 100%;
    justify-content: space-between;
    padding: 0 10px 0 20px;
    .operable-btn {
      line-height: 20px;
      align-items: center;
      div {
        margin: 0 10px;
        line-height: 20px;
        font-size: 14px;
      }
      .delete-text {
        color: #f67270;
      }
      .copy-text {
        color: #708df5;
      }
    }
  }
  .down-loda-list {
    color: #708df5;
  }
  .risk-form-item {
    margin: 20px 0;
  }
  .risk-form-item:last-child {
    margin-bottom: 0;
  }
  .add-risk-btn {
    margin-top: 20px;
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
        width: 120px;
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
  border: solid 1px #dcdfe6;
  border-radius: 4px;
  position: relative;
  padding: 10px;
  margin-top: 20px;
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
    white-space: initial;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
  }
  .select-applet-img {
    width: 100%;
    margin-top: 10px;
    border-radius: 4px;
  }
}
.select-material-list {
  // width: 140px;
  height: 70px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
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
.com-marketing-plan > div {
  margin-bottom: 30px;
}
.title-item {
  box-sizing: border-box;
  line-height: 20px;
  height: 36px;
  background: #ebebf1;
  padding: 8px 20px;
  justify-content: space-between;
}
.risk-status-icon {
  width: 14px;
  height: 14px;
  margin: 0 5px;
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
.showRiskContent {
  line-height: 20px;
}
.position-relative {
  position: relative;
}
.lable-title {
  // width: 90px;
  vertical-align: top;
  margin-top: 5px;
  text-align: center;
  font-size: 14px;
  color: #606266;
  max-height: 40px;
  line-height: 20px !important;
  min-width: 20px;
  max-width: 108px;
  position: absolute;
  padding-left: 10px;
  margin-left: 4px;
}
.lable-title-berfore::before {
  content: "*";
  color: #f67170;
  position: absolute;
  left: 2px;
  top: 50%;
  -webkit-transform: translateY(-45%);
  transform: translateY(-45%);
}
.risk-content-wrap {
  width: calc(100% - 90px);
  margin-left: 90px;
}
.pointer-cursor {
  cursor: pointer;
}
.show-material-picture {
  margin-right: 10px;
  margin-bottom: 10px;
  background-size: 100% 100%;
  border-radius: 6px;
  float: left;
  overflow: hidden;
  cursor: pointer;
  width: 100px;
  height: 100px;
  position: relative;
  .mask {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.5);

    .double-btn {
      width: 60px;
      position: absolute;
      top: 0;
      left: 0;
      bottom: 0;
      right: 0;
      margin: auto;
      .el-icon-zoom-in {
        cursor: pointer;
        color: white;
        font-size: 20px;
        position: absolute;
        top: 1px;
        left: 0;
        line-height: 100px !important;
      }
      .has-preview {
        cursor: pointer;
        color: white;
        font-size: 20px;
        position: absolute;
        top: 0;
        right: 0px;
        line-height: 100px !important;
      }
    }
  }
}
// .el-popover {
//   width: 110px !important;
// }
.m-r-3 {
  margin-right: 3px;
}
.border-bottom {
  border-bottom: 1px solid #dcdfe6;
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
.red {
  color: #f67270;
}
.show-content-text {
  box-sizing: border-box;
  height: 36px;
  font-size: 14px;
  color: #303133;
  line-height: 36px;
  padding-top: 3px;
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

.m-l-20 {
  margin-left: 20px;
}
.m-b-30 {
  margin-bottom: 30px;
}
.m-t-15 {
  margin-top: 15px;
}
.m-t-5 {
  margin-top: 5px;
}
.m-t-0 {
  margin-top: 0px;
}
.p-b-15 {
  padding-bottom: 15px;
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
</style>
