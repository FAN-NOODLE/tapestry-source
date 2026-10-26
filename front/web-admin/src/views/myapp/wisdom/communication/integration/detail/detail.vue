<template>
  <layout-form>
    <template slot="content">
      <div class="top">
        <el-image :src="formInfo.cover||defaultCover" class="cover" fit="contain" :preview-src-list="[formInfo.cover]"></el-image>
        <div class="info">
          <div class="title">{{ formInfo.title }}</div>
          <div class="name">{{ formInfo.member_name }}</div>
          <div class="time">发布时间: {{ (formInfo.send_time||0) | formatDate('yyyy/mm/dd hh:mm:ss') }}</div>
          <div class="time" v-if="formInfo.email_loop_url">邮件链接: {{ formInfo.email_loop_url }}</div>
        </div>
        <div class="status">
          <el-button type="success">{{ ['暂未发布','已发布'][formInfo.status] }}</el-button>
        </div>
      </div>
      <div class="line">温馨提示：点击标题可跳转查看传讯内容</div>
      <div class="list">
        <div class="item" v-for="(item, index) in formInfo.content" :key="index">
          <div class="dp-name">{{ item.department_name }}</div>
          <div class="item-inner">
            <p class="m-b-5" v-for="info in item.info_list" :key="info.id" @click="handleLink(info)">
              <a class="link">{{ info.title }}</a>
            </p>
          </div>
        </div>
      </div>
      <!-- 通知人群 -->
      <div class="user">
        <div class="user-title">通知人群:</div>
        <div ref="userList" class="user-list" :style="{maxHeight: isExpand ? 'none' : '100px'}">
          <el-tag class="m-r-10 m-b-10" v-for="item in userList" :key="item.id" size="mini"><i class="el-icon-user-solid m-r-5"></i>{{item.title}}</el-tag>
        </div>
        <div v-if="showMoreBtn" class="more" @click="handleShowMore">{{isExpand?'收起':'展开'}}</div>
      </div>
      <copy-text ref="copyText" :value="formInfo"></copy-text>
    </template>
    <template slot="footer">
      <el-button @click="handleCancel">返回</el-button>
      <el-button v-if="formInfo.type==1" type="primary" @click="handleDownLoad">下载周传讯PDF</el-button>
      <el-button type="primary" :loading="zipLoading" @click="handleDownLoadZip">下载全部附件</el-button>
      <el-button v-if="formInfo.type==1" type="primary" @click="handleCopyText">复制邮件文本</el-button>
      <temp-two ref="detail" :value="formInfo"></temp-two>
    </template>
  </layout-form>
</template>
<script lang="ts">
import { Vue, Component } from 'vue-property-decorator'
import Api from '@/api/wisdom/communication';
import { IFormInfo, defaultFormInfo } from '../components/edit-form.vue'
import CopyText from '../components/copy-text.vue'
import TempTwo from '../components/temp-two.vue'
import { saveFileToZip } from '@/utils/zip';
import { defaultCover } from '../../const';

@Component({
  components: {
    CopyText,
    TempTwo
  }
})
export default class IntegrationDetail extends Vue {
  /** 表单数据 */
  formInfo: IFormInfo = JSON.parse(JSON.stringify(defaultFormInfo))

  zipLoading = false

  defaultCover = defaultCover

  showMoreBtn = false
  isExpand = false
  userList:any[] = []

  mounted() {
    this.getDetail()
    
  }

  handleShowMore() {
    this.isExpand = !this.isExpand
  }

  /**
   * 获取详情
   */
  async getDetail() {
    const res:any = await Api.getMergeDetailApi({
      id: Number(this.$route.query.id)
    })
    if (res.code == 200) {
      res.data.content = res.data.content_list
      
      this.formInfo = res.data

      if(this.formInfo.crowd_type === 1) {
        if(res.data.send_name_list && res.data.send_name_list.length) {
          this.userList = res.data.send_name_list.map(item => ({title: item}))
        } else {
          await this.getWhitelistDetail()
        }
      } else {
        this.userList = res.data.crowd_list
      }
      
      this.$nextTick(() => {
        const dom  = this.$refs.userList as HTMLElement
        this.showMoreBtn = dom.scrollHeight > 100
      })
    } else {
      this.$message.error(res.message);
    }
  }
  getWhitelistDetail() {
    const promise = this.formInfo.crowd_list.map(item => Api.getWhitelistDetailApi({
      id: item.id
    }))

    return Promise.all(promise).then(res => {
      this.userList = res.reduce((prev, item) => {
        prev = prev.concat(item.data.member_list.map(item => {
          item.title = item.member_name
          return item
        }))
        return prev
      }, [])
    })
  }
  

  handleCancel() {
    this.$router.back()
  }
  /**
   * 编辑
   */
  handleLink(info) {
    this.$router.push({
      name: 'myapp_wisdom_communication_content_detail',
      query: {
        id: String(info.id)
      }
    })
  }

  /**
   * 下载文件pdf
   */
  handleDownLoad() {
    // (this.$refs.detail as any).handleDownloadPdf()
    if(this.formInfo.pdf_url) {
      (this.$refs.detail as any).handleDownloadFile()
    } else {
      (this.$refs.detail as any).handleDownloadPdf()
    }
  }

  /**
   * 下载文件zip
   */
  handleDownLoadZip(){
    this.zipLoading = true
    const { title, content } = this.formInfo
    let files = []
    content.forEach((item: any) => {
      item.info_list.forEach((item: any) => {
        const attach = item.attachment || []
        files = files.concat(attach.map(val => {
          return {
            name: val.attachment_name,
            content: val.url
          }
        }))
      })
    })
    if(files.length === 0) {
      this.$message.error('您没有可以下载的附件～')
      this.zipLoading = false
      return
    }
    
    saveFileToZip(files, `${title}-附件.zip`).then(() => {
      this.$message.success('下载成功')
    }).catch(() => {
      this.$message.error('下载失败')
    }).finally(() => {
      this.zipLoading = false
    })
  }

  /**
   * 复制文本
   */
  handleCopyText() {
    (this.$refs.copyText as any).copy()
  }

}
</script>

<style lang="scss" scoped>
.top {
  padding: 20px;
  background-color: #f5f6f7;
  display: flex;
  position: relative;
  .cover {
    width: 160px;
    height: 90px;
    background-color: #f6f7f8;
  }

  .info {
    padding-left: 20px;
    .title {
      font-size: 14px;
      font-weight: bold;
    }

    .name {
      margin-top: 5px;
      color: #409EFF;
    }

    .time {
      margin-top: 5px;
      font-size: 14px;
      color: #666;
    }
  }

  .status {
    position: absolute;
    top: 0;
    right: 0;
  }

}
.line {
  color: #999;
  padding: 5px;
  border-bottom: 1px dashed #999;
  margin-top: 10px;
  font-size: 12px;
}

.list {
  margin-top: 20px;

  .item {
    margin-bottom: 40px;

    .dp-name {
      font-size: 14px;
      font-weight: bold;
    }

    .item-inner {
      margin-top: 5px;
      .link {
        color: #409EFF;
        cursor: pointer;
      }
    }

  }
}

.user {
  padding: 15px;
  background-color: #f6f7f8;

  .user-title {
    font-weight: bold;
  }

  .user-list {
    margin-top: 10px;
    overflow: hidden;
    max-height: none;
  }

  .more {
    text-align: center;
    border-top: 1px solid #eee;
    padding-top: 10px;
    color: #999;
    cursor: pointer;
  }
}
</style>