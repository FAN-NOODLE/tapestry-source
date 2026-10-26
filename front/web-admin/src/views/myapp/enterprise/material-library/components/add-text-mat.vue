<template>
  <div class="add-edit-group">
    <el-dialog
      :title="title"
      :visible.sync="show_dialog"
      width="700px"
      @close="closeDialog"
    >
      <div class="add-box">
        <cy-form
          class="main-box"
          ref="cform"
          :rules="rules"
          :model="ruleForm"
          label-width="90px"
        >
          <cy-form-item label="素材分组">
            <el-cascader
              v-model="ruleForm.group_id"
              :options="groupData"
              :props="{
                label: 'group_name',
                value: 'id',
                checkStrictly: true,
                emitPath: false
              }"
              :style="{
                width: '400px'
              }"
            ></el-cascader>
          </cy-form-item>
          <cy-form-item label="素材标签" prop="category_ids" trigger="change">
            <el-select
              v-model="ruleForm.category_ids"
              multiple
              placeholder="请选择素材标签"
              style="width:400px"
            >
              <el-option
                v-for="item in category_options"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              >
              </el-option>
            </el-select>
            <el-button
              style="font-size:20px;transform: translateY(3px);color: #000000;"
              class="m-l-20"
              type="text"
              icon="el-icon-circle-plus"
              @click="openAddCate"
            ></el-button>
          </cy-form-item>
          <cy-form-item
            label="有效期"
            prop="datePicker"
            required
            trigger="change"
          >
            <el-date-picker
              v-model="ruleForm.datePicker"
              :default-time="['00:00:00', '23:59:59']"
              type="datetimerange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              clearable
              value-format="timestamp"
              @change="datePickerChange"
            >
            </el-date-picker>
          </cy-form-item>
          <cy-form-item label="发布内容" prop="text" required>
            <div class="publish-content-box">
              <!-- 表情、文本选择区域-->
              <div class="select-box-insert">
                <add-emoji @select="addEmoji" class="add-emoji">
                  <template slot="reference">
                    <el-button class="insert-btn" size="small"
                      ><img
                        class="btn-img"
                        src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/face_icon.png"
                        alt=""
                      /><span>插入表情</span></el-button
                    >
                  </template>
                </add-emoji>
              </div>
              <div>
                <el-input
                  ref="textarea"
                  class="textarea-box"
                  placeholder="请输入消息内容"
                  v-model="ruleForm.text"
                  type="textarea"
                  maxlength="500"
                  show-word-limit
                >
                </el-input>
              </div>
            </div>
          </cy-form-item>
        </cy-form>
      </div>
      <span slot="footer" class="dialog-footer">
        <el-button @click="closeDialog" size="small">取 消</el-button>
        <el-button type="primary" @click="submit" size="small">保 存</el-button>
      </span>
    </el-dialog>
    <add-category
      v-if="addVisible"
      @closeVisible="closeVisible"
      @submit="addCateSubmit"
    ></add-category>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import addEmoji from "../../components/add-emoji.vue";
import addCategory from "./add-category.vue";
import Api from "@/api/enterprise/material-library";
import Times from "./datePicker";
interface T {
  group_id: number;
  content: string;
  id: number;
}
@Component({
  components: {
    addEmoji,
    addCategory
  }
})
export default class extends Vue {
  @Prop() title!: string;
  @Prop() groupName!: string;
  @Prop() groupId!: number;
  @Prop() editTextDate?: T;
  //------------------------------------------------------data-----------------------------------------------------------/
  private show_dialog = true;
  private groupData: any[] = [];
  private originData: any[] = [{ group_name: "未分组", id: -1 }]; //下拉列表
  private rules = {};
  private ruleForm = {
    group_id: 0,
    name: "",
    text: "",
    category_ids: [], //素材标签
    datePicker: Times.times,
    start_time: Times.start, //有效期开始时间
    end_time: Times.end //有效期结束时间
  };
  // 素材标签
  private category_options: any[] = [];
  private addVisible = false; // 控制添加素材标签弹窗

