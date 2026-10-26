<template>
  <div class="store-list">
    <layout-list>
      <template slot="search">
        <div class="top-time" v-if="false">
          <el-button type="text" @click="addStore(row)">更新员工信息</el-button>
          <span class="time">上次更新时间: {{ refresh_time }}</span>
        </div>
        <search-item text="员工姓名：">
          <cy-input
            size="small"
            prefix-icon="el-icon-search"
            placeholder="请输入员工姓名或工号"
            clearable
            v-model="keyword"
          ></cy-input>
        </search-item>
        <search-item>
          <el-button type="primary" @click="search">搜索</el-button>
          <el-button @click="reset">重置</el-button>
        </search-item>
      </template>
      <template slot="operate">
        <el-button
          :disabled="select_staff.length === 0"
          type="primary"
          @click="allTansferStore"
          >批量关联</el-button
        >
      </template>
      <template slot="table">
        <el-table
          ref="multipleTable"
          tooltip-effect="light"
          :data="staff_list"
          stripe
          max-height="600"
          style="width:100%;"
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="55"> </el-table-column>
          <el-table-column label="员工姓名" min-width="120">
            <template slot-scope="{ row }">
              {{ row.name }}
            </template>
          </el-table-column>
          <el-table-column label="员工ID" width="180">
            <template slot-scope="{ row }">
              {{ row.employee_code }}
            </template>
          </el-table-column>
          <el-table-column label="手机号" width="150">
            <template slot-scope="{ row }">
              {{
                row.mobile
                  ? row.mobile.replace(row.mobile.substring(3, 7), "****")
                  : "--"
              }}
            </template>
          </el-table-column>
          <el-table-column label="性别" width="125">
            <template slot-scope="{ row }">
              {{ row.gender == 1 ? "男" : row.gender == 2 ? "女" : "未定义" }}
            </template>
          </el-table-column>
          <el-table-column
            label="工作地点"
            min-width="185"
            prop="job_location"
          ></el-table-column>
          <el-table-column fixed="right" label="操作" width="120">
            <template slot-scope="{ row }">
              <el-button type="text" @click="tansferStore(row)"
                >关联门店</el-button
              >
            </template>
          </el-table-column>
          <template slot="empty" v-if="staff_list.length == 0">
            <img
              class="empty-img m-t-40"
              src="@/assets/images/empty-box.png"
              alt=""
            />
            <div class="empty-box">
              <div class="desc">暂无员工</div>
            </div>
          </template>
        </el-table>
      </template>
      <template slot="footer" v-if="staff_list.length !== 0">
        <el-pagination
          class="float-right m-t-20"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page.sync="meta.page"
          :page-sizes="[10, 20, 30, 40]"
          :page-size="meta.size"
          layout="total, sizes, prev, pager, next, jumper"
          :total="meta.total"
        >
        </el-pagination>
      </template>
    </layout-list>
    <tansfer-staff
      @close="tansferClose"
      @success="tansferSuccess"
      :selectList="select_list"
      v-if="tansfer_show"
    ></tansfer-staff>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LayoutList from "@/components/public/layout-list/index.vue";
import Api from "@/api/enterprise/business";
import TansferStaff from "./../component/tansfer-staff.vue";
@Component({
  components: {
    LayoutList,
    TansferStaff
  }
})
export default class extends Vue {
  // ---------------------------------------------------------------此页面所需变量---------------------------------------------------------------
  // 分页配置
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  private refresh_time = "2202-01-20 16:22:00";
  private keyword = "";
  private staff_list = [];
  private select_list = [];
  private select_staff = [];
  private tansfer_show = false;
  /************************************* 方法 ************************************************ */
  private getStaffList() {
    const params = {
      page: this.meta.page,
      size: this.meta.size,
      mobile_or_name: this.keyword,
      no_department: 1
    };
    Api.getStaffListApi(params).then(res => {
      const { list, ...meta } = res.data;
      this.staff_list = list;
      this.meta = meta;
    });
  }
  // 搜索
  private search() {
    this.meta.page = 1;
    this.getStaffList();
  }
  private reset() {
    this.keyword = "";
    this.search();
  }
  private addStore() {
    console.log("新增门店");
  }
  private allTansferStore() {
    this.select_list = JSON.parse(JSON.stringify(this.select_staff));
    this.tansfer_show = true;
  }
  private tansferSuccess() {
    this.tansferClose();
    this.search();
  }
  private tansferStore(row) {
    this.select_list = [row];
    this.tansfer_show = true;
  }
  private transferAll() {
    this.tansfer_show = true;
  }
  private tansferClose() {
    this.tansfer_show = false;
  }
  handleSelectionChange(val) {
    this.select_staff = val;
  }
  // 分页
  handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.getStaffList();
  }
  handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getStaffList();
  }

  // ---------------------------------------------------------------生命周期---------------------------------------------------------------

  // 实例挂载完
  mounted() {
    console.log("进入页面");
    this.getStaffList();
  }
}
</script>
<style lang="scss" scoped>
.top-time {
  margin-bottom: 20px;
  .time {
    margin-left: 10px;
    color: #606266;
  }
}
</style>
