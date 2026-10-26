//  单一状态树
const state = {
  MENU: null,
  MENU_TOP: null,
  MENU_LEFT: null
};

// getters 暴露出去的获取 state 的计算属性
const getters = {
  getMenu: state => {
    return state.MENU;
  },
  getMenuTop: state => {
    return state.MENU_TOP;
  },
  getMenuLeft: state => {
    return state.MENU_LEFT;
  }
};

// mutations 私有化方法，只能被 actions 调用， 可以直接变更状态
const mutations = {
  // 设置全量 menu
  setMenu: (state, data) => {
    state.MENU = data;
  },
  setMenuTop: (state, data) => {
    state.MENU_TOP = data;
  },
  setMenuLeft: (state, data) => {
    state.MENU_LEFT = data;
  }
};

// actions 暴露到外部的方法 不直接变更状态
const actions = {
  setMenu: ({ commit }, data) => {
    commit("setMenu", data);
  },
  setMenuTop: ({ commit }, data) => {
    commit("setMenuTop", data);
  },
  setMenuLeft: ({ commit }, data) => {
    commit("setMenuLeft", data);
  }
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
};
