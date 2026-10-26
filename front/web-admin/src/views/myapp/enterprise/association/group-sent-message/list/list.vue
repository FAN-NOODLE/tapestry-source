<template>
  <div class="supernatant">
    <layout-list>
      <template slot="buttons">
        <el-button
          type="primary"
          @click="addMessage()"
          class="font-14"
          size="small"
          >新增群发消息</el-button
        >
      </template>
      <template slot="search">
        <search-item text="群发主题：">
          <cy-input
            size="small"
            prefix-icon="el-icon-search"
            placeholder="请输入群发标题"
            clearable
            class="custom-name"
            v-model="search.title"
          ></cy-input>
        </search-item>
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
          <cy-select
            v-model="search.sendStatus"
            placeholder=""
            :optionsList="search.sendStatusList"
          >
          </cy-select>
        </search-item>
        <search-item text="最近发送时间：">
          <el-date-picker
            style="height:36px"
            v-model="search.senfTime"
            type="datetimerange"
            range-separator="-"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            :default-time="['00:00:00', '23:59:59']"
            value-format="yyyy-MM-dd HH:mm:ss"
          ></el-date-picker>
        </search-item>
        <search-item text="群发类型：">
          <cy-select
            v-model="search.send_type"
            placeholder=""
            :optionsList="search.massStatusList"
          >
          </cy-select>
        </search-item>
        <div>
          <el-button
            type="primary"
            class="font-14"
            @click="clickSearch()"
            size="small"
            >查询</el-button
          >
          <el-button @click="reset()" class="font-14" size="small"
            >重置</el-button
          >
        </div>
      </template>
      <template slot="table">
        <el-table
          :data="groupList"
          stripe
          style="width: 100%"
          tooltip-effect="light"
        >
          <el-table-column
            label="群发主题"
            width="150"
            prop="name"
            show-overflow-tooltip
          >
          </el-table-column>
          <el-table-column
            label="群发类型"
            show-overflow-tooltip
            width="200"
            prop="send_type"
          >
            <template slot-scope="scope">
              <span>{{
                scope.row.send_type == 1
                  ? "群发客户"
                  : scope.row.send_type == 2
                  ? "群发客户群"
                  : "--"
              }}</span>
            </template>
          </el-table-column>
          <el-table-column
            label="执行者"
            show-overflow-tooltip
            width="200"
            prop="members"
          >
          </el-table-column>
          <el-table-column
            prop="content"
            label="发送内容"
            show-overflow-tooltip
            min-width="300"
          >
            <template slot-scope="scope">
              <div class="aligns-center">
                <el-popover
                  placement="top-start"
                  width="250"
                  trigger="hover"
                  :disabled="scope.row.content_obj.text.content.length <= 20"
                >
                  <div
                    style="overflow:hidden"
                    v-html="
                      scope.row.content_obj.text.content
                        .replace(/\n/g, '<br/>')
                        .replace(/\s/g, '&nbsp;')
                    "
                  ></div>
                  <span
                    slot="reference"
                    v-if="scope.row.content_obj.text.content.length <= 20"
                    class="m-l-10 goods-name el-popover-ref"
                    >{{ scope.row.content_obj.text.content }}</span
                  >
                  <span slot="reference" v-else class="m-l-10 goods-name">{{
                    scope.row.content_obj.text.content.substr(0, 20) + "..."
                  }}</span>
                </el-popover>
                <span v-if="'image' in scope.row.content_obj">
                  [图片]
                </span>
                <span v-if="'link' in scope.row.content_obj">
                  [链接]
                </span>
                <span v-if="'miniprogram' in scope.row.content_obj">
                  [小程序]
                </span>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="发送状态" min-width="150" prop="status">
            <template slot-scope="scope">
              <span v-if="scope.row.created_at > 0">{{
                scope.row.status == 1
                  ? "定时未开始"
                  : scope.row.status == 2
                  ? "未发送"
                  : scope.row.status == 3
                  ? "部分未发送"
                  : scope.row.status == 4
                  ? "已发送"
                  : scope.row.status == 5
                  ? "已失败"
                  : "任务下发中"
              }}</span>
            </template>
          </el-table-column>
          <el-table-column label="创建时间" min-width="200">
            <template slot-scope="scope">
              <span v-if="scope.row.created_at > 0">{{
                scope.row.created_at | formatDate("yyyy-mm-dd hh:mm:ss")
              }}</span>
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column
            label="发送时间"
            min-width="200"
            prop="last_send_msg"
          >
            <template slot-scope="scope">
              <span v-if="scope.row.last_send_msg > 0">{{
                scope.row.last_send_msg | formatDate("yyyy-mm-dd hh:mm:ss")
              }}</span>
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="150" fixed="right">
            <template slot-scope="scope">
              <el-button type="text" @click="goDetails(scope.row.id)">
                查看
              </el-button>
              <el-dropdown class="m-l-20" @command="handleCommand">
                <el-button type="text"
                  >更多<i
                    aria-hidden="true"
                    class="el-icon-arrow-down el-icon--right"
                  ></i
                ></el-button>
                <el-dropdown-menu slot="dropdown">
                  <el-dropdown-item
                    :command="
                      beforeHandleCommand('detail_' + scope.row.id, scope)
                    "
                    >明细</el-dropdown-item
                  >
                  <el-dropdown-item
                    v-if="scope.row.status == 1"
                    :command="
                      beforeHandleCommand('edit_' + scope.row.id, scope)
                    "
                    class="
                    color-edit"
                    >编辑</el-dropdown-item
                  >
                </el-dropdown-menu>
              </el-dropdown>
            </template>
          </el-table-column>
          <div slot="empty" class="table-empty-label">
            <img
              class="empty-img"
              :src="
                searchLoading
                  ? require('@/assets/images/empty-search.png')
                  : require('@/assets/images/empty-box.png')
              "
              alt=""
            />
            <div>
              {{ searchLoading ? "没有找到相关群发消息" : "暂无群发消息" }}
            </div>
          </div>
        </el-table>
      </template>
      <template slot="footer" v-if="meta.total > 0">
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
import { groupListApi, memberlist } from "@/api/enterprise/association";
import { PaginationDataDefault, ApiResult, Meta } from "@/types";
import { publicMethod } from "@/utils/tools";
@Component({})
export default class extends Vue {
  //----------------data ------------------/
  private search = {
    title: "",
    staff: 0,
    staffList: [{ member_id: 0, name: "全部" }],
    sendStatus: 0,
    sendStatusList: [
      { label: "全部", value: 0 },
      { label: "定时未开始", value: 1 },
      { label: "未发送", value: 2 },
      { label: "部分未发送", value: 3 },
      { label: "已发送", value: 4 },
      { label: "已失败", value: 5 },
      { label: "任务下发中", value: 6 }
    ],
    send_type: 0,
    massStatusList: [
      { label: "全部", value: 0 },
      { label: "群发客户", value: 1 },
      { label: "群发客户群", value: 2 }
    ],
    senfTime: 0
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
  private created() {
    let seaData: any = sessionStorage.getItem("groupleaveData");
    seaData = JSON.parse(seaData);
    if (seaData) {
      this.search.title = seaData.name;
      this.search.sendStatus = seaData.status;
      this.search.staff = seaData.member_id;
      this.search.senfTime = seaData.senfTime;
      this.meta.page = seaData.page;
      this.meta.size = seaData.size;
      sessionStorage.removeItem("groupleaveData");
    }
  }
  private mounted() {
    this.loadData();
    this.memberlistFn();
  }
  //----------------method ------------------/
  // 执行群主下拉列表
  private async memberlistFn() {
    const {
      code,
      message = "",
      data = PaginationDataDefault
    }: ApiResult = await memberlist();
    if (code !== 200) {
      this.$message.error(message);
      return;
    }
    this.search.staffList = [...this.search.staffList, ...data];
  }
  // 点击查看
  private goDetails(id: number) {
    this.groupleaveData();
    this.pathFn(
      "myapp_enterprise_association_group-sent-message_details",
      id + ""
    );
  }
  // 下拉菜单数据
  private beforeHandleCommand(key, row) {
    return {
      key,
      row
    };
  }

  // 点击操作-更多按钮
  private handleCommand(e: any) {
    const command = e.key;
    const type = e.row.row.send_type;
    console.log(e);
    if (command.indexOf("edit") > -1) {
      this.groupleaveData();
      const id = command.replace("edit_", "");
      this.pathFn("myapp_enterprise_association_group-sent-message_edit", id);
    }
    if (command.indexOf("detail") > -1) {
      this.groupleaveData();
      const id = command.replace("detail_", "");
      this.pathFn(
        "myapp_enterprise_association_group-sent-message_group-detail",
        id,
        type
      );
    }
  }
  // 新增群发消息
  private addMessage() {
    this.groupleaveData();
    this.pathFn(
      "myapp_enterprise_association_group-sent-message_add-group-message"
    );
  }
  private groupleaveData() {
    const obj = {
      page: this.meta.page,
      size: this.meta.size,
      name: this.search.title,
      status: this.search.sendStatus,
      member_id: this.search.staff,
      senfTime: this.search.senfTime
    };
    sessionStorage.setItem("groupleaveData", JSON.stringify(obj));
  }
  // 跳转路由
  private pathFn(path: string, queryVal?: string, type?: any) {
    this.$router.push({
      name: path,
      query: {
        id: queryVal,
        type: type
      }
    });
  }
  // 查询
  private clickSearch() {
    this.meta.page = 1;
    this.meta.size = 10;
    this.searchLoading = true;
    this.loadData();
  }
  // 重置
  private reset() {
    this.meta.page = 1;
    this.meta.size = 10;
    this.searchLoading = false;
    this.search.title = "";
    this.search.sendStatus = 0;
    this.search.send_type = 0;
    this.search.staff = 0;
    this.search.senfTime = 0;
    this.loadData();
  }
  // 列表数据
  private async loadData() {
    const { page, size } = this.meta;
    const params = {
      page,
      size,
      name: this.search.title,
      status: this.search.sendStatus,
      member_id: this.search.staff,
      send_type: this.search.send_type,
      start_at: this.search.senfTime
        ? publicMethod(this.search.senfTime[0]) || 0
        : 0,
      end_at: this.search.senfTime
        ? publicMethod(this.search.senfTime[1]) || 0
        : 0
    };
    const {
      code,
      message = "",
      data = PaginationDataDefault
    }: ApiResult = await groupListApi(params);
    if (code !== 200) {
      this.$message.error(message);
      return;
    }
    this.groupList = data.list || [];
    this.meta.total = data.total;
    this.meta.page = data.page;
  }
  // 分页
  handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.loadData();
  }
  handleCurrentChange(val: number) {
    this.meta.page = val;
    this.loadData();
  }
}
</script>
<style lang="scss" scoped>
// /deep/.el-dropdown-menu__item.color-edit {
// color: #606266;
// }
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
}
</style>
