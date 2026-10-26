<template>
  <el-dialog
    :title="title"
    :visible.sync="syncVisible"
    width="1000px"
    :class="['select-person', appendToBody && 'body-select-person']"
    :append-to-body="appendToBody"
    :destroy-on-close="true"
    :close-on-click-modal="false"
  >
    <el-row class="content">
      <!-- 左侧员工树 -->
      <el-col :span="14" class="left-content">
        <div class="box-left p-r-5">
          <!-- 搜索框模块-->
          <search-item>
            <el-input
              :placeholder="placeholder"
              v-model.trim="searchInput"
              style="width: 300px;"
              clearable
              size="small"
            >
              <el-select
                v-model="selectType"
                slot="prepend"
                style="width: 120px;"
              >
                <el-option
                  v-for="item in selectDefault"
                  :label="item.label"
                  :value="item.value"
                  :key="item.vlaue"
                ></el-option>
              </el-select>
            </el-input>
          </search-item>

          <div class="search-btn m-l-10">
            <el-button size="small" @click="handleSearch" type="primary"
              >查询</el-button
            >
            <el-button size="small" class="m-l-10" @click="handleReset"
              >重置</el-button
            >
          </div>
        </div>

        <div class="store-sort">
          <div class="sort-title">门店分类：</div>
          <div>
            <el-checkbox-group v-model="storeType" @change="storeTypeChange">
              <el-checkbox :label="1">零售</el-checkbox>
              <el-checkbox :label="2">奥莱</el-checkbox>
            </el-checkbox-group>
          </div>
        </div>

        <div class="tree-box" v-if="syncVisible">
          <div class="scroll-tree">
            <el-tree
              :key="temp"
              v-show="!searchInput"
              :props="defaultProps"
              :expand-on-click-node="false"
              node-key="id"
              ref="tree"
              lazy
              check-strictly
              :load="loadNode"
              :default-expanded-keys="treeApi.expand_key"
              @node-click="handleNodeClick"
            >
              <person slot-scope="{ data, node }" :data="data" :node="node" />
            </el-tree>

            <!-- 搜索结果页 -->
            <div
              v-show="searchInput"
              v-for="(item, index) in searchList"
              :key="index"
              class="search-item"
              @click="selectSearchItem(item)"
            >
              <person :data="item" :node="item" />
            </div>
          </div>
        </div>
      </el-col>

      <!-- 右侧已选员工 -->
      <el-col :span="10" class="right-content">
        <!-- 提示文案 -->
        <div class="right-title">
          已选择的部门或{{level===1?'门店':'成员'}}
        </div>
        <div class="right-tip">
          （保存时 会自动合并已选部门下{{level===1?'门店':'成员'}}）
        </div>

        <!-- 右侧列表 -->
        <div class="right-persons">
          <div
            v-for="(item, index) in selectMember"
            :key="index"
            class="person-item"
          >
            <person
              :data="{ ...item, mem_subscribe_status: true }"
              :del="true"
              @handleClick="seletedListApi.delNode(item, index)"
            />
          </div>
        </div>
      </el-col>
    </el-row>
    <span slot="footer" class="dialog-footer">
      <el-button size="small" @click="cancel">取 消</el-button>
      <el-button size="small" type="primary" @click="submit">确 定</el-button>
    </span>
  </el-dialog>
</template>

<script lang="ts">
import { Component, Vue, PropSync, Watch, Prop } from "vue-property-decorator";
import SelectInput from "@/components/business/select-input/index.vue";
import { flatten, cloneDeep } from "lodash";
import {
  Member,
  DpParams,
  DpType,
  PersonType,
  PersonParams,
  treeApiType,
  guideNameListType
} from "./type";
import Person from "./person.vue";
import API from "@/api/select-member";
import { ApiResult } from "@/types";

// 搜索类型
const SEARCHTYPE = [
  { key: 3, name: "员工ID" },
  { key: 4, name: "员工手机号" },
  { key: 1, name: "员工姓名" },
  { key: 2, name: "部门名称" },
];

// 群活码搜索条件类型
type SearchType = {
  codeType: typeof SEARCHTYPE;
  value: string;
  checkedkey: number;
  placeholder: string;
  show: boolean; // 判断是否显示搜索内容
};

// tree 配置项
const DETAULTPROPS = {
  isLeaf: "leaf",
  children: "children",
  label: "name"
};

