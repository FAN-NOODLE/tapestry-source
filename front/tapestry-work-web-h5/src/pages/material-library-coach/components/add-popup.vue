<template>
  <div>
    <div
      class="add-icon"
      :class="[!isApplets ? '' : 'bottom']"
      @click="visible=true"
    ></div>
    <van-popup
      v-model="visible"
      position="bottom"
    >
      <header class="add-modal-header">请选择新增素材类型</header>
      <ul>
        <li
          class="cate"
          v-for="(item, index) in list"
          :key="index"
          @click="handleAdd(item.route)"
        >
          <img :src="item.link" />
        </li>
      </ul>
    </van-popup>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import { Popup } from "vant";
@Component({
  components: {
    [Popup.name]: Popup,
  },
})
export default class AddPopup extends Vue {
  visible = false;

  list = [
    {
      link: require("@/assets/material-library/text.png"),
      route: "text-library",
    },
    {
      link: require("@/assets/material-library/image.png"),
      route: "image-library",
    },
    {
      link: require("@/assets/material-library/video.png"),
      route: "video-library",
    },
  ];

  /**
   * 是否来自消息通知
   */
   get isApplets() {
    return this.$route.query.is_applets;
  }

  // 添加跳转标识
  handleAdd(path) {
    this.$router.push({
      name: path,
    });
  }
}
</script>
<style lang="scss" scoped>

.add-icon {
  width: 36px;
  height: 36px;
  background: url("../../../assets/material-library/add-icon.png") no-repeat;
  background-size: 100% 100%;
  position: fixed;
  bottom: 126px;
  right: 12px;
  z-index: 2;

  &.bottom {
    bottom: 76px;
  }
}
.add-modal-header {
  height: 20px;
  font-size: 14px;
  font-family: PingFangSC-Medium, PingFang SC;
  font-weight: 700;
  color: #000000;
  line-height: 20px;
  margin: 16px auto 20px;
  text-align: center;
}

.cate {
  padding: 4px 12px 0;
  img {
    width: 100%;
    height: auto;
  }
}
</style>