<template>
  <div class="friend-circle-edit">
    <layout-form>
      <template slot="content">
        <div class="content-box">
          <cy-form
            class="main-box"
            ref="cform"
            :rules="rules"
            :model="info"
            label-width="120px"
          >
            <cy-form-item
              label="任务名称"
              prop="template_name"
              :rules="rules.template_name"
              required
            >
              <el-input
                v-model.trim="info.template_name"
                clearable
                maxlength="20"
                placeholder="请输入任务名称"
                show-word-limit
              ></el-input>
            </cy-form-item>
            <!-- 任务描述，200字 -->
            <cy-form-item
              label="任务描述"
              prop="template_desc"
              :rules="rules.template_desc"
              required
            >
              <el-input
                ref="textarea"
                class="textarea-box"
                placeholder="请输入任务描述内容"
                v-model.trim="info.template_desc"
                :rows="4"
                type="textarea"
                maxlength="200"
                show-word-limit
              >
              </el-input>
            </cy-form-item>
            <cy-form-item
              v-if="ep.isCoach"
              label="任务类型"
              prop="group_send_type"
              :rules="rules.task_category"
              required
            >
              <el-radio-group v-model="info.task_category" style="width: 608px;" @change="handleCategoryChange">
                <el-radio :label="0">沟通任务</el-radio>
                <el-radio :label="1">朋友圈任务</el-radio>
              </el-radio-group>
            </cy-form-item>
            <cy-form-item
              v-if="!ep.isCoach || info.task_category === 0"
              label="任务群发类型"
              prop="group_send_type"
              :rules="rules.group_send_type"
              required
            >
              <el-radio-group v-model="info.group_send_type" style="width: 608px;">
                <el-radio :label="1">企微一键群发（从群发助手发送）</el-radio>
                <el-radio :label="2">灵活沟通转发（从任务中心发送）</el-radio>
              </el-radio-group>
              <div class="message-info">※每位客户/每个客户群每月最多可接收条数为当月天数，企微一键群发和灵活沟通转发共享额度，超过接收上限的客户/客户群将无法再收到群发消息。（当前为每位客人每天接收1条群发消息。可从企业微信-群发助手-设置群发规则中调整其它使用规则）</div>
            </cy-form-item>
            <cy-form-item
              label="群发顾客"
              required
              v-if="!ep.isCoach || info.task_category === 0 && info.group_send_type === 1"
            >
              <el-radio-group
                v-model="info.allow_select"
                style="width: 608px;"
              >
                <el-radio :label="0">导购不可编辑</el-radio>
                <el-radio :label="1">导购可编辑</el-radio>
              </el-radio-group>
            </cy-form-item>
            <!-- 任务完成标准 -->
            <cy-form-item
              label="任务完成标准"
              required
              :rules="rules.finish_rate"
              v-if="ep.isCoach && info.task_category === 0 && info.group_send_type === 2"
              prop="finish_rate"
            >
              <el-row class="flex align-center">
                <span class="m-r-5">顾客名单触达率达到</span>
                <el-input v-model="info.finish_rate" style="width:100px" size="small"></el-input>
                <span class="m-l-5">%</span>
              </el-row>
            </cy-form-item>
            <cy-form-item
              label="任务内容"
              required
              prop="template_content"
              :rules="rules.template_content"
            >
              <task-content v-model="info.template_content" :task_category="info.task_category" ></task-content>
              <div
                class="error-msg"
                v-if="showError"
              >请填写任务内容</div>
            </cy-form-item>
            <cy-form-item label="配置备选模版">
              <el-switch
                v-model="info.open_back_up"
                :active-value="1"
                :inactive-value="0"
              ></el-switch>
              <span class="tips m-l-5">开启后，可通过配置时间，定时生效备用模板中的内容</span>
            </cy-form-item>
            <div v-if="info.open_back_up === 1" style="padding-left: 120px;">
              <!-- 任务群发类型 -->
              <cy-form-item
                label="生效时间"
                prop="template_back.effective_time_type"
                :rules="rules.effective_time_type"
                required
              >
                <el-radio-group v-model="info.template_back.effective_time_type" style="width: 608px;">
                  <el-radio :label="1">按时间点</el-radio>
                  <el-radio :label="2">按时间段</el-radio>
                </el-radio-group>
                <el-form-item v-show="info.template_back.effective_time_type === 1" prop="template_back.start_time" :rules="rules.start_time">
                  <span class="m-r-5">在</span>
                  <el-date-picker
                    v-model="info.template_back.start_time"
                    type="date"
                    placeholder="选择时间"
                    value-format="timestamp"
                    style="width: 200px;"
                  ></el-date-picker>
                  <span class="m-l-5">日，开始生效</span>
                </el-form-item>
                <el-form-item v-if="info.template_back.effective_time_type === 2" style="width: 800px;" prop="template_back.end_time" :rules="rules.effective_time">
                  <span class="m-r-5">在</span>
                  <el-date-picker
                    v-model="effective_time"
                    type="daterange"
                    value-format="timestamp"
                    style="width: 200px;"
                    :default-time="['00:00:00', '23:59:59']"
                    range-separator="至"
                    start-placeholder="开始时间"
                    end-placeholder="结束时间"
                  ></el-date-picker>
                  <span class="m-l-5">期间生效</span>
                </el-form-item>
              </cy-form-item>
              <cy-form-item
                v-if="info.task_category === 0"
                label="任务群发类型"
                prop="group_send_type"
                :rules="rules.group_send_type"
                required
              >
                <el-radio-group v-model="info.template_back.group_send_type" style="width: 608px;">
                  <el-radio :label="1">企微一键群发（从群发助手发送）</el-radio>
                  <el-radio :label="2">灵活沟通转发（从任务中心发送）</el-radio>
                </el-radio-group>
              </cy-form-item>
              <!-- 群发顾客 -->
              <cy-form-item
                label="群发顾客"
                required
                v-if="info.task_category === 0 && info.template_back.group_send_type === 1"
              >
                <el-radio-group
                  v-model="info.template_back.allow_select"
                  style="width: 608px;"
                >
                  <el-radio :label="0">导购不可编辑</el-radio>
                  <el-radio :label="1">导购可编辑</el-radio>
                </el-radio-group>
              </cy-form-item>
              <cy-form-item
                label="任务完成标准"
                required
                :rules="rules.finish_rate"
                prop="finish_rate"
                v-if="info.task_category === 0 && info.template_back.group_send_type === 2"
              >
                <el-row class="flex">
                  <span class="m-r-5">顾客名单触达率达到</span>
                  <el-input v-model="info.template_back.finish_rate" style="width:100px"></el-input>
                  <span class="m-l-5">%</span>
                </el-row>
              </cy-form-item>
              <!-- 任务内容 -->
              <cy-form-item
                label="任务内容"
              >
                <task-content v-model="info.template_back.template_content" :task_category="info.task_category" ></task-content>
                <div class="message-info">※若未配置备选模板的任务素材内容，则默认发送原模板任务素材内容</div>
              </cy-form-item>
            </div>
          </cy-form>
        </div>
      </template>
      <template slot="footer">
        <el-button
          type="primary"
          @click="save()"
        >保 存</el-button>
        <el-button @click="cancel()">取 消</el-button>
      </template>
    </layout-form>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Watch } from "vue-property-decorator";
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
import TaskContent from "../components/task-content.vue";
const TYPE_MAP = {
  "text": 1,
  "picture": 2,
  "video": 4,
  "link": 6,
  "applet": 7,
}

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
function dealContent(template_content: any, template_content_type: string) {
  let tempObj: any = {};
  switch (template_content_type) {
    case "card":
      tempObj = {
        card: template_content.card
      };
      break;
    case "picture":
      tempObj = {
        image: {
          imgList: []
        }
      };
      template_content.picture.forEach((imgItem: any) => {
        tempObj.image.imgList.push({
          pic_url: imgItem.url,
          media_id: String(imgItem.id)
        });
      });
      break;
    case "video":
      tempObj = {
        video: {
          media_id: String(template_content.video.id),
          video_url: template_content.video.url
        }
      };
      break;
    case "applet":
      tempObj = {
        applets: {
          ...template_content.applet
        }
      };
      break;
    case "link":
      tempObj = {
        link: {
          ...template_content.link,
          pic_url: template_content.link.cover_url
        }
      };
      break;
  }
  return tempObj;
}
function parseTemplateContent(template_content: any) { 
  const params: any = {}
  const template_content_type = getMaterialType(
    template_content
  );
  // 处理素材数据
  params.template_content = dealContent(
    template_content,
    template_content_type
  );
  // 需要将text文本拼接上
  params.template_content.text = {
    content: template_content.text.content
  };
  params.template_content_type = TYPE_MAP[template_content_type];
    
  return params;
}

