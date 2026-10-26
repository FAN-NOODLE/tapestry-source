import { Toast } from "vant";
/**
 * @description: 日期格式化
 * @param {time} 可以是时间戳也可以是日期类型
 * @return {formatStr} 时间格式yyyy-MM-dd hh:mm:ss
 */
function dateFormat(time: number, formatStr: string): string {
  let date: any;
  if (typeof time == "number" && time) {
    date = new Date(time * 1000);
  } else {
    if (!(date instanceof Date)) {
      return "-";
    }
  }

  const o: object = {
    // 键值对形式的数组。只能通过加强的for循环来迭代取值
    "M+": date.getMonth() + 1, // month
    "d+": date.getDate(), // day
    "h+": date.getHours(), // hour
    "m+": date.getMinutes(), // minute
    "s+": date.getSeconds(), // second
    "q+": Math.floor((date.getMonth() + 3) / 3), // quarter (季度)
    S: date.getMilliseconds(), // millisecond (毫秒)
  };

  if (/(y+)/.test(formatStr)) {
    formatStr = formatStr.replace(
      RegExp.$1,
      (date.getFullYear() + "") //  RegExp.$1 : 取正则表达式中第一个分组匹配到的内容
        .substr(4 - RegExp.$1.length)
    );
  }

  for (const k in o) {
    // 加强的for循环。k为键值对数组o中的键，故o[k]为对应的值
    if (new RegExp("(" + k + ")").test(formatStr)) {
      formatStr = formatStr.replace(
        RegExp.$1,
        RegExp.$1.length == 1 ? o[k] : ("00" + o[k]).substr(String(o[k]).length)
      );
    }
  }
  return formatStr;
}
// 将时间段格式化成 时:分:秒
function formatTime(time: number): string {
  const h = Math.floor(time / 3600);
  const m = Math.floor((time % 3600) / 60);
  const s = Math.floor(time % 60);
  return `${h >= 10 ?h:`0${h}`}:${m >= 10 ?m:`0${m}`}:${s >= 10 ?s:`0${s}`}`;
}
function limitMinDate(): number {
  const str: string = dateFormat((Date.now() - 86400000) / 1000, "yyyy-MM-dd");
  return new Date(str + " 23:59:59").getTime();
}
function startAndEndOfday(day: number): any {
  const str: string = dateFormat(day, "yyyy-MM-dd");
  return {
    start: new Date(str + " 00:00:00"),
    end: new Date(str + " 23:59:59"),
  };
}
/**
 * @description:数字千分化
 * @param {number} num
 * @return {*}
 */
function NumberThousand(num: number, isMoney = false): string {
  let format: string | number = String(num).replace(
    /(\d{1,3})(?=(\d{3})+(?:$|\.))/g,
    "$1,"
  );
  format = isMoney ? `¥${num}` : format;
  return format;
}
/**
 * @description: 防抖函数 一段时间内只执行一次 控制优先执行还是滞后执行
 * @return {*}
 */
function debounce(fn: Function, wait = 500, immediate = false) {
  if (typeof fn !== "function") throw "argument[0] must be the function";
  let timeout;
  const debounced = function () {
    const context = this,
      args = arguments;
    if (timeout) clearTimeout(timeout);
    if (immediate) {
      const callNow = !timeout;
      timeout = setTimeout(function () {
        timeout = null;
      }, wait);
      if (callNow) fn.apply(context, args);
    } else {
      timeout = setTimeout(function () {
        fn.apply(context, args);
      }, wait);
    }
  };
  debounced.cancel = function () {
    clearTimeout(timeout);
    timeout = null;
  };
  return debounced;
}

/**
 * @description: 节流函数
 * @param {Function} fn
 * @param {number} threshhold
 * @return {*}
 */
