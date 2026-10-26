<template>
  <div v-if="detail.store_info" class="container">
    <!-- 任务详情 -->
    <div class="detail">
      <div class="flex-center">
        <div class="title">门店：{{detail.store_info.store_code}}</div>
        <div class="audit-status" :data-status="detail.store_info.finish_status">{{auditStatusLabel}}</div>
      </div>
      <div class="tips">
        <span>提交时间：{{ dateFormat(detail.store_info.submit_time||Date.now(), 'yyyy年MM月dd日 hh:mm:ss') }}</span>
        <span v-if="detail.store_info.finish_status === 3 || detail.store_info.finish_status === 4" class="grade">{{ detail.store_info.grade }}</span>
      </div>
      <div v-if="detail.store_info.finish_status === 4" class="tips">
        <span>初审VM：{{ detail.store_info.special_name }}</span>
      </div>
    </div>
    <!-- 陈列内容 -->
    <div class="content">
      <div class="label">陈列内容</div>
      <div class="list">
        <div class="item" v-for="(item, index) in detail.task_info.template_info" :key="index">
          <div class="item-title" :class="{required: item.is_required}">{{item.template_name}}</div>
          <div class="image-list">
            <!-- 图片列表 -->
            <div 
              class="image-wrap"
              v-for="(image, imageIndex) in item.image_list" 
              :key="imageIndex" >
              <image-item 
                
                :item="image"
                :disabled-delete="true"
                @comment="handleComment(index, imageIndex)"
                >
              </image-item>
              <div 
                class="image-tag" 
                :class="{ 
                  cancel: image._status === 1,
                  disabled: !image.is_draft && image.status === 1
                }"
                @click="handleToggleMark(image)">{{image._status === 1? '取消标记' : '标记不通过'}}</div>
            </div>
            <!-- 不能编辑 & 没有上传图片 -->
            <span v-if="item.image_list.length === 0" class="is-empty">未上传内容</span>
          </div>
        </div>
      </div>
    </div>
    <!-- 整体评价 -->
    <div class="remark">
      <div class="label">VM留言</div>
      <div class="remark__content">
        <textarea ref="input" class="textarea" v-model="detail.store_info.review_remark" @click="handleInputClick" placeholder="输入备注内容" :disabled="!isCorporateCheck() && !isSpecialCheck()" :maxlength="maxlength"></textarea>
        <div class="maxlength">{{ detail.store_info.review_remark.length }}/{{ maxlength }}</div>
      </div>
    </div>
    <!-- footer -->
    <div class="footer">
      <!-- 返回 -->
      <div class="back" @click="handleBack">
        <div class="back-icon"></div>
        <div class="back-text">返回</div>
      </div>
      <!-- 预览 -->
      <div class="preview" :class="{disabled: false}" @click="handlePreview">
        <div class="preview-icon"></div>
        <div class="preview-text">预览长图</div>
      </div>
      <div v-if="isSpecialCheck()" class="btn-group">
        <button class="btn plain" :class="{disabled: false}" @click="handleSave(1)">存草稿</button>
        <button class="btn" :class="{disabled: false}" @click="handleSave(2)">提交</button>
      </div>
      <div v-if="isCorporateCheck()" class="btn-group">
        <button class="btn plain" :class="{disabled: false}" @click="handleSubmit(1)">不通过</button>
        <button class="btn" :class="{disabled: markCount>0}" @click="handleSubmit(2)">通过</button>
      </div>
    </div>
    <!-- 评分 -->
    <grade-popup ref="gradePopup" @confirm="handleSubmitConfirm"></grade-popup>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Watch } from "vue-property-decorator";
import { auditStatusOptions } from "./const";
import { getAuthApi, getStoreTaskInfoApi, getTaskInfoApi, saveVmCheckApi, specialSaveApi } from "@/api/display-task";
import { Image as VanImage } from "vant";
import ImageItem from "./components/image-item.vue";
import GradePopup from './components/grade-popup.vue'
import { getBaseInfoApi } from '@/api/common'
import wx from 'weixin-js-sdk'
@Component({
  components: {
    ImageItem,
    VanImage,
    GradePopup
  },
})
export default class DisplayTaskPreview extends Vue {
  /** 备注内容 */
  content: string = "";
  /** 最大长度 */
  maxlength: number = 200;
  /** 评分 */
  grade: number = null;
  scrollTop = 0

