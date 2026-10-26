<template>
  <div class="add-manual-group">
    <el-dialog
      title="新增标签组"
      :visible="true"
      width="500px"
      @close="closeDialog"
    >
      <cy-form
        class="main-box"
        ref="cform"
        :rules="rules"
        :model="ruleForm"
        label-width="73px"
      >
        <cy-form-item label="标签组" prop="groupName" required>
          <cy-input
            v-model="ruleForm.groupName"
            clearable
            maxlength="10"
            placeholder="请输入标签组名称"
            show-word-limit
            @keyup.native="checkSpace(1)"
          ></cy-input>
        </cy-form-item>
      </cy-form>
      <div>
        <div
          class="add-label-btn"
          :class="labelList.length >= 20 ? 'disable-add' : ''"
        >
          <div @click="addLabel">
            <i aria-hidden="true" class="el-icon-circle-plus"></i> 新建标签
          </div>
          <div class="add-label-tip">本次最多可新增20个标签</div>
        </div>
        <el-scrollbar
          wrap-class="scrollbar-wrapper"
          style="height:345px;width: 100% !important;"
        >
          <div
            class="item-label"
            v-for="(item, index) in labelList"
            :key="index"
          >
            <cy-input
              v-model="item.name"
              placeholder="请输入标签名称"
              clearable
              maxlength="10"
              show-word-limit
              @keyup.native="checkSpace(2, item, index)"
              @input="partRepeate2(item, index)"
            ></cy-input>
            <i
              aria-hidden="true"
              class="el-icon-close m-l-10"
              @click="deleteItem(index)"
            ></i>
            <div class="error-msg" v-if="item.error">
              {{ item.error }}
            </div>
          </div>
        </el-scrollbar>
      </div>
      <span slot="footer" class="dialog-footer">
        <el-button size="small" @click="closeDialog">取 消</el-button>
        <el-button size="small" type="primary" @click="submit">保 存</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import API from "@/api/customer";
@Component({})
export default class extends Vue {
  @Prop({ default: false }) centerDialogVisible?: boolean;
  private checkGroupName = (rule: any, value: any, callback: any) => {
    if (!value) {
      return callback(new Error("请输入标签组名称"));
    }
    callback();
  };
  /************data**************/
  private rules = {
    groupName: [{ validator: this.checkGroupName }]
  };
  private ruleForm = {
    groupName: ""
  };
  private labelList: any = [
    {
      name: ""
    },
    {
      name: ""
    },
    {
      name: ""
    },
    {
      name: ""
    },
    {
      name: ""
    }
  ];
  private btnDisable = 0; // 是否可以点击提交
  /************methods**************/
  // 禁止输入空格
  private checkSpace(type: number, item: any, index: number) {
    if (type == 1) {
      this.ruleForm.groupName = this.ruleForm.groupName.replace(/\s*/g, "");
    } else {
      this.$set(this.labelList[index], "name", item.name.replace(/\s*/g, ""));
    }
  }
  // 校验标签名称是否重复或未填
  private partRepeate2(itm: any, idx: number) {
    this.btnDisable = 0;
    itm.error = "";
    this.labelList.map((item: any, index: number) => {
      if (
        itm.name.trim() &&
        item.name.trim() == itm.name.trim() &&
        idx !== index
      ) {
        itm.error = "该标签名称已存在";
        this.btnDisable = 1;
      }
    });
  }
  // 点击完成按钮
  private submit() {
    (this.$refs["cform"] as any).validate((valid: any) => {
      if (valid) {
        const errorItem = this.labelList.find((x: any) => x.error);
        if (errorItem && errorItem.error == "该标签名称已存在") {
          this.btnDisable = 1;
        } else {
          this.btnDisable = 0;
        }
        if (this.btnDisable) {
          this.$message({
            message: "标签名不可重复",
            type: "warning",
            duration: 2000
          });
          return;
        }
        this.tagCreateGroup();
      } else {
        this.$message.warning("请输入标签组");
        return;
      }
    });
  }
  // 新增标签组接口
  tagCreateGroup() {
    const nameArray: any = [];
    this.labelList.map((item: any) => {
      if (item.name.trim()) {
        nameArray.push(item.name.trim());
      }
    });
    const params = {
      name: this.ruleForm.groupName.trim(),
      tag_names: nameArray
    };
    API.tagCreateGroupApi(params).then((res: any) => {
      console.log(res);
      if (res.code === 200) {
        console.log("新增成功");
        this.$emit("submit");
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 关闭弹框
  private closeDialog() {
    this.$emit("closeDialog");
  }
  // 点击新增标签按钮
  private addLabel() {
    const item = {
      name: ""
    };
    if (this.labelList.length < 20) {
      this.labelList.unshift(item);
    }
  }
  // 删除单列标签
  private deleteItem(index: number) {
    this.labelList.splice(index, 1);
  }
}
</script>
<style lang="scss" scoped>
.add-manual-group {
  .item-label {
    margin-top: 30px;
    .error-msg {
      padding-top: 0;
      margin-top: 5px;
      line-height: 18px;
      max-height: 36px;
      color: #f67170;
      font-size: 12px;
      text-overflow: -o-ellipsis-lastline;
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      line-clamp: 2;
      -webkit-box-orient: vertical;
    }
    ::v-deep.el-icon-close {
      color: #c0c4cc;
      cursor: pointer;
    }
    ::v-deep.el-icon-close:hover {
      color: #6f3112;
    }
  }
  .add-label-btn {
    cursor: pointer;
    display: flex;
    font-size: 14px;
    color: $color-main;
    ::v-deep.el-icon-circle-plus {
      color: $color-main;
      cursor: pointer;
    }
    .add-label-tip {
      color: $color-text-second;
      margin-left: 10px;
    }
  }
  .add-label-btn:hover {
    color: #6f3112;
    ::v-deep.el-icon-circle-plus {
      color: #6f3112;
      cursor: pointer;
    }
  }
  .add-label-btn:active {
    color: #6f3112;
    ::v-deep.el-icon-circle-plus {
      color: #6f3112;
      cursor: pointer;
    }
  }
  .disable-add,
  .disable-add:hover,
  .disable-add:active {
    cursor: not-allowed;
    color: #c0c4cc;
    ::v-deep.el-icon-circle-plus {
      color: $color-main;
      color: #c0c4cc;
    }
  }
}
</style>
