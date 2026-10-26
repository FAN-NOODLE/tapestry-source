<template>
  <div class="scroll">
    <div ref="scroll" class="scroll-content" :style="{ maxHeight: maxHeight + 'px' }">
      <div ref="content" class="scroll-item">
        <slot></slot>
      </div>
    </div>
    <div ref="scrollbar" class="scroll-bar">
      <div ref="scrollthumb" class="scroll-thumb"></div>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Model, Prop, Vue, Watch } from 'vue-property-decorator';


@Component({})
export default class Scrollbar extends Vue {
  @Prop({ type: String, default: '' }) maxHeight: string

  init() {
    const scroll = this.$refs.scroll as HTMLElement;
    const content = this.$refs.content as HTMLElement;
    const scrollbar = this.$refs.scrollbar as HTMLElement;
    const scrollthumb = this.$refs.scrollthumb as HTMLElement;
    // 计算thumb的高度
    
    scroll.addEventListener('scroll', () => {
      const min = 0
      const max = scrollbar.clientHeight - scrollthumb.clientHeight;
      scrollthumb.style.top = `${Math.max(min, Math.min(scroll.clientHeight / content.clientHeight * scroll.scrollTop, max))}px`;
      // scrollbar.style.left = `${scroll.scrollLeft}px`;
    });

    // 监听content尺寸变化
    const observer = new ResizeObserver(() => {
      
      if(scroll.clientHeight < content.clientHeight) {
        scrollbar.style.display = 'block';
      } else {
        scrollbar.style.display = 'none';
      }
      const h = scroll.clientHeight / content.clientHeight * scrollbar.clientHeight
      scrollthumb.style.height = `${h||0}px`;
      console.log('滚动条', scroll.clientHeight, content.clientHeight, scrollbar.clientHeight, h);
      
    });
    observer.observe(content);
  }

  mounted() {
    this.$nextTick(() => {
      this.init()
    })
  }
}

</script>

<style lang="scss" scoped>
.scroll {
  width: 100%;
  height: 100%;
  position: relative;
  overflow: hidden;

  .scroll-content {
    width: 100%;
    height: 100%;
    overflow-y: auto;
    &::-webkit-scrollbar {
      display: none;
    }
  }


  .scroll-bar {
    position: absolute;
    width: 4px;
    // height: 100%;
    background: rgba(0, 0, 0, 0.1);
    right: 0;
    top: 0;
    bottom: 10px;
  }

  .scroll-thumb { 
    position: absolute;
    left: 0;
    width: 100%;
    height: 0;
    background: rgba(0, 0, 0, 0.2);
    border-radius: 4px;
  }
}
</style>