<template>
  <div class="com-open-video">
    <transition name="fade">
      <div class="img-view">
        <span
          class="close-btn"
          @click="closeModel"
        ><i
            class="el-icon-close"
            aria-hidden="true"
          ></i></span>
        <div class="video-size">
          <video
            :src="video_url"
            controls="controls"
            autoplay
          />
        </div>
      </div>
    </transition>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";

@Component({
  components: {}
})
export default class extends Vue {
  @Prop() video_url?: any;
  // ------------------生命周期-----------
  closeModel() {
    this.$emit("closeModel");
  }
}
</script>
<style lang="scss" scoped>
.com-open-video {
  .img-view {
    width: 100%;
    height: 100%;
    position: fixed;
    z-index: 9999;
    top: 0;
    left: 0;
    background: rgba(0, 0, 0, 0.7);
    overflow: hidden;
  }

  .close-btn {
    font-size: 20px;
    position: fixed;
    width: 40px;
    height: 40px;
    right: 20px;
    top: 12px;
    text-align: center;
    line-height: 40px;
    color: #bbb;
    cursor: pointer;
    z-index: 99991;
  }

  .video-size {
    width: 610px;
    height: 346px;
    margin: 0 auto;
  }

  video {
    width: 100%;
    height: 100%;
    margin: 200px 0;
  }
}</style>
