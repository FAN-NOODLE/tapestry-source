<template>
  <div class="friend-circle-setting">
    <layout-form type="mobile">
      <template slot="left">
        <div class="friend-circle-left">
          <div class="content-left">
            <div class="phone-box">
              <div class="friend-top">
                <img
                  class="mobile-phone-top"
                  src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/mobile-header.png"
                  alt=""
                />
                <!-- 背景 -->
                <img
                  class="top-bg-img"
                  v-if="info.background"
                  :src="info.background"
                  alt=""
                  style="width:375px;height:320px"
                />
                <img
                  class="top-bg-img"
                  v-else
                  src="https://wesociastg.blob.core.chinacloudapi.cn/wesocial-uat/uploads/2021/3/12/f0ea071d94dd85390d45c7d249f72951.jpeg"
                  alt=""
                  style="width:375px;height:320px"
                />
                <div class="head-portrait-box">
                  <!-- 头像 -->
                  <img
                    :src="
                      info.avatar
                        ? info.avatar
                        : 'https://wesociastg.blob.core.chinacloudapi.cn/wesocial-uat/uploads/2021/3/12/5f136e2f41e2a7e63540a03b1f4d6f56.png'
                    "
                    alt=""
                    style="width:71px;height:71px"
                  />
                </div>
                <div class="staff-name">
                  {{ info.name ? info.name : "此处展示员工名称" }}
                </div>
                <div class="signature">
                  {{ info.sign }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>
      <template slot="right">
        <div class="friend-circle-right">
          <cy-form
            class="main-box"
            ref="cform"
            :rules="rules"
            :model="info"
            label-width="100px"
          >
            <cy-form-item
              label="朋友圈背景"
              prop="background"
              required
              trigger="change"
            >
              <upload-image
                class="upload-template"
                :urls="info.background"
                accept="image/jpeg,image/png"
                @change="picSuccess1"
                type="file"
                :before-upload="beforeUpload"
              ></upload-image>
              <div class="upload-tip">{{ tipsNote }}</div>
            </cy-form-item>
            <cy-form-item label="头像" prop="avatar" required trigger="change">
              <upload-image
                class="upload-template"
                :urls="info.avatar"
                accept="image/jpeg,image/png"
                @change="picSuccess2"
                :before-upload="beforeUpload"
              ></upload-image>
              <div class="upload-tip">{{ tipsNote }}</div>
            </cy-form-item>
            <cy-form-item label="名称" prop="name" required>
              <el-input
                type="text"
                placeholder="请输入名称，如“张三”"
                v-model="info.name"
                maxlength="15"
                show-word-limit
                @input="change"
                value=""
                style="height:36px;width:280px;"
              >
              </el-input>
            </cy-form-item>
            <cy-form-item label="签名" prop="sign">
              <el-input
                type="textarea"
                placeholder="请输入签名，如“加我领取优惠哦”"
                v-model="info.sign"
                maxlength="30"
                show-word-limit
                @input="change"
                value=""
                style="position:relative;width:280px;"
              >
              </el-input>
            </cy-form-item>
            <cy-form-item label="朋友圈链接" prop="link">
              <div class="friends-links">
                <span class="work-break">{{ info.link }}</span>
                <el-button
                  class="btn-copy m-l-10"
                  size="mini"
                  type="text"
                  v-clipboard:copy="info.link"
                  v-clipboard:success="copyfn"
                  >复制</el-button
                >
              </div>
            </cy-form-item>
          </cy-form>
        </div>
      </template>
      <template slot="footer">
        <el-button type="primary" @click="submit" size="small">保存</el-button>
        <el-button plain @click="cacel" size="small">取消</el-button>
      </template>
    </layout-form>
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
  private rules = {};
  private info: any = {
    avatar: "https://wesociastg.blob.core.chinacloudapi.cn/wesocial-uat/uploads/2021/3/12/5f136e2f41e2a7e63540a03b1f4d6f56.png", // 头像
    background: "https://wesociastg.blob.core.chinacloudapi.cn/wesocial-uat/uploads/2021/3/12/f0ea071d94dd85390d45c7d249f72951.jpeg", // 背景图
    name: "", // 名称
    mark: "", // 链接加密mark
    sign: "" // 签名
  };
  private user_name: any = ""; // 员工企微姓名
  private mark: any = ""; // 链接加密mark
  private tipsNote = "支持jpg、png格式图片，大小＜2M，尺寸建议为1:1";
  /**************************************************data结束**************************************************/
  /***********************************************methods开始**************************************************/
  // 上传框
  private picSuccess1(res: any) {
    this.$set(this.info, "background", res);
    (this.$refs["cform"] as any).validateField("background");
  }
  private picSuccess2(res: any) {
    this.$set(this.info, "avatar", res);
    (this.$refs["cform"] as any).validateField("avatar");
  }
  // 复制成功提醒
  private copyfn() {
    this.$message({
      message: "复制成功",
      type: "success"
    });
  }
  // 字符数检测
  private change() {
    const userinputs = document.getElementsByClassName("el-input__count");
    if (this.info.name.length == 15) {
      (userinputs[0] as HTMLElement).style.color = "red";
    } else {
      (userinputs[0] as HTMLElement).style.color = "#909399";
    }
    if (this.info.sign.length == 30) {
      (userinputs[1] as HTMLElement).style.color = "red";
    } else {
      (userinputs[1] as HTMLElement).style.color = "#909399";
    }
  }
  // 上传文件限制规范
  private beforeUpload(file: any) {
    console.log("file", file);
    if (file.type != "image/jpeg" && file.type != "image/png") {
      this.$message.error("请上传JPG、PNG格式的图片");
      return false;
    }
    if (file.size > 2 * 1024 * 1024) {
      this.$message.error("请上传2MB大小以内的图片");
      return false;
    }
  }
  // 获取详情
  private getDetail() {
    const params = {
      mark: this.mark
    };
    API.momentAppearanceGetApi(params).then((res: any) => {
      if (res.code === 200) {
        this.info = {
          avatar: res.data.avatar
            ? res.data.avatar
            : "https://wesociastg.blob.core.chinacloudapi.cn/wesocial-uat/uploads/2021/3/12/5f136e2f41e2a7e63540a03b1f4d6f56.png", // 头像
          background: res.data.background
            ? res.data.background
            : "https://wesociastg.blob.core.chinacloudapi.cn/wesocial-uat/uploads/2021/3/12/f0ea071d94dd85390d45c7d249f72951.jpeg", // 背景图
          name: res.data.name || this.user_name, // 名称
          mark: res.data.mark, // 链接加密mark
          link: res.data.link, // 链接
          sign: res.data.sign // 签名
        };
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 点击保存按钮
  private submit() {
    (this.$refs["cform"] as any).validate((valid: any) => {
      if (valid) {
        this.save();
      }
    });
  }
  // 保存接口
  private save() {
    const params = {
      ...this.info
    };
    API.momentAppearanceEditApi(params).then((res: any) => {
      if (res.code === 200) {
        this.$message.success("保存成功");
        this.$router.go(-1);
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 点击取消按钮
  private cacel() {
    this.$router.go(-1);
  }
  /**************************************************methods结束***********************************************/
  /**************************************************生命周期结束***********************************************/
  private created() {
    const { mark, user_name } = this.$route.query;
    this.mark = mark;
    this.user_name = user_name;
  }
  private mounted() {
    this.getDetail();
  }
  private destroyed() {
    console.log("destroyed");
  }
}
</script>
<style lang="scss" scoped>
/deep/.in-item {
  min-width: 500px;
  max-width: calc(100% - 20px);
  width: auto;
}
/deep/textarea {
  resize: none;
  height: 55px !important;
  padding: 8px 15px 8px 15px;
  font-family: PingFangSC-Medium, PingFang SC;
}
/deep/.el-textarea .el-input__count {
  right: 15px;
  bottom: -2px;
  background-color: transparent;
}
.work-break {
  word-break: break-all;
}
.friend-circle-setting {
  .friend-circle-left {
    width: 375px;
    height: 667px;
    box-sizing: border-box;
    border: 1px solid #ebeef5;
    img {
      width: 375px;
      height: 667px;
    }
    .content-left {
      // flex: 0 0 405px;
      // width: 405px;
      .phone-box {
        width: 375px;
        height: 667px;
      }
      .friend-top {
        background-color: #dcdfe6;
        width: 375px;
        height: 320px;
        position: relative;
        .mobile-phone-top {
          width: 375px;
          height: 22px;
          position: relative;
          z-index: 1;
        }
        .top-bg-img {
          position: absolute;
          left: 0;
        }
        .head-portrait-box {
          width: 71px;
          height: 71px;
          position: absolute;
          right: 15px;
          bottom: -35px;
          // background-color: #f5f6fa;
        }
        .staff-name {
          position: absolute;
          font-size: 16px;
          color: #ffffff;
          right: 106px;
          bottom: 8px;
        }
        .signature {
          position: absolute;
          font-size: 14px;
          color: #606266;
          right: 15px;
          left: 15px;
          top: 364px;
          text-align: right;
        }
      }
      .friend-middle {
        display: flex;
        position: relative;
        padding-top: 106px;
        overflow: hidden;
        .middle-date {
          font-size: 28px;
          line-height: 40px;
          padding-left: 10px;
        }
        .middle-img {
          margin-left: 20px;
          margin-right: 5px;
          background: #dcdfe6;
          .img {
            width: 75px;
            height: 75px;
            display: block;
          }
        }
        .middle-content {
          span {
            color: #303133;
            font-size: 14px;
            overflow: hidden;
            text-overflow: ellipsis;
            display: -webkit-box;
            -webkit-line-clamp: 2;
            -webkit-box-orient: vertical;
            width: 188px;
            line-height: 20px;
            word-break: break-all;
          }
        }
        .img-num {
          font-size: 14px;
          color: #909399;
          position: absolute;
          left: 166px;
          bottom: 0px;
        }
      }
    }
  }
  .friend-circle-right {
    color: #606266;
    .upload-tip {
      color: #909399;
      font-size: 12px;
      line-height: 18px;
      margin-top: 5px;
    }
    .friends-links {
      position: relative;
      top: 3px;
      color: #606266;
      .btn-copy {
        font-size: 14px;
      }
    }
  }
}
</style>
