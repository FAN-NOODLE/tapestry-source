export default {
  /**
   * @name 保留两位小数的千分位过滤器
   * @description 后端接口获取的金钱都以分为单位，处理前必为整数
   * */

  key: "FormatYuanMoney",
  func: function(money: number | string): string | number {
    let _money: number | string = Number(money);
    if (!_money) return 0;
    _money = _money.toFixed(2);
    const intPart = _money.substr(0, _money.length - 3);
    const floatPart = _money.substr(_money.length - 3);
    return `¥${intPart
      .toString()
      .replace(/(\d)(?=(?:\d{3})+$)/g, "$1,")}${floatPart}`;
  }
};
