import Api from "@/api/public/Common";
// 手机号范围校验规则
export const checkMobile = (rule, value, callback) => {
  const reg = /^1[3456789]\d{9}$/;
  if (!reg.test(value)) {
    callback(new Error("手机号有误，请重新输入"));
  } else {
    callback();
  }
};
// 邮箱格式验证
export const checkEmail = (rule, value, callback) => {
  // eslint-disable-next-line no-useless-escape
  const reg = /^([A-Za-z0-9_\-\.])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,4})$/;
  if (!reg.test(value)) {
    callback(new Error("邮箱有误，请重新输入"));
  } else {
    // 验证邮箱是否已被使用
    Api.checkEmail({ email: value }).then(res => {
      // const { data } = res;
      if (res.status == true) {
        callback(new Error("该邮箱已注册，请重新输入或直接登录"));
      } else {
        callback();
      }
    });
  }
};
