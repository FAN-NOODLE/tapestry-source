//  单一状态树
const state = {
  // 记录本次已发送成功的数据
  nowSentData: [],
  // 树型组件 单选当前数据
  nowChecked: {},
  // 选中项
  selection: [],
  /** 素材详情 */
  detail: null,
};

// getters 暴露出去的获取 state 的计算属性
const getters = {
  getNowSentData: (state: { nowSentData: any }) => {
    return state.nowSentData;
  },
  getNowChecked: (state: { nowChecked: any }) => {
    return state.nowChecked;
  },

  getSelection: (state: { selection: any }) => {
    return state.selection;
  },
  getDetail: (state: { detail: any }) => {
    return state.detail;
  },
};

// mutations 私有化方法，只能被 actions 调用， 可以直接变更状态
const mutations = {
  setNowSentData(state: { nowSentData: any }, data: any) {
    state.nowSentData = data;
  },
  setNowChecked(state: { nowChecked: any }, data: any) {
    state.nowChecked = data;
  },
  setSelection(state: { selection: any }, data: any) {
    state.selection = data;
  },
  setDetail(state: { detail: any }, data: any) {
    state.detail = data;
  },
};

// actions 暴露到外部的方法 不直接变更状态
const actions = {
  setNowSentData({ commit }: any, data: any) {
    commit("setNowSentData", data);
  },
  setNowChecked({ commit }: any, data: any) {
    commit("setNowChecked", data);
  },
  setSelection({ commit }: any, data: any) {
    commit("setSelection", data);
  },
  setDetail({ commit }: any, data: any) {
    commit("setDetail", data);
  },
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations,
};
