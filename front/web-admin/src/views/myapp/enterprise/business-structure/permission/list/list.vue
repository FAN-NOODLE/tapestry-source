<template>
  <layout-list>
    <template slot="buttons">
      <el-button type="primary" size="small" @click="handleAdd">新增职位</el-button>
    </template>
    <template slot="search">
      <search-item text="员工查询：">
        <el-input
          :placeholder="placeholder"
          class="search-input"
          size="medium"
          v-model="searchParams.search_key"
          clearable
        >
          <template slot="prepend">
            <el-select
              v-model="searchParams.search_type"
              class="input-with-select"
              @change="handleSelectChange"
            >
              <el-option
                v-for="item in type_list"
                :key="item.value"
                :value="item.value"
                :label="item.label"
              ></el-option>
            </el-select>
          </template>
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
      <el-table :data="list" stripe>
        <el-table-column prop="id" label="职位ID" width="100">
        </el-table-column>
        <el-table-column prop="name" label="员工名称" width="120">
        </el-table-column>
        <el-table-column prop="userid" label="员工ID" width="140">
        </el-table-column>
        <el-table-column prop="position" label="职位" width="180">
          <template slot-scope="scope">
            <span>{{scope.row.position || '--'}}</span>
          </template>
        </el-table-column>
        <el-table-column prop="path_str" label="管理层级">
          <template slot-scope="scope">
            <span v-if="!scope.row.path_name_arr || scope.row.path_name_arr.length===0">--</span>
            <el-tag class="m-r-10 m-b-10" v-for="(item, index) in scope.row.path_name_arr||[]" :key="index">{{item}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="store_num" label="管理店数" width="100">
        </el-table-column>
        <el-table-column
          prop="zip"
          label="操作"
          width="180"
          fixed="right"
        >
          <template slot-scope="scope">
            <el-button
              type="text"
              @click="handleEdit(scope.row)"
            >编辑</el-button>
            <el-button
              type="text"
              @click="handleView(scope.row)"
            >查看</el-button>
            <el-button
              type="text"
              @click="handleDelete(scope.row)"
            >删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </template>
    <template slot="footer">
      <el-pagination
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="meta.page"
        :page-sizes="[10, 20, 30, 40]"
        :page-size="meta.size"
        layout="total, sizes, prev, pager, next, jumper"
        :total="meta.total"
      >
      </el-pagination>
    </template>
  </layout-list>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LayoutList from "@/components/public/layout-list/index.vue";
import Api from "@/api/enterprise/permission";
import { ApiResult } from "@/types";

@Component({
  components: {
    LayoutList,
  }
})
export default class PermissionList extends Vue {
  /** 筛选类型 */
  type_list = [
    {
      value: 1,
      label: '员工名称'
    },
    {
      value: 2,
      label: '员工ID'
    },
  ];

  /** 筛选条件 */
  searchParams = {
    search_type: 1,
    search_key: ''
  };

  /** 分页参数 */
  meta = {
    page: 1,
    size: 10,
    total: 0
  };

  /** 列表数据 */
  list: any[] = [];
  
  /** placeholder */
  get placeholder() {
    const name = this.type_list.find((item) => item.value === this.searchParams.search_type)?.label
    return '请输入' + name
  }
  
  created() {
    this.getData()
  }

  async getData() {
    try {
      // 请求参数
      const params: any = {
        page: this.meta.page,
        size: this.meta.size,
      }

      // 筛选条件
      if(this.searchParams.search_type === 1) {
        params.user_name = this.searchParams.search_key
      } else if(this.searchParams.search_type === 2) {
        params.user_id = this.searchParams.search_key
      }

      const {
        code,
        message = "",
        data
      }: ApiResult<any> = await Api.getManagerAuthListApi(params)
      if (code !== 200) {
        this.$message(message);
        return false;
      }
      this.list = data.list;
      this.meta.total = data.total;
    } catch (error) {
      console.log(error);
    }
  }

  /**
   * 清空输入内容
   */
  handleSelectChange() {
    this.searchParams.search_key = ''
  }

  /**
   * 查询
   */
  handleSearch() {
    this.meta.page = 1
    this.getData()
  }

  /**
   * 重置
   */
  handleReset() {
    this.searchParams = {
      search_type: 1,
      search_key: ''
    }
    this.handleSearch()
  }

  /**
   * 分页
   * @param val 
   */
  handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.getData();
  }

  /**
   * 分页
   * @param val 
   */
  handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getData();
  }

  handleAdd() {
    this.$router.push({
      name: 'myapp_enterprise_business-structure_permission_add',
    })
  }

  /**
   * 编辑
   * @param index 
   * @param row 
   */
  handleEdit(row: any) {
    this.$router.push({
      name: 'myapp_enterprise_business-structure_permission_edit',
      query: {
        id: row.id
      }
    })
  }
  /**
   * 查看
   * @param index 
   * @param row 
   */
  handleView(row: any) {
    this.$router.push({
      name: 'myapp_enterprise_business-structure_permission_detail',
      query: {
        id: row.id
      }
    })
  }
  /**
   * 删除
   * @param index 
   * @param row 
   */
  handleDelete(row: any) {
    this.$confirm(
      "确定删除该员工的职位，删除后该员工数据看板相关管理者查看权限将失效？",
      "删除提醒",
      {
        confirmButtonText: "确认",
        cancelButtonText: "取消",
        confirmButtonClass: "button-danger",
        type: "warning"
      }
    )
    .then(() => {
      Api.deleteManagerAuthApi({ id: row.id }).then(() => {
        this.$message.success('删除成功!')
        this.getData();
      })
    })
    .catch(() => {
      console.log("取消,什么都不做");
    });
  }
}
</script>
<style lang="scss" scoped>
.search-input {
  width: 280px !important;
}
.input-with-select {
  width: 110px;
}
</style>