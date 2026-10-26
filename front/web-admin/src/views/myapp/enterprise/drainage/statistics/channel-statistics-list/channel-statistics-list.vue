<template>
  <layout-list>
    <!-- 搜索框内容 -->
    <template slot="search">
      <div class="user-search">
        <search-item text="时间：">
          <el-radio-group
            v-model="timeData"
            @change="timeChange"
          >
            <el-radio-button label="今天"></el-radio-button>
            <el-radio-button label="昨天"></el-radio-button>
            <el-radio-button label="近7天"></el-radio-button>
            <el-radio-button label="近30天"></el-radio-button>
          </el-radio-group>
        </search-item>
        <search-item text="">
          <div class="data-time">
            <el-date-picker
              v-model="timeChoose"
              @change="timeChooseChange"
              type="daterange"
              class="dtimer-choose"
              format="yyyy-MM-dd"
              value-format="timestamp"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              :default-time="['00:00:00', '23:59:59']"
            >
            </el-date-picker>
          </div>
        </search-item>
        <search-item text="渠道：">
          <cy-select
            v-model="channelChoose"
            placeholder="请选择渠道来源"
            class="labelselectd-left"
            search-input
            @change="selectChange"
            :optionsList="channelChooseList"
            :options="{ label: 'name', value: 'id' }"
          >
          </cy-select>
        </search-item>
      </div>
      <div class="search-button">
        <el-button
          class="user-remarksbutton"
          @click="userSearch"
          type="primary"
          size="small"
        >查询</el-button>
        <el-button
          class="user-remarksbuttons"
          @click="resetFn"
          plain
          size="small"
        >重置</el-button>
      </div>
    </template>
    <!-- 客户信息数据统计显示 -->
    <template slot="tabs">
      <div class="customer-module">
        <div class="customer-moduleBox">
          <div
            v-for="item in aiplArr"
            :key="item.group_id"
            class="customer-all"
          >
            <div class="show">
              <div class="show-top">
                {{ item.group_name }}
                <el-popover
                  placement="top-start"
                  trigger="hover"
                  :content="item.group_tips"
                >
                  <svg-icon
                    slot="reference"
                    class="help-icon"
                    icon-class="group-help"
                    style=""
                  ></svg-icon>
                </el-popover>
              </div>
              <div class="number">{{ item.value || 0 }}</div>
            </div>
          </div>
        </div>
      </div>
    </template>
    <!-- 标签页 -->
    <template slot="tabs">
      <el-tabs
        v-model="activeName"
        @tab-click="tabChangeFn"
      >
        <el-tab-pane
          label="按日期查看"
          name="first"
        ></el-tab-pane>
        <el-tab-pane
          label="按渠道查看"
          name="second"
        ></el-tab-pane>
      </el-tabs>
    </template>
    <!-- 导出excel -->
    <template slot="operate">
      <el-button
        size="small"
        class="excel-button"
        :disabled="list.length <= 0"
        @click="statisticsCountExport"
      >导出Excel</el-button>
    </template>
    <!-- 表格 -->
    <template slot="table">
      <el-table
        :data="list"
        stripe
        :cell-style="{ height: '54px' }"
      >
        <!-- isTabShow === 0为按时间查看 -->
        <el-table-column
          prop="time_or_name"
          label="日期"
          v-if="isTabShow === 0"
          min-width="150"
        >
          <template slot-scope="scope">
            <span>{{ scope.row.time_or_name }}</span>
          </template>
        </el-table-column>
        <!-- isTabShow === 1为按渠道查看 -->
        <el-table-column
          prop="time_or_name"
          label="渠道"
          min-width="150"
          v-if="isTabShow === 1"
        >
          <template slot-scope="scope">
            <span>{{ scope.row.time_or_name }}</span>
          </template>
        </el-table-column>
        <el-table-column
          prop="exist_relation_customer"
          align="center"
          width="150"
        >
          <template slot="header">
            <span class="m-r-5">留存好友数</span>
            <el-popover
              placement="top-start"
              trigger="hover"
            >
              <div>
                {{
                  activeName == "first"
                  ? "至当天，仍保留好友关系的外部联系人数"
                  : "至筛选时间的结束时间为止，该渠道仍保留好友关系的外部联系人数"
                }}
              </div>
              <svg-icon
                slot="reference"
                class="help-icon"
                icon-class="group-help"
              ></svg-icon>
            </el-popover>
          </template>
        </el-table-column>
        <el-table-column
          prop="increase_customer_count"
          align="center"
          width="150"
        >
          <template slot="header">
            <span class="m-r-5">新增好友数</span>
            <el-popover
              placement="top-start"
              trigger="hover"
            >
              <div>
                {{
                  activeName == "first"
                  ? "当天，新增的外部联系人加总去重"
                  : "筛选时间内，该渠道新增的外部联系人加总去重"
                }}
              </div>
              <svg-icon
                slot="reference"
                class="help-icon"
                icon-class="group-help"
              ></svg-icon>
            </el-popover>
          </template>
        </el-table-column>
        <el-table-column
          prop="newly_added_customer_count"
          align="center"
          width="150"
        >
          <template slot="header">
            <span class="m-r-5">净增好友数</span>
            <el-popover
              placement="top-start"
              trigger="hover"
            >
              <div>
                {{
                  activeName == "first"
                  ? "当天，新增好友-流失好友"
                  : "筛选时间内，该渠道的新增好友-该渠道的流失好友"
                }}
              </div>
              <svg-icon
                slot="reference"
                class="help-icon"
                icon-class="group-help"
              ></svg-icon>
            </el-popover>
          </template>
        </el-table-column>
        <el-table-column
          prop="loss_customer_count"
          align="center"
          width="150"
        >
          <template slot="header">
            <span class="m-r-5">流失好友数</span>
            <el-popover
              placement="top-start"
              trigger="hover"
            >
              <div>
                {{
                  activeName == "first"
                  ? "当天，跟所有员工解除好友关系的好友数"
                  : "筛选时间内，与该渠道员工解除好友关系的好友数"
                }}
              </div>
              <svg-icon
                slot="reference"
                class="help-icon"
                icon-class="group-help"
              ></svg-icon>
            </el-popover>
          </template>
        </el-table-column>
        <el-table-column
          prop="customer_del_member"
          align="center"
          width="170"
        >
          <template slot="header">
            <span class="m-r-5">好友删除员工数</span>
            <el-popover
              placement="top-start"
              trigger="hover"
            >
              <div>
                {{
                  activeName == "first"
                  ? "当天，好友主动删员工个数"
                  : "筛选时间内，该渠道中好友主动删员工个数"
                }}
              </div>
              <svg-icon
                slot="reference"
                class="help-icon"
                icon-class="group-help"
              ></svg-icon>
            </el-popover>
          </template>
        </el-table-column>
        <el-table-column
          prop="member_del_customer"
          align="center"
          width="170"
        >
          <template slot="header">
            <span class="m-r-5">员工删除好友数</span>
            <el-popover
              placement="top-start"
              trigger="hover"
            >
              <div>
                {{
                  activeName == "first"
                  ? "当天，员工主动删好友个数"
                  : "筛选时间内，该渠道中员工主动删好友个数"
                }}
              </div>
              <svg-icon
                slot="reference"
                class="help-icon"
                icon-class="group-help"
              ></svg-icon>
            </el-popover>
          </template>
        </el-table-column>
        <div
          slot="empty"
          class="table-empty-label"
        >
          <img
            class="empty-img"
            :src="searchLoad
                  ? 'https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/empty-search.png'
                  : require('@/assets/images/empty-box.png')
                "
            alt=""
          />
          <p>{{ searchLoad ? "查询无结果" : "暂无数据" }}</p>
        </div>
      </el-table>
    </template>
    <!-- 分页 -->
    <template slot="footer">
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
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LayoutList from "@/components/public/layout-list/index.vue";
import API from "@/api/enterprise/drainage-coach";
@Component({
  components: {
    LayoutList
  }
})
export default class extends Vue {
  /************************************************* Data开始 *************************************************/
  private timeData = "近7天"; // 时间el-radio-group默认
  private timeChoose: Array<number> | number = 0; // 时间el-date-picker默认
  private channelChoose = 0; // 渠道默认选择全部
  private channelChooseList: Array<any> = []; // 渠道列表
  private channelId = 0; // 渠道ID
  // 客户信息统计数据以及提醒
  private aiplArr: Array<any> = [
    {
      group_name: "留存好友数",
      group_tips: "至筛选时间的结束时间为止，仍保留好友关系的外部联系人数",
      value: 0,
      type: "exist_relation_customer"
    },
    {
      group_name: "新增好友数",
      group_tips: "筛选时间内，新增的外部联系人",
      value: 0,
      type: "increase_customer_count"
    },
    {
      group_name: "净增好友数",
      group_tips: "筛选时间内，新增的外部联系人减去流失好友",
      value: 0,
      type: "newly_added_customer_count"
    },
    {
      group_name: "流失好友数",
      group_tips:
        "筛选时间内，员工删除外部联系人或者外部联系人删除员工，或者互删",
      value: 0,
      type: "loss_customer_count"
    },
    {
      group_name: "好友删除员工数",
      group_tips: "筛选时间内，好友主动删除与员工关系",
      value: 0,
      type: "customer_del_member"
    },
    {
      group_name: "员工删除好友数",
      group_tips: "筛选时间内，员工主动删除与好友关系",
      value: 0,
      type: "member_del_customer"
    }
  ];
  private activeName = "first"; // tabs分页默认
  private isTabShow = 0; // 时间或者渠道名称
  private list: Array<any> = []; // 表格数据
  private searchList: Array<any> = []; // 表格查询结果数据
  private searchLoad = false; // 是否正在查询
  // 分页
  private meta = {
    page: 1,
    size: 10,
    total: 100
  };
  /************************************************* Methods开始 *************************************************/
  //查询
  private userSearch() {
    this.searchLoad = true;
    this.meta.page = 1;
    this.channelsDataSum();
    this.tableData();
  }
  //重置
  private resetFn() {
    this.searchLoad = false;
    this.timeChoose = 0;
    this.channelChoose = 0;
    this.timeData = "近7天";
    this.activeName = "first";
    this.isTabShow = 0;
    this.channelId = 0;
    this.timeChange(this.timeData);
    this.channelsDataSum();
    this.meta = {
      page: 1,
      size: 10,
      total: 0
    };
    this.tableData();
  }

