<template>
  <div class="select-materials-dialog">
    <div
      v-if="allDialog"
      width="1000px"
    >
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
            @click="touchSearch()"
            size="small"
          >查询</el-button>
          <el-button
            @click="reset()"
            size="small"
          >重置</el-button>
          <el-button
            type="primary"
            @click="addTextmat()"
            class="addClass"
            size="small"
          >上传文件</el-button>
        </div>
        <div class="content">
          <div class="menu-left">
            <div class="sub-wrap">
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
            <div class="right-cont-top">
              <div
                :class="item.isChecked ? 'cont-item selected-item' : 'cont-item'
                  "
                v-for="(item, index) in tableList"
                :key="index"
                @click="selectedItem(item, index)"
              >
                <div
                  class="file-type picimg"
                  :class="`file-bg${item.file_type}`"
                >
                  {{
                    item.file_type != 5 ? customIcon[item.file_type - 1] : ""
                  }}
                  <i
                    v-if="item.file_type == 5"
                    class="icon el-icon-document-remove"
                    aria-hidden="true"
                  ></i>
                </div>
                <span v-if="item.name.length < 8">{{ item.name }}</span>
                <el-tooltip
                  v-else
                  class="item"
                  effect="light"
                  :content="item.name"
                  placement="top"
                >
                  <span>{{ item.name.slice(0, 7) }}...</span>
                </el-tooltip>
                <el-checkbox
                  class="item-select"
                  v-if="item.isChecked"
                  :value="true"
                ></el-checkbox>
              </div>
              <div
                slot="empty"
                v-if="!emptyFlag && tableList.length == 0"
                style="margin-top: 50px;text-align: center;"
              >
                <img
                  class="empty-img"
                  src="@/assets/images/empty-box.png"
                  alt=""
                />
                <div class="empty-box">
                  <div class="desc">暂无素材</div>
                </div>
              </div>
              <div
                slot="empty"
                v-if="emptyFlag && tableList.length == 0"
                style="margin-top: 50px;text-align: center;"
              >
                <img
                  class="empty-img"
                  src="@/assets/images/empty-search.png"
                  alt=""
                />
                <div class="empty-box">
                  <div class="desc">没有找到相关文件素材</div>
                </div>
              </div>
            </div>
            <div v-if="meta.total > 0">
              <el-pagination
                @size-change="handleSizeChange"
                @current-change="handleCurrentChange"
                :current-page.sync="meta.page"
                :page-sizes="[10, 20, 30, 40]"
                :page-size="10"
                layout="prev, pager, next"
                :total="meta.total"
              >
              </el-pagination>
            </div>
            <div class="fotter">
              <el-button
                @click="cancelDialog"
                size="small"
              >取 消</el-button>
              <el-button
                type="primary"
                @click="submitPar"
                size="small"
              >确 认</el-button>
            </div>
          </div>
        </div>
      </div>
    </div>
    <el-dialog
      title="选择图片素材"
      :visible.sync="selectDialogVisible"
      :before-close="cancelDialog"
      width="1000px"
      v-else
    >
      <div class="header">
        <div>
          <span
            @click="goback"
            :class="{ 'chosoe-clas': !addtexmat }"
          >选择文件素材</span>
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
            @click="touchSearch()"
            size="small"
          >查询</el-button>
          <el-button
            @click="reset()"
            size="small"
          >重置</el-button>
          <el-button
            type="primary"
            @click="addTextmat()"
            class="addClass"
            size="small"
          >上传文件</el-button>
        </div>
        <div class="content">
          <div class="menu-left">
            <div class="sub-wrap">
              <el-tree
                :data="groupData"
                node-key="id"
                default-expand-all
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
            <div class="right-cont-top">
              <div
                :class="item.isChecked ? 'cont-item selected-item' : 'cont-item'
                  "
                v-for="(item, index) in tableList"
                :key="index"
                @click="selectedItem(item, index)"
              >
                <div
                  class="file-type picimg"
                  :class="`file-bg${item.file_type}`"
                >
                  {{
                    item.file_type != 5 ? customIcon[item.file_type - 1] : ""
                  }}
                  <i
                    v-if="item.file_type == 5"
                    class="icon el-icon-document-remove"
                    aria-hidden="true"
                  ></i>
                </div>
                <span v-if="item.name.length < 8">{{ item.name }}</span>
                <el-tooltip
                  v-else
                  class="item"
                  effect="light"
                  :content="item.name"
                  placement="top"
                >
                  <span>{{ item.name.slice(0, 7) }}...</span>
                </el-tooltip>
                <el-checkbox
                  class="item-select"
                  v-if="item.isChecked"
                  :value="true"
                ></el-checkbox>
              </div>
              <div
                slot="empty"
                v-if="!emptyFlag && tableList.length == 0"
                style="margin-top: 50px;text-align: center;"
              >
                <img
                  class="empty-img"
                  src="@/assets/images/empty-box.png"
                  alt=""
                />
                <div class="empty-box">
                  <div class="desc">暂无素材</div>
                </div>
              </div>
              <div
                slot="empty"
                v-if="emptyFlag && tableList.length == 0"
                style="margin-top: 50px;text-align: center;"
              >
                <img
                  class="empty-img"
                  src="@/assets/images/empty-search.png"
                  alt=""
                />
                <div class="empty-box">
                  <div class="desc">没有找到相关文件素材</div>
                </div>
              </div>
            </div>
            <div v-if="meta.total > 0">
              <el-pagination
                @size-change="handleSizeChange"
                @current-change="handleCurrentChange"
                :current-page.sync="meta.page"
                :page-sizes="[10, 20, 30, 40]"
                :page-size="10"
                layout="prev, pager, next"
                :total="meta.total"
              >
              </el-pagination>
            </div>
            <div class="fotter">
              <el-button
                @click="cancelDialog"
                size="small"
              >取 消</el-button>
              <el-button
                type="primary"
                @click="submitPar"
                size="small"
              >确 认</el-button>
            </div>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Mixins, Component, Prop } from "vue-property-decorator";
