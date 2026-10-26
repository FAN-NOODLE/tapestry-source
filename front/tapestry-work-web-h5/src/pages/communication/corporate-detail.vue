<template>
  <div v-if="detail" class="detail">
    <div class="info">
      <div class="info__content">
        <div class="title">{{ detail.title }}</div>
        <div class="remark">{{ detail.member_name }}</div>
        <div class="time">{{ formatTime(detail.send_time) }}</div>
        <div class="num">
          <div class="num-item">
            <div class="num-item__icon"></div>
            <div class="num-item__value">{{detail.view_people_num||0}}</div>
          </div>
          <div class="num-item">
            <div class="num-item__icon"></div>
            <div class="num-item__value">{{detail.view_num||0}}</div>
          </div>
        </div>
      </div>
      <van-image class="info__cover" lazy-load :src="detail.cover || defaultCover" fit="contain"></van-image>
    </div>
    <div class="content">
      <div class="item" v-for="(item, index) in detail.content_list" :key="index">
        <div class="item-title">{{ item.department_name }}</div>
        <div class="item-content">
          <div class="item-content__text" v-for="(info, infoIndex) in item.info_list" :key="infoIndex" @click="handleClick(info)">{{ info.title }}</div>
        </div>
      </div>
    </div>

    <div v-if="showDownloadBtn(1)||showDownloadBtn(2)" class="footer">
      <!-- <div class="btn" @click="handleDownload()">下载周传讯PDF</div> -->
      <div v-if="showDownloadBtn(1)" class="btn" @click="handleDownload()">一键下载周传讯及全部附件</div>
      <div v-else-if="showDownloadBtn(2)" class="btn" @click="handleDownload()">下载传讯及全部附件</div>
    </div>
    <!-- 下载弹框 -->
    <com-confirm v-model="showDownloadPopup" @confirm="handleDownloadConfirm"></com-confirm>
    <!-- 下载成功提示 -->
    <com-toast v-model="showDownloadToast"></com-toast>
    <temp-two ref="detail" :data="detail" @success="handleDownloadSuccess"></temp-two>
  </div>
</template>
<script lang="ts">
import { Vue, Component } from 'vue-property-decorator';
import { Image } from "vant";
import { defaultCover } from './const'
import ComConfirm from './components/com-confirm.vue'
import ComToast from './components/com-toast.vue'
import TempTwo from './components/temp-two.vue'
import { getManagerInfoApi, getMergeDetailApi, saveMergePdfApi, sendFileByEmailApi } from '@/api/communication';
import wx from 'weixin-js-sdk'

interface IInfoItem {
  id: number;
  title: string;
}

interface IContentItem {
  department_name: string;
  info_list: IInfoItem[];
}

interface IDetail {
  /** RO传讯ID */
  id: number;
  /** 第几周 */
  week_num: string;
  /** 标题 */
  title: string;
  /** 图url */
  cover: string;
  /** 发送人群类型 1白名单 2直接选人 */
  crowd_type: number;
  /** 白名单id或者member_id集合 */
  crowd_list: any[];
  /** 发送类型  0定时发布，1保存后发布，2暂不发布 */
  send_type: number;
  /** 定时发布时间 */
  send_time: number;
  /** 整合传讯内容 */
  content_list: IContentItem[];
  /** 创建人 */
  member_name: string;
  /** 状态 */
  status: number;
  /** pdf地址 */
  pdf_url: string;
  /** 阅读人数 */
  view_people_num: number;
  /** 阅读次数 */
  view_num: number;
  /** 类型 */
  type: number;
}

@Component({
  components: {
    [Image.name]: Image,
    ComConfirm,
    ComToast,
    TempTwo
  }
})
export default class CorporateDetail extends Vue {
  /** 默认封面 */
  defaultCover = defaultCover
  /** 下载确认 */
  showDownloadPopup = false
  /** 下载城通提示 */
  showDownloadToast = false
  /** 详情 */
  detail: IDetail|null = null;
  email = ''

  /** 用户信息 */
  get userInfo(): IUserInfo {
    return this.$store.getters["user/getUserInfo"] || {}
  }

  async created() {
    if(this.userInfo.app_source !== 2) {
      await this.getAuth()
    }
    // this.getDetail()
    if(this.$route.query.from !== '1') {
      // 当前需要直接退出的页面
      window.history.pushState("_android_back", "", location.href);
      window.addEventListener("popstate", this.popstate, false);
      // 神策数据上报 -- 进入传讯
      this.$track.trackEnterTransmittal({
        transmittal_entrance: 2
      })
    }
    
  }
  activated() {
    if(!this.detail || Number(this.$route.query.id) !== this.detail.id) {
      console.log(1111);
      this.detail = null
      this.getDetail()
    }
  }
  destroyed() {
    // 记得清除!!!!不然单页面，其它页面也可以用
    window.removeEventListener("popstate", this.popstate, false);
  }
  deactivated() {
    // 记得清除!!!!不然单页面，其它页面也可以用
    window.removeEventListener("popstate", this.popstate, false);
  }
  

  popstate() {
    wx.closeWindow(); // 微信网页退出
  }


  /**
   * @description 获取用户权限
   */
   async getAuth() {
    try {
      const {code, message, data}:IApiResult = await getManagerInfoApi({})
      if (code!== 200) {
        this.$toast.fail(message);
        return;
      }
      this.$store.commit("user/setUserInfo", {
        app_source: 2,
        ...this.userInfo, 
        ...data,
      })
    } catch (error) {
      console.log(error);
    }
  }

  get showDownloadBtn() {
    return (type) => {
      if(!this.detail) return false
      if(type===1) {
        return this.detail.type === 1
      }
      if(type===2) {
        const {content_list} = this.detail
        return content_list.length>1 || content_list[0].info_list.length>1
      }

      return false
    }
  }

