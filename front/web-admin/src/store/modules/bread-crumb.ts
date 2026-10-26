import {
  VuexModule,
  Module,
  Action,
  Mutation,
  getModule
} from "vuex-module-decorators";
import store from "@/store";
import { IBCNodeItem } from "@/types";
import { getBCNodes } from "@/utils/keys-map";
import { Route } from "vue-router";
import MenuModule from "@/store/modules/menu";
import HelpModule from "@/store/modules/help";

interface IStatesInterface {
  BREADCRUMB: Array<IBCNodeItem>;
}
interface IParams {
  to: Route;
  nodeItem: IBCNodeItem;
}
@Module({ name: "bread-crumb", dynamic: true, store })
class BreadCrumb extends VuexModule implements IStatesInterface {
  private keys = JSON.parse(localStorage.getItem("KEYS") || "{}"); //权限集合
  public BREADCRUMB: Array<IBCNodeItem> = [];
  get getBreadCrumb(): Array<IBCNodeItem> {
    return this.BREADCRUMB;
  }

  /*
   *  面包屑展示规则
   *  1. 一条数据：直接展示
   *  2. 两条数据：只展示第二条
   *  3. 三条以上: 第三条开始逐级展示
   *
   *  type : 1 菜单
   *  type : 2 视图
   *  type : 3 按钮
   *  type : 4 菜单+视图
   *
   *  1. 一条数据：直接展示：一级
   *  2. 两条数据：
   *      1. 当最后一条数据type===4: 展示 二级
   *      2. 当最后一条数据type===2: 展示 一级+二级（非第一条数据）
   *  3. 三条数据：
   *      1. 当最后一条数据type===4: 展示 三级
   *      2. 当最后一条数据type===2: 展示 二级+三级（非第一条数据）
   *  4. 四条数据：
   *      1. 展示 三级+四级（非第一条数据）
   *
   * 最后一个子节点type===4 展示：当前级别
   * 最后一个子节点type===2 展示：父级+子级（非第一条数据）
   *
   */

  static getBCShowRule = (list: Array<IBCNodeItem>) => {
    MenuModule.setRoot(list[0].key);
    const lastNode = list[list.length - 1]; // 获取最后一个节点，即当前页面节点
    let bcList: Array<IBCNodeItem> = [];
    if (list.length === 1) {
      bcList = list;
    } else if (lastNode.type === 4) {
      // 最后一个子节点type===4 展示：当前级别
      bcList.push(lastNode);
    } else {
      //最后一个子节点type===2 展示：父级+子级（非第一条数据）
      bcList = list.slice(-2); // 倒数第二个元素开始截取
      if (lastNode.key === bcList[0].children[0].key) {
        // 子级（非第一条数据）
        bcList.pop(); // 删除最后一条
      }
    }
    MenuModule.setLeaf(bcList[0].key);
    return bcList;
  };

  @Mutation
  private SET_BREADCRUMB(params: IParams) {
    const {
      to: { query },
      nodeItem
    } = params;
    if (this.keys !== "{}") {
      const list = getBCNodes(this.keys, nodeItem) as Array<IBCNodeItem>;
      this.BREADCRUMB = list.length > 0 ? BreadCrumb.getBCShowRule(list) : [];
      if (query.bread)
        this.BREADCRUMB[
          this.BREADCRUMB.length - 1
        ].name = query.bread as string;
    }
    HelpModule.setHelp(JSON.parse(localStorage.getItem("HELP") || "1"));
  }
  @Mutation
  private CLEAR_BREADCRUMB() {
    this.BREADCRUMB = [];
  }

  @Action({ rawError: true })
  public setBreadCrumb(params: IParams) {
    this.context.commit("SET_BREADCRUMB", params);
  }
  @Action({ rawError: true })
  public clearBreadCrumb() {
    this.context.commit("CLEAR_BREADCRUMB");
  }
}

export default getModule(BreadCrumb);
