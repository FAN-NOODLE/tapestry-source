<template>
  <div class="supernatant">
    <layout-list>
      <template slot="buttons">
        <el-button
          type="primary"
          @click="goToCreateStrategyFn"
          size="small"
        >创建营销策略</el-button>
      </template>
      <template slot="search">
        <p class="time-show">
          数据更新时间：<span v-if="updateTime">
            {{ updateTime | formatDate("yyyy-mm-dd hh:mm:ss") }}</span><span v-else>--</span>
          <el-button
            class="user-refresh"
            icon="el-icon-refresh"
            size="mini"
            type="text"
            @click="reload"
          >刷新</el-button>
        </p>
        <search-item text="策略名称：">
          <cy-input
            width-size="small"
            size="small"
            placeholder="请输入策略名称"
            prefix-icon="el-icon-search"
            v-model="searchBox.strategyName"
          ></cy-input>
        </search-item>
        <search-item text="客户分组：">
          <cy-select
            width-size="small"
            size="small"
            placeholder="请选择客户分组"
            v-model="searchBox.strategyItem"
            :optionsList="searchBox_opt.strategyItems"
            :options="{ label: 'label', value: 'value' }"
            search-input
          >
          </cy-select>
        </search-item>
        <!-- <search-item text="执行方式：">
          <cy-select
            width-size="small"
            size="small"
            placeholder="请选择执行方式"
            v-model="searchBox.runMode"
            :optionsList="searchBox_opt.runModes"
            :options="{ label: 'label', value: 'value' }"
            @change="runModeChangeFn"
          >
          </cy-select>
        </search-item> -->
        <search-item
          text="执行时间："
          v-if="searchBox.runMode == 2"
        >
          <el-date-picker
            v-model="searchBox.runTime"
            type="datetime"
            placeholder="选择日期时间"
          >
          </el-date-picker>
        </search-item>
        <search-item
          text="执行时间："
          v-else
        >
          <el-date-picker
            v-model="searchBox.runTimeRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            :default-time="['00:00:00', '23:59:59']"
          >
          </el-date-picker>
        </search-item>
        <!-- <search-item text="触达渠道：">
          <cy-select
            v-model="searchBox.ReachWayArr"
            width-size="small"
            size="small"
            placeholder="请选择触达渠道"
            :optionsList="searchBox_opt.ReachWays"
            :options="{ label: 'label', value: 'value' }"
          >
          </cy-select>
        </search-item>
        <search-item text="触达规则：">
          <cy-select
            width-size="small"
            size="small"
            placeholder="请选择触达规则"
            v-model="searchBox.ReachRule"
            :optionsList="searchBox_opt.ReachRules"
            :options="{ label: 'label', value: 'value' }"
          >
          </cy-select>
        </search-item> -->
        <search-item text="策略状态：">
          <cy-select
            width-size="small"
            size="small"
            placeholder="请选择状态"
            v-model="searchBox.status"
            :optionsList="searchBox_opt.statusList"
            :options="{ label: 'label', value: 'value' }"
          >
          </cy-select>
        </search-item>
        <div>
          <el-button
            type="primary"
            @click="searchMarketDataFn"
            size="small"
          >查询</el-button>
          <el-button
            @click="resetFn"
            size="small"
          >重置</el-button>
        </div>
      </template>
      <template slot="table">
        <el-table
          :data="list"
          stripe
        >
          <el-table-column
            prop="strategy_name"
            label="策略名称"
            min-width="310"
            :show-overflow-tooltip="true"
            fixed="left"
          >
          </el-table-column>
          <el-table-column
            prop="customer_group_name"
            label="客户分组"
            min-width="150"
          >
          </el-table-column>
          <el-table-column
            prop="trigger_count"
            label="预计触达人数"
            :show-overflow-tooltip="true"
            min-width="150"
          >
          </el-table-column>
          <el-table-column
            prop="strategy_method_text"
            label="执行方式"
            min-width="150"
          >
          </el-table-column>
          <el-table-column
            label="执行时间"
            min-width="280"
            :show-overflow-tooltip="true"
          >
            <template slot-scope="scope">
              <span>{{
                formatDateFn(scope.row.start_time, scope.row.end_time)
              }}</span>
            </template>
          </el-table-column>
          <el-table-column
            prop="channels_text"
            label="触达渠道"
            min-width="270"
          >
          </el-table-column>
          <el-table-column
            prop="strategy_policy_text"
            label="触达规则"
            min-width="130"
          >
          </el-table-column>
          <el-table-column
            prop="strategy_status_text"
            label="状态"
            min-width="100"
          >
            <template slot-scope="{ row }">
              <el-tag
                :type="row.strategy_status == 1
                      ? ''
                      : row.strategy_status == 2
                        ? 'success'
                        : row.strategy_status == 3
                          ? 'danger'
                          : 'info'
                    "
                :effect="row.strategy_status == 3 ? 'light' : 'plain'"
              >
                {{ row.strategy_status_text }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column
            label="操作"
            width="196"
            fixed="right"
          >
            <template slot-scope="scope">
              <el-button
                type="text"
                size="small"
                @click="goItemInfoFn(scope.row)"
              >数据</el-button>
              <el-button
                type="text"
                size="small"
                @click="
                  lookItemInfoFn(
                    scope.row.id,
                    scope.row.strategy_status,
                    scope.row.channels_text
                  )
                  "
              >查看</el-button>
              <el-dropdown
                v-if="scope.row.strategy_status == 1 ||
                  (scope.row.strategy_status == 2 &&
                    scope.row.strategy_method == 1) ||
                  scope.row.strategy_status == 3
                  "
                class="m-l-20"
                @command="handleCommand"
              >
                <el-button
                  size="small"
                  type="text"
                >更多<i
                    aria-hidden="true"
                    class="el-icon-arrow-down el-icon--right"
                  ></i></el-button>
                <el-dropdown-menu slot="dropdown">
                  <el-dropdown-item
                    v-if="scope.row.strategy_status == 1 ||
                      scope.row.strategy_status == 3
                      "
                    :command="'edit_' +
                      scope.row.id +
                      '_' +
                      scope.row.strategy_status +
                      '_' +
                      scope.row.channels_text
                      "
                    style="color:#000000"
                  >编辑</el-dropdown-item>
                  <el-dropdown-item
                    v-if="scope.row.strategy_status == 1 ||
                      scope.row.strategy_status == 3
                      "
                    :command="'open_' + scope.row.id"
                    class="color-edit"
                  >开始</el-dropdown-item>
                  <el-dropdown-item
                    v-if="scope.row.strategy_status == 2 &&
                      (scope.row.strategy_method != 1 ||
                        scope.row.strategy_method != 3)
                      "
                    :command="'pause_' + scope.row.id"
                    class="color-edit"
                  >暂停</el-dropdown-item>
                  <el-dropdown-item
                    class="color-red"
                    v-if="scope.row.strategy_status == 1"
                    :command="'delete_' + scope.row.id"
                  >删除</el-dropdown-item>
                </el-dropdown-menu>
              </el-dropdown>
              <!-- <el-button
                type="text"
                size="small"
                v-if="
                  scope.row.strategy_status == 1 ||
                    scope.row.strategy_status == 3
                "
                @click="itemStartFn(scope.row.id)"
                >开始</el-button
              >
              <el-button
                type="text"
                size="small"
                v-if="scope.row.strategy_status == 2"
                @click="itemPauseFn(scope.row.id)"
                >暂停</el-button
              >
              <el-button
                type="text"
                size="small"
                v-if="
                  scope.row.strategy_status == 1 ||
                    scope.row.strategy_status == 3
                "
                @click="
                  editItemFn(
                    scope.row.id,
                    scope.row.strategy_status,
                    scope.row.channels_text
                  )
                "
                >编辑</el-button
              >
              <el-button
                type="text"
                size="small"
                style="color:#f56c6c;"
                @click="delRowFn(scope.row.id)"
                v-if="scope.row.strategy_status == 1"
                >删除</el-button
              > -->
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
            <div>{{ searchLoading ? "没有找到该策略" : "暂无策略" }}</div>
          </div>
        </el-table>
      </template>
      <template slot="footer">
        <el-pagination
          @size-change="handleSizeChangeFn"
          @current-change="handleCurrentChangeFn"
          :current-page="currentPage"
          :page-sizes="[10, 20, 30, 40]"
          :page-size="currentSize"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
        >
        </el-pagination>
      </template>
    </layout-list>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LayoutList from "@/components/public/layout-list/index.vue";
import { dateFormat } from "../../util/index";
import API from "@/api/customer";
@Component({
  components: {
    LayoutList
  }
})
export default class extends Vue {
  private updateTime = Date.now() / 1000; //数据更新时间
  private searchLoading = false;
  private searchBox_opt = {
    strategyItems: [],
    runModes: [
      {
        label: "自动长期执行",
        value: 1
      },
      {
        label: "手动定期执行",
        value: 3
      },
      {
        label: "手动立即执行",
        value: 2
      }
    ],
    ReachWays: [
      {
        label: "短信",
        value: 1
      },
      {
        label: "企微群发",
        value: 2
      },
      {
        label: "自动发券",
        value: 3
      }
    ],
    ReachRules: [
      {
        label: "全部渠道",
        value: 2
      },
      {
        label: "仅通知一个渠道",
        value: 1
      }
    ],
    statusList: [
      {
        label: "投放中",
        value: 2
      },
      {
        label: "未开始",
        value: 1
      },
      {
        label: "暂停中",
        value: 3
      },
      {
        label: "已结束",
        value: 4
      }
    ]
  };
  private searchBox = {
    strategyName: "",
    strategyItem: "",
    runMode: "",
    runTimeRange: ["", ""],
    runTime: "",
    ReachWayArr: "",
    status: "",
    ReachRule: ""
  };
  private list: Array<any> = [];
  private currentPage = 1;
  private total = 0;
  private currentSize = 10;
  private packParamsFn() {
    const packParams: object = {
      strategy_name: this.searchBox.strategyName,
      customer_group_id: this.searchBox.strategyItem || 0,
      strategy_method: this.searchBox.runMode || 0,
      start_time:
        this.searchBox.runMode == "3"
          ? new Date(this.searchBox.runTime).getTime() / 1000
          : this.searchBox.runTimeRange
          ? new Date(this.searchBox.runTimeRange[0]).getTime() / 1000 || 0
          : 0,
      end_time:
        this.searchBox.runMode == "3"
          ? new Date(this.searchBox.runTime).getTime() / 1000
          : this.searchBox.runTimeRange
          ? new Date(this.searchBox.runTimeRange[1]).getTime() / 1000 || 0
          : 0,
      strategy_type: 1,
      strategy_status: this.searchBox.status || 0,
      strategy_policy: this.searchBox.ReachRule || 0,
      channels: parseInt(this.searchBox.ReachWayArr),
      page: this.currentPage,
      size: this.currentSize
    };
    return packParams;
  }
  //刷新
  private reload() {
    this.getMarketDataFn();
  }
  formatDateFn(start: any, end: any) {
    if (!start) return "-";
    if (start != end) {
      return (
        dateFormat(start, "yyyy.MM.dd hh:mm") +
        "-" +
        dateFormat(end, "yyyy.MM.dd hh:mm")
      );
    } else {
      return dateFormat(start, "yyyy.MM.dd hh:mm");
    }
  }
  searchMarketDataFn() {
    this.currentPage = 1;
    this.searchLoading = true;
    this.getMarketDataFn();
  }
  async getMarketDataFn() {
    let list: any = [];
    let total: any = 0;
    let updateTime: any;
    try {
      const data: object = this.packParamsFn();
      const res = await API.marketStrategyListApi(data);
      list = res.data.list;
      total = res.data.total;
      updateTime = Date.now() / 1000;
    } catch (error) {
      console.log(error);

      console.warn("获取策略人群信息失败！");
      list = [];
    }
    this.list = list;
    this.updateTime = updateTime;
    this.total = total;
  }

  handleCommand(command: string) {
    const [key, id, status, text] = command.split("_");
    switch (key) {
      case "edit":
        this.editItemFn(id, status, text);
        break;
      case "open":
        this.itemStartFn(id);
        break;
      case "pause":
        this.itemPauseFn(id);
        break;
      case "delete":
        this.delRowFn(id);
        break;
      default:
        break;
    }
  }

  itemStartFn(id: any) {
    this.$confirm("是否需要开始该营销策略", "开始确认", {
      confirmButtonText: "确认开始",
      confirmButtonClass: "button-primary",
      cancelButtonText: "我再想想",
      type: "info"
    })
      .then(async () => {
        const params: object = {
          id: parseInt(id),
          strategy_status: 2
        };
        const res: any = await API.marketingUpdateStatusApi(params);
        if (res.code === 200) {
          this.$message({
            message: "操作成功",
            type: "success",
            duration: 2000
          });
          this.currentPage = 1;
          this.getMarketDataFn();
        } else {
          this.$message.error(res.message);
        }
      })
      .catch(() => {
        this.$message({
          type: "info",
          message: "已取消操作"
        });
      });
  }
  itemPauseFn(id: any) {
    this.$confirm("是否需要暂停该营销策略", "暂停确认", {
      confirmButtonText: "确认暂停",
      confirmButtonClass: "button-primary",
      cancelButtonText: "我再想想",
      type: "info"
    })
      .then(async () => {
        const params: object = {
          id: parseInt(id),
          strategy_status: 3
        };
        const res: any = await API.marketingUpdateStatusApi(params);
        if (res.code === 200) {
          this.$message({
            message: "操作成功",
            type: "success",
            duration: 2000
          });
          this.currentPage = 1;
          this.getMarketDataFn();
        } else {
          this.$message.error(res.message);
        }
      })
      .catch(() => {
        this.$message({
          type: "info",
          message: "已取消操作"
        });
      });
  }
  editItemFn(id: string, status: string, item: string) {
    const param: any = JSON.stringify(this.searchBox);
    this.$router.push({
      name: "myapp_customer_auto-market_market-strategy_edit-stratage",
      query: {
        id: id,
        status: status,
        way: item,
        param,
        page: String(this.currentPage),
        size: String(this.currentSize)
      }
    });
  }
  lookItemInfoFn(id: string, status: string, item: string) {
    const param: any = JSON.stringify(this.searchBox);
    this.$router.push({
      name: "myapp_customer_auto-market_market-strategy_look-stratage",
      query: {
        id: id,
        way: item,
        param,
        page: String(this.currentPage),
        size: String(this.currentSize)
      }
    });
  }
  goItemInfoFn(row: any): void {
    const param: any = JSON.stringify(this.searchBox);
    this.$router.push({
      name: "myapp_customer_auto-market_market-strategy_data-track",
      query: {
        id: row.id,
        way: row.channels_text,
        param,
        page: String(this.currentPage),
        size: String(this.currentSize)
      }
    });
  }
  goToCreateStrategyFn(): void {
    const param: any = JSON.stringify(this.searchBox);
    this.$router.push({
      name: "myapp_customer_auto-market_market-strategy_set-stratage",
      query: {
        param,
        page: String(this.currentPage),
        size: String(this.currentSize)
      }
    });
  }
  delRowFn(id: number | string): void {
    this.$confirm("是否需要删除该营销策略", "删除确认", {
      confirmButtonText: "确认删除",
      confirmButtonClass: "button-danger",
      cancelButtonText: "我再想想",
      type: "warning"
    })
      .then(async () => {
        const params: object = {
          id: Number(id) // id为数字类型
        };
        const res: any = await API.marketStrategyRecordDelApi(params);
        if (res.code === 200) {
          this.$message({
            message: "删除成功",
            type: "success",
            duration: 2000
          });
          this.getMarketDataFn();
        } else {
          this.$message.error(res.message);
        }
      })
      .catch(() => {
        this.$message({
          type: "info",
          message: "已取消删除"
        });
      });
  }
  private resetFn() {
    this.searchBox = {
      strategyName: "",
      strategyItem: "",
      runMode: "",
      runTimeRange: ["", ""],
      runTime: "",
      ReachWayArr: "",
      status: "",
      ReachRule: ""
    };
    this.currentPage = 1;
    this.searchLoading = false;
    this.getMarketDataFn();
  }
  private runModeChangeFn(checked: number) {
    if (checked == 2) {
      this.searchBox.runTimeRange = ["", ""];
    } else {
      this.searchBox.runTime = "";
    }
  }

  handleSizeChangeFn(val: number) {
    this.currentPage = 1;
    this.currentSize = val;
    this.getMarketDataFn();
  }
  handleCurrentChangeFn(val: number) {
    this.currentPage = val;
    this.getMarketDataFn();
  }
  async getStrategyItemsFn() {
    const items: any = [];
    try {
      const res = await API.marketStrategyGroupListApi();
      const list = res.data;
      list.forEach((ele: { name: any; id: any }) => {
        items.push({
          label: ele.name,
          value: ele.id
        });
      });
    } catch (error) {
      console.warn("获取策略人群信息失败！");
    }
    this.searchBox_opt.strategyItems = items;
  }
  created() {
    if (this.$route.params.page) {
      this.searchBox = JSON.parse(this.$route.params.param);
      this.currentPage = +this.$route.params.page;
      this.currentSize = +this.$route.params.size;
    }
    this.getStrategyItemsFn();
    this.getMarketDataFn();
  }
}
</script>
<style lang="scss" scoped>
/deep/.el-dropdown-menu__item.color-red {
  color: $color-danger;
}

.time-show {
  color: $color-text-normal;
  font-size: 14px;
  line-height: 20px;
  margin-bottom: 20px;

  span {
    font-size: 14px;
    color: $color-main;
    line-height: 20px;
    margin-left: 10px;
    cursor: pointer;
  }
}

.table-empty-label {
  display: flex;
  height: 340px;
  flex-direction: column;

  img {
    width: 124px;
    margin: 80px auto 0;
  }

  p {
    line-height: 20px;
    margin-top: 14px;
  }
}

.user-refresh {
  margin-left: 15px;
}

/deep/.el-date-editor {
  height: 36px !important;
}</style>
