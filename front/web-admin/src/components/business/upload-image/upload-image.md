# 上传图片组件

> 说明：上传图片的专用组件，可单张上传，可批量多张上传并拖拽排序

## 组件参数说明

| 参数 | 类型  | 说明 | 可选值  | 默认值  |    是否必填 |
|----------------|------- |-------------------- |----------------- |------------- | ------------------ |
| img-type | String | 单张上传or多张上传 | single(单张)/batch(多张) | single |  否    | 
| urls | String/Array | 绑定的图片url,单张上传时为字符串，多张上传时为数组 | — | - | 是  |
| multiple | Boolean | 是否开启批量上传，只在多张上传时生效 | — | false | 否  |
| on-success | Function | 上传成功的回调函数 | | - |  否    | 
| before-upload | function(file) | 上传之前的钩子函数，参数为上传的文件，若返回 false 或者返回 Promise 且被 reject，则停止上传 | — | - | 否  |
| on-remove | function() | 文件列表移除文件时的钩子 | — | - | 否  |
| disabled | Boolean | 是否禁用 | — | false | 否  |
| size | string | 大小型号 | medium(宽高100px)/small(宽高50px) | medium | 否  |
| tips | string | 位于上传器的文案提示 | - | '上传图片' | 否  |

## 组件事件说明

| 事件 | 说明  | 回调参数 | 
|-----|------- |---------- |
| @change |  绑定urls值发生变化会执行  | res: urls | 
| @fail |  上传失败执行  | res: urls |


#### 效果图

<f-image src="https://img-crs.vchangyi.com/组件说明样式资源/管理后台/ts/upload-image.png" width="500px" />

#### 代码参考
 ``` html
 <template>
  <div class="upload-demo">
    <div class="upload">
      <h3 class="title">上传图片组件 单张demo</h3>
      <upload-image
        class="upload-template"
        :urls="form_data.avator"
        @change="picSuccess"
      ></upload-image>
    </div>
    <div class="upload">
      <h3 class="title">上传图片组件 多张可拖拽排序demo</h3>
      <upload-image
        class="upload-template"
        img-type="batch"
        :urls="form_data.urls"
        @change="update"
        multiple
      ></upload-image>
    </div>
    <div class="upload">
      <h3 class="title">size为small的上传图片</h3>
      <upload-image
        class="upload-template"
        size="small"
        :urls="form_data.avator"
        @change="picSuccess"
      ></upload-image>
    </div>
    <div class="upload">
      <h3 class="title">size为small的上传多张图片</h3>
      <upload-image
        class="upload-template"
        img-type="batch"
        size="small"
        :urls="form_data.urls"
        @change="update"
        multiple
      ></upload-image>
    </div>
  </div>
</template>
<script lang="ts">
interface Normal {
  [propName: string]: any;
}
import { Component, Mixins } from "vue-property-decorator";
import UploadImage from "@/components/business/upload-image/index.vue";
import UploadValidate from "@/mixins/UploadValidate";
import { Loading } from "element-ui";
@Component({
  components: {
    UploadImage
  }
})
export default class extends Mixins(UploadValidate) {
  // data 变量
  private images = [];
  private remove(i: number) {
    console.log(i, "删除的索引");
  }
  private form_data: Normal = {
    avator: "",
    order_num: "",
    selecter: "",
    urls: []
  };
  private picSuccess(res: any) {
    this.$set(this.form_data, "avator", res);
  }
  private update(res: Array<string>) {
    this.$set(this.form_data, "urls", res);
    Loading.service({ fullscreen: true }).close();
  }
  private picRemove() {
    console.log("我删除了");
  }
}
</script>
<style lang="scss" scoped>
// 公共颜色

.upload-demo {
  margin: 30px;
  .upload {
    margin-bottom: 30px;
    .title {
      display: inline-block;
      vertical-align: middle;
      margin-right: 50px;
    }
    .upload-template {
      display: inline-block;
      vertical-align: middle;
      margin: 0;
    }
  }
}
</style>

 ```