  /** 任务详情 */
  get detail(): IDisplayTaskDetail {
    return this.$store.getters["displayTask/getDetail"]
  }

  /** 是否更新 */
  get update(): boolean {
    return this.$store.getters["displayTask/getUpdate"]
  }

  /** 审核状态 */
  get auditStatusLabel() {
    if(!this.detail.store_info) return ''
    return auditStatusOptions.find((item) => item.value === this.detail.store_info.finish_status).label
  }

  /** 用户信息 */
  get userInfo(): IUserInfo {
    return this.$store.getters["user/getUserInfo"]
  }

  /** 标记数量 */
  get markCount() {
    return this.detail.task_info.template_info.reduce((prev, item)=> {
      prev += item.image_list.filter(image => image.status === 0 && image._status === 1).length

      return prev
    }, 0)
  }

  get isSpecialCheck() {
    return() => {
      return this.userInfo.is_special === 1 && this.detail.store_info.is_special_check
    }
  }
  get isCorporateCheck() {
    return() => {
      return this.userInfo.is_special !== 1 && this.detail.store_info.is_corporate_check
    }
  }

  @Watch("detail", { deep: true })
  onDetailChange(value: IDisplayTaskDetail) {
    this.$store.commit("displayTask/setDetail", value)
  }

  async mounted() {
    this.$toast.loading({
      mask: true,
      message: "加载中..."
    });
    try {
      if(this.userInfo.app_source !== 1) {
        await this.getAuth()
      }
      await this.getBaseInfo()
      if(!this.detail.store_info) {
        await this.getStoreTaskInfo()
      }
      if(!this.detail.task_info) {
        await this.getDetail()
      }
      document.title = this.detail.store_info.store_name +  "陈列任务"
      this.content = this.detail.store_info.review_remark
    } catch (error) {
      console.log(error);
    } finally {
      this.$toast.clear()
    }

    if(this.$route.query.from !== '1') {
      // 神策数据上报 -- 进入VM
      this.$track.trackEnterDisplayTask({
        vm_entrance: 2
      })
    }

    this.$nextTick(() => {
      const scrollTop = localStorage.getItem("scrollTop") || 0
      if (Number(scrollTop) > 0) {
        document.documentElement.scrollTo({
          top: Number(scrollTop),
          behavior: "auto"
        });
      }
      localStorage.removeItem("scrollTop")
    })
  }

  /**生命周期 */
  private async created() {
    // 当前需要直接退出的页面
    if(this.$route.query.from !== '1') {
      // 用户使用 Android 设备
      window.history.pushState("_android_back", "", location.href);
      window.addEventListener("popstate", this.popstate, false);
    }
  }
  private async destroyed() {
    // 记得清除!!!!不然单页面，其它页面也可以用
    if(this.$route.query.from !== '1') {
      window.removeEventListener("popstate", this.popstate, false);
    }
  }

  popstate() {
    try {
      wx.closeWindow(); // 微信网页退出
    } catch (error) {
      console.log(1111, error);
    }
  }


