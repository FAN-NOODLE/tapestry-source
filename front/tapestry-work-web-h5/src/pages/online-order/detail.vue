<template>
  <div v-if="detail" class="container">
    <div class="content">
      <div v-if="!detail.receiver_phone" class="section flex-center">
        <van-image width="20" height="20" :src="require('@/assets/online-order/icon_address.png')"></van-image>
        <div class="">
          <div class="address">待客户填写收件信息</div>
        </div>
      </div>
      <div v-else class="section flex">
        <van-image width="20" height="20" :src="require('@/assets/online-order/icon_address.png')"></van-image>
        <div class="flex-1">
          <div class="address">{{ detail.receiver_address }}</div>
          <div class="flex info">
            <span class="name">{{detail.receiver_name}}</span>
            <span class="phone">{{detail.receiver_phone}}</span>
          </div>
        </div>
      </div>
      <template v-for="item,index in statusList">
        <div v-if="item.list.length>0" class="section" :key="index">
          <div class="status" :class="`status-${item.value}`">{{item.label}}</div>
          <div class="list">
            <div class="item" v-for="sku, skuIndex in item.list" :key="skuIndex">
              <van-image 
                width="80px"
                height="80px"
                fit="contain"
                lazy-load
                :src="sku.image_url||emptyImage"
              ></van-image>
              <div class="item-info">
                <div class="item-code">{{sku.sku_info}}</div>
                <div class="item-name">{{sku.name}}</div>
                <div class="item-stock">{{sku.quantity}}件</div>
              </div>
            </div>
          </div>
          <div v-if="item.value === 3" class="box">
            <div class="row">
              <div class="label">顺丰单号</div>
              <div class="value flex">
                <span>{{detail.tracking_number}}</span>
                <van-button size="mini" color="#000000" plain @click="$helper.copyMessage(detail.tracking_number)">复制</van-button>
              </div>
            </div>
            <div class="row">
              <div class="label">发货时间</div>
              <div class="value">{{$helper.dateFormat(detail.shipped_at, 'yyyy-MM-dd hh:mm:ss')}}</div>
            </div>
          </div>

          <div v-if="item.value === 4" class="box">
            <div class="row">
              <div class="label">拒单时间</div>
              <div class="value">{{$helper.dateFormat(detail.operate_time, 'yyyy-MM-dd hh:mm:ss')}}</div>
            </div>
          </div>
        </div>
      </template>

      <div class="section">
        <div class="row">
          <div class="label">销售单号</div>
          <div class="value">{{detail.order_sn}}</div>
        </div>
        <div class="row">
          <div class="label">门店代码</div>
          <div class="value">{{detail.store_code}}</div>
        </div>
        <div class="row">
          <div class="label">门店名称</div>
          <div class="value">{{detail.store_name}}</div>
        </div>
        <div class="row">
          <div class="label">创建人</div>
          <div class="value">{{detail.member_name}}</div>
        </div>
        <div class="row">
          <div class="label">创建时间</div>
          <div class="value">{{ $helper.dateFormat(detail.created_at, "yyyy-MM-dd hh:mm:ss")}}</div>
        </div>
      </div>
    </div>
    <div class="footer">
      <van-button plain color="#000000" @click="handleBack">返回列表</van-button>
      <van-button v-if="(detail.status == 1 || detail.status == 2) && isOwn" plain color="#000000" @click="handleCancelOrder"><div class="icon-tips" @click.stop="showTips=true"></div>取消订单</van-button>
      <van-button v-if="detail.status == 1 && isOwn" color="#000000" @click="handleOpenProfile">通知客户填写</van-button>
    </div>
    <!-- 预加载图片 -->
    <img class="preview-image" :src="detail.share_poster_url" alt="">
    <order-preview-popup v-model:show="showSuccessPopup" :content="detail.share_poster_url" cancel-text="取消" @close="handleCancel" @confirm="handleConfirm"></order-preview-popup>
    <!-- 取消订单 -->
    <order-delete-popup v-model="showDeletePopup" title="确定取消当前销售单吗？" okText="确定" @confirm="handleCancelOrderConfirm"></order-delete-popup>
    <!-- 取消订单提示 -->
    <van-popup v-model="showTips" position="bottom" @click-overlay="showTips=false">
      <div class="popup">
        <div class="popup-close" @click="showTips=false"></div>
        <div class="popup-title">取消订单提示</div>
        <div class="popup-content">
          <div class="item">
            <div class="item-title">提示一</div>
            <div class="item-content">取消订单按钮仅在订单处于待填写、待发货状态时显示。</div>
          </div>
          <div class="item">
            <div class="item-title">提示二</div>
            <div class="item-content">当仓库开始处理订单时，不能取消会提示“仓库正在处理当前订单，不能进行取消”。</div>
          </div>
        </div>
      </div>
    </van-popup>
  </div>
