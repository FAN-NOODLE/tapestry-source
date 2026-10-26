<template>
  <div v-if="detail.task_info" class="container">
    <!-- 任务详情 -->
    <div class="detail">
      <div class="flex-center">
        <div class="title">{{detail.task_info.display_name}}</div>
        <div class="audit-status" :data-status="detail.store_info.finish_status">{{auditStatusLabel}}</div>
      </div>
      <div class="tips">
        <span>任务接收时间：{{ dateFormat(detail.task_info.created_at, 'yyyy年MM月dd日 hh:mm:ss') }}</span>
        <span v-if="show(1)" class="grade">{{ detail.store_info.grade }}</span>
      </div>
      <div v-if="detail.store_info.finish_status === 4" class="tips">
        <span>初审VM：{{ detail.store_info.special_name }}</span>
      </div>
      <div class="flex-center line">
        <span class="status" :data-status="detail.task_info.status">{{statusLabel}}</span>
        <span class="time">执行时间：{{ dateFormat(detail.task_info.start_time) }}-{{ dateFormat(detail.task_info.end_time) }} </span>
      </div>
      <div class="desc" :class="{
        show: expend,
      }">
        <span class="ellipse">{{ detail.task_info.content.length > 150 && !expend ? detail.task_info.content.slice(0, 150)+'...':detail.task_info.content }}</span>
        <span v-if="detail.task_info.content.length>150" class="more" @click="toggleExpend">{{expend?'收起':'全部'}}</span>
      </div>
      <div class="link">
        <div v-for="(item, index) in detail.task_info.class_url" :key="index" class="link-item flex-center">
          <div class="link-text">{{item.class_name}}</div>
          <div class="link-btn flex-center" @click="handleToLink(item)">去学习</div>
        </div>
      </div>
      <div v-if="show(2) && beforeStartTime<5*24*60*60*1000" class="countdown flex-center start">
        <van-count-down :time="beforeStartTime" format="距开始DD天HH时mm分" />
      </div>
      <div v-if="show(3) && beforeEndTime<5*24*60*60*1000" class="countdown flex-center">
        <van-count-down :time="beforeEndTime" format="距结束DD天HH时mm分" />
      </div>
    </div>
    <!-- VM留言 -->
    <div v-if="detail.store_info.review_remark" class="remark">
      <div class="remark-title">VM留言：</div>
      <div class="remark-content">{{ detail.store_info.review_remark }}</div>
    </div>
    <template v-if="$route.query.manager && detail.store_info.finish_status === 0">
      <empty-page 
        picture="noDisplayTask"
        paddingTop="60"
        paddingBottom="60"
        text="暂无上传内容"
      ></empty-page>
    </template>
    <!-- 默认 -->
    <template v-else>
      <!-- 所属VM -->
      <div ref="selectVm" class="select-vm m-t-15">
        <div class="select-vm__label">
          <span>选择所属VM</span>
          <span v-if="!show(5)">{{ detail.store_info.owner_member_name || '--' }}</span>
        </div>
        <!-- 未提交|未通过 -->
        <select-vm 
          v-if="show(5)" 
          v-model="detail.store_info.owner_member_id" 
          :list="vmList" 
          :disabled="show(2)"
          @change="handleVmChange">
        </select-vm>
      </div>
      <!-- 陈列内容 -->
      <div ref="content" class="content">
        <div class="label">上传陈列内容
          <span class="tips">（图片大小10M内）</span>
        </div>
        <div class="list">
          <div ref="item" class="item" v-for="(item, index) in detail.task_info.template_info" :key="index">
            <div class="item-title" :class="{required: item.is_required}">{{item.template_name}}</div>
            <div class="image-list">
              <!-- 图片列表 -->
              <image-item 
                v-for="(image, imageIndex) in item.image_list" 
                :key="imageIndex" :item="image" 
                :disabled="!show(4)"
                :disabled-delete="!!$route.query.manager || image._status === 1"
                @delete="handleDeleteImage(item, imageIndex)"
                @comment="handleComment(index, imageIndex)"
                @auto-save="handleAutoSaveImage(image, item)"
                >
              </image-item>
              <!-- 上传 -->
              <upload-image 
                v-if="show(5) && item.image_list.length<20" 
                :maxLength="20"
                v-model="item.image_list"
                :disabled="show(2)"
                @change="handleUploadChange($event, item)"/>
              <!-- 不能编辑 & 没有上传图片 -->
              <span v-if="!show(5) && item.image_list.length === 0" class="is-empty">未上传内容</span>
            </div>
          </div>
        </div>
      </div>
    </template>
    
    <!-- footer -->
    <div
      :class="[
        'footer',
        {
          stop: !show(4) && !show(2)
        }
      ]">
      <!-- 返回 -->
      <div class="back" @click="handleBack">
        <div class="back-icon"></div>
        <div class="back-text">返回</div>
      </div>
      <!-- 提示 -->
      <div v-if="show(2)" class="btn-group">
        <button class="btn" :class="{disabled: true}">执行时间开始后可上传内容</button>
      </div>
      <!-- 预览 -->
      <div v-else-if="!($route.query.manager && detail.store_info.finish_status === 0)" class="preview" :class="{disabled: false}" @click="handlePreview">
        <div class="preview-icon"></div>
        <div class="preview-text">预览长图</div>
      </div>
      <!-- 操作 -->
      <div v-if="show(4)" class="btn-group">
        <button class="btn plain" :class="{disabled: false}" @click="handleSubmit(0)">存为草稿</button>
        <button class="btn" :class="{disabled: isDisabled}" @click="handleSubmit(1)">提交</button>
      </div>
    </div>
    <!-- 删除确认 -->
    <confirm-popup ref="deletePopup" title="删除图片" @confirm="handleDeleteConfirm">
      <div>点击确定后，该图片和备注同时删除</div>
      <div class="check" :class="{checked: checked}" @click="handleCheck">不再提醒</div>
    </confirm-popup>
    <!-- 确认提交 -->
    <confirm-popup ref="popup" title="确认提交" @confirm="handleSubmitConfirm">
      <div class="text-align-left">点击确定后，将内容提交给VM【{{detail.store_info.owner_member_name}}】，审核期间内容不能修改。</div>
    </confirm-popup>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Watch } from "vue-property-decorator";
