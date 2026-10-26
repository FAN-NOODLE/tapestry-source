<template>
  <div class="subscribe-detail">
    <layout-form>
      <template slot="content">
        <div class="row">
          <img
            :src="
              detailInfo.avatar ||
                'https://oss-crs.vchangyi.com/wechat/default-avatar-icon.png'
            "
            alt=""
          />
          <div class="user-info">
            <div class="customer-name">
              客户姓名：{{ detailInfo.user_name }}
            </div>
            <div class="user-info-box">
              <div class="info-item">
                邀约状态：{{ statusMap[detailInfo.status] || "-" }}
              </div>
              <div class="info-item">
                邀请来源：{{
                  orderOriginList[detailInfo.order_origin]
                    ? orderOriginList[detailInfo.order_origin].label
                    : "--"
                }}
              </div>
              <div class="info-item">
                是否成单：{{
                  ifOrNotLists[detailInfo.is_pay]
                    ? ifOrNotLists[detailInfo.is_pay].label
                    : "--"
                }}
              </div>
              <div class="info-item">
                是否评价：{{
                  ifOrNotLists[detailInfo.is_comment]
                    ? ifOrNotLists[detailInfo.is_comment].label
                    : "--"
                }}
              </div>
              <div class="info-item">
                顾客类型：{{
                  customerStatusList[detailInfo.customer_status]
                    ? customerStatusList[detailInfo.customer_status].label
                    : "--"
                }}
              </div>
              <div class="info-item" v-if="epNameDetail !== 'coach'">
                签到状态：{{ detailInfo.is_sign_in_str || "--" }}
              </div>
            </div>
          </div>
        </div>
        <el-form ref="form" label-width="100px">
          <div class="form-title">邀约单详情</div>
          <el-form-item label="创建时间：">{{
            detailInfo.create_time_str || "--"
          }}</el-form-item>
          <el-form-item :label="`邀约服务：`">{{
            detailInfo.reserve_name || "--"
          }}</el-form-item>
          <el-form-item :label="`邀约时间：`">{{
            detailInfo.reserve_time_str || "--"
          }}</el-form-item>
          <el-form-item :label="`邀约门店：`">{{
            detailInfo.store_name || "--"
          }}</el-form-item>
          <el-form-item :label="`邀约导购：`">{{
            detailInfo.member_name || "--"
          }}</el-form-item>
          <el-form-item label="导购留言：">{{
            detailInfo.reserve_remark || "--"
          }}</el-form-item>
          <el-form-item label="客户留言：">{{
            detailInfo.user_remark || "--"
          }}</el-form-item>
        </el-form>

        <el-row ref="form" label-width="120px">
          <div class="form-title">关联订单详情</div>
          <el-row>
            <el-col :span="6">
              <span>顾客类型：{{
                ["--", "新客", "老客"][detailInfo.customer_status]
              }}</span>
            </el-col>
            <el-col :span="6">
              <span>是否到店购买：{{
                ["--", "是", "否"][detailInfo.is_store]
              }}</span>
            </el-col>
            <el-col :span="6">
              <span>关联订单数：{{ detailInfo.relation_order_num||'0' }}单</span>
            </el-col>
            <el-col :span="6">
              <span>订单总金额：¥{{
                Number(detailInfo.reserve_order_price / 100).toLocaleString()
              }}</span>
            </el-col>
          </el-row>
        </el-row>

        <el-form v-if="showBlock(1)" ref="form" label-width="120px">
          <div class="form-title">邀约评价详情</div>
          <el-form-item label="评价时间：">{{
            detailInfo.comment_time | formatDate("yyyy-mm-dd hh:mm:ss")
          }}</el-form-item>
          <el-form-item label="活动安排：">
            <div class="center">
              <StarComment
                :list="setStarList(detailInfo.service_score)"
              ></StarComment>
            </div>
          </el-form-item>
          <el-form-item label="整体服务：">
            <div class="center">
              <StarComment
                :list="setStarList(detailInfo.store_score)"
              ></StarComment>
            </div>
          </el-form-item>
          <el-form-item label="顾客评价：">{{
            detailInfo.remark || "--"
          }}</el-form-item>
        </el-form>
        <div v-if="showBlock(2)">
          <div class="form-title">操作日志</div>
          <el-timeline :reverse="false" v-if="detailInfo.log_list_arr">
            <el-timeline-item
              v-for="(activity, index) in detailInfo.log_list_arr"
              :key="index"
              :timestamp="activity.create_at_str"
            >
              <div
                class="flex-start"
                v-for="(item, index) in activity"
                :key="index"
              >
                <div class="min-w-56">{{ item.key }}</div>
                <div
                  class="m-l-50"
                  v-if="item.key === '活动安排' || item.key === '整体服务'"
                >
                  <StarComment :list="setStarList(item.value)"></StarComment>
                </div>
                <div class="m-l-50" v-else>{{ item.value || "--" }}</div>
              </div>
            </el-timeline-item>
          </el-timeline>
        </div>
        <!-- 订单记录 -->
        <div v-if="showBlock(3)" class="m-t-20">
          <div class="form-title">订单记录</div>
          <div class="order-list">
            <div v-for="(item, index) in detailInfo.reserve_order_sales" :key="index">
              <div class="item" v-for="(goods, goodsIndex) in item.reserve_order_goods" :key="goodsIndex">
                <el-row class="item-order-num">
                  <el-col :span="12">订单编号：{{ goods.relation_order_id }}</el-col>
                  <el-col :span="12" class="text-right">金额：¥{{ Number(goods.goods_price/100).toLocaleString() }}</el-col>
                </el-row>
                <el-row style="display: flex;align-items: center;">
                  <el-col :span="4">
                    <el-image class="img-goods" :src="goods.goods_url" fit="contain">
                      <template slot="error">
                        <div class="image-slot"></div>
                      </template>  
                    </el-image>
                  </el-col>
                  <el-col :span="5">
                    <span>商品名称：{{ goods.goods_name }}</span>
                  </el-col>
                  <el-col :span="5">
                    <span>款号：{{ goods.style }}</span>
                  </el-col>
                  <el-col :span="5">
                    <span>购买日期：{{ item.pay_time|formatDate("yyyy-mm-dd hh:mm:ss") }}</span>
                  </el-col>
                  <el-col :span="5">
                    <span>店铺：{{ item.store_code }}</span>
                  </el-col>
                </el-row>
              </div>
            </div>
          </div>
        </div>
      </template>
    </layout-form>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Api from "@/api/store-invite/index";
