<template> 
  <div class="com-popover">
    <div class="com-popover-reference" ref="content" @click="toggle">
      <slot name="reference"></slot>
    </div>
    <div ref="popoverRef">
      <div 
        class="com-popover-content" 
        :class="[`placement-${placement}`, `theme-${theme}`]"
        v-show="show"
        ref="popoverContent"
      >
        <div class="com-popover-arrow" ref="arrow"></div>
        <div class="com-popover-inner">
          <slot></slot>
        </div>
      </div>
      <div class="com-popover-mask" @click="close" v-if="show"></div>
    </div>
  </div>
</template>
<script lang="ts"> 
import { Component, Prop, Vue, Ref } from "vue-property-decorator";
import { Button } from "vant";

@Component({
  components: {
    [Button.name]: Button,
  },
})
export default class ComPopover extends Vue {
  show: boolean = false;

  @Prop({
    type: String,
    default: "bottom",
    validator: (value: string) => {
      return ["top", "top-start", "top-end", "bottom", "bottom-start", "bottom-end"].includes(value);
    }
  }) placement!: string;
  
  @Prop({
    type: String,
    default: "dark",
    validator: (value: string) => {
      return ["dark", "light"].includes(value);
    }
  }) theme!: string;
  
  @Prop({
    type: Number,
    default: 0
  }) offsetX!: number;
  
  @Prop({
    type: Number,
    default: 0
  }) offsetY!: number;
  
  toggle() {
    this.show = !this.show;
    if (this.show) {
      this.$nextTick(() => {
        this.calculatePosition();
      });
    }
  }
  
  open() {
    this.show = true;
    this.$nextTick(() => {
      this.calculatePosition();
    });
  }
  
  close() {
    this.show = false;
  }
  
  calculatePosition() {
    if (!this.$refs.popoverContent) return;
    
    const reference = this.$el.querySelector(".com-popover-reference") as HTMLElement;
    if (!reference) return;
    
    const referenceRect = reference.getBoundingClientRect();
    const popoverRect = (this.$refs.popoverContent as any).getBoundingClientRect();
    
    let left = 0;
    let top = 0;
    
    switch (this.placement) {
      case "top":
        left = referenceRect.left + referenceRect.width / 2 - popoverRect.width / 2 + this.offsetX;
        top = referenceRect.top - popoverRect.height - 10 + this.offsetY;
        break;
      case "top-start":
        left = referenceRect.left + this.offsetX;
        top = referenceRect.top - popoverRect.height - 10 + this.offsetY;
        break;
      case "top-end":
        left = referenceRect.right - popoverRect.width + this.offsetX;
        top = referenceRect.top - popoverRect.height - 10 + this.offsetY;
        break;
      case "bottom":
        left = referenceRect.left + referenceRect.width / 2 - popoverRect.width / 2 + this.offsetX;
        top = referenceRect.bottom + 10 + this.offsetY;
        break;
      case "bottom-start":
        left = referenceRect.left + this.offsetX;
        top = referenceRect.bottom + 10 + this.offsetY;
        break;
      case "bottom-end":
        left = referenceRect.right - popoverRect.width + this.offsetX;
        top = referenceRect.bottom + 10 + this.offsetY;
        break;
      default:
        left = referenceRect.left + referenceRect.width / 2 - popoverRect.width / 2 + this.offsetX;
        top = referenceRect.top - popoverRect.height - 10 + this.offsetY;
    }
    
    // 确保 popover 不会超出屏幕
    const viewportWidth = window.innerWidth;
    const viewportHeight = window.innerHeight;
    
    if (left < 0) {
      left = 24;
    } else if (left + popoverRect.width > viewportWidth) {
      left = viewportWidth - popoverRect.width - 10;
    }
    
    if (top < 0) {
      top = 24;
    } else if (top + popoverRect.height > viewportHeight) {
      top = viewportHeight - popoverRect.height - 10;
    }
    
    (this.$refs.popoverContent as any).style.left = `${left}px`;
    (this.$refs.popoverContent as any).style.top = `${top}px`;

    // 将 this.$refs.popoverContent插入body
    document.body.appendChild(this.$refs.popoverRef as any);

    // 确保箭头不会超出屏幕
    const contentRect = (this.$refs.content as any).getBoundingClientRect();
    console.log(contentRect);
    
    (this.$refs.arrow as any).style.left = `${ contentRect.left - left+ contentRect.width / 2}px`;
  }
}
</script>
<style lang="scss" scoped> 
.com-popover { 
  position: relative;
  display: inline-block;
}

