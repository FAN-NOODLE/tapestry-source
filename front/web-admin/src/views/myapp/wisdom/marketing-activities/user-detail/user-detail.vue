<template>
  <div class="user-detail">
    <!-- 基本信息 -->
    <div class="cy-title">基本信息</div>
    <div class="user-info">
      <el-image fit="cover"
        :src="
          userInfo.img_url
            ? userInfo.img_url
            : 'https://oss-crs.vchangyi.com/wechat/default-avatar-icon.png'
        "
        class="user-avatar m-b-20"
      >
      </el-image>
      <div class="message-grid">
        <div class="message-item">
          <span class="item min-first">{{ userInfo.nickname || "-" }}</span>
          <span class="item" v-if="type == 2"
            >参与时间:{{
              userInfo.updated_at | formatDate("yyyy-mm-dd hh:mm:ss")
            }}</span
          >
          <span class="item" v-else
            >参与时间:{{
              userInfo.created_at | formatDate("yyyy-mm-dd hh:mm:ss")
            }}</span
          >
          <span class="item item-right" v-if="type == 1"
            >报名审核状态:
            <span class="c-blue">{{
              userInfo.audit_status == 1
                ? "待审核"
                : userInfo.audit_status == 2
                ? "审核通过"
                : userInfo.audit_status == 3
                ? "审核未通过"
                : "--"
            }}</span></span
          >
          <span class="item item-right" v-else-if="type == 2"
            >签到天数: <span class="c-blue">{{ userInfo.total }}</span></span
          >
          <span class="item item-right" v-else></span>
        </div>
        <div class="mt-30 message-item">
          <span class="item min-first"
            >unionid: {{ userInfo.union_id || "-" }}</span
          >
          <span class="item">导购: {{ userInfo.name }}</span>
          <span class="item item-right"></span>
        </div>
      </div>
    </div>
    <!-- 内容信息 -->
    <div class="cy-title">内容信息</div>
    <!-- 报名活动 、 问卷活动 -->
    <div class="sign-up" v-if="type === 1 || type === 3">
      <div v-for="(item, index) in activity_content" :key="index">
        <div class="title">{{ index + 1 }}、{{ item.tit }}</div>
        <p class="desc" v-if="item.mark">{{ item.mark }}</p>
        <!-- 图片上传 -->
        <div v-if="item.type == 5">
          <div class="pic-wraper">
            <el-image fit="cover"
              v-for="(pic, pIndex) in item.value"
              :key="pIndex"
              style="width: 100px; height: 100px"
              :src="pic"
              @click="checkNowImage(item.value, pIndex)"
            >
            </el-image>
          </div>
        </div>
        <!-- 单选、多选、下拉、性别、时间、地址 -->
        <div v-else-if="[1, 3, 6, 7, 9, 12].indexOf(item.type) > -1">
          <div class="message-box">
            {{ item.valueName }}
          </div>
          <!-- 单选多选时可能会有图片 -->
        </div>
        <div v-else-if="item.type == 10">
          <div class="message-box">
            {{
              item.value.substring(0, 3) + "****" + item.value.substring(7, 11)
            }}
          </div>
        </div>
        <div v-else>
          <div class="message-box">
            {{ item.value }}
          </div>
        </div>
      </div>
    </div>
    <!-- 签到活动 -->
    <div class="sign-in" v-if="type === 2">
      <div class="sign-in-item">
        <div class="left-name">签到日期</div>
        <div class="right-content calendar-box">
          <ele-calendar :signInlist="signInlist" @checkDate="checkDate" />
        </div>
      </div>
      <div class="sign-in-item mt-30" v-if="activity_content.length > 0">
        <div class="left-name">签到内容:</div>
        <div class="right-content">
          <div v-for="(item, index) in activity_content" :key="index">
            <div class="title">{{ index + 1 }}、{{ item.tit }}</div>
            <p class="desc" v-if="item.mark">{{ item.mark }}</p>
            <!-- 图片上传 -->
            <div v-if="item.type == 5">
              <div class="pic-wraper">
                <el-image fit="cover"
                  v-for="(pic, pIndex) in item.value"
                  :key="pIndex"
                  style="width: 100px; height: 100px"
                  :src="pic"
                  @click="checkNowImage(item.value, pIndex)"
                >
                </el-image>
              </div>
            </div>
            <!-- 单选、多选、下拉、性别、时间、地址 -->
            <div v-else-if="[1, 3, 6, 7, 9, 12].indexOf(item.type) > -1">
              <div class="message-box">
                {{ item.valueName }}
              </div>
              <!-- 单选多选时可能会有图片 -->
            </div>
            <div v-else-if="item.type == 10">
              <div class="message-box">
                {{
                  item.value.substring(0, 3) +
                    "****" +
                    item.value.substring(7, 11)
                }}
              </div>
            </div>
            <div v-else>
              <div class="message-box">
                {{ item.value }}
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="sign-in-item mt-30" v-if="remark && remark.length > 0">
        <div class="left-name">留资内容:</div>
        <div class="right-content">
          <div v-for="(item, index) in remark" :key="index">
            <div class="title mt-30">{{ index + 1 }}、{{ item.tit }}</div>
            <p class="desc">{{ item.tit }}</p>
            <div class="message-box" v-if="item.type === 10">
              {{
                item.value.substring(0, 3) +
                  "****" +
                  item.value.substring(7, 11)
              }}
            </div>
            <div
              class="message-box"
              v-else-if="[9, 12].indexOf(item.type) > -1"
            >
              {{ item.valueName }}
            </div>
            <div class="message-box" v-else>
              {{ item.value }}
            </div>
          </div>
        </div>
      </div>
    </div>
    <div class="table_footer">
      <el-button
        size="medium"
        v-if="type == 1 && userInfo.audit_status === 1"
        class="audit-btn"
        @click="review(2)"
      >
        审核通过
      </el-button>
      <el-button
        size="medium"
        v-if="type == 1 && userInfo.audit_status === 1"
        class="audit-btn"
        @click="review(3)"
      >
        审核未通过
      </el-button>
      <el-button size="small" @click="goBack">
        返回
      </el-button>
    </div>
    <!-- 图片视频弹框 -->
    <preview
      v-if="preview_show"
      :urls="checkImgArr"
      @close="preview_show = false"
    ></preview>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import eleCalendar from "./components/selectCalendar.vue";
