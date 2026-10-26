import {
  VuexModule,
  Module,
  Action,
  Mutation,
  getModule
} from "vuex-module-decorators";
import store from "@/store";
import actions from "@/shared/actions";
import { helpIsShow } from "@/utils/help-data";
interface IHCListItem {
  title: string; //标题
  released_at: number; //时间戳
  content: string; //公告内容
}
interface IHCItem {
  key: string;
  type: string;
  content: string | Array<IHCListItem>;
}
interface IStatesInterface {
  HELPSHOW: boolean;
  HELPCONTENT: IHCItem | string;
}

@Module({ name: "help", dynamic: true, store })
class Help extends VuexModule implements IStatesInterface {
  public HELPSHOW = true;
  public HELPCONTENT: IHCItem = { key: "key", type: "text", content: "默认" };

  get getHelp(): boolean {
    return this.HELPSHOW;
  }
  get getHelpContent(): IHCItem | string {
    return this.HELPCONTENT;
  }

  /**
   * @param status 1展开 2:隐藏
   **/

  @Mutation
  private SET_HELP(config: any) {
    // localStorage.setItem("HELP", status);
    // const config = helpIsShow();
    // actions.setGlobalState({ config:  });
    this.HELPSHOW = config.isShow;
    // this.HELPSHOW = false;
    this.HELPCONTENT = config.config;
  }

  @Action({ rawError: true })
  public async setHelp(status: string) {
    localStorage.setItem("HELP", status);
    const config = await helpIsShow();
    actions.setGlobalState({ config });
    this.context.commit("SET_HELP", config);
  }
}

export default getModule(Help);
