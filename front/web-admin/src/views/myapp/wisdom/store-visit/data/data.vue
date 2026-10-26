<template>
  <div class="store-visit-data">
    <layout-list>
      <template slot="search">
        <el-form :inline="true" class="form-inline">
          <el-form-item label="区域:">
            <el-select
              v-model="listQuery.area_id"
              filterable
              placeholder="请选择区域"
            >
              <el-option label="全部" :value="0"></el-option>
              <el-option
                v-for="(item, index) in areaList"
                :key="index"
                :label="item.name"
                :value="item.id"
              ></el-option>
            </el-select>
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
          <el-form-item label="统计时间:">
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
          <el-form-item>
            <el-button type="primary" @click="handleSearch">查询</el-button>
            <el-button @click="handleReset">重置</el-button>
            <el-button @click="handleExport">导出</el-button>
          </el-form-item>
        </el-form>
      </template>

      <template slot="tabs">
        <el-tabs v-model="activeTab" @tab-click="handleTabChange">
          <el-tab-pane label="按门店汇总" name="1"></el-tab-pane>
          <el-tab-pane label="按渠道汇总" name="2"></el-tab-pane>
          <el-tab-pane label="按大区汇总" name="3"></el-tab-pane>
          <el-tab-pane label="按小区汇总" name="4"></el-tab-pane>
        </el-tabs>
      </template>

      <template slot="operate">
        <div class="summary-info">
          <div class="summary-item">
            <span class="label">总巡店数: </span>
            <span class="value">{{ totalStores }}</span>
          </div>
        </div>
      </template>

      <template slot="table">
        <el-table :data="list" stripe highlight-current-row>
          <el-table-column label="门店编号" prop="store_code" min-width="120" v-if="activeTab == '1'">
            <template slot-scope="{ row }">{{row.store_code || "-"}}</template>
          </el-table-column>
          <el-table-column
            label="门店名称"
            prop="name"
            min-width="180"
            v-if="activeTab == '1'"
          ></el-table-column>
          <el-table-column
            label="渠道"
            prop="name"
            min-width="180"
            v-if="activeTab == '2'"
          ></el-table-column>
          <el-table-column
            label="大区名称"
            prop="name"
            min-width="180"
            v-if="activeTab == '3'"
          ></el-table-column>
          <el-table-column
            label="小区名称"
            prop="name"
            min-width="180"
            v-if="activeTab == '4'"
          ></el-table-column>
          <el-table-column
            label="大区名称"
            prop="parent_name"
            min-width="180"
            v-if="activeTab == '4'"
          ></el-table-column>
          <el-table-column label="部门" prop="dept_name" min-width="120" v-if="activeTab == '1'">
            <template slot-scope="{ row }">{{ row.dept_name || "-" }}</template>
          </el-table-column>
          <el-table-column
            label="巡店数"
            prop="visit_num"
            min-width="100"
          ></el-table-column>
          <el-table-column
            label="巡店完成率"
            prop="visit_rate"
            min-width="120"
          ></el-table-column>
        </el-table>
      </template>

      <template slot="footer" v-if="activeTab == '1'">
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
  </div>
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
export default class StoreVisitData extends Vue {
  deptList: string[] = [];
  areaList: any[] = [];
  list: any[] = [];
  total = 0;
  totalStores = 0;
  totalVisits = 0;
  activeTab = "1";
  dateRange: any = null;

  listQuery = {
    page: 1,
    size: 10,
    dept_name: "",
    area_id: 0,
    start_time: 0,
    end_time: 0,
    type: 1,
  };

  created() {
    this.getAreaList();
    this.getDeptList();
    this.getList();
  }

  /**
   * 获取区域列表
   */
  getAreaList() {
    StoreVisitApi.getAreaList({}).then((res: any) => {
      if (res.code === 200) {
        this.areaList = res.data || [];
      }
    });
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
   * 获取表格数据
   */
  getList() {
    StoreVisitApi.getStoreVisitDashboard(this.listQuery)
      .then((res: any) => {
        if (res.code === 200) {
          this.list = res.data.list || [];
          this.total = res.data.total || 0;
          this.totalStores = res.data.num || 0
        } else {
          this.$message.error(res.message);
        }
      })
      .catch(() => {});
  }

  /**
   * 处理标签页变化
   */
  handleTabChange(tab: any) {
    this.listQuery.type = parseInt(tab.name);
    this.listQuery.page = 1;
    this.getList();
  }

  /**
   * 处理日期变化
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
    this.listQuery = {
      page: 1,
      size: 10,
      dept_name: "",
      area_id: 0,
      start_time: 0,
      end_time: 0,
      type: parseInt(this.activeTab),
    };
    this.getList();
  }

  /**
   * 处理页面大小变化
   */
  handleSizeChange(val: number) {
    this.listQuery.size = val;
    this.getList();
  }

  /**
   * 处理页码变化
   */
  handleCurrentChange(val: number) {
    this.listQuery.page = val;
    this.getList();
  }

  /**
   * 导出数据
   */
  async handleExport() {
    // StoreVisitApi.exportDashboard({
    //   ...this.listQuery,
    // }).then((res: any) => {
    //   downloadFile(res.data.url);
    // });
    const res:any = await StoreVisitApi.exportDashboard(this.listQuery)
    if(res.code == 200) {
      cosUploader.getObjectUrl(res.data.url).then(res => {
        downloadFile(res, {
          filename: `巡店数据明细_${dayjs().format('YYYYMMDD')}.xlsx`
        })
      })
    } else {
      this.$message.error('导出失败')
    }
  }
}
</script>

<style lang="scss" scoped>
.store-visit-data {
  .summary-info {
    display: flex;
    margin-bottom: 10px;

    .summary-item {
      margin-right: 20px;
      font-size: 14px;

      .label {
        color: #606266;
      }

      .value {
        font-weight: bold;
        color: #409eff;
        font-size: 16px;
      }
    }
  }

  .form-inline {
    display: flex;
    flex-wrap: wrap;
  }
}
</style>
