import Store from "../store/index";
import Ciphertext from "./Ciphertext";

let result = {}; // login response
let super_account = true, // is super account
  server_type = 1, // 应用权限
  pass = "", // result pass
  all_keys = Store.state.APP.APP_MOD_KEYS, // all keys
  parse = Ciphertext("decode"); // 解密

// 根据 role 生成
let createBaseKeys = () => {
  // 超管
  if (super_account) {
    return Store.state.APP.BASE_MOD_KEYS;
  } else {
    // 子管理员
    return new Set(parse(result.powers).split(","));
  }
};

// 根据 serve_type 生成 app keys
let createAppkeys = () => {
  let tmp = null,
    app = Store.state.APP;
  switch (server_type) {
    case 0: // 有 shop 和 crs 权限
      tmp = new Set([
        ...app.BASE_MOD_KEYS,
        ...app.CRS_MOD_KEYS,
        ...app.SHOP_MOD_KEYS,
        ...app.MOD_APP_INIT_KEYS
      ]);
      break;
    case 1: // 有 shop 权限
      tmp = new Set([
        ...app.BASE_MOD_KEYS,
        ...app.SHOP_MOD_KEYS,
        ...app.MOD_APP_INIT_KEYS
      ]);
      break;
    case 2: // 有 crs 权限
      tmp = new Set([...app.BASE_MOD_KEYS, ...app.CRS_MOD_KEYS]);
      break;
  }
  return tmp;
};

// 生成子应用 keys (暂时只有 shop keys)
let createChildAppKeys = () => {
  let apps = result.power,
    tmp_set = new Set();
  if (apps && apps.shop_power && server_type !== 2) {
    apps.shop_power.map(app => {
      let tmp_f_key = app.power_alias;
      let len = tmp_f_key.split("_").length;
      all_keys.forEach(key => {
        let tmp_key = key;
        let tmp = tmp_key.split("_");
        tmp.length = len;
        tmp = tmp.join("_");
        if (tmp === app.power_alias) {
          tmp_set.add(key);
        }
      });
    });
  }
  return tmp_set;
};

// 过滤应用 key
let filter = tmp => {
  let filter = null;
  if (server_type === 0) {
    return;
  }
  if (server_type === 1) {
    filter = new Set([
      ...Store.state.APP.CRS_MOD_KEYS,
      ...Store.state.APP.APP_MOD_KEYS
    ]);
  }
  if (server_type === 2) {
    filter = new Set([
      ...Store.state.APP.SHOP_MOD_KEYS,
      ...Store.state.APP.APP_MOD_KEYS
    ]);
  }
  tmp.forEach(key => {
    if (filter.has(key)) {
      tmp.delete(key);
    }
  });
};

/**
 * @description
 *  step 1. 判断账号是否为超管。
 *  step 2. 判断应用权限 shop/crs/shop&crs
 *  step 3. 解析 shop || crs 子应用权限
 *    A 超管流程：
 *    step 1: 根据 step 2 解析的应用权限生成超管 keys
 *    step 2: 合并 step 3 生成的 shop || crs 子应用 keys
 *    step 3: 返回 keys
 *    B 普通管理员流程：
 *    step 1: 判断是否能拿到用户 keys 拿不到就使用默认 keys
 *    step 2: 合并 step 3 生成的 shop || crs 子应用 keys
 *    step 3: 根据 step 2 应用权限判断是否过滤并过滤 old keys
 *    step 4: 返回 keys
 *  step 4: 解析权限完成
 * @param {object} res login response
 * @returns {string} 权限加密字符串
 */
export default res => {
  result = Object.assign({}, res);
  super_account = parseInt(result.role) === 0;
  server_type = parseInt(result.serve_type);
  let base = createBaseKeys(); // create base key
  let app = super_account ? createAppkeys() : new Set(); // create app key
  let child_app_key = super_account ? createChildAppKeys() : new Set(); // create child app keys
  pass = new Set([...base, ...app, ...child_app_key]);
  if (!super_account) {
    filter(pass);
  }
  return new Promise(resolve => {
    resolve(Array.from(pass));
  });
};
