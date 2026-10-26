<template>
  <div v-if="detail" class="detail" :style="$route.params.type=='ft'?{paddingBottom: '20px'}:{}">
    <div v-if="$route.params.type != 'sa'" class="header">
      <div v-if="detail.status==3||detail.status==4||detail.status==5" class="audit">
        <div class="auditor">审核者：{{detail.review_name}}</div>
        <div class="audit-time">审核时间：{{formatTime(detail.review_time)}}</div>
        <div v-if="detail.comment" class="audit-remark">
          <div class="audit-remark-label">留言：</div>
          <div class="audit-remark-value">{{ detail.comment ||'--'}}</div>
        </div>
      </div>
      <div class="header-inner">
        <div class="title">
          <span>{{ detail.department_name }}</span>
          <span> | </span>
          <span :class="{type: detail.type === 3}">{{ detail.type_str }}</span>
        </div>
        <div class="time">
          <span>{{ formatTime(detail.submit_time) }}</span>
          <span class="submitter">提交者：{{ detail.member_name }}</span>
        </div>
      </div>
      <div class="status" :data-status="detail.status">{{ detail.status_str }}</div>
    </div>
    <div class="content">
      <div v-if="detail.type===3 && detail.theme_title" class="theme-title">{{ detail.theme_title }}</div>
      <div class="title">{{ detail.title }}</div>
      <div v-if="$route.params.type=='sa'" class="num">
        <div class="num-item">
          <div class="num-item__icon"></div>
          <div class="num-item__value">{{detail.view_people_num||0}}</div>
        </div>
        <div class="num-item">
          <div class="num-item__icon"></div>
          <div class="num-item__value">{{detail.view_num||0}}</div>
        </div>
      </div>
      <div class="row">
        <div class="row-label width-auto">至：</div>
        <div class="row-value">
          <div class="until-name">{{ detail.until_name }}</div>
        </div>
      </div>
      <div class="row">
        <div class="row-label width-auto">内容：</div>
        <div class="row-value">
          <div class="main" v-html="formatContent(detail.content)"></div>
        </div>
      </div>
      <div v-if="detail.perform" class="row">
        <div class="row-label">执行事务：</div>
        <div class="row-value">
          <div class="main" v-html="formatContent(detail.perform)"></div>
        </div>
      </div>
      <div v-if="detail.date_time" class="row">
        <div class="row-label">时间：</div>
        <div class="row-value">{{ detail.date_time }}</div>
      </div>
      <div v-if="detail.contact_info && detail.contact_info.length" class="row">
        <div class="row-label">联系人：</div>
        <div class="row-value">
          <div class="contact-item" v-for="(item,index) in detail.contact_info" :key="index">
            <div v-if="item.department_name">{{ item.department_name }}</div>
            <div class="contact-item__value">
              <span>{{ item.name }}</span>
              <span v-if="item.tel_phone || item.email"> / {{ [item.tel_phone,item.email].filter(Boolean).join(' / ') }}</span>
              <span v-if="item.remark">（{{ item.remark }}）</span>
            </div>
          </div> 
        </div> 
      </div>
      <!-- 周传讯显示 -->
      <div v-if="detail.type==1&&detail.page_url !== '2'" class="bottom">
        <div class="row-text">负责跟进同事：</div>
        <div class="row-text">执行行动日期：</div>
        <div class="row-text">完成行动日期：</div>
        <div class="line"></div>
      </div>
      <!-- 附件 -->
      <div v-if="detail.attachment.length" class="attach">
        <div class="attach-title">相关附件</div>
        <div class="attach-list">
          <div class="attach-item" v-for="(item, index) in detail.attachment" :key="index"> 
            <div class="attach-item-label">{{item.name}}{{index+1}}</div>
            <div class="attach-item-type" :data-type="item.type"></div>
            <div class="attach-item-name">{{item.attachment_name}}</div>
            <div class="attach-item-btn" @click="handlePreview(item)">
              <div class="icon-preview"></div>
              <div class="text">预览</div>
            </div>
            <div class="attach-item-btn" @click="handleDownload(item)">
              <div class="icon-download"></div>
              <div class="text">下载</div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <!-- RO审核 -->
    <div v-if="$route.params.type=='ro' && detail.status==1" class="footer">
      <div class="btn plain" @click="handleCancel">取消</div>
      <div class="btn plain" @click="handleAudit(0)">不通过</div>
      <div class="btn" @click="handleAudit(1)">通过</div>
    </div>
    <!-- SA视角 -->
    <div v-if="$route.params.type=='sa'" class="footer">
      <!-- <div class="btn" @click="handleDownload()">下载传讯文件</div> -->
      <div class="btn" @click="handleDownload()">下载传讯及全部附件</div>
    </div>
    <!-- 水印 -->
    <!-- <watermark :content="watermarkContent" :gapX="20" :gapY="20" :opacity="0.2"></watermark> -->
    <!-- 下载弹框 -->
    <com-confirm v-model="showDownloadPopup" @confirm="handleDownloadConfirm"></com-confirm>
    <!-- 下载成功提示 -->
    <com-toast v-model="showDownloadToast"></com-toast>
    <!-- 审核 -->
    <audit-confirm v-model="showAuditPopup" :type="auditType" @confirm="handleAuditConfirm"></audit-confirm>
    <!-- pdf -->
    <temp-one ref="detail" :data="detail" @success="handleDownloadSuccess"></temp-one>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import ComConfirm from './components/com-confirm.vue'
