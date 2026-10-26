//  单一状态树
const state = {
  resCustomerDetail: null,
};

// getters 暴露出去的获取 state 的计算属性
const getters = {
  getResCustomerDetail: (state: { resCustomerDetail: string }) => {
    return state.resCustomerDetail;
  },
};

// mutations 私有化方法，只能被 actions 调用， 可以直接变更状态
const mutations = {
  setResCustomerDetail(state: { resCustomerDetail: string }, data: string) {
    state.resCustomerDetail = data;
  },
};

// actions 暴露到外部的方法 不直接变更状态
const actions = {
  setResCustomerDetail({ commit }: any, data: string) {
    commit("setResCustomerDetail", data);
  },
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations,
};
