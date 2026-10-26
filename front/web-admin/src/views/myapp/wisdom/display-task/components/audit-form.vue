<template>
  <el-form label-width="100px">
    <el-form-item label="上传内容：">
      <div class="section" v-for="(item, index) in value.template_info" :key="index">
        <div class="section-title">{{item.template_name}}</div>
        <div v-if="item.image_list.length" class="section-content">
          <!-- 循环展示 -->
          <el-row class="section-row" :gutter="20" v-for="(imageItem, imageIndex) in item.image_list" :key="imageIndex">
            <el-col :span="6">
              <div class="image">
                <el-image fit="cover" :src="imageItem.url" :preview-src-list="value.image_list.map(val=> val.url)"></el-image>
                <div v-if="imageItem._status === 1" class="image-tag"></div>
              </div>
            </el-col>
            <el-col :span="18">
              <!-- 我的意见 -->
              <el-form-item v-if="!readonly" label="我的意见：">
                <div class="m-b-10 upload">
                  <upload-image-button
                    class="upload-template"
                    accept="image/jpg,image/jpeg,image/png"
                    :urls="[]"
                    img-type="batch"
                    :before-upload="handleBeforeUpload"
                    @update="handleImageChange($event, imageItem)"
                    :limit="10"
                    multiple
                    size="small"
                  >
                    <el-button class="btn" icon="el-icon-picture" @click="handleAddImage(imageItem)">+图片</el-button>
                  </upload-image-button>
                  <div class="btn-group">
                    <el-button class="btn btn-tag" :disabled="imageItem.status === 1" :icon="imageItem._status===1?'el-icon-close':'el-icon-warning-outline'" @click="handleTag(imageItem)">{{imageItem._status===1?'取消标记':'标记不通过'}}</el-button>
                  </div>
                  <div class="upload-image-list">
                    <div class="upload-image-item" v-for="(image, imageIndex) in imageItem._image_list" :key="imageIndex">
                      <el-image fit="cover" :src="image.url" :preview-src-list="imageItem._image_list.map(item => item.url)"></el-image>
                      <div class="upload-image-delete" @click="handleDeleteImage(imageItem, imageIndex)"><i class="el-icon-close"></i></div>
                    </div>
                  </div>
                </div>
                <div class="m-b-20">
                  <el-input
                    v-model="imageItem._remark"
                    size="medium"
                    :rows="4"
                    type="textarea"
                    maxlength="200"
                    show-word-limit
                    placeholder="请输入修改意见"
                  >
                  </el-input>
                  <div v-if="imageItem._remark" class="right">
                    <el-button type="primary" size="mini" @click="handleSaveComment(imageItem)">保存</el-button>
                  </div>
                </div>
              </el-form-item>
              <!-- 店铺反馈 -->
              <el-form-item label="店铺备注：">
                <span v-if="imageItem.comment.length == 0">无</span>
                <!-- 列表 -->
                <div class="item" v-for="(commentItem, commentIndex) in imageItem.comment" :key="commentIndex">
                  <div class="item-user">
                    <el-avatar :src="commentItem.avatar||defaultAvatar"></el-avatar>
                    <div class="item-user-info">
                      <span class="item-user-name">{{commentItem.member_name}}</span>
                      <span v-if="commentItem.store_code && commentItem.is_vm == 0" class="item-user-store">{{commentItem.store_code||'--'}}</span>
                      <span v-if="commentItem.is_self" class="item-user-tag">我</span>
                      <span v-else-if="commentItem.is_vm" class="item-user-tag">VM</span>
                    </div>
                    <div class="item-user-time">{{ commentItem.comment_time | formatDate("yyyy-mm-dd hh:mm:ss") }}</div>

                  </div>
                  <div class="item-content">
                    <div class="item-text">
                      <div v-if="commentItem.content.image_url" class="item-text__image">
                        <el-image fit="cover" v-for="(val, i) in commentItem.content.image_url" :key="i" :src="val.url" :preview-src-list="commentItem.content.image_url.map(item => item.url)"></el-image>
                      </div>
                      <div><span v-if="commentItem.reply_name">回复 <span class="reply-name">{{ commentItem.reply_name }}：</span></span>{{ commentItem.content.text||'' }}</div>
                      <div v-if="!readonly" class="item-text__footer">
                        <!-- 回复备注 -->
                        <el-popover
                          v-if="!commentItem.is_self"
                          width="400"
                          placement="bottom"
                          trigger="manual"
                          v-model="commentItem._reply_show">
                          <div 
                            class="input" 
                            :style="{
                              '--text-indent': commentItem._prefix_width
                            }">
                            <el-input
                              size="medium"
                              v-model="commentItem._reply_content"
                              :rows="5"
                              type="textarea"
                              maxlength="200"
                              show-word-limit
                              placeholder="请输入回复内容"
                            >
                            </el-input>
                            <span :ref="'inputPrefix'+commentIndex" class="input-prefix" v-if="commentItem._reply_name">回复 <span class="reply-name">{{ commentItem._reply_name }}：</span></span>
                          </div>
                          <div class="m-t-10" style="text-align: right;">
                            <el-button size="mini" @click="handleReplyCancel(commentItem)">取消</el-button>
                            <el-button type="primary" size="mini" @click="handleReplySave(commentItem, imageItem)">保存</el-button>
                          </div>
                          <el-button slot="reference" class="reply" type="text" size="mini" icon="el-icon-s-comment" @click="handleReply(commentItem, commentIndex)">回复</el-button>
                        </el-popover>
                        <!-- 删除备注 -->
                        <el-popconfirm
                          v-if="commentItem.is_self && !commentItem.id"
                          confirm-button-text='确定'
                          cancel-button-text='取消'
                          icon="el-icon-info"
                          icon-color="#000"
                          title="确定删除吗？"
                          @onConfirm="handleDeleteComment(imageItem, commentIndex)">
                          <el-button slot="reference" class="reply" type="text" size="mini" icon="el-icon-delete">删除</el-button>
                        </el-popconfirm>
                        <!-- 编辑备注 -->
                        <el-popover
                          v-if="commentItem.is_self && !commentItem.id"
                          width="400"
                          placement="bottom"
                          trigger="manual"
                          v-model="commentItem._edit_show">

                          <div class="input">
                            <el-input
                              size="medium"
                              v-model="commentItem._edit_content"
                              :rows="5"
                              type="textarea"
                              maxlength="200"
                              show-word-limit
                              placeholder="请输入回复内容"
                            >
                            </el-input>
                            <span class="input-prefix" v-if="commentItem.reply_name">回复 <span class="reply-name">{{ commentItem.reply_name }}：</span></span>
                          </div>
                          <div class="m-t-10" style="text-align: right;">
                            <el-button size="mini" @click="handleEditCancel(commentItem)">取消</el-button>
                            <el-button type="primary" size="mini" @click="handleEditSave(commentItem)">保存</el-button>
                          </div>
                          <el-button slot="reference" class="reply" type="text" size="mini" icon="el-icon-edit" @click="handleEdit(commentItem)">编辑</el-button>
                        </el-popover>
                      </div>
                    </div>
                  </div>
                </div>
              </el-form-item>
            </el-col>
          </el-row>
        </div>
        <div v-else class="section-content">
          <span style="color:#999">未上传图片</span>
        </div>
      </div>
    </el-form-item>
    <el-form-item label="整体评价：">
      <el-input
        size="medium"
        v-model="value.review_remark"
        :rows="7"
        type="textarea"
        maxlength="500"
        show-word-limit
        placeholder="请输入整体评价"
        :readonly="readonly"
      >
      </el-input>
    </el-form-item>
    <el-form-item v-if="value.finish_status === 3" label="评分：">
      {{ value.grade }}
    </el-form-item>
  </el-form>
