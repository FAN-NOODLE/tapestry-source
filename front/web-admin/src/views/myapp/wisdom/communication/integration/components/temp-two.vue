<template>
  <!-- pdf排版 -->
  <div v-if="data" ref="pdf" class="pdf" :style="{width: 680+'px'}">
    <img class="right" width="139px" src="https://cdn-cos.coach.com.cn/uploads/admin/pdf_1.png" alt="">
    <div class="fs24 red lh12 right bold">每周传讯</div>
    <div class="fs14 lh12 right">{{ getWeek }}</div>
    <div class="fs14 lh12 right">{{ (data.send_time||Date.now()/1000)|formatDate('yyyy年mm月dd日')}}</div>
    <!-- 目录 -->
    <div class="directory page">
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
      <img style="margin-top: 10px;width:690px;" src="https://cdn-cos.coach.com.cn/uploads/admin/pdf_3.png" alt="">
    </div>
    <!-- 传讯内容 -->
    <div class="content">
      <!-- 部门组
        :data-pdfmake="JSON.stringify(index<data.content_list.length-1?{pageBreak: 'after'}:{})"
       -->
      <div 
        v-for="(item, index) in data.content_list" 
        :key="index">
        <div class="fs13 red lh11" style="margin-top: 10px;">{{item.department_name}}</div>
        <!-- 传讯列表 -->
        <div v-for="(info, infoIndex) in item.info_list" :key="infoIndex" class="content-item page">
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
          <img v-if="info.type === 1 && info.page_url!=='2' && infoIndex == item.info_list.length -1" class="center" style="margin-top: 5px;width:690px;" src="https://cdn-cos.coach.com.cn/uploads/admin/pdf_4.png" alt="">
          <!-- 周传讯 -->
          <img v-else-if="info.type === 1 && info.page_url!=='2'" class="center" style="margin-top: 5px;width:690px;" src="https://cdn-cos.coach.com.cn/uploads/admin/pdf_5.png" alt="">
          <!-- 最后一个 -->
          <img v-else-if="infoIndex == item.info_list.length -1" class="center" style="margin-top: 5px;width:690px;" src="https://cdn-cos.coach.com.cn/uploads/admin/pdf_3.png" alt="">
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop, Watch } from 'vue-property-decorator'
import { IFormInfo } from './edit-form.vue'
import UserModule from "@/store/modules/user";
import DownloadModule, {Task} from "@/store/modules/download";

import { SavePdfFromHtml } from '@/utils/download/pdf'

import styles from '../../styles'
import { getImageData } from '@/utils/download/util';
import Api from '@/api/wisdom/communication';
import { downloadFile } from '@/utils/download';
import { cosUploader } from '@/utils/upload';

@Component
export default class TempTwo extends Vue {
  /** 表单数据 */
  @Prop({ type: Object, default: () => (null) }) value: IFormInfo;
  data: IFormInfo = null;

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
      .replaceAll('<table> </table>', '')
      .replaceAll('<table></table>', '')
  }

  get getWeek() {
    // const year = new Date().getFullYear().toString().slice(2)
    
    return `FY26 WK${this.data.week_num}`
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

  get userInfo() {
    return UserModule.getUserInfo
  }
  // 获取pdf info
  get getPdfInfo() {
    const { account_name } = this.userInfo
    const { title, member_name } = this.data
    return {
      title: title,
      author: account_name,
      subject: title,
      keywords: title,
      creator: member_name,
      producer: account_name,
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
      return `中国区每周传讯-第${week_num}周_${id||0}.pdf`
    }
  }

  @Watch('value', {deep: true})
  onDataChange(value) {
    this.data = value
  }

  /**
   * 下载pdf文件
   */
  async handleDownloadPdf() {
    const element = this.$refs.pdf as HTMLElement
    const that = this
    const { id, title, week_num } = this.data
    const { base64 } = await getImageData('https://cdn-cos.coach.com.cn/uploads/admin/pdf_1.png')
    const task = new Task({
      title: '企业传讯-下载PDF文件',
      description: '下载传讯PDF',
    }).onRun(function() {
      return new Promise((resolve) => {
        // 一秒更新一次任务进度
        const pdfLoader = new SavePdfFromHtml({
          element,
          filename: that.filename(),
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
          },
        })
        pdfLoader.progress((value)=> {
          this.percent = value.percent
        }).create().then(res=> {
          that.handleDownLoadSuccess(res);
          resolve(res);
        })
      });
    })
    DownloadModule.addTask(task)
  }

  /**
   * 下载成功
   */
  async handleDownLoadSuccess(e) {
    cosUploader.getObjectUrl(e.url).then(res => {
      downloadFile(res, {
        filename: this.filename(true)
      })
    })
    const res:any = await Api.saveMergePdfApi({
      merge_id: this.data.id,
      pdf_url: e.url
    })
    if (res.code == 200) {
      this.$message.success('保存成功');
    } else {
      this.$message.error(res.message);
    }
  }

  /**
   * 下载文件 pdf
   */
  handleDownloadFile() {
    const { pdf_url } = this.data
    const that = this
    const task = new Task({
      title: '企业传讯-下载PDF文件',
      description: '下载传讯PDF',
    }).onRun(async function() {
      // 下载
      const url = await cosUploader.getObjectUrl(pdf_url)
      const res = await downloadFile(url, {
        filename: that.filename(true),
        onProgress: (e) => {
          this.percent = (e/100).toFixed(1)
        },
      })
      that.$message.success('保存成功');
      return res
    })
    DownloadModule.addTask(task)
  }

  async handlePreview() {
    const element = this.$refs.pdf as HTMLElement
    const that = this
    const { id, title, week_num } = this.data
    const { base64 } = await getImageData('https://cdn-cos.coach.com.cn/uploads/admin/pdf_1.png')
    const task = new Task({
      title: '企业传讯-预览传讯PDF',
      description: '预览传讯PDF',
    }).onRun(function() {
      return new Promise((resolve) => {
        // 一秒更新一次任务进度
        const pdfLoader = new SavePdfFromHtml({
          element,
          filename: that.filename(),
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
          },
        })
        pdfLoader.progress((value)=> {
          this.percent = value.percent
        }).preview().then(res=> {
          resolve(res);
        })
      });
    })
    DownloadModule.addTask(task)
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
  // white-space: pre-wrap;
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
  right: -100%;
}
</style>