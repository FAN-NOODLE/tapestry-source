<template>
  <div>
    <layout-form type="mobile">
      <template slot="left">
        <cy-form
          ref="tform"
          :rules="{}"
          :model="pageConfig"
          label-width="120px"
        >
          <cy-form-item label="活动信息展现样式" prop="active_type">
            <el-radio-group v-model="pageConfig.type" @change="changeType">
              <el-radio :label="1">文字</el-radio>
              <el-radio :label="2">图片</el-radio>
            </el-radio-group>
          </cy-form-item>
          <div v-if="pageConfig.type === 2">
            <cy-form-item label="活动信息图片" prop="active_pic">
              <upload-image
                :urls="pageConfig.active_pic"
                accept="image/jpeg,image/png"
                :before-upload="
                  (file) => {
                    return beforeUpload(file, 1);
                  }
                "
                @change="uploadSuccess1"
              ></upload-image>
              <div class="tip">
                图片格式为jpg、png文件；图片尺寸：750*960px
              </div>
              <div class="warm-tip">
                温馨提示：<br />
                1、图片内容需包含”活动标题、活动封面、活动有效期及活动规则”<br />
                2、为了页面整体美观，建议上传的图片底部采取渐变色设计。
              </div>
            </cy-form-item>
          </div>
          <div v-else>
            <cy-form-item label="文字颜色配置">
              <div class="color-content">
                <div class="color-item">
                  <span>活动标题颜色:</span>
                  <el-color-picker
                    class="color-pick"
                    v-model="pageConfig.color1"
                  ></el-color-picker>
                </div>
                <div class="color-item">
                  <span>活动有效期颜色:</span>
                  <el-color-picker
                    class="color-pick"
                    v-model="pageConfig.color2"
                  ></el-color-picker>
                </div>
                <div class="color-item">
                  <span>活动规则颜色:</span>
                  <el-color-picker
                    class="color-pick"
                    v-model="pageConfig.color3"
                  ></el-color-picker>
                </div>
              </div>
            </cy-form-item>
            <cy-form-item label="背景颜色">
              <el-color-picker
                class="color-pick"
                v-model="pageConfig.color4"
              ></el-color-picker>
            </cy-form-item>
          </div>
          <cy-form-item label="按钮颜色配置">
            <div class="color-content">
              <div class="color-item">
                <span>实心按钮颜色:</span>
                <el-color-picker
                  class="color-pick"
                  v-model="pageConfig.btn_color1"
                ></el-color-picker>
              </div>
              <div class="color-item">
                <span>实心按钮文字颜色:</span>
                <el-color-picker
                  class="color-pick"
                  v-model="pageConfig.btn_color2"
                ></el-color-picker>
              </div>
              <div class="color-item">
                <span>空心按钮颜色:</span>
                <el-color-picker
                  class="color-pick"
                  v-model="pageConfig.btn_color3"
                ></el-color-picker>
              </div>
              <div class="color-item">
                <span>空心按钮文字颜色:</span>
                <el-color-picker
                  class="color-pick"
                  v-model="pageConfig.btn_color4"
                ></el-color-picker>
              </div>
            </div>
          </cy-form-item>
          <cy-form-item label="分享标题">
            <cy-input
              :trim="true"
              placeholder="请输入分享内容标题"
              :maxlength="20"
              :showWordLimit="true"
              v-model="pageConfig.share_tit"
              widthSize="medium"
            ></cy-input>
            <div class="tip">
              分享标题为空，默认使用活动标题内容，输入文本建议在20个字内！
            </div>
          </cy-form-item>
          <cy-form-item
            label="分享图片上传"
            prop="share_pic"
            style="width: 600px"
          >
            <upload-image
              :urls="pageConfig.share_pic"
              accept="image/jpeg,image/png"
              :before-upload="
                (file) => {
                  return beforeUpload(file, 2);
                }
              "
              @change="uploadSuccess2"
            ></upload-image>
            <div class="tip">
              分享图片为空，默认使用活动封面图片；
            </div>
            <div class="tip">
              图片格式为jpg、png文件；图片尺寸：800*600px
            </div>
          </cy-form-item>
        </cy-form>
      </template>
      <template slot="right">
        <div class="phone-content">
          <div class="phone-title">页面效果预览</div>
          <div class="phone">
            <div class="phone-header">
              <img src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/mobile-header.png" alt="" />
            </div>
            <div class="top-bar">
              <i class="el-icon-arrow-left back-icon"></i>
              <span class="active-title">活动名称</span>
              <img src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/header-icon.png" alt="" />
            </div>
            <div class="phone-body">
              <div class="body-header" :style="bgStyle">
                <div v-if="pageConfig.type === 1">
                  <div class="active-title">
                    <div :style="{ color: pageConfig.color1 }">
                      {{ base_info_copy.title }}
                    </div>
                  </div>
                  <div class="body-card">
                    <div
                      class="active-time"
                      :style="{ color: pageConfig.color2 }"
                    >
                      活动时间：{{
                        base_info_copy.start_time
                          | formatDate("yyyy/mm/dd hh:mm")
                      }}
                      ~
                      {{
                        base_info_copy.end_time | formatDate("yyyy/mm/dd hh:mm")
                      }}
                    </div>
                    <div class="active-banner">
                      <img :src="base_info_copy.cover_url" alt="" />
                    </div>
                    <div
                      class="active-mark"
                      :style="{ color: pageConfig.color3 }"
                    >
                      {{ base_info_copy.activity_rule }}
                      <i class="el-icon-arrow-right arrow-icon"></i>
                    </div>
                  </div>
                </div>
              </div>
              <div class="sign-in-box">
                <div
                  class="btn-msg"
                  :style="{
                    background: pageConfig.btn_color1,
                    color: pageConfig.btn_color2,
                  }"
                >
                  {{
                    base_info_copy.activity_type === 1
                      ? "报名"
                      : base_info_copy.activity_type === 2
                      ? "签到"
                      : "问卷"
                  }}信息
                </div>
                <img
                  class="message-bnt"
                  src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/marketing-campaign/title-bg.png"
                />
              </div>
              <div class="form-content">
                <sign-in-content
                  v-if="base_info_copy.activity_type === 2 && isShow"
                  :startTime="base_info_copy.sign_start"
                  :endTime="base_info_copy.sign_end"
                ></sign-in-content>
                <div class="signPage" v-else>
                  <div
                    class="activity-detail-info"
                    v-if="templateList && templateList.length > 0"
                  >
                    <div
                      class="widget"
                      v-for="(item, index) in templateList"
                      :key="index"
                    >
                      <!-- <p class="title-p">{{ item.mark }}</p> -->
                      <p class="title-new">
                        <span v-if="item.isRequired" class="title-one">*</span>
                        <b class="title-b" style="margin-right: 4px">/</b>
                        {{ item.tit }}
                        <b class="title-b" style="margin-left: 4px">/</b>
                        <b class="title-line"></b>
                        <span class="title-two">{{ ++index }}</span>
                      </p>
                      <p class="vice-title-p" v-if="item.mark">
                        {{ item.mark }}
                      </p>
                      <!-- 单选 -->
                      <div v-if="item.type == 1" class="widget">
                        <div
                          class="custom-radio"
                          v-for="(items, i) in item.options"
                          :key="i"
                        >
                          <div>
                            <!-- <el-image fit="cover"
                              v-if="i == 0"
                              src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/marketing-campaign/choice.png"
                              alt=""
                            /> -->
                            <img
                              src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/marketing-campaign/no-choice.png"
                              alt=""
                            />
                            <!-- <span v-if="i == 0" style="color: #caa361">{{
                              items.option
                            }}</span> -->
                            <span>{{ items.option }}</span>
                          </div>
                          <div
                            class="pic-wraper"
                            v-if="items.pic && items.pic.length > 0"
                          >
                            <el-image fit="cover"
                              v-for="(img, pIndex) in items.pic"
                              :key="pIndex"
                              style="width: 90px; height: 90px; margin: 2px 10px"
                              :src="img"
                            >
                            </el-image>
                          </div>
                        </div>
                        <div class="custom-radio" v-if="item.isOther">
                          <div>
                            <!-- <el-image fit="cover"
                              v-if="i == 0"
                              src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/marketing-campaign/choice.png"
                              alt=""
                            /> -->
                            <img
                              src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/marketing-campaign/no-choice.png"
                              alt=""
                            />
                            <!-- <span v-if="i == 0" style="color: #caa361">{{
                              items.option
                            }}</span> -->
                            <span>其他</span>
                          </div>
                        </div>
                      </div>
                      <!-- 单行文本输入 -->
                      <div v-if="item.type == 2" class="widget">
                        <div class="input-style">
                          <span class="hint-text">请输入</span>
                        </div>
                      </div>
                      <!-- 多选 -->
                      <div v-if="item.type == 3" class="widget">
                        <div
                          class="custom-radio"
                          v-for="(items, i) in item.options"
                          :key="i"
                        >
                          <div>
                            <!-- <el-image fit="cover"
                              v-if="i == 0 || i == 1"
                              src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/marketing-campaign/choice.png"
                              alt=""
                            /> -->
                            <img
                              src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/marketing-campaign/no-choice.png"
                              alt=""
                            />
                            <!-- <span
                              v-if="i == 0 || i == 1"
                              style="color: #caa361"
                              >{{ items.option }}</span
                            > -->
                            <span>{{ items.option }}</span>
                          </div>
                          <div
                            class="pic-wraper"
                            v-if="items.pic && items.pic.length > 0"
                          >
                            <el-image fit="cover"
                              v-for="(img, pIndex) in items.pic"
                              :key="pIndex"
                              style="width: 90px; height: 90px; margin: 2px 10px"
                              :src="img"
                            >
                            </el-image>
                          </div>
                        </div>
                        <div class="custom-radio" v-if="item.isOther">
                          <div>
                            <!-- <el-image fit="cover"
                              v-if="i == 0"
                              src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/marketing-campaign/choice.png"
                              alt=""
                            /> -->
                            <img
                              src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/marketing-campaign/no-choice.png"
                              alt=""
                            />
                            <!-- <span v-if="i == 0" style="color: #caa361">{{
                              items.option
                            }}</span> -->
                            <span>其他</span>
                          </div>
                        </div>
                      </div>
                      <!-- 多行文本输入 -->
                      <div v-if="item.type == 4" class="widget">
                        <div
                          class="input-style"
                          style="height: 85px; align-items: flex-start"
                        >
                          <span class="hint-text">请输入</span>
                        </div>
                      </div>
                      <!-- 上传图片 -->
                      <div v-if="item.type == 5" class="widget">
                        <div class="upload-img">
                          <el-image fit="cover">
                            <div slot="error" class="image-slot">
                              <i class="el-icon-plus"></i>
                            </div>
                          </el-image>
                        </div>
                      </div>
                      <!-- 下拉选择 -->
                      <div v-if="item.type == 6" class="widget">
                        <div class="input-style">
                          <span class="hint-text">请选择</span>
                        </div>
                      </div>
                      <!-- 时间选择 -->
                      <div v-if="item.type == 7" class="widget">
                        <div class="input-style">
                          <span class="hint-text">请选择</span>
                        </div>
                      </div>
                      <!-- 输入姓名 -->
                      <div v-if="item.type == 8" class="widget">
                        <div class="input-style">
                          <span class="hint-text">请输入您的真实姓名</span>
                        </div>
                      </div>
                      <!-- 性别单选 -->
                      <div v-if="item.type == 9" class="widget">
                        <div
                          class="custom-radio"
                          v-for="(item, index) in item.options"
                          :key="index"
                        >
                          <!-- <el-image fit="cover"
                            v-if="index == 0"
                            src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/marketing-campaign/choice.png"
                            alt=""
                          /> -->
                          <img
                            src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/marketing-campaign/no-choice.png"
                            alt=""
                          />
                          <!-- <span v-if="index == 0" style="color: #caa361">{{
                            item.option
                          }}</span> -->
                          <span>{{ item.option }}</span>
                        </div>
                      </div>
                      <!-- 手机号 -->
                      <div v-if="item.type == 10" class="widget">
                        <div class="input-style">
                          <span class="hint-text">请输入您的手机号</span>
                        </div>
                      </div>
                      <!-- 邮箱 -->
                      <div v-if="item.type == 11" class="widget">
                        <div class="input-style">
                          <span class="hint-text">请输入您的邮箱地址</span>
                        </div>
                      </div>
                      <!-- 地址 -->
                      <div v-if="item.type == 12" class="widget">
                        <div class="select-style">
                          <div class="input-style" style="width: 30%">
                            <span class="hint-text">选择省</span>
                            <div class="select-icon"></div>
                          </div>
                          <div class="input-style" style="width: 30%">
                            <span class="hint-text">选择市</span>
                            <div class="select-icon"></div>
                          </div>
                          <div class="input-style" style="width: 30%">
                            <span class="hint-text">选择区</span>
                            <div class="select-icon"></div>
                          </div>
                        </div>
                        <div class="input-style" v-if="item.isAddress">
                          <span class="hint-text">请输入详细地址</span>
                        </div>
                      </div>
                      <!-- 邮编 -->
                      <div v-if="item.type == 13" class="widget">
                        <div class="input-style">
                          <span class="hint-text"
                            >请输入您所在区域的邮政编码</span
                          >
                        </div>
                      </div>
                    </div>
                  </div>
                  <div class="no-content" v-else>
                    暂无配置表单
                  </div>
                </div>
              </div>
              <div class="footer-btn">
                <div
                  class="btn"
                  :style="{
                    background: pageConfig.btn_color1,
                    color: pageConfig.btn_color2,
                  }"
                >
                  立即发送
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>
    </layout-form>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop, Watch } from "vue-property-decorator";
