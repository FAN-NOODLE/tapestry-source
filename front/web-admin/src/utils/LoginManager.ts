import { IUserInfo } from "@/types";
import Api from "@/api/public/public";
const isProduction = process.env.NODE_ENV === "production";
/**
 * 登陆流程相关控制
 *
 */
// const cookiePath = {
//   "https://www.vchangyi.com": "/login",
//   "https://crs.vchangyi.com": "/",
//   "https://crs-test.vchangyi.com": "/",
//   "https://crs-dev.vchangyi.com": "/"
// };
class LoginManager {
  static MOCK_ACCOUNT = "yanwenqing@vchangyi.com"; //本地静默登录账号
  static PASS_WORD = "QxjrJ0"; //本地静默登录密码
  // 获取登陆数据
  static getLoginData: () => IUserInfo = function() {
    try {
      return JSON.parse(localStorage.getItem("USERINFO") || "");
    } catch (error) {
      return null;
    }
  };
  // 设置登陆数据
  static setLoginData: (data: object) => void = function(data: object) {
    if (data) {
      localStorage.setItem("USERINFO", JSON.stringify(data));
    }
  };
  // 获取hash
  static getAppHash: () => string = function() {
    try {
      return JSON.parse(localStorage.getItem("LOGIN_HASH") || "");
    } catch (error) {
      return null;
    }
  };
  // 获取真实hash
  static getRealHash: () => string = function() {
    try {
      return JSON.parse(localStorage.getItem("REAL_HASH") || "");
    } catch (error) {
      return null;
    }
  };
  // 退出登录
  static backLogin: () => void = function() {
    // await
    Api.logout({})
      .then(() => {
        localStorage.clear();
        const date: Date = new Date();
        date.setTime(date.getTime() - 10000);
        const keys = document.cookie.match(/[^ =;]+(?=)/g);
        if (keys) {
          for (let i = keys.length; i--; ) document.cookie = keys[i] + "=; expire=" + date.toUTCString() + "; path=/";
        }
        if(isProduction){
          window.location.href = `${window.location.origin}/login`;
        }else{
          window.location.href = `${window.location.origin}/login`;
        }
      })
      .catch((err) => {
        console.log(err);
      });
  };
  // 获取token
  static getToken: () => string | null = function() {
    try {
      return JSON.parse(localStorage.getItem("TOKEN") || "");
    } catch (error) {
      return null;
    }
  };
}

export default LoginManager;