  // 时间el-radio-group变化联动el-date-picker
  private timeChange(value: string) {
    const date = new Date();
    switch (value) {
      case "今天":
        this.getDateRange(date, 0);
        break;
      case "昨天":
        this.getDateRange(date, 1);
        break;
      case "近7天":
        this.getDateRange(date, 6);
        break;
      case "近30天":
        this.getDateRange(date, 29);
    }
  }
  // 计算时间差
  private getDateRange(dateNow: any, intervalDays: number) {
    const oneDayTime = 24 * 3600 * 1000; // 毫秒数
    // 当前时间戳减去一天的毫秒数(6近一周的时间，30近一个月的时间)
    // 过去时间
    const lastDay =
      new Date(new Date().toLocaleDateString()).getTime() -
      intervalDays * oneDayTime;
    // 当前时间
    const nowDays =
      new Date(dateNow.toLocaleDateString()).getTime() + 86400000 * 1 - 1000;
    // const lastDay = new Date(dateNow.getTime() - intervalDays * oneDayTime);
    if (intervalDays === 1) {
      const lastDays =
        new Date(dateNow.toLocaleDateString()).getTime() + 86400000 * 0 - 1000;
      this.timeChoose = [lastDay, lastDays];
    } else {
      this.timeChoose = [lastDay, nowDays];
    }

    return this.timeChoose;
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
  //日期时间选择器
  private timeChooseChange() {
    this.timeData = "";
  }
  // 渠道选择
  private selectChange(itemId: number) {
    this.channelId = itemId;
  }
  // tabs变化
  private tabChangeFn(value: object) {
    this.isTabShow = Number(value["index"]);
    this.meta = {
      page: 1,
      size: 10,
      total: 0
    };
    this.tableData();
  }
  //获取每页数据量
  handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.tableData();
  }
  //获取当前第几页
  handleCurrentChange(val: number) {
    this.meta.page = val;
    this.tableData();
  }
  /************************************************* API接口开始 *************************************************/
  //渠道导出Excel
  private statisticsCountExport() {
    const params = {
      type: this.isTabShow + 1,
      start_time: (this.timeChoose && this.timeChoose[0] / 1000) || undefined,
      end_time: (this.timeChoose && this.timeChoose[1] / 1000) || undefined,
      channels_id: this.channelId || undefined
    };

    API.channelStatisticsCountExportApi(params).then((res: any) => {
      const blob = new Blob([res], {
        type: "application/zip;charset=utf-8"
      });
      let fileName = "";
      const nowTime = new Date();
      fileName = "渠道统计_" + this.formateDate(nowTime) + ".xlsx";

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
  //渠道数据接口API
  private channeData() {
    const params = {
      page: 1,
      size: 9999
    };
    API.channelsListApi(params).then((res: any) => {
      if (res.code == 200) {
        const obj = { id: 0, name: "全部" };
        if (res.data.list) {
          res.data.list.unshift(obj);
          this.channelChooseList = res.data.list;
        }
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 列表数据接口API
  private tableData() {
    parseInt(this.timeChoose[1]);
    const params = {
      start_time: (this.timeChoose && this.timeChoose[0] / 1000) || undefined,
      end_time: (this.timeChoose && this.timeChoose[1] / 1000) || undefined,
      channels_id: this.channelId || undefined,
      page: this.meta.page,
      size: this.meta.size,
      type: this.isTabShow + 1
    };
    API.channelStatisticsListApi(params).then((res: any) => {
      if (res.code === 200) {
        if (res.data) {
          this.list = res.data.list || [];
          this.meta = {
            page: res.data.page,
            size: res.data.size,
            total: res.data.total
          };
        }
      }
    });
  }
  private channelsDataSum() {
    const params = {
      start_time: (this.timeChoose && this.timeChoose[0] / 1000) || undefined,
      end_time: (this.timeChoose && this.timeChoose[1] / 1000) || undefined,
      channels_id: this.channelId || undefined,
      type: this.isTabShow + 1
    };
    API.channelsDataSumApi(params).then((res: any) => {
      if (res.code === 200) {
        if (res.data) {
          this.aiplArr.forEach((item: any) => {
            item.value = res.data[item.type];
          });
          // this.aiplArr[0].value = res.data.exist_relation_customer;
          // this.aiplArr[1].value = res.data.increase_customer_count;
          // this.aiplArr[2].value = res.data.newly_added_customer_count;
          // this.aiplArr[3].value = res.data.loss_customer_count;
          // this.aiplArr[4].value = res.data.customer_del_member;
          // this.aiplArr[5].value = res.data.member_del_customer;
        }
      }
    });
  }
  /************************************************* API接口结束 *************************************************/
  /************************************************* 生命周期 *************************************************/
  private mounted() {
    this.getDateRange(new Date(), 6);
    this.tableData();
    this.channelsDataSum();
    this.channeData();
  }
}
</script>
<style lang="scss" scoped>
.user-search {
  .data-time {
    .dtimer-choose {
      height: 36px;
      padding-top: 0;
      padding-bottom: 0;
    }
  }
}

// .search-button {
//   .user-remarksbutton {
//     width: 68px;
//     height: 36px;
//     background-color: #708df5;
//     color: #ffffff;
//     border: 1px solid transparent;
//     border-color: rgba(112, 141, 245, 0.3);
//   }
//   .user-remarksbuttons {
//     width: 68px;
//     height: 36px;
//     color: #606266;
//     background-color: #ffffff;
//     border: 1px solid #dcdfe6;
//   }
// }
.customer-module {
  margin: 30px auto 0px;

  .customer-moduleBox {
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
    align-items: center;
    padding-left: 0px;
    width: 100%;
    padding-left: 0.5%;

    .customer-all {
      margin-right: 1%;
      list-style: none;
      background-color: #f5f6fa;
      width: 24%;
      height: 137px;
      box-sizing: border-box;
      margin-bottom: 1%;
    }
  }

  div.show {
    border: 1px sold black;
    background: $color-background;
    height: 100%;
    border-radius: 5px;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    position: relative;
    margin: 0 auto;

    /deep/.el-button--mini {
      font-size: 14px;
    }

    span.tips {
      position: absolute;
      top: 10px;
      right: 10px;
      z-index: 1;
      color: $color-main;
      font-size: 12px;
      line-height: 18px;
    }

    .show-top {
      line-height: 15px;
      padding-top: 20px;
      font-size: 14px;
      font-weight: bold;
      color: #303133;

      .help-icon {
        cursor: pointer;
        position: relative;
        top: -1.4px;
      }
    }

    div.number {
      font-size: 26px;
      line-height: 36px;
      font-weight: bold;
      padding: 20px 0;
      color: #303133;
    }
  }
}

.excel-button {
  color: #303133;
  font-size: 14px;
  width: 103px;
}

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

/deep/.slide-box {
  position: absolute;
  top: 200px;
}</style>
