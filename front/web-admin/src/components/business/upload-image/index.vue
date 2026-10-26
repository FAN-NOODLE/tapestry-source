<template>
  <div class="upload-image">
    <component
      :is="imgType"
      :urls="urls"
      :accept="accept"
      @change="onChange"
      @update="onUpdate"
      @updateAI="onUpdateAI"
      @fail="onFail"
      :before-upload="beforeUpload"
      :on-success="onSuccess"
      :on-remove="onRemove"
      :multiple="multiple"
      :size="size"
      :limit="limit"
      :disabled="disabled"
    ></component>
  </div>
</template>
<script lang="ts">
/**
 * @author LeeYunxiang
 * @desc 不用card-list的组件，可以不用再绑数组，只适用于传一张图片的情景
 * @param {Function} onSuccess 成功回调
 * @param {Function} onError 失败回调
 * @param {Function} beforeUpload 上传前回调
 * @param {Function} onRemove 移除回调
 * @param {String} value 图片url
 */
import { Component, Vue, Prop } from "vue-property-decorator";
import { IFileParam, ImgType } from "./type";
import UploadImageBatch from "./upload-image-batch.vue";
import UploadImageSingle from "./upload-image-single.vue";

@Component({
  components: {
    Batch: UploadImageBatch,
    Single: UploadImageSingle
  }
})
export default class ComUploadImage extends Vue {
  /**
   * @param imgType
   * @description 单张上传为single 多张上传为batch
   * */
  @Prop({ default: "single" }) imgType?: ImgType;
  /**
   * @param urls
   * @description 绑定的图片url imgType = 'single'为string类型 为 batch为array类
   * */
  @Prop() urls!: string | Array<string>;
  /**
   * @param multiple
   * @description 是否开启批量上传，只有多张上传时生效
   * */
  @Prop({ default: false, type: Boolean }) multiple!: boolean;
  @Prop(Function) onSuccess!: Function;
  @Prop(Function) beforeUpload?: Function;
  @Prop(Function) onRemove?: Function;
  @Prop({ default: "" }) tips?: string;
  @Prop({ default: "" }) accept?: string;
  @Prop({ default: false, type: Boolean }) disabled!: boolean;
  @Prop({ default: "medium" }) size?: string;
  @Prop({ default: 0 }) limit?: number;
  private onChange(res: string | Array<string>) {
    this.$emit("change", res);
  }
  private onUpdate(res: Array<IFileParam>) {
    this.$emit("update", res);
  }
  private onFail() {
    this.$emit("fail");
  }
  private onUpdateAI(data: any) {
    this.$emit("updateAI", data);
  }
}
</script>
