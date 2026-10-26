<template>
  <div class="compass-ys">
    <layout-list>
      <template slot="search">
        <search-item text="日期：">
          <el-radio-group
            class="search-info m-t-2"
            v-model="time_type"
            size="medium"
            @change="changeTime"
          >
            <el-radio-button
              class="search-data"
              v-for="(item, index) in time_list"
              :key="index"
              :label="item.label"
            >
              {{ item.name }}</el-radio-button>
          </el-radio-group>
          <el-date-picker
            class="m-l-20 choice-picker"
            v-model="times"
            type="datetimerange"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            size="medium"
            :default-time="['00:00:00', '23:59:59']"
            :clearable="false"
            @change="timeChange"
            value-format="timestamp"
          >
          </el-date-picker>
        </search-item>
        <search-item text="导购：">
          <search-select-member
            :selectedList="hign_choose"
            @visible-change="highSearch"
            @remove-member="removeGuide"
          >
          </search-select-member>
        </search-item>
        <div class="operate-btn">
          <el-button
            type="primary"
            size="small"
            @click="search"
          >查询</el-button>
          <el-button
            size="small"
            class="reset"
            @click="reset"
          >重置</el-button>
        </div>
      </template>
    </layout-list>
    <!-- 社群数据 -->
    <layout-card class="margin-top-30">
      <template slot="title">
        <div>社群数据</div>
      </template>
      <template slot="content">
        <div class="privateGuide">
          <div
            class="card"
            v-for="(item, index) in communityData"
            :key="index"
          >
            <div class="title">
              {{ item.name }}
              <el-popover
                trigger="hover"
                :content="item.desc"
              >
                <svg-icon
                  style="margin-left: 5px;"
                  icon-class="help-funnel"
                  slot="reference"
                ></svg-icon>
              </el-popover>
            </div>
            <div class="price">{{ item.value }}</div>
            <svg-icon
              class="icon-tip"
              :icon-class="item.icon"
            ></svg-icon>
          </div>
        </div>
      </template>
    </layout-card>
    <!-- 好友数据 -->
    <layout-card class="margin-top-30">
      <template slot="title">
        <div>好友数据</div>
      </template>
      <template slot="content">
        <div class="privateGuide">
          <div
            class="card"
            v-for="(item, index) in buddyData"
            :key="index"
          >
            <div class="title">
              {{ item.name }}
              <el-popover
                trigger="hover"
                :content="item.desc"
              >
                <svg-icon
                  style="margin-left: 5px;"
                  icon-class="help-funnel"
                  slot="reference"
                ></svg-icon>
              </el-popover>
            </div>
            <div class="price">{{ item.value }}</div>
            <svg-icon
              class="icon-tip"
              :icon-class="item.icon"
            ></svg-icon>
          </div>
        </div>
      </template>
    </layout-card>
    <!-- 社群数据明细 -->
    <layout-card class="margin-top-30 list-content">
      <template slot="title">
        <div>社群数据明细</div>
      </template>
      <template slot="content">
        <el-radio-group
          class="search-info"
          v-model="community_time_type"
          size="small"
          @change="changeCommunityTime"
        >
          <el-radio-button
            class="search-data"
            v-for="(item, index) in time_list"
            :key="index"
            :label="item.label"
          >
            {{ item.name }}</el-radio-button>
        </el-radio-group>
        <el-button
          size="small"
          class="export-btn"
          :disabled="communityList.length === 0"
          @click="exportCommunity"
        >导出Excel</el-button>
        <div class="table-content">
          <el-table
            :data="communityList"
            stripe
            tooltip-effect="light"
          >
            <el-table-column
              fixed="left"
              prop="group_name"
              label="社群名称"
              width="180"
            >
            </el-table-column>
            <el-table-column
              prop="cumulative_group_number"
              label="累计群成员数"
            >
            </el-table-column>
            <el-table-column
              prop="number_of_newcomers"
              label="新入群人数"
            >
            </el-table-column>
            <el-table-column
              prop="leave_group_number"
              label="退群人数"
            >
            </el-table-column>
            <el-table-column
              prop="leave_group_percentage"
              label="退群率"
            >
            </el-table-column>
            <el-table-column
              prop="ba_number"
              label="导购人数"
            >
            </el-table-column>
            <div
              slot="empty"
              class="table-empty-label-page"
            >
              <img
                class="empty-img"
                :src="require('@/assets/images/empty-box.png')"
                alt=""
              />
              <div>暂无社群数据</div>
            </div>
          </el-table>
          <div class="table_footer">
            <el-pagination
              v-if="communityList && communityList.length > 0"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
              :current-page="community_meta.page"
              :page-sizes="[10, 20, 30, 40]"
              :page-size="community_meta.size"
              layout="total, sizes, prev, pager, next, jumper"
              :total="community_meta.total"
            >
            </el-pagination>
          </div>
        </div>
      </template>
    </layout-card>
    <!-- 好友数据明细 -->
    <layout-card class="margin-top-30 list-content">
      <template slot="title">
        <div>好友数据明细</div>
      </template>
      <template slot="content">
        <el-radio-group
          class="search-info"
          v-model="buddy_time_type"
          size="small"
          @change="changeBuddyTime"
        >
          <el-radio-button
            class="search-data"
            v-for="(item, index) in time_list"
            :key="index"
            :label="item.label"
          >
            {{ item.name }}</el-radio-button>
        </el-radio-group>
        <el-button
          size="small"
          class="export-btn"
          :disabled="buddyList.length === 0"
          @click="exportBuddy"
        >导出Excel</el-button>
        <div class="table-content">
          <el-table
            :data="buddyList"
            stripe
            tooltip-effect="light"
          >
            <el-table-column
              prop="shopping_guide_id"
              label="导购ID"
            ></el-table-column>
            <el-table-column
              label="导购"
              min-width="130"
            >
              <template slot-scope="scope">
                <img
                  :src="scope.row.avatar
                    ? scope.row.avatar
                    : 'https://oss-crs.vchangyi.com/wechat/default-avatar-icon.png'
                    "
                  alt=""
                  width="30"
                  height="30"
                  class="table-text-img"
                /><span v-text="scope.row.shopping_guide_name"></span>
              </template>
            </el-table-column>
            <el-table-column
              prop="belonging_store"
              label="所属门店"
            >
            </el-table-column>
            <el-table-column
              prop="new_add_buddy_number"
              label="新增好友数"
            >
            </el-table-column>
            <!-- <el-table-column
              label="新增好友会员绑定率"
              width="110"
            >
              <template slot-scope="scope">
                {{ scope.row.new_buddy_member_binding_percentage || "0%" }}
              </template>
            </el-table-column> -->
            <el-table-column label="好友流失率">
              <template slot-scope="scope">
                {{ scope.row.buddy_churn_percentage || "0%" }}
              </template>
            </el-table-column>
            <el-table-column
              prop="grand_total_number"
              label="累计好友数"
            >
            </el-table-column>
            <el-table-column
              prop="buddy_into_percen_number"
              label="好友入群数"
            >
            </el-table-column>
            <div
              slot="empty"
              class="table-empty-label-page"
            >
              <img
                class="empty-img"
                :src="require('@/assets/images/empty-box.png')"
                alt=""
              />
              <div>暂无好友数据</div>
            </div>
          </el-table>
          <div class="table_footer">
            <el-pagination
              v-if="buddyList && buddyList.length > 0"
              @size-change="secondHandleSizeChange"
              @current-change="secondHandleCurrentChange"
              :current-page="buddy_meta.page"
              :page-sizes="[10, 20, 30, 40]"
              :page-size="buddy_meta.size"
              layout="total, sizes, prev, pager, next, jumper"
              :total="buddy_meta.total"
            >
            </el-pagination>
          </div>
        </div>
      </template>
    </layout-card>
    <!-- 选择导购 -->
    <select-member
      :show.sync="member_show"
      :value="hign_choose"
      @cancel="member_show = false"
      @confirm="submitDialog"
    ></select-member>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LayoutCard from "@/views/myapp/home/components/layout-card/index.vue";