import { CountDown } from 'vant';
import SelectVm from './components/select-vm.vue';
import UploadImage from "./components/upload-image.vue";
import ImageItem from "./components/image-item.vue";
import ConfirmPopup from "./components/confirm-popup.vue";
import { auditStatusOptions, statusOptions } from "./const";
import { getTaskInfoApi, getVmListApi, saveTaskApi, getStoreTaskInfoApi, lookTaskApi, saveImageApi, deleteImageApi, getAuthApi } from '@/api/display-task'
import { getRedirectParams } from "@/utils/welearning";
import EmptyPage from "@/components/base/empty-page.vue";
import wx from 'weixin-js-sdk'
@Component({
  components: {
    [CountDown.name]: CountDown,
    SelectVm,
    UploadImage,
    ImageItem,
    ConfirmPopup,
    EmptyPage
  },
  beforeRouteLeave(to, from, next) {
    this.scrollTop = document.documentElement.scrollTop;
    next()
  }
})
export default class DisplayTaskDetail extends Vue {
  /** 临时缓存 */
  cache: {
    item: ITempItem;
    index: number;
  } = {
    item: null,
    index: -1,
  }
  /** 下次不再提醒 */
  checked: boolean = false
  /** vm列表 */
  vmList: IVmItem[] = []
  scrollTop = 0
  expend = false

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

  /** 剩余多少时间开始 */
  get beforeStartTime() {
    return this.detail.task_info.start_time * 1000 - Date.now()
  }

  /** 剩余多少时间结束 */
  get beforeEndTime() {
    return this.detail.task_info.end_time * 1000 - Date.now()
  }

  /** 是否禁用 */
  get isDisabled() {
    // 是否选了vm
    const hasVm = !!this.detail.store_info.owner_member_id
    // 没有素材
    const isEmpty = this.detail.task_info.template_info.every(item => item.image_list.length === 0)
    // 必传素材没有传
    const hasRequiredEmpty = this.detail.task_info.template_info.some(item => item.is_required && item.image_list.length === 0)
    
    return !hasVm || isEmpty || hasRequiredEmpty
  }

  get isEmptyIndex() {
    const notRequired = this.detail.task_info.template_info.every(item => !item.is_required)
    const isEmpty = this.detail.task_info.template_info.every(item => item.image_list.length === 0)
    const index = this.detail.task_info.template_info.findIndex(item => item.is_required && item.image_list.length === 0)
    if(isEmpty) {
      return notRequired ? 0 : index
    }
    return index
  }

