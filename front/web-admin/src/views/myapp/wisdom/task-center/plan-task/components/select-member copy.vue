
<template>
  <el-dialog
    class="dialog"
    :title="title"
    :visible.sync="visible"
    width="1000px"
  >
    <div class="select-member">
      <div class="select-member__content">
        <div class="sticky">
          <!-- 搜索条件 -->
          <div class="search">
            <el-input
              size="medium"
              v-model="search.value"
              :placeholder="`请输入${placeholder}`"
              style="width: 300px;"
              clearable
            >
              <el-select
                v-model="search.type"
                slot="prepend"
                style="width: 120px;"
              >
                <el-option
                  v-for="item in selectDefault"
                  :label="item.label"
                  :value="item.value"
                  :key="item.value"
                ></el-option>
              </el-select>
            </el-input>
            <el-button
              class="m-l-10"
              size="medium"
              type="primary"
              @click="handleSearch()"
            >查询</el-button>
            <el-button
              size="medium"
              class="m-l-10"
              @click="handleReset()"
            >重置</el-button>
          </div>
          <div v-if="!isSearchMember" class="cate-store">
            <div class="cate-item">门店分类</div>
            <div
              class="cate-item m-l-10"
              v-for="(item, index) in select_cate"
              :key="index"
            >
              <el-checkbox
                v-if="index !== 0"
                :disabled="item.store_list.length == 0"
                @change="handleSelectByCate(item)"
                :indeterminate="item.indeterminate"
                v-model="item.checked"
                >{{ item.name }}</el-checkbox
              >
            </div>
          </div>
        </div>
        <!-- 列表 -->
        <div v-if="isSearchMember" class="tree">
          <el-tree
            ref="tree"
            :data="searchList" 
            :props="defaultProps"
            :node-key="'value'"
            :default-checked-keys="checkedKeys"
            :show-checkbox="true"
            :expand-on-click-node="false"
            :check-on-click-node="false"
            @check-change="handleCheckChange">
            <template slot-scope="{ node, data }">
              <div class="custom-tree-node">
                <i class="el-icon-user-solid"></i>
                <div class="m-l-10">{{ data.label }}</div>
              </div>
            </template>
          </el-tree>
        </div>
        <div v-else class="tree">
          <el-tree
            ref="tree"
            :data="departmentTree" 
            :props="defaultProps"
            :node-key="'value'"
            :default-checked-keys="checkedKeys"
            :show-checkbox="true"
            :expand-on-click-node="false"
            :check-on-click-node="false"
            :load="loadNode"
            :filter-node-method="filterNode"
            @check="handleCheck"
            @node-click="handleNodeClick"
            @node-expand="handleNodeExpand">

            <template slot-scope="{ node, data }">
              <div class="custom-tree-node">
                <img
                  class="group-img"
                  src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/group-icon.png"
                  alt=""
                  v-if="data.is_dp"
                />
                <i v-else class="el-icon-user-solid"></i>
                <div class="m-l-5">{{ data.label }}</div>
              </div>
            </template>
          </el-tree>
        </div>
      </div>
      <div class="select-member__result">
        <div class="result-title">
          <div>已选人员({{ selection.length }})</div>
          <el-button v-if="selection.length" type="text" @click="handleDeleteAll"><i class="el-icon-delete"></i>清空已选</el-button>
        </div>
        <div class="result-list">
          <div class="result-item" v-for="(item, index) in selection" :key="index">
            <div class="item-name">
              <i class="el-icon-user-solid"></i>
              {{ item.label }}
            </div>
            <div class="item-remove" @click="handleRemove(item, index)">
              <i class="el-icon-delete"></i>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div slot="footer" class="dialog-footer">
      <el-button @click="handleCancel">取消</el-button>
      <el-button type="primary" @click="handleConfirm">确定</el-button>
    </div>
  </el-dialog>
</template>

<script lang="ts">
// 1. 人员选择
// 2. 实现零售和奥莱的快速选择
// 3. 实现单选与多选
// 4. 实现部门快速定位查询
// 5. 实现导购姓名快速查询
// 6. 实现数据反选
// 7. 实现已选数据展示
import { Component, Vue, Prop, PropSync, Watch } from "vue-property-decorator";
import API from "@/api/select-member";
import { forEachTree } from '@/utils/tree-tool'
import { debounce } from 'lodash'

@Component({
  name: "SelectMember",
})
export default class SelectMember extends Vue {
  @PropSync('show', { type: Boolean, default: false }) visible: boolean
  @Prop({ type: String, default: "选择员工" }) title: string
  @Prop({ type: String, default: "member_id" }) rowKey: string
  @Prop({ type: Array, default: ()=>[] }) value: any[]
  
