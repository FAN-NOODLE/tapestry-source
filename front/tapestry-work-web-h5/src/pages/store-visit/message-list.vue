<template>
  <div class="message">
    <div class="message-title">{{$t('storeVisit.tips.25')}}</div>
    <van-list class="list" 
        v-if="list.length>0"
        v-model="loading"
        :finished="finished"
        :finished-text="finishedText"
        :immediate-check="false"
        @load="getData">
      <div ref="message" class="message-item" :id="`id_${item.id}`" v-for="item, index in list" :key="index">
        <div class="message-item-inner" @touchstart="handleTouchStart($event, index)" @touchend="handleTouchEnd($event, index)">
          <div class="message-item-head">
            <van-image class="user-avatar" :src="item.member_info.avatar"></van-image>
            <div class="user-name">{{item.member_info.name}}</div>
            <div class="user-dept">{{item.member_info.dept_name}}</div>
          </div>
          <div class="message-item-content">
            <div class="text">{{item.content}}</div>
            <div class="imgs">
              <van-image v-for="media, mediaIndex in item.media_list" :key="mediaIndex" class="img" :src="media.video_image_url || media.url" @click.stop="handlePreview($event, media, item.media_list)"></van-image>
            </div>
          </div>
          <div class="message-item-footer">
            <div class="time">{{$helper.dateFormat(item.created_at, 'yyyy.MM.dd hh:mm:ss')}}</div>
            <div class="btn" @click.stop="handleReply(item, index)">{{ $t('storeVisit.popup.btn.17') }}</div>
          </div>
        </div>
        <div v-if="item.children" class="reply-list">
          <div class="reply-item" :id="`id_${replyItem.id}`" v-for="replyItem, replyIndex in item.children" @touchstart.stop="handleTouchStart($event, index, replyIndex)" @touchend.stop="handleTouchEnd($event, index, replyIndex)">
            <div class="reply-item-header">
              <div class="name">{{replyItem.member_info.name}}</div>
              <div class="dept">{{replyItem.member_info.dept_name}}</div>
              <div class="icon-reply"></div>
              <div class="name">{{replyItem.reply_member_info.name}}</div>
              <div class="dept">{{replyItem.reply_member_info.dept_name}}</div>
            </div>
            <div class="reply-item-content">
              <div class="text">{{replyItem.content}}</div>
              <div class="imgs">
                <van-image v-for="media, mediaIndex in replyItem.media_list" :key="mediaIndex" class="img" :src="media.video_image_url || media.url" @click.stop="handlePreview($event, media, item.media_list)"></van-image>
              </div>
            </div>
            <div class="reply-item-footer">
              <div class="time">{{$helper.dateFormat(replyItem.created_at, 'yyyy.MM.dd hh:mm:ss')}}</div>
              <div class="btn" @click.stop="handleReply(replyItem, index, replyIndex)">{{ $t('storeVisit.popup.btn.17') }}</div>
            </div>
          </div>
        </div>
      </div>
    </van-list>
    <div v-if="list.length==0" class="message-empty">{{ $t('storeVisit.empty.6') }}</div>
    <!-- 发送评论 -->
    <send-message-popup v-model="showSendMessagePopup" :placeholder="placeholder" @confirm="handleSendMessageConfirm"></send-message-popup>
    <!-- 评论删除/复制/回复 -->
    <sheetmenu-popup v-model="showSheetmenuPopup" :disabled-delete="disabledDelete" @confirm="handleSheetmenuConfirm"></sheetmenu-popup>
  </div>
</template>
<script lang="ts"> 
import { Component, Prop, Vue, Watch } from 'vue-property-decorator';
import SendMessagePopup from './components/send-message-popup.vue';
import SheetmenuPopup from './components/sheetmenu-popup.vue';
import { deleteCommentApi, getCommentListApi, saveCommentApi } from '@/api/store-visit';
import PreviewVideoPopup from './components/PreviewVideoPopup';
import { ImagePreview } from 'vant';
import dayjs from 'dayjs';

@Component({
  components: {
    SendMessagePopup,
    SheetmenuPopup
  },
})
export default class Message extends Vue {
  @Prop({ type: Object, default: () => {} }) detail: StoreVisit.IDetail;
  /** 发送留言弹窗 */
  showSendMessagePopup = false;
  /**  评论操作弹框弹窗 */
  showSheetmenuPopup = false;
  list: StoreVisit.ICommentItem[] = []
  
