<template lang="pug">
  .slide-bar-wrap
    .slider-bar
      .slide-top
        i.el-icon-close.close(@click="close")
      .slider-header
        span 完成拼图验证
        .refresh-block
          i.el-icon-refresh
          span 换一张
      slide-verify.slide-plugin(
        ref="slideblock" 
        :imgs="imgs" 
        sliderText="向右滑动完成拼图" 
        @success="success"
        @fail="onFail")
</template>
<script>
export default {
  name: "SlideBar",
  data() {
    return {
      patch_left: 0,
      imgs: [
        require("@/assets/slide/bg-slide-1.png"),
        require("@/assets/slide/bg-slide-2.png"),
        require("@/assets/slide/bg-slide-3.png")
      ]
    };
  },
  methods: {
    success() {
      // this.$message.success("验证成功");
      this.$emit("success");
    },
    onFail() {
      this.$message.error("验证失败");
      this.$refs["slideblock"].reset();
      this.$emit("fail");
    },
    close() {
      this.$emit("close");
    },
    // 更改插件内部的样式
    updateStone() {
      let mask = document.createElement("div");
      let icon = document.createElement("i");
      mask.className = "slide-icon-mask";
      icon.className = "el-icon-arrow-right slide-right";
      document
        .querySelector(".slide-verify-slider-mask-item")
        .appendChild(mask);
      mask.appendChild(icon);
      mask.addEventListener("mouseover", () => {
        icon.className = "el-icon-arrow-right slide-right-white";
      });
      mask.addEventListener("mouseleave", () => {
        icon.className = "el-icon-arrow-right slide-right";
      });
    }
  },
  mounted() {
    this.updateStone();
  }
};
</script>
<style lang="sass" scoped>
/deep/.slide-verify-slider-mask-item
  border-radius: 50%
/deep/.slide-verify-slider
  border-radius: 20px
  background-color: rgb(243,249,255) !important
  // margin-top: 6px !important
/deep/.slide-verify-slider-mask
  border-radius: 20px 0 0 20px
/deep/.slide-verify-slider-text
  color: #1989FA !important
/deep/.slide-verify-slider-mask-item-icon
  display: none
/deep/.slide-icon-mask
  width: 100%
  height: 100%
/deep/.slide-right
  color: #1989FA
  font-weight: 600
/deep/.slide-right-white
  color: #FFFFFF
  font-weight: 600
/deep/.slide-verify-refresh-icon
  top: -37px !important
  right: 0 !important
  width: 62px !important
  z-index: 201
  opacity: 0
.slide-bar-wrap
  width: 100%
  height: 100%
  background-color: rgba(0,0,0,0.7)
  position: fixed
  z-index: 100
.slider-bar
  width: 420px
  height: 300px
  background-color: #FFFFFF
  position: absolute
  border-radius: 8px
  left: 0
  right: 0
  top: 0
  bottom: 0
  margin: auto
  .slide-top
    position: relative
    height: 16px
    padding-top: 10px
    .close
      position: absolute
      right: 15px
      top: 15px
      cursor: pointer
  .slider-header
    margin-bottom: 10px
    margin-left: 50px
    position: relative
    .refresh-block
      position: absolute
      top: 0
      right: 60px
      cursor: pointer
      color: #1989FA
      z-index: 200
      span
        color: #1989FA
    span
      color: #a0a3a8
      font-size: 14px
  .slide-plugin
    margin-left: 50px
  .slider-footer
    padding-top: 10px
    text-align: center
    color: #a0a3a8
    font-size: 14px
</style>
