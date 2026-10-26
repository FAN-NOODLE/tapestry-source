<template>
  <div class="content">
    <!-- 左侧组织员工树 -->
    <div class="tree-box">
      <div class="tree-search clearfix">
        <div class="search-ipt pull-left">
          <cy-input
            class="input-with-select"
            v-model.trim="filter_text"
            size="small"
            width-size="small"
            :placeholder="`请输入${placeholder}`"
            clearable="clearable"
            @clear="reset"
          >
            <el-select
              style="width:100px"
              v-model="search_type"
              slot="prepend"
              @change="changeFilter"
            >
              <el-option
                v-for="item in serach_types"
                :label="item.name"
                :value="item.key"
                :key="item.key"
              >
              </el-option>
            </el-select>
          </cy-input>
        </div>
        <div class="search-btn pull-left m-l-20">
          <el-button
            size="small"
            type="primary"
            @click="changeFilter"
          >
            查询
          </el-button>
          <el-button
            size="small"
            @click="reset"
          >重置</el-button>
        </div>
      </div>
      <div class="tree">
        <el-tree
          v-show="is_show_tree"
          lazy
          ref="tree"
          :props="default_props"
          :load="loadNode"
          node-key="third_id"
          :check-strictly="true"
          :expand-on-click-node="AllExpandClick"
          @node-click="handleNodeClick"
          :default-expanded-keys="expand_key"
        >
          <div
            slot-scope="{ node, data }"
            class="custom-tree-node"
          >
            <!-- 部门(is_dp：true 代表是部门) -->
            <div
              v-if="data.is_dp"
              class="tree-item-box"
            >
              <img
                src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/group.png"
                alt=""
              />
              <span class="label">{{ node.label
              }}{{ data.mobile ? `(${data.mobile})` : "" }}</span>
              <em
                v-if="node.checked"
                class="el-icon-check pull-left m-l-10 select"
                :class="data.is_dp || data.mem_subscribe_status == 1
                    ? 'color-main'
                    : 'color-dark'
                  "
              ></em>
            </div>
            <!-- 员工 普通版显示(mem_subscribe_status: 1 已加入企业 其他：未加入 不可选) -->
            <div v-if="!data.is_dp && !data.is_more && !isGuide">
              <div
                v-if="data.mem_subscribe_status === 1"
                class="tree-item-box"
              >
                <img
                  src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/member_normal.png"
                  alt=""
                />
                <span class="label">{{ node.label
                }}{{ data.mobile ? `(${data.mobile})` : "" }}</span>
                <em
                  v-if="node.checked"
                  class="el-icon-check pull-left m-l-10 select"
                  :class="data.is_dp || data.mem_subscribe_status == 1
                      ? 'color-main'
                      : 'color-dark'
                    "
                ></em>
              </div>
              <div
                v-else
                class="tree-item-box"
              >
                <img
                  src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/member_dark.png"
                  alt=""
                />
                <span class="disabled label">{{ node.label
                }}{{ data.mobile ? `(${data.mobile})` : "" }}</span>
                <em
                  v-if="node.checked"
                  class="el-icon-check pull-left m-l-10 select"
                  :class="data.is_dp || data.mem_subscribe_status == 1
                      ? 'color-main'
                      : 'color-dark'
                    "
                ></em>
              </div>
            </div>
            <!-- 员工 导购版显示(is_guide：0 不是导购 1 已经是导购 不可选) -->
            <div v-if="!data.is_dp && !data.is_more && isGuide">
              <div
                v-if="data.mem_subscribe_status === 1 && data.is_guide === 1"
                class="tree-item-box"
              >
                <img
                  src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/member_normal.png"
                  alt=""
                />
                <span class="label">{{ node.label
                }}{{ data.mobile ? `(${data.mobile})` : "" }}</span>
                <em
                  v-if="node.checked"
                  class="el-icon-check pull-left m-l-10 select"
                  :class="data.is_dp || data.mem_subscribe_status == 1
                      ? 'color-main'
                      : 'color-dark'
                    "
                ></em>
              </div>
              <div
                v-else-if="data.mem_subscribe_status !== 1 || data.is_guide === 0
                  "
                class="tree-item-box"
              >
                <img
                  src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/member_dark.png"
                  alt=""
                />
                <span class="disabled label">{{ node.label
                }}{{ data.mobile ? `(${data.mobile})` : "" }}</span>
                <!-- <em
                  v-if="node.checked || data.is_guide === 1"
                  class="el-icon-check"
                ></em> -->
              </div>
            </div>
            <!-- 点击查看更多 -->
            <div v-if="!data.is_dp && data.is_more">
              <span
                class="remind-info"
                @click="addMore(node)"
              >{{ node.label
              }}{{ data.mobile ? `(${data.mobile})` : "" }}</span>
            </div>
          </div>
        </el-tree>
        <!-- 搜索时隐藏树，展示列表 -->
        <div
          class="filter-tree"
          v-show="!is_show_tree"
        >
          <div
            class="filter-box"
            v-if="filter_obj.length && search_type !== 'department'"
          >
            <div class="filter-title">员工</div>
            <div
              class="item-box"
              v-for="(item, index) in filter_obj"
              :key="index"
              @click="listClick(item)"
            >
              <div
                class="item tree-item-box"
                v-if="item.mem_subscribe_status === 1 && item.is_guide === 1"
              >
                <img
                  src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/member_normal.png"
                  alt=""
                />
                <span class="label">{{ item.user_name }} [{{ item.dept_name }},
                  {{ item.mobile || "" }}]</span>
              </div>
              <div
                class="item tree-item-box"
                v-else-if="item.mem_subscribe_status !== 1 || item.is_guide !== 1
                  "
              >
                <img
                  src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/member_dark.png"
                  alt=""
                />
                <span class="disabled label">{{ item.user_name }} [{{ item.dept_name }},
                  {{ item.mobile || "" }}]</span>
                <!-- <em
                  class="el-icon-check"
                  v-if="item.is_guide === 1"
                ></em> -->
              </div>
            </div>
          </div>
          <div
            class="filter-box"
            v-if="filter_obj.length && search_type === 'department'"
          >
            <div class="filter-title">部门</div>
            <div
              class="item-box"
              v-for="(item, index) in filter_obj"
              :key="index"
              @click="listClick(item)"
            >
              <div class="item tree-item-box">
                <img
                  src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/group.png"
                  alt=""
                />
                <span class="label">{{ item.dp_name
                }}{{ item.mobile ? `(${item.mobile})` : "" }}</span>
              </div>
            </div>
          </div>
          <div
            class="default"
            v-if="!filter_obj.length"
          >
            <span>暂无数据</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 右侧选择的数据 -->
    <div class="selecte-box">
      <div class="filter-title">已选员工</div>
      <div class="selecte-content">
        <div
          class="item-box"
          v-for="(item, index) in select_arr"
          :key="index"
        >
          <template v-if="!item.is_dp">
            <div class="item">
              <div class="tree-item-box">
                <img
                  src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/member_normal.png"
                  alt=""
                />
                <span class="label">{{ item.dp_name
                }}{{ item.mobile ? `(${item.mobile})` : "" }}</span>
              </div>
            </div>
            <em
              @click="deleteStaff(item)"
              class="el-icon-circle-close font-14 color-main pull-right select tree-item-box"
            ></em>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop, Watch } from "vue-property-decorator";
