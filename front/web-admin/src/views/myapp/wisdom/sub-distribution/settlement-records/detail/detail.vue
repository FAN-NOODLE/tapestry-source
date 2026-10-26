<template>
  <layout-list v-if="isLoaded" class="set-detail">
    <template class="wrap" slot="tabs">
      <div class="settlementInfo mar-b-30">
        <el-row class="title base">
          <el-col :span="24">结算单基本信息</el-col>
        </el-row>
        <el-row>
          <el-col class="" :span="8">
            <span class="info-tag">结算单号</span>
            <span v-text="infoList.id"></span>
          </el-col>
          <el-col class="" :span="8">
            <span class="info-tag">结算状态</span>
            <el-tag
              v-if="infoList.settlement_status === 0"
              type="warning"
              size="small"
            >
              待结算
            </el-tag>
            <el-tag
              v-if="infoList.settlement_status === 1"
              type="success"
              size="small"
            >
              已结算
            </el-tag>
            <el-tag
              v-if="infoList.settlement_status === 2"
              type="info"
              size="small"
            >
              已取消
            </el-tag>
          </el-col>
          <el-col class="" :span="8">
            <span class="info-tag">订单编号</span>
            <span v-text="infoList.order_no"></span>
          </el-col>
        </el-row>
        <el-row>
          <el-col class="" :span="8">
            <span class="info-tag">订单状态</span>
            <el-tag
              v-if="infoList.order_status === 20"
              type="info"
              size="small"
            >
              未支付
            </el-tag>
            <el-tag
              v-if="infoList.order_status === 21"
              type="info"
              size="small"
            >
              已取消
            </el-tag>
            <el-tag
              v-if="infoList.order_status === 23"
              type="info"
              size="small"
            >
              待支付
            </el-tag>
            <el-tag
              v-if="infoList.order_status === 25"
              type="info"
              size="small"
            >
              支付中
            </el-tag>
            <el-tag
              v-if="infoList.order_status === 30"
              type="info"
              size="small"
            >
              待确认
            </el-tag>
            <el-tag
              v-if="infoList.order_status === 40"
              type="info"
              size="small"
            >
              待发货
            </el-tag>
            <el-tag
              v-if="
                infoList.order_status === 41 ||
                  infoList.order_status === 61 ||
                  infoList.order_status === 81
              "
              type="info"
              size="small"
            >
              售后中
            </el-tag>
            <el-tag
              v-if="infoList.order_status === 60"
              type="text"
              size="small"
            >
              待收货
            </el-tag>
            <el-tag
              v-if="infoList.order_status === 80"
              type="info"
              size="small"
            >
              已签收
            </el-tag>
            <el-tag
              v-if="infoList.order_status === 83"
              type="info"
              size="small"
            >
              部分退款
            </el-tag>
            <el-tag
              v-if="infoList.order_status === 100"
              type="info"
              size="small"
            >
              已完成
            </el-tag>
            <el-tag
              v-if="infoList.order_status === 200"
              type="info"
              size="small"
            >
              退款中
            </el-tag>
            <el-tag
              v-if="infoList.order_status === 201"
              type="info"
              size="small"
            >
              退款失败
            </el-tag>
            <el-tag
              v-if="infoList.order_status === 205"
              type="info"
              size="small"
            >
              已退款
            </el-tag>
          </el-col>
          <el-col class="" :span="8">
            <span class="info-tag">支付时间</span>
            <span>{{
              infoList.order_pay_time | formatDate("yyyy-mm-dd hh:mm:ss")
            }}</span>
          </el-col>
          <el-col class="mess" :span="8">
            <span class="info-tag">分销合伙人</span>
            <img
              :src="
                infoList.partner_avatar
                  ? infoList.partner_avatar
                  : 'https://oss-crs.vchangyi.com/wechat/default-avatar-icon.png'
              "
              alt=""
              width="30"
              height="30"
              class="table-text-img"
            />
            <span v-text="infoList.partner_name"></span>
          </el-col>
        </el-row>
        <el-row>
          <el-col class="" :span="8">
            <span class="info-tag">分销合伙人佣金</span>
            <span v-text="infoList.partner_commission / 100"></span>
          </el-col>
          <el-col :span="8">
            <span class="info-tag">所属导购</span>
            <img
              :src="
                infoList.worker_avatar
                  ? infoList.worker_avatar
                  : 'https://oss-crs.vchangyi.com/wechat/default-avatar-icon.png'
              "
              alt=""
              width="30"
              height="30"
              class="table-text-img"
            />
            <span v-text="infoList.worker_name"></span>
          </el-col>
          <el-col class="" :span="8">
            <span class="info-tag">导购佣金</span>
            <span v-text="infoList.promotion_commission / 100"></span>
          </el-col>
        </el-row>
        <div class="line"></div>
      </div>
      <div class="distributionGoods">
        <div class="title mar-b-10">分销商品</div>
        <el-table
          :data="detailList"
          stripe
          class="mar-t-15 list-last-line-none min-height-100"
          :header-cell-style="{
            fontSize: '14px',
            color: '#606266',
            fontWeight: '400',
            padding: '11px 0 10px 0'
          }"
        >
          <el-table-column
            prop="sku_id"
            label="SKU ID"
            min-width="170"
            fixed="left"
          ></el-table-column>
          <el-table-column label="商品" min-width="300">
            <template slot-scope="scope">
              <div class="settlementBox">
                <img :src="scope.row.spu_img" alt="" />
                <div class="attribute m-l-5">
                  <div class="skuname">{{ scope.row.spu_name }}</div>
                  <div
                    v-if="
                      scope.row.attribute && scope.row.attribute.length === 1
                    "
                  >
                    <span v-for="item in scope.row.attribute" :key="item">
                      {{ item }}
                    </span>
                  </div>
                  <div v-else>
                    <span v-for="item in scope.row.attribute" :key="item">
                      {{ item }};
                    </span>
                  </div>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column
            label="购买数量"
            min-width="160"
            prop="purchase_quantity"
          >
          </el-table-column>
          <el-table-column label="商品价格(元)" min-width="120">
            <template slot-scope="scope">
              <span>{{ scope.row.order_goods_price / 100 }}</span>
            </template>
          </el-table-column>
          <el-table-column label="实收金额(元)" min-width="120">
            <template slot-scope="scope">
              <span>{{ scope.row.goods_pay_price / 100 }}</span>
            </template>
          </el-table-column>
          <el-table-column
            prop="promotion_commission"
            label="导购佣金（元）"
            min-width="150"
          >
            <template slot-scope="{ row }">{{
              row.promotion_commission / 100
            }}</template>
          </el-table-column>
          <el-table-column
            prop="partner_commission"
            label="分销合伙人佣金（元）"
            min-width="180"
          >
            <template slot-scope="{ row }">{{
              row.partner_commission / 100
            }}</template>
          </el-table-column>
          <div slot="empty">
            <layout-empty :message="'暂无余额记录'"></layout-empty>
          </div>
        </el-table>
      </div>
    </template>
    <template slot="footer">
      <!-- <div class="line down-line"></div> -->
      <el-button size="mini" @click="backList()" slot="back">返回 </el-button>
    </template>
  </layout-list>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import settlementRecordsApi from "@/api/wisdom/settlementRecords";
