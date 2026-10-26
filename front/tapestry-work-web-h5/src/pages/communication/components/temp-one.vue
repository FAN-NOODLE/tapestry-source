<template>
  <!-- pdf排版 -->
  <div ref="pdf" class="pdf" :style="{width: '690px'}">
    <img class="center" width="139px" height="52px" src="https://cdn-cos.coach.com.cn/uploads/admin/pdf_1.png" alt="">
    <table data-pdfmake="{'layout': 'noBorders'}" style="margin-top: 20px;">
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
          <pre class="main fs10 lh12" v-html="formatContent(data.content)"></pre>
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
    <img v-if="data.type==1&&data.page_url !== '2'" style="width: 690px;margin-top: 10px;" src="https://cdn-cos.coach.com.cn/uploads/admin/pdf_4.png" alt="">
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from 'vue-property-decorator'
import { SavePdfFromHtml } from '@/utils/download/pdf'
import styles from '../styles'
import { saveMergePdfApi } from '@/api/communication';


@Component
export default class TempOne extends Vue {
  /** 表单数据 */
  @Prop({ type: Object, default: () => (null) }) data: any;

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
  }

  get userInfo(): IUserInfo {
    return this.$store.getters["user/getUserInfo"]
  }
  // 获取pdf info
  get getPdfInfo() {
    const { name:account_name } = this.userInfo
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
      const { title, id } = this.data
      if(useTitle) {
        return `${title.slice(0, 80)}.pdf`
      }
      return `${title}_${id}.pdf`
    }
  }

  /**
   * 下载pdf文件
   */
  handleDownloadPdf() {
    const element = this.$refs.pdf as HTMLElement
    const toast = this.$toast.loading({
      message: '生成中(0%)',
      duration: 0,
      forbidClick: true,
      loadingType: 'spinner'
    })
    // 一秒更新一次任务进度
    const pdfLoader = new SavePdfFromHtml({
      element,
      filename: this.filename(),
      styles: styles as any,
      info: this.getPdfInfo
    })
    pdfLoader.progress(({percent}) => {
      toast.message = `生成中(${(percent*100).toFixed(0)}%)`
    }).create().then(res=> {
      toast.message = '正在提交'
      this.handleDownloadSuccess(res);
    }).catch(err => {
      this.$toast.fail('生成失败')
    })
  }
    /**
   * 下载成功
   * @param value 
   */
  async handleDownloadSuccess(e) {
    this.$emit('success', {
      ...e,
      rename: this.filename(true)
    });
    saveMergePdfApi({
      info_id: this.data.id,
      pdf_url: e.url
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
  white-space: pre-wrap;
  text-align: justify;
  line-height: 1.2;
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

  ::v-deep div {
    font-family: STSongti !important;
    font-size: 12px !important;
  }
}
.pdf {
  position: fixed;
  top: 100px;
  right: -10000px;
  background-color: #fff;

}
</style>