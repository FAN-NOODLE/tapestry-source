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
              <div v-if="info.fileData.length > 0">
                <p v-for="(v, i) in info.fileData" :key="i" class="line">
                  <span class="item-avatar-left"></span>
                  <span class="item-content-right str">
                    <div
                      class="select-picture-box"
                      v-if="v.fileType && v.fileType == 1"
                    >
                      <div class="item-picture">
                        <img class="item-img" :src="v.url" alt="" />
                      </div>
                    </div>

                    <div class="select-link-box" v-if="v.fileType == 2">
                      <div class="select-link-title">
                        {{ v.title }}
                      </div>
                      <div class="select-link-bottom">
                        <!-- 为了解决flex布局下超出省略失效问题，所以嵌套-->
                        <div>
                          <span class="select-link-description">{{
                            v.description
                          }}</span>
                        </div>
                        <img
                          class="select-link-cover"
                          :src="v.cover_url"
                          alt=""
                        />
                      </div>
                    </div>

                    <div class="select-applet-box" v-if="v.fileType == 3">
                      <div class="select-applet-name">
                        <img
                          class="minilogo"
                          src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/mini-applet-icon.png"
                          alt=""
                        />

                        <span class="smallName">{{ v.app_name }}</span>
                      </div>
                      <div class="select-applet-title">
                        {{ v.title }}
                      </div>
                      <img
                        class="select-applet-img"
                        :src="v.cover_url"
                        alt=""
                      />
                    </div>

                    <div class="" v-if="v.fileType == 4">
                      <div class="video-img-s">
                        <video :src="v.url" @click="openVideo(v.url)" />
                        <img
                          @click="openVideo(v.url)"
                          class="video-play"
                          src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/play_icon.png"
                          alt=""
                        />
                      </div>
                    </div>
                    <div class="select-file-box" v-if="v.fileType == 5">
                      <div class="file-type" :class="`file-bg${v.file_type}`">
                        {{
                          v.file_type != 5 ? customIcon[v.file_type - 1] : ""
                        }}
                        <i
                          v-if="v.file_type == 5"
                          class="icon el-icon-document-remove"
                          aria-hidden="true"
                        ></i>
                      </div>
                      <div class="file-title">
                        <el-popover
                          placement="top-start"
                          v-if="v.name.length >= 8"
                          title=""
                          width="200px"
                          trigger="hover"
                          :content="v.name"
                        >
                          <h3 slot="reference">{{ v.name.slice(0, 8) }}...</h3>
                        </el-popover>
                        <h3 v-else>{{ v.name }}</h3>
                      </div>
                      <div class="date-time">
                        <span>{{
                          v.created_at | formatDate("yyyy-mm-dd hh:mm:ss")
                        }}</span>
                      </div>
                    </div>
                  </span>
                </p>
              </div>
            </div>
          </div>
          <div class="content-right">
            <el-form
              class="main-box"
              :class="{ 'main-box-label': info.content.fileData }"
              ref="cform"
              :rules="rules"
              :model="info"
              label-width="90px"
            >
              <!-- 条件范围 -->
              <el-form-item label="条件范围" prop="conditionRange">
                <el-radio v-model="info.conditionRange" label="1"
                  >第三方数据</el-radio
                >
              </el-form-item>
              <!-- 条件设置 -->
              <el-form-item
                v-for="(item, index) in allConditionList"
                :key="item.code"
                :label="getLabelName(item, index)"
                prop="store_type"
                :rules="getLabelName(item, index) === '门店类型' ?{ required: true, message: '请选择门店类型' }:[]"
              >
                <el-select
                  style="width: 400px"
                  v-model="selectCondition[item.code]"
                  :multiple="item.code !== 'reserve_status'"
                  :placeholder="`请选择${item.name}`"
                  @change="handleChange($event, item.code)"
                >
                  <el-option
                    v-for="cell in item.type_data_list"
                    :key="cell.code"
                    :label="cell.name"
                    :value="cell.code"
                  ></el-option>
                </el-select>
              </el-form-item>
              <el-form-item label="消息内容" prop="text" required>
                <!-- :required="info.content.fileData ? false : true" -->
                <div class="publish-content-box">
                  <div class="right-select">
                    <el-table
                      :data="getTextList()"
                      border
                      @row-click="rowClick"
                    >
                      <el-table-column label="可替换文本" width="150">
                        <template slot="header">
                          <span class="m-r-5">可替换文本</span>
                          <el-popover placement="top-start" trigger="hover">
                            <div>
                              使用说明：<br />
                              点击可替换文本，将内容添加到消息内容对话框中，系统在发送欢迎语<br />
                              将对应的文本替换为实际门店/导购/顾客信息。基于定义如下：<br />
                              [顾客昵称] ：顾客的微信昵称<br />
                              [会员姓名]：顾客在CRM留存的真实姓名<br />
                              [门店名称]：添加顾客时，导购所归属的门店名称<br />
                              [导购姓名]：导购企业微信名称<br />
                              <!-- [预约服务名称]：预约服务名称<br />[预约时间]：预约时间 -->
                            </div>
                            <i
                              slot="reference"
                              class="icon el-icon-question"
                            ></i>
                          </el-popover>
                        </template>
                        <template slot-scope="scope">
                          <div class="td-item">
                            {{ scope.row }}
                          </div></template
                        >
                      </el-table-column>
                    </el-table>
                  </div>
                  <div class="publish-content-wrap">
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
                        @blur="handleInputBlur"
                        v-model="info.text"
                        :rows="4"
                        type="textarea"
                        maxlength="500"
                        show-word-limit
                      >
                      </el-input>
                    </div>
                  </div>
                </div>
              </el-form-item>
              <el-form-item
                label="消息素材"
                prop="fileData"
                :required="info.text ? false : true"
                trigger="change"
              >
                <!-- 图片、视频、链接、小程序选择区域-->
                <div class="material-select-box">
                  <div class="m-b-10 select-btn">
                    <el-popover placement="top" width="340" trigger="hover">
                      <div class="select-material-list">
                        <div
                          class="item-select"
                          v-for="(item, index) in materialSelectList"
                          @mouseover="selectStyle(item)"
                          @mouseout="outStyle(item)"
                          @click="openDialog(item.type)"
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
                        type="primary"
                        class="select-material-btn"
                        size="small"
                        ><span>设置消息素材</span></el-button
                      >
                    </el-popover>
                    <p v-if="ep.isCoach" class="tips">
                      支持最多9个不同类型的素材，【本地上传】可支持上传图片及视频素材，最大不超过10M
                    </p>
                    <p v-else class="tips">
                      支持最多9个不同类型的素材
                    </p>
                  </div>
                  <div v-show="info.fileData.length > 0">
                    <div>
                      <span class="m-r-20"
                        >已选<span class="high"
                          >{{ info.fileData.length }}/9</span
                        >个素材</span
                      >

                      <el-popconfirm
                        icon="el-icon-info"
                        confirm-button-text="确认"
                        cancel-button-type="Primary"
                        @onConfirm="info.fileData = []"
                        @confirm="info.fileData = []"
                        title="你确认移除全部素材吗？"
                      >
                        <el-button
                          type="text"
                          size="small"
                          class="del-btn"
                          slot="reference"
                          >全部移除</el-button
                        >
                      </el-popconfirm>
                    </div>
                    <el-table
                      class="droptable"
                      :data="info.fileData"
                      stripe
                      row-key="sid"
                      :row-class-name="
                        info.fileData.length > 1 ? 'pointer' : ''
                      "
                    >
                      <el-table-column label="序号" min-width="130">
                        <template slot-scope="scope">
                          {{ scope.$index + 1 }}
                          <!-- <img
                            v-if="info.fileData.length > 1"
                            class="hover-icon"
                            src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/customer/drag_icon.png"
                            alt=""
                          /> -->
                        </template>
                      </el-table-column>
                      <el-table-column label="素材名称" min-width="220">
                        <template slot-scope="scope">
                          <div class="aligns-center">
                            <el-popover
                              placement="top-start"
                              width="250"
                              trigger="hover"
                              :disabled="scope.row.title.length <= 15"
                            >
                              <div
                                style="overflow:hidden"
                                v-html="
                                  scope.row.title
                                    .replace(/\n/g, '<br/>')
                                    .replace(/\s/g, '&nbsp;')
                                "
                              ></div>
                              <span
                                v-if="scope.row.title.length <= 15"
                                class="m-l-10 goods-name el-popover-ref"
                                slot="reference"
                                >{{ scope.row.title }}</span
                              >
                              <span
                                slot="reference"
                                v-else
                                class="m-l-10 goods-name el-popover-ref"
                                >{{
                                  scope.row.title.substr(0, 15) + "..."
                                }}</span
                              >
                            </el-popover>
                          </div>
                        </template>
                      </el-table-column>
                      <el-table-column label="素材类型" min-width="120">
                        <template slot-scope="{ row }">
                          {{ row.fileType | typeNameFilter }}
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
                            @onConfirm="deleteRowFn(scope.$index)"
                            @confirm="deleteRowFn(scope.$index)"
                            title="你确认删除该素材吗？"
                          >
                            <el-button
                              type="text"
                              size="small"
                              class="del-btn"
                              slot="reference"
                              >移除</el-button
                            >
                          </el-popconfirm>
                        </template>
                      </el-table-column>
                    </el-table>
                  </div>
                </div>
              </el-form-item>
            </el-form>
          </div>
        </div>
      </template>
      <template slot="footer">
        <div class="dialog-footer">
          <el-button type="primary" @click="save()" size="small"
            >保 存</el-button
          >
          <el-button @click="cancel()" size="small">取 消</el-button>
        </div>
      </template>
    </layout-form>
    <!-- 选择卡等 -->
    <select-card
      v-if="cardDialog"
      @channelDialog="channelDialog"
      @submitCard="submitCard"
    >
    </select-card>
    <!-- 选择文本素材弹窗-->
    <select-text-material
      v-if="textDialog"
      @textData="getTextData"
    ></select-text-material>
    <!-- 选择素材弹窗 -->
    <select-material
      v-if="dialog"
      :chooseInfo="info.fileData"
      :file="true"
      :type="infoType"
      @submit="getFileData"
    ></select-material>
    <!-- 预览视频弹框 -->
    <open-video
      v-if="is_open_video"
      :video_url="videoUrl"
      @closeModel="cancelVideo"
    ></open-video>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Watch } from "vue-property-decorator";