  placeholder = 'storeVisit.tips.24';
  /** 回复评论id */
  replyCommentId = 0
  /** 当前需要处理的评论 */
  currentComment = {
    item: null,
    index: -1,
    subIndex: -1
  }

  /** 分页 */
  meta: IMeta = {
    page: 1,
    size: 10,
  }
  /** 是否正在加载数据  */
  loading = false;
  /** 是否没有更多数据 */
  finished = false;

  timer: any = null

  get finishedText() {
    return this.list.length > 0 ? this.$t('storeVisit.empty.5'):''
  }

  /** 用户信息 */
  get userInfo():IUserInfo {
    return this.$store.getters["user/getUserInfo"]
  }

  get disabledDelete() {
    if(!this.currentComment.item) return true
    const { member_id } = this.currentComment.item.member_info
    return this.userInfo.member_id !== member_id
  }


  get trackInfo() {
    return {
      SV_id: this.detail.id,
      store_name: this.detail.store_name,
      store_code: this.detail.store_code,
      SV_initiator: this.detail.creator.name,
      SV_initiator_id: this.detail.creator.user_id,
      SV_initiator_department: this.detail.creator.dept_name,
      visit_time: this.detail.date + ' ' + this.detail.time_span,
      report_time: dayjs(this.detail.updated_at*1000).format('YYYY-MM-DD HH:mm:ss'),
    }
  }

  previewVideoPopup = null

  async mounted() {
    await this.getData()

    // 让页面滚动到对应锚点
    const comment_id = this.$route.query.comment_id as string
    if (comment_id) {
      await this.$nextTick()
      setTimeout(() => {
        document.getElementById(`id_${comment_id}`).scrollIntoView({
          behavior: "smooth",
        })
      }, 100)
    }
  }

  destroyed() {
    PreviewVideoPopup.close()
  }

  async getData() {

    const params = {
      id: this.detail.id,
      ...this.meta,
    }

    try {
      const {code, message, data}: IApiResult = await getCommentListApi(params)
      if (code !== 200) {
        this.finished = true;
        this.loading = false;
        this.$toast.fail(message);
        return;
      }
      this.loading = false;
      this.list = this.meta.page == 1 ? data.list : this.list.concat(data.list)
      this.finished = this.list.length === data.total;
      this.meta.page = this.meta.page + 1

    } catch (error) {
      console.log(error);
      this.finished = true;
      this.loading = false;
      this.list = []
    }
  }

  /**
   * 发布评论
   */
  handleSendMessage() {
    this.replyCommentId = 0
    this.placeholder = 'storeVisit.tips.24'
    this.showSendMessagePopup = true;
  }

  /**
   * 取消发布
   */
  handleSendMessageCancel() {
    this.showSendMessagePopup = false;
  }

  /**
   * 发布评论
   * @param data 
   */
  async handleSendMessageConfirm(data) {
    this.showSendMessagePopup = false;
    // 执行发送逻辑
    console.log(data);
    const params = {
      ...data,
      id: this.detail.id,
    }
    if(this.replyCommentId) {
      params.reply_comment_id = this.replyCommentId
    }

    const res =  await saveCommentApi(params)
    console.log(res);
    // 神策数据上报 -- 评论
    this.$track.trackComments({
      ...this.trackInfo,
      comment_id: String(res.data),
      is_picture: params.media_list.length > 0,
      is_text: !!params.content,
      comment_type: params.reply_comment_id ?'回复': '发表评论'
    })
    
    this.meta.page = 1
    this.getData()
  }


  /**
   * 底部菜单
   */
  handleSheetmenu() {
    this.showSheetmenuPopup = true
  }

  /**
   * 取消底部菜单
   */
  handleSheetmenuCancel() {
    this.showSheetmenuPopup = false
  }

  /**
   * 确认底部菜单
   */
  handleSheetmenuConfirm(value) {
    this.showSheetmenuPopup = false
    switch (value) {
      case 1:
        // 回复
        this.handleReply(this.currentComment.item)
        break;
      case 2:
        // 复制
        this.handleCopy(this.currentComment.item)
        
        break;
      case 3:
        // 删除
        this.handleDelete()
        break;
    }
  }

