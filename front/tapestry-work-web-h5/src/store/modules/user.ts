//  单一状态树

const TOKEN = localStorage.getItem('token')
const state = {
  USERINFO: {
    user_id: "",
  },
  TOKEN,
  externalId: "",
  customerId: 0,
  enterprise_mark: "",
  role: ''
};

// getters 暴露出去的获取 state 的计算属性
const getters = {
  getUserInfo: (state: { USERINFO: any }) => {
    return state.USERINFO;
  },
  getRole: (state: { role: any }) => {
    return state.role;
  },
  getCustomerId: (state: { customerId: number }) => {
    return state.customerId;
  },
  getExternalId: (state: { externalId: string }) => {
    return state.externalId;
  },
  getToken: (state: { TOKEN: string }) => {
    return state.TOKEN;
  },
  getEnterpriseMark: (state: { enterprise_mark: string }) => {
    return state.enterprise_mark;
  },
};

// mutations 私有化方法，只能被 actions 调用， 可以直接变更状态
const mutations = {
  setUserInfo(state: { USERINFO: any }, data: any) {
    state.USERINFO = data;
  },
  setCustomerId: (state: { customerId: number }, data: any) => {
    state.customerId = data;
  },
  setRole: (state: { role: number }, data: any) => {
    state.role = data;
  },
  setExternalId(state: { externalId: string }, data: string) {
    state.externalId = data;
  },
  setToken(state: { TOKEN: string }, data: string) {
    state.TOKEN = data;
  },
  setEnterpriseMark(state: { enterprise_mark: string }, data: string) {
    state.enterprise_mark = data;
  },
};

// actions 暴露到外部的方法 不直接变更状态
const actions = {
  setUserInfo({ commit }: any, data: any) {
    commit("setUserInfo", data);
  },
  setExternalId({ commit }: any, data: string) {
    commit("setExternalId", data);
  },
  setRole({ commit }: any, data: string) {
    commit("setRole", data);
  },
  setToken({ commit }: any, data: string) {
    commit("setToken", data);
  },
  setEnterpriseMark({ commit }: any, data: string) {
    commit("setEnterpriseMark", data);
  },
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations,
};
