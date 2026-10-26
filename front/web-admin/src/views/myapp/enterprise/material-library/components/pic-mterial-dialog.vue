<template>
  <div class="select-materials-dialog">
    <el-dialog
      title="选择图片素材"
      :visible.sync="selectDialogVisible"
      :before-close="cancelDialog"
      width="1000px"
      class="dialog-class"
    >
      <div class="header">
        <div>
          <span
            @click="goback"
            :class="{ 'chosoe-clas': !addtexmat }"
          >选择图片素材</span><span
            v-if="addtexmat"
            :class="{ 'addText-clas': addtexmat }"
          >>上传图片</span>
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
            <el-input
              size="small"
              prefix-icon="el-icon-search"
              placeholder="请输入素材关键词"
              clearable
              style="width: 160px;"
              v-model="search.keyWord"
              @keyup.enter.native="touchSearch"
            ></el-input>
          </search-item>
          <search-item v-if="ep.isCoach" text="商品名称：">
            <el-input
              size="small"
              prefix-icon="el-icon-search"
              placeholder="请输入商品名称"
              clearable
              style="width: 160px;"
              v-model="search.product_name"
              @keyup.enter.native="touchSearch"
            ></el-input>
          </search-item>
          <search-item v-if="ep.isCoach" text="商品款号：">
            <el-input
              size="small"
              prefix-icon="el-icon-search"
              placeholder="请输入商品款号"
              clearable
              style="width: 160px;"
              v-model="search.sku_item"
              @keyup.enter.native="touchSearch"
            ></el-input>
          </search-item>
          <search-item v-if="ep.isCoach" text="商品颜色：">
            <el-input
              size="small"
              prefix-icon="el-icon-search"
              placeholder="请输入商品颜色"
              clearable
              style="width: 160px;"
              v-model="search.sku_color"
              @keyup.enter.native="touchSearch"
            ></el-input>
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
          >上传图片</el-button>
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
                
                <div class="picimg-box" v-loading="item.loading">
                  <el-image
                    class="picimg"
                    :src="item.url"
                    fit="cover"
                  />
                  <div v-if="!item.sku_list || item.sku_list.length > 0" class="tag-box">
                    <!-- 商品标签 -->
                    <div v-if="item.sku_list.length <=2" class="picimg-text">
                      <span>款号：{{ item.sku_list.join('、') }}</span>
                    </div>
                    <el-popover v-else placement="bottom" width="200px" trigger="hover">
                      
                      <div class="picimg-text" slot="reference">
                        <span>款号：{{ item.sku_list.join('、') }}</span>
                      </div>
                      <div class="tags w">
                        <el-tag v-for="(item, index) in item.sku_list" :key="index" type="gray" size="small">{{ item }}</el-tag>
                      </div>
                    </el-popover>
                  </div>
                </div>
                <el-tooltip
                  class="item"
                  effect="light"
                  :content="item.name"
                  placement="bottom"
                >
                  <span class="text">{{ item.name }}</span>
                </el-tooltip>
                <img
                  v-if="item.isChecked"
                  class="item-select"
                  src="@/assets/images/select-checked.png"
                  alt=""
                />
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
                  <div class="desc">暂无图片素材</div>
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
                  <div class="desc">没有找到相关图片素材</div>
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
          :model="ruleForm"
          :rules="rules"
          label-width="90px"
        >
          <search-item text="素材分组：">

            <el-cascader
              v-model="ruleForm.sucai_id"
              :options="newGroupData"
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
          </search-item>
          <cy-form-item
            label="本地图片"
            required
            prop="urls"
            trigger="change"
          >
            <upload-image
              class="upload-template"
              :urls="ruleForm.urls"
              accept="image/jpeg,image/png"
              img-type="batch"
              :before-upload="beforeUpload"
              @change="update"
              :limit="limiteNum"
              multiple
            >
            </upload-image>
            <span>支持JPG、PNG格式，图片小于2MB</span>
          </cy-form-item>
        </cy-form>
        <div class="fotter">
          <el-button
            type="primary"
            @click="childSure"
            size="small"
          >确 认</el-button>
        </div>
      </div>
      <div v-if="!addtexmat" slot="footer" class="fotter">
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
// todo:增加获取图片信息
import { getImageWh } from "@/utils/get-image-info";
import { getCosCompressUrl } from "@/utils/upload/utils";
interface FormIn {
  sucai_id: number | undefined;
  name: string;
  urls: Array<any>;
}
@Component({
  components: {
    UploadImage,
    UploadDrag,
  },
})
export default class extends Mixins(UploadValidate) {
  @Prop({ default: 9, type: Number }) limiteNum?: number; // 限制选择图片人数, 0为不限制
  // todo:是不是需要限制图片boolen
  @Prop({ default: false, type: Boolean }) isFriend?: boolean;
  @Prop({ default: '', type: Number }) groupSendType?: number;
  private rules = {};
  private tableList: any = []; // 表格数据
  private selectDialogVisible = true;
  private addtexmat = false;
  private emptyFlag = false;
  private ruleForm: FormIn = {
    sucai_id: 0,
    name: "",
    urls: [],
  };
  // 搜索区域
  private search = {
    keyWord: "",
    sku_item: "",
    product_name: "",
    sku_color: ""
  };
  // 是否全选
  private checked = false;
  private isActive = 0;
  // 列表移动分组数据
  private groupData: Array<any> = [];
  private originData: any[] = [
    { group_name: "全部分组", id: 0 },
    { group_name: "未分组", id: -1 },
  ];
  // 列表选中数据
  private selectedMaterialArray: Array<any> = [];
  //分页
  private meta = {
    size: 10,
    page: 1,
    total: 0,
  };
  // 以下为---------------------------------------------------新增文本data--------------------------------------------------------

