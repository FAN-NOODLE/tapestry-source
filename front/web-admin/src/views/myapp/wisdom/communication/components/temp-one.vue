<template>
  <!-- pdf排版 -->
  <div ref="pdf" class="pdf" :style="{width: '690px'}">
    <img class="center" style="width: 139px;" src="https://cdn-cos.coach.com.cn/uploads/admin/pdf_1.png" alt="">
    <table data-pdfmake="{'layout': 'noBorders'}"  style="margin-top: 20px;">
      <tr v-if="data.type===3 && data.theme_title">
        <td colspan="2">
          <div class="fs16 red center bold">{{data.theme_title}}</div>
        </td>
      </tr>
      <tr>
        <td width="80px">
          <div class="fs10 lh12">标题:</div>
        </td>
        <td>
          <div class="fs10 lh12 blue bold">{{ data.title }}</div>
        </td>
      </tr>
      <tr>
        <td width="80px">
          <div class="fs10 lh12">至:</div>
        </td>
        <td>
          <div class="fs10 lh12">{{ data.until_name }}</div>
        </td>
      </tr>
      <tr>
        <td width="80px">
          <div class="fs10 lh12">内容:</div>
        </td>
        <td>
          <div class="main fs10 lh12" v-html="formatContent(data.content)"></div>
        </td>
      </tr>
      <tr v-if="data.perform">
        <td width="80px">
          <div class="fs10 lh12">执行事务:</div>
        </td>
        <td>
          <div class="main fs10 lh12" v-html="formatContent(data.perform)"></div>
        </td>
      </tr>
      <tr v-if="data.date_time">
        <td width="80px">
          <div class="fs10 lh12">日期:</div>
        </td>
        <td>
          <div class="fs10 lh12">{{ data.date_time }}</div>
        </td>
      </tr>
      <tr v-if="data.contact_info && data.contact_info.length">
        <td width="80px">
          <div class="fs10 lh12">联系人:</div>
        </td>
        <td>
          <div class="fs10 lh12">
            <div v-for="(item, index) in data.contact_info" :key="index">
              <div v-if="item.department_name">{{ item.department_name }}</div>
              <div>
                <span>{{ item.name }}</span>
                <span v-if="item.tel_phone"> / </span>
                <span v-if="item.tel_phone">{{item.tel_phone}}</span>
                <span v-if="item.email"> / </span>
                <a class="blue" v-if="item.email" :href="`email:${item.email}`">{{item.email}}</a>
                <span v-if="item.remark">（{{ item.remark }}）</span>
              </div>
            </div>
          </div>
        </td>
      </tr>
      <tr v-if="data.attachment.length">
        <td width="80px">
          <div class="fs10 lh12">附件:</div>
        </td>
        <td>
          <div class="fs10 lh12">
            <div v-for="(item, index) in data.attachment" :key="index">
              <span>{{ item.name }} {{ index+1 }} - </span>
              <span>{{ item.attachment_name }}</span>
            </div>
          </div>
        </td>
      </tr>
    </table>
    <!-- 周传讯底部 -->
    <img v-if="data.type==1 && data.page_url!=='2'" style="width: 690px;margin-top: 10px;" src="https://cdn-cos.coach.com.cn/uploads/admin/pdf_4.png" alt="" />
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from 'vue-property-decorator'
import { IFormInfo } from './edit-form.vue'
import UserModule from "@/store/modules/user";
import DownloadModule, {Task} from "@/store/modules/download";
import { SavePdfFromHtml } from '@/utils/download/pdf'
import styles from '../styles';
import Api from '@/api/wisdom/communication';
import { downloadFile } from '@/utils/download';
import dayjs from 'dayjs'
import { cosUploader } from '@/utils/upload';

@Component
export default class TempOne extends Vue {
  /** 表单数据 */
  @Prop({ type: Object, default: () => (null) }) data: IFormInfo;

  get formatContent() {
    // const html = this.data.content||'';
    
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
      .replaceAll('<p>', '<div>')
      .replaceAll('<p ', '<div ')
      .replaceAll('</p>', '</div>')
      .replaceAll('<table> </table>', '')
      .replaceAll('<table></table>', '')
  }

  get filename() {
    return (useTitle=false) => {
      if(!this.data) return ''
      const { id, title } = this.data
      if(useTitle) {
        return `${title}_${dayjs().format('YYYYMMDD')}.pdf`
      }
      
      return `${title}_${id||0}.pdf`
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

  /**
   * 下载pdf文件
   */
  handleDownloadPdf() {
    const element = this.$refs.pdf as HTMLElement
    const that = this
    const task = new Task({
      title: '传讯内容-下载文件',
      description: '下载传讯PDF',
    }).onRun(function() {
      return new Promise((resolve) => {
        // 一秒更新一次任务进度
        const pdfLoader = new SavePdfFromHtml({
          element,
          filename: that.filename(),
          styles: styles as any,
          info: that.getPdfInfo
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
   * @param value 
   */
  async handleDownLoadSuccess(e) {
    // 下载
    cosUploader.getObjectUrl(e.url).then(res => {
      downloadFile(res, {
        filename: this.filename(true)
      })
    })
    const res:any = await Api.saveMergePdfApi({
      info_id: this.data.id,
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

  handlePreview() {
    const element = this.$refs.pdf as HTMLElement
    // 一秒更新一次任务进度
    const pdfLoader = new SavePdfFromHtml({
      element,
      filename: this.filename(),
      styles: styles as any,
      info: this.getPdfInfo
    })
    pdfLoader.preview()
  }

  handleCreate() {
    const element = this.$refs.pdf as HTMLElement
    // 一秒更新一次任务进度
    const pdfLoader = new SavePdfFromHtml({
      element,
      filename: this.filename(),
      styles: styles as any,
      info: this.getPdfInfo
    })
    return pdfLoader.create()
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
  font-family: STSongti;
  font-size: 12px;
  line-height: 1.5;
  position: fixed;
  top: 100px;
  right: -100%;
  background-color: #fff;
  .logo {
    text-align: center;
  }
  .icon {
    text-align: center;
    overflow: hidden;
    img {
      display: inline;
    }
  }
  .title {
    text-align: center;
    font-size: 16px;
    color: #FF0D0D;
    margin-top: 20px;
  }

  .table {
    font-size: 10px;
    .tr {
      line-height: 2;
    }
    .td:nth-child(1) {
      width: 80px;
    }
  
    .cell {
      line-height: 2;
    }
  }

}
</style>