/**
 * 选择员工组件
 */
@Component({
  components: { Person, SelectInput }
})
export default class SelectPerson extends Vue {
  // 判断父属性显示影藏，双向绑定
  @PropSync("visible", { type: Boolean }) syncVisible!: boolean;

  // 默认选中数据
  @Prop({ default: () => [], type: Array }) selectedData?: Array<
    guideNameListType
  >;

  // 选择角色：0-导购；1-门店
  @Prop({default: 0, type: Number }) level: number;

  // 是否添加到body下
  @Prop({ default: false, type: Boolean }) appendToBody?: boolean;

  // 自定义标题
  @Prop({ type: String, default: "添加员工（支持多选）" }) title?: string;
  private ocList = []; //零售
  private ocfList = []; //奥莱
  temp = Date.now()
  // 搜索类型
  private selectDefault = [
    { label: "员工ID", value: "userid" },
    { label: "员工手机号", value: "mobile" },
    { label: "员工姓名", value: "mem_username" },
    { label: "部门名称", value: "department" },
  ];
  private selectType = "userid"; // 搜索类型默认选中
  private searchInput = ""; // 搜索内容

  get placeholder() {
    const label = this.selectDefault.find((item) => item.value === this.selectType).label
    return `请输入${label}`;
  }
  // private mounted() {}
  @Watch("selectType")
  private onChangeSelectType() {
    // 监听清空搜索数据
    this.searchInput = "";
    this.searchList = [];
  }

  @Watch("visible")
  changeVisible(visible: boolean) {
    if (visible) {
      // 初始化数据
      this.treeApi.expand_key = [];
      const selectedData = (this.selectedData as Array<guideNameListType>).map(
        (item) => {
          console.log(item);
          return {
            ...item,
            user_id: item.user_id || item.userid
          };
        }
      );

      this.selectMember = cloneDeep(selectedData);
      console.log("SSSSelectMember: ", this.selectMember, this.selectedData);
      this.searchList = [];
      this.handleReset();
    }
  }

  private mounted() {
    this.getStoreDataList();
  }

  /** 获取奥莱 零售门店数据
   *  1. 勾选奥莱  将 --selectMember-- 赋值为 --奥莱所有数据--； 零售同理
   *  2. 加一个或者减少一个 奥莱/零售以外的数据，奥莱/零售 取消勾选; *****
   *  3. 反选： selectedData 是否包含  oc_list / ocf_list 中全部数据 若包含-> 则勾选奥莱/零售 反之 不勾选
   * */
  private getStoreDataList() {
    // oc_list:零售 / ocf_list:奥莱
    API.getStoreDataList().then((res) => {
      if (res.data) {
        this.ocList = res.data.list?.oc_list;
        this.ocfList = res.data.list?.ocf_list;
      }
    });
  }

  private compareStoreData() {
    // 返回true/false
  }

  // 选中员工的ids
  get selectMemberIds() {
    const idsArray: Array<Array<number | string>> = [[], []];

    const [dpIds, personIds] = idsArray;

    this.selectMember.forEach((item) => {
      if (item.is_dp) {
        dpIds.push(item.id);
      } else {
        personIds.push(item.id);
      }
    });
    return idsArray;
  }

  // tree 配置项
  private defaultProps = DETAULTPROPS;

  // 搜索条件
  get searchInfo() {
    const isStore = this.level===1;
    const options = SEARCHTYPE.filter((item) => isStore && item.key === 2 || !isStore)
    const checked = isStore ? 2 : 3
    const placeholder = options.find((item) => item.key === checked)?.name
    return {
      show: false,
      codeType: options,
      value: "",
      checkedkey: checked,
      placeholder: `请输入${placeholder}`
    }
  }

  // 选中员工集合
  private selectMember: Array<Member> = [];

  // 搜索返回值
  private searchList: Array<Member> = [];

  // 门店类型
  private storeType = [];

  // 树操作
  private treeApi: treeApiType = {
    expand_key: [],
    setNode: this.setChoice,
    append: this.appendTree
  };

  // 选中列表
  private seletedListApi = {
    delNode: this.removeSelect,
    addNode: this.addSelect
  };

