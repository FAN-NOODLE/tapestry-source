<template>
  <div>
    <layout-list>
      <template slot="search">
        <div class="notice">
          <span>群名：</span
          ><span class="group-name">{{ groupData.name }}</span>
        </div>
        <div class="notice notice-nowarp">
          <div>
            <span>群主：</span><span>{{ groupData.owner || "--" }}</span>
          </div>
          <div>
            <span>群ID：</span><span>{{ groupData.chat_id || "--" }}</span>
          </div>
          <div>
            <span>创建时间：</span
            ><span>{{
              groupData.create_time | formatDate("yyyy-mm-dd hh:mm:ss")
            }}</span>
          </div>
        </div>
        <div class="notice">
          <span>群管理员：</span>
          <div>
            <span>{{
              groupData.admins ? groupData.admins.join("|") : "--"
            }}</span>
          </div>
        </div>
        <div class="notice-text">
          <span>群公告：</span><span>{{ groupData.notice || "--" }}</span>
        </div>
        <div class="line"></div>
        <div class="notice">
          <div class="warp-item">
            <span>所属渠道：</span
            ><span>{{ groupData.channel_name || "--" }}</span>
          </div>
          <div class="warp-item">
            <span>群类型：</span
            ><span>{{ groupData.chat_type_name || "--" }}</span>
          </div>
          <div class="warp-item">
            <span>群运营者：</span><span>{{ groupData.operator || "--" }}</span>
          </div>
          <div class="warp-item">
            <span>运营者所属部门：</span
            ><span>{{
              groupData.operator_departments
                ? groupData.operator_departments.join("--")
                : "--"
            }}</span>
          </div>
          <div class="warp-item">
            <span>群目标：</span
            ><span>{{ groupData.chat_target || "--" }}</span>
          </div>
        </div>
      </template>
      <template slot="table">
        <div class="group-num">群成员明细（共{{ groupList.length }}人）</div>
        <div class="export-data">
          <el-button @click="exportExcel">导出</el-button>
        </div>
        <div>
          <!-- stripe  是否为斑马纹 table -->
          <el-table
            tooltip-effect="light"
            :data="showData"
            stripe
            max-height="600"
            style="width:100%;"
          >
            <el-table-column
              prop="name"
              label="群成员"
              align="center"
              min-width="100"
            ></el-table-column>
            <el-table-column
              prop="group_nick_name"
              label="群昵称"
              align="center"
              min-width="100"
            ></el-table-column>
            <el-table-column
              prop="unionid"
              align="center"
              label="微信unionid"
              min-width="130"
            ></el-table-column>
            <el-table-column
              prop="vip_status"
              align="center"
              label="是否会员"
              min-width="100"
            ></el-table-column>
            <el-table-column
              prop="vip_code"
              label="会员ID"
              align="center"
              min-width="100"
            ></el-table-column>
            <el-table-column
              prop="mobile"
              align="center"
              label="手机号"
              min-width="160"
            ></el-table-column>
            <el-table-column
              prop="guiders"
              align="center"
              label="已添加导购"
              min-width="200"
            >
              <template slot-scope="scope">
                <el-popover
                  placement="top-start"
                  width="250"
                  trigger="hover"
                  :disabled="
                    scope.row.guiders && scope.row.guiders.length <= 10
                  "
                >
                  <div v-html="scope.row.guiders"></div>
                  <span
                    v-if="scope.row.guiders.length <= 10"
                    slot="reference"
                    >{{ scope.row.guiders }}</span
                  >
                  <span
                    slot="reference"
                    v-else
                    class="m-l-10 goods-name el-popover-ref"
                    >{{ scope.row.guiders.substr(0, 10) + "..." }}</span
                  >
                </el-popover>
              </template>
            </el-table-column>
            <el-table-column
              prop="join_time"
              align="center"
              label="进群时间"
              min-width="200"
            ></el-table-column>
            <el-table-column
              prop="join_scene"
              align="center"
              label="进群方式"
              width="190"
            ></el-table-column>

            <div slot="empty" v-if="groupList.length == 0">
              <img
                class="empty-img"
                src="@/assets/images/empty-box.png"
                alt=""
              />
              <div class="empty-box">
                <div class="desc">暂无群成员</div>
              </div>
            </div>
          </el-table>
          <el-pagination
            v-if="groupList.length > 10"
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
        </div>
      </template>
      <template slot="footer">
        <el-button @click="$router.go(-1)">返回</el-button>
      </template>
    </layout-list>
  </div>
