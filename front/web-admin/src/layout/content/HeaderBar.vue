<template>
  <!-- 顶部栏主要为顶部面包屑/切换平台/下载中心/帮助中心/用户中心 -->
  <div class="header-bar">
    <!-- 面包屑部分 -->
    <div class="bread-crumb">
      <el-breadcrumb
        class="app-breadcrumb m-l-30"
        separator-class="el-icon-arrow-right"
        aria-hidden="true"
      >
        <el-breadcrumb-item
          class="crs-bread-item"
          v-for="item in levelList"
          :key="item.path"
        >
          <span
            class="link"
            @click="goPage(item)"
          > {{ item.name }}</span>
        </el-breadcrumb-item>
      </el-breadcrumb>
    </div>
    <!-- 切换平台/下载中心/帮助中心/用户中心 -->
    <div class="icon-bar pull-right">
      <div class=" pull-left m-r-20">
        <download-task></download-task>
      </div>
      <div
        class="icon-button pull-left m-r-20"
        @click="goDownload"
      >
        <el-tooltip
          content="下载中心"
          :open-delay="1000"
          placement="bottom"
        >
          <svg-icon
            class="header-icon"
            icon-class="download-icon"
          ></svg-icon>
        </el-tooltip>
      </div>
      <!-- 暂时隐藏 -->
      <div
        class="icon-button pull-left m-r-20"
        v-if="false"
      >
        <svg-icon
          class="header-icon"
          icon-class="message"
        ></svg-icon>
      </div>
      <div
        class="icon-button pull-left m-r-20"
        @click="reloadPage"
      >
        <el-tooltip
          class="item"
          effect="dark"
          content="刷新"
          placement="bottom"
          :open-delay="1000"
        >
          <svg-icon
            style="width:16px;height:16px;"
            class="header-icon"
            icon-class="refresh"
          > </svg-icon>
        </el-tooltip>
      </div>
      <div class="icon-line pull-left m-r-20 m-t-20"></div>
      <logo-hover
        @selectEnterprise="selectEnterprise"
        @logout="logout"
      ></logo-hover>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Watch } from "vue-property-decorator";