  //----------------------------------------------------methods-----------------------------------------------------------/
  // 获取素材详情
  private sourceDetail() {
    Api.sourceDetail({
      source_id: (this.editTextDate as T).id,
      type: 1
    }).then((res: any) => {
      if (res.code == 200) {
        const data = res.data[0];
        let categoryIds = [];
        if (data.category_ids.trim()) {
          categoryIds = data.category_ids
            .split(",")
            .map(item => parseInt(item));
        }
        if (data.start_time) {
          this.ruleForm.start_time = data.start_time;
          this.ruleForm.end_time = data.end_time;
          this.ruleForm.datePicker = [
            data.start_time * 1000,
            data.end_time * 1000
          ];
        }

        this.ruleForm.category_ids = categoryIds;
        this.ruleForm.group_id = data.group_id || -1;
        if (!data.group_id) {
          // this.ruleForm.name = "全部分组";
          this.ruleForm.name = "未分组";
        } else if (data.group_id === -1) {
          this.ruleForm.name = "未分组";
        } else {
          this.ruleForm.name = data.group_name;
        }
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 添加表情
  async addEmoji(emj: any) {
    const myField = (this.$refs.textarea as any).$refs.textarea;
    let strTem: string = emj || "";
    const surplus: number = 500 - this.ruleForm.text.length;
    if (surplus - strTem.length < 0) {
      strTem = strTem.substring(0, surplus);
      this.$message.warning("消息内容长度不能超过500字符，超出部分则被截断。");
    }
    if (myField.selectionStart || myField.selectionStart === 0) {
      const startPos = myField.selectionStart;
      const endPos = myField.selectionEnd;
      this.ruleForm.text =
        myField.value.substring(0, startPos) +
        strTem +
        myField.value.substring(endPos, myField.value.length);
      await this.$nextTick();
      myField.focus();
      myField.setSelectionRange(endPos + strTem.length, endPos + strTem.length);
    } else {
      this.ruleForm.text += strTem;
    }
  }
  selectNode(node: any) {
    this.ruleForm.name = node.group_name;
    this.ruleForm.group_id = node.id;
    (this.$refs.selectBlur as any).blur();
  }
  // 点击完成
  private submit() {
    (this.$refs["cform"] as any).validate((valid: any) => {
      if (valid) {
        this.textSave();
      }
    });
  }
  // 保存
  private textSave() {
    if (this.title == "编辑文本素材") {
      this.sourceUpdate();
    } else {
      this.addtextSave();
    }
  }
  // 编辑接口
  private sourceUpdate() {
    Api.sourceUpdate({
      source_type: 1,
      id: (this.editTextDate as T).id,
      content: this.ruleForm.text,
      group_id: this.ruleForm.group_id || 0,
      start_time: this.ruleForm.start_time,
      end_time: this.ruleForm.end_time,
      category_ids: this.ruleForm.category_ids
    }).then((res: any) => {
      if (res.code == 200) {
        this.$message.success("编辑文本成功");
        this.closeDialog();
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 新增文本接口
  private addtextSave() {
    Api.textSaveApi({
      content: this.ruleForm.text,
      group_id: this.ruleForm.group_id || 0,
      start_time: this.ruleForm.start_time,
      end_time: this.ruleForm.end_time,
      category_ids: this.ruleForm.category_ids
    }).then((res: any) => {
      if (res.code == 200) {
        this.$message.success("新增文本成功");
        this.closeDialog();
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 获取分组名
  inquireName(row: any) {
    console.log(6666666);
    Api.inquireGroupName({ group_id: row.group_id }).then((res: any) => {
      if (res.code == 200) {
        this.ruleForm.name = res.data.group_name;
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 素材分组接口
  grouplist() {
    Api.sourceGroupListApi({}).then((res: any) => {
      if (res.code == 200) {
        const data = res.data;
        this.groupData = [...this.originData, ...data];
        if (this.title == "编辑文本素材") {
          this.ruleForm.group_id = (this.editTextDate as T).group_id;
          // this.inquireName(this.ruleForm);
          this.ruleForm.text = (this.editTextDate as T).content || "";
        }
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 关闭弹框
  private closeDialog() {
    this.$emit("closeDialog");
  }
  // 打开添加素材标签弹窗
  private openAddCate() {
    this.addVisible = true;
  }
  // 关闭取消弹窗
  private closeVisible() {
    this.addVisible = false;
  }
  // 新增素材标签成功
  private addCateSubmit() {
    this.getCategory();
  }

  // 获取数据
  async getCategory() {
    const param: object = {
      page: 1,
      size: 999
    };
    const res: any = await Api.categoryList(param);
    if (res.code == 200) {
      this.category_options = res.data.list || [];
    } else {
      this.$message.error(res.message);
    }
  }
  // 有效期选择
  private datePickerChange(val: any) {
    if (val) {
      // 因为 parseInt 的参数类型为 string ,所以处理后的时间戳只能先转成 string 然后再取整
      this.ruleForm.start_time = parseInt((val[0] / 1000).toString());
      this.ruleForm.end_time = parseInt((val[1] / 1000).toString());
    } else {
      this.ruleForm.start_time = 0;
      this.ruleForm.end_time = 0;
    }
  }
  private mounted() {
    this.ruleForm.group_id = this.groupId || -1;
    this.ruleForm.name = this.groupId ? this.groupName : "未分组";
    this.grouplist();
    this.getCategory();
    if (this.title === "编辑文本素材") {
      this.sourceDetail();
    }
  }
}
</script>
<style lang="scss" scoped>
.add-box {
  max-height: 460px;
  overflow: auto;
}
.publish-content-box {
  width: 100%;
  border: solid 1px #dcdfe6;
  min-height: 200px;
  .select-box-insert {
    height: 60px;
    margin: 0 15px;
    display: flex;
    align-items: center;
    .add-emoji {
      display: inline-flex;
      justify-content: center;
      margin-right: 10px;
      /deep/span {
        display: inline-flex;
        align-items: center;
      }
    }
    .insert-btn {
      border: none;
      border-radius: 14px !important;
      display: inline-flex;
      align-items: center;
      justify-content: center;
      padding: 0;
      font-size: 12px;
      width: 87px;
      height: 28px !important;
      background-color: $color-background;
      /deep/span {
        display: inline-flex;
        align-items: center;
      }
      .btn-img {
        width: 14px;
        height: 14px;
        margin-right: 5px;
      }
    }
  }
  .textarea-box {
    width: calc(100% - 30px) !important;
    margin-left: 15px;
    border-top: solid 1px #ebeef5;
    /deep/textarea {
      border: none;
      resize: none;
      height: 255px !important;
      padding: 15px 0;
      font-family: PingFangSC-Medium, PingFang SC;
    }
    /deep/.el-input__count {
      right: 0;
    }
  }
}
</style>