import ComToast from './components/com-toast.vue'
import AuditConfirm from './components/audit-confirm.vue'
import Watermark from '@/components/base/watermark.vue'
import { checkApi, getContentInfoApi, getManagerInfoApi, sendFileByEmailApi } from '@/api/communication'
import TempOne from './components/temp-one.vue'
import wx from 'weixin-js-sdk'
import { ImagePreview } from "vant";

import wxUtils from "@/utils/wx";
import { cosUploader } from '@/utils/upload'

interface IAttachmentItem {
  /** 附件label名称 */
  name: string;
  /** 附件名称 */
  attachment_name: string;
  /** 附件url */
  url: string;
  /** 附件类型 */
  type: string;
  /** 附件大小 */
  size: number;
}

/** 联系人 */
interface IContactItem {
  /** 部门名称 */
  department_name: string;
  /** 姓名 */
  name: string;
  /** 邮箱 */
  email: string;
  /** 手机号 */
  tel_phone: string;
  /** 备注 */
  remark: string;
}

interface IDetail {
  id: number;
  /** 传讯类型 */
  type: number;
  /** 传讯类型名称 */
  type_str?: string;
  /** 部门id */
  department_name: string;
  /** 标题 */
  title: string;
  /** 主题 */
  theme_title: string;
  /** 至 */
  until_name: string;
  /** 内容 */
  content: string;
  /** 执行事务 */
  perform: string;
  /** 时间 */
  date_time: string;
  /** 联系人 */
  contact_info: IContactItem[];
  /** 附件 */
  attachment: IAttachmentItem[];
  /** 审核人 */
  review_name?: string;
  /** 审核时间 */
  review_time?: number;
  /** 审核留言 */
  comment: string;
  /** 状态 */
  status: number;
  /** 状态名称 */
  status_str?: string;
  /** 页脚设置 */
  page_url?: string;
  /** 导购ID */
  member_id?: number;
  /** 导购名称 */
  member_name?: string;
  /** 提交时间 */
  submit_time?: number;
  /** PDF文件地址 */
  pdf_url: string;
  /** 阅读人数 */
  view_people_num: number;
  /** 阅读次数 */
  view_num: number;
}

