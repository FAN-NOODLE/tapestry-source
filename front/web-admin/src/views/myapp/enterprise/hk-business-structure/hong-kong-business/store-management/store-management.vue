<template>
  <div class="store-management">
    <layout-list>
      <template slot="search">
        <div class="sync-time">
          数据同步时间：
          <span v-if="lastSyncTime">
            {{ lastSyncTime | formatDate("yyyy-mm-dd hh:mm:ss") }}
          </span>
          <span v-else>--</span>
        </div>
        <div class="search-row">
          <search-item text="门店信息：">
            <el-input
              v-model.trim="query.search_str"
              size="small"
              clearable
              placeholder="请输入门店名称或编码"
            ></el-input>
          </search-item>
          <search-item text="门店类型：">
            <el-select v-model="query.store_type" size="small">
              <el-option
                v-for="item in storeTypeOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              ></el-option>
            </el-select>
          </search-item>
          <search-item text="区域：">
            <el-select v-model="query.area" size="small">
              <el-option
                v-for="item in areaOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              ></el-option>
            </el-select>
          </search-item>
          <search-item text="门店状态：">
            <el-select v-model="query.closed_status" size="small">
              <el-option label="全部" :value="0"></el-option>
              <el-option label="营业中" :value="1"></el-option>
              <el-option label="已关店" :value="2"></el-option>
            </el-select>
          </search-item>
        </div>
        <div class="action-row">
          <el-button type="primary" size="small" @click="handleSearch">搜索</el-button>
          <el-button size="small" @click="handleReset">重置</el-button>
          <el-button
            type="primary"
            size="small"
            :loading="syncing"
            @click="handleSync"
          >
            同步企微门店
          </el-button>
        </div>
      </template>
      <template slot="table">
        <el-table :data="list" stripe>
          <el-table-column prop="wechat_department_id" label="部门 ID" width="120"></el-table-column>
          <el-table-column prop="store_code" label="门店编号" min-width="130"></el-table-column>
          <el-table-column prop="store_name" label="门店名称" min-width="180"></el-table-column>
          <el-table-column label="区域" width="100">
            <template slot-scope="{ row }">{{ getAreaName(row.area) }}</template>
          </el-table-column>
          <el-table-column prop="inventory_name" label="英文名" min-width="160">
            <template slot-scope="{ row }">{{ row.inventory_name || "--" }}</template>
          </el-table-column>
          <el-table-column label="门店状态" width="100">
            <template slot-scope="{ row }">
              <el-tag v-if="row.closed_at" size="mini" type="danger">已关店</el-tag>
              <el-tag v-else size="mini">营业中</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="门店类型" width="100">
            <template slot-scope="{ row }">{{ getStoreTypeName(row.store_type) }}</template>
          </el-table-column>
          <el-table-column prop="city" label="所属城市" min-width="130">
            <template slot-scope="{ row }">{{ row.city || "--" }}</template>
          </el-table-column>
          <el-table-column prop="email" label="门店邮箱" min-width="170">
            <template slot-scope="{ row }">{{ row.email || "--" }}</template>
          </el-table-column>
          <el-table-column prop="address" label="门店地址" min-width="180">
            <template slot-scope="{ row }">{{ row.address || "--" }}</template>
          </el-table-column>
          <el-table-column fixed="right" label="操作" width="120">
            <template slot-scope="{ row }">
              <el-button type="text" @click="handleEdit(row)">编辑</el-button>
              <el-button
                v-if="!row.closed_at"
                type="text"
                class="danger-button"
                @click="handleCloseStore(row)"
              >
                关店
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
    <store-edit-dialog
      v-model="showEditDialog"
      :store="currentStore"
      @success="getStoreList"
    ></store-edit-dialog>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LayoutList from "@/components/public/layout-list/index.vue";
import StoreEditDialog from "../components/store-edit-dialog.vue";
import Api from "@/api/enterprise/hong-kong-business";

