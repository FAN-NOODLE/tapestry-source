<template>
  <div>
    <el-dialog
      title="企微群发"
      :visible.sync="dialogVisible"
      :close-on-click-modal="false"
      @close="closeDialogFn"
    >
      <div class="container">
        <el-form ref="form" :model="form" label-width="90px">
          <el-form-item label="消息内容" prop="message">
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
                  v-model="form.message"
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
                  v-if="info.content.fileType == 0"
                >
                  <div class="select-material-list">
                    <div
                      class="item-select"
                      v-for="(item, index) in materialSelectList"
                      @mouseover="selectStyle(item)"
                      @mouseout="outStyle(item)"
                      @click="openDialog(item)"
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
                  v-if="info.content.fileData && info.content.fileType == 2"
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
                  v-if="info.content.fileData && info.content.fileType == 3"
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
                  v-if="info.content.fileData && info.content.fileType == 4"
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
        <el-button @click="closeDialogFn" size="small">取 消</el-button>
        <el-button
          type="primary"
          class="marginL10"
          @click="submitFn"
          size="small"
          >确 认
        </el-button>
      </span>
    </el-dialog>
    <!-- 选择文本素材弹窗-->
    <select-text-material
      v-if="textDialog"
      @textData="getTextData"
    ></select-text-material>
    <!-- 选择图片素材弹窗-->
    <select-picture-material
      v-if="pictureDialog"
      :limiteNum="1"
      @submit="getPictureData"
    ></select-picture-material>
    <!-- 选择小程序素材弹窗-->
    <select-applet-material
      v-if="appletDialog"
      @submit="getAppletData"
    ></select-applet-material>
    <!-- 选择链接素材弹窗-->
    <select-link-material
      v-if="urlDialog"
      @submit="getUrlData"
    ></select-link-material>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import {
  ITextData,
  IPictureData,
  ILinkData,
  IAppData,
  IAllData
} from "@/api/customer/select-wechat-msg";
// 添加表情
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
  /*********************************** data ************************************/
  private dialogVisible = true; // 企微弹窗
  // 文本内容
  private form = {
    message: ""
  };
  // 文件类型名 图片/链接/小程序
  private typeName = "";
  private textDialog = false; // 文本素材库弹窗
  private pictureDialog = false; // 图片素材库弹窗
  private urlDialog = false; // 链接素材库弹窗
  private appletDialog = false; // 小程序素材库弹窗
  // 图片/链接/小程序内容
  private info: any = {
    content: {
      text: "", // 图片/链接/小程序名称
      fileData: null, // 原始数据
      fileType: 0 // 图片：2/链接：3/小程序：4
    }
  };
  // 图片/链接/小程序素材
  private materialSelectList: Array<any> = [
    {
      type: 2,
      hoverIcon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-img-icon.png",
      icon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-img-hover-icon.png",
      tip: "图片",
      dialog: "pictureDialog"
    },
    {
      type: 3,
      hoverIcon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-link-icon.png",
      icon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-link-hover-icon.png",
      tip: "链接",
      dialog: "urlDialog"
    },
    {
      type: 4,
      hoverIcon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-mini-icon.png",
      icon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-mini-hover-icon.png",
      tip: "小程序",
      dialog: "appletDialog"
    }
  ];
  /*********************************** 生命周期 ************************************/
  private mounted() {
    // 初始化数据
    this.getFormInfo(this.msgToDialog);
  }
  /*********************************** method ************************************/
  // 添加表情
  async addEmojiFn(emj: string) {
    console.log(emj, "emj", typeof emj);
    const myField = (this.$refs.textarea as any).$refs.textarea;
    // myField.selectionStart是当前光标的位置 允许在任何位置插入
    if (myField.selectionStart || myField.selectionStart === 0) {
      // startPos 开始位置 endPos 结束位置
      const startPos = myField.selectionStart;
      const endPos = myField.selectionEnd;
      this.form.message =
        myField.value.substring(0, startPos) + // myField.value.substring(0, startPos) 这里是文本框内容的 显示0-startPos之间的内容
        emj + // emj为传入的文本内容或者表情
        myField.value.substring(endPos, myField.value.length); // myField.value.substring(endPos, myField.value.length) 这里是在结束位置插入的内容 显示endPos-文本总长度之间的内容
      await this.$nextTick(); // 获取更新后的DOM
      myField.focus(); // 执行内置方法 获得焦点
      myField.setSelectionRange(endPos + emj.length, endPos + emj.length); // 设定当前选中文本的起始和结束位置。
    } else {
      this.form.message += emj; // 开始位置在最后一位时，执行拼接
    }
    // 超出文本框最大字符数进行截取
    if (this.form.message.length > 500) {
      this.form.message = this.form.message.slice(0, 500);
      this.$message.info("消息内容长度不能超过500字符，超出部分则被截断。");
    }
  }
  // 打开文本弹框
  openTextTemplateFn() {
    this.textDialog = true;
  }
  // 打开图片/链接/小程序弹框
  private openDialog(originData: object) {
    this[originData["dialog"]] = true; // 打开指定的素材库弹窗
    this.typeName = originData["type"]; // 储存类型(图片/链接/小程序)
    this.info.content.fileType = originData["type"]; // 储存type值
  }
  // 获取图片/链接/小程序数据
  // getFileData(data: any) {
  //   console.log(data);
  //   // 这里传入的data是从素材库组件中传入的，是获取到了图片/链接/小程序的值
  //   // 判断this.info.content.fileType，就是判断type值，可以显示指定的样式并获取数据
  //   switch (this.info.content.fileType) {
  //     case 2:
  //       this.getPictureData(data);
  //       break;
  //     case 3:
  //       this.getUrlData(data);
  //       break;
  //     case 4:
  //       this.getAppletData(data);
  //       break;
  //     default:
  //       break;
  //   }
  //   if (!data) this.info.content.fileType = 0; // 当data值为空时 置空type值 即不显示企微弹框内图片/链接/小程序的样式
  // }

  // 鼠标划入事件(滑动效果)
  private selectStyle(item: object) {
    const _this = this;
    this.$nextTick(function() {
      this.materialSelectList.forEach(function(item: any) {
        _this.$set(item, "active", false);
      });
      this.$set(item, "active", true);
    });
  }
  // 鼠标划出事件(滑动效果)
  private outStyle(item: object) {
    this.$set(item, "active", false);
  }
  // 数据结构筛选判断
  getDataFromArrayOrObject(data: IAllData) {
    // 判断数据类型
    // Object.prototype.toString.call(data) 这里可以判断data对象到底属于那种内置类型
    // 这里存在的类型只有两种数组以及自定义类型
    if (Object.prototype.toString.call(data) == "[object Array]") {
      return data[0]; // 数组中获取数组下标为零的对象 当前功能只能做单选 所以选择第零位
    } else if (Object.prototype.toString.call(data) == "[object Object]") {
      return data; // 自定义类型则返回当前的data值
    } else {
      return null; // 其他类型返回一律返回null值
    }
  }
  /*******************
   * 获取数据
   * 功能：执行此函数时，首先会关闭指定的（文本/图片/链接/小程序）弹窗
   * 当存在返回的原始数据时，将数据进行判断，储存所需要的数据
   * 在这里需要文本/图片/链接/小程序的原始数据即接口返回的数据，图片/链接/小程序的名称标题内容
   *******************/
  // 获取文本内容
  getTextData(data: ITextData) {
    // data为获得的文本的原始内容
    this.textDialog = false;
    if (data) {
      this.addEmojiFn(this.getDataFromArrayOrObject(data).content);
    }
  }
  // 获取图片内容
  getPictureData(data: Array<IPictureData>) {
    // data为获得的图片的原始内容
    this.pictureDialog = false;
    if (data) {
      this.info.content.fileData = this.getDataFromArrayOrObject(data);
      this.info.content.text = this.getDataFromArrayOrObject(data).name;
    } else {
      this.info.content.fileType = 0;
    }
  }
  // 获取链接内容
  getUrlData(data: ILinkData) {
    // data为获得的链接的原始内容
    this.urlDialog = false;
    if (data) {
      this.info.content.fileData = this.getDataFromArrayOrObject(data);
      this.info.content.text = data["title"];
    } else {
      this.info.content.fileType = 0;
    }
  }
  // 获取小程序内容
  getAppletData(data: IAppData) {
    // data为获得的小程序的原始内容
    this.appletDialog = false;
    if (data) {
      this.info.content.fileData = this.getDataFromArrayOrObject(data);
      this.info.content.text = data["app_name"];
    } else {
      this.info.content.fileType = 0;
    }
  }
  /**************************************** 获取数据结束 ****************************************/
  // 删除图片/链接/小程序
  deleteMaterial() {
    // 删除功能需要置空info内的数据
    this.info.content.fileType = 0; // type值
    this.info.content.fileData = null; // 接口返回的原始数据
    this.info.content.text = ""; // 图片/链接/小程序的名称标题内容
    this.typeName = ""; // 类型（图片/链接/小程序）置空
  }
  // 关闭整体弹窗
  closeDialogFn() {
    this.$emit("addWeChatMsgDialog");
  }
  // 确认保存
  async submitFn() {
    if (this.$refs["form"] as any) {
      this.info.content.fileData == null
        ? (this.info.content.text = "")
        : false;
      if (
        this.info.content.text === undefined ||
        this.info.content.text === ""
      ) {
        let fileDataName = "";
        switch (this.info.content.fileType) {
          case 2:
            fileDataName = this.info.content.fileData.name;
            break;
          case 3:
            fileDataName = this.info.content.fileData.title;
            break;
          case 4:
            fileDataName = this.info.content.fileData.app_name;
            break;
        }
        this.info.content.text = fileDataName;
      }
      if (this.form.message !== "" || this.info.content.fileData) {
        const message: Array<object> = [
          {
            name: this.form.message,
            type: 1,
            typeName: this.form.message == "" ? "" : "文本"
          },
          {
            name: this.info.content.text || "",
            type: this.info.content.fileType || 0,
            typeName: this.typeName,
            data: this.info.content.fileData || {}
          }
        ];
        this.$emit("addWeChatMsgDialog", message);
      } else {
        this.$message({
          message: "消息内容不能为空",
          type: "warning",
          duration: 2000
        });
        return false;
      }
    }
  }
  // 数据初始化
  private getFormInfo(item: object) {
    if (item) {
      item["list"].forEach((element: any) => {
        if (item["list"].length === 2) {
          if (element.type === 1) {
            this.form.message = element.originData.name;
          } else {
            this.info.content.text = element.originData.data.name;
            this.info.content.fileData = element.originData.data;
            this.info.content.fileType = element.type;
          }
        } else if (item["list"].length === 1) {
          this.form.message = element.originData.content.text;
          this.info.content.fileType = element.originData.content.type
            ? element.originData.content.type + 1
            : 0;
          const fileData: object = {
            fileData: null,
            text: ""
          };
          switch (element.originData.content.type) {
            case 1:
              fileData["fileData"] = element.originData.content.picture;
              fileData["text"] = element.originData.content.picture.name;
              break;
            case 2:
              fileData["fileData"] = element.originData.content.link;
              fileData["text"] = element.originData.content.link.title;
              break;
            case 3:
              fileData["fileData"] = element.originData.content.applet;
              fileData["text"] = element.originData.content.applet.app_name;
          }
          this.info.content.fileData = fileData["fileData"];
          this.info.content.text = fileData["text"];
        } else {
          return false;
        }
      });
    }
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
    white-space: pre-wrap;
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
::v-deep .el-form-item {
  margin-bottom: 0px;
}
</style>