  /**
   * 点击获取焦点
   */
  handleInputClick() {
    (this.$refs.input as any).focus()
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
   * @description 格式化时间
   * @param time 
   * @param format 
   */
  dateFormat(time: number, format = "yyyy/MM/dd") {
    return this.$helper.dateFormat(time, format);
  }

  /**
   * @description 获取门店任务信息
   */
  async getStoreTaskInfo() {
    try {
      const params = {
        id: Number(this.$route.query.id),
        store_id: Number(this.$route.query.store_id),
      }
      const { code, message, data }:IApiResult<IStoreInfo> = await getStoreTaskInfoApi(params)
      if(code!== 200) {
        this.$toast.fail(message)
        return
      }
      // 如果VM对不上去到任务详情
      if(data.owner_member_user_id !== this.userInfo.user_id && data.special_user_id !== this.userInfo.user_id || ![1, 4].includes(data.finish_status)) {
          this.$router.replace({
            path: 'display-task-detail-vm',
            query: {
              id: String(params.id),
            }
          });
        return
      }
      if(this.detail.task_info) {
        this.detail.task_info.template_info.forEach(item => {
          item.image_list = data.image_list
            .filter(val => val.identifier === item.template_name).map(val => {
              return {
                ...val,
                _status: val.status,
                comment: val.comment.map(comment => {
                  const images = comment.content.image_url?.map(url => ({
                    url: url as unknown as string,
                  }))||[]
                  return {
                   ...comment,
                   content: images.length ? {image_url: images} : {text: comment.content.text}
                  }
                })
              }
            })
        });
      }
      this.$store.commit("displayTask/setDetail", {
        ...this.detail,
        store_info: data,
      })
    } catch (error) {
      console.log(error);
      
    }
  }

    /**
   * @description 获取详情
   */
  async getDetail() {
    try {
      const params = {
        id: Number(this.$route.query.id),
      }
      const { code, message, data }: IApiResult<ITaskInfo> = await getTaskInfoApi(params)
      
      if(code!== 200) {
        this.$toast.fail(message)
        return
      }

      data.template_info.forEach(item => {
        item.image_list = this.detail.store_info.image_list
          .filter(val => val.identifier === item.template_name).map(val => {
            return {
              ...val,
              _status: val.status,
              comment: val.comment.map(comment => {
                const images = comment.content.image_url?.map(url => ({
                  url: url as unknown as string,
                }))||[]
                return {
                 ...comment,
                 content: images.length ? {image_url: images} : {text: comment.content.text}
                }
              })
            }
          })
      });
      // 数据格式化
      this.$store.commit("displayTask/setDetail", {
        ...this.detail,
        task_info: data
      })
      
    } catch (error) {
      console.log(error);
    }
  }

  /** 
   * @description 预览图片
   */
  handleComment(idx: number, index: number) {
    if(!this.isCorporateCheck() && !this.isSpecialCheck()) return
    this.$router.push({
      path: 'display-task-comment',
      query: {
        idx: String(idx),
        index: String(index),
        is_vm: '1'
      }
    })
    const scrollTop = document.documentElement.scrollTop;
    localStorage.setItem('scrollTop', scrollTop.toString())
    
  }

  /**
   * @description 标记不通过/取消标记
   * @param image 
   */
  handleToggleMark(image: IImageItem) {
    if(!this.isCorporateCheck() && !this.isSpecialCheck()) return
    if(!image.is_draft && image.status === 1) return
    image._status = image._status === 1? 0: 1;
  }

  /**
   * @description 返回
   */
  handleBack() {
    this.$router.back();
  }

  /**
   * @description 预览
   */
  handlePreview() {
    this.$router.push({
      path: 'display-task-preview',
    })
  }

  /**
   * @description 提交审核
   */
  handleSubmit(type: number) {
    if(type === 2) {
      if(this.markCount > 0) return;
      (this.$refs.gradePopup as any).showPopup(this.detail.store_info.grade)
    } else {
      this.submit(type)
    }
  }

  /**
   * @description 提交审核
   * @param type 1: 保存草稿 2: 提交审核
   */
  handleSave(type: number) {
    if(type === 2) {
      // if(this.markCount > 0) return;
      (this.$refs.gradePopup as any).showPopup()
    } else {
      this.submit(type)
    }
  }

  /**
   * @description 审核提交确认
   */
  handleSubmitConfirm(value:number) {
    this.grade = value;
    this.submit(2)
  }

  /**
   * @description 审核提交
   * @param type 
   */
  async submit(review_type: number) {
    
    try {
      const params = {
        allot_id: this.detail.store_info.allot_id,
        id: this.detail.task_info.id,
        review_type: review_type,
        review_remark: this.detail.store_info.review_remark.trim(),
        grade: this.grade,
        image_list: []
      }
      params.image_list = this.detail.task_info.template_info.reduce((prev, item) => {
        return prev.concat(item.image_list.map(image => {
          return {
            ...image,
            status: image._status,
            comment: image.comment.map(comment => {
              const images = comment.content.image_url?.map(val => val.url) ??[]
              return {
               ...comment,
                content: images.length ? {image_url: images} : {text: comment.content.text}
              }
            })
          }
        }))
      }, [])
      const { code, message, data }:IApiResult<ITaskInfo> = this.userInfo.is_special ? await specialSaveApi(params) : await saveVmCheckApi(params)
      if(code!==200) {
        this.$toast.fail(message)
        return
      }
      this.$store.commit('displayTask/setUpdateList', {
        update: true,
        updateDetail: true,
        message: review_type === 2 ? '审核完成' : this.userInfo.is_special ? '保存成功':'审核完成'
      })
      if(review_type === 2) {
        // 神策数据上报 -- 操作陈列审核通过（vm视角）
        this.$track.trackMarkPass({
          display_store: this.detail.store_info.store_code,
          display_score: this.grade
        })
      } else {
        // 神策数据上报 -- 操作陈列审核不通过（vm视角）
        this.$track.trackMarkNotPass({
          display_store: this.detail.store_info.store_code
        })
      }
      this.$router.back()
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
}
</script>
<style lang="scss" scoped>
.container {
  background-color: #F7F7F7;
  min-height: 100vh;
  box-sizing: border-box;
  padding-bottom: calc(56px + constant(safe-area-inset-bottom)); /** iOS 11.0+ */
  padding-bottom: calc(56px + env(safe-area-inset-bottom)); /** iOS 11.0+ */
    
  .detail {
    padding: 15px;
    background-color: #fff;
    margin-bottom: 15px;
    .title {
      font-size: 16px;
      font-weight: 600;
      color: #000000;
      line-height: 22px;
      display: flex;
      padding-right: 15px;
      flex: 1;
    }
    .audit-status {
      font-weight: 400;
      font-size: 10px;
      width: 60px;
      height: 15px;
      border-radius: 15px;
      display: flex;
      align-items: center;
      justify-content: center;
      // margin-right: 5px;
      &[data-status="0"] {
        color: #000;
        background: rgba(0,0,0,0.1);
      }
      &[data-status="1"],
      &[data-status="4"] {
        color: #000;
        background: #D1BFA7;
      }
      &[data-status="2"] {
        color: #E02020;
        background: #FFD0D0;
      }
      &[data-status="3"] {
        color: #000;
        background: #BBD1A7;
      }
    }

    .tips {
      font-size: 12px;
      font-weight: 400;
      color: #A6A6A6;
      line-height: 17px;
      margin-top: 10px;
      position: relative;

      .grade {
        position: absolute;
        right: 0;
        top: 40%;
        transform: translateY(-50%);
        font-size: 22px;
        font-weight: 600;
        color: #70A560;
        line-height: 30px;
        display: flex;
        align-items: center;
        justify-content: center;
        &::before {
          content: "评分:";
          font-size: 12px;
          font-weight: 400;
          margin-right: 2px;
          padding-top: 5px;
        }
      }
    }

   .line {
      width: 345px;
      height: 37px;
      background: rgba(119, 119, 119, 0.1);
      justify-content: flex-start;
      padding: 0 15px;
      margin-top: 10px;

      .status {
        font-size: 12px;
        font-weight: bold;
        margin-right: 10px;
        &[data-status="0"] {
          color: #C6AD7F;
        }
        &[data-status="1"] {
          color: #70A560;
        }
        &[data-status="2"] {
          color: #A6A6A6;
        }
      }
      .time {
        font-size: 12px;
        font-weight: 400;
        color: #000000;
        line-height: 17px;
      }
    }

    .desc {
      font-size: 12px;
      font-weight: 400;
      color: #777777;
      line-height: 17px;
      margin-top: 10px;
      text-align: justify;
      word-break: break-all;
    }

    .link {
      margin-top: 15px;

      &-item {
        border: 1px solid #A6A6A6;
        box-sizing: border-box;
        padding: 10px 15px;

        &:not(:last-child) {
          margin-bottom: 15px;
        }
      }

      &-text {
        font-size: 12px;
        font-weight: 400;
        color: #000000;
        line-height: 17px;
        flex: 1;
        padding-right: 15px;
        word-break: break-all;
      }

      &-btn {
        width: 62px;
        height: 24px;
        background: #000000;
        font-size: 11px;
        font-weight: 400;
        color: #FFFFFF;
        line-height: 17px;
      }
    }

    .countdown {
      margin-top: 15px;
      .van-count-down {
        line-height: 18px;
        font-size: 12px;
        color: #70A560;
      }

      &.start .van-count-down {
        color: #C6AD7F;
      }
    }
  }

  .content {
    background-color: #fff;
    padding: 15px;
    margin-bottom: 15px;

    .label {
      font-size: 14px;
      font-weight: 500;
      color: #000000;
      line-height: 20px;
    }

    .list {
      margin-top: 15px;
      padding:  0 5px;

      .item {
        .item-title {
          font-size: 12px;
          font-weight: 400;
          color: #000;
          line-height: 17px;

          &.required::before {
            content: "*";
            color: #E02020;
          }
        }

        .image-list {
          margin-top: 10px;
          display: flex;
          flex-wrap: wrap;
          .image-wrap {
            width: 100px;
            margin-right: 15px;
            margin-bottom: 15px;
            .image-item {
              margin-bottom: 10px;
            }
            
            &:nth-child(3n) {
              margin-right: 0;
            }

            .image-tag {
              font-size: 12px;
              font-weight: 400;
              color: #000000;
              line-height: 17px;
              display: flex;
              align-items: center;
              justify-content: center;

              &::before {
                content: "";
                width: 13px;
                height: 13px;
                margin-right: 5px;
                background: url("~@/assets/display-task/icon_mark.png") center/100% no-repeat;
              }

              &.cancel::before {
                background-image: url("~@/assets/display-task/icon_mark_cancel.png");
              }
              &.disabled {
                color: #A6A6A6;
                &::before {
                  background-image: url("~@/assets/display-task/icon_mark_dis.png");
                }
              }


            }
          }

          .is-empty {
            font-size: 12px;
            font-weight: 400;
            color: #A6A6A6;
            line-height: 17px;
            margin-bottom: 15px;
          }
        }
      }
    }
  }
  
  .remark {
    background-color: #fff;
    padding: 15px;
    margin-bottom: 15px;

    .label {
      font-size: 14px;
      font-weight: 500;
      color: #000000;
      line-height: 20px;
    }

    &__content {
      width: 345px;
      height: 122px;
      background: #F6F5F5;
      padding: 15px;
      margin-top: 10px;
      position: relative;

      .textarea {
        background-color: transparent;
        font-size: 12px;
        width: 100%;
        height: 100%;

        &::placeholder {
          color: #A6A6A6;
        }
      }

      .maxlength {
        position: absolute;
        right: 15px;
        bottom: 10px;
        font-size: 12px;
        color: #A6A6A6;
        // background: #F6F5F5;
      }
    }
  }

  .footer {
    display: flex;
    align-items: center;
    justify-content: space-around;
    position: fixed;
    bottom: 0;
    width: 100vw;
    background-color: #fff;
    padding: 10px 20px;
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
      margin-left: 10px;

      &-icon {
        width: 20px;
        height: 20px;
        background: url("~@/assets/display-task/icon_back.png") center/100% no-repeat;
      }

      &-text {
        margin-top: 2px;
      }
    }
    .preview {
      display: flex;
      align-items: center;
      justify-content: center;
      flex-direction: column;
      font-size: 10px;
      font-weight: 400;
      color: #000000;
      line-height: 14px;
      margin-left: 30px;

      &-icon {
        width: 20px;
        height: 20px;
        background: url("~@/assets/display-task/icon_preview.png") center/100% no-repeat;
      }

      &-text {
        margin-top: 2px;
      }

      &.disabled {
        color: #A6A6A6;

        .preview-icon {
          background: url("~@/assets/display-task/icon_preview_gray.png") center/100% no-repeat;
        }
      }
    }

    .btn-group {
      display: flex;
      align-items: center;
      flex: 1;
      margin-left: 40px;

      .btn {
        flex: 1;
        background: #000000;
        color: #fff;
        font-size: 14px;
        font-weight: 400;
        line-height: 20px;
        display: flex;
        align-items: center;
        justify-content: center;
        height: 36px;

        &.disabled {
          background: #A6A6A6;
        }

        &.plain {
          background: #fff;
          color: #000000;
          border: 1px solid #000000;

          &.disabled {
            background: #fff;
            color: #A6A6A6;
            border: 1px solid #A6A6A6;
          }
        }

      }
    }

    &.stop {
      padding: 10px 40px;
      .back,
      .preview {
        margin: 0
      }
    }
  }
}
</style>