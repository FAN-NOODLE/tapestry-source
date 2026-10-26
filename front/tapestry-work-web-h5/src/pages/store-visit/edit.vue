<template>
  <div v-if="detail" class="container">
    <div class="header">
      <div class="name">{{detail.creator.dept_name}}{{ $t('storeVisit.title.3') }}</div>
      <div class="title">
        <div class="store-name">{{detail.store_name}}</div>
      </div>
      <div class="store-info">
        <div class="flex">
          <div class="label">{{ $t('storeVisit.item.label.1') }}:</div>
          <div class="value">{{detail.date.replaceAll('-', '.')}} {{detail.time_span}}</div>
        </div>
        <div class="flex">
          <div class="label">{{ $t('storeVisit.item.label.2') }}:</div>
          <div v-if="detail.store_manager.length" class="value flex">
            <div class="name" v-for="(m, index) in detail.store_manager" :key="index">{{ m.name }}</div>
          </div>
          <div v-else class="value">--</div>
        </div>
      </div>
    </div>
    <div id="item--1" class="content">
      <div v-if="!isEdit" class="title">
        <span>{{ $t('storeVisit.tips.32') }}</span>
        <div v-if="isCreator(detail)"  class="edit" @click="handleEdit">
          <span>{{ $t('storeVisit.tips.33') }}</span>
          <div class="icon-edit"></div>
        </div>
      </div>
      <div v-if="isEdit" class="item add" @click="handleSelectFeedback()">
        <div class="icon-plus"></div>
        <div>{{ $t('storeVisit.btn.4') }}</div>
      </div>
      <div class="tips" v-if="recordList.length==0">{{ $t('storeVisit.tips.1') }}</div>
      <div class="list">
        <div class="item" :id="`item-${index}`" :class="{expand: false}" v-for="item,index in recordList" :key="item.id">
          <div class="item-head">
            <div class="item-title">
              <span v-if="!isEdit || !item.is_custom">{{ item.type_name }}</span>
              <input v-else v-model="item.type_name" ref="input" class="item-input" type="text" placeholder="自定义类型" @click="handleInputClick(index)" @blur="handleInputBlur(index)" />
            </div>
            <!-- <div v-else class="item-title">{{item.type_name}}</div> -->
            <tool-popup v-if="isEdit">
              <div class="icon-more"></div>
              <template #content>
                <div class="tool-item" @click="handleSelectFeedback(index)">
                  <div class="icon-edit"></div>
                  <div class="text">{{ $t('storeVisit.btn.6') }}</div>
                </div>
                <div class="tool-item" @click="handleDelete(index)">
                  <div class="icon-delete"></div>
                  <div class="text">{{ $t('storeVisit.btn.7') }}</div>
                </div>
              </template>
            </tool-popup>
          </div>
          <div class="item-expand">
            <div class="item-content">
              <div class="item-part">
                <div class="item-label">{{ $t('storeVisit.tips.2') }}</div>
                <div class="item-value">
                  <div class="file-list">
                    <div v-for="media,mediaIndex in item.media_list" :key="mediaIndex" class="file-item" @click="handlePreview(media, item.media_list)">
                      <img class="file-item-img" :src="media.media_type == 1 ? media.content||media.url : media.poster || media.video_image_url" alt=""></img>
                      <!-- <video v-else class="file-item-img" :src="media.content||media.url" :poster="media.poster || media.video_image_url"></video> -->
                      <div v-if="media.status === 'uploading'" class="file-item-cover">
                        <div class="icon-upload-white"></div>
                        <com-progress class="progress" :value="media.progress*100"></com-progress>
                      </div>
                      <div v-if="media.status !== 'uploading' && isEdit" class="icon-close-circle" @click.stop="handleDeleteMedia(index, mediaIndex)"></div>
                      <div v-if="media.status !== 'uploading' && media.media_type == 2 && media.video_second" class="file-item-time">{{ $helper.formatTime(media.video_second) }}</div>
                    </div>
                    
                    <com-upload v-if="isEdit" class="file-item add" v-model="item.media_list" :use-sheet="!isIos" accept="image/jpg,image/jpeg,image/png,video/*" @upload-success="handleUploadSuccess($event, index)">
                      <div class="icon-upload"></div>
                    </com-upload>
                  </div>
                </div>
              </div>
              <div class="item-part">
                <div class="item-label">{{ $t('storeVisit.tips.3') }}</div>
                <div class="item-value">
                  <div class="textarea" @click="handleEditor(index)">
                    <!-- <com-editor v-model="item.content" @blur="handleInputBlur(index)"></com-editor> -->
                     <div v-if="item.content" v-html="item.content"></div>
                     <div v-else class="placeholder">{{ $t('storeVisit.placeholder.2') }}</div>
                  </div>
                </div>
              </div>
            </div>
            <div class="item-footer">
              <div class="time">
                <div class="time-label">{{ $t('storeVisit.tips.4') }}：</div>
                <div class="time-value">{{formatTime(item.created_at*1000)}}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div v-if="isEdit" class="footer">
      <van-button color="#000000" @click="handleSave">{{ $t('storeVisit.btn.21') }}</van-button>
    </div>
    <select-feedback-popup v-model="showSelectFeedbackPopup" @cancel="handleSelectFeedbackCancel" @confirm="handleSelectFeedbackConfirm"></select-feedback-popup>
    <com-editor-popup v-model="showEditorPopup" :content="content" @confirm="handleEditorConfirm" @cancel="handleEditorCancel"></com-editor-popup>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Watch } from "vue-property-decorator";
