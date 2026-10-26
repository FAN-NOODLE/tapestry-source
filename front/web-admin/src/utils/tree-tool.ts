type TreeNode<T> = {
  [P in keyof T]: T[P];
} & {
  id?: number;
  pid?: number;
  children?: Array<TreeNode<T>>;
};

type TreeOptions<T> = {
  id?: keyof T;
  pid?: keyof T;
  children?: keyof T;
};

/**
 * 数组转树
 * @param list 数组
 * @param options 配置
 * @returns
 */
export function createTreeFromList<T>(
  list: TreeNode<T>[],
  options: TreeOptions<TreeNode<T>> = {},
): TreeNode<T>[] {
  const nodeMap = new Map();
  const result: TreeNode<T>[] = [];
  const { id = 'id', pid = 'pid', children = 'children' } = options;

  // 如果数组为空，直接返回
  if (list.length === 0) return [];

  for (const node of list) {
    nodeMap.set(node[id], node);
  }

  for (const node of list) {
    const parent = nodeMap.get(node[pid]);
    if (parent) {
      if (parent[children]) {
        parent[children].push(node);
      } else {
        parent[children] = [node];
      }
    } else {
      result.push(node);
    }
  }

  return result;
}

/**
 * 树转数组
 * @param tree 树
 * @param options 配置
 * @returns
 */
export function treeToList<T>(
  tree: TreeNode<T>[],
  options: TreeOptions<TreeNode<T>> = {},
): TreeNode<T>[] {
  const { children = 'children' } = options;
  const result: TreeNode<T>[] = [...tree];

  for (let i = 0; i < result.length; i++) {
    if (!result[i][children]) continue;
    result.splice(i + 1, 0, ...(result[i][children] as TreeNode<T>[]));
  }
  return result;
}

/**
 * 插入子节点
 * @param tree
 * @param node
 * @param options
 */
export function insertTreeNode<T>(
  tree: TreeNode<T>[],
  node: TreeNode<T>,
  options: TreeOptions<TreeNode<T>> = {},
): TreeNode<T>[] {
  const { pid = 'pid' } = options;
  if (!node[pid]) {
    // 如果父ID不存在，放在最顶层
    tree.push(node);
  } else {
    // 否则查找对应的父节点
    // eslint-disable-next-line @typescript-eslint/no-use-before-define
    const pNode = findTreeNode<T>(tree, (item) => item.id == node[pid]);
    if (pNode) {
      if (pNode.children) {
        pNode.children.push(node);
      } else {
        pNode.children = [node];
      }
    }
  }
  return tree;
}

/**
 * 移除树节点
 * @param tree
 * @param callback
 * @param options
 * @returns
 */
export function removeTreeNode<T>(
  tree: TreeNode<T>[],
  callback: (node: TreeNode<T>) => boolean,
  options: TreeOptions<TreeNode<T>> = {},
): void {
  const { children = 'children' } = options;
  // 如果树没有节点，返回
  if (tree.length === 0) return;

  for (let i = 0; i < tree.length; i++) {
    const item = tree[i];
    // 使用回调判断要移除的节点，不考虑是否有子节点
    if (callback(item)) {
      tree.splice(i, 1);
      break;
    }
    // 如果存在子节点，遍历查询子节点
    if (tree[i][children]) {
      removeTreeNode(tree[i][children] as TreeNode<T>[], callback, options);
    }
  }
}

/**
 * 筛选树节点
 * @param tree
 * @param callback
 * @param options
 * @returns
 */
export function filterTreeNode<T>(
  tree: TreeNode<T>[],
  callback: (treeNode: TreeNode<T>) => boolean,
  options: TreeOptions<TreeNode<T>> = {},
): TreeNode<T>[] {
  const result: TreeNode<T>[] = [...tree];
  const { children = 'children' } = options;

  function filter(list: TreeNode<T>[]): TreeNode<T>[] {
    return list
      .map((node) => ({ ...node }))
      .filter((node) => {
        // 筛选子节点
        if (node[children]) {
          node.children = filter(node[children] as TreeNode<T>[]);
        }
        return callback(node) || (node[children] && (node[children] as TreeNode<T>[]).length);
      });
  }

  return filter(result);
}

/**
 * 获取树节点
 * @param tree
 * @param callback
 * @param options
 * @returns
 */
export function findTreeNode<T>(
  tree: TreeNode<T>[],
  callback: (node: TreeNode<T>) => boolean,
  options: TreeOptions<TreeNode<T>> = {},
): TreeNode<T> | null {
  const { children = 'children' } = options;
  for (const node of tree) {
    if (callback(node)) return node;
    node[children] && tree.push(...(node[children] as TreeNode<T>[]));
  }
  return null;
}

/**
 * 遍历树
 * @param tree 树
 * @param callback
 * @param options
 */
export function forEachTree<T>(
  tree: TreeNode<T>[],
  callback: (node: TreeNode<T>) => void,
  options: TreeOptions<TreeNode<T>> = {},
): void {
  const list: TreeNode<T>[] = [...tree];
  const { children = 'children' } = options;
  for (let i = 0; i < list.length; i++) {
    callback(list[i]);
    list[i][children] && list.splice(i + 1, 0, ...(list[i][children] as TreeNode<T>[]));
  }
}

/**
 * 通过节点ID获取父节点
 * @param tree
 * @param id
 * @returns
 */
export function getParentNodeById<T>(tree: TreeNode<T>[], id: string | number): TreeNode<T> | null {
  const map = new Map();
  for (let i = 0; i < tree.length; i++) {
    map.set(tree[i].id, tree[i]);
  }
  return map.get(map.get(id).pid) || null;
}
