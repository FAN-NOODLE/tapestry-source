<template>
  <div class="page-loading">
    <div
      class="load-img"
    >
      <img
        class="image"
        src="https://cdn-cos.coach.com.cn/uploads/2023/v3.2.0.1/loading.png"
      />
    </div>
    <div class="text">亲，数据努力加载中，请稍后~</div>
  </div>
</template>
<script lang="ts">
import { Prop, Component, Vue } from "vue-property-decorator";
@Component({
  components: {},
})
export default class extends Vue {

}
</script>

<style lang="sass" scoped>
.page-loading
  position: fixed
  z-index: 100
  top: 0
  right: 0
  bottom: 0
  left: 0
  display: flex
  flex-direction: column
  width: 243px
  height: 146px
  margin: auto
  padding: 40px 15px 0 15px
  border-radius: 8px
  background: #fff
  box-shadow: 0 0 7px 0 rgba(0,0,0,.1)

  align-items: center
  .load-img
    transform: scale(0.8)
    .image
      width: 40px
      height: 40px
      animation: bounceIn 1s infinite
  .text
    font-size: 12px
    margin-top: 23px
    line-height: 17px
    color: #7f7f7f 
@keyframes bounceIn
  0%,
  20%,
  40%,
  60%,
  80%,
  to
    animation-timing-function: cubic-bezier(.215, .61, .355, 1)
  0%
    transform: rotate(30deg)
  10%
    transform: rotate(60deg)
  20%
    transform: rotate(90deg)
  30%
    transform: rotate(120deg)
  40%
    transform: rotate(150deg)
  50%
    transform: rotate(180deg)
  60%
    transform: rotate(210deg)
  70%
    transform: rotate(240deg)
  80%
    transform: rotate(270deg)
  90%
    transform: rotate(300deg)
  100%
    transform: rotate(330deg)
  to
    transform: rotate(360deg)

</style>