function parseResult(templateContents: any) {

  // 这里需要重新组装数据
  let newContents: any = {};
  newContents = templateContents;
  // 处理数据
  if (newContents.text.content) {
    newContents.text = {
      content: newContents.text.content
    };
  }
  if (newContents.image) {
    newContents.picture = newContents.image.imgList.map((item: any) => {
      return {
        ...item,
        url: item.pic_url
      };
    });
  } else {
    newContents.picture = []
  }
  if (newContents.video) {
    newContents.video = {
      ...newContents.video,
      url: newContents.video.video_url
    };
  }
  if (newContents.applets) {
    newContents.applet = newContents.applets;
  }
  if (newContents.link) {
    newContents.link = {
      ...newContents.link,
      cover_url: newContents.link.pic_url
    };
  }
  console.log(111111, newContents);
  
  return newContents;
}


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
    openVideo,
    TaskContent
  },
  mixins: [LocalUpload],
})
export default class extends Vue {
  /**************************************************data开始**************************************************/
  is_edit = 0; // 是否是编辑
  id = 0; // 朋友圈id
  edit_member_id = 0; // 当前员工的member_id,编辑单人的单条朋友圈时，必传
  moment_topic_member_id = 0; // 编辑单人的单条朋友圈时，必传
  mark: any = ""; // 链接加密mark
  materialSelectList: Array<any> = TASK_MATERIAL_SELECT_LIST_1;
  materialSelectList2: Array<any> = TASK_MATERIAL_SELECT_LIST_2;
  contentEmpty = false; // 判断发布内容是否为空
  showError = false; // 控制发布内容错误提示是否展示
  currentMaterialType = ""; // 当前选择素材的类型
  showMaterialModal = ""; // 控制当前选择素材的弹框显隐
  isOpenVideo = false; //是否打开视频
  videoUrl = ""; // 当前播放视频的url
  rules = {
    template_name: [
      { required: true, message: "请输入模板名称", trigger: "blur" }
    ],
    template_desc: [
      { required: true, message: "请输入模板描述", trigger: "blur" },
      { max: 200, message: "模板描述不能超过200个字符", trigger: "blur" }
    ],
    group_send_type: [
      { required: true, message: "请选择群发类型", trigger: "change" }
    ],
    finish_rate: [{ validator: (rule, value, callback) => {
      if (!/^[1-9][0-9]?$|100/.test(value) || Number(value)>100) {
        callback(new Error("请输入1-100的整数"));
      } else {
        callback();
      }
    }, trigger: 'blur' }],
    effective_time_type: [
      { required: true, message: "请选择生效时间", trigger: "change" }
    ],
    start_time: [
      { required: true, message: "请选择生效时间", trigger: "change" }
    ],
    effective_time: [
      { required: true, message: "请选择生效时间", trigger: "change" }
    ],
    template_content: [
      { 
        validator: (rule, value, callback) => {
          if (
            value.text && value.text.content ||
            (value.picture && value.picture.length > 0) ||
            (value.link && value.link.link_url) ||
            (value.video && value.video.url) ||
            (value.applet && value.applet.app_name)
          ) {
            callback();
          } else {
            callback(new Error("请填写任务内容"));
          }
      }, message: "请填写任务内容" }
    ],
  };
  appearanceInfo = {
    avatar: "", // 头像
    background: "", // 背景图
    name: "", // 名称
    sign: "" // 签名
  };
  info: any = {
    template_name: "",
    template_desc: "",
    template_content_type: 1,
    group_send_type: 1,
    task_category: 0,
    allow_select: 0,
    template_content: {
      text: {
        content: ""
      },
      video: {},
      applet: {},
      link: {},
      picture: []
    },
    // 任务完成标准
    finish_rate: 100,
    open_back_up: 0,
    template_back: {
      finish_rate: 100,
      group_send_type : 1,
      allow_select: 0,
      effective_time_type: 1, // 1-按照时间点 2-按时间范围
      start_time: '',
      end_time: '',
      template_content_type: 1,
      template_content: {
        text: {
          content: ""
        },
        video: {},
        applet: {},
        link: {},
        picture: []
      },
    }
  };
  showSelectText = false; // 控制选择文本弹框显隐