</template>
<script lang="ts">
import { Component, Vue, Prop, Watch, Model } from 'vue-property-decorator';
import UploadImageButton from "@/components/business/upload-image/upload-image-button.vue";
import UserModule from "@/store/modules/user";

export interface ICommentItem {
  /** 备注id */
  id?: number;
  /** 备注的员工id */
  member_id: number;
  /** 备注的员工工号 */
  member_user_id: string;
  /** 备注的员工昵称 */
  member_name: string;
  /** 员工头像 */
  avatar: string;
  /** 是否为VM人员，0为否，1为是 */
  is_vm: number;
  /** 是否自己备注，0否，1是 */
  is_self: number;
  /** 创建时间 */
  comment_time: number;
  /** 备注内容 */
  content: {
    text?: string;
    image_url?: IImageItem[];
  };
  /** 门店编码 */
  store_code: string;
  /** 0 默认展示，1为删除，2为编辑 */
  action: number;
  /** 回复备注的员工昵称 */
  reply_name?: string;
  // 以下是内部变量
  _reply_name?: string;
  _reply_content?: string;
  _reply_show?: boolean;
  _edit_content?: string;
  _edit_show?: boolean;
  _prefix_width?: string;
}
export interface IImageItem {
  /** 图片id */
  id?: number;
  /** 图片地址 */
  url: string;
  /** 图片的标识名称 */
  identifier?: string;
  /** 是否通过 0 通过 1 未通过 */
  status?: number; 
  /** VM的备注数 */
  vm_num?: number;
  /** 备注的总数 */
  num?: number;
  /** 是否查看过VM备注的气泡 */
  is_look?: number;
  /** 备注 */
  comment?: ICommentItem[];
  /** 上传内容 */
  _content?: string;
  /** 上传状态 uploading | done | error | 0 | 1 */
  _status?: string|number;
  /** 上传状态文案 */
  _message?: string;
  _image_list: any[];
  _comment_time?: number;
  _remark?: string;
}

