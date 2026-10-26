<template>
  <div v-if="data" class="detail-form">
    <!-- 顶部基础信息 -->
    <div class="info">
      <el-row>
        <el-col :span="5">
          <span>{{ data.department_name }} | {{ data.type_name }}</span>
        </el-col>
        <el-col :span="8">
          <span>提交时间：{{ (data.submit_time||0) | formatDate('yyyy/mm/dd hh:mm:ss') }}</span>
        </el-col>
        <el-col :span="5">
          <span>提交者：{{ data.member_name||'--' }}</span>
        </el-col>
        <el-col :span="2" :offset="4">
          <el-tag :type="{
            1: 'warning',
            3: 'danger',
            4: 'success'
          }[data.status]||''">{{ data.status_name }}</el-tag>
        </el-col>
      </el-row>
      <el-row v-if="data.status === 3 || data.status === 4">
        <el-col :span="12">
          <span>审核者：{{ data.review_name||'--' }}</span>
        </el-col>
        <el-col :span="12">
          <span>审核时间：{{ data.review_time | formatDate('yyyy/mm/dd hh:mm:ss') }}</span>
        </el-col>
        <el-col v-if="data.status === 3" :span="24">
          <span class="remark">留言：{{ data.comment||'--' }}</span>
        </el-col>
      </el-row>
    </div>
    <!-- 主题内容排版 -->
    <div>
      <div v-if="data.type===3" class="title">{{ data.theme_title }}</div>
      <div class="middle" >
        <table :border="0" style="width: 100%;border-collapse: collapse;">
          <tr>
            <td style="width: 100px;padding: 10px">标题：</td>
            <td style="padding: 10px">{{ data.title }}</td>
          </tr>
          <tr>
            <td style="width: 100px;padding: 10px">至：</td>
            <td style="padding: 10px">{{ data.until_name }}</td>
          </tr>
          <tr>
            <td style="width: 100px;padding: 10px">内容：</td>
            <td style="padding: 10px"><div class="main" v-html="formatContent(data.content)"></div></td>
          </tr>
          <tr v-if="data.perform">
            <td style="width: 100px;padding: 10px">执行事务：</td>
            <td style="padding: 10px">
              <div class="main" v-html="formatContent(data.perform)"></div>
            </td>
          </tr>
          <tr v-if="data.date_time">
            <td style="width: 100px;padding: 10px">日期：</td>
            <td style="padding: 10px">{{ data.date_time }}</td>
          </tr>
          <tr v-if="data.contact_info && data.contact_info.length">
            <td style="width: 100px;padding: 10px">联系人：</td>
            <td style="padding: 10px">
              <div v-for="(item, index) in data.contact_info" :key="index">
                <div v-if="item.department_name">{{ item.department_name }}</div>
                <div>
                  <span>{{ item.name }}</span>
                  <span v-if="item.tel_phone"> / </span>
                  <span v-if="item.tel_phone" :href="'tel:'+item.tel_phone">{{item.tel_phone}}</span>
                  <span v-if="item.email"> / </span>
                  <a v-if="item.email" :href="'mailto:' + item.email" target="_top">{{ item.email }}</a>
                  <span v-if="item.remark">（{{ item.remark }}）</span>
                </div>
              </div>
            </td>
          </tr>
        </table>
        <!-- 页脚图片 -->
        <div v-if="data.type===1 && data.page_url!=='2'" style="border-top: 1px dashed #ccc;">
          <div class="page-footer">
            <div class="page-footer__item" style="width:30%">
              <span class="page-footer__item__label">负责跟进同事：</span>
              <span class="page-footer__item__value">负责跟进同事</span>
            </div>
            <div class="page-footer__item" style="width:30%">
              <span class="page-footer__item__label">执行行动日期：</span>
              <span class="page-footer__item__value">负责跟进同事</span>
            </div>
            <div class="page-footer__item" style="width:30%">
              <span class="page-footer__item__label">完成行动日期：</span>
              <span class="page-footer__item__value">负责跟进同事</span>
            </div>
          </div>
        </div>
      </div>
    </div>
    <!-- 附件相关 -->
    <div v-if="data.attachment.length" class="attachment">
      <h3 class="m-t-10">相关附件</h3>
      <el-row class="item" v-for="(item, index) in data.attachment" :key="item.name">
        <div class="item-label">{{item.name}}{{index+1}}</div>
        <div class="item-value">
          <div class="item-name">
            <svg-icon class="icon" :icon-class="item.type"></svg-icon>
            <div class="name">{{item.attachment_name}}</div>
            <el-button type="text" @click="handleDownload(item)">下载</el-button>
          </div>
        </div>
      </el-row>
    </div>
    
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from 'vue-property-decorator';
import { IFormInfo } from './edit-form.vue'
import { downloadFile } from '@/utils/download';
import { cosUploader } from '@/utils/upload';

function convertToPixels(mm, dpi) {  
    // 1英寸等于25.4毫米  
    const inchesPerMm = 1 / 25.4;  
    // 将毫米转换为英寸，然后乘以DPI得到像素  
    return Math.round(mm * inchesPerMm * dpi);  
}  
  
function getA4SizeInPixels(dpi) {  
    // A4纸的尺寸（毫米）  
    const widthMm = 210;  
    const heightMm = 297;  
      
    // 转换为像素  
    const widthPx = convertToPixels(widthMm, dpi);  
    const heightPx = convertToPixels(heightMm, dpi);  
      
    return {  
        width: widthPx,  
        height: heightPx  
    };  
}  

@Component
export default class DetailForm extends Vue {
  /** 表单数据 */
  @Prop({ type: Object, default: () => (null) }) data: IFormInfo;

  a4SizeInPixels = getA4SizeInPixels(96)

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

  }
  
  /**
   * 下载附件
   */
  handleDownload(item) {
    cosUploader.getObjectUrl(item.url).then((url) => {
      // 下载
      downloadFile(url, {
        filename: item.attachment_name
      })
    })
  }
}
</script>

<style lang="scss" scoped>
.detail-form {
  width: 1080px;
  max-width: 100%;
  .main {
    // white-space: pre-wrap;
    text-align: justify;
    line-height: 1.5;
    ::v-deep p {
      margin: 10px 0 !important;
      line-height: 1.5 !important;
    }
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
  .info {
    padding: 15px 20px;
    background-color: #f6f7f8;
    margin-bottom: 10px;
    color: #666;
    span {
      line-height: 32px;
      &.remark {
        word-break: break-all;
        white-space: pre-wrap;
        line-height: 18px;
      }
    }
  }
  .title {
    line-height: 30px;
  }
  .middle {
    border: 1px dashed #ccc;
    position: relative;
    td {
      border-bottom: 1px dashed #ccc;
  
      &:not(:last-child) {
        border-right: 1px dashed #ccc;
      }
    }
    tr:last-child td {
      border-bottom: none;
    }
  
    .page-footer {
      height: 40px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 0 20px;
      margin: 20px;
      border-bottom: 2px dotted #000;
  
      .page-footer__item__value {
        color: transparent;
        text-decoration: underline #000;
  
      }
    }
  }
  
  .label {
    width: 200px !important;
  }
  .item {
    display: flex;
    align-items: center;
    margin-top: 10px;
    .item-label {
      width: 80px;
    }
    .item-value {
      flex: 1;
      padding: 10px 15px;
      background-color: #f6f7f8;
      overflow: hidden;
    }
  
    .item-name {
      display: flex;
      align-items: center;
  
      .icon {
        font-size: 20px;
      }
  
      .name {
        flex: 1;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        margin: 0 5px;
      }
    }
  }
}



</style>