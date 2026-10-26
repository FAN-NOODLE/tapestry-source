<template>
  <div>
    <!-- 请使用微信打开 -->
    <div v-if="status==3" class="container">
      <div class="result">
        <van-image width="50" height="50" :src="require('@/assets/online-order/icon_warning.png')"></van-image>
        <div class="title">温馨提示</div>
        <div class="tip">当前环境内无法访问当前页面，请在微信环境内访问</div>
      </div>
    </div>
    <!-- 已作废 -->
    <div v-if="status==2" class="container">
      <div class="result">
        <van-image width="50" height="50" :src="require('@/assets/online-order/icon_warning.png')"></van-image>
        <div class="title">订单已作废</div>
        <div class="code">销售单号：{{ detail.order_sn }}</div>
        <div class="tip">当前订单已超过地址填写时间，如有疑问请联系导购</div>
      </div>
    </div>
    <!-- 已提交 -->
    <div v-if="status==1" class="container">
      <div class="result">
        <van-image width="50" height="50" :src="require('@/assets/online-order/icon_success.png')"></van-image>
        <div class="title">订单已提交</div>
        <div class="code">销售单号：{{ detail.order_sn }}</div>
        <div class="tip">当前订单配送地址已提交，我们会尽快处理订单</div>
      </div>
    </div>
    <!-- 未作废/未提交 -->
    <div v-if="status==0" class="container">
      <div class="content">
        <div class="section">
          <div class="section-title required">填写收件信息</div>
          <div class="section-content">
            <add-address v-model="address"></add-address>
          </div>
        </div>
        <div class="section">
          <div class="section-title">商品信息</div>
          <div class="section-content">
            <div class="list">
              <div class="item" v-for="item, index in detail.sku_list" :key="index">
                <van-image 
                  width="80px"
                  height="80px"
                  fit="contain"
                  lazy-load
                  :src="item.image_url||emptyImage"
                ></van-image>
                <div class="item-info">
                  <div class="item-code">
                    <span>{{item.sku_info}}</span>
                    <span>x{{item.quantity}}</span>
                  </div>
                  <div class="item-price">{{ item.name }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="section">
          <div class="row">
            <div class="label">销售单号</div>
            <div class="value">{{detail.order_sn}}</div>
          </div>
          <div class="row">
            <div class="label">创建时间</div>
            <div class="value">{{$helper.dateFormat(detail.created_at,  "yyyy-MM-dd hh:mm:ss")}}</div>
          </div>
        </div>
      </div>
      <div class="footer">
        <van-button block color="#000000" @click="handleSubmit">提交</van-button>
      </div>
      <!-- 用户踢脚确认弹框 -->
      <van-popup v-model="visible" @close="hidePopup">
        <div class="popup">
          <div class="popup-close" @click="hidePopup"></div>
          <div class="popup-content">
            <div class="title">请核对您填写的收货信息一经提交无法修改！</div>
          </div>
          <div class="popup-footer">
            <van-button color="#000000" plain @click="hidePopup">返回修改</van-button>
            <van-button color="#000000" @click="handleConfirm">继续提交</van-button>
          </div>
        </div>
      </van-popup>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop } from 'vue-property-decorator'
import AddAddress from './components/add-address.vue';
import { getOrderDetailApi, updateOrderApi } from '@/api/online-order';
import wx from 'weixin-js-sdk'
import wxUtils from "@/utils/wx";
@Component({
  components: {
    AddAddress
  }
})
export default class OnlineOrderInfo extends Vue {
  emptyImage = require('@/assets/material-library/v4.2.0/icon_empty.png')
  visible = false
  status = -1

  address = {
    name: '',
    phone: '',
    county: '',
    countyCode: '',
    address: '',
    zipCode: '',
    province: '',
    provinceCode: '',
    city: '',
    cityCode: '',
    street: '',
    streetCode: '',
    idCard: ''
  }

  detail:any = {}



  async mounted() {
    // try {
    //   await wxUtils.wxConfig();
    //   console.log(3333);
      
    // } catch (error) {
    //   console.log(111);
      
    //   console.log(error);
    // }

    // 当前环境
    var isComWx = /wxwork/i.test(navigator.userAgent); // 是否企业微信
    if(isComWx) {
      this.status = 3
    } else {
      this.getData()
    }
    // 用户使用 Android 设备
    window.history.pushState("_android_back", "", location.href);
    window.addEventListener("popstate", this.popstate, false);
    
    
    
  }
  destroyed() {
    // 记得清除!!!!不然单页面，其它页面也可以用
    window.removeEventListener("popstate", this.popstate, false);
  }
  popstate() {
    wx.closeWindow(); // 微信网页退出
  }
  

