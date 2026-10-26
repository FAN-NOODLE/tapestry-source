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
                  placeholder="请输入员工姓名"
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
                  :data="treeData"
                  :filter-node-method="filterNode"
                  :props="defaultProps"
                  node-key="dp_id"
                  :default-expanded-keys="[this.rootId.id]"
                  :default-checked-keys="this.smIdList"
                  ref="tree"
                  check-strictly
                  @node-click="handleNodeClick"
                >
                  <div
                    class="custom-tree-node clearfix"
                    slot-scope="{ node, data }"
                  >
                    <template>
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
                        v-else-if="data.mem_subscribe_status === 1 && !data.disabled
                          "
                        src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/staff_icon.png"
                        alt=""
                      />
                      <!-- icon -> 不可操作的员工 -->
                      <img
                        class="pull-left"
                        v-else-if="data.disabled"
                        src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/member_dark.png"
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
                  </div>
                </el-tree>
              </div>
            </div>
          </div>
          <!-- 右边已选成员或部门展示模块-->
          <div class="tree-right p-l-20 pull-right">
            <div class="al-select p-b-10 font-14 color-text-normal">
              已选择的员工
            </div>
            <div class="al-box">
              <div
                class="clearfix custom-tree-node al-select-item"
                v-for="item in filterData(selectMember)"
                :key="item.dp_id"
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
                  v-if="!item.disabled"
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
import { Component, Vue, Prop } from "vue-property-decorator";
import API from "@/api/select-member";
@Component({
  name: "App"
})
export default class extends Vue {
  @Prop({
    default: () => {
      return {
        id: 0,
        name: ""
      };
    },
    type: Object
  })
  rootId; //指定根节点的id
  @Prop() hasSelectArr!: Array<object>;
  @Prop({ default: "添加员工", type: String }) dialogTitle?: string; // 限制选择员工人数, 0为不限制

  private centerDialogVisible = true;
  // 搜索类型
  private selectDefault = [{ label: "员工姓名", value: "member" }];
  private selectType = "member"; // 搜索类型默认选中
  private searchInput = ""; // 搜索内容
  // tree 配置项
  private defaultProps = {
    isLeaf: "leaf",
    children: "children",
    label: "name"
  };
  private selectMember: Array<any> = []; // 选中项
  private treeData: Array<any> = []; // 树结构
  private smIdList: Array<string> = []; //店长id集合

  //根据搜索框对树节点人员进行过滤
  private filterNode(value, data) {
    console.log("filterNode", data);
    if (!value) return true;
    return !data.is_dp && data.name.indexOf(value) !== -1;
  }
  // 过滤展示右边的已选（部门）列表
  private filterData(data: any) {
    return data.filter((item: any) => !item.is_dp);
  }
  /* 根据id查询部门人员-不分页 */
  private async getMemberFromId(id: number) {
    try {
      const params = {
        dp_ids: [id]
      };
      const res = await API.getMemberFromId(params);
      if (res["code"] == 200) {
        res.data.forEach((member) => {
          member.dp_id = member.userid;
          member.dp_parent_id = id;
          member.leaf = true;
          member.is_dp = false;
          member.mem_subscribe_status = member.status; //激活状态: 0=刚创建 1=已激活，2=已禁用，4=未激活，5=退出企业
          member.disabled = member.type == 4; //导购类型：4：店长，其他：门店导购 店长默认选中且不能点击
          //找到店长，并存起来，默认选中
          if (member.type == 4) {
            this.smIdList.push(member.userid || member.user_id);
            this.selectMember.push(member);
          }
        });
        return res.data || [];
      }
      return [];
    } catch (e) {
      console.error("getMemberFromId", e);
      return [];
    }
  }

