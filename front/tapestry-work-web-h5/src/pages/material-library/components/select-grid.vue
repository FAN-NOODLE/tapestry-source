<template>
  <div class="grid">
    <div 
      v-for="(item, index) in [0, 1]" 
      :key="index" 
      :class="['grid-item', { active: value==item}]"
      @click="handleGideChange(item)"></div>
    <div v-if="showFirstTips" class="grid__tips" @click="handleHideTips"></div>
  </div>
</template>
<script lang="ts">
import { Component, Vue, PropSync, Prop, Watch } from "vue-property-decorator";

@Component({
  model: {
    prop: "value",
    event: "change"
  }
})
export default class SelectGrid extends Vue {
  @Prop({default: 0, type: Number}) value: number
  @Prop({default: false, type: Boolean}) showFirstTips: boolean

  /**
   * @description 隐藏首次提示
   */
    handleHideTips() {
    this.showFirstTips = false
    this.$emit('hide')
  }
  /**
   * @description 切换宫格
   */
  handleGideChange(value) {
    this.$emit('change', value)
    this.handleHideTips()
  }
}
</script>

<style lang="scss" scoped>
.grid {
  height: 20px;
  padding: 0 20px;
  display: flex;
  align-items: center;
  margin-top: 15px;
  position: relative;

  &-item {
    width: 20px;
    height: 20px;
    margin-right: 30px;
    position: relative;

    &:nth-child(1) {
      &::before {
        content: "";
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        width: 20px;
        height: 20px;
        z-index: 2;
        opacity: 1;
        background: url("https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/icon_grade_4.png") no-repeat center center/100%;
      }
      &::after {
        content: "";
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        width: 20px;
        height: 20px;
        z-index: 1;
        opacity: 0;
        background: url("https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/icon_grade_4_1.png") no-repeat center center/100%
      }
    }
    &:nth-child(2) {
      &::before {
        content: "";
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        width: 21px;
        height: 23px;
        z-index: 2;
        opacity: 1;
        background: url("https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/icon_grade_9.png") no-repeat center center/100%;
      }
      &::after {
        content: "";
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        width: 21px;
        height: 23px;
        z-index: 1;
        opacity: 0;
        background: url("https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/icon_grade_9_1.png") no-repeat center center/100%;
      }
    }
      
    &.active {
      &::before {
        z-index: 1;
        opacity: 0
      }
      &::after {
        z-index: 2;
        opacity: 1;
      }
    }
  }
  .grid__tips {
    z-index: 4;
    position: absolute;
    top: 15px;
    left: 60px;
    width: 182px;
    height: 69px;
    background: url("https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/icon_popup_tips.png") no-repeat center center/100%
  }
}
</style>