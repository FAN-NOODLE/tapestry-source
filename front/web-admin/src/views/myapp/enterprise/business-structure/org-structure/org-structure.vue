<template>
  <layout-list>
    <template slot="search">
      <search-item text="门店查询：">
        <el-input
          :placeholder="placeholder"
          class="search-input"
          size="medium"
          v-model.trim="searchParams.keyword"
          clearable
        >
          <el-select
            slot="prepend"
            v-model="searchParams.type"
            class="input-with-select"
            @change="handleSelectChange"
          >
            <el-option
              v-for="item in type_list"
              :key="item.key"
              :value="item.key"
              :label="item.name"
            ></el-option>
          </el-select>
        </el-input>
      </search-item>
      <el-button size="small" type="primary" @click="handleSearch">
        查询
      </el-button>
      <el-button class="mar-l-20" size="small" @click="handleReset">
        重置
      </el-button>
    </template>
    <template slot="table">
      <div class="wrap">
        <structure-tree ref="tree" :list="treeList" @table="handleTableChange" @update="getData"></structure-tree>
        <structure-table :title="title" :treeList="treeList" :data="listData" @update="getData"></structure-table>
      </div>
    </template>
  </layout-list>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LayoutList from "@/components/public/layout-list/index.vue";
import StructureTable from './components/structure-table.vue'
import StructureTree, { ITreeNode } from './components/structure-tree.vue'
import Api from '@/api/enterprise/org-structure';

@Component({
  components: {
    LayoutList,
    StructureTable,
    StructureTree
  }
})
export default class extends Vue {
  
  /** 筛选类型 */
  type_list = [
    {
      key: 1,
      name: '门店编号'
    },
    {
      key: 2,
      name: '门店名称'
    },
  ];

  /** 筛选条件 */
  searchParams = {
    type: 1,
    keyword: ''
  };

  searchKeyword = ''
  searchType = 1

  title = ''

  treeList: ITreeNode[] = []

  /** 列表数据 */
  list: ITreeNode[] = [];
  
  get placeholder() {
    return '请输入'
  }

  get listData() {
    return this.list.filter(item => {
      if(this.searchType === 1) {
        return item.store_code.includes(this.searchKeyword)
      } else {
        return item.store_name.includes(this.searchKeyword)
      }
    })
  }


  created() {
    this.getData()
  }

  /**
   * 获取数据
   */
  getData() {
    Api.getDashboardTree({
      has_children: true
    }).then((res: any) => {
      this.treeList = res.data
    })
  }

  /**
   * 清空输入内容
   */
  handleSelectChange() {
    this.searchKeyword = ''
    this.searchType = this.searchParams.type
  }

  /** 
   * 表格数据
   */
  handleTableChange({list, title}) {
    this.list = list
    this.title = title
  }

  /** 
   * 搜索
   */
  handleSearch() {
    this.searchKeyword = this.searchParams.keyword;
    this.searchType = this.searchParams.type;
    (this.$refs.tree as any).clearCurrentNodeKey();
    this.getData()
  }

  handleReset() {
    this.searchKeyword = ''
    this.searchType = 1
    this.searchParams = {
      type: 1,
      keyword: ''
    }
    this.getData()
  }
}
</script>
<style lang="scss" scoped>
::v-deep {
  .search-item {
    margin-bottom: 0;
  }
}
.search-input {
  width: 280px !important;
}
.input-with-select {
  width: 110px;
}
.wrap {
  display: flex;
}
</style>