  get ep() {
    return this.$store.getters['enterprise']
  }

  get effective_time () {
    const { start_time, end_time } = this.info.template_back
    if (!start_time || !end_time) {
      return []
    }
    return [start_time, end_time]
  }
  set effective_time (val: Array<string>) {
    const [start_time, end_time] = val
    this.info.template_back.start_time = start_time
    this.info.template_back.end_time = end_time
  }

  /**************************************************data结束**************************************************/
  /***********************************************methods开始**************************************************/
  
  // 关闭预览视频
  // 处理文本内容
  public get textStr(): string {
    return this.info.template_content.text.content
      .replace(/\n/g, "<br/>")
      .replace(/\s/g, "&nbsp;");
  }
  // 点击发布按钮
  save() {
    (this.$refs["cform"] as any).validate((valid: any) => {
      if (valid) {
        if (this.contentEmpty) {
          // this.showError = false;
          return;
        }
        if (
          this.info.template_content_type == 1 &&
          this.info.template_content.text.content.length > 500
        ) {
          this.$message.warning("文本内容最多500个字符");
          return;
        }
        const params = {
          template_name: this.info.template_name,
          template_desc: this.info.template_desc,
          template_content: this.info.template_content,
          template_content_type: 1,
          group_send_type: this.info.group_send_type,
          finish_rate: Number(this.info.finish_rate),
          ...parseTemplateContent(this.info.template_content),
          open_back_up: this.info.open_back_up,
          allow_select: this.info.allow_select,
          template_back: this.info.template_back ? {
            ...this.info.template_back,
            ...parseTemplateContent(this.info.template_back.template_content),
            start_time: this.info.template_back.start_time ? parseInt(this.info.template_back.start_time/1000 + '') : 0,
            end_time: this.info.template_back.effective_time_type === 2 && this.info.template_back.end_time ? parseInt(this.info.template_back.end_time/1000 + '') : 0,
            finish_rate: Number(this.info.template_back.finish_rate),
          } : {},
          task_category: this.info.task_category,
        };
        // 判断是新增还是编辑
        if (this.$route.query.id) {
          params["id"] = Number(this.$route.query.id);
          Api.editTaskTemp(params)
            .then((res: any) => {
              if (res.code == 200) {
                this.$message.success("保存成功");
                this.$router.go(-1);
              } else {
                this.$message.error(res.msg);
              }
            })
            .catch((err: any) => {
              this.$message.error(err.msg);
            });
        } else {
          Api.createTaskTemp(params)
            .then((res: any) => {
              if (res.code == 200) {
                this.$message.success("保存成功");
                this.$router.go(-1);
              } else {
                this.$message.error(res.msg);
              }
            })
            .catch((err: any) => {
              this.$message.error(err.msg);
            });
        }
      } else {
        return;
      }
    });
  }
  @Watch("info.template_content", { deep: true })
  handleWatch(content: any) {
    if (
      content.text && content.text.content ||
      (content.picture && content.picture.length > 0) ||
      (content.link && content.link.link_url) ||
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
  /**
   * @description:判断当前素材类型
   * @return {*}
   */
  getMaterialType(template_content: any) {
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
  // 获取任务模板详情
  getTempInfo() {
    const { id } = this.$route.query;
    Api.getTaskTempDetail({ id: Number(id) }).then((res: any) => {
      res.data.template_back.start_time = res.data.template_back.start_time ? res.data.template_back.start_time * 1000 : ''
      res.data.template_back.end_time = res.data.template_back.end_time ? res.data.template_back.end_time * 1000 : ''
      res.data.template_back.effective_time_type = res.data.template_back.end_time ? 2 : 1
      if(res.data.template_back.template_content) {
        res.data.template_back.template_content = parseResult(JSON.parse(res.data.template_back.template_content))
      } else {
        res.data.template_back.template_content = {
          text: {
            content: ""
          },
          video: {},
          applet: {},
          link: {},
          picture: []
        }
      }
      res.data.template_back.group_send_type = res.data.template_back.group_send_type || 1
      
      this.info = res.data;
      const { template_content } = res.data;
      if (template_content) {
        
        this.info.template_content = parseResult(JSON.parse(template_content));
        console.log(this.info.template_content);
        
        // 判断当前素材类型
        this.currentMaterialType = getMaterialType(
          this.info.template_content
        );
        if (this.currentMaterialType == "text") {
          this.currentMaterialType = "";
        }
        console.log(this.currentMaterialType);
      }
    });
  }
  onContentChange(val: any) {
    console.log(val);
    console.log(this.info.template_back.template_content);
    
  }
  handleCategoryChange() {
    this.info.template_content = {
      text: {
        content: ""
      },
      video: {},
      applet: {},
      link: {},
      picture: []
    }
    this.info.template_back.template_content = {
      text: {
        content: ""
      },
      video: {},
      applet: {},
      link: {},
      picture: []
    }
    this.contentEmpty = true
  }
  cancel() {
    this.$router.go(-1);
  }
  created() {
    const { is_edit, id } = this.$route.query;
    this.is_edit = Number(is_edit) || 0;
    this.id = Number(id) || 0;
    if (!this.is_edit) {
      // 路由名字
      (this.$route.meta as any).title = "新建任务模板";
      this.contentEmpty = true;
      this.showError = false;
    }
    if (this.is_edit && id) {
      // 路由名字
      (this.$route.meta as any).title = "编辑任务模板";
      this.getTempInfo();
    }
    this.materialSelectList = this.ep.isCoach ? TASK_MATERIAL_SELECT_LIST_1 : TASK_MATERIAL_SELECT_LIST_1.slice(0, 4)
  }
}
</script>
<style lang="scss" scoped>
div {
  box-sizing: border-box;
}

.align-center {
  display: flex;
  align-items: center;
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
