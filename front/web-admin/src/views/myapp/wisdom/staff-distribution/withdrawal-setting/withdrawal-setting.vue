<template>
  <div class="withdrawal-setting">
    <template>
      <layout-list title_name="提现设置" class="main-con">
        <template slot="search">
          <el-dialog title="保存确认" :visible.sync="dialogVisible" width="30%">
            <span>
              您填写的项目编号{{ form.transfer_no }}，对应的项目名称为{{
                project_name
              }}，请确认是否正确
            </span>
            <template slot="footer">
              <el-button @click="dialogVisible = false">取消</el-button>
              <el-button type="primary" @click="onSave">确定</el-button>
            </template>
          </el-dialog>
          <el-form
            class="demo-ruleForm"
            :model="form"
            :rules="rules"
            label-position="right"
            ref="form"
            label-width="83px"
            :validate-on-rule-change="false"
          >
            <el-form-item
              class="header-item"
              label="提现方式"
              prop="transfer_no"
              required
            >
              <div class="content">
                <el-radio v-model="form.payType" label="1">第三方转账</el-radio>
                <el-input
                  v-model="form.transfer_no"
                  placeholder="请输入转账项目编号"
                  size="medium"
                ></el-input>
              </div>
              <div class="tips" style="margin-left: 125px;">
                请联系畅移商务获取转账项目编号
              </div>
            </el-form-item>
            <el-form-item label="可提现人员" prop="duration_time">
              <div class="content" style="margin-top:10px">
                <el-checkbox-group v-model="form.type" style="height:17px">
                  <el-checkbox
                    v-for="item in user_list"
                    :key="item.key"
                    :label="item.key"
                  >
                    {{ item.name }}
                  </el-checkbox>
                </el-checkbox-group>
              </div>
              <div class="tips">
                设置为可提现人员的类型可在小程序端点击提现并进入提现流程
              </div>
            </el-form-item>
          </el-form>
          <div class="line"></div>
        </template>
        <template slot="footer">
          <el-button type="primary" @click="onCheck('form')" size="small">
            保存
          </el-button>
        </template>
      </layout-list>
    </template>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Api from "@/api/wisdom/dataCompass";
import UserModule from "@/store/modules/user";
@Component
export default class extends Vue {
  private dialogVisible = false;
  private user_list = [
    { name: "企业导购", key: 3 },
    { name: "社群导购", key: 2 },
    { name: "门店导购", key: 1 },
    { name: "企业管理", key: 4 },
    { name: "分销合伙人", key: 5 }
  ];
  private form = {
    payType: "1",
    transfer_no: "",
    type: []
  };
  private rules = {
    transfer_no: [
      { required: true, message: "请输入转帐项目编号", trigger: "blur" },
      { required: true, message: "请输入转帐项目编号", trigger: "change" }
    ]
  };
  private project_name = "";

  get USER() {
    return UserModule.getUserInfo;
  }

  // methods
  private getDetail() {
    Api.getSettingDetail().then((res: any) => {
      if (res.data) {
        this.form.transfer_no = res.data.transfer_no || "";
        this.form.type = res.data.type || [];
      }
    });
  }

  private setUserList() {
    if (this.USER.version === 3) {
      this.user_list = [
        { name: "企业导购", key: 3 },
        { name: "社群导购", key: 2 },
        { name: "企业管理", key: 4 },
        { name: "分销合伙人", key: 5 }
      ];
    }
  }

  // 保存按钮事件
  private onCheck(form: string) {
    // 校验ref标签引用
    (this.$refs[form] as any).validate((valid: boolean) => {
      if (valid) {
        Api.getNameByCode({
          transfer_no: this.form.transfer_no
        }).then(res => {
          if (res.data && res.data.project_name) {
            this.project_name = res.data.project_name;
            this.dialogVisible = true;
          }
        });
      }
    });
  }

  // 保存确认函数
  private onSave() {
    if (!this.project_name) {
      this.$message.error("");
      return;
    }
    const params = {
      type: this.form.type,
      transfer_no: this.form.transfer_no
    };
    Api.saveSetting(params).then(res => {
      if (res.data) {
        this.dialogVisible = false;
        this.$message.success("保存成功");
      }
    });
  }

  private mounted() {
    this.setUserList();
    this.getDetail();
  }
}
</script>
<style lang="scss" scoped>
.header-item {
  margin-bottom: 0;
}
.main-con /deep/ .search-content {
  background-color: #fff;
  padding: 0;
}
.main-con /deep/ .buttons-content {
  display: none;
}
.main-con /deep/ .tabs-content {
  display: none;
}
.main-con /deep/ .operate-content {
  display: none;
}
.main-con /deep/ .slide-box {
  display: none;
}
.main-con /deep/ .table-content {
  display: none;
}
/deep/.el-form-item__error {
  position: absolute;
  top: 48% !important;
  left: 10% !important;
}
/deep/ .el-form-item {
  margin-bottom: 20px;
}
/deep/ .el-input {
  width: 300px;
}
/deep/ .el-form-item__content {
  line-height: 18px;
}
.tips {
  margin-top: 5px;
}
.line {
  width: 100%;
  height: 1px;
  background-color: #dcdfe6;
}
.content {
  display: flex;
  align-items: center;
}
.tips {
  font-size: 12px;
  font-family: PingFangSC-Regular;
  font-weight: 400;
  color: #999;
  display: block;
}
</style>