  /**
   * 搜索
   */
  private async handleSearch() {
    // 扁平化ids
    const ids = flatten(this.selectMemberIds);

    if (this.selectType === "department") {
      // 搜索部门
      this.getDp({ dp_name: this.searchInput }).then((res) => {
        this.searchList = res.map((item: Member) => ({
          ...item,
          checked: ids.indexOf(item.id) !== -1
        }));
      });
    } else {
      // 搜索员工
      const list = await this.getPeople(-1, {});
      this.searchList = list.map((item: PersonType) => ({
        is_dp: false,
        name: item.user_name,
        id: item.member_id,
        user_id: item.user_id,
        leaf: true,
        mem_subscribe_status: item.mem_subscribe_status === 1,
        checked: ids.indexOf(item.member_id) !== -1
      }));
    }
  }

  /**
   * 搜索重置
   */
  private handleReset() {
    this.$nextTick(() => {
      this.searchInput = "";
      this.searchList = [];
      this.selectType = 'userid';
      this.temp = Date.now();
    });
  }

  /**
   * 点击搜索内容
   */
  private selectSearchItem(item: Member) {
    this.handleNodeClick(item, item);
  }

  /**
   * 右侧点击移除
   * @param { Member } 右侧员工
   * @param { Number } 下标
   */
  private removeSelect(data: Member, index?: number) {
    // 删除右侧item
    if (!index) {
      index = this.selectMember.findIndex((item) => item.id === data.id);
    }
    //如果删除的节点是奥莱2或者零售1
    const finxIndexLs = this.storeType.findIndex((item) => item == 1);
    if (data.store_type == 1 && finxIndexLs != -1) {
      this.storeType.splice(finxIndexLs, 1);
    }
    const finxIndexOl = this.storeType.findIndex((item) => item == 2);
    if (data.store_type == 2 && finxIndexOl != -1) {
      this.storeType.splice(finxIndexOl, 1);
    }

    this.selectMember.splice(index, 1);

    // 处理选中标示
    if (this.searchInfo.show) {
      // 处理搜索列表
      this.searchList.forEach((item) => {
        if (item.id === data.id) {
          item.checked = false;
        }
      });
    } else {
      // 处理tree列表
      this.treeApi.setNode();
    }
  }

  /**
   * 右侧添加选中
   */
  private addSelect(data: Member) {
    this.selectMember.push(data);
    this.treeApi.setNode();
  }

  // 取消 奥莱/零售 选择
  private cancelChecked(type) {
    if (this.storeType) {
      const data = this.storeType.filter((item) => item !== type);
      this.storeType = data;
      console.log(this.storeType);
    }
  }
  // 整理出 已选数据/奥莱/零售的 ID集合
  private idDataList(list) {
    const data = [];
    list?.forEach((item) => {
      data.push(item.id);
    });
    return data;
  }
  //手动勾选tree数据时，判断是否选择了 奥莱/零售 的全部数据 若全部选择  则自动勾选 奥莱/零售复选框
  private setOcOcfChecked(list, selectMember) {
    const listData = this.idDataList(list);
    const selectMemberData = this.idDataList(selectMember);

    console.log(listData.every((item) => selectMemberData.includes(item)));
    console.log([...listData], [...selectMemberData]);
    if (selectMember) {
      return listData.every((item) => selectMemberData.includes(item));
    } else {
      return true;
    }
  }

  /**
   * 点击tree回调
   */
  private handleNodeClick(data: Member, node: any) {
    console.log("click node -> ", data, node);
    if (!data.mem_subscribe_status)
      return this.$message.warning(`员工${data.name}还未加入企业, 请重新选择`);

    // 取消选中
    if (node.checked) {
      node.checked = false;
      this.seletedListApi.delNode(data);

      // 使用该data在 **this.ocList-奥莱 / this.ocfList-零售** 中做筛选，如果包含，则取消选中
      // 1 - ocList:零售 / 2 - ocfList:奥莱
      if (this.ocList.findIndex((item) => item.id === data.id) !== -1) {
        // 零售取消勾选
        this.cancelChecked(1);
      }
      if (this.ocfList.findIndex((item) => item.id === data.id) !== -1) {
        // 奥莱取消勾选
        this.cancelChecked(2);
      }

      console.log("CCCancel: ", data);
    } else {
      // 设置选中
      node.checked = true;
      this.seletedListApi.addNode(data);

      // this.selectMember 对比 this.ocList/this.ocfList， 若为包含关系 则设置 奥莱/零售勾选
      const typeList = [...this.storeType];
      if (this.setOcOcfChecked(this.ocList, this.selectMember)) {
        if (!typeList.includes(1)) {
          typeList.push(1);
        }
      }
      if (this.setOcOcfChecked(this.ocfList, this.selectMember)) {
        if (!typeList.includes(2)) {
          typeList.push(2);
        }
      }
      console.log([...typeList]);
      this.storeType = typeList;
    }
  }

