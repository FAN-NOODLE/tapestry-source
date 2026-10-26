<template>
  <!--分组弹窗 -->
  <van-popup
    v-model="visible"
    :close-on-click-outside="false"
    :close-on-click-action="false"
    position="bottom"
    class="material-popup"
  >
    <div class="popup-container">
      <div class="popup-header">
        <div class="popup-header-name">请选择发送图片</div>
        <div
          class="popup-header-close"
          @click.stop="handleClose"
        ></div>
      </div>
      <div class="popup-content">
        <div class="radio-buttons">
          <van-checkbox
            :value="isCheckedAll"
            @click="handleSelectAll"
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
            v-for="(item, index) in list"
            :key="index"
          >
            <div class="card">
              <div
                class="sel-tag"
                @click="handleSelect(item)"
              >
                <van-icon :name="isChecked(item) ? 'checked' : 'circle'" />
              </div>
              <van-image
                lazy-load
                class="pic there-cols"
                :src="item.url"
                fit="contain"
                @click="handlePreview(index)"
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
              <div>已选：{{ selection.length }}</div>
              <div
                class="max-num-left"
                v-if="selection.length == 9"
              >最多可选9个</div>
            </div>
            <!-- todo:清空已选按钮 -->
            <div
              class="clear-btn"
              @click="handleClear"
              v-if="selection.length > 0"
            >
              清空已选
            </div>
            <div
              class="button-right"
              :class="[!selection.length ? 'dis-btn' : '']"
              @click="handleSend(selection)"
            >
              {{ selection.length > 1 ? '组合发送' : '发送' }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Vue, PropSync, Prop, Watch } from "vue-property-decorator";
import { Popup, Image, Row, Col, Icon, Checkbox } from "vant";
import { intersectionBy } from "lodash"; // 引入lodash
// todo: 文本省略组件
import { VueEllipsis } from 'vue-ellipsis-component';
// todo:图片预览弹窗
@Component({
  components: {
    [Popup.name]: Popup,
    [Image.name]: Image,
    [Row.name]: Row,
    [Col.name]: Col,
    [Icon.name]: Icon,
    [Checkbox.name]: Checkbox,
    VueEllipsis,
  },
})
export default class SelectMergeImagePopup extends Vue {
  @PropSync("show", { type: Boolean }) visible;
  @Prop({ type: Object, default: () => ({}) }) item: any;
  @Prop({ type: Array }) selection: any[];
  @Prop({ type: Boolean }) isFirstOpen;
  
  sendDisabled = false;
  
  // 记录本次已发送素材
  get sentData(): any {
    return this.$store.getters["materialLibrary/getNowSentData"];
  }

  /** 是否是组合图片 */
  get isMergeImage() {
    if(!this.item) return false
    return this.item.merge_image && this.item.merge_image.length > 1
  }

  /** 素材列表 */
  get list () {
    if(!this.item) return []
    if(this.isMergeImage) {
      return this.item.merge_image
    } else {
      return [this.item]
    }
  }

  get status() {
    const selectedIds = this.selection.map(item => item.id)
    
    const checkedLength = this.list.filter(item => selectedIds.includes(item.id))
    if(checkedLength.length === this.list.length) {
      return 2
    } else if(checkedLength.length > 0) {
      return 1
    } else {
      return 0
    }
  }

  /** 是否 全选 */
  get isCheckedAll() {
    return this.status === 2
  }

  get isChecked() {
    return (item)=> {
      const ids = this.selection.map(item => item.id)
      return ids.includes(item.id)
    }
  }

  /**
   * @description: 选中的图片
   * @return {*}
   */
  handleSelect(item: any) {
    const ids = this.selection.map(item => item.id)
    // 判断是否全部选中
    this.$emit("select", {
      status: ids.includes(item.id) ? 2: 0,
      item: {type: 2, ...item},
      isMergeImage: false
    });
  }
  /**
   * @description: 全选
   * @return {*}
   */
  handleSelectAll() {
    this.$emit('select', {
      item: this.item, 
      status: this.status,
      isMergeImage: this.isMergeImage
    })
  }
  /**
   * @description: 组合发送
   * @return {*}
   */
  handleSend(arr: any) {
    if (arr.length === 0) {
      return
    }
    this.visible = false
    // 组合发送
    this.$emit("send");
  }
  handleClose() {
    this.visible = false
    this.$emit("close");
  }
  /**
   * @description: 清空已选
   * @return {*}
   */
  handleClear() {
    this.$emit("clear");
  }
  /**
   * @description: 预览图片
   * @return {*}
   */
  handlePreview(index: number) {
    this.$emit('preview', {
      index,
      item: this.item
    })
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
