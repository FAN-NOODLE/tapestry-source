
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
              @keyup.enter.native="handleSearch()"
              clearable
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
          <!-- <div v-if="!isSearchMember" class="cate-store">
            <div class="cate-item">门店分类</div>
            <div
              class="cate-item m-l-10"
              v-for="(item, index) in select_cate"
              :key="index"
            >
              <el-checkbox
                v-if="index !== 0"
                @change="handleSelectByCate(item)"
                :indeterminate="item.isIndeterminate"
                v-model="item.checked"
                >{{ item.name }}</el-checkbox
              >
            </div>
          </div> -->
        </div>
        <div class="tree">
          <el-tree
            ref="tree"
            :data="departmentTree" 
            :props="defaultProps"
            :node-key="'value'"
            :default-checked-keys="checkedKeys"
            :default-expanded-keys="[thirdIds[0]]"
            :show-checkbox="true"
            :expand-on-click-node="true"
            :check-on-click-node="false"
            :check-strictly="true"
            :filter-node-method="filterNode"
            @check="handleCheck">

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
          <div>已选部门/员工({{ selection.length }})</div>
          <el-button v-if="selection.length" type="text" @click="handleDeleteAll"><i class="el-icon-delete"></i>清空已选</el-button>
        </div>
        <div class="result-list">
          <div class="result-item" v-for="(item, index) in selection" :key="index">
            <div class="result-item__head">
              <div class="item-name">
                <img
                  class="item-img"
                  src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/group-icon.png"
                  alt=""
                  v-if="item.is_dp"
                />
                <i v-else class="el-icon-user-solid"></i>
                {{ item.label }}
              </div>
              <div class="item-remove" @click="handleRemove(item, index)">
                <i class="el-icon-delete"></i>
              </div>
            </div>
            <div class="result-item__body"> 
              <div class="result-item" v-for="m in selectDpMap[item.value]" :key="m.value">
                <div class="result-item__head">
                  <div class="item-name">
                    <i class="el-icon-user-solid"></i>
                    {{ m.label }}
                  </div>
                </div>
              </div>
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

@Component({
  name: "SelectMember",
})
export default class SelectMember extends Vue {
  @PropSync('show', { type: Boolean, default: false }) visible: boolean
  @Prop({ type: String, default: "选择部门/员工" }) title: string
  @Prop({ type: String, default: "member_id" }) rowKey: string
  @Prop({ type: Array, default: ()=>[] }) value: any[]

  isIndeterminate = false
  

  defaultProps = {
    children: "children",
    label: "label",
    value: "value",
    isLeaf: function(data,node) {
      return data.leaf
    },
  };
  // 组织架构树
  dpIds: any[] = []
  thirdIds: any[] = []
  dpThirdMap: any = {}
  dpIdNameMap: any = {}
  // 组织架构树
  departmentTree: any[] = []
  // 当前选中数据
  selection: any[] = []
  /** 默认选中key */
  checkedKeys: number[] = []
  /** 获取组织架构树 */
  selectDpMap = {}
  memberList: any[] = []

  search = {
    type: 'member',
    value: ''
  }
  searchList: any[] =[]
  select_cate: Array<any> = [
    {
      name: "其他",
      store_type: 0,
      checked: false,
      isIndeterminate: false,
      dpIds: [],
      memberIds: []
    },
    {
      name: "奥莱",
      store_type: 1,
      checked: false,
      isIndeterminate: false,
      dpIds: [],
      memberIds: []
    },
    {
      name: "零售",
      checked: false,
      isIndeterminate: false,
      store_type: 2,
      dpIds: [],
      memberIds: []
    }
  ]; // 门店门类

  isSearchMember = false

  @Watch('show') 
  handleValueChange(val: any) {
    if(val) {
      this.selection = this.value
      const memberIds = this.value.map(item => item.value)
      this.checkedKeys = memberIds;
      this.$nextTick(() => {
        (this.$refs.tree as any).setCheckedKeys(memberIds)
      })
      // 重置搜索条件
      this.handleReset()
      this.getDpMember()
    }
  }

  async mounted() {
    await this.getDepartmentTree()
    await this.getAllMemberList()
    this.setCateStore()
  }