import ClockPopup from "./components/clock-popup.vue";
import GuidePopup from "./components/guide-popup.vue";
import ComUpload from "./components/com-upload.vue";
import ComProgress from "./components/com-progress.vue";
import ComEditor from "./components/com-editor.vue";
import ComEditorPopup from "./components/com-editor-popup.vue";
import SelectFeedbackPopup from "./components/select-feedback-popup.vue";
import SelectStorePopup from "./components/SelectStorePopup/index.vue";
import ToolPopup from "./components/tool-popup.vue";
import { completeVisitApi, deleteVisitRecordApi, deleteVisitRecordMediaApi, getFirstOpenStateApi, getStoreVisitInfoApi, getVisitRecordListApi, saveVisitRecordApi, saveVisitRecordMediaApi, startVisitApi } from "@/api/store-visit";
import dayjs from "dayjs";
import { getVideoCover, getVideoDuration } from "@/utils/upload/utils";
import { ImagePreview } from 'vant';
import PreviewVideoPopup from "./components/PreviewVideoPopup";
import { base64ToImage } from "@/utils/tools";
import { cosUploader } from "@/utils/upload";
import { getFileData } from "@/utils/download";
const feedbacks = [
  '店铺陈列',
  '店铺升级',
  '市场活动',
  '市场拓展',
  '产品/库存',
  '销售反馈',
  '服务技巧',
  '店铺运营',
  '系统反馈',
]

function detectDevice() {
  const userAgent = navigator.userAgent;

  if (/iPad|iPhone|iPod/.test(userAgent)) {
    return 'ios';
  } else if (/Android/.test(userAgent)) {
    return 'android';
  } else {
    return 'other';
  }
}

@Component({
  components: {
    ClockPopup,
    GuidePopup,
    ComUpload,
    ComProgress,
    SelectFeedbackPopup,
    ToolPopup,
    SelectStorePopup,
    ComEditor,
    ComEditorPopup
  }
})
export default class StoreVisitStart extends Vue {
  
  /** 选择反馈类型 */
  showSelectFeedbackPopup = false
  /** 文本编辑弹框 */
  showEditorPopup = false
  /** 巡店记录详情 */
  detail:StoreVisit.IDetail = null
  /** 巡店反馈记录列表 */
  recordList: StoreVisit.IFeedbackRecordItem[] = []

  current = -1

  isEdit = false
  
  
  get isIos() {
    return detectDevice() === 'ios'
  }
  /** 是否已经签到 */
  get disabledClock() {
    return this.detail && this.detail.to_store_time>0
  }

  /** 用户信息 */
  get userInfo():IUserInfo {
    return this.$store.getters["user/getUserInfo"]
  }

  get formatTime() {
    return (value) => dayjs(value).format('YYYY.M.D HH:mm')
  }