import StarComment from "../../component/star-comment.vue";
import {
  customerStatusList,
  orderOriginList,
  ifOrNotLists,
  statusMap,
  swStatusMap
} from "../../contants";
@Component({
  components: {
    StarComment
  }
})
export default class extends Vue {
  // 生命周期
  private statusMap = {};
  private detailInfo = {};
  private customerStatusList = customerStatusList;
  private orderOriginList = orderOriginList;
  private ifOrNotLists = ifOrNotLists;
  private id: any = 0;
  private epNameDetail = localStorage.getItem("epNameDetail");

  get showBlock() {
    return (block) => {
      switch (block) {
        case 1:
        case 2:
          return this.epNameDetail !== 'coach'
        case 3:
          return this.epNameDetail === 'coach'
      
        default:
          break;
      }
    }
  }

  get imageUrl() {
    return (item) => {
      let url = 'https://s7d2.scene7.com/is/image/Coach/'
      url += item.style + '_' + item.color + '_a0'
      return url
    }
  }

  private created() {
    console.log("created");
    // todo:sw单独邀约状态map
    this.statusMap = this.epNameDetail === "sw" ? swStatusMap : statusMap;
  }
  private mounted() {
    console.log("mounted");
    this.id = this.$route.query.id || 0;
    this.detail();
  }

  // 设置星级显示
  setStarList(score) {
    score = score ? Number(score) : 0;
    // 评级星级
    const checkedStar = new Array(score).fill({
      checked: true
    });
    // 灰色星级
    const uncheckedStar = new Array(5 - score).fill({
      checked: false
    });

    const starList = checkedStar.concat(uncheckedStar);
    return starList;
  }
  private detail() {
    const epName = localStorage.getItem("epNameDetail");
    Api.orderDetailApi({ id: Number(this.id), enterprise_name: epName }).then(
      res => {
        this.detailInfo = res.data;
        // localStorage.removeItem("epNameDetail");
      }
    );
  }
  private destroyed() {
    console.log("destroyed");
  }
}
</script>
<style lang="scss" scoped>
/deep/.el-form-item {
  margin-bottom: 0;
}
.center {
  height: 40px;
  display: flex;
  align-items: center;
}
.row {
  margin-bottom: 30px;
  align-items: center;
  display: flex;
  color: $color-text-main;
  .title {
    margin-right: 20px;
    color: $color-text-normal;
  }
  .user-info {
    margin-left: 30px;
  }
  .customer-name {
    font-size: 16px;
    font-weight: bold;
  }

  .user-info-box {
    display: flex;
    margin-top: 20px;
    .info-item {
      margin-right: 30px;
    }
  }
  .tag {
    margin-right: 20px;
  }
}
img {
  width: 100px;
  height: 100px;
  border-radius: 50%;
}
.form-title {
  font-size: 16px;
  font-weight: bold;
  margin: 20px 0;
}
.flex-start {
  display: flex;
  justify-content: flex-start;
  padding: 10px 0;
  align-items: center;
}
.m-l-50 {
  margin-left: 50px;
}
.min-w-56 {
  min-width: 56px;
  max-width: 56px;
}

.img-goods {
  width: 100px;
  height: 100px;

}

.image-slot {
  width: 100px;
  height: 100px;
  background: url('~@/assets/images/icon_empty.png') center/100% no-repeat;
}

.text-right {
  text-align: right;
}

.item {
  padding: 20px;
  background-color: #f7f7f7;
  margin-bottom: 15px;
}

.item-order-num {
  font-size: 14px;
  font-weight: bold;
  margin-bottom: 10px;
}
</style>
