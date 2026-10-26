<template>
  <van-popup v-model="visible" position="bottom" round @click-overlay="handleCancel">
    <div class="popup">
      <div class="popup-content">
        <div class="inner" @click="handleClick">
          <quill-editor v-model="innerValue"
            ref="myQuillEditor"
            :options="editorOption()"
            @blur="onEditorBlur($event)"
            @focus="onEditorFocus($event)"
            @ready="onEditorReady($event)"
            @change="onEditorChange($event)">
          </quill-editor>
        </div>
          
      </div>
      <div class="popup-footer">
        <div id="tool" class="toolbar">
          <button class="ql-bold"></button>
          <button class="ql-italic"></button>
          <button class="ql-underline"></button>
          <div class="line"></div>
          <button class="ql-list" value="ordered"></button>
          <button class="ql-list" value="bullet"></button>
          <!-- <div class="line"></div>
          <button class="ql-color ql-color-picker"></button> -->
        </div>
        <div class="btn" @click="handleConfirm">完成</div>
      </div>
    </div>
  </van-popup>
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
export default class ComEditorPopup extends Vue {
  @Model('update:show', { type: Boolean, default: false }) show: boolean;
  @Prop({ type: String, default: '' }) content: string;
  
  innerValue = ''
  
  get editorOption() {
    return () => {
      return {
        // some quill options
        modules: {
          toolbar: '#tool'
          // toolbar: [
          //   ['bold', 'italic', 'underline'],
          //   [{ 'list': 'ordered' }, { 'list': 'bullet' }],
          // ]
        },
        placeholder: this.$t('storeVisit.tips.36') as string,
      }
    }
  }


  visible: boolean = false;
  quill = null;
  
  @Watch('show', { immediate: true })
  onShowChange(val: boolean) {
    this.visible = val;
    if(val) {
      this.innerValue = this.content;
      this.$nextTick(() => {
        // 触发编辑器的ready事件
        if (this.quill) {
          console.log(this.quill);
          
          this.quill.focus();
          // 设置光标到编辑器末尾
          const length = this.quill.getLength();
          this.quill.setSelection(length, 0);
        }
      });
    }
  }

  @Watch('content', { immediate: true })
  onValueChanged(value) {
    this.innerValue = value
  }


  /**
   * 取消
   */
   handleCancel() {
    this.$emit('update:show', false);
    this.$emit('cancel');
  }

  /**
   * 确认
   */
  handleConfirm = this.$helper.debounce(() => {
    // 触发确认事件

    this.emitConfirm();
  }, 0, true)
  
  /**
   * 确认事件
   */
  emitConfirm() {
    this.$emit('update:show', false);
    this.$emit('confirm', this.innerValue);
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
    quill.pasteHTML(this.innerValue)
    this.quill = quill
    this.$emit('ready', quill)
  }
  onEditorChange({ quill, html, text }) {
    this.innerValue = html
  }

  handleClick() {
    // 触发编辑器的focus事件
    if (this.quill) {
      this.quill.focus();
      // 设置光标到编辑器末尾
      // const length = this.quill.getLength();
      // this.quill.setSelection(length, 0);
      // this.quill.scrollIntoView();
    }
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
  height: 100%;
}
::v-deep .ql-editor {
  border: none;
  padding: 0;
  height: 100%;
}
::v-deep .quill-editor {
  height: 100%;
}

::v-deep .ql-editor ol, 
::v-deep .ql-editor ul {
  padding-left: 0;
}
::v-deep .ql-bold,
::v-deep .ql-italic,
::v-deep .ql-underline,
::v-deep .ql-list {
  width: 20px !important;
  height: 20px !important;
  padding: 0 !important;
  display: block !important;
  position: relative;
  margin-right: 8px;

  svg {
    display: none;
  }

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    opacity: 1;
    transition: all 0.2s;
  }
  &::after {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    opacity: 0;
    transition: all 0.2s;
  }

  &.ql-active {
    &::before {
      opacity: 0;
    }
    &::after {
      opacity: 1;
    }
  }

}

::v-deep .ql-bold {
  &::before {
    background: url('~@/assets/store-visit/toolbar/1_0.png') center/100% no-repeat;
  }
  &::after {
    background: url('~@/assets/store-visit/toolbar/1_1.png') center/100% no-repeat;
  }
}
::v-deep .ql-italic {
  &::before {
    background: url('~@/assets/store-visit/toolbar/2_0.png') center/100% no-repeat;
  }
  &::after {
    background: url('~@/assets/store-visit/toolbar/2_1.png') center/100% no-repeat;
  }
}
::v-deep .ql-underline {
  &::before {
    background: url('~@/assets/store-visit/toolbar/3_0.png') center/100% no-repeat;
  }
  &::after {
    background: url('~@/assets/store-visit/toolbar/3_1.png') center/100% no-repeat;
  }
}
::v-deep .ql-list[value='ordered'] {
  &::before {
    background: url('~@/assets/store-visit/toolbar/4_0.png') center/100% no-repeat;
  }
  &::after {
    background: url('~@/assets/store-visit/toolbar/4_1.png') center/100% no-repeat;
  }
}
::v-deep .ql-list[value='bullet'] {
  &::before {
    background: url('~@/assets/store-visit/toolbar/5_0.png') center/100% no-repeat;
  }
  &::after {
    background: url('~@/assets/store-visit/toolbar/5_1.png') center/100% no-repeat;
  }
}

.line {
  width: 12px;
  height: 20px;
  background: url('~@/assets/store-visit/toolbar/line.png') center/4px no-repeat;
  margin-right: 8px;
}

.popup {
  .popup-content {
    padding: 16px 16px 8px;
    .inner {
      height: 104px;
      padding: 12px 16px ;
      background-color: rgba(233, 236, 239, 1);
      box-sizing: border-box;
      border-radius: 8px;
    }
  }

  .popup-footer { 
    display: flex;
    height: 44px;
    padding: 0 16px;


    .toolbar {
      flex: 1;
      display: flex;
      align-items: center;
    }

    .btn {
      width: 32px;
      height: 100%;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 12px;
      font-weight: 600;
      color: rgba(13, 110, 253, 1);

      &.disabled {
        color: rgba(153, 153, 153, 1);
      }
    }
  }
}
</style>