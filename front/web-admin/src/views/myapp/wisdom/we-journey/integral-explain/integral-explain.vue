<template>
  <layout-form>
    <template slot="content">
      <div class="title">积分规则说明</div>
      <div id="editContent"></div>
      <div class="footer-rank">
        <el-button type="primary" @click="save">保存</el-button>
        <el-button @click="goBack">取消</el-button>
      </div>
    </template>
  </layout-form>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { ruleDetailAPi, ruleSaveAPi } from "@/api/we-journey/Integral";
import { ApiResult } from "@/types";
// 引入 wangEditor
import wangEditor from "wangeditor";
@Component({
  components: {
  }
})
export default class IntegralExplain extends Vue {
  /** 详情 */
  private value = "";
  // 富文本编辑器
  private editorSms: any = null;
   //编辑器菜单
  private editorMenus = [
    "head",
    "bold",
    "fontSize",
    "fontName",
    "italic",
    "underline",
    "strikeThrough",
    "indent",
    "lineHeight",
    "foreColor",
    "backColor",
    // "link",
    "list",
    "justify",
    "quote",
    "emoticon",
    // "image",
    // "video",
    "table",
    // "code",
    "splitLine",
    "undo",
    "redo"
  ];
  /** 保存 */
  private async save() {
    if (this.value === "<p><br></p>")
      return this.$message.error("内容不能为空");

    const res: ApiResult = await ruleSaveAPi({
      content: this.value
    });

    if (res.code !== 200) return this.$message.error(res.message as string);

    this.$message.success("保存成功");
  }

  /** 获取详情 */
  private async getDetail() {
    const res: ApiResult = await ruleDetailAPi();

    if (res.code !== 200) return this.$message.error(res.message as string);

    this.value = res.data.content;
    this.editorSms.txt.html(this.value);
  }

  /** 返回 */
  private goBack() {
    this.$router.go(-1);
  }

  /** 初始化 */
  private created() {
    this.getDetail();
  }
  private mounted(){
      this.editorSms = new wangEditor(`#editContent`);
      this.editorSms.config.zIndex = 500;
      this.editorSms.config.menus = this.editorMenus;
      this.editorSms.create();
  }
}
</script>
<style lang="scss" scoped>
.title {
  display: flex;
  justify-content: center;
  margin-bottom: 15px;
  font-size: 14px;
  font-weight: 600;
}

.footer-rank {
  margin-top: 20px;
}
</style>
