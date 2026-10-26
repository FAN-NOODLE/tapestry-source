import MenuModule from "@/store/modules/menu";
import AppModule from "../store/modules/app";
import Myapp from "@/views/myapp/myapp";
import Layout from "@/layout/layout.vue";
import { RouteConfig } from "vue-router";
interface Normalarray {
  path: string;
  [index: number]: any;
}
interface Normalobject {
  path: string;
  [propName: string]: any;
}
interface NodeItem {
  children?: Array<NodeItem>;
  path: string;
  name: string;
  ename?: string;
  cname?: string;
  icon: string;
  type: number;
  root?: string;
  meta: any;
  key: string;
  file_path: string;
  hide?: boolean;
  show: boolean;
  component: any;
  fid: string;
  kid: string;
  redirect?: string;
}
const BASE_URL = "/myapp";
const BASE_KID = "myapp";
const routers: Array<Normalarray> = [...Myapp];
const result = (
  ary: Array<Normalarray>,
  resolve: Function,
  target: Array<Normalarray>
) => {
  return ary.length === 0 ? resolve(target) : ary.length !== 0;
};
// 设置 node path
const setPath = (node: NodeItem) => {
  node.path = node.ename ? node.ename : "";
  return node;
};
// 设置 node 是否显示
const setShow = (node: NodeItem) => {
  node.show = node.hide ? false : true;
  return node;
};

// 设置 node 的文件路径以及 key
const setFilePathAndKey = (node: NodeItem) => {
  if (node.fid) {
    node.file_path = `${node.fid}/${node.ename}`;
    node.key = `${node.kid}_${node.ename}`;
  } else {
    node.file_path = `${BASE_URL}/${node.ename}`;
    node.key = `${BASE_KID}_${node.ename}`;
  }
  // 根据 config 生成 router name
  node.name = node.key;
  return node;
};
// 设置 component
const setComponent = (node: NodeItem) => {
  // lazy load module
  node.component = () =>
    import(
      /* webpackChunkName: "[request]" */ `@/views${node.meta.file_path}/${node.meta.ename}.vue`
    );
  return node;
};
const setMeta = (node: NodeItem) => {
  // 路由存储字段
  node.meta = {
    title: node.cname,
    ename: node.ename,
    icon: node.icon,
    type: node.type,
    file_path: node.file_path,
    hide: node.hide,
    key: node.key,
    root: node.root || "",
    model: false
  };
  return node;
};
const setRedirect = (node: NodeItem) => {
  if (node.children && node.children.length > 0) {
    node.redirect = `${node.file_path}/${node.children[0].ename}`;
  }
  return node;
};
const saveKeys = (node: NodeItem) => {
  // 判断 type 将值存入vuex (包括 page 内部子权限)
  AppModule.setKeys(node.key);
  // 临时变量用于 临时存储 node button keys 存储完后立即删除
  delete node.ename;
  delete node.cname;
  delete node.hide;
  delete node.root;
  return node;
};

/**
 * @method
 * @description 深度优先遍历
 * @param {Object/Array} obj
 * @param {String} childName
 * @param {Function} callback
 */
const deepMapRouter = (data: Array<Normalarray>) => {
  return new Promise((resolve: Function) => {
    // 判断 参数类型
    if (!data) resolve({});
    const node: Normalobject = {
      path: ""
    };
    // 类型判断
    node.children = data;
    // data.constructor === Array ? (node.children = data) : (node = data);
    const stack: Array<Normalarray> = []; //同来存放将来要访问的节点
    stack.push(node);
    while (result(stack, resolve, data)) {
      const child: Normalobject | any = stack.pop(); //正在访问的节点
      const childrens = child.children;
      // 存在 ename 证明其必然存在 key
      if (child.ename) {
        // 解析路由
        const setList = [
          setPath,
          setShow,
          setFilePathAndKey,
          setMeta,
          setRedirect,
          saveKeys,
          setComponent
        ];
        const arr = setList.reduce((x, f) => f(x), child);
        console.log(arr);
      }
      if (childrens && childrens.length > 0) {
        childrens.map((children: Normalobject) => {
          children.$f = child; // 存入父节点供子节点使用
          children.fid = child.file_path || child.ename; // 记录 children base file_path
          children.kid = child.key || child.ename; // 设置 children base key
          stack.push(children);
        });
      }
    }
  });
};

// 初始化
(async function init() {
  await deepMapRouter(routers); // 全量 Menu List, 用于 Menu 渲染
  // 全量权限数据
  const realRouter: Array<RouteConfig> = [...routers];
  MenuModule.setMenu(realRouter);
})();

export default {
  path: "/myapp",
  redirect: "/myapp/enterprise",
  component: Layout,
  children: routers
};
