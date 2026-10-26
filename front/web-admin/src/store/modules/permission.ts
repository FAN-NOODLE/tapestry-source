import {
  VuexModule,
  Module,
  Action,
  Mutation,
  getModule
} from "vuex-module-decorators";
import store from "@/store";
import { IPermissionItem, IRoleItem } from "@/types";
interface StatesInterface {
  PERMISSION: Array<IPermissionItem>;
  ROLES: Array<IRoleItem>;
  EQ_NAME: string;
}

@Module({ name: "permission", dynamic: true, store })
class Permission extends VuexModule implements StatesInterface {
  public PERMISSION: Array<IPermissionItem> = [];
  public ROLES: Array<IRoleItem> = [];
  public EQ_NAME = "";
  /** 获取品牌名称 */
  get getEpName(): string {
    
    if(!this.EQ_NAME) {
      return localStorage.getItem("EQ_NAME")|| ''
    }
    
    return this.EQ_NAME;
  }
  /** 判断品牌 */
  get enterprise() {
    return {
      isCoach: this.getEpName === 'coach',
      isSW: this.getEpName === 'sw',
      isKS: this.getEpName === 'ks'
    };
  }
  get getPermission(): Array<IPermissionItem> {
    let _permission = this.PERMISSION;
    if (!this.PERMISSION || this.PERMISSION.length === 0) {
      const list = localStorage.getItem("PERMISSION") || "[]";
      _permission = JSON.parse(list);
    }
    return _permission;
  }
  get getRoleList(): Array<IRoleItem> {
    let _roles = this.ROLES;
    if (!this.ROLES || this.ROLES.length === 0) {
      const list = localStorage.getItem("ROLES") || "[]";
      _roles = JSON.parse(list);
    }
    return _roles;
  }

  @Mutation
  private SET_PERMISSION(data: Array<IPermissionItem>) {
    localStorage.setItem("PERMISSION", JSON.stringify(data));
    this.PERMISSION = data;
  }
  @Mutation
  private SET_ROLES(data: Array<IRoleItem>) {
    localStorage.setItem("ROLES", JSON.stringify(data));
    this.ROLES = data;
  }
 /** 设置品牌名称 */
  @Mutation
  private SET_EQ_NAME(name: string) {
    this.EQ_NAME = name;
  }
  @Action({ rawError: true })
  public setPermission(params: Array<IPermissionItem>) {
    this.context.commit("SET_PERMISSION", params);
  }
  @Action({ rawError: true })
  public setRoles(params: Array<IRoleItem>) {
    this.context.commit("SET_ROLES", params);
  }
  @Action({ rawError: true })
  public setEqName(name: string) {
    this.context.commit("SET_EQ_NAME", name);
    localStorage.setItem("EQ_NAME", name);
  }
}
export default getModule(Permission);
