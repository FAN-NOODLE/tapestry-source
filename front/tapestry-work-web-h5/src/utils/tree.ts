
export interface TreeNode {
  children?: TreeNode[];
  parent?: TreeNode;
  [key: string]: any;
}

/**
 * 将树转换为列表
 * @param treeNode 树节点
 * @param level 级别，默认为0
 * @returns 列表项数组
 */
export function treeToList(treeNode: TreeNode, level: number = 0): TreeNode[] {
  // 创建结果数组，并添加树节点对象
  const result: TreeNode[] = [treeNode];

  // 如果树节点有子节点，则递归调用treeToList函数，并将结果添加到result数组中
  if (treeNode.children) {
    for (const child of treeNode.children) {
      result.push(...treeToList(child, level + 1));
    }
  }

  // 返回列表项数组
  return result;
}


/**
 * 递归函数：遍历树并添加 level 和 isLeaf 属性
 * @param node 
 * @param level 
 * @returns 
 */
export function traverseTree(node: TreeNode, level: number = 0, callback: (node: TreeNode) => TreeNode): TreeNode {
  // 添加 level 属性
  node.level = level;

  // 添加 isLeaf 属性
  node.isLeaf = !node.children || node.children.length === 0;

  node.parent = node.parent || null

  // 递归处理子节点
  if (node.children) {
    
    node.children.forEach(child => {
      child.parent = node;
      child = traverseTree(child, level + 1, callback)
    });
  }

  return callback(node);
}

/**
 * 遍历树形列表
 * @param list 树形列表
 * @returns 遍历后的树形列表
 */
export function traverseTreeList<T extends TreeNode>(list: T[], callback: (node: T) => T) {
  return list.map((item) => traverseTree(item, 0, callback)) as T[];
}

export function bfs(root: TreeNode, callback: (node: TreeNode) => void) {
  const stack: TreeNode[] = [root];
  while (stack.length) {
    const node = stack.shift() as TreeNode;
    callback(node);
    if (node.children) {
      stack.push(...node.children);
    }
  }
}

export function dfs(root: TreeNode, callback: (node: TreeNode) => void) {
  const stack: TreeNode[] = [root];
  while (stack.length) {
    const node = stack.pop() as TreeNode;
    callback(node);
    if(node.children) {
      stack.push(...node.children);
    }
  }
}

export function searchTree(node: TreeNode, condition: (node: TreeNode) => boolean): TreeNode[] {
  let result: TreeNode[] = [];

  // 检查当前节点是否满足条件
  if (condition(node)) {
    result.push(node);
  }

  // 递归处理子节点
  if (node.children) {
    for (const child of node.children) {
      result = result.concat(searchTree(child, condition));
    }
  }

  return result;
}

export function hasMatchingChild(node: TreeNode, condition: (child: TreeNode) => boolean): boolean {
  if (node.children) {
    // 遍历当前节点的子节点
    for (const child of node.children) {
      // 检查子节点是否符合条件
      if (condition(child)) {
        return true; // 如果找到符合条件的子节点，返回 true
      }

      // 递归检查子节点的子节点
      if (hasMatchingChild(child, condition)) {
        return true; // 如果在子节点的子节点中找到符合条件的节点，返回 true
      }
    }
  }

  return condition(node); // 如果没有符合条件的子节点，返回 false
}


export function searchTreeList(tree: TreeNode[], condition: (node: TreeNode) => boolean): TreeNode[] {
  const list = tree.reduce((prev, item) => {
    return prev.concat(treeToList(item));
  }, [])
  return list.filter(item => hasMatchingChild(item, condition))
}


/**
 * 获取树节点
 * @param tree
 * @param callback
 * @param options
 * @returns
 */
export function findTreeNode<T>(
  tree: TreeNode[],
  callback: (node: TreeNode) => boolean,
): TreeNode | null {
  for (const node of tree) {
    if (callback(node)) return node;
    node.children && tree.push(...(node.children as TreeNode[]));
  }
  return null;
}


/**
 * 遍历树
 * @param tree 树
 * @param callback
 * @param options
 */
export function forEachTree(
  tree: TreeNode[],
  callback: (node: TreeNode) => void,
): void {
  const list: TreeNode[] = [...tree];
  for (let i = 0; i < list.length; i++) {
    callback(list[i]);
    list[i].children && list.splice(i + 1, 0, ...(list[i].children as TreeNode[]));
  }
}


// 广度优先遍历
export function breadthFirstTraversal(
  tree: TreeNode[],
  callback: (node: TreeNode) => void,
): void {
  const queue: TreeNode[] = [...tree];
  while (queue.length) {
    const node = queue.shift()!;
    callback(node);
    node.children && queue.push(...(node.children as TreeNode[]));
  }
}