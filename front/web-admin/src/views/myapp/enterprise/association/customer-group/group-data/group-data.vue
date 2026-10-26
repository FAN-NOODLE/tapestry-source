<template>
  <div>
    <layout-list>
      <template slot="search"
        ><div class="notice">
          <span>群名：</span
          ><span class="group-name">{{ groupData.name || "--" }}</span>
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
              groupData.admin ? groupData.admin.join("|") : "--"
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
        <div class="group-num">数据概览</div>
        <div class="search-time">
          <div>
            <span>时间：</span>
            <el-date-picker
              style="height:36px"
              v-model="time"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              :default-time="['00:00:00', '23:59:59']"
              :picker-options="startTimeOptSetting"
              value-format="timestamp"
            ></el-date-picker>
            <el-button
              style="height:36px"
              type="primary"
              class="m-l-20"
              @click="searchData"
              >查询</el-button
            >
            <el-button style="height:36px" @click="reSearch">重置</el-button>
          </div>
          <div>当前数据更新时间：{{ refreshTime }}</div>
        </div>
        <div>
          <div class="dashboard-content">
            <div class="dashboard">
              <div class="item">
                <div class="item-title">
                  群留存成员数
                  <el-popover
                    placement="top-start"
                    width="200"
                    trigger="hover"
                    content="当前查询时间下，该社群存当前存在的群成员总数"
                  >
                    <i slot="reference" class="el-icon-question"></i>
                  </el-popover>
                </div>
                <div class="item-num">{{ target.valid_uv | thousandBit }}</div>
              </div>
              <div class="item">
                <div class="item-title">
                  新入群人数
                  <el-popover
                    placement="top-start"
                    width="200"
                    trigger="hover"
                    content="当前查询时间下，该社群新入群的群成员数"
                  >
                    <i slot="reference" class="el-icon-question"></i>
                  </el-popover>
                </div>
                <div class="item-num">{{ target.add_uv | thousandBit }}</div>
              </div>
            </div>
            <div class="dashboard">
              <div class="item">
                <div class="item-title">
                  累计退群人数
                  <el-popover
                    placement="top-start"
                    width="200"
                    trigger="hover"
                    content="群创建日起至当前查询时间下该社群每日退群人数之和"
                  >
                    <i slot="reference" class="el-icon-question"></i>
                  </el-popover>
                </div>
                <div class="item-num">
                  {{ target.all_loss_uv | thousandBit }}
                </div>
              </div>
              <div class="item">
                <div class="item-title">
                  新增退群人数
                  <el-popover
                    placement="top-start"
                    width="200"
                    trigger="hover"
                    content="当前查询时间下，新增退出社群的人数"
                  >
                    <i slot="reference" class="el-icon-question"></i>
                  </el-popover>
                </div>
                <div class="item-num">
                  {{ target.new_loss_uv | thousandBit }}
                </div>
              </div>
              <div class="item">
                <div class="item-title">
                  退群率
                  <el-popover
                    placement="top-start"
                    width="200"
                    trigger="hover"
                    content="当前查询时间下，新增退群人数与退群前累计群留存成员成员数的占比：新增退群人数/（新增退群人数+群留存成员数）*100%"
                  >
                    <i slot="reference" class="el-icon-question"></i>
                  </el-popover>
                </div>
                <div class="item-num">
                  {{
                    target.new_loss_uv + target.valid_uv
                      ? Math.round(
                          (target.new_loss_uv /
                            (target.new_loss_uv + target.valid_uv)) *
                            1000
                        ) / 10
                      : 0
                  }}%
                </div>
              </div>
            </div>
          </div>
          <div class="dashboard-content">
            <div v-if="is_vip" class="dashboard">
              <div class="item">
                <div class="item-title">
                  白金卡及以下用户数
                  <el-popover
                    placement="top-start"
                    width="200"
                    trigger="hover"
                    content="当前查询时间下，群内会员等级为（白金卡、金卡、名媛、GC）的用户数之和"
                  >
                    <i slot="reference" class="el-icon-question"></i>
                  </el-popover>
                </div>
                <div class="item-num">{{ target.glod_down | thousandBit }}</div>
              </div>
              <div class="item">
                <div class="item-title">
                  白金卡及以下用户入群率
                  <el-popover
                    placement="top-start"
                    width="200"
                    trigger="hover"
                    content="当前查询时间下，白金卡及以下用户数在群留存成员中的占比：（白金卡及以下用户数/群留存成员数（去除员工）)*100%"
                  >
                    <i slot="reference" class="el-icon-question"></i>
                  </el-popover>
                </div>
                <div class="item-num">
                  {{
                    target.member_uv + target.general_uv
                      ? Math.round(
                          (target.glod_down /
                            (target.member_uv + target.general_uv)) *
                            1000
                        ) / 10
                      : 0
                  }}%
                </div>
              </div>
            </div>
            <div class="dashboard dashboard-hide">
              <div class="item">
                <div class="item-title">
                  会员总人数
                  <el-popover
                    placement="top-start"
                    width="200"
                    trigger="hover"
                    content="当前查询时间下，群内的会员总人数"
                  >
                    <i slot="reference" class="el-icon-question"></i>
                  </el-popover>
                </div>
                <div class="item-num">{{ target.member_uv | thousandBit }}</div>
              </div>
              <div class="item">
                <div class="item-title">
                  非会员总人数
                  <el-popover
                    placement="top-start"
                    width="200"
                    trigger="hover"
                    content="当前查询时间下，群内非会员的总人数"
                  >
                    <i slot="reference" class="el-icon-question"></i>
                  </el-popover>
                </div>
                <div class="item-num">
                  {{ target.general_uv | thousandBit }}
                </div>
              </div>
              <div class="item">
                <div class="item-title">
                  群会员率
                  <el-popover
                    placement="top-start"
                    width="200"
                    trigger="hover"
                    content="当前查询时间下，群内会员总人数在群留存成员中的占比：（会员总人数/群留存成员数（去除员工））*100%"
                  >
                    <i slot="reference" class="el-icon-question"></i>
                  </el-popover>
                </div>
                <div class="item-num">
                  {{
                    target.member_uv + target.general_uv
                      ? Math.round(
                          (target.member_uv /
                            (target.member_uv + target.general_uv)) *
                            1000
                        ) / 10
                      : 0
                  }}%
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="chart-content">
          <div class="chart-search">
            <div class="chart-tit">会员等级</div>
            <cy-select
              v-model="grade"
              placeholder="全部"
              :optionsList="memberType"
              :options="{ label: 'name', value: 'type' }"
              @change="changePie"
            >
            </cy-select>
          </div>
          <div class="m-t-20">
            <div id="pie-charts" class="pieChartBox"></div>
          </div>
        </div>
        <div class="export-data">
          <el-button @click="downloadExcel">导出</el-button>
        </div>
        <div class="member-table">
          <el-tabs v-model="active" type="card" @tab-click="changeTab">
            <el-tab-pane label="新入群用户明细" name="1">
              <el-table
                tooltip-effect="light"
                :data="showData"
                stripe
                max-height="600"
                style="width:100%;"
              >
                <el-table-column
                  prop="join_time1"
                  label="日期"
                  align="center"
                  min-width="120"
                ></el-table-column>
                <el-table-column
                  prop="invite_username"
                  label="邀请人昵称"
                  align="center"
                  min-width="130"
                ></el-table-column>
                <el-table-column
                  prop="invite_userid"
                  label="邀请人ID"
                  align="center"
                  min-width="130"
                ></el-table-column>
                <el-table-column
                  prop="name"
                  align="center"
                  label="入群成员昵称"
                  min-width="130"
                ></el-table-column>
                <el-table-column
                  prop="unionid"
                  align="center"
                  label="入群成员ID"
                  min-width="130"
                ></el-table-column>
                <el-table-column
                  prop="vip_code"
                  label="会员ID"
                  align="center"
                  min-width="130"
                ></el-table-column>
                <el-table-column
                  prop="vip_level"
                  align="center"
                  label="会员等级"
                  min-width="100"
                ></el-table-column>
                <el-table-column
                  prop="join_time2"
                  align="center"
                  label="入群时间"
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
              </el-table></el-tab-pane
            >
            <el-tab-pane label="退群用户明细" name="2">
              <el-table
                tooltip-effect="light"
                :data="showData"
                stripe
                max-height="600"
                style="width:100%;"
              >
                <el-table-column
                  prop="leave_time1"
                  label="日期"
                  align="center"
                  min-width="120"
                ></el-table-column>
                <el-table-column
                  prop="name"
                  align="center"
                  label="退群成员昵称"
                  min-width="130"
                ></el-table-column>
                <el-table-column
                  prop="unionid"
                  align="center"
                  label="退群成员ID"
                  min-width="130"
                ></el-table-column>
                <el-table-column
                  prop="vip_code"
                  label="会员ID"
                  align="center"
                  min-width="130"
                ></el-table-column>
                <el-table-column
                  prop="vip_level"
                  align="center"
                  label="会员等级"
                  min-width="100"
                ></el-table-column>
                <el-table-column
                  prop="leave_time2"
                  align="center"
                  label="退群时间"
                  min-width="200"
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
              </el-table></el-tab-pane
            >
          </el-tabs>
          <!-- stripe  是否为斑马纹 table -->
          <el-pagination
            v-if="groupList.length > 0"
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
import {
  chatMemberApi,
  chatInfoApi,
  chatDataApi,
  getRefreshTime
} from "@/api/enterprise/association";
import dateFormat from "@/filter/FormatDate";
@Component({
  filters: {
    thousandBit(num: number | string) {
      if (!num) return num;
      num = num + "";
      return (
        num &&
        (num.toString().indexOf(".") !== -1
          ? num.toString().replace(/(\d)(?=(\d{3})+\.)/g, function($0, $1) {
              return $1 + ",";
            })
          : num.toString().replace(/(\d)(?=(\d{3})+$)/g, function($0, $1) {
              return $1 + ",";
            }))
      );
    }
  }
})
export default class extends Vue {
  // -----------------------------------data
  private memberType = [
    { name: "新入群用户", type: "new" },
    { name: "群留存用户", type: "valid" }
  ];
  private grade = "valid";
  private time: any = [];
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  private groupList = [];
  private pie_data = {
    legend: {
      top: "0",
      left: "center"
    },
    // color: ["#71D4D4", "#A3D2f6", "#CEDD79", "#FACFA9", "#C3B6E6"],
    series: [
      {
        type: "pie",
        radius: [0, 180],
        center: ["50%", "55%"],
        itemStyle: {
          borderColor: "#fff",
          borderWidth: 1
        },
        label: {
          formatter: "{tit|{b}（人）}\n{number|{c}}",
          rich: {
            number: {
              fontSize: 14,
              lineHeight: 20,
              align: "center"
            },
            tit: {
              fontSize: 16,
              fontWeight: 600,
              color: "#333",
              align: "center"
            }
          }
        },
        labelLine: {
          show: true,
          length: "15px",
          length2: "30px"
        },
        data: []
      },
      {
        type: "pie",
        radius: [0, 180],
        center: ["50%", "55%"],
        itemStyle: {
          borderColor: "#fff",
          borderWidth: 1
        },
        label: {
          position: "inside",
          formatter: function(paramy: any) {
            return `{number|${paramy.percent.toFixed(1)}%}`;
          },
          rich: {
            number: {
              fontSize: 20,
              color: "#333",
              fontWeight: 600
            }
          }
        },
        labelLine: {
          show: true,
          length: "15px",
          length2: "30px"
        },
        data: []
      }
    ]
  };
  private groupData = {};
  private active = "1";
  private target: any = {};
  private chartData: any = {};
  private is_vip = false;
  // 刷新时间
  private refreshTime = "";
  private startTimeOptSetting = {
    disabledDate(time: Date) {
      const now: number =
        +new Date(new Date(new Date().toLocaleDateString()).getTime()) - 1;
      return time.getTime() > now;
    }
  };
  // -----------------------------------methods
  // 切换查询表格
  private changeTab() {
    this.meta.page = 1;
    this.meta.size = 10;
    this.getMemberList();
  }
  // 查询
  private searchData() {
    this.meta.page = 1;
    this.meta.size = 10;
    this.getPieData();
    this.getMemberList();
  }
  // 切换饼图
  private changePie() {
    this.pie_data.series[0].data = this.chartData[this.grade] || [];
    this.pie_data.series[1].data = this.chartData[this.grade] || [];
    this.init();
  }
  // 获取饼图和总览列表数据
  private getPieData() {
    chatDataApi({
      chat_id: String(this.$route.query.chatId),
      sts: this.time[0] / 1000,
      ets: parseInt(String(this.time[1] / 1000))
    }).then((res: any) => {
      if (res.code == 200) {
        console.log(res);
        this.target = res.data.target;
        this.chartData = res.data.chart;
        this.pie_data.series[0].data = res.data.chart[this.grade] || [];
        this.pie_data.series[1].data = res.data.chart[this.grade] || [];
        this.init();
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 群成员信息
  private getMemberList() {
    chatMemberApi({
      chat_id: String(this.$route.query.chatId),
      member_type: Number(this.active),
      sts: this.time[0] / 1000,
      ets: parseInt(String(this.time[1] / 1000))
    }).then(res => {
      console.log(res);
      if (res.data.items) {
        res.data.items.map((v: any) => {
          if (this.active === "1") {
            const join_scene =
              v.join_scene == 3
                ? "扫二维码入群"
                : v.join_scene == 1 || v.join_scene == 2
                ? `由 ${v.invite_username || "群成员"} 邀请入群`
                : v.join_scene == 0 && "未知";
            v.join_time1 = dateFormat.func(v.join_time, "yyyy-mm-dd");
            v.join_time2 = dateFormat.func(v.join_time, "yyyy-mm-dd hh:mm:ss");
            v.invite_username = v.invite_username || "--";
            v.invite_userid = v.invite_userid || "--";
            v.vip_code = v.vip_code || "--";
            v.vip_level = v.vip_level || "--";
            v.join_scene = join_scene;
          } else {
            v.leave_time1 = dateFormat.func(v.leave_time, "yyyy-mm-dd");
            v.leave_time2 = dateFormat.func(
              v.leave_time,
              "yyyy-mm-dd hh:mm:ss"
            );
            v.vip_code = v.vip_code || "--";
            v.vip_level = v.vip_level || "--";
          }
        });
      }
      this.groupList = res.data.items || [];
      this.meta.total = res.data.items ? res.data.items.length : 0;
    });
  }
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
  private handleSizeChange(e: number) {
    console.log(e);
    this.meta.page = 1;
    this.meta.size = e;
  }
  private handleCurrentChange(e: number) {
    console.log(e);
    this.meta.page = e;
  }
  // 实例图表
  private init() {
    const chart = this.$echarts.init(document.getElementById("pie-charts"));
    chart.setOption(this.pie_data, true);
    window.onresize = function() {
      chart.resize();
    };
  }
  // 重置
  private reSearch() {
    this.time = [
      +new Date(new Date(new Date().toLocaleDateString()).getTime()) - 86400000,
      +new Date(new Date(new Date().toLocaleDateString()).getTime()) - 1
    ];
    this.meta.page = 1;
    this.meta.size = 10;
    this.active = "1";
    this.grade = "valid";
    this.getMemberList();
    this.getPieData();
  }
  // 生命周期
  private created() {
    const enterprise_id = localStorage.getItem("ENTERPRISE_ID") || "0";
    if (enterprise_id == "3314240094228544") {
      this.is_vip = true;
    }
  }
  private mounted() {
    this.reSearch();
    this.groupDetail();
    this.getMemberList();
    this.getPieData();
    // 获取数据更新时间
    this.getRefreshTime();
  }
  // 获取数据更新时间
  private getRefreshTime() {
    getRefreshTime({}).then((res: any) => {
      if (res.code === 200) this.refreshTime = res.data;
    });
  }
  // 计算属性
  private get showData() {
    return this.groupList.slice(
      (this.meta.page - 1) * this.meta.size,
      this.meta.page * this.meta.size
    );
  }
  // 导出excel
  private downloadExcel() {
    import("@/utils/Export2Excel").then(excel => {
      const tHeader = {
        header1: [
          "日期",
          "邀请人昵称",
          "邀请人ID",
          "入群成员昵称",
          "入群成员ID",
          "会员ID",
          "会员等级",
          "入群时间",
          "进群方式"
        ],
        header2: [
          "日期",
          "退群成员昵称",
          "退群成员ID",
          "会员ID",
          "会员等级",
          "退群时间"
        ]
      };
      const data: any = [];
      this.groupList.map((v: any) => {
        if (this.active === "1") {
          data.push([
            v.join_time1,
            v.invite_username,
            v.invite_userid,
            v.name,
            v.unionid,
            v.vip_code + "",
            v.vip_level,
            v.join_time2,
            v.join_scene
          ]);
        } else {
          data.push([
            v.leave_time1,
            v.name,
            v.unionid,
            v.vip_code + "",
            v.vip_level,
            v.leave_time2
          ]);
        }
      });
      excel.export_json_to_excel({
        header: tHeader[`header${this.active}`],
        data,
        filename: this.active === "1" ? "新用户入群明细" : "退群用户明细",
        sheetname: this.active === "1" ? "新用户入群明细" : "退群用户明细",
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
.search-time {
  display: flex;
  justify-content: space-between;
  padding: 0 20px;
  margin: 30px 0;
  line-height: 36px;
}
.dashboard-content {
  display: flex;
  &:nth-child(1) {
    margin-bottom: 20px;
  }
}
.dashboard {
  flex: 3;
  background-color: #f5f6fa;
  border-radius: 4px;
  padding: 40px 20px;
  display: flex;
  &:nth-child(1) {
    flex: 2;
    margin-right: 20px;
  }
  .item {
    flex: 1;
    position: relative;
    text-align: center;
    .item-title {
      font-size: 17px;
    }
    .item-num {
      font-size: 30px;
      font-weight: 600;
      line-height: 40px;
      margin-top: 15px;
    }
  }
}

.dashboard-hide {
  &:nth-child(1) {
    margin-right: 0px;
  }
}
.chart-content {
  padding: 15px;
  border-radius: 4px;
  box-shadow: 0px 0 2px #aaaaaa;
  margin-top: 20px;
  .chart-search {
    display: flex;
    justify-content: space-between;
    .chart-tit {
      font-size: 22px;
      line-height: 36px;
      font-weight: 600;
    }
  }
}
.pieChartBox {
  height: 500px;
  #pieChart-box {
    width: 100%;
    height: 100%;
  }
}
// /deep/ #crs-scrm .member-table .el-tabs__header {
//   margin: 0 !important;
// }
</style>