import Api from "@/api/wisdom/Personnel";
@Component
export default class extends Vue {
  @Prop({ default: true }) AllExpandClick?: boolean; // 是否允许点击父级展开下级
  @Prop({ default: true }) isGuide?: boolean; // 是否显示导购
  private checked_one = 1; // 导购的单选多选 0多选 1单选
  private expand_key: Array<any> = []; // 默认展开的key
  private filter_text = ""; // 搜索的字段
  private search_type = 'userid'; // 搜索类型
  private serach_types = [
    { name: "员工ID", key: 'userid' },
    { name: "手机号", key: 'mobile' },
    { name: "员工姓名", key: 'mem_username' },
    // { name: "部门名称", key: 'department' }
  ];
  private default_props = { label: "dp_name", isLeaf: "leaf" }; // 配置参数
  private select_arr: Array<any> = []; // 已选择的员工部门数据
  private is_show_tree = true; // 是否显示tree
  private filter_obj: Array<any> = []; // 筛选出的数据

  get placeholder() {
    return this.serach_types.find((item: any) => item.key === this.search_type)?.name || "";
  }

  @Watch("select_arr")
  onChangeSel(val: any) {
    this.$emit("selChange", val);
  }

  /*** methods ***/
  // 搜索框搜索
  private changeFilter() {
    this.is_show_tree = !this.filter_text ? true : false;
    this.filter_obj = [];
    if (!this.filter_text) return;
    if(this.search_type === 'department') {
      this.searchDp(this.filter_text);
    } else {
      this.searchPeople(this.filter_text);
    }
  }
  // 搜索框重置
  private reset() {
    this.filter_text = "";
    this.changeFilter();
  }

