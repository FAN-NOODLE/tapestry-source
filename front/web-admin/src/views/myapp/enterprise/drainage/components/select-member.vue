<template>
  <div id="select-member">
    <el-dialog
      :title="dialogTitle"
      :visible.sync="centerDialogVisible"
      width="1000px"
      @close="cancel"
    >
      <div class="main-box">
        <div class="box clearfix">
          <div class="tree-left p-r-5 pull-left">
            <!-- 搜索框模块-->
            <div class="search-ctn clearfix">
              <div class="search-ipt pull-left">
                <cy-input
                  size="small"
                  width-size="medium"
                  :placeholder="placeholder"
                  v-model.trim="searchInput"
                  class="input-with-select"
                  clearable
                >
                  <el-select
                    v-model="selectType"
                    slot="prepend"
                  >
                    <el-option
                      v-for="item in selectDefault"
                      :label="item.label"
                      :value="item.value"
                      :key="item.vlaue"
                    ></el-option>
                  </el-select>
                </cy-input>
              </div>
              <div class="search-btn pull-left m-l-10">
                <el-button
                  type="primary"
                  size="small"
                  @click="search()"
                >查询</el-button>
                <el-button
                  size="small"
                  class="m-l-10"
                  @click="reset()"
                >重置</el-button>
              </div>
            </div>
            <!-- 左边组织结构和员工展示模块-->
            <div class="tree-box">
              <div class="scroll-tree">
                <el-tree
                  v-show="!searchInput"
                  :props="defaultProps"
                  node-key="value"
                  ref="tree"
                  lazy
                  :check-strictly="true"
                  :expand-on-click-node="false"
                  :check-on-click-node="true"
                  :show-checkbox="false"

                  :load="loadNode"
                  :default-expanded-keys="expand_key"
                  @node-click="handleNodeClick"
                >
                  <div
                    class="custom-tree-node clearfix"
                    slot-scope="{ node, data }"
                  >
                    <template v-if="!data.is_more">
                      <!-- 员工 普通版显示(mem_subscribe_status: 1 已加入企业 其他：未加入 不可选)
                    员工 导购版显示(is_guide：0 不是导购 1 已经是导购 不可选) -->
                      <!-- icon -> 部门 -->
                      <img
                        class="pull-left"
                        v-if="data.is_dp"
                        src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/group-ic.png"
                        alt=""
                      />
                      <!-- icon -> 正常员工 -->
                      <img
                        class="pull-left"
                        v-else-if="data.mem_subscribe_status == 1"
                        src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/staff_icon.png"
                        alt=""
                      />
                      <!-- icon -> 未加入企业 -->
                      <img
                        class="pull-left"
                        v-else
                        src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/member_dark.png"
                        alt=""
                      />
                      <!-- 非已加入企业员工姓名置灰-->
                      <div
                        class="pull-left m-l-10"
                        :class="data.is_dp || data.mem_subscribe_status == 1
                            ? ''
                            : 'color-dark'
                          "
                      >
                        {{ node.label
                        }}{{ data.mobile ? `(${data.mobile})` : "" }}
                      </div>
                      <!-- 已选员工对勾icon-->
                      <i
                        v-if="node.checked"
                        class="el-icon-check pull-left m-l-10 select"
                        :class="data.is_dp || data.mem_subscribe_status == 1
                            ? 'color-main'
                            : 'color-dark'
                          "
                        aria-hidden="true"
                      ></i>
                    </template>
                    <!-- 点击加载更多按钮-->
                    <template v-else>
                      <el-button
                        class="m-l-10"
                        type="text"
                      >{{
                        data.label
                      }}</el-button>
                    </template>
                  </div>
                </el-tree>
                <!-- search-content -->
                <div
                  v-show="searchInput"
                  class="custom-tree-node m-t-20"
                >
                  <div
                    class="clearfix search-item p-l-20"
                    v-for="item in searchList"
                    :key="item.id"
                    @click="selectSearchItem(item)"
                  >
                    <template v-if="!item.is_more">
                      <!-- icon -> 部门 -->
                      <img
                        class="pull-left"
                        v-if="selectType === 'department'"
                        src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/group-ic.png"
                        alt=""
                      />
                      <!-- icon -> 正常员工 -->
                      <img
                        class="pull-left"
                        v-else-if="item.mem_subscribe_status === 1"
                        src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/staff_icon.png"
                        alt=""
                      />
                      <!-- icon -> 未加入企业 -->
                      <img
                        class="pull-left"
                        v-else
                        src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/member_dark.png"
                        alt=""
                      />
                      <!-- 非已加入企业员工姓名置灰-->
                      <div
                        class="pull-left m-l-10"
                        :class="item.is_dp || item.mem_subscribe_status === 1
                            ? ''
                            : 'color-dark'
                          "
                      >
                        {{ item.name }} [{{ item.dept_name }},
                        {{ item.mobile || "" }}]
                      </div>
                      <!-- 已选员工对勾icon-->
                      <i
                        v-if="item.checked"
                        class="el-icon-check pull-left m-l-10 select"
                        :class="item.is_dp || item.mem_subscribe_status === 1
                            ? 'color-main'
                            : 'color-dark'
                          "
                        aria-hidden="true"
                      ></i>
                    </template>
                    <!-- 点击加载更多按钮-->
                    <template v-else>
                      <el-button type="text">
                        {{ item.name }}
                      </el-button>
                    </template>
                  </div>
                  <div
                    v-if="searchList.length == 0 && !searchLoading"
                    class="default"
                  >
                    {{
                      selectType == "department"
                      ? "没有找到相关部门"
                      : "没有找到相关员工"
                    }}
                  </div>
                </div>
              </div>
            </div>
          </div>
          <!-- 左边已选成员或部门展示模块-->
          <div class="tree-right p-l-20 pull-right">
            <div class="al-select p-b-10 font-14 color-text-normal">
              {{ hideSelectDepartment ? "已选择的员工" : "已选择的部门或成员" }}
            </div>
            <div class="al-box">
              <div
                class="clearfix custom-tree-node al-select-item"
                v-for="item in filterData(selectMember)"
                :key="item.value"
              >
                <!--部门icon -->
                <img
                  class="pull-left"
                  v-if="item.is_dp"
                  src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/group-ic.png"
                  alt=""
                />
                <!--已加入企业的正常员工icon -->
                <img
                  class="pull-left"
                  v-else-if="item.mem_subscribe_status === 1"
                  src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/staff_icon.png"
                  alt=""
                />
                <!--非正常员工icon -->
                <img
                  class="pull-left"
                  v-else
                  src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/member_dark.png"
                  alt=""
                />
                <div class="pull-left m-l-10 choose-name">
                  {{ item.name }}{{ item.mobile ? `(${item.mobile})` : "" }}
                </div>
                <!-- 删除已选人员的icon-->
                <i
                  @click="removeSelect(item)"
                  class="el-icon-circle-close font-14 color-main pull-right select"
                  aria-hidden="true"
                ></i>
              </div>
            </div>
          </div>
        </div>
      </div>
      <span
        slot="footer"
        class="dialog-footer"
      >
        <el-button
          size="small"
          @click="cancel"
        >取 消</el-button>
        <el-button
          size="small"
          type="primary"
          @click="submit"
        >确 定</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Watch, Prop } from "vue-property-decorator";