  // 点击tree回调
  private handleNodeClick(data: any, node: any) {
    console.log("click node -> ", data, node);
    if (!data.is_dp) {
      // 正常员工 && 可点击
      if (Number(data.mem_subscribe_status) === 1 && !data.disabled) {
        if (node.checked) {
          node.checked = false;
          // 取消选中，则在右侧列表清理掉该员工
          this.selectMember = this.selectMember.filter(
            (item) => item.dp_id !== data.dp_id
          );
          // 全选部门后，取消某一员工，该部门的对勾得取消掉
          const dp_index = this.selectMember.findIndex(
            (v) => v.dp_id === data.dp_parent_id
          );
          if (dp_index > -1) this.selectMember.splice(dp_index, 1);
        } else {
          node.checked = true;
          this.selectMember.push(data);
          // 如果在全选状态下,取消了某一兄弟选项，则他们父级的对勾得取消掉
          const allChecked = node.parent.childNodes.every((key: any) => {
            return key.checked == true;
          });
          if (allChecked) {
            node.parent.checked = true;
            this.selectMember.push(node.parent.data);
          }
        }
        this.setChoice();
      } else {
        if (data.disabled) {
          this.$message.warning(`不能操作店长${data.name}`);
          return;
        }
        this.$message.warning(`员工${data.name}还未加入企业, 请重新选择`);
        return;
      }
    } else if (data.is_dp) {
      node.checked = !node.checked;
      //选中父节点
      if (node.checked) {
        this.setChild(data);
      } else {
        //取消选中父节点
        this.delChild(data);
      }
      console.log("handleNodeClick", data, node.checked);
    }
    return false;
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
      const index = this.selectMember.findIndex((v) => v.dp_id === item.dp_id);
      if (index <= -1 && !item.is_dp && item.mem_subscribe_status === 1) {
        this.selectMember.push(item);
      }
    });
  }
  // 点击父级删除子级
  delChild(data: any) {
    const dp_index = this.selectMember.findIndex((v) => v.dp_id === data.dp_id);
    if (dp_index > -1) this.selectMember.splice(dp_index, 1);
    //如果父节点是部门 同时包含有子节点
    if (data.is_dp && data.children) {
      data.children.forEach((item: any) => {
        const index = this.selectMember.findIndex(
          (v) => v.dp_id === item.dp_id
        );
        //如果有选中项
        if (index > -1 && !item.disabled) this.selectMember.splice(index, 1);
      });
    } else if (data.is_dp && !data.children) {
      this.selectMember = this.selectMember.filter(
        (v) => v.dp_parent_id !== data.dp_id
      );
    }
    this.setChoice();
  }
  // 右侧已选移出按钮
  private removeSelect(data: any) {
    this.selectMember = this.selectMember.filter(
      (item) => item.dp_id !== data.dp_id
    );
    // 删除了某部门的员工时，取消某部门的对勾
    const dp_index = this.selectMember.findIndex(
      (v) => v.dp_id === data.dp_parent_id
    );
    if (dp_index > -1) this.selectMember.splice(dp_index, 1);
    this.setChoice();
  }
  // 设置树选中
  private setChoice() {
    const key_arr = this.selectMember.map((v) => v.dp_id);
    (this.$refs.tree as HTMLFormElement).setCheckedKeys(key_arr);
  }
  // 搜索框事件 搜索
  private async search() {
    (this.$refs.tree as HTMLFormElement).filter(this.searchInput);
  }
  // 搜索框事件 重置
  private reset() {
    this.searchInput = "";
    (this.$refs.tree as HTMLFormElement).filter("");
  }
  private submit() {
    const memberList = this.filterData(this.selectMember);
    this.$emit("submit", memberList);
  }
  private cancel() {
    this.$emit("close", false);
  }
  async mounted() {
    const member = await this.getMemberFromId(this.rootId.id);
    //当前门店
    const root = {
      is_dp: true,
      name: this.rootId.name,
      dp_id: this.rootId.id,
      children: [...member]
    };
    this.treeData = [root];
    if (this.hasSelectArr && this.hasSelectArr.length) {
      //解除this.selectMember与this.hasSelectArr双向绑定关系避免选人的时候将影响父组件元素
      this.selectMember = JSON.parse(JSON.stringify(this.hasSelectArr));
      this.selectMember.forEach((item: any) => {
        item.dp_parent_id = this.rootId.id;
        item.userid = item.user_id;
        item.mem_subscribe_status = 1;
        item.dp_id = item.user_id;
        if (this.smIdList.includes(item.user_id)) {
          item.disabled = true;
        }
      });
      console.log("mounted", this.selectMember);
    }
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
