<template>
  <van-popup v-model="show" class="popup" position="bottom" get-container="body" :overlay-style="{background: 'none'}">
    <div class="popup-content">
      <textarea
        v-if="show" 
        id="input"
        ref="textarea"
        class="textarea" 
        contenteditable
        v-model.trim="content" 
        placeholder="输入备注内容" 
        :maxlength="maxlength" 
        >
      </textarea>
      <div v-if="prefix" ref="prefix" class="prefix">回复<span class="prefix-name">{{ prefix }}</span>：</div>
      <div class="maxlength">{{ content.length }}/{{ maxlength }}</div>
    </div>
    <div class="popup-footer">
      <div class="btn plain" @click="handleCancel">取消</div>
      <div class="btn" :class="{disabled: content.length === 0}" @click="handleConfirm">保存</div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Vue, Prop } from 'vue-property-decorator'
import { Popup } from 'vant';
@Component({
  components: {
    [Popup.name]: Popup
  }
})
export default class CommentPopup extends Vue {
  @Prop({ default: 200 }) maxlength: number
  show: boolean = false
  content: string = ''
  prefix: string = ''

  showPopup(value?: string, prefix='') {
    this.show = true
    if (value) {
      this.content = value
    }

    if (prefix) {
      this.prefix = prefix
    }
    this.$nextTick(() => {
      const dom = this.$refs.prefix as HTMLElement
      const textarea = this.$refs.textarea as HTMLElement
      if(dom && textarea) {
        const width = dom.clientWidth;
        textarea.style.textIndent = `${width}px`
      }
      document.getElementById('input').focus();
      document.documentElement.scrollTo({
        top: document.body.scrollHeight,
        behavior: "auto"
      });
      // 让页面滚动到最底部
      // window.scrollTo(0, document.body.scrollHeight);
    })
  }

  hidePopup() {
    this.show = false
    this.content = ''
    this.prefix = ''
  }

  

  handleConfirm() {
    if(this.content.length === 0) {
      return
    }
    this.$emit('confirm', this.content, this.prefix)
    this.hidePopup()
  }

  handleCancel() {
    this.$emit('cancel')
    this.hidePopup()
  }
}
</script>
<style lang="scss" scoped>
.popup {
  padding: 15px;
  padding-bottom: calc(15px + constant(safe-area-inset-bottom)); /** iOS 11.0+ */
  padding-bottom: calc(15px + env(safe-area-inset-bottom)); /** iOS 11.0+ */

  &-content {
    width: 345px;
    height: 122px;
    background: #F6F5F5;
    padding: 15px;
    position: relative;
    .textarea {
      background-color: transparent;
      font-size: 12px;
      width: 100%;
      height: 100%;
      position: relative;
      user-select: auto;
      -webkit-user-select: auto;

      &::placeholder {
        color: #A6A6A6;
      }
    }

    .prefix {
      position: absolute;
      top: 15px;
      left: 15px;
      pointer-events: none; /* 防止伪元素阻碍用户与textarea交互 */
      font-size: 12px;
      color: #A6A6A6;

      &-name {
        color: #777777;
      }
    }

    .maxlength {
      position: absolute;
      right: 15px;
      bottom: 10px;
      font-size: 12px;
      color: #A6A6A6;
      // background: #F6F5F5;
    }
  
  }

  &-footer {
    margin-top: 15px;
    display: flex;
    align-items: center;
    .btn {
      flex: 1;
      background: #000000;
      color: #fff;
      font-size: 14px;
      font-weight: 400;
      line-height: 20px;
      height: 40px;
      display: flex;
      align-items: center;
      justify-content: center;

      &.plain {
        background: #fff;
        color: #000000;
        border: 1px solid #000000;
        margin-right: 15px;
      }

      &.disabled {
        background: #A6A6A6;
      }

    }
  }
}
</style>