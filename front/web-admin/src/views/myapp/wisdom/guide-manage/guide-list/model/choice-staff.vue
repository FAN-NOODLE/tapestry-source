<template>
  <div class="contaniner">
    <el-dialog
      class="dialog-box"
      width="1000px"
      :title="dialog_name"
      :visible.sync="dialogVisible"
      :close-on-click-modal="false"
      @close="closeDialog"
    >
      <div class="dialog-body">
        <el-steps
          class="step-wrap"
          v-if="cur_step === 1"
          :active="activeStep"
          finish-status="success"
          :space="600"
        >
          <el-step title="选择企业管理人员"></el-step>
          <el-step title="选择管理范围"></el-step>
        </el-steps>

        <div v-if="activeStep === 0">
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
                    :placeholder="`请输入${serach_types[search_type - 1].name}`"
                    clearable="clearable"
                  >
                    <el-select
                      class="search-select"
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
                        v-if="data.mem_subscribe_status === 1 && data.is_guide === 0
                          "
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
                        v-else-if="data.mem_subscribe_status !== 1 || data.is_guide === 1
                          "
                        class="tree-item-box"
                      >
                        <img
                          src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/member_dark.png"
                          alt=""
                        />
                        <span class="disabled label">{{ node.label
                        }}{{ data.mobile ? `(${data.mobile})` : "" }}</span>
                        <em
                          v-if="node.checked || data.is_guide === 1"
                          class="el-icon-check"
                        ></em>
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
                    v-if="filter_obj.length && search_type === 1"
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
                        v-if="item.mem_subscribe_status === 1 && item.is_guide === 0
                          "
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
                        v-else-if="item.mem_subscribe_status !== 1 || item.is_guide === 1
                          "
                      >
                        <img
                          src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/member_dark.png"
                          alt=""
                        />
                        <span class="disabled label">{{ item.user_name }} [{{ item.dept_name }},
                          {{ item.mobile || "" }}]</span>
                        <em
                          class="el-icon-check"
                          v-if="item.is_guide === 1"
                        ></em>
                      </div>
                    </div>
                  </div>
                  <div
                    class="filter-box"
                    v-if="filter_obj.length && search_type === 2"
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
              <div class="filter-title">{{ choice_name }}</div>
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
        </div>
      </div>

      <select-member-item
        v-if="activeStep === 1"
        ref="selectMemberItem"
        :tabName="tabName"
        :selectArr="selArr"
        @submit="getData"
      ></select-member-item>

      <div
        slot="footer"
        class="dialog-footer"
      >
        <!-- 社群导购、企业导购 -->
        <template v-if="cur_step !== 1">
          <el-button
            size="small"
            @click="closeDialog"
          >
            取 消
          </el-button>
          <el-button
            size="small"
            type="primary"
            @click="submitDialog"
          >
            确 认
          </el-button>
        </template>

        <!-- 企业管理 -->
        <template v-else>
          <!-- 第一步 -->
          <template v-if="activeStep === 0">
            <el-button
              size="small"
              @click="closeDialog"
            >
              取消
            </el-button>
            <el-button
              size="small"
              type="primary"
              @click="nextCom(1)"
            >
              下一步
            </el-button>
          </template>

          <!-- 第二步 -->
          <template v-if="activeStep === 1">
            <el-button
              v-if="isEdit"
              size="small"
              @click="closeDialog"
            >
              取消
            </el-button>
            <el-button
              v-else
              size="small"
              @click="nextCom(0)"
            >
              上一步
            </el-button>
            <el-button
              size="small"
              type="primary"
              @click="submitDialog"
            >
              确 认
            </el-button>
          </template>
        </template>
      </div>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop, Emit } from "vue-property-decorator";
import UserModule from "@/store/modules/user";
import Api from "@/api/wisdom/Personnel";
import SelectMemberItem from "@/components/sy-components/select-member-item.vue";

@Component({
  components: {
    SelectMemberItem,
  },
})
export default class extends Vue {
  @Prop({ default: false }) isEdit?: boolean; // true编辑，false新增
  @Prop({ default: 0 }) cur_step?: number; // 是否有多个步骤，0否 1是
  @Prop({ default: 0 }) checked_one?: number; // 导购的单选多选
  @Prop({ default: "添加员工" }) dialog_name?: string; // 弹框名称
  @Prop({ default: "" }) remind_name?: string; // 提示文字
  @Prop({ default: false }) visible?: boolean; // 控制模态框显隐
  @Prop({ default: () => [] }) selectArr?: Array<any>; // 选择的员工
  @Prop({ default: () => [] }) selectedRange?: Array<any>; // 已选择的管理范围
  @Prop({ default: 2 }) selectedRangeType?: number; // 已选管理范围类型：2 区域，4部门
  @Prop({ default: false }) allowNull?: boolean; // 允许选择空保存
  @Prop({ default: false }) AllExpandClick?: boolean; // 是否允许点击父级展开下级
  @Prop({ default: true }) AllowSelectDp?: boolean; // 是否允许选择部门
  @Prop({ default: false }) isGuide?: boolean; // 是否显示导购
  @Prop({ default: "已选择的部门或成员" }) choice_name?: string; // 右侧选择标题
  @Prop({ default: false }) IsChoiceChild?: boolean; // 是否选择父级选中子级

  /*** data ***/
  private expand_key: Array<any> = []; // 默认展开的key
  private search_type = 1; // 搜索类型
  private serach_types = [
    { name: "员工姓名", key: 1 },
    { name: "部门名称", key: 2 },
  ];
  private default_props = { label: "dp_name", isLeaf: "leaf" }; // 配置参数
  private filter_text = ""; // 搜索的字段
  private is_show_tree = true; // 是否显示tree
  private filter_obj: Array<any> = []; // 筛选出的数据
  private select_arr: Array<any> = []; // 已选择的员工部门数据
  private activeStep = 0; // 步骤条
  private selArr: Array<any> = []; // 已选择的管理范围
  private rangeType: string | number = "";

