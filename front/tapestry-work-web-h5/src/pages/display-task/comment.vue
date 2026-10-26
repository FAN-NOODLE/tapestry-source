<template>
  <div v-if="imageList.length" class="container">
    <!-- 图片展示区域 -->
    <com-swipe 
      class="swipe" 
      :class="{'has-comment': imageList[current].num>0}" 
      :initialSwipe="initialSwipe" 
      :current="current" 
      :list="imageList" 
      @change="handleSwipeChange"></com-swipe>
    <!-- 评论展示区域 -->
    <div v-if="imageList[current].comment.length" class="content">
      <template v-for="(item, index) in imageList[current].comment" >
        <comment-item 
          v-if="item.action !== 1"
          :key="index"
          :item="item"
          :is_vm="is_vm"
          :currentImage="imageList[current]"
          :disabled="show(3)"
          @delete="handleDelete(index)"
          @edit="handleEdit(index)"
          @reply="handleReply(index)"
          @delete-image="handleDeleteImage($event, index)">
        </comment-item>
      </template>
    </div>
    <!-- 底部操作 -->
    <div class="footer">
      <div class="back" @click="handleBack">
        <div class="back-icon"></div>
        <div class="back-text">返回</div>
      </div>
      <!-- 标记为未通过 -->
      <div 
        v-if="is_vm === 1" class="mark" 
        :class="{
          cancel:imageList[current]._status === 1,
          disabled: !imageList[current].is_draft && imageList[current].status === 1
        }" 
        @click="handleToggleMark">
        <div class="mark-icon"></div>
        <div class="mark-text">{{imageList[current]._status === 1 ? '取消标记' : '标记不通过'}}</div>
      </div>
      <!-- 添加备注 -->
      <div v-if="show(1)" class="input" @click="handleInput">输入备注内容</div>
      <!-- 添加图片备注 -->
      <upload-image v-if="is_vm === 1" @change="handleUploadChange">
        <div class="image-remark" @click="handleInput">
          <div class="image-remark-icon"></div>
          <div class="image-remark-text">图片备注</div>
        </div>
      </upload-image>
      
    </div>
    <!-- 评论弹框 -->
    <comment-popup ref="commentPopup" @confirm="handleInputConfirm"></comment-popup>
    <!-- 删除确认 -->
    <confirm-popup ref="deletePopup" title="确认删除" @confirm="handleDeleteConfirm">
      <div>确认删除当前备注？</div>
    </confirm-popup>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Watch } from "vue-property-decorator";
import ComSwipe from "./components/com-swipe.vue";
import CommentItem from "./components/comment-item.vue";
import CommentPopup from "./components/comment-popup.vue";
import ConfirmPopup from "./components/confirm-popup.vue";
import UploadImage from "./components/upload-image.vue";
import { getBaseInfoApi } from '@/api/common'
import { getAuthApi, saveCommentApi, deleteCommentApi } from '@/api/display-task'
@Component({
  components: {
    ComSwipe,
    CommentItem,
    CommentPopup,
    ConfirmPopup,
    UploadImage
  }
})
export default class DisplayTaskComment extends Vue {
  /** 图片索引 */
  current = 0
  initialSwipe = 0
  /** 编辑下标 */
  editIndex =  -1;
  /** 删除下标 */
  deleteIndex = -1;
  /** 是否是vm */
  is_vm = 0;
  commentImage: IImageItem[] = []
  /** 图片列表 */
  imageList: IImageItem[] = []


  /** 任务详情 */
  get detail(): IDisplayTaskDetail {
    return this.$store.getters["displayTask/getDetail"]
  }

  set detail(value: IDisplayTaskDetail) {
    this.$store.commit("displayTask/setDetail", value)
  }

  /** 用户信息 */
  get userInfo(): IUserInfo {
    return this.$store.getters["user/getUserInfo"] || {}
  }

  /** 控制页面区域展示 */
  get show() {
    const { finish_status } = this.detail.store_info
    // const { status } = this.detail.task_info
    return (key: number, item?: ICommentItem) => {
      switch (key) {
        case 1:
          // VM身份
          if(this.is_vm == 1) return true
          // 导购身份 任务待提交 | 任务审核未通过
          return this.is_vm == 0 && [0, 2].includes(finish_status)
        case 2:
          // 导购身份 && VM的备注 || VM身份 && 备注已经保存
          return this.is_vm == 0 && item.is_vm === 1 || this.is_vm === 1 && item.id
        case 3:
          // VM身份
          if(this.is_vm == 1) return false
          // 导购身份 任务审核中 || 审核已通过
          return this.is_vm == 0 && [1, 3].includes(finish_status)
        default:
          return true
      }
    }
  }


  @Watch("detail", { deep: true })
  onDetailChange(value: IDisplayTaskDetail) {
    this.$store.commit("displayTask/setDetail", value)
  }

  @Watch("imageList", { deep: true })
  onImageListChange(value: IImageItem[]) {
    const { idx } = this.$route.query
    this.detail.task_info.template_info[Number(idx)].image_list = value
    this.$store.commit("displayTask/setDetail", this.detail)
  }

