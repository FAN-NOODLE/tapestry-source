<template>
  <div class="select-materials-dialog">
    <el-dialog
      :visible.sync="selectDialogVisible"
      width="1000px"
      :before-close="cancelDialog"
      @click="cancelDialog"
    >
      <div class="header">
        <div>
          <span
            @click="goback"
            :class="{ 'chosoe-clas': !addtexmat }"
          >选择文本素材</span><span
            v-if="addtexmat"
            :class="{ 'addText-clas': addtexmat }"
          >>新增文本</span>
        </div>
        <span @click="cancelDialog"><i
            aria-hidden="true"
            class="el-icon-close"
          ></i></span>
      </div>
      <div
        class="div-box"
        v-if="!addtexmat"
      >
        <div class="search">
          <search-item>
            <cy-input
              size="small"
              prefix-icon="el-icon-search"
              placeholder="请输入素材关键词"
              clearable
              v-model="search.keyWord"
            ></cy-input>
          </search-item>
          <search-item>
            <el-checkbox v-model="checked" @change="touchSearch">仅显示自己上传的素材</el-checkbox>
          </search-item>
          <el-button
            type="primary"
            size="small"
            @click="touchSearch()"
          >查询</el-button>
          <el-button
            size="small"
            @click="reset()"
          >重置</el-button>
          <el-button
            size="small"
            type="primary"
            @click="addTextmat()"
            class="addClass"
          >新增文本</el-button>
        </div>
        <div class="content">
          <div class="menu-left">
            <div class="sub-wrap">
              <!-- <div
                v-for="(item, index) in groupData"
                :key="index"
                @click="setMenuFn(index, item.id)"
                :class="{ 'check-color': index == isActive, dd: true }"
              >
                <span v-if="item.group_name.length <= 7">{{
                  item.group_name
                }}</span>
                <el-tooltip
                  v-else
                  class="item"
                  effect="light"
                  :content="item.group_name"
                  placement="top"
                >
                  <span>{{ item.group_name.slice(0, 7) }}...</span>
                </el-tooltip>
              </div> -->
              <el-tree
                :data="groupData"
                node-key="id"
                :default-expand-all="false"
                :expand-on-click-node="false"
                @node-click="setMenuFn"
              >
                <span
                  style="width:210px;display:flex;justify-content: space-between;align-items: center;"
                  slot-scope="{ data }"
                >
                  <span v-if="data.group_name.length < 6">{{
                    data.group_name
                  }}</span>
                  <span v-if="data.group_name.length >= 6">
                    {{ data.group_name.slice(0, 5) }}...
                  </span>
                </span>
              </el-tree>
            </div>
          </div>
          <div class="right-content">
            <el-table
              :data="tableList"
              ref="multipleTable"
              stripe
              tooltip-effect="light"
            >
              <el-table-column width="38">
                <template slot-scope="scope">
                  <el-radio
                    :label="scope.row.id"
                    v-model="radioId"
                    @change="change(scope.row)"
                  ></el-radio>
                </template>
              </el-table-column>

              <el-table-column
                label="文本内容"
                min-width="650"
                prop="content"
                show-overflow-tooltip
              >
              </el-table-column>
              <div
                slot="empty"
                v-if="!emptyFlag && tableList.length == 0"
                style="margin-top: 50px;"
              >
                <img
                  class="empty-img"
                  src="@/assets/images/empty-box.png"
                  alt=""
                />
                <div class="empty-box">
                  <div class="desc">暂无文本素材</div>
                </div>
              </div>
              <div
                slot="empty"
                v-if="emptyFlag && tableList.length == 0"
                style="margin-top: 50px;"
              >
                <img
                  class="empty-img"
                  src="@/assets/images/empty-search.png"
                  alt=""
                />
                <div class="empty-box">
                  <div class="desc">没有找到相关文本素材</div>
                </div>
              </div>
            </el-table>
            <div v-if="meta.total > 0">
              <el-pagination
                @size-change="handleSizeChange"
                @current-change="handleCurrentChange"
                :current-page.sync="meta.page"
                :page-sizes="[10, 20, 30, 40]"
                :page-size="20"
                layout="prev, pager, next"
                :total="meta.total"
              >
              </el-pagination>
            </div>
            <div class="fotter">
              <el-button
                size="small"
                @click="cancelDialog"
              >取 消</el-button>
              <el-button
                size="small"
                type="primary"
                @click="submitPar"
              >确 认</el-button>
            </div>
          </div>
        </div>
      </div>
      <div
        class="div-box"
        v-if="addtexmat"
      >
        <cy-form
          class="main-box"
          ref="cform"
          :rules="rules"
          :model="ruleForm"
          label-width="90px"
        >
          <!-- <cy-form-item label="素材分组" prop="name">
            <cy-select
              v-model="ruleForm.id"
              size="medium"
              placeholder="请选择素材分组"
              search-input
              :popperAppendToBody="true"
            >
              <cy-option
                :label="item.group_name"
                :value="item.id"
                :key="index"
                v-for="(item, index) in newGroupData"
              ></cy-option>
            </cy-select>
          </cy-form-item> -->

          <search-item text="素材分组：">
            <el-select
              style="width:280px"
              v-model="ruleForm.name"
              ref="selectBlur"
            >
              <el-option
                :value="ruleForm.id"
                :label="ruleForm.name"
                style="height: auto;padding: 0px;"
              >
                <el-tree
                  :data="newGroupData"
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
          <cy-form-item
            label="文本内容"
            prop="text"
            required
          >
            <div class="publish-content-box">
              <!-- 表情、文本选择区域-->
              <div class="select-box-insert">
                <add-emoji
                  @select="addEmoji"
                  class="add-emoji"
                >
                  <template slot="reference">
                    <el-button
                      class="insert-btn"
                      size="small"
                    ><img
                        class="btn-img"
                        src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/face_icon.png"
                        alt=""
                      /><span>插入表情</span></el-button>
                  </template>
                </add-emoji>
              </div>
              <div>
                <el-input
                  ref="textarea"
                  class="textarea-box"
                  placeholder="请输入文本内容"
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
        <div class="fotter">
          <el-button
            size="small"
            type="primary"
            @click="childSure"
          >确 认</el-button>
        </div>
      </div>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
