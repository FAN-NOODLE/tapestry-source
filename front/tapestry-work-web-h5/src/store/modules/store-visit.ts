import dayjs from "dayjs";

//  单一状态树
const state = {
  /** 重新加载列表数据 */
  update: false,
  toast: '',
  updateId: 0,
  title: ''
};

// getters 暴露出去的获取 state 的计算属性
const getters = {
  getUpdate: (state: any) => {
    return state.update;
  },
  getToast: (state: any) => {
    return state.toast;
  },
  getUpdateId: (state: any) => {
    return state.updateId;
  },
  getTitle: (state: any) => {
    return state.title;
  },
};

// mutations 私有化方法，只能被 actions 调用， 可以直接变更状态
const mutations = {
  setUpdate(state: any, data: any) {
    state.update = data;
  },
  setToast(state: any, data: any) {
    state.toast = data;
  },
  setUpdateId(state: any, data: any) {
    state.updateId = data;
  },
  setTitle(state: any, data: any) {
    state.title = data;
  },
};

// actions 暴露到外部的方法 不直接变更状态
const actions = {
  setUpdate({ commit }: any, data: any) {
    commit("setUpdate", data);
  },
  setToast({ commit }: any, data: any) {
    commit("setToast", data);
  },
  setUpdateId({ commit }: any, data: any) {
    commit("setUpdateId", data);
  },
  setTitle({ commit }: any, data: any) {
    commit("setTitle", data);
  },
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations,
};
