export default {
  data() {
    const validateUsername = (rule, value, callback) => {
      const reg = /^1[23456789]\d{9}$/;
      if (reg.test(value)) {
        return callback(new Error("请不要用手机号作为用户名"));
      }
      callback();
    };
    return {
      SESSION_NAME: "accounts-edit-form",
      base_rules: {
        username: [
          {
            validator: validateUsername,
            trigger: ["blur", "change"]
          }
        ],
        mobile: [
          {
            required: true,
            trigger: ["blur", "change"],
            message: "请输入账号手机号"
          }
        ],
        email: [
          {
            required: true,
            trigger: ["blur", "change"],
            message: "请输入账号邮箱"
          }
        ]
      },
      inner_no_choice: ["guide_enterprise_guide", "guide_enterprise_admin"],
      inner_no_choice_admin: [
        "guide_enterprise_guide",
        "guide_enterprise_admin"
      ]
    };
  },
  methods: {
    // 返回
    goBack() {
      if (this.isChange() && !this.$route.query.id) {
        this.$confirm(
          "检测到未保存的内容，是否保存至下次新建继续编辑？",
          "离开页面提示",
          {
            distinguishCancelAndClose: true,
            confirmButtonText: "保存并离开",
            confirmButtonClass: "button-danger",
            cancelButtonText: "直接离开",
            type: ""
          }
        )
          .then(() => {
            this.catchBack(false);
          })
          .catch(action => {
            if (action === "cancel") {
              // 清除缓存
              this.catchBack(true);
            }
          });
      } else {
        this.$router.go(-1);
      }
    },
    // 返回并清除缓存
    catchBack(flag) {
      if (flag) {
        sessionStorage.removeItem(this.SESSION_NAME);
      } else {
        sessionStorage.setItem(
          this.SESSION_NAME,
          JSON.stringify(this.formData)
        );
      }
      this.$router.go(-1);
    },
    isChange() {
      console.log("当前======" + JSON.stringify(this.formData));
      //是否修改
      if (
        JSON.stringify(this.formData) ===
        sessionStorage.getItem(this.SESSION_NAME)
      ) {
        return false;
      } else {
        return true;
      }
    },
    setAndGetSesssion() {
      if (sessionStorage.getItem(this.SESSION_NAME)) {
        this.formData = JSON.parse(sessionStorage.getItem(this.SESSION_NAME));
      } else {
        sessionStorage.setItem(
          this.SESSION_NAME,
          JSON.stringify(this.formData)
        );
      }
    }
  }
};
