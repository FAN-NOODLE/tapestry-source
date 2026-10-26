<template>
  <!--分组弹窗 -->
  <van-popup
    v-model="visible"
    :close-on-click-outside="false"
    :close-on-click-action="false"
    position="center"
    @click-overlay="handleClose"
  >
    <video
      v-if="visible"
      :src="item.url"
      :controls="true"
      :poster="item.video_image_url"
      class="video"
      playsinline="true"
      x-webkit-airplay="allow"
      x5-playsinline
      webkit-playsinline
      object-fit="contain"
      ref="videoRef"
      @loadedmetadata="handleLoaded"
    ></video>
  </van-popup>
</template>
<script lang="ts">
import { Component, Vue,PropSync, Prop,Watch, Model } from "vue-property-decorator";
import { Popup } from "vant";
@Component({
  components: {
    [Popup.name]: Popup,
  },
})
export default class PreviewVideoPopup extends Vue {
  @Model("update:show", { type: Boolean }) show;
  @Prop({ type: Object, default: () => ({}) }) item: any;

  player = null
  visible = false;

  @Watch("show", { immediate: true })
  onVisibleChange(val: boolean) {
    this.visible = val
  }

  destroyed() {
    
    this.player = null;
    // this.$emit("update:show", false);
  }

  /**
   * 视频加载完成
   */
  handleLoaded() {
    this.player = this.$refs.videoRef;
    this.player.play();
  }

  handleClose() {
    this.visible = false
    this.$emit("close");
  }
  
}
</script>
<style lang="scss" scoped>
.van-popup {
  background-color: transparent;
}
.video {
  width: 100vw;
  height: 50vh;
  background-color: #000;
  display: block;
}
</style>