  /** 控制页面区域展示 */
  get show() {
    const { finish_status } = this.detail.store_info
    const { status } = this.detail.task_info
    // 来自数据看板
    
    return (key: number) => {
      const isManager = this.$route.query.manager
      if(isManager) {
        return false
      }
      switch (key) {
        case 1:
          // 任务审核已通过 显示评分
          return finish_status === 3
        case 2:
          // 任务未开始 不能操作
          return status === 0
        case 3:
          // 任务进行中 显示进行中倒计时
          return status === 1
        case 4:
          // 任务进行中 && 任务待提交 || 任务审核未通过 可以进行操作
          return status == 1 && finish_status === 0 || finish_status === 2
        case 5:
          // 任务未开始 || 任务进行中 && 任务待提交 || 任务审核未通过 可以进行操作
          return status === 0 || status === 1 && finish_status === 0 || finish_status === 2
        default:
          return false
      }

    }
  }


  /** 用户信息 */
  get userInfo(): IUserInfo {
    return this.$store.getters["user/getUserInfo"] || {}
  }

  @Watch("detail", { deep: true })
  onDetailChange(value: IDisplayTaskDetail) {
    if(!this.hasUploading()) {
      this.$store.commit("displayTask/setDetail", value)
    }
  }

  async mounted() {
    if(this.userInfo.app_source !== 1) {
      await this.getAuth()
    }
    this.getVmList()
    this.$toast.loading({
      message: "加载中...",
      forbidClick: true,
    });
    await this.getStoreTaskInfo()
    await this.getDetail()
    this.$toast.clear()

    this.checked = localStorage.getItem('delete_image_dialog') === '1'
    if(this.$route.query.from !== '1') {
      // 神策数据上报 -- 进入VM
      this.$track.trackEnterDisplayTask({
        vm_entrance: 2
      })
    }

    // 神策数据上报 -- 查看陈列任务
    this.$track.trackViewDisplayTask({
      dispalytask_title: this.detail.task_info.display_name,
      displaytask_complete_state: this.detail.store_info.finish_status,
      displaytask_state: this.detail.task_info.status,
      display_store: this.detail.store_info.store_code,
    })

    console.log(11111, this.detail.task_info);
    
  }

  async activated() {
    if(!this.detail.store_info) {
      this.scrollTop = 0
      await this.getStoreTaskInfo()
      await this.getDetail()
    }

    this.$nextTick(() => {
      if (this.scrollTop > 0) {
        console.log('activated', this.scrollTop);
  
        document.documentElement.scrollTo({
          top: this.scrollTop,
          behavior: "auto"
        });
      }
    })

    if(this.$route.query.from !== '1') {
      window.addEventListener("popstate", this.popstate, false);
    }
  }

