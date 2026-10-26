<template>
  <div class="comment-item">
    <van-image class="comment-item__avatar" round :src="item.avatar||defaultAvatar"></van-image>
    <div class="comment-item__body">
      <div v-if="is_vm == 0" class="comment-item__info">
        <span class="comment-item__name">{{item.member_name}}</span>
        <span v-if="item.store_code && item.is_vm == 0" class="comment-item__store">{{item.store_code}}</span>
        <span v-if="item.is_self == 1 && item.is_vm == 0" class="comment-item__tag gray">我</span>
        <span v-if="item.is_vm == 1" class="comment-item__tag">VM</span>
        <span class="comment-item__time">{{dateFormat(item.comment_time, 'yyyy/MM/dd hh:mm:ss')}}</span>
      </div>
      <div v-else class="comment-item__info">
        <span class="comment-item__name">{{item.member_name}}</span>
        <span v-if="item.store_code && item.is_vm == 0" class="comment-item__store">{{item.store_code}}</span>
        <span v-if="item.is_vm == 1" class="comment-item__tag">{{item.is_self == 1 ? '我':'VM'}}</span>
        <span class="comment-item__time">{{dateFormat(item.comment_time, 'yyyy/MM/dd hh:mm:ss')}}</span>
      </div>
      <div class="comment-item__content">
        <div v-if="item.content.text" class="comment-item__text">
          <pre class="comment-item__pre"><span v-if="item.reply_name" class="prefix">回复 <span class="prefix-name">{{ item.reply_name }}</span>：</span>{{item.content.text}}</pre>
        </div>
        <div v-if="item.content.image_url" class="comment-item__image-list">
          <image-item 
            v-for="(image, imageIndex) in item.content.image_url" 
            :key="imageIndex" 
            :item="image" 
            :disabled="!show(4)"
            :disabled-comment="true"
            @delete="handleDeleteImage(imageIndex)"
            @comment="handlePreviewImage(item.content.image_url, imageIndex)">
          </image-item>
        </div>
        <div v-if="show(1) || show(2) || show(3)" class="comment-item__footer">
          <div v-if="show(1)" class="btn" @click="handleDelete">
            <div class="btn-icon icon-delete"></div>
            <div class="delete-text">删除</div>
          </div>
          <div v-if="show(3)" class="btn" @click="handleEdit">
            <div class="btn-icon icon-edit"></div>
            <div class="delete-text">编辑</div>
          </div>
          <div v-if="show(2)" class="btn" @click="handleReply">
            <div class="btn-icon icon-reply"></div>
            <div class="delete-text">回复</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
import { Image as VanImage, ImagePreview } from "vant";
import ImageItem from "./image-item.vue";

@Component({
  components: {
    VanImage,
    ImageItem
  }
})
export default class CommentItem extends Vue {
  @Prop({ type: Object, required: true }) item: ICommentItem
  @Prop({ type: Boolean, default: false }) disabled: boolean
  @Prop({ type: Number, default: 0 }) is_vm: number
  @Prop({ type: Object, default: () => ({}) }) currentImage: any
  /** 用户信息 */
  get userInfo() {
    return this.$store.getters["user/getUserInfo"]
  }
  get show() {
    const { id, is_self, is_vm, content } = this.item
    const is_text = !!content.text
    return (key: number) => {
      switch (key) {
        case 1:
          // 是否显示删除
          if(this.is_vm === 1) {
            // VM身份 
            // 自己添加的备注 && 未保存 && 是文本 && 不是禁用 就可以编辑
            return is_vm === 1 && is_self === 1 && is_text && !this.disabled
          } else if(this.is_vm == 0) {
            // 导购身份 
            // 自己添加的备注 && 是文本 && 不是禁用 就可以编辑
            return is_vm === 0 && is_self === 1 && is_text &&!this.disabled
          } else {
            return false
          }
        case 2:
          // 是否显示回复按钮
          if(this.is_vm == 1) {
            return is_vm === 0  && !this.disabled
          } else if(this.is_vm == 0) {
            return is_vm === 1  && !this.disabled
          } else  {
            return false
          }
        case 3:
          // 是否显示编辑
          if(this.is_vm === 1) {
            // VM-special 草稿 自己添加的备注 && 是文本 && 不是禁用 就可以编辑
            if(this.userInfo.is_special && this.currentImage.is_draft) {
              return is_vm === 1 && is_self === 1 && is_text && !this.disabled
            }
            
            // VM身份 自己添加的备注 && 未保存 && 是文本 && 不是禁用 就可以编辑
            return is_vm === 1 && is_self === 1 && is_text && !this.disabled && !id
          } else {
            // 导购身份 自己添加的备注 && 是文本 && 不是禁用 就可以编辑
            return is_vm === 0 && is_self === 1 && is_text && !this.disabled
          }
        case 4:
          // 是否显示图片备注的删除
          // VM身份添加的备注 && 未保存 && 是文本 && 不是禁用 就可以编辑
          if(this.is_vm == 1) {
            return is_vm === 1 && is_self === 1 && !id && !this.disabled
          } else if(this.is_vm == 0) {
            return false
          } else {
            return false
          }
          
        default:
          break;
      }
    }
  }
  defaultAvatar = 'https://cdn-cos.coach.com.cn/uploads/common/icon_logo.png'

