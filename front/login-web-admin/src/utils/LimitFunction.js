export const _limitFrequence = (fn, num = 1) => {
  let result;
  let time = 0;
  return function() {
    if (time < num) {
      time++;
      result = fn.apply(this, arguments);
    } else {
      console.log("不要连续点击哦");
    }
    return result;
  };
};