  /**
   * 获取部门
   * @returns { Member } 节点信息
   */
  private async getDp(params: DpParams) {
    try {
      const res: ApiResult = await API.depListApi(params);
      return res.data.map((item: DpType) => {
        // 重置字段
        return {
          is_dp: true,
          name: item.dp_name,
          id: item.dp_id,
          leaf: false,
          mem_subscribe_status: true,
          third_id: item.third_id,
          store_type: item.store_type
        };
      });
    } catch (error) {
      console.error(error);
    }
  }

  /**
   * 根据id获取部门人员
   */
  private async getPeople(id: string | number, node: any) {
    try {
      // 组装请求参数
      const params: any = {
        page: 1,
        size: 99999,
        search: { 
          dp_id: id, 
          [this.selectType]: this.searchInput
        }
      };

      // 添加搜索情况
      if (id === -1) {
        delete params.search.dp_id;
      }
      const res: ApiResult = await API.treePeople(params);

      const { list } = res.data;

      // 搜索处理
      if (id === -1) return list;

      // 循环插入
      list.forEach((item: PersonType) => {
        const newItem = {
          is_dp: false,
          name: item.user_name,
          id: item.member_id,
          user_id: item.user_id,
          leaf: true,
          mem_subscribe_status: item.mem_subscribe_status === 1
        };
        this.treeApi.append(newItem, node);
      });
    } catch (e) {
      console.error(e);
    }
  }

  /**
   * tree load事件
   */
  private async loadNode(
    node: { level: number; data: Member },
    resolve: Function
  ) {
    // 根节点
    if (node.level === 0) {
      const dp_arr = await this.getDp({ is_top: 1 });

      // 默认展开第一级菜单
      if (dp_arr && dp_arr.length > 0)
        this.treeApi.expand_key.push(dp_arr[0].id);

      // 返回
      resolve(dp_arr);
    } else if (node.level >= 1) {
      // 分节点
      const params = { third_id: node.data.third_id };
      const list = await this.getDp(params);
      resolve(list);
      // 获取人员
      if(this.level === 0) {
        this.getPeople(node.data.id, node);
      }
    }
    // 设置选中
    this.treeApi.setNode();
  }

  /**
   * 添加到节点中
   */
  private appendTree(node: Member, target: Member) {
    (this.$refs.tree as any)["append"](node, target);
  }

  /**
   * 设置树选中
   */
  private setChoice() {
    // 扁平化数组
    const key_arr = flatten(this.selectMemberIds);
    (this.$refs.tree as HTMLFormElement).setCheckedKeys(key_arr);
  }

  /**
   * 确认, 点击提交
   * @emit { callBackType }
   */
  private async submit() {
    // 处理选中为空的情况
    if (!this.selectMember.length) {
      this.syncVisible = false;
      return this.$emit("submit", { guide_name_list: [], rawData: [] });
    }

    this.$emit("submit", {
      rawData: this.selectMember
    });
    this.syncVisible = false;
  }

  /**
   * 奥莱 /零售 门店勾选
   */
  private storeTypeChange(val) {
    console.log(132123,val);
    
    // 1 - oc_list:零售 / 2 - ocf_list:奥莱
    const list = [...val];
    let selectedList = [];
    if (list.findIndex((item) => item === 1) !== -1) {
      selectedList = selectedList.concat(this.ocList);
    }
    if (list.findIndex((item) => item === 2) !== -1) {
      selectedList = selectedList.concat(this.ocfList);
    }

    // 先从已选数据里删除已有数据
    // 获取id 集合过滤当前数据
    const oldListIds = [...this.ocList, ...this.ocfList].map((item) => item.id);
    const selectMemberTemp = this.selectMember.filter((item) => {
      return (item.is_dp && !oldListIds.includes(item.id)) || !item.is_dp;
    });
    // 把奥莱/零售添加到现有selectMember中
    this.selectMember = [...selectMemberTemp, ...selectedList];
    this.treeApi.setNode(); //反选树
    console.log("selectedList: ", selectedList, this.selectMember);
  }

