<template>
  <div class="supernatant">
    <layout-list>
      <template slot="search">
        <search-item text="执行者：">
          <cy-select
            v-model="search.staff"
            placeholder="全部"
            search-input
            :optionsList="search.staffList"
            :options="{ label: 'name', value: 'member_id' }"
          >
          </cy-select>
        </search-item>
        <search-item text="发送状态：">
          <el-select
            v-model="search.task_status"
            clearable
            placeholder="请选择"
          >
            <el-option
              v-for="(item, index) in statusList"
              :key="index"
              :label="item.label"
              :value="item.value"
            >
            </el-option>
          </el-select>
        </search-item>
        <!-- <search-item :text="$route.query.type == 1 ? '客户昵称：' : '群名称：'">
          <cy-input
            size="small"
            prefix-icon="el-icon-search"
            :placeholder="
              $route.query.type == 1 ? '请输入客户昵称' : '请输入群名称'
            "
            clearable
            class="custom-name"
            v-model="search.title"
          ></cy-input>
        </search-item> -->
        <!-- <search-item text="最近发送时间：">
          <el-date-picker
            v-model="search.senfTime"
            type="datetimerange"
            range-separator="-"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            :default-time="['00:00:00', '23:59:59']"
            value-format="yyyy-MM-dd HH:mm:ss"
          ></el-date-picker>
        </search-item> -->
        <!-- <div> -->
        <el-button
          type="primary"
          @click="clickSearch()"
        >查询</el-button>
        <el-button @click="reset()">重置</el-button>
        <!-- </div> -->
      </template>
      <template slot="table">
        <el-table
          :data="groupList"
          stripe
          style="width: 100%"
          tooltip-effect="light"
        >
          <el-table-column
            show-overflow-tooltip
            label="执行者"
            width="200"
            prop="name"
          >
          </el-table-column>
          <el-table-column
            show-overflow-tooltip
            label="发送状态"
            width="100"
            prop="task_status"
          >
            <template slot-scope="scope">
              <span v-if="scope.row.status == 1">未发送</span>
              <span v-else-if="scope.row.status == 2">已发送</span>
              <span v-else-if="scope.row.status == 3">发送失败</span>
              <span v-else>- -</span>
            </template>
          </el-table-column>
          <el-table-column
            v-if="$route.query.type == 1"
            show-overflow-tooltip
            label="发送成功客户数"
            width="200"
            prop="success_count"
          >
            <template slot-scope="scope">
              <span v-if="scope.row.success_count">{{
                scope.row.success_count
              }}</span>
              <span v-else>0</span>
            </template>
          </el-table-column>
          <el-table-column
            v-if="$route.query.type == 1"
            show-overflow-tooltip
            label="发送失败客户数"
            width="200"
            prop="failed_count"
          >
            <template slot-scope="scope">
              <span v-if="scope.row.failed_count">{{
                scope.row.failed_count
              }}</span>
              <span v-else>0</span>
            </template>
          </el-table-column>
          <el-table-column
            v-if="$route.query.type == 1"
            show-overflow-tooltip
            label="未发送客户数"
            width="200"
            prop="not_send_count"
          >
            <template slot-scope="scope">
              <span v-if="scope.row.not_send_count">{{
                scope.row.not_send_count
              }}</span>
              <span v-else>0</span>
            </template>
          </el-table-column>
          <el-table-column
            v-if="$route.query.type == 2"
            show-overflow-tooltip
            label="发送成功客户群数"
            width="200"
            prop="success_count"
          >
            <template slot-scope="scope">
              <span v-if="scope.row.success_count">{{
                scope.row.success_count
              }}</span>
              <span v-else>0</span>
            </template>
          </el-table-column>
          <el-table-column
            v-if="$route.query.type == 2"
            show-overflow-tooltip
            label="发送失败客户群数"
            width="200"
            prop="failed_count"
          >
            <template slot-scope="scope">
              <span v-if="scope.row.failed_count">{{
                scope.row.failed_count
              }}</span>
              <span v-else>0</span>
            </template>
          </el-table-column>
          <el-table-column
            v-if="$route.query.type == 2"
            show-overflow-tooltip
            label="未发送客户群数"
            width="200"
            prop="not_send_count"
          >
            <template slot-scope="scope">
              <span v-if="scope.row.not_send_count">{{
                scope.row.not_send_count
              }}</span>
              <span v-else>0</span>
            </template>
          </el-table-column>
          <!-- <el-table-column
            :label="$route.query.type == 1 ? '客户昵称' : '群名称'"
            width="300"
            prop="chat_name"
            show-overflow-tooltip
          >
            <template slot-scope="scope">
              <span v-if="scope.row.chat_name">{{ scope.row.chat_name }}</span>
              <span v-else>未命名</span>
            </template>
          </el-table-column> -->

          <el-table-column
            prop="member_send_time"
            label="发送时间"
            min-width="200"
          >
            <template slot-scope="scope">
              <span v-if="scope.row.member_send_time > 0">{{
                scope.row.member_send_time | formatDate("yyyy-mm-dd hh:mm:ss")
              }}</span>
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column
            label="操作"
            width="150"
            fixed="right"
          >
            <template slot-scope="scope">
              <el-button
                type="text"
                @click="goDetails(scope.row)"
              >
                查看详情
              </el-button>
            </template>
          </el-table-column>
          <div
            slot="empty"
            class="table-empty-label"
          >
            <img
              class="empty-img"
              :src="searchLoading
                  ? 'https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/empty-search.png'
                  : require('@/assets/images/empty-box.png')
                "
              alt=""
            />
            <div>{{ searchLoading ? "没有找到相关明细" : "暂无明细" }}</div>
          </div>
        </el-table>
      </template>
      <template
        slot="footer"
        v-if="meta.total > 0"
      >
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
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import {
  taskDetailApi,
  memberlist,
  taskMemberDetailApi
} from "@/api/enterprise/association";
import { ApiResult, Meta } from "@/types";
@Component({})
export default class extends Vue {
  //----------------data ------------------/
  private statusList: any[] = [
    {
      label: "全部",
      value: 0
    },
    {
      label: "未发送",
      value: 2
    },
    {
      label: "已发送",
      value: 1
    }
  ];
  private search = {
    title: "",
    staff: 0,
    staffList: [{ member_id: 0, name: "全部" }],
    senfTime: 0,
    task_status: 0
  };

