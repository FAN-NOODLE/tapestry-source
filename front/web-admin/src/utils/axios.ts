"use strict";

import axios from "axios";
import LoginAboat from "@/utils/LoginManager";
import { Loading, Message } from "element-ui";
import { VUE_APP_API } from "@/utils/env-config";

let loadCount = 0; // 记录当前请求数，目的显示隐藏loading
// loading 实例
let loadingInstance: any;

// loading 配置文件
const options = {
  fullscreen: true,
  background: "rgba(0,0,0,0.1)",
};
axios.defaults.headers.post["Content-Type"] = "application/x-www-form-urlencoded;charset=UTF-8";
// api 配置
const config = {
  baseURL: VUE_APP_API, // gw-scrm
  timeout: 60000, // Timeout 为了防止导出过大文件失败，所以把时间设为1分钟
  headers: {},
  withCredentials: true, // Check cross-site Access-Control
};

const _axios = axios.create(config);

_axios.interceptors.request.use(
  function(_config) {
    if (_config.data && _config.data.apiVersion && _config.data.apiVersion === "old") {
      _config.baseURL = VUE_APP_API + "gw-console/gw-console/v1/";
      delete _config.data.apiVersion;
    } else if (_config.data && _config.data.apiVersion && _config.data.apiVersion === "brand") {
      _config.baseURL = VUE_APP_API + "gw-console/gw-console/brand/";
      delete _config.data.apiVersion;
    } else {
      _config.baseURL = VUE_APP_API + "gw-scrm/gw-scrm/";
    }
    if (_config.data && _config.data.apiVersion && _config.data.apiVersion === "yang") {
      _config.baseURL = "/api";
      _config.data.enterprise_id = 3306120225798016;
      delete _config.data.apiVersion;
    }
    if (loadCount++ === 0 && !(_config as any).hideLoading) {
      loadingInstance =  Loading.service(options);
    }
    // 如果是本地静默登录，不传hash和token
    if (LoginAboat.getRealHash()) {
      _config.headers["Enterprise-Hash"] = LoginAboat.getRealHash();
    }
    if (LoginAboat.getToken()) {
      _config.headers["Auth-Token"] = LoginAboat.getToken();
    }
    return _config;
  },
  function(error) {
    --loadCount <= 0 && loadingInstance && loadingInstance.close();
    return Promise.reject(error);
  }
);

// Add a response interceptor
_axios.interceptors.response.use(
  function(response) {
    --loadCount <= 0 && loadingInstance && loadingInstance.close();
    if ((response.status === 200 || response.status === 210) && response.data.error === "ERROR") {
      Message({
        showClose: true,
        type: "error",
        duration: 2000,
        message: response.data.message,
      });
      return Promise.reject(response.data);
    } else {
      return response.data;
    }
  },
  function(error) {
    --loadCount <= 0 && loadingInstance && loadingInstance.close();
    const { response } = error;
    loadingInstance && loadingInstance.close();
    if (
      parseInt(response.data.code) === 451004 ||
      parseInt(response.data.code) === 9999 ||
      parseInt(response.data.code) === 450000
    ) {
      const timeOut = setTimeout(() => {
        // LoginAboat.backLogin();
        localStorage.clear();
        window.location.href = `${window.location.origin}/login`;
        clearTimeout(timeOut);
      }, 1500);
      Message({
        showClose: true,
        type: "error",
        duration: 2000,
        message: "已过期, 请重新登录",
      });
    } else if (parseInt(response.data.code) === 700015) {
      const timeOut = setTimeout(() => {
        LoginAboat.backLogin();
        clearTimeout(timeOut);
      }, 1500);
      Message({
        showClose: true,
        type: "error",
        duration: 2000,
        message: "该管理员已被禁用",
      });
    } else {
      Message({
        showClose: true,
        type: "error",
        duration: 2000,
        message: response.data.message,
      });
      return Promise.reject(error.response.data);
    }
  }
);
export default _axios;
