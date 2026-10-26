<template>
  <div class="friend-circle-detail">
    <div class="friend-circle-left">
      <layout-form>
        <template slot="content">
          <el-button type="primary" @click="goAdd" size="small"
            >发布朋友圈</el-button
          >
          <div class="cut-line"></div>
          <div class="p-b-5" v-if="friendList.length > 0">
            <div
              class="count-detail"
              v-for="(item, index) in friendList"
              :key="index"
            >
              <!-- 朋友圈顶部信息 -->
              <div>
                <div class="detail-top">
                  <img
                    v-if="item.user_avatar"
                    class="avatar"
                    :src="
                      item.moment_user_avatar
                        ? item.moment_user_avatar
                        : 'https://wesociastg.blob.core.chinacloudapi.cn/wesocial-uat/uploads/2021/3/12/5f136e2f41e2a7e63540a03b1f4d6f56.png'
                    "
                    alt=""
                  />
                  <div class="user-message">
                    <p>
                      {{ item.moment_user_name ? item.moment_user_name : "--" }}
                    </p>
                    <span
                      >{{
                        item.publish_type == 1
                          ? item.created_at
                          : item.publish_time
                            | formatDate("yyyy-mm-dd hh:mm:ss")
                      }}
                    </span>
                  </div>
                  <div class="user-work">
                    <el-dropdown @command="handleCommand">
                      <div class="operation-box">
                        <img src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/menu.png" alt="" />
                      </div>
                      <el-dropdown-menu slot="dropdown">
                        <el-dropdown-item :command="'edit_' + index"
                          >编辑</el-dropdown-item
                        >
                        <el-dropdown-item
                          class="color-red"
                          :command="'delete_' + index"
                          >删除</el-dropdown-item
                        >
                      </el-dropdown-menu>
                    </el-dropdown>
                  </div>
                </div>
                <div class="user-count">
                  <span v-html="textStr(item.content_json.text)"></span>
                </div>
              </div>
              <div>
                <!-- 图片 -->
                <div
                  class="detail-image"
                  v-if="item.content_json.picture.length > 0"
                >
                  <div
                    :class="{
                      'image-top': item.content_json.picture.length === 4,
                      'image-second':
                        item.content_json.picture.length > 4 &&
                        item.content_json.picture.length < 7,
                      'image-three':
                        item.content_json.picture.length > 6 &&
                        item.content_json.picture.length < 10,
                      'image-four': item.content_json.picture.length <= 3
                    }"
                  >
                    <img
                      v-for="(item, index) in item.content_json.picture"
                      :key="index"
                      :src="item.url"
                      alt=""
                    />
                  </div>
                </div>
                <!-- 链接 -->
                <div
                  class="detail-url"
                  v-if="
                    item.content_json.link && item.content_json.link.cover_url
                  "
                >
                  <div class="url-top">
                    <div class="top-second">
                      <img alt="" :src="item.content_json.link.cover_url" />
                      <span>{{ item.content_json.link.title }}</span>
                    </div>
                  </div>
                </div>
                <!-- 小程序 -->
                <div
                  class="detail-applets"
                  v-if="
                    item.content_json.applet &&
                      item.content_json.applet.app_name
                  "
                >
                  <div class="applets-top">
                    <div class="appletstop-first">
                      <img
                        src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/mini-applet-icon.png"
                        alt=""
                      />
                      <span class="app-name">{{
                        item.content_json.applet.app_name
                      }}</span>
                    </div>
                    <div class="appletstop-second">
                      <span class="app-title">{{
                        item.content_json.applet.title
                      }}</span>
                    </div>
                    <div class="appletstop-three">
                      <img :src="item.content_json.applet.cover_url" alt="" />
                    </div>
                    <div class="appletstop-four">
                      <img
                        src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/mini-icon.png"
                        alt=""
                      />
                      <span>小程序</span>
                    </div>
                  </div>
                </div>
                <!-- 视频 -->
                <div
                  class="detail-video"
                  v-if="item.content_json.video && item.content_json.video.url"
                >
                  <div class="video-count">
                    <img
                      class="play-icon"
                      src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/play_icon.png"
                      alt=""
                    />
                    <video
                      width="145"
                      height="180"
                      x5-video-player-fullscreen="true"
                      x5-playsinline
                      playsinline
                      webkit-playsinline
                      muted
                      preload="auto"
                      :src="item.content_json.video.url"
                    ></video>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div v-else class="empty-label">
            <img src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/default/default_noList.png" alt="" />
            <p>暂无朋友圈</p>
          </div>
        </template>
      </layout-form>
    </div>
    <div class="friend-circle-right">
      <layout-form>
        <template slot="content">
          <h3>形象设置</h3>
          <div class="one-row">
            <span>朋友圈背景</span>
            <img
              v-if="info.background"
              :src="info.background"
              alt=""
              class="count-rowimg"
            />
            <img
              v-else
              src="https://wesociastg.blob.core.chinacloudapi.cn/wesocial-uat/uploads/2021/3/12/f0ea071d94dd85390d45c7d249f72951.jpeg"
              alt=""
              class="count-rowimg"
            />
          </div>
          <div class="two-row">
            <span>头像</span>
            <img
              v-if="info.avatar"
              :src="info.avatar"
              alt=""
              class="count-rowimg"
            />
            <img
              v-else
              src="https://wesociastg.blob.core.chinacloudapi.cn/wesocial-uat/uploads/2021/3/12/5f136e2f41e2a7e63540a03b1f4d6f56.png"
              alt=""
              class="count-rowimg"
            />
          </div>
          <div class="three-row">
            <span>名称</span>
            <p>
              {{ info.name ? info.name : this.user_name }}
            </p>
          </div>
          <div class="four-row">
            <span>签名</span>
            <p>{{ info.sign ? info.sign : "暂无签名" }}</p>
          </div>
          <div class="five-row">
            <span>朋友圈链接</span>
            <div>
              {{ info.link }}
              <el-button
                size="mini"
                type="text"
                v-clipboard:copy="info.link"
                v-clipboard:success="copyfn"
                >复制</el-button
              >
            </div>
          </div>
        </template>
        <template slot="footer">
          <el-button type="primary" @click="routerLinkfn" size="small"
            >编辑</el-button
          >
        </template>
      </layout-form>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LayoutForm from "@/components/public/layout-form/index.vue";
