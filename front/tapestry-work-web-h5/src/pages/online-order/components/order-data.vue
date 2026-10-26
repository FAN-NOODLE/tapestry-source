<template>
  <div class="content">
    <com-filter></com-filter>
    <div class="section">
      <div class="title" @click="handleViewOrder">
        <div class="icon"></div>
        <div class="text">全部订单</div>
        <div class="num">{{Number(orderData.all_num||0).toLocaleString()}}</div>
        <div class="icon-right"></div>
      </div>
    </div>
    <div class="section">
      <div class="title">
        <div class="icon-2"></div>
        <div class="text">订单数据</div>
        <div class="flex" @click="showOrderData = true">
          <div class="more">查看全部</div>
          <div class="icon-right-3"></div>
        </div>
      </div>
      <van-row gutter="10">
        <van-col span="24">
          <div class="box">
            <div class="box-label">有效订单
              <div class="icon-tips" @click="showTips=true"></div>
            </div>
            <div class="box-value">{{(Number(orderData.pending_num||0) + Number(orderData.send_num||0)).toLocaleString()}}</div>
          </div>
        </van-col>
        <van-col span="12">
          <div class="box">
            <div class="box-label">待发货</div>
            <div class="box-value">{{Number(orderData.pending_num||0).toLocaleString()}}</div>
          </div>
        </van-col>
        <van-col span="12">
          <div class="box">
            <div class="box-label">已发货</div>
            <div class="box-value">{{Number(orderData.send_num||0).toLocaleString()}}</div>
          </div>
        </van-col>
      </van-row>
    </div>
    <div class="section">
      <div class="title">
        <div class="icon-2"></div>
        <div class="text">商品数据</div>
        <div class="flex" @click="showGoodsData = true">
          <div class="more">查看全部</div>
          <div class="icon-right-3"></div>
        </div>
      </div>
      <van-row gutter="10">
        <van-col span="12">
          <div class="box">
            <div class="box-label">待发货件数</div>
            <div class="box-value">{{Number(goodsData.pending_num||0).toLocaleString()}}</div>
          </div>
        </van-col>
        <van-col span="12">
          <div class="box">
            <div class="box-label">已发货件数</div>
            <div class="box-value">{{Number(goodsData.send_num||0).toLocaleString()}}</div>
          </div>
        </van-col>
      </van-row>
    </div>
    <!-- 订单数据弹框 -->
    <van-popup v-model:show="showOrderData" position="bottom">
      <div class="popup">
        <div class="popup-title">订单数据</div>
        <div class="popup-close" @click.stop="showOrderData=false"></div>
        <div class="popup-content">
          <van-row gutter="10">
            <van-col span="24">
              <div class="box">
                <div class="box-label">订单总数</div>
                <div class="box-value">{{Number(orderData.all_num||0).toLocaleString()}}</div>
              </div>
            </van-col>
            <van-col span="12">
              <div class="box">
                <div class="box-label">待发货</div>
                <div class="box-value">{{Number(orderData.pending_num||0).toLocaleString()}}</div>
              </div>
            </van-col>
            <van-col span="12">
              <div class="box">
                <div class="box-label">已发货</div>
                <div class="box-value">{{Number(orderData.send_num||0).toLocaleString()}}</div>
              </div>
            </van-col>
            <van-col span="8">
              <div class="box">
                <div class="box-label">待填写</div>
                <div class="box-value">{{Number(orderData.wait_num||0).toLocaleString()}}</div>
              </div>
            </van-col>
            <van-col span="8">
              <div class="box">
                <div class="box-label">已拒单</div>
                <div class="box-value">{{Number(orderData.refuse_num||0).toLocaleString()}}</div>
              </div>
            </van-col>
            <van-col span="8">
              <div class="box">
                <div class="box-label">已作废</div>
                <div class="box-value">{{Number(orderData.cancel_num||0).toLocaleString()}}</div>
              </div>
            </van-col>
          </van-row>
        </div>
      </div>
    </van-popup>
    <!-- 商品数据弹框 -->
    <van-popup v-model:show="showGoodsData" position="bottom">
      <div class="popup">
        <div class="popup-title">商品数据</div>
        <div class="popup-close" @click.stop="showGoodsData=false"></div>
        <div class="popup-content">
          <van-row gutter="10">
            <van-col span="24">
              <div class="box">
                <div class="box-label">商品总件数</div>
                <div class="box-value">{{Number(goodsData.all_num||0).toLocaleString()}}</div>
              </div>
            </van-col>
            <van-col span="12">
              <div class="box">
                <div class="box-label">待发货件数</div>
                <div class="box-value">{{Number(goodsData.pending_num||0).toLocaleString()}}</div>
              </div>
            </van-col>
            <van-col span="12">
              <div class="box">
                <div class="box-label">已发货件数</div>
                <div class="box-value">{{Number(goodsData.send_num||0).toLocaleString()}}</div>
              </div>
            </van-col>
            <van-col span="8">
              <div class="box">
                <div class="box-label">待填写件数</div>
                <div class="box-value">{{Number(goodsData.wait_num||0).toLocaleString()}}</div>
              </div>
            </van-col>
            <van-col span="8">
              <div class="box">
                <div class="box-label">已作废件数</div>
                <div class="box-value">{{Number(goodsData.cancel_num||0).toLocaleString()}}</div>
              </div>
            </van-col>
            <van-col span="8">
              <div class="box">
                <div class="box-label">已拒单件数</div>
                <div class="box-value">{{Number(goodsData.refuse_num||0).toLocaleString()}}</div>
              </div>
            </van-col>
          </van-row>
        </div>
      </div>
    </van-popup>
    <!-- 提示弹框 -->
    <van-popup v-model:show="showTips" position="center">
      <div class="popup">
        <div class="popup-close" @click.stop="showTips=false"></div>
        <div class="popup-content">
          <div class="tips">有效订单 = 待发货订单 + 已发货订单</div>
          <div class="btn" @click="showTips = false">我知道了</div>
        </div>
      </div>
    </van-popup>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import ComFilter from './com-filter.vue';

