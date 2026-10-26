<template>
  <div :id="editorId" class="rich-text-editor">
  </div>
</template>
<script lang="ts">
import { cosUploader } from "@/utils/upload";
import { forEach } from "jszip";
import { Component, Prop, Vue, Watch } from "vue-property-decorator";
import wangEditor from "wangeditor";

@Component({
  model: {
    prop: "value",
    event: "change"
  }
})
export default class RichTextEditor extends Vue {
  /** 富文本内容 */
  @Prop({ type: String, default: "" }) value!: string
  @Prop({ type: String, default: "editContent" }) editorId!: string
  /** 富文本编辑器 */
  private editor: any = null;

  mounted() {
    this.$nextTick(() => {
      this.init()
    })
  }

  /**
   * 初始化
   */
  init() {
    this.editor = new wangEditor(`#${this.editorId}`);
    // 属性配置
    this.editor.config.zIndex = '99'
    this.editor.config.placeholder = '请输入传讯内容'
    this.editor.config.uploadImgShowBase64 = false;
    this.editor.config.showLinkImg = false;
    this.editor.config.autofocus = false;
    this.editor.config.indentation = '1em'
    // 自定义上传
    this.editor.config.customUploadImg = function (files: File[], insertFn: any) {  // TS 语法
      files.forEach((file: File) => {
        cosUploader.uploadFile(file).then((data) => {
          insertFn(data.url)
        }).catch(() => {
          this.$message.error('上传失败')
        })
      })
    }

    // 工具栏配置
    this.editor.config.menus = [
      "bold",
      // 缩进
      "underline",
      "indent",
      "link",
      "list",
      "image",
      "table",
      "undo",
      "redo"
    ]

    // 监听内容变化
    this.editor.config.onchange = () => {
      this.value = this.editor.txt.html();
      this.$emit('change', this.value);
    };
    // 创建编辑器
    this.editor.create();
    // 自定义高度
    const list = document.getElementsByClassName('w-e-text-container') as any
    list.forEach((item:HTMLElement) => {
      item.style.height='auto'
      item.style.minHeight='300px'
    })
    
    setTimeout(() => {
      this.editor.txt.html(this.value)
    },1000)
    console.log(this.editor.config);
    
  }
}
</script>

<style lang="scss" scoped>
.rich-text-editor {
  width: 80%;
  position: relative;
}
::v-deep {
  .w-e-toolbar {
    z-index: 1000 !important;
  }
}
</style>