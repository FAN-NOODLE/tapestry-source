<template>
  <div class="container">
    <div class="content">
      <div class="scroll">
        <div class="store">
          <div class="store-code">{{storeInfo.store_code}}：{{storeInfo.store_name}}</div>
        </div>
        <div id="input1" class="section">
          <div class="section-title">销售单号</div>
          <div class="section-content">
            <van-field ref="input" v-model="formData.order_sn" class="input" contenteditable autofocus placeholder="请输入正确的销售单号，确保订单处理无误。" @click="handleClick" @focus="handleFocus('input1')"></van-field>
          </div>
        </div>
        <div id="input2" class="section">
          <div class="section-title">添加商品</div>
          <div class="section-content">
            <add-goods ref="addGoods" v-model:value="formData.sku_list"></add-goods>
          </div>
        </div>
        <div class="section">
          <div class="title">网仓配送时效</div>
          <div class="label">工作日</div>
          <div class="value">
            <div>08:00-16:00下单，将于当日18:00前发货</div>
            <div class="mt-5">16:00-次日08:00下单，将于次日10:00前发货</div>
          </div>
          <div class="label">周末及公众节假日</div>
          <div class="value">
            <div>每日发货一次</div>
          </div>
        </div>
      </div>
    </div>
    <div class="footer">
      <van-button block color="#000000" @click="handleSubmit">提交订单</van-button>
    </div>
    <!-- 确认弹框 -->
    <order-confirm-popup v-model:show="showConfirmPopup" :orderSn="formData.order_sn" :list="formData.sku_list" @cancel="showConfirmPopup=false" @confirm="handleSubmitConfirmBefore"></order-confirm-popup>
    <!-- 提交成功 -->
    <order-success-popup v-model:show="showSuccessPopup" @close="handleCancel" @confirm="handleConfirm"></order-success-popup>
    <!-- 销售单号重复提示 -->
    <van-popup v-model="showRepeatPopup" @click-overlay="showRepeatPopup=false">
      <div class="popup">
        <div class="popup-close" @click="showRepeatPopup=false"></div>
        <div class="popup-content">
          <van-image width="50" height="50" :src="require('@/assets/online-order/icon_success.png')"></van-image>
          <div class="title">重要提醒</div>
          <div>【销售单号：
            <span style="color:#F64040">{{repeatInfo.order_sn}}</span>
            】</div>
          <div class="text">
            <div>监测到“{{ repeatInfo.order_time }}”已提交过</div>
            <div>请再次核对是否要提交订单！</div>
          </div>
        </div>
        <div class="popup-footer">
          <van-button color="#000000" @click="handleRepeatConfirm">已核对，确定提交</van-button>
          <van-button color="#000000" plain @click="showRepeatPopup=false">返回修改</van-button>
        </div>
      </div>
    </van-popup>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import AddAddress from './components/add-address.vue';
import AddGoods from './components/add-goods.vue';
import OrderConfirmPopup from './components/order-confirm-popup.vue';
import OrderSuccessPopup from './components/order-success-popup.vue';
import { storeInfoApi, userInfoApi } from '@/api/guide-invite';
import { checkOrderRepeatApi, createOrderApi, getImageMediaApi } from '@/api/online-order';
import wxUtils from "@/utils/wx";
import wx from "weixin-js-sdk";

@Component({
  components: {
    OrderSuccessPopup,
    OrderConfirmPopup,
    AddGoods,
    AddAddress
  }
})
export default class OnlineOrderCreate extends Vue {
  formData: any = {
    user_name: '',
    // 销售编码
    order_sn: '',
    // 商品列表
    sku_list: [],
    /** 外部联系人 */
    external_user_id: '',
    /** 顾客唯一标识 */
    union_id: '',
    /** 部门ID */
    third_id: '',
    media_id: ''
  }
  // 确认弹框
  showConfirmPopup = false
  // 成功弹框
  showSuccessPopup = false

  showRepeatPopup = false

  storeInfo:any = {}
  submitInfo:any = {}
  repeatInfo:any = {}

