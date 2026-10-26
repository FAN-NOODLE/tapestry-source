export default {
  /**
   * @name 保留两位小数的千分位过滤器
   * @description 后端接口获取的金钱都以分为单位，处理前必为整数
   * */

  key: "formatMoneyByDot",
  func: function(money: number | string): string | number {
    const _money: number | string = Number(money);
    if (!_money) return 0;
    return `${_money.toString().replace(/(\d)(?=(?:\d{3})+$)/g, "$1,")}`;
  }
};
