import dayjs from "dayjs";

//  单一状态树
const state = {
  /** 筛选条件 */
  filterParams: {
    leave_member_ids: [],
    member_ids: [],
    third_ids: [],
    third_name_list: [], 
    member_name_list: [],
    date_type: 0,
    start_time: 0,
    end_time: 0
  },
  /** 是否更新 */
  update: false,
};

// getters 暴露出去的获取 state 的计算属性
const getters = {
  getFilterParams: (state: any) => {
    return state.filterParams;
  },
  getUpdate: (state: any) => {
    return state.update;
  },
};

// mutations 私有化方法，只能被 actions 调用， 可以直接变更状态
const mutations = {
  setFilterParams(state: any, data: any) {
    state.filterParams = data;

  },
  setUpdate(state: any, data: any) {
    state.update = data;
  },
};

// actions 暴露到外部的方法 不直接变更状态
const actions = {
  setFilterParams({ commit }: any, data: any) {
    commit("setFilterParams", data);
  },
  setUpdate({ commit }: any, data: any) {
    commit("setUpdate", data);
  },
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations,
};
