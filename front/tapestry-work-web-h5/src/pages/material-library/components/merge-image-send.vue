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
      <div class="popup-header">
        <div class="popup-header-name">请选择发送图片</div>
        <div
          class="popup-header-close"
          @click.stop="closePopup"
        ></div>
      </div>
      <div class="popup-content">
        <div class="radio-buttons">
          <van-checkbox
            v-model="checkAll"
            @click="checkImageAll"
          >
            <div class="content-title">全选</div>
            <template #icon="props">
              <div :class="props.checked ? 'active-radio' : 'radio'"></div>
            </template>
          </van-checkbox>
        </div>
        <van-row
          class="card-wrapper"
          :gutter="10"
        >
          <van-col
            :span="8"
            v-for="(item, index) in merageImgList"
            :key="index"
          >
            <div class="card">
              <div
                class="sel-tag"
                @click="checkedItem(item)"
              >
                <van-icon :name="chooseArr.findIndex((i) => i.id === item.id) > -1
                  ? 'checked'
                  : 'circle'
                  " />
              </div>
              <van-image
                lazy-load
                class="pic there-cols"
                :src="item.url"
                fit="contain"
                @click="previewImage(index)"
              />
              <div class="title">
                <vue-ellipsis
                  :visible-line="2"
                  :text="item.name.split('.')[0]"
                >
                </vue-ellipsis>
              </div>
            </div>
          </van-col>
        </van-row>
      </div>
      <div class="popup-footer">
        <!-- 组合发送 -->
        <div class="bottom-view clearfix">
          <div class="button-group">
            <div class="num-left">
              <div>已选：{{ chooseArr.length }}</div>
              <div
                class="max-num-left"
                v-if="chooseArr.length == 9"
              >最多可选9个</div>
            </div>
            <!-- todo:清空已选按钮 -->
            <div
              class="clear-btn"
              @click="clearChoose"
              v-if="chooseArr.length > 0"
            >
              清空已选
            </div>
            <div
              class="button-right"
              :class="[!chooseArr.length ? 'dis-btn' : '']"
              @click="combinationSend(chooseArr)"
            >
              {{ chooseArr.length > 1 ? '组合发送' : '发送' }}
            </div>
          </div>
        </div>
      </div>
    </div>
    <!-- todo: 图片预览弹窗 -->
    <image-preview-popup
      v-if="showImagePopup"
      :showPopup.sync="showImagePopup"
      :imgList="merageImgList"
      :chooseArr="chooseArr"
      :isFirstOpen="isFirstOpen"
      :showIndex="showIndex"
      @checkedItem="checkedItem"
      @closePopup="showImagePopup = false"
    />
  </van-popup>
