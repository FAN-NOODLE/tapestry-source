<template>
  <div class="welcome">
    <layout-form>
      <template slot="content">
        <div class="chain-header">
          <div class="chain-header-bg">配置外链</div>
          <div class="chain-header-main">
            <el-form
              ref="cform"
              :rules="chainrules"
              :model="chainForm"
              label-width="90px"
            >
              <el-form-item label="一级渠道" prop="channelFirst">
                <cy-select
                  required
                  width-size="small"
                  size="small"
                  placeholder="请选择一级渠道"
                  v-model="chainForm.channelFirst"
                  @change="channelFirstChange"
                  :optionsList="channelFirstData"
                  :options="{ label: 'name', value: 'id' }"
                  :disabled="pageType !== 0"
                  search-input
                  clearable
                  filterable
                >
                </cy-select>
              </el-form-item>
              <el-form-item label="二级渠道" prop="channelSecond">
                <cy-select
                  required
                  width-size="small"
                  size="small"
                  placeholder="请选择二级渠道"
                  v-model="chainForm.channelSecond"
                  @change="channelSecondChange"
                  :optionsList="channelSecondData"
                  :options="{ label: 'name', value: 'id' }"
                  :disabled="pageType !== 0"
                  search-input
                  clearable
                  filterable
                >
                </cy-select>
              </el-form-item>
              <el-form-item>
                <el-button
                  size="small"
                  type="primary"
                  :disabled="chainForm.shortChain !== ''"
                  @click="universalShortChain('cform')"
                  >生成通用短外链</el-button
                >
                <div
                  class="short-chain-tip"
                  :class="{ 'has-chain-tip': chainForm.shortChain !== '' }"
                >
                  生成通用短外链，有中间页并可配置，在微信生态内/外均可唤醒小程序商城
                </div>
              </el-form-item>
              <el-form-item v-if="chainForm.shortChain">
                <div class="short-chain-txt">{{ chainForm.shortChain }}</div>
              </el-form-item>
            </el-form>
          </div>
        </div>
        <div class="chain-wrap">
          <div
            class="chain-header-bg chain-header-bg1"
            @click="setNoneFn()"
            :class="{ 'chain-header-active': setNone }"
          >
            生成外链中间页（非必填）
            <i></i>
            <span>{{ !setNone ? "收起" : "展开" }}</span>
          </div>
          <div class="create-wrap" :class="{ 'create-wrap-active': setNone }">
            <div class="sms-wrap">
              <div class="short-msg">
                {{ showSms ? "短信模版" : "中间页模版" }}
              </div>
              <div v-if="showSms" class="sms-show">
                <div class="sms-bg-wrap">
                  <p class="time">下午 6:30</p>
                  <div class="clear-fix">
                    <div
                      class="content-msg pull-left"
                      v-html="msgTextStr"
                      :attr="smsSignName"
                      ref="editWordLenDom"
                    ></div>
                  </div>
                  <div class="short-card">
                    <img
                      v-if="passParams.brand_img"
                      :src="passParams.brand_img"
                      class="card-brand-logo"
                    />
                    <div
                      v-if="passParams.card_title || passParams.card_desc"
                      class="card-title-wrap"
                    >
                      <div v-if="passParams.card_title" class="card-title">
                        {{ passParams.card_title }}
                      </div>
                      <div v-if="passParams.card_desc" class="card-tip">
                        {{ passParams.card_desc }}
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <div v-else class="browser-wrap">
                <div class="browser-header">
                  <span>vchangyi.com</span>
                </div>
                <div
                  class="browser-body"
                  :style="
                    `background:url(${passParams.background_img}) no-repeat center;background-size: 100% auto;`
                  "
                >
                  <div
                    class="browser-body-box"
                    v-if="passParams.spu_data.length > 0 || nowSwitchValue"
                  ></div>
                  <template v-if="nowSwitchValue">
                    <div class="browser-btn-wrap">
                      <div
                        class="browser-btn"
                        :style="{
                          color: passParams.button_font_color,
                          background: passParams.button_color
                        }"
                      >
                        {{ passParams.button_name }}
                      </div>
                    </div>
                    <div
                      v-if="passParams.location_desc"
                      class="browser-btn-tip"
                      :style="{ color: passParams.location_font_color }"
                    >
                      {{ passParams.location_desc }}
                    </div>
                  </template>
                  <div
                    class="two-spu-wrap"
                    v-if="passParams.spu_data.length > 0"
                  >
                    <div
                      v-for="(item, index) in passParams.spu_data"
                      :key="index"
                      class="two-spu-item"
                    >
                      <img :src="item.spu_img" class="one-spu-img" />
                      <div class="one-spu-info">
                        <div class="one-spu-name">{{ item.spu_name }}</div>
                        <div class="one-spu-sku">{{ item.sku_attribute }}</div>
                        <div class="one-spu-price-wrap">
                          <span class="one-spu-price-icon">¥</span>
                          <span class="one-spu-price">{{
                            item.sku_price
                          }}</span>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div v-if="passParams.active_img" class="banner-wrap">
                    <img
                      :src="passParams.active_img"
                      alt=""
                      class="banner-wrap-img"
                    />
                  </div>
                </div>
              </div>
              <div class="short-msg-bottom">
                <div
                  class="short-msg-img"
                  :class="{ 'short-msg-img-active': showSms }"
                  @click="preViewPage(1)"
                ></div>
                <div class="short-msg-watch">
                  {{ showSms ? "短信卡片预览图" : "中间页预览图" }}
                </div>
                <div
                  class="short-msg-img short-msg-img-right"
                  :class="{ 'short-msg-img-active': !showSms }"
                  @click="preViewPage(2)"
                ></div>
              </div>
            </div>
            <div class="content-right">
              <div class="content-right-title">
                配置短信卡片(部分安卓机可显示)
              </div>
              <el-form label-width="110px" :labelPosition="'right'">
                <el-form-item label="品牌图片">
                  <upload-image
                    v-if="pageType !== 2"
                    class="upload-template"
                    :urls="passParams.brand_img"
                    accept="image/jpeg,image/png"
                    :before-upload="beforeUploadLogo"
                    @change="changeTextMessageBackground"
                  >
                  </upload-image>
                  <span v-if="pageType !== 2" class="img-desc"
                    >支持jpg、png格式，请上传尺寸460*260px图片＜500kb</span
                  >
                  <template v-if="pageType === 2">
                    <img
                      v-if="passParams.brand_img !== ''"
                      :src="passParams.brand_img"
                      class="view-page-img"
                    />
                    <div v-else class="view-page-img2">——</div>
                  </template>
                </el-form-item>
                <el-form-item label="卡片标题">
                  <cy-input
                    trim
                    v-if="pageType !== 2"
                    v-model.trim="passParams.card_title"
                    maxlength="10"
                    minlength="1"
                    placeholder="请输入卡片标题"
                    :show-word-limit="true"
                    width-size="medium"
                    @focus="preViewPage(1)"
                  ></cy-input>
                  <div v-else>{{ passParams.card_title }}</div>
                </el-form-item>
                <el-form-item label="卡片活动说明">
                  <cy-input
                    trim
                    placeholder="请输入卡片活动说明"
                    v-if="pageType !== 2"
                    v-model.trim="passParams.card_desc"
                    maxlength="16"
                    minlength="1"
                    :show-word-limit="true"
                    width-size="medium"
                    @focus="preViewPage(1)"
                  ></cy-input>
                  <div v-else>{{ passParams.card_desc }}</div>
                </el-form-item>
              </el-form>
              <div class="hr" />
              <div class="content-right-title">配置外链中间页</div>
              <el-form label-width="80px" style="margin-left: 35px">
                <el-form-item label="背景图片">
                  <upload-image
                    v-if="pageType !== 2"
                    class="upload-template"
                    :urls="passParams.background_img"
                    accept="image/jpeg,image/png"
                    :before-upload="beforeUploadBgImg"
                    @change="changeUploadBgImg"
                  >
                  </upload-image>
                  <span v-if="pageType !== 2" class="img-desc"
                    >支持jpg、png 格式，请上传尺寸750*1334px图片＜500kb</span
                  >
                  <template v-if="pageType === 2">
                    <img
                      v-if="passParams.background_img !== ''"
                      :src="passParams.background_img"
                      class="view-page-img"
                    />
                    <div v-else class="view-page-img2">——</div>
                  </template>
                </el-form-item>
                <el-form-item label="按钮显示">
                  <el-switch
                    v-model="nowSwitchValue"
                    @change="changeSwitch"
                    :disabled="pageType === 2"
                    active-color="#6F3112"
                  >
                  </el-switch>
                  <span style="margin-left: 10px;display:inline-block;">{{
                    nowSwitchValue ? "已显示" : "已隐藏"
                  }}</span>
                  <div v-if="nowSwitchValue" class="switch-container">
                    <el-form label-width="90px">
                      <el-form-item
                        label="按钮名称"
                        :class="{ 'is-error': isErrorInfo.button_name_error }"
                      >
                        <cy-input
                          trim
                          v-if="pageType !== 2"
                          v-model.trim="passParams.button_name"
                          maxlength="5"
                          minlength="1"
                          :show-word-limit="true"
                          width-size="medium"
                          @focus="focusChange()"
                        ></cy-input>
                        <div v-else>
                          {{ passParams.button_name }}
                        </div>
                        <div
                          class="error-info"
                          v-if="isErrorInfo.button_name_error"
                        >
                          {{ isErrorInfo.setErrorCopy }}
                        </div>
                      </el-form-item>
                      <el-form-item label="按钮样式" style="margin-bottom: 0">
                        <el-button
                          v-if="pageType !== 2"
                          type="text"
                          @click="resetColor(1)"
                          class="item-btn-txt"
                          style="margin-bottom: 15px;"
                          >恢复默认</el-button
                        >
                        <div
                          class="switch-form-txt switch-form-txt-1"
                          :class="{
                            'is-error-btn-color': isErrorInfo.button_color_error
                          }"
                        >
                          <span class="item-span-color">按钮颜色</span>
                          <el-color-picker
                            @change="focusChange()"
                            :disabled="pageType === 2"
                            v-model="passParams.button_color"
                          ></el-color-picker>
                        </div>
                        <div
                          class="error-info"
                          v-if="isErrorInfo.button_color_error"
                        >
                          {{ isErrorInfo.setErrorCopy }}
                        </div>
                        <div
                          class="switch-form-txt"
                          :class="{
                            'is-error-btn-color':
                              isErrorInfo.button_font_color_error
                          }"
                        >
                          <span class="item-span-color">文字颜色</span>
                          <el-color-picker
                            @change="focusChange()"
                            :disabled="pageType === 2"
                            v-model="passParams.button_font_color"
                          ></el-color-picker>
                        </div>
                        <div
                          class="error-info"
                          v-if="isErrorInfo.button_font_color_error"
                        >
                          {{ isErrorInfo.setErrorCopy }}
                        </div>
                      </el-form-item>
                    </el-form>
                  </div>
                </el-form-item>
                <el-form-item label="跳转说明" v-if="nowSwitchValue">
                  <cy-input
                    trim
                    :class="{
                      'is-error-btn-input': isErrorInfo.location_desc_error
                    }"
                    v-if="pageType !== 2"
                    v-model.trim="passParams.location_desc"
                    maxlength="30"
                    minlength="1"
                    :show-word-limit="true"
                    width-size="medium"
                    @focus="focusChange()"
                  ></cy-input>
                  <div v-else>{{ passParams.location_desc }}</div>
                  <div
                    class="jump-tip-item"
                    :class="{
                      'is-error-btn-color':
                        isErrorInfo.location_font_color_error
                    }"
                  >
                    <span class="item-span-color item-span-color-1"
                      >文字颜色</span
                    >
                    <el-color-picker
                      @change="focusChange()"
                      :disabled="pageType === 2"
                      v-model="passParams.location_font_color"
                    ></el-color-picker>
                    <el-button
                      v-if="pageType !== 2"
                      type="text"
                      @click="resetColor(2)"
                      class="item-btn-txt item-btn-txt-1"
                      >恢复默认</el-button
                    >
                  </div>
                  <div
                    class="error-info"
                    style="margin-top: 5px"
                    v-if="
                      isErrorInfo.location_desc_error ||
                        isErrorInfo.location_font_color_error
                    "
                  >
                    {{ isErrorInfo.setErrorCopy }}
                  </div>
                </el-form-item>
                <el-form-item
                  label="附加链接1"
                  label-width="85px"
                  :class="{ 'is-error': isErrorInfo.additionalLink1Error }"
                >
                  <cy-cascader
                    ref="cascader1"
                    v-if="pageType !== 2"
                    @focus="focusChange()"
                    v-model="outerLinkPage.additionalLink1"
                    :options="additionalLink1"
                    :props="additionalLink1Props"
                    @change="
                      changeAdditionalLink(outerLinkPage.additionalLink1, 1)
                    "
                    placeholder="请选择附加链接"
                    clearable
                  ></cy-cascader>
                  <div
                    v-if="outerLinkPage.showAddBtn && pageType !== 2"
                    @click="additionalLink(1)"
                    class="add-link-btn"
                  >
                    <div class="add-link-btn-1"></div>
                    <div class="add-link-btn-1 add-link-btn-2"></div>
                  </div>
                  <div
                    v-if="
                      pageType === 2 &&
                        outerLinkPage.additionalLink1Name !== '/'
                    "
                  >
                    {{ outerLinkPage.additionalLink1Name }}
                  </div>
                  <div
                    class="error-info"
                    v-if="isErrorInfo.additionalLink1Error"
                  >
                    {{ isErrorInfo.additionalLink1ErrorCopy }}
                  </div>
                </el-form-item>
                <el-form-item
                  label-width="85px"
                  v-if="outerLinkPage.showLink2"
                  label="附加链接2"
                  :class="{ 'is-error': isErrorInfo.additionalLink2Error }"
                >
                  <cy-cascader
                    @focus="focusChange()"
                    @change="
                      changeAdditionalLink(outerLinkPage.additionalLink2, 2)
                    "
                    v-if="pageType !== 2"
                    v-model="outerLinkPage.additionalLink2"
                    :options="additionalLink2"
                    :props="additionalLink2Props"
                    placeholder="请选择附加链接"
                    clearable
                  ></cy-cascader>
                  <div
                    v-if="outerLinkPage.showLink2 && pageType !== 2"
                    @click="additionalLink(0)"
                    class="add-link-btn"
                  >
                    <div class="add-link-btn-1"></div>
                  </div>
                  <div v-if="pageType === 2">
                    {{ outerLinkPage.additionalLink2Name }}
                  </div>
                  <div
                    class="error-info"
                    v-if="isErrorInfo.additionalLink2Error"
                  >
                    {{ isErrorInfo.additionalLink2ErrorCopy }}
                  </div>
                </el-form-item>
                <el-form-item
                  label-width="82px"
                  v-if="showBannerImage"
                  label="banner图"
                  type="array"
                  :class="{ 'is-error': isErrorInfo.active_img_error }"
                >
                  <upload-image
                    v-if="pageType !== 2"
                    class="upload-template"
                    :urls="passParams.active_img"
                    accept="image/jpeg,image/png"
                    :before-upload="beforeUploadBanner"
                    @change="changeUploadBanner"
                  >
                  </upload-image>
                  <span v-if="pageType !== 2" class="img-desc"
                    >建议上传尺寸为280*120px的图片以保持最佳效果</span
                  >
                  <template v-if="pageType === 2">
                    <img
                      v-if="passParams.active_img !== ''"
                      :src="passParams.active_img"
                      class="view-page-img"
                    />
                    <div v-else class="view-page-img2">——</div>
                  </template>
                  <div class="error-info" v-if="isErrorInfo.active_img_error">
                    保存失败，请检查填写内容
                  </div>
                </el-form-item>
              </el-form>
            </div>
          </div>
        </div>
      </template>
      <template v-if="pageType === 2" slot="footer" class="dialog-footer">
        <el-button size="small" @click="cancel()">返 回</el-button>
        <el-button size="small" type="primary" @click="edit()">编 辑</el-button>
      </template>
      <template v-else slot="footer" class="dialog-footer">
        <el-button size="small" type="primary" @click="save()">保 存</el-button>
        <el-button size="small" @click="cancel()">取 消</el-button>
      </template>
    </layout-form>
    <back-pop v-if="isPop" @goBack="goBackFn"></back-pop>
  </div>
