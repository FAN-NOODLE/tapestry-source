/*
 * @Author: your name
 * @Date: 2020-10-13 17:46:25
 * @LastEditTime: 2020-11-06 14:30:10
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /crs-admin-base/src/store/modules/user.ts
 */
import {
  VuexModule,
  Module,
  Action,
  Mutation,
  getModule
} from "vuex-module-decorators";
import { IUserInfo } from "@/types";
import store from "@/store";

interface IStatesInterface {
  USERINFO: IUserInfo;
}
@Module({ name: "user", dynamic: true, store })
class User extends VuexModule implements IStatesInterface {
  public USERINFO: any = "";
  get getUserInfo(): any {
    let _root = this.USERINFO;
    if (!this.USERINFO) {
      const user_info = localStorage.getItem("USERINFO") || JSON.stringify({});
      _root = JSON.parse(user_info);
    }
    return _root;
  }

  @Mutation
  private SET_USERINFO(data: IUserInfo) {
    localStorage.setItem("USERINFO", JSON.stringify(data));
    this.USERINFO = data;
  }

  @Action({ rawError: true })
  public setUserInfo(params: IUserInfo) {
    this.context.commit("SET_USERINFO", params);
  }
}

export default getModule(User);