  /** 是否可以修改门店 */
  get disabledChangeStore() {
    return !(this.$route.query.third_id && this.userInfo.is_store_visit_manager)
  }

  get content () {
    if(this.current < 0) {
      return ''
    }
    return this.recordList[this.current]?.content
  }


  /**
   * 创建人
   */
  get isCreator() {
    return (item) => item.creator.member_id === this.userInfo.member_id
  }

  get trackInfo () {
    return {
      SV_id: this.detail.id,
      store_name: this.detail.store_name,
      store_code: this.detail.store_code,
      visit_time: `${this.detail.date} ${this.detail.time_span}`,
      is_synced: !!this.detail.is_report,
    }
  }

  
  mounted() {
    this.init()
  }
  destroyed() {
    PreviewVideoPopup.close()
  }
  /**
   * 初始化
   */
  async init() {
    const { id } = this.$route.query
    // 如果有ID，先获取详情
    await this.getDetail(Number(id))
    // 获取巡店反馈记录
    await this.getVisitRecordList()

  }
  /**
   * 获取详情
   */
  async getDetail(id:number) {
    const res:IApiResult<StoreVisit.IDetail> = await getStoreVisitInfoApi({
      id: Number(id)
    })
    if(res.code === 200) {
      this.detail = res.data
    }
  }

  /**
   * 巡店反馈记录
   */
  async getVisitRecordList() {
    const res:IApiResult<StoreVisit.IFeedbackRecordItem[]> = await getVisitRecordListApi({
      id: Number(this.$route.query.id),
    })

    if(res.code === 200) {
      this.recordList = res.data.map(item => {
        item.is_custom = !feedbacks.includes(item.type_name)
        return item
      })
    }
  }


  /**
   * 打开选择反馈
   */
  handleSelectFeedback(index =  -1) {
    if(!this.isEdit) return


    const hasEmpty = this.recordList.some(item => item.content.trim() === '' && !item.media_list.length)
    if (hasEmpty && index === -1) {
      this.$toast(this.$t('common.toast.9'))
      return
    }

    
    this.current = index
    this.showSelectFeedbackPopup = true
  }
  /**
   * 选择反馈取消
   */
  handleSelectFeedbackCancel() {
    this.showSelectFeedbackPopup = false
  }

  /**
   * 选择反馈确认
   * @param value 
   */
  async handleSelectFeedbackConfirm(e) {
    console.log(this.current);
    
    this.showSelectFeedbackPopup = false
    if(this.current === -1) {
      const params = {
        id: 0,
        visit_id: this.detail.id,
        type_name: e.label,
        content: "",
        is_custom: e.label === '自定义分类'
      }
      const id = await this.saveVisitRecord(params)
      this.recordList.unshift({
        ...params,
        id,
        isExpand: true,
        created_at: dayjs().unix(),
        media_list: []
      })
    } else {
      const item = this.recordList[this.current]
      saveVisitRecordApi({
        id: item.id,
        visit_id: item.visit_id,
        content: item.content,
        type_name: e.label
      })
      this.recordList[this.current].type_name = e.label
      this.recordList[this.current].is_custom = e.label === '自定义分类'
    }

    this.current = -1
  }

  /**
   * 保存巡店反馈
   * @param params 
   */
  async saveVisitRecord(params) {
    const res: IApiResult = await saveVisitRecordApi(params)

    if(res.code === 200) {
      return res.data
    }
  }

  handleInputBlur(index) {
    const item = this.recordList[index]
    saveVisitRecordApi({
      id: item.id,
      visit_id: item.visit_id,
      content: item.content,
      type_name: item.type_name
    })
  }


  /**
   * 获取焦点
   * @param index 
   */
  handleInputClick(index) {
    (this.$refs.input as any)[index].focus()
  }

  /**
   * 展开收起
   * @param index 
   */
  handleExpand(index) {
    this.$set(this.recordList[index], 'isExpand', !this.recordList[index].isExpand)
  }

