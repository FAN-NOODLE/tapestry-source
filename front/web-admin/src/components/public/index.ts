import Upload from "./cy-upload/index.vue";
import Table from "./cy-table/index.vue";
import Input from "./cy-input/index.vue";
import { Form, FormItem } from "./cy-form";
import { Select, Option } from "./cy-select";
import CasCader from "./cy-cascader/index.vue";
import LayoutForm from "./layout-form/index.vue";
import LayoutList from "./layout-list/index.vue";
import SearchItem from "./search-item/index.vue";
/**
 * 整理组件
 * */

const components = {
  "cy-upload": Upload,
  "cy-table": Table,
  "cy-input": Input,
  "cy-select": Select,
  "cy-option": Option,
  "cy-form": Form,
  "cy-form-item": FormItem,
  "cy-cascader": CasCader,
  "layout-form": LayoutForm,
  "layout-list": LayoutList,
  "search-item": SearchItem,
};

/**
 * 全局注册
 * */

export default {
  install(Vue: any) {
    for (const key in components) {
      Vue.component(key, components[key]);
    }
  }
};
