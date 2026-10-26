<template>
  <div
    class="default"
    :style="{
      paddingTop: paddingTop + 'px',
      paddingBottom: paddingBottom + 'px',
    }"
  >
    <img
      :class="{
        system: picture == 'system',
        nodata: picture == 'noData',
        noContent: picture == 'noContent',
        emptyList: picture == 'emptyList',
        noInfo: picture == 'noInfo',
        noDisplayTask: picture == 'noDisplayTask' || picture == 'noSearchData',
      }"
      :src="imageUrl"
      alt=""
    />
    <div
      v-if="picture == 'noInfo'"
      class="noInfo-p"
      :style="{ bottom: 118 + Number(paddingBottom) + 'px' }"
    >
      <div>这里什么都没有哦～</div>
      <div>请通过分享卡片进入指定页面</div>
    </div>
    <p
      v-else
      :class="{ nopadding: picture == 'emptyList' || picture == 'noData' }"
    >
      {{ text }}
    </p>
  </div>
</template>
<script lang="ts">
const system = require("@/assets/default/default.png");
const noData = require("@/assets/material-library/no-data.png");
const noContent = require("@/assets/default/no-content.png");
const emptyList = require("@/assets/default/empty-list.png");
const noInfo = require("@/assets/default/no-info.png");
const noDisplayTask = require("@/assets/display-task/no-data.png");
const noSearchData = require("@/assets/display-task/no-search-data.png");
const onlineOrder = require("@/assets/online-order/empty.png");
import { Prop, Component, Vue } from "vue-property-decorator";
@Component({
  components: {},
})
export default class extends Vue {
  @Prop({ default: "noData", type: String }) picture!: string;
  @Prop({ default: "暂无数据", type: String }) text!: string;
  @Prop({ default: "100", type: String }) paddingTop!: string;
  @Prop({ default: "0", type: String }) paddingBottom!: string;
  // 计算属性
  get imageUrl(): string {
    return {
      system,
      noData,
      noContent,
      emptyList,
      noInfo,
      noDisplayTask,
      noSearchData,
      onlineOrder
    }[this.picture] || noContent 
  }
}
</script>

<style lang="less" scoped>
.default {
  width: 100%;
  height: 50vh;
  background: transparent;
  text-align: center;
  padding-top: 60px;
  box-sizing: border-box;
  position: relative;
  img.system {
    width: 166px;
    height: 160px;
    margin: 0 auto;
  }
  img.nodata {
    width: 90px;
    height: 90px;
    margin: 0 auto;
  }
  img.noDisplayTask {
    width: 86px;
    height: 86px;
    margin: 0 auto;

    &~p {
      padding: 0;
      color: #4a4a4a;
    }
  }
  img.emptyList {
    width: 200px;
    height: 200px;
    margin: 0 auto;
  }
  img.noInfo {
    width: 100%;
    height: 480px;
    margin: 0 auto;
  }
  .noInfo-p {
    position: absolute;
    left: 0;
    right: 0;
    margin: auto;
    div {
      height: 24px;
      font-size: 17px;
      font-weight: 600;
      color: #33363c;
      line-height: 24px;
      margin-top: 10px;
    }
  }
  p {
    padding-top: 26px;
    color: #a1a1a1;
    line-height: 20px;
    font-size: 13px;
    &.nopadding {
      padding-top: 0;
    }
  }
}
</style>
