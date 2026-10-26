<template>
  <div class="preview-image">
    <i class="el-icon-close close-btn" aria-hidden="true" @click="close"></i>
    <i
      v-if="left_show"
      class="el-icon-arrow-left left-btn"
      aria-hidden="true"
      @click="left"
    ></i>
    <i
      v-if="right_show"
      class="el-icon-arrow-right right-btn"
      aria-hidden="true"
      @click="right"
    ></i>
    <img class="img-list" alt="" :src="urls[pre_index]" />
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
@Component({})
export default class PreviewImage extends Vue {
  @Prop() urls!: Array<string>;
  @Prop({ default: 0 }) index?: number;

  mounted() {
    this.pre_index = this.index || 0;
    this.pre_index === 0 && (this.left_show = false);
    this.pre_index === this.urls.length - 1 && (this.right_show = false);
  }
  private left_show = true;
  private right_show = true;
  private pre_index = 0;

  private left() {
    const index = this.pre_index;
    this.right_show = true;
    this.pre_index = index - 1;
    if (index === 1) {
      this.left_show = false;
    }
  }
  private right() {
    const r = this.urls.length;
    this.left_show = true;
    const index = this.pre_index + 1;
    this.pre_index = index;
    if (r - 1 == index) {
      this.right_show = false;
    }
  }
  private close() {
    this.$emit("close");
  }
}
</script>
<style lang="scss" scoped>
.preview-image {
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  position: fixed;
  z-index: 9999;
  top: 0;
  left: 0;
  overflow: hidden;
  .btn {
    position: fixed;
    font-size: 56px;
    cursor: pointer;
    color: #fff;
  }
  .close-btn {
    @extend .btn;
    position: fixed;
    font-size: 40px;
    width: 40px;
    height: 40px;
    right: 20px;
    top: 12px;
    line-height: 40px;
    z-index: 10000;
  }
  .right-btn {
    @extend .btn;
    top: 50%;
    right: 25px;
    margin-top: -70px;
    width: 55px;
    height: 122px;
    line-height: 122px;
  }

  .left-btn {
    @extend .btn;
    top: 50%;
    left: 25px;
    margin-top: -70px;
    width: 55px;
    height: 122px;
    line-height: 122px;
  }
  .img-list {
    position: absolute;
    width: 100%;
    height: 100%;
    object-fit: contain;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    margin: auto;
  }
}
</style>
