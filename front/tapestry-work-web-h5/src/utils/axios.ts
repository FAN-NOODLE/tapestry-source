import axios, {
  AxiosError,
  AxiosInstance,
  AxiosRequestConfig,
  AxiosResponse,
} from "axios";
import Store from "@/store/index";
import router from "vue-router";
import { loginFnc } from "@/utils/login";
import helper from "@/utils/helper";
import { VUE_APP_API_ROOT, VUE_APP_API_VERSION } from "@/utils/env-config";
// todo: 增加接口异常上报
import { report  } from "./robot";
//unknown类型的值只能赋值给any或者unknown
const store: any = Store;
// 请求头
axios.defaults.headers.post["Content-Type"] =
  "application/x-www-form-urlencodedcharset=UTF-8";
axios.defaults.headers["Type-Origin"] = "H5";
const VERSION: string | undefined = VUE_APP_API_VERSION || "";
const APIHEAD: string | undefined = VUE_APP_API_ROOT || "";

// const IS_LOCAL = /^(172\.|127\.|localhost)/.test(location.hostname);
const IS_LOCAL = false;
// api 配置
//  process.env.NODE_ENV == "test"? "http://localhost:20000/gw-work-scrm"
// IS_LOCAL ? "/api/" + VERSION :
const config: AxiosRequestConfig = {
  baseURL:  APIHEAD + VERSION,
  timeout: 10000, // 超时时间
  withCredentials: true, // 表示跨域请求时是否需要使用凭证
};
const _axios: AxiosInstance = axios.create(config);

// 请求数据拦截
_axios.interceptors.request.use(
  (config) => {
    config.headers["Enterprise-Hash"] = store.state.common.ENTERPRISE_HASH;
    if (store.state.user.TOKEN) {
      config.headers["Authorization"] = "Bearer " + store.state.user.TOKEN;
    }

    if(Store.getters['lang/getLang'] === 'en-US') {
      config.headers["Accept-Language"] = 'en';
    } else {
      config.headers["Accept-Language"] = 'zh';
    }

    if (config.data && config.data.apiVersion == "new") {
      const oldURL = APIHEAD + VERSION;
      config.baseURL = oldURL.replace(
        "gw-work-scrm/gw-work-scrm/v1",
        "rest-api"
      );
      delete config.data.apiVersion;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);
// 相应数据拦截
_axios.interceptors.response.use(
  (response: AxiosResponse) => {
    // todo: 增加接口异常上报
    if(response.data.code !== 0 && response.data.code !== 200){
      // report({
      //   title: '企微侧边栏h5请求异常',
      //   url: `${response.config.baseURL}/${response.config.url}`,
      //   time: new Date().getTime(),
      //   header: response.config.headers,
      //   data: response.config.data,
      //   errMsg: response.data,
      // })
    }
    if (response.status === 200 && response.data.error === "ERROR") {
      return Promise.reject(response.data);
    } else if (
      response.data.code <= -1 &&
      response.data.code >= -5 &&
      IS_LOCAL
    ) {
      sessionStorage.setItem("pageUrl", window.location.href);
      // token 如果失效了跳转登录
      console.log(process.env.NODE_ENV);
      window.location.href = window.location.origin + "/login";
      return response.data;
    } else if (response.data.code == -2) {
      store.state.common.IS_LOGINED &&
        loginFnc({ code: helper.getUrlParam("code") }, router.prototype);
      store.state.common.IS_LOGINED = false;
      return response.data;
    } else if (response.data.code == 40029) {
      // code失效了,需要删除到url带的code和state参数,然后重新去获取code
      const url = helper.urlDel("code");
      store.state.common.IS_LOGINED = false;
      window.document.location.replace(url);
      return response.data;
    } else {
      return response.data;
    }
  },
  (error: AxiosError) => {
    if (
      error.response &&
      error.response.status === 500 &&
      error.response &&
      error.response.data.error === "ERROR"
    ) {
      if (parseInt(error.response.data.code) === 200004) {
        setTimeout(() => {
          window.location.href = window.location.origin + "/login";
        }, 1500);
      }
    }
    return Promise.reject(error.response && error.response.data);
  }
);

export default _axios;