import UploadImage from "@/components/business/upload-image/index.vue";
import Api from "@/api/wisdom/marketing-activities";
import signInContent from "./signInContent.vue";
@Component({
  components: {
    UploadImage,
    signInContent,
  },
})
export default class extends Vue {
  @Prop({ default: () => 0 }) activityId?: number;
  @Prop({ default: () => 0 }) newActivityId?: number;
  @Prop() base_info_copy: any;
  @Prop() config: any;
  @Prop() content: any;
  @Watch("config")
  configChange(data: any) {
    try {
      const attribute = JSON.parse(data.attribute);
      this.pageConfig.type = data.attribute_type;
      if (this.pageConfig.type === 1) {
        this.pageConfig.color1 = attribute.word_color.title;
        this.pageConfig.color2 = attribute.word_color.period;
        this.pageConfig.color3 = attribute.word_color.rule;
        this.pageConfig.color4 = attribute.word_color.backgroup;
      } else {
        this.pageConfig.active_pic = attribute.img.url;
      }
      this.pageConfig.btn_color1 = attribute.button_color.solid_button;
      this.pageConfig.btn_color2 = attribute.button_color.solid_word;
      this.pageConfig.btn_color3 = attribute.button_color.hollow_button;
      this.pageConfig.btn_color4 = attribute.button_color.hollow_word;
      this.pageConfig.share_tit = data.share_title;
      this.pageConfig.share_pic = data.share_pic;
    } catch (error) {
      console.log(error);
    }
  }
  @Watch("content")
  contentChange(data: any) {
    try {
      this.templateList = JSON.parse(data).templateList || [];
    } catch (error) {
      this.templateList = [];
    }
  }
  private isShow = false;
  @Watch("base_info_copy")
  baseInfoChange(data: any) {
    this.isShow = false;
    this.base_info_copy = data;
    if (
      this.base_info_copy.sign_start &&
      this.base_info_copy.sign_start !== 0
    ) {
      setTimeout(() => {
        this.isShow = true;
      });
    }
  }
  private statusType = Number(this.$route.query.status_type);
  private pageConfig = {
    type: 1,
    active_pic: "",
    color1: "#FFFFFF",
    color2: "#808080",
    color3: "#33363C",
    color4: "#AF8B4E",
    btn_color1: "#AF8B4E",
    btn_color2: "#EFEFEF",
    btn_color3: "#AF8B4E",
    btn_color4: "#EFEFEF",
    share_tit: "",
    share_pic: "",
  };
  private radio = "1";
  private genderRadio = "1";
  private result = ["a", "b"];
  private templateList: Array<any> = [];
  private bgStyle: any = {
    background: this.pageConfig.color4,
  };
  private sign_start = 0;
  private sign_end = 0;
  @Watch("pageConfig.color4")
  onChangeColor() {
    this.bgStyle = {
      background: this.pageConfig.color4,
    };
  }
  @Watch("pageConfig.active_pic")
  onChangeImg() {
    this.bgStyle = {
      background: `url(${this.pageConfig.active_pic}) no-repeat`,
    };
  }
  private changeType() {
    if (this.pageConfig.type === 1) {
      this.bgStyle = {
        background: this.pageConfig.color4,
      };
    } else {
      this.bgStyle = {
        background: this.pageConfig.active_pic
          ? `url(${this.pageConfig.active_pic}) no-repeat`
          : "#FFFFFF",
      };
    }
  }
  private = {};
  // 上传文件限制规范
  private async beforeUpload(file: any, uploadtype: number) {
    const that = this;
    if (file.type != "image/jpeg" && file.type != "image/png") {
      this.$message.error("请上传JPG、PNG格式的图片");
      return new Promise((res, rej) => {
        return rej();
      });
    }
    let isImg = false;
    let fixedWidth = 0,
      fixedHeight = 0;
    if (uploadtype === 1) {
      fixedWidth = 750;
      fixedHeight = 960;
    } else if (uploadtype === 2) {
      fixedWidth = 800;
      fixedHeight = 600;
    }
    await new Promise(function(resolve: any, reject: any) {
      const _URL = window.URL || window.webkitURL;
      const image = new Image();
      image.onload = function() {
        isImg = true;
        const nowWidth = image.width;
        const nowHeight = image.height;
        const valid2 = nowWidth == fixedWidth && nowHeight == fixedHeight;
        valid2 ? resolve() : reject();
      };
      image.src = _URL.createObjectURL(file);
      setTimeout(() => {
        if (!isImg) {
          reject();
        }
      }, 300);
    }).then(
      () => {
        return file;
      },
      () => {
        (that as any).$message.error(
          `上传的尺寸不符合或文件有误, 建议${uploadtype === 1 ? "750" : "800"}*
          ${uploadtype === 1 ? "960" : "600"}px
          是最佳体验!`
        );
        return Promise.reject();
      }
    );
  }
  private uploadSuccess1(res: any) {
    this.pageConfig.active_pic = res;
  }
  private uploadSuccess2(res: any) {
    this.pageConfig.share_pic = res;
  }
  private submit() {
    if (this.pageConfig.type === 2 && !this.pageConfig.active_pic) {
      this.$message.warning("请先上传活动信息图片");
      return;
    }
    if (!this.pageConfig.share_pic) {
      this.pageConfig.share_pic = this.base_info_copy.cover_url;
    }
    if (this.pageConfig.share_tit == "") {
      this.pageConfig.share_tit = this.base_info_copy.title;
    }
    const params: any = {
      base_info: {
        activity_type: this.base_info_copy.activity_type,
      },
      attribute_type: this.pageConfig.type,
      attribute: {
        button_color: {
          solid_button: this.pageConfig.btn_color1,
          solid_word: this.pageConfig.btn_color2,
          hollow_button: this.pageConfig.btn_color3,
          hollow_word: this.pageConfig.btn_color4,
        },
      },
      share_title: this.pageConfig.share_tit,
      share_pic: this.pageConfig.share_pic,
    };
    if (this.pageConfig.type === 1) {
      params.attribute.word_color = {
        title: this.pageConfig.color1,
        period: this.pageConfig.color2,
        rule: this.pageConfig.color3,
        backgroup: this.pageConfig.color4,
      };
    } else {
      params.attribute.img = {
        url: this.pageConfig.active_pic,
      };
    }
    if (this.newActivityId !== 0 && this.newActivityId !== this.activityId) {
      params.id = this.newActivityId;
    } else {
      params.id = this.activityId;
    }
    if (this.statusType === 1 || this.statusType === 3) {
      params.steps = 3;
      Api.activityAddSubmit(params)
        .then((res: any) => {
          if (res.httpCode === 200) {
            this.$emit("submitSuccess", res.data.id || 0);
          } else {
            this.$message.warning("保存失败，请稍后重试！");
          }
        })
        .catch((err: any) => {
          console.log("activityOneSubmit", err);
        });
    } else if (this.statusType === 2) {
      params.steps = 3;
      Api.activityUpdata(params)
        .then((res: any) => {
          if (res.httpCode === 200) {
            this.$emit("submitSuccess", res.data.id || 0);
          } else {
            this.$message.warning("保存失败，请稍后重试！");
          }
        })
        .catch((err: any) => {
          console.log("activityOneSubmit", err);
        });
    }
  }
}
</script>
<style lang="scss" scoped>
.color-content {
  display: flex;
  flex-wrap: wrap;
  font-size: 15px;
  color: #888;
  .color-item {
    width: 50%;
    margin-bottom: 15px;
    &:nth-child(3) {
      margin: 0;
    }
    .color-pick {
      vertical-align: middle;
      margin-left: 5px;
    }
  }
}
.tip {
  font-size: 13px;
  color: #c0c4cc;
}
.warm-tip {
  font-size: 12px;
  background: #f4f7fa;
  line-height: 14px;
  color: #c0c4cc;
  padding: 10px;
}
/deep/ .layout-form .wrap .container .left {
  position: relative;
  width: calc(100% - 406px - 36px);
  padding-left: 30px;
}
/deep/ .layout-form .wrap .container .right {
  width: 382px;
  position: relative;
  min-height: 100%;
}
/deep/ .layout-form .wrap {
  padding-right: 0;
  padding-left: 0;
}
.phone-title {
  text-align: center;
  font-size: 18px;
  font-weight: 600;
  line-height: 60px;
}
.phone-content {
  background-color: #f5f6fa;
}
.phone {
  width: 375px;
  background-color: #fff;
  margin: 0 auto;
  .phone-header {
    width: 100%;
    height: 36px;
    img {
      width: 375px;
      margin: 11px 0;
    }
  }
  .top-bar {
    display: flex;
    justify-content: space-between;
    line-height: 32px;
    margin: 10px 0;
    .back-icon {
      width: 87px;
      font-size: 21px;
      line-height: 32px;
      padding-left: 12px;
    }
    .active-title {
      font-size: 17px;
      font-weight: 500;
    }
    img {
      width: 87px;
      margin-right: 7px;
    }
  }
  .phone-body {
    max-height: 600px;
    min-height: 600px;
    overflow: auto;
    .body-header {
      height: 360px;
      background-size: 100% 100% !important;
      .active-title {
        height: 66px;
        padding: 0 40px;
        font-size: 24px;
        font-weight: 500;
        text-align: center;
        display: flex;
        justify-content: center;
        align-items: center;
        line-height: 1.1;
      }
      .body-card {
        width: 345px;
        height: 273px;
        margin: 0 auto;
        background: linear-gradient(
          180deg,
          rgba(255, 255, 255, 0.3) 0%,
          rgba(255, 255, 255, 0.6) 100%
        );
        box-shadow: 0px 0px 5px 0px rgba(175, 139, 78, 0.3);
        border-radius: 12px;
        .active-time {
          width: 308px;
          height: 30px;
          margin: 0 auto;
          text-align: center;
          font-weight: 400;
          color: #808080;
          font-size: 12px;
          line-height: 30px;
          background: #ffffff;
          box-shadow: 0px 0px 5px 0px rgba(175, 139, 78, 0.3),
            0px 2px 8px 0px rgba(0, 0, 0, 0.09);
          border-radius: 0px 0px 12px 12px;
        }
        .active-banner {
          width: 315px;
          height: 158px;
          background: #ffffff;
          box-shadow: 0px 0px 5px 0px rgba(175, 139, 78, 0.3);
          border-radius: 12px;
          margin: 15px auto;
          overflow: hidden;
          img {
            width: 100%;
            height: 100%;
          }
        }
        .active-mark {
          padding: 0 15px;
          padding-right: 30px;
          height: 40px;
          font-size: 14px;
          font-weight: 400;
          color: #33363c;
          line-height: 20px;
          text-shadow: 0px 0px 5px rgba(175, 139, 78, 0.3);
          display: -webkit-box;
          -webkit-box-orient: vertical;
          -webkit-line-clamp: 2;
          overflow: hidden;
          position: relative;
          .arrow-icon {
            position: absolute;
            right: 15px;
            font-size: 14px;
            bottom: 2px;
          }
        }
      }
    }
    .sign-in-box {
      position: relative;
      margin: 25px 0;
      padding: 0 40px;
      font-size: 16px;
      .btn-msg {
        position: absolute;
        left: 0;
        right: 0;
        height: 30px;
        line-height: 30px;
        border-radius: 20px;
        background-color: red;
        color: #ffffff;
        width: 140px;
        text-align: center;
        margin: auto;
      }
      .message-bnt {
        width: 298px;
      }
    }
    .footer-btn {
      padding: 0 20px;
      margin: 20px 0;
      .btn {
        font-size: 16px;
        width: 100%;
        height: 40px;
        line-height: 40px;
        text-align: center;
        border-radius: 20px;
      }
    }
  }
}
.signPage {
  width: 100%;
  background: #ffffff;
  .activity-detail-info {
    padding: 0 20px;
    display: flex;
    flex-direction: column;
    .widget {
      margin: 0px 0px 20px 0px;
      .custom-radio {
        display: flex;
        flex-direction: column;
        font-size: 16px;
        color: #33363c;
        line-height: 22px;
        margin: 10px 0px;
        img {
          width: 16px;
          height: 16px;
          margin-right: 11px;
        }
      }
      .title-p {
        font-size: 16px;
        font-weight: bold;
        color: #33363c;
        line-height: 26px;
      }
      // 新的标题样式
      .title-new {
        display: flex;
        align-items: center;
        font-size: 15px;
        font-family: PingFangSC-Medium, PingFang SC;
        font-weight: bold;
        color: #33363c;
        line-height: 26px;
        .title-one {
          font-size: 18px;
          font-family: YouSheBiaoTiHei;
          color: #fa4046;
          line-height: 23px;
          margin-right: 3px;
        }
        .title-b {
          font-size: 20px;
          color: #33363c;
        }
        .title-line {
          width: 60px;
          height: 0px;
          border-bottom: 1px dashed #33363c;
          margin: 0px 7px;
        }
        .title-two {
          font-size: 22px;
          font-family: DINAlternate-Bold, DINAlternate;
          font-weight: bold;
          color: #caa361;
          line-height: 26px;
        }
      }
      .vice-title-p {
        font-size: 14px;
        color: #808080;
        line-height: 24px;
      }
      .upload-img {
        width: 98px;
        height: 98px;
        background: #f6f6f9;
        border-radius: 10px;
        margin-top: 15px;
        display: flex;
        justify-content: center;
        align-items: center;
      }
      .hint-text {
        font-size: 14px;
        color: #96979e;
        line-height: 20px;
      }
      .input-style {
        width: 315px;
        height: 44px;
        background: #f6f6f9;
        border-radius: 10px;
        margin-top: 15px;
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 10px;
      }
      .select-style {
        width: 100%;
        display: flex;
        justify-content: space-between;
        align-items: center;
      }
      .select-icon {
        width: 0;
        height: 0;
        border-right: 5px solid transparent;
        border-left: 5px solid transparent;
        border-top: 5px solid #96979e;
      }
    }
  }
}
.no-content {
  width: 100%;
  height: 300px;
  display: flex;
  justify-content: center;
  align-items: center;
}
/deep/ .el-form-item__label {
  width: 146px !important;
}
/deep/.label {
  max-width: 146px !important;
}
/deep/ .el-form-item__content {
  margin-left: 146px !important;
}
</style>