import preview from "@/components/business/preview/index.vue";
import Api from "@/api/wisdom/marketing-activities";
import { dateFormat } from "../../../customer/auto-market/util/index";
@Component({
  components: {
    eleCalendar,
    preview,
  },
})
export default class extends Vue {
  // data
  private user_id = Number(this.$route.query.user_id);
  private activityId = Number(this.$route.query.activityId);
  private activity_content: Array<any> = [];
  private mp_user_id = 0;
  private remark: Array<any> = []; // 留资

  private type = 0;
  private userInfo: any = {};
  private signInlist: Array<any> = [];
  private checkedInList: Array<any> = [];
  private checkImgArr: Array<any> = [];
  private preview_show = false;
  // 生命周期
  private mounted() {
    this.getDateil();
  }

  private getDateil() {
    const params = {
      id: this.activityId,
      user_id: this.user_id,
    };
    Api.userDetail(params)
      .then((res: any) => {
        if (res.httpCode === 200) {
          try {
            if (res.data.sign_up) {
              this.type = 1;
              this.userInfo = res.data.sign_up;
              this.activity_content = JSON.parse(
                res.data.sign_up.activity_content
              ).templateList;
            } else if (res.data.sign_in) {
              this.type = 2;
              this.userInfo = res.data.sign_in;
              const start = res.data.sign_in.sign_start;
              const end = res.data.sign_in.sign_end;
              const signInTimelist = this.getSignInPeriod(start, end); // 签到周期
              this.checkedInList = res.data.sign_in.records; // 已签到
              signInTimelist.forEach((item: any) => {
                if (this.checkedInList.some((v: any) => v.date === item)) {
                  this.signInlist.push({
                    time: item,
                    isSignIn: true,
                  });
                } else {
                  this.signInlist.push({
                    time: item,
                    isSignIn: false,
                  });
                }
              });
              // 显示签到最后一天的内容
              if (this.checkedInList[this.checkedInList.length - 1].content) {
                this.activity_content = JSON.parse(
                  this.checkedInList[this.checkedInList.length - 1].content
                ).templateList;
              }

              // 有留资内容
              if (res.data.sign_in.remark) {
                this.remark = JSON.parse(res.data.sign_in.remark);
              }
            } else if (res.data.question) {
              this.type = 3;
              this.userInfo = res.data.question;
              this.activity_content = JSON.parse(
                res.data.question.activity_content
              ).templateList;
            }
          } catch (error) {
            this.activity_content = [];
            this.remark = [];
          }
          console.log("this.activity_content", this.activity_content);
        }
      })
      .catch((err: any) => {
        console.log("userDetail", err);
      });
  }