  /**
   * 取消选中
   */
  private cancel() {
    // 清除已选内容
    this.selectMember = [];
    this.syncVisible = false;
  }
}
</script>

<style lang="scss" scoped>
.select-person {
  /deep/.el-tree-node {
    margin: 2px 0;
  }
  .content {
    min-height: 200px;
    .left-content {
      padding-right: 15px;
      .search-btn {
        display: inline-block;
      }
      .search-item {
        cursor: pointer;
      }
      .tree-box {
        height: 402px;
        overflow: hidden;
        overflow-y: auto;
      }
      .scroll-tree {
        width: 100%;
        min-height: 100%;
      }
    }
    .right-content {
      border-left: 1px solid $color-third-line;
      padding-left: 10px;
      .right-tip {
        color: red;
        font-size: 12px;
        margin-bottom: 5px;
      }
      .right-persons {
        height: 418px;
        overflow: hidden;
        overflow-y: auto;
        .person-item {
          padding: 0 5px;
        }
      }
    }
  }
  .buttons {
    margin-bottom: 15px;
  }
}

// 解决将el-dialog 添加到body 中, 全局样式丢失的问题
.body-select-person {
  ::v-deep {
    .el-dialog {
      border-radius: 10px;
      -webkit-border-radius: 10px !important;
      -moz-border-radius: 10px !important;
      -ms-border-radius: 10px !important;
      -o-border-radius: 10px !important;
      top: 50%;
      transform: translateY(-50%);
      margin-top: 0 !important;
    }
    .el-dialog__header {
      padding: 30px 30px 0;
      height: 26px;
      font-size: 18px;
      font-family: PingFangSC-Semibold, PingFang SC;
      font-weight: 600;
      line-height: 26px;
    }
    .el-dialog__headerbtn {
      top: 28px;
      right: 28px;
    }
    .el-dialog__headerbtn .el-dialog__close {
      font-weight: 600;
      color: #c0c4cc;
      font-size: 20px;
    }
    .el-dialog__headerbtn:focus .el-dialog__close,
    .el-dialog__headerbtn:hover .el-dialog__close {
      color: $color-main;
    }
    .el-dialog__body {
      padding: 20px 30px;
    }
    .el-dialog__footer {
      padding: 0 30px 30px;
    }
    .el-table thead {
      color: $color-text-normal;
      font-weight: 500;
    }
    .el-table th {
      height: 36px;
      padding: 0;
      background-color: #ebebf1;
    }
    .el-table td {
      padding: 17px 0;
      border-bottom: 0px solid #ebeef5;
    }
    .el-table .caret-wrapper {
      width: 10px;
      margin-left: 5px;
    }
    .el-table .sort-caret {
      left: 0px;
    }
    .el-table th > .cell {
      padding-left: 20px;
      padding-right: 20px;
    }
    .el-table--striped .el-table__body tr.el-table__row--striped td {
      background: $color-background;
    }
    .el-table--striped .el-table__body tr.el-table__row--striped.hover-row td {
      background: #ebebf1;
    }
    .el-table__body .el-table__row.hover-row td {
      background-color: #ebebf1;
    }
    .el-table tbody tr:hover > td {
      background-color: #ebebf1 !important;
    }
    .el-table .cell {
      line-height: 20px;
      padding-left: 20px;
      padding-right: 20px;
    }
    .el-table .radio .cell {
      padding-right: 10px !important;
    }
    .cy-t-popper {
      max-width: 280px;
    }
    .el-table__body-wrapper::-webkit-scrollbar {
      width: 6px;
      height: 6px;
    }
    .el-table__body-wrapper::-webkit-scrollbar-track {
      background-color: transparent;
    }
    .el-table__body-wrapper::-webkit-scrollbar-thumb {
      background-color: rgba(0, 0, 0, 0.3);
    }
    .el-table__body-wrapper::-webkit-scrollbar-thumb:hover {
      background-color: rgba(0, 0, 0, 0.3);
    }
    .el-table__body-wrapper::-webkit-scrollbar-thumb:active {
      background-color: rgba(0, 0, 0, 0.3);
    }
  }
}

.store-sort {
  display: flex;
  .sort-title {
    margin-right: 30px;
  }
}
</style>