@Component({
  components: { LayoutList, StoreEditDialog }
})
export default class StoreManagement extends Vue {
  query = {
    search_str: "",
    area: "",
    store_type: -1,
    closed_status: 0
  };
  areaOptions = [
    { label: "全部", value: "" },
    { label: "香港", value: "hk" },
    { label: "澳门", value: "mo" }
  ];
  storeTypeOptions = [
    { label: "全部", value: -1 },
    { label: "未设置", value: 0 },
    { label: "奥莱", value: 1 },
    { label: "零售", value: 2 }
  ];
  meta = { page: 1, size: 10, total: 0 };
  list: any[] = [];
  lastSyncTime = 0;
  syncing = false;
  showEditDialog = false;
  currentStore: any = {};

  created() {
    this.getStoreList();
    this.getLastSyncTime();
  }

  getStoreTypeName(value: number) {
    const option = this.storeTypeOptions.find(item => item.value === value);
    return option ? option.label : "--";
  }

  getAreaName(value: string) {
    const option = this.areaOptions.find(item => item.value === value);
    return option && option.value ? option.label : "--";
  }

  async getStoreList() {
    const params: any = {
      ...this.query,
      page: this.meta.page,
      size: this.meta.size
    };
    if (!params.area) {
      delete params.area;
    }
    const res: any = await Api.getStoreList(params);
    if (res.code !== 200) {
      this.$message.error(res.message || "获取门店列表失败");
      return;
    }
    this.list = res.data.list || [];
    this.meta = {
      page: res.data.page || 1,
      size: res.data.size || this.meta.size,
      total: res.data.total || 0
    };
  }

  async getLastSyncTime() {
    const res: any = await Api.getLastSyncTime();
    if (res.code === 200) {
      this.lastSyncTime = res.data.last_sync_time || 0;
    }
  }

  handleSearch() {
    this.meta.page = 1;
    this.getStoreList();
  }

  handleReset() {
    this.query = { search_str: "", area: "", store_type: -1, closed_status: 0 };
    this.handleSearch();
  }

  async handleSync() {
    this.syncing = true;
    try {
      const res: any = await Api.syncStores();
      if (res.code !== 200) {
        this.$message.error(res.message || "同步失败");
        return;
      }
      const data = res.data || {};
      await this.$alert(
        `新增门店：${data.created_count || 0}；未变化：${data.skipped_count || 0}；自动关店：${data.closed_count || 0}`,
        "同步完成",
        { confirmButtonText: "确定" }
      );
      this.meta.page = 1;
      this.getStoreList();
      this.getLastSyncTime();
    } finally {
      this.syncing = false;
    }
  }

  async handleEdit(row: any) {
    const res: any = await Api.getStoreDetail({ id: row.id });
    if (res.code !== 200) {
      this.$message.error(res.message || "获取门店详情失败");
      return;
    }
    this.currentStore = res.data || {};
    this.showEditDialog = true;
  }

  handleCloseStore(row: any) {
    this.$confirm(
      "关店后不可恢复，该门店已有职位权限将保留但不再生效。确认关店吗？",
      "关店提醒",
      {
        confirmButtonText: "确认关店",
        cancelButtonText: "取消",
        confirmButtonClass: "button-danger",
        type: "warning"
      }
    ).then(async () => {
      const res: any = await Api.closeStore({ id: row.id });
      if (res.code !== 200) {
        this.$message.error(res.message || "关店失败");
        return;
      }
      this.$message.success("关店成功");
      this.getStoreList();
    }).catch(() => undefined);
  }

  handleSizeChange(size: number) {
    this.meta.page = 1;
    this.meta.size = size;
    this.getStoreList();
  }

  handleCurrentChange(page: number) {
    this.meta.page = page;
    this.getStoreList();
  }
}
</script>

<style lang="scss" scoped>
.sync-time {
  margin-bottom: 18px;
  color: $color-text-normal;
}

.search-row {
  display: flex;
  align-items: center;
  flex-wrap: wrap;

  /deep/ .el-input,
  /deep/ .el-select {
    width: 220px;
  }
}

.action-row {
  margin-top: 16px;
}

.danger-button {
  color: $color-danger;
}
</style>
