/**
 * @description 引用svg全局注册组件
 * */

import Vue from "vue";
import SvgIcon from "./SvgCom.vue";

Vue.component("svg-icon", SvgIcon);

const req = require.context("./svg", false, /\.svg$/);
const requireAll = (requireContext: any) =>
  requireContext.keys().map(requireContext);
requireAll(req);