  /**
   * @description 格式化时间
   * @param time 
   * @param format 
   */
  dateFormat(time: number, format = "yyyy/MM/dd") {
    return this.$helper.dateFormat(time, format);
  }

  /**
   * @description 删除评论
   */
  handleDelete() {
    this.$emit("delete", this.item);
  }

  /** 
   * @description 删除图片
   * @param index 
   */
  handleDeleteImage(index: number) {
    this.$emit("delete-image", index);
  }

  /**
   * 预览图片
   * @param list 
   * @param current 
   */
  handlePreviewImage(list: IImageItem[], current: number) {
    ImagePreview({
      images: list.map(item => item.url),
      startPosition: current,
    });
  }

  /**
   * @description 编辑评论
   */
  handleEdit() {
    this.$emit("edit", this.item);
  }

  /**
   * @description 回复评论
   */
  handleReply() {
    this.$emit("reply", this.item);
  }
}
</script>
<style lang="scss" scoped>
.prefix {
  pointer-events: none; /* 防止伪元素阻碍用户与textarea交互 */
  font-size: 12px;

  &-name {
    color: #777777;
    font-weight: 500;
  }
}
.comment-item {
  font-size: 14px;
  color: #000000;
  display: flex;
  align-items: flex-start;
  margin-bottom: 15px;

  &__avatar {
    width: 35px;
    height: 35px;
    margin-right: 20px;
    flex-shrink: 0;
  }

  &__body {
    flex: 1;
  }

  &__info {
    display: flex;
    align-items: center;
    white-space: nowrap;
  }

  &__name {
    font-weight: 500;
    line-height: 20px;
    margin-right: 4px;
    max-width: 30vw;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  &__store {
    color: #777777;
    line-height: 17px;
    margin-right: 4px;
  }

  &__tag {
    font-size: 10px;
    padding: 0 10px;
    border-radius: 17px;
    background: #C6AD7F;
    color: #fff;
    margin-right: 4px;

    &.gray {
      background: #A6A6A6;
    }
  }

  &__time {
    font-size: 10px;
    color: #A6A6A6;
    flex: 1;
    text-align: right;
  }

  &__content {
    margin-top: 6px;
    background: #F6F5F5;
    box-shadow: -2px 3px 8px 0px rgba(0,0,0,0.1);
    padding: 15px;
    position: relative;
    &::before {
      content: "";
      width: 26px;
      height: 31px;
      position: absolute;
      left: -20px;
      top: 5px;
      background: url('~@/assets/display-task/icon_tr.png') no-repeat center/100%;
      // background: #F6F5F5;
    }
    &::after {
      content: "";
      width: 4px;
      height: 80%;
      position: absolute;
      left: 0;
      top: 50%;
      transform: translateY(-50%);
      background: #F6F5F5;
    }
  }

  &__pre {
    font-size: 12px;
    color: #777777;
    white-space: pre-wrap;
  }

  &__image-list {
    display: flex;
    flex-wrap: wrap;

    .image-item:nth-child(3n) {
      margin-right: 15px;
    }

    .image-item:nth-last-child(1) {
      margin-bottom: 0;
    }
    .image-item:nth-last-child(2) {
      margin-bottom: 0;
    }
  }

  &__footer {
    margin-top: 10px;
    display: flex;
    align-items: center;
    justify-content: flex-end;

    .btn {
      display: flex;
      align-items: center;
      justify-content: flex-end;
      margin-left: 15px;
      font-size: 12px;
      color: #000;
      &-icon {
        width: 20px;
        height: 20px;
        margin-right: 5px;

        &.icon-delete {
          background: url('~@/assets/display-task/icon_remove.png') no-repeat center/100%;
        }

        &.icon-edit {
          background: url('~@/assets/display-task/icon_edit.png') no-repeat center/100%;
        }

        &.icon-reply {
          margin-top: -2px;
          background: url('~@/assets/display-task/icon_reply.png') no-repeat center/100%;
        }
      }

    }
  }
}
</style>