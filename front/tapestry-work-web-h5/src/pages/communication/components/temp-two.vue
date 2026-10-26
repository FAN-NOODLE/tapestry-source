<template>
  <!-- pdf排版 -->
  <div v-if="data" ref="pdf" class="pdf" :style="{width: 680+'px'}">
    <img class="right" width="139px" src="https://cdn-cos.coach.com.cn/uploads/admin/pdf_1.png" alt="">
    <div class="fs24 red lh1 right bold">每周传讯</div>
    <div class="fs14 lh1 right">{{ getWeek }}</div>
    <div class="fs14 lh1 right">{{ formatTime(data.send_time)}}</div>
    <!-- 目录 -->
    <div class="directory">
      <div class="red fs13 lh12 bold">目录</div>
      <div class="directory-content">
        <table data-pdfmake="{'layout': 'noBorders'}">
          <tr>
            <td style="width: 50%;">
              <div class="directory-column">
                <div class="part" v-for="(item, index) in column(1)" :key="index">
                  <div class="fs13 underline bold">{{item.department_name}}</div>
                  <ul v-for="(info, infoIndex) in item.info_list" :key="infoIndex">
                    <li class="fs10">{{info.title}}</li>
                  </ul>
                </div>
              </div>
            </td>
            <td style="width: 50%;">
              <div class="directory-column">
                <div class="part" v-for="(item, index) in column(2)" :key="index">
                  <div class="fs13 underline bold">{{item.department_name}}</div>
                  <ul v-for="(info, infoIndex) in item.info_list" :key="infoIndex">
                    <li class="fs10">{{info.title}}</li>
                  </ul>
                </div>
              </div>
            </td>
          </tr>
        </table>
      </div>
      <img style="margin-top: 10px;width: 690px;" src="https://cdn-cos.coach.com.cn/uploads/admin/pdf_3.png" alt="">
    </div>
    <!-- 传讯内容 -->
    <div class="content">
      <!-- 部门组 :data-pdfmake="JSON.stringify(index<data.content_list.length-1?{pageBreak: 'after'}:{})" -->
      <div v-for="(item, index) in data.content_list" :key="index" :data-pdfmake="JSON.stringify(index<data.content_list.length-1?{pageBreak: 'after'}:{})">
        <div class="fs13 red lh11" style="margin-top: 10px;">{{item.department_name}}</div>
        <!-- 传讯列表 -->
        <div v-for="(info, infoIndex) in item.info_list" :key="infoIndex" class="content-item">
          <!-- 传讯内容 -->
          <table 
            :data-pdfmake="JSON.stringify({layout: 'noBorders'})" 
            :style="{
              width: '690px',
              marginTop: infoIndex == 0 ? '0' : '20px'
            }">
            <!-- 紧急传讯主题 -->
            <tr v-if="info.type==3 && info.until_name">
              <td colspan="2">
                <div class="fs16 lh12 red center bold">{{info.until_name}}</div>
              </td>
            </tr>
            <tr>
              <td width="80px">
                <div class="fs10 lh12">标题:</div>
              </td>
              <td>
                <div class="fs10 lh12 blue bold">{{ info.title }}</div>
              </td>
            </tr>
            <tr>
              <td width="80px">
                <div class="fs10 lh12">至:</div>
              </td>
              <td>
                <div class="fs10 lh12">{{ info.until_name }}</div>
              </td>
            </tr>
            <tr>
              <td width="80px">
                <div class="fs10 lh12">内容:</div>
              </td>
              <td>
                <div class="main fs10 lh12" v-html="formatContent(info.content)"></div>
              </td>
            </tr>
            <tr v-if="info.perform">
              <td width="80px">
                <div class="fs10 lh12">执行事务:</div>
              </td>
              <td>
                <div class="main fs10 lh12" v-html="formatContent(info.perform)"></div>
              </td>
            </tr>
            <tr v-if="info.date_time">
              <td width="80px">
                <div class="fs10 lh12">日期:</div>
              </td>
              <td>
                <div class="fs10 lh12">{{ info.date_time }}</div>
              </td>
            </tr>
            <tr v-if="info.contact_info && info.contact_info.length">
              <td width="80px">
                <div class="fs10 lh12">联系人:</div>
              </td>
              <td>
                <div class="fs10 lh12">
                  <div v-for="(contact, index) in info.contact_info" :key="index">
                    <div v-if="contact.department_name">{{ contact.department_name }}</div>
                    <div>
                      <span>{{ contact.name }}</span>
                      <span v-if="contact.tel_phone"> / </span>
                      <span v-if="contact.tel_phone">{{contact.tel_phone}}</span>
                      <span v-if="contact.email"> / </span>
                      <a v-if="contact.email" :href="`email:${contact.email}`">{{contact.email}}</a>
                      <span v-if="contact.remark">（{{ contact.remark }}）</span>
                    </div>
                  </div>
                </div>
              </td>
            </tr>
            <tr v-if="info.attachment.length">
              <td width="80px">
                <div class="fs10 lh12">附件:</div>
              </td>
              <td>
                <div class="fs10 lh12">
                  <div v-for="(attach, index) in info.attachment" :key="index">
                    <span>{{ attach.name }} {{ index+1 }} - </span>
                    <span>{{ attach.attachment_name }}</span>
                  </div>
                </div>
              </td>
            </tr>
          </table>
          <!-- 部门传讯的最后一个 -->
          <img v-if="info.type === 1&&info.page_url !== '2' && infoIndex == item.info_list.length -1" class="center" style="margin-top: 5px;width:690px;" src="https://cdn-cos.coach.com.cn/uploads/admin/pdf_4.png" alt="">
          <!-- 周传讯 -->
          <img v-else-if="info.type === 1&&info.page_url !== '2'" class="center" style="margin-top: 5px;width:690px;" src="https://cdn-cos.coach.com.cn/uploads/admin/pdf_5.png" alt="">
          <!-- 最后一个 -->
          <img v-else-if="infoIndex == item.info_list.length -1" class="center" style="margin-top: 5px;width:690px;" src="https://cdn-cos.coach.com.cn/uploads/admin/pdf_3.png" alt="">
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from 'vue-property-decorator'

