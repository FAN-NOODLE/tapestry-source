<template>
  <layout-list>
    <template slot="search">
      <el-form :inline="true" class="form-inline">
        <el-form-item label="门店:">
          <el-input
            v-model="listQuery.store_name"
            placeholder="请输入门店名称、编号"
            clearable
          ></el-input>
        </el-form-item>
        <el-form-item label="部门:">
          <el-select
            v-model="listQuery.dept_name"
            placeholder="请选择部门"
            clearable
          >
            <el-option label="全部" value=""></el-option>
            <el-option
              v-for="(item, index) in deptList"
              :key="index"
              :label="item"
              :value="item"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="状态:">
          <el-select v-model="listQuery.status" placeholder="请选择状态">
            <el-option label="全部" :value="0"></el-option>
            <el-option label="巡店预约" :value="1"></el-option>
            <el-option label="巡店中" :value="2"></el-option>
            <el-option label="巡店完成" :value="3"></el-option>
            <el-option label="关闭" :value="4"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="创建时间:">
          <el-date-picker
            v-model="dateRange"
            type="datetimerange"
            format="yyyy-MM-dd HH:mm:ss"
            value-format="timestamp"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            :default-time="['00:00:00', '23:59:59']"
            @change="handleDateChange"
          ></el-date-picker>
        </el-form-item>
        <el-form-item label="巡店时间:">
          <el-date-picker
            v-model="visitDateRange"
            type="datetimerange"
            format="yyyy-MM-dd HH:mm:ss"
            value-format="timestamp"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            :default-time="['00:00:00', '23:59:59']"
            @change="handleVisitDateChange"
          ></el-date-picker>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
          <el-button @click="handleExport">导出</el-button>
        </el-form-item>
      </el-form>
    </template>
    <template slot="table">
      <el-table
        :data="list"
        element-loading-text="加载中..."
        stripe
        highlight-current-row
      >
        <el-table-column
          label="巡店记录ID"
          min-width="150"
          prop="id"
          show-overflow-tooltip
        ></el-table-column>
        <el-table-column
          label="巡店门店"
          min-width="180"
          prop="store_name"
          show-overflow-tooltip
        ></el-table-column>
        <el-table-column
          label="门店编号"
          min-width="120"
          prop="store_code"
          show-overflow-tooltip
        ></el-table-column>
        <el-table-column label="巡店时间" min-width="180">
          <template slot-scope="{ row }">
            <span>{{ row.date }} {{ row.time_span }}</span>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" min-width="180">
          <template slot-scope="{ row }">
            <span>{{ row.created_at | formatDate("yyyy-mm-dd hh:mm:ss") }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" min-width="120">
          <template slot-scope="{ row }">
            <div class="tag" :class="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </div>
          </template>
        </el-table-column>
        <el-table-column label="完成时间" min-width="180">
          <template slot-scope="{ row }">
            <span
              v-if="row.updated_at && (row.status === 3 || row.status === 4)"
              >{{ row.updated_at | formatDate("yyyy-mm-dd hh:mm:ss") }}</span
            >
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="100">
          <template slot-scope="{ row }">
            <!-- <el-button type="text" size="small" @click="handleDetail(row)"
              >查看</el-button
            > -->
            <el-button
              v-if="row.status !== 3 && row.status !== 4"
              type="text"
              size="small"
              @click="handleChangeStatus(row)"
              >关闭</el-button
            >
          </template>
        </el-table-column>
      </el-table>
    </template>
    <template slot="footer">
      <el-pagination
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="listQuery.page"
        :page-sizes="[10, 20, 30, 50]"
        :page-size="listQuery.size"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
      ></el-pagination>
    </template>
  </layout-list>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LayoutList from "@/components/public/layout-list/index.vue";
import StoreVisitApi from "@/api/wisdom/StoreVisit";
import { cosUploader } from '@/utils/upload';
import { downloadFile } from '@/utils/download';
import dayjs from "dayjs";

@Component({
  components: {
    LayoutList,
  },
})
export default class StoreVisitList extends Vue {
  list: any[] = [];
  total = 0;
  deptList: string[] = [];
  dateRange: any = null;
  visitDateRange: any = null;
  listQuery = {
    page: 1,
    size: 10,
    dept_name: "",
    store_name: "",
    start_time: 0,
    end_time: 0,
    status: 0,
    visit_start_time: 0,
    visit_end_time: 0,
  };

  created() {
    this.getDeptList();
    this.getList();
  }

  /**
   * 获取部门列表
   */
  getDeptList() {
    StoreVisitApi.getDeptNameList({}).then((res: any) => {
      if (res.code === 200) {
        this.deptList = res.data || [];
      }
    });
  }

  /**
   * 获取巡店列表
   */
  getList() {
    StoreVisitApi.getStoreVisitList(this.listQuery)
      .then((res: any) => {
        if (res.code === 200) {
          this.list = res.data.list || [];
          this.total = res.data.total || 0;
        } else {
          this.$message.error(res.message);
        }
      })
      .catch(() => {});
  }

  /**
   * 获取状态文本
   */
  getStatusText(status: number) {
    const statusMap = {
      1: "巡店预约",
      2: "巡店中",
      3: "巡店完成",
      4: "关闭",
    };
    return statusMap[status] || "-";
  }

  /**
   * 获取状态类型
   */
  getStatusType(status: number) {
    const statusMap = {
      1: "warning",
      2: "primary",
      3: "success",
      4: "info",
    };
    return statusMap[status] || "";
  }

  /**
   * 搜索
   */
  handleSearch() {
    this.listQuery.page = 1;
    this.getList();
  }

  /**
   * 重置
   */
  handleReset() {
    this.dateRange = null;
    this.visitDateRange = null;
    this.listQuery = {
      page: 1,
      size: 10,
      dept_name: "",
      store_name: "",
      start_time: 0,
      end_time: 0,
      status: 0,
      visit_start_time: 0,
      visit_end_time: 0,
    };
    this.getList();
  }

  /**
   * 改变每页数量
   */
  handleSizeChange(val: number) {
    this.listQuery.size = val;
    this.getList();
  }

  /**
   * 改变页码
   */
  handleCurrentChange(val: number) {
    this.listQuery.page = val;
    this.getList();
  }

  /**
   * 创建时间范围变化
   */
  handleDateChange(val: any) {
    if (val) {
      this.listQuery.start_time = Math.floor(val[0] / 1000);
      this.listQuery.end_time = Math.floor(val[1] / 1000);
    } else {
      this.listQuery.start_time = 0;
      this.listQuery.end_time = 0;
    }
  }

  /**
   * 巡店时间范围变化
   */
  handleVisitDateChange(val: any) {
    if (val) {
      this.listQuery.visit_start_time = Math.floor(val[0] / 1000);
      this.listQuery.visit_end_time = Math.floor(val[1] / 1000);
    } else {
      this.listQuery.visit_start_time = 0;
      this.listQuery.visit_end_time = 0;
    }
  }

  /**
   * 查看详情
   */
  handleDetail(row: any) {
    this.$router.push({
      path: `/myapp/wisdom/store-visit/data`,
      query: {
        id: row.id,
      },
    });
  }

  /**
   * 关闭巡店
   */
  handleChangeStatus(row: any) {
    this.$confirm(`确定要关闭该巡店记录吗？`, "提示", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning",
    })
      .then(() => {
        StoreVisitApi.closeStoreVisit({
          id: row.id,
        }).then((res: any) => {
          if (res.code === 200) {
            this.$message({
              type: "success",
              message: "操作成功!",
            });
            this.getList();
          } else {
            this.$message.error(res.message);
          }
        });
      })
      .catch(() => {
        // 取消操作
      });
  }

  /**
   * 导出
   */
  async handleExport() {
    // StoreVisitApi.exportStoreVisit({
    //   ...this.listQuery,
    // }).then((res: any) => {
    //   downloadFile(res.data.url)
    // });
    const res:any = await StoreVisitApi.exportStoreVisit(this.listQuery)
    if(res.code == 200) {
      cosUploader.getObjectUrl(res.data.url).then(res => {
        downloadFile(res, {
          filename: `巡店数据_${dayjs().format('YYYYMMDD')}.xlsx`
        })
      })
    } else {
      this.$message.error('导出失败')
    }
  }
}
</script>

<style lang="scss" scoped>
.form-inline {
  display: flex;
  flex-wrap: wrap;
}
.tag {
  display: flex;
  align-items: center;
  &.warning {
    color: #fe5500;
  }
  &.primary {
    color: #00a854;
  }
  &.success {
    color: #108de9;
  }
  &.info {
    color: #606266;
  }
  &::before {
    content: "";
    display: block;
    margin-right: 4px;
    width: 12px;
    height: 12px;
    border-radius: 50%;
    background: currentColor;
    transform: scale(0.5);
  }
}
</style>
