<template>
  <div>
    <el-dialog
      title="企微群发"
      :visible.sync="dialogVisible"
      :close-on-click-modal="false"
      @close="closeDialogFn"
    >
      <div class="container">
        <el-form ref="form" :model="form" :rules="rules" label-width="90px">
          <el-form-item label="消息内容" required prop="message">
            <div class="publish-content-box">
              <!-- 表情、文本选择区域-->
              <div class="select-box-insert">
                <add-emoji @select="addEmojiFn" class="add-emoji">
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
                  @click="openTextTemplateFn"
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
                  v-model.trim="form.message"
                  :rows="5"
                  type="textarea"
                  maxlength="500"
                  show-word-limit
                >
                </el-input>
              </div>
              <!-- 图片、视频、链接、小程序选择区域-->
              <div class="material-select-box">
                <el-popover
                  placement="top-start"
                  width="180"
                  trigger="hover"
                  v-if="!info.content.fileData"
                >
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
                    slot="reference"
                    class="select-material-btn"
                    size="small"
                    ><img
                      class="material-btn-img"
                      src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/add_icon.png"
                      alt=""
                    /><span>图片/链接/小程序</span></el-button
                  >
                </el-popover>
                <div
                  class="select-picture-box"
                  v-if="info.content.fileData && info.content.fileType == 1"
                >
                  <div class="item-picture">
                    <img
                      class="item-img"
                      :src="info.content.fileData.url"
                      alt=""
                    />
                    <div class="delete-img-box" @click="deleteMaterial">
                      <img
                        src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-delete-icon.png"
                        alt=""
                      />
                    </div>
                  </div>
                </div>

                <div
                  class="select-link-box"
                  v-if="info.content.fileData && info.content.fileType == 2"
                >
                  <div class="select-link-title">
                    {{ info.content.fileData.title }}
                  </div>
                  <div class="select-link-bottom">
                    <!-- 为了解决flex布局下超出省略失效问题，所以嵌套-->
                    <div>
                      <span class="select-link-description">{{
                        info.content.fileData.description
                      }}</span>
                    </div>
                    <img
                      class="select-link-cover"
                      :src="info.content.fileData.cover_url"
                      alt=""
                    />
                  </div>
                  <div class="delete-link-box" @click="deleteMaterial">
                    <img
                      src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-delete-icon.png"
                      alt=""
                    />
                  </div>
                </div>

                <div
                  class="select-applet-box"
                  v-if="info.content.fileData && info.content.fileType == 3"
                >
                  <div class="select-applet-name">
                    <img
                      class="minilogo"
                      src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/mini-applet-icon.png"
                      alt=""
                    />

                    <span class="smallName">{{
                      info.content.fileData.app_name
                    }}</span>
                  </div>
                  <div class="select-applet-title">
                    {{ info.content.fileData.title }}
                  </div>
                  <img
                    class="select-applet-img"
                    :src="info.content.fileData.cover_url"
                    alt=""
                  />
                  <div class="delete-applet-box" @click="deleteMaterial">
                    <img
                      src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-delete-icon.png"
                      alt=""
                    />
                  </div>
                </div>
              </div>
            </div>
          </el-form-item>
        </el-form>
      </div>
      <span slot="footer" class="dialog-footer">
        <el-button
          type="primary"
          class="marginL10"
          @click="submitFn"
          size="small"
          >确 定
        </el-button>
        <el-button @click="closeDialogFn" size="small">取 消</el-button>
      </span>
    </el-dialog>
    <!-- 选择文本素材弹窗-->
    <select-text-material
      v-if="textDialog"
      @textData="getTextData"
    ></select-text-material>
    <!-- 选择文本素材弹窗-->
    <select-text-material
      v-if="textDialog"
      @textData="getTextData"
    ></select-text-material>
    <!-- 选择图片素材弹窗-->
    <select-picture-material
      v-if="pictureDialog"
      :limiteNum="1"
      @submit="getFileData"
    ></select-picture-material>
    <!-- 选择小程序素材弹窗-->
    <select-applet-material
      v-if="appletDialog"
      @submit="getFileData"
    ></select-applet-material>
    <!-- 选择链接素材弹窗-->
    <select-link-material
      v-if="urlDialog"
      @submit="getFileData"
    ></select-link-material>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import addEmoji from "@/views/myapp/enterprise/components/add-emoji.vue";
