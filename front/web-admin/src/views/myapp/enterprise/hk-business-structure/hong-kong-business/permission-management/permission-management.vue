<template>
  <layout-list>
    <template slot="buttons">
      <el-button type="primary" size="small" @click="handleAdd">新增职位</el-button>
    </template>
    <template slot="search">
      <search-item text="员工查询：">
        <el-input
          v-model.trim="searchParams.search_key"
          class="search-input"
          size="small"
          clearable
          :placeholder="'请输入' + currentSearchLabel"
        >
          <el-select
            slot="prepend"
            v-model="searchParams.search_type"
            @change="searchParams.search_key = ''"
          >
            <el-option label="员工名称" :value="1"></el-option>
            <el-option label="员工 ID" :value="2"></el-option>
          </el-select>
        </el-input>
      </search-item>
      <el-button type="primary" size="small" @click="handleSearch">查询</el-button>
      <el-button size="small" @click="handleReset">重置</el-button>
    </template>
    <template slot="table">
      <el-table :data="list" stripe>
        <el-table-column prop="id" label="职位 ID" width="100"></el-table-column>
        <el-table-column prop="name" label="员工名称" min-width="120"></el-table-column>
        <el-table-column prop="userid" label="员工 ID" min-width="140"></el-table-column>
        <el-table-column prop="position" label="职位" min-width="190"></el-table-column>
        <el-table-column label="管理范围" min-width="260">
          <template slot-scope="{ row }">
            <span v-if="!row.path_name_arr || !row.path_name_arr.length">--</span>
            <el-tag
              v-for="(name, index) in row.path_name_arr || []"
              :key="index"
              class="scope-tag"
            >
              {{ name }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="store_num" label="管理店数" width="100"></el-table-column>
        <el-table-column fixed="right" label="操作" width="180">
          <template slot-scope="{ row }">
            <el-button type="text" @click="handleEdit(row)">编辑</el-button>
            <el-button type="text" @click="handleView(row)">查看</el-button>
            <el-button type="text" class="danger-button" @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </template>
    <template slot="footer">
      <el-pagination
        v-if="meta.total"
        :current-page="meta.page"
        :page-size="meta.size"
        :page-sizes="[10, 20, 30, 40]"
        :total="meta.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      ></el-pagination>
    </template>
  </layout-list>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LayoutList from "@/components/public/layout-list/index.vue";
import Api from "@/api/enterprise/hong-kong-business";

@Component({ components: { LayoutList } })
export default class PermissionManagement extends Vue {
  searchParams = { search_type: 1, search_key: "" };
  meta = { page: 1, size: 10, total: 0 };
  list: any[] = [];

  get currentSearchLabel() {
    return this.searchParams.search_type === 1 ? "员工名称" : "员工 ID";
  }

  created() {
    this.getList();
  }

  async getList() {
    const params: any = { page: this.meta.page, size: this.meta.size };
    if (this.searchParams.search_type === 1) {
      params.user_name = this.searchParams.search_key;
    } else {
      params.user_id = this.searchParams.search_key;
    }
    const res: any = await Api.getPermissionList(params);
    if (res.code !== 200) {
      this.$message.error(res.message || "获取职位权限列表失败");
      return;
    }
    this.list = res.data.list || [];
    this.meta = {
      page: res.data.page || 1,
      size: res.data.size || this.meta.size,
      total: res.data.total || 0
    };
  }

  handleSearch() {
    this.meta.page = 1;
    this.getList();
  }

  handleReset() {
    this.searchParams = { search_type: 1, search_key: "" };
    this.handleSearch();
  }

  handleAdd() {
    this.$router.push({
      name: "myapp_enterprise_hk-business-structure_hong-kong-business_permission-add"
    });
  }

  handleEdit(row: any) {
    this.$router.push({
      name: "myapp_enterprise_hk-business-structure_hong-kong-business_permission-edit",
      query: { id: row.id }
    });
  }

  handleView(row: any) {
    this.$router.push({
      name: "myapp_enterprise_hk-business-structure_hong-kong-business_permission-detail",
      query: { id: row.id }
    });
  }

  handleDelete(row: any) {
    this.$confirm(
      "确定删除该员工的职位权限吗？删除后相关管理权限将失效。",
      "删除提醒",
      {
        confirmButtonText: "确认",
        cancelButtonText: "取消",
        confirmButtonClass: "button-danger",
        type: "warning"
      }
    ).then(async () => {
      const res: any = await Api.deletePermission({ id: row.id });
      if (res.code !== 200) {
        this.$message.error(res.message || "删除失败");
        return;
      }
      this.$message.success("删除成功");
      if (this.list.length === 1 && this.meta.page > 1) {
        this.meta.page -= 1;
      }
      this.getList();
    }).catch(() => undefined);
  }

  handleSizeChange(size: number) {
    this.meta.page = 1;
    this.meta.size = size;
    this.getList();
  }

  handleCurrentChange(page: number) {
    this.meta.page = page;
    this.getList();
  }
}
</script>

<style lang="scss" scoped>
.search-input {
  width: 340px;

  /deep/ .el-input-group__prepend {
    width: 90px;
  }
}

.scope-tag {
  margin: 0 8px 8px 0;
}

.danger-button {
  color: $color-danger;
}
</style>