import UploadImage from "@/components/business/upload-image/index.vue";
import API from "@/api/enterprise/friend-circle";
@Component({
  components: {
    LayoutForm,
    UploadImage
  }
})
export default class extends Vue {
  /**************************************************data开始**************************************************/
  private friendList: any = ""; //朋友圈列表
  private detailImage: Array<any> = []; //朋友圈图片储存
  private detailUrl: Array<any> = []; //朋友圈地址储存
  private detailApplets: Array<any> = []; //朋友圈小程序储存
  private detailVideo: Array<any> = []; //朋友圈视频储存
  private mark: any = ""; // 链接加密mark
  private member_id: any = ""; // 当前员工id
  private user_name: any = ""; // 当前员工姓名
  private info = {
    id: 0,
    avatar: "", // 头像
    background: "", // 背景图
    member_id: 0, // 员工id
    name: "", // 名称
    mark: "", // 链接加密mark
    sign: "", // 签名
    updated_at: ""
  };
  /**************************************************data结束**************************************************/
  /***********************************************methods开始**************************************************/
  // 点击三点点
  private handleCommand(command: string) {
    if (command.indexOf("edit") > -1) {
      const index = Number(command.replace("edit_", ""));
      console.log(index, this.friendList[index]);
      this.$router.push({
        name: "myapp_enterprise_friend-circle_friend-circle-edit",
        query: {
          mark: this.mark,
          is_edit: "1",
          moment_topic_id: this.friendList[index].moment_topic_id,
          moment_topic_member_id: this.friendList[index].moment_topic_member_id,
          member_id: this.friendList[index].member_id,
          user_name: this.friendList[index].user_name
        }
      });
    }
    if (command.indexOf("delete") > -1) {
      this.delete(Number(command.replace("delete_", "")));
    }
  }
  // 删除朋友圈按钮
  private delete(index: number) {
    this.$confirm("确认删除该条朋友圈吗？", "删除确认", {
      confirmButtonText: "确认删除",
      confirmButtonClass: "button-danger",
      cancelButtonText: "我再想想",
      type: "warning"
    })
      .then(() => {
        const params = {
          member_id: this.friendList[index].member_id,
          moment_topic_id: this.friendList[index].moment_topic_id,
          moment_topic_member_id: this.friendList[index].moment_topic_member_id
        };
        API.memberMomentDeleteApi(params).then((res: any) => {
          if (res.code === 200) {
            this.$message({
              message: "删除成功",
              type: "success",
              duration: 2000
            });
            this.friendDetailLeft();
          } else {
            this.$message.error(res.message);
          }
        });
      })
      .catch(() => {
        this.$message({
          type: "info",
          message: "已取消删除"
        });
      });
  }
  // 复制成功提醒
  private copyfn() {
    this.$message({
      message: "复制成功",
      type: "success"
    });
  }
  // 页面跳转
  private routerLinkfn() {
    this.$router.push({
      name: "myapp_enterprise_friend-circle_friend-circle-setting",
      query: {
        mark: this.mark,
        user_name: this.user_name
      }
    });
  }
  // 点击发布朋友圈，给该员工新增一条朋友圈
  private goAdd() {
    this.$router.push({
      name: "myapp_enterprise_friend-circle_friend-circle-edit",
      query: {
        mark: this.mark,
        member_id: this.member_id,
        user_name: this.user_name
      }
    });
  }