  handleReply(item: StoreVisit.ICommentItem, index=-1, replyIndex=-1) {
    this.showSendMessagePopup = true
    this.placeholder = `${this.$t('storeVisit.popup.btn.17')} ${item.member_info.name}`
    this.replyCommentId = item.id
  }

  handleCopy(item) {
    //复制文本到剪切板
    if(!item.content) {
      this.$toast(this.$t('common.toast.11'))
    }

    const content = item.content;
    let oInput = document.createElement("input");
    oInput.value = content;
    document.body.appendChild(oInput);
    oInput.select(); // 选择对象;
    document.execCommand("Copy"); // 执行浏览器复制命令
    oInput.remove();
    this.$toast(this.$t('common.toast.12'));

    // 神策数据上报 -- 评论
    this.$track.trackComments({
      ...this.trackInfo,
      comment_id: String(item.id),
      is_picture: item.media_list.length > 0,
      is_text: !!item.content,
      comment_type: '复制',
    })
  }

  async handleDelete() {
    const { item, index, subIndex } = this.currentComment
    const res:IApiResult = await deleteCommentApi({
      comment_id: item.id
    })
    if(res.code === 200) {
      if(subIndex === -1) {
        this.list.splice(index, 1)
      } else {
        this.list[index].children.splice(subIndex, 1)
      }

      // 神策数据上报 -- 评论
      this.$track.trackComments({
        ...this.trackInfo,
        comment_id: String(item.id),
        is_picture: item.media_list.length > 0,
        is_text: !!item.content,
        comment_type: '删除',
      })
    }
  }

  /**
   * 触摸开始
   * @param e 
   * @param index 
   */
  handleTouchStart(e:any, index, subIndex = -1) {
    
    
    // 给目标元素加一个背景颜色
    if(subIndex === -1) {
      this.$refs.message[index].style.background = 'linear-gradient(0deg, rgba(255, 255, 255, 0.3), rgba(255, 255, 255, 0.3)),linear-gradient(0deg, #E9ECEF, #E9ECEF)'
    } else {
      const doms:any[] = (this.$refs.message[index] as HTMLElement).getElementsByClassName('reply-item') as any
      doms[subIndex].style.background = 'linear-gradient(0deg, #E9ECEF, #E9ECEF),linear-gradient(0deg, rgba(255, 255, 255, 0.3), rgba(255, 255, 255, 0.3)),linear-gradient(0deg, rgba(0, 0, 0, 0.05), rgba(0, 0, 0, 0.05))'
    }
    clearTimeout(this.timer)
    this.timer = setTimeout(() => {
      this.handleSheetmenu()

      if(subIndex === -1) {
        this.currentComment = {
          item: this.list[index],
          index,
          subIndex
        }
      } else {
        this.currentComment = {
          item: this.list[index].children[subIndex],
          index,
          subIndex
        }
      }
    }, 500);

  }

  /**
   * 触摸结束
   * @param e 
   * @param index 
   * @param subIndex 
   */
  handleTouchEnd(e:any, index, subIndex = -1) {
    clearTimeout(this.timer)
    
    
    if(subIndex === -1) {
      this.$refs.message[index].style.background = 'none'
    } else {
      const doms:any[] = (this.$refs.message[index] as HTMLElement).getElementsByClassName('reply-item') as any
      doms[subIndex].style.background = 'linear-gradient(0deg, rgba(255, 255, 255, 0.3), rgba(255, 255, 255, 0.3)), linear-gradient(0deg, var(--Grey-200, #E9ECEF), var(--Grey-200, #E9ECEF))'
    }
  }