  //通过起止时间戳返回日期格式的数据
  private getSignInPeriod(start: number, end: number): Array<any> {
    const startTime = dateFormat(start, "yyyy-MM-dd");
    const endTime = dateFormat(end, "yyyy-MM-dd");
    const timeLead = Date.parse(endTime) - Date.parse(startTime);
    const dateLead = timeLead / 86400 / 1000;
    let arr: Array<any> = [];
    for (let i = 0; i <= dateLead; i++) {
      const time = Date.parse(startTime) / 1000 + i * 86400;
      const newTime = dateFormat(time, "yyyy-MM-dd");
      arr = [...arr, newTime];
    }
    return arr;
  }

  private checkDate(dataItem: any) {
    this.checkedInList.forEach((item: any) => {
      if (item.content && item.date === dataItem) {
        this.activity_content = JSON.parse(item.content).templateList;
      }
    });
  }

  // 查看图片
  private checkNowImage(item: any, index: number) {
    this.checkImgArr = [
      ...item.slice(index, item.length),
      ...item.slice(0, index),
    ];
    this.preview_show = true;
  }

  private review(type: number) {
    const params = {
      id: this.user_id,
      audit_status: type,
      mp_user_id: this.userInfo.mp_user_id,
    };
    Api.auditUser(params)
      .then((res: any) => {
        if (res.httpCode === 200) {
          this.$message.success("已完成审核");
          this.goBack();
        }
      })
      .catch((err: any) => {
        console.log("review", err);
      });
  }

  private goBack() {
    this.$router.go(-1);
  }
}
</script>
<style lang="scss" scoped>
.user-detail {
  background-color: #fff;
  padding: 20px 30px;
  border-radius: 10px;
  .user-info {
    display: flex;
    padding: 0 20px;
    .user-avatar {
      width: 70px;
      height: 70px;
      border-radius: 50%;
      margin-right: 10px;
    }
    .message-grid {
      margin-left: 10px;
      width: 100%;
      .message-item {
        display: flex;
        font-size: 16px;
        width: 100%;
        justify-content: space-between;
        .item {
          flex: 1;
          text-align: left;
        }
        .min-first {
          min-width: 440px;
        }
        .item-right {
          text-align: right;
        }
      }
    }
  }
  .sign-up,
  .questionnaire {
    padding: 0 20px;
    .title {
      font-size: 16px;
    }
    .desc {
      margin: 4px 0;
      color: #bbbbbb;
    }
    .message-box {
      font-size: 14px;
      width: 520px;
      min-height: 40px;
      line-height: 40px;
      padding-left: 20px;
      margin: 8px 0;
      background: #f8f9fc;
    }
  }
  .sign-in {
    .sign-in-item {
      display: flex;
      .left-name {
        width: 80px;
        text-align: right;
        margin-right: 20px;
      }
      .right-content {
        flex: 1;
        border: 1px solid #b9b9b9;
        padding: 20px;
        .title {
          font-size: 16px;
        }
        .desc {
          margin: 4px 0;
          color: #bbbbbb;
        }
        .message-box {
          font-size: 14px;
          width: 520px;
          min-height: 40px;
          line-height: 40px;
          padding-left: 20px;
          margin: 8px 0;
          background: #f8f9fc;
        }
      }
    }
  }
  .table_footer {
    margin-top: 30px;
    .audit-btn {
      background: #708df5;
      color: #fff;
    }
  }
}
.cy-title {
  background: #f2f2f2;
  height: 40px;
  line-height: 40px;
  padding: 0 10px;
  border-radius: 2px;
  margin: 10px 0 30px 0;
}
.mt-30 {
  margin-top: 30px;
}
.c-blue {
  color: #9cb0f8;
}
.calendar-box {
  border: none !important;
  padding-top: 0 !important;
}
/deep/ .el-date-picker-calendar {
  width: 520px;
}
.el-image {
  margin: 4px;
}
</style>