</template>
<script lang="ts">
import { cancelOrderApi, getImageMediaApi, getOrderInfoApi } from '@/api/online-order';
import { Component, Vue } from 'vue-property-decorator';
import wxUtils from "@/utils/wx";
import wx from "weixin-js-sdk";
import OrderPreviewPopup from './components/order-preview-popup.vue';
import { storeInfoApi } from '@/api/guide-invite';
import OrderDeletePopup from './components/order-delete-popup.vue';

@Component({
  components: {
    OrderPreviewPopup,
    OrderDeletePopup
  }
})
export default class OnlineOrderDetail extends Vue {
  emptyImage = require('@/assets/material-library/v4.2.0/icon_empty.png')
  statusList = [
    {
      label: '待填写',
      value: 1,
      list: []
    },
    {
      label: '待发货',
      value: 2,
      list: []
    },
    {
      label: '已发货',
      value: 3,
      list: []
    },
    {
      label: '已拒单',
      value: 4,
      list: []
    },
    {
      label: '已作废',
      value: 5,
      list: []
    }
  ]
  get statusText() {
    const current = this.statusList.find(item => item.value === this.detail.status);
    if (!current) return ''
    return current.label
  }

  detail: any = null
  showSuccessPopup = false
  showDeletePopup = false
  showTips = false

  isOwn = false

  async mounted(){
    try {
      wxUtils.wxAgentConfig();
    } catch (error) {
      console.log(error);
    }
    await this.getDetail()
    this.getStoreInfo()
  }

  /**
   * @description:获取门店信息
   */
   async getStoreInfo() {
    const res:any =  await storeInfoApi()
    if (res.code === 200) {
      this.isOwn = this.detail.member_id === res.data.members_id
    } else {
      this.$toast(res.message);
    }
  }

  async getDetail() {

    
    const res:any = await getOrderInfoApi({
      order_id: Number(this.$route.query.order_id)
    })
    // await getImageMediaApi({id: res.data.id})

    if(res.code === 200) {
      this.detail = res.data

      this.statusList.forEach(item => {
        item.list = res.data.sku_list.filter(sku => sku.status === item.value)
      })
    }
  }

  handleBack() {
    // 神策数据上报 -- 订单详情点击
    this.$track.trackOrderDetailClick({
      sales_order_id: this.detail.order_sn,
      order_status: this.statusText,
      order_creator: this.detail.member_name,
      order_belonged_store: this.detail.store_name,
      customer_id: this.detail.external_user_id,
      button_name: '返回列表'
    })
    this.$router.back()
  }

  handleCancelOrder() {
    // 神策数据上报 -- 订单详情点击
    this.$track.trackOrderDetailClick({
      sales_order_id: this.detail.order_sn,
      order_status: this.statusText,
      order_creator: this.detail.member_name,
      order_belonged_store: this.detail.store_name,
      customer_id: this.detail.external_user_id,
      button_name: '取消订单'
    })
    this.showDeletePopup = true
  }

  async handleCancelOrderConfirm(){
    this.$toast.loading({mask:true})
    const res:any = await cancelOrderApi({
      order_id: this.detail.id
    })
    this.$toast.clear()
    if(res.code === 200) {
      this.showDeletePopup = false
      this.$toast('已取消')
      this.getDetail()
    } else {
      this.$toast(res.message)
    }
  }

  handleOpenProfile() {
    // 神策数据上报 -- 订单详情点击
    this.$track.trackOrderDetailClick({
      sales_order_id: this.detail.order_sn,
      order_status: this.statusText,
      order_creator: this.detail.member_name,
      order_belonged_store: this.detail.store_name,
      customer_id: this.detail.external_user_id,
      button_name: '通知客户填写'
    })
    this.showSuccessPopup = true
  }

  handleCancel() {
    this.showSuccessPopup = false

    // 神策数据上报 -- 发送顾客弹框点击
    this.$track.trackToCustomerPopupClick({
      sales_order_id: this.detail.order_sn,
      pop_name: '订单链接已复制，可发给客户进行填写',
      button_name: '取消',
      customer_id: this.detail.external_user_id
    })
  }


