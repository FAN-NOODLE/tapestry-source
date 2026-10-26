<template>
  <div class="image-item" @click="handleClick">
    <div class="image-item__image">
      <van-image v-if="item._content" class="image-top" :src="item._content" fit="cover" :preview-src-list="[item.url]" lazy-load></van-image>
      <van-image 
        v-if="item.url" 
        class="image-current" 
        :class="{
          hide: !!item._content
        }"
        :src="item.url+'?imageMogr2/thumbnail/500x'" 
        fit="cover" 
        :preview-src-list="[item.url]" lazy-load @load="handleLoad">
        <template v-slot:loading>
          <van-loading type="spinner" size="20" />
        </template>
      </van-image>
    </div>
    <!-- 底部信息 -->
    <div v-if="show(1)" class="image-item__info">
      <div class="image-item__info-title">{{item.num?`${item.num}条备注`:'添加备注'}}</div>
      <div v-if="item.is_look === 1 && !disabledDelete" class="image-item__info-desc">{{item.vm_num> 999 ? '999+' : item.vm_num}}</div>
    </div>
    <!-- 评审失败角标 -->
    <div v-if="show(2)" class="image-item__fail"></div>
    <!-- 加载中 -->
    <div v-if="show(3)" class="image-item__mask" @click.stop="handleClickMask">
      <van-loading size="24" color="#fff"></van-loading>
      <span>{{ item._message }}</span>
    </div>
    <!-- 删除角标 -->
    <div v-if="show(4)" class="image-item__delete" @click.stop="handleDelete"></div>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
import { Image as VanImage, Loading } from 'vant';
import ConfirmPopup from "./confirm-popup.vue";

@Component({
  components: {
    VanImage,
    Loading,
    ConfirmPopup
  }
})
export default class ImageItem extends Vue {
  @Prop({ type: Object, default: ()=>({})}) item: IImageItem;
  @Prop({ type: Boolean, default: false }) disabled: boolean;
  @Prop({ type: Boolean, default: false }) disabledDelete: boolean;
  @Prop({ type: Boolean, default: false }) disabledComment: boolean;

  get show() {
    const { disabled, disabledDelete, disabledComment, item: { _upload_status, _status, num, status } } = this;
    return (key: number) => {
      switch (key) {
        case 1:
          // 是否显示下面的备注
          return (!disabled || num) && _upload_status !== 'uploading' && !disabledComment
        case 2:
          // 是否显示审核不通过的角标
          return _status === 1
        case 3:
          // 是否显示审核不通过的角标
          return _upload_status === 'uploading'
        case 4:
          // 是否显示删除按钮
          return _upload_status !== 'uploading' && status !== 1 && !disabled && !disabledDelete
      
        default:
          break;
      }
    }
  }

 
  /**
   * @description 删除图片
   */
  handleDelete() {
    this.$emit('delete');
  }

  /**
   * @description 点击图片
   */
  handleClick() {
    this.$emit('comment')
  }

  /**
   * @description 点击遮罩层
   */
  handleClickMask() {
    this.$emit('click-mask');
  }

  handleLoad() {
    console.log(111111);
    
    this.item._content = ''
    this.$emit('auto-save');
  }
}
</script>
<style lang="scss" scoped>
.image-item {
  width: 100px;
  height: 100px;
  margin-right: 15px;
  position: relative;
  margin-bottom: 15px;

  &__image {
    width: 100%;
    height: 100%;
    position: relative;
    .van-image {
      width: 100%;
      height: 100%;
    }

    .image-top {
      position: absolute;
      top: 0;
      left: 0;
      z-index: 1;
    }
    .image-current {
      position: absolute;
      top: 0;
      left: 0;
      opacity: 1;

      &.hide {
        opacity: 0;
      }
    }
  }

  &:nth-child(3n) {
    margin-right: 0;
  }


  &__mask {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(50, 50, 51, 0.88);
    display: flex;
    align-items: center;
    justify-content: center;
    flex-direction: column;
    color: #fff;
    font-size: 12px;
    z-index: 1;
  }

  &__delete {
    width: 18px;
    height: 18px;
    background: url("~@/assets/display-task/icon_delete.png") center/100% no-repeat;
    position: absolute;
    top: -5px;
    right: -5px;
    z-index: 1;
  }

  &__info {
    width: 100%;
    height: 20px;
    position: absolute;
    bottom: 0;
    background-color: rgba($color: #000000, $alpha: 0.7);
    z-index: 1;
    &-title {
      font-size: 12px;
      color: #fff;
      text-align: center;
      line-height: 20px;
    }

    &-desc {
      position: absolute;
      top: -8px;
      right: 15px;
      transform: translateX(50%);
      background: #E02020;
      border-radius: 8px;
      height: 16px;
      padding: 0 8px;
      line-height: 16px;
      font-size: 10px;
      color: #fff;
    }
  }

  &__fail {
    position: absolute;
    left: -1px;
    top: -1px;
    width: 40px;
    height: 40px;
    background: url("~@/assets/display-task/icon_fail.png") center/100% no-repeat;
    z-index: 1;
  }
}

</style>