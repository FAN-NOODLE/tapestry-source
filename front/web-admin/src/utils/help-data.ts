import BreadCrumbModule from "@/store/modules/bread-crumb";
import { getConfig } from "@/utils/help-config";

/**
 * @description 这个方法是判断当前的目录
 * */

export async function hasHelp() {
  try {
    const list = BreadCrumbModule.getBreadCrumb; //获取当前面包展示的面包屑list
    const node = list[list.length - 1]; // 得到最后一个面包屑节点
    const config =
      node.children.length > 0 //当面包屑有子节点时获取第一个子节点 原因：当前展示的view视图为第一个子节点
        ? await getConfig(node.children[0].key) // 得到当前view视图的key 通过key得到对应文案
        : await getConfig(node.key); // 得到当前view视图的key 通过key得到对应文案]
    return config !== undefined;
  } catch (error) {
    return false;
  }
}
// TODO 注释
export async function helpIsShow() {
  const list = BreadCrumbModule.getBreadCrumb; // 获取当前面包展示的面包屑list
  const node = list[list.length - 1]; // 得到最后一个面包屑节点
  const config =
    node.children.length > 0 //当面包屑有子节点时获取第一个子节点 原因：当前展示的view视图为第一个子节点
      ? await getConfig(node.children[0].key) // 得到当前view视图的key 通过key得到对应文案
      : await getConfig(node.key); // 得到当前view视图的key 通过key得到对应文案

  // 获取当前本地存储状态 默认开启
  const status = localStorage.getItem("HELP") || "1";
  return {
    isShow: (status === undefined || status === "1") && config !== undefined,
    config:
      config !== undefined
        ? config
        : { key: "key", type: "text", content: "默认" }
  };
}
