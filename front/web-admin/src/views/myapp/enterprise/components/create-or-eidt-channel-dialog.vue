<template>
  <div>
    <el-dialog
      :title="dialogTitle[dialogType]"
      :visible="true"
      @close="closeDialogFn"
      :close-on-click-modal="false"
    >
      <select-tag-dialog
        v-if="tagDialog"
        :updateSelectedTags="selectedTags"
        @selectedTags="selectedTagsFn"
      ></select-tag-dialog>
      <div class="section">
        <el-form
          ref="formObject"
          :rules="rules"
          :model="formObject"
          label-width="80px"
          label-position="right"
          size="medium"
        >
          <el-form-item label="渠道名称" prop="name" string required>
            <cy-input
              v-model.trim="formObject.name"
              maxlength="30"
              minlength="1"
              placeholder="请输入渠道名称"
              :show-word-limit="true"
              width-size="medium"
              clearable
            ></cy-input>
          </el-form-item>
          <el-form-item label="标签">
            <el-button
              type="primary"
              size="small"
              @click="openTagFn"
              v-if="showTagsList && showTagsList.length == 0"
            >
              选择标签
            </el-button>
            <p class="tip">通过该渠道添加的客户，将自动打上标签</p>
          </el-form-item>
          <el-form-item label="">
            <el-table
              ref="multipleTable"
              :data="tagsToTableData"
              v-if="!(showTagsList && showTagsList.length == 0)"
            >
              <el-table-column label="标签">
                <template slot-scope="scope">
                  <el-tag
                    v-for="(item, index) in scope.row.list"
                    class="tableTag"
                    :key="index"
                    >{{ item }}</el-tag
                  >
                </template>
              </el-table-column>
              <el-table-column label="操作" width="72" fixed="right">
                <template>
                  <el-button @click="openTagFn" type="text" size="small"
                    >修改</el-button
                  >
                </template>
              </el-table-column>
            </el-table>
          </el-form-item>
        </el-form>
      </div>
      <span slot="footer" class="dialog-footer">
        <el-button size="small" @click="closeDialogFn">取 消</el-button>
        <el-button size="small" type="primary" @click="submitFn"
          >保 存</el-button
        >
      </span>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import selectTagDialog from "./select-tag-dialog.vue";
import { buildChannelLabelValue } from "./channel-label-value";
import API from "@/api/enterprise/drainage-coach";
@Component({
  components: {
    selectTagDialog
  }
})
export default class extends Vue {
  @Prop({ default: () => null }) selectedChannel?: any;
  @Prop({ default: () => "1" }) dialogType?: string;

  public get tagsToTableData(): any {
    return [
      {
        list: this.showTagsList
      }
    ];
  }

  /************data**************/
  private dialogTitle = {
    "1": "新增渠道",
    "2": "编辑渠道"
  };
  /************methods**************/
  private formObject = {
    name: ""
  };
  private selectedTags: Array<any> = [];
  private showTagsList: Array<string> = [];
  private tagDialog = false;
  private validateMap = {
    name(rule: any, value: any, callback: Function) {
      if (!value) {
        callback(new Error("请输入渠道名称"));
      } else {
        callback();
      }
    }
  };
  private rules = {
    name: {
      validator: this.validateMap.name,
      trigger: ["blur", "change", "focus"]
    }
  };
  openTagFn() {
    this.tagDialog = true;
  }
  // 点击完成
  packageData(type?: string) {
    const data: any = {
      name: this.formObject.name,
      label: []
    };
    if (this.selectedChannel) {
      data.id = this.selectedChannel.id;
    }
    data.label = this.selectedTags.map((v: any, i: number) => {
      const labelArr: Array<any> = v.split(",");
      const req: any = {
        label_id: Number(labelArr[1]),
        label_group_id: Number(labelArr[0]),
        wx_tag_group_id: labelArr[2] || "",
        wx_tag_id: labelArr[3] || ""
      };
      if (type) {
        req.label_name = this.showTagsList[i];
      }
      return req;
    });
    return data;
  }
  submitFn() {
    (this.$refs["formObject"] as any).validate(async (valid: boolean) => {
      if (valid) {
        let res: any;
        const data: any = this.packageData();
        if (this.dialogType == "1") {
          res = await API.createChannelsApi(data);
          if (res.code !== 200) {
            this.$message.error(res.message);
            return false;
          }
          this.$message.success("保存成功!");
        } else {
          res = await API.channelsEditApi(data);
          if (res.code !== 200) {
            this.$message.error(res.message);
            return false;
          }
          this.$message.success("保存成功!");
        }
        this.$emit("channelDialog");
      }
    });
  }
  // 关闭弹框方法
  closeDialogFn() {
    this.$emit("channelDialog");
  }
  selectedTagsFn(list: any) {
    console.log(list);
    
    this.tagDialog = false;
    if (list) {
      this.selectedTags = list.ids;
      this.showTagsList = list.showList;
    }
  }
  setInitData(data: any) {
    data.forEach((e: any) => {
      this.showTagsList.push(e.label_name);
      this.selectedTags.push(buildChannelLabelValue(e));
    });
  }
  async initPageFn() {
    if (this.dialogType == "1") {
      let local: any = window.localStorage.getItem("CHANNEL_CREATE_LOCAL_DATA");
      local = local ? window.JSON.parse(local) : {};
      if (local.name) {
        this.formObject.name = local.name;
      }
      if (local.label && Object.keys(local.label).length > 0) {
        //showTagsList
        this.setInitData(local.label);
      }
    } else {
      const res: any = await API.channelsDetailApi({
        id: this.selectedChannel.id
      });
      if (res.code == 200) {
        this.formObject.name = res.data.name;
        this.setInitData(res.data.labels);
      } else {
        this.$message.error(res.message);
      }
    }
  }
  created() {
    this.initPageFn();
  }
}
</script>
<style lang="scss" scoped>
.section {
  .el-tag {
    background: #f7f7f7;
    color: #000000;
    border: 1px solid #dcdfe6;
  }
  height: 400px;
  .tip {
    color: #909399;
  }
}
.tableTag {
  margin: 0 10px 10px 0;
}
</style>