</template>
<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
import { chatMemberApi, chatInfoApi } from "@/api/enterprise/association";
import dateFormat from "@/filter/FormatDate";
@Component({})
export default class extends Vue {
  // -----------------------------------data
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  private groupList = [];
  private groupData = {};
  // -----------------------------------methods
  // 获取群信息
  private groupDetail() {
    chatInfoApi({
      chat_id: String(this.$route.query.chatId)
    }).then((res: any) => {
      console.log(res);
      if (res.code == 200) {
        this.groupData = res.data;
      }
    });
  }
  // 群成员信息
  private getMemberList() {
    chatMemberApi({
      chat_id: String(this.$route.query.chatId),
      member_type: 0,
      ets: Number(this.$route.query.ets)
    }).then(res => {
      console.log(res);
      if (res.data.items) {
        res.data.items.map((v: any) => {
          const join_scene =
            v.join_scene == 3
              ? "扫二维码入群"
              : v.join_scene == 1 || v.join_scene == 2
              ? `由 ${v.invite_username || "群成员"} 邀请入群`
              : v.join_scene == 0 && "未知";

          v.vip_status = v.vip_status ? "是" : "否";
          v.join_scene = join_scene;
          v.vip_code = v.vip_code || "--";
          v.mobile = v.mobile
            ? v.mobile.slice(0, 3) + "****" + v.mobile.slice(7, 11)
            : "--";
          v.join_time = dateFormat.func(v.join_time, "yyyy-mm-dd hh:mm:ss");
          v.guiders = v.guiders ? v.guiders.join(" | ") : "";
        });
        this.groupList = res.data.items;
        this.meta.total = res.data.items ? res.data.items.length : 0;
      }
    });
  }
  private handleSizeChange(e: number) {
    console.log(e);
    this.meta.page = 1;
    this.meta.size = e;
  }
  private handleCurrentChange(e: number) {
    console.log(e);
    this.meta.page = e;
  }
  // 计算属性
  private get showData() {
    return this.groupList.slice(
      (this.meta.page - 1) * this.meta.size,
      this.meta.page * this.meta.size
    );
  }
  // 生命周期
  private mounted() {
    this.groupDetail();
    this.getMemberList();
    console.log(this.$route.query.chatId);
  }
  // 导出excel
  private exportExcel() {
    import("@/utils/Export2Excel").then(excel => {
      const tHeader: any = [
        "群成员",
        "群昵称",
        "微信unionid",
        "是否会员",
        "会员ID",
        "会员等级",
        "手机号",
        "已添加导购",
        "进群时间",
        "进群方式"
      ];
      const data: any = [];
      this.groupList.map((v: any) => {
        data.push([
          v.name,
          v.group_nick_name,
          v.unionid,
          v.vip_status,
          v.vip_code + "",
          v.vip_level,
          v.mobile,
          v.guiders,
          v.join_time,
          v.join_scene
        ]);
      });
      excel.export_json_to_excel({
        header: tHeader,
        data,
        filename: "群成员列表明细",
        sheetname: "群成员列表明细",
        autoWidth: true,
        bookType: "xlsx"
      });
    });
  }
}
</script>
<style lang="scss">
.line {
  width: 100%;
  height: 0;
  border-top: 1px dashed #aaaaaa;
  margin: 15px 0;
}
.group-num {
  padding: 5px 20px;
  background-color: #f5f6fa;
  border-radius: 4px;
}
.export-data {
  display: flex;
  margin: 15px 0;
  justify-content: flex-end;
}
.el-table__empty-block {
  min-height: 300px !important;
}
.float-right {
  float: right;
}
.notice {
  display: flex;
  flex-wrap: wrap;
  margin-bottom: 20px;
  &:nth-last-child(1) {
    margin-bottom: 0;
  }
  .group-name {
    font-size: 18px;
    font-weight: 600;
  }
}
.notice-nowarp {
  justify-content: space-between;
}
.warp-item {
  margin-right: 100px;
  line-height: 26px;
  &:nth-last-child(1) {
    margin-right: 0;
  }
}
.notice-text {
  line-height: 26px;
}
</style>