import { dateFormat } from "@/views/myapp/customer/auto-market/util/index";
import API from "@/api/enterprise/drainage-coach";
import selectCard from "./selectCardDialog.vue";
import AddEmoji from "../../../components/add-emoji.vue";
import {
  Credential,
  getFileCredentialApi
} from "@/api/enterprise/sensitive-word";
import { ApiResult } from "@/types";
import COSSDK from "cos-js-sdk-v5";
// 选择文本素材弹窗
import selectTextMaterial from "../../../material-library/components/text-material-dialog.vue";
import selectMaterial from "@/components/business/select-material/select-material.vue";
import openVideo from "../../../components/com-open-video.vue";
import Sortable from "sortablejs";
import LocalUpload from '@/mixins/LocalUpload';
import { MATERIAL_SELECT_LIST } from "@/contanst";
interface SelectInterface {
  reserve_status?: string;
  [propName: string]: any;
}
@Component({
  filters: {
    typeNameFilter(val: number) {
      switch (val) {
        case 1:
          return "图片";
        case 2:
          return "链接";
        case 3:
          return "小程序";
        case 4:
          return "视频";
        case 5:
          return "文件";
        default:
          break;
      }
    }
  },
  components: {
    selectCard,
    AddEmoji,
    selectTextMaterial,
    selectMaterial,
    openVideo
  },
  mixins: [LocalUpload],
})
export default class extends Vue {
  private personaliseWelcomeId: any = 0;
  private personaliseWelcomeType: any = 0;
  public get msgTextStr(): string {
    return this.info.text.replace(/\n/g, "<br/>").replace(/\s/g, "&nbsp;");
  }
  private timer: any = undefined;
  private updateTime = dateFormat(
    Number.parseInt(String(Date.now() / 1000)),
    "hh:mm"
  );
  private customIcon = ["W", "X", "P", "PDF"]; //自定义icon
  private info: any = {
    conditionRange: "1", //条件范围
    conditionSelect: "1", //条件选择
    conditionSetting: {}, //条件设置
    card: "", // 选择卡等
    storeType: "", // 门店类型
    customerType: "", // 用户
    content: {
      text: "",
      fileType: 0
    },
    text: "",
    fileData: [],
    store_type: ''
  };
  private selectCondition = {} as SelectInterface;
  private allConditionList = [];
  private materialSelectList: any = JSON.parse(JSON.stringify(MATERIAL_SELECT_LIST));
  private rules = {
    text: [
      {
        required: true,
        message: "请输入欢迎语",
        trigger: "blur"
      }
    ],
    fileData: [
      {
        required: false,
        validator: (rule: any, value: any[], callback: (error?: Error) => void) =>
          this.validateFileData(rule, value, callback),
        message: "请选择素材",
        trigger: "blur"
      }
    ],
    store_type: [
      {
        required: true,
        message: "请选择门店类型",
        trigger: "change"
      }
    ],
  };
  private cardDialog = false;
  private textDialog = false;
  private pictureDialog = false;
  private urlDialog = false;
  private appletDialog = false;
  private infoType = 1;
  private dialog = false;
  private videoUrl = "";
  private is_open_video = false;
  private cos: any = null;
  private credential?: Credential;
  private vipList = [];
  private cursorIndex = 0;
  get ep() {
    return this.$store.getters['enterprise']
  }
  // 初始化腾讯对象存储Cos
  private async initCos() {
    const {
      code,
      message = "",
      data
    }: ApiResult<Credential> = await getFileCredentialApi();
    if (code !== 200) {
      this.$message(message);
      return false;
    }
    this.credential = data;
    this.cos = new COSSDK({
      getAuthorization: (options: any, callback: (credentials: any) => void) => {
        console.log(options);
        callback({
          TmpSecretId: data?.secret_id,
          TmpSecretKey: data?.secret_key,
          XCosSecurityToken: data?.session_token,
          StartTime: data?.start_time,
          ExpiredTime: data?.expired_time
        });
      }
    });

    return true;
  }
  private mounted() {
    this.timer = setInterval(() => {
      this.updateTime = dateFormat(
        Number.parseInt(String(Date.now() / 1000)),
        "hh:mm"
      );
    }, 1000);
    this.personaliseWelcomeType = this.$route.query.type
      ? Number(this.$route.query.type)
      : 0;
    this.personaliseWelcomeId = this.$route.query.id
      ? Number(this.$route.query.id)
      : 0;

    if (this.personaliseWelcomeType === 1) {
      this.getDrainageDetail();
    } else {
      this.getTypeList();
    }
    // 阻止默认行为
    document.body.ondrop = function(event) {
      event.preventDefault();
      event.stopPropagation();
    };
    // this.rowDrop();
    this.materialSelectList = this.ep.isCoach ? MATERIAL_SELECT_LIST : MATERIAL_SELECT_LIST.slice(0, 5)
  }
  private getTextList() {
    let textList = ["顾客昵称", "会员姓名", "门店名称", "导购姓名"];
    if (
      this.selectCondition.reserve_status &&
      this.selectCondition.reserve_status === "yes_reserve"
    ) {
      textList = [...textList, "预约服务名称", "预约时间"];
    }
    return textList;
  }
  private getLabelName(item: any, index: number) {
    let value = "门店类型";
    if (index === 0) value = "条件设置";
    if (item.code === "reserve_status") value = "预约状态";
    return item.name || value;
  }
  handleChange(value, key) {
    console.log(value, key);
    this.info[key] = value;
    console.log(this.info);
    
  }
  private handleInputBlur(e) {
    this.cursorIndex = e.srcElement.selectionStart;
  }
  @Watch("info.text")
  private handleTextInput(value: string) {
    if (value && value.trim()) {
      this.$nextTick(() => {
        (this.$refs["cform"] as any).clearValidate(["fileData"]);
      });
    }
  }
  private validateFileData(rule: any, value: any[], callback: (error?: Error) => void) {
    if (this.info.text && this.info.text.trim()) {
      callback();
      return;
    }
    if (value && value.length > 0) {
      callback();
      return;
    }
    callback(new Error(rule.message));
  }
  private rowClick(e) {
    const value = this.info.text;
    if (this.cursorIndex < this.info.text.length) {
      this.info.text = `${value.substring(
        0,
        this.cursorIndex
      )}{${e}}${value.substring(this.cursorIndex, value.length)}`;
    } else {
      this.info.text += `{${e}}`;
    }
    (this.$refs.textarea as any).focus();
  }
  private async getTypeList(obj?: Record<string, any>) {
    const res: any = await API.getAllTypeListApi({});
    if (res.code === 200) {
      this.allConditionList = res.data || [];
      this.selectCondition = {};
      this.allConditionList.forEach(item => {
        let value: any = [];
        // coach不需要默认值
        if(!this.ep.isCoach) {
          if (item.code === "reserve_status") {
            value = "no_reserve";
          }
        }
        if (obj && obj[item.code]) {
          value = obj[item.code];
        }
        this.$set(this.selectCondition, item.code, value);
      });
    } else {
      this.$message.error(res.message);
    }
  }
  // 预览视频
  private openVideo(url: string) {
    this.videoUrl = url;
    this.is_open_video = true;
  }
  // 关闭预览视频
  private cancelVideo() {
    this.is_open_video = false;
  }
  destroyed() {
    clearInterval(this.timer);
  }

