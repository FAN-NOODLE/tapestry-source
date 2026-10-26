<template lang="pug">
  .account-overdue
    .account-content
      img.close(src="@/assets/close-icon.png" @click="close")
      .header-title 到期提示
      .body-content
        .body-text 你订购的畅鲟管理系统已到期，如需恢复使用，请续费。
        .gray-text 
          span 官方咨询电话：
          span.phone-official 400-860-6961
      el-button.btn(size="medium" @click="purchase") 立即续费 
</template>
<script>
export default {
  methods: {
    purchase() {
      window.open(
        "http://xueyuan.vchangyi.com/bin/view/CRS%E5%8A%9F%E8%83%BD%E8%AF%B4%E6%98%8E/%E6%99%BA%E6%85%A7%E5%AF%BC%E8%B4%AD/%E6%99%BA%E6%85%A7%E5%AF%BC%E8%B4%AD%E7%AE%A1%E7%90%86%E5%90%8E%E5%8F%B0/%E6%99%BA%E6%85%A7%E5%AF%BC%E8%B4%AD%E4%BB%98%E8%B4%B9%E8%AE%A2%E9%98%85/",
        "_blank"
      );
    },
    close() {
      this.$emit("close");
    }
  }
};
</script>
<style lang="sass" scoped>
.account-overdue
  width: 100%
  height: 100vh
  position: fixed
  z-index: 100
  background-color: rgba(0,0,0,0.7)
  .account-content
    width: 420px
    height: 209px
    position: absolute
    top: 0
    left: 0
    right: 0
    bottom: 0
    margin: auto
    padding: 30px 30px 30px 30px
    border-radius: 10px
    background-color: #FFFFFF
    box-sizing: border-box
    .close
      width: 14px
      height: 14px
      float: right
      cursor: pointer
      margin-top: 6px
    .header-title
      display: inline-block
      color: #303133
      font-size: 18px
      line-height: 26px
      font-weight: bold
    .body-content
      padding-top: 20px
      margin-bottom: 20px
      &:after
        display: block
        content: ""
        clear: both
      .body-text
        color: #606266
        font-size: 14px
        line-height: 20px
      .gray-text
        color: #606266
        line-height: 26px
        font-size: 14px
        margin-top: 5px
        .phone-official
          font-weight: 600
          font-size: 16px
          line-height: 20px
          color: #303133
    // 立即续费按钮
    .btn
      color: #fff
      background-color: #708DF5
      border-color: #708DF5
      width: 96px
      height: 36px
      float: right
      border-radius: 3px
    .btn:hover, .btn:focus
      color: #fff
      background-color: #708DF5
      border-color: #708DF5
</style>
