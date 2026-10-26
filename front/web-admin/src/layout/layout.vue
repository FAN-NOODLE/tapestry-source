<template>
  <el-container class="base-crs-container">
    <!-- 选择切换企业弹窗 -->
    <select-enterprise
      :show="isShowSelectEnt"
      :on-close="() => {
        isShowSelectEnt = false;
      }
        "
    ></select-enterprise>
    <!-- 侧边栏菜单部分 -->
    <el-aside class="base-crs-aside">
      <asider-bar></asider-bar>
    </el-aside>
    <el-container class="crs-wrap">
      <!-- 头部栏 -->
      <el-header
        class="crs-header"
        :class="scroll ? 'add-shadow' : ''"
      >
        <header-bar
          :isHelp="isHelp"
          @update="showHelp"
          @selectEnterprise="selectEnterprise"
        ></header-bar>
      </el-header>
      <!-- 主容器 -->
      <el-main
        class="crs-main"
        :class="HELPSHOW ? 'short-width' : ''"
        ref="crsMain"
      >
        <div id="crs-container">
          <router-view></router-view>
        </div>
        <!-- 底部版本号 -->
        <!-- <div class="footer-version">{{ version_number }}</div> -->
        <!--TODO未来基座逻辑处理地方目前注释-->
        <!-- <router-view></router-view> -->
      </el-main>
      <!-- 帮助中心 -->
      <transition name="slide-fade">
        <article
          class="crs-article"
          v-if="HELPSHOW"
        >
          <help-center></help-center>
        </article>
      </transition>
    </el-container>
  </el-container>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

import { HeaderBar, HelpCenter, AsiderBar } from "./content";
import { addEvent } from "@/utils/event";
import SelectEnterprise from "@/components/business/select-enterprise/index.vue";
import HelpStore from "@/store/modules/help";
import { hasHelp } from "@/utils/help-data";
import Api from "@/api/public/public";
@Component({
  components: {
    HeaderBar,
    HelpCenter,
    AsiderBar,
    SelectEnterprise
  }
})
export default class extends Vue {
  private index = 0;
  private scroll = false;
  private isOpenQIANKUN = false;
  private isShowSelectEnt = false;
  private isHelp = true;
  private version_number = "";
  get HELPSHOW() {
    return HelpStore.getHelp;
  }
  private created() {
    this.getCurrentVersionInfo();
  }
  mounted() {
    // 动态添加阴影
    addEvent("scroll", (this.$refs.crsMain as any).$el, () => {
      this.scroll = (this.$refs.crsMain as any).$el.scrollTop > 0;
    });
  }
  private getCurrentVersionInfo() {
    Api.currentVersionInfo({}).then((res: any) => {
      if (res.data) {
        this.version_number = res.data.version_info.them;
        console.log("当前版本号", this.version_number);
      }
    });
  }
  private findIndex(path: string, list: Array<any>): void {
    for (let i = 0; i < list.length; i++) {
      if (list[i].meta && path === list[i].meta.file_path) {
        this.index = i;
        break;
      } else {
        if (list[i].children && list[i].children.length > 0) {
          this.findIndex(path, list[i].children);
        }
      }
    }
  }
  private async showHelp() {
    this.isHelp = await hasHelp();
    console.log("showHelp中的状态" + this.isHelp);
    if (!this.isHelp) {
      // console.log(this.isHelp)
      return;
    }
    this.HELPSHOW ? HelpStore.setHelp("2") : HelpStore.setHelp("1");
  }
  private selectEnterprise() {
    this.isShowSelectEnt = true;
  }
}
</script>
<style lang="scss" scoped>
.base-crs-container {
  width: 100%;
  height: 100vh;
  scrollbar-width: none;

  .base-crs-aside {
    width: 140px !important;
    max-height: 100vh;

    ::-webkit-scrollbar {
      display: none;
    }
  }

  .crs-wrap {
    position: relative;
    background-color: $color-background;
    border-radius: 15px 0 0 0;
    width: 100%;

    .crs-header {
      height: 50px !important;
      padding: 0 !important;
      position: absolute;
      top: 0;
      width: 100%;
      z-index: 1000;
    }

    .add-shadow {
      box-shadow: 0px 4px 10px 0px rgba($color-dark, 0.13);
    }

    .crs-main {
      &::-webkit-scrollbar {
        width: 0;
      }

      width: 100%;
      height: 100%;
      padding: 80px 30px 30px;
      box-sizing: border-box;
      display: inline-block;
      transition: width 0.3s cubic-bezier(1, 0.5, 0.8, 1);

      .footer-version {
        position: relative;
        bottom: -20px;
        text-align: center;
      }
    }

    .short-width {
      width: calc(100% - 270px);
      transition: width 0.3s cubic-bezier(1, 0.5, 0.8, 1);
    }

    .crs-article {
      position: absolute;
      width: 240px;
      max-height: calc(100% - 82px);
      overflow: auto;
      right: 30px;
      top: 80px;
      box-shadow: 0px 0px 20px 1px rgba(0, 0, 0, 0.08);
      background: #ffffff;
      font-size: 18px;
      text-align: center;
      line-height: 500px;

      .triangle {
        position: absolute;
        top: -8px;
        right: 64px;
        width: 0;
        height: 0;
        border-right: 8px solid transparent;
        border-left: 8px solid transparent;
        border-bottom: 8px solid #ffffff;
      }
    }
  }
}
</style>
<style lang="scss" scoped>
.slide-fade-enter-active {
  transition: all 0.4s ease;
}

.slide-fade-leave-active {
  transition: all 0.4s cubic-bezier(1, 0.5, 0.8, 1);
}

.slide-fade-enter,
.slide-fade-leave-to {
  transform: translateX(270px);
  opacity: 0;
}
</style>
