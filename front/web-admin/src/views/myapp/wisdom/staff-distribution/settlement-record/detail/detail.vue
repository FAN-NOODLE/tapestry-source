<template>
  <div class="record-detail" v-if="page_init">
    <layout-list>
      <template slot="table">
        <div class="settlement-info m-b-30">
          <el-row class="title">
            <el-col :span="24">结算单基本信息</el-col>
          </el-row>
          <el-row :gutter="20" class="info">
            <el-col :span="8">
              <span class="staf-title">结算单号</span>
              <span>{{ detail.id }}</span>
            </el-col>
            <el-col :span="8">
              <span class="staf-title">结算状态</span>
              <span>
                <el-tag
                  v-if="detail.settlement_status === 0"
                  type="success"
                  size="small"
                >
                  待结算
                </el-tag>
                <el-tag
                  v-if="detail.settlement_status === 1"
                  type="info"
                  size="small"
                >
                  已结算
                </el-tag>
                <el-tag
                  v-if="detail.settlement_status === 2"
                  type="info"
                  size="small"
                >
                  已取消
                </el-tag>
              </span>
            </el-col>
            <el-col :span="8">
              <span class="staf-title">订单编号</span>
              <span>{{ detail.order_no }}</span>
            </el-col>
          </el-row>
          <el-row :gutter="20" class="info">
            <el-col :span="8">
              <span class="staf-title">订单状态</span>
              <span>
                <el-tag v-if="detail.order_status === 40" size="small">
                  待发货
                </el-tag>
                <el-tag
                  v-if="detail.order_status === 60"
                  type="warning"
                  size="small"
                >
                  待收货
                </el-tag>
                <el-tag
                  v-if="detail.order_status === 80"
                  type="info"
                  size="small"
                >
                  已签收
                </el-tag>
                <el-tag
                  v-if="detail.order_status === 61"
                  type="info"
                  size="small"
                >
                  售后中
                </el-tag>
                <el-tag
                  v-if="detail.order_status === 200"
                  type="info"
                  size="small"
                >
                  退款中
                </el-tag>
                <el-tag
                  v-if="detail.order_status === 205"
                  type="info"
                  size="small"
                >
                  已退款
                </el-tag>
                <el-tag
                  v-if="detail.order_status === 21"
                  type="info"
                  size="small"
                >
                  已取消
                </el-tag>
                <el-tag
                  v-if="detail.order_status === 100"
                  type="info"
                  size="small"
                >
                  已完成
                </el-tag>
                <el-tag
                  v-if="detail.order_status === 201"
                  type="info"
                  size="small"
                >
                  退款失败
                </el-tag>
                <el-tag
                  v-if="detail.order_status === 83"
                  type="info"
                  size="small"
                >
                  部分退款
                </el-tag>
              </span>
            </el-col>
            <el-col :span="8">
              <span class="staf-title">支付时间</span>
              <span>{{
                detail.order_pay_time | formatDate("yyyy-mm-dd hh:mm:ss")
              }}</span>
            </el-col>
            <el-col :span="8">
              <span class="staf-title">推广导购</span>
              <span class="imgBox">
                <div class="nameBox">
                  <img
                    v-if="detail.promotion_avatar"
                    :src="detail.promotion_avatar"
                    alt=""
                  />
                  <span class="m-l-5">{{
                    detail.promotion_name ? detail.promotion_name : "--"
                  }}</span>
                </div>
              </span>
            </el-col>
          </el-row>
          <el-row :gutter="20" class="info" v-if="USER.version !== 3">
            <el-col :span="8">
              <span class="staf-title">推广佣金</span>
              <span>{{ detail.promotion_commission / 100 }}元</span>
            </el-col>
            <el-col v-if="USER.version !== 3" :span="8">
              <span class="staf-title">服务导购</span>
              <span class="imgBox">
                <div class="nameBox">
                  <img
                    v-if="detail.service_avatar"
                    :src="detail.service_avatar"
                    alt=""
                  />
                  <span class="m-l-5">{{
                    detail.service_name ? detail.service_name : "--"
                  }}</span>
                </div>
              </span>
            </el-col>
            <el-col v-if="USER.version !== 3" :span="5">
              <span class="staf-title">服务佣金</span>
              <span>{{ detail.service_commission / 100 }}元</span>
            </el-col>
          </el-row>
        </div>

        <div class="distribution-goods">
          <div class="title">分销商品</div>
          <el-table
            :data="goods_list"
            stripe
            class="m-t-15 list-last-line-none min-height-100"
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
              min-width="150"
            ></el-table-column>
            <el-table-column label="商品" min-width="200">
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
              min-width="110"
              prop="purchase_quantity"
            ></el-table-column>
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
            <el-table-column label="推广佣金(元)" min-width="120">
              <template slot-scope="scope">
                <span>{{ scope.row.promotion_commission / 100 }}</span>
              </template>
            </el-table-column>
            <el-table-column
              v-if="USER.version !== 3"
              label="服务佣金(元)"
              min-width="120"
            >
              <template slot-scope="scope">
                <span>{{ scope.row.service_commission / 100 }}</span>
              </template>
            </el-table-column>
            <div slot="empty">
              <layout-empty :message="'暂无分销商品'"></layout-empty>
            </div>
          </el-table>
        </div>
      </template>
    </layout-list>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import UserModule from "@/store/modules/user";
import SettlementApi from "@/api/wisdom/Settlement";
import LayoutEmpty from "@/components/sy-components/layout-empty.vue";

@Component({
  components: {
    LayoutEmpty
  }
})
export default class extends Vue {
  /*** data ***/
  private page_init = false;
  private detail: Array<any> = []; // 结算记录明细
  private goods_list: Array<any> = []; // 结算商品信息

  /*** 计算属性 ***/
  get USER() {
    return UserModule.getUserInfo;
  }

  getDetail() {
    const id = this.$route.query.id;
    SettlementApi.getSettlementDetailApi({ id }).then(res => {
      this.detail = res.data.detail;
      this.goods_list = res.data.settlement_goods;
      this.page_init = true;
    });
  }
  private mounted() {
    this.getDetail();
  }
}
</script>
<style lang="scss" scoped>
.record-detail {
  color: #303133;
}
.settlement-info {
  padding-bottom: 20px;
  border-bottom: 1px solid $color-third-line;
}
.title {
  line-height: 32px;
  font-size: 22px;
  font-weight: 700;
}
.info {
  margin-top: 25px;
  font-size: 14px;
  span {
    margin-right: 10px;
  }
}
.staf-title {
  display: inline-block;
  width: 60px;
  color: #606266;
  text-align: right;
}
.imgBox {
  display: inline-block;
  vertical-align: middle;
}
.nameBox {
  display: flex;
  align-items: center;
  img {
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
</style>
