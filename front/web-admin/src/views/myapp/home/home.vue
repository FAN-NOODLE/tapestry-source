<template>
  <div class="supernatant">
    <!-- 私域客户 -->
    <layout-card class="m-b-30">
      <template slot="title">
        <div>私域客户</div>
      </template>
      <template slot="content">
        <!-- 企微好友数 -->
        <div class="m-b-30">
          <el-radio-group
            v-model="private_date"
            @change="changePrivateDate"
            class="label-no-padding"
          >
            <el-radio-button :label="private_date_categroy.today">今日</el-radio-button>
            <el-radio-button :label="private_date_categroy.yesterday">昨日</el-radio-button>
            <el-radio-button :label="private_date_categroy.week">近7天</el-radio-button>
            <el-radio-button :label="private_date_categroy.month">近30天</el-radio-button>
          </el-radio-group>
          <div class="m-t-30">
            <div class="second-title m-b-0">
              企微好友数
            </div>
          </div>
          <div class="private-guide">
            <div class="card">
              <div class="title">
                留存好友数
                <el-popover
                  trigger="hover"
                  :content="hover_content.stay_QW_Friend"
                >
                  <svg-icon
                    class="m-l-5"
                    icon-class="help-funnel"
                    slot="reference"
                  ></svg-icon>
                </el-popover>
              </div>
              <div
                class="price"
                v-if="qwFriend.friend"
              >
                {{ qwFriend.friend.re_uv | formatMoneyByDot }}
              </div>
              <div class="desc">
                同时段环比：
                <span v-if="qwFriend.friend">
                  <span v-if="qwFriend.friend.re_uv_percent">
                    {{ qwFriend.friend.re_uv_percent }}
                  </span>
                  <span v-else>
                    --
                  </span>
                </span>
                <span v-if="qwFriend.friend">
                  <svg-icon
                    icon-class="level-up"
                    v-if="qwFriend.friend.re_uv_direction === arrowType.level_up"
                  ></svg-icon>
                  <svg-icon
                    icon-class="level-down"
                    v-else-if="qwFriend.friend.re_uv_direction === arrowType.level_down"
                  ></svg-icon>
                </span>
              </div>
              <svg-icon
                class="icon-tip"
                icon-class="stay-fan"
              ></svg-icon>
            </div>
            <div class="card">
              <div class="title">
                新增好友数
                <el-popover
                  trigger="hover"
                  :content="hover_content.add_QW_Friend"
                >
                  <svg-icon
                    class="m-l-5"
                    icon-class="help-funnel"
                    slot="reference"
                  ></svg-icon>
                </el-popover>
              </div>
              <div
                class="price"
                v-if="qwFriend.friend"
              >
                {{ qwFriend.friend.add_uv | formatMoneyByDot }}
              </div>
              <div class="desc">
                同时段环比：
                <span v-if="qwFriend.friend">
                  <span v-if="qwFriend.friend.add_uv_percent">
                    {{ qwFriend.friend.add_uv_percent }}
                  </span>
                  <span v-else>
                    --
                  </span>
                </span>
                <span v-if="qwFriend.friend">
                  <svg-icon
                    icon-class="level-up"
                    v-if="qwFriend.friend.add_uv_direction === arrowType.level_up"
                  ></svg-icon>
                  <svg-icon
                    icon-class="level-down"
                    v-else-if="qwFriend.friend.add_uv_direction === arrowType.level_down"
                  ></svg-icon>
                </span>
              </div>
              <svg-icon
                class="icon-tip icon-81-81"
                icon-class="newadd-fan"
              ></svg-icon>
            </div>
            <div class="card">
              <div class="title">
                流失好友数
                <el-popover
                  trigger="hover"
                  :content="hover_content.flow_QW_Friend"
                >
                  <svg-icon
                    class="m-l-5"
                    icon-class="help-funnel"
                    slot="reference"
                  ></svg-icon>
                </el-popover>
              </div>
              <div
                class="price"
                v-if="qwFriend.friend"
              >
                {{ qwFriend.friend.loss_uv | formatMoneyByDot }}
              </div>
              <div class="desc">
                同时段环比：
                <span v-if="qwFriend.friend">
                  <span v-if="qwFriend.friend.loss_uv_percent">
                    {{ qwFriend.friend.loss_uv_percent }}
                  </span>
                  <span v-else>
                    --
                  </span>
                </span>
                <span v-if="qwFriend.friend">
                  <svg-icon
                    icon-class="level-up"
                    v-if="qwFriend.friend.loss_uv_direction === arrowType.level_up"
                  ></svg-icon>
                  <svg-icon
                    icon-class="level-down"
                    v-else-if="qwFriend.friend.loss_uv_direction === arrowType.level_down"
                  ></svg-icon>
                </span>
              </div>
              <svg-icon
                class="icon-tip icon-81-81"
                icon-class="flow-fan"
              ></svg-icon>
            </div>
          </div>
        </div>
        <!-- 群成员数 -->
        <div class="m-b-30">
          <div class="second-title">
            群成员数
          </div>
          <div class="flex-c-c">
            <div class="company-num">
              <span style="font-weight: 500;">企微群数</span>
              <el-popover
                trigger="hover"
                :content="hover_content.qw_group_all"
              >
                <svg-icon
                  icon-class="help-funnel"
                  class="help-home-icon"
                  slot="reference"
                ></svg-icon>
              </el-popover>
              <span
                class="m-l-20"
                v-if="qwFriend.group"
              >{{ qwFriend.group.group_num | formatMoneyByDot }}</span>
              <span
                class="m-l-40"
                v-if="qwFriend.group"
              >
                同时段环比
                <span>
                  <span v-if="qwFriend.group.group_num_uv_percent">
                    {{ qwFriend.group.group_num_uv_percent }}
                  </span>
                  <span v-else>
                    --
                  </span>
                </span>
                <span v-if="qwFriend.group">
                  <svg-icon
                    icon-class="level-up"
                    v-if="qwFriend.group.group_num_uv_direction === arrowType.level_up"
                  ></svg-icon>
                  <svg-icon
                    icon-class="level-down"
                    v-else-if="qwFriend.group.group_num_uv_direction === arrowType.level_down"
                  ></svg-icon>
                </span>
              </span>
            </div>
          </div>
          <div class="private-guide">
            <div class="card">
              <div class="title">
                留存群成员
                <el-popover
                  trigger="hover"
                  :content="hover_content.stay_group_Friend"
                >
                  <svg-icon
                    class="m-l-5"
                    icon-class="help-funnel"
                    slot="reference"
                  ></svg-icon>
                </el-popover>
              </div>
              <div
                class="price"
                v-if="qwFriend.group"
              >
                {{ qwFriend.group.re_group_uv | formatMoneyByDot }}
              </div>
              <div class="desc">
                同时段环比：
                <span v-if="qwFriend.group">
                  <span v-if="qwFriend.group.re_group_uv_percent">
                    {{ qwFriend.group.re_group_uv_percent }}
                  </span>
                  <span v-else>
                    --
                  </span>
                </span>
                <span v-if="qwFriend.group">
                  <svg-icon
                    icon-class="level-up"
                    v-if="qwFriend.group.re_group_uv_direction === arrowType.level_up"
                  ></svg-icon>
                  <svg-icon
                    icon-class="level-down"
                    v-else-if="qwFriend.group.re_group_uv_direction === arrowType.level_down"
                  ></svg-icon>
                </span>
              </div>
              <svg-icon
                class="icon-tip"
                icon-class="stay-fan"
              ></svg-icon>
            </div>
            <div class="card">
              <div class="title">
                新增群成员
                <el-popover
                  trigger="hover"
                  :content="hover_content.add_group_Friend"
                >
                  <svg-icon
                    class="m-l-5"
                    icon-class="help-funnel"
                    slot="reference"
                  ></svg-icon>
                </el-popover>
              </div>
              <div
                class="price"
                v-if="qwFriend.group"
              >
                {{ qwFriend.group.add_group_uv | formatMoneyByDot }}
              </div>
              <div class="desc">
                同时段环比：
                <span v-if="qwFriend.group">
                  <span v-if="qwFriend.group.add_group_uv_percent">
                    {{ qwFriend.group.add_group_uv_percent }}
                  </span>
                  <span v-else>
                    --
                  </span>
                </span>
                <span v-if="qwFriend.group">
                  <svg-icon
                    icon-class="level-up"
                    v-if="qwFriend.group.add_group_uv_direction === arrowType.level_up"
                  ></svg-icon>
                  <svg-icon
                    icon-class="level-down"
                    v-else-if="qwFriend.group.add_group_uv_direction === arrowType.level_down"
                  ></svg-icon>
                </span>
              </div>
              <svg-icon
                class="icon-tip icon-81-81"
                icon-class="newadd-fan"
              ></svg-icon>
            </div>
            <div class="card">
              <div class="title">
                流失群成员
                <el-popover
                  trigger="hover"
                  :content="hover_content.flow_group_Friend"
                >
                  <svg-icon
                    class="m-l-5"
                    icon-class="help-funnel"
                    slot="reference"
                  ></svg-icon>
                </el-popover>
              </div>
              <div
                class="price"
                v-if="qwFriend.group"
              >
                {{ qwFriend.group.loss_group_uv | formatMoneyByDot }}
              </div>
              <div class="desc">
                同时段环比：
                <span v-if="qwFriend.group">
                  <span v-if="qwFriend.group.loss_group_uv_percent">
                    {{ qwFriend.group.loss_group_uv_percent }}
                  </span>
                  <span v-else>
                    --
                  </span>
                </span>
                <span v-if="qwFriend.group">
                  <svg-icon
                    icon-class="level-up"
                    v-if="qwFriend.group.loss_group_uv_direction === arrowType.level_up"
                  ></svg-icon>
                  <svg-icon
                    icon-class="level-down"
                    v-else-if="qwFriend.group.loss_group_uv_direction === arrowType.level_down"
                  ></svg-icon>
                </span>
              </div>
              <svg-icon
                class="icon-tip icon-81-81"
                icon-class="flow-fan"
              ></svg-icon>
            </div>
          </div>
        </div>
      </template>
    </layout-card>
    <!-- 导购业绩排行榜 -->
    <layout-card class="m-b-30">
      <template slot="title">
        <div>导购排行榜</div>
      </template>
      <template slot="content">
        <el-radio-group
          v-model="shoppingGuideDate"
          @change="refreshTable"
          class="label-no-padding"
        >
          <el-radio-button :label="private_guide_date.today">今日</el-radio-button>
          <el-radio-button :label="private_guide_date.yesterday">昨日</el-radio-button>
          <el-radio-button :label="private_guide_date.week">近7天</el-radio-button>
          <el-radio-button :label="private_guide_date.month">近30天</el-radio-button>
        </el-radio-group>
        <el-button
          size="small"
          class="export-btn"
          :disabled="shopperList.length === 0"
          @click="exportShopper"
        >导出Excel</el-button>
        <el-table
          :header-cell-style="{ background: '#F5F7FA' }"
          ref="shopperGuideTable"
          class="border-btm-none shopper-guide-table"
          :data="shopperList"
          show-summary
          height="610"
          :summary-method="getSummaries"
          stripe
          @sort-change="sortChange"
          :default-sort="{ prop: 'yj_fee', order: 'descending' }"
        >
          <!-- <el-table-column label="业绩排名" prop="id"> </el-table-column> -->
          <el-table-column
            label="导购"
            prop="shopping_guide_name"
          > </el-table-column>
          <el-table-column
            label="企微好友数"
            sortable
            prop="qw_uv"
            align="center"
          >
            <template slot="header">
              <el-popover
                placement="top-start"
                trigger="hover"
                :content="table_hover_content.enterprise_friend"
              >
                <span slot="reference">企微好友数</span>
              </el-popover>
            </template>
            <template slot-scope="scope">
              <span> {{ scope.row.num }} </span>
              <span class="inline-block-width-25"></span>
            </template>
          </el-table-column>
          <el-table-column
            label="门店"
            prop="belonging_store"
          > </el-table-column>
          <div
            slot="empty"
            class="table-empty-firstlabel"
          >
            <img
              src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/default/default_noList.png"
              alt=""
            />
            <p>暂无数据</p>
          </div>
        </el-table>
      </template>
    </layout-card>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LayoutCard from "./components/layout-card/index.vue";