  /**
   * 删除
   * @param index 
   */
  async handleDelete(index) {
    console.log(index);
    const item = this.recordList[index]
    await deleteVisitRecordApi({
      id: item.id
    })
    this.recordList.splice(index, 1)

    
    // 神策数据上报 -- 删除反馈
    this.$track.trackDeleteFeedback({
      ...this.trackInfo,
      feedback_id: String(item.id),
      feedback_category: item.type_name,
      is_customized: !feedbacks.includes(item.type_name),
      is_picture: item.media_list.length > 0,
      is_text: !!(item.content && item.content.length > 0),
    })
  }

  /**
   * 上传成功后
   * @param value 
   */
  async handleUploadSuccess(e, index) {
    const { id, visit_id, type_name } = this.recordList[index]
    
    const params = {
      id: id,
      visit_id: visit_id,
      type_name: type_name,
      media_type: e.media_type,
      url: e.url,
      video_second: e.video_second,
      video_image_url: e.video_image_url,
    }

    const res:IApiResult = await saveVisitRecordMediaApi(params)
    if (res.code === 200) {
      this.$set(e, 'id', res.data.media_id)
      if(!params.id) {
        this.$set(this.recordList[index], 'id', res.data.id)
      }
    }
  }

  /**
   * 删除媒体
   * @param index 
   * @param mediaIndex 
   */
  async handleDeleteMedia(index, mediaIndex) {
    const currentMedia = this.recordList[index].media_list[mediaIndex];
    if(currentMedia.id) {
      // 请求接口 删除素材
      await deleteVisitRecordMediaApi({
        media_id: currentMedia.id
      })
    }
    this.recordList[index].media_list.splice(mediaIndex, 1);
  }


  /**
   * 文本编辑
   * @param index 
   */
  handleEditor(index) {
    if(!this.isEdit) return
    this.showEditorPopup = true
    this.current = index;
  }
  /**
   * 文本编辑确认
   * @param e 
   */
  handleEditorConfirm(e) {
    this.showEditorPopup = false
    this.recordList[this.current].content = e

    // 页面滚动到
    this.$nextTick(() => {
      const item = document.getElementById(`item-${this.current === 0 ? -1 : this.current}`);
      if (item) {
        item.scrollIntoView({ behavior: 'smooth', block: 'start' });
      }
    });

    this.handleInputBlur(this.current)
  }


  handleEditorCancel() {
    this.showEditorPopup = false
    
    // 页面滚动到
    this.$nextTick(() => {
      const item = document.getElementById(`item-${this.current === 0 ? -1 : this.current}`);
      if (item) {
        item.scrollIntoView({ behavior: 'smooth', block: 'start' });
      }
    });
  }

  /**
   * 预览
   * @param index 
   */
  handlePreview(item, media_list) {
    
    if(item.media_type === 2) {
      PreviewVideoPopup.open({
        item
      })
    } else {
      const list = media_list.filter(item => item.media_type === 1).map(item => item.url)
      const mediaIndex = list.findIndex(media => media === item.url)
      ImagePreview({
        images: list,
        startPosition: mediaIndex,
      })
    }
  }

  handleEdit() {
    this.isEdit = true
  }

  /**
   * 保存
   */
  async handleSave() {
    const hasUploading = this.recordList.some(item => item.media_list.some(media => media.status === 'uploading'))
    if(hasUploading) {
      this.$toast(this.$t('common.toast.10'))
      return 
    }
    const hasEmpty = this.recordList.some(item => item.content.trim() === '' && !item.media_list.length)
    if (hasEmpty) {
      this.$toast(this.$t('common.toast.9'))
      return
    }
    this.isEdit = false
    this.track()
    // this.$router.back()
  }
  track() {
    const list = this.recordList.map(item => {
      return {
        feedback_id: String(item.id),
        feedback_category: item.type_name,
        is_customized: item.is_custom,
        is_picture: item.media_list.length > 0,
        is_text: !!(item.content && item.content.length > 0),
      }
    })
    console.log(this.recordList, list)
    // 神策数据上报 -- 提交巡店记录
    this.$track.trackSubmitReport({
      ...this.trackInfo,
      feedback_quantity: list.length ,
      feedback_category_list: list.map(item => item.feedback_category) ,
      is_synced: !!this.detail.is_report,
      button_name: '保存并提交',
    })

    list.forEach(item => {
      // 神策数据上报 -- 提交巡店记录详情
      this.$track.trackSubmitReportDetail({
        ...this.trackInfo,
        ...item,
        is_synced: !!this.detail.is_report,
        button_name: '保存并提交',
      })
    })
  }
  
}
</script>

