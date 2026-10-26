<template>
  <div class="events-detail">
    <div class="cy-title">
      基本信息
    </div>
    <div class="tabs-firstrow">
      <el-image fit="cover"
        :src="userInfo.avatar
            ? userInfo.avatar
            : 'https://oss-crs.vchangyi.com/wechat/default-avatar-icon.png'
          "
        class="user-avatar m-b-20"
      >
      </el-image>
      <div class="message-grid">
        <div class="message-item">
          <span>{{ userInfo.wx_name || "-" }}</span>
          <span>状态: <span class="c-blue">{{ userInfo.status_text }}</span></span>
        </div>
        <div class="mt-30">手机号码: {{ userInfo.wx_phone || "-" }}</div>
      </div>
    </div>
    <div class="cy-title">
      内容信息
    </div>
    <el-row>
      <el-col
        :span="2"
        class="col-left"
      >报名内容：</el-col>
      <el-col :span="22">
        {{ userInfo.content }}
      </el-col>
    </el-row>
    <el-row>
      <el-col
        :span="2"
        class="col-left"
      >上传图片：</el-col>
      <el-col :span="22">
        <el-image fit="cover"
          :src="userInfo.picture_url"
          class="user-content m-b-20"
        >
        </el-image>
      </el-col>
    </el-row>
    <div class="cy-title">
      收货信息
    </div>
    <el-table
      :data="list"
      stripe
      tooltip-effect="light"
    >
      <el-table-column
        prop="receiver_name"
        label="收货人"
        width="120"
      >
      </el-table-column>
      <el-table-column
        prop="receiver_phone"
        label="收货人手机号"
        width="160"
      >
      </el-table-column>
      <el-table-column
        prop="receiver_address"
        label="收货地址"
        min-width="240"
      >
      </el-table-column>
      <div
        slot="empty"
        class="table-empty-label-page"
      >
        <img
          class="empty-img"
          :src="require('@/assets/images/empty-box.png')"
          alt=""
        />
        <div>暂无收货信息</div>
      </div>
    </el-table>
    <div class="table_footer">
      <el-button
        size="medium"
        class="pass"
        :disabled="isPass"
        :class="[isPass ? 'nopass' : 'pass']"
        @click="pass"
      >审核通过</el-button>
      <el-button
        size="medium"
        @click="goBack"
      >返回</el-button>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Api from "@/api/wisdom/researchActivities";
@Component
export default class extends Vue {
  // data
  private userId = 0;
  private userInfo: any = {};
  private list: Array<any> = [];
  private isPass = false;
  // 生命周期
  private mounted() {
    const { id } = this.$route.query;
    this.userId = Number(id);
    this.getUserInfo();
  }
  private getUserInfo() {
    Api.getUserInfo({ activity_user_id: this.userId }).then((res: any) => {
      if (res.code === 200) {
        this.userInfo = res.data;
        if (this.userInfo.status !== 1) {
          this.isPass = true;
        }
        this.list = [
          {
            receiver_name: this.userInfo.receiver_name,
            receiver_phone: this.userInfo.receiver_phone,
            receiver_address: this.userInfo.receiver_address
          }
        ];
      }
    });
  }
  private pass() {
    if (this.isPass) {
      return;
    }
    Api.verifyUser({ activity_user_id: this.userId }).then((res: any) => {
      if (res.code === 200) {
        this.$message.success("审核通过");
        this.goBack();
      }
    });
    console.log("pass");
  }
  private goBack() {
    this.$router.go(-1);
  }
}
</script>
<style lang="scss" scoped>
.events-detail {
  background-color: #fff;
  padding: 20px 30px;
  border-radius: 10px;
}

.cy-title {
  background: #f2f2f2;
  height: 40px;
  line-height: 40px;
  padding: 0 10px;
  border-radius: 2px;
  margin: 10px 0 30px 0;
}

.tabs-firstrow {
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
      font-size: 18px;
      width: 100%;
      justify-content: space-between;
    }
  }
}

.mt-30 {
  margin-top: 30px;
}

.c-blue {
  color: #9cb0f8;
}

.col-left {
  text-align: right;
}

.el-row {
  margin-bottom: 20px;
}

.table_footer {
  margin-top: 20px;

  .pass {
    background: #708df5;
    color: #fff;
  }

  .nopass {
    background: #d7d7d7;
    color: #fff;
    cursor: no-drop;
  }
}

.user-content {
  width: 300px;
  height: 400px;
}</style>
