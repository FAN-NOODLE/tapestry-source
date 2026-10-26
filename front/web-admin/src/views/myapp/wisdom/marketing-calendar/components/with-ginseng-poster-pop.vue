<template>
  <!-- 添加小程序带参海报弹窗 -->
  <div class="with-ginseng-poster-pop">
    <el-dialog
      width="1000px"
      :title="dialog_name"
      :visible.sync="dialogVisible"
      :close-on-click-modal="false"
      @close="closeDialog"
    >
      <div class="content-wrap">
        <div class="left">
          <div class="applet-page">
            <div
              class="poster-bacImg-wrap"
              v-if="resource_url"
              @mouseover.prevent="is_hover = true"
              @mouseleave.prevent="is_hover = false"
            >
              <img class="poster-bacImg" :src="resource_url" alt="" />
              <div class="mask" v-show="is_hover">
                <i
                  class="el-icon-zoom-in"
                  @click.prevent="preview_show = true"
                  aria-hidden="true"
                ></i>
                <i
                  class="el-icon-delete has-preview"
                  @click.prevent="remove"
                  aria-hidden="true"
                ></i>
              </div>
            </div>
            <div class="upload-img-wrap" v-else>
              <cy-upload
                class="upload-template"
                accept="image/png, image/jpeg"
                type="image"
                :on-success="onSuccessFn"
                :beforeUpload="beforeUploadFn"
                :onRemove="removeFn"
                :onError="errorFn"
                :action="newUploadImg"
              >
                <div>
                  <img
                    class="tip-img"
                    src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/marketing-calendar/img_upload.png"
                    alt=""
                  />
                  <div class="title">上传海报背景图片</div>
                  <div class="img-tip">
                    请上传750*1220px尺寸的图片，格式jpg、png；上传前请仔细查看作图示例
                  </div>
                </div>
              </cy-upload>
            </div>
            <img
              class="qr-code"
              src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/marketing-calendar/img_qr.png"
              alt=""
            />
          </div>
          <div class="tip cursor" @click="checkSampleGraph">
            查看作图示例<i class="el-icon-arrow-right"></i>
          </div>
        </div>
        <div class="right">
          <div>
            <div class="title">小程序带参海报</div>
            <div>
              <div class="item-title">商城小程序名称</div>
              <div>{{ app_name }}</div>
            </div>
            <div>
              <div class="item-title">商城APPID</div>
              <div>{{ shopId }}</div>
            </div>
            <div class="url-wrap">
              <div class="item-title"><span>*</span>识别小程序码跳转路径</div>
              <div>
                <div class="flex-align-center">
                  <cy-input
                    size="small"
                    width-size="medium"
                    placeholder="请输入识别海报小程序码后用户的跳转页面路径"
                    v-model.trim="applet_path"
                    class="input-with-select"
                    @change="appletPathChange"
                  ></cy-input>
                  <span class="get-sun-code-btn" @click="getSunCode"
                    >生成预览码</span
                  >
                </div>
                <div v-if="sunCodeWrongTip" class="tip margin-tip">
                  {{ sunCodeWrongContent }}
                </div>
                <div class="tip m-b-5" :class="!sunCodeWrongTip ? 'm-t-7' : ''">
                  务必填写当前企业商城小程序的有效页面地址
                </div>
                <div class="tip">
                  建议填写：商品详情页、活动页、直播页、一码多券口令券、商城个人中心页面地址
                </div>
                <!-- 太阳码 -->
                <div
                  v-if="sunCodeUrl && !sunCodeWrongTip"
                  class="flex-align-center sun-code-wrap"
                >
                  <img class="sun-code-img" :src="sunCodeUrl" alt="" />

                  <div>
                    <div class="flex-align-center">
                      <i class="sun-code-icon el-icon-success"></i>
                      <span class="sun-code-tip">生成成功</span>
                    </div>
                    <div class="tip">预览码不带导购参，切勿直接使用</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div>
            <div class="title appletTitle">小程序卡片</div>
            <div class="tip">
              非必填项；设置此项后，营销日历添加在导购聊天工具栏的情况下，导购可以从聊天工具栏以链接形式发送小程序页面
            </div>
            <div>
              <div class="item-title applet-card-title">小程序卡片标题</div>
              <cy-input
                size="small"
                width-size="medium"
                placeholder="请输入小程序卡片标题"
                v-model.trim="applet_card_title"
                class="input-with-select"
                maxlength="20"
                show-word-limit
              ></cy-input>
            </div>
            <div>
              <div class="item-title applet-card-title">小程序卡片封面</div>
              <div>
                <upload-image
                  class="upload-template"
                  :urls="applet_card_cover_url"
                  accept="image/jpeg, image/png"
                  :before-upload="beforeUpload"
                  @change="update"
                  :action="newUploadImg"
                >
                </upload-image>
                <div class="img-tip">
                  请上传分辨率为400 × 320的图片，支持jpg、png格式，大小500kb以下
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div slot="footer" class="dialog-footer">
        <el-button size="small" @click="closeDialog">取 消</el-button>
        <el-button class="comfirm-btn" size="small" @click="save"
          >确认</el-button
        >
      </div>
    </el-dialog>
    <!-- 小程序带参海报实例图 -->
    <el-dialog
      width="1000px"
      :visible.sync="innerVisible"
      :close-on-click-modal="false"
      class="sample-graph-wrap"
    >
      <div slot="title" class="title-wrap">
        <div class="last-step cursor" @click="closeSampleGraph">
          添加小程序带参海报
        </div>
        <i class="el-icon-arrow-right"></i>
        <div class="title">查看作图示例</div>
      </div>
      <div>
        <img
          src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/marketing-calendar/sample-graph-image.png"
          alt=""
        />
      </div>
      <div slot="footer">
        <el-button size="small" @click="closeSampleGraph">返 回</el-button>
      </div>
    </el-dialog>

    <preview
      v-if="preview_show"
      :urls="[resource_url]"
      @close="preview_show = false"
    ></preview>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Emit, Mixins } from "vue-property-decorator";