  selectDefault = [
    { label: "员工ID", value: "userid" },
    { label: "员工手机号", value: "mobile" },
    { label: "员工姓名", value: "mem_username" },
    { label: "部门名称", value: "department" },
  ];

  defaultProps = {
    children: "children",
    label: "label",
    value: "value",
    isLeaf: function(data,node) {
      return data.leaf
    },
  };
  // 组织架构树
  departmentTree: any[] = []
  // 当前选中数据
  selection: any[] = []
  /** 默认选中key */
  checkedKeys: number[] = []

  search = {
    type: 'userid',
    value: ''
  }
  searchList: any[] =[]
  select_cate: Array<any> = [
    {
      name: "其他",
      checked: false,
      store_list: []
    },
    {
      name: "奥莱",
      checked: false,
      store_list: []
    },
    {
      name: "零售",
      checked: false,
      store_list: []
    }
  ]; // 门店门类

  isSearchMember = false

  get placeholder() {
    return this.selectDefault.find(item => item.value === this.search.type)?.label || ''
  }

  @Watch('value') 
  handleValueChange(val: any) {
    this.selection = val
    const memberIds = val.map(item => item.value)
    this.checkedKeys = memberIds
  }
  @Watch('selection', {deep: true}) 
  handleSelectionChange(val: any) {
    const dpIds = [...new Set(val.map(item => item.department_id))]
    console.log(dpIds, this.select_cate);
    if(dpIds.length === 0) {
      this.select_cate = this.select_cate.map(item => {
        item.checked = false
        item.store_list.forEach(store => {
          (this.$refs.tree as any).setChecked(store.department_id, false);
        })
        return item
      })
    }
    
  }

  mounted() {
    this.selection = this.value
    const memberIds = this.value.map(item => item.value)
    this.checkedKeys = memberIds
    
    this.getDepartmentTree()
  }

  /** 设置分类状态 */
  setCateStore(list) {
    list.forEach(item => {
      this.select_cate[item.store_type].store_list.push(item);
      if(item.children_list) {
        this.setCateStore(item.children_list);
      }
    });
    this.select_cate.forEach(item => {
      item.checked = item.store_list.every(store => this.checkedKeys.length>0 && this.checkedKeys.includes(store.department_id))
    })
  }
  /** 通过分类选择数据 */
  async handleSelectByCate(item) {
    const ids = item.store_list.map(i => i.value);
    console.log(item.checked);
    
    ids.forEach(id => {
      (this.$refs.tree as any).setChecked(id, item.checked)
    });

    this.handleCheck()
  }
  /**
   * 获取组织架构树
   */
  async getDepartmentTree() {
    const res:any = await API.getDepartmentTreeApi({})
    if(res.code === 200) {
      forEachTree<any>(res.data.list, (row) => {
        row.label = row.dp_store_name
        row.children = row.children_list
        row.value = row.department_id
        row.is_dp = true
        row.leaf = false
      })
      this.departmentTree = res.data.list
      this.setCateStore(res.data.list)
    }
  }
  /**
   * 获取成员列表
   * @param ids 
   */
  async getMemberList(ids) {
    const res:any = await API.dpMembers({
      departments: ids,
    })
    if (res.code === 200) {
      return res.data.map(item => {
        return {
          ...item,
          value: item[this.rowKey],
          label: item.user_name,
          is_dp: false,
          leaf: true
        }
      })
    }
  }
  /**
   * 获取成员列表
   * @param ids 
   */
  async getMembers(search) {
    const params = {
      page: 1, 
      size: 500,
      search: search
    }
    const res: any = await API.treePeopleNew(params)
    if (res.code === 200) {
      return (res.data.list||[]).map(item => {
        return {
          ...item,
          value: item[this.rowKey],
          label: item.user_name,
          is_dp: false,
          leaf: true
        }
      })
    }
  }
  /** 
   * 懒加载节点
   */
  loadNode(node, resolve) {
    console.log(node);
    resolve(node.data.children || node.data)
  }
  filterNode(value, data) {
    return data.label.indexOf(value) !== -1
  }
  /** 
   * 节点展开
   */
  async handleNodeExpand(data, node) {
    
    const ids = node.childNodes.map(item => item.data.value)

    if(data.is_dp) {
      const list = await this.getMembers({
        dp_id: data.value
      })
      if(list.length) {
        list.forEach((item: any) => {
          item.dp_parent_id = data.id;
          item.label = item.user_name; //  + (item.mobile? `(${item.mobile})`:'')
          item.value = item[this.rowKey];
          item.leaf = true;
          item.is_dp = false;
          if(!ids.includes(item.value)) {
            // 插入节点
            (this.$refs.tree as any).append(item, node);
            // 设置选中
            // 获取默认选中状态
            if(node.checked || this.checkedKeys.includes(item.value)) {
              (this.$refs.tree as any).setChecked(item.value, true);
            }
          }
        });
        if(!node.isLeaf) {
          node.expanded = true
        }
      } else {
        node.isLeaf = true;
      }
    }
  }