<style lang="scss" scoped>
.icon-plus {
  width: 20px;
  height: 20px;
  background: url("~@/assets/store-visit/icon_plus.png") center/20px no-repeat;
}
.icon-caret-down {
  width: 48px;
  height: 48px;
  background: url("~@/assets/store-visit/icon_caret_down.png") center/16px no-repeat;
}
.icon-down-white {
  width: 24px;
  height: 24px;
  background: url("~@/assets/store-visit/icon_down_white.png") center/100% no-repeat;
}
.icon-clock {
  width: 24px;
  height: 24px;
  background: url("~@/assets/store-visit/icon_clock.png") center/100% no-repeat;
}
.icon-clock-gray {
  width: 24px;
  height: 24px;
  background: url("~@/assets/store-visit/icon_clock_1_gray.png") center/100% no-repeat;
}
.icon-feedback {
  width: 20px;
  height: 20px;
  background: url("~@/assets/store-visit/icon_feedback.png") center/100% no-repeat;
}
.icon-more {
  width: 16px;
  height: 16px;
  background: url("~@/assets/store-visit/icon_more.png") center/100% no-repeat;
}
.icon-upload {
  width: 98px;
  height: 98px;
  background: url("~@/assets/store-visit/icon_upload.png") center/24px no-repeat;
}
.icon-upload-white {
  width: 24px;
  height: 24px;
  background: url("~@/assets/store-visit/icon_upload_white.png") center/24px no-repeat;
}
.icon-close-circle {
  width: 24px;
  height: 24px;
  background: url("~@/assets/store-visit/icon_close_circle.png") center/14px no-repeat;
}



