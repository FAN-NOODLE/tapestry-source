"use strict";

import Vue from "vue";
import axios from "axios";
import { Loading, Message } from "element-ui";

// loading 配置文件
const options = {
  fullscreen: true,
  background: "rgba(0,0,0,0.1)"
};
// loading 实例
let loadingInstance = null;

axios.defaults.headers.post["Content-Type"] =
  "application/x-www-form-urlencoded;charset=UTF-8";

// api 配置

let config = {
  // // baseURL: process.env.VUE_APP_API_ROOT + process.env.VUE_APP_API_VERSION,
  baseURL: "",
  timeout: 10000, // Timeout
  withCredentials: true // Check cross-site Access-Control
};

function _getAPIROOT(config) {
  const api =
    config && config.data && config.data.apiVersion === "old"
      ? "gw-console/"
      : "gw-scrm/";

  const APIS = [
    {
      name: "DEV",
      origin: "tapestry-dev.vchangyi.com",
      domain: "https://tapestry-dev.vchangyi.com"
    },
    {
      name: "UAT",
      origin: "tapestry-web-uat.coach.com.cn",
      domain: "https://tapestry-api-uat.coach.com.cn"
    },
    {
      name: "PROD",
      origin: "we.coach.com.cn",
      domain: "https://webk.coach.com.cn"
    },
    {
      name: "SW_UAT",
      origin: "tapestry-web-uat.stuartweitzman.cn",
      domain: "https://tapestry-api-uat.stuartweitzman.cn"
    },
    {
      name: "SW_PROD",
      origin: "we.stuartweitzman.cn",
      domain: "https://webk.stuartweitzman.cn"
    },
    {
      name: "KS_PROD",
      origin: "we.katespade.cn",
      domain: "https://webk.katespade.cn"
    },
    {
      name: "KS_UAT",
      origin: "tapestry-web-uat.katespade.cn",
      domain: "https://tapestry-api-uat.katespade.cn"
    }
  ];
  const API = APIS.find(item => {
    return location.href.indexOf(item.origin) > -1;
  });
  const Domain = API ? API.domain : APIS[0].domain;

  const VUE_APP_API_ROOT = `${Domain}/${api}`;

  return VUE_APP_API_ROOT;
}

const _axios = axios.create(config);

_axios.interceptors.request.use(
  function(config) {
    config.baseURL = _getAPIROOT(config);
    if (config && config.data && config.data.apiVersion)
      delete config.data.apiVersion;
    if (!config.loadingHide) {
      loadingInstance = Loading.service(options);
    }
    if (localStorage.getItem("TOKEN")) {
      config.headers["Auth-Token"] = localStorage.getItem("TOKEN");
    }
    return config;
  },
  function(error) {
    return Promise.reject(error);
  }
);

// Add a response interceptor
_axios.interceptors.response.use(
  function(response) {
    loadingInstance && loadingInstance.close();
    if (
      (response.status === 200 || response.status === 210) &&
      response.data.error === "ERROR"
    ) {
      Message({
        showClose: true,
        type: "error",
        duration: 2000,
        message: response.data.message
      });
      return Promise.reject(response.data);
    } else {
      return response.data;
    }
  },
  function(error) {
    const { response } = error;
    loadingInstance && loadingInstance.close();
    Message({
      showClose: true,
      type: "error",
      duration: 2000,
      message: response.data.message
    });
    if (parseInt(response.data.code) === 200004) {
      setTimeout(() => {
        LoginAboat.backLogin();
      }, 1500);
    }
    return Promise.reject(error.response.data);
  }
);

Plugin.install = function(Vue, options) {
  Vue.axios = _axios;
  window.axios = _axios;
  Object.defineProperties(Vue.prototype, {
    axios: {
      get() {
        return _axios;
      }
    },
    $axios: {
      get() {
        return _axios;
      }
    }
  });
};

Vue.use(Plugin);

export default Plugin;
