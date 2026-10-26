<template>
  <div class="layout-list">
    <layout-card>
      <template slot="content">
        <div class="top-title m-b-20">将微信小程序绑定店铺，提升交易量</div>
        <div class="flex-between-c">
          <div class="card-content flex-c-s m-r-20">
            <img
              src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/mall-set/wechat.png"
              alt=""
              class="m-t-34"
            />
            <p class="desc-top m-t-15">
              微信注册小程序
            </p>
            <p class="desc">前往微信公众号注册小程序</p>
            <el-button
              type="primary"
              size="small"
              class="font-14 m-b-30 m-t-60"
              @click="registerWechat"
              >立即注册</el-button
            >
          </div>
          <div class="card-content flex-c-s m-r-20">
            <img
              src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/mall-set/fast.png"
              alt=""
              class="m-t-34"
            />
            <p class="desc-top m-t-15">
              快速注册小程序
            </p>
            <p class="desc">绑定已认证公众号，即可一键快速注册小程序，</p>
            <p class="desc">无需额外费用，立即得到审核结果</p>
            <el-button
              type="primary"
              size="small"
              class="font-14 m-b-30 m-t-40 m-t-20"
              @click="registerWechat"
              >去绑定认证公众号</el-button
            >
          </div>
          <div class="card-content flex-c-s ">
            <img
              src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/mall-set/registration.png"
              alt=""
              class="m-t-34"
            />
            <p class="desc-top m-t-15 ">
              代注册小程序
            </p>
            <p class="desc-last">1.填写审核资料</p>
            <p class="desc-last">2.小时内完成身份认证</p>
            <p class="desc-last">3.设置邮箱和密码</p>
            <el-button
              type="primary"
              size="small"
              class="font-14 m-b-30 m-t-20 last"
              @click="replaceRegister"
              >立即注册</el-button
            >
          </div>
        </div>
      </template>
    </layout-card>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Api from "@/api/setting/setting";
import LayoutCard from "./componments/layout-card/index.vue";
@Component({
  components: { LayoutCard }
})
export default class extends Vue {
  // 详情接口数据
  private detailData = {};
  // 详情接口没有数据
  private isfirstRegister = false;
  //  详情接口的状态
  private status = 0;
  // 路由跳转至微信公众号平台
  private registerWechat() {
    window.open("https://mp.weixin.qq.com", "_blank");
  }
  // 路由跳转至代注册小程序页面
  private replaceRegister() {
    if (JSON.stringify(this.detailData) != "{}") {
      // 如果详情信息不为空，返回申请结果界面
      console.log(this.status);
      this.$router.push({
        name: "myapp_setting_authorization-settings_store-set_register-result",
        // 携带参数isShowResult去判断审核的状态
        query: {
          isShowResult: this.status as any
          // isShowResult: "1",
        }
      });
    } else {
      this.$router.push({
        name: "myapp_setting_authorization-settings_store-set_replace-register"
      });
    }
  }
  // 调取详情接口
  private getregisterDetail() {
    Api.registerDetail({}).then(res => {
      if (JSON.stringify(res.data) != "{}") {
        // 详情的data
        this.detailData = res.data;
        // 申请的状态
        this.status = res.data[0].status;
      }
    });
  }
  mounted() {
    // 调取详情接口
    this.getregisterDetail();
  }
}
</script>
<style lang="scss" scoped>
div {
  user-select: none;
}
.top-title {
  color: #303133;
  font-size: 16px;
  font-weight: 500;
}
.flex-between-c {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.flex-c-s {
  display: flex;
  justify-content: center;
  flex-direction: column;
  align-items: center;
}
.card-content {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  flex: 1;
  img {
    width: 120px;
    height: 120px;
  }
  .desc-top {
    font-size: 16px;
    font-family: PingFangSC-Regular, PingFang SC;
    font-weight: 400;
    color: #303133;
    text-align: center;
    margin-bottom: 10px;
  }
  .desc {
    font-size: 12px;
    font-family: PingFangSC-Regular, PingFang SC;
    font-weight: 400;
    color: #909399;
    text-align: center;
  }
  .desc-last {
    width: 137px;
    font-size: 12px;
    font-weight: 400;
    color: #909399;
    text-align: left;
    margin-left: 40px;
  }
  .el-button {
    width: 180px;
    height: 40px !important;
    color: #fff;
    background-color: #708cf6;
  }
  .last {
    color: #606266;
    background-color: #fff;
  }
  .m-b-30 {
    margin-bottom: 30px;
  }
  .m-t-20 {
    margin-top: 20px;
  }
  .m-t-40 {
    margin-top: 40px;
  }
  .m-t-60 {
    margin-top: 60px;
  }
}
.m-t-34 {
  margin-top: 34px;
}
</style>