import API from "@/api/select-member";
@Component({
  name: "App"
})
export default class extends Vue {
  @Prop() hasSelectArr!: Array<any>;
  @Prop({ default: false, type: Boolean }) hideSelectDepartment?: boolean; //是否在右侧隐藏选择的部门
  @Prop({ default: 0, type: Number }) limiteNum?: number; // 限制选择员工人数, 0为不限制
  @Prop({ default: "添加员工", type: String }) dialogTitle?: string; // 限制选择员工人数, 0为不限制

  private centerDialogVisible = true;
  // 搜索类型
  private selectDefault = [
    { label: "员工ID", value: "userid" },
    { label: "员工手机号", value: "mobile" },
    { label: "员工姓名", value: "mem_username" },
    { label: "部门名称", value: "department" }
  ];
  private selectType = "userid"; // 搜索类型默认选中
  private searchInput = ""; // 搜索内容
  // tree 配置项
  private defaultProps = {
    isLeaf: "leaf",
    children: "children",
    label: "name",
  };
  private expand_key: Array<string> = []; // 默认展示
  private searchList: Array<any> = []; // 搜索返回值
  private selectMember: Array<any> = []; // 选中项
  private searchLoading = true; // 搜索框loading

  get placeholder() {
    const placeholder = this.selectDefault.find(item => item.value === this.selectType)?.label || "";
    return `请输入${placeholder}`
  }

