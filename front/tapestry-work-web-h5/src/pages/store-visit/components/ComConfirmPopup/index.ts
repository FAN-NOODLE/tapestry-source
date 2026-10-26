import Vue from "vue";
import ComConfirmPopup from "./index.vue";

const DialogConstructor = Vue.extend(ComConfirmPopup);

type ConfirmPopupInstance = Vue & {
  show: boolean;
};
function confirm(options): Promise<{action: string}> {
  return new Promise((resolve, reject) => {
    const instance = new DialogConstructor({
      propsData: {
        ...options,
      },
      
    }) as ConfirmPopupInstance

    instance.$mount(); // 挂载实例
    document.body.appendChild(instance.$el); // 将实例添加到 DOM 中

    instance.$on('confirm', (res: any) => {
      document.body.removeChild(instance.$el); // 移除 DOM
      instance.$destroy(); // 销毁实例
      resolve(res); // 解析 Promise
    });

    instance.$on('close', () => {
      document.body.removeChild(instance.$el); // 移除 DOM
      instance.$destroy(); // 销毁实例
      reject(); // 拒绝 Promise
    });

    instance.show = true; // 显示 Dialog
  });
}

export default {
  confirm
}