  /**
   * 预览
   * @param index 
   */
  handlePreview(e, item, media_list) {
    e.stopPropagation();
    e.preventDefault();
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

  
}
</script>

<style lang="scss" scoped>
.icon-reply {
  width: 16px;
  height: 16px;
  background: url("~@/assets/store-visit/icon_reply.png") center/100% no-repeat;
}
.message {
  padding-bottom: 100px;
  user-select: none;
  background-color: rgba(248, 249, 250, 1);
  .message-title {
    font-size: 14px;
    font-weight: 600;
    text-align: center;
    padding: 24px 0;
    background-color: rgba(248, 249, 250, 1);
  }

  .message-empty {
    text-align: center;
    font-size: 10px;
    color: rgba(151, 151, 151, 1);
  }

  .message-item { 
    padding: 12px 0;

    &:not(:last-child) {
      margin-bottom: 16px;
    }

    .message-item-head { 
      display: flex;
      align-items: center;
      height: 28px;
      padding: 0 24px;

      .user-avatar {
        width: 24px;
        height: 24px;
        border-radius: 50%;
        margin-right: 8px;
        overflow: hidden;
      }

      .user-name {
        font-size: 14px;
        font-weight: 700;
      }

      .user-dept {
        margin-left: 4px;
        padding: 0 4px;
        border-radius: 2px;
        font-size: 10px;
        line-height: 20px;
        height: 20px;
        background-color: rgba(206, 212, 218, 1);
      }
    }

    .message-item-content {
      padding: 0 24px 0 56px;
      margin-top: 8px;

      .text {
        font-size: 12px;
        font-weight: 400;
        line-height: 20px;
        white-space: wrap;
        word-break: break-word;
      }
      .imgs {
        display: flex;
        flex-wrap: wrap;

        .img {
          width: 84px;
          height: 84px;
          border-radius: 4px;
          background-color: rgba(0, 0, 0, 0.2);
          margin-right: 8px;
          overflow: hidden;

          &:nth-child(3n) {
            margin-right: 0;
          }

          &:nth-child(n+4) {
            margin-top: 8px;
          }
        }
      }

      .text~.imgs {
        margin-top: 8px;
      }
    }

    .message-item-footer {
      padding: 0 24px 0 56px;
      margin-top: 8px;
      display: flex;
      justify-content: space-between;

      .time {
        font-size: 12px;
        color: rgba(153, 153, 153, 1);
      }

      .btn {
        color: rgba(153, 153, 153, 1);
        font-size: 12px;
        font-weight: 600;
      }
    }

    &.active {
      background: linear-gradient(0deg, rgba(255, 255, 255, 0.3), rgba(255, 255, 255, 0.3)), linear-gradient(0deg, var(--Grey-200, #E9ECEF), var(--Grey-200, #E9ECEF));

    }

    .reply-list {
      padding: 0 8px 0 40px;

      .reply-item  {
        margin-top: 16px;
        padding: 8px 16px;
        border-radius: 8px;
        background: linear-gradient(0deg, rgba(255, 255, 255, 0.3), rgba(255, 255, 255, 0.3)), linear-gradient(0deg, var(--Grey-200, #E9ECEF), var(--Grey-200, #E9ECEF));
        background-size: 100% 100%;
        background-repeat: no-repeat;
      }
      .reply-item-header {
        display: flex;
        align-items: center;
        height: 24px;
        font-size: 12px;
        .name {
          font-size: 14px;
          font-weight: 700;
        }

        .dept {
          margin-left: 4px;
          padding: 0 4px;
          border-radius: 2px;
          font-size: 10px;
          line-height: 16px;
          height: 16px;
          background-color: rgba(206, 212, 218, 1);
        }

        .icon-reply {
          margin-left: 4px;
        }
      }
      .reply-item-content {
        margin-top: 8px;
        .text {
          font-size: 12px;
          font-weight: 400;
          line-height: 20px;
        }
        .imgs {
          display: flex;
          flex-wrap: wrap;

          .img {
            width: 84px;
            height: 84px;
            border-radius: 4px;
            background-color: rgba(0, 0, 0, 0.2);
            margin-right: 8px;
            overflow: hidden;

            &:nth-child(3n) {
              margin-right: 0;
            }

            &:nth-child(n+4) {
              margin-top: 8px;
            }
          }
        }

        .text~.imgs {
          margin-top: 8px;
        }
      }
      .reply-item-footer {
        margin-top: 8px;
        display: flex;
        justify-content: space-between;

        .time {
          font-size: 12px;
          color: rgba(153, 153, 153, 1);
        }

        .btn {
          color: rgba(153, 153, 153, 1);
          font-size: 12px;
          font-weight: 600;
        }
      }
    }
  }
}

::v-deep .van-list__loading, 
::v-deep .van-list__finished-text, 
::v-deep .van-list__error-text {
  font-size: 10px;
  color: rgba(151, 151, 151, 1);
}
</style>