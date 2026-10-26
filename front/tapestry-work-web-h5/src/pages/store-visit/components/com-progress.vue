<template>
  <div 
    class="progress-circle" 
    :style="{
      '--size': `${size}px`, 
      '--width': width,
      '--color': color,
      '--background-color': backgroundColor,
    }">
    <svg viewBox="0 0 36 36">
      <circle class="background" cx="18" cy="18" r="15" />
      <circle class="progress" cx="18" cy="18" r="15" transform="rotate(-90, 18, 18)" />
    </svg>
    <!-- <div class="progress-circle__background"></div>
    <div class="progress-circle__progress" :style="{'--progress': `${360/100 * progress}deg`}"></div> -->
  </div>
</template>

<script lang="ts">
import { Component, Model, Prop, Vue, Watch } from "vue-property-decorator";

@Component({

})
export default class ComProgress extends Vue {
  @Prop({ type: Number, default: 0 }) value: number;
  @Prop({ type: Number }) size: number;
  @Prop({ type: Number, default: 3 }) width: number;
  @Prop({ type: String, default: '#FFFFFF' }) color: string;
  @Prop({ type: String, default: '#DEE2E6' }) backgroundColor: string;

  @Watch('value', { immediate: true })
  onProgressChange() {
    this.$nextTick(() => {
      this.update();
    })
  }

  update() {
    const progressCircle = document.querySelector('.progress-circle .progress') as SVGPathElement;
    if(!progressCircle) return
    const totalLength = progressCircle.getTotalLength();

    // 设置 stroke-dasharray 的总长度
    progressCircle.style.strokeDasharray = totalLength+'';

    // 更新进度
    function updateProgress(value) {
      const offset = totalLength * (1 - value / 100);
      progressCircle.style.strokeDashoffset = offset+'';
    }

    // 示例：设置进度为 75%
    updateProgress(this.value);

  }
}
</script>

<style lang="scss" scoped>
.progress-circle {
  width: var(--size, 30px);
  height: var(--size, 30px);
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.progress-circle svg {
  width: 100%;
  height: 100%;
  vertical-align: middle;
}

.progress-circle circle {
  fill: none;
  stroke-width: var(--width);
  stroke-linecap: round;
}

.progress-circle .background {
  stroke: var(--background-color);
}

.progress-circle .progress {
  stroke: var(--color);
  stroke-dasharray: 100;
  stroke-dashoffset: 100;
  transition: stroke-dashoffset 0.5s ease;
}

</style>