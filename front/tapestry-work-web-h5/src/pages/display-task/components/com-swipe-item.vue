<template>
  <div
    :ref="hash"
    class="preview-wrap"
    @click="handleClick"
    @touchmove="handleTouchMove"
    @touchstart="handleTouchStart"
    @touchend="handleTouchEnd"
    @dblclick="handleDblClick"
  >
    <div ref="draggable" class="draggable" :style="style">
      <div class="preview-wrap__image" @click.stop="()=>{}">
        <div v-if="item._status === 1" class="preview-wrap__mask"></div>
        <img ref="image" :src="item.url" @load="handleLoad"/>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
import { Image as VanImage } from "vant";
import crypto from "crypto";
function sha256(str) {
    const hash = crypto.createHash('sha256');
    hash.update(str);
    return hash.digest('hex').substring(0, 16);;
}

// Example usage
const inputString = 'Hello, World!';
const hash = sha256(inputString);
console.log('SHA-256 Hash:', hash);
@Component({
  components: {
    VanImage
  }
})
export default class ComSwipe extends Vue {
  @Prop({ type: Object, required: true }) item: IImageItem
  /** 双指距离 */
  distance: number = null
  /** 缩放比例 */
  scale: number = 1
  /** 开始点 */
  start = { x: 0, y: 0 }
  /** 偏移量 */
  offset = { x: 0, y: 0 }
  /** 是否正在缩放 */
  scaling = false
  /** 图片样式 */
  get style() {
    return {
      transform: `scale(${this.scale}) translate(${this.offset.x}px, ${this.offset.y}px)`,
    };
  }

  get hash() {
    return sha256(JSON.stringify(this.item))
  }
  handleClick() {
    this.$router.back()
  }
  /**
   * 图片加载完成后处理
   */
  handleLoad(e:Event) {
    const wrap = this.$refs[this.hash] as HTMLElement
    const image = e.currentTarget as HTMLElement
    const wrapHeight = wrap.clientHeight
    const wrapWidth = wrap.clientWidth
    const imageHeight = image.clientHeight
    const imageWidth = image.clientWidth
    // 判断图片是否是横屏
    const rate = imageWidth / imageHeight
    const wrapRate = wrapWidth / wrapHeight
    console.log(e.currentTarget);
    
    if(rate > wrapRate) {
      // 图片宽度大于屏幕宽度
      const w = Math.min(wrapWidth, imageWidth)
      image.style.width = `${w}px`
      image.style.height = `${w * imageHeight / imageWidth}px`
    } else {
      // 图片高度大于屏幕高度
      const h = Math.min(wrapHeight, imageHeight)
      image.style.height = `${h}px`
      image.style.width = `${h * imageWidth / imageHeight}px`
    }
  
  }

  handleTouchStart(event: TouchEvent) {
    if(event.touches.length === 1) {
      this.scaling = false;
      const touch = event.touches[0];
      this.start = { x: touch.clientX - this.offset.x, y: touch.clientY - this.offset.y };
    } else if(event.touches.length === 2) {
      this.distance = Math.hypot(
        event.touches[0].pageX - event.touches[1].pageX,
        event.touches[0].pageY - event.touches[1].pageY
      );
      this.scaling = true;
    }
  }
  /**
   * 手势
   * @param event 
   */
  handleTouchMove(event: TouchEvent) {
    
    if(event.touches.length === 1) {
      // 如果正在缩放
      if(this.scaling) {
        return
      }
      // 如果大于1阻止默认事件
      if(this.scale > 1) {
        event.stopPropagation();
      } else {
        return
      }
      
      event.preventDefault();
      var touch = event.touches[0];
      let x = touch.clientX - this.start.x;
      let y = touch.clientY - this.start.y;
      // 拖动边界
      const dom = this.$refs.draggable as HTMLElement;
      const maxOffsetX = (dom.clientWidth * this.scale - dom.clientWidth) / 2
      const maxOffsetY = (dom.clientHeight * this.scale - dom.clientHeight) / 2

      if(x < 0- maxOffsetX) {
        x = 0- maxOffsetX;
      } else if(x > maxOffsetX) {
        x = maxOffsetX;
      }

      if(y < 0- maxOffsetY) {
        y = 0- maxOffsetY;
      } else if(y > maxOffsetY) {
        y = maxOffsetY;
      }
      
      this.offset.x = x;
      this.offset.y = y;
    
    } else if (event.touches.length === 2) {
      event.stopPropagation();
      event.preventDefault();
      const distance = Math.hypot(
        event.touches[0].pageX - event.touches[1].pageX,
        event.touches[0].pageY - event.touches[1].pageY
      );

      if (this.distance) {
        this.scale = this.scale * distance / this.distance;
        this.distance = distance;
        // 如果缩放小于1 
        if(this.scale <= 1) {
          this.offset = { x: 0, y: 0 }
        }
      }
    }
  }

  handleTouchEnd(event: TouchEvent) {
    if(event.touches.length === 2) {
      this.scaling = false;
    }
  }

  handleDblClick() {
    this.scale = 1;
  }
}
</script>
<style lang="scss" scoped>
.preview-wrap {
  width: 100%;
  height: 100%;
  position: relative;
  overflow: hidden;

  .draggable {
    width: 100%;
    height: 100%;
    transform-origin: center center;
    position: relative;
  }

  &__image {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    img {
      width: auto;
      display: block;
      padding: 0;
      margin: 0;
    }
  }

  &__mask {
    position: absolute;
    left: 0;
    top: 0;
    width: 55px;
    height: 55px;
    background: url("~@/assets/display-task/icon_fail_big.png") center/100% no-repeat;
  }
}
</style>
