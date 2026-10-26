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

    <!-- 视频素材表单 -->
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
        <div class="content2" @click="openRange">
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
        <div class="label">视频名称 <i class="require">*</i></div>
        <div class="content">
          <van-field
            class="video-name"
            v-model="formData.videoName"
            maxlength="50"
            placeholder="请输入视频名称，限制50个字符"
          />
        </div>
      </div>

      <div class="item">
        <div class="label">素材内容 <i class="require">*</i></div>
        <div class="content3">
          <single-video-upload @change="uploadChange" />
        </div>
      </div>
      <!-- 上传视频限制 desc -->
      <ul class="desc">
        <!-- <li>1. 视频格式仅支持*.mp4、*.m4v</li> -->
        <li>1. 视频不可超过10M</li>
        <li>2. 视频上传时会进行压缩，10M以内视频可选择原视频上传</li>
        <li>3. 建议视频时长控制在30秒内</li>
      </ul>
    </div>
    <div class="btn-group">
      <div class="cancel" @click="cancelAdd">取消</div>
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
import { Field, Loading, Popup } from "vant";
import materialGroup from "./components/material-group.vue";
import materialTag from "./components/material-tag.vue";
import materialTime from "./components/material-time.vue";
import singleVideoUpload from "./components/single-video-upload.vue";
import { saveMaterial } from "@/api/marketing-promotion/index";

@Component({
  components: {
    [Loading.name]: Loading,
    [Popup.name]: Popup,
    [Field.name]: Field,
    "material-group": materialGroup,
    "material-tag": materialTag,
    "material-time": materialTime,
    "single-video-upload": singleVideoUpload,
  },
  filters: {
    timeFormat: (value: string) => {
      return value.split(" ")[0];
    },
  },
})
export default class extends Vue {
  private saving = false;
  /**素材分组显示隐藏*/
  private showPopup = false;
  /** 添加素材标签 */
  private addTagPopup = false;
  /** 添加素材有效期 */
  private timeRangePopup = false;
  private formData = {
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
    /** 视频名称 */
    videoName: "",
    /** 视频 */
    sources: [],
  };
  private created() {
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
  private get validate() {
    if (
      this.formData.sources.length === 0 ||
      !this.formData.videoName.trim() ||
      this.formData.timeRange.length < 2
    ) {
      return false;
    }
    return true;
  }
  /**标签名称 */
  private get tagsName() {
    let tagName = "";
    this.formData.tags.names.forEach((name) => {
      tagName += `${name};`;
    });
    return tagName.slice(0, -1);
  }
  /**打开素材分组 */
  private openGroup() {
    this.showPopup = true;
  }
  /** 关闭选择素材分组 */
  private closeGroup() {
    this.showPopup = false;
  }
  /**选择素材分组 */
  private setGroup(e) {
    console.log("setGroup", e);
    // const regx = />/g;
    let { group, groupName } = e;
    this.formData.group.id = group;
    this.formData.group.name = groupName;
    this.closeGroup();
  }
  /**打开标签 */
  private openTag() {
    this.addTagPopup = true;
  }
  /**关闭标签选择 */
  private closeTag() {
    this.addTagPopup = false;
  }
  /**确认标签选择 */
  private setTag(e) {
    const { names, ids } = e;
    this.formData.tags.ids = ids;
    this.formData.tags.names = names;
    this.closeTag();
  }
  /**打开有效期 */
  private openRange() {
    this.timeRangePopup = true;
  }
  /**取消有效期 */
  private closeRange() {
    this.timeRangePopup = false;
  }
  /** 设置有效期 */
  private setRange(e) {
    const { start, end } = e;
    this.formData.timeRange = [start, end];
    this.closeRange();
  }
  /** 字符串时间 转 秒 */
  private dateToSeconds(date: string) {
    let st = new Date(date.replace(/-/g, "/")).getTime();
    return Math.floor(st / 1000);
  }
  /** 保存视频素材 */
  private handleSubmit() {
    if (!this.validate) return;
    const start_time = this.dateToSeconds(this.formData.timeRange[0]);
    const end_time = this.dateToSeconds(this.formData.timeRange[1]);
    // 输入框50个字的名称 拼接视频格式 提交
    this.formData.sources[0].name = this.formData.videoName.trim(); // name
    const params: ISaveMaterial = {
      source_type: 5, //1文本素材，2图片素材，5视频素材
      group_id:
        typeof this.formData.group.id === "number"
          ? this.formData.group.id
          : -1,
      category_ids: this.formData.tags.ids,
      start_time,
      end_time,
      sources: this.formData.sources,
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
        console.log('保存视频素材', data);
        (data||[]).forEach((item:String) => {
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
  private cancelAdd() {
    this.$router.go(-1);
  }
  /** 上传成功的文件列表 */
  private uploadChange(e) {
    console.log("uploadChange", e);
    this.formData.sources = e.map((file) => {
      return {
        url: file.url,
        is_origin: file.original ? 1 : 0,
        name: "",
        ext: file.type,
        size: file.size,
        is_transcode: file.type !== "mp4" ? 1 : 0,
        cover_url: file.poster,
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
          background: url("../../assets/material-library/arrow-icon.png")
            no-repeat;
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
      .content3 {
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
    }
  }
  .desc {
    padding-top: 10px;
    font-family: PingFangSC-Regular, PingFang SC;
    font-weight: 400;
    font-size: 10px;
    color: #9e9e9e;
    li {
      line-height: 20px;
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
