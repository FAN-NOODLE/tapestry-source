<template>
  <div ref="watermark" :class="['watermark', {'watermark--full': fullPage}]" :style="style">
    <div ref="wrap" class="watermark__wrapper">
      <svg
          :viewBox="`0 0 ${svgWidth} ${svgHeight}`"
          :width="svgWidth"
          :height="svgHeight"
          xmlns="http://www.w3.org/2000/svg"
          xmlns:xlink="http://www.w3.org/1999/xlink"
          :style="{
            padding: `0 ${gapX}px ${gapY}px 0`,
            opacity: opacity,
          }"
        >
          <foreignObject x="0" y="0" :width="width" :height="height">
            <div xmlns="http://www.w3.org/1999/xhtml" :style="{
              transformOrigin: 'center',
              transform: `rotate(${rotate}deg)`,
              width: `${svgWidth}px`,
              height: `${svgHeight}px`,
              display: 'flex',
              justifyContent: 'center',
              alignItems: 'center',
              color: textColor,
              fontSize: fontSize,
            }">
              <slot>
                <span :style="{ 
                  color: textColor,
                  fontSize: fontSize
                 }">{{content||watermarkContent}}</span>
              </slot>
            </div>
          </foreignObject>
        </svg>
    </div>
  </div>
</template>
<script lang="ts">
import { Prop, Component, Vue, Watch } from "vue-property-decorator";

const makeSvgToBlobUrl = (svgStr: string) => {
  // svg MIME type: https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types
  const svgBlob = new Blob([svgStr], {
    type: 'image/svg+xml',
  });
  return URL.createObjectURL(svgBlob);
};

@Component
export default class Watermark extends Vue {
  @Prop({ type: String, default: "" }) content!: string;
  @Prop({ type: Number, default: 180 }) width!: number;
  @Prop({ type: Number, default: 180 }) height!: number;
  @Prop({ type: Number, default: 0 }) gapX!: number;
  @Prop({ type: Number, default: 0 }) gapY!: number;
  @Prop({ type: String, default: "#A6A6A6" }) textColor!: string;
  @Prop({ type: String, default: "12px" }) fontSize!: string;
  @Prop({ type: [Number, String], default: 1 }) opacity!: number|string;
  @Prop({ type: [Number, String], default: 3000 }) zIndex!: number|string;
  @Prop({ type: [Number, String], default: -45 }) rotate!: number|string;
  @Prop({ type: Boolean, default: true }) fullPage!: boolean;

  watermarkUrl = ''
  /** 用户信息 */
  get userInfo(): IUserInfo {
    return this.$store.getters["user/getUserInfo"]
  }

  get watermarkContent() {
    if(!this.userInfo) return ""
    const { name, user_id } = this.userInfo
    if(name && user_id) {
      return `${name||'导购姓名'}-${user_id||'导购工号'}`
    } else {
      return ""
    }
  }

  get style() {
    const style: any = {};
    if (this.zIndex !== undefined) {
      style.zIndex = +this.zIndex;
    }
    // 背景图片
    if(this.watermarkUrl) {
      style.backgroundImage = `url(${this.watermarkUrl})`;
    }
    style.backgroundSize = `${this.svgWidth}px ${this.svgHeight}px`;
    style.backgroundPosition = `-${this.svgWidth/2}px -${this.svgHeight/2}px`;
    return style
  }

  get svgWidth() {
    return this.width + this.gapX;
  }

  get svgHeight() {
    return this.height + this.gapY;
  }

  @Watch("userInfo", { immediate: true, deep: true })
  handleUserInfoChange() {
    this.$nextTick(() => {
      this.init();
    })
  }
  
  @Watch("contnet")
  handleContentChange() {
    this.$nextTick(() => {
      this.init();
    })
  }

  mounted() {
    this.$nextTick(() => {
      this.init();
    })
  }

  init() {
    const html = (this.$refs.wrap as HTMLElement).innerHTML
    this.watermarkUrl = makeSvgToBlobUrl(html)
    // 将水印节点插入到body中
    document.body.appendChild(this.$refs.watermark as HTMLElement)
  }
}
</script>
<style lang="scss" scoped>
.watermark {
  position: absolute;
  height: 100%;
  width: 100%;
  left: 0;
  top: 0;
  background-repeat: repeat;
  pointer-events: none;

  &__wrapper {
    display: none;
  }

  &--full {
    position: fixed;
  }
}
</style>