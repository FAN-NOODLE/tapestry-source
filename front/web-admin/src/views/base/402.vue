<template>
  <div class="page-notfound">
    <img
      class="no-page-img"
      src="https://oss-crs.vchangyi.com/wechat/no-contacts-icon.png"
      alt=""
    />
    <div class="no-page-title">抱歉，您没有访问该页面的权限</div>
    <div class="operation-box">
      <div class="text-one">您可以：</div>
      <div class="text-two" @click="goBack">返回上一页</div>
      <div class="text-three" @click="goHome">前往平台首页</div>
    </div>
  </div>
</template>
<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
@Component({ name: "402" })
export default class extends Vue {
  private goBack() {
    this.$router.go(-1);
  }
  private goHome() {
    this.$router.push({
      name: "myapp_home"
    });
  }
}
</script>
<style lang="scss" scoped>
.page-notfound {
  background-color: $color-background;
  position: relative;
  height: 100vh;
  display: flex;
  align-content: center;
  justify-content: center;
  flex-direction: column;
  text-align: center;
  .no-page-img {
    display: block;
    height: 275px;
    width: 360px;
    left: 0;
    right: 0;
    margin: auto;
    margin-bottom: 0;
    margin-top: 0;
  }
  .no-page-title {
    color: $color-text-normal;
    font-size: 20px;
    display: block;
  }
  .operation-box {
    font-size: 0;
    height: 14px;
    line-height: 1;
    display: inline-block;
    margin-top: 20px;

    .text-one {
      color: $color-normal;
      display: inline-block;
      font-size: 14px;
    }

    .text-two {
      color: $color-main;
      display: inline-block;
      padding: 0 30px;
      cursor: pointer;
      font-size: 14px;
    }

    .text-three {
      color: $color-main;
      display: inline-block;
      padding: 0 30px;
      cursor: pointer;
      font-size: 14px;
      border-left: 0.5px solid #bdbdbd;
    }
  }
}
</style>