  /**
   * 时间格式化
   * @param time 
   */
  formatTime(time: number) {
    return this.$helper.dateFormat(time, "yyyy/MM/dd hh:mm");
  }

  /**
   * 获取详情
   */
  async getDetail() {
    try {
      this.$toast.loading({
        duration: 0,
        forbidClick: true,
        message: '加载中...'
      })
      const {code, message, data}: IApiResult = await getMergeDetailApi({
        id: Number(this.$route.query.id)
      })
      this.$toast.clear();
      if (code!== 200) {
        this.$toast.fail(message);
        return;
      }
      this.detail = data
      console.log(data);
      
    } catch (error) {
      this.$toast.clear();
      console.log(error);
    }
  }


  /**
   * 下载
   */
  handleDownload() {
    this.showDownloadPopup = true;
    
    // 神策数据上报 -- 点击下载周传讯pdf
    this.$track.trackClickAllTransmittalDownload({
      transmittal_type: this.detail.type
    })
  }
  /**
   * 下载确认
   */
  async handleDownloadConfirm(value:string) {
    this.email = value;
    // if(this.detail.pdf_url) {
    this.$toast.loading({
      message: '正在提交',
      duration: 0,
      forbidClick: true,
      loadingType: 'spinner'
    })
    this.handleDownloadSuccess({
      url: this.detail.pdf_url,
      rename: (this.$refs.detail as any).filename(true)
    })
    // } else {
    //   (this.$refs.detail as any).handleDownloadPdf()
    // }
    // (this.$refs.detail as any).handleDownloadPdf()
  }

  async handleDownloadSuccess(e) {
    try {
      const params:any = {
        merge_id: this.detail.id,
        email: this.email,
        file_url: e.url,
        rename: e.rename
      }
      const {code, message, data}: IApiResult = await sendFileByEmailApi(params)
      if (code!== 200) {
        this.$toast.fail(message);
        return;
      }
      this.showDownloadToast = true;
      this.showDownloadPopup = false
      this.$toast.clear()
      // 神策数据上报 -- 下载传讯附件
      this.$track.trackDownloadTransmittal({
        transmittal_title: this.detail.title,
        transmittal_publisher: this.detail.member_name,
        release_time: this.detail.send_time,
      })
    } catch (error) {
      console.log(error);
    }
  }
  
  /**
   * 查看传讯详情
   * @param info 
   */
  handleClick(info: IInfoItem) {
    this.$router.push({
      path: 'communication-detail/sa',
      query: {
        id: info.id + '',
        from: '1'
      }
    });
  }
}

</script>
<style lang="scss" scoped>
.detail {
  min-height: 100vh;
  padding: 15px;
  background-color: #f7f7f7;
  padding-bottom: calc(56px + contain(safe-area-inset-bottom));
  padding-bottom: calc(56px + env(safe-area-inset-bottom));


  .info {
    display: flex;
    padding: 0 5px 25px;
    border-bottom: 1px dashed rgba(166, 166, 166, 0.4);
    &__cover {
      width: 108px;
      height: 108px;
      background: rgba(0, 0, 0, 0.2);
      // box-shadow: 0px 1px 3px 0px rgba(0,0,0,0.1);

      .van-image {
        width: 100%;
        height: 100%;
      }
    }

    &__content {
      flex: 1;
      overflow: hidden;
      padding-right: 20px;

      .title {
        font-weight: 600;
        font-size: 14px;
        color: #000000;
        line-height: 20px;
        padding-top: 5px;
      }

      .remark {
        height: 17px;
        font-weight: 400;
        font-size: 12px;
        color: #000000;
        line-height: 17px;
        margin-top: 10px;
      }

      .time {
        height: 17px;
        font-weight: 400;
        font-size: 12px;
        color: #A6A6A6;
        line-height: 17px;
        margin-top: 4px;
      }

      .num {
        display: flex;
        align-items: center;
        margin-left: -6px;

        .num-item {
          display: flex;
          align-items: center;
          font-size: 14px;
          color: #777777;
          margin-right: 18px;

          &__icon {
            width: 24px;
            height: 24px;
            background: url("~@/assets/communication/icon_user.png") center/cover no-repeat;
          }
        }
        .num-item:nth-child(2) .num-item__icon {
          background-image: url("~@/assets/communication/icon_view.png");
        }
      }
    }
  }

  .item {
    padding-top: 20px;

    .item-title {
      font-weight: 600;
      font-size: 14px;
      color: #A6A6A6;
      line-height: 20px;
      padding: 0 5px;
    }

    .item-content {

      &__text {
        padding: 15px 20px 15px 5px;
        font-weight: 400;
        font-size: 14px;
        color: #000000;
        line-height: 20px;
        border-bottom: 1px dashed rgba(166, 166, 166, 0.4);
        position: relative;

        &::after {
          content: "";
          position: absolute;
          right: 0;
          top: 50%;
          transform: translateY(-50%);
          width: 15px;
          height: 15px;
          background: url('~@/assets/communication/icon_right.png') center/100% no-repeat;
        }
      }
    }
  }

  .footer {
    font-size: 12px;
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    padding: 10px 20px;
    z-index: 2;
    display: flex;
    align-items: center;
    background-color: #fff;
    padding-bottom: calc(10px + contain(safe-area-inset-bottom));
    padding-bottom: calc(10px + env(safe-area-inset-bottom));

    .btn {
      flex: 1;
      height: 36px;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #000000;
      color: #fff;

      &.plain {
        background: #fff;
        color: #000;
        border: 1px solid #000;
      }

      &:not(:first-child) {
        margin-left: 5px
      }
    }
  }
  
}
</style>