.container {
  min-height: 100vh;
  background: rgba(248, 249, 250, 1);
  padding-bottom: calc(env(safe-area-inset-bottom) + 64px);
  padding-bottom: calc(constant(safe-area-inset-bottom) + 64px);

  .header {
    color: #FFFFFF;
    padding: 16px 24px;
    background: linear-gradient(180deg, #3975C6 0%, #1870E8 100%);

    .name {
      font-weight: 600;
      font-size: 12px;
      position: relative;
    }

    .title {
      font-size: 14px;
      font-weight: 600;
      display: flex;
      align-items: center;
      justify-content: space-between; 
      margin-top: 4px;
    }
    .store-info  { 
      margin-top: 4px;
      .flex {
        display: flex;
        line-height: 16px;
        .label {
          color: rgba(222, 226, 230, 1);
          font-size: 12px;
        }
        .value {
          color: rgba(255, 255, 255, 1);
          font-size: 12px;
          margin-left: 8px;
          flex: 1;
          overflow: hidden;
          &.flex {
            display: flex;
            flex-wrap: wrap;
          }

          .name {
            margin-right: 8px;
            &:not(:last-child)::after {
              content: ',';
            }
          }
          
        }

        &:not(:last-child) {
          margin-bottom: 4px;
        }
      }
    }
  }
  
  .content {
    padding: 24px 16px;
    .title {
      font-size: 16px;
      font-weight: 600;
      padding: 0 12px;
      display: flex;
      align-items: center;
      justify-content: space-between; 
      
      .edit {
        display: flex;
        align-items: center;
        font-size: 14px;
        color: rgba(153, 153, 153, 1);

        .icon-edit {
          margin-left: 4px;
          width: 24px;
          height: 24px;
          background: url("~@/assets/store-visit/icon_edit.png") center/100% no-repeat;
        }
      }
    }

    .tips {
      font-size: 10px;
      color: rgba(153, 153, 153, 1);
      text-align: center;
      margin-top: 24px;
    }

    .list {
      margin-top: 24px;
    }

    .item {
      padding: 24px 16px;
      background-color: rgba(255, 255, 255, 1);
      margin-bottom: 8px;
      border-radius: 12px;

      .item-head {
        position: relative;
        display: flex;
        align-items: center;
        justify-content: space-between;
        .item-title {
          font-size: 14px;
          font-weight: 600;

          .item-input {
            width: 100%;
            height: 20px;
            padding-right: 24px;

            &:focus {
              outline: none;
            }

            &::placeholder {
              color: rgba(153, 153, 153, 0.3);
            }

            &:disabled {
              background-color: transparent;
            }
          }
        }

        .icon-caret-down {
          position: absolute;
          right: -16px;
          top: 50%;
          transform: translateY(-50%);
          transition: all  0.3s ease-in-out;
        }
      }
      .item-content {
        padding: 24px 0;

        .item-label {
          color: rgba(153, 153, 153, 1);
          font-size: 12px;
          line-height: 16px;
        }

        .item-value {
          margin-top: 12px;
        }

        .item-part:not(:last-child) {
          margin-bottom: 20px;
        }

        .file-list {
          display: flex;
          flex-wrap: wrap;

          .file-item {
            width: 98px;
            height: 98px;
            background-color: rgba(222, 226, 230, 1);
            border-radius: 4px;
            overflow: hidden;
            position: relative;
            margin-right: 8px;

            &:nth-child(n+4) {
              margin-top: 8px;
            }
            &:nth-child(3n) {
              margin-right: 0;
            }

            .file-item-img {
              width: 100%;
              height: 100%;
              display: block;
              object-fit: cover;
            }

            .file-item-cover {
              position: absolute;
              top: 0;
              left: 0;
              width: 100%;
              height: 100%;
              z-index: 1;
              background-color: rgba(0, 0, 0, 0.5);
              display: flex;
              align-items: center;
              justify-content: center;

              .progress {
                position: absolute;
                top: 50%;
                left: 50%;
                transform: translate(-50%, -50%);
              }
            }

            .icon-close-circle {
              position: absolute;
              top: 0;
              right: 0;
            }

            .file-item-time {
              position: absolute;
              top: 4px;
              left: 4px;
              height: 16px;
              border-radius: 23px;
              padding: 0 4px;
              line-height: 16px;
              font-size: 10px;
              font-weight: 700;
              color: #FFFFFF;
              background-color: #0D6EFD;
            }

            &.add {
              display: flex;
              align-items: center;
              justify-content: center;
            }
          }
        }

        .textarea {
          padding: 12px 16px;
          background-color: #F8F9FA;
          border-radius: 8px;
          min-height: 148px;
          font-size: 12px;
          line-height: 20px;

          .placeholder {
            color: rgba(153, 153, 153, 1);
          }
          ::v-deep ul {
              list-style-type: disc;
              padding-left: 20px;
          }
          ::v-deep ol {
              list-style-type: decimal;
              padding-left: 20px;
          }
        }
      }

      .item-footer {
        padding-top: 24px;
        border-top: 1px solid rgba(248, 249, 250, 1);
        line-height: 16px;
        display: flex;
        align-items: center;
        justify-content: space-between;

        .time {
          font-size: 12px;
          color: rgba(153, 153, 153, 1);
          display: flex;
          align-items: center;

          .time-label {
            font-weight: 600;
          }
        }
      }


      &.add {
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 14px;
        font-weight: 600;
        padding: 16px;

        .icon-plus {
          margin-right: 8px;
        }
      }

      &.expand {
        .item-head {
          .icon-caret-down {
            transform: translateY(-50%) rotate(0deg);
          }
        }
        .item-content {
          display: none;
        }
        .item-footer {
          display: none;
        }
      }
    }
  }
  

  .footer {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    z-index: 10;
    padding: 12px 24px;
    display: flex;
    padding-bottom: calc(env(safe-area-inset-bottom) + 12px);
    padding-bottom: calc(constant(safe-area-inset-bottom) + 12px);

    .van-button {
      height: 44px;
      border-radius: 8px;
      font-size: 14px;
      font-weight: 600;
      flex: 1;
      &.van-button--disabled {
        opacity: 1;
        background: rgba(222, 226, 230, 1) !important;
        border-color: rgba(222, 226, 230, 1) !important;
      }
    }
  }
}
</style>
