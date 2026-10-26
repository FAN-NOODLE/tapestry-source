<template>
  <layout-list>
    <!-- 头部 -->
    <template slot="search">
      <div class="content-item content-topBox">
        <span>活动类型：{{ activity_info.type_text || "-" }}</span>
        <span>活动标题：{{ activity_info.title || "-" }}</span>
        <span>活动状态：{{ activity_info.status_text || "-" }}</span>
      </div>
      <div>
        <div>
          时间：
          <el-date-picker
            v-model="searchTime"
            type="daterange"
            range-separator="-"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            value-format="yyyy-MM-dd"
            @blur="changeSelector()"
            :picker-options="pickerOptions"
          >
          </el-date-picker>
        </div>
        <div class="top-box" style="width:100%;">
          <div
            v-for="(item, index) in topOneList"
            :key="index"
            :class="[
              'content-box',
              index == topOneList.length - 1 ? '' : 'last-content-box',
            ]"
            style="width:20%;"
          >
            <p>
              <span class="title-name">{{ item.name }}</span>
              <el-popover
                :ref="`popoverA${index}`"
                placement="top"
                width="100"
                trigger="hover"
                visible-arrow="true"
                :content="item.desc"
              >
                <svg-icon
                  slot="reference"
                  class="help-icon"
                  icon-class="group-help"
                ></svg-icon>
              </el-popover>
            </p>
            <p class="sum">{{ statistical[item.key] || "0" }}</p>
          </div>
        </div>
      </div>
    </template>
    <!-- 列表 -->
    <template slot="table">
      <div class="table-title">参与数据明细</div>
      <div style="padding:20px 0px;">
        <search-item text="参与时间：">
          <el-date-picker
            v-model="joinTime"
            type="daterange"
            range-separator="-"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            value-format="yyyy-MM-dd"
            :picker-options="pickerOptions"
          >
          </el-date-picker>
        </search-item>
        <search-item text="来源类型：">
          <cy-select
            v-model="sourceType"
            size="medium"
            :optionsList="sourceStatusList"
            :options="{ label: 'label', value: 'value' }"
          >
          </cy-select>
        </search-item>
        <search-item text="报名审核状态：" v-if="type == 1">
          <cy-select
            v-model="auditType"
            size="medium"
            :optionsList="auditStatusList"
            :options="{ label: 'label', value: 'value' }"
          >
          </cy-select>
        </search-item>
        <search-item text="是否完成签到任务：" v-if="type == 2">
          <cy-select
            v-model="signType"
            size="medium"
            :optionsList="signStatusList"
            :options="{ label: 'label', value: 'value' }"
          >
          </cy-select>
        </search-item>
        <el-button size="small" type="primary" @click="searchEvn">
          查询
        </el-button>
        <el-button class="mar-l-20" size="small" @click="reset">
          重置
        </el-button>
        <el-button
          size="small"
          @click="exportFile"
          :disabled="dataList.length == 0"
        >
          数据导出
        </el-button>
      </div>
      <el-table stripe :data="dataList" class="mar-t-20">
        <el-table-column width="125" prop="nickname" label="参与用户">
        </el-table-column>
        <el-table-column min-width="180" prop="union_id" label="union_id">
        </el-table-column>
        <el-table-column min-width="180" label="参与时间">
          <template slot-scope="scope">
            <span v-if="type == 2">
              {{
                Number(scope.row.updated_at) | formatDate("yyyy-mm-dd hh:mm:ss")
              }}
            </span>
            <span v-else>
              {{
                Number(scope.row.created_at) | formatDate("yyyy-mm-dd hh:mm:ss")
              }}
            </span>
          </template>
        </el-table-column>
        <el-table-column min-width="140" label="来源类型">
          <template slot-scope="scope">
            <span>
              {{ scope.row.source_type === 1 ? "企业社群" : "个人聊天" }}
            </span>
          </template>
        </el-table-column>
        <el-table-column min-width="130" label="导购">
          <template slot-scope="scope">
            <span>
              {{
                (scope.row.member && scope.row.member.name) || scope.row.name
              }}
            </span>
          </template>
        </el-table-column>
        <el-table-column min-width="130" label="所属部门">
          <template slot-scope="scope">
            <el-popover
              placement="top-start"
              width="250"
              trigger="hover"
              :disabled="scope.row.area.length <= 10"
            >
              <div style="overflow:hidden" v-text="scope.row.area"></div>
              <span
                v-if="scope.row.area.length <= 10"
                class="m-l-10 goods-name el-popover-ref"
                slot="reference"
                >{{ scope.row.area }}</span
              >
              <span
                slot="reference"
                v-else
                class="m-l-10 goods-name el-popover-ref"
                >{{ scope.row.area.substr(0, 10) + "..." }}</span
              >
            </el-popover>
          </template>
        </el-table-column>
        <!-- 报名活动 -->
        <el-table-column label="报名审核状态" v-if="type == 1">
          <template slot-scope="scope">
            <span>
              {{
                scope.row.audit_status === 1
                  ? "待审核"
                  : scope.row.audit_status === 2
                  ? "审核通过"
                  : scope.row.audit_status === 3
                  ? "审核未通过"
                  : "--"
              }}
            </span>
          </template>
        </el-table-column>
        <!-- 签到活动 -->
        <el-table-column min-width="130" label="签到天数" v-if="type == 2">
          <template slot-scope="scope">
            <span>
              {{ scope.row.total }}
            </span>
          </template>
        </el-table-column>
        <el-table-column
          label="是否完成签到任务"
          min-width="180"
          v-if="type == 2"
        >
          <template slot-scope="scope">
            <span>
              {{ scope.row.is_complete === 1 ? "否" : "是" }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="zip" label="操作" width="150" fixed="right">
          <template slot-scope="scope">
            <span class="c-blue" @click="viewChildEvn(scope.row.id)"
              >查看详情
            </span>
          </template>
        </el-table-column>
        <div slot="empty" class="table-empty-label">
          <img src="@/assets/images/empty-box.png" alt="" />
          <p>暂无数据</p>
        </div>
      </el-table>
    </template>
    <!-- 分页 -->
    <template slot="footer">
      <el-pagination
        v-if="dataList && dataList.length > 0"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="meta.page"
        :page-sizes="[10, 20, 30, 40]"
        :page-size="meta.size"
        class="page-right font-weight-400"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
      ></el-pagination>
    </template>
    <template slot="footer">
      <el-button class="range-choose-btn" @click="cancel">返回</el-button>
    </template>
  </layout-list>
</template>
<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import Api from "@/api/wisdom/marketing-activities";
import { dateFormat } from "../../../../customer/auto-market/util/index";
@Component
export default class extends Vue {
  @Prop({ default: false }) isGoBack!: boolean;
  @Prop({ default: () => null }) activity_info?: any;
  @Prop({ default: () => null }) statistical?: any;
  @Prop({ default: () => [] }) dataList?: Array<any>;
  @Prop({ default: () => 0 }) total?: number;
  private type = 1; // 活动类型 1、报名 2、签到 3、问卷
  private searchTime: any = "";
  private topStartTime: any = "";
  private topEndTime: any = "";
  private sourceType = 0; //来源
  private auditType = 0; //审核
  private signType = 0; //签到'
  private joinTime: any = "";
  private topOneList = [
    {
      name: "活动PV (次)",
      desc: "筛选时间内，当前活动分享页面被打开的总次数（次数不去重）",
      key: "pv",
    },
    {
      name: "活动UV (人)",
      desc: "筛选时间内，打开当前活动分享页面的用户总人数（人数去重）",
      key: "uv",
    },
    {
      name: "活动参与人数 (人）",
      desc: "筛选时间内，打开当前活动分享页并成功参与的用户总人数（人数去重）",
      key: "jv",
    },
    {
      name: "活动参与社群数 (个）",
      desc:
        "筛选时间内，导购分享活动到社群，被社群成员打开并成功参与活动的社群数（社群数去重）",
      key: "gv",
    },
    {
      name: "活动参与率 (%)",
      desc: "筛选时间内，当前活动的活动参与率（活动参与人数/活动UV人数）*100%",
      key: "join_rate",
    },
  ];
  private sourceStatusList: Array<any> = [
    {
      label: "全部",
      value: 0,
    },
    {
      label: "企业社群",
      value: 1,
    },
    {
      label: "个人聊天",
      value: 2,
    },
  ];
  private auditStatusList: Array<any> = [
    {
      label: "全部",
      value: 0,
    },
    {
      label: "待审核",
      value: 1,
    },
    {
      label: "审核通过",
      value: 2,
    },
    {
      label: "审核未通过",
      value: 3,
    },
  ];
  private signStatusList: Array<any> = [
    {
      label: "全部",
      value: 0,
    },
    {
      label: "是",
      value: 2,
    },
    {
      label: "否",
      value: 1,
    },
  ];
  private list = this.dataList;
  private searchLoading = false;
  private meta: any = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10, // 每页数量
  };
  private pickerOptions = {};
  private mounted() {
    const _this = this;
    this.searchTime = this.getActiveSearchTime();
    this.joinTime = this.getActiveSearchTime();
    this.type = this.activity_info.type;
    this.meta.total = this.total;
    this.pickerOptions = {
      disabledDate(time: any) {
        return (
          time.getTime() <
            _this.activity_info.start_time * 1000 - (24 * 3600 * 1000 - 1) ||
          time.getTime() > _this.activity_info.end_time * 1000
        );
      },
    };
    if (this.isGoBack) {
      const params = JSON.parse(localStorage.getItem("params") || "");
      this.meta.page = params.page;
      this.meta.size = params.page_size;
      this.sourceType = params.source_type || 0;
      this.auditType = params.audit_status || 0;
      this.signType = params.is_finish_sign || 0;
      if (params.start_time !== 0) {
        this.$set(this, "joinTime", [
          dateFormat(params.start_time, '"yyyy-MM-dd"'),
          dateFormat(params.end_time, '"yyyy-MM-dd"'),
        ]);
      } else {
        this.$set(this, "joinTime", "");
      }
      this.getChildTaskList();
      localStorage.removeItem("params");
    }
  }

  private getActiveSearchTime() {
    return [
      this.getActiveDayTime(this.activity_info.start_time),
      this.getActiveDayTime(this.activity_info.end_time),
    ];
  }

  private getActiveDayTime(time: number) {
    const date = time ? new Date(time * 1000) : new Date();
    const year = date.getFullYear().toString();
    const month =
      date.getMonth() + 1 < 10
        ? "0" + (date.getMonth() + 1).toString()
        : (date.getMonth() + 1).toString();
    const day =
      date.getDate() + 1 < 10
        ? "0" + date.getDate().toString()
        : date.getDate().toString();
    return year + "-" + month + "-" + day;
  }

  // 根据时间查询头部的信息
  private changeSelector() {
    const params = {
      type: this.type,
      start_time: new Date(`${this.searchTime[0]} 00:00:00`).getTime() / 1000,
      end_time: new Date(`${this.searchTime[1]} 23:59:59`).getTime() / 1000,
    };
    this.$emit("getHeaderDetail", this.type, params);
  }
  // 查询列表
  private searchEvn() {
    this.meta.page = 1;
    this.meta.total = 0;
    // 是否正在搜索
    this.searchLoading = true;
    this.getChildTaskList();
  }

  handleSizeChange(val: number) {
    this.meta.size = val;
    this.getChildTaskList();
  }
  handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getChildTaskList();
  }
  // 重置
  private reset() {
    // 初始化数据
    this.sourceType = 0;
    this.auditType = 0;
    this.signType = 0;
    this.joinTime = [
      this.getActiveDayTime(this.activity_info.start_time),
      this.getActiveDayTime(this.activity_info.end_time),
    ];
    this.meta.page = 1;
    this.meta.size = 10;
    // 是否正在搜索
    this.searchLoading = false;
    this.getChildTaskList();
  }
  // 数据导出
  private exportFile() {
    const params: any = {
      id: Number(this.$route.query.id),
      page: this.meta.page,
      page_size: this.meta.size,
      start_time:
        (this.joinTime && new Date(`${this.joinTime[0]} 00:00:00`).getTime()) /
          1000 || 0,
      end_time:
        (this.joinTime && new Date(`${this.joinTime[1]} 23:59:59`).getTime()) /
          1000 || 0,
      source_type: this.sourceType,
    };
    if (this.type === 1) {
      params.audit_status = this.auditType;
    }
    if (this.type === 2) {
      params.is_finish_sign = this.signType;
    }
    Api.exportStatistics(params).then((res: any) => {
      const blob = new Blob([res], {
        type: "application/zip;charset=utf-8",
      });
      const fileName =
        this.activity_info.title +
        new Date().toLocaleString().substr(0, 10) +
        ".xlsx";
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
  private getSearchTime() {
    const year = this.searchTime.getFullYear();
    let month = this.searchTime.getMonth() + 1;
    let day = this.searchTime.getDate();
    if (month < 10) {
      month = "0" + month;
    }
    if (day < 10) {
      day = "0" + day;
    }
    const nowDate = `${year}-${month}-${day} 23:59:59`;
    return nowDate;
  }
  // 数据明细
  private getChildTaskList() {
    const params: any = {
      page: this.meta.page,
      page_size: this.meta.size,
      source_type: this.sourceType,
    };
    if (this.joinTime) {
      params.start_time =
        new Date(`${this.joinTime[0]} 00:00:00`).getTime() / 1000;
      params.end_time =
        new Date(`${this.joinTime[1]} 23:59:59`).getTime() / 1000;
    }
    if (this.type === 1) {
      params.audit_status = this.auditType;
    }
    if (this.type === 2) {
      params.is_finish_sign = this.signType;
    }
    this.$emit("getList", this.type, params);
  }
  // 查看用户详情
  private viewChildEvn(id: any) {
    const params: any = {
      page: this.meta.page,
      page_size: this.meta.size,
      start_time:
        (this.joinTime && new Date(`${this.joinTime[0]} 00:00:00`).getTime()) /
          1000 || 0,
      end_time:
        (this.joinTime && new Date(`${this.joinTime[1]} 23:59:59`).getTime()) /
          1000 || 0,
      source_type: this.sourceType,
    };
    if (this.type === 1) {
      params.audit_status = this.auditType;
    }
    if (this.type === 2) {
      params.is_finish_sign = this.signType;
    }
    localStorage.setItem("params", JSON.stringify(params));
    this.$router.push({
      name: "myapp_wisdom_marketing-activities_user-detail",
      query: {
        activityId: this.activity_info.id,
        user_id: id,
      },
    });
  }
  private cancel() {
    this.$router.go(-1);
  }
}
</script>
<style lang="scss" scoped>
.content-item {
  font-size: 16px;
  display: flex;
  justify-content: space-between;
  line-height: 30px;
  margin-bottom: 20px;
}
.content-topBox {
  border-bottom: 1px dashed #009999;
  padding-bottom: 15px;
}
.top-box {
  display: flex;
  justify-content: space-around;
  align-items: center;
  padding: 15px 0px;
  .content-box {
    p {
      text-align: center;
      line-height: 26px;
      font-size: 17px;
    }
  }
  .last-content-box {
    border-right: 1px solid #009999;
  }
}
.table-title {
  width: 99%;
  background: #f2f2f2;
  padding: 8px;
  margin: 10px 0px;
}
.c-blue {
  cursor: pointer;
  color: #708cf6;
  margin-right: 8px;
}
.help-icon {
  margin-left: 5px;
  cursor: pointer;
  width: 13px;
}
.sum {
  font-size: 24px !important;
  font-weight: 600;
  margin-top: 10px;
}
.title-name {
  font-size: 14px !important;
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
</style>
