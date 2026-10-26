// 检测 map end
let result = (ary, resolve, target) => {
  if (ary.length === 0) {
    resolve(target);
  }
  return ary.length !== 0;
};

export default {
  // 找到菜单的第一个子叶节点
  getFirstChildMenu: function(node) {
    if (
      node.children &&
      node.children.length > 0 &&
      node.children[0].type !== "button"
    ) {
      let tmp = node.children[0];
      return this.getFirstChildMenu(tmp);
    } else {
      return node;
    }
  },
  // 使用 promise 封装的 deep map
  deepMapPromise: ({ data, childField = "child", callback }) => {
    let target = data;
    return new Promise(resolve => {
      if (target && childField) {
        let nodes = [],
          node = {};
        // 类型判断
        target.constructor === Array
          ? (node[childField] = target)
          : (node = target);
        let stack = []; //同来存放将来要访问的节点
        stack.push(node);
        while (result(stack, resolve, target)) {
          let item = stack.pop(); //正在访问的节点
          nodes.push(item);
          if (callback) {
            // 回调函数处理 item
            callback(item);
          }
          let childrens = item[childField];
          if (childrens) {
            for (let i = childrens.length - 1; i >= 0; i--) {
              stack.push(childrens[i]);
            }
          }
        }
      }
    });
  }
};
