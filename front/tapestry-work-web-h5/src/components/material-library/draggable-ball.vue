<template>
    <div 
        class="collect-ball" 
        ref="draggable" 
        :style="styleObject"
        @touchstart="onTouchStart"
        @touchmove="onTouchMove"
        @touchend="onTouchEnd"
        @click="clickCallback">
        <div class="collect-icon"><slot name="icon"></slot></div>
        <div class="collect-text" v-show="!isScrolling"><slot name="text"></slot></div>
    </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'

@Component
export default class DraggableBall extends Vue {
    //接收 width 和 height
    @Prop({ type: Number, default: 115 }) width!: number
    @Prop({ type: Number, default: 32 }) height!: number
    
    /** ---------------------- 页面滚动事件 --------------------------- */
    isScrolling = false; // 是否正在滚动
    scrollTimeout: number | null = null; // 滚动延迟计时器
    // 监听滚动事件
    onScroll() {
        // console.log('scroll')
        this.isScrolling = true
        this.boxWidth = 34;
        this.borderRadius = '34px';
        this.posX = this.posX + this.boxWidth / 2 < this.screenMiddle ? 10 : window.innerWidth - this.boxWidth - 10
        if (this.scrollTimeout) {
            clearTimeout(this.scrollTimeout)
        }
        this.scrollTimeout = window.setTimeout(() => {
        this.isScrolling = false
        this.boxWidth = this.width;
        this.posX = this.posX + this.boxWidth / 2 < this.screenMiddle ? 0 : window.innerWidth - this.boxWidth
        if(this.posX === 0) {
            this.borderRadius = '0 34px 34px 0'
        } else {
            this.borderRadius = '34px 0 0 34px'
        }
        }, 300)
    }

    /** ---------------------- 拖拽事件 --------------------------- */
    boxWidth = 0
    boxHeight = 0
    borderRadius = ''
    startX = 0
    startY = 0
    posX = 0
    posY = 0
    
    isMoving = false
    screenMiddle = window.innerWidth / 2

    get styleObject() {
        return {
            position: 'fixed',
            left: `${this.posX}px`,
            top: `${this.posY}px`,
            width: `${this.boxWidth}px`,
            height: `${this.boxHeight}px`,
            borderRadius: this.borderRadius,
            zIndex: 999,
            transition: this.isMoving ? 'none' : 'all 0.3s ease',
            userSelect: 'none',
        }
    }
    onTouchStart(e: TouchEvent) {
        const touch = e.touches[0]
        this.startX = touch.clientX - this.posX
        this.startY = touch.clientY - this.posY
        this.isMoving = true
    }

    onTouchMove(e: TouchEvent) {
        if (!this.isMoving) return

        const touch = e.touches[0]
        let newX = touch.clientX - this.startX
        let newY = touch.clientY - this.startY

        const maxX = window.innerWidth - this.boxWidth
        const maxY = window.innerHeight - this.boxHeight

        // 限制不能超出屏幕
        newX = Math.max(0, Math.min(newX, maxX))
        newY = Math.max(0, Math.min(newY, maxY))

        this.posX = newX
        this.posY = newY

        document.body.style.overflow = 'hidden' // 禁止滚动
        document.body.style.touchAction = 'none'
    }

    onTouchEnd() {
        this.isMoving = false
        this.posX = this.posX + this.boxWidth / 2 < this.screenMiddle ? 0 : window.innerWidth - this.boxWidth

        if(this.posX === 0) {
            this.borderRadius = '0 34px 34px 0'
        } else {
            this.borderRadius = '34px 0 0 34px'
        }

        document.body.style.overflow = '' // 恢复滚动
        document.body.style.touchAction = ''
    }
    clickCallback() {
        this.$emit('click')
    }

    async mounted() {
        this.boxWidth =  this.width
        this.boxHeight = this.height
        this.borderRadius = `${this.boxHeight}px 0 0 ${this.boxHeight}px`
        this.startX = window.innerWidth - this.boxWidth
        this.startY = window.innerHeight - this.boxHeight - 200
        this.posX = this.startX
        this.posY = this.startY
        window.addEventListener('scroll', this.onScroll)
    }
    beforeDestroy(): void {
        window.removeEventListener('scroll', this.onScroll)
    }
}
</script>

<style lang="scss" scoped>
// 拖拽球，悬浮球
.collect-ball {
    position: fixed;
    bottom: 100px;
    right: 0px;
    height: 32px;
    background: #ffffff;
    border-radius: 34px 0 0 34px;
    display: flex;
    justify-content: center;
    transition: all 0.2s;
    align-items: center;
    z-index: 10;
    font-size: 12px;
    color: #000000;
    box-shadow: -1px 2px 7px 0px rgba(0,0,0,0.2);;
    .collect-icon {
        width: 15px;
        height: 32px;
        // background: url('~@/assets/material-library/v5.1.5/m0529-fanhui.png') no-repeat center center/8px 12px;
        display: flex;
        justify-content: center;
        align-items: center;
        .van-icon {
            display: block;
        }
    }
    .collect-text {
        transition: all 0.4s;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        text-align: center;
        margin-left: 5px;
    }
    }
</style>