import Api from "@/api/enterprise/material-library";
import { isNull } from "@/utils/tools";
import UploadImage from "@/components/business/upload-image/index.vue";
import UploadValidate from "@/mixins/UploadValidate";
import UploadDrag from "@/components/business/upload-drag/index.vue";
interface FormIn {
  sucai_id: number | undefined;
  name: string;
  urls: Array<any>;
}
@Component({
  components: {
    UploadImage,
    UploadDrag
  }
})
export default class extends Mixins(UploadValidate) {
  @Prop({ default: false, type: Boolean }) allDialog?: boolean; // 是否是放在tab组件下全部类型
  private rules = {};
  private tableList: any = []; // 表格数据
  private selectDialogVisible = true;
  private addtexmat = false;
  private emptyFlag = false;
  private ruleForm: FormIn = {
    sucai_id: 0,
    name: "",
    urls: []
  };
  private customIcon = ["W", "X", "P", "PDF"]; //自定义icon
  // 搜索区域
  private search = {
    keyWord: ""
  };
  // 是否全选
  private checked = false;
  private isActive = 0;
  // 列表移动分组数据
  private groupData: Array<any> = [];
  private originData: any[] = [
    { group_name: "全部分组", id: 0 },
    { group_name: "未分组", id: -1 }
  ];
  // 列表选中数据
  private selectedMaterialArray: Array<any> = [];
  //分页
  private meta = {
    size: 10,
    page: 1,
    total: 0
  };
  // 以下为---------------------------------------------------新增文本data--------------------------------------------------------

