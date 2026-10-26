const state = {
  PERMISSION: null,
  ROLES: null
};

const getters = {
  getPermission: state => {
    if (!state.PERMISSION) {
      state.PERMISSION = JSON.parse(localStorage.getItem("PERMISSION"));
    }
    return state.PERMISSION;
  },
  getRoleList: state => {
    if (!state.ROLES) {
      state.ROLES = JSON.parse(localStorage.getItem("ROLES"));
    }
    return state.ROLES;
  }
};

const mutations = {
  setPermission(state, data) {
    localStorage.setItem("PERMISSION", JSON.stringify(data));
    state.PERMISSION = data;
  },
  setRoleList(state, data) {
    localStorage.setItem("ROLES", JSON.stringify(data));
    state.ROLES = data;
  }
};

const actions = {
  setPermission({ commit }, data) {
    commit("setPermission", data);
  },
  setRoleList({ commit }, data) {
    commit("setRoleList", data);
  }
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
};
