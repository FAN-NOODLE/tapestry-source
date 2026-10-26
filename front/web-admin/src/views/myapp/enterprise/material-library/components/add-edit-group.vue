<template>
  <div class="add-edit-group">
    <el-dialog
      :title="title"
      :visible.sync="show_dialog"
      width="500px"
      @close="closeDialog"
    >
      <cy-form
        class="main-box"
        ref="cform"
        :rules="rules"
        :model="ruleForm"
        label-width="90px"
        @submit.native.prevent
      >
        <cy-form-item label="分组名称" prop="name" required>
          <cy-input
            v-model="ruleForm.name"
            maxlength="15"
            show-word-limit
          ></cy-input>
        </cy-form-item>
      </cy-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="closeDialog" size="small">取 消</el-button>
        <el-button type="primary" @click="submit" size="small">保 存</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import Api from "@/api/enterprise/material-library";
@Component({})
export default class extends Vue {
  @Prop() title!: string;
  @Prop() rowData?: object;
  @Prop() filterObj?: any; // 筛选条件
  //------------------------------------------------------data-----------------------------------------------------------/
  private show_dialog = true;
  private groupData = [];
  private rules = {};
  private ruleForm = {
    name: ""
  };
  private groupId = 0;
  private total = 0;
  //----------------------------------------------------methods-----------------------------------------------------------/

  // 添加分组接口
  togroupAdd() {
    this.groupAdd();
    // if (this.total > 0 && (this.rowData as any).children === null) {
    //   this.$confirm(
    //     "当前分组下已有素材内容，添加分组后，所有素材内容将自动迁移到新分组下。",
    //     "分组提示",
    //     {
    //       confirmButtonText: "确定",
    //       cancelButtonText: "取消",
    //       type: "warning"
    //     }
    //   )
    //     .then(() => {
    //       this.groupAdd();
    //     })
    //     .catch(() => {
    //       this.$message({
    //         type: "info",
    //         message: "已取消添加"
    //       });
    //     });
    // } else {
    //   this.groupAdd();
    // }
  }
  groupAdd() {
    let parent_id = 0;
    if (this.title == "添加一级分组") {
      parent_id = 0;
    } else {
      parent_id = (this.rowData as any).id;
    }
    Api.sourceGroupAddApi({
      group_name: this.ruleForm.name,
      parent_id
    }).then((res: any) => {
      if (res.code == 200) {
        this.$message.success("添加成功");
        this.closeDialog();
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 编辑分组接口
  groupEdit() {
    Api.sourceGroupEditApi({
      id: (this.rowData as any).id,
      group_name: this.ruleForm.name
    }).then((res: any) => {
      if (res.code == 200) {
        this.$message.success("修改成功");
        this.closeDialog();
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 点击完成
  private submit() {
    (this.$refs["cform"] as any).validate((valid: any) => {
      if (valid) {
        if (this.title == "添加分组" || this.title == "添加一级分组") {
          this.togroupAdd();
        } else {
          this.groupEdit();
        }
      }
    });
  }
  // 关闭弹框
  private closeDialog() {
    this.$emit("closeDialog");
  }
  // 列表数据
  getlist() {
    const params = {
      page: 1, //页码
      size: 10, //每页数量
      type: 1
    };
    // 对象合并
    Object.assign(params, this.filterObj);
    Api.consoleList(params).then((res: any) => {
      if (res.code == 200) {
        // const data = res.data.list;
        // this.tableList = data;
        this.total = res.data.total;
      } else {
        this.$message.error(res.message);
      }
    });
  }
  private mounted() {
    // if (Object.keys(this.rowData as any).length > 0) {
    //   this.ruleForm.name = (this.rowData as any).group_name;
    // }
    console.log(this.rowData);
    if (this.title == "添加分组") {
      this.getlist();
    }
    if (this.title == "添加分组" || this.title == "添加一级分组") {
      this.ruleForm.name = "";
      return;
    }
    if (this.title == "编辑分组") {
      this.ruleForm.name = (this.rowData as any).group_name;
      return;
    }
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
