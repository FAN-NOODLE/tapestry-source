<template>
  <!--分组弹窗 -->
  <van-popup
    v-model="visible"
    :close-on-click-outside="false"
    :close-on-click-action="false"
    position="top"
  >
    <video
      v-if="visible"
      :src="item.url"
      :controls="true"
      class="video"
      playsinline="true"
      x-webkit-airplay="allow"
      x5-playsinline
      webkit-playsinline
      object-fit="contain"
      ref="videoRef"
      @loadedmetadata="handleLoaded"
      @click="handleClose"
    ></video>
  </van-popup>
</template>
<script lang="ts">
import { Component, Vue,PropSync, Prop,Watch } from "vue-property-decorator";
import { Popup } from "vant";
@Component({
  components: {
    [Popup.name]: Popup,
  },
})
export default class EditTextPopup extends Vue {
  @PropSync("show", { type: Boolean }) visible;
  @Prop({ type: Object, default: () => ({}) }) item: any;

  player = null

  @Watch("visible", { immediate: true })
  onVisibleChange(val: boolean) {
    if (val) {
      this.$nextTick(() => {
        // (this.$refs.videoRef as any).play();
      });
    }
  }

  destroyed() {
    this.player && this.player.dispose();
    this.player = null;
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
}
</style>