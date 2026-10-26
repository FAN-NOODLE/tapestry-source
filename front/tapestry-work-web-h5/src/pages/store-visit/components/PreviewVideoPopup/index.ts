import Vue from "vue";
import PreviewVideoPopup from "./index.vue";

const DialogConstructor = Vue.extend(PreviewVideoPopup);

type PreviewVideoPopupInstance = Vue & {
  show: boolean;
};

let instance: PreviewVideoPopupInstance;

function open(options): any {
   instance = new DialogConstructor({
    propsData: {
      ...options,
    },
    
  }) as PreviewVideoPopupInstance

  instance.$mount(); // 挂载实例
  document.body.appendChild(instance.$el); // 将实例添加到 DOM 中

  instance.$on('close', () => {
    document.body.removeChild(instance.$el); // 移除 DOM
    instance.$destroy(); // 销毁实例
  });

  instance.show = true; // 显示 Dialog
  return instance
}

function close() {
  if(instance) {
    instance.show = false;
  }
}

export default {
  open,
  close
}