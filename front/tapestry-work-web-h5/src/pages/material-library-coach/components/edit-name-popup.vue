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
      <div class="popup-footer">
        <div class="input-buttons">
          <van-field
            v-model="fileName"
            class="van-input"
            type="text"
            clearable
            ref="fileNameInput"
          />
          <div class="button" @click="enterFileName">
            确定
          </div>
        </div>
      </div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Vue,PropSync, Prop,Watch } from "vue-property-decorator";
import { Popup, Image, Field } from "vant";
@Component({
  components: {
    [Popup.name]: Popup,
    [Image.name]: Image,
    [Field.name]: Field,
  },
})
export default class MaterialGroup extends Vue {
  @PropSync("showPopup", { type: Boolean }) show;
  @Watch("show")
  onShowChange(val:Boolean){
    if(val){
      this.$nextTick(() => {
        (this.$refs.fileNameInput as HTMLElement)?.focus();
      });
    }
  }
  @Prop(Number) index!: any;
  @Prop(Object) file!: any;
  @Watch('file', { immediate: true, deep: true })
  onFileChange(val: any) {
    this.fileName = val.name
  }
  fileName = "";
  closePopup() {
    this.$emit("closePopup");
  }
  /**
   * @description: 确定
   * @return {*}
   */
  enterFileName() {
    // 验证文件名只能输入汉字，字母，数字，下划线，中划线
    const reg = /^[\u4e00-\u9fa5_a-zA-Z0-9-]+$/;
    if (!reg.test(this.fileName)) {
      this.$toast("文件名输入不合法!");
      return;
    }
    this.$emit("enterFileName", {
      name: this.fileName,
      index: this.index,
    });
  }
}
</script>
<style lang="scss" scoped>
::-webkit-scrollbar {
  display: none;
  /* Chrome Safari */
}

.material-popup {
  padding: 0;
  border-radius: 0;

  .popup-container {
    .popup-footer {
      padding: 15px;

      .input-buttons {
        display: flex;
        justify-content: space-between;

        .van-input {
          width: 261px;
          height: 44px;
          background: #F4F7FB;
          font-size: 14px;
        }

        .button {
          width: 84px;
          height: 44px;
          line-height: 44px;
          font-size: 14px;
          color: #ffffff;
          text-align: center;
          background-color: #000;
        }
      }
    }
  }
}
</style>
