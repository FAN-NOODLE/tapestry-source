<template>
  <div :class="selfItem.index === cardLine ? 'flex-s-c' : ''">
    <div class="line" v-if="selfItem.index === cardLine"></div>
    <div
      class="card"
      :style="
        `background: linear-gradient(360deg, ${selfItem.background[0]} 0%, ${selfItem.background[1]} 100%);`
      "
    ></div>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
@Component
export default class Card extends Vue {
  @Prop() item!: object;
  private selfItem = {};
  private cardLine = 6; //有线卡片标志位
  private cardWrap = 3; //卡片换行标志位
  created() {
    this.selfItem = this.item || {};
  }
  mounted() {
    console.log(123);
  }
}
</script>
<style lang="scss" scoped>
.flex-s-c {
  display: flex;
  justify-content: flex-start;
  align-items: center;
}
.line {
  height: 147px;
  width: 1px;
  background: #ebeef5;
  margin-right: 30px;
}
.card {
  background: #ccc;
  border-radius: 10px;
  width: 100%;
  height: 147px;
}
</style>
