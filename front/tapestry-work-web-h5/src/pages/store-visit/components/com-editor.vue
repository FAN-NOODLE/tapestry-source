<template>
  <div class="">
    <quill-editor v-model="content"
      ref="myQuillEditor"
      :options="editorOption"
      @blur="onEditorBlur($event)"
      @focus="onEditorFocus($event)"
      @ready="onEditorReady($event)"
      @change="onEditorChange($event)">
    </quill-editor>
  </div>
</template>
<script lang="ts">
// require styles

import 'vue-quill-editor/node_modules/quill/dist/quill.core.css';
import 'vue-quill-editor/node_modules/quill/dist/quill.snow.css';
import 'vue-quill-editor/node_modules/quill/dist/quill.bubble.css';
import { quillEditor } from 'vue-quill-editor'
import { Component, Model, Prop, Vue, Watch } from "vue-property-decorator";

@Component({
  components: {
    quillEditor
  }
})
export default class ComEditor extends Vue {
  @Model('change', { type: String, default: '' }) value: string;
  content = ''
  editorOption = {
    // some quill options
    modules: {
      toolbar: [
        ['bold', 'italic', 'underline'],
        [{ 'list': 'ordered' }, { 'list': 'bullet' }],
      ]
    },
    placeholder: '输入内容...'
  }

  @Watch('value', { immediate: true, deep: true })
  onValueChanged(value) {
    this.content = value
  }

  onEditorBlur(quill) {
    console.log('editor blur!', quill)
    this.$emit('blur', quill)
  }
  onEditorFocus(quill) {
    console.log('editor focus!', quill)
    this.$emit('focus', quill)
  }
  onEditorReady(quill) {
    console.log('editor ready!', quill)
    quill.pasteHTML(this.content)
    this.$emit('ready', quill)
  }
  onEditorChange({ quill, html, text }) {
    console.log('editor change!', quill.hasFocus())
    this.$emit('change', html)
  }
}
</script>

<style lang="scss" scoped>
::v-deep .ql-toolbar.ql-snow {
  border: none;
  padding: 0;
  // height: 40px;
  display: flex;
  flex-wrap: wrap;
  border-bottom: 1px solid #ebeef5;
  padding-bottom: 8px;
}
::v-deep .ql-container.ql-snow {
  border: none;
  padding: 0;
  margin-top: 8px;
}
::v-deep .ql-editor {
  border: none;
  padding: 0;
}
</style>