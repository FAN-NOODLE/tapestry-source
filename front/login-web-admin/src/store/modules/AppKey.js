//  单一状态树
const state = {
  ALL_MOD_KEYS: new Set(), // 全部应用 keys
  BASE_MOD_KEYS: new Set(), // 基础应用 keys
  CRS_MOD_KEYS: new Set(), // crs 应用 keys
  SHOP_MOD_KEYS: new Set(), // shop 应用 keys
  APP_MOD_KEYS: new Set(), // child app keys
  /* 初始应用 */
  MOD_APP_INIT_KEYS: new Set([
    "myapp_app",
    "myapp_app_myapplication",
    "myapp_app_application"
  ]) // 基础应用
};

// getters 暴露出去的获取 state 的计算属性
const getters = {
  getKeys: state => {
    return state;
  }
};

// mutations 私有化方法，只能被 actions 调用， 可以直接变更状态
const mutations = {
  // 设置 keys
  setKeys: (state, obj) => {
    if (typeof obj.data === "string") {
      state[obj.type].add(obj.data);
    } else {
      state[obj.type] = new Set([...state[obj.type], ...obj.data]);
    }
  }
};

// actions 暴露到外部的方法 不直接变更状态
const actions = {
  setKeys: ({ commit }, type, data) => {
    commit("setKeys", type, data);
  }
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
};
