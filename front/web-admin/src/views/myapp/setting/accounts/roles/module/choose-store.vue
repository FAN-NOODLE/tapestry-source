<template>
  <el-dialog
    title="选择角色可查看的门店范围"
    :visible.sync="is_show"
    width="800px"
    @open="$nextTick(()=&gt;openDialog())"
    @close="$nextTick(()=&gt;closeDialog())"
  >
    <div class="dialog-store-body">
      <el-row>
        <el-col :span="24">
          <el-input
            class="store-input"
            placeholder="搜索门店或部门"
            v-model="filter_text"
            size="small"
          ></el-input>
          <el-button class="store-btn" size="small" @click="search"
            >查询</el-button
          >
          <el-button class="store-btn" size="small" @click="reset"
            >重置</el-button
          >
        </el-col>
      </el-row>
      <el-row class="show-tree">
        <el-col :span="24">
          <el-tree
            v-if="is_show_tree"
            :props="props"
            :load="loadNode"
            node-key="dp_id"
            :expand-on-click-node="false"
            :check-strictly="true"
            :default-expanded-keys="first_key"
            lazy="lazy"
            ref="tree2"
            ><span class="custom-tree-node" slot-scope="{ node, data }"
              ><span>{{ node.label }}</span
              ><span
                ><span class="is-store" v-if="data.is_store === 1">门店</span>
                <el-tooltip
                  v-if="data.is_store === 0"
                  effect="dark"
                  content="当前部门不是门店，不可选择"
                  placement="right"
                >
                  <el-checkbox
                    v-model="data.checked"
                    :disabled="data.is_store === 0"
                  ></el-checkbox>
                </el-tooltip>
                <el-checkbox
                  v-if="data.is_store === 1"
                  :checked="data.checked"
                  :disabled="data.is_store === 0"
                  @change="handleCheckChange(data)"
                ></el-checkbox>
              </span>
            </span>
          </el-tree>
        </el-col>
      </el-row>
    </div>
    <div class="dialog-footer" slot="footer">
      <el-button @click="toCancel" size="small">取消</el-button>
      <el-button type="primary" @click.native="toTrue" size="small"
        >确定</el-button
      >
    </div>
  </el-dialog>
</template>

