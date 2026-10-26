/**
 * 防抖
 * @param {Function} func -需要函数防抖的函数
 * @param {Number} delay -延迟时间
 **/
export class Debounce {
  public _debounce = (func: Function, delay = 500): Function => {
    let timeout: any;
    let arr: any;
    return () => {
      clearTimeout(timeout); // 如果持续触发，那么就清除定时器，定时器的回调就不会执行。
      timeout = setTimeout(() => {
        func.apply(this, arr);
      }, delay);
    };
  };
}

/**
 * @description: 限流
 * @param {Function} func 需要函数节流的函数
 * @param {Function} delay 多长时间不能再次点击
 */
export class Throttle {
  public _throttle = (func: Function, delay = 500): Function => {
    let last = 0;
    let arr: any;
    return () => {
      const now = Date.now();
      if (now - last >= delay) {
        // 当前时间 - 上次点击的时间 >= 间隔时间
        func.apply(this, arr);
        last = now;
      }
    };
  };
}
