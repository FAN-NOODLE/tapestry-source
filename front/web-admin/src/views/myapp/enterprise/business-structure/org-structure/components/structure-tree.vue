<template>
  <div class="structure-tree">
    <!-- <el-button class="add" type="primary" size="mini" @click="handleAdd">新增第一层级</el-button> -->
    <el-tree 
      ref="tree"
      node-key="id"
      default-expand-all
      highlight-current
      :expand-on-click-node="false"
      :props="defaultProps"
      :data="treeList"
      @node-click="handleNodeClick"
    >
      <template slot-scope="{ node, data }">
        <div class="node">
          <div class="node-name">
            <div class="name">{{ data.name }}</div>
          </div>
          <div v-if="false" @click.stop="()=> {}">
            <el-dropdown v-if="currentNodeKey === data.id" @command="handleCommand($event, data)">
              <span class="node-btn">
                <i class="el-icon-more"></i>
              </span>
              <el-dropdown-menu slot="dropdown">
                <el-dropdown-item command="a">重命名</el-dropdown-item>
                <el-dropdown-item command="b" v-if="data.type < 3">添加子节点</el-dropdown-item>
                <el-dropdown-item command="c">向下添加节点</el-dropdown-item>
                <el-dropdown-item command="d">向上添加节点</el-dropdown-item>
                <!-- <el-dropdown-item command="f" :disabled="data.children && data.children.length>0" >移除节点</el-dropdown-item> -->
              </el-dropdown-menu>
            </el-dropdown>
          </div>
        </div>
      </template>
    </el-tree>
    <structure-add ref="structureAdd" v-model="visible" @confirm="handleConfirm"></structure-add>
  </div>
</template>
<script lang="ts">
import { Component, Prop, Vue, Watch } from "vue-property-decorator";
import Api from '@/api/enterprise/org-structure';
import StructureAdd from "./structure-add.vue";
import { filterTreeNode, findTreeNode, forEachTree } from '@/utils/tree-tool';
export interface ITreeNode {
  children: ITreeNode[];
  id: number;
  name: string;
  parent_id: number;
  parent_name: string;
  store_code: string;
  store_name: string;
  third_id: number;
  type: number;
}
@Component({
  components: {
    StructureAdd
  }
})
export default class extends Vue {
  @Prop({ default: () => [] }) list: ITreeNode[]
  visible = false
  /** 默认属性 */
  defaultProps = {
    isLeaf: "leaf",
    children: "children",
    label: "name",
  }
  /** 当前节点key */
  currentNodeKey = 0;

  currentNode: any = null

  get treeList() {
    return filterTreeNode<ITreeNode>(this.list, node => node.type !== 4, {pid:'parent_id'})
  }

  @Watch("list")
  onListChange(value) {
    const list = JSON.parse(JSON.stringify(value))
    // 设置默认选中节点
    this.$nextTick(() => {
      if(this.currentNodeKey) {
        (this.$refs.tree as any).setCurrentKey(this.currentNodeKey)
        const currentNode = findTreeNode<ITreeNode>(list, node => node.id === this.currentNodeKey)
        this.handleTable(currentNode.children||[], currentNode.name)
      } else {
        this.handleTable(list)
      }
    })
  }

  /**
   * 节点点击
   * @param value 
   */
  handleNodeClick(value: any) {
    this.currentNodeKey = value.id
    const list = JSON.parse(JSON.stringify(this.list))
    const currentNode = findTreeNode<ITreeNode>(list, node => node.id === value.id)
    this.handleTable(currentNode.children||[], value.name)
  }

  /**
   * 操作命令
   * @param value 
   * @param data 
   */
  handleCommand(value: string, data: any) {
    console.log(value, data);
    const params: Record<string, any> = {}
    switch (value) {
      case 'a':
        // 重命名
        params.id = data.id
        params.name = data.name
        break;
      case 'b':
        // 添加子节点
        params.parent_id = data.id
        params.name = ''
        break;
      case 'c':
        // 向下添加节点
        params.brother_id = data.id
        params.name = ''
        params.is_above = false
        break;
      case 'd':
        // 向上添加节点
        params.brother_id = data.id
        params.name = ''
        params.is_above = true
        break;
      case 'f':
        // 向上添加节点
        this.deleteData(data.id) 
        return;
    
      default:
        break;
    }

    (this.$refs.structureAdd as any).show(value, params)
  }

  /**
   * 添加
   */
  handleAdd() {
    (this.$refs.structureAdd as any).show('n', {
      name: ''
    })
  }

  /**
   * 确认
   */
  handleConfirm(value: any) {
    const { key, ...params} = value
    console.log(key, params);
    this.update(params);
  }

  /**
   * 修改
   */
  update(params: any) {
    Api.createDashboardTree(params)
      .then((res: any) => {
        if(res.code === 200) {
          this.visible = false
          this.$emit('update')
        } else {
          this.$message.error(res.message)
        }
      })
  }
  /**
   * 删除
   */
  deleteData(id) {
    this.$confirm('确定删除当前节点吗？', '删除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(() => {
      Api.delDashboardTree({id})
      .then((res: any) => {
        if(res.code === 200) {
          this.currentNodeKey = 0;
          this.$emit('update')
        } else {
          this.$message.error(res.message)
        }
      })
    })
  }


  /**
   * 获取门店列表
   * @param list 
   * @param title 
   */
  handleTable(list?: ITreeNode[], title = '全部门店') {
    const table: ITreeNode[] = []
    forEachTree<ITreeNode>(list, node => {
      if(node.type === 4) {
        table.push(node)
      }
    })
    // table.s
    table.sort((a: ITreeNode, b: ITreeNode) => a.id - b.id)
    this.$emit('table', {
      list: table,
      title
    })
  }

  /**
   * 清空当前节点key
   */
  clearCurrentNodeKey() {
    this.currentNodeKey = 0;
    (this.$refs.tree as any).setCurrentKey(null)
    this.handleTable(JSON.parse(JSON.stringify(this.list)))
  }
}
</script>
<style lang="scss" scoped>

.structure-tree {
  width: 300px;
  margin-right: 20px;
  background-color: #f9f9f9;


  .el-tree {
    background-color: #f9f9f9;
    min-height: 500px;
    max-height: 700px;
    overflow-y: scroll;
  }

  

  .add {
    width: 100%;
  }

  .node {
    display: flex;
    align-items: center;
    width: 100%;
    padding-right: 10px;
    overflow: hidden;

    .node-name {
      flex: 1;
      padding-right: 10px;
      overflow: hidden;
      .name {
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
      }
    }

    .node-btn {
      font-size: 20px;
      color: #f7f7f7;
    }
  }
}
</style>