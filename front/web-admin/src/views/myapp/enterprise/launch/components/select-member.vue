
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
              placeholder="请输入部门名称/员工姓名/员工ID/手机号"
              style="width: 300px;"
              clearable
              @keyup.enter.native="handleSearch()"
            ></el-input>
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
        </div>
        <div class="tree">
          <el-tree
            ref="tree"
            :data="treeList" 
            :props="defaultProps"
            :node-key="'id'"
            :default-checked-keys="checkedKeys"
            :default-expand-all="false"
            :show-checkbox="true"
            :check-strictly="true"
            :expand-on-click-node="true"
            :check-on-click-node="false"
            :filter-node-method="filterNode"
            
            @check="handleCheck">

            <template slot-scope="{ node, data }">
              <div class="custom-tree-node">
                <img
                  class="group-img"
                  src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/group-icon.png"
                  alt=""
                  v-if="data.type !==5"
                />
                <i v-else class="el-icon-user-solid"></i>
                <div class="m-l-5">{{ data.name }}</div>
              </div>
            </template>
          </el-tree>
        </div>
      </div>
      <div class="select-member__result">
        <div class="result-title">
          <div>已选({{ selection.length }})</div>
          <el-button v-if="selection.length" type="text" @click="handleDeleteAll"><i class="el-icon-delete"></i>清空已选</el-button>
        </div>
        <div class="result-list">
          <div class="result-item" v-for="(item, index) in selection" :key="index">
            <div class="item-name">
              <img
                class="item-img"
                src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/group-icon.png"
                alt=""
                v-if="item.type !==5"
              />
              <i v-else class="el-icon-user-solid"></i>
              {{ item.name }}
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
// 8. 使用非强制关联父子级
import { Component, Vue, Prop, PropSync, Watch } from "vue-property-decorator";
import API from "@/api/select-member";
import { forEachTree } from '@/utils/tree-tool'
import { getMemberByThirdsApi } from "@/api/enterprise/activity";

@Component({
  name: "SelectMember",
})
export default class SelectMember extends Vue {
  @PropSync('show', { type: Boolean, default: false }) visible: boolean
  @Prop({ type: String, default: "选择" }) title: string
  @Prop({ type: String, default: "member_id" }) rowKey: string
  @Prop({ type: Array, default: ()=>[] }) value: any[]

  isIndeterminate = false

  defaultProps = {
    children: "children",
    label: "name",
    value: "id",
  };
  
  // 当前选中数据
  selection: any[] = []
  /** 默认选中key */
  checkedKeys: number[] = []

  search = {
    type: 'member',
    value: ''
  }

  // 新组织架构树
  treeList:any[] = []

  @Watch('show') 
  handleValueChange(val: any) {
    if(val) {
      const selection:any[]  = []
      const ids = this.value.map(item => item.id)
      this.checkedKeys = ids;
      forEachTree<any>(this.treeList, (row) => {
        if(ids.includes(row.id)) {
          selection.push(row)
        }
      })
      this.selection = selection
      this.handleReset()
    }
  }

  async mounted() {
    await this.getDepartmentTree()
    this.getAllMemberList()
  }

  /**
   * 获取组织架构树
   */
  async getDepartmentTree() {
    const res:any = await API.getDashboardTreeApi({
      has_children: true
    })
    if(res.code === 200) {
      // 获取所有门店thirdId
      this.treeList = res.data
    }
  }
  async getAllMemberList() {

    const thirdIds:number[] = []
    forEachTree<any>(this.treeList, (row) => {
      if(row.third_id) {
        thirdIds.push(row.third_id)
      }
    })
    // 获取门店下的导购
    const res:any = await getMemberByThirdsApi({
      third_ids: thirdIds
    })
    if(res.code === 200) {
      
      const memberList = res.data.map(item => {
        item.type = 5
        return item
      })
      
      forEachTree<any>(this.treeList, (row) => {
        // console.log(row);
        if(row.type === 4) {
          row.children = memberList.filter(item => item.third_id === row.third_id)
        }
      })
    }
  }
  filterNode(value, data) {
    return !value || data.name.indexOf(value) !== -1 || data.mobile === value || data.userid === value;
  }

  async handleCheck(value, node){
    this.selection = node.checkedNodes
  }

  handleRemove(item, index) {
    (this.$refs.tree as any).setChecked(item.id, false);
    this.selection.splice(index, 1);
    if(this.selection.length === 0) {
      (this.$refs.tree as any).setCheckedKeys([])
    }
  }

  handleDeleteAll() {
    this.selection = [];
    (this.$refs.tree as any).setCheckedKeys([]);
  }

  async handleSearch() {
    if(this.search.value.trim() === '') return
    (this.$refs.tree as any).filter(this.search.value)
  }
  handleReset() {
    this.search.value = '';
    // 重置树选中状态
    (this.$refs.tree as any)?.filter('');
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

      .item-img {
        width: 14px;
      }
    }
  }
}
</style>