  private newGroupData: Array<any> = [];
  private newOriginData: any[] = [{ group_name: "未分组", id: -1 }]; //左侧分组
  private images: any = {
    url: "",
    name: "",
    size: "",
    ext: ""
  };
  private isUploading = false;
  //----------------------------------------------------methods-----------------------------------------------------------/
  private update(res: Array<string>) {
    this.ruleForm.urls = res;
    this.isUploading = false;
  }
  selectNode(node: any) {
    this.ruleForm.name = node.group_name;
    this.ruleForm.sucai_id = node.id;
    (this.$refs.selectBlur as any).blur();
  }
  // 点击×号关闭弹框
  private cancelDialog() {
    this.$emit("submit");
  }
  // 点击上传图片
  private addTextmat() {
    const newLink = `${origin}/web-admin/myapp/enterprise/material-library/material-list/list?isAddTask=1&type=6`;
    window.open(newLink, "_blank");
  }
  private goback() {
    if (this.isUploading) {
      this.$message.warning("图片正在上传中");
      return;
    }
    this.addtexmat = false;
  }
  // 外弹框点击确认
  private submitPar() {
    if (this.selectedMaterialArray.length == 0) {
      this.$message.warning("请选择图片素材");
      return;
    }
    this.$emit("submit", this.selectedMaterialArray);
  }
  // 列表数据
  getlist() {
    Api.consoleList({
      page: this.meta.page,
      size: this.meta.size,
      just_me: this.checked ? 1 : 0,
      // 分组ID  99 全部 -1 未分组
      source_group_id: this.ruleForm.sucai_id,
      keyword: this.search.keyWord,
      type: 6,
      is_use: 1
    }).then((res: any) => {
      if (res.code == 200) {
        const data = res.data.list || [];
        this.meta.total = res.data.total;
        if (data.length <= 0 || isNull(data)) {
          if (this.search.keyWord != "" || this.checked == true) {
            this.emptyFlag = true;
          } else {
            this.emptyFlag = false;
          }
          this.tableList = [];
          return;
        } else {
          // 选中者高亮
          const hasSelectArray = this.selectedMaterialArray.map(
            item => item.id
          );
          if (hasSelectArray && hasSelectArray.length > 0) {
            data.forEach((item: any) => {
              item.isChecked = hasSelectArray.indexOf(item.id) > -1 ? 1 : 0;
            });
          }
          this.tableList = data;
        }
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
    this.getlist();
  }
  // 点击重置
  private reset() {
    this.search.keyWord = "";
    this.checked = false;
    this.ruleForm.sucai_id = 0;
    this.isActive = 0;
    this.getlist();
  }
  // 点击每项图片
  private selectedItem(item: any, index: number) {
    // 仅能选择单张，则只高亮选中的
    this.tableList.map((value: any, i: number) => {
      this.$set(this.tableList[i], "isChecked", i == index ? 1 : 0);
    });
    this.selectedMaterialArray = [];
    this.selectedMaterialArray.push(item);
  }
  // 点击列表分组项
  private setMenuFn(data: any) {
    this.meta.page = 1;
    // this.ruleForm.sucai_id = id;
    // this.isActive = index;
    this.ruleForm.sucai_id = data.id;
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
.select-materials-dialog {
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

    .menu-left {
      margin-right: 10px;
      background: $color-background;
      width: 200px;
      height: 366px;
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

      .right-cont-top {
        height: 310px;
        clear: both;
        overflow-y: scroll;
      }

      .cont-item:nth-child(5n) {
        margin-right: 0px;
      }

      .cont-item {
        margin-right: 20px;
        margin-bottom: 15px;
        float: left;
        width: 128px;

        span {
          display: inline-block;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
          width: 100%;
        }

        .picimg {
          display: block;
          width: 128px;
          height: 128px;
          margin-bottom: 10px;
          border-radius: 4px;
          line-height: 128px;
          cursor: pointer;
          object-fit: cover;
          border: 1px solid transparent;
          box-sizing: border-box;
        }
      }
    }

    .selected-item {
      position: relative;

      .picimg {
        border: 1px solid #000000 !important;
      }

      .item-select {
        position: absolute;
        right: 0;
        top: -3px;
      }
    }
  }

  .fotter {
    position: absolute;
    right: 30px;
    bottom: 30px;
  }
}

.div-box {
  height: 481px;
}

.file-bg1 {
  background: #25589a;
}

.file-bg2 {
  background: #02843f;
}

.file-bg3 {
  background: #d34729;
}

.file-bg4 {
  background: #0e88e7;
}

.file-bg5 {
  background: #c9c9c9;
}

.file-type {
  width: 100%;
  height: 155px;
  line-height: 155px;
  border-radius: 4px;
  margin: 10px 0;
  color: #fff;
  text-align: center;
  font-size: 40px;

  .icon {
    font-size: 40px;
  }
}

/deep/.el-pagination {
  text-align: right;
  margin-top: 20px;
  margin-bottom: 20px;
}</style>
