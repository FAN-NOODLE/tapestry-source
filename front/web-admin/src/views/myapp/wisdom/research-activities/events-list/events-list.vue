<template>
  <div class="events-list">
    <div class="cy-title">
      数据明细
    </div>
    <div class="search_content">
      <search-item text="参与时间：">
        <el-date-picker
          v-model="signTime"
          :default-time="['00:00:00', '23:59:59']"
          value-format="timestamp"
          type="datetimerange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          clearable
        >
        </el-date-picker>
      </search-item>
      <search-item text="审核状态：">
        <cy-select
          width-size="small"
          size="small"
          v-model="status"
        >
          <cy-option
            label="全部"
            :value="0"
          ></cy-option>
          <cy-option
            label="待审核"
            :value="1"
          ></cy-option>
          <cy-option
            label="审核通过"
            :value="2"
          ></cy-option>
        </cy-select>
      </search-item>
      <div class="search_btn">
        <el-button
          type="primary"
          @click="search"
          size="small"
        >查询</el-button>
        <el-button
          @click="reset"
          size="small"
        >重置</el-button>
        <el-button
          @click="exportExcel"
          size="small"
        >数据导出</el-button>
      </div>
    </div>
    <div class="table-content">
      <el-table
        :data="list"
        stripe
        tooltip-effect="light"
      >
        <el-table-column
          prop="wx_name"
          label="用户昵称"
          width="160"
        >
        </el-table-column>
        <el-table-column
          prop="receiver_name"
          label="收货人姓名"
          width="120"
        >
        </el-table-column>
        <el-table-column
          prop="receiver_phone"
          label="收货人手机号"
          width="160"
        >
        </el-table-column>
        <el-table-column
          prop="receiver_address"
          label="收货地址"
          min-width="240"
        >
        </el-table-column>
        <el-table-column
          label="参与时间"
          min-width="140"
        >
          <template slot-scope="scope">
            <span>
              {{ scope.row.created_at | formatDate("yyyy-mm-dd hh:mm:ss") }}
            </span>
          </template>
        </el-table-column>
        <el-table-column
          prop="status_text"
          label="审核状态"
          width="100"
        >
        </el-table-column>
        <el-table-column
          label="操作"
          width="120"
          fixed="right"
        >
          <template slot-scope="scope">
            <el-button
              type="text"
              size="small"
              class="marginR20"
              @click="checkDetail(scope.row)"
            >查看详情</el-button>
          </template>
        </el-table-column>
        <div
          slot="empty"
          class="table-empty-label-page"
        >
          <img
            class="empty-img"
            :src="require('@/assets/images/empty-box.png')"
            alt=""
          />
          <div>暂无活动</div>
        </div>
      </el-table>
      <div class="table_footer">
        <el-button
          size="medium"
          @click="goBack"
        >返回</el-button>
        <el-pagination
          v-if="list && list.length > 0"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page="meta.page"
          :page-sizes="[10, 20, 30, 40]"
          :page-size="meta.size"
          layout="total, sizes, prev, pager, next, jumper"
          :total="meta.total"
        >
        </el-pagination>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Api from "@/api/wisdom/researchActivities";
@Component
export default class extends Vue {
  // data
  private activity_id = 0;
  private signTime = 0;
  private status = 0;
  private searchLoading = false;
  private list: Array<any> = [];
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  // 生命周期
  private mounted() {
    this.activity_id = Number(this.$route.query.id);
    this.getList();
  }
  private getList() {
    const params = {
      page: this.meta.page,
      size: this.meta.size,
      activity_id: this.activity_id,
      start_time: (this.signTime && this.signTime[0] / 1000) || undefined,
      end_time: (this.signTime && this.signTime[1] / 1000) || undefined,
      status: this.status
    };
    Api.getUserList(params).then((res: any) => {
      if (res.code === 200) {
        this.list = res.data.list;
        this.meta.total = res.data.total;
      }
    });
  }

  private search() {
    this.meta.page = 1;
    this.searchLoading = true;
    this.getList();
  }

  private reset() {
    this.signTime = 0;
    this.status = 0;
    this.searchLoading = false;
    this.getList();
  }

  // 时间格式转换
  private formateDate(time: any) {
    const year = time.getFullYear();
    let month = time.getMonth() + 1;
    let day = time.getDate();

    month = month < 10 ? "0" + month : month;
    day = day < 10 ? "0" + day : day;

    return year + "-" + month + "-" + day;
  }

  private exportExcel() {
    const params = {
      activity_id: this.activity_id,
      start_time: (this.signTime && this.signTime[0] / 1000) || undefined,
      end_time: (this.signTime && this.signTime[1] / 1000) || undefined,
      status: this.status
    };
    Api.exportUserData(params).then((res: any) => {
      const blob = new Blob([res], {
        type: "application/x-xls;charset=utf-8"
      });
      let fileName = "";
      const nowTime = new Date();
      fileName = "活动列表_" + this.formateDate(nowTime) + ".xlsx";

      //下载的文件名称及其后缀，后缀要和后台保持的一致

      if ("download" in document.createElement("a")) {
        // 非IE下载
        const elink = document.createElement("a");
        elink.download = fileName;
        elink.style.display = "none";
        elink.href = URL.createObjectURL(blob);
        document.body.appendChild(elink);
        elink.click();
        URL.revokeObjectURL(elink.href); // 释放URL 对象
        document.body.removeChild(elink);
      } else {
        // IE10+下载
        (navigator as any).msSaveBlob(blob, fileName);
      }
    });
  }

  // 每页
  private handleSizeChange(val: any) {
    this.meta.page = 1;
    this.meta.size = val;
    this.getList();
  }

  // 当前页
  private handleCurrentChange(val: any) {
    this.meta.page = val;
    this.getList();
  }

  private checkDetail(row: any) {
    this.$router.push({
      name: "myapp_wisdom_research-activities_events-detail",
      query: {
        id: row.id
      }
    });
  }

  private goBack() {
    this.$router.go(-1);
  }
}
</script>
<style lang="scss" scoped>.cy-title {
  background: #f2f2f2;
  height: 40px;
  line-height: 40px;
  padding: 0 10px;
  border-radius: 2px;
  margin-bottom: 30px;
}

.events-list {
  background-color: #fff;
  padding: 20px;
  border-radius: 10px;

  .search_btn {
    display: inline-block;
  }
}

.table_footer {
  display: flex;
  margin-top: 20px;
  justify-content: space-between;
}

.table-empty-label-page {
  display: flex;
  height: 340px;
  flex-direction: column;

  img {
    width: 124px;
    margin: 80px auto 0;
  }

  div {
    line-height: 20px;
    margin-top: 14px;
  }
}</style>