  private dialogVisible = true; // 控制模态框显示/隐藏
  private tabName = ""; // 第二步 默认选中的tab名称

  get USER() {
    return UserModule.getUserInfo;
  }

  /*** methods ***/
  // 搜索框搜索
  private changeFilter() {
    this.is_show_tree = !this.filter_text ? true : false;
    this.filter_obj = [];
    if (!this.filter_text) return;
    this.search_type === 1
      ? this.searchPeople(this.filter_text)
      : this.searchDp(this.filter_text);
  }

  // 搜索框重置
  private reset() {
    this.filter_text = "";
    this.changeFilter();
  }

  // 点击下一步
  private nextCom(step: number) {
    if (this.select_arr.length <= 0) {
      this.$message.warning("请选择企业成员");
      return;
    }
    this.activeStep = step;
  }

  // 搜索员工
  private searchPeople(value: string) {
    const params = {
      page: 1,
      size: 150,
      search: { mem_username: value, filter_type: 1 },
    };
    Api.TreePeople(params).then((res: any) => {
      const list = res.data.list.map((item: any) => {
        item.dept_name = item.dp_name;
        return item;
      });
      this.filter_obj = list;
      console.log();
      // list.forEach(item => {
      //   if (
      //     (!this.isGuide && item.mem_subscribe_status === 1) ||
      //     (this.isGuide &&
      //       item.mem_subscribe_status === 1 &&
      //       item.is_guide === 0)
      //   ) {
      //     this.filter_obj.push(item);
      //   }
      // });
    });
  }

  // 获取管理门店的管理范围
  private getData(data: any) {
    this.selArr = data.arr;
    this.rangeType = data.type;
  }

  // 搜索部门
  searchDp(value: string) {
    Api.depListApi({ dp_name: value }).then((res: any) => {
      this.filter_obj = res.data.list;
    });
  }

  // 搜索列表点击
  private async listClick(item: any) {
    if (this.search_type === 1 && item.is_guide === 1) return;

    if (this.search_type === 1) {
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
    }
    const data = this.select_arr.find((v) => v.third_id === item.third_id);
    if (!data) this.select_arr.push(item);
    this.setChoice();
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

  // 点击组织节点
  private handleNodeClick(data: any, node: any) {
    console.log(
      data,
      node,
      "点击组织节点",
      node.parent && node.parent.parent && node.parent.parent.childNodes
    );
    if (data.is_more) return;
    if (!data.is_dp && data.mem_subscribe_status !== 1) {
      this.$message.warning("员工还未加入企业, 请重新选择");
      return;
    }
    if (this.isGuide && data.is_guide === 1) {
      this.$message.warning("该员工已经是导购了");
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

  // 关闭模态框的回调
  @Emit("close")
  private closeDialog() {
    this.dialogVisible = false;
  }

  // 选择完成提交数据
  @Emit("submit")
  private submitDialog() {
    if (this.cur_step === 1) {
      (this.$refs.selectMemberItem as any).save();
      if (!this.selArr.length) {
        this.$message.warning("请选择管理范围");
        return;
      }
    } else {
      if (!this.select_arr.length) {
        this.$message.warning("请选择成员");
        return;
      }
    }
    const memberList: Array<any> = [];
    this.select_arr.forEach((item) => {
      if (!item.is_dp) {
        memberList.push({
          mem_uid: item.mem_uid,
          user_name: item.user_name,
          mobile: item.mobile,
          avatar: item.avatar,
          user_id: item.user_id,
        });
      }
    });
    const data = {
      isEdit: this.isEdit,
      memberList,
    };
    if (this.cur_step === 1) {
      data["rangeList"] = this.selArr;
      data["rangeType"] = this.rangeType;
    }
    return data;
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
    (this.$refs.tree as any).setCheckedKeys(key_arr);
  }

  private mounted() {
    if (this.selectArr?.length) {
      this.select_arr = this.selectArr;
      this.select_arr.forEach((item) => {
        item.third_id = item.value;
        item.dp_name = item.label;
        item.is_dp = item.type === 1 ? true : false;
      });
    }
    if (this.selectedRange?.length) {
      this.selArr = this.selectedRange;
    }
    this.activeStep = this.isEdit ? 1 : 0;
    this.tabName =
      this.USER.version === 3
        ? "SelectMemberByDp"
        : this.selectedRangeType === 4
          ? "SelectMemberByDp"
          : "SelectMemberByArea";
  }
}
</script>

<style lang="scss" scoped>
/deep/.el-dialog {
  max-height: 685px !important;
}

/deep/ .el-dialog__body {
  padding: 0 30px;
}

/deep/.el-step__title.is-wait {
  color: #606266;
}

/deep/.el-step__icon.is-icon {
  color: #606266;
}

/deep/.el-step__icon.is-icon {
  width: 22px !important;
}

/deep/.el-input-group__append,
/deep/.el-input-group__prepend {
  background-color: $color-white;
}

.dialog-body {
  .step-wrap {
    margin-bottom: 20px;
  }
}

.el-icon-check {
  font-weight: 600;
  font-size: 17px;
  opacity: 0.5;
}

.input-width {
  width: 240px;
  margin-left: 5px;
}

.search-select {
  width: 110px;
}

.dialog-box {}

.remind-info {
  color: #409eff;
  margin-left: 23px;
}

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
    height: 380px;
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

.el-select-dropdown__item.selected {
  font-weight: normal;
}
</style>