import BreadCrumbModule from "@/store/modules/bread-crumb";
import { IBCNodeItem } from "@/types";
import LogoHover from "../components/logo-hover.vue";
import { hasHelp } from "@/utils/help-data";
import { VUE_APP_API_ROOT } from "@/utils/env-config";
import LoginAboat from "@/utils/LoginManager";
import DownloadTask from './DownloadTask.vue'
@Component({
  components: { LogoHover,DownloadTask },
})
export default class HeaderBar extends Vue {
  private see = false;
  private pollingST: any = null;
  get BREADCRUMB() {
    console.log(BreadCrumbModule.getBreadCrumb)
    return BreadCrumbModule.getBreadCrumb;
  }
  @Watch("$route")
  async firstNameChange() {
    this.isHelp = await hasHelp();
    this.see = false;
    this.getBreadcrumb();
  }
  private isHelp = false;
  private levelList: Array<any> = [];
  private getBreadcrumb() {
    this.levelList = this.BREADCRUMB;
  }
  // todo:自动下载接口
  private async getDownload() {
    // 轮询下载接口
    try {
      const res: any = await this.$axios.request({
        headers: {
          "Enterprise-Hash": LoginAboat.getAppHash(),
          "Auth-Token": LoginAboat.getToken()
        },
        method: "post",
        url: `${VUE_APP_API_ROOT}/download-center/get-auto-detail`,
      })
      if (res.data.code === 200) {
        const filePath = res.data.data.path;
        // 判断filePath是否存在，判断有后缀.的文件
        if (filePath && filePath.indexOf(".") > -1) {
          window.location.href = res.data.data.path;
        }
      }
      if(window.pluginWebUpdateNotice_){
        window.pluginWebUpdateNotice_.checkUpdate()
      }
    } catch (error) {
      console.log(error);
    }
  }
  // 轮询方法
  polling() {
    this.pollingST = setInterval(() => {
      this.getDownload()
    }, 60000)
  }
  // 退出登录清除定时器
  logout() {
    console.log("退出登录清除定时器");
    clearInterval(this.pollingST)
  }
  private goDownload() {
    this.$router.push({
      name: "myapp_setting_download-center_list"
    })
    this.see = false;
  }
  private async reloadPage() {
    this.$alert("是否刷新页面？", "提示", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning"
    })
      .then(() => {
        // 清除轮询
        clearInterval(this.pollingST)
        window.location.replace(window.location.href)
      })
      .catch(() => {
        console.log("取消刷新");
      });
  }
  private selectEnterprise() {
    this.$emit("selectEnterprise");
  }
  private goPage(item: IBCNodeItem) {
    const curKey = ("myapp" + location.pathname.split("/myapp")[1]).replace(
      /\//g,
      "_"
    ); // 获取当前页面的key
    if (item.children.length > 0 && item.children[0].key !== curKey) {
      this.$router.push({
        name: item.children[0].key
      })
    }
  }
  async mounted() {
    this.isHelp = await hasHelp();
    console.log("mounted中的状态" + this.isHelp);
    this.getBreadcrumb();
    this.polling();
  }
  destroyed(){
    // 清除轮询
    clearInterval(this.pollingST)
  }
}
</script>
<style lang="scss" scoped>
.header-bar {
  height: 100%;
  background-color: #ffffff;
  border-radius: 15px 0 0 0;

  .bread-crumb {
    float: left;

    .app-breadcrumb {
      padding-top: 12px;

      .crs-bread-item {
        cursor: pointer;

        .link {
          font-size: 14px !important;
          color: #606266 !important;
          font-weight: 400 !important;
          line-height: 26px !important;
        }

        &:hover {
          .link {
            color: $color-main !important;
          }
        }
      }

      .crs-bread-item:last-child {
        .link {
          font-size: 18px !important;
          color: #292a30 !important;
          font-weight: 600 !important;
        }
      }
    }
  }

  .icon-bar {
    height: 100%;
    line-height: 50px;

    .header-icon {
      width: 20px;
      height: 20px;
      vertical-align: middle;
      cursor: pointer;
    }

    .icon-button {
      cursor: pointer;

      .title {
        display: inline-block;
        color: #606266;
        vertical-align: middle;
        margin-left: 3px;
      }
    }

    .icon-line {
      height: 9.22px;
      border-right: 1px solid #d8d8d8;
    }

    .brand-logo {
      width: 30px;
      height: 30px;
      margin-top: 8px;
      border-radius: 50%;
      background-size: cover;
      background-color: #222222;
      background-image: url("https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/logo.jpg");
    }
  }
}

.down-tip {
  right: 110px !important;
}

.tip {
  padding: 4px 6px 10px 10px;
  position: absolute;
  top: 47px;
  right: 66px;
  // width: 210px;
  // padding: 0px 8px;
  // left: 1203px;
  // top: 47px;
  width: 158px;
  height: 45px;
  line-height: 20px;
  border-radius: 4px 4px 4px 4px;
  background: #ffffff;
  text-align: center;
  box-shadow: 0px -2px 12px 0px rgba(0, 0, 0, 0.1);
  border: 1px solid #e4e7ed;
}

.trianglr {
  position: absolute;
  right: 31.5px;
  top: -12.5px;
  width: 1px;
  border-top: 5.9px solid transparent;
  border-right: 5.9px solid transparent;
  border-left: 5.9px solid transparent;
  border-bottom: 5.9px solid #e4e7ed;
  z-index: 99;
}

.trianglr-inner {
  position: absolute;
  right: -5.1px;
  top: -4.5px;
  width: 0px;
  height: 0px;
  border-top: 5.4px solid transparent;
  border-right: 5.5px solid transparent;
  border-left: 5.5px solid transparent;
  border-bottom: 5.7px solid #fff;
  z-index: 100;
}

.word {
  // padding: 10px;
  // position: absolute;
  // top: 2px;
  // right: 6px;
  width: 158px;
  height: 45px;
  font-size: 14px;
  // background: pink;
  line-height: 25px;
  color: #606266;
  text-align: left;
  font-family: Arial-regular;
}

.item {
  margin: 4px;
}
</style>
