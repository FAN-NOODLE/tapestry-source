<template>
  <layout-form v-if="formInfo">
    <template slot="content">
      <audit-form v-if="!isInit" v-model="formInfo" @change="handleChange"></audit-form>
      <!-- 审核确认弹框 -->
      <el-dialog
        title="提示"
        :visible.sync="dialogVisible"
        width="400px">
        <span class="el-icon-warning m-r-5"></span>
        <span>确认保存当前操作，并通过该门店提交的内容？</span>
        <el-form ref="form" :model="formInfo" :rules="rules" class="m-t-10" label-width="60px">
          <el-form-item label="评分" prop="grade">
            <el-input 
              v-model.trim="formInfo.grade"
              size="medium"
              placeholder="请输入0~100的整数"
            ></el-input>
          </el-form-item>
        </el-form>
        <span slot="footer" class="dialog-footer">
          <el-button @click="handleSubmitCancel">取 消</el-button>
          <el-button type="primary" @click="handleSubmitConfirm()">确 定</el-button>
        </span>
      </el-dialog>
    </template>
    <template slot="footer">
      <div>
        <el-button @click="handleCancel">取消</el-button>
        <el-button @click="handleSubmit(1)">不通过</el-button>
        <el-button type="primary" @click="handleSubmit(2)" :disabled="markCount>0">通过</el-button>
      </div>
    </template>
  </layout-form>
</template>
<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import Api from '@/api/wisdom/displayTask';
import AuditForm, { IStoreInfo, ICommentItem, IImageItem } from '../components/audit-form.vue'
import UserModule from "@/store/modules/user";
import { IFormInfo } from '../components/edit-form.vue';

@Component({
  components: {
    AuditForm
  },
  beforeRouteLeave (to, from, next) {
    this.leaveConfirm(() => {
      next()
    })
  }
})
export default class DisplayTaskAudit extends Vue {
  taskDetail: IFormInfo = null
  /** 表单数据 */
  formInfo: IStoreInfo = null
  /** 表单校验规则 */
  rules = {
    grade: [
      {
        required: true,
        trigger: 'change',
        validator(rule: any, value: string, callback: Function) {
          const regex = /^(0|[1-9][0-9]{0,1}|100)$/;
          if (!value) {
            callback(new Error('请输入评分'))
          } else if (!regex.test(Number(value)+'')) {
            callback(new Error('请输入0~100之间的整数'))
          } else {
            callback()
          }
        }
      }
    ]
  }
  /** 提交确认 */
  dialogVisible = false
  /** 是否初始化 */
  isInit = true
  /** 内容是否改变 */
  isChanged = false

  defaultAvatar = 'https://cdn-cos.coach.com.cn/uploads/common/icon_logo.png'
  get markCount() {
    return this.formInfo.template_info.reduce((prev, item)=> {
      prev += item.image_list.filter(image => image.status === 0 && image._status === 1).length

      return prev
    }, 0)
  }

  get userInfo() {
    return UserModule.getUserInfo
  }

  async mounted() {
    await this.getTaskDetail()
    this.getDetail()
  }

  getDetail() {
    Api.getAllotInfoApi({
      id: Number(this.$route.query.id),
      store_id: Number(this.$route.query.store_id)
    }).then((res: any) => {
      if(res.code === 200) {
        console.log(res.data);
        
        res.data.template_info = this.taskDetail.template_info.map((item) => {
          const imageList = res.data.image_list
            .filter(image => image.identifier === item.template_name)
            .map(value => {
              value._status = value.status
              value._image_list = []
              value.comment = value.comment.map(comment => {
                const images = comment.content.image_url?.map(url => ({
                  url: url as unknown as string,
                }))||[]
                return {
                  ...comment,
                  _reply_content: '',
                  _reply_name: '',
                  _reply_show: false,
                  _prefix_width: 0,
                  content: images.length ? {image_url: images} : {text: comment.content.text}
                }
              })
              return value
            })
          return {
            ...item,
            image_list: imageList
          }
        })
        res.data.grade = res.data.grade || null
        this.formInfo = res.data
        this.isInit = false
      }
    })
  }

  async getTaskDetail() {
    const res: any = await  Api.getDisplayTaskDetailApi({
      id: Number(this.$route.query.id)
    })
    if (res.code == 200) {
      this.taskDetail = res.data
    } else {
      this.$message.error(res.message);
    }
  }

  /** 
   * 表单变化
   */
  handleChange() {
    this.isChanged = true
  }

  /**
   * @description 取消
   */
  handleCancel() {
    this.leaveConfirm(() => {
      this.isChanged = false
      this.$router.back()
    })
  }

  leaveConfirm(callback: Function) {
    if(this.isChanged) {
      const message = '你确认要放弃编辑吗？放弃后编辑的信息将不会保存。'
      const title = '离开页面提示'
      this.$confirm(message, title, {
        confirmButtonText: "确认",
        cancelButtonText: "取消"
      }).then(() => {
        callback && callback()
      })
    } else {
      callback && callback()
    }
    
  }

  /**
   * @description 提交
   */
  handleSubmit(type: number) {
    console.log(type);
    // 通过
    if(type == 2) {
      this.dialogVisible = true
      return
    } 
    // 不通过
    const message = '确认保存当前操作，并不通过该门店提交的内容？'
    const title = '提示'
    this.$confirm(message, title, {
      confirmButtonText: "确认",
      cancelButtonText: "取消"
    }).then(() => {
      this.submit(1);
    })
  }