</template>

<script lang="ts">
let _this: any = null;
import API from "@/api/customer";
import UploadImage from "@/components/business/upload-image/index.vue";
import { cloneDeep, isEqual } from "lodash";
import backPop from "./back-pop.vue";

function getBeforeUploadLogoFunc(
  legalWidth: number,
  legalHeight: number,
  imgType: string
) {
  return async function beforeUploadLogo1(file: any) {
    // 上传图片前处理函数
    const isJPG = file.name
      .substring(file.name.lastIndexOf(".") + 1, file.name.length)
      .toLocaleLowerCase();
    const KB = 1024;
    const isLt5K = file.size < 500 * KB;
    if (isJPG !== "jpg" && isJPG !== "png") {
      (_this as any).$message.error(`上传的${imgType}只能是 jpg、png 格式!`);
      return Promise.reject();
    }
    if (!isLt5K) {
      (_this as any).$message.error(`上传的${imgType}大小不能超过 500Kb!`);
      return Promise.reject();
    }
    const isSize = await new Promise(function(resolve, reject) {
      const _URL = window.URL || window.webkitURL;
      const image = new Image();
      image.onload = function() {
        const valid = image.width == legalWidth && image.height == legalHeight;
        valid ? resolve(true) : reject();
      };
      image.src = _URL.createObjectURL(file);
    }).then(
      () => {
        return file;
      },
      () => {
        (_this as any).$message.error(
          `上传的${imgType}尺寸不符合, 只能是${legalWidth}*${legalHeight}!`
        );
        return Promise.reject();
      }
    );
    if (!isJPG || !isLt5K) {
      throw new Error("大小不合法");
    }
  };
}
const nowId = 0;
export default {
  name: "",
  components: {
    UploadImage,
    backPop
  },
  data() {
    const validateChannel1 = (rule: any, value: any, callback: any) => {
      if (value === "") {
        callback(new Error("请选择一级渠道"));
      } else if (
        (this as any).channelFirstTemp &&
        value !== (this as any).channelFirstTemp
      ) {
        if ((this as any).chainForm.shortChain !== "") {
          (this as any).chainForm.shortChain = "";
          callback(new Error("一级渠道已改变, 请重新生成"));
        }
      } else {
        callback();
      }
    };
    const validateChannel2 = (rule: any, value: any, callback: any) => {
      if (value === "") {
        callback(new Error("请选择二级渠道"));
      } else if (
        (this as any).channelSecondTemp &&
        value !== (this as any).channelSecondTemp
      ) {
        if ((this as any).chainForm.shortChain !== "") {
          (this as any).chainForm.shortChain = "";
          callback(new Error("二级渠道已改变, 请重新生成"));
        }
      } else {
        callback();
      }
    };
    return {
      isPop: false,
      isErrorInfo: {
        button_name_error: false,
        button_color_error: false,
        button_font_color_error: false,
        location_desc_error: false,
        location_font_color_error: false,
        additionalLink1Error: false,
        additionalLink2Error: false,
        active_img_error: false,
        buttonSetErrorCopy: "请填写按钮名称",
        locationSetErrorCopy: "请填写跳转说明",
        additionalLink1ErrorCopy: "请选择二级渠道不为空的附加链接",
        additionalLink2ErrorCopy: "请选择二级渠道不为空的附加链接"
      },
      showAddBtn: false,
      nowSwitchValue: 1,
      pageType: 2, // 0 新建 1 编辑 2 查看
      editId: "", // 查看/编辑的外链的id
      channelFirstData: [],
      channelSecondData: [],
      channelSecondTemp: "",
      channelFirstTemp: "",
      chainForm: {
        channelFirst: "",
        channelSecond: "",
        shortChain: ""
      },
      chainrules: {
        channelFirst: [
          {
            required: true,
            message: "请选择一级渠道",
            trigger: "change"
          },
          {
            validator: validateChannel1,
            trigger: "change"
          }
        ],
        channelSecond: [
          {
            required: true,
            message: "请选择二级渠道",
            trigger: "change"
          },
          {
            validator: validateChannel2,
            trigger: "change"
          }
        ]
      },
      isAlertSaveTip: true,
      setNone: false,
      showSms: true,
      msgTextStr: "【畅移科技】点击短外链进入小程序商城",
      smsSignName: "",
      additionalLink1: [], // 附加链接1的一级渠道数据
      additionalLink1LevelSecond: [], // 附加链接1的二级渠道
      additionalLink2: [], // 附加链接2的一级渠道数据
      additionalLink2LevelSecond: [], // 附加链接1的二级渠道
      additionalLink1Props: {
        label: "name",
        value: "id",
        disabled: "disabled",
        lazy: true,
        lazyLoad(node: any, resolve: any) {
          (_this as any).resolveChannelData(
            node && node.data && node.data.id,
            resolve,
            1
          );
        }
      },
      additionalLink2Props: {
        label: "name",
        value: "id",
        lazy: true,
        lazyLoad(node: any, resolve: any) {
          (_this as any).resolveChannelData(
            node && node.data && node.data.id,
            resolve,
            2
          );
        }
      },
      outerLinkPage: {
        additionalLink1: [],
        additionalLink1Name: "",
        additionalLink2: [],
        additionalLink2Name: "",
        showLink2: false,
        bannerImg: ""
      },
      deepChainForm: null,
      deepPassParams: null,
      deepOuterLinkPage: null,
      passParams: {
        page_url: "",
        one_level_name: "",
        two_level_name: "",
        short_url: "",
        brand_img: "",
        card_title: "",
        card_desc: "",
        background_img: "",
        is_hide: 0,
        button_name: "打开小程序",
        button_color: "#6F3112",
        button_font_color: "#fffffe",
        location_desc: "将自动为您打开小程序或点击上方按钮",
        location_font_color: "#6F3112",
        additional_type: 1,
        spu_data: [],
        active_img: "",
        active_page_url: ""
      },
      evenAddition1: [],
      evenAddition2: []
    };
  },
  created() {
    _this = this;
    // 获取一级渠道的数据
    (this as any).getChannelData().then((res: any) => {
      (this as any).channelFirstData = res.data.list;
      (this as any).additionalLink1 = res.data.list;
      (this as any).additionalLink2 = res.data.list;
    });
  },
  mounted() {
    const type = Number((this as any).$route.query.type);
    if (type === 0 && !localStorage.getItem("LOCALCHANNEL"))
      (this as any).deepCloneFn();
    (this as any).getFormData();
    (this as any).showSwitch();
  },
  methods: {
    showSwitch() {
      (this as any).nowSwitchValue =
        parseInt((this as any).passParams.is_hide) === 0;
    },
    changeSwitch(data: any) {
      (this as any).nowSwitchValue = data;
      (this as any).preViewPage(2);
    },
    focusChange() {
      (this as any).preViewPage(2);
      (this as any).clearError();
    },
    /**
     * 获取二级渠道的pageUrl方法
     *
     * @method getPageUrl
     */
    getPageUrl(id: number, data: any) {
      const nowObject = data.find((item: any) => item.id === id);
      if (nowObject === undefined) return "";
      return nowObject.page_url;
    },
    /***
     * 获取整个表单的数据
     * @params node 点击的渠道
     */
    async getFormData() {
      const type = Number((this as any).$route.query.type);
      (this as any).editId = Number((this as any).$route.query.id);
      (this as any).pageType = type;
      const params: object = {
        id: (this as any).editId
      };
      let nowData: object = {};
      if (type !== 0) {
        try {
          const res: any = await API.getFormData(params);
          if (res.code === 200) {
            nowData = res.data;
          } else {
            (this as any).$message.warning(res.message);
          }
        } catch (error) {
          console.log("error", error);
        }
      }
      if (type === 0) {
        if (localStorage.getItem("LOCALCHANNEL") !== null) {
          const nowDraftData: any = localStorage.getItem("LOCALCHANNEL");
          const nowDraftData2 = {
            two_level_id: JSON.parse(nowDraftData).chainForm.channelSecond,
            one_level_id: JSON.parse(nowDraftData).chainForm.channelFirst,
            short_url: JSON.parse(nowDraftData).chainForm.shortChain,
            config_params: JSON.parse(nowDraftData).passParams
          };
          (this as any).valueToData(nowDraftData2);
        }
      } else {
        (this as any).valueToData(nowData);
      }
    },
    /**
     * 给data中的数据赋值
     *
     * @method valueToData
     * @params data object 如果是新增那么优先从缓存中拿数据如果是编辑那么是从接口中取数据
     */
    valueToData(data: any) {
      (this as any).channelFirstTemp = data.one_level_id;
      (this as any).channelSecondTemp = data.two_level_id;
      let nowAdditionalLink1 = [];
      let nowAdditionalLink2: any = [];
      if (data.config_params.spu_data.length > 0) {
        nowAdditionalLink1 = [
          data.config_params.spu_data[0].spu_one_level_id,
          data.config_params.spu_data[0].spu_two_level_id
        ];
        if (data.config_params.spu_data.length > 1) {
          nowAdditionalLink2 = [
            data.config_params.spu_data[1].spu_one_level_id,
            data.config_params.spu_data[1].spu_two_level_id
          ];
        }
      } else {
        nowAdditionalLink1 = [
          data.config_params.active_one_level_id,
          data.config_params.active_two_level_id
        ];
      }
      (this as any).outerLinkPage = {
        additionalLink1: nowAdditionalLink1,
        additionalLink2: nowAdditionalLink2,
        additionalLink1Name: "",
        additionalLink2Name: "",
        showLink2: data.config_params.spu_data.length > 1,
        showAddBtn: data.config_params.spu_data.length > 1
      };
      (this as any).evenAddition1 = nowAdditionalLink1;
      (this as any).evenAddition2 = nowAdditionalLink2;
      (this as any).chainForm = {
        channelFirst: data.one_level_id,
        channelSecond: data.two_level_id,
        shortChain: data.short_url
      };
      (this as any).shortMsgTxt(data.short_url);
      (this as any).nowSwitchValue =
        data.config_params.is_hide === 1 ? false : true;
      (this as any).passParams = {
        page_url: data.page_url,
        one_level_name: data.one_level_name,
        two_level_name: data.two_level_name,
        two_level_id: data.two_level_id,
        brand_img: data.config_params.brand_img,
        card_title: data.config_params.card_title,
        card_desc: data.config_params.card_desc,
        background_img: data.config_params.background_img,
        button_name: data.config_params.button_name,
        button_color: data.config_params.button_color,
        button_font_color: data.config_params.button_font_color,
        location_desc: data.config_params.location_desc,
        location_font_color: data.config_params.location_font_color,
        additional_type: data.config_params.additional_type,
        spu_data: data.config_params.spu_data,
        active_img: data.config_params.active_img,
        active_page_url: data.config_params.active_page_url,
        active_one_level_id: data.config_params.active_one_level_id,
        active_two_level_id: data.config_params.active_two_level_id
      };
      (this as any).backfillData();
    },
    /**
     * 用于拿到一级渠道ID请求二级渠道并回填数据
     *
     * @method backfillData
     */
    backfillData() {
      (this as any)
        .getChannelData((this as any).chainForm.channelFirst)
        .then((res: any) => {
          (this as any).channelSecondData = res.data.list;
        });
      Promise.all([(this as any).getChannelData()])
        .then(result => {
          (this as any)
            .getChannelData((this as any).outerLinkPage.additionalLink1[0])
            .then((res: any) => {
              const nowData1 = res.data.list.filter(
                (item: any) => item.url_type !== 1
              );
              (this as any).additionalLink1LevelSecond = nowData1;
              const firstChannelName = (this as any).getNameById(
                (this as any).outerLinkPage.additionalLink1[0],
                (this as any).channelFirstData
              );
              const secondChannelName = (this as any).getNameById(
                (this as any).outerLinkPage.additionalLink1[1],
                res.data.list
              );
              (this as any).passParams.page_url = (this as any).getPageUrl(
                (this as any).outerLinkPage.additionalLink1[1],
                (this as any).channelSecondData
              );
              (this as any).outerLinkPage.additionalLink1Name =
                firstChannelName + "/" + secondChannelName;
              (this as any).additionalLink1 = (this as any).mergeData(
                (this as any).additionalLink1,
                nowData1,
                (this as any).outerLinkPage.additionalLink1[0]
              );
              const nowLink =
                (this as any).outerLinkPage.additionalLink1 &&
                (this as any).outerLinkPage.additionalLink1[1];
              const nowChoose = res.data.list.find(
                (item: any) => item.id === nowLink
              );
              if (nowChoose !== undefined) {
                const url_type = nowChoose.url_type === 2;
                (this as any).outerLinkPage.showAddBtn = url_type;
                if (url_type) {
                  (this as any).clearData();
                }
              }
            });
          (this as any)
            .getChannelData((this as any).outerLinkPage.additionalLink2[0])
            .then((res: any) => {
              let nowData2 = res.data.list.filter(
                (item: any) => item.url_type !== 1
              );
              nowData2 = (this as any).disabledNonCommercial(nowData2);
              (this as any).additionalLink2LevelSecond = nowData2;
              const firstChannelName = (this as any).getNameById(
                (this as any).outerLinkPage.additionalLink2[0],
                (this as any).channelFirstData
              );
              const secondChannelName = (this as any).getNameById(
                (this as any).outerLinkPage.additionalLink2[1],
                res.data.list
              );
              (this as any).outerLinkPage.additionalLink2Name =
                firstChannelName + "/" + secondChannelName;
              (this as any).additionalLink2 = (this as any).mergeData(
                (this as any).additionalLink2,
                nowData2,
                (this as any).outerLinkPage.additionalLink2[0]
              );
            });
          (this as any).deepCloneFn();
        })
        .catch(error => {
          console.log(error, 66666); // 失败了，打出 '失败'
        });
    },
    /**
     * 过滤二级渠道数据
     *
     * @method filterData
     */
    filterData(data: any) {
      return data.filter((item: any) => item.is_hide === 0);
    },
    /**
     * 拿到一级渠道数据resolve到二级渠道中
     * @params id 一级渠道的ID
     * @params resolve 方法自动获取二级渠道的时候选择
     * @params type 1为附加链接1 2为附加链接2
     * @method resolveChannelData
     */
    resolveChannelData(id: number, resolve: any, type: number) {
      const addtionData =
        type === 1
          ? (this as any).additionalLink1
          : (this as any).additionalLink2;
      const arr = addtionData.filter((item: any) => item.id === id);
      (this as any).getChannelData(id).then((res: any) => {
        const nowData = res.data.list.filter(
          (item: any) => item.url_type !== 1
        );
        if (type === 1) {
          if ((this as any).chainForm.shortChain !== "") {
            (this as any).additionalLink1LevelSecond = nowData.filter(
              (item: any) => item.id !== (this as any).chainForm.channelSecond
            );
          } else {
            (this as any).additionalLink1LevelSecond = nowData;
          }
        } else {
          (this as any).additionalLink2LevelSecond = (this as any).disabledNonCommercial(
            nowData
          );
          if (id === (this as any).outerLinkPage.additionalLink1[0]) {
            (this as any).additionalLink2LevelSecond = nowData.filter(
              (item: any) =>
                item.id !== (this as any).outerLinkPage.additionalLink1[1]
            );
          } else {
            (this as any).additionalLink2LevelSecond = nowData;
          }
        }
        const mapData =
          type === 1
            ? (this as any).additionalLink1LevelSecond
            : (this as any).additionalLink2LevelSecond;
        const nowList = mapData.map((item: any) => {
          item.leaf = 1;
          return item;
        });
        if (
          ((this as any).outerLinkPage.additionalLink1[0] === id &&
            type === 1) ||
          ((this as any).outerLinkPage.additionalLink2[0] === id && type === 2)
        ) {
          resolve();
          return false;
        }
        if (
          ((this as any).evenAddition1[0] === id && type === 1) ||
          ((this as any).evenAddition2 === id && type === 2)
        ) {
          resolve();
          return false;
        }
        resolve(nowList);
      });
    },
    /**
     *  重新请求数据并赋值
     *
     * @method toRequestForData
     */
    toRequestForData() {
      (this as any).getChannelData().then((res: any) => {
        (this as any).additionalLink1 = res.data.list;
        (this as any).additionalLink2 = res.data.list;
        (this as any).outerLinkPage.additionalLink1 = [];
        ((this as any).$refs
          .cascader1 as any).$refs.cascader.$refs.panel.activePath = [];
      });
    },
    /**
     * 禁用非商品的类型: 只有二级渠道会出现商品
     *
     * @method disabledNonCommercial
     */
    disabledNonCommercial(data: any) {
      return data.map((item: any) => {
        if (item.url_type === 2) {
          item.disabled = false;
        } else {
          item.disabled = true;
        }
        return item;
      });
    },
    /**
     * 获取渠道数据
     *
     * @method getChannelData
     * @params id 一级渠道的ID ID为0 表示一级渠道否则是二级渠道
     */
    getChannelData(id: number) {
      return new Promise((resolve, reject) => {
        const params: object = {
          enterprise_id:
            (this as any).$store.state.user.USERINFO.enterprise_id || "",
          pid: id || 0
        };
        const res: any = API.geChannelName(params);
        if (res.code === 200) {
          return res.data.list;
        }
        resolve(res);
      });
    },
    /**
     * 合并数据，将二级渠道合并到一级渠道
     *
     * @method mergeData
     * @params arr 一级渠道的数据
     * @params pushArr 二级渠道的相关数据
     * @params id 一级渠道的ID
     */
    mergeData(arr: Array<any>, pushArr: Array<any>, id: number) {
      return arr.map((item: any) => {
        if (item.id === id) {
          pushArr.map(i => {
            i.leaf = "leaf";
            return i;
          });
          item.children = pushArr;
        } else {
          item.children = [];
        }
        return item;
      });
    },
    // 生成外链区域隐藏
    setNoneFn() {
      (this as any).setNone = !(this as any).setNone;
    },
    /***
     * 一级渠道change事件
     */
    channelFirstChange(id: number) {
      (this as any).chainForm.channelSecond = "";
      (this as any).passParams.one_level_name = (this as any).getNameById(
        id,
        (this as any).channelFirstData
      );
      (this as any).getChannelData(id).then((res: any) => {
        (this as any).channelSecondData = (this as any).filterData(
          res.data.list
        );
      });
    },
    channelSecondChange(id: number) {
      (this as any).passParams.page_url = (this as any).getPageUrl(
        (this as any).chainForm.channelSecond,
        (this as any).channelSecondData
      );
      (this as any).passParams.two_level_name = (this as any).getNameById(
        id,
        (this as any).channelSecondData
      );
    },
    /**
     * 通过id来拿name
     * @method getNameById
     * @params id 附加链接的一级渠道ID
     * @params nowId 表示二级渠道ID
     */
    getNameById(id: any, nowData: any) {
      const nowObject =
        nowId === 0
          ? nowData.find((item: any) => item.id === id)
          : nowData.find((item: any) => item.id === nowId);
      if (nowObject === undefined) return "";
      return nowObject.name;
    },
    // 获取通用短链
    universalShortChain(formName: string) {
      (this as any).channelFirstTemp = (this as any).chainForm.channelFirst;
      (this as any).channelSecondTemp = (this as any).chainForm.channelSecond;
      (this as any).$refs[formName].validate((valid: any) => {
        if (valid) {
          (this as any).getGeneralShortChain();
        } else {
          console.log("error submit!!");
          return false;
        }
      });
    },
    /**
     * 获取通用短外链
     *
     * @method getGeneralShortChain
     */
    async getGeneralShortChain() {
      const params: object = {
        longURL: (this as any).getPageUrl(
          (this as any).chainForm.channelSecond,
          (this as any).channelSecondData
        ),
        short_type: 1
      };
      try {
        const res: any = await API.getGeneralShortChain(params);
        if (res.code === 200) {
          (this as any).chainForm.shortChain = res.data.shortURL;
          (this as any).shortMsgTxt((this as any).chainForm.shortChain);
          const flag = true;
          if (flag) {
            (this as any).copyLink((this as any).chainForm.shortChain);
            (this as any).channelFirstTemp = (this as any).chainForm.channelFirst;
            (this as any).channelSecondTemp = (this as any).chainForm.channelSecond;
            (this as any).clearData();
            (this as any).passParams.spu_data = [];
            (this as any).toRequestForData();
          } else {
            (this as any).$message.warning("短外链已生成, 且复制成功！");
          }
        } else {
          (this as any).$message.warning(res.message);
        }
      } catch (error) {
        console.log("error", error);
      }
    },
    /**
     * 复制链接
     *
     * @method copyLink
     * @param {data 要复制的那条数据的信息}
     */
    copyLink(data: string) {
      const url: string = data;
      if (url === "") return false;
      const oInput: HTMLTextAreaElement = document.createElement("textarea");
      oInput.value = url;
      document.body.appendChild(oInput);
      oInput.select(); // 选择对象;
      document.execCommand("Copy"); // 执行浏览器复制命令
      (this as any).$message({
        message:
          "已生成通用短外链，有中间页可配置，在微信生态内/外均可唤醒小程序商城",
        type: "success",
        duration: 2000
      });
      oInput.remove();
    },
    // 短信内容
    shortMsgTxt(link: string) {
      (this as any).msgTextStr = `【畅移科技】点击短外链进入小程序商城<a href="javascript:;">${link}</a>`;
    },
    // 短信卡片图片上传验
    beforeUploadLogo: getBeforeUploadLogoFunc(460, 260, "品牌图片"),
    // 中间页背景图
    beforeUploadBgImg: getBeforeUploadLogoFunc(750, 1334, "背景图片"),
    // banner图
    beforeUploadBanner: getBeforeUploadLogoFunc(280, 120, "banner图"),
    changeTextMessageBackground(res: any) {
      (this as any).preViewPage(1);
      (this as any).passParams.brand_img = res;
    },
    changeUploadBgImg(file: any) {
      (this as any).preViewPage(2);
      (this as any).passParams.background_img = file;
    },
    changeUploadBanner(file: any) {
      (this as any).clearError();
      (this as any).preViewPage(2);
      (this as any).passParams.active_img = file;
    },
    beforeUploadOuterLink(file: any) {
      console.log("beforeUploadOuterLink", file);
    },
    updateOuterLink(file: any) {
      console.log("updateOuterLink", file);
    },
    // 预览短信卡片 or 中间页模版
    preViewPage(id: number) {
      (this as any).showSms = id === 1;
    },
    // 恢复默认
    resetColor(type: number) {
      switch (type) {
        case 1:
          (this as any).passParams.button_font_color = "#fffffe";
          (this as any).passParams.button_color = "#6F3112";
          break;
        case 2:
          (this as any).passParams.location_font_color = "#6F3112";
          break;
      }
    },
    // 增减附加链接1
    additionalLink(type: number) {
      if (type === 1) {
        (this as any).outerLinkPage.showLink2 = true;
      } else {
        (this as any).outerLinkPage.showLink2 = false;
        (this as any).passParams.spu_data.splice(1, 1);
        (this as any).outerLinkPage.additionalLink2 = [];
      }
    },
    /**
     * 获取商品详情页的相关参数
     *
     * @method getGoodsParams
     * @params pageUrl string 页面的URL
     * @params type number 1表示附加链接1 2表示附加链接2
     */
    async getGoodsParams(pageUrl: string, type: number) {
      const nowChooseData =
        type === 1
          ? (this as any).outerLinkPage.additionalLink1
          : (this as any).outerLinkPage.additionalLink2;
      const params: object = {
        page_url: pageUrl
      };
      try {
        const res: any = await API.getGoodsParams(params);
        if (res.code === 200) {
          const nowData = res.data;
          (this as any).$set((this as any).passParams.spu_data, type - 1, {
            spu_id: nowData.spu_id,
            spu_one_level_id: nowChooseData[0],
            spu_two_level_id: nowChooseData[1],
            sku_attribute: nowData.sku_attribute, // 商品规格
            sku_price: nowData.sku_price,
            spu_img: nowData.spu_img,
            spu_name: nowData.spu_name,
            spu_page_url: pageUrl
          });
        } else {
          (this as any).$message.warning(res.message);
        }
      } catch (error) {
        console.log("error", error);
      }
    },
    // 附加链接改变
    changeAdditionalLink(data: any, type: number) {
      let nowChoose: any = {};
      if (type === 1) {
        (this as any).outerLinkPage.additionalLink2 = [];
        if (data.length === 0) {
          (this as any).clearData();
        }
        if (data.length === 1) {
          (this as any).$message.error("请选择带有二级渠道的附加链接");
          (this as any).isErrorInfo.additionalLink1ErrorCopy =
            "请选择带有二级渠道的附加链接";
          (this as any).isErrorInfo.additionalLink1Error = true;
        }
        if ((this as any).filterChooseData()) {
          (this as any).outerLinkPage.showAddBtn = true;
        } else {
          (this as any).outerLinkPage.showAddBtn = false;
          (this as any).outerLinkPage.showLink2 = false;
        }
        nowChoose = (this as any).additionalLink1LevelSecond.find(
          (item: any) => item.id === data[1]
        );
        if ((this as any).passParams.spu_data.length > 0) {
          (this as any).outerLinkPage.showLink2 = false;
          (this as any).passParams.spu_data.splice(1, 1);
        }
      } else {
        if (data.length === 1) {
          (this as any).$message.error("请选择带有二级渠道的附加链接");
          (this as any).isErrorInfo.additionalLink2ErrorCopy =
            "请选择带有二级渠道的附加链接";
          (this as any).isErrorInfo.additionalLink2Error = true;
        }
        nowChoose = (this as any).additionalLink2LevelSecond.find(
          (item: any) => item.id === data[1]
        );
      }
      if (nowChoose === undefined) {
        (this as any).getChannelData(data[0]).then((res: any) => {
          const nowFilterData = res.data.list.filter(
            (item: any) => item.url_type !== 1
          );
          if (type === 1) {
            (this as any).additionalLink1LevelSecond = nowFilterData;
          } else {
            (this as any).additionalLink2LevelSecond = nowFilterData;
          }
          const nowChoose2 = nowFilterData.find(
            (item: any) => item.id === data[1]
          );
          if ((this as any).filterChooseData()) {
            (this as any).outerLinkPage.showAddBtn = true;
            (this as any).clearData();
            if (nowChoose2) {
              (this as any).getGoodsParams(nowChoose2.page_url, type);
            }
          } else {
            (this as any).passParams.spu_data = [];
            (this as any).passParams.active_one_level_id = data[0];
            (this as any).passParams.active_two_level_id = data[1];
            (this as any).passParams.active_page_url = (this as any).getPageUrl(
              data[1],
              (this as any).additionalLink1LevelSecond
            );
          }
        });
      } else {
        if ((this as any).filterChooseData()) {
          (this as any).outerLinkPage.showAddBtn = true;
          (this as any).clearData();
          if (nowChoose) {
            (this as any).getGoodsParams(nowChoose.page_url, type);
          }
        } else {
          (this as any).passParams.spu_data = [];
          (this as any).passParams.active_one_level_id = data[0];
          (this as any).passParams.active_two_level_id = data[1];
          (this as any).passParams.active_page_url = (this as any).getPageUrl(
            data[1],
            (this as any).additionalLink1LevelSecond
          );
        }
      }
    },
    clearData(): void {
      (this as any).passParams.active_img = "";
      (this as any).passParams.active_one_level_id = "";
      (this as any).passParams.active_two_level_id = "";
      (this as any).passParams.active_page_url = "";
    },
    /**
     * 取消
     *
     * @method cancel
     */
    cancel() {
      const { pageType } = this as any;
      if (pageType === 2) {
        (this as any).$router.push({
          name:
            "myapp_customer_auto-market_externalchain-manage_externalchain-list"
        });
      } else {
        (this as any).deepCompare(pageType);
      }
    },
    getName(id: any, nowData: any) {
      const nowObject = nowData.find((item: any) => item.id === id);
      return nowObject.name;
    },
    getUrlType(id: any, nowData: any) {
      const nowObject = nowData.find((item: any) => item.id === id);
      return nowObject.url_type;
    },
    /**
     * 保存的时候传递给后台的参数
     *
     * @method saveParams
     */
    saveParams() {
      return {
        url_type: (this as any).getUrlType(
          (this as any).chainForm.channelSecond,
          (this as any).channelSecondData
        ),
        id: Number((this as any).$route.query.id),
        operator: (this as any).$store.state.user.USERINFO.contact_name || "",
        enterprise_id:
          (this as any).$store.state.user.USERINFO.enterprise_id || "",
        page_url: (this as any).getPageUrl(
          (this as any).chainForm.channelSecond,
          (this as any).channelSecondData
        ),
        one_level_name: (this as any).getName(
          (this as any).chainForm.channelFirst,
          (this as any).channelFirstData
        ),
        one_level_id: (this as any).chainForm.channelFirst,
        two_level_name: (this as any).getName(
          (this as any).chainForm.channelSecond,
          (this as any).channelSecondData
        ),
        two_level_id: (this as any).chainForm.channelSecond,
        short_url: (this as any).chainForm.shortChain,
        brand_img: (this as any).passParams.brand_img,
        card_title: (this as any).passParams.card_title,
        card_desc: (this as any).passParams.card_desc,
        background_img: (this as any).passParams.background_img || "",
        is_hide: (this as any).nowSwitchValue === true ? 0 : 1,
        button_name: (this as any).passParams.button_name,
        button_color: (this as any).passParams.button_color,
        button_font_color: (this as any).passParams.button_font_color,
        location_desc: (this as any).passParams.location_desc,
        location_font_color: (this as any).passParams.location_font_color,
        additional_type: (this as any).passParams.additional_type,
        spu_data: (this as any).passParams.spu_data,
        active_img: (this as any).passParams.active_img,
        active_one_level_id: Number(
          (this as any).passParams.active_one_level_id
        ),
        active_two_level_id: Number(
          (this as any).passParams.active_two_level_id
        ),
        active_page_url: (this as any).passParams.active_page_url
      };
    },
    // 保存配置
    save() {
      if (
        (this as any).chainForm.shortChain === "" &&
        ((this as any).chainForm.shortChain !== "" ||
          (this as any).chainForm.channelSecond !== "")
      ) {
        (this as any).$message.warning("您还没有生成通用短外链");
        return false;
      }
      (this as any).$refs["cform"].validate((valid: any) => {
        if (valid) {
          if ((this as any).chainForm.shortChain === "") {
            (this as any).$message.warning("您还没有生成通用短外链");
            return false;
          } else {
            (this as any).configBtnAndCopy();
          }
        } else {
          if ((this as any).chainForm.channelFirst === "") {
            (this as any).$message.error("请选择一级渠道");
            return false;
          }
          if ((this as any).chainForm.channelSecond === "") {
            (this as any).$message.error("请选择二级渠道");
            return false;
          }
          if ((this as any).chainForm.shortChain === "") {
            (this as any).$message.error("渠道已改变, 请重新生成短外链");
            return false;
          }
          console.log("error submit!");
        }
      });
    },
    /**
     * 判断有木有配置按钮和跳转文案
     * @method configBtnAndCopy
     */
    configBtnAndCopy() {
      if ((this as any).nowSwitchValue) {
        if ((this as any).passParams.button_name === "") {
          (this as any).$message.error("保存失败，请检查填写内容");
          (this as any).isErrorInfo.button_name_error = true;
          return false;
        }
        if (!(this as any).passParams.button_color) {
          (this as any).$message.error("保存失败，请检查填写内容");
          (this as any).isErrorInfo.button_color_error = true;
          return false;
        }
        if (!(this as any).passParams.button_font_color) {
          (this as any).$message.error("保存失败，请检查填写内容");
          (this as any).isErrorInfo.button_font_color_error = true;
          return false;
        }
        if ((this as any).passParams.location_desc === "") {
          (this as any).$message.error("保存失败，请检查填写内容");
          (this as any).isErrorInfo.location_desc_error = true;
          return false;
        }
        if (!(this as any).passParams.location_font_color) {
          (this as any).$message.error("保存失败，请检查填写内容");
          (this as any).isErrorInfo.location_font_color_error = true;
          return false;
        }
      }
      (this as any).judgeCoincidence();
    },
    clearError(): void {
      setTimeout(() => {
        (this as any).isErrorInfo = {
          button_name_error: false,
          button_color_error: false,
          button_font_color_error: false,
          location_desc_error: false,
          location_font_color_error: false,
          additionalLink1Error: false,
          additionalLink2Error: false,
          active_img_error: false,
          additionalLink1ErrorCopy: "请选择二级渠道不为空的附加链接",
          additionalLink2ErrorCopy: "请选择二级渠道不为空的附加链接",
          setErrorCopy: "保存失败，请检查填写内容"
        };
      }, 0);
    },
    /**
     * 判断渠道们有木有重合如果重合报错
     *
     * @method judgeCoincidence
     */
    judgeCoincidence() {
      const second = (this as any).chainForm.channelSecond;
      if ((this as any).passParams.spu_data.length > 0) {
        const spu_two_level_id1 = (this as any).passParams.spu_data[0]
          .spu_two_level_id;
        const spu_two_level_id2 =
          (this as any).passParams.spu_data.length > 1
            ? (this as any).passParams.spu_data[1].spu_two_level_id
            : 0;
        if (second === spu_two_level_id1 || second === spu_two_level_id2) {
          (this as any).$message.error("不可添加相同的短外链与附加链接");
          if (second === spu_two_level_id1) {
            (this as any).isErrorInfo.additionalLink1ErrorCopy =
              "不可添加相同的短外链与附加链接";
            (this as any).isErrorInfo.additionalLink1Error = true;
          } else {
            (this as any).isErrorInfo.additionalLink2ErrorCopy =
              "不可添加相同的短外链与附加链接";
            (this as any).isErrorInfo.additionalLink2Error = true;
          }
          return false;
        }
        if (spu_two_level_id1 === spu_two_level_id2) {
          (this as any).$message.error("不可添加2个相同的附加链接");
          (this as any).isErrorInfo.additionalLink1ErrorCopy =
            "不可添加2个相同的附加链接";
          (this as any).isErrorInfo.additionalLink2ErrorCopy =
            "不可添加2个相同的附加链接";
          (this as any).isErrorInfo.additionalLink1Error = true;
          (this as any).isErrorInfo.additionalLink2Error = true;
          return false;
        }
        (this as any).isHasBtnAndBg();
      } else {
        const active_two_id = (this as any).passParams.active_two_level_id;
        if (active_two_id === second) {
          (this as any).$message.error("不可添加相同的短外链与附加链接");
        } else {
          (this as any).isHasBtnAndBg();
        }
      }
    },
    /**
     * 判断是否没配置按钮和背景图
     *
     * @method isHasBtnAndBg
     */
    isHasBtnAndBg() {
      const data: any = (this as any).saveParams();
      if (data.background_img === "" && data.is_hide === 1) {
        (this as any)
          .$confirm(
            "微信外部可打开 微信内部打开为空白页无法跳转到小程序",
            "提示",
            {
              confirmButtonText: "确认",
              cancelButtonText: "取消",
              confirmButtonClass: "button-danger",
              type: "warning"
            }
          )
          .then(() => {
            (this as any).saveData();
          })
          .catch(() => {
            console.log("取消,什么都不做");
          });
      } else {
        (this as any).saveData();
      }
    },
    async saveData() {
      const data: any = (this as any).saveParams();
      if (data.spu_data.length === 0) {
        if (
          isNaN(data.active_two_level_id) &&
          !isNaN(data.active_one_level_id)
        ) {
          (this as any).$message.error("请选择带有二级渠道的附加链接");
          (this as any).isErrorInfo.additionalLink1Error = true;
          return false;
        }
        if (
          parseInt(data.active_one_level_id) > 0 &&
          parseInt(data.active_two_level_id) > 0 &&
          data.active_img === ""
        ) {
          (this as any).$message.error("保存失败，请检查填写内容");
          (this as any).isErrorInfo.active_img_error = true;
          return false;
        }
      } else {
        if (
          isNaN((this as any).outerLinkPage.additionalLink1[1]) &&
          !isNaN((this as any).outerLinkPage.additionalLink1[0])
        ) {
          (this as any).$message.error("请选择带有二级渠道的附加链接");
          (this as any).isErrorInfo.additionalLink1Error = true;
          return false;
        }
        if (
          isNaN((this as any).outerLinkPage.additionalLink2[1]) &&
          !isNaN((this as any).outerLinkPage.additionalLink2[0])
        ) {
          (this as any).$message.error("请选择带有二级渠道的附加链接");
          (this as any).isErrorInfo.additionalLink2Error = true;
          return false;
        }
      }
      try {
        const res: any = await API.saveFormData(data);
        if (res.code === 200) {
          (this as any).$message.success("保存成功");
          window.localStorage.removeItem("LOCALCHANNEL");
          (this as any).$router.push({
            name:
              "myapp_customer_auto-market_externalchain-manage_externalchain-list"
          });
        } else {
          (this as any).$message.warning(res.message);
        }
      } catch (error) {
        console.log(error);
        console.warn("获取渠道数据失败！");
      }
    },
    edit() {
      (this as any).$router.push({
        name:
          "myapp_customer_auto-market_externalchain-manage_externalchain-edit",
        query: {
          id: (this as any).editId,
          type: 1
        }
      });
    },
    filterChooseData() {
      if ((this as any).outerLinkPage.additionalLink1 == undefined)
        return false;
      const nowLink =
        (this as any).outerLinkPage.additionalLink1 &&
        (this as any).outerLinkPage.additionalLink1[1];
      const nowChoose = (this as any).additionalLink1LevelSecond.find(
        (item: any) => item.id === nowLink
      );
      if (nowChoose == undefined) return false;
      return nowChoose.url_type === 2;
    },
    goBackFn(b: string) {
      const { pageType } = this as any;
      if (b == "save") {
        const beforeSaveChainForm = cloneDeep((this as any).chainForm);
        const beforePassParams = cloneDeep((this as any).passParams);
        const beforeSaveOuterLinkPage = cloneDeep((this as any).outerLinkPage);
        const localChannel: object = {
          chainForm: beforeSaveChainForm,
          passParams: beforePassParams,
          outerLinkPage: beforeSaveOuterLinkPage
        };
        localStorage.setItem("LOCALCHANNEL", JSON.stringify(localChannel));
        (this as any).$router.push({
          name:
            "myapp_customer_auto-market_externalchain-manage_externalchain-list"
        });
      } else if (b == "cancel") {
        (this as any).$router.push({
          name:
            "myapp_customer_auto-market_externalchain-manage_externalchain-list"
        });
      }
      (this as any).isAlertSaveTip = false;
      (this as any).isPop = false;
    },
    /***
     * 复制页面的初始化数据
     */
    deepCloneFn() {
      (this as any).deepChainForm = cloneDeep((this as any).chainForm);
      (this as any).deepPassParams = cloneDeep((this as any).passParams);
      (this as any).deepOuterLinkPage = cloneDeep((this as any).outerLinkPage);
    },
    /***
     * 判断表单数据是否改变
     * @param {type}
     *  0: 新建 将更改的数据保存到本地
     *  1: 编辑 将更改的数据保存到数据库, 调用save()方法
     */
    deepCompare(type: number) {
      const { deepChainForm, deepPassParams, deepOuterLinkPage } = this as any;
      const beforeSaveChainForm = cloneDeep((this as any).chainForm);
      const beforePassParams = cloneDeep((this as any).passParams);
      const beforeSaveOuterLinkPage = cloneDeep((this as any).outerLinkPage);
      const pass =
        isEqual(deepChainForm, beforeSaveChainForm) &&
        isEqual(deepPassParams, beforePassParams) &&
        isEqual(deepOuterLinkPage, beforeSaveOuterLinkPage);
      if (pass) {
        (this as any).$router.push({
          name:
            "myapp_customer_auto-market_externalchain-manage_externalchain-list"
        });
      } else if (type === 0) {
        (this as any).isPop = true;
      } else {
        /***
         * 取消修改 或 调用save()方法,将修改的数据上传
         */
        (this as any)
          .$confirm(
            "你确认放弃编辑吗? 放弃后编辑的信息将不会被保存",
            "离开页面提示",
            {
              confirmButtonText: "确认",
              cancelButtonText: "取消",
              confirmButtonClass: "button-danger",
              type: "warning"
            }
          )
          .then(() => {
            console.log("确认放弃编辑");
            (this as any).$router.push({
              name:
                "myapp_customer_auto-market_externalchain-manage_externalchain-list"
            });
          })
          .catch(() => {
            console.log("取消, 继续编辑, 什么都不做");
          });
      }
    }
  },
  computed: {
    showBannerImage() {
      const nowData =
        (this as any).outerLinkPage.additionalLink1[0] > 0 &&
        (this as any).outerLinkPage.additionalLink1.length > 0;
      return (
        ((this as any).passParams.spu_data.length === 0 &&
          !(this as any).filterChooseData() &&
          (this as any).pageType !== 2 &&
          nowData) ||
        ((this as any).pageType === 2 &&
          (this as any).passParams.active_img !== "")
      );
    }
  }
};
</script>

