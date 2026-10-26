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
              <p class="line" v-if="info.content.text.length > 0">
                <span class="item-avatar-left"></span>
                <span
                  class="item-content-right str"
                  v-html="info.content.text"
                ></span>
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
                              item.url
                            }}</span>
                          </div>
                          <img
                            class="select-link-cover"
                            :src="item.pic_url"
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
            <cy-form class="main-box" label-width="90px" :rules="rules">
              <cy-form-item label="群发主题">{{ info.name }}</cy-form-item>
              <cy-form-item label="执行群主">
                <el-table
                  :data="stat_list"
                  style="width: 100%"
                  max-height="500"
                >
                  <el-table-column prop="stat_list" label="员工">
                    <template slot-scope="scope">
                      <div
                        class="item-staff"
                        v-for="(item, index) in scope.row.members"
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
                    <span>111</span>
                  </el-table-column>
                </el-table>
              </cy-form-item>
              <cy-form-item label="群发类型">
                <span class="m-r-5">{{
                  info.send_type === 1
                    ? "群发客户"
                    : info.send_type === 2
                    ? "群发客户群"
                    : "--"
                }}</span>
              </cy-form-item>
              <cy-form-item label="发布时间">
                <span class="m-r-5">{{
                  info.type == 1 ? "立即发布" : "定时发布"
                }}</span>
                <span v-if="info.type == 2">
                  {{ info.send_time | formatDate("yyyy-mm-dd hh:mm:ss") }}</span
                >
              </cy-form-item>
              <cy-form-item label="消息内容">
                <div class="publish-content-box">
                  <div class="message-text-div" v-if="info.content.text">
                    <span>{{ info.content.text }} </span>
                  </div>
                </div>
              </cy-form-item>
              <cy-form-item label="消息素材">
                <el-table
                  :data="info.content.fileData"
                  stripe
                  style="width: 100%"
                  tooltip-effect="light"
                >
                  <el-table-column label="发送顺序">
                    <template slot-scope="{ $index }">
                      <span>{{ $index + 1 }} </span>
                    </template>
                  </el-table-column>
                  <el-table-column label="素材名称" min-width="120">
                    <template slot-scope="{ row }">
                      <span>{{ row.title ? row.title : "--" }} </span>
                    </template>
                  </el-table-column>
                  <el-table-column label="素材类型" prop="fileType">
                    <template slot-scope="{ row }">
                      <span>{{
                        materialSelectList[row.fileType - 1].tip
                      }}</span>
                    </template>
                  </el-table-column>
                </el-table>
              </cy-form-item>
            </cy-form>
          </div>
        </div>
      </template>
      <template slot="footer" class="dialog-footer">
        <el-button @click="goback()" size="small">返回</el-button>
        <el-button
          type="primary"
          size="small"
          @click="edit()"
          v-if="info.status == 1"
          >编辑</el-button
        >
      </template>
    </layout-form>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { taskInfoApi } from "@/api/enterprise/association";
import { ApiResult } from "@/types";
// 时间戳转换
import { dateFormat } from "@/views/myapp/customer/auto-market/util/index";
@Component({
  components: {}
})
export default class extends Vue {
  // 消息ID
  private messageId: any = 0;
  private rules = {};
  private timer: any = undefined;
  private info: any = {
    status: 0,
    name: "",
    type: 1, //1-立即 2-定时
    send_time: 0,
    content: {
      text: "",
      fileData: null,
      fileType: 0
    }
  };
  private customIcon = ["W", "X", "P", "PDF"]; //自定义icon
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
  private stat_list: Array<any> = [
    {
      members: []
    }
  ];
  // --------------------------生命周期---------------------
  private destroyed() {
    clearInterval(this.timer);
  }
  private mounted() {
    this.timer = setInterval(() => {
      this.updateTime = dateFormat(
        Number.parseInt(String(Date.now() / 1000)),
        "hh:mm"
      );
    }, 1000);
    this.messageId = this.$route.query.id ? Number(this.$route.query.id) : null;
    this.infoData();
  }
  // --------------------------methods---------------------
  // 当前时间 hh:mm
  private updateTime = dateFormat(
    Number.parseInt(String(Date.now() / 1000)),
    "hh:mm"
  );
  //详情赋值
  private analyzingData(data: any) {
    this.info.name = data.name;
    this.info.type = data.type;
    this.info.status = data.status;
    this.info.send_type = data.send_type;
    this.info.send_time = data.send_time;
    this.info.content.text = data.content.text.content;
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
  private edit() {
    this.$router.push({
      name: "myapp_enterprise_association_group-sent-message_edit",
      query: {
        id: this.messageId
      }
    });
  }
  // 详情接口
  private async infoData() {
    const { code, message = "", data }: ApiResult = await taskInfoApi(
      this.messageId
    );
    if (code !== 200) {
      this.$message.error(message);
      return;
    }
    this.analyzingData(data);
    window.sessionStorage.setItem("editdefaultData", JSON.stringify(this.info));
    this.stat_list[0].members = data.members;
    this.info.member_ids = data.members.map((item: any) => item.member_id);
  }
  // 返回
  private goback() {
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
      }
    }
  }
}

// /deep/ .content-right .el-table tr td .cell {
//   height: 105px !important;
// }
/deep/ .el-table tr td:first-child .cell {
  overflow-y: scroll !important;
}
.message-text-div {
  margin: 10px;
  padding-bottom: 10px;
}
.message-text-material-div {
  border-bottom: solid 1px #ebeef5;
}
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
</style>