  async mounted() {
    if(this.$route.query.userid) {
      this.formData.external_user_id = this.$route.query.userid
      this.$store.commit("user/setExternalId", this.$route.query.userid);
      this.getUserInfo()
    }
    
    try {
        if(!this.formData.external_user_id) {
          await wxUtils.wxAgentConfig();
          await this.getWxContent();
        }
        this.getUserInfo()
      } catch (error) {
        console.log(error);
      }
    
    this.getStoreInfo()

    // 当前需要直接退出的页面
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
  handleClick() {
    (this.$refs.input as any).focus()
  }

  handleFocus(id) {
    // 页面滚动

    this.$nextTick(() => {
      const catShopElement = document.getElementById(id);
      if (catShopElement) {
        document.documentElement.scrollTo({
          top: catShopElement.offsetTop, // 滚动到元素的顶部
          behavior: "smooth"
        });
      }

    })
  }

  async getWxContent() {
    if (!window.wx.invoke) return;
    const { entry } = await wxUtils.getContext() as any
    if (entry == "single_chat_tools") {
      const { userId } = await wxUtils.getCurExternalContact() as any
      // 神策数据上报 -- 点击侧边栏
      this.$track.trackClickItem({
        feature: this.$route.path,
        customer_id: userId
      })
      this.formData.external_user_id = userId;
    } if(entry == "group_chat_tools") {
      await wxUtils.getCurExternalChat()
    }
  }

  /**
   * @description:获取用户信息
   */
   async getUserInfo() {
    try {
      const { code, data, message }: any = await userInfoApi({
        external_userid: this.formData.external_user_id,
      });
      if (code === 200) {
        this.formData.user_name = data.name || data.nick
        this.formData.union_id = data.union_id	
      } else {
        this.$toast(message);
      }
    } catch (error) {
      console.error(error, "获取用户userInfo错误");
    }
  }


  /**
   * @description:获取门店信息
   */
  async getStoreInfo() {
    const res:any =  await storeInfoApi()
    if (res.code === 200) {
      this.storeInfo = res.data;
      this.formData.third_id = res.data.department_third_id
    } else {
      this.$toast(res.message);
    }
  }
  async sendChatMessage() {
    const res = await getImageMediaApi({ id: this.formData.id })
    
    const message = {
      msgtype: "image",
      image: {
        mediaid: res.data.media_id
      },
    };
    window.wx.invoke("sendChatMessage", message)

  }

  /**
   * 提交订单
   */
  handleSubmit() {
    if(this.showConfirmPopup || this.showSuccessPopup) return

    if(!this.formData.order_sn.trim()) {
      this.$toast('请输入销售单号');
      return 
    }
    // const isLoading = this.formData.sku_list.some(item => item.loading)
    // if(isLoading) {
    //   this.$toast('网仓库存查询中，请稍后提交')
    //   return
    // }
    if(!this.formData.sku_list.length) {
      this.$toast('请添加商品');
      return 
    }
    this.showConfirmPopup = true
  }
  /**
   * 提交订单前 判断销售单号是否重复
   * @param info 
   */ 
  async handleSubmitConfirmBefore(info) {
    const res:any = await checkOrderRepeatApi({order_sn: this.formData.order_sn})
    if(res.code === 200) {
      if(res.data.order_sn) {
        this.showRepeatPopup = true
        this.submitInfo = info
        this.repeatInfo = res.data
      } else {
        this.handleSubmitConfirm(info)
      }
      
    } else {
      this.handleSubmitConfirm(info)
    }
  }
  /**
   * 提交订单 确认
   */
  async handleSubmitConfirm(info) {
    const params = JSON.parse(JSON.stringify(this.formData))
    
    this.$toast.loading({mask:true})
    // 去除空格
    params.order_sn = params.order_sn.replaceAll(' ', '')
    params.sku_list = params.sku_list.map(item => {
      if(item.loading) item.status = 0
      else if(item.dcl_on_hand ===-1) item.status = 4
      else if(item.dcl_on_hand ===0) item.status = 1
      else if(item.dcl_on_hand < item.quantity) item.status = 2
      else item.status = 3
      return {
        image_url: item.image,
        name: item.product_name,
        sku_info: item.product_code,
        upc_code: item.upc_code,
        quantity: item.quantity,
        status: item.status,
        sku_number: item.sku_number,
        style_color: item.style_color,
        dept: item.department,
      }
    })
    const res:any = await createOrderApi(params)
    if(res.code === 200) {
      this.formData.id = res.data.id
      this.showConfirmPopup = false
      this.showSuccessPopup = true
      
      // 神策数据上报 -- 订单生成
      this.$track.trackOrderCreated({
        sales_order_id: params.order_sn,
        has_reject_risk: info.has_reject_risk,
        risk_type: info.risk_type,
        customer_id: params.external_user_id
      })

      params.sku_list.forEach(item => {
        // 神策数据上报 -- 订单生成明细
        this.$track.trackOrderCreatedDetail({
          sales_order_id: params.order_sn,
          commodity_name: item.name,
          sku_id: item.sku_info,
          commodity_quantity: item.quantity,
          inventory_inquiry_status: ['查询中', '无库存', '库存不足', '有库存', '商品库存未知'][item.status],
          customer_id: params.external_user_id
        })
      });

    } else {
      this.$toast.fail(res.message)
    }
  }

  /**
   * 发送消息
   */
  async handleConfirm() {
    // 神策数据上报 -- 发送顾客弹框点击
    this.$track.trackToCustomerPopupClick({
      sales_order_id: this.formData.order_sn,
      pop_name: '订单已生成，请发给客户填写收货信息',
      button_name: '发送至客户',
      customer_id: this.formData.external_user_id
    })
    await this.sendChatMessage()
    this.showConfirmPopup = false
    this.formData.order_sn = '' 
    this.formData.sku_list = [];
    (this.$refs.addGoods as any).handleClear()
  }

  /**
   * 继续创建
   */
  handleCancel() {
    // 神策数据上报 -- 发送顾客弹框点击
    this.$track.trackToCustomerPopupClick({
      sales_order_id: this.formData.order_sn,
      pop_name: '订单已生成，请发给客户填写收货信息',
      button_name: '继续创建',
      customer_id: this.formData.external_user_id
    })
    this.showConfirmPopup = false
    this.formData.order_sn = '' 
    this.formData.sku_list = [];
    (this.$refs.addGoods as any).handleClear()
  }

  handleRepeatConfirm() {
    this.showRepeatPopup = false
    this.handleSubmitConfirm(this.submitInfo)
  }
}
</script>
<style lang="scss" scoped>
::v-deep .van-field__control {
  user-select: auto;
  -webkit-user-select: auto;
}
::v-deep .van-field__body {
  height: 100%;
}
.van-button {
  border-radius: 0;
}
.btn {
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #000000;
  color: #fff;
  font-size: 14px;

  &.plain {
    background: #fff;
    color: #000;
    border: 1px solid #000;
  }
}
.container {
  width: 100vw;
  height: 100vh;
  background: #F7F7F7;
  display: flex;
  flex-direction: column;

  .content {
    flex: 1;
    overflow: hidden;
    .scroll {
      height: 100%;
      overflow-y: scroll;
      padding: 15px;
      box-sizing: border-box;
    }

    .store {
      display: flex;
      align-items: flex-start;
      margin-bottom: 15px;
      line-height: 20px;

      &::before {
        content: "";
        width: 20px;
        height: 20px;
        background: url(~@/assets/online-order/icon_store.png) center/100% no-repeat;
        margin-right: 8px;
      }

      .line {
        width: 1px;
        height: 13px;
        border: 1px solid #D1BFA7;
        margin: 0 12px;
      }

      .store-code, 
      .store-name {
        font-size: 14px;
        color: #000;
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

        &::after {
          content: "*";
          color: #F64040;
          margin-left: 2px;
        }
      }

      .section-content {
        padding-top: 10px;

        .input {
          border: 1px solid #000;
          padding: 0 10px;
          font-size: 14px;
          font-weight: bold;
          height: 40px;
          width: 100%;
          border-radius: 0;
          user-select: all;

          ::placeholder {
            color: #A6A6A6;
          }
        }

      }
      .label {
        font-weight: 400;
        font-size: 12px;
        color: #000000;
        line-height: 17px;
        margin-top: 15px;
        margin-bottom: 10px;
        position: relative;
        padding-left: 7px;

        &::before {
          content: "";
          position: absolute;
          left: 0;
          top: 50%;
          transform: translateY(-50%);
          width: 2px;
          height: 10px;
          background-color: #000;
        }
      }

      .value {
        font-weight: 400;
        font-size: 12px;
        color: #000000;
        line-height: 17px;

        .mt-5 {
          margin-top: 5px;
        }
      }

      .title {
        font-weight: 500;
        font-size: 12px;
        color: #000000;
        line-height: 17px;
        text-align: center;
      }
    }
  }

  .footer {
    // position: fixed;
    // bottom: 0;
    // left: 0;
    // right: 0;
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
}

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
    overflow: hidden;
    font-size: 14px;
    line-height: 20px;
    .van-image {
      display: block;
      margin: auto;
    }
    .title {
      width: 160px;
      margin: 15px auto;
      font-weight: 500;
      font-size: 16px;
      color: #000000;
      line-height: 22px;
    }
    .text {
      margin-top: 5px;
    }
  }

  .popup-footer {
    padding: 20px 15px;
    overflow: hidden;

    .van-button {
      width: 100%;
      height: 36px;
      margin-top: 0;
      display: block;

      &:nth-child(2) {
        margin-top: 10px;
      }
    }
  }
}

</style>