// 上传图片
import UploadImage from "@/components/business/upload-image/index.vue";
import UploadValidate from "@/mixins/UploadValidate";
import preview from "@/components/business/preview/index.vue";
import { miniProgram, sunCode, ISunCode } from "@/api/wisdom/marketingPlan";
import { IUploadImageBeforeParam } from "@/api/public/public-interface";
import { VUE_APP_API_ROOT } from "@/utils/env-config";

@Component({
  components: {
    UploadImage,
    preview,
  },
})
export default class BaseSetting extends Mixins(UploadValidate) {
  @Prop({ default: false }) visible?: boolean; // 控制模态框显隐
  @Prop({ default: {} }) appletPosterMess?: object;
  private resource_url = ""; // 海报背景图
  private dialog_name = "添加小程序页面"; // 模态框名称
  private dialogVisible = true; // 添加小程序带参海报模态框显示/隐藏
  private app_name = "三星体"; // 小程序名称
  private shopId = "112345678"; // 商城id
  private applet_card_title = ""; // 小程序卡片标题
  private applet_path = ""; // 识别海报小程序码后用户的跳转页面路径
  private temPath = ""; // 生成太阳码的小程序跳转路径
  private applet_card_cover_url = ""; // 小程序卡片封面
  private appletImg = "";
  private testFlag = false;
  private innerVisible = false; // 小程序带参海报示例图弹窗 显示/隐藏
  private preview_show = false;
  private is_hover = false;
  // 上传图片
  private newUploadImg = `${VUE_APP_API_ROOT}common/upload-img`;
  private sunCodeUrl = ""; // 生成太阳码路径
  private sunCodeWrongTip = false; // 生成预览码失败flag
  private sunCodeWrongContent = "请填入正确的小程序路径链接后重新生成";
  private saveSunCodeFlag = false; // 生成太阳码失败保存校验太阳码falg
  /* 生命周期 */
  private mounted() {
    this.init();
    // 从新建页面的小程序素材点击进来
    if (
      this.appletPosterMess &&
      (this.appletPosterMess as object)["resource_url"]
    ) {
      const {
        resource_url,
        app_name,
        applet_card_title,
        applet_card_cover_url,
        applet_path,
        shopId,
        sun_code,
      } = this.appletPosterMess as any;
      this.resource_url = resource_url;
      this.app_name = app_name;
      this.applet_card_title = applet_card_title;
      this.applet_card_cover_url = applet_card_cover_url;
      this.applet_path = applet_path;
      this.temPath = applet_path;
      this.shopId = shopId;
      this.sunCodeUrl = sun_code;
      this.saveSunCodeFlag = sun_code ? true : false;
    }
  }
  // 关闭模态框的回调
  @Emit("close")
  private closeDialog() {
    this.dialogVisible = false;
  }
  // 保存
  @Emit("save")
  private save() {
    if (!this.resource_url || !this.applet_path) {
      const message = !this.resource_url
        ? "请上传海报背景图片"
        : "请输入跳转页面路径";
      this.$message.error(message);
      return;
    }
    if (!this.saveSunCodeFlag) {
      this.sunCodeWrongContent = "请生成小程序预览码";
      this.sunCodeWrongTip = true;
      this.$message.error("请生成小程序预览码");
      return;
    }
    this.dialogVisible = false;
    const appletMess = {
      resource_url: this.resource_url,
      app_name: this.app_name,
      applet_card_title: this.applet_card_title,
      applet_card_cover_url: this.applet_card_cover_url,
      applet_path: this.applet_path,
      shopId: this.shopId,
      sun_code: this.sunCodeUrl,
    };
    return appletMess;
  }
  // 小程序路径变化
  private appletPathChange() {
    // 如果修改了小程序路径，清空太阳码
    if (this.applet_path != this.temPath) {
      this.sunCodeWrongTip = false;
      this.saveSunCodeFlag = false;
      this.sunCodeUrl = "";
    }
  }
  // 查看海报示例图
  private checkSampleGraph() {
    this.innerVisible = true;
  }
  // 关闭海报示例图
  private closeSampleGraph() {
    this.innerVisible = false;
  }
  // 小程序卡片封面
  // 上传文件限制规范
  private beforeUpload(file: File) {
    this.testFlag = false;
    const obj: any = {
      ext: "",
      name: "",
      size: "",
      url: "",
    };
    obj.name = file.name;
    obj.ext = file.type;
    obj.size = file.size;
    this.appletImg = obj;

    const _URL = window.URL;
    const img = new Image();
    img.src = _URL.createObjectURL(file);
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
    if (file.size > 0.5 * 1024 * 1024) {
      this.$message.error("请上传500kb大小以内的图片");
      return false;
    }
    img.src = _URL.createObjectURL(file);
    const is_size = new Promise((resolve, reject) => {
      const img = new Image();
      img.src = window.URL.createObjectURL(file);
      img.onload = () => {
        // 限制宽高必须为 4:3 比例
        // const valid = img.width / img.height == 4 / 3;
        const valid = img.width == 400 && img.height == 320;
        if (valid) {
          resolve(file);
        } else {
          this.$message.error("请上传分辨率为400 × 320的图片");
          reject();
        }
      };
    });
    return is_size;
  }
  // 上传
  private update(res: string) {
    this.testFlag = true;
    this.applet_card_cover_url = res;
  }
  // 海报背景图
  // 上传文件前
  private beforeUploadFn(file: File) {
    const width = 750;
    const height = 1220;
    const _URL = window.URL;
    const img = new Image();
    const fileType = file
      ? file.name.substring(file.name.lastIndexOf(".") + 1)
      : "";
    if (
      file.type != "image/jpg" &&
      file.type != "image/png" &&
      fileType != "png" &&
      fileType != "jpg"
    ) {
      this.$message.error("请上传JPG、PNG格式的图片");
      return false;
    }
    if (file.size > 0.5 * 1024 * 1024) {
      this.$message.error("请上传500KB大小以内的图片");
      return false;
    }
    img.src = _URL.createObjectURL(file);
    const is_size = new Promise((resolve, reject) => {
      const img = new Image();
      img.src = window.URL.createObjectURL(file);
      img.onload = () => {
        // 限制宽高必须为 750*1220 像素
        const valid = img.width == width && img.height == height;
        // 限制宽高必须为 n:n 比例
        // const valid = img.width == img.height;
        if (valid) {
          resolve(file);
        } else {
          this.$message.error("请上传尺寸为750x1220的图片");
          reject();
        }
      };
    });
    return is_size;
  }
  // 上传错误
  private errorFn() {
    console.log("报错");
  }
  //上传成功
  private onSuccessFn(res: any) {
    const { url } = res.data;
    this.resource_url = url;
  }
  //移除上传的文件
  private removeFn() {
    console.log("移除上传的文件");
  }
  /* 
    图片上传尺寸校验
    file: 图片信息
  */
  checkImageSize(file: any) {
    const width = 750;
    const height = 1220;
    const _URL = window.URL;
    const img = new Image();
    img.src = _URL.createObjectURL(file);
    const size = new Promise((resolve, reject) => {
      const img = new Image();
      img.src = window.URL.createObjectURL(file);
      img.onload = () => {
        // 限制宽高必须为 750*1220 像素
        const valid = img.width == width && img.height == height;
        // 限制宽高必须为 n:n 比例
        // const valid = img.width == img.height;
        if (valid) {
          resolve(file);
        } else {
          reject();
        }
      };
    });
    return size;
  }
  // 移除
  private remove() {
    this.resource_url = "";
  }
  private init() {
    miniProgram()
      .then((res: any) => {
        if (res.code == 200 && res.data) {
          this.app_name = res.data[0].name;
          this.shopId = res.data[0].app_id;
        }
      })
      .catch(() => {
        console.log("miniProgram");
      });
  }
  private getSunCode() {
    if (!this.applet_path) {
      this.$message.error("请输入跳转页面路径");
      return;
    }
    const params: ISunCode = {
      applet_path: this.applet_path, // 小程序路径
    };
    this.temPath = this.applet_path;
    sunCode(params)
      .then((res: any) => {
        if (res.code == 200) {
          this.sunCodeUrl = res.data && res.data.url ? res.data.url : "";
          this.sunCodeWrongTip = false;
          this.saveSunCodeFlag = true;
        }
      })
      .catch((err: any) => {
        console.log("生成太阳码sunCode:", err);
        this.sunCodeWrongContent = "请填入正确的小程序路径链接后重新生成";
        this.sunCodeWrongTip = true;
        this.saveSunCodeFlag = false;
      });
  }
}
</script>