import Store from "@/store/index";
import { getDashboardApi, getOrderNumApi } from '@/api/online-order';

@Component({
  components: {
    ComFilter
  }
})
export default class OrderData extends Vue {
  orderData: any = {};
  goodsData: any = {};
  showOrderData: boolean = false;
  showGoodsData: boolean = false;
  showTips: boolean = false;
  get filterParams() {
    return this.$store.getters["onlineOrder/getFilterParams"];
  }

  mounted() {
    this.getData();
  }

  activated() {
    this.getData();
  }

  async getData() {
    
    const res:any = await getDashboardApi(this.filterParams)
    if(res.code === 200) {
      this.orderData = res.data.order
      this.goodsData = res.data.goods
    }
  }

  handleViewOrder() {
    this.$router.push({
      path: 'online-order-list',
      query: {
        from: '1'
      }
    })
    // 神策数据上报 -- tab切换
    this.$track.trackTabSwitch({
      tab_name: '全部订单'
    })
  }
}
</script>
<style lang="scss" scoped>
.content {
  padding: 15px;
}

.section {
  background-color: #fff;
  margin-top: 15px;
  padding: 15px;

  .title {
    display: flex;
    align-items: center;
    .icon {
      width: 20px;
      height: 20px;
      background: url('~@/assets/online-order/icon_1.png') center/100% no-repeat
    }

    .icon-2 {
      width: 20px;
      height: 20px;
      background: url('~@/assets/online-order/icon_2.png') center/100% no-repeat
    }

    .text {
      font-weight: 500;
      font-size: 16px;
      color: #000000;
      margin-left: 4px;
      flex: 1;
    }

    .num {
      font-size: 16px;
      color: #000000;
      margin-right: 4px;
    }
    .flex {
      align-items: center;
    }
    .more {
      font-size: 14px;
      color: #666666;
      margin-right: 4px;
    }

    .icon-right {
      width: 9px;
      height: 14px;
      background: url('~@/assets/online-order/icon_right_black.png') center/100% no-repeat
    }
    .icon-right-3 {
      width: 9px;
      height: 14px;
      background: url('~@/assets/online-order/icon_right_3.png') center/100% no-repeat
    }
  }
}

.box {
  height: 86px;
  background: #F6F6F6;
  margin-top: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;

  .box-label {
    font-size: 14px;
    color: #000000;
    line-height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;

    .icon-tips {
      margin-left: 2px;
    }

  }
  .box-value {
    font-size: 28px;
    color: #4A4A4A;
    line-height: 34px;
    font-weight: 500;
  }
}

.popup {
  position: relative;
  .popup-close {
    width: 48px;
    height: 48px;
    position: absolute;
    right: 0;
    top: 0;
    background: url(~@/assets/online-order/icon_popup_close.png) center/18px no-repeat;
  }
  .popup-title {
    height: 52px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 500;
    font-size: 16px;
    color: #000000;
  }
  .popup-content {
    padding: 0 15px 40px;

    .tips {
      font-size: 14px;
      color: #000000;
      line-height: 20px;
      width: 280px;
      padding: 40px 0 20px;
      text-align: center;
    }

    .btn {
      font-size: 14px;
      color: #000000;
      width: 150px;
      height: 34px;
      line-height: 34px;
      text-align: center;
      border: 1px solid #000000;
      margin: auto;
    }
  }
}


.icon-tips {
  position: relative;
  z-index: 11;
  width: 20px;
  height: 20px;
  background: url('~@/assets/online-order/icon_tips_2.png') center/100% no-repeat;
}

</style>