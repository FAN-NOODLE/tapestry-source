/*
 * 权限遍历相关函数封装
 * getBCNodes 获取面包屑数据
 * TODO
 * 获取菜单Tree
 */

import { IBCNodeItem } from "@/types";
// 根据关键id获取所有节点
export function getAllChildNode(data: Array<IBCNodeItem>, id: number) {
  const nodes = data.filter(item => {
    return item.pid === id;
  });

  return nodes;
}
// 获取面包屑数据
export function getBCNodes(data: Array<IBCNodeItem>, node: IBCNodeItem) {
  const result: Array<IBCNodeItem> = [] || undefined;
  if (!data) {
    return; //如果key传空，直接返回
  }
  const forFn = (list: Array<IBCNodeItem>, nodeChild: IBCNodeItem) => {
    const pNode = list.find(item => {
      return item.id === nodeChild.pid;
    });

    if (pNode !== undefined) {
      const childNodes = getAllChildNode(list, pNode.id);
      pNode.children = childNodes;
      result.push(pNode);
      // console.log(pNode);
      forFn(list, pNode);
    }
  };
  // 获取当前路由所有子节点
  const nodes = getAllChildNode(data, node.id);
  node.children = nodes;
  result.push(node);
  forFn(data, node);

  return result.reverse();
}

// 根据关键key获取节点
export function getKeyByBCNodeItem(data: Array<IBCNodeItem>, key: string) {
  if(data instanceof Array){
    const node = data.find(item => {
      return item.key === key;
    });

    return node;
  }
}