import { SavePdfFromHtml } from '@/utils/download/pdf'
import styles from '../styles'
import { downloadFile } from '@/utils/download';
import { getImageData } from '@/utils/download/util';
import { saveMergePdfApi } from '@/api/communication';

@Component
export default class TempTwo extends Vue {
  /** 表单数据 */
  @Prop({ type: Object, default: () => (null) }) data: any;

  /** 用户信息 */
  get userInfo(): IUserInfo {
    return this.$store.getters["user/getUserInfo"]
  }

  get formatContent() {
    // 将html中的p标签替换成div标签
    return (html:string) => html
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
      .replaceAll('<p>', '<div>')
      .replaceAll('<p ', '<div ')
      .replaceAll('</p>', '</div>')
  }

  get getWeek() {
    const year = new Date().getFullYear().toString().slice(2)
    return `FY${year} WK${this.data.week_num}`
  }

  get column() {
    if(!this.data || !this.data.content_list) return () => []
    return (type) => {
      const len = this.data.content_list.length
      if(type == 1) {
        return this.data.content_list.slice(0, Math.round(len/2))
      } else {
        return this.data.content_list.slice(Math.round(len/2), len)
      }
    }
  }
  // 获取pdf info
  get getPdfInfo() {
    const { name } = this.userInfo
    const { title, member_name } = this.data
    return {
      title: title,
      author: name,
      subject: title,
      keywords: title,
      creator: member_name,
      producer: name,
      creationDate: new Date(),
      modDate: new Date(),
    }
  }


  get filename() {
    return (useTitle=false) => {
      if(!this.data) return ''
      const { week_num, id } = this.data
      if(useTitle) {
        return `中国区每周传讯-第${week_num}周.pdf`
      }
      return `中国区每周传讯-第${week_num}周_${id}.pdf`
    }
  }


  /**
   * 格式化时间
   * @param time 
   */
  formatTime(time: number) {
    return this.$helper.dateFormat(time, "yyyy年MM月dd日");
  }

  /**
   * 下载pdf文件
   */
  async handleDownloadPdf() {
    const element = this.$refs.pdf as HTMLElement
    const that = this
    const { week_num } = this.data
    const toast = this.$toast.loading({
      message: '生成中(0%)',
      duration: 0,
      forbidClick: true,
      loadingType: 'spinner'
    })
    const { base64 } = await getImageData('https://cdn-cos.coach.com.cn/uploads/admin/pdf_1.png')
    // 一秒更新一次任务进度
    const pdfLoader = new SavePdfFromHtml({
      element,
      filename: this.filename(),
      styles: styles as any,
      info: that.getPdfInfo,
      footer: function(currentPage, pageSize) {
        return [
          {
              text: `第${week_num}周 - ` + currentPage.toString() + '/' + pageSize.toString(),
              alignment: 'center',
              margin: [0, 10, 0, 0],
              fontSize: 10,
              font: 'SongTiSC'
          },
        ] as any
      },
      background: function(currentPage, pageSize) {
        return [
          {
            image: base64, // 图像路径
            alignment: 'right', // 右对齐
            margin: [0, 800, 20, 0], // 右边距为20单位
            width: 70, // 宽度为100单位
          }
        ] as any
      }
    })
    pdfLoader.progress(({percent}) => {
      toast.message = `生成中(${(percent*100).toFixed(0)}%)`
    }).create().then(res=> {
      toast.message = '正在提交'
      this.handleDownLoadSuccess(res);
    }).catch(err=> {
      console.log(err);
      this.$toast.fail('下载失败')
    })
  }

  /**
   * 下载成功
   */
  async handleDownLoadSuccess(e) {
    await saveMergePdfApi({
      merge_id: this.data.id,
      pdf_url: e.url
    })
    this.$emit('success', {
      ...e,
      rename: this.filename(true)
    })
  }

}
</script>

<style lang="scss" scoped>
@font-face {
  font-family: STSongti;
  src: url('https://cdn-cos.coach.com.cn/uploads/fonts/STSongti-SC-Regular.ttf') format('truetype');
  font-weight: normal;
}
.main {
  font-family: STSongti;
  font-size: 12px;
  line-height: 1.2;
  white-space: pre-wrap;
  text-align: justify;
  ::v-deep table {
    border: 1px solid #000;
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
.pdf {
  position: fixed;
  top: 100px;
  right: -1000px;
}
</style>