  // private mounted() {}
  @Watch("selectType")
  private onChangeSelectType() {
    // 监听清空搜索数据
    this.searchList = [];
  }
  @Watch("searchInput")
  private onChangeSearchInput() {
    // 监听清空搜索数据
    this.searchInput.length === 0 && (this.searchList = []);
  }
  // 过滤展示右边的已选列表
  private filterData(data: any) {
    if (this.hideSelectDepartment) {
      return data.filter((item: any) => !item.is_dp);
    } else {
      return data;
    }
  }

  //获取没有部门的员工，与根节点平级展示
  private async getNoDp() {
    try {
      const params = {
        page: 1,
        size: 9999,
        get_no_dp: 1,
        search: {
          // dp_id: "",
          mem_username: ""
        }
      };
      const res = await API.treePeople(params);
      const list = res.data.no_dp_list;
      list.forEach((item: any) => {
        console.log(2222, item);
        
        item.name = item.user_name;
        item.value = item.user_id;
        item.mobile = item.mobile || "";
        // 设置叶子字段
        item.leaf = true;
        // 设置部门
        item.is_dp = false;
      });
      return list;
    } catch (error) {
      console.error(error);
    }
  }
  // 获取部门
  private async getDp(params: any) {
    try {
      const res = await API.depListApi(params);
      const data = res.data || [];
      data.forEach((item: any) => {
        // 设置部门字段
        item.is_dp = true;
        item.name = item.dp_name;
        item.value = item.third_id;
      });
      return data;
    } catch (error) {
      console.error(error);
    }
  }
  // 根据id获取部门人员
  private async getPeople(id: string | number, node: any) {
    // 组装请求参数
    const params = {
      page: node.data.page ? node.data.page + 1 : 1,
      size: node.data.size ? node.data.size : 500,
      search: { dp_id: id }
    };
    const res = await API.treePeople(params);
    const { list, ...meta } = res.data;
    if (list.length <= 0) return;
    // 加载更多按钮组装
    if (meta.page * meta.size <= meta.total) {
      const data = {
        label: "点击加载更多",
        user_id: `${id}*`,
        is_more: true,
        leaf: true,
        ...meta
      };
      list.push(data);
    }
    // 插入到tree方法
    const treeAppend = node.data.is_more ? "insertBefore" : "append";
    // 循环插入
    list.forEach((item: any) => {
      // item.checked = !!node.checked;
      item.dp_parent_id = id;
      item.name = item.user_name;
      item.value = item.user_id;
      item.leaf = true;
      item.is_dp = false;
      item.mobile = item.mobile || "";
      (this.$refs.tree as any)[treeAppend](item, node);
    });
    // 没有更多则删除更多按钮
    if (node.data.is_more) (this.$refs.tree as any).remove(node);
  }
  // tree load事件
  private async loadNode(node: any, resolve: any) {
    // 根节点
    if (node.level === 0) {
      
      const dp_arr = await this.getDp({ is_top: 1 });
      if (dp_arr && dp_arr.length > 0) this.expand_key.push(dp_arr[0].value);
      const hum_arr = (await this.getNoDp()) || [];
      resolve([...dp_arr, ...hum_arr]);
      // 分节点
    } else if (node.level >= 1) {
      const params = { third_id: node.data.third_id };
      const list = await this.getDp(params);
      resolve(list);
      console.log("loadNode", node.checked);
      this.getPeople(node.data.dp_id, node);
    }
    
  }

