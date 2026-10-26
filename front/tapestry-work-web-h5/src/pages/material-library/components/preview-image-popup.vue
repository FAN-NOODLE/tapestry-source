<template>
  <!--分组弹窗 -->
  <van-popup
    v-model="visible"
    @close="handleClose"
    :close-on-click-outside="false"
    :close-on-click-action="false"
    class="material-popup"
  >
    <div
      class="popup-container"
      @click="isShowTips = false"
    >
      <div class="popup-header">
        <div class="top-num">{{ activeIndex + 1 }}/{{ list.length }}</div>
        <div
          class="popup-header-close"
          @click="handleClose"
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
          ref="swipe"
          :show-indicators="false"
          style="height: 420px;"
          @change="onChange"
          :loop="false"
          :activeIndex="activeIndex"
          :initial-swipe="showIndex"
        >
          <van-swipe-item
            v-for="(item, index) in list"
            :key="index"
          >
            <video
              v-if="item.type === 5"
              ref="video"
              :src="item.url"
              :controls="true"
              class="video"
              playsinline="true"
              x-webkit-airplay="allow"
              x5-playsinline
              webkit-playsinline
              object-fit="contain"
              :poster="item.poster"
            ></video>
            <van-image
              v-else
              lazy-load
              class="swiper-pic"
              :src="item.url"
              fit="contain"
            />
            <div
              class="sel-tag"
              @click="handleSelect(item)"
            >
              <van-icon :name="isChecked(item)
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
import { getVideoPoster } from "@/utils/tools";
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
export default class PreviewImagePopup extends Vue {
  @PropSync("show", { type: Boolean }) visible;
  @Prop(Object) item!: any;
  // 选中的图片
  @Prop(Array) selection!: any;
  // 是否第一次打开
  @Prop(Boolean) isFirstOpen!: any;
  // 显示第几个图片,默认第一个
  @Prop({ type: Number, default: 0 }) showIndex!: any;


  /** 是否是组合图片 */
  get isMergeImage() {
    if(!this.item) return false
    return this.item.merge_image && this.item.merge_image.length > 1
  }

  /** 素材列表 */
  get list () {
    if(!this.item) return []
    if(this.isMergeImage || this.item.type==8) {
      return this.item.merge_image
    } else {
      return [this.item]
    }
  }

  get activeName() {
    if(!this.list[this.activeIndex]) return ''
    return this.list[this.activeIndex].name
  }

  get isChecked() {
    return (item)=> {
      const ids = this.selection.map(item => item.id)
      return ids.includes(item.id)
    }
  }

  @Watch('show')
  onShowIndexChange(val: any) {
    if(val) {
      this.activeIndex = this.showIndex;
      this.$nextTick(() => {
        const videoRef = this.$refs.video as HTMLVideoElement[]
        videoRef?.forEach(async element => {
          const res:any = await getVideoPoster(element)
          element.poster = res.base64
        });
        (this.$refs.swipe as any).swipeTo(this.showIndex);
      })
    }
  }
  @Watch('list', { immediate: true, deep: true })
  onImgListChange(val: any) {
    if (val.length > 1 && this.isFirstOpen) {
      this.isShowTips = true
    } else {
      this.isShowTips = false
    }
  }
  activeIndex = 0
  isShowTips = false

  mounted() {

    this.$nextTick(() => {
      this.trackAssetView()
    })
  }

  /** 预览上报 */
  trackAssetView() {
    const { activeIndex, list, item } = this
    const value = list[activeIndex]
    console.log('神策', list)
    console.log('神策', item)
    const params = {
      source_id: item.id,
      task_id: '素材库无任务ID',
      is_official: 1,
      asset_type_modify: item.type,
      official_url: '非官网商品素材',
      official_upc: '非官网商品素材',
      official_product: '非官网商品素材',
      official_name: '非官网商品素材',
      current_page: this.$route.query.path||this.$route.path,
      entry_point: 2,
      library_type: 2,
      is_top: !!item.is_top,
      is_favorite: !!item.is_collect
    }

    if (item.type === 2) {
      // 图片素材
      params.source_id = value.id
      params.is_favorite = !!value.is_collect
    }

    if (item.type === 8) {
      // 商品素材
      params.source_id = value.id
      params.is_official = 2
      params.asset_type_modify = value.type
      params.official_url = value.url
      params.official_upc = value.upc_code
      params.official_product = value.product_code
      params.official_name = value.product_name
    }

    // 神策数据上报 -- 查看素材
    this.$track.trackAssetView(params)
  }


  handlePlay() {
    this.trackAssetView()
  }

  handleClose() {
    this.visible = false
    this.$emit("close");
  }
  /**
   * @description: change事件
   * @return {*}
   */
  onChange(index: number) {
    this.activeIndex = index
    this.isShowTips = false
    if(this.list[index].type!==5) {
      this.trackAssetView()
    }
  }
  /**
   * @description: 选中
   * @return {*}
   */
  handleSelect(item: any) {
    const ids = this.selection.map(item => item.id)
    // 判断是否全部选中
    this.$emit("select", {
      status: ids.includes(item.id) ? 2: 0,
      item: { type: 2, ...item },
      isMergeImage: false
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
      .video {
        width: 100%;
        height: 100%;
        // background-color: #000;
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