  //行拖拽
  rowDrop() {
    const tbody: any = document.querySelector(
      ".droptable .el-table__body-wrapper tbody"
    );
    const _this = this;
    Sortable.create(tbody, {
      onEnd(res: any) {
        const currRow = _this.info.fileData.splice(res.oldIndex, 1)[0];
        _this.info.fileData.splice(res.newIndex, 0, currRow);
      }
    });
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

  // 打开选择文本弹窗
  private openSelectText() {
    this.textDialog = true;
  }
  private openDialog(type: number) {
    if(type === 6) {
      this.handleUpload()
      return
    }
    this.dialog = true;
    this.infoType = type;
    console.log(type);
    // this[dialog] = true;
    // this.info.fileType = type;
  }
  async handleUpload() {
    const len = this.info.fileData.length
    if(len === 9) {
      this.$message("最多上传9个文件");
      return
    }
    try {
      // 上传多个媒体文件
      const results = await (this as any).uploadMedias({
        multiple: true,
        maxCount: 9 - len,
        maxSize: 10 * 1024 * 1024, // 10MB
      });
      console.log('多文件上传结果:', results);
      // this.info.fileData = results;
      results.forEach((item: any) => {
        this.info.fileData.push({
          url: item.url,
          title: item.name,
          fileType: item.mediaType === 'image' ? 1 : 4,
          name: item.name,
          size: item.size,
          duration: item.duration,
          cover_url: item.cover_url,
        })
      });

    } catch (error) {
      console.error('上传失败:', error);
      // this.$message.error(error.message);
    }
  }
  // 移除单条素材
  private deleteRowFn(i: number) {
    this.info.fileData.splice(i, 1);
  }
  getTextData(data: any) {
    this.textDialog = false;
    if (data) {
      this.addEmoji(this.getDataFromArrayOrObject(data).content);
    }
  }
  getDataFromArrayOrObject(data: any) {
    if (Object.prototype.toString.call(data) == "[object Array]") {
      return data[0];
    } else if (Object.prototype.toString.call(data) == "[object Object]") {
      return data;
    } else {
      return null;
    }
  }
  getPictureData(data: any) {
    this.pictureDialog = false;
    if (data) {
      this.info.content.fileData = this.getDataFromArrayOrObject(data);
    }
  }
  getUrlData(data: any) {
    this.urlDialog = false;
    if (data) {
      this.info.content.fileData = this.getDataFromArrayOrObject(data);
    }
  }
  getAppletData(data: any) {
    this.appletDialog = false;
    if (data) {
      this.info.content.fileData = this.getDataFromArrayOrObject(data);
    }
  }

  getFileData(data: any) {
    console.log(3333, data);
    this.dialog = false;
    if (!data) return;
    this.$nextTick(() => {
      data.length > 0 && (this.$refs["cform"] as any).clearValidate(["text"]);
      data.length > 0 &&
        (this.$refs["cform"] as any).clearValidate(["fileData"]);
    });
    let id = 1;
    data.map((v: any) => {
      v.sid = id++;
    });
    this.info.fileData = JSON.parse(JSON.stringify(data));
    console.log(data);
  }
  deleteMaterial(i: number) {
    this.info.fileData.splice(i, 1);
  }
  // 添加表情
  async addEmoji(emj: any) {
    const myField = (this.$refs.textarea as any).$refs.textarea;
    let strTem: string = emj || "";
    const surplus: number = 500 - this.info.text.length;
    if (surplus - strTem.length < 0) {
      strTem = strTem.substring(0, surplus);
      this.$message.warning("消息内容长度不能超过500字符，超出部分则被截断。");
      return;
    }
    if (myField.selectionStart || myField.selectionStart === 0) {
      const startPos = myField.selectionStart;
      const endPos = myField.selectionEnd;
      this.info.text =
        myField.value.substring(0, startPos) +
        strTem +
        myField.value.substring(endPos, myField.value.length);
      await this.$nextTick();
      myField.focus();
      myField.setSelectionRange(endPos + strTem.length, endPos + strTem.length);
    } else {
      this.info.text += strTem;
    }
  }

  // 选择卡等
  private selectCardEtc() {
    this.cardDialog = true;
  }

  private channelDialog() {
    this.cardDialog = false;
  }

  private submitCard(selectCard: any) {
    this.info.conditionSetting = selectCard;
    this.cardDialog = false;
  }

  // 删除卡等
  private deleteCard() {
    this.info.conditionSetting = {};
  }

  //详情赋值
  analyzingData(data: any) {
    this.getTypeList(data);
    let id = 1;
    data.content_text.forEach((v: any) => {
      if ("text" in v) {
        this.info.text = v.text.content;
      }
      if ("image" in v) {
        this.info.fileData.push({
          url: v.image.pic_url,
          fileType: 1,
          title: v.image.pic_url,
          sid: id++
        });
      } else if ("link" in v) {
        this.info.fileData.push({
          fileType: 2,
          cover_url: v.link.pic_url,
          title: v.link.title,
          description: v.link.desc,
          link_url: v.link.url,
          sid: id++
        });
      } else if ("miniprogram" in v) {
        this.info.fileData.push({
          app_id: v.miniprogram.appid,
          path: v.miniprogram.page,
          title: v.miniprogram.title,
          app_name: v.miniprogram.app_name,
          cover_url: v.miniprogram.cover_url,
          fileType: 3,
          sid: id++
        });
      } else if ("video" in v) {
        this.info.fileData.push({
          id: v.video.media_id,
          name: v.video.name,
          url: v.video.cover_url,
          h_size: v.video.h_size,
          fileType: 4,
          title: v.video.name,
          sid: id++
        });
      } else if ("file" in v) {
        this.info.fileData.push({
          id: v.file.id,
          name: v.file.name,
          url: v.file.cover_url,
          created_at: v.file.created_at,
          file_type: v.file.file_type,
          fileType: 5,
          title: v.file.name,
          sid: id++
        });
      }
    });
  }
  private async getDrainageDetail() {
    const params = {
      id: this.personaliseWelcomeId
    };
    const res: any = await API.drainageDetailApi(params);
    if (res.code === 200) {
      this.info.conditionSetting = {
        id: res.data.welcomes_external_type_data_id,
        name: res.data.type_name
      };

      const { reserve_status } = res.data.condition_content
      this.analyzingData({
        ...res.data,
        ...res.data.condition_content,
        reserve_status:
          reserve_status && reserve_status.length ? reserve_status[0] : ""
      });

      this.info.store_type = res.data.condition_content.store_type;
    } else {
      this.$message.error(res.message);
    }
  }
  private packageData() {
    const data: any = [];
    data.push({
      text: {
        content: this.info.text
      }
    });
    if (this.info.fileData.length > 0) {
      this.info.fileData.forEach((v: any) => {
        switch (v.fileType) {
          case 1:
            data.push({
              image: {
                pic_url: v.url
              }
            });
            break;
          case 2:
            data.push({
              link: {
                pic_url: v.cover_url,
                title: v.title,
                desc: v.description,
                url: v.link_url
              }
            });
            break;
          case 3:
            data.push({
              miniprogram: {
                appid: v.app_id,
                page: v.path,
                title: v.title,
                app_name: v.app_name,
                cover_url: v.cover_url
              }
            });
            break;
          case 4:
            data.push({
              video: {
                name: v.name,
                h_size: v.h_size,
                cover_url: v.url
              }
            });
            break;
          case 5:
            data.push({
              file: {
                created_at: v.created_at,
                name: v.name,
                cover_url: v.url,
                file_type: v.file_type
              }
            });
            break;
          default:
            break;
        }
      });
    }
    const condition:any = {}
    for (const key in this.selectCondition) {
      if (Object.prototype.hasOwnProperty.call(this.selectCondition, key)) {
        const element = this.selectCondition[key];
        if(element.length) {
          condition[key] = element
        }
      }
    }
    if(this.selectCondition.reserve_status) {
      condition.reserve_status = [this.selectCondition.reserve_status]
    } else {
      delete condition.reserve_status
    }
    const params: any = {
      content: data,
      condition_content: {
        ...condition,
      }
    };
    if (this.personaliseWelcomeId) {
      params.id = this.personaliseWelcomeId;
    }
    return params;
  }

  private save() {
    // for (const key in this.selectCondition) {
    //   if (this.selectCondition[key].length === 0) {
    //     const item = this.allConditionList.find(item => item.code === key);
    //     this.$message.warning(`请选择${item.name}`);
    //     return;
    //   }
    // }
    console.log(this.selectCondition);
    
    (this.$refs["cform"] as any).validate(async (valid: any) => {
      if (valid) {
        const data: any = this.packageData();
        const res: any = await API.drainageSaveApi(data);
        if (res.code == 200) {
          this.$message.success("保存成功");
          this.$router.go(-1);
        } else {
          this.$message.error(res.message);
        }
      }
    });
  }

  private cancel() {
    this.$router.go(-1);
  }
}
</script>
<style lang="scss">
.welcome {
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
</style>
<style lang="scss" scoped>
div {
  box-sizing: border-box;
}
.select-material-list {
  width: 180px;
  height: 70px;
  display: flex;
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
.video-img-s {
  width: 145px;
  height: 180px;
  border-radius: 4px;
  position: relative;
  .video-play {
    position: absolute;
    left: 50%;
    top: 50%;
    transform: translate(-50%, -50%);
  }
  video {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}
.select-vedio-box {
  width: 225px;
  height: 83px;
  border: solid 1px #dcdfe6;
  border-radius: 4px;
  position: relative;
  padding: 10px;
  display: flex;
  .delete-vedio-box {
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
  .video-img {
    width: 50px;
    height: 50px;
    border-radius: 4px;
    margin-right: 5px;
    position: relative;
    .video-play {
      position: absolute;
      left: 50%;
      top: 50%;
      transform: translate(-50%, -50%);
    }
    video {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
  }
  .item-info {
    flex: 1;
    width: 152px;
    font-size: 12px;
    font-family: PingFangSC-Regular, PingFang SC;
    font-weight: 400;
    color: #909399;
    .item-title {
      height: 34px;
      line-height: 17px;
      text-overflow: -o-ellipsis-lastline;
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2; //行数需设置
      line-clamp: 2;
      -webkit-box-orient: vertical;
    }
    .item-time {
      margin-top: 5px;
    }
  }
}
.select-file-box {
  width: 155px;
  padding: 0 10px;
  border-radius: 4px;
  border: 1px solid $color-second-line;
  position: relative;
  .file-down {
    width: 100%;
    height: 100%;
    position: absolute;
    top: 0;
    left: 0;
    background: rgba(0, 0, 0, 0.3);
    cursor: pointer;
    display: none;
    .file-down-icon {
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      font-size: 30px;
      color: #ffffff;
    }
  }
  &:hover .file-down {
    display: block;
  }
  .delete-file-box {
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
  .file-title {
    // line-height: 20px;
    height: 20px;
    color: #707277;
    h3 {
      font-size: 14px;
      color: $color-text-main;
      text-align: right;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
      height: 20px;
      line-height: 20px;
    }
  }
  .file-type {
    width: 100%;
    height: 155px;
    line-height: 155px;
    border-radius: 4px;
    margin: 10px 0;
    color: #fff;
    text-align: center;
    font-size: 40px;
    .icon {
      font-size: 40px;
    }
  }
  .date-time {
    text-align: right;
    font-size: 12px;
    color: #707277;
    line-height: 20px;
    margin: 5px 0 10px 0;
  }
  .btns {
    border-top: 1px solid $color-second-line;
    text-align: right;
    margin: 2px 0;
  }
  .file-bg5 {
    background-color: #c9c9c9;
  }
  .del-btn {
    margin-left: 20px;
    color: #f67270;
  }

  .item-select {
    position: absolute;
    right: 0;
    top: 0;
    width: 14px;
    height: 14px;
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
        position: relative;
        padding-right: 180px;
        .right-select {
          position: absolute;
          right: 0;
          top: 0;
          width: 151px;
          /deep/.el-table td {
            padding: 0 0 !important;
          }
          .td-item {
            padding: 13px 0;
            cursor: pointer;
          }
          .not-click {
            cursor: not-allowed;
          }
        }

        .publish-content-wrap {
          border: solid 1px #dcdfe6;
          border-radius: 4px;
          position: relative;
        }
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
            height: 252px !important;
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
        }
      }
    }
  }
}
.video-play {
  width: 25px;
  height: 25px;
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
.tips {
  font-size: 12px;
  font-weight: 400;
  color: #909399;
}
.hover-icon {
  width: 18px;
  height: 18px;
  vertical-align: top;
}
.del-btn {
  font-size: 14px;
  font-weight: 400;
  color: #f67270;
}
.high {
  color: #f67270;
}

::v-deep(.el-form-item__label) {
  text-align: right;
}
</style>
<style>
.pointer {
  cursor: move !important;
}
</style>
