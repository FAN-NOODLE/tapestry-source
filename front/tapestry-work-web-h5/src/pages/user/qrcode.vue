<template>
  <div v-if="!loading" class="container">
    <div v-if="!showStoreList">
      <div class="header">
        <img v-if="detail.banner_url" class="banner" :src="detail.banner_url" alt="">
        <img v-else class="banner" src="@/assets/user/default.png" alt="">
        <div class="shadow">
          <img class="shadow-img" src="@/assets/user/icon_shadow.png" alt="">
          <img class="logo" src="@/assets/user/icon_logo.png" alt="">
        </div>
      </div>
      <div v-if="detail.qr_code_member_info && detail.qr_code_member_info.qr_code" class="content">
        <div v-if="detail.activity_name" class="title">
          <div class="text">{{ detail.activity_name }}</div>
        </div>
        <div v-if="detail.qr_code_member_info?.store_name" class="address">
          <img class="address-img" src="@/assets/user/icon_address_3.png" alt="">
          <span>{{detail.qr_code_member_info?.store_name}}</span>
        </div>
        <div v-if="detail.qr_code_member_info?.sa_name" class="user">
          <img class="avatar" :src="detail.qr_code_member_info?.avatar" alt="">
          <div class="nickname">{{ detail.qr_code_member_info?.sa_name }}</div>
        </div>
  
        <div v-if="detail.qr_code_member_info?.qr_code" class="qrcode">
          <img class="qrcode-img" :src="detail.qr_code_member_info?.qr_code" alt="" @touchstart="handleTouchStart" @touchend="handleTouchEnd">
        </div>
        <div v-if="detail.qr_code_member_info?.qr_code" class="tips">即刻长按识别联系顾问，尊享专属服务</div>
        <div v-if="detail.qr_code_member_info?.qr_code" class="btn" @click="handleSelectOtherMember">选择其他门店服务顾问</div>
      </div>
      <div v-else class="content">
        <div class="title">
          <div class="text">活动名称</div>
        </div>
        <div class="address">
          <img class="address-img" src="@/assets/user/icon_address_3.png" alt="">
          <span>门店名称</span>
        </div>
        <!-- <div v-if="detail.qr_code_member_info?.sa_name" class="user">
          <img class="avatar" :src="detail.qr_code_member_info?.avatar" alt="">
          <div class="nickname">{{ detail.qr_code_member_info?.sa_name }}</div>
        </div> -->
  
        <div class="qrcode">
          <img class="qrcode-img" src="@/assets/user/demo.png" alt="">
        </div>
        <div class="tips">即刻长按识别联系顾问，尊享专属服务</div>
        <div class="btn">选择其他门店服务顾问</div>
      </div>
      <van-popup
        v-model="visible"
        position="center"
      >
        <div class="popup">
          <div class="popup-header">{{ title }}</div>
          <div class="popup-content">{{ message }}</div>
          <div class="popup-footer">
            <van-button class="btn" block @click="handleBack">知 道 了</van-button>
          </div>
        </div>
      </van-popup>
    </div>
    <!-- <StoreList v-else :form="searchForm" /> -->
  </div>
</template>
<script lang="ts"> 
import { Component, Vue } from "vue-property-decorator";
import wxUtil from '@/utils/wx';
import { getActivityDetailApi, getQrcodeByUnionidApi } from "@/api/user";
import { sensorsLogin } from "@/utils/sensors";
import { getDeviceInfo } from "@/utils/tools";
import StoreList from './store.vue'

@Component({
  components: {
    StoreList
  }
})
export default class QrcodePage extends Vue {
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
  visible = false;
  title = '活动暂未开始';
  message = '抱歉，当前活动暂未开始，请关注最新活动动态～'
  timer: NodeJS.Timeout | null = null;
  loading = true

  showStoreList = false

  get trackInfo() {
    return (data) => {
      return {
        display_campaign_name: data.activity_name,
        store_name: data.qr_code_member_info?.store_name || '',
        store_code: data.qr_code_member_info?.store_code || '',
        SA_name: data.qr_code_member_info?.sa_name || '',
        SA_id: data.qr_code_member_info?.user_id || '',
        QR_code_id: data.qr_code_member_info?.qr_code || '',
        leading_words: '即刻长按识别联系顾问，尊享专属服务',
        SA_added_status: [
          '已添加过投放范围内的导购',
          '选择其它服务顾问',
          '未添加过投放范围内的导购，有绑定过投放范围内的绑定柜台',
          '未添加过投放范围内的导购，且没有绑定对应的管理柜台'
        ][data.status] || '',
        campaign_status: {
          0: '进行中',
          1: '已失效',
          2: '未开始',
          3: '已结束'
        }[data.error_code] || '其他',
      }
    }
  }
  async mounted() {
    // const host = window.location.origin
    // const hash = window.location.pathname.split("/")[2];
    // const IDENTIFIER = window.location.pathname.split("/")[1] || "";
    // const url = `${host}/${IDENTIFIER}/${hash}${this.$route.fullPath}`
    // if(!this.$route.query.s) {
    //   const url = window.location.href
    //   window.location.href = url + '&s=1'
    // }

    // await wxUtil.wxConfig()
    if(!this.$route.query.unionid) {
      this.$toast('请传入unionid参数')
      return 
    }

    this.searchForm.unionid = String(this.$route.query.unionid)
    this.searchForm.state = (this.$route.query.scene) as string
    if(this.$route.query.third_id) {
      this.searchForm.third_id = Number(this.$route.query.third_id as string)
    }
    await this.getActivityDetail()
    // 登录神策上报
    sensorsLogin(this.searchForm.unionid)
    this.getData();
  }

