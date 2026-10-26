// /**
//  * @author maliang
//  * @description 同过 map 模块配置文件，生成并注册路由，采用深度遍历，能适应无限层级子路由情况。
//  */

// import Myapp from "@/views/myapp/myapp.js";
// import Layout from "@/views/Layout";
// import Store from "../store/index";

// const BASE_URL = "/myapp", // 基础 url
//   BASE_KID = "myapp"; // 基础 ID
// let routers = Myapp.concat(); // 深拷贝 modules 配置

// // deep map end
// let result = (ary, resolve, target) =>
//   ary.length === 0 ? resolve(target) : ary.length !== 0;

// // 设置 node path
// let setPath = node => {
//   node.path = node.ename;
//   return node;
// };
// // 设置 node 是否显示
// let setShow = node => {
//   node.show = node.hide ? false : true;
//   return node;
// };
// // 设置 node 的文件路径以及 key
// let setFilePathAndKey = node => {
//   if (node.fid) {
//     node.file_path = `${node.fid}/${node.ename}`;
//     node.key = `${node.kid}_${node.ename}`;
//   } else {
//     node.file_path = `${BASE_URL}/${node.ename}`;
//     node.key = `${BASE_KID}_${node.ename}`;
//   }
//   // 根据 config 生成 router name
//   node.name = node.key;
//   return node;
// };

// // 设置 component
// let setComponent = node => {
//   // lazy load module
//   node.component = () => import(`@/views${node.file_path}/${node.ename}.vue`);
//   return node;
// };

// // 设置 router meta
// let setMeta = node => {
//   // 路由存储字段
//   node.meta = {
//     title: node.cname,
//     is_check_auth: node.is_authorization,
//     router_key: node.key
//   };
//   return node;
// };

// // 处理 node 按钮
// let setButton = node => {
//   // 页面按钮生成 key
//   node.button_keys = [];
//   if (node.buttons && node.buttons.length > 0) {
//     node.buttons.map(child => {
//       child.is_authorization = true;
//       if (!child.key) {
//         child.key = `${node.key}_${child.ename}`;
//         child.path = "/BUTTON_NULL"; // vue-router node must use path
//         child.type = "button";
//         node.button_keys.push(child.key);
//       }
//     });
//   }
//   /* （为了性能优化，减少循环次数）这种情况下 router 会将页面内部按钮作默认为路由, 但没有真实路由, 所以理论上不影响使用。 */
//   if (node.key && node.buttons && node.buttons.length > 0) {
//     if (node.children) {
//       node.children = node.children.concat(node.buttons);
//     } else {
//       node.children = node.buttons;
//     }
//   }
//   return node;
// };

// let saveKeys = node => {
//   // 判断 type 将值存入vuex (包括 page 内部子权限)
//   switch (node.type) {
//     case "base": // 基础权限
//       Store.dispatch("APP/setKeys", {
//         type: "BASE_MOD_KEYS",
//         data:
//           node.button_keys.length > 0
//             ? node.button_keys.concat(node.key)
//             : node.key
//       });
//       break;
//     case "crs": // crs 权限
//       Store.dispatch("APP/setKeys", {
//         type: "CRS_MOD_KEYS",
//         data:
//           node.button_keys.length > 0
//             ? node.button_keys.concat(node.key)
//             : node.key
//       });
//       break;
//     case "shop": // shop 权限
//       Store.dispatch("APP/setKeys", {
//         type: "SHOP_MOD_KEYS",
//         data:
//           node.button_keys.length > 0
//             ? node.button_keys.concat(node.key)
//             : node.key
//       });
//       break;
//     case "app": // app 权限
//       Store.dispatch("APP/setKeys", {
//         type: "APP_MOD_KEYS",
//         data:
//           node.button_keys.length > 0
//             ? node.button_keys.concat(node.key)
//             : node.key
//       });
//       break;
//   }
//   Store.dispatch("APP/setKeys", {
//     type: "ALL_MOD_KEYS",
//     data: node.key
//   });
//   // 临时变量用于 临时存储 node button keys 存储完后立即删除
//   delete node.button_keys;
//   return node;
// };

// /**
//  * @method
//  * @description 深度优先遍历
//  * @param {Object/Array} obj
//  * @param {String} childName
//  * @param {Function} callback
//  */
// let deepMapRouter = data => {
//   return new Promise(resolve => {
//     // 判断 参数类型
//     if (!data) resolve({});
//     let node = {};
//     // 类型判断
//     data.constructor === Array ? (node.children = data) : (node = data);
//     let stack = []; //同来存放将来要访问的节点
//     stack.push(node);
//     while (result(stack, resolve, data)) {
//       let child = stack.pop(); //正在访问的节点
//       let childrens = child.children;
//       // 存在 ename 证明其必然存在 key
//       if (child.ename) {
//         // 解析路由
//         [
//           setPath,
//           setShow,
//           setFilePathAndKey,
//           setMeta,
//           setButton,
//           saveKeys,
//           setComponent
//         ].reduce((x, f) => f(x), child);
//       }
//       if (childrens && childrens.length > 0) {
//         childrens.map(children => {
//           children.$f = child; // 存入父节点供子节点使用
//           children.fid = child.file_path || child.ename; // 记录 children base file_path
//           children.kid = child.key || child.ename; // 设置 children base key
//           stack.push(children);
//         });
//       }
//     }
//   });
// };

// // 初始化
// (async function init() {
//   await deepMapRouter(routers); // 全量 Menu List, 用于 Menu 渲染
//   // 初始化菜单数据
//   Store.dispatch("MENU/setMenuTop", routers.concat());
//   // 全量权限数据
//   Store.dispatch("MENU/setMenu", routers.concat());
// })();

// export default {
//   path: "/myapp",
//   component: Layout,
//   redirect: routers[0].file_path,
//   children: routers
// };