import LayoutForm from "@/components/public/layout-form/index.vue";
import LayoutEmpty from "@/components/sy-components/layout-empty.vue";
import LayoutList from "@/components/public/layout-list/index.vue";
import UserModule from "@/store/modules/user";
@Component({
  components: {
    LayoutForm,
    LayoutEmpty,
    LayoutList
  }
})
export default class extends Vue {
  // data数据
  private detailList = []; //余额明细
  private infoList: object = {}; //导购基本信息
  private tooltip = false;
  private isLoaded = false; // 数据是否加载完成

  // 计算属性
  get USER() {
    return UserModule.getUserInfo;
  }

  // methods
  // 返回列表页
  private backList() {
    this.$router.go(-1);
  }

  // 获取明细数据
  private getdetailList() {
    const params = {
      id: this.$route.query.id
    };
    settlementRecordsApi.getSettlementDetailApi(params).then(res => {
      this.detailList = res.data.settlement_goods;
      this.infoList = res.data.detail;
      this.isLoaded = true;
    });
  }

  // 生命周期
  private mounted() {
    this.getdetailList();
  }
}
</script>
<style lang="scss" scoped>
.set-detail {
  color: #303133;
}
.mar-t-5 {
  margin-top: -3px;
}

.skuname {
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
  font-size: 13px;
  line-height: 18px;
  height: 36px;
  width: 200px;
}

.clearfix:after {
  content: "";
  display: block;
  height: 0;
  clear: both;
  visibility: hidden;
}

.clearfix {
  zoom: 1 span {
    float: left;
    padding: 0 5px;
    line-height: 20px;
    color: #409eff;
    border: 1px solid #409eff;
    border-radius: 3px;
    background-color: rgba(64, 158, 255, 0.1);
    margin: 5px;
  }
}

.detail /deep/ .footer-content .el-pagination {
  float: none;
}

.layout-list {
  padding: 20px;
}

.table-text-img {
  border-radius: 100%;
  vertical-align: middle;
  margin-right: 4px;
}

.goBack {
  font-size: 14px;
  font-weight: normal;
  cursor: pointer;
  margin-right: 20px;
  padding-right: 20px;
  border-right: 1px solid #bfbfbf;
  transform: translateY(-5px);
}

.goBack,
.base {
  display: inline-block;
  vertical-align: middle;
}

.base {
  margin-top: -7px;
}

.info-tag {
  margin-right: 10px;
  color: #909399;
}

/deep/.content-block {
  padding: 20px;
}

/deep/.search-content {
  display: none;
}

/deep/.buttons-content {
  display: none;
}

.detail /deep/ .operate-content {
  display: none;
}

.detail /deep/ .slide-box {
  display: none;
}

.title {
  font-weight: 700;
  font-size: 18px;
  line-height: 26px;
  margin-bottom: 20px;
}

.info {
  margin-top: 20px;
  font-size: 14px;

  span {
    margin-right: 10px;
    display: inline-block;
    vertical-align: text-bottom;
  }
}
.el-col {
  margin-bottom: 20px;
}

.info .mess {
  display: flex;
  align-items: center;
  justify-content: flex-start;
}

.line {
  width: 100%;
  height: 1px;
  margin: 0 0 20px;
  background-color: #dcdfe6;
}

.down-line {
  margin: 20px 0;
}

.imgBox {
  display: inline-block;
  vertical-align: middle;
}

.nameBox {
  display: flex;

  align-items: center img {
    width: 30px;
    height: 30px;
    border-radius: 50%;
  }

  span {
    width: 85px;
    flex-shrink: 0;
  }
}

.settlementBox {
  display: flex;
  align-items: center;
  img {
    width: 45px;
    height: 45px;
  }
  span {
    color: $color-text-second;
  }
}

.page-right {
  text-align: right;
}
</style>
