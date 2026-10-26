/**
 * @description: 日期格式化
 * @param {time} 可以是时间戳也可以是日期类型
 * @return {formatStr} 时间格式yyyy-MM-dd hh:mm:ss
 */
export function dateFormat(time: any, formatStr: string) {
  let date: any;
  if (typeof time == "number" && time) {
    date = new Date(time * 1000);
  } else {
    if (!(time instanceof Date)) {
      throw new TypeError(`time:${time}不是一个有效日期类型`);
      // return time;
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
    S: date.getMilliseconds() // millisecond (毫秒)
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
  const str: string = dateFormat((Date.now() - 86400000) / 1000, "yyyy-MM-dd");
  return new Date(str + " 23:59:59").getTime();
}
export function startAndEndOfday(day: number): any {
  const str: string = dateFormat(day, "yyyy-MM-dd");
  return {
    start: new Date(str + " 00:00:00"),
    end: new Date(str + " 23:59:59")
  };
}
export function sendTimeDate(date: any) {
  if (date) {
    const nowDate =
      new Date(date).getFullYear() +
      "-" +
      (new Date(date).getMonth() + 1 < 10
        ? "0" + (new Date(date).getMonth() + 1)
        : new Date(date).getMonth() + 1) +
      "-" +
      (new Date(date).getDate() < 10
        ? "0" + new Date(date).getDate()
        : new Date(date).getDate());
    const nowTime =
      (new Date(date).getHours() < 10
        ? "0" + new Date(date).getHours()
        : new Date(date).getHours()) +
      ":" +
      (new Date(date).getMinutes() < 10
        ? "0" + new Date(date).getMinutes()
        : new Date(date).getMinutes()) +
      ":" +
      (new Date(date).getSeconds() < 10
        ? "0" + new Date(date).getSeconds()
        : new Date(date).getSeconds());
    return nowDate + " " + nowTime;
  } else {
    return "";
  }
}
