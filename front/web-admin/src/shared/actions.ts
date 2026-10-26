function emptyAction(callback: any, fireImmediately?: boolean) {
  // 警告：提示当前使用的是空 Action
  console.warn("Current execute action is empty!", callback, fireImmediately);
}
function emptyActionSet(state: Record<string, any>) {
  // 警告：提示当前使用的是空 Action
  console.warn("Current execute action is empty!", state);
}
class Actions {
  // 默认值为空 Action
  actions = {
    onGlobalStateChange: (callback: any, fireImmediately?: boolean) => {
      emptyAction(callback, fireImmediately);
    },
    setGlobalState: (state: Record<string, any>) => {
      emptyActionSet(state);
    }
  };

  /**
   * 设置 actions
   */
  setActions(action: any) {
    this.actions = action;
  }

  /**
   * 映射
   */
  onGlobalStateChange(callback: any, fireImmediately?: boolean) {
    return this.actions.onGlobalStateChange(callback, fireImmediately);
  }

  /**
   * 映射
   */
  setGlobalState(state: Record<string, any>) {
    return this.actions.setGlobalState(state);
  }
}

const actions = new Actions();
export default actions;
