<template>
  <!--分组弹窗 -->
  <van-popup
    v-model="show"
    @close="groupCancel"
    :close-on-click-outside="false"
    :close-on-click-action="false"
    round
    position="bottom"
    class="material-popup"
  >
    <div class="group-wrap">
      <van-sidebar v-model="activeSideKey" class="v-side-bar">
        <van-sidebar-item
          v-for="item in sidleListGroup"
          :key="item.id"
          :title="item.name"
          @click="recursionOtherGroupData(item)"
          class="v-side-bar-item"
        />
      </van-sidebar>
      <div class="right-tree">
        <cyTree :options="groupOption" v-if="groupOption.length > 0"></cyTree>
      </div>
    </div>
    <div class="btn-view">
      <div class="btn-li reset" @click="groupCancel">取消</div>
      <div
        class="btn-li confirm"
        :class="[nowChecked.id ? '' : 'empty']"
        @click="groupCheck"
      >
        确定
      </div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Vue, PropSync } from "vue-property-decorator";
import { Field, Loading, Popup, Sidebar, SidebarItem } from "vant";
import { getListForStore } from "@/api/marketing-promotion/index";
import "vs-tree/dist/vs-tree.css";
import { install } from "vs-tree";
import cyTree from "@/components/base/cy-tree/cy-tree.vue";

Vue.use(install);
Vue.use(Sidebar);
Vue.use(SidebarItem);
/** 子分组 */
interface Igr {
  name: string;
  id: number;
  children: null | Igr[];
  source_lvl: number;
}
/** 主分组 */
interface IGroup {
  created_at: number;
  deleted_at: number;
  enterprise_id: number;
  group_name: string;
  id: number;
  is_show: number;
  orders: number;
  parent_id: number;
  updated_at: number;
  children: Array<Igr>;
}

@Component({
  components: {
    cyTree,
    [Loading.name]: Loading,
    [Popup.name]: Popup,
    [Field.name]: Field,
  },
})
export default class MaterialGroup extends Vue {
  @PropSync("showPopup", { type: Boolean }) show;
  /** 当前点击的分组 */
  private get nowChecked() {
    return this.$store.getters["materialLibrary/getNowChecked"];
  }
  /**当前popup侧边栏激活项 */
  private activeSideKey = -1;
  /** 所有分组 */
  private sidleListGroup: Array<Igr> = [
    { name: "未分组", id: -1, children: null, source_lvl: 1 },
  ];
  /**分组下的树结构子分组 */
  private groupOption: Array<Igr> = [];
  private stairGroupName = ""; //左侧一级名称
  private headerGroupName = ""; //当前选中的所有分组名称
  private groupVal = 0; // 所选的子分组