  // 搜索员工
  private searchPeople(value: string) {
    const params = {
      page: 1,
      size: 150,
      search: { [this.search_type]: value, filter_type: 1 },
    };
    Api.TreePeople(params).then((res: any) => {
      const list = res.data.list.map((item: any) => {
        item.dept_name = item.dp_name;
        return item;
      });
      this.filter_obj = list;
    });
  }
  // 删除已选员工
  private deleteStaff(data: any) {
    this.select_arr = this.select_arr.filter(
      (item: any) => item.third_id !== data.third_id
    );
    const dp_index = this.select_arr.findIndex(
      (v: any) => v.third_id === data.dp_parent_id
    );
    if (dp_index > -1) this.select_arr.splice(dp_index, 1);
    this.setChoice();
  }
  // 搜索部门
  searchDp(value: string) {
    Api.depListApi({ dp_name: value }).then((res: any) => {
      this.filter_obj = res.data.list;
    });
  }

  //获取没有部门的员工，与根节点平级展示
  private getNoDp() {
    return new Promise((resolve) => {
      const params = { page: 1, size: 600, search: { filter_type: 1 } };
      Api.TreePeople(params).then((res: any) => {
        const list = res.data.list;
        if (list.length <= 0) {
          resolve(list);
          return;
        }
        list.forEach((item: any) => {
          item.dp_name = item.user_name;
          item.dp_id = item.user_id;
          item.third_id = item.user_id;
          item.mobile = item.mobile || "";
          item.leaf = true;
          item.is_dp = false;
        });
        resolve(list);
      });
    });
  }

  // 懒加载树节点
  private async loadNode(node: any, resolve: any) {
    if (node.level === 0) {
      Promise.all([this.getDp(), this.getNoDp()]).then((values: any) => {
        const datas = [...values[0], ...values[1]];
        if (values[0].length > 0) {
          this.expand_key = [values[0][0].third_id];
        }
        resolve(datas);
      });
      this.setChoice();
    }
    if (node.level >= 1) {
      const params = { dp_id: node.data.dp_id };
      const list = await this.getDp(params);
      resolve(list);
      this.getPeople(node.data.dp_id, node);
      this.setChoice();
    }
  }

  // 获取部门
  private getDp(params?: any) {
    return new Promise((resolve) => {
      Api.depListApi(params).then((res: any) => {
        const data = res.data.list;
        if (!data || data.length <= 0) {
          resolve(data);
          return;
        }
        data.forEach((item: any) => {
          // if (!item.is_child_department) item.leaf = true;
          item.third_id = item.dp_id;
          item.is_dp = true;
        });
        resolve(data);
      });
    });
  }

  // 设置树选中
  private setChoice() {
    const key_arr = this.select_arr.map((v: any) => v.third_id);
    if (this.$refs.tree) {
      (this.$refs.tree as any).setCheckedKeys(key_arr);
    }
  }