@Component({
  components: {
    ComConfirm,
    ComToast,
    AuditConfirm,
    Watermark,
    TempOne
  }
})
export default class CommunicationDetail extends Vue {
  /** 详情 */
  detail: IDetail|null = null
  /** 下载确认 */
  showDownloadPopup = false
  /** 下载城通提示 */
  showDownloadToast = false
  /** 审核确认 */
  showAuditPopup = false
  /** 审核类型 */
  auditType = 0
  /** 要下载的素材 */
  downloadInfo:IAttachmentItem = null
  email = ''

  /** 用户信息 */
  get userInfo(): IUserInfo {
    return this.$store.getters["user/getUserInfo"] || {}
  }

  get watermarkContent() {
    if(!this.userInfo) return ""
    const { name, user_id } = this.userInfo
    return `${name||'导购姓名'}-${user_id||'导购工号'}`
  }

  get formatContent() {
    // const html = this.detail.content||'';
    // 将html中的p标签替换成div标签
    return (value:string) => value
      .replaceAll('<br/><br/>', '<br/>')
      .replaceAll('&nbsp;', '')
      .replaceAll('\n', '')
      // 将所有的图片标签使用p标签包裹
      .replace(/(<img[^>]*>)/g, ($1) => {
        return `</p><p>${$1}</p><p>`
      })
      .replaceAll('<p><br/></p>', '')
      .replaceAll('<li>', '<li><p>')
      .replaceAll('</li>', '</p></li>')

  }

