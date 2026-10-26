<template>
  <div class="store-list">
    <layout-list>
      <template slot="search">
        <search-item>
          <cy-input
            size="large"
            prefix-icon="el-icon-search"
            placeholder="请输入文件名称"
            clearable
            v-model="query.file_name"
          ></cy-input>
        </search-item>
        <search-item>
          <el-date-picker
            v-model="time"
            @change="selectChangeFn"
            value-format="timestamp"
            type="datetimerange"
            range-separator="至"
            :default-time="['00:00:00', '23:59:59']"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
          >
          </el-date-picker>
        </search-item>
        <search-item>
          <el-button type="primary" @click="search">搜索</el-button>
          <el-button @click="reset">重置</el-button>
        </search-item>
      </template>
      <template slot="table">
        <el-table
          tooltip-effect="light"
          :data="list"
          stripe
          max-height="600"
          style="width:100%;"
        >
          <el-table-column
            label="编号"
            prop="id"
            min-width="120"
          ></el-table-column>
          <el-table-column
            prop="file_name"
            label="文件名"
            min-width="200"
          ></el-table-column>
          <el-table-column
            prop="error_count"
            label="错误数据量"
            min-width="150"
          ></el-table-column>
          <el-table-column label="同步时间" width="185">
            <template slot-scope="{ row }">
              {{ row.created_at | formatDate("yyyy-mm-dd hh:mm:ss") }}
            </template>
          </el-table-column>
          <el-table-column fixed="right" label="操作" width="120">
            <template slot-scope="{ row }">
              <el-button
                v-if="row.error_count > 0"
                type="text"
                @click="dowmFile(row)"
                >下载</el-button
              >
            </template>
          </el-table-column>
          <template slot="empty" v-if="list.length == 0">
            <img
              class="empty-img m-t-40"
              src="@/assets/images/empty-box.png"
              alt=""
            />
            <div class="empty-box">
              <div class="desc">暂无数据</div>
            </div>
          </template>
        </el-table>
      </template>
      <template slot="footer" v-if="list.length !== 0">
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
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LayoutList from "@/components/public/layout-list/index.vue";
import Api from "@/api/enterprise/business";
import { downFilePost } from "@/utils/tools";
@Component({
  components: {
    LayoutList
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
  private time = "";
  private query = {
    file_name: "",
    start_time: 0,
    end_time: 0
  };
  private list = [];
  /************************************ 方法 ********************************************* */
  private getList() {
    const params = {
      ...this.query,
      page: this.meta.page,
      size: this.meta.size
    };
    Api.syncListApi(params).then(res => {
      const { total, page, size, list } = res.data;
      this.list = list;
      this.meta = {
        total,
        page,
        size
      };
    });
  }
  private dowmFile(row: any) {
    downFilePost(Api.downloadFileApi, { file_log_id: row.id }, row.file_name);
  }
  private selectChangeFn(e) {
    this.query.start_time = !e ? 0 : e[0] / 1000;
    this.query.end_time = !e ? 0 : e[1] / 1000;
  }
  // 搜索
  private search() {
    this.meta.page = 1;
    this.getList();
  }
  private reset() {
    this.query.file_name = "";
    this.time = "";
    this.selectChangeFn(null);
    this.search();
  }
  // 分页
  handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.getList();
  }
  handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getList();
  }

  // ---------------------------------------------------------------生命周期---------------------------------------------------------------

  // 实例挂载完
  mounted() {
    console.log("进入页面");
    this.getList();
  }
}
</script>
<style lang="scss" scoped></style>