  /**
   * @description 提交确认
   */
  handleSubmitConfirm() {
    (this.$refs.form as any).validate((valid) => {
      if (valid) {
        this.submit(2)
      } else {
        return false;
      }
    });
  }

  /**
   * @description 提交确认
   * @param type 1:不通过 2:通过
   */
  handleSubmitCancel() {
    this.dialogVisible = false
    this.$message.success('取消提交')
  }
  /**
   * @description 提交
   * @param type 1:不通过 2:通过
   */
  submit(type: number) {
    this.dialogVisible = false
    const { id, allot_id, review_remark, grade } = this.formInfo
  
    const params = {
      id: id,
      allot_id: allot_id,
      review_remark: review_remark.trim(),
      grade: Number(grade),
      review_type: type,
      image_list: []
    }
    params.image_list = this.formInfo.template_info.reduce((prev, item) => {
      return prev.concat(item.image_list.map(image => {
        // 处理图片备注
        const commentList = image.comment.map(comment => {
          const images = comment.content.image_url?.map(val => val.url) ??[]
          const res = {
            ...comment,
            content: images.length ? {image_url: images} : {text: comment.content.text}
          }
          this.removeInfo(res)
          return res
        })

        if(image._image_list?.length) {
          Object.entries(image._image_list.reduce((prev, it) => {
            // this.removeInfo(val)
            if(prev[it._comment_time]) {
              prev[it._comment_time].push(it.url)
            } else {
              prev[it._comment_time] = [it.url]
            }
            return prev
          }, {})).map(([key, value]) => {
            const {account_name, admin_id} = this.userInfo
            const comment = {
              id: 0,
              member_id: admin_id,
              member_user_id: '',
              member_name: account_name,
              avatar: this.defaultAvatar,
              store_code: this.formInfo.store_code,
              action: 0,
              is_vm: 1,
              is_self: 1,
              comment_time: Number(key),
              content: {
                image_url: value as []
              },
              reply_name: '',
            }
            commentList.push(comment)
          })
        }
        // 排序
        commentList.sort((a, b) => {
          return a.comment_time - b.comment_time
        })

        const result = {
         ...image,
          status: image._status,
          comment: commentList
        }

        this.removeInfo(result)
        return result
      }))
    }, [])

    Api.saveAllotApi(params).then((res: any) => {
      if(res.code === 200) {
        this.$message.success('审核完成')
        this.isChanged = false
        this.$router.back()
      } else {
        this.$message.error(res.message)
      }
    })
  }

  /** 
   * 移除内部变量
   */
  removeInfo(data: Record<string, unknown>) {
    for (const key in data) {
      if (Object.prototype.hasOwnProperty.call(data, key)) {
        if(key.startsWith('_')) {
          delete data[key]
        }
      }
    }
  }

  /**
   * 创建备注
   * @param value 
   * @param reply_name 
   */
   createComment(value: string | IImageItem[], reply_name = ''): ICommentItem {
    const {account_name, admin_id} = this.userInfo
    const comment: ICommentItem = {
      id: 0,
      member_id: admin_id,
      member_user_id: '',
      member_name: account_name,
      avatar: 'https://img-crs.vchangyi.com/store-manage-wxapp/sy_register_photo.png',
      store_code: this.formInfo.store_code,
      action: 0,
      is_vm: 1,
      is_self: 1,
      comment_time: parseInt(Date.now()/1000+ ''),
      content: {},
      reply_name: reply_name,
      _edit_content: '',
      _edit_show: false
    }
    if(typeof value ==='string') {
      comment.content.text = value
    } else {
      comment.content.image_url = value
    }
    return comment
  }
}
</script>
<style lang="scss" scoped>
.section {
  &:not(:last-child) {
    margin-bottom: 20px;
  }

  &-title {
    position: relative;
    padding-left: 10px;
    font-size: 16px;
    margin-bottom: 10px;
    &::before {
      content: "";
      width: 5px;
      height: 20px;
      position: absolute;
      top: 50%;
      left: 0;
      transform: translateY(-50%);
      background-color: #000;
    }
  }

  &-row {
    &:not(:last-child) {
      margin-bottom: 10px;
    }
  }

  .btn {

    border: none;
    background-color: #f5f6fa;
    /deep/.el-icon-picture,
    /deep/.el-icon-warning-outline {
      font-size: 16px;
      vertical-align: middle;
    }
  }
}

.upload {
  position: relative;
  .btn-tag {
    position: absolute;
    top: 0;
    left: 105px;
  }
}

.item {
  &:not(:last-child) {
    margin-bottom: 20px;
  }
  .item-user {
    display: flex;
    align-items: center;
    .el-avatar {
      margin-right: 10px;
    }

    .item-user-info {
      line-height: normal;
      flex: 1;
    }
  
    .item-user-name {
      font-weight: bold;
    }
    .item-user-store {
      font-size: 12px;
      color: #999;
      margin-top: 5px;
    }

    .item-user-time {
      color: #999;
      line-height: normal;
    }

    
  }
  .item-content {
    margin-top: 10px;
    .item-text {
      padding: 10px;
      background-color: #f7f7f7;
      line-height: 1.5;
      color: #666;
      position: relative;
      text-align: justify;
      border-radius: 4px;
      .reply {
        height: 20px !important;
        line-height: 20px !important;
        padding: 0 5px!important;  
      }
    }
  }
}
</style>