// 选择文本素材弹窗
import selectTextMaterial from "@/views/myapp/enterprise/material-library/components/text-material-dialog.vue";
// 选择图片素材弹窗
import selectPictureMaterial from "@/views/myapp/enterprise/material-library/components/pic-mterial-dialog.vue";
// 选择小程序素材弹窗
import selectAppletMaterial from "@/views/myapp/enterprise/material-library/components/small-mterial-dialog.vue";
// 选择链接素材弹窗
import selectLinkMaterial from "@/views/myapp/enterprise/material-library/components/link-material-dialog.vue";
@Component({
  components: {
    addEmoji,
    selectTextMaterial,
    selectPictureMaterial,
    selectAppletMaterial,
    selectLinkMaterial
  }
})
export default class extends Vue {
  @Prop() msgToDialog!: any;
  /************data**************/
  private dialogVisible = true;
  private pictureDialog = false;
  private urlDialog = false;
  private appletDialog = false;
  private form = {
    message: this.msgToDialog.content || ""
  };
  private info: any = {
    member_ids: [],
    content: {
      text: "",
      fileData: null,
      fileType: 0
    }
  };
  private typeName = "";
  private rules = {
    message: [
      { required: true, message: "请输入消息内容", trigger: "blur" },
      { min: 1, max: 500, message: "长度在 1 到 500 个字符", trigger: "blur" }
    ]
  };
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
    }
  ];
  private textDialog = false;
  /************method**************/
  // 添加表情
  async addEmojiFn(emj: any) {
    const myField = (this.$refs.textarea as any).$refs.textarea;
    if (myField.selectionStart || myField.selectionStart === 0) {
      const startPos = myField.selectionStart;
      const endPos = myField.selectionEnd;
      this.form.message =
        myField.value.substring(0, startPos) +
        emj +
        myField.value.substring(endPos, myField.value.length);
      await this.$nextTick();
      myField.focus();
      myField.setSelectionRange(endPos + emj.length, endPos + emj.length);
    } else {
      this.form.message += emj;
    }
  }
  openTextTemplateFn() {
    this.textDialog = true;
  }
  getTextData(data: any) {
    this.textDialog = false;
    if (data) {
      this.addEmojiFn(this.getDataFromArrayOrObject(data).content);
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
  async submitFn() {
    (this.$refs["form"] as any).validate((valid: boolean) => {
      if (valid) {
        const message: object = {
          name: this.form.message,
          type: 1,
          typeName: this.typeName == "" ? "文本" : this.typeName
        };
        this.$emit("replaceWeChatMsgDialog", message);
      } else {
        return false;
      }
    });
  }
  closeDialogFn() {
    this.$emit("replaceWeChatMsgDialog");
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
  getPictureData(data: any) {
    this.pictureDialog = false;
    if (data) {
      this.info.content.fileData = this.getDataFromArrayOrObject(data);
    }
    console.log(data);
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
    switch (this.info.content.fileType) {
      case 1:
        this.getPictureData(data);
        break;
      case 2:
        this.getUrlData(data);
        break;
      case 3:
        this.getAppletData(data);
        break;
      default:
        break;
    }
    if (!data) this.info.content.fileType = 0;
  }
  // 鼠标划出事件
  private outStyle(item: any) {
    this.$set(item, "active", false);
  }
  deleteMaterial() {
    this.info.content.fileType = 0;
    this.info.content.fileData = null;
  }
  private openDialog(dialog: string, type: number) {
    this[dialog] = true;
    this.info.content.fileType = type;
    switch (type) {
      case 1:
        this.typeName = "图片";
        break;
      case 2:
        this.typeName = "链接";
        break;
      case 3:
        this.typeName = "小程序";
        break;
    }
    console.log(dialog);
  }
}
</script>
<style lang="scss" scoped>
.select-material-list {
  width: 180px;
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
.publish-content-box {
  width: 100%;
  border: solid 1px #dcdfe6;
  border-radius: 4px;
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
}
.marginL10 {
  margin-left: 10px;
}
</style>
