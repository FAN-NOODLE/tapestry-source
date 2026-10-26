//  单一状态树

const TOKEN = localStorage.getItem('token')
const state = {
  lang: localStorage.getItem('lang') || 'zh-CN',
};

// getters 暴露出去的获取 state 的计算属性
const getters = {
  getLang: (state: any) => {
    return state.lang;
  },
};

// mutations 私有化方法，只能被 actions 调用， 可以直接变更状态
const mutations = {
  setLang(state: any, data: any) {
    state.lang = data;
    localStorage.setItem('lang', data);
  },
};

// actions 暴露到外部的方法 不直接变更状态
const actions = {
  setLang({ commit }: any, data: any) {
    commit("setLang", data);
  },
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations,
};
