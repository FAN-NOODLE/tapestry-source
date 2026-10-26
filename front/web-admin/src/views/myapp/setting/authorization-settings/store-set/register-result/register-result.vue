<template>
  <div class="supernatant">
    <layout-form>
      <template>
        <!-- 审核完成 -->
        <template v-if="isShowResult === wx_app_status.success">
          <template slot="content">
            <div class="flex-between-c">
              <div class="card-content flex-c-s ">
                <img
                  src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/mall-set/register-end.png"
                  alt=""
                  class="m-b-14"
                />
                <p class="desc ">
                  小程序申请已通过，请完善资料后进行授权操作<el-button
                    type="text"
                    class="text-button"
                    @click="goConsummate"
                    >去完善</el-button
                  >
                </p>
                <div>
                  <a :href="shop_auth_detail.auth_url" class="go-accredit">
                    <el-button
                      type="primary"
                      size="small"
                      class="font-14 m-t-30"
                      >进入微信授权页面</el-button
                    >
                  </a>
                </div>
              </div>
            </div>
          </template>
        </template>
        <!-- 审核中 -->
        <template v-if="isShowResult === wx_app_status.loading">
          <template slot="content">
            <div class="flex-between-c">
              <div class="card-content flex-c-s ">
                <img
                  src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/mall-set/registering.png"
                  alt=""
                  class="m-b-14"
                />
                <p class="desc ">
                  信息已提交，请留意微信消息进行后续操作，
                </p>
                <p class="desc ">
                  24小时内完成认证，并完善小程序信息
                </p>
                <el-button
                  type="primary"
                  size="small"
                  class="font-14 m-t-30 el-w-138 "
                  @click="dialogVisible = true"
                  >如何完成认证？</el-button
                >
                <el-dialog
                  title="如何完成认证"
                  :visible.sync="dialogVisible"
                  class="replace-register-dialog"
                >
                  <div>
                    <p class="m-b-20">
                      通过法人&企业主体校验，平台向法人微信下发模板消息。法人需在24
                      小时内点击消息，进行身份证信息与人脸识别信息收集。
                    </p>
                    <img
                      src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/mall-set/dialog.png"
                      alt=""
                      width="100%"
                    />
                  </div>
                  <span slot="footer" class="dialog-footer">
                    <el-button
                      type="primary"
                      @click="dialogVisible = false"
                      class="el-w-36"
                      >确 定</el-button
                    >
                  </span>
                </el-dialog>
              </div>
            </div>
          </template>
        </template>
        <!-- 审核失败-->
        <template v-if="isShowResult === wx_app_status.fail">
          <template slot="content">
            <div class="flex-between-c">
              <div class="card-content flex-c-s ">
                <img
                  src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/mall-set/register-defeat .png"
                  alt=""
                  class="m-b-14"
                />
                <p class="desc ">
                  审核失败：
                </p>
                <p
                  class="desc"
                  v-if="
                    detailData && detailData[0] && detailData[0].fail_reason
                  "
                >
                  {{ detailData[0].fail_reason }}
                </p>

                <el-button
                  type="primary"
                  size="small"
                  class="font-14 m-t-30  "
                  @click="goRegister"
                  >重新提交申请</el-button
                >
              </div>
            </div>
          </template>
        </template>
      </template>
    </layout-form>
  </div>
</template>
<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
import UploadImage from "@/components/business/upload-image/index.vue";
//调取数据接口
import Api from "@/api/setting/setting";
@Component({})
export default class FormDemo extends Vue {
  // 注册状态
  private wx_app_status = {
    success: 1, //成功
    loading: 2, //审核中
    fail: 3 //失败
  };
  // 授权状态下小程序的数据
  private shop_auth_detail = {};
  //微信授权回调地址
  private back_url = "";
  private auth_url = "";
  // 审核状态
  private status = {};
  // 展示结果
  private isShowResult = {};
  // 弹窗
  private dialogVisible = false;
  // 详情接口数据
  private detailData = {};
  // 审核完成--------去完善弹窗
  private goConsummate() {
    this.$router.push({
      name: "myapp_setting_authorization-settings_store-set_store-sets",
      query: {
        isShowModal: "1" // 跳转首页------1表示显示弹窗
      }
    });
  }

  // 审核失败status=fail-------重新提交申请返回注册界面
  private goRegister() {
    //第二次返回注册界面
    this.$router.push({
      name: "myapp_setting_authorization-settings_store-set_replace-register"
    });
  }
  // 调取详情接口
  private getregisterDetail() {
    Api.registerDetail({}).then(res => {
      if (JSON.stringify(res.data) != "{}") {
        // 详情的data
        this.detailData = res.data;
        this.isShowResult = res.data[0].status;
      }
    });
  }
  // 获取商城小程序授权详情
  private getshopAuthDetail() {
    Api.shopAuthDetail({
      back_url: this.back_url
    }).then(res => {
      this.shop_auth_detail = res.data;
    });
  }
  private listenPageRefresh() {
    const that = this;
    window.onbeforeunload = function(e: any) {
      e = e || window.event;
      if (e) {
        //调取API
        that.getregisterDetail();
      }
    };
  }
  mounted() {
    const back_url = this.$router.resolve({
      path: "/myapp/setting/authorization-settings/store-set/empty"
    }).href;
    this.back_url = window.location.origin + back_url;
    // 调取详情接口
    this.getregisterDetail();
    // 获取商城小程序授权详情接口
    this.getshopAuthDetail();
    // 页面展示isShowResult传过来的状态
    this.isShowResult = Number(this.$route.query.isShowResult);
    console.log("数据类型" + typeof "isShowResult");
    console.log(this.isShowResult);
    //监听页面刷新
    this.listenPageRefresh();
  }
}
</script>
<style lang="scss" scoped>
// 已完成样式
div {
  user-select: none;
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
  width: 1260px;
  height: 520px;
  background: #ffffff;
  border-radius: 10px;
  flex: 1;
  .img {
    width: 120px;
    height: 120px;
  }
  .desc {
    font-size: 14px;
    font-family: PingFangSC-Regular, PingFang SC;
    font-weight: 400;
    color: #909399;
    text-align: center;
  }
  .text-button {
    margin-left: 5px;
  }
  .m-t-10 {
    margin-top: 10px;
  }
  .m-t-30 {
    margin-top: 30px;
  }
  .m-b-14 {
    margin-bottom: 14px;
  }
}
// 弹窗样式调整
/deep/.el-dialog__body {
  padding: 10px 30px 18px 30px !important;
}
/deep/.replace-register-dialog > .el-dialog {
  max-width: 971px !important;
  max-height: 777px !important;
}
/deep/.el-dialog {
  .p {
    margin-top: -10px;
    font-size: 14px;
    font-weight: 400;
    color: #606266;
  }
  .el-button {
    width: 68px !important;
  }
}
// 输入框和选择框样式调整
.cy-select,
.cy-input {
  width: 400px !important;
}
/deep/.el-w-138 {
  width: 138px !important;
}
/deep/.el-w-36 {
  height: 36px !important;
  line-height: 13px;
  text-align: center;
}

// 输入框样式调整
/deep/.el-form-item__label {
  width: 114px !important;
}
/deep/.el-form-item__content {
  margin-left: 114px !important;
}
/deep/.desc.m-t-5.color-text-second {
  margin-left: 114px !important;
}
/deep/.fixed-beyound-short {
  padding-bottom: 0px !important;
}
/deep/.footer {
  bottom: 0px !important;
  padding: 0px !important;
}
/deep/.content {
  padding: 30px !important;
}
</style>