  /** 选中节点 */
  handleCheckChange() {
    this.selectMember = (this.$refs.tree as any).getCheckedNodes()
    console.log(this.selectMember);
    
  }
  // 点击tree回调
  private handleNodeClick(data: any, node: any) {
    console.log(data.value, node.checked);
    if(!data.is_dp && Number(data.mem_subscribe_status) !== 1) {
      this.$message.warning(`员工${data.name}还未加入企业, 请重新选择`);

      return false
    }
    
    if(node.checked) {
      this.selectMember.push(data)
    } else {
      this.selectMember = this.selectMember.filter(item => item.value !== data.value)
    }
  }
  // 设置选择父级选中下面的子级
  setChild(data: any) {
    if (data.is_dp && data.children) {
      this.selectMember.push(data);
      this.setArr(data.children);
    }
    this.setChoice();
  }
  // 将可选子级节点全部选中
  setArr(arr: any) {
    arr.forEach((item: any) => {
      const index = this.selectMember.findIndex(v => v.value === item.value);
      if (
        index <= -1 &&
        !item.is_dp &&
        item.mem_subscribe_status === 1 &&
        !item.is_more
      ) {
        this.selectMember.push(item);
      }
    });
  }
  // 点击父级删除子级
  delChild(data: any) {
    const dp_index = this.selectMember.findIndex(v => v.value === data.value);
    if (dp_index > -1) this.selectMember.splice(dp_index, 1);
    if (data.is_dp && data.children) {
      data.children.forEach((item: any) => {
        const index = this.selectMember.findIndex(v => v.value === item.value);
        if (index > -1) this.selectMember.splice(index, 1);
      });
    } else if (data.is_dp && !data.children) {
      this.selectMember = this.selectMember.filter(
        v => v.dp_parent_id !== data.value
      );
    }
    this.setChoice();
  }
  // 右侧已选移出按钮
  private removeSelect(data: any) {
    this.selectMember = this.selectMember.filter(
      item => item.value !== data.value
    );
    // 搜索后右边点取消左边同步
    this.searchList.map(v => {
      if (v.value == data.value || v.value == data.dp_parent_id) {
        v.checked = false;
      }
    });
    // 删除了某部门的员工时，取消某部门的对勾
    const dp_index = this.selectMember.findIndex(
      v => v.value === data.dp_parent_id
    );
    if (dp_index > -1) this.selectMember.splice(dp_index, 1);
    this.setChoice();
  }
  // 设置树选中
  private setChoice() {
    const key_arr = this.selectMember.map(v => v.value);
    console.log(key_arr);
    
    (this.$refs.tree as any).setCheckedKeys(key_arr);
  }
  // 搜索框事件 搜索
  private async search() {
    if (this.searchInput.length === 0) return;
    this.searchLoading = true;
    this.selectType === "department" ? this.searchDp() : this.searchMember({});
  }
  // 搜索框事件 重置
  private reset() {
    this.searchInput = "";
    this.selectType = "userid";
  }
  // 搜索员工
  private async searchMember(obj: any) {
    try {
      const params = {
        page: obj.page ? obj.page + 1 : 1,
        size: obj.size ? obj.size : 15,
        search: { [this.selectType]: this.searchInput }
      };
      const {
        data: { list, ...meta }
      } = await API.treePeople(params);
      list.forEach((item: any) => {
        item.value = item.user_id;
        item.name = item.user_name;
        this.selectMember.some(
          (v: any) => v.value === item.value && (item.checked = true)
        );
      });
      if (list.length && meta.page * meta.size < meta.total) {
        list.push({
          value: Date.parse(Date() + ""),
          is_more: true,
          ...meta,
          name: "点击加载更多"
        });
      }
      // 清除上次的加载更多按钮
      obj.page && this.searchList.pop();
      // 如果是点击加载更多，就拼接数据
      if (obj.is_more) {
        this.searchList = [...this.searchList, ...list];
      } else {
        this.searchList = [...list];
      }
      this.searchLoading = false;
    } catch (error) {
      console.error(error);
      this.searchLoading = false;
    }
  }
  // 搜索部门
  private async searchDp() {
    try {
      const res = await API.depListApi({
        dp_name: this.searchInput
      });
      const list = res.data || [];
      list.forEach((item: any) => {
        item.name = item.dp_name;
        item.is_dp = true;
        this.selectMember.some(
          (v: any) => v.value === item.value && (item.checked = true)
        );
      });
      this.searchList = list;
      this.searchLoading = false;
    } catch (error) {
      console.error(error);
      this.searchLoading = false;
    }
  }
  // 点击搜索内容
  private selectSearchItem(item: any) {
    if (item.is_dp) {
      this.selectSearchDp(item);
    } else if (item.is_more) {
      this.searchMember(item);
    } else {
      this.selectSearchMember(item);
    }
  }
  // 选择搜索的部门
  private async selectSearchDp(item: any) {
    if (item.checked) {
      item.checked = false;
      this.selectMember = this.selectMember.filter(
        v => v.value !== item.value && v.dp_parent_id !== item.value
      );
    } else {
      // 选中该部门下面的所有直属员工
      // const params = {
      //   page: 1,
      //   size: 999,
      //   search: { dp_id: item.dp_id }
      // };
      // const res = await API.treePeople(params);
      // const arr = res.data.list;
      // if (arr.length > 0) {
      //   item.checked = true;
      //   arr.forEach((v: any) => {
      //     v.dp_parent_id = item.dp_id;
      //     v.name = v.user_name;
      //     v.value = v.user_id;
      //     v.leaf = true;
      //     v.is_dp = false;
      //     v.dept_name = v.dept_name || "";
      //     v.mobile = v.mobile || "";
      //   });
      //   this.setArr(arr);
        this.selectMember.push(item);
      // } else {
      //   this.$message.warning("该部门下没有可选的员工");
      // }
    }
    this.setChoice();
  }
  // 选择搜索的人员
  private selectSearchMember(item: any) {
    if (Number(item.mem_subscribe_status) === 1) {
      if (item.checked) {
        item.checked = false;
        this.selectMember = this.selectMember.filter(
          v => v.value !== item.value
        );
      } else {
        item.checked = true;
        this.selectMember.push(item);
      }
      this.setChoice();
    } else {
      this.$message.warning(`员工${item.name}还未加入企业, 请重新选择`);
      return;
    }
  }
  private submit() {
    const memberList = this.filterData(this.selectMember);
    // 如果产品需求对选人数量有限制，则提示
    if (this.limiteNum && memberList.length > this.limiteNum) {
      this.$message.warning(`最多只能选择${this.limiteNum}个员工`);
      return;
    }
    //this.centerDialogVisible = false;
    this.$emit("submit", memberList);
  }
  private cancel() {
    this.$emit("close", false);
  }
  private mounted() {
    if (this.hasSelectArr && this.hasSelectArr.length) {
      //解除this.selectMember与this.hasSelectArr双向绑定关系避免选人的时候将影响父组件元素
      this.selectMember = this.hasSelectArr.concat([]);
    }
    setTimeout(() => {
      this.setChoice();
    }, 0);
  }
}
</script>

