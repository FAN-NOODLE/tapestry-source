/* eslint-disable @typescript-eslint/ban-types */
/**
 * @author Leeyunxiang
 * @description 兼容各浏览器添加事件的惰性函数
 * @param {String} type 事件类型(如click)
 * @param {Element} element 事件触发者
 * @param {Function} callback 事件的触发回调函数
 * */

export let addEvent: any = function(
  type: string,
  element: any,
  callback: Function
) {
  if (element.addEventListener) {
    // Ie9+和chrome下 绑定鼠标滚轮事件
    addEvent = function(_type: string, _element: any, _callback: Function) {
      element.addEventListener(_type, _callback, false);
    };
  } else if (element.attachEvent as Function) {
    // IE8浏览器及一下
    addEvent = function(_type: string, _element: any, _callback: Function) {
      element.attachEvent(`on${_type}`, _callback);
    };
  } else {
    element[`on${type}`] = callback;
  }
  return addEvent(type, element, callback);
};

/**
 * @author Leeyunxiang
 * @description 兼容火狐浏览器滚轮事件的惰性函数
 * @param {Element} element 事件触发者
 * @param {Function} callback 事件的触发回调函数
 * */

export let Wheel: any = function(element: any, callback: Function) {
  // 非火狐浏览器
  if ((window as any).onmousewheel !== undefined) {
    Wheel = function(_element: any, _callback: Function) {
      addEvent("mousewheel", _element, (e: any) => {
        const ev = e || window.event;
        const { wheelDelta } = ev;
        _callback(ev, wheelDelta);
      });
    };
  } else {
    Wheel = function(_element: any, _callback: Function) {
      addEvent("DOMMouseScroll", _element, (e: any) => {
        const ev = e || window.event;
        const data = (ev.detail / -3) * 120; //正负整数取反
        _callback(ev, data);
      });
    };
  }
  return Wheel(element, callback);
};
