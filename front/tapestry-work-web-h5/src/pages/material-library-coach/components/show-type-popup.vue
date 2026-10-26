<template>
  <!--分组弹窗 -->
  <van-popup
    v-model="show"
    @close="closePopup"
    :close-on-click-outside="false"
    :close-on-click-action="false"
    round
    position="bottom"
    class="material-popup"
  >
    <div class="popup-container">
      <van-sticky>
        <div class="popup-header">
          <div class="popup-header-name">展示方式</div>
          <div
            class="popup-header-close"
            @click="closePopup"
          ></div>
        </div>
      </van-sticky>
      <div class="popup-content">
        <van-image
          style="margin-top: 0;"
          class="type-image"
          fit="cover"
          :src="showtype1"
        />
        <van-image
          class="type-image"
          fit="cover"
          :src="showtype2"
        />
      </div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Vue, PropSync } from "vue-property-decorator";
import { Popup, Image, Sticky } from "vant";
@Component({
  components: {
    [Popup.name]: Popup,
    [Image.name]: Image,
    [Sticky.name]: Sticky,
  },
})
export default class MaterialGroup extends Vue {
  @PropSync("showPopup", { type: Boolean }) show;
  showtype1 = `https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/showType1.png?time=${new Date().getTime()}`
  showtype2 =  `https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/showType2.png?time=${new Date().getTime()}` 
  closePopup() {
    this.$emit("closePopup");
  }

}
</script>
<style lang="scss" scoped>
.material-popup {
  padding: 0;
  border-radius: 0;

  .popup-container {
    .popup-header {
      height: 44px;
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 0 16px;

      .popup-header-name {
        width: 100%;
        text-align: center;
        font-size: 16px;
        font-weight: 500;
        color: #333333;
        margin-left: 20px;
      }

      .popup-header-close {
        width: 20px;
        height: 20px;
        background: url("https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/icon_close.png") no-repeat center;
        background-size: 20px 20px;
      }
    }

    .popup-content {
      padding: 0 15px;
      height: 500px;
      overflow-y: auto;

      .type-image {
        margin-top: 5px;
        height: 340px;
      }
    }
  }
}</style>