.com-popover-reference {
  cursor: pointer;
}

.com-popover-mask {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 3000;
}

.com-popover-content {
  position: fixed;
  z-index: 3001;
  max-width: 100vw;
  min-width: 120px;
  padding: 8px 0;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.15);
  transition: all 0.3s ease;
  
  &.theme-dark {
    background-color: #333;
    color: #fff;
  }
  
  &.theme-light {
    background-color: #fff;
    color: #333;
  }
  
  .com-popover-arrow {
    position: absolute;
    width: 0;
    height: 0;
    border: 6px solid transparent;
    
    // &::after {
    //   content: '';
    //   position: absolute;
    //   width: 0;
    //   height: 0;
    //   border: 5px solid transparent;
    // }
  }
  
  // Top placements
  &.placement-top {
    .com-popover-arrow {
      bottom: -12px;
      left: 50%;
      transform: translateX(-50%);
      
      &::after {
        bottom: 1px;
        left: 50%;
        transform: translateX(-50%);
      }
    }
  }
  
  &.placement-top-start {
    .com-popover-arrow {
      bottom: -12px;
      left: 16px;
      
      &::after {
        bottom: 1px;
        left: 50%;
        transform: translateX(-50%);
      }
    }
  }
  
  &.placement-top-end {
    .com-popover-arrow {
      bottom: -12px;
      right: 16px;
      
      &::after {
        bottom: 1px;
        left: 50%;
        transform: translateX(-50%);
      }
    }
  }
  
  // Bottom placements
  &.placement-bottom {
    .com-popover-arrow {
      top: -12px;
      left: 50%;
      transform: translateX(-50%);
      
      &::after {
        top: 1px;
        left: 50%;
        transform: translateX(-50%);
      }
    }
  }
  
  &.placement-bottom-start {
    .com-popover-arrow {
      top: -12px;
      left: 16px;
      
      &::after {
        top: 1px;
        left: 50%;
        transform: translateX(-50%);
      }
    }
  }
  
  &.placement-bottom-end {
    .com-popover-arrow {
      top: -12px;
      right: 16px;
      
      &::after {
        top: 1px;
        left: 50%;
        transform: translateX(-50%);
      }
    }
  }
  
  // Theme colors for arrows
  &.theme-dark {
    &.placement-top .com-popover-arrow {
      border-top-color: #333;
      &::after {
        border-top-color: #333;
      }
    }
    &.placement-top-start .com-popover-arrow {
      border-top-color: #333;
      &::after {
        border-top-color: #333;
      }
    }
    &.placement-top-end .com-popover-arrow {
      border-top-color: #333;
      &::after {
        border-top-color: #333;
      }
    }
    &.placement-bottom .com-popover-arrow {
      border-bottom-color: #333;
      &::after {
        border-bottom-color: #333;
      }
    }
    &.placement-bottom-start .com-popover-arrow {
      border-bottom-color: #333;
      &::after {
        border-bottom-color: #333;
      }
    }
    &.placement-bottom-end .com-popover-arrow {
      border-bottom-color: #333;
      &::after {
        border-bottom-color: #333;
      }
    }
  }
  
  &.theme-light {
    &.placement-top .com-popover-arrow {
      border-top-color: #fff;
      &::after {
        border-top-color: #fff;
      }
    }
    &.placement-top-start .com-popover-arrow {
      border-top-color: #fff;
      &::after {
        border-top-color: #fff;
      }
    }
    &.placement-top-end .com-popover-arrow {
      border-top-color: #fff;
      &::after {
        border-top-color: #fff;
      }
    }
    &.placement-bottom .com-popover-arrow {
      border-bottom-color: #fff;
      &::after {
        border-bottom-color: #fff;
      }
    }
    &.placement-bottom-start .com-popover-arrow {
      border-bottom-color: #fff;
      &::after {
        border-bottom-color: #fff;
      }
    }
    &.placement-bottom-end .com-popover-arrow {
      border-bottom-color: #fff;
      &::after {
        border-bottom-color: #fff;
      }
    }
  }
  
  .com-popover-inner {
    padding: 8px 12px;
    font-size: 14px;
    line-height: 1.5;
  }
}
</style>