  // 点击部门并部门下的员工
  private getPeople(id: any, node: any) {
    console.log(node, 2222);
    const params = {
      page: 1,
      size: 50,
      search: { dp_id: id, filter_type: 1 },
    };
    Api.TreePeople(params).then((res: any) => {
      const { list, page, size, total } = res.data;
      const list_data = list;
      if (list_data.length <= 0) return;
      if (list_data.length < total) {
        const data = {
          user_name: "点击加载更多",
          user_id: `${id}*`,
          is_more: true,
          meta: { page, size, total },
        };
        list_data.push(data);
      }

      list_data.forEach((item: any) => {
        item.dp_parent_id = id;
        item.dp_name = item.user_name;
        item.dp_id = item.user_id;
        item.third_id = item.user_id;
        item.mobile = item.mobile || "";
        item.leaf = true;
        item.is_dp = false;
        (this.$refs.tree as any).append(item, node);
      });
      if (this.checked_one === 0) {
        if (node.checked) this.setArr(list_data);
        this.setChoice();
      }
    });
  }

  private setArr(arr: Array<any>) {
    arr.forEach((item) => {
      const index = this.select_arr.findIndex(
        (v: any) => v.third_id === item.third_id
      );
      if (
        index <= -1 &&
        !item.is_dp &&
        item.mem_subscribe_status === 1 &&
        item.is_guide === 0 &&
        !item.is_more
      ) {
        this.select_arr.push(item);
      }
    });
  }

  // 点击组织节点
  private handleNodeClick(data: any, node: any) {
    // console.log(
    //   data,
    //   node,
    //   "点击组织节点",
    //   node.parent && node.parent.parent && node.parent.parent.childNodes
    // );
    if (data.is_more) return;
    if (!data.is_dp && data.mem_subscribe_status !== 1) {
      this.$message.warning("员工还未加入企业, 请重新选择");
      return;
    }
    if (this.isGuide && data.is_guide === 0) {
      this.$message.warning("该员工不是导购");
      return;
    }
    if (this.checked_one === 1) {
      if (this.select_arr.length === 1) {
        this.$message.warning("当前只能选择一位人员！");
        return false;
      }
    }
    node.checked = !node.checked;
    // 选中/取消外层对应的导购
    if (
      node.parent &&
      node.parent.parent &&
      node.parent.parent.childNodes &&
      node.parent.parent.childNodes.length > 0
    ) {
      node.parent.parent.childNodes.forEach((item: any) => {
        if (item.data.third_id == node.data.third_id) {
          item.checked = !item.checked;
        }
      });
    }
    if (node.checked) {
      // if (!data.children) {
      //   this.setChild(data);
      // }
      if ((this.checked_one == 1 && !data.children) || this.checked_one !== 1) {
        this.setChild(data);
      }
    } else {
      this.delChild(data);
    }
  }

  // 点击父级删除子级
  private delChild(data: any) {
    const dp_index = this.select_arr.findIndex(
      (v: any) => v.third_id === data.third_id
    );
    if (dp_index > -1) this.select_arr.splice(dp_index, 1);
    if (data.is_dp && data.children) {
      data.children.forEach((item: any) => {
        const index = this.select_arr.findIndex(
          (v: any) => v.third_id === item.third_id
        );
        if (index > -1) this.select_arr.splice(index, 1);
      });
    } else if (data.is_dp && !data.children) {
      this.select_arr = this.select_arr.filter(
        (v: any) => v.dp_parent_id !== data.third_id
      );
    }
    this.setChoice();
  }

