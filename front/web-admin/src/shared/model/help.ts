/*
 * @Author: your name
 * @Date: 2020-11-03 16:21:34
 * @LastEditTime: 2020-11-10 10:16:27
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /crs-admin/src/shared/model/help.ts
 */
class Help {
  private status = false;

  /**
   * 获取帮助中心状态
   */
  getHelp() {
    // 子应用独立运行时，在 localStorage 中获取 token
    // return this.status;
  }
  /**
   * 获取帮助中心状态
   */
  pushState(name: string, query?: any) {
    // 子应用独立运行时，在 localStorage 中获取 token
    console.log(name);
    // return this.status;
  }
}
export default Help;
