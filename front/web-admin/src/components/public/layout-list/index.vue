<template>
  <div class="layout-list" ref="LayoutList">
    <!-- 按钮部分，用于操作的按钮 -->
    <div
      class="buttons-content"
      v-show="hasButtons"
      :class="{ 'has-top-tabs': hasTopTabs }"
    >
      <slot name="buttons"></slot>
    </div>
    <!-- tab切换部分，用于整体tab切换 -->
    <div class="tabs-content" v-show="hasTopTabs">
      <slot name="top-tabs"></slot>
    </div>
    <!-- 搜索操作部分，用于表格筛选 -->
    <div
      class="search-content"
      v-show="hasSearch"
      :class="{ 'has-tabs': hasTabs }"
    >
      <slot name="search"></slot>
    </div>
    <!-- tab切换部分，用于整体tab切换 -->
    <div class="tabs-content" v-show="hasTabs">
      <slot name="tabs"></slot>
    </div>
    <!-- 批量操作部分，用于表格上方的操作按钮插槽 -->
    <div class="operate-content" v-show="hasOperate">
      <slot name="operate"></slot>
    </div>
    <div class="slide-box">
      <!-- 批量操作部分，用于表格上方的操作按钮插槽 -->
      <div class="slide-menu" v-show="hasSlide">
        <slot name="slide"></slot>
      </div>
      <!-- 表格部分，用于展示表格 -->
      <div class="slide-table" v-show="hasSlideTable" ref="table">
        <slot name="slide-table"></slot>
      </div>
    </div>
    <!-- 素材分组可左右拖动，手动调整宽度 -->
    <div class="resize-box">
      <!-- 批量操作部分，用于表格上方的操作按钮插槽 -->
      <div class="column-left" v-show="hasColumnLeft">
        <div class="resize-bar" @mouseup="barmouseup" ref="resizeBar"></div>
        <div class="resize-line"></div>
        <div class="resize-cont">
          <slot name="column-left"></slot>
        </div>
      </div>

      <!-- 表格部分，用于展示表格 -->
      <div
        class="column-right"
        :style="{ width: 'calc(100% - ' + leftWidth + 'px)' }"
        v-show="hasColumnRight"
      >
        <slot name="column-right"></slot>
      </div>
    </div>
    <!-- 表格部分，用于展示表格 -->
    <div class="table-content" ref="table">
      <slot name="table"></slot>
    </div>
    <!-- 表格底部部分，用于分页器 -->
    <div class="footer-content">
      <slot name="footer"></slot>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Watch } from "vue-property-decorator";
@Component
export default class LayoutList extends Vue {
  private hasTabs = !!this.$slots["tabs"];
  private hasTopTabs = !!this.$slots["top-tabs"];
  private hasButtons = !!this.$slots["buttons"];
  private hasSearch = !!this.$slots["search"];
  private hasOperate = !!this.$slots["operate"];
  private hasSlide = !!this.$slots["slide"];
  private hasSlideTable = !!this.$slots["slide-table"];
  private hasColumnLeft = !!this.$slots["column-left"];
  private hasColumnRight = !!this.$slots["column-right"];
  private over = false;
  private leftWidth = "214";
  private barmouseup() {
    this.leftWidth = (this.$refs.resizeBar as any).offsetWidth;
  }
}
</script>
<style lang="scss" scoped>
.layout-list {
  background-color: $color-white;
  padding: 30px;
  box-shadow: 0px 0px 20px 1px rgba(0, 0, 0, 0.08);
  // border-radius: 10px;
  .buttons-content {
    margin-bottom: 20px;
  }
  .has-top-tabs {
    margin-bottom: 7px;
  }
  .search-content {
    padding: 20px;
    margin-bottom: 20px;
    background-color: $color-background;
    border-radius: 4px;
  }
  .has-tabs {
    margin-bottom: 7px;
  }
  .operate-content {
    margin-bottom: 20px;
    color: $color-text-normal;
  }
  .table-content {
    margin-bottom: 20px;
  }
  .slide-box {
    display: flex;
  }
  .resize-box {
    display: flex;
    .column-right {
      margin-bottom: 20px;
      margin-left: 10px;
      overflow: hidden;
    }
    .column-left {
      position: relative;
      .resize-bar {
        width: 214px;
        height: 700px;
        resize: horizontal;
        cursor: ew-resize;
        cursor: col-resize;
        opacity: 0;
        overflow: scroll;
      }
      .resize-bar:hover ~ .resize-line,
      .resize-bar:active ~ .resize-line {
        border-left: 3px dashed transparent;
      }
      .resize-bar::-webkit-scrollbar {
        width: 200px;
        height: inherit;
      }
      /* 拖拽线 */
      .resize-line {
        position: absolute;
        right: 0;
        top: 0;
        bottom: 0;
        pointer-events: none;
        border-left: 3px solid transparent;
      }
      .resize-cont {
        margin-bottom: 20px;
        position: absolute;
        top: 0;
        right: 5px;
        bottom: 0;
        left: 0;
        overflow-x: hidden;
      }
    }
  }

  .slide-menu {
    margin-bottom: 20px;
    margin-right: 10px;
    width: 20%;
  }
  .slide-table {
    margin-bottom: 20px;
    width: 80%;
  }
  .footer-content {
    &:after {
      content: "";
      display: block;
      clear: both;
    }
    ::v-deep {
      .el-pagination {
        float: right;
      }
    }
  }
}
::v-deep {
  .clone-header {
    position: fixed;
    display: inline-block;
    top: 50px;
    left: 180px;
    text-align: left;
    text-indent: 20px;
    height: 36px;
    line-height: 36px;
    background-color: #ebebf1;
    color: #606266;
    z-index: 30;
    .el-table__body-wrapper {
      display: none;
      // z-index: 40;
    }
    .el-table__fixed-right {
      z-index: 40 !important;
    }
    .el-tabs__header {
      margin-bottom: 0 !important;
    }
  }
}
/* Firefox只有下面一小块区域可以拉伸 */
@supports (-moz-user-select: none) {
  .resize-bar:hover ~ .resize-line,
  .resize-bar:active ~ .resize-line {
    border-left: 1px solid #bbb;
  }
  .resize-bar:hover ~ .resize-line::after,
  .resize-bar:active ~ .resize-line::after {
    content: "";
    position: absolute;
    width: 16px;
    height: 16px;
    bottom: 0;
    right: -8px;
    // background-size: 100% 100%;
  }
}
</style>
