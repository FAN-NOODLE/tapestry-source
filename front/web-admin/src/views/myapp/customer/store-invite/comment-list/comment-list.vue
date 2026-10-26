<template>
  <div class="comment-list">
    <layout-list>
      <template slot="search">
        <search-item text="评分时间： ">
          <el-date-picker
            v-model="search.score_time"
            type="datetimerange"
            prefix-icon="el-icon-date"
            format="yyyy-MM-dd HH:mm:ss"
            value-format="timestamp"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            :default-time="['00:00:00', '23:59:59']"
            style="height:36px"
          >
          </el-date-picker>
        </search-item>
        <search-item text="预约时间： ">
          <el-date-picker
            v-model="search.invite_time"
            type="datetimerange"
            prefix-icon="el-icon-date"
            format="yyyy-MM-dd HH:mm:ss"
            value-format="timestamp"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            :default-time="['00:00:00', '23:59:59']"
            style="height:36px"
          >
          </el-date-picker>
        </search-item>
        <search-item text="预约门店： ">
          <SelectStore ref="selectStore"></SelectStore>
        </search-item>
        <search-item text="预约导购： ">
          <SelectStaff ref="selectStaff"></SelectStaff>
        </search-item>
        <search-item>
          <el-button type="primary" @click="goSearch" size="small"
            >查询</el-button
          >
          <el-button plain @click="reset" size="small">重置</el-button>
          <el-button plain @click="exportList" size="small">导出</el-button>
        </search-item>
      </template>
      <template slot="table">
        <el-table
          tooltip-effect="light"
          :data="list"
          ref="multipleTable"
          @selection-change="selectChange"
        >
          <!-- <el-table-column type="selection" width="60"></el-table-column> -->
          <el-table-column label="编号" prop="content" min-width="180">
            <template slot-scope="scope">
              {{ scope.row.reserve_number }}
            </template>
          </el-table-column>
          <el-table-column label="顾客昵称" prop="content" min-width="150">
            <template slot-scope="scope">
              {{ scope.row.user_name }}
            </template>
          </el-table-column>
          <el-table-column label="服务导购" prop="content" min-width="150">
            <template slot-scope="scope">
              {{ scope.row.member_name }}
            </template>
          </el-table-column>
          <el-table-column label="预约门店" prop="content" min-width="150">
            <template slot-scope="scope">
              {{ scope.row.store_name || "--" }}
            </template>
          </el-table-column>
          <el-table-column label="预约时间" prop="content" min-width="190">
            <template slot-scope="scope">
              {{ scope.row.reserve_time | formatDate("yyyy-mm-dd hh:mm:ss") }}
            </template>
          </el-table-column>

          <el-table-column
            v-if="epName == 'coach'"
            label="服务评分"
            prop="content"
            min-width="110"
          >
            <template slot-scope="scope">
              {{ scope.row.service_score }}
            </template>
          </el-table-column>
          <el-table-column
            v-if="epName == 'coach'"
            label="门店评分"
            prop="content"
            min-width="110"
          >
            <template slot-scope="scope">
              {{ scope.row.store_score }}
            </template>
          </el-table-column>
          <el-table-column
            v-if="epName == 'ks' || epName == 'sw'"
            label="活动安排"
            prop="content"
            min-width="150"
          >
            <template slot-scope="scope">
              {{ scope.row.service_score }}
            </template>
          </el-table-column>
          <el-table-column
            v-if="epName == 'ks' || epName == 'sw'"
            label="整体服务"
            prop="content"
            min-width="150"
          >
            <template slot-scope="scope">
              {{ scope.row.store_score }}
            </template>
          </el-table-column>
          <!-- <el-table-column
            v-if="epName == 'ks'"
            label="服务流程"
            prop="content"
            min-width="150"
          >
            <template slot-scope="scope">
              {{ scope.row.flow_score }}
            </template>
          </el-table-column> -->
          <el-table-column label="评分时间" prop="content" min-width="190">
            <template slot-scope="scope">
              {{ scope.row.updated_at | formatDate("yyyy-mm-dd hh:mm:ss") }}
            </template>
          </el-table-column>
          <el-table-column label="顾客评价" width="300">
            <template slot-scope="scope">
              {{ scope.row.remark || "--" }}
            </template>
          </el-table-column>
          <div slot="empty" v-if="!emptyFlag" style="margin-top: 50px;">
            <img class="empty-img" src="@/assets/images/empty-box.png" alt="" />
            <div class="empty-box">
              <div class="desc">暂无评论</div>
            </div>
          </div>
          <div slot="empty" v-if="emptyFlag" style="margin-top: 50px;">
            <img
              class="empty-img"
              src="@/assets/images/empty-search.png"
              alt=""
            />
            <div class="empty-box">
              <div class="desc">没有找到相关评论</div>
            </div>
          </div>
        </el-table>
      </template>
      <template slot="footer" v-if="list.length > 0">
        <el-pagination
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
import SelectStore from "../component/select-store.vue";
import SelectStaff from "../component/select-staff.vue";
import Api from "@/api/store-invite/index";
import { downFilePost } from "@/utils/tools";
@Component({
  components: {
    SelectStore,
    SelectStaff
  }
})
export default class extends Vue {
  private search = {
    score_time: [],
    invite_time: []
  };
  private meta = {
    page: 1,
    size: 10,
    total: 5
  };
  private list = [];
  private epName = "";
  private multipleSelection = [];
  private emptyFlag = false;
  // 生命周期
  private created() {
    console.log("created");
  }
  private mounted() {
    console.log("mounted");
    this.goSearch();
  }
  private goSearch() {
    console.log("搜索");
    this.meta.page = 1;
    this.getList();
  }
  private exportList() {
    let params = {};
    if (this.multipleSelection.length === 0) {
      const { page, size, ...obj } = this.getParams();
      params = obj;
    } else {
      params = { ids: this.multipleSelection };
    }
    downFilePost(Api.commentExportApi, params, "评论列表");
    // Api.serverExportApi()
  }
  // 分页
  private handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.getList();
  }
  private handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getList();
  }
  private getParams() {
    const selectStore: any = this.$refs.selectStore;
    const selectStaff: any = this.$refs.selectStaff;
    const params = {
      store_ids: selectStore.list.map(item => {
        return item.dp_third_id;
      }),
      member_ids: selectStaff.list.map(item => {
        return item.member_id;
      }),
      page: this.meta.page,
      size: this.meta.size,
      reserve_start_time: this.search.invite_time
        ? (this.search.invite_time[0] || 0) / 1000
        : 0,
      reserve_end_time: this.search.invite_time
        ? (this.search.invite_time[1] || 0) / 1000
        : 0,
      score_start_time: this.search.score_time
        ? (this.search.score_time[0] || 0) / 1000
        : 0,
      score_end_time: this.search.score_time
        ? (this.search.score_time[1] || 0) / 1000
        : 0
    };
    return params;
  }
  // 当选择项发生变化时会触发该事件，数组
  selectChange(val: any) {
    this.multipleSelection = [];
    val.map((item: any) => {
      this.multipleSelection.push(item.id);
    });
  }
  private reset() {
    console.log("重置");
    this.search = {
      score_time: [],
      invite_time: []
    };
    (this.$refs.selectStore as any).clear();
    (this.$refs.selectStaff as any).clear();
    this.meta.page = 1;
    this.getList();
  }
  private getList() {
    Api.getEpConfig().then(res => {
      // // 企业名称：coach;ks;sw
      const { name } = res.data;
      this.epName = name;
      const params = this.getParams();
      Api.userCommonApi(params).then(res => {
        const { total, list } = res.data;
        this.list = list;
        this.meta.total = total;
      });
    });
  }
  private destroyed() {
    console.log("destroyed");
  }
}
</script>
<style lang="scss" scoped></style>