  // 是否正在搜索
  private searchLoading = false;
  // 分页配置
  private meta: Meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  private groupList: Array<any> = []; // 列表数组
  //----------------生命周期 ------------------/
  private mounted() {
    console.log("111111");
    this.tableData();
    this.memberlist();
  }
  //----------------method ------------------/
  private async memberlist() {
    const { code, message = "", data }: ApiResult = await memberlist(
      +this.$route.query.id
    );
    if (code !== 200) {
      this.$message.error(message);
      return;
    }
    this.search.staffList = [...this.search.staffList, ...data];
  }
  // 搜索查询
  private clickSearch() {
    this.meta.page = 1;
    this.meta.size = 10;
    this.searchLoading = true;
    this.tableData();
  }
  // 重置
  private reset() {
    this.search.title = "";
    this.search.staff = 0;
    this.search.task_status = 0;
    this.search.senfTime = 0;
    this.meta.page = 1;
    this.meta.size = 10;
    this.searchLoading = false;
    this.tableData();
  }
  // 列表数据
  private async tableData() {
    const { page, size } = this.meta;
    const params = {
      page,
      size,
      id: +this.$route.query.id,
      member_id: this.search.staff,
      task_status: this.search.task_status
      // chat_name: this.search.title,
      // start_time: new Date(this.search.senfTime[0]).getTime() / 1000 || 0,
      // end_time: new Date(this.search.senfTime[1]).getTime() / 1000 || 0
    };
    /** 员工明细 */
    const { code, message = "", data }: ApiResult = await taskMemberDetailApi(
      params
    );
    if (code !== 200) {
      this.$message.error(message);
      return;
    }
    this.groupList = data.list;
    this.meta.total = data.total;
    // else if (this.$route.query.type == "2") {
    //   /** 群发明细 */
    //   const {
    //     code,
    //     message = "",
    //     data,
    //   }: ApiResult = await taskDetailApi(params);
    //   if (code !== 200) {
    //     this.$message.error(message);
    //     return;
    //   }
    //   this.groupList = data.list;
    //   this.meta.total = data.total;
    // }
  }
  // 分页
  handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.tableData();
  }
  handleCurrentChange(val: number) {
    this.meta.page = val;
    this.tableData();
  }
  // 点击查看
  private goDetails(row: any) {
    this.$router.push({
      name: "myapp_enterprise_association_group-sent-message_info-detail",
      query: {
        msg_id: row.msgid,
        member_id: row.member_id,
        type: this.$route.query.type
      }
    });
  }
}
</script>
<style lang="scss" scoped>.table-empty-label {
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