  /** 设置分类状态 */
  setCateStore() {
    forEachTree<any>(this.departmentTree, (row) => {
      if(row.is_dp) {
        this.select_cate[row.store_type].dpIds.push(row.value)
      } else {
        this.select_cate[row.store_type].memberIds.push(row.value)
      }
    })
  }
  updateCateChecked() {
    const selectionIds = this.selection.map(item => item.value)
    this.select_cate.forEach(item => {
      if(selectionIds.length) {
        item.checked = selectionIds.every(id => item.memberIds.includes(id) || item.dpIds.includes(id))
        item.isIndeterminate = selectionIds.some(id => item.memberIds.includes(id) || item.dpIds.includes(id))
        console.log(item);
      } else {
        item.checked = false
        item.isIndeterminate = false
      }
    })
  }
  /** 通过分类选择数据 */
  async handleSelectByCate(item) {
    item.isIndeterminate = false
    // 选中员工
    item.memberIds.forEach(id => {
      (this.$refs.tree as any).setChecked(id, item.checked)
    });
    // 选中门店
    item.dpIds.forEach(id => {
      (this.$refs.tree as any).setChecked(id, item.checked)
    });
    // 更新选中结果
    this.selection = (this.$refs.tree as any).getCheckedNodes().filter(item => !item.is_dp);
  }
  /**
   * 获取组织架构树
   */
  async getDepartmentTree() {
    const res:any = await API.getDepartmentTreeApi({})
    if(res.code === 200) {
      const dpIds:any[] = []
      const thirdIds:any[] = []
      forEachTree<any>(res.data.list, (row) => {
        // console.log(row);
        row.label = row.dp_store_name
        row.search_name = row.search_name ||row.dp_store_name
        row.children = row.children_list
        row.value = row.third_id
        row.is_dp = true
        row.leaf = false
        // 便利子节点给每个子节点增加search_name
        row.children.forEach(child => {
          child.search_name = child.search_name || row.dp_store_name + child.dp_store_name
        })

        // 记录所有的部门ID
        dpIds.push(row.department_id)
        thirdIds.push(row.third_id)
        this.dpThirdMap[row.department_id] = row.third_id
        this.dpIdNameMap[row.third_id] = row.search_name
      })
      this.dpIds = dpIds
      this.thirdIds = thirdIds
      this.departmentTree = res.data.list
    }
  }
  async getAllMemberList() {
    const result = await this.getMemberList(this.dpIds)
    this.memberList = result
    forEachTree<any>(this.departmentTree, (row) => {
      if(row.is_dp) {
        const members = result.filter(item => item.department_id === row.department_id).map(item => {
          return {
            store_type: row.store_type,
            search_name: row.search_name + item.user_name,
            ...item
          }
        })
        row.children = [...row.children, ...members]
        row.leaf = row.children.length === 0
      }
    })
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
          label: item.user_name + `(${[item.user_id,item.position].filter(Boolean).join('，')})`,
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
  filterNode(value, data) {
    // if(!data.is_dp) {
    //   const [dpId] = Object.entries(this.dpIdNameMap).find(item => (item[1] as string).toLowerCase().indexOf(value.toLowerCase()) !== -1)||[]
      
    //   if(dpId && dpId == data.department_id) {
    //     return true
    //   }
    // }
    // 包含部门名称、用户名、用户ID、手机号
    return String(data.search_name||'').toLowerCase().indexOf(value.toLowerCase()) !== -1 || data.user_id === value || data.mobile === value
  }

  async handleCheck(e, node){
    console.log(e, node);
    const selectionIds = this.selection.map(item => item.value)
    if(selectionIds.includes(e.value)) {
      this.selection = this.selection.filter(item => item.value !== e.value)
    } else {
      this.selection.push(e)
    }
    
    
    this.getDpMember()
    this.updateCateChecked()
  }

  handleRemove(item, index) {
    (this.$refs.tree as any).setChecked(item.value, false);
    this.selection.splice(index, 1);
    if(this.selection.length === 0) {
      (this.$refs.tree as any).setCheckedKeys([])
    }
    this.updateCateChecked()
    if(this.selectDpMap[item.value]) {
      delete this.selectDpMap[item.value]
    }
  }

  handleDeleteAll() {
    this.selection = [];
    (this.$refs.tree as any).setCheckedKeys([]);
    this.selectDpMap = {}
    this.updateCateChecked()
  }

  async handleSearch() {
    // if(this.search.value.trim() === '') return
    (this.$refs.tree as any).filter(this.search.value)
  }
  handleReset() {
    this.search.value = '';
    // 重置树选中状态
    this.$nextTick(() => {
      (this.$refs.tree as any).filter('');
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

  getDpMember() {
    const selectionIds = this.selection.map(v => v.value)
    const salesPosition = ['Sr Sales Associate', 'Sales Associate', 'Senior Sales', 'Senior Sales Associate', 'Seasonal Sales Associate', 'eSA'].map(item => item.toLowerCase())
    this.selection.forEach(item => {
      if(item.is_dp) {
        // 获取当前部门下面的人，但是要排除no sales 以及已选的人
        const members = this.memberList.filter(member => {
          return salesPosition.includes(member.position.toLowerCase()) && item.value === this.dpThirdMap[member.department_id]
        })
        this.selectDpMap[item.value] = members
      }
    })
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

        .result-item__head {
          height: 40px;
          border-bottom: 1px solid #f3f4f5;
          display: flex;
          justify-content: space-between;
          align-items: center;
        }

        .item-img {
          width: 14px;
        }

        .result-item__body {
          padding-left: 20px;
          background: #f6f7f8;
        }
      }
    }
  }
}
</style>
