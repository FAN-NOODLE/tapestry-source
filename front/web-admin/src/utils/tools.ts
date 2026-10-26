export const cloneDeep = (obj: any) => {
  if (obj === null || typeof obj !== "object") return obj;
  const newObj = obj instanceof Array ? [] : {};
  for (const key in obj) {
    if (Object.prototype.hasOwnProperty.call(obj, key)) {
      newObj[key] = typeof obj[key] === "object" ? cloneDeep(obj[key]) : obj[key];
    }
  }
  return newObj;
};
// 判断空对象、空数组、空字符串、0
export const isNull = (data: any) => {
  if (
    typeof data == "undefined" ||
    data == null ||
    data == "" ||
    (typeof data === "object" && data instanceof Array && data.length == 0)
  ) {
    return true;
  } else {
    if (typeof data === "object" && !(data instanceof Array)) {
      if (Object.keys(data).length == 0) {
        return true;
      } else {
        return false;
      }
    }
  }
};
// 将字符串空格回车转转义成html可识别的内容
export const trimEnter: (text: string) => string = function(text: string) {
  return `${text.replace(/\n|\r\n/g, "<br/>").replace(/\s/g, "&nbsp;&nbsp;")}`;
};
export const publicMethod = (time: any) => {
  //把时间日期转成时间戳
  return Math.floor(new Date(time).getTime() / 1000);
};
/**
 * @description: 日期格式化
 * @param {time} 可以是时间戳也可以是日期类型
 * @return {formatStr} 时间格式yyyy-MM-dd hh:mm:ss
 */
export function dateFormat(time: any, formatStr: string): string | null {
  let date: any;
  if (typeof time == "number" && time) {
    date = new Date(time * 1000);
  } else {
    if (!(time instanceof Date)) {
      throw new TypeError(`time:${time}不是一个有效日期类型`);
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
export function limitMinDate(): number {
  const str: string | null = dateFormat((Date.now() - 86400000) / 1000, "yyyy-MM-dd");
  return new Date(str + " 23:59:59").getTime();
}
export function startAndEndOfday(day: number): any {
  const str: string | null = dateFormat(day, "yyyy-MM-dd");
  return {
    start: new Date(str + " 00:00:00"),
    end: new Date(str + " 23:59:59"),
  };
}

export function getYesterDays() {
  // 起止日期数组
  const current = new Date();
  // 昨天开始时间
  const millisecond = 1000 * 60 * 60 * 24;
  // 近七天的开始时间
  const yesterday = new Date(current.getTime() - 1 * millisecond);
  // debugger;
  const param = [yesterday.setHours(0, 0, 0, 0), yesterday.setHours(23, 59, 59, 999)];
  return param;
}

export function getSevenDays() {
  // 起止日期数组
  const current = new Date();
  // 昨天开始时间
  const millisecond = 1000 * 60 * 60 * 24;
  // 七天前的开始时间
  const sevenDaysStart = new Date(current.getTime() - 6 * millisecond);
  const param = [sevenDaysStart.setHours(0, 0, 0, 0), current.getTime()];
  return param;
}

export function getThirtyDays() {
  // 起止日期数组
  // 获取当前时间
  const current = new Date();
  // 一天的毫秒数
  const millisecond = 1000 * 60 * 60 * 24;
  // 近30天的开始时间
  const thirtyDay = new Date(current.getTime() - 29 * millisecond);
  const param = [thirtyDay.setHours(0, 0, 0, 0), current.getTime()];
  return param;
}

// 一段时间的时间戳-例昨天、七天内、30天时间戳
export function timeForMat(count: number) {
  switch (count) {
    case 1:
      const yesArr = getYesterDays();
      console.log("yesArr", yesArr);
      return {
        t1: yesArr[1],
        t2: yesArr[0],
      };
    case 7:
      const sevenArr = getSevenDays();
      console.log("sevenArr", sevenArr);
      return {
        t1: sevenArr[1],
        t2: sevenArr[0],
      };
    case 30:
      const thirtyArr = getThirtyDays();
      console.log("thirtyArr", thirtyArr);
      return {
        t1: thirtyArr[1],
        t2: thirtyArr[0],
      };
    default:
      // 拼接时间
      const time1 = new Date();
      const Y1 = time1.getFullYear();
      const M1 = time1.getMonth() + 1 > 10 ? time1.getMonth() + 1 : "0" + (time1.getMonth() + 1);
      const D1 = time1.getDate() > 10 ? time1.getDate() : "0" + time1.getDate();
      const timer1 = Y1 + "-" + M1 + "-" + D1; // 当前时间
      const time2 = new Date();
      time2.setTime(time2.getTime() - 24 * 60 * 60 * 1000 * count);
      const Y2 = time2.getFullYear();
      const M2 = time2.getMonth() + 1 > 10 ? time2.getMonth() + 1 : "0" + (time2.getMonth() + 1);
      const D2 = time2.getDate() > 10 ? time2.getDate() : "0" + time2.getDate();
      const timer2 = Y2 + "-" + M2 + "-" + D2; // 之前的7天或者30天
      return {
        t1: timer1 + " 23:59:59",
        t2: timer2 + " 00:00:00",
      };
  }
}
export const downFilePost = (fuc: any, params: any, name: string, suffix = ".xlsx") => {
  fuc(params).then((res: any) => {
    const blob = new Blob([res], {
      type: "application/zip;charset=utf-8",
    });
    // let fileName = "";
    const fileName = name + new Date().toLocaleString().substr(0, 10) + suffix;

    //下载的文件名称及其后缀，后缀要和后台保持的一致

    if ("download" in document.createElement("a")) {
      // 非IE下载
      const elink = document.createElement("a");
      elink.download = fileName;
      elink.style.display = "none";
      elink.href = URL.createObjectURL(blob);
      document.body.appendChild(elink);
      elink.click();
      URL.revokeObjectURL(elink.href); // 释放URL 对象
      document.body.removeChild(elink);
    } else {
      // IE10+下载
      (navigator as any).msSaveBlob(blob, fileName);
    }
  });
};
