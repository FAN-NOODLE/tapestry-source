<template>
  <div class="move-group">
    <el-dialog
      title="移动分组"
      :visible.sync="show_dialog"
      width="500px"
      @close="closeDialog"
    >
      <search-item text="素材分组：">
        <el-select style="width:280px" v-model="ruleForm.name" ref="selectBlur">
          <el-option
            :value="ruleForm.value"
            :label="ruleForm.name"
            style="height: auto;padding: 0px;"
          >
            <el-tree
              :data="groupData"
              node-key="id"
              default-expand-all
              :expand-on-click-node="false"
              @node-click="selectNode"
            >
              <span slot-scope="{ data }">
                <span>{{ data.group_name }}</span>
              </span>
            </el-tree>
          </el-option>
        </el-select>
      </search-item>
      <span slot="footer" class="dialog-footer">
        <el-button @click="closeDialog" size="small">取 消</el-button>
        <el-button type="primary" @click="submit" size="small">确 认</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import Api from "@/api/enterprise/material-library";
@Component({})
export default class extends Vue {
  @Prop() multipleSelection!: any; // 所选素材
  //------------------------------------------------------data-----------------------------------------------------------/
  private show_dialog = true;
  private rules = {};
  private ruleForm = {
    name: "",
    value: ""
    // group_id: ""
  };
  private groupData: any = []; //下拉列表
  //----------------------------------------------------methods-----------------------------------------------------------/
  selectNode(node: any) {
    this.ruleForm.name = node.group_name;
    this.ruleForm.value = node.id;
    (this.$refs.selectBlur as any).blur();
    console.log(node);
  }
  // 分组接口
  moveGroup() {
    Api.moveGroupApi({
      group_id: this.ruleForm.value,
      source_ids: this.multipleSelection
    }).then((res: any) => {
      if (res.code == 200) {
        this.$message.success("分组移动成功");
        this.$emit("closeDialog", true);
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 分组接口
  grouplist() {
    Api.getSourceGroupListApi({}).then((res: any) => {
      if (res.code == 200) {
        const data = res.data;
        this.groupData = data;
        // const obj1 = { group_name: "未分组", id: -1 };
        // (this.groupData as any).unshift(obj1);
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 点击完成
  private submit() {
    // (this.$refs["cform"] as any).validate((valid: any) => {
    if (this.ruleForm.name !== "") {
      this.moveGroup();
    } else {
      // this.$message.warning("请选择分组");
      return;
    }
    // });
  }

  // 关闭弹框
  private closeDialog() {
    this.$emit("closeDialog", false);
  }
  private mounted() {
    this.grouplist(); // 获取标签组下拉列表
  }
}
</script>
<style lang="scss" scoped>
/deep/.el-dialog {
  max-height: 325px;
}
/deep/.el-dialog__body {
  max-height: 156px;
}
</style>
