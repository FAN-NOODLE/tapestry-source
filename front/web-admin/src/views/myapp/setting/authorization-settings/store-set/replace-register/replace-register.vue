<template>
  <div class="supernatant">
    <layout-form>
      <!-- 代注册---------信息填写界面 -->
      <template>
        <template slot="content">
          <cy-form ref="cform" :rules="rules" :model="ruleForm">
            <cy-form-item
              label="法人姓名"
              prop="legal_persona_name"
              desc="法人姓名需与绑定银行卡的姓名保持一致"
              required
            >
              <cy-input
                placeholder="请填写法人姓名"
                v-model="ruleForm.legal_persona_name"
              ></cy-input>
            </cy-form-item>
            <cy-form-item
              label="法人微信"
              prop="legal_persona_wechat"
              desc="需要绑定法人银行卡号的微信号"
              required
            >
              <cy-input
                placeholder="请填写法人微信号"
                v-model="ruleForm.legal_persona_wechat"
              ></cy-input>
            </cy-form-item>
            <cy-form-item
              label="企业名称"
              prop="name"
              desc="企业名称需要与工商部门登记信息一致"
              required
            >
              <cy-input
                placeholder="请填写工商营业执照上的企业名称"
                v-model="ruleForm.name"
              ></cy-input>
            </cy-form-item>
            <cy-form-item label="企业代码类型" prop="code_type" required>
              <div class="template">
                <cy-select
                  v-model="ruleForm.code_type"
                  size="medium"
                  placeholder="选择企业代码类型"
                  trigger="change"
                  required
                >
                  <cy-option
                    v-for="item in ruleForm.code_types"
                    :key="item.id"
                    :label="item.label"
                    v-model="item.id"
                  ></cy-option>
                </cy-select>
              </div>
            </cy-form-item>
            <cy-form-item label="企业代码" prop="code" required>
              <cy-input
                placeholder="请填写企业代码"
                v-model="ruleForm.code"
              ></cy-input>
            </cy-form-item>
          </cy-form>
        </template>
        <template slot="footer">
          <el-button
            type="primary"
            size="small"
            class="font-14"
            @click="goCommit(ruleForm)"
            >提交申请</el-button
          >
          <el-button size="small" class="font-14" @click="goBack()"
            >取消</el-button
          >
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
  // 页面状态
  private status = 0;
  // 默认第一次注册 --------详情的data没有数据
  private isfirstRegister = true;
  // 弹窗
  private dialogVisible = false;
  // 详情接口数据
  private detailData = {};
  // 第二次进入的id
  private id = "";
  // 输入框数据
  private ruleForm = {
    id: "", //区分是第一次申请还是第二次申请
    legal_persona_name: "", //法人姓名
    legal_persona_wechat: "", //法人微信
    name: "", //企业名称
    code_type: "", //企业代码类型
    code: "", //企业代码
    //企业代码类型
    code_types: [
      {
        id: 1,
        label: "统一社会信用代码（18 位）"
      },
      {
        id: 2,
        label: "组织机构代码（9 位 xxxxxxxx-x）"
      },
      {
        id: 3,
        label: "营业执照注册号(15 位)"
      }
    ]
  };
  // 校验规则
  private rules = {};
  //返回取消
  private goBack() {
    this.$router.push({
      name: "myapp_setting_authorization-settings_store-set_register-app"
    });
  }

  // 提交
  private goCommit() {
    // 详情为数据为空的时候-------------------第一次申请
    (this.$refs.cform as any).validate((valid: boolean) => {
      if (valid) {
        let params;
        if (this.isfirstRegister) {
          params = {
            legal_persona_name: this.ruleForm.legal_persona_name, //法人姓名
            legal_persona_wechat: this.ruleForm.legal_persona_wechat, //法人微信
            name: this.ruleForm.name, //企业名称
            code_type: this.ruleForm.code_type, //企业代码类型
            code: this.ruleForm.code //企业代码
          };
        } else {
          this.isfirstRegister = false;
          params = {
            id: this.id, //第二次返回id
            legal_persona_name: this.ruleForm.legal_persona_name, //法人姓名
            legal_persona_wechat: this.ruleForm.legal_persona_wechat, //法人微信
            name: this.ruleForm.name, //企业名称
            code_type: this.ruleForm.code_type, //企业代码类型
            code: this.ruleForm.code //企业代码
          };
        }
        Api.registerApply(params).then(res => {
          this.isfirstRegister = false;
          // 默认第一次注册成功跳转审核中
          this.$router.push({
            name:
              "myapp_setting_authorization-settings_store-set_register-result",
            query: {
              isShowResult: "2" // 审核结果页面------2表示等待审核结果
            }
          });
        });
      }
    });
  }

  // 调取详情接口
  private getregisterDetail() {
    Api.registerDetail({}).then(res => {
      if (JSON.stringify(res.data) != "{}") {
        // 详情的data
        this.detailData = res.data;
        // 详情的data没有数据
        this.isfirstRegister = false;
        // 申请的状态
        this.status = res.data[0].status;
        // 第二次申请或者申请失败重新申请时返回的id    第一次不返
        this.id = res.data[0].id;
        // 把第一次申请的数据直接传上到输入框中
        this.ruleForm.legal_persona_name = res.data[0].legal_persona_name;
        (this.ruleForm.legal_persona_wechat = res.data[0].legal_persona_wechat), //法人微信
          (this.ruleForm.name = res.data[0].name), //企业名称
          (this.ruleForm.code_type = res.data[0].code_type), //企业代码类型
          (this.ruleForm.code = res.data[0].code); //企业代码
      } else {
        this.isfirstRegister = true;
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
