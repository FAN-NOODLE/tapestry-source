<template>
  <div ref="content" class="copy-text">
    <p>各店同事：</p>
    <br/>
    <p>请见如下本周传讯目录</p>
    <br/>
    <p v-if="value.email_loop_url">邮件链接：{{ value.email_loop_url }}</p>
    <br/>
    <div v-for="(item, index) in value.content_list" :key="index">
      <p><u>{{item.department_name}}</u></p>
      <ul style="padding-left: 50px;">
        <li v-for="(info, infoIndex) in item.info_list" :key="infoIndex">{{info.title}}</li>
      </ul>
      <br/>
    </div>
    <p>中国区每周传讯-第{{value.week_num}}周.pdf</p>
    <p v-for="(item, index) in attachment" :key="index">Attach {{index+1}} - {{item}}</p>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from 'vue-property-decorator';

@Component
export default class CopyText extends Vue {
  @Prop({ type: Array, default: () => ({}) }) value: any

  get attachment() {
    let list = []
    if(!this.value.content_list) return list
    this.value.content_list.forEach((item: any) => {
      item.info_list.forEach((item: any) => {
        const attach = item.attachment || []
        list = list.concat(attach.map(val => val.attachment_name))
      })
    })
    return list
  }

  selectElementContents(el: HTMLElement) {
    if (document.createRange && window.getSelection) {
      const range = document.createRange();
      const selection = window.getSelection();
      selection?.removeAllRanges();
      try {
        range.selectNodeContents(el);
        selection?.addRange(range);
      } catch (e) {
        range.selectNode(el);
        selection?.addRange(range);
      }

      document.execCommand('copy');
      //取消文本选中状态
      window.getSelection()?.empty();
      this.$message({
        message: '已复制到剪贴板',
        type: 'success',
      });
    }
  }

  copy() {
    this.selectElementContents(this.$refs.content as HTMLElement);
  }
}
</script>

<style lang="scss" scoped>
.copy-text {
  position: fixed;
  right: 10000px;
}
</style>