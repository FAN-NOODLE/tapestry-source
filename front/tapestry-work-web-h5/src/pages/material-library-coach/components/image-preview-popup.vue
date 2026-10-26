<template>
  <!--分组弹窗 -->
  <van-popup
    v-model="show"
    @close="closePopup"
    :close-on-click-outside="false"
    :close-on-click-action="false"
    class="material-popup"
  >
    <div
      class="popup-container"
      @click="isShowTips = false"
    >
      <div class="popup-header">
        <div class="top-num">{{ activeIndex }}/{{ imgList.length }}</div>
        <div
          class="popup-header-close"
          @click="closePopup"
        >
          <van-image
            class="close-pic"
            src="https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/close-white.png"
            fit="cover"
          />

        </div>
      </div>
      <div class="popup-content">
        <div
          class="popup-tips"
          v-if="isShowTips"
        >
          <van-image
            class="tips-pic"
            src="https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/next-icon.png"
            fit="cover"
          />
        </div>
        <van-swipe
          :show-indicators="false"
          style="height: 420px;"
          @change="onChange"
          :loop="false"
          :initial-swipe="showIndex"
        >
          <van-swipe-item
            v-for="(item, index) in imgList"
            :key="index"
          >
            <van-image
              lazy-load
              class="swiper-pic"
              :src="item.url"
              fit="contain"
            />
            <div
              class="sel-tag"
              @click="checkedItem(item)"
            >
              <van-icon :name="chooseArr.findIndex((i) => i.id === item.id) > -1
                ? 'checked'
                : 'circle'
                " />
            </div>
          </van-swipe-item>
        </van-swipe>
      </div>
      <div class="popup-footer">
        <div class="image-name">
          <vue-ellipsis
            :visible-line="2"
            :text="activeName.split('.')[0]"
          >
          </vue-ellipsis>
        </div>
      </div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Vue, PropSync, Prop, Watch } from "vue-property-decorator";
import { Popup, Image, Swipe, SwipeItem, Icon } from "vant";
// todo: 文本省略组件
import { VueEllipsis } from 'vue-ellipsis-component';
import { watch } from "vue";
@Component({
  components: {
    [Popup.name]: Popup,
    [Image.name]: Image,
    [Swipe.name]: Swipe,
    [SwipeItem.name]: SwipeItem,
    [Image.name]: Image,
    [Icon.name]: Icon,
    VueEllipsis
  },
})
export default class MaterialGroup extends Vue {
  [x: string]: any;
  @PropSync("showPopup", { type: Boolean }) show;
  @Prop(Array) imgList!: any;
  // 选中的图片
  @Prop(Array) chooseArr!: any;
  // 是否第一次打开
  @Prop(Boolean) isFirstOpen!: any;
  // 显示第几个图片,默认第一个
  @Prop(Number) showIndex!: any;
  @Watch('showIndex', { immediate: true, deep: true })
  onShowIndexChange(val: any) {
    this.activeIndex = val + 1
    this.activeName = this.imgList[val].name
  }
  @Watch('imgList', { immediate: true, deep: true })
  onImgListChange(val: any) {
    if (val.length > 1 && this.isFirstOpen) {
      this.isShowTips = true
    } else {
      this.isShowTips = false
    }
  }
  activeIndex = 1
  activeName = ''
  isShowTips = false
  closePopup() {
    this.$emit("closePopup");
  }
  /**
   * @description: change事件
   * @return {*}
   */
  onChange(index: number) {
    this.activeIndex = index + 1
    this.activeName = this.imgList[index].name
    this.isShowTips = false
  }
  /**
   * @description: 选中
   * @return {*}
   */
  checkedItem(item: any) {
    this.$emit("checkedItem", {
      ...item,
      type: 2,
      is_deal: 2,
    });
  }
}
</script>
<style lang="scss" scoped>
::-webkit-scrollbar {
  display: none;
  /* Chrome Safari */
}

.material-popup {
  background: none;

  .popup-container {
    width: 331px;

    .popup-header {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .popup-header-close {
        position: relative;
        bottom: 8px;
        right: 0;

        .close-pic {
          width: 12px;
          height: 12px;
        }
      }

      .top-num {
        width: 100%;
        font-size: 16px;
        font-weight: 500;
        color: #FFFFFF;
        text-align: center;
      }
    }

    .popup-content {
      position: relative;
      .swiper-pic{
        width: 100%;
        height: 100%;
      }
      .popup-tips {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        width: 200px;
        height: 200px;
        z-index: 9;
      }

      .van-swipe-item {
        display: flex;
        align-items: center;
        justify-content: center;
        background: #F0F0F0;
      }


      .sel-tag {
        position: absolute;
        width: 28px;
        height: 28px;
        display: flex;
        justify-content: center;
        align-items: center;
        box-shadow: -2px 2px 8px 0px rgba(0, 0, 0, 0.1);
        border-radius: 0px 0px 0px 10px;
        right: 0;
        top: -1px;
        line-height: 22px;
        text-align: center;
        background-color: #ffffff;
        z-index: 1;

        .van-icon {
          font-size: 16px;
          vertical-align: middle;
          color: #4a4a4a;
        }
      }
    }

    .popup-footer {
      display: flex;
      justify-content: center;
      margin-top: 20px;

      .button {
        width: 165px;
        height: 40px;
        font-size: 16px;
        color: #ffffff;
        line-height: 34px;
        text-align: center;
        background-color: #000;
        line-height: 40px;
      }

      .image-name {
        height: 46px;
        overflow: hidden;
        font-size: 16px;
        font-weight: 500;
        color: #FFFFFF;
        margin-bottom: 10px;
        text-align: center;
      }
    }
  }
}</style>