  /** 获取分组以及分组树的选项 */
  private async loadGroupOptions() {
    const { code, data = [] }: IApiResult<Array<IGroup>> =
      await getListForStore();
    if (code !== 200) {
      return;
    }
    let sidleList = this.recursionGroupData(data, 0);
    this.sidleListGroup = [...this.sidleListGroup, ...sidleList];
    console.log("this.sli :>> ", this.sidleListGroup);
  }
  /**重组子分组下数据(添加全部) */
  private recursionGroupData(data: Array<IGroup | Igr>, level: number) {
    // sidleList左侧一级分组数据
    let sidleList = [];
    data.forEach((item: IGroup) => {
      let newChildren = this.recursionGroupData(item.children || [], level + 1);
      if (!level) {
        newChildren.unshift({
          name: "全部",
          id: item.id,
          children: null,
          source_lvl: level + 1,
          is_show: item.is_show,
        });
      }

      sidleList.push({
        ...item,
        name: item.group_name,
        children: newChildren,
        source_lvl: level + 1,
      });
    });
    return sidleList;
  }
  /**点击左侧分组，右侧tree重组数据赋值 */
  private recursionOtherGroupData(option: Igr) {
    this.activeSideKey = option.id;
    this.stairGroupName = option.name;
    this.groupOption = [];
    //点击的分组为 全部分组 未分组
    if (option.id == 0 || option.id == -1) {
      this.groupVal = option.id; //子分组为全部
      this.headerGroupName = this.screenGroupName(option);
      this.$store.commit("materialLibrary/setNowChecked", {
        id: 0,
        groupName: this.headerGroupName,
        source_lvl: 1,
      });

      this.$emit("confirm", {
        group: this.groupVal,
        groupName: this.headerGroupName,
      });
    } else {
      setTimeout(() => {
        if (option.children == null) {
          this.groupOption = [];
        } else {
          this.groupOption = option.children;
        }
      });
      // 选择左侧，重置选择的node
      this.$store.commit("materialLibrary/setNowChecked", {});
    }
  }
  /** 拼各级组名 */
  private screenGroupName(data: Igr) {
    if (data.source_lvl === 1) {
      //左侧主分组
      return this.stairGroupName;
    } else if (data.source_lvl === 2) {
      //右侧子分组
      return `${this.stairGroupName}>${data.name}`;
    } else if (data.source_lvl === 3) {
      //子分组的子集名称
      let secondGroupName = "";
      this.groupOption.forEach((itemTwo: Igr) => {
        if (itemTwo.children != null) {
          itemTwo.children.forEach((itemThree: Igr) => {
            if (itemThree.id == this.groupVal) {
              secondGroupName = itemTwo.name;
              return;
            }
          });
        }
      });
      return `${this.stairGroupName}>${secondGroupName}>${data.name}`;
    } else if (data.source_lvl === 4) {
      //子分组的子集名称
      let secondGroupName = "";
      let threeGroupName = "";
      this.groupOption.forEach((itemTwo: Igr) => {
        if (itemTwo.children != null) {
          itemTwo.children.forEach((itemThree: Igr) => {
            if (itemThree.children != null) {
              itemThree.children.forEach((itemFour: Igr) => {
                if (itemFour.id == this.groupVal) {
                  threeGroupName = itemThree.name;
                  return;
                }
              });
              secondGroupName = itemTwo.name;
              return;
            }
          });
        }
      });
      return `${this.stairGroupName}>${secondGroupName}>${threeGroupName}>${data.name}`;
    }
  }
  /** 分组筛选取消 */
  private groupCancel() {
    this.$emit("cancel");
    this.handleReset()
  }
  /** 分组筛选确定 */
  private groupCheck() {
    const node = this.$store.getters["materialLibrary/getNowChecked"];
    console.log("node :>> ", node);
    if (Object.keys(node).length < 1) return;
    this.getClick(node);
  }
  /**点击tree */
  private getClick(node: any) {
    this.groupVal = node.id;
    this.headerGroupName = this.screenGroupName(node);

    this.$emit("confirm", {
      group: this.groupVal,
      groupName: this.headerGroupName,
    });
    this.handleReset();
  }
  /**清空已选择 */
  private handleReset() {
    /**当前popup侧边栏激活项 */
    this.activeSideKey = -1;
    this.groupOption = [];
    this.stairGroupName = ""; //左侧一级名称
    this.headerGroupName = ""; //当前选中的所有分组名称
    // 选择左侧，重置选择的node
    this.$store.commit("materialLibrary/setNowChecked", {});
    this.groupVal = 0; // 所选的子分组
  }
  /**生命周期 */
  private mounted() {
    this.loadGroupOptions();
  }
}
</script>
<style lang="scss" scoped>
::-webkit-scrollbar {
  display: none; /* Chrome Safari */
}
.material-popup {
  height: 453px;
  padding: 0;
  display: flex;
  justify-content: space-between;
  flex-direction: column;
  border-radius: 0;
  .group-wrap {
    flex-grow: 1;
    display: flex;
    overflow: hidden;
    padding-top: 4px;
  }
  .btn-view {
    flex-shrink: 1;
    width: 100%;
    display: flex;
    background-color: #ffffff;
    padding: 15px;
    .btn-li {
      flex: 2;
      border: 1px solid #4a4a4a;
      color: #4a4a4a;
      font-size: 14px;
      text-align: center;
      line-height: 36px;
    }
    .confirm {
      margin-left: 15px;
      color: #ffffff;
      background-color: #4a4a4a;
    }
    .empty {
      color: #ffffff;
      background-color: #cccccc;
      border-color: #cccccc;
    }
  }
  .right-tree {
    height: 100%;
    overflow-y: auto;
    width: 100%;
  }
  .v-side-bar {
    width: 112px;
    flex-shrink: 0;
    background: #fff;
    overflow-y: auto;
    height: 100%;
    border-right: 2px solid #eeeeee;
    .v-side-bar-item {
      height: 60px;
      display: flex;
      align-items: center;
      background: #fff;
      .van-sidebar-item__text {
        letter-spacing: 0;
        overflow: hidden;
        display: -webkit-box;
        text-overflow: ellipsis;
        -webkit-line-clamp: 2; /*要显示的行数*/
        -webkit-box-orient: vertical;
        font-size: 14px;
      }
    }
    .van-sidebar-item--select {
      font-weight: 600;
      color: #000000;
      border-radius: 0;
      background: rgba(244, 247, 251, 1);
    }
    .van-sidebar-item--select::before {
      display: none;
    }
  }
  .vs-tree {
    width: calc(100% - 112px);
    max-height: 100% !important;
    padding: 0px 13px;
    ::v-deep.vs-expand-icon::after {
      width: 16px;
      height: 16px;
      background-size: 100% 100%;
      background-image: url("../../../assets/material-library/show-icon.png");
      margin-top: -2px;
      margin-right: 3px;
    }
    ::v-deep.vs-tree-text {
      margin-left: 3px;
      font-size: 14px;
      color: #777777;
      line-height: 20px;
    }
    ::v-deep.vs-tree-node {
      height: 40px;
      color: #777777;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
    ::v-deep.vs-tree-node:hover {
      height: 40px;
    }
    ::v-deep.selected {
      height: 40px;
      background: rgba(245, 244, 247, 1);
      font-weight: 900;
      padding-left: 5px;
    }
    ::v-deep.expand-empty {
      display: none;
    }
    ::v-deep.selected .vs-tree-text {
      color: #666668;
    }
  }
}
</style>
