<template>
  <div class="supernatant">
    <layout-list>
      <template slot="buttons">
        <el-button
          :disabled="groupList.length == 0"
          type="primary"
          @click="exportInfo()"
          class="font-14"
          size="small"
        >导出明细</el-button>
      </template>
      <template slot="table">
        <el-table
          v-if="$route.query.type == 1"
          :data="groupList"
          stripe
          style="width: 100%"
          tooltip-effect="light"
        >
          <el-table-column
            show-overflow-tooltip
            label="顾客昵称"
            width="300"
            prop="user_name"
          >
          </el-table-column>
          <el-table-column
            label="发送状态"
            width="400"
            prop="status"
            show-overflow-tooltip
          >
            <template slot-scope="scope">
              <span v-if="scope.row.status == 0">未发送</span>
              <span v-else-if="scope.row.status == 1">已发送</span>
              <span v-else-if="scope.row.status == 2">因客户不是好友导致发送失败</span>
              <span v-else-if="scope.row.status == 3">因客户已经收到其他群发消息导致发送失败</span>
              <span v-else>- -</span>
            </template>
          </el-table-column>

          <el-table-column
            prop="send_time"
            label="发送时间"
            min-width="300"
          >
            <template slot-scope="scope">
              <span v-if="scope.row.send_time > 0">{{
                scope.row.send_time | formatDate("yyyy-mm-dd hh:mm:ss")
              }}</span>
              <span v-else>--</span>
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
        <el-table
          v-else
          :data="groupList"
          stripe
          style="width: 100%"
          tooltip-effect="light"
        >
          <el-table-column
            show-overflow-tooltip
            label="客户群ID"
            min-width="200"
            prop="chat_id"
          >
          </el-table-column>
          <el-table-column
            show-overflow-tooltip
            label="客户群名称"
            width="200"
            prop="chat_name"
          >
          </el-table-column>
          <el-table-column
            label="发送状态"
            width="200"
            prop="status"
            show-overflow-tooltip
          >
            <template slot-scope="scope">
              <span v-if="scope.row.status == 0">未发送</span>
              <span v-else-if="scope.row.status == 1">已发送</span>
              <span v-else-if="scope.row.status == 2">因客户不是好友导致发送失败</span>
              <span v-else-if="scope.row.status == 3">因客户已经收到其他群发消息导致发送失败</span>
              <span v-else>- -</span>
            </template>
          </el-table-column>

          <el-table-column
            prop="channel_name"
            label="发送时间"
            min-width="300"
          >
            <template slot-scope="scope">
              <span v-if="scope.row.send_time > 0">{{
                scope.row.send_time | formatDate("yyyy-mm-dd hh:mm:ss")
              }}</span>
              <span v-else>--</span>
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
import { Vue, Component } from "vue-property-decorator";
import { taskSendApi } from "@/api/enterprise/association";
import { ApiResult, Meta } from "@/types";
@Component({})
export default class extends Vue {
  //----------------data ------------------/
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
    this.tableData();
  }
  //----------------method ------------------/
  // 列表数据
  private async tableData() {
    const { page, size } = this.meta;
    const { member_id, msg_id } = this.$route.query;
    const params = {
      page,
      size,
      member_id: Number(member_id),
      msg_id
    };
    const { code, message = "", data }: ApiResult = await taskSendApi(
      params,
      {}
    );
    if (code !== 200) {
      this.$message.error(message);
      return;
    }
    this.groupList = data.list && data.list.length ? data.list : [];
    this.meta.total = data.total;
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
  // 导出
  exportInfo() {
    const { member_id, msg_id } = this.$route.query;
    const params = {
      page: 1,
      size: 999999,
      member_id: Number(member_id),
      export: 1,
      msg_id
    };
    const responseType = {
      responseType: "blob"
    };
    taskSendApi(params, responseType).then((res: any) => {
      const blob = new Blob([res], {
        type: "application/x-xls;charset=utf-8"
      });
      console.log(blob);
      const fileName: string =
        "明细" + new Date().toLocaleString().substr(0, 10) + ".xls";
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
}
</script>
<style lang="scss" scoped>
.table-empty-label {
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
