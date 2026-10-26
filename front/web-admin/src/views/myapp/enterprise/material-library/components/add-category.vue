<template>
  <div>
    <el-dialog
      :title="dialogTitle[dialogType]"
      @close="closeVisible"
      :visible.sync="dialogVisible"
      width="400px"
    >
      <el-input
        type="textarea"
        placeholder="请输入标签名称"
        v-model="cate_name"
        maxlength="50"
        autosize
        show-word-limit
      >
      </el-input>
      <span slot="footer" class="dialog-footer">
        <el-button @click="closeVisible">取 消</el-button>
        <el-button type="primary" @click="submit">确 定</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import API from "@/api/enterprise/material-library";
@Component({})
export default class extends Vue {
  @Prop({ default: () => null }) selectedCate?: any;
  @Prop({ default: () => "1" }) dialogType?: string;
  /************data**************/
  private dialogTitle = {
    "1": "添加素材标签",
    "2": "编辑素材标签"
  };
  private cate_name = "";
  private dialogVisible = true;

  private closeVisible() {
    this.$emit("closeVisible");
  }
  async submit() {
    if (!this.cate_name) {
      this.$message.warning("请输入标签名称");
      return;
    }
    let res: any = null;
    if (this.dialogType == "1") {
      res = await API.categoryCreate({ name: this.cate_name });
    } else {
      res = await API.categoryUpdate({
        name: this.cate_name,
        id: this.selectedCate.id
      });
    }

    if (res.code == 200) {
      this.$message.success(this.dialogType == "1" ? "添加成功" : "修改成功");
      this.closeVisible();
      this.$emit("submit");
    } else {
      this.$message.error(res.message);
    }
  }
  private mounted() {
    if (this.dialogType == "2") {
      this.cate_name = this.selectedCate.name;
    }
  }
}
</script>
<style lang="scss" scoped>
.emj-boxs {
  width: 400px;
  height: 300px;
}
.emj-box {
  height: 250px;
  overflow: auto;
  font-size: 16px;
  user-select: none;
}
.emjs {
  display: inline-block;
  width: 28px;
  height: 28px;
  padding: 3px;
  text-align: center;
  cursor: pointer;
}
.text-head-button {
  font-size: 14px;
  font-family: PingFangSC-Regular, PingFang SC;
  font-weight: 400;
  color: rgba(96, 98, 102, 1);
  cursor: pointer;
}
.image {
  width: 14px;
  cursor: pointer;
  height: 14px;
}
</style>