export interface ITempItem {
  is_required: boolean;
  template_name: string;
  image_list?: IImageItem[];
}

export interface IStoreInfo {
  /** 任务ID */
  id: number;
  /** 分配表id */
  allot_id: number;
  /** 任务名称 */
  display_name: string;
  /** 任务完成状态,0为未提交，1为待审核，2为未通过，3为已通过 */
  finish_status: number;
  /** VM陈列的图片 */
  image_list: IImageItem[];
  /** 所属VM 导购ID */
  owner_member_id: number;
  /** 所属VM 导购昵称 */
  owner_member_name: string;
  /** 所属VM 导购工号 */
  owner_member_user_id: string;
  /** 评分 */
  grade: number;
  /** 整体评价 */
  review_remark: string;
  /** 门店编码 */
  store_code: string;
  /** 门店名称 */
  store_name: string;
  /** 门店ID */
  store_id?: string;
  /** 提交时间 */
  submit_time: number;
  /** 审核时间 */
  review_time: number;
  /** 图片分组 */
  template_info: ITempItem[];
}

@Component({
  components: {
    UploadImageButton
  }
})
export default class DisplayTaskAuditForm extends Vue {
  @Prop({ default: false, type: Boolean }) readonly: boolean;
  @Model('change', { default: ()=>(null), type: Object }) value: IStoreInfo;
  get userInfo() {
    return UserModule.getUserInfo
  }

  @Watch('value', { deep: true })
  onChange(value: IStoreInfo) {
    this.$emit('change', value);
  }
  

  defaultAvatar = 'https://cdn-cos.coach.com.cn/uploads/common/icon_logo.png'

  /**
   * 上传文件限制规范
   * @param file 
   */
  handleBeforeUpload(file: any) {
    if (file.type != "image/jpeg" && file.type != "image/png") {
      this.$message.error("请上传JPG、PNG格式的图片");
      return false;
    }
    // const name = file.name
    //   .substring(file.name.lastIndexOf(".") + 1, file.name.length)
    //   .toLocaleLowerCase();
    // if (name === "jpeg") {
    //   this.$message.warning("企业微信不支持JPEG的图片");
    //   return false;
    // }
    if (file.size > 10 * 1024 * 1024) {
      this.$message.error("请上传10MB大小以内的图片");
      return false;
    }
  }

  /**
   * 添加图片
   * @param imageItem 
   */
  handleAddImage(imageItem: IImageItem) {
    imageItem._comment_time = parseInt(Date.now() / 1000 + '')
  }

  /**
   * 删除图片
   */
  handleDeleteImage(imageItem: IImageItem, index: number) {
    imageItem._image_list.splice(index, 1)
  }
  
  /**
   * 上传图片
   * @param value 
   */
  handleImageChange(value: any[], imageItem: IImageItem) {
    imageItem._image_list.push(...value.map(item => ({
      url: item.url,
      _comment_time: imageItem._comment_time
    })))
    // imageItem.comment.push(this.createComment(value))
  }

  /**
   * @description 标记通过
   * @param item 
   */
  handleTag(item: IImageItem) {
    if(item.status === 1) return
    item._status = item._status === 1 ? 0 : 1;
  }