  // 设置选择父级选中下面的子级
  private setChild(data: any) {
    if (data.is_dp && data.children) {
      this.select_arr.push(data);
      this.setArr(data.children);
    } else {
      if (this.checked_one === 1) {
        if (typeof data.is_child_department == "undefined") {
          this.select_arr.push(data);
        }
      } else {
        this.select_arr.push(data);
      }
    }
    this.setChoice();
  }
  // 搜索列表点击
  private async listClick(item: any) {
    // if (this.search_type === 1 && item.is_guide === 1) return;

    if (this.search_type !== 'department') {
      if (!item.is_dp && item.mem_subscribe_status !== 1) return;

      item.is_dp = false;
      item.third_id = item.user_id;
      item.dp_name = item.user_name;
    } else {
      item.is_dp = true;
      item.third_id = item.dp_id;
    }
    if (item.is_dp) {
      const params = {
        page: 1,
        size: 50,
        search: { dp_id: item.third_id, filter_type: 1 },
      };
      const res = await Api.TreePeople(params);
      const arr = res.data.list;
      arr.forEach((v: any) => {
        v.dp_parent_id = item.third_id;
        v.dp_name = v.user_name;
        v.dp_id = v.user_id;
        v.third_id = v.user_id;
        v.mobile = v.mobile || "";
        v.leaf = true;
        v.is_dp = false;
      });
      this.setArr(arr);
    } else {
      if (this.checked_one === 1) {
        if (this.select_arr.length === 1) {
          this.$message.warning("当前只能选择一位人员！");
          return false;
        }
      }
    }

    const data = this.select_arr.find((v) => v.third_id === item.third_id);
    if (!data) this.select_arr.push(item);
    this.setChoice();
  }
  // 点击查看更多按钮
  private addMore(node: any) {
    const parent = node.parent;
    const params = {
      page: node.data.meta.page + 1,
      size: 50,
      search: { dp_id: parent.data ? parent.data.dp_id : "", filter_type: 1 },
    };
    Api.TreePeople(params).then((res: any) => {
      const { list, page, size, total } = res.data;
      const list_data = list;
      if (list_data.length <= 0) return;
      node.data.meta = { page, size, total };
      list_data.forEach((item: any) => {
        item.dp_parent_id = parent.data.dp_id;
        item.dp_name = item.user_name;
        item.dp_id = item.user_id;
        item.third_id = item.user_id;
        item.mobile = item.mobile || "";
        item.leaf = true;
        item.is_dp = false;
        (this.$refs.tree as any).insertBefore(item, node);
      });
      const child_num = parent.childNodes.filter(
        (v: any) => !v.data.is_dp && !v.data.is_more
      );
      if (child_num.length === node.data.meta.total) {
        (this.$refs.tree as any).remove(node);
      }
      if (parent.checked) this.setArr(list_data);
      this.setChoice();
    });
  }
  // 生命周期
}
</script>
<style lang="scss" scoped>
.content {
  display: flex;

  .tree-box {
    flex: 1;
    border-right: 1px solid #ecedef;
    padding-right: 5px;

    .tree {
      .disabled {
        opacity: 0.5;
      }

      .filter-box {
        .filter-title {
          margin: 8px 0 10px 0;
          width: 100%;
          color: #787878;
        }

        .item-box {
          display: flex;
          justify-content: space-between;
          padding-right: 20px;
          height: 26px;
          line-height: 26px;
          color: #787878;
          cursor: pointer;
        }
      }

      .filter-tree {
        height: 100%;

        .default {
          display: flex;
          justify-content: center;
          align-items: center;
          height: 100%;
        }
      }
    }
  }

  .selecte-box {
    padding: 0 0 12px 12px;
    width: 35%;

    .title {
      margin-bottom: 20px;
    }

    .item-box {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding-right: 12px;
      line-height: 26px;

      .delete {
        cursor: pointer;
      }
    }
  }

  .tree,
  .selecte-content {
    height: 340px;
    overflow: auto;

    &::-webkit-scrollbar {
      /*滚动条整体样式*/
      width: 5px;
      /*高宽分别对应横竖滚动条的尺寸*/
      height: 1px;
    }

    &::-webkit-scrollbar-thumb {
      /*滚动条里面小方块*/
      border-radius: 2px;
      // box-shadow: inset 0 0 5px rgba(0, 0, 0, .2)
      background: #e4e7ed;
    }

    &::-webkit-scrollbar-track {
      /*滚动条里面轨道*/
      // box-shadow: inset 0 0 5px rgba(0, 0, 0, .2)
      border-radius: 2px;
      background: white;
    }
  }

  .selecte-content {
    margin-top: 12px;

    .label {
      max-width: 231px;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
  }

  .tree-item-box {
    display: flex;
    align-items: center;

    img {
      width: 14px;
    }

    .label {
      margin: 0 5px 0 10px;
    }
  }
}
</style>