import UserModule from "@/store/modules/user";
// import SelectMember from "@/components/sy-components/select-member.vue";
import SelectMember from "@/views/myapp/wisdom/communication/components/select-member.vue";
import SearchSelectMember from "@/components/sy-components/search-select-member.vue";
import Api from "@/api/wisdom/dataCompass";
@Component({
  components: {
    LayoutCard,
    SelectMember,
    SearchSelectMember,
  },
})
export default class extends Vue {
  // data
  private times: any = "";
  private time_type: any = 0;
  private time_list = [
    { name: "今日", label: 0 },
    { name: "昨日", label: 3 },
    { name: "近7天", label: 1 },
    { name: "近30天", label: 2 },
  ];
  // 计算属性
  get USER() {
    return UserModule.getUserInfo;
  }
  private member_show = false;
  private rangeTabName = 0;
  private hign_choose: Array<any> = [];
  private guide_workers: Array<any> = [];
  // 社群数据
  private communityData: Array<any> = [
    {
      name: "群总数",
      desc: "截止当前查询条件下，有效的群（不包括删除的群）累计总数",
      key: "group_number",
      icon: "stay-fan",
      value: 0,
    },
    {
      name: "累计群成员数",
      desc: "截止当前查询条件下，有效的群（不包括删除的群）累计群成员总数",
      key: "cumulative_group_number",
      icon: "stay-fan",
      value: 0,
    },
    {
      name: "新入群人数",
      desc: "截止当前查询条件下，在查询时间内入群（有效群）的用户的总数",
      key: "number_of_newcomers",
      icon: "new-people",
      value: 0,
    },
    {
      name: "退群人数",
      desc: "截止当前查询条件下，在查询时间累计退群（有效群）的用户的总数",
      key: "leave_group_number",
      icon: "number-of-withdrawals",
      value: 0,
    },
    {
      name: "退群率",
      desc: "截止当前查询条件下，退群人数与退群前累计群成员数的占比",
      key: "leave_group_percentage",
      icon: "percentage",
      value: 0,
    },
    {
      name: "加好友人数",
      desc: "截止当前查询条件下，已添加导购好友且已入群的累计客户数总和",
      key: "add_friend_number",
      icon: "add-friend",
      value: 0,
    },
    {
      name: "加好友率",
      desc:
        "截止当前查询条件下，添加导购好友且已入群的客户数与累计群成员数的占比",
      key: "add_friend_percentage",
      icon: "percentage",
      value: 0,
    },
    // {
    //   name: "微信会员绑定人数",
    //   desc: "截止当前查询条件下，已入群且已绑定微信会员的人数",
    //   key: "wx_member_binding_number",
    //   icon: "wx-bind",
    //   value: 0,
    // },
    // {
    //   name: "微信会员绑定率",
    //   desc: "截止当前查询条件下，群内绑定微信会员人数与累计群成员数的占比",
    //   key: "wx_member_binding_percentage",
    //   icon: "percentage",
    //   value: 0,
    // },
  ];
  // 好友数据
  private buddyData: Array<any> = [
    {
      name: "新增好友数",
      desc: "截止当前查询条件下，新增的好友总人数",
      key: "new_add_buddy_number",
      icon: "add-friend",
      value: 0,
    },
    {
      name: "累计好友数",
      desc: "截止当前查询条件下，导购的累计总好友数",
      key: "grand_total_number",
      icon: "stay-fan",
      value: 0,
    },
    {
      name: "好友入群数",
      desc: "截止当前查询条件下，已入群的好友数",
      key: "buddy_into_number",
      icon: "new-people",
      value: 0,
    },
    {
      name: "好友流失人数",
      desc: "截止当前查询条件下，删除或拉黑好友的人数",
      key: "buddy_churn_number",
      icon: "number-of-withdrawals",
      value: 0,
    },
    {
      name: "好友流失率",
      desc: "截止当前查询条件下，删除或拉黑好友与好友总人数的比例",
      key: "buddy_churn_percentage",
      icon: "percentage",
      value: 0,
    },
    // {
    //   name: "新增好友会员绑定率",
    //   desc: "截止当前查询条件下，新增的好友总人数与新增好友中已绑定微信会员的比例",
    //   key: "new_buddy_member_binding_number",
    //   icon: "percentage",
    //   value: 0,
    // },
    {
      name: "导购总人数",
      desc: "截止当前查询条件下，导购总人数（管理者视角具有该数据）",
      key: "ba_number",
      icon: "wx-bind",
      value: 0,
    },
    // {
    //   name: "微信会员绑定人数",
    //   desc: "截止当前查询条件下，好友中绑定微信会员的人数",
    //   key: "wx_member_binding_number",
    //   icon: "wx-bind",
    //   value: 0,
    // },
    // {
    //   name: "微信会员绑定率",
    //   desc: "截止当前查询条件下，好友中绑定微信会员的人数与全部好友的比例",
    //   key: "wx_member_binding_percentage",
    //   icon: "percentage",
    //   value: 0,
    // },
  ];
  private community_time_type: any = 0;
  private community_time: any = "";
  private communityList = []; // 社群列表数据
  private community_meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10, // 每页数量
  };

  private buddy_time_type: any = 0;
  private buddy_time: any = "";
  private buddyList = []; // 好友列表数据
  private buddy_meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10, // 每页数量
  };
  // 生命周期
  private mounted() {
    this.times = this.getDates(this.time_type);
    this.community_time = this.getDates(this.community_time_type);
    this.buddy_time = this.getDates(this.buddy_time_type);
    this.getData();
    this.getCommunityList();
    this.getBuddyList();
  }
  // 选择自定义时间，清空时间类型
  private timeChange() {
    this.time_type = "";
  }
  // 切换时间类型，清空自定义时间
  private changeTime(data: any) {
    this.times = this.getDates(data);
  }
  //获取时间
  private getDates(type: any) {
    const getDate = (n: number, m: number) => {
      const date = new Date();
      const year = date.getFullYear();
      const month = date.getMonth() + 1;
      const day = date.getDate();
      const nowDay = new Date(`${year}/${month}/${day}`).getTime();
      return [nowDay - 86400000 * n, nowDay - 1000 - 86400000 * m];
    };

    const arr = [
      () => {
        return getDate(0, -1);
      },
      () => {
        return getDate(6, -1);
      },
      () => {
        return getDate(29, -1);
      },
      () => {
        return getDate(1, 0);
      },
    ];
    return arr[type]();
  }

  // 员工输入框即高级搜索
  private highSearch() {
    this.member_show = true;
  }
  // 关闭选择员工弹窗
  private closeMember() {
    this.hign_choose.forEach((item: any) => {
      const idx = item.name.indexOf("(");
      if (idx > -1) {
        item.name = item.name.substring(0, idx);
      }
    });
    this.member_show = false;
  }
  // 高级搜索完成
  private submitDialog(val: any) {
    // if (!val.list || !val.type) return;
    // this.rangeTabName = val.type;
    this.member_show = false;
    this.hign_choose = val.map(item => {
      return {
        ...item,
        name: item.label,
        id: item.user_id
      }
    });
    this.guide_workers = val.map((v: any) => v.user_id);
  }

  // 删除已选导购
  private removeGuide(val: any) {
    const idx = this.hign_choose.findIndex((item: any) => item.user_id === val);
    idx > -1 && this.hign_choose.splice(idx, 1);
    this.guide_workers = this.hign_choose.map((v: any) => v.user_id);
  }
  // 搜索
  private search() {
    this.getData();
  }
  // 重置
  private reset() {
    this.time_type = 0;
    this.times = this.getDates(this.time_type);
    this.rangeTabName = 0;
    this.hign_choose = [];
    this.guide_workers = [];
    this.getData();
  }

  // 获取基本数据
  private getData() {
    const params = {
      start_time: this.times && this.times[0] / 1000,
      end_time: this.times && this.times[1] / 1000,
      // shopping_guide: this.hign_choose.map((item) => item.value),
      user_ids: this.guide_workers,
    };
    Api.communityDataSummary(params).then((res: any) => {
      if (res.code === 200) {
        this.communityData.forEach((item: any) => {
          item.value = res.data.community_data[item.key];
        });
        this.buddyData.forEach((item: any) => {
          item.value = res.data.buddy_data[item.key];
        });
      }
    });
  }
  // 社群列表切换时间
  private changeCommunityTime(data: any) {
    this.community_time = this.getDates(data);
    this.community_meta.page = 1;
    this.getCommunityList();
  }

  // 社群列表每页
  private handleSizeChange(val: any) {
    this.community_meta.page = 1;
    this.community_meta.size = val;
    this.getCommunityList();
  }

  // 社群列表当前页
  private handleCurrentChange(val: any) {
    this.community_meta.page = val;
    this.getCommunityList();
  }

  // 社群列表
  private getCommunityList() {
    const params = {
      page_num: this.community_meta.page,
      page_size: this.community_meta.size,
      start_time: this.community_time && this.community_time[0] / 1000,
      end_time: this.community_time && this.community_time[1] / 1000,
    };
    Api.communityData(params).then((res: any) => {
      if (res.code === 200) {
        this.communityList = res.data.list;
        this.community_meta.total = res.data.total;
      }
    });
  }

  // 导出社群数据
  private exportCommunity() {
    const params = {
      start_time: this.community_time && this.community_time[0] / 1000,
      end_time: this.community_time && this.community_time[1] / 1000,
    };
    Api.communityDataExport(params).then((res: any) => {
      const blob = new Blob([res], {
        type: "application/x-xls;charset=utf-8",
      });
      let fileName = "";
      const nowTime = new Date();
      fileName = "社群数据明细_" + this.formateDate(nowTime) + ".xlsx";

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

  // 好友列表切换时间
  private changeBuddyTime(data: any) {
    this.buddy_time = this.getDates(data);
    this.buddy_meta.page = 1;
    this.getBuddyList();
  }

  // 好友列表每页
  private secondHandleSizeChange(val: any) {
    this.buddy_meta.page = 1;
    this.buddy_meta.size = val;
    this.getBuddyList();
  }

  // 好友列表当前页
  private secondHandleCurrentChange(val: any) {
    this.buddy_meta.page = val;
    this.getBuddyList();
  }

  // 好友列表
  private getBuddyList() {
    const params = {
      page_num: this.buddy_meta.page,
      page_size: this.buddy_meta.size,
      start_time: this.buddy_time && this.buddy_time[0] / 1000,
      end_time: this.buddy_time && this.buddy_time[1] / 1000,
    };
    Api.friendData(params).then((res: any) => {
      if (res.code === 200) {
        this.buddyList = res.data.list;
        this.buddy_meta.total = res.data.total;
      }
    });
  }

  // 导出好友数据
  private exportBuddy() {
    const params = {
      start_time: this.buddy_time && this.buddy_time[0] / 1000,
      end_time: this.buddy_time && this.buddy_time[1] / 1000,
    };
    Api.friendDataExport(params).then((res: any) => {
      const blob = new Blob([res], {
        type: "application/x-xls;charset=utf-8",
      });
      let fileName = "";
      const nowTime = new Date();
      fileName = "好友数据明细_" + this.formateDate(nowTime) + ".xlsx";

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

  // 时间格式转换
  private formateDate(time: any) {
    const year = time.getFullYear();
    let month = time.getMonth() + 1;
    let day = time.getDate();

    month = month < 10 ? "0" + month : month;
    day = day < 10 ? "0" + day : day;

    return year + "-" + month + "-" + day;
  }
}
</script>
<style lang="scss" scoped>
.search-select-member {
  width: 280px;
  display: flex;
  display: inline-block;
  vertical-align: middle;
  align-items: center;
  overflow: hidden;
  border: 1px solid $color-second-line;
  background-color: $color-white;

  /deep/.el-input__inner {
    border: none;
    border-radius: 0;
    color: $color-text-second;
  }

  .select-collapse {
    /deep/.el-input__suffix {
      display: none;
    }
  }
}

/deep/ .search-content {
  margin-bottom: 0 !important;
}

.margin-top-30 {
  margin-top: 30px;
}

.m-t-2 {
  margin-top: -2px !important;
}

.privateGuide {
  display: flex;
  justify-content: flex-start;
  align-items: center;
  margin-top: 20px;
  flex-wrap: wrap;

  .card {
    width: 30%;
    min-width: 30%;
    max-width: 30%;
    margin-right: 30px;
    margin-top: 30px;
    position: relative;

    div {
      line-height: 1.5;
    }

    background: #ffffff;
    border: 1px solid #000000;
    padding: 24px 30px;
    box-sizing: border-box;

    .title {
      font-size: 16px;
      color: #606266;
    }

    .price {
      color: #606266;
      font-weight: bolder;
      font-size: 26px;
      position: relative;
      top: 12px;
    }

    .desc {
      margin-top: 10px;
      color: #606266;
      font-size: 14px;
    }

    .icon-tip {
      width: 92px;
      height: 92px;
      position: absolute;
      bottom: 0;
      right: 0;
    }
  }
}

/deep/ .has-title {
  font-size: 14px !important;
  font-weight: 400 !important;
  margin-bottom: -20px !important;
}

.table-empty-label-page {
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

.list-content {
  /deep/.content {
    margin-top: 40px !important;
  }
}

.export-btn {
  margin-left: 50px;
}

.table-content {
  margin-top: 30px;

  .table_footer {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}

.table-text-img {
  border-radius: 100%;
  vertical-align: middle;
  margin-right: 4px;
}
</style>