import addEmoji from "../../components/add-emoji.vue";
import Api from "@/api/enterprise/material-library";
import { isNull } from "@/utils/tools";
interface FormIn {
  id: number | undefined;
  name: string;
  text: string;
}
@Component({
  components: {
    addEmoji
  }
})
export default class extends Vue {
  //------------------------------------------------------data-----------------------------------------------------------/
  private tableList: any = []; // 表格数据
  private selectDialogVisible = true;
  private addtexmat = false;
  private emptyFlag = false;
  private radioId = null;
  private selection = {}; // 选中的数据
  // 搜索区域
  private search = {
    keyWord: ""
  };
  // 是否全选
  private checked = false;
  private groupData: Array<any> = [];
  private originData: any[] = [
    { group_name: "全部分组", id: 0 },
    { group_name: "未分组", id: -1 }
  ];
  private newGroupData: Array<any> = [];
  private newOriginData: any[] = [{ group_name: "未分组", id: -1 }]; //左侧分组
  private isActive = 0;
  //分页
  private meta = {
    size: 10,
    page: 1,
    total: 20
  };

  // 以下为---------------------------------------------------新增文本data--------------------------------------------------------
  private rules = {};
  private ruleForm: FormIn = {
    id: 0,
    name: "",
    text: ""
  };
  //----------------------------------------------------methods-----------------------------------------------------------/
  // 单选框绑定值变化时触发的事件（选中的 Radio label 值）
  change(row: any) {
    this.radioId = row.id;
    this.selection = row;
  }
  selectNode(node: any) {
    this.ruleForm.name = node.group_name;
    this.ruleForm.id = node.id;
    (this.$refs.selectBlur as any).blur();
  }
  // 点击×号关闭弹框
  private cancelDialog() {
    this.$emit("textData");
  }
  private addTextmat() {
    // if (this.ruleForm.id == 0 || this.ruleForm.id == -1) {
    //   this.ruleForm.id = undefined;
    // }
    // this.addtexmat = true;
    const newLink = `${origin}/web-admin/myapp/enterprise/material-library/material-list/list?isAddTask=1`;
    window.open(newLink, "_blank");
  }
  private goback() {
    this.addtexmat = false;
  }
  // 外弹框点击确认
  private submitPar() {
    if (Object.keys(this.selection).length == 0) {
      this.$message.warning("请选择文本素材");
      return;
    }
    console.log(this.selection);
    this.$emit("textData", this.selection);
  }
  // 内弹框点击确定
  private childSure() {
    (this.$refs["cform"] as any).validate((valid: any) => {
      if (valid) {
        this.textSave();
      }
    });
  }
  // 文本素材新增编辑
  textSave() {
    Api.textSaveApi({
      // ID 新增不传，编辑必传
      content: this.ruleForm.text,
      group_id: this.ruleForm.id || -1
    }).then((res: any) => {
      if (res.code == 200) {
        this.$message.success("新增文本成功");
        this.$emit("textData", res.data);
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 列表数据
  getlist() {
    Api.consoleList({
      page: this.meta.page,
      size: 20,
      keyWord: this.search.keyWord,
      source_group_id: this.ruleForm.id,
      // 是否仅展示自己上传的素材 1 是
      just_me: this.checked ? 1 : 0,
      type: 1,
      is_use: 1
    }).then((res: any) => {
      if (res.code == 200) {
        const data = res.data.list || [];
        if (data.length <= 0 || isNull(data)) {
          if (this.search.keyWord != "" || this.checked == true) {
            this.emptyFlag = true;
          } else {
            this.emptyFlag = false;
          }
        }
        this.meta.total = res.data.total;
        this.tableList = data;
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 分组接口
  grouplist() {
    Api.sourceGroupListApi({}).then((res: any) => {
      if (res.code == 200) {
        const data = res.data;
        this.groupData = [...this.originData, ...data];
        this.newGroupData = [...this.newOriginData, ...data];
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 点击查询
  private touchSearch() {
    this.meta.page = 1;
    this.getlist();
  }
  // 点击重置
  private reset() {
    this.search.keyWord = "";
    this.checked = false;
    this.ruleForm.id = 0;
    this.isActive = 0;
    this.meta.page = 1;
    this.getlist();
  }
  private setMenuFn(data: any) {
    this.meta.page = 1;
    // this.ruleForm.id = id;
    // this.isActive = index;
    this.ruleForm.id = data.id;
    this.ruleForm.name = data.group_name;
    this.getlist();
  }
  // 分页
  handleSizeChange(val: number) {
    this.meta.size = val;
    this.meta.page = 1;
    this.getlist();
  }
  handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getlist();
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
  mounted() {
    this.getlist();
    this.grouplist();
  }
}
</script>
<style lang="scss" scoped>
/deep/.el-input__inner {
  border-color:#DCDFE6 !important;
  &:focus {
    border-color: #000 !important;
  }
}
.text-span {
  display: inline-block;
  text-align: left;
  width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/deep/.el-pagination {
  text-align: right;
  margin-top: 20px;
}

::v-deep .el-dialog__body {
  height: 600px;
}

.header {
  display: flex;
  justify-content: space-between;
  text-align: center;
  margin-bottom: 20px;
  font-size: 16px;
  cursor: pointer;

  .chosoe-clas {
    color: #33363c;
    font-weight: 600;
  }

  .addText-clas {
    color: #33363c;
    font-weight: 600;
  }
}

/deep/.el-dialog__header {
  display: none;
}

.select-materials-dialog {
  .search {
    margin-bottom: 5px;
    position: relative;

    .addClass {
      position: absolute;
      right: 0;
    }
  }

  .content {
    display: flex;
    height: 420px;

    .menu-left {
      margin-right: 10px;
      background: $color-background;
      width: 200px;
      height: 320px;
      overflow-y: scroll;

      .sub-wrap {
        .dd {
          display: flex;
          align-items: center;
          justify-content: space-between;
          height: 20px;
          line-height: 20px;

          span {
            padding-left: 20px;
          }

          img {
            width: 3px;
            height: 17px;
          }
        }

        .dd.check-color span {
          border-left: 4px solid #708df5;
          color: #708df5;
        }

        img {
          display: inline-block;
          width: 20px;
          height: 20px;
          margin-right: 10px;
          border-radius: 50%;
        }
      }

      .dd {
        margin-inline-start: 0px;
        height: 20px;
        line-height: 20px;
        // background-color: rgba(255, 255, 255, 1);
        color: $color-text-main;
        font-size: 14px;
        padding: 10px 0;
        cursor: pointer;
      }

      .dd:hover {
        background: #ebebf1;
      }
    }

    .right-content {
      flex: 1;
      width: 731px;

      /deep/.el-table__body-wrapper {
        height: 285px;
        overflow-y: scroll;
      }
    }
  }

  .fotter {
    position: absolute;
    right: 30px;
    bottom: 30px;
  }
}

// 新增文本
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

.div-box {
  height: 481px;
}</style>