  async created() {
    const { idx, index, is_vm } = this.$route.query
    this.current = Number(index) || 0
    this.initialSwipe = Number(index) || 0
    document.title = this.detail.task_info.template_info[Number(idx)].template_name
    this.is_vm = is_vm ? Number(is_vm) : 0
    this.imageList = this.detail.task_info.template_info[Number(idx)].image_list
    
    // 直接进入的时候获取用户信息
    await this.getAuth()
    this.getBaseInfo()
  }
  mounted() {
    window.addEventListener('resize', this.keyboard);
  }
  destroyed() {
    window.removeEventListener('resize', this.keyboard);
  }
  keyboard() {
    if (document.activeElement.tagName === 'INPUT' || document.activeElement.tagName === 'TEXTAREA') {
      document.activeElement.scrollIntoView({ behavior: 'smooth', block: 'center', inline: 'center' });
    }
  }
  /**
   * @description 获取用户权限
   */
  async getAuth() {
    try {
      const {code, message, data}:IApiResult = await getAuthApi()
      if (code!== 200) {
        this.$toast.fail(message);
        return;
      }
      this.$store.commit("user/setUserInfo", {
        app_source: 1,
        ...this.userInfo, 
        ...data,
      })
    } catch (error) {
      console.log(error);
    }
  }
  /**
   * @description 获取用户权限
   */
  async getBaseInfo() {
    try {
      const {code, message, data}:IApiResult = await getBaseInfoApi()
      if (code!== 200) {
        this.$toast.fail(message);
        return;
      }
      this.$store.commit("user/setUserInfo", {
        ...this.userInfo, 
        ...data
      })
    } catch (error) {
      console.log(error);
    }
  }

  /**
   * @description 图片滑动
   * @param index 
   */
  handleSwipeChange(index: number) {
    console.log(123131, index);
    
    this.current = index
  }

  /**
   * @description 返回
   * @returns
   */
  handleBack() {
    this.$router.go(-1)
  }

  /** 
   * @description 输入备注内容
   */
  handleInput() {
    (this.$refs.commentPopup as any).showPopup()
  }

  /** 
   * @description 输入备注内容确认
   */
  handleInputConfirm(value: string, prefix: string) {
    // 编辑
    if(this.editIndex !== -1) {
      this.handleEditConfirm(value)
      return 
    } 
    // 新增
    this.handleAddConfirm(value, prefix)
    
  }

  /**
   * @description 构建评论
   * @param value 
   */
  createComment(value: string | IImageItem[], reply_name: string): ICommentItem {
    const { member_id, user_id, member_name, avatar } = this.userInfo
    const { store_code } = this.detail.store_info
    const comment: ICommentItem = {
      id: 0,
      member_id: member_id,
      member_user_id: user_id,
      member_name: member_name,
      avatar: avatar,
      store_code: store_code,
      action: 0,
      is_vm: this.is_vm,
      is_self: 1,
      comment_time: parseInt(Date.now()/1000+ ''),
      content: {},
      reply_name: reply_name||''
    }
    if(typeof value ==='string') {
      comment.content.text = value
    } else {
      comment.content.image_url = value
    }

    return comment
  }

  /**
   * 新增评论
   * @param value 
   */
  handleAddConfirm(value: string|IImageItem[], reply_name: string) {
    const comment = this.createComment(value, reply_name)
    
    this.imageList[this.current].num++;
    this.imageList[this.current].comment.push(comment)
    // 不是vm
    if(this.is_vm === 0) {
      const image = this.imageList[this.current]
      this.handleAutoSaveComment(image, comment, (id: number) => {
        comment.id = id
      })
    }
    // 滚动页面到底部
    this.$nextTick(() => {
      setTimeout(() => {
        document.documentElement.scrollTo({
          top: document.body.scrollHeight,
          behavior: "smooth"
        });
      }, 100)
    })
  }

  /** 
   * 导购自动保存评论
   */
  handleAutoSaveComment(image: IImageItem, comment: ICommentItem, callback:Function) {
    const display_task_id = this.detail.task_info.id
    const { allot_id } = this.detail.store_info
    const allot_data_id = image.id
    const params = {
      allot_id,
      display_task_id,
      allot_data_id,
      ...comment
    }
    saveCommentApi(params).then(res => {
      callback(res.data.id)
    })
  }
  
  /**
   * @description 删除评论
   * @param index 
   */
  handleDelete(index: number) {
    this.deleteIndex = index;
    (this.$refs.deletePopup as any).showPopup()
  }

  /** 
   * @description 删除确认
   */
  handleDeleteConfirm() {
    const index = this.deleteIndex
    const id = this.imageList[this.current].comment[index].id
    this.imageList[this.current].comment.splice(index, 1)
    this.imageList[this.current].num--;
    if(id) {
      this.handleAutoDeleteComment(id)
    }
    
  }

  /**
   * 自动删除备注
   * @param id 
   */
  handleAutoDeleteComment(id: number) {
    deleteCommentApi({id})
  }