  handleConfirm() {
    // const url = location.href.split('?')[0].replace('online-order-detail', 'online-order-info')
    // const text = `您的订单已生成，请点击链接填写收货地址。${url}?order_sn=${this.detail.order_sn}`
    // 复制内容到粘贴板
    // this.$helper.copyMessage(text)


    // 神策数据上报 -- 发送顾客弹框点击
    this.$track.trackToCustomerPopupClick({
      sales_order_id: this.detail.order_sn,
      pop_name: '订单链接已复制，可发给客户进行填写',
      button_name: '发送至客户',
      customer_id: this.detail.external_user_id
    })

    const message = {
      userIds: "", //参与会话的企业成员列表，格式为userid1;userid2;...，用分号隔开。
      externalUserIds: this.detail.external_user_id, // 参与会话的外部联系人列表，格式为userId1;userId2;…，用分号隔开。
      groupName: "", // 会话名称。单聊时该参数传入空字符串""即可。
      chatId: "", // 若要打开已有会话，需指定此参数。如果是新建会话，chatId必须为空串
    };

    // 调用企微接口
    window.wx.invoke("openEnterpriseChat", message, (res) => {
      if (res.err_msg == "openEnterpriseChat:ok") {
        //发送成功
        console.log("成功回调");
        this.showSuccessPopup = false
      }
    });
  }
  
}
</script>
<style lang="scss" scoped>
.container {
  width: 100vw;
  min-height: 100vh;
  background: #F7F7F7;
  padding-bottom: calc(64px + contain(safe-area-inset-bottom));
  padding-bottom: calc(64px + env(safe-area-inset-bottom));

  .content {
    padding: 15px;

    .section {
      padding: 15px;
      background: #FFFFFF;
      margin-bottom: 15px;
      position: relative;

      .flex-1 {
        flex: 1;
        overflow: hidden;
        margin-left: 5px;
      }

      .address {
        font-weight: 600;
        font-size: 14px;
        color: #000000;
        line-height: 20px;
      }

      .info {
        margin-top: 10px;
      }

      .name,
      .phone {
        font-weight: 400;
        font-size: 14px;
        color: #000000;
        line-height: 20px;
        margin-right: 15px;
      }

      .status {
        position: absolute;
        top: 15px;
        left: -6px;
        width: 66px;
        height: 32px;
        background: url('~@/assets/online-order/icon_s_1_l.png') center/100% no-repeat;
        font-weight: 400;
        font-size: 14px;
        color: #000000;
        padding-bottom: 4px;
        display: flex;
        align-items: center;
        justify-content: center;

        &.status-1 {
          background-image: url('~@/assets/online-order/icon_s_0_l.png');
          color: #fff;
        }
        &.status-2 {
          background-image: url('~@/assets/online-order/icon_s_1_l.png');
        }
        &.status-3 {
          background-image: url('~@/assets/online-order/icon_s_2_l.png');
        }
        &.status-4 {
          background-image: url('~@/assets/online-order/icon_s_3_l.png');
        }
        &.status-5 {
          background-image: url('~@/assets/online-order/icon_s_4_l.png');
        }
      }

      .list {
        padding-top: 40px;
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
          }

          .item-name {
            font-weight: 400;
            font-size: 12px;
            color: #96979E;
            line-height: 17px;
            margin-top: 5px;
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

      .box {
        padding-top: 15px;
        border-top: 1px solid #F2F2F2;
        margin-top: 15px;
      }

      .row {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 15px;
        &:last-child {
          margin-bottom: 0;
        }

        .label {
          font-size: 14px;
          color: #96979E;
          line-height: 20px;
        }

        .value {
          font-size: 14px;
          color: #000;
          line-height:20px;

          .van-button {
            height: 20px;
            margin-left: 5px;
          }
        }
      }
    }
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

    .van-button {
      flex: 1;
      padding: 0;
      border: 1px solid #000 !important;
      box-sizing: border-box;
      border-radius: 0;

      &:nth-child(2) {
        margin-left: 15px;
      }
      &:nth-child(3) {
        margin-left: 15px;
      }

      .van-button__text {
        display: flex;
        align-items: center;
        font-size: 14px;
      }
    }

    .icon-tips {
      position: relative;
      z-index: 11;
      width: 20px;
      height: 20px;
      background: url('~@/assets/online-order/icon_tips.png') center/100% no-repeat;
    }
  }
}

.popup {
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
    padding: 0 30px;

    .item {
      margin-bottom: 40px;
    }

    .item-title {
      font-weight: 500;
      font-size: 14px;
      color: #000000;
      line-height: 20px;
    }

    .item-content {
      font-size: 12px;
      color: #676767;
      line-height: 17px;
      margin-top: 10px;
    }
  }

}

.preview-image {
  width: 1px;
  height: 1px;
  position: fixed;
  top: 0;
  right: -100px;
}
</style>