  private newGroupData: Array<any> = [];
  private newOriginData: any[] = [{ group_name: "未分组", id: -1 }]; //左侧分组
  private images: any = {
    url: "",
    name: "",
    size: "",
    ext: "",
  };
  private imageArr = [];
  private isUploading = false;

  get ep() {
    return this.$store.getters['enterprise']
  }
  //----------------------------------------------------methods-----------------------------------------------------------/
  private update(res: any) {
    console.log(res);
    // this.ruleForm.urls = res;
    this.$set(this.ruleForm, "urls", res);
    this.ruleForm.urls = res;
    this.isUploading = false;
  }
  selectNode(node: any) {
    this.ruleForm.name = node.group_name;
    this.ruleForm.sucai_id = node.id;
    (this.$refs.selectBlur as any).blur();
  }
  // 点击号关闭弹框
  private cancelDialog() {
    this.$emit("submit");
  }
  // 点击上传图片
  private addTextmat() {
    if (this.ruleForm.sucai_id == 0 || this.ruleForm.sucai_id == -1) {
      this.ruleForm.sucai_id = undefined;
    }
    // this.addtexmat = true;
    this.ruleForm.urls = [];
    const newLink = `${origin}/web-admin/myapp/enterprise/material-library/create-image?from=2`;
    window.open(newLink, "_blank");
  }
  private goback() {
    // if (this.isUploading) {
    //   this.$message.warning("图片正在上传中");
    //   return;
    // }
    this.addtexmat = false;
  }
  // 外弹框点击确认
  private submitPar() {
    if (this.selectedMaterialArray.length == 0) {
      this.$message.warning("请选择图片素材");
      return;
    }
    const nowArr = this.selectedMaterialArray.some((item) => {
      return (
        item &&
        item.name &&
        item.name
          .substring(item.name.lastIndexOf(".") + 1, item.name.length)
          .toLocaleLowerCase() === "jpeg"
      );
    });
    if (nowArr) {
      this.$message.warning("不能选择带有jpeg格式的图片");
      return false;
    }
    this.$emit("submit", this.selectedMaterialArray);
  }
  // 内弹框点击确定
  private childSure() {
    this.ruleForm.urls.forEach((item: any, index: any) => {
      (this.imageArr[index] as any).url = item;
    });
    // console.log("this.isUploading---", this.isUploading);
    if (this.isUploading) {
      this.$message.warning("图片正在上传中");
      return;
    }
    (this.$refs["cform"] as any).validate((valid: any) => {
      if (valid) {
        this.imageAdd();
      }
    });
  }
  // 上传文件限制规范
  private beforeUpload(file: any) {
    const obj: any = {
      ext: "",
      name: "",
      size: "",
      url: "",
    };
    obj.name = file.name;
    obj.ext = file.type;
    obj.size = file.size;
    (this.imageArr as any).push(obj);
    this.isUploading = true;
    const name = file.name
      .substring(file.name.lastIndexOf(".") + 1, file.name.length)
      .toLocaleLowerCase();
    if (name !== "jpg" && name !== "png") {
      this.$message.error("请上传JPG、PNG格式的图片");
      return false;
    }
    if (file.size > 2 * 1024 * 1024) {
      this.$message.error("请上传2MB大小以内的图片");
      return false;
    }
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
      sku_item: this.search.sku_item,
      product_name: this.search.product_name,
      sku_color: this.search.sku_color,
      is_use: 1,
      type: 2,
    }).then((res: any) => {
      if (res.code == 200) {
        const data = res.data.list;
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
            (item) => item.id
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
  // 图片上传接口
  imageAdd() {
    const params = {
      group_id: this.ruleForm.sucai_id || -1,
      images: this.imageArr,
    };
    Api.imageAddApi(params).then((res: any) => {
      if (res.code == 200) {
        console.log("成功！");
        this.selectedMaterialArray = res.data;
        this.$emit("submit", this.selectedMaterialArray);
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
    this.search.sku_item = "";
    this.search.product_name = "";
    this.search.sku_color = "";
    this.checked = false;
    this.ruleForm.sucai_id = 0;
    this.isActive = 0;
    this.getlist();
  }
  // 点击新增
  private addBtn() {
    console.log("新增");
  }
  // 点击每项图片
  private async selectedItem(item: any, index: number) {
    // todo: groupSendType 为1时查询是否存在media_id
    this.$set(item, "loading", true);
    if (this.groupSendType == 1) {
      item.url = getCosCompressUrl({
        url: item.url,
        maxWidth: 1080,
        maxHeight: 1440,
      });
      
      const mediaRes: any = await Api.getMediaIdApi({
        source_id: item.id,
        source_type: 2,
        url: item.url,
        type: this.isFriend ? 1 : 0,
      });
      if (mediaRes.code != 200) {
        this.$message.error('该素材资源不符合要求，请重新选择');
        this.$set(item, "loading", false);
        return;
      }
    }
    // 仅能选择单张，则只高亮选中的
    if (this.limiteNum == 1) {
      this.tableList.map((value: any, i: number) => {
        this.$set(this.tableList[i], "isChecked", i == index ? 1 : 0);
      });
      this.selectedMaterialArray = [];
      this.selectedMaterialArray.push(item);
      return;
    }
    // 可选多张，选中者高亮
    if (!item.isChecked) {
      //todo: 判断是否超过9张
      if (this.selectedMaterialArray.length > 8) {
        this.$message.error(`最多只能选择9张图片`);
        return;
      }
      this.selectedMaterialArray.push(item);
      this.$set(this.tableList[index], "isChecked", 1);
    } else {
      this.selectedMaterialArray.splice(
        this.selectedMaterialArray.indexOf(item),
        1
      );
      this.$set(this.tableList[index], "isChecked", 0);
    }
    this.$set(item, "loading", false);
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

/deep/.el-tag.el-tag--gray {
  color: #666 !important;
  background-color: #F6F7F8 !important;
  border-color: #F6F7F8 !important;
}
/deep/.el-tag.el-tag--black {
  color: #fff !important;
  background-color: #000 !important;
  border-color: #000 !important;
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
      height: 420px;
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
          border-left: 4px solid #404040;
          color: #404040;
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
        height: 380px;
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

        .text {
          display: block;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
          width: 100%;
        }
        .picimg-box {
          position: relative;
          margin-bottom: 5px;
          border-radius: 4px;
          overflow: hidden;
        }

        .picimg {
          display: block;
          width: 128px;
          height: 128px;
          cursor: pointer;
        }

      }
    }

    .selected-item {
      position: relative;

      .picimg {
        border: 1px solid #404040;
        box-sizing: border-box;
      }

      .item-select {
        position: absolute;
        right: 0;
        top: 0;
        width: 14px;
        height: 14px;
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
  // height: 481px;
}

/deep/.el-pagination {
  text-align: right;
  margin-top: 20px;
  margin-bottom: 20px;
}
</style>
<style lang="scss"> 

.dialog-class .el-dialog {
  max-height: 800px !important;
}

.dialog-class .el-tree {
  padding: 0 !important;
}


.tag-box {
  position: absolute;
  bottom: 0px;
  left: 0;
  right: 0;
  height: 24px;
  background-color: rgba($color: #000000, $alpha: 0.5);
  line-height: 24px;

  .picimg-text {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    font-size: 10px;
    color: #fff;
    padding: 0 10px;
    text-align: center;

  }

}
.tags {
  margin-top: 5px;
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
  &.h {
    height: 32px;
  }

  &.w {
    max-width: 300px;
  }

  .more {
    color: #666;
    font-size: 12px;
  }
}
</style>
