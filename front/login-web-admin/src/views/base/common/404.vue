<template lang="pug">
  .not-find
    img.not-image(src="@/assets/404.png")
    .no-page-title 抱歉，您查找的页面不存在
    p.empty-msg
      span 您可以：
      el-button(type="text" @click="goBack") 返回登录页
</template>
<script>
export default {
  name: "NotFind",
  methods: {
    goBack() {
      this.$router.push({
        name: "login"
      });
    }
  }
};
</script>
<style lang="sass" scoped>
.not-find
  width: 100%
  height: 100vh
  text-align: center
  .not-image
    padding-top: 130px
    margin: 0 auto
    width: 480px
    height: 380px
    display: block
  .no-page-title
    color: #606266
    font-size: 20px
    display: block
  .empty-msg
    text-align: center
    font-size: 14px
    line-height: 20px
    padding-top: 15px
    span
      color: #909399
      font-size: 14px
</style>
