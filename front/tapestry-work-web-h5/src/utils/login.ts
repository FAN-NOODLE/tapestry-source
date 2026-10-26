import { loginApi, redirectUrlApi } from "@/api/common";
import Store from "@/store/index";
import VueRouter from "vue-router";
import { Toast } from "vant";
import { sensorsInit, sensorsLogin, setSensorsProfile } from "./sensors";
import Vue from "vue";
import { getStoreVisitAuthApi } from "@/api/store-visit";
import { getManagerInfoApi } from "@/api/communication";
import { getAuthApi } from "@/api/display-task";
const store: any = Store;

async function getStoreVisitAuth() {
  try {
    const {code, message, data}:IApiResult = await getStoreVisitAuthApi({})
    if (code!== 200) {
      return;
    }
    Store.commit("user/setUserInfo", {
      ...Store.getters["user/getUserInfo"], 
      ...data,
    })
  } catch (error) {
    console.log(error);
  }
}

/**
 * @description 获取用户权限
 */
async function getTransmitalAuth() {
  try {
    const {code, message, data}:IApiResult = await getManagerInfoApi({})
    if (code!== 200) {
      return;
    }
    
    Store.commit("user/setUserInfo", {
      app_source: 2,
      ...Store.getters["user/getUserInfo"], 
      ...data,
    })
  } catch (error) {
    console.log(error);
  }
}


/**
 * @description 获取用户权限
 */
async function getVMAuth() {
  try {
    const {code, message, data}:IApiResult = await getAuthApi()
    if (code!== 200) {
      this.$toast.fail(message);
      return;
    }
    Store.commit("user/setUserInfo", {
      app_source: 1,
      ...Store.getters["user/getUserInfo"], 
      ...data,
    })
  } catch (error) {
    console.log(error);
  }
}


async function loginFnc(query: any, router: VueRouter) {
  const userCode = localStorage.getItem("userCode");
  // 判断有无 code
  if (query.hasOwnProperty("code") && !!query.code && userCode !== query.code) {
    // 微信重定向后，用code换取token
    console.log(query, router, "我是Code");
    console.log(query, "我是Code");
    console.log(window.location);

    localStorage.setItem("userCode", query.code);
    return codeloginFnc(query, router);
  } else {
    // 第一次登陆，微信重定向后，链接上拼接了code
    return redirectloginFnc(router);
  }
}

async function redirectloginFnc(router: VueRouter) {
  try {
    const uri = window.location.href.replace(
      /(code=(\w+)&?)|(state=(\w+)&?)/g,
      ""
    );
    const {
      code,
      message = "",
      data,
    }: IApiResult = await redirectUrlApi({
      redirect_uri: uri,
      hash: store.state.common.ENTERPRISE_HASH,
    });
    console.log("接口code", code);
    // 登录成功
    if (code === 200) {
      // 授权
      if (data.url) {
        if (
          process.env.NODE_ENV !== "test" &&
          process.env.NODE_ENV !== "development"
        ) {
          console.log("url", data.url);
          // 删除路由记录
          // history.replaceState(null, "", data.url);
          // window.document.location.replace(data.url);
          window.location.href = data.url;
          // 新窗口打开
          // window.open(data.url);
        } else {
          console.log(
            "开发者记得在src/store/modules/user.ts文件里写个有效的token来模拟登陆"
          );
          console.log(
            "模拟token可在浏览器访问该链接获取(记得根据实际开发环境修改链接的企业id、员工user_id访问)：https://crs-api-test.vchangyi.com/gw-work-scrm/v1/auth/mock?userid=luolan&enterprise_id=114"
          );
        }
      } else {
        console.log("redirectloginFnc登录成功");
        // 用户信息存入全局
        Store.commit("user/setUserInfo", data.member);
        Store.commit("user/setIsLogined", true);
        // 获取门店访问权限
        const isStoreVisit = location.href.indexOf('store-visit') !== -1
        if(isStoreVisit) {
          await getStoreVisitAuth();
        }

        // 传讯
        const isTransmital = location.href.indexOf('communication') !== -1
        if(isTransmital) {
          await getTransmitalAuth();
        }

        // vm
        const isVM = location.href.indexOf('display-task') !== -1
        if(isVM) {
          await getVMAuth();
        }


        // // 神策数据上报 -- 初始化
        // const sensors = sensorsInit()
        // // 神策数据上报 -- 初始化并挂载全局
        // Vue.prototype.sensors = sensors
        // 神策数据上报 -- 登录
        sensorsLogin(data.member.user_id)
        // 神策数据上报 -- 设置用户属性
        setSensorsProfile({
          store_name: Store.getters["user/getUserInfo"].store_name,
          store_code: Store.getters["user/getUserInfo"].store_code || Store.getters["user/getUserInfo"].store_name
        })
      }
    } else {
      //  todo: 未授权
      Toast("登录失败,请刷新页面重试~");
    }
  } catch (err) {
    console.log("redirectloginFnc", err);
    //  todo: 失败
    Toast("你没网络了哦，请检查下网络设置~");
  }
  return;
}

async function codeloginFnc(query: any, router: VueRouter) {
  try {
    const { code, data }: IApiResult = await loginApi({
      code: query.code,
      hash: store.state.common.ENTERPRISE_HASH,
    });
    if (code === 200) {
      console.log("codeloginFnc登录成功");
      // 用户信息存入全局
      Store.commit("user/setUserInfo", data.member);
      Store.commit("user/setToken", data.token);
      Store.commit("common/setIsLogined", true);
      Store.commit("user/setEnterpriseMark", data.brand_code);

      // 获取门店访问权限
      const isStoreVisit = location.href.indexOf('store-visit') !== -1
      if(isStoreVisit) {
        await getStoreVisitAuth();
      }

      const isTransmital = location.href.indexOf('communication') !== -1
      if(isTransmital) {
        await getTransmitalAuth();
      }

      // vm
      const isVM = location.href.indexOf('display-task') !== -1
      if(isVM) {
        await getVMAuth();
      }

      // // 神策数据上报 -- 初始化
      // const sensors = sensorsInit()
      // // 神策数据上报 -- 初始化并挂载全局
      // Vue.prototype.sensors = sensors
      // 神策数据上报 -- 登录
      sensorsLogin(data.member.user_id)
      // 神策数据上报 -- 设置用户属性
      setSensorsProfile({
        store_name: Store.getters["user/getUserInfo"].store_name,
        store_code: Store.getters["user/getUserInfo"].store_code || Store.getters["user/getUserInfo"].store_name
      })
    }
  } catch (err) {
    console.log("codeloginFnc", err);
  }
  return;
}

export { loginFnc, redirectloginFnc };
