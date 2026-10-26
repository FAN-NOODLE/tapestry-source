/*
 * @Author: your name
 * @Date: 2020-11-03 16:21:34
 * @LastEditTime: 2020-11-05 19:33:19
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /crs-admin/src/shared/index.ts
 */
import Help from "@/shared/model/help";
class SharedModule {
  static help = new Help();

  /**
   * 重载 shared
   */
  static overloadShared(help: Help) {
    SharedModule.help = help;
  }

  /**
   * 获取 shared 实例
   */
  static getShared() {
    return SharedModule.help;
  }
}

export default SharedModule;