<style lang="scss" scoped>
::v-deep.el-select .el-input {
  width: 110px;
}

::v-deep.el-input-group__append,
::v-deep.el-input-group__prepend {
  background-color: $color-white;
}

::v-deep.el-tree {
  box-shadow: none;
}

.box {
  width: 100%;
}

.tree-left {
  width: 570px;
}

.search-ipt {
  width: 400px;
}

.tree-left {
  position: relative;

  &:after {
    content: "";
    width: 1px;
    height: 100%;
    position: absolute;
    top: 0;
    right: 0;
    background-color: $color-third-line;
  }
}

.custom-tree-node {
  position: relative;
  width: 100%;

  img {
    width: 14px;
    height: 14px;
    position: relative;
    top: 3px;
  }

  .select {
    font-size: 12px;
    position: relative;
    top: 5px;
    font-weight: 900;
  }

  .default {
    position: absolute;
    left: 0;
    top: 180px;
    width: 100%;
    text-align: center;
  }
}

::v-deep.el-tree-node__content {
  height: 30px;
}

.color-dark {
  color: $color-placeholder;
}

.tree-box {
  height: 402px;
  overflow: hidden;
  overflow-y: auto;
  margin-top: 20px;
}

.scroll-tree {
  width: 100%;
  min-height: 100%;
}

.search-item {
  height: 30px;
  line-height: 30px;
  cursor: pointer;

  &:hover {
    background-color: rgba($color-main, 0.1);
  }

  img {
    top: 8px;
  }

  .select {
    top: 9px;
  }
}

.tree-right {
  width: 323px;
  height: 458px;
  overflow-y: auto;

  .al-select {
    font-family: PingFangSC-Medium, PingFang SC;
    font-weight: 500;
  }

  .select {
    font-weight: 500;
    cursor: pointer;
    padding: 8px;
    top: 0;

    &:hover {
      font-weight: 900;
    }
  }

  .al-select-item {
    height: 30px;
    line-height: 30px;

    img {
      top: 8px;
    }
  }

  .choose-name {
    height: 30px;
    width: 260px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}

#crs-scrm .el-tree {
  padding: 0;
}

#crs-scrm .el-button+.el-button {
  margin-left: 10px !important;
}</style>
