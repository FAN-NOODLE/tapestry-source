<template>
  <div>
    <div class="footer-bar flex-between-c" v-show="show">
      <div class="flex-c-c width-82-p">
        <div class="left flex-c-c">
          <div>
            <p class="title-color">您正在试用天鲟CRS管理系统，如需</p>
            <p class="title-color">正常使用，请订购。</p>
            <p class="desc-color m-t-10">
              <span>
                使用到期时间：
              </span>
              <span>
                {{
                  USERINFO.account_expire_at | formatDate("yyyy-mm-dd hh:mm:ss")
                }}
              </span>
            </p>
          </div>
          <div>
            <el-button
              type="primary"
              class="maincolor"
              size="small"
              @click="payNow"
              >立即订购</el-button
            >
          </div>
        </div>
        <div class="line"></div>
        <div class="middle flex-c-c">
          <div class="img-con border-none">
            <img src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/default/business-contact.png" alt="" />
          </div>
          <div>
            <p class="title-color">添加官方顾问微信</p>
            <p class="desc-color m-t-10">1对1专业咨询服务</p>
          </div>
        </div>
        <div class="line"></div>
        <div class="right flex-c-c">
          <div class="img-con">
            <img src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/default/footer.png" alt="" />
          </div>
          <div>
            <p class="title-color">更多经典行业案例</p>
            <p class="desc-color m-t-10">
              全球100+知名企业首选，畅移诚邀优质品牌入驻
            </p>
          </div>
        </div>
      </div>
      <div class="packup" @click="packUp">
        <span class="white">收起</span>
        <i aria-hidden="true" class="el-icon-arrow-right white"></i>
      </div>
    </div>
    <div @click="packUp" v-show="!show" class="posf">
      <span class="white">展开</span>
      <i aria-hidden="true" class="el-icon-arrow-left white"></i>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import UserModule from "@/store/modules/user";

@Component
export default class extends Vue {
  get USERINFO() {
    return UserModule.getUserInfo;
  }
  private show = true;
  //收起
  private packUp() {
    this.show = !this.show;
  }
  //立即订购
  private payNow() {
    window.open(
      "http://xueyuan.vchangyi.com/bin/view/CRS%E5%8A%9F%E8%83%BD%E8%AF%B4%E6%98%8E/%E6%99%BA%E6%85%A7%E5%AF%BC%E8%B4%AD/%E6%99%BA%E6%85%A7%E5%AF%BC%E8%B4%AD%E7%AE%A1%E7%90%86%E5%90%8E%E5%8F%B0/%E6%99%BA%E6%85%A7%E5%AF%BC%E8%B4%AD%E4%BB%98%E8%B4%B9%E8%AE%A2%E9%98%85/",
      "_blank"
    );
  }
}
</script>
<style lang="scss" scoped>
div {
  user-select: none;
}
.footer-bar {
  width: 100%;
  padding: 30px 30px 30px 56px;
  background: rgba(18, 28, 66, 1);
  box-sizing: border-box;
  opacity: 0.85;
}
.flex-c-c {
  display: flex;
  justify-content: space-around;
  align-items: center;
}
.title-color {
  font-size: 16px;
  color: #fff;
  font-weight: bolder;
}
.desc-color {
  font-size: 14px;
  color: #c2c4ce;
}

.maincolor {
  background: #708df5 !important;
  border: 1px solid #708df5 !important;
  color: #fff;
  margin-left: 20px;
}
.line {
  height: 70px;
  width: 1px;
  background: #c2c4ce;
}
.m-t-10 {
  margin-top: 10px;
}
.img-con {
  width: 70px;
  height: 70px;
  border: 1px dashed #fff;
  margin-right: 15px;
  img {
    width: 100%;
    height: 100%;
  }
}
.border-none {
  border: none;
}
.width-82-p {
  width: 85%;
}
.white {
  color: #fff;
}
.posr {
  position: relative;
}
.flex-between-c {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.packup {
  cursor: pointer;
  margin-right: 30px;
  height: 71px;
  line-height: 71px;
}
.width-none {
  width: 0%;
  transition: width 0.5s linear;
}
.width-100 {
  width: 100%;
  transition: width 0.5s linear;
}
.posf {
  position: fixed;
  bottom: 0;
  right: 0;
  height: 71px;
  line-height: 71px;
  width: 5%;
  background: rgba(18, 28, 66, 1);
  opacity: 0.85;
  text-align: center;
  padding: 30px 30px 30px 56px;
  cursor: pointer;
}
</style>
