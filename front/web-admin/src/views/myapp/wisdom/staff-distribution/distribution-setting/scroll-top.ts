export function scrollTop(
  el: any,
  from: number,
  to: number,
  duration = 500,
  endCallback?: Function
) {
  if (!window.requestAnimationFrame) {
    window.requestAnimationFrame =
      (window as any).webkitRequestAnimationFrame ||
      window["mozRequestAnimationFrame"] ||
      window["msRequestAnimationFrame"] ||
      function(callback) {
        return window.setTimeout(callback, 1000 / 60);
      };
  }
  const difference = Math.abs(from - to);
  const step = Math.ceil((difference / duration) * 50);

  function scroll(start: number, end: number, steps: number) {
    if (start === end) {
      endCallback && endCallback();
      return;
    }

    let d = start + steps > end ? end : start + steps;
    if (start > end) {
      d = start - steps < end ? end : start - steps;
    }

    if (el === window) {
      window.scrollTo(d, d);
    } else {
      el.scrollTop = d;
    }
    window.requestAnimationFrame(() => scroll(d, end, steps));
  }
  scroll(from, to, step);
}
