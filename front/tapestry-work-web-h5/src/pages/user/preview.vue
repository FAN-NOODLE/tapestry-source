<template>
  <div class="container">
    <div class="header">
      <img v-if="detail.banner_url" class="banner" :src="detail.banner_url" alt="">
      <img v-else class="banner" src="@/assets/user/default.png" alt="">
      <div class="shadow">
        <img class="shadow-img" src="@/assets/user/icon_shadow.png" alt="">
        <img class="logo" src="@/assets/user/icon_logo.png" alt="">
      </div>
    </div>
    <div v-if="detail.qr_code_member_info" class="content">
      <div v-if="detail.activity_name" class="title">
        <div class="text">{{ detail.activity_name }}</div>
      </div>
      <div v-if="detail.qr_code_member_info?.store_name" class="address">
        <img class="address-img" src="@/assets/user/icon_address.png" alt="">
        <span>{{detail.qr_code_member_info?.store_name}}</span>
      </div>
      <div v-if="detail.qr_code_member_info?.sa_name" class="user">
        <img class="avatar" :src="detail.qr_code_member_info?.avatar" alt="">
        <div class="nickname">{{ detail.qr_code_member_info?.sa_name }}</div>
      </div>

      <div class="qrcode">
        <img class="qrcode-img" src="@/assets/user/demo.png" alt="">
      </div>
      <div class="tips">即刻长按识别联系顾问，尊享专属服务</div>
      <div class="btn">选择其他服务顾问</div>
    </div>
  </div>
</template>
<script lang="ts"> 
import { Component, Vue } from "vue-property-decorator";

@Component({})
export default class  extends Vue {
  searchForm = {
    unionid: '',
    state: '',
    third_id: 0,
    lat: 0,
    lng: 0
  }
  activityDetail:any = null
  detail: ActivityQrResponse = {
    error_code: -1,
    activity_name: "",
    banner_url: "",
    qr_code_member_info: undefined,
    store_list: []
  }
  
  
  mounted() {
    console.log(decodeURIComponent(this.$route.query.url as string));
    
    this.detail = {
      error_code: 0,
      activity_name: "预览活动",
      banner_url: this.$route.query.url ? decodeURIComponent(this.$route.query.url as string) :'',
      qr_code_member_info: {
        store_name: "门店名称",
        member_id: 0,
        user_id: "",
        qr_code: "",
        third_id: 0,
        sa_name: "",
        avatar: ""
      },
      store_list: []
    }
  }
  

}
</script>


<style lang="scss" scoped>
.container {
  background-color: #fff;
  min-height: 100vh;
}
.header {
  max-height: 279px;
  overflow: hidden;
  position: relative;
  .shadow {
    position: absolute;
    bottom:  0; 
    left: 0;
    right: 0;
  }

  .shadow-img {
    display: block;
    height: auto;
    width: 100vw;
  }

  .logo {
    position: absolute;
    top: -25px;
    left: 50%;
    transform: translateX(-50%);
    width: 71px;
    height: 71px;
  }
}
.banner {
  display: block;
  width: 100vw;
}

.content {
  padding-top: 20px;

  .title {
    height: 17px;
    font-size: 17px;
    font-weight: bold;
    display: flex;
    align-items: center;
    justify-content: center;

    .text {
      text-align: center;
      max-width: 200px;
    }
    
    &::before {
      content: "";
      width: 30px;
      height: 1px;
      background: #404042;
      margin-right: 10px;
    }
    &::after {
      content: "";
      width: 30px;
      height: 1px;
      background: #404042;
      margin-left: 10px;
    }
  }

  .address {
    display: flex;
    align-items: flex-start;
    justify-content: center;
    margin-top: 24px;
    font-size: 16px;
    font-weight: 600;
    padding: 0 10px;
    .address-img {
      width: 20px;
      height: 20px;
    }
  }
  .qrcode { 
    width: 168px;
    height: 168px;
    background: url('~@/assets/user/icon_box.png') center/cover no-repeat;
    margin: auto;
    margin-top: 16px;
    box-sizing: border-box;
    padding: 9px;

    .qrcode-img { 
      width: 100%;
      display: block;
    }
  }

  .tips {
    font-size: 14px;
    text-align: center;
    color: #949494;
    margin-top: 24px;
  }

  .user {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 24px;
    margin-top: 24px;

    .avatar  {
      width: 20px;
      height: 20px;
      border-radius: 2px;
      margin-right: 4px;
    }

    .nickname {
      font-size: 16px;
      color: #333;
    }
  }

  .btn {
    color: #0085EA;
    font-size: 14px;
    padding: 20px;
    text-align: center;
  }
}
::v-deep .van-icon {
  position: absolute !important;
}
.popup {
  width: 325px;
  padding: 30px 26px;
  box-sizing: border-box;
  background-color: #fff;

  .popup-header {
    font-size: 16px;
    font-weight: bold;
    line-height: 24px;
  }

  .popup-content { 
    font-size: 14px;
    line-height: 22px;
    margin-top: 24px;
    color: #666666;
  }

  .popup-footer {
    display: flex;
    align-items: center;
    justify-content: center;
    margin-top: 24px;
    .btn {
      background-color: #000;
      color: #fff;
      font-size: 14px;
    }
  }
}
</style>