import { VuexModule, Module, Action, Mutation, getModule } from "vuex-module-decorators";
import store from "@/store";

interface IStatesInterface {
  GROUPS: [];
  GROUPU: [];
  GROUPNAME: "";
  LABELEDIT: object;
}

@Module({ name: "customer", dynamic: true, store })
class Customer extends VuexModule implements IStatesInterface {
  public GROUPS: any = [];
  public GROUPU: any = [];
  public LABELEDIT: any = null;
  public GROUPNAME: any = "";
  public CUSTOMERGROUPS: any = null;
  get getGroupS(): any {
    return this.GROUPS;
  }
  get getGroupU(): any {
    return this.GROUPU;
  }
  get getGroupName(): any {
    return this.GROUPNAME;
  }

  get getLabelEdit(): any {
    return this.LABELEDIT || JSON.parse(sessionStorage.getItem("CUSTOMERLABELEDIT") || JSON.stringify({}));
  }
  get getCustomerGroups(): any {
    return this.CUSTOMERGROUPS || JSON.parse(sessionStorage.getItem("CUSTOMERGROUPS") || JSON.stringify({}));
  }
  @Mutation
  private SET_GROUP_S(data: any) {
    this.GROUPS = data;
  }
  @Mutation
  private SET_GROUP_U(data: any) {
    this.GROUPU = data;
  }
  @Mutation
  private SET_GROUP_NAME(data: any) {
    this.GROUPNAME = data;
  }
  @Mutation
  private SET_LABELEDIT(data: any) {
    sessionStorage.setItem("CUSTOMERLABELEDIT", JSON.stringify(data));
    this.LABELEDIT = data;
  }

  @Action({ rawError: true })
  public setGroupS(params: any) {
    this.context.commit("SET_GROUP_S", params);
  }
  @Action({ rawError: true })
  public setGroupU(params: any) {
    this.context.commit("SET_GROUP_U", params);
  }
  @Action({ rawError: true })
  public setGroupName(params: any) {
    this.context.commit("SET_GROUP_NAME", params);
  }
  @Action({ rawError: true })
  public setLabelEdit(params: any) {
    this.context.commit("SET_LABELEDIT", params);
  }
  @Action({ rawError: true })
  public setCustomerGroups(params: any) {
    this.context.commit("SET_CUSTOMERGROUPS", params);
  }
}

export default getModule(Customer);