  async created() {
    if(this.userInfo.app_source !== 2) {
      await this.getAuth()
    }
    await this.getDetail()
    
    if(this.$route.query.from !== '1') {
      // 当前需要直接退出的页面
      window.history.pushState("_android_back", "", location.href);
      window.addEventListener("popstate", this.popstate, false);
      // 神策数据上报 -- 进入传讯
      this.$track.trackEnterTransmittal({
        transmittal_entrance: 2
      })
    }

    // 神策数据上报 -- 浏览传讯内容
    this.$track.trackViewTransmittal({
      transmittal_title: this.detail.title,
      transmittal_publisher: this.detail.member_name,
      release_time: this.detail.submit_time,
      transmittal_entrance: this.$route.query.from === '1' ? this.$route.params.type:4
    })

    try {
      wxUtils.wxAgentConfig();
    } catch (error) {
      console.log(error);
    }
    
  }
  destroyed() {
    if(this.$route.query.from !== '1') {
      // 记得清除!!!!不然单页面，其它页面也可以用
      window.removeEventListener("popstate", this.popstate, false);
    }
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

  /**
   * 格式化时间
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
      const {code, message, data}: IApiResult = await getContentInfoApi({
        id: Number(this.$route.query.id)
      })
      if (code!== 200) {
        this.$toast.fail(message);
        return;
      }
      this.detail = data
      this.$nextTick(() => {
        // 获取所有的img
        const imgs = document.querySelectorAll(".main img")
        imgs.forEach(item => {
          item.addEventListener("click", this.handleImageClick)
        })
      })
    } catch (error) {
      console.log(error);
    }
  }

  // 给图片添加监听事件
  handleImageClick(e) {
    const { src } = e.target
    console.log(src);
    
    // this.$toast.loading({
    //   message: '图片加载中',
    //   duration: 0,
    //   forbidClick: true,
    //   loadingType: 'spinner'
    // })
    ImagePreview({
      images: [src],
      startPosition: 0,
    });

    // window.wx.invoke("previewImage", {
    //   current: src,
    //   urls: [src]
    // }, (res) => {
    //   console.log(res);
      
    //   this.$toast.clear()
    // })
  }

  /**
   * 预览
   */
  async handlePreview(item: IAttachmentItem) {
    this.$toast.loading({
      message: '文件加载中',
      duration: 0,
      forbidClick: true,
      loadingType: 'spinner'
    })
    try {
      const url = await cosUploader.getObjectUrl(item.url)
      window.wx.invoke("previewFile", {
        url: url,
        name: item.attachment_name.replaceAll('%', '％'),
        size: item.size
      }, (res) => {
        console.log(res);
        this.$toast.clear()
        if(res.err_msg !== 'previewFile:ok') {
          this.$toast.fail('文件打开失败')
        }
      });
    } catch (error) {
      console.error('获取文件预览地址失败', error)
      this.$toast.clear()
      this.$toast.fail('文件打开失败')
    }
  }

  /**
   * 下载
   */
  handleDownload(item?:IAttachmentItem) {
    this.showDownloadPopup = true
    this.downloadInfo = item
    if(item) {
      // 神策数据上报 -- 点击下载传讯文件
      this.$track.trackClickTransmittalDownload({
        transmittal_type: this.detail.type
      })
    } else {
      // 神策数据上报 -- 点击下载传讯及所有附件
      this.$track.trackClickAllAttachmentDownload({
        transmittal_type: this.detail.type
      })
    }
  }
  /**
   * 下载确认
   */
  async handleDownloadConfirm(value:string) {
    this.email = value;
    if(this.downloadInfo) {
      this.$toast.loading({
        message: '正在提交',
        duration: 0,
        forbidClick: true,
        loadingType: 'spinner'
      })
      this.handleDownloadSuccess(this.downloadInfo)
    } else {

      // (this.$refs.detail as any).handleDownloadPdf()
      if(this.detail.pdf_url) {
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
      } else {
        (this.$refs.detail as any).handleDownloadPdf()
      }
    }
  }

  async handleDownloadSuccess(e) {
    try {
      const params:any = {
        id: this.detail.id,
        email: this.email,
        file_url: e.url,
        rename: e.rename||e.attachment_name,
        attachment: e.attachment_name?1:0
      }
      // debugger
      const {code, message, data}: IApiResult = await sendFileByEmailApi(params)
      if (code!== 200) {
        this.$toast.fail(message);
        return;
      }
      this.showDownloadToast = true
      this.showDownloadPopup = false
      this.$toast.clear()
      // 神策数据上报 -- 下载传讯附件
      this.$track.trackDownloadTransmittal({
        transmittal_title: this.detail.title,
        transmittal_publisher: this.detail.member_name,
        release_time: this.detail.submit_time,
      })
      
    } catch (error) {
      console.log(error);
    }
  }

  /**
   * 取消
   */
  handleCancel() {
    this.$router.go(-1)
  }

  /**
   * 审核
   */
  handleAudit(type: number) {
    this.auditType = type
    this.showAuditPopup = true
  }
  /**
   * 审核确认
   */
  async handleAuditConfirm(e) {
    try {
      const {code, message, data}: IApiResult = await checkApi({
        id: Number(this.$route.query.id),
        status: [3, 4][e.type],
        comment: e.value
      })
      if (code!== 200) {
        this.$toast.fail(message);
        return;
      }
      this.showAuditPopup = false
      this.$toast.success('审核完成')
      // 神策数据上报 -- 审核传讯
      this.$track.trackAuditTransmittal({
        transmittal_title: this.detail.title,
        transmittal_publisher: this.detail.member_name,
        release_time: this.detail.submit_time,
        audit_type: e.type
      })
      this.getDetail()
    } catch (error) {
      console.log(error);
    }
  }
}

</script>
<style lang="scss" scoped>
::v-deep ol {
  list-style-type: decimal;
  li {
    margin-left: 2em;
  }
}
::v-deep ul {
  list-style-type: disc;
  li {
    margin-left: 2em;
  }
}
.detail {
  min-height: 100vh;
  font-size: 14px;
  background-color: #f7f7f7;
  overflow: hidden;
  box-sizing: border-box;
  padding-bottom: calc(56px + contain(safe-area-inset-bottom));
  padding-bottom: calc(56px + env(safe-area-inset-bottom));


  .header {
    width: 345px;
    margin: auto;
    background: rgba(166, 166, 166, 0.1);
    margin-top: 15px;
    position: relative;

    .audit {
      padding: 10px 15px;
      border-bottom: 1px dashed rgba(166, 166, 166, 0.1);

      .auditor {
        font-weight: 500;
        font-size: 12px;
        color: #000000;
        line-height: 17px;
      }

      .audit-time {
        font-weight: 400;
        font-size: 10px;
        color: #777777;
        line-height: 14px;
        margin-top: 5px;
      }

      .audit-remark {
        width: 315px;
        margin-top: 10px;
        background-color: rgba(166, 166, 166, 0.4);
        padding: 10px;
        display: flex;

        .audit-remark-label {
          width: 33px;
          font-weight: 500;
          font-size: 11px;
          color: #000000;
          line-height: 16px;
        }

        .audit-remark-value {
          flex: 1;
          font-weight: 400;
          font-size: 10px;
          color: #000000;
          line-height: 14px;
          white-space: pre-wrap;
          word-break: break-all;
          overflow: hidden;
        }
      }
    }

    .header-inner {
      padding: 10px 15px;
    }

    .title {
      height: 20px;
      font-weight: 400;
      font-size: 14px;
      color: #000000;
      line-height: 20px;

      .type {
        color: #FF8F52;
      }
    }

    .time {
      font-weight: 300;
      font-size: 12px;
      color: #777777;
      line-height: 17px;
      margin-top: 5px;
    }

    .submitter {
      font-weight: 500;
      font-size: 12px;
      color: #777777;
      line-height: 17px;
      margin-left: 5px;
    }

    .status {
      width: 40px;
      height: 21px;
      box-shadow: 0px 1px 3px 0px rgba(0,0,0,0.1);
      border-radius: 100px 0px 0px 100px;
      position: absolute;
      top: 10px;
      right: -7px;
      font-size: 10px;
      color: #000;
      display: flex;
      align-items: center;
      justify-content: center;

      &::after {
        content: '';
        width: 7px;
        height: 6px;
        position: absolute;
        bottom: -6px;
        right: 0;
      }

      &[data-status="0"] {
        background: #A6A6A6;
        &::after {
          background: linear-gradient(to right bottom, #787777, #787777 50%, transparent 51%, transparent);
        }
      }
      &[data-status="1"] {
        background: #C7AD7F;
        &::after {
          background: linear-gradient(to right bottom, #8E713E, #8E713E 50%, transparent 51%, transparent);
        }
      }
      &[data-status="3"] {
        background: #E6B8B8;

        &::after {
          background: linear-gradient(to right bottom, #B27575, #B27575 50%, transparent 51%, transparent);
        }
      }
      &[data-status="4"] {
        background: #BBD1A7;

        &::after {
          background: linear-gradient(to right bottom, #7A9F59, #7A9F59 50%, transparent 51%, transparent);
        }
      }
      &[data-status="2"] {
        background: #E1E1E1;

        &::after {
          background: linear-gradient(to right bottom, #C0C0C0, #C0C0C0 50%, transparent 51%, transparent);
        }
      }
      &[data-status="5"] {
        background: #E1E1E1;
        opacity: 0.5;

        &::after {
          background: linear-gradient(to right bottom, #C0C0C0, #C0C0C0 50%, transparent 51%, transparent);
        }
      }
    }
  }

  .content {
    padding: 15px 20px;
    position: relative;
    user-select: auto;

    .bottom {
      position: relative;
      padding-bottom: 20px;
      font-size: 11px;
      .line {
        width: 345px;
        border-bottom: 1px dashed #A6A6A6;
        opacity: 0.4;
        position: absolute;
        bottom: 0;
        left: 50%;
        transform: translateX(-50%);
      }
    }

    .theme-title {
      font-weight: 600;
      font-size: 16px;
      color: #3875C7;
      line-height: 22px;
      margin-bottom: 10px
    }

    .title {
      font-weight: 600;
      font-size: 16px;
      color: #000000;
      line-height: 22px;
    }

    .until-name {
      font-weight: 400;
      font-size: 14px;
      color: #000000;
      line-height: 20px;
      text-align: left;
      font-style: normal;
    }

    .main {
      font-weight: 400;
      font-size: 14px;
      color: #000000;
      line-height: 26px;
      text-align: justify;
      word-break: break-all;
      white-space: pre-wrap;
      ::v-deep img {
        max-width: 100% !important;
      }
      ::v-deep table {
        border: 1px solid #000;
        width: 100%;
      }
      ::v-deep td,
      ::v-deep th {
        &:not(:last-child) {
          border-right: 1px solid #000;
        }
      }

      ::v-deep tr:not(:last-child) td,
      ::v-deep tr:not(:last-child) th {
        border-bottom: 1px solid #000;
      }
    }

    .row {
      display: flex;
      margin-top: 10px;
      font-size: 14px;

      &-label {
        width: 70px;
        height: 20px;
        line-height: 20px;
        color: #A6A6A6;

        &.width-auto {
          width: auto;
          min-width: 40px;
        }
      }

      &-value {
        flex: 1;
        overflow: hidden;
        line-height: 20px;
        word-break: break-all;
        text-align: justify;
      }
    }

    .row-text {
      height: 20px;
      display: flex;
      align-items: center;
      margin-top: 10px;
      color: #A6A6A6;
      justify-content: space-between;

      &::after {
        content: "";
        height: 20px;
        width: 250px;
        border-bottom: 1px solid #A6A6A6;
        opacity: 0.4;
      }
    }

    .attach {
      padding: 10px 0;
      .attach-title {
        height: 20px;
        font-weight: 600;
        font-size: 14px;
        color: #A6A6A6;
        line-height: 20px;
        margin-left: -5px;
      }

      .attach-list {
        margin-top: 20px;
      }

      .attach-item {
        width: 335px;
        height: 68px;
        background: #FFFFFF;
        border-radius: 10px;
        padding: 12px 15px;
        display: flex;
        align-items: center;
        position: relative;
        margin-bottom: 20px;

        &-label {
          position: absolute;
          left: 10px;
          top: 0;
          transform: translateY(-50%);
          height: 15px;
          line-height: 15px;
          height: 15px;
          background: #D8D8D8;
          border-radius: 3px;
          padding: 0 5px;
          font-size: 10px;
        }

        &-type {
          width: 39px;
          height: 44px;
          margin-right: 20px;
          
          &[data-type="img"] {
            background: url('~@/assets/communication/attach/image.png') center/100% no-repeat;
          }
          &[data-type="word"] {
            background: url('~@/assets/communication/attach/word.png') center/100% no-repeat;
          }
          &[data-type="excel"] {
            background: url('~@/assets/communication/attach/excel.png') center/100% no-repeat;
          }
          &[data-type="ppt"] {
            background: url('~@/assets/communication/attach/ppt.png') center/100% no-repeat;
          }
          &[data-type="pdf"] {
            background: url('~@/assets/communication/attach/pdf.png') center/100% no-repeat;
          }
          &[data-type="video"] {
            background: url('~@/assets/communication/attach/video.png') center/100% no-repeat;
          }
        }

        &-name {
          flex: 1;
          font-weight: 400;
          font-size: 10px;
          color: #000000;
          line-height: 14px;
          text-align: left;
          font-style: normal;
          overflow: hidden;
          // 两行隐藏
          display: -webkit-box;
          -webkit-line-clamp: 2;
          -webkit-box-orient: vertical;
          text-overflow: ellipsis;
        }

        &-btn {
          margin-left: 20px;
          font-size: 12px;
          display: flex;
          align-items: center;
          justify-content: center;
          flex-direction: column;
          color: #3875C7;

          .icon-download {
            width: 16px;
            height: 16px;
            background: url('~@/assets/communication/icon_save.png') center/100% no-repeat;
          }
          .icon-preview {
            width: 16px;
            height: 16px;
            background: url('~@/assets/communication/icon_preview.png') center/100% no-repeat;
          }

          .text {
            margin-top: 5px;
            font-weight: 400;
            font-size: 12px;
            line-height: 17px;
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
</style>