  /**
   * @description 回复
   * @param reply 
   */
  handleReply(item: ICommentItem, index: number) {
    item._reply_show = true
    item._reply_name = item.member_name;
    // 回复
    this.$nextTick(() => {
      const dom = this.$refs[`inputPrefix${index}`][0] as HTMLElement
      item._prefix_width = dom.clientWidth + 'px'
    })
    
  }

  /** 
   * @description 取消
   */
  handleReplyCancel(item: ICommentItem) {
    item._reply_show = false
    item._reply_content = '';
    item._reply_name = ''
  }

  /** 
   * @description 保存
   */
  handleReplySave(item: ICommentItem, imageItem: IImageItem) {
    if(item._reply_content.trim()=== '') {
      this.$message.warning('请输入备注')
      return
    }
    const comment = this.createComment(item._reply_content.trim(), item._reply_name)
    imageItem.comment.push(comment)
    item._reply_show = false
    item._reply_content = ''
    item._reply_name = ''
  }

    /**
   * @description 回复
   * @param reply 
   */
  handleEdit(item: ICommentItem) {
    item._edit_show = true
    item._edit_content = item.content.text
  }

  /** 
   * @description 取消
   */
  handleEditCancel(item: ICommentItem) {
    item._edit_content = "";
    item._edit_show = false
  }

  /** 
   * @description 保存
   */
  handleEditSave(item: ICommentItem) {
    if(item._edit_content.trim()=== '') {
      this.$message.warning('请输入备注')
      return
    }
    item.content.text = item._edit_content.trim()
    item._edit_show = false
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
      avatar: this.defaultAvatar,
      store_code: this.value.store_code,
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

  /**
   * @description 保存备注
   * @param item 
   */
  handleSaveComment(item: IImageItem) {
    if(item._remark.trim()=== '') {
      this.$message.warning('请输入备注')
      return
    }
    item.comment.push(this.createComment(item._remark))
    item._remark = ''
  }

  /**
   * @description 删除备注
   */
  handleDeleteComment(item: IImageItem, index: number) {
    item.comment.splice(index, 1)
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
      padding-bottom: 20px;
      border-bottom: 1px solid #eee;
      margin-bottom: 20px;
    }

    .image {
      position: relative;

      &-tag {
        position: absolute;
        top: 0;
        left: 0;
        width: 60px;
        height: 60px;
        background: url('~@/assets/images/icon_fail_big.png') center/100% no-repeat;
      }
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
  .btn-group {
    position: absolute;
    top: 0;
    left: 105px;
  }

  .upload-image-list {
    display: flex;
    flex-wrap: wrap;

    .upload-image-item {
      width: 60px;
      height: 60px;
      border-radius: 4px;
      margin-bottom: 10px;
      margin-right: 10px;
      position: relative;
      cursor: pointer;
    }
    .el-image {
      width: 60px;
      height: 60px;
      border-radius: 4px;
    }

    .upload-image-delete {
      position: absolute;
      top: 0;
      right: 0;
      transform: translate(50%, -50%);
      width: 15px;
      height:15px;
      border-radius: 50%;
      background-color: #000;
      color: #fff;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 12px;
    }
  }

}
.right {
  text-align: right;
  margin-top: 5px;
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
      margin-left: 5px;
    }

    .item-user-tag {
      padding: 0 10px;
      border-radius: 20px;
      margin-left: 5px;
      background: #C6AD7F;
      color: #fff;
      font-size: 10px;
      line-height: 20px;
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
      .reply-name {
        font-weight: 500;
      }
      .reply {
        height: 20px !important;
        line-height: 20px !important;
        padding: 0 5px!important;  
        margin-left: 10px;
      }

      &__footer {
        text-align: right;
        margin-top: 5px;
      }

      &__image {
        .el-image {
          width: 100px;
          height: 100px;
          margin-right: 10px;
          margin-bottom: 10px;
        }
      }
    }
  }
}

.input {
  position: relative;
  color: #999;

  &-prefix {
    position: absolute;
    top: 6PX;
    left: 15PX;
  }

  .reply-name {
    font-weight: 500;
  }

  ::v-deep .el-textarea__inner {
    text-indent: var(--text-indent, 0);
  }
}
</style>