  async getActivityDetail() {
    const res = await getActivityDetailApi({
      state: this.searchForm.state
    })
    this.activityDetail = res.data
    localStorage.setItem('activityDetail', JSON.stringify(res.data))
  }


  async getData() {
    this.$toast.loading({
      mask: true,
      message: '加载中...',
      duration: 0
    })
    const res = await getQrcodeByUnionidApi(this.searchForm)

    this.$toast.clear()
    let status = this.$route.query.third_id ? 1 : -1
    if(res.data.error_code === 0) {
      if(res.data.qr_code_member_info?.user_id) {
        status = 0 // 页面显示了导购 -- 已添加过投放范围内的导购
      } else if(res.data.qr_code_member_info?.qr_code) {
        status = 2 // 页面上没有导购但是有二维码 --  未添加过投放范围内的导购，有绑定过投放范围内的绑定柜台
      } else {
        status = 3 // 页面上没有导购并且没有二维码 -- 未添加过投放范围内的导购，且没有绑定对应的管理柜台
      }
    }

    // 导购ID - user_id
    // 已选门店 - this.$route.query.third_id
    // 已添加过投放范围内的导购 user_id
    // 选择其它服务顾问 this.$route.query.third_id
    // 未添加过投放范围内的导购，有绑定过投放范围内的绑定柜台 qr_code && !user_id
    // 未添加过投放范围内的导购，且没有绑定对应的管理柜台 !qr_code
    // 神策数据上报 - 记录页面浏览
    res.data.status = status
    this.$track.trackAddSAPageView(this.trackInfo(res.data))
    this.detail = res.data
    if(res.data.error_code === 0) {
      if(!res.data.qr_code_member_info?.qr_code) {
        console.log(11111, localStorage.getItem('lbs'));
        
        if(localStorage.getItem('lbs')) {
          this.$router.push({
            path: '/user/store',
            query: {
              unionid: this.searchForm.unionid,
              scene: this.searchForm.state,
            }
          })
          // const url = window.location.href.replace('/user/qrcode', '/user/store').split('?')[0]
          // window.location.href = url + `?unionid=${this.searchForm.unionid}&scene=${this.searchForm.state}&t=${new Date().getTime()}`
      
        } else {
          this.loading = false
          this.handleGetLocation()
        }
      } else {
        this.loading = false
        this.detail = res.data
      }
    } else if(res.data.error_code === 1) {
      this.loading = false
      this.visible = true
      this.title = '活动已失效'
      this.message = '抱歉，当前活动已失效，请关注最新活动动态～'
    } else if(res.data.error_code === 2) {
      this.loading = false
      this.visible = true
      this.title = '活动暂未开始'
      this.message = '抱歉，当前活动暂未开始，请关注最新活动动态～'
    } else if(res.data.error_code === 3) {
      this.loading = false
      this.visible = true
      this.title = '活动已结束'
      this.message = '抱歉，当前活动已结束，请期待下次为您服务！'
    } else {
      this.loading = false
      this.visible = true
      this.title = '活动已失效'
      this.message = '抱歉，当前活动已失效，请关注最新活动动态～'
    }

  }

  handleSelectOtherMember() {
    this.$track.trackChooseOtherSA(this.trackInfo(this.detail))
    this.$router.push({
      path: '/user/store',
      query: {
        unionid: this.searchForm.unionid,
        scene: this.searchForm.state,
        other: '1'
      }
    })
    // const url = window.location.href.replace('/user/qrcode', '/user/store').split('?')[0]
    // window.location.href = url + `?unionid=${this.searchForm.unionid}&scene=${this.searchForm.state}&other=1&t=${new Date().getTime()}`
  }

  // 长安0.5秒后 上报数据
  handleTouchStart() {
    if (this.timer) {
      clearTimeout(this.timer);
    }
    this.timer = setTimeout(() => {
      this.$track.trackScanQRCode(this.trackInfo(this.detail))
    }, 500)
  }
  handleTouchEnd() {
    if (this.timer) {
      clearTimeout(this.timer);
    }
  }
  handleBack() {
    this.visible = false;
    (window.wx as any).miniProgram.navigateBack();
  }
  
  handleGetLocation() {
    // const host = window.location.origin
    // const hash = window.location.pathname.split("/")[2];
    // const IDENTIFIER = window.location.pathname.split("/")[1] || "";
    // const url = `${host}/${IDENTIFIER}/${hash}${this.$route.fullPath}`
    
    wxUtil.getLocation().then((res:any) => {

      this.searchForm.lat = Number(res.latitude)
      this.searchForm.lng = Number(res.longitude)

      localStorage.setItem('lbs', '1')
      
      this.$track.trackLBSAuth({
        button_type: '授权',
      })
      // const url = window.location.href.replace('/user/qrcode', '/user/store').split('?')[0]

      // window.location.href = url + `?unionid=${this.searchForm.unionid}&scene=${this.searchForm.state}&lat=${res.latitude}&lng=${res.longitude}&t=${new Date().getTime()}`
      this.$router.push({
        path: '/user/store',
        query: {
          unionid: this.searchForm.unionid,
          scene: this.searchForm.state,
          lat: res.latitude,
          lng: res.longitude
        }
      })

      
    }).catch(err => {
      console.error(err);
      this.$track.trackLBSAuth({
        button_type: '拒绝',
      })
      localStorage.setItem('lbs', '2')
      this.$router.push({
        path: '/user/store',
        query: {
          unionid: this.searchForm.unionid,
          scene: this.searchForm.state
        }
      })
      // const url = window.location.href.replace('/user/qrcode', '/user/store').split('?')[0]

      // window.location.href = url + `?unionid=${this.searchForm.unionid}&scene=${this.searchForm.state}&t=${new Date().getTime()}`
      
    });
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