<style lang="scss" scoped>
/deep/.el-dialog {
  height: 620px;
}
.content-wrap {
  height: 450px;
  padding-top: 10px;
  display: flex;
  .left {
    width: 255px;
    margin-right: 30px;
    text-align: center;
    .applet-page > div {
      width: 100%;
      height: 100%;
    }
    .applet-page {
      height: 415px;
      background: #eeeff3;
      position: relative;
      .upload-img-wrap {
        padding: 0 20px;
        box-sizing: border-box;
      }
      .poster-bacImg-wrap {
        position: relative;
        .mask {
          position: absolute;
          top: 0;
          left: 0;
          right: 0;
          bottom: 0;
          background-color: rgba(0, 0, 0, 0.5);
        }
        .el-icon-zoom-in {
          cursor: pointer;
          color: white;
          font-size: 20px;
          position: absolute;
          top: 50%;
          left: 35%;
          line-height: var(--height) !important;
        }
        .has-preview {
          cursor: pointer;
          color: white;
          font-size: 20px;
          position: absolute;
          top: 50%;
          right: 35%;
          line-height: var(--height) !important;
        }
      }
      .poster-bacImg {
        width: 100%;
      }
      .tip-img {
        margin-top: 110px;
        width: 54px;
        height: 34px;
      }
      .title {
        font-size: 14px;
        font-weight: 400;
        color: #606266;
        line-height: 20px;
      }
      .img-tip {
        font-size: 12px;
        font-family: PingFangSC-Regular, PingFang SC;
        font-weight: 400;
        color: #909399;
        line-height: 18px;
        margin-top: 10px;
      }
      .qr-code {
        width: 68px;
        height: 84px;
        position: absolute;
        bottom: 15px;
        right: 15px;
      }
    }

    .tip {
      font-size: 12px;
      font-weight: 400;
      color: #6b8bfd;
      line-height: 18px;
      margin-top: 15px;
    }
  }
  .right > div > div {
    display: flex;
    margin-bottom: 30px;
    font-size: 14px;
    font-family: PingFangSC-Regular, PingFang SC;
    font-weight: 400;
    color: #303133;
    line-height: 20px;
  }
  .right {
    width: 655px;
    overflow: auto;
    .title {
      padding: 10px 20px;
      background: #ebeef5;
      font-size: 14px;
      font-family: PingFangSC-Medium, PingFang SC;
      font-weight: 500;
      color: #606266;
      line-height: 20px;
      margin-bottom: 20px;
    }
    .appletTitle {
      margin-bottom: 5px;
    }
    .item-title {
      text-align: right;
      width: 118px;
      color: #606266;
      margin-right: 20px;
      span {
        color: #f67270;
      }
    }
    .tip {
      font-size: 12px;
      font-family: PingFangSC-Regular, PingFang SC;
      font-weight: 400;
      color: #909399;
      line-height: 20px;
    }
    .url-wrap {
      display: flex;
      .get-sun-code-btn {
        cursor: pointer;
        font-size: 14px;
        color: #708df5;
        margin-left: 20px;
      }
    }
    .margin-tip {
      margin: 7px 0 5px;
      color: #f67270;
    }
    .applet-card-title {
      width: 98px;
    }
    .img-tip {
      font-size: 12px;
      font-family: PingFangSC-Regular, PingFang SC;
      font-weight: 400;
      color: #909399;
      line-height: 18px;
      margin-top: 5px;
    }
  }
}
.comfirm-btn {
  color: #fff;
  background: #708df5;
}
.sample-graph-wrap {
  .title-wrap {
    display: flex;
    align-items: center;
    .last-step {
      font-size: 14px;
      color: #606266;
      line-height: 20px;
    }
    i {
      color: #606266;
      font-size: 12px;
      margin: 0 5px;
    }
    .title {
      font-size: 18px;
      font-family: PingFangSC-Semibold, PingFang SC;
      font-weight: 600;
      color: #292a30;
      line-height: 26px;
    }
  }
  img {
    width: 100%;
  }
}
.post-img-wrap {
  width: 100%;
  height: 100%;
  opacity: 0.2;
}
.preview-code-wrap {
  display: flex;
  align-items: center;
}
.sun-code-wrap {
  margin-top: 20px;
  .sun-code-img {
    width: 110px;
    height: 110px;
    margin-right: 20px;
  }
  .sun-code-icon {
    font-size: 16px;
    color: #48ba9b;
    margin-right: 10px;
  }
  .sun-code-tip {
    color: #48ba9b;
    font-size: 14px;
    margin-bottom: 5px;
  }
}
.flex-align-center {
  display: flex;
  align-items: center;
}
.cursor {
  cursor: pointer;
}
.p-l-20 {
  padding-left: 20px;
}
.p-r-20 {
  padding-right: 20px;
}
.m-b-5 {
  margin-bottom: 5px;
}
.m-t-7 {
  margin-top: 7px;
}
</style>