  /**
   * @description 编辑评论
   * @param index 
   */
  handleEdit(index: number) {
    this.editIndex = index;
    (this.$refs.commentPopup as any).showPopup(this.imageList[this.current].comment[index].content.text)
  }

  /**
   * @description 编辑确认
   * @param value 
   */
  handleEditConfirm(value: string) {
    const image = this.imageList[this.current]
    const item = this.imageList[this.current].comment[this.editIndex]
    this.imageList[this.current].comment[this.editIndex].content.text = value
    this.editIndex = -1
    
    if(this.is_vm === 1) {
      if(item.id) {
        item.action = 2
      }
      return
    }
    this.handleAutoSaveComment(image, item, () => {} )
  }

  /**
   * @description 回复评论
   * @param index 
   */
  handleReply(index: number) {
    const { member_name } = this.imageList[this.current].comment[index];
    (this.$refs.commentPopup as any).showPopup('', member_name)
  }

  /**
   * @description 标记不通过/取消标记
   * @param image 
   */
  handleToggleMark() {
    if(!this.imageList[this.current].is_draft && this.imageList[this.current].status === 1) return
    this.imageList[this.current]._status = this.imageList[this.current]._status === 1? 0: 1;
    
  }

  /**
   * 上传图片
   * @param value 
   */
  handleUploadChange(value: IImageItem[]) {
    this.handleAddConfirm(value, '')
  }

  /**
   * @description 删除图片
   * @param index 
   */
  handleDeleteImage(imageIndex: number, commentIndex: number) {
    const comment = this.imageList[this.current].comment[commentIndex]
    if(comment.content.image_url.length === 1) {
      this.imageList[this.current].comment.splice(commentIndex, 1)
      this.imageList[this.current].num--;
    } else {
      comment.content.image_url.splice(imageIndex, 1)
    }
    
  }
}
</script>
<style lang="scss" scoped>
.container {
  min-height: 100vh;
  background-color: #fff;
  box-sizing: border-box;
  padding-bottom: calc(56px + constant(safe-area-inset-bottom)); /** iOS 11.0+ */
  padding-bottom: calc(56px + env(safe-area-inset-bottom)); /** iOS 11.0+ */

  .swipe {
    background-color: rgba($color: #000000, $alpha: 0.8);
    height: calc(100vh - 56px - constant(safe-area-inset-bottom));
    height: calc(100vh - 56px - env(safe-area-inset-bottom));

    &.has-comment {
      height: calc(100vh - 148px - constant(safe-area-inset-bottom));
      height: calc(100vh - 148px - env(safe-area-inset-bottom));
    }
  }

  .content {
    padding: 20px 15px;
  }

  .footer {
    width: 100vw;
    padding: 10px 20px;
    background: #FFFFFF;
    box-shadow: 0px -3px 10px 0px rgba(117,117,117,0.06);
    position: fixed;
    bottom: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    padding-bottom: calc(10px + constant(safe-area-inset-bottom)); /** iOS 11.0+ */
    padding-bottom: calc(10px + env(safe-area-inset-bottom)); /** iOS 11.0+ */
    z-index: 100;
    .back {
      display: flex;
      align-items: center;
      justify-content: center;
      flex-direction: column;
      font-size: 10px;
      font-weight: 400;
      color: #000000;
      line-height: 14px;

      &-icon {
        width: 20px;
        height: 20px;
        background: url("~@/assets/display-task/icon_back.png") center/100% no-repeat;
      }

      &-text {
        margin-top: 2px;
      }
    }
    .mark {
      display: flex;
      align-items: center;
      justify-content: center;
      flex-direction: column;
      font-size: 10px;
      font-weight: 400;
      color: #000000;
      line-height: 14px;
      margin-left: 10px;

      &-icon {
        width: 20px;
        height: 20px;
        background: url("~@/assets/display-task/icon_mark.png") center/16px 16px no-repeat;
      }

      &-text {
        margin-top: 2px;
      }

      &.cancel {
        .mark-icon {
          background-image: url("~@/assets/display-task/icon_mark_cancel.png");
        }
      }
      &.disabled {
        color: #a6a6a6;
        .mark-icon {
          background-image: url("~@/assets/display-task/icon_mark_dis.png");
        }
      }
    }
    .image-remark {
      display: flex;
      align-items: center;
      justify-content: center;
      flex-direction: column;
      font-size: 10px;
      font-weight: 400;
      color: #000000;
      line-height: 14px;
      margin-left: 10px;

      &-icon {
        width: 20px;
        height: 20px;
        background: url("~@/assets/display-task/icon_pic.png") center/100% no-repeat;
      }

      &-text {
        margin-top: 2px;
      }
    }

    .input {
      flex: 1;
      height: 38px;
      background: #F6F6F6;
      font-size: 12px;
      font-weight: 400;
      color: #A1A1A1;
      line-height: 17px;
      padding: 10px 15px;
      margin-left: 15px;
    }
  }
}
</style>