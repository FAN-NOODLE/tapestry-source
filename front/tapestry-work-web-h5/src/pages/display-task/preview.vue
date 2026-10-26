<template>
  <div v-if="detail.store_info" class="container">
    <!-- 任务详情 -->
    <div class="detail">
      <div class="flex-center">
        <div class="title">门店：{{detail.store_info.store_name}}_{{detail.store_info.store_code}}</div>
        <div v-if="detail.store_info.finish_status > 0" class="audit-status" :data-status="detail.store_info.finish_status">{{auditStatusLabel}}</div>
      </div>
      <div class="tips">
        <span>提交时间：{{ dateFormat(detail.store_info.submit_time||(Date.now()/1000), 'yyyy年MM月dd日 hh:mm:ss') }}</span>
        <span v-if="detail.store_info.finish_status === 3" class="grade">{{ detail.store_info.grade }}</span>
      </div>
    </div>
    <!-- VM留言 -->
    <div v-if="detail.store_info.review_remark" class="remark">
      <div class="remark-title">VM留言：</div>
      <div class="remark-content">{{ detail.store_info.review_remark }}</div>
    </div>
    <!-- 陈列内容 -->
    <div class="content">
      <div class="title">陈列内容</div>
      <div class="section" v-for="(item, index) in detail.task_info.template_info" :key="index">
        <div class="section-title">{{ item.template_name }}</div>
        <div v-if="item.image_list.length>0" class="image-list">
          <div class="image-item" v-for="(image, index) in item.image_list" :key="index">
            <!-- 图片 -->
            <div class="image-item__inner">
              <van-image :src="image.url" lazy-load></van-image>
              <div v-if="image._status === 1" class="image-item__fail"></div>
            </div>
            <!-- 评论 -->
            <comment-item 
              v-for="(comment, commentIndex) in image.comment" 
              :key="commentIndex" 
              :item="comment" 
              :is_vm="is_vm"
              disabled>
            </comment-item>
          </div>
        </div>
        <div v-if="item.image_list.length === 0" class="is-empty">未上传内容</div>
      </div>
    </div>
    <!-- footer -->
    <div class="footer">
      <div class="btn" @click="handleBack">返回</div>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { auditStatusOptions, statusOptions } from "./const";
import CommentItem from "./components/comment-item.vue";
import { Image as VanImage } from "vant";
import { getStoreTaskInfoApi, getTaskInfoApi } from "@/api/display-task";
@Component({
  components: {
    CommentItem,
    VanImage,
  }
})
export default class DisplayTaskPreview extends Vue {
  is_vm = 0
  timer = null
  /** 任务详情 */
  get detail(): IDisplayTaskDetail {
    return this.$store.getters["displayTask/getDetail"]
  }

  /** 任务状态 */
  get statusLabel() {
    return statusOptions.find((item) => item.value === this.detail.task_info.status).label
  }

  /** 审核状态 */
  get auditStatusLabel() {
    return auditStatusOptions.find((item) => item.value === this.detail.store_info.finish_status).label
  }

  async mounted() {
    const { is_vm, other } = this.$route.query
    this.is_vm = is_vm ? Number(is_vm) : 0
    if(!this.detail.store_info) {
      await this.getStoreTaskInfo()
    }
    if(!this.detail.task_info) {
      await this.getDetail()
    }
    
    if(other) {
      // 神策数据上报 -- 查看他店陈列图
      this.$track.trackViewStoreDisplay({
        dispalytask_title: this.detail.task_info.display_name,
        displaytask_complete_state: this.detail.store_info.finish_status,
        displaytask_state: this.detail.task_info.status,
        display_store: this.detail.store_info.store_code,
      })
    }
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
   * @description 格式化时间
   * @param time 
   * @param format 
   */
  dateFormat(time: number, format = "yyyy/MM/dd") {
    return this.$helper.dateFormat(time, format);
  }

  /**
   * @description 返回
   */
  handleBack() {
    this.timer && clearTimeout(this.timer)
    this.timer = setTimeout(() => {
      this.$router.back();
    }, 10);
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
    .title {
      font-size: 16px;
      font-weight: 600;
      color: #000000;
      line-height: 22px;
      display: flex;
      padding-right: 15px;
      flex: 1;
      // overflow: hidden;
      // text-overflow: ellipsis;
      // white-space: nowrap;
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
      margin-right: 5px;
      &[data-status="0"] {
        color: #000;
        background: rgba(0,0,0,0.1);
      }
      &[data-status="1"] {
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
        white-space: nowrap;
        &::before {
          content: "评分:";
          font-size: 12px;
          font-weight: 400;
          margin-right: 2px;
          padding-top: 5px;
          white-space: nowrap;
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

  .remark {
    background: rgba(0,0,0,0.1);
    box-shadow: 0px 3px 17px 0px rgba(54,68,58,0.06);
    padding: 15px 20px;
    margin-top: 15px;

    &-title {
      font-size: 14px;
      font-weight: 500;
      color: #000000;
      line-height: 20px;
    }
    &-content {
      padding: 10px 10px 5px 10px;
      font-size: 12px;
      font-weight: 400;
      color: #000000;
      line-height: 17px;
      text-align: justify;
      word-break: break-all;
    }
  }

  .content {
    background-color: #fff;
    padding: 15px 20px;
    margin-top: 15px;
    .title {
      font-size: 14px;
      font-weight: 500;
    }

    .section {
      margin-top: 15px;

      .section-title {
        font-size: 12px;
        font-weight: 400;
        color: #000000;
        line-height: 17px;
      }

      .image-item {
        margin-top: 15px;
        .van-image {
          background-color: #F7F7F7;
          display: block;
        }

        .comment-item {
          margin-top: 15px;
          margin-bottom: 0;
        }

        &__inner {
          position: relative;
        }

        &__fail {
          position: absolute;
          left: 0;
          top: 0;
          width: 55px;
          height: 55px;
          background: url("~@/assets/display-task/icon_fail_big.png") center/100% no-repeat;
        }
      }

      .is-empty {
        font-size: 12px;
        font-weight: 400;
        color: #A6A6A6;
        line-height: 17px;
        margin-top: 15px;
      }
    }
  }

  .footer {
    position: fixed;
    bottom: 0;
    width: 100vw;
    background-color: #fff;
    padding: 10px 30px;
    padding-bottom: calc(10px + constant(safe-area-inset-bottom)); /** iOS 11.0+ */
    padding-bottom: calc(10px + env(safe-area-inset-bottom)); /** iOS 11.0+ */
    z-index: 100;

    .btn {
      font-size: 14px;
      font-weight: 400;
      line-height: 20px;
      display: flex;
      align-items: center;
      justify-content: center;
      height: 36px;
      background: #fff;
      color: #000000;
      border: 1px solid #000000;
    }

  }
}
</style>