import NewIntroduction from "./components/new-introduction.vue";
import MarketService from "./components/market-service.vue";
import UserModule from "@/store/modules/user";
import ShareModule from "@/shared/index";
// import versionPrompt from "./components/version-prompt.vue";
import Api from "@/api/home/home";
import { timeForMat, publicMethod } from "@/utils/tools";
interface SortDate {
  name: string;
  age: number;
}
interface IPurchase {
  user_count: number;
  days: number;
}
@Component({
  components: {
    LayoutCard,
    NewIntroduction,
    MarketService
    // versionPrompt,
  }
})
export default class extends Vue {
  //获取用户信息
  get USER() {
    return UserModule.getUserInfo;
  }
  get SHARE() {
    return ShareModule.getShared();
  }
  //告警阈值控制
  private isShowNotice = false;
  private strockNum = 0;
  private checkList = [];
  // 页面变量begin
  private sort_change = {
    ascending: 1, //升序
    descending: 0 //降序
  };
  private priceLength = {
    //定义金钱样式更改最大值
    five: 99999
  };
  //底部年份
  private year = new Date().getFullYear();
  //定义同时段环比箭头方向
  private arrowType = {
    level_up: 1, //向上红色箭头
    level_down: 2 //向下绿色箭头
  };
  //小程序默认授权
  private disposeState = {
    on: 1, //已配置
    off: 2 //未配置
  };
  //导购排行榜hover解释
  private table_hover_content = {
    shop_guide: "归属该导购的客户的实收金额-归属该导购的客户的退款金额",
    contribute_per: "导购业绩/总业绩*100%",
    bind_custom: "跟该导购有业绩绑定关系也是企微好友",
    enterprise_friend: "该导购的留存的企业微信好友数量",
    discount: "归属该员工的订单优惠券面额/含券支付金额*100%"
  };
  //页面所有问号hover提示信息
  private hover_content = {
    //企微好友数
    stay_QW_Friend: "筛选时间内，仍保留好友关系的外部联系人数",
    add_QW_Friend: "筛选时间内，新增的外部联系人",
    flow_QW_Friend: "筛选时间内，员工删除外部联系人或者外部联系人删除员工，或者互删",
    // 群成员数
    qw_group_all: "筛选时间内，仍保留的企微群，每个企微员工企微群加总",
    stay_group_Friend: "筛选时间内，每个企微员工仍存在的企微群成员加总去重",
    add_group_Friend: "筛选时间内，每个企微员工企微群成员中入群人数加总去重",
    flow_group_Friend: "筛选时间内，外部联系人主动退出群或者员工将外部联系人移除出群",
    //公众号粉丝数
    stay_fans: "筛选时间内，仍保留关注公众号的粉丝数",
    add_fans: "筛选时间内，新增关注公众号的粉丝数",
    flow_fans: "筛选时间内，取消关注公众号的粉丝",
    //商城变现
    purchase: "所有已支付订单的商家端实收金额",
    refundAmount: "所有退款完成的订单金额，不包含在途退款订单",
    payordernum: "所有已支付完成的订单数",
    customPrice: "实收金额/支付买家数",
    visitor: "访问小程序的独立客户数",
    visitor_buyer: "所有已提交订单的客户数",
    visitor_payer: "所有已支付完成的客户数",
    conversion_cycle: "客户首次访问到首次支付的平均天数",
    repurchase_cycle: "成交2次及以上的客户，复购的平均天数",
    visit_order_inversion: "下单买家数/访客数*100%",
    order_pay_inversion: "支付买家数/下单买家数*100%",
    visit_pay_inversion: "支付买家数/访客数*100%",
    all_customer: "商城历史至今的所有访客加总去重",
    all_pay_customer: "商城历史至今的所有支付买家数加总去重",
    all_fg_customer: "商城历史至今的所有支付次数>=2次的客户数",
    all_visit_pay_inversion: "总支付买家数/总访客数*100%",
    all_fg_pay_inversion: "总复购客户/总支付买家数*100%"
  };
  // interface优化
  //企微好友数默认筛选今日
  private private_date = 0;
  private private_date_categroy = {
    today: 0, //今日
    yesterday: 1, //昨日
    week: 7, //近7天
    month: 30 //近30天
  };
  //导购业绩排行榜默认筛选今日
  private shoppingGuideDate = 0;
  private private_guide_date = {
    today: 0, //今日
    yesterday: 1, //昨日
    week: 7, //近7天
    month: 30 //近30天
  };
  //商城变现默认筛选今日
  private shop_date = 0;
  private private_shop_date = {
    today: 0, //今日
    yesterday: 1, //昨日
    week: 7, //近7天
    month: 30 //近30天
  };
  //公众号
  private fansData = [];
  private fansCheck: any = "";
  private value = "";
  // 导购表格数据
  private shopperList: Array<any> = [];
  private shopperTotal: any = {};
  //首页基本配置
  private baseConfig: any = {};
  //私域客户数据
  private qwFriend: any = {};
  //私域粉丝数据
  private qwFans: any = {};
  //商城变现数据
  private shopFee: any = {};
  private currentStoreName = "全部门店";
  //商城总转化数据
  private allShopFee: any = {};
  //门店列表
  private shopLists: Array<any> = [];
  private storeId = 0; //默认获取全部门店 0
  private version = {
    old_store_shop: 1, //旧版商城
    store_shop: 2, //中心商城有门店
    store_noshop: 3, //中心商城无门店
    regular_chain_store: 4 //直营连锁门店
  };
  private store_platform = 0; //代表商城门店
  private purchase_data: IPurchase = {
    user_count: 0,
    days: 0
  };
  // 获取账号购买信息
  private getBuyNotice() {
    Api.buyNotice().then((res: any) => {
      const { user_count, days } = res.data;
      Object.assign(this.purchase_data, { user_count, days });
      // (user_count < 1 || !user_count) && this.isNeedShowAlarmNotice();
    });
  }
  // 页面变量end
  private shopchange(data: any) {
    this.storeId = data;
    this.shopLists.forEach((item: any) => {
      if (item.id === data) {
        this.currentStoreName = item.short_name;
      }
    });
    this.getShopFee({
      day: this.shop_date,
      store_id: this.storeId
    });
  }
  private changePrivateDate(data: any) {
    const params = {
      day: data
    };
    this.getQwFriend(params);
  }
  private sortChange(column: any) {
    console.log(column);
    let direction: any;
    let params: any;
    const columnName = column.prop;
    if (column.order === "ascending") {
      direction = this.sort_change.ascending;
    } else if (column.order === "descending") {
      direction = this.sort_change.descending;
    }
    switch (columnName) {
      case "yj_fee":
        params = {
          yj_sort: direction
        };
        this.getShoperList(params);
        break;
    }
  }
  private getQwFriend(data: any) {
    Api.getQwFriend(data).then((res) => {
      if (res.data) {
        this.qwFriend = res.data;
        if (res.data.fans) {
          this.qwFans = res.data.fans[0] && res.data.fans[0].v ? res.data.fans[0].v : [];
          this.fansCheck = res.data.fans[0] && res.data.fans[0].name ? res.data.fans[0].name : "";
          this.fansData = res.data.fans;
        }
      }
    });
  }
  private getShopFee(data: any) {
    Api.getShopFee(data).then((res) => {
      if (res.data) {
        this.shopFee = res.data;
      }
    });
  }
  private getShoperList(data: any) {
    const params = {
      day: this.shoppingGuideDate
    };
    if (data) {
      Object.assign(params, data);
    }
    Api.getShopers(params).then((res) => {
      if (res.data && res.data.list && res.data.list.length > 0) {
        this.shopperTotal = res.data.total;
        this.shopperList = res.data.list;
      }
    });
  }
  private refreshTable() {
    this.getGuideRank();
  }
  //导购业绩表格总计
  private getSummaries() {
    this.$nextTick(() => {
      (this.$refs.shopperGuideTable as any).doLayout();
    });
    const shopperSum: Array<any> = [];
    shopperSum[0] = "合计";
    shopperSum[1] = `${this.shopperList.length}个导购`;
    shopperSum[2] = this.shopperTotal.yj_fee;
    shopperSum[3] = this.shopperTotal.contribution;
    shopperSum[4] = this.shopperTotal.uv;
    shopperSum[5] = this.shopperTotal.qw_uv;
    shopperSum[6] = this.shopperTotal.discount ? this.shopperTotal.discount : "0%";
    return shopperSum;
  }
  private getNormalCheck() {
    Api.getNormalCheck({}).then((res) => {
      this.baseConfig = res.data;
    });
  }
  // 获取导购业绩排行榜
  private async getGuideRank() {
    const start_time: number = publicMethod(timeForMat(this.shoppingGuideDate).t2);
    const end_time: number = publicMethod(timeForMat(this.shoppingGuideDate).t1);
    const params: any = {
      type: 2, // 2表示查100条
      start_time,
      end_time
    };
    const res: any = await Api.customerRankData(params);
    if (res.code === 200) {
      this.shopperList = res.data.list;
    }
  }
  private goPurchase() {
    const protocol = "http";
    window.open(
      `${protocol}://xueyuan.vchangyi.com/bin/view/CRS%E5%8A%9F%E8%83%BD%E8%AF%B4%E6%98%8E/%E6%99%BA%E6%85%A7%E5%AF%BC%E8%B4%AD/%E6%99%BA%E6%85%A7%E5%AF%BC%E8%B4%AD%E7%AE%A1%E7%90%86%E5%90%8E%E5%8F%B0/%E6%99%BA%E6%85%A7%E5%AF%BC%E8%B4%AD%E4%BB%98%E8%B4%B9%E8%AE%A2%E9%98%85/`,
      "_blank"
    );
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

  // 导出导购排行榜
  private exportShopper() {
    const start_time: number = publicMethod(timeForMat(this.shoppingGuideDate).t2);
    const end_time: number = publicMethod(timeForMat(this.shoppingGuideDate).t1);
    const params: any = {
      // type: 1,
      type: 2, // 2表示查100条
      start_time,
      end_time
    };
    Api.shopperListExport(params).then((res: any) => {
      const blob = new Blob([res], {
        type: "application/x-xls;charset=utf-8"
      });
      let fileName = "";
      const nowTime = new Date();
      fileName = "导购排行榜_" + this.formateDate(nowTime) + ".xlsx";

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
  mounted() {
    this.getBuyNotice();
    this.getNormalCheck();
    //获取私域客户数据
    this.getQwFriend({
      day: 0
    });
    //  获取导购业绩排行
    this.getGuideRank();
  }
}
</script>
<style lang="scss" scoped>
.supernatant {
  min-width: 700px;
}

.overdue-tip {
  margin-bottom: 30px;
  font-size: 14px;
  background: #fef1f1;
  width: 100%;
  border-radius: 4px;
  padding: 20px;
  display: flex;
  box-sizing: border-box;
  justify-content: center;
  align-items: center;

  .el-icon-warning {
    color: #f67270;
  }

  p:first-child {
    color: #303133;
  }

  p:last-child {
    color: #606266;
    margin-left: 22px;
  }

  .paybtn {
    background: $color-main;
    color: #fff;
    font-size: 14px;
    margin-left: 20px;
  }
}

.alarm-notice-container {
  display: flex;
  justify-content: space-between;
  margin-bottom: 30px;
  color: #606266;
  width: 100%;
  height: 50px;
  background: #fff;
  border-radius: 4px;
  line-height: 50px;
  padding: 0 20px;
  box-sizing: border-box;

  .alarm-type {
    display: inline-block;
    padding-right: 5px;
    padding-left: 5px;
    color: #f6ab3c;
  }

  .btn-group {
    margin-left: 20px;

    span {
      cursor: pointer;
      display: inline-block;
    }

    span:first-child {
      padding-right: 20px;
    }
  }

  .right-flex {
    display: flex;

    .el-icon-close {
      font-size: 20px;
      height: 50px;
      display: inline-block;
      line-height: 50px;
      color: #dcdfe6;
      margin-left: 30px;
      cursor: pointer;
    }
  }

  .el-icon-warning {
    color: #f6ab3c;
    position: relative;
    top: 1px;
    font-size: 16px;
  }
}

.no-auth-tip {
  margin-bottom: 30px;
  color: #606266;
  width: 100%;
  height: 50px;
  text-align: center;
  background: #fef7ec;
  border-radius: 4px;
  line-height: 50px;

  .el-icon-warning {
    color: #f6ab3c;
    width: 16px;
    height: 16px;
  }
}

.suspension {
  width: 100%;
  height: 87px;
  background-color: #fff7e6;

  .suspension-content {
    max-width: 1200px;
    padding: 20px 0;
    margin: 0 auto;
    margin-left: 30px;

    .icon {
      vertical-align: top;
    }

    .content-text {
      display: inline-block;
      margin-left: 12px;
      margin-top: -2px;

      .title {
        line-height: 22px;
        font-weight: 600;
        margin-bottom: 5px;
        color: #303133;
        font-size: 16px;
      }

      .text {
        font-size: 14px;
        color: #606266;
        line-height: 20px;
        font-weight: 400;
      }
    }

    .purchase {
      vertical-align: top;
      margin-left: 20px;
      margin-top: 7px;
    }
  }
}

.crs-footer {
  color: #909399;
  text-align: center;

  .footer-text {
    margin-bottom: 11px;
  }

  .dev-member {
    width: 90px;
    margin: 0 auto;
    cursor: pointer;

    .dev-icon {
      width: 15px;
      height: 15px;
      vertical-align: middle;
    }
  }

  .footer-text,
  .dev-member {
    line-height: 17px;
    font-size: 12px;

    span {
      vertical-align: middle;
    }
  }
}

.shopper-guide-table {
  width: 100%;
  margin: 20px 0 0 0;

  ::v-deep {
    th {
      padding: 0 !important;
      margin: 0 !important;
      height: 64px !important;
      line-height: 64px !important;
    }

    .el-table__header {
      th {
        height: 36px !important;
        line-height: 36px !important;
      }
    }

    td {
      padding: 0 !important;
      margin: 0 !important;
      height: 64px !important;
      line-height: 64px !important;
    }

    .el-table__footer-wrapper {
      td:nth-child(1) {
        padding: 0 !important;
      }

      td:nth-child(2) {
        padding: 0 !important;
      }

      td {
        padding-right: 25px !important;
      }
    }
  }
}

.shop-cash {
  display: flex;
  justify-content: justify-content;
  flex-wrap: wrap;

  .shop-left,
  .shop-right {
    padding: 30px;
    background: #f5f6fa;
    border-radius: 8px;
    margin-bottom: 30px;
  }

  .shop-left {
    flex: 6;
    margin-right: 30px;

    .shop-left-title {
      color: #606266;
      font-size: 14px;
      margin-bottom: 20px;
      font-weight: bolder;
    }

    .body-card {
      display: flex;
      flex-wrap: nowrap;
      justify-content: space-around;
      align-items: center;
      border: 1px solid #dcdfe6;
      padding: 20px;
      border-radius: 10px;

      .single-card {
        color: #606266;
        box-sizing: border-box;

        p:first-child {
          display: flex;
          align-items: center;

          span:first-child {
            margin-right: 5px;
          }
        }

        p:last-child {
          margin-top: 10px;
        }
      }
    }

    .bottom {
      margin-top: 40px;
    }

    .funnel {
      height: 100px;
      display: flex;
      justify-content: flex-start;
      align-items: center;
      margin-bottom: 10px;

      .funnel-png-container {
        position: relative;
        width: 200px;
        padding: 0 10px;
        display: flex;
        align-items: center;
        justify-content: center;

        .funnel-png {
          position: relative;

          .desc {
            width: 100%;
            position: absolute;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);

            .funnel-title {
              color: #fff;
              display: flex;
              justify-content: center;
              align-items: center;
            }
          }
        }

        .funnel-png-1 {
          width: 200px;
          height: 100%;

          img {
            width: 100%;
            height: 100%;
          }
        }

        .funnel-png-2 {
          width: 166px;
          height: 100%;

          img {
            width: 100%;
            height: 100%;
          }
        }

        .funnel-png-3 {
          width: 132px;
          height: 100%;

          img {
            width: 100%;
            height: 100%;
          }
        }

        .inversion-line {
          position: absolute;
          width: 110px;
          height: 90px;
          top: 26px;
          right: -125px;
          border-top: 2px solid #dcdfe6;
          border-right: 2px solid #dcdfe6;
          border-top-right-radius: 6px;
          z-index: 2;

          .order-inversion {
            width: 40px;
            height: 20px;
            border-top: 2px solid #dcdfe6;
            border-right: 2px solid #dcdfe6;
            border-top-right-radius: 6px;
            margin-top: 20px;
            margin-left: -10px;
          }

          .order-inversion-top {
            width: 50px;
            height: 20px;
            border-bottom: 2px solid #dcdfe6;
            border-right: 2px solid #dcdfe6;
            border-bottom-right-radius: 6px;
            margin-left: -10px;
            margin-top: -10px;
          }

          .order-inversion-desc {
            padding: 10px;
            box-sizing: border-box;
            margin-left: -30px;

            p {
              text-align: center;
            }

            p:first-child {
              color: #909399;
              font-size: 12px;
            }

            p:last-child {
              color: #606266;
              font-weight: bolder;
              font-size: 16px;
            }
          }
        }

        .bordernone {
          border: none;
          border-radius: none;
          right: -115px;

          .width-55 {
            width: 55px;
          }

          .padding-10 {
            padding: 0px 0;
          }
        }

        .border-b-r {
          position: absolute;
          width: 159px;
          height: 119px;
          border-bottom: 2px solid #dcdfe6;
          border-right: 2px solid #dcdfe6;
          border-bottom-right-radius: 6px;
          bottom: 30px;
          right: -125px;

          .border-buy {
            position: absolute;
            bottom: 30px;
            left: 15px;
            width: 70px;
            height: 30px;
            border-bottom: 2px solid #dcdfe6;
            border-right: 2px solid #dcdfe6;
            border-bottom-right-radius: 6px;
          }
        }

        .pay-invation {
          position: absolute;
          top: -45px;
          right: -56px;

          p {
            text-align: center;
          }

          p:first-child {
            color: #909399;
            font-size: 12px;
          }

          p:last-child {
            color: #606266;
            font-weight: bolder;
            font-size: 16px;
          }
        }
      }
    }
  }

  .shop-right {
    flex: 4;

    .shop-right-title {
      color: #606266;
      font-size: 14px;
      font-weight: bolder;
      margin-bottom: 20px;
    }

    .shop-right-time {
      color: #606266;
      font-size: 14px;
      margin-bottom: 58px;
    }

    .custom-invation {
      .custon-invation-card {
        position: relative;
        padding: 15px 0 0 30px;
        box-sizing: border-box;
        width: 100%;
        height: 77px;
        background: linear-gradient(180deg, #86a9ff 0%, #708df5 100%);
        border-radius: 4px;

        img {
          width: 100%;
          height: 100%;
        }

        .descinfo {
          p {
            color: #fff;
            line-height: 1.8;
          }

          p:first-child {
            font-size: 16px;
            font-weight: bolder;
          }

          p:last-child {
            font-size: 14px;
          }
        }
      }

      .periodcon {
        display: flex;
        justify-content: center;
        align-items: center;
        margin: 30px 0;

        .arrow {
          position: relative;

          p {
            width: 0;
            height: 0;
            border-left: 20px solid transparent;
            border-right: 20px solid transparent;
            border-top: 25px solid #dcdfe6;
          }
        }

        .descinfo {
          margin-left: 15px;
          line-height: 1.5;

          .font-normal {
            color: #909399;
            font-size: 12px;
          }

          .font-bolder {
            font-weight: bolder;
            font-size: 16px;
            color: #606266;
            margin-left: 5px;
          }
        }
      }
    }
  }
}

.second-title {
  color: #303133;
  font-size: 14px;
  margin-bottom: 15px;
  font-weight: bolder;

  span {
    font-weight: normal;
  }
}

.company-num {
  color: #000000;
  font-size: 14px;
  padding: 8px 27px;
  background: #ebebeb;
}

.private-guide {
  display: flex;
  justify-content: flex-start;
  align-items: center;
  margin-top: 15px;

  .card {
    width: 30%;
    margin-right: 30px;
    position: relative;

    span {
      line-height: 0;
    }

    div {
      line-height: 1.5;
    }

    background: #ffffff;
    border: 1px solid #000000;
    padding: 24px 30px;
    box-sizing: border-box;

    .title {
      display: flex;
      align-items: center;
      font-size: 16px;
      color: #606266;
    }

    .price {
      color: #606266;
      font-weight: bolder;
      font-size: 26px;
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

.select-header {
  position: absolute;
  left: 0;
  top: 50%;
  transform: translate(0%, -50%);
}

.border-btm-none::before {
  display: none;
}

.font-funnel {
  font-size: 16px;
  color: #fff;
}

.hidden {
  visibility: hidden;
}

.font-bolder {
  font-weight: bolder;
}

.posr {
  position: relative;
}

.flex-start-center {
  display: flex;
  justify-content: flex-start;
  align-items: center;
  flex-wrap: nowrap;
}

.m-l-3 {
  margin-left: 3px;
}

.padding-bottom-none {
  padding-bottom: 0px;
}

.help-home-icon {
  width: 14px;
  height: 14px;
  margin-top: -2px;
  margin-left: 4px;
}

.help-icon {
  width: 12px;
  height: 12px;
}

div {
  user-select: none;
}

.main-color {
  color: $color-main;
}

.flex-c-c {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  flex-wrap: wrap;
}

.pointer {
  cursor: pointer;
}

.el-radio-group {
  ::v-deep.el-radio-button {
    .el-radio-button__inner {
      height: 36px !important;
      font-size: 12px;
    }
  }
}

.shop-left {
  ::v-deep.cy-select-small {
    width: 160px !important;
  }
}

.m-b-0 {
  margin-bottom: 0px;
}

.inline-block-width-25 {
  display: inline-block;
  width: 20px !important;
}

.center {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.width-400 {
  width: 400px;
}

.min-width-200 {
  min-width: 300px;
}

.table-empty-firstlabel {
  text-align: center;

  img {
    width: 132px;
  }
}

// 新增的首页底部提示部分
.footer-tips {
  padding: 10px !important;

  .content {
    .title {
      font-weight: 600;
      color: $color-text-main;
      font-size: 12px;
      line-height: 18px;
      margin-bottom: 5px;
    }

    .names {
      font-size: 12px;
      color: $color-text-normal;
      line-height: 18px;
    }
  }
}

.icon-81-81 {
  width: 81px !important;
  height: 81px !important;
}

#crs-store {
  .label-no-padding {
    ::v-deep {
      .el-radio-button__inner {
        padding: 12px 20px !important;
      }
    }
  }
}

.dialog-content {
  padding: 0 20px;

  h3,
  h4 {
    text-align: center;
  }

  .content-item {
    margin-top: 20px;

    .item-desc {
      font-size: 12px;
    }
  }
}

.dialog-footer {
  width: 100%;
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.export-btn {
  margin-left: 50px;
}
</style>