</template>
<script lang="ts">
import { Component, Vue, PropSync, Prop, Watch } from "vue-property-decorator";
import { Popup, Image, Row, Col, Icon, Checkbox } from "vant";
import { intersectionBy } from "lodash"; // 引入lodash
// todo: 文本省略组件
import { VueEllipsis } from 'vue-ellipsis-component';
// todo:图片预览弹窗
import imagePreviewPopup from "./image-preview-popup.vue"
@Component({
  components: {
    [Popup.name]: Popup,
    [Image.name]: Image,
    [Row.name]: Row,
    [Col.name]: Col,
    [Icon.name]: Icon,
    [Checkbox.name]: Checkbox,
    VueEllipsis,
    imagePreviewPopup
  },
})
export default class MaterialGroup extends Vue {
  @PropSync("showPopup", { type: Boolean }) show;
  @Prop({ type: Array }) merageImgList;
  @Prop({ type: Array }) chooseArr;
  @Prop({ type: Boolean }) isFirstOpen;
  @Watch("chooseArr", { immediate: true, deep: true })
  onChooseArrChange(val: any) {
    this.checkAllFn()
  }
  sendDisabled = false;
  checkAll = false
  showImagePopup = false
  showIndex = 0
  // 记录本次已发送素材
  public get sentData(): any {
    return this.$store.getters["materialLibrary/getNowSentData"];
  }
  /**
   * @description: 判断全选
   * @return {*}
   */
  checkAllFn() {
    // 与chooseArr对比是否包含merageImgList全部数据
    const arr = intersectionBy(this.chooseArr, this.merageImgList, "id");
    if (arr.length === this.merageImgList.length) {
      this.checkAll = true
    } else {
      this.checkAll = false
    }
  }
  /**
   * @description: 选中的图片
   * @return {*}
   */
  checkedItem(item: any) {
    // 判断是否全部选中
    this.$emit("checkedItem", {
      ...item,
      type: 2,
      is_deal: 2,
    });
    this.checkAllFn()
  }
  /**
   * @description: 全选
   * @return {*}
   */
  checkImageAll() {
    this.checkAll = !this.checkAll
    console.log(this.checkAll)
    // 判断是否全部选中
    if (this.checkAll) {
      // 取消全选
      for (let i = 0; i < this.merageImgList.length; i++) {
        this.checkedItem(this.merageImgList[i])
      }
    } else {
      // 循环选中
      for (let i = 0; i < this.merageImgList.length; i++) {
        // 判断this.chooseArr是否已经有了
        const index = this.chooseArr.findIndex((f) => f.id === this.merageImgList[i].id)
        if (index === -1) {
          this.checkedItem(this.merageImgList[i])
        }
      }
    }
  }
  /**
   * @description: 发送图片
   * @return {*}
   */
  sendChatMessage(item: any) {
    this.$emit("closePopup");
    // 发送图片
    this.$emit("sendChatMessage", item);
  }
  /**
   * @description: 组合发送
   * @return {*}
   */
  combinationSend(arr: any) {
    if (arr.length === 0) {
      return
    }
    this.$emit("closePopup");
    // 组合发送
    this.$emit("combinationSend");
  }
  closePopup() {
    this.$emit("closePopup");
  }
  /**
   * @description: 清空已选
   * @return {*}
   */
  clearChoose() {
    this.$emit("clearChoose");
    this.checkAll = false
  }
  /**
   * @description: 预览图片
   * @return {*}
   */
  previewImage(index:number) {
    this.showIndex = index
    this.showImagePopup = true
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
  display: flex;
  justify-content: space-between;
  flex-direction: column;
  border-radius: 0;

  .radio-buttons {
    font-size: 12px;
    font-weight: 400;
    color: #000000;
    margin-bottom: 10px;
    margin-left: 5px;

    .content {
      &-title {
        font-size: 12px;
        font-weight: 400;
        color: #000000;
        line-height: 17px;
      }

      &-sub {
        font-size: 10px;
        font-weight: 400;
        color: #9E9E9E;
        line-height: 14px;
      }
    }

    .radio {
      position: relative;
      top: 2px;
      width: 15px;
      height: 15px;
      background: #FFFFFF;
      opacity: 0.4;
      border: 1px solid #000000;
      border-radius: 50%;
    }

    .active-radio {
      position: relative;
      top: 2px;
      width: 15px;
      height: 15px;
      border: 1px solid #000000;
      border-radius: 50%;

      &::before {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        content: '';
        width: 9px;
        height: 9px;
        background: #000000;
        border-radius: 50%;
      }
    }
  }

  .popup-container {
    height: 100%;
    overflow: hidden;

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
      padding: 0 23px;
      margin-bottom: 84px;

      .card-wrapper {
        .card {
          position: relative;
          margin-bottom: 10px;
          background: #ffffff;
          box-shadow: 0px 6px 12px 0px rgba(222, 224, 227, 0.5);
          padding-bottom: 7px;

          .sel-tag {
            position: absolute;
            width: 28px;
            height: 28px;
            display: flex;
            align-items: center;
            justify-content: center;
            box-shadow: -2px 2px 8px 0px rgba(0, 0, 0, 0.1);
            border-radius: 0px 0px 0px 10px;
            right: 0;
            top: -1px;
            line-height: 22px;
            text-align: center;
            background-color: #ffffff;
            z-index: 1;

            .van-icon {
              font-size: 16px;
              vertical-align: middle;
              color: #4a4a4a;
            }
          }

          .file-sent-tag {
            top: inherit;
            bottom: 86px;
            left: 0;
          }

          .pic {
            position: relative;
            margin-bottom: 8px;
            background: #f0f0f0;

            &.video {
              background: #d8d8d8;
              font-size: 32px;
            }
          }

          .pic-num {
            position: absolute;
            top: 15px;
            left: 15px;
            display: flex;
            align-items: center;
            justify-content: center;
            width: 36px;
            height: 19px;
            background: rgba(0, 0, 0, 0.5);
            border-radius: 9px;
            font-size: 10px;
            font-weight: 400;
            color: #FFFFFF;
            z-index: 999;
            padding: 0 5px;
          }

          .two-cols {
            width: 150px;
            height: 150px;
            margin: 7px;
          }

          .there-cols {
            width: 94px;
            height: 94px;
            margin: 6px 5px 0 5px;
          }

          .title {
            position: relative;
            font-size: 12px;
            color: #000000;
            height: 32.6px;
            overflow: hidden;
            padding: 0 11px;
            text-align: center;
          }

          .btn {
            height: 32px;
            display: flex;
            justify-content: center;
            align-items: center;
            cursor: pointer;
            border: 1px solid #000000;
            color: #000000;
            border-radius: 0;
            font-size: 13px;
            margin: 14px 7px 0px;

            &.disabled {
              font-size: 12px;
              background: #f4f7fb;
              color: #7f7f7f;
              border-color: transparent;
            }
          }
        }

        .fileCard {
          height: 209px !important;

          .filePic {
            padding: 25px 0 17px;

            .fileImg {
              display: block;
              width: 71px;
              height: 94px;
              margin: 0 auto;
            }
          }
        }
      }
    }

    .popup-footer {

      .bottom-view {
        position: fixed;
        bottom: -1px;
        width: 100%;
        height: 50px;
        left: 0;
        padding: 8px;
        background-color: #ffffff;
        font-size: 14px;
        z-index: 9;
        box-shadow: 0px 1px 6px 0px rgba(0, 0, 0, 0.15);

        .button-group {
          position: relative;
          display: flex;
          align-items: center;
          width: 100%;
          height: 100%;

          .num-left {
            position: absolute;
            left: 0;

            div {
              margin: 0 4px;
              font-size: 16px;
            }

            .max-num-left {
              font-size: 10px;
              color: #E02020;
            }
          }

          .clear-btn {
            position: absolute;
            right: 180px;
            font-size: 12px;
            font-weight: 500;
            color: #000000;
            text-decoration: solid underline #000000;
          }

          .button-right {
            position: absolute;
            right: 0;
            width: 165px;
            height: 34px;
            color: #ffffff;
            line-height: 34px;
            text-align: center;
            background-color: #000;

            &.dis-btn {
              background-color: #cccccc;
            }
          }
        }
      }
    }
  }
}
</style>
