import Vue from "vue";
import SelectStorePopup from "./index.vue";
import { i18n } from '@/locales'

const DialogConstructor = Vue.extend(SelectStorePopup);

type SelectStorePopupInstance = Vue & {
  show: boolean;
};

let instance: SelectStorePopupInstance;


function open(options): Promise<{action: string, data:any}> {
   instance = new DialogConstructor({
    propsData: {
      ...options,
    },
    i18n
  }) as SelectStorePopupInstance

  return new Promise((resolve, reject) => { 

    instance.$mount(); // 挂载实例
    document.body.appendChild(instance.$el); // 将实例添加到 DOM 中

    instance.$on('cancel', (res: any) => {
      document.body.removeChild(instance.$el); // 移除 DOM
      instance.$destroy(); // 销毁实例
      options.onCancel && options.onCancel(res);
      resolve({
        action: 'cancel',
        data: res,
      });
    });

    instance.$on('select', (res: any) => {
      // document.body.removeChild(instance.$el); // 移除 DOM
      // instance.$destroy(); // 销毁实例
      options.onSelect && options.onSelect(res);
    });

    instance.$on('confirm', (res: any) => {
      document.body.removeChild(instance.$el); // 移除 DOM
      instance.$destroy(); // 销毁实例
      options.onConfirm && options.onConfirm(res);
      resolve({
        action: 'confirm',
        data: res,
      });
    });

    instance.$on('close', () => {
      document.body.removeChild(instance.$el); // 移除 DOM
      instance.$destroy(); // 销毁实例
      options.onClose && options.onClose();
      resolve({
        action: 'close',
        data: null,
      })
    });

    instance.show = true; // 显示 Dialog
    
  })
  
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