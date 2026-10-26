// 格式化时间
export default {
  key: "formatDate",
  func: (time: number, type: string) => {
    if (time == 0) {
      return "--";
    }
    const date = new Date(time * 1000);
    const year = date.getFullYear();
    /* *
    在日期格式中，月份是从0开始的，因此要加0使用三元表达式在小于10的前面加0，以达到格式统一  如 09:11:05
    * */
    const month =
      date.getMonth() + 1 < 10
        ? `0${date.getMonth() + 1}`
        : date.getMonth() + 1;
    const day = date.getDate() < 10 ? `0${date.getDate()}` : date.getDate();
    const hours =
      date.getHours() < 10 ? `0${date.getHours()}` : date.getHours();
    const minutes =
      date.getMinutes() < 10 ? `0${date.getMinutes()}` : date.getMinutes();
    const seconds =
      date.getSeconds() < 10 ? `0${date.getSeconds()}` : date.getSeconds();
    let format_date = "";
    if (type.search("yyyy-mm-dd") !== -1) {
      if (type.search("hh:mm:ss") !== -1) {
        // 拼接
        format_date = `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
      } else {
        // 拼接
        format_date = `${year}-${month}-${day}`;
      }
    } else if (type.search("yyyy年mm月dd日") !== -1) {
      if (type.search("hh:mm:ss") !== -1) {
        // 拼接
        format_date = `${year}年${month}月${day}日 ${hours}:${minutes}:${seconds}`;
      } else {
        // 拼接
        format_date = `${year}年${month}月${day}日`;
      }
    } else if (type.search("yyyy/mm/dd") !== -1) {
      if (type.search("hh:mm:ss") !== -1) {
        // 拼接
        format_date = `${year}/${month}/${day} ${hours}:${minutes}:${seconds}`;
      } else {
        // 拼接
        format_date = `${year}/${month}/${day}`;
      }
    }
    return format_date;
  }
};
