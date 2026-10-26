//  单一状态树
const state = {
  current_input: -1, // 当前使用输入组件
  current_area: -1, // 当前使用区域
  drag_temp: null, // 拖拽临时元素
  component_catch: null, // 创建组件时的临时缓存
  component_first_catch: null, // 页面的第一个组件
  component_index: 0 // 组件 id
};

const getters = {
  // 查看 current_input
  getCurrentInput: state => {
    return state.current_input;
  },
  // 查看 current_input
  getCurrentArea: state => {
    return state.current_area;
  },
  // 查看 input_id
  getComponentId: state => {
    return state.component_index;
  },
  // 查看 drag_temp
  getDragTemp: state => {
    return state.drag_temp;
  },
  // 查看 component_catch
  getComponentCatch: state => {
    return state.component_catch;
  },
  // 查看 first component_catch
  getComponentFirstCatch: state => {
    return state.component_first_catch;
  }
};

const mutations = {
  // 设置当前显示的 input 输入组
  setCurrentInput: (state, data) => {
    state.current_input = data;
  },
  // 设置当前显示的 input 输入组
  setCurrentArea: (state, data) => {
    if (state.current_area !== data) {
      state.current_area = data;
      state.current_input = -1;
    }
  },
  // input id 自加
  componentIdAdd: state => {
    state.component_index++;
  },
  // 设置拖拽元素
  setDragTmpObj: (state, data) => {
    state.drag_temp = data;
  },
  // 设置 component 缓存
  setComponentCatch: (state, data) => {
    state.component_catch = data;
  },
  // 设置 component first 缓存
  setComponentFirstCatch: (state, data) => {
    state.component_first_catch = data;
  }
};

const actions = {
  // 设置当前显示的 input 输入组
  setCurrentInput: ({ commit }, data) => {
    commit("setCurrentInput", data);
  },
  // 设置当前显示的 input 输入组
  setCurrentArea: ({ commit }, data) => {
    commit("setCurrentArea", data);
  },
  // input id 自加
  componentIdAdd: ({ commit }) => {
    commit("componentIdAdd");
  },
  // 设置当前显示的 input 输入组
  setDragTmpObj: ({ commit }, data) => {
    commit("setDragTmpObj", data);
  },
  // 设置 component 缓存
  setComponentCatch: ({ commit }, data) => {
    commit("setComponentCatch", data);
  },
  // 设置 component first 缓存
  setComponentFirstCatch: ({ commit }, data) => {
    commit("setComponentFirstCatch", data);
  }
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
};
