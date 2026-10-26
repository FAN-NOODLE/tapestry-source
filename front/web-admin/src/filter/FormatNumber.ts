// 格式化数字三位加一逗号
export default {
  key: "formatNumber",
  func: (value: number, length = 2): string => {
    let intPart = String(value);
    let floatPart = "";
    //判断是否是小数
    if (intPart.includes(".")) {
      console.log(value, 222, length);
      const arr = value.toFixed(length).split(".");
      intPart = arr[0];
      floatPart = "." + arr[1];
    }
    intPart = intPart.replace(/(\d)(?=(?:\d{3})+$)/g, "$1,");
    return intPart + floatPart;
  }
};
