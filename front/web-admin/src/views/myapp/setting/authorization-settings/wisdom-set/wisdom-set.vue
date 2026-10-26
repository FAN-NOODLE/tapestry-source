<template>
  <div class="layout-list">
    <layout-card>
      <template slot="content">
        <div class="overflow-hidden">
          <div class="card">
            <!-- v-if="third_auth_detail.can_install === wisdom.canInstall" -->

            <div class="left flex-s-c">
              <div class="img-container">
                <img
                  src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/wisdom.png"
                  alt=""
                />
              </div>
              <div class="desc m-l-10">
                <div class="title">智慧导购</div>
                <div class="info">
                  企业微信管理员扫码安装，安装成功后可在企业微信工作台使用
                </div>
              </div>
            </div>
            <div class="right">
              <el-button
                v-if="third_auth_detail.installed === wisdom.install"
                size="small"
                class="font-14"
                @click="WisdomInstall"
              >安装</el-button>
              <el-button
                v-if="third_auth_detail.installed === wisdom.installed"
                size="small"
                class="font-14"
                disabled
              >已安装</el-button>
            </div>
          </div>
        </div>
      </template>
    </layout-card>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LayoutCard from "./components/layout-card/index.vue";

import Api from "@/api/setting/setting";

@Component({
  components: {
    LayoutCard
  }
})
export default class extends Vue {
  private wisdom = {
    install: 2,
    installed: 1,
    canInstall: 1
  };
  private wisdom_set = 1; //是否有智慧导购设置栏
  private unbindDesc = "还没有授权项哦~";
  private third_auth_detail: any = {}; //第三方授权详情
  private WisdomInstall() {
    window.location.href = this.third_auth_detail.auth_url;
  }
  private getWisdomDetail() {
    Api.installThirdApp({
      back_url: window.location.href,
      wisdom_set: this.wisdom_set
    }).then(res => {
      this.third_auth_detail = res.data;
    });
  }
  mounted() {
    this.getWisdomDetail();
  }
}
</script>
<style lang="scss" scoped>
div {
  user-select: none;
}

.card {
  padding: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #f5f6fa;
  margin-bottom: 20px;
  border-radius: 4px;

  .left {
    .img-container {
      width: 50px;
      height: 50px;
      border-radius: 4px;
      overflow: hidden;

      img {
        width: 100%;
        height: 100%;
        border-radius: 4px;
      }
    }

    .desc {
      display: flex;
      flex-direction: column;
      justify-content: space-between;
      color: #606266;
      font-size: 14px;

      .title {
        font-size: 14px;
      }

      .info {
        font-size: 12px;
      }
    }
  }
}

.card:last-child {
  margin-bottom: 0px;
}

.el-icon-warning {
  font-size: 20px;
  color: #f6ad42;
}

.flex-s {
  display: flex;
  justify-content: flex-start;
}

.flex-s-c {
  display: flex;
  justify-content: flex-start;
  align-items: center;
}

.m-t-150 {
  margin-top: 150px;
}

.overflow-hidden {
  overflow: hidden;
  overflow-y: scroll;
}</style>
