//  单一状态树
const state = {
  IS_LOGINED: false, // 是否登录
  ENTERPRISE_HASH: "", // 企业hash
  // ENTERPRISE_ID: "", // 企业Id
  pageType: 0,
  groupChatId: "",
};

// getters 暴露出去的获取 state  的计算属性
const getters = {
  getIsLogined: (state: { IS_LOGINED: any }) => {
    return state.IS_LOGINED;
  },
  getEnterpriseHash: (state: { ENTERPRISE_HASH: any }) => {
    return state.ENTERPRISE_HASH;
  },
  // getEnterpriseId: (state: { ENTERPRISE_ID: any }) => {
  //   return state.ENTERPRISE_ID;
  // },
  getPageType: (state: { pageType: any }) => {
    return state.pageType;
  },
  getGroupChatId: (state: { groupChatId: any }) => {
    return state.groupChatId = state.groupChatId || sessionStorage.getItem('groupChatId');
  },
};

// mutations 私有化方法，只能被 actions 调用， 可以直接变更状态
const mutations = {
  setIsLogined(state: { IS_LOGINED: any }, data: any) {
    state.IS_LOGINED = data;
  },
  setEnterpriseHash(state: { ENTERPRISE_HASH: any }, data: any) {
    state.ENTERPRISE_HASH = data;
  },
  // setEnterpriseId(state: { ENTERPRISE_ID: any }, data: any) {
  //   state.ENTERPRISE_ID = data;
  // },
  setPageType(state: { pageType: any }, data: any) {
    state.pageType = data;
  },
  setGroupChatId: (state: { groupChatId: any }, data: any) => {
    sessionStorage.setItem('groupChatId', data)
    state.groupChatId = data;
  },
};

// actions 暴露到外部的方法 不直接变更状态
const actions = {
  setIsLogined({ commit }: any, data: any) {
    commit("setIsLogined", data);
  },
  setEnterpriseHash({ commit }: any, data: any) {
    commit("setEnterpriseHash", data);
  },
  // setEnterpriseId({ commit }: any, data: any) {
  //   commit("setEnterpriseId", data);
  // },
  setPageType({ commit }: any, data: any) {
    commit("setPageType", data);
  },
  setGroupChatId: ({ commit }: any, data: any) => {
    commit("setGroupChatId", data);
  },
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations,
};