<style lang="scss">
.welcome {
  .layout-form .footer {
    z-index: 10;
  }
  .el-table__body-wrapper::-webkit-scrollbar {
    width: 6px !important;
  }

  .el-form-item__label {
    height: 36px;
  }

  .el-input--small .el-input__inner {
    height: 36px;
    line-height: 36px;
  }

  .el-form-item {
    margin-bottom: 30px;
  }

  .cy-cascader {
    width: 400px !important;
  }

  .is-error-btn {
    border: 1px solid #f67170 !important;
  }

  .is-error-btn-color {
    .el-color-picker__trigger {
      border: 1px solid #f67170;

      span {
        border: none;
      }
    }
  }

  .is-error-btn-input {
    .el-input--small .el-input__inner {
      border: 1px solid #f67170;
    }
  }

  .is-error {
    .avatar-uploader .el-upload {
      border: 1px solid #f67170;
    }
  }

  .error-info {
    color: #f67170 !important;
    margin-top: 6px;
    font-size: 12px;
    line-height: 12px;
  }

  .el-form-item__label {
    padding: 0 20px 0 0;
  }
}
</style>
<style lang="scss" scoped>
div {
  box-sizing: border-box;
}

.welcome {
  .chain-header {
    margin-bottom: 30px;

    .chain-header-main {
      // margin: 30px 0;
      .short-chain-tip {
        height: 20px;
        line-height: 20px;
        font-size: 12px;
        color: #909399;
        margin-top: 10px;
      }
      .has-chain-tip {
        margin-bottom: -20px;
      }

      .short-chain-txt {
        width: 400px;
        line-height: 20px;
        padding: 20px;
        font-size: 14px;
        font-family: PingFangSC-Regular, PingFang SC;
        color: #606266;
        border-radius: 4px;
        background-color: #f5f6fa;
      }
    }
  }

  .chain-wrap {
    .create-wrap {
      display: flex;
      max-height: 1500px;
      transition: max-height 0.6s;

      .sms-wrap {
        flex: 0 0 375px;
        margin-right: 30px;

        .short-msg {
          // width: 56px;
          height: 20px;
          margin-bottom: 5px;
          font-size: 14px;
          color: #909399;
          line-height: 20px;
        }

        .sms-show {
          width: 375px;
          height: 667px;
          padding: 100px 0 50px 0px;
          background: url("https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/customer/phone.png") no-repeat 100%
            top;

          .sms-bg-wrap {
            width: 100%;
            height: 100%;
          }

          .time {
            text-align: center;
            line-height: 60px;
            color: $color-text-normal;
            font-size: 14px;
          }

          .clear-fix {
            overflow: hidden;
          }

          .content-msg {
            padding: 10px 20px;
            font-size: 12px;
            background: #ebebf1;
            margin: 0px 14px;
            line-height: 18px;
            border-radius: 36px;
            color: $color-text-normal;
            max-width: 245px;
            min-height: 18px;
            word-wrap: break-word;
            white-space: normal;
            word-break: break-all;
          }

          .content-msg::after {
            content: "\300c\9000\8ba2\56de\0054\300d";
          }

          .content-msg::before {
            content: attr(attr);
          }

          .short-card {
            width: 230px;
            // height: 191px;
            margin: 17px 14px 0px;
            overflow: hidden;
            border-radius: 15px;
            background: #f7f7f7;

            .card-brand-logo {
              display: block;
              width: 100%;
              height: 130px;
            }

            .card-title-wrap {
              padding: 10px;

              .card-title {
                width: 160px;
                height: 22px;
                line-height: 22px;
                font-size: 16px;
                font-weight: bold;
                color: #606266;
                white-space: nowrap;
                overflow: hidden;
                text-overflow: ellipsis;
              }

              .card-tip {
                width: 192px;
                height: 17px;
                line-height: 17px;
                font-size: 12px;
                color: #909493;
                white-space: nowrap;
                overflow: hidden;
                text-overflow: ellipsis;
              }
            }
          }
        }

        .short-msg-bottom {
          display: flex;
          justify-content: space-between;
          align-items: center;
          width: 375px;
          height: 30px;
          margin-top: 23px;

          .short-msg-watch {
            height: 20px;
            font-size: 14px;
            font-weight: bold;
            color: #606266;
          }

          .short-msg-img {
            width: 30px;
            height: 30px;
            cursor: pointer;
            border-radius: 100%;
            background: #6f3112
              url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAEAAAABACAQAAAAAYLlVAAAAfklEQVR42u2YoQ2AQBAE1+KhD4LBor4D+qCEt98QFkdH6EVBCxvITAMzya06CQAA/oVnNzfPKX31w5LV20dWb+9ZvV2y+ooePXr06NH/Ur9G9ZLPqF7yng4o4RPER0gCCSSQQAIJX0sY0wlb+knVS7mEK7CBN2Hw5E4AAPAlbqdNj8I27uGFAAAAAElFTkSuQmCC")
              no-repeat center;
            background-size: 20px auto;
          }

          .short-msg-img-right {
            width: 30px;
            height: 30px;
            background: #6f3112
              url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAEAAAABACAQAAAAAYLlVAAABPklEQVR42u3YzSpEYQCH8dNEE000mdEc5KOJJmVhpmxsZiFlY6O5DPdg7QZcgoWIRGJIpolII82VjAt4bBTmTFl5n835X8B5fpvz8Z4o+mPs06FHi83IGLt8rxE+X+T3QhOog0ogBptw4BOOZAIZjm3CECc2YZhTm5DlzCaMcG4TRrmwCTkubcIYVzZhnGubkOfGJkxwaxMK3NmESe5tQokHmxDTsgnTtG3CDI82YZYnmzDPs01Y4MUmlHm1CYt0EoAeayEJS7wlCM2wp4gK7wlCHJawTLcPULMBq24eSm6+6ebD3YYD8+EeRGk+zaf5NP+1HTe/7ea33PyGm6+7+XU3L7/vq25+xc1X3HzZzc+5+alwn9qHA/LFkH8+ksuHPGjVEvlc2INmvi+fjUKPvR/5TGSMBl0+aFP4n+t/AptcLlVuEXgUAAAAAElFTkSuQmCC")
              no-repeat center;
            background-size: 14px auto;
          }

          .short-msg-img-active {
            background-color: #c0c4cc;
            cursor: not-allowed;
          }
        }
      }

      .browser-wrap {
        flex: 0 0 375px;

        // margin-right: 30px;
        .browser-header {
          width: 375px;
          height: 93px;
          padding-top: 55px;
          text-align: center;
          background: url("https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/customer/phone_header.png")
            no-repeat;
          background-size: 375px 93px;
        }
        .browser-body {
          width: 100%;
          height: 575px;
          overflow: auto;
          position: relative;
        }
        .browser-body-box {
          position: absolute;
          top: 0;
          left: 0;
          width: 100%;
          height: 100%;
          background: rgba(0, 0, 0, 0.6);
          z-index: 1;
        }

        .browser-body {
          position: relative;
          height: 574px;
          padding: 0px 15px;
          padding-top: 203px;
          background: #aaa;

          .browser-btn-wrap {
            display: flex;
            justify-content: center;
            position: relative;
            z-index: 2;

            .browser-btn {
              display: flex;
              justify-content: center;
              align-items: center;
              width: 162px;
              height: 44px;
              border-radius: 4px;
            }
          }

          .browser-btn-tip {
            position: absolute;
            left: 15px;
            top: 258px;
            width: 345px;
            z-index: 2;
            line-height: 18px;
            text-align: center;
            font-size: 12px;
          }

          .one-spu-wrap {
            position: absolute;
            left: 15px;
            top: 422px;
            display: flex;
            width: 345px;
            height: 131px;
            padding: 15px 15px 0;
            border-radius: 14px;
            background: #fff;
          }

          .two-spu-wrap {
            position: absolute;
            left: 15px;
            bottom: 15px;
            z-index: 2;
            margin: 0 auto;
            width: 345px;
            padding: 15px;
            box-sizing: border-box;
            border-radius: 14px;
            background-color: #fff;

            .two-spu-item {
              display: flex;
              margin-bottom: 15px;

              &:last-child {
                margin-bottom: 0;
              }
            }
          }

          .banner-wrap {
            position: absolute;
            margin: 0 auto;
            bottom: 20px;
            z-index: 2;
            width: 345px;
            border-radius: 10px;
            background: #ebebf1;

            .banner-wrap-img {
              display: block;
              width: 100%;
              height: auto;
              border-radius: 10px;
            }
          }

          .one-spu-img {
            width: 100px;
            height: 100px;
            margin-right: 10px;
          }

          .one-spu-info {
            font-size: 12px;

            .one-spu-name {
              width: 205px;
              height: 40px;
              line-height: 20px;
              font-size: 14px;
              font-weight: bold;
              color: #333333;
              text-overflow: -o-ellipsis-lastline;
              overflow: hidden;
              text-overflow: ellipsis;
              display: -webkit-box;
              -webkit-line-clamp: 2;
              line-clamp: 2;
              -webkit-box-orient: vertical;
            }

            .one-spu-sku {
              width: 205px;
              height: 24px;
              line-height: 24px;
              margin: 4px 0px 12px 0;
              padding-left: 4px;
              color: #a6a6a6;
              background: #f7f7f7;
              text-overflow: -o-ellipsis-lastline;
              overflow: hidden;
              text-overflow: ellipsis;
              display: -webkit-box;
              -webkit-line-clamp: 1;
              line-clamp: 1;
              -webkit-box-orient: vertical;
            }

            .one-spu-price-wrap {
              height: 24px;
              line-height: 24px;
              color: #fa4046;

              .one-spu-price {
                font-size: 18px;
                font-weight: bold;
              }
            }
          }
        }
      }

      .content-right {
        flex: 1;

        // border-left: solid 1px #dcdfe6;
        // padding-left: 30px;
        /deep/.in-item {
          width: auto;
          margin-bottom: 20px;
        }

        /deep/ .label {
          font-size: 12px;
          padding-left: 0px;
        }

        .content-right-title {
          margin-bottom: 30px;
          color: #303133;
          font-size: 14px;
          font-weight: bold;
          line-height: 20px;
        }

        .switch-container {
          width: 514px;
          padding: 20px 0;
          background: #f5f6fa;
          border-radius: 3px;

          .switch-form-txt {
            display: flex;
            align-items: center;
          }

          .switch-form-txt-1 {
            margin-bottom: 15px;
          }
        }

        .add-link-btn {
          position: relative;
          display: inline-block;
          width: 12px;
          height: 12px;
          margin-left: 15px;
          cursor: pointer;
          border-radius: 100%;
          border: 1px solid #6f3112;

          .add-link-btn-1 {
            position: absolute;
            left: 0;
            right: 0;
            top: 0;
            bottom: 0;
            margin: auto;
            width: 6px;
            height: 1px;
            background-color: #6f3112;
          }

          .add-link-btn-2 {
            transform: rotate(90deg);
          }
        }

        .jump-tip-item {
          display: flex;
          align-items: center;
          height: 36px;
          margin-top: 15px;

          .item-span-color-1 {
            margin-right: 18px;
          }

          .item-btn-txt-1 {
            margin-left: 20px;
          }
        }

        .item-btn-txt {
          color: #6f3112;
          font-size: 12px;
        }

        .item-span-color {
          width: 48px;
          height: 15px;
          line-height: 15px;
          margin-right: 10px;
          font-size: 12px;
          color: #606266;
        }

        .img-desc {
          height: 14px;
          line-height: 14px;
          font-size: 12px;
          color: #909399;
        }

        .view-page-img {
          display: inline-block;
          width: 100px;
          height: 100px;
          border-radius: 4px;
        }
        .view-page-img2 {
          width: 100px;
          height: 20px;
        }
      }
    }

    .create-wrap-active {
      max-height: 0;
      overflow: hidden;
      transition: max-height 0.6s;
    }
  }

  .chain-header-bg {
    width: 100%;
    height: 40px;
    line-height: 40px;
    font-size: 18px;
    font-weight: bold;
    padding-left: 25px;
    margin-bottom: 30px;
    cursor: pointer;
    background: #ebeef5;
  }

  .chain-header-bg1 {
    position: relative;
    cursor: pointer;
    i {
      width: 9px;
      height: 40px;
      position: absolute;
      top: 0;
      right: 20px;
      display: block;
      background: #ebeef5
        url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABIAAAAUCAYAAACAl21KAAAAYUlEQVR4AWMYeiA2Nk0OhCk2JCEp7S4Ig9iUGvIfhGGGUWLICQhGGEaWIaFpafwgjG4YyYbA5LAZRrIh+Awj2hCSDItPTKvGbghhw0B6USQTktPzUQ0hbBhIDzUT/ygYBQB1/Ho9S897CwAAAABJRU5ErkJggg==")
        no-repeat 0px 16px;
      background-size: 9px auto;
      transition: transform 0.3s, -webkit-transform 0.3s;
    }
    span {
      font-size: 14px;
      color: #606266;
      line-height: 40px;
      font-weight: normal;
      position: absolute;
      top: 0;
      right: 34px;
    }
  }
  .chain-header-active {
    cursor: pointer;
    i {
      width: 9px;
      height: 40px;
      display: block;
      position: absolute;
      top: 0;
      right: 20px;
      transition: transform 0.3s, -webkit-transform 0.3s;
      background: #ebeef5
        url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABIAAAAUCAYAAACAl21KAAAAYUlEQVR4AWMYeiA2Nk0OhCk2JCEp7S4Ig9iUGvIfhGGGUWLICQhGGEaWIaFpafwgjG4YyYbA5LAZRrIh+Awj2hCSDItPTKvGbghhw0B6USQTktPzUQ0hbBhIDzUT/ygYBQB1/Ho9S897CwAAAABJRU5ErkJggg==")
        no-repeat 0px 16px;
      background-size: 9px auto;
      transform: translate(180deg);
    }
  }

  .hr {
    width: 100%;
    height: 1px;
    margin: 23px 0;
    background: #ebeef5;
  }
}
</style>