<script>
import Api from "@/api/setting/Role";
import Apis from "@/api/organization/Personnel";
export default {
  name: "ChooseStore",
  props: {
    value: {
      type: Object
    },
    check_list: {
      type: Array
    },
    visible: {
      type: Boolean,
      default() {
        return false;
      }
    }
  },
  data() {
    return {
      is_show: false, // 是否显示
      query: {
        mem_username: "",
        mem_mobile: "",
        dp_id: ""
      },
      store_list: [],
      list: [],
      filter_text: "",
      checked_key: [],
      props: {
        label: "dp_name",
        children: [],
        isLeaf: "leaf"
      },
      depa_list: [],
      is_show_tree: false,
      first_key: [],
      parent_key: null
    };
  },
  watch: {
    visible() {
      this.is_show = true;
      this.store_list = [...this.check_list];
    }
  },
  methods: {
    //动态加载节点
    loadNode(node, resolve) {
      // 一级节点处理
      if (node.level === 0) {
        this.setChecked(this.depa_list);
        this.isChecked(this.checked_key, this.depa_list);
        return resolve(this.depa_list);
      }
      // 其余节点处理
      if (node.level >= 1) {
        this.getIndex(node, resolve);
      }
    },

    // 加载叶子节点数据函数
    getIndex(node, resolve) {
      const self = this;
      const id = node.data.dp_id;
      Api.departmentList({
        dp_id: id
      }).then(res => {
        const data = res.data.list;
        self.setChecked(data);
        self.isChecked(self.checked_key, data);
        resolve(data);
      });
    },

    //获取组织列表
    departmentList(index) {
      const self = this;
      return new Promise((resolve, reject) => {
        if (index === 1) {
          Apis.depPageList({
            dp_id: 0,
            dp_name: self.filter_text,
            page: 1,
            size: 100
          })
            .then(res => {
              self.depa_list = res.data.list;
              resolve();
            })
            .catch(() => {
              reject();
            });
        } else {
          Api.departmentList({
            dp_id: self.parent_key,
            dp_name: self.filter_text
          })
            .then(res => {
              self.depa_list = res.data.list;
              resolve();
            })
            .catch(() => {
              reject();
            });
        }
      });
    },

    //第一次进入时获取组织列表
    firstLists() {
      const self = this;
      return new Promise((resolve, reject) => {
        Api.departmentList({})
          .then(res => {
            self.depa_list = res.data.list;
            self.first_key = [];
            if (res.data.list.length > 0) {
              self.first_key.push(res.data.list[0].dp_id);
            }
            resolve();
          })
          .catch(() => {
            reject();
          });
      });
    },

    // 搜索
    async search() {
      if (this.filter_text == "") {
        this.parent_key = null;
      } else {
        this.parent_key = this.first_key[0];
      }
      this.is_show_tree = false;
      await this.departmentList(1);
      this.$nextTick(() => {
        this.is_show_tree = true;
      });
    },

    // 重置
    async reset() {
      this.parent_key = null;
      this.filter_text = "";
      this.is_show_tree = false;
      await this.departmentList(2);
      this.$nextTick(() => {
        this.is_show_tree = true;
      });
    },

    // 取消
    toCancel() {
      this.is_show = false;
    },

    // 确定
    toTrue() {
      this.$emit("callBack", this.store_list);
      this.is_show = false;
    },

    filterNode(value, data) {
      if (!value) return true;
      return data.dp_name.indexOf(value) !== -1;
    },

    //节点选择改变时触发
    handleCheckChange(data) {
      data.checked = !data.checked;
      if (data.checked) {
        const indexs = this.store_list
          .map(item => item.dp_id)
          .indexOf(data.dp_id);
        if (indexs === -1) {
          this.store_list.push(data);
        }
      } else {
        const index = this.store_list
          .map(item => item.dp_id)
          .indexOf(data.dp_id);
        if (index !== -1) {
          this.store_list.splice(index, 1);
        }
      }
    },

    //弹框打开时触发
    async openDialog() {
      await this.firstLists();
      const temp = [];
      this.store_list.forEach(item => {
        temp.push(item.dp_id);
      });
      this.checked_key = [...temp];
      this.is_show_tree = true;
    },

    //弹框关闭时触发
    closeDialog() {
      this.filter_text = "";
      setTimeout(() => {
        this.is_show_tree = false;
      }, 300);
    },

    //循环为每个部门加上checked属性
    setChecked(arr) {
      arr.forEach(item => {
        if (!item.is_child_department) item.leaf = true;
        item.checked = false;
      });
    },

    //找到选中的节点
    isChecked(arr, arrs) {
      for (let i = arr.length; i--; ) {
        // let data = arrs.find(item => item.dp_id === arr[i]);
        const index = arrs.map(item => item.dp_id).indexOf(arr[i]);
        if (index !== -1) {
          arrs[index].checked = true;
        }
      }
    }
  },
  created() {
    // this.firstLists();
  }
};
</script>

<style lang="sass" scoped>
.tip
  font-size: 12px
  line-height: 14px
  position: absolute
  top: 60px
.select-box
  width: 300px
.el-dialog__body
  padding: 0px 20px
  .dialog-store-body
    margin-top: -30px
  .el-row
    .el-col
      padding: 10px
    .el-form
      padding-bottom: 10px
      .el-form-item
        margin-bottom: 0
  .custom-tree-node
    flex: 1
    display: flex
    align-items: center
    justify-content: space-between
    font-size: 14px
    padding-right: 8px
  .is-store
    margin-right: 50px
  .store-btn
    margin-left: 15px
    display: inline-block
  .store-input
    width: 300px
    display: inline-block
  .show-tree
    max-height: 548px
    overflow: auto
</style>
<style lang="sass">
.dialog-store-body
  .el-tree-node__content
    height: 35px
</style>
