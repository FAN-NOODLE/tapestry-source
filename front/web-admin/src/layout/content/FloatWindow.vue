<template>
  <div
    ref="floatWindowRef"
    class="float-window"
    :style="style"
    @mousedown="handleMouseDown"
  >
    <slot></slot>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component
export default class FloatWindow extends Vue {
  @Prop({ type: [Number, String], default: 60 }) top: string|number
  @Prop({ type: [Number, String], default: 'auto' }) left: string|number
  @Prop({ type: [Number, String], default: 60 }) right: string|number
  @Prop({ type: [Number, String], default: 'auto' }) bottom: string|number

  current = {
    top: 0,
    left: 0,
  }

  get style() {
    return {
      top: this.current.top + 'px',
      left: this.current.left + 'px',
    }
  }

  mounted() {
    this.$nextTick(() => {
      this.computedPosition()
    })
  }

  computedPosition() {
    const dom = this.$refs.floatWindowRef as HTMLElement
    let left = 0;
    const offsetRight = this.right === 'auto' ? 0 : parseInt(this.right as string);
    if (this.left === 'auto') {
      left = document.body.clientWidth - dom.clientWidth - offsetRight;
    } else {
      left = parseInt(this.left as string);
    }
    let top = 0;
    if (this.top === 'auto') {
      top =
        document.body.clientHeight -
        dom.clientHeight -
        parseInt(this.bottom as string);
    } else {
      top = parseInt(this.top as string);
    }

    this.current =  {
      top: top,
      left: left,
    };
  }

  /**
   * 鼠标按下
   * @param e
   */
  handleMouseDown(e: MouseEvent) {
    const that = this
    const dom = this.$refs.floatWindowRef as HTMLElement;
    const start = {
      x: e.clientX,
      y: e.clientY,
    };

    /**
     * 鼠标移动
     * @param e
     */
    function handleMouseMove(e: MouseEvent) {
      if (!dom) return;
      const dx = e.clientX - start.x;
      const dy = e.clientY - start.y;
      start.x = e.clientX;
      start.y = e.clientY;

      that.current.left += dx;
      that.current.top += dy;

      // 规定窗口不能超出窗口元素
      const parent = document.body;
      const maxX = parent.clientWidth - dom.clientWidth;
      const maxY = parent.clientHeight - dom.clientHeight;
      if (that.current.left < 0) {
        that.current.left = 0;
      } else if (that.current.left > maxX) {
        that.current.left = maxX;
      }
      if (that.current.top < 0) {
        that.current.top = 0;
      } else if (that.current.top > maxY) {
        that.current.top = maxY;
      }
    }
    /**
     * 鼠标抬起
     * @param e
     */
    function handleMouseUp() {
      // 如果没有移动则不处理
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mouseup', handleMouseUp);
    }

    document.addEventListener('mousemove', handleMouseMove);
    document.addEventListener('mouseup', handleMouseUp);
  }
}
</script>

<style lang="scss" scoped>
.float-window {
  position: fixed;
  z-index: 1000;
}
</style>