  //获取朋友圈形象(右侧)
  private friendDetailRight() {
    const params = {
      mark: this.mark
    };
    API.momentAppearanceGetApi(params)
      .then((res: any) => {
        if (res.code === 200) {
          this.info = res.data;
        }
      })
      .catch((error: any) => {
        console.log(error);
      });
  }
  //获取朋友圈(左侧)
  private friendDetailLeft() {
    API.momentMemberTopicApi({
      mark: this.mark,
      page: 1,
      size: 9999
    }).then((res: any) => {
      if (res.code === 200) {
        this.friendList = res.data.list;
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 处理文本内容
  public textStr(text: string): string {
    return text.replace(/\n/g, "<br/>").replace(/\s/g, "&nbsp;");
  }
  /**************************************************methods结束***********************************************/
  /**************************************************生命周期结束***********************************************/
  private created() {
    const { mark, user_name, member_id } = this.$route.query;
    this.mark = mark;
    this.user_name = user_name;
    this.member_id = member_id;
  }
  private mounted() {
    this.friendDetailLeft();
    this.friendDetailRight();
  }
  private destroyed() {
    console.log("destroyed");
  }
}
</script>

<style lang="scss" scoped>
/deep/.el-dropdown-menu__item.color-red {
  color: $color-danger;
}
.friend-circle-detail {
  display: flex;
  .friend-circle-left {
    background: #ffffff;
    width: 700px;
    margin-right: 30px;
    border-radius: 10px;
  }
  .friend-circle-right {
    flex: 1;
    width: 500px;
    // background: #ffffff;
    border-radius: 10px;
    .one-row {
      display: flex;
      span {
        color: $color-text-normal;
        margin-right: 20px;
      }
      .count-rowimg {
        width: 100px;
        height: 100px;
        border-radius: 6px;
      }
    }
    .two-row {
      display: flex;
      margin-top: 20px;
      span {
        color: $color-text-normal;
        margin-left: 42px;
        margin-right: 20px;
      }
      .count-rowimg {
        width: 100px;
        height: 100px;
        border-radius: 6px;
      }
    }
    .three-row {
      display: flex;
      margin-top: 30px;
      color: $color-text-main;
      span {
        color: $color-text-normal;
        margin-left: 42px;
        margin-right: 20px;
      }
    }
    .four-row {
      display: flex;
      margin-top: 30px;
      color: $color-text-main;
      span {
        color: $color-text-normal;
        white-space: nowrap;
        margin-left: 42px;
        margin-right: 20px;
      }
    }
    .five-row {
      display: flex;
      margin-top: 30px;
      color: $color-text-main;
      span {
        color: $color-text-normal;
        white-space: nowrap;
      }
      div {
        margin-left: 20px;
        position: relative;
        top: -5px;
        word-break: break-all;
      }
    }
  }
  h3 {
    font-weight: normal;
    margin-bottom: 30px;
  }
  .cut-line {
    border-bottom: solid 1px #ebeef5;
    width: 100%;
    height: 1px;
    margin: 30px 0px;
  }
  .empty-label {
    text-align: center;
    margin-top: 100px;
    img {
      width: 124px;
    }
  }
  .count-detail {
    margin-bottom: 35px;
    .detail-top {
      display: flex;
      position: relative;
      .avatar {
        width: 70px;
        height: 70px;
        border-radius: 50%;
        margin-top: 5px;
      }
      .user-message {
        position: relative;
        top: 20px;
        margin-left: 10px;
        p {
          display: line-block;
          color: #292a30;
        }
        span {
          display: line-block;
          margin-top: 20px;
          color: #606266;
        }
      }
      .user-work {
        width: 20px;
        height: 20px;
        position: absolute;
        right: 0;
        margin-top: 20px;
        z-index: 1;
        .operation-box {
          width: 30px;
          height: 17px;
          text-align: center;
          img {
            width: 3px;
            height: 17px;
          }
        }
      }
    }
    .user-count {
      margin-top: 20px;
      span {
        color: #292a30;
        word-break: break-all;
      }
    }
    .detail-image {
      margin-top: 20px;
      .image-top {
        width: 260px;
        display: flex;
        flex-wrap: wrap;
        margin-top: 10px;
        margin-left: 10px;
        img {
          width: 100px;
          height: 100px;
          margin-top: 10px;
          margin-left: 10px;
        }
      }
      .image-second {
        width: 330px;
        display: flex;
        flex-wrap: wrap;
        img {
          width: 100px;
          height: 100px;
          margin-left: 10px;
          margin-top: 10px;
        }
      }
      .image-three {
        width: 330px;
        display: flex;
        flex-wrap: wrap;
        img {
          width: 100px;
          height: 100px;
          margin-left: 10px;
          margin-top: 10px;
        }
      }
      .image-four {
        width: 330px;
        display: flex;
        img {
          width: 100px;
          height: 100px;
          margin-right: 10px;
        }
      }
    }
    .detail-url {
      margin-top: 20px;
      .url-top {
        width: 240px;
        height: 70px;
        background: #f5f6fa;
        // .top-first {
        //   width: 250px;
        //   height: 40px;
        //   margin-left: 10px;
        //   padding-top: 10px;
        // }
        .top-second {
          width: 240px;
          height: 40px;
          display: flex;
          span {
            margin-right: 10px;
            color: #292a30;
            font-weight: bold;
            width: 200px;
            height: 40px;
            display: inline-block;
            line-height: 40px;
            margin-top: 15px;
            margin-left: 10px;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
          }
          img {
            margin-top: 10px;
            margin-left: 10px;
            width: 50px;
            height: 50px;
          }
        }
      }
    }
    .detail-applets {
      margin-top: 20px;
      .applets-top {
        width: 200px;
        height: 320px;
        background: #f5f6fa;

        .appletstop-first {
          padding-top: 10px;
          margin-left: 10px;
          display: flex;
          align-items: center;
          img {
            width: 30px;
            height: 30px;
            border-radius: 50%;
          }
          .app-name {
            margin-left: 5px;
            color: #909399;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
            width: 142px;
            display: inline-block;
          }
        }
        .appletstop-second {
          width: 168px;
          max-height: 40px;
          margin-top: 10px;
          margin-left: 10px;
          color: #606266;
          .app-title {
            overflow: hidden;
            text-overflow: ellipsis;
            display: -webkit-box;
            -webkit-line-clamp: 2;
            -webkit-box-orient: vertical;
          }
        }
        .appletstop-three {
          width: 180px;
          height: 180px;
          margin-top: 10px;
          margin-left: 10px;
          img {
            width: 180px;
            height: 180px;
          }
        }
        .appletstop-four {
          margin-top: 10px;
          margin-left: 10px;
          img {
            width: 15px;
            height: 14px;
          }
          span {
            margin-left: 5px;
            color: #606266;
          }
        }
      }
    }
    .detail-video {
      margin-top: 20px;
      .video-count {
        width: 145px;
        height: 180px;
        position: relative;
        width: fit-content;
        .play-icon {
          position: absolute;
          left: 0;
          right: 0;
          bottom: 0;
          top: 0;
          margin: auto;
          width: 25px;
          height: 25px;
          z-index: 99;
        }
        video {
          display: inline-block;
          vertical-align: baseline;
          object-fit: fill;
        }
      }
    }
  }
}
</style>