  handleNodeClick(data, node) {
    if(node.expanded) {
      node.expanded = false
      return
    }
    this.handleNodeExpand(data, node)
  }

  async handleCheck(){
    const checkedNodes = (this.$refs.tree as any).getCheckedNodes();
    const dps = checkedNodes.filter(item => item.is_dp);
    const members = checkedNodes.filter(item => !item.is_dp);
    const dpIds = dps.map(item => item.value);
    if(dpIds.length) {
      const result = await this.getMemberList(dpIds)
      const resultIds = result.map(item => item.value);
      members.forEach(item => {
        if(!resultIds.includes(item.value)) {
          result.push(item);
        }
      });
      this.selection = result
    } else {
      this.selection = members;
    }
    
  }

  handleCheckChange(data, checked) {
    if(checked) {
      this.selection.push(data);
    } else {
      this.selection = this.selection.filter(item => item.value !== data.value);
    }
    
  }

  handleRemove(item, index) {
    (this.$refs.tree as any).setChecked(item.value, false);
    (this.$refs.tree as any).setChecked(item.department_id, false);
    this.selection.splice(index, 1);
  }

  handleDeleteAll() {
    this.selection.forEach(item => {
      (this.$refs.tree as any).setChecked(item.value, false);
      (this.$refs.tree as any).setChecked(item.department_id, false);
    })
    this.selection = [];
  }

  async handleSearch() {
    if(this.search.value.trim() === '') return
    if(this.search.type !== 'department') {
      const list = await this.getMembers({
        [this.search.type]: this.search.value
      })
      this.searchList = list
      this.isSearchMember = true
      const ids = this.selection.map(item => item.value)
      this.$nextTick(() => {
        // 设置选中项
        list.forEach(item => {
          console.log(ids.includes(item.value));
          if(ids.includes(item.value)) {
            (this.$refs.tree as any).setChecked(item, true);
          }
        });
      })
    } else {
      this.isSearchMember = false;
      (this.$refs.tree as any).filter(this.search.value)
    }
  }
  handleReset() {
    this.search = {
      type: 'userid',
      value: ''
    }
    this.searchList = []
    this.isSearchMember = false;
    // 重置树选中状态
    this.$nextTick(() => {
      (this.$refs.tree as any).filter('');
      const dpIds = [...new Set(this.selection.map(item => item.department_id))]
      
      // departmentTree
      forEachTree<any>(this.departmentTree, (row) => {
        console.log(row);
        if(dpIds.includes(row.department_id)) {
          (this.$refs.tree as any).setChecked(row, true)
        }
      })
    })
  }

  /**
   * 确认
   */
  handleConfirm() {
    this.$emit('confirm', this.selection)
  }

  /**
   * 取消
   */
  handleCancel() {
    this.$emit('cancel')
  }
}
</script>

<style lang="scss" scoped>
#crs-scrm .dialog ::v-deep .el-dialog {
  max-height: none !important;
}

.select-member {
  height: 60vh;
  display: flex;

  .select-member__content {
    flex: 1;
    height: 100%;
    overflow-y: scroll;
    
    .sticky {
      position: sticky;
      top: -1px;
      z-index: 10;
      background-color: #fff;
    }

  }
  .cate-store {
    display: flex;
    margin-top: 15px;
  }
  .custom-tree-node {
    display: flex;
    align-items: center;
    padding-left: 10px;

    .group-img {
      width: 14px;
    }
  }
  .select-member__result {
    width: 400px;
    height: 100%;
    overflow-y: scroll;
    border-left: 1px solid #f3f4f5;
    .result-title {
      position: sticky;
      top: -1px;
      background: #fff;
      padding: 10px;
      z-index: 1;
      display: flex;
      justify-content: space-between;
    }
    .result-list {
      padding: 0 10px;
      .result-item {
        height: 40px;
        border-bottom: 1px solid #f3f4f5;
        display: flex;
        justify-content: space-between;
        align-items: center;
      }
    }
  }
}
</style>
