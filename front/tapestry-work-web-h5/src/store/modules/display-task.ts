//  单一状态树
const state = {
  // 详情信息
  detail: {
    store_info: null,
    task_info: null
  },
  /** 是否更新 */
  update: false,
  /** 更新列表？ */
  update_list: {
    update: false,
    current: -1
  },
  search_params: {
    /** 任务ID */
    ids: [],
    /** 门店节点ID */
    area_ids: [],
    /** 任务名称 */
    task_name: '',
    /** 任务数量 */
    task_num: 0,
    /** 区域数量 */
    area_num: 0,
    /** 门店数量 */
    store_num: 0,
    /** 全部区域门店 */
    is_total: false
  }
};

// getters 暴露出去的获取 state 的计算属性
const getters = {
  getDetail: (state: any) => {
    return state.detail;
  },
  getUpdate: (state: any) => {
    return state.update;
  },
  getUpdateList: (state: any) => {
    return state.update_list;
  },
  getSearchParams: (state: any) => {
    return state.search_params;
  },
};

// mutations 私有化方法，只能被 actions 调用， 可以直接变更状态
const mutations = {
  setDetail(state: any, data: any) {
    state.detail = data;
  },
  setUpdate(state: any, data: any) {
    state.update = data;
  },
  setUpdateList(state: any, data: any) {
    state.update_list = data;
  },
  setSearchParams(state: any, data: any) {
    state.search_params = data;
  },
};

// actions 暴露到外部的方法 不直接变更状态
const actions = {
  setDetail({ commit }: any, data: any) {
    commit("setDetail", data);
  },
  setUpdate({ commit }: any, data: any) {
    commit("setUpdate", data);
  },
  setUpdateList({ commit }: any, data: any) {
    commit("setUpdateList", data);
  },
  setSearchParams({ commit }: any, data: any) {
    commit("setSearchParams", data);
  },
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations,
};
