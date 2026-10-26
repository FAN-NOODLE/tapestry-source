<template>
  <div class="content">
    <layout-form :title_name="title_name">
      <template slot="content">
        <cy-form
          class="demo-ruleForm"
          ref="base"
          :model="formData"
          label-width="104px"
          size="medium"
          :rules="base_rules"
        >
          <cy-form-item
            label="原始密码"
            prop="originalPass"
            required
          >
            <cy-input
              type="password"
              show-password
              autocomplete="new-password"
              widthSize="medium"
              v-model.trim="formData.originalPass"
              placeholder="请输入账号密码"
            ></cy-input>
          </cy-form-item>
          <cy-form-item
            label="新密码"
            prop="newPass"
            required
          >
            <cy-input
              type="password"
              show-password
              autocomplete="new-password"
              widthSize="medium"
              v-model.trim="formData.newPass"
              placeholder="请输入账号密码"
            ></cy-input>
          </cy-form-item>
          <cy-form-item
            label="确认密码"
            prop="againPass"
            required
          >
            <cy-input
              type="password"
              show-password
              autocomplete="new-submitPass"
              widthSize="medium"
              v-model.trim="formData.againPass"
              placeholder="请输入账号密码"
            ></cy-input>
          </cy-form-item>
        </cy-form>
      </template>
      <template slot="footer">
        <div>
          <el-button
            type="primary"
            size="small"
            @click="submitForm()"
          >
            保存
          </el-button>
          <el-button
            size="small"
            @click="goBack()"
          >取消</el-button>
        </div>
      </template>
    </layout-form>
  </div>
</template>
<script>
import Api from "@/api/setting/Account";
import LoginAboat from "@/utils/LoginManager";
export default {
  name: "ChangePassword",
  data() {
    const validatePass = (rule, value, callback) => {
      const reg = /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])[a-zA-Z0-9]{8,}$/;
      if (!value) {
        return callback(new Error("请输入登录密码"));
      }
      if (!reg.test(value)) {
        return callback(
          new Error("密码至少包含8个字符，一个大写字母，一个小写字母和一个数字")
        );
      }
      callback();
    };
    const validateAgainPass = (rule, value, callback) => {
      const reg = /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])[a-zA-Z0-9]{8,}$/;
      if (!value) {
        return callback(new Error("请再次输入密码"));
      }
      if (value != this.formData.newPass) {
        return callback(new Error("请与新密码保持一直"));
      }
      if (!reg.test(value)) {
        return callback(
          new Error("密码至少包含8个字符，一个大写字母，一个小写字母和一个数字")
        );
      }
      callback();
    };
    return {
      base_rules: {
        originalPass: [
          {
            required: true,
            trigger: ["blur", "change"],
            message: "请输入原始密码"
          }
        ],
        newPass: [
          {
            validator: validatePass,
            trigger: ["blur", "change"]
          }
        ],
        againPass: [
          {
            validator: validateAgainPass,
            trigger: ["blur", "change"]
          }
        ]
      },
      title_name: "修改密码",
      formData: {
        originalPass: "",
        newPass: "",
        againPass: ""
      }
    };
  },
  methods: {
    // 提交
    submitForm() {
      this.$refs["base"].validate(valid => {
        if (valid) {
          Api.editPassword({
            old_ps: this.$md5(this.formData.originalPass + ""), //原始密码
            new_ps: this.$md5(this.formData.newPass + ""), //新密码
            new_ps_two: this.$md5(this.formData.againPass + "") //确认密码
          }).then(res => {
            if (res.code == 200) {
              // 成功后退出
              this.$message({
                showClose: true,
                message: "修改成功，请重新登录",
                type: "success"
              });
              LoginAboat.backLogin();
            }
          });
        } else {
          console.log("输入有误请完善 :>> ");
        }
      });
    },
    //返回
    goBack() {
      this.$router.go(-1);
    }
  }
};
</script>
<style lang="scss" scoped></style>