function throttle(fn, wait = 1200, immediate) {
  if (typeof fn !== 'function') throw "argument[0] must be the function";
  let last = 0, timeout;
  const throttle = function () {
      const now = Date.now();
      if (last && now < last + wait) {
          if (!immediate) {
              clearTimeout(timeout);
              let context = this, args = arguments,
              remain=last + wait - now;
              if(remain<0) remain=0;
              timeout = setTimeout(function () {
                  last = now;
                  timeout = null;
                  fn.apply(context, args);
              }, remain);//执行最近一次函数
          }
          return;
      }
      if (!immediate && timeout) return;
      last = now;
      fn.apply(this, arguments);
  }
  throttle.cancel = function () {
      clearTimeout(timeout);
      timeout = null;
      last = 0;
  }
  return throttle;
};
function getUrlParam(key: string) {
  let value = "";
  try {
    value = decodeURIComponent(
      (new RegExp("[?|&]" + key + "=" + "([^&;]+?)(&|#|;|$)").exec(
        location.href
      ) || "")[1].replace(/\+/g, "%20")
    );
  } catch (error) {}
  return value;
}
// 对url进行某个参数的删除，并返回url
function urlDel(name: string) {
  console.log(typeof window.location, "window.location");
  const loca = window.location;
  const baseUrl = loca.origin + loca.pathname + "?";
  const query = loca.search.substr(1);
  if (query.indexOf(name) > -1) {
    const obj = {};
    const arr: any = query.split("&");
    for (let i = 0; i < arr.length; i++) {
      arr[i] = arr[i].split("=");
      obj[arr[i][0]] = arr[i][1];
    }
    delete obj[name];
    if (name == "code") {
      delete obj["state"];
    }
    const url =
      baseUrl +
      JSON.stringify(obj)
        .replace(/[\"\{\}]/g, "")
        .replace(/\:/g, "=")
        .replace(/\,/g, "&");
    return url;
  } else {
    return window.location.href;
  }
}
function copyMessage(msg: string): void {
  try {
    const element = document.createElement("textarea");
    element.value = msg;
    if (
      "append" in document.body &&
      typeof document.body.append == "function"
    ) {
      document.body.append(element);
    } else {
      (document.body as any).appendChild(element);
    }
    element.select();
    document.execCommand("Copy");
    element.remove();
  } catch (error) {
    console.log(error);
    Toast({
      duration: 1000,
      message: "复制失败",
      className: "toast-one-row",
    });
    return;
  }

  Toast({
    duration: 1000,
    message: "复制成功",
    className: "toast-one-row",
  });
}
// 是否KS企业
function isKS() {
  // 企业hash：KS-uat 522474ecff3b03fb41d8360eb259da58
  // 企业hash：KS-prod 0db7b1ee65bb6b547e687c9ae54f4f56
  let hashArrs = [
    "522474ecff3b03fb41d8360eb259da58",
    "0db7b1ee65bb6b547e687c9ae54f4f56",
  ];
  return hashArrs.some((ele) => {
    return location.href.indexOf(ele) > -1;
  });
  // return true;
}
// 是否SW企业
function isSW() {
  // 企业hash：SW-uat 9c217b24b3fc5dcf685a7c566d191272
  // 企业hash：SW-prod 7cdd2af2e96d8617437846067440a3ec
  let hashArrs = [
    "9c217b24b3fc5dcf685a7c566d191272",
    "7cdd2af2e96d8617437846067440a3ec",
  ];
  return hashArrs.some((ele) => {
    return location.href.indexOf(ele) > -1;
  });
  // return true;
}
// 是否coach企业
function isCoach() {
  // 企业hash：coach-uat 9c217b24b3fc5dcf685a7c566d191272
  // 企业hash：coach-prod 7cdd2af2e96d8617437846067440a3ec
  let hashArrs = [
    "457a684cb1477448d1798ef36a510102",
    "f7c9d26cf6441e9dbad633bab7be9af6",
  ];
  return hashArrs.some((ele) => {
    return location.href.indexOf(ele) > -1;
  });
  // return true;
}
// 获取小程序appid
function getAppid() {
  let appid = "";
  // coach uat
  if (location.href.indexOf("457a684cb1477448d1798ef36a510102") > -1) {
    appid = "wxbd40b04da240ef2b";
    // coach prod
  } else if (location.href.indexOf("f7c9d26cf6441e9dbad633bab7be9af6") > -1) {
    appid = "wx19ba52f691c4d289";
    // sw uat
  } else if (location.href.indexOf("9c217b24b3fc5dcf685a7c566d191272") > -1) {
    appid = "wxab53091260791a35";
    // sw prod
  } else if (location.href.indexOf("7cdd2af2e96d8617437846067440a3ec") > -1) {
    appid = "wx5abbf56188275030";
    // ks uat
  } else if (location.href.indexOf("522474ecff3b03fb41d8360eb259da58") > -1) {
    appid = "wxbc38547936d58845";
    // ks prod
  } else if (location.href.indexOf("0db7b1ee65bb6b547e687c9ae54f4f56") > -1) {
    appid = "wxce2ea306aa80fe90";
  }
  return appid;
}
const Helper: helper = {
  formatTime: formatTime,
  dateFormat: dateFormat,
  limitMinDate: limitMinDate,
  startAndEndOfday: startAndEndOfday,
  NumberThousand: NumberThousand,
  debounce: debounce,
  throttle: throttle,
  getUrlParam: getUrlParam,
  urlDel: urlDel,
  copyMessage: copyMessage,
  getAppid: getAppid,
  isKS: isKS,
  isSW: isSW,
  isCoach: isCoach,
};
export default Helper;
