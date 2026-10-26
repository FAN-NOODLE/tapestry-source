/*
 * @Author: your name
 * @Date: 2020-10-13 17:45:46
 * @LastEditTime: 2020-11-06 14:29:48
 * @LastEditors: your name
 * @Description: In User Settings Edit
 * @FilePath: /crs-admin-base/src/store/modules/menu.ts
 */
import {
  VuexModule,
  Module,
  Action,
  Mutation,
  getModule
} from "vuex-module-decorators";
import store from "@/store";

interface IStatesInterface {
  MENU: any;
  ROOT: string;
  LEAF: string;
}

@Module({ name: "menu", dynamic: true, store })
class Menu extends VuexModule implements IStatesInterface {
  public MENU: any = null;
  public ROOT = "";
  public LEAF = "";
  get getMenu(): any {
    return this.MENU;
  }
  get getRoot(): string {
    let _root = this.ROOT;
    if (!this.ROOT) {
      const menu_root =
        sessionStorage.getItem("MENUROOT") || JSON.stringify("home");
      _root = JSON.parse(menu_root);
    }
    return _root;
  }
  get getLeaf(): string {
    let _leaf = this.LEAF;
    if (!this.LEAF) {
      const menu_leaf =
        sessionStorage.getItem("MENUROOT") || JSON.stringify("home");
      _leaf = JSON.parse(menu_leaf);
    }
    return _leaf;
  }

  @Mutation
  private SET_MENU(data: any) {
    this.MENU = data;
  }
  @Mutation
  private SET_ROOT(data: any) {
    sessionStorage.setItem("MENUROOT", JSON.stringify(data));
    this.ROOT = data;
  }
  @Mutation
  private SET_LEAF(data: any) {
    sessionStorage.setItem("MENULEAF", JSON.stringify(data));
    this.LEAF = data;
  }

  @Action({ rawError: true })
  public setMenu(params: any) {
    this.context.commit("SET_MENU", params);
  }
  @Action({ rawError: true })
  public setRoot(params: any) {
    this.context.commit("SET_ROOT", params);
  }
  @Action({ rawError: true })
  public setLeaf(params: any) {
    this.context.commit("SET_LEAF", params);
  }
}

export default getModule(Menu);
