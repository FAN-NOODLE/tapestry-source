import {
  VuexModule,
  Module,
  Action,
  Mutation,
  getModule
} from "vuex-module-decorators";
import store from "@/store";

interface IStatesInterface {
  [propName: string]: any;
}
interface IObj {
  data: any;
  type: any;
}

@Module({ name: "app", dynamic: true, namespaced: true, store })
class App extends VuexModule implements IStatesInterface {
  public BASE_MOD_KEYS = new Set();

  get getAppKey(): any {
    return this.BASE_MOD_KEYS;
  }

  @Mutation
  private SET_KEYS(val: string | Array<any>) {
    if (typeof val === "string") {
      this.BASE_MOD_KEYS.add(val);
    } else {
      this.BASE_MOD_KEYS = new Set([...this.BASE_MOD_KEYS, ...val]);
    }
    // this.BASE_MOD_KEYS.add(val);
  }

  @Action({ rawError: true })
  public setKeys(data: any) {
    this.context.commit("SET_KEYS", data);
  }
}

export default getModule(App);
