<template>
  <div class="video-library">
    <!--素材分组-->
    <material-group
      :showPopup.sync="showPopup"
      @cancel="closeGroup"
      @confirm="setGroup"
    />

    <!-- 素材标签 -->
    <material-tag
      :showPopup.sync="addTagPopup"
      :ids="formData.tags.ids"
      :names="formData.tags.names"
      @cancel="closeTag"
      @confirm="setTag"
    ></material-tag>

    <!-- 素材有效期 -->
    <material-time
      :showPopup.sync="timeRangePopup"
      @cancel="closeRange"
      @confirm="setRange"
    ></material-time>
    <!-- todo:展示类型 -->
    <show-type-popup
      v-if="isCoach"
      :showPopup.sync="showTypePopup"
      @closePopup="showTypePopup = false"
    />
    <!-- 图片素材表单 -->
    <div class="form">
      <div class="item">
        <div class="label">素材分组</div>
        <div
          class="content"
          :class="[formData.group.name ? '' : 'empty']"
          @click="openGroup"
        >
          <p>
            {{ formData.group.name ? formData.group.name : "请选择素材分组" }}
          </p>
          <i class="arrow"></i>
        </div>
      </div>

      <div class="item">
        <div class="label">素材标签</div>
        <div
          class="content"
          :class="[formData.tags.ids.length > 0 ? '' : 'empty']"
          @click="openTag"
        >
          <p>
            {{ formData.tags.ids.length > 0 ? tagsName : "请选择素材标签" }}
          </p>
          <i class="arrow"></i>
        </div>
      </div>

      <div class="item">
        <div class="label">素材有效期 <i class="require">*</i></div>
        <div
          class="content2"
          @click="openRange"
        >
          <span :class="[formData.timeRange[0] ? '' : 'empty']">
            <template v-if="formData.timeRange[0]">
              {{ formData.timeRange[0] | timeFormat }}
            </template>
            <template v-else> 请选择开始时间 </template>
          </span>
          <i>-</i>
          <span :class="[formData.timeRange[1] ? '' : 'empty']">
            <template v-if="formData.timeRange[1]">
              {{ formData.timeRange[1] | timeFormat }}
            </template>
            <template v-else> 请选择结束时间 </template>
          </span>
        </div>
      </div>
      <div class="item">
        <div class="label">素材内容 <i class="require">*</i></div>
        <div class="content3">
          <!-- todo:展示方式 -->
          <div
            class="show-type"
            v-if="isCoach"
          >
            <div class="left-label">
              <div class="left-label-name">展示方式</div>
              <img
                class="left-label-icon"
                src="https://cdn-cos.coach.com.cn/uploads/2023/v3.5.0/icon_tips.png"
                alt=""
                @click="showTypePopup = true"
              >
            </div>
            <div class="radio-buttons">
              <van-radio-group
                direction="horizontal"
                v-model="formData.show_type"
                @change="showTypeChange"
              >
                <van-radio
                  name="1"
                  class="content"
                >
                  <div class="content-title">单张展示</div>
                  <div class="content-sub">可传20张</div>
                  <template #icon="props">
                    <div :class="props.checked ? 'active-radio' : 'radio'"></div>
                  </template>
                </van-radio>
                <van-radio
                  name="2"
                  class="content"
                >
                  <div class="content-title">合并展示</div>
                  <div class="content-sub">可传9张</div>
                  <template #icon="props">
                    <div :class="props.checked ? 'active-radio' : 'radio'"></div>
                  </template>
                </van-radio>
              </van-radio-group>
            </div>
          </div>
          <multi-image-upload
            :maxLength="formData.show_type === '1' ? 20 : 9"
            @change="uploadChange"
          />
        </div>
      </div>
    </div>
    <div class="btn-group">
      <div
        class="cancel"
        @click="cancelAdd"
      >取消</div>
      <div
        class="submit"
        :class="[validate ? '' : 'disabled']"
        @click="handleSubmit"
      >
        上传
        <van-loading
          class="inline-block"
          v-if="saving"
          size="14px"
          color="#fff"
        />
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { Field, Loading, Popup, RadioGroup, Radio } from "vant";
import materialGroup from "./components/material-group.vue";
import materialTag from "./components/material-tag.vue";
import materialTime from "./components/material-time.vue";
import multiImageUpload from "./components/multi-image-upload.vue";
import { saveMaterial } from "@/api/marketing-promotion/index";
// todo: 品牌判断
import Helper from "@/utils/helper";
// todo: 展示类型弹窗
import ShowTypePopup from "./components/show-type-popup.vue";
@Component({
  components: {
    [Loading.name]: Loading,
    [Popup.name]: Popup,
    [Field.name]: Field,
    [RadioGroup.name]: RadioGroup,
    [Radio.name]: Radio,
    "material-group": materialGroup,
    "material-tag": materialTag,
    "material-time": materialTime,
    "multi-image-upload": multiImageUpload,
    ShowTypePopup
  },
  filters: {
    timeFormat: (value: string) => {
      return value.split(" ")[0];
    },
  },
})
export default class extends Vue {
  [x: string]: any;
  isCoach = Helper.isCoach();
  showTypePopup = false;
  /** 图片上传列表状态全部成功 */
  finished = false;
  saving = false;
  /**素材分组显示隐藏*/
  showPopup = false;
  /** 添加素材标签 */
  addTagPopup = false;
  /** 添加素材有效期 */
  timeRangePopup = false;
  formData = {
    /**素材分组 */
    group: {
      id: "",
      name: "",
    },
    /** 素材标签 */
    tags: {
      ids: [],
      names: [],
    },
    /** 素材有效期 */
    timeRange: [],
    /** 图片 */
    sources: [],
    //todo: 展示方式
    show_type: '1'
  };
  created() {
    // 素材有效期默认当前时间+1年
    const now = new Date();
    const year = now.getFullYear();
    let month: any = now.getMonth() + 1;
    let day: any = now.getDate();
    if (month < 10) {
      month = `0${month}`;
    }
    if (day < 10) {
      day = `0${day}`;
    }
    const nextYear = year + 1;
    this.formData.timeRange = [
      `${year}-${month}-${day} 00:00:00`,
      `${nextYear}-${month}-${day} 23:59:59`,
    ];
  }
  /**表单必填是否完成 */
  get validate() {
    if (
      this.formData.sources.length === 0 ||
      this.formData.timeRange.length < 2 ||
      !this.finished
    ) {
      return false;
    }
    return true;
  }
  /**标签名称 */
  get tagsName() {
    let tagName = "";
    this.formData.tags.names.forEach((name) => {
      tagName += `${name};`;
    });
    return tagName.slice(0, -1);
  }
  /**打开素材分组 */
  openGroup() {
    this.showPopup = true;
  }
  /** 关闭选择素材分组 */
  closeGroup() {
    this.showPopup = false;
  }
  /**选择素材分组 */
  setGroup(e) {
    console.log("setGroup", e);
    // const regx = />/g;
    let { group, groupName } = e;
    this.formData.group.id = group;
    this.formData.group.name = groupName;
    this.closeGroup();
  }
  /**打开标签 */
  openTag() {
    this.addTagPopup = true;
  }
  /**关闭标签选择 */
  closeTag() {
    this.addTagPopup = false;
  }
  /**确认标签选择 */
  setTag(e) {
    const { names, ids } = e;
    this.formData.tags.ids = ids;
    this.formData.tags.names = names;
    this.closeTag();
  }
  /**打开有效期 */
  openRange() {
    this.timeRangePopup = true;
  }
  /**取消有效期 */
  closeRange() {
    this.timeRangePopup = false;
  }
  /** 设置有效期 */
  setRange(e) {
    const { start, end } = e;
    console.log(e);
    this.formData.timeRange = [start, end];
    this.closeRange();
  }
  /** 字符串时间 转 秒 */
  dateToSeconds(date: string) {
    let st = new Date(date.replace(/-/g, "/")).getTime();
    return Math.floor(st / 1000);
  }
  /** 展示方式改变事件 */
  showTypeChange(val: String) {
    // 合并展示不能找过9张
    if (val === '2' && this.formData.sources.length > 9) {
      this.$toast('合并展示仅支持上传9张图片');
    }
  }
  /** 保存图片素材 */
  handleSubmit() {
    if (!this.validate) return;
    // 合并展示至少上传2张图片
    if (this.formData.show_type === '2' && this.formData.sources.length < 2) {
      this.$toast('合并展示最少上传2张');
      return false;
    }
    // 单张展示不能超过20张
    if (this.formData.show_type === '1' && this.formData.sources.length > 20) {
      this.$toast('单张展示仅支持上传20张图片');
      return false;
    }
    // 合并展示不能找过9张
    if (this.formData.show_type === '2' && this.formData.sources.length > 9) {
      this.$toast('合并展示仅支持上传9张图片');
      return false;
    }
    const start_time = this.dateToSeconds(this.formData.timeRange[0]);
    const end_time = this.dateToSeconds(this.formData.timeRange[1]);

    const params: ISaveMaterial = {
      source_type: 2, //1文本素材，2图片素材，5视频素材
      group_id:
        typeof this.formData.group.id === "number"
          ? this.formData.group.id
          : -1,
      category_ids: this.formData.tags.ids,
      start_time,
      end_time,
      sources: this.formData.sources,
      show_type: Number(this.formData.show_type)
    };
    console.log(this.formData, start_time, "提交的数据", params);
    this.saving = true;
    this.$toast.loading({
      duration: 0,
      message: "上传中...",
      forbidClick: true,
    });
    saveMaterial(params)
      .then((res: IApiResult) => {
        const { code, message = "", data } = res;
        this.saving = false;
        // 手动清除 Toast
        this.$toast.clear();
        if (code != 200) return this.$toast.fail(message);
        this.$toast("素材上传成功");
        //  神策数据上报 -- 上传素材
        console.log('保存图片素材', data);
        (data || []).forEach((item: String) => {
          this.$track.trackAssetUpload({
            asset_id: item
          })
        });
        setTimeout(() => {
          this.cancelAdd();
        }, 2000);
      })
      .catch((e) => {
        console.error(e);
        // 手动清除 Toast
        this.$toast.clear();
        this.saving = false;
      });
  }
  /**返回上一级 */
  cancelAdd() {
    this.$router.go(-1);
  }
  /** 上传成功的文件列表 */
  uploadChange(e) {
    const { fileList, finished } = e;
    console.log("uploadChange", e);
    this.finished = finished;
    this.formData.sources = fileList.map((file) => {
      return {
        url: file.url,
        is_origin: file.original ? 1 : 0,
        name: file.name,
        ext: file.type,
        size: file.size,
        // 如果图片是jpeg格式，is_transcode设置为1，否则设置为0
        is_transcode: file.type === "jpeg" ? 1 : 0,
        sort: file.sort,
      };
    });
  }
}
</script>
<style lang="scss" scoped>
.video-library {
  background: rgba(244, 247, 251, 1);
  padding: 15px 15px 70px;
  min-height: 100vh;

  .form {
    padding: 15px 15px 7px;
    background: #ffffff;

    .item {
      padding-bottom: 15px;
      border-bottom: 1px solid #eeeeee;
      margin-bottom: 15px;

      &:last-of-type {
        border-color: transparent;
        padding-bottom: 0px;
      }

      .show-type {
        display: flex;
        justify-content: space-around;
        align-items: center;
        width: 315px;
        height: 61px;
        background: rgba(216, 216, 216, 0.3);
        padding: 15px;
        margin-bottom: 15px;

        .left-label {
          display: flex;
          width: 30%;
          font-size: 14px;
          font-weight: 400;
          color: #000000;

          &-icon {
            position: relative;
            top: 5px;
            left: 5px;
            width: 12px;
            height: 12px;
          }
        }

        .radio-buttons {
          width: 70%;
          font-size: 12px;
          font-weight: 400;
          margin-left: 30px;
          color: #000000;

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
      }

      .label {
        font-size: 14px;
        font-family: PingFangSC-Regular, PingFang SC;
        font-weight: 500;
        color: #000000;
        line-height: 20px;
        margin-bottom: 15px;

        .require {
          color: #f64040;
        }
      }

      .content {
        display: flex;
        justify-content: space-between;
        font-size: 14px;
        font-family: PingFangSC-Medium, PingFang SC;
        color: #000000;
        line-height: 20px;
        font-weight: 700;

        i {
          display: block;
          width: 20px;
          height: 20px;
          background: url("../../assets/material-library/arrow-icon.png") no-repeat;
          background-size: 100% 100%;
        }

        &.empty {
          color: #999999;
          font-weight: 400;
        }
      }

      .content2 {
        font-size: 14px;
        font-family: PingFangSC-Regular, PingFang SC;
        color: #000000;

        i {
          margin: 0 10px;
          color: #000000;
          font-weight: 400;
        }

        span.empty {
          color: #999999;
          font-weight: 400;
        }
      }

      &:last-of-type {
        margin-bottom: 0;
      }
    }
  }

  .btn-group {
    z-index: 9;
    width: 100%;
    bottom: 0;
    left: 0;
    position: fixed;
    padding: 15px;
    background: #fff;
    display: flex;
    justify-content: space-between;
    font-size: 14px;
    font-weight: 500;
    font-family: PingFangSC-Regular, PingFang SC;

    .cancel {
      width: 165px;
      height: 40px;
      line-height: 40px;
      border: 1px solid #000000;
      color: #000000;
      text-align: center;
    }

    .submit {
      width: 165px;
      height: 40px;
      line-height: 40px;
      background: #000000;
      color: #ffffff;
      text-align: center;

      &.disabled {
        color: #ffffff;
        background: #cccccc;
      }

      .inline-block {
        display: inline-block;
      }
    }
  }
}

::v-deep {
  .video-name .van-field__control::-webkit-input-placeholder {
    /* WebKit browsers */
    color: #999;
    font-weight: normal;
  }

  .video-name.van-cell.van-field {
    background: #fff;
    padding: 0;
    font-size: 14px;

    .van-field__control {
      color: #000;
    }
  }
}
</style>