  async getData() {
    const res = await getOrderDetailApi({
      id: Number(this.$route.query.order_id)
    })
    
    this.status = res.data.order_status
    this.detail = {
      ...res.data.order_info,
      order_sn: res.data.order_sn,
    }
  }

  handleSubmit() {
    this.visible = true
  }

  hidePopup() {
    this.visible = false
  }

  async handleConfirm() {
    this.visible = false
    const {phone,name,province, city, county, address} = this.address
    if(!phone || !name || !province || !city || !county || !address) {
      this.$toast('请输入完整的收货信息')
      return
    }

    try {
      this.$toast.loading({mask:true})
      const res:any = await updateOrderApi({
        order_id: Number(this.$route.query.order_id),
        receiver_phone: phone,
        receiver_name: name,
        receiver_address: address,
        province,
        city,
        county
      })
      this.$toast.clear()
      if(res.code===200) {
        this.status = 1
      } else {
        this.$toast(res.message)
        this.getData()
      }
      
    } catch (error) {
      this.$toast.clear()
      this.$toast('网络异常，请稍后重试')
    }
  }
}
 
</script>

<style lang="scss" scoped>
.popup {
  width: 325px;
  background: #FFFFFF;

  position: relative;
  .popup-close {
    width: 48px;
    height: 48px;
    position: absolute;
    right: 0;
    top: 0;
    background: url(~@/assets/online-order/icon_popup_close.png) center/18px no-repeat;
  }
  
  .popup-content {
    text-align: center;
    padding: 40px 0 0;
    .title {
      margin-top: 15px;
      font-weight: 500;
      font-size: 16px;
      color: #000000;
      line-height: 22px;
      margin: auto;
      width: 176px;
    }
  }

  .popup-footer {
    padding: 20px 15px;
    display: flex;
    align-items: center;

    .van-button {
      flex: 1;
      height: 36px;
      &:nth-child(2) {
        margin-left: 10px;
      }
    }
  }
}
.container {
  width: 100vw;
  min-height: 100vh;
  background: #F7F7F7;
  padding-bottom: calc(64px + contain(safe-area-inset-bottom));
  padding-bottom: calc(64px + env(safe-area-inset-bottom));

  .content {
    padding: 15px;
  }


  .footer {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    padding: 10px 15px;
    z-index: 2;
    display: flex;
    align-items: center;
    background-color: #fff;
    padding-bottom: calc(10px + contain(safe-area-inset-bottom));
    padding-bottom: calc(10px + env(safe-area-inset-bottom));

    .btn {
      flex: 1;
      &:not(:first-child) {
        margin-left: 15px
      }
    }
  }

  .section {
    padding: 15px;
    background: #FFFFFF;
    margin-bottom: 15px;
    .section-title {
      font-weight: 500;
      font-size: 16px;
      color: #000000;
      line-height: 22px;

      &.required::after {
        content: "*";
        color: #F64040;
        margin-left: 2px;
      }
    }
    .section-content {
      margin-top: 10px;
    }
    .row {
      display: flex;
      align-items: center;

      &:not(:last-child) {
        margin-bottom: 15px;
      }

      .label {
        font-weight: 400;
        font-size: 14px;
        color: #96979E;
        line-height: 20px;
        flex: 1;
      }
      .value {
        font-weight: 400;
        font-size: 14px;
        color: #000000;
        line-height: 20px;
      }
    }
  }

  .item {
    display: flex;
    position: relative;

    &:not(:last-child) {
      margin-bottom: 15px;
    }

    .van-image {
      background: #EFEFEF;
    }

    .item-info {
      margin-left: 10px;
      flex: 1;
      overflow: hidden;

      .item-code {
        font-weight: 400;
        font-size: 14px;
        color: #000000;
        line-height: 20px;
        display: flex;
        align-items: center;
        justify-content: space-between;
      }

      .item-price {
        font-weight: 400;
        font-size: 12px;
        line-height: 17px;
        margin-top: 20px;
      }

      .item-stock {
        font-weight: 400;
        font-size: 12px;
        color: #000000;
        line-height: 17px;
        margin-top: 14px;

        &.error {
          color: #F64040;
        }
      }
    }
  }

  .result {
    padding-top: 78px;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;


    .title {
      font-weight: 500;
      font-size: 16px;
      color: #000000;
      line-height: 22px;
      margin-top: 15px;
    }

    .code {
      font-weight: 400;
      font-size: 14px;
      color: #000000;
      line-height: 20px;
      margin-top: 5px;
    }

    .tip {
      font-weight: 400;
      font-size: 14px;
      color: #A6A6A6;
      line-height: 20px;
      margin-top: 10px;
    }
  }
}
</style>