  deactivated() {
    if(this.$route.query.from !== '1') {
      window.removeEventListener("popstate", this.popstate, false);
    }
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

  toggleExpend() {
    this.expend = !this.expend;
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
      // 数据格式化
      data.template_info.forEach(item => {
        item.image_list = this.detail.store_info.image_list
          .filter(val => val.identifier === item.template_name).map(val => {
            return {
              ...val,
              _upload_status: '',
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
      // data.status = 2
      this.$store.commit("displayTask/setDetail", {
        ...this.detail,
        task_info: data
      })
      this.$store.commit('displayTask/setUpdate', false);
    } catch (error) {
      console.log(error);
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
      
      this.$store.commit("displayTask/setDetail", {
        ...this.detail,
        store_info: data,
      })

    } catch (error) {
      console.log(error);
      
    }
  }

  /**
   * @description 获取VM列表
   */
  async getVmList() {
    try {
      const { code, message, data }: IApiResult<IVmItem[]> = await getVmListApi()
      if(code!== 200) {
        this.$toast.fail(message)
        return
      }
      this.vmList = data
    } catch (error) {
      console.log(error);
    }
  }

  /**
   * @description 去学习
   * @param item 
   */
  async handleToLink(item: IClassUrlItem) {
    const redirect = getRedirectParams(item.url)
    const uri = await redirect.uri()
    // 神策数据上报 -- 课程学习
    this.$track.trackCourseLearing({
      course_title: item.class_name
    })
    console.log(uri);
    
    location.href = uri;
  }

  /**
   * @description 选择所属VM
   */
  handleVmChange(value: number) {
    const current = this.vmList.find(item => item.member_id === value);
    this.detail.store_info.owner_member_name = current.member_name
    this.detail.store_info.owner_member_user_id = current.user_id
  }

  /**
   * @description 上传图片
   * @param value 
   * @param index 
   */
  handleUploadChange(value: any, item: ITempItem) {
    // console.log(11111, value);
    
    // 格式化数据
    item.image_list = value
    // this.debounceSave()
  }

  /**
   * @description 删除图片
   */
  handleDeleteImage(item: ITempItem, index: number) {
    if(!this.checked) {
      (this.$refs.deletePopup as any).showPopup();
      this.cache.item = item
      this.cache.index = index
    } else {
      
      this.handleAutoDeleteImage(item, index)
    }
  }

  /**
   * @description 删除图片确认
   */
  handleDeleteConfirm() {
    this.handleAutoDeleteImage(this.cache.item, this.cache.index)
    // this.cache.item.image_list.splice(this.cache.index, 1);
    // 不再提醒
    if(this.checked) {
      localStorage.setItem('delete_image_dialog', '1')
    }
  }

  /**
   * @description 取消删除
   */
  handleDeleteCancel() {
    this.cache.item = null
    this.cache.index = -1
  }

  /** 
   * @description 预览图片/添加备注
   */
  async handleComment(idx: number, index: number) {
    if(this.$route.query.manager === '1'){
      this.$router.push({
        path: 'display-task-comment',
        query: {
          idx: idx+'',
          index: index+'',
          is_vm: '2'
        }
      })
      return
    }
    this.$router.push({
      path: 'display-task-comment',
      query: {
        idx: idx+'',
        index: index+'',
        is_vm: '0'
      }
    })
    const current = this.detail.task_info.template_info[idx].image_list[index]

    if(current.is_look === 1) {
      const latest = current.comment.filter(item => item.is_vm === 1).map(item => item.id);
      try {
        const { code, message, data }: IApiResult<any> = await lookTaskApi({
          allot_id: this.detail.store_info.allot_id,
          image_id: current.id,
          comment_id: Math.max(...latest),
        })
        if(code!== 200) {
          this.$toast.fail(message)
          return
        }
        current.is_look = 0
        
      } catch (error) {
        console.log(error);
      }

    }
  }

  /**
   * @description 下次不再提醒
   */
  handleCheck() {
    this.checked = !this.checked
  }

  /**
   * @description 预览长图
   */
  handlePreview() {
    // if(this.isDisabled) return
    this.$router.push({
      path: 'display-task-preview',
    })
  }

  /**
   * @description 返回
   */
  handleBack() {
    this.$router.back()
  }

  /**
   * @description 提交
   * @param show_type 
   */
  handleSubmit(show_type: number) {
    if(show_type === 1) {
      // 校验
      if(this.isDisabled) {
        if(!this.detail.store_info.owner_member_id) {
          this.$toast.fail('请选择所属VM');
          (this.$refs.selectVm as HTMLElement).scrollIntoView({ behavior: 'smooth' })
        } else if(this.isEmptyIndex === 0) {
          this.$toast.fail('请完成必传模块内容上传后提交');
          (this.$refs.content as HTMLElement).scrollIntoView({ behavior: 'smooth' });
        } else if(this.isEmptyIndex > 0) {
          this.$toast.fail('请完成必传模块内容上传后提交');
          (this.$refs.item[this.isEmptyIndex] as HTMLElement).scrollIntoView({ behavior: 'smooth' });
        }
        return 
      }
      (this.$refs.popup as any).showPopup();
    } else {
      this.save(0)
    }
  }
  /**
   * @description 提交确认
   */
  handleSubmitConfirm() {
    this.save(1)
  }

  /**
   * 保存数据
   * @param submit_type  0 保存 1 提交
   */
  async save(submit_type: number, toast: boolean = true) {
    const { store_info, task_info } = JSON.parse(JSON.stringify(this.detail))
    
    if(this.hasUploading()) {
      this.$toast.fail('请等待图片上传完成')
      return
    }
    const params = {
      id: task_info.id,
      allot_id: store_info.allot_id,
      owner_member_id: store_info.owner_member_id,
      owner_member_name: store_info.owner_member_name,
      owner_member_user_id: store_info.owner_member_user_id,
      submit_type: submit_type,
      image_list: [],
    }
    params.image_list = task_info.template_info.reduce((prev, item) => {
      return prev.concat(item.image_list.map(image => {
        return {
          id: image.id || 0,
          identifier: image.identifier || item.template_name,
          url: image.url || '',
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
    
    try {
      const { code, message }: IApiResult<any> = await saveTaskApi(params)
      if(code!== 200) {
        this.$toast.fail(message)
        return
      }
      // 更新列表
      if(submit_type === 1) {
        this.$store.commit('displayTask/setUpdateList', {
          update: true,
          current: 1
        })
      }
      if(submit_type === 0) {
        toast && this.$toast.success('保存成功');
      } else {
        // 神策数据上报 -- 提交陈列任务
        this.$track.trackSubmitDisplayTask({
          dispalytask_title: task_info.display_name,
          displaytask_complete_state: store_info.finish_status,
          displaytask_state: task_info.status,
          is_resubmit: store_info.finish_status === 0 ? 0 : 1
        })
        this.$router.back()
      }
    } catch (error) {
      console.log(error);
    }
    
  }

  /**
   * 是否有正在上传的图片
   */
  hasUploading() {
    if(!this.detail.task_info) return false
    return this.detail.task_info.template_info.some(item => {
      return item.image_list.some(image => {
        return image._upload_status === 'uploading' || !image.url || !image.id;
      })
    })
  }

  /**
   * 图片上传成功后自动保存记录
   * @param image 
   * @param item 
   */
  handleAutoSaveImage(image:IImageItem, item: ITempItem ) {
    if(!image.id) {
      const { id } = this.detail.task_info
      const { allot_id } = this.detail.store_info
      const params = {
        display_task_id: id,
        allot_id,
        url: image.url,
        id: 0,
        identifier: item.template_name
      }
      saveImageApi(params).then((res: IApiResult)=> {
        if(res.code === 200) {
          // image.id = res.data.id
          image._upload_status = 'done';
          image._message = '上传成功';
          this.$set(image, 'id', res.data.id)
        }
      })
    }
  }

  /**
   * 删除图片
   * @param item 
   * @param index 
   */
  handleAutoDeleteImage(item: ITempItem, index: number) {
    const image = item.image_list[index];
    item.image_list.splice(index, 1);
    deleteImageApi({
      id: image.id
    })
  }

}
</script>
<style lang="scss" scoped>
.check {
  margin-top: 15px;
  color: #777;
  display: flex;
  justify-content: center;
  align-items: center;

  &::before {
    content: "";
    display: block;
    width: 16px;
    height: 16px;
    border: 1px solid #000000;
    border-radius: 50%;
    margin-right: 10px;
  }

  &.checked::before {
    background: url("~@/assets/display-task/icon_checked.png") center/28px 28px no-repeat;
  }
}
.text-align-left {
  text-align: left;
}
.container {
  background-color: #F7F7F7;
  min-height: 100vh;
  overflow: hidden;
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
      position: relative;

      
      .more {
        color: #A1A1A1;
        font-size: 12px;
        font-weight: 400;
        line-height: 17px;
        text-align: right;
        cursor: pointer;
        transition: all .3s;
        padding-left: 10px;
        &::after {
          content: "";
          margin-left: 5px;
          display: inline-block;
          width: 14px;
          height: 9px;
          background: url('~@/assets/display-task/icon_down.png') right center/14px 9px no-repeat;
        }
      }

      &.show .more::after {
        transform: rotate(180deg);
      }
    
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
    margin-bottom: 15px;

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

  .select-vm {
    background-color: #fff;
    margin-bottom: 15px;
    padding: 15px;
    position: relative;
    
    &__label {
      font-size: 14px;
      font-weight: 500;
      color: #000000;
      line-height: 20px;

      &::before {
        content: "*";
        color: #E02020;
      }

      span:nth-child(2) {
        position: absolute;
        top: 50%;
        right: 15px;
        transform: translateY(-50%);
        font-size: 12px;
        font-weight: 400;
        color: #000000;
        line-height: 17px;
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

      &::before {
        content: "*";
        color: #E02020;
      }
    }

    .tips {
      font-size: 12px;
      font-weight: 400;
      color: #A6A6A6;
      line-height: 17px;
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
      padding-bottom: calc(10px + constant(safe-area-inset-bottom)); /** iOS 11.0+ */
      padding-bottom: calc(10px + env(safe-area-inset-bottom)); /** iOS 11.0+ */
    
      .back,
      .preview {
        margin: 0
      }
    }
  }
}
</style>