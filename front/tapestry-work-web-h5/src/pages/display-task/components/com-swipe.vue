<template>
  <van-swipe v-if="list.length" :loop="false" :initial-swipe="initialSwipe" @change="handleChange">
    <van-swipe-item v-for="(item, index) in list" :key="index">
      <com-swipe-item :item="item" :index="index"></com-swipe-item>
    </van-swipe-item>
    <template #indicator>
      <div class="custom-indicator">{{ current + 1 }}/{{list.length}}</div>
    </template>
  </van-swipe>
</template>
<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
import { Swipe, SwipeItem } from 'vant';
import ComSwipeItem from './com-swipe-item.vue'

@Component({
  components: {
    [Swipe.name]: Swipe,
    [SwipeItem.name]: SwipeItem,
    ComSwipeItem
  }
})
export default class ComSwipe extends Vue {
  @Prop({ type: Array, required: true }) list: IImageItem[] 
  @Prop({ type: Number, default: 0 }) current: number;
  @Prop({ type: Number, default: 0 }) initialSwipe: number;

  handleChange(index: number) {
    // this.current = index;
    console.log(12312312);
    
    this.$emit('change', index)
  }
}
</script>
<style lang="scss" scoped>

.custom-indicator {
  position: absolute;
  left: 50%;
  bottom: 20px;
  transform: translateX(-50%);
  background: rgba(0, 0, 0, 0.6);
  font-size: 14px;
  width: 83px;
  height: 27px;
  border-radius: 14px;
  text-align: center;
  color: #fff;
  line-height: 27px;
}
</style>