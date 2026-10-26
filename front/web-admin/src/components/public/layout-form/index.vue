<template>
  <div class="layout-form" ref="layout" :style="{ '--width': `${width}px` }">
    <!-- 常规 -->
    <div
      class="content"
      :class="{ 'pad-b-30': !hasSlotFooter }"
      ref="content"
      v-if="type === 'normal'"
    >
      <div class="container">
        <slot name="content"></slot>
      </div>
      <div
        class="footer"
        v-if="$slots.footer"
        :class="{
          'fixed-beyound-long':
            (_isOpenQIANKUN && !HELPSHOW) ||
            (!_isOpenQIANKUN && !showInQiankun),
          'fixed-beyound-short':
            (_isOpenQIANKUN && HELPSHOW) || (!_isOpenQIANKUN && showInQiankun)
        }"
        ref="footer"
      >
        <div class="slot-content">
          <slot name="footer"></slot>
        </div>
      </div>
    </div>
    <!-- 一般用于带有手机的左右布局 -->
    <div
      class="wrap"
      :class="{ 'pad-b-30': !hasSlotFooter }"
      ref="content"
      v-if="type === 'mobile'"
    >
      <div class="container clearfix">
        <div class="left pull-left">
          <slot name="left"></slot>
          <div class="border"></div>
        </div>
        <div class="right pull-left">
          <div class="border"></div>
          <slot name="right"></slot>
        </div>
      </div>
      <div
        class="footer"
        v-if="$slots.footer"
        :class="{
          'fixed-beyound-long':
            (_isOpenQIANKUN && !HELPSHOW) ||
            (!_isOpenQIANKUN && !showInQiankun),
          'fixed-beyound-short':
            (_isOpenQIANKUN && HELPSHOW) || (!_isOpenQIANKUN && showInQiankun)
        }"
        ref="footer"
      >
        <div class="slot-content">
          <slot name="footer"></slot>
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Mixins, Prop } from "vue-property-decorator";
import HelpStore from "@/store/modules/help";
import action from "@/shared/actions";
import Qiankun from "@/mixins/QianKunMixin";
import { IConfig } from "@/types";
@Component
export default class LayoutForm extends Mixins(Qiankun) {
  @Prop({ default: "normal" }) type?: string;
  get HELPSHOW() {
    return HelpStore.getHelp;
  }
  private beyound = false; //是否超出一屏
  private width = 0; //当超出一屏幕需要让footer固定在底部,需要拿一下相等的宽度
  private hasSlotFooter = !!this.$slots["footer"];
  private footerSticky(): void {
    const { offsetHeight, offsetTop, offsetWidth } = this.$refs.content as any;
    const screen = window.innerHeight || document.documentElement.clientHeight;
    this.width = offsetWidth;
    this.beyound = offsetHeight + offsetTop >= screen;
  }
  private showInQiankun = true; //帮助中心在qiankun微前端环境中是否展示状态
  mounted() {
    if (!this._isOpenQIANKUN) {
      action.onGlobalStateChange((config: IConfig) => {
        const { isShow } = config.config;
        this.showInQiankun = isShow;
      }, true);
    }
    // this.$nextTick(() => {
    //   this.footerSticky();
    // });
    // addEvent("resize", window, () => {
    //   this.$nextTick(() => {
    //     this.footerSticky();
    //   });
    // });
  }
  // beforeDestroy() {
  //   window.onresize = null;
  // }
}
</script>
<style lang="scss" scoped>
.pad-b-30 {
  padding-bottom: 30px !important;
}
.layout-form {
  .content {
    background-color: #ffffff;
    border-radius: 10px;
    padding: 30px 30px 0 30px;
  }
  .wrap {
    @extend .content;
    .container {
      .left {
        width: 406px;
        position: relative;
        min-height: 100%;
        .border {
          width: 1px;
          height: 100%;
          position: absolute;
          right: 0;
          top: 0;
          background-color: $color-third-line;
        }
      }
      .right {
        position: relative;
        width: calc(100% - 406px - 36px);
        padding-left: 30px;
        .border {
          width: 1px;
          height: 100%;
          position: absolute;
          left: -1px;
          top: 0;
          background-color: $color-third-line;
        }
      }
    }
  }
  .footer {
    width: 100%;
    position: sticky;
    bottom: -30px;
    z-index: 100;
    overflow: auto;
    padding: 20px 0 30px;
    background-color: #ffffff;
    box-sizing: border-box;
    .slot-content {
      background-color: #ffffff;
      border-top: 1px solid $color-third-line;
      padding-top: 20px;
    }
  }
  .fixed-beyound-long {
    left: 30px;
    padding-bottom: 30px;
    transition: width 0.3s cubic-bezier(1, 0.5, 0.8, 1);
  }
  .fixed-beyound-short {
    left: 30px;
    padding-bottom: 30px;
    transition: width 0.3s cubic-bezier(1, 0.5, 0.8, 1);
  }
}
</style>
