<template>
  <div class="select-materials-dialog">
    <el-dialog
      :visible.sync="selectDialogVisible"
      width="1000px"
      @close="cancel"
    >
      <div class="header">
        <div>
          <span @click="goback" :class="{ 'chosoe-clas': defaultShow }"
            >选择链接素材</span
          ><span v-if="!defaultShow" :class="{ 'add-clas': !defaultShow }"
            >>新增链接</span
          >
        </div>
        <span @click="cancelDialog"
          ><i aria-hidden="true" class="el-icon-close"></i
        ></span>
      </div>
      <!-- 默认展示 -->
      <div v-if="defaultShow">
        <div class="search">
          <search-item text="搜索：">
            <el-input
              size="small"
              prefix-icon="el-icon-search"
              placeholder="请输入素材关键词"
              clearable
              style="width: 160px;"
              v-model="search.keyWord"
            ></el-input>
          </search-item>
          <search-item text="商品名称：">
            <el-input
              size="small"
              prefix-icon="el-icon-search"
              placeholder="请输入商品名称"
              clearable
              style="width: 160px;"
              v-model="search.product_name"
            ></el-input>
          </search-item>
          <search-item text="商品款号：">
            <el-input
              size="small"
              prefix-icon="el-icon-search"
              placeholder="请输入商品款号"
              clearable
              style="width: 160px;"
              v-model="search.sku_item"
            ></el-input>
          </search-item>
          <search-item text="商品颜色：">
            <el-input
              size="small"
              prefix-icon="el-icon-search"
              placeholder="请输入商品颜色"
              clearable
              style="width: 160px;"
              v-model="search.sku_color"
            ></el-input>
          </search-item>
          <search-item>
            <el-checkbox v-model="checked">仅显示自己上传的素材</el-checkbox>
          </search-item>
          <el-button type="primary" @click="touchSearch()" size="small"
            >查询</el-button
          >
          <el-button @click="reset()" size="small">重置</el-button>
          <el-button
            type="primary"
            @click="addBtn()"
            class="addClass"
            size="small"
            >新增链接</el-button
          >
        </div>
        <div class="content">
          <div class="menu-left">
            <div class="sub-wrap">
              <!-- <dd
                v-for="(item, index) in groupData"
                :key="index"
                @click="setMenuFn(index, item.id)"
                :class="{ 'check-color': index == isActive }"
              >
                <span v-if="item.group_name.length <= 9">{{
                  item.group_name
                }}</span>
                <el-popover
                  placement="top-start"
                  v-else
                  title=""
                  width="100%"
                  trigger="hover"
                  :content="item.group_name"
                >
                  <p slot="reference" class="text-span">
                    {{ item.group_name.slice(0, 9) }}...
                  </p>
                </el-popover>
              </dd> -->
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
          <div class="empty-content" v-if="linkList.length == 0">
            <img
              class="empty-img"
              v-if="!emptyFlag"
              src="@/assets/images/empty-box.png"
              alt=""
            />
            <div class="empty-box" v-if="!emptyFlag">
              <div class="desc">暂无链接素材</div>
            </div>
            <img
              class="empty-img"
              v-if="emptyFlag"
              src="@/assets/images/empty-search.png"
              alt=""
            />
            <div class="empty-box" v-if="emptyFlag">
              <div class="desc">没有找到相关链接素材</div>
            </div>
          </div>
          <div class="right-content" v-else>
            <div
              :class="
                isChecked === item.id
                  ? 'right-item selected-item'
                  : 'right-item'
              "
              v-for="(item, index) in linkList"
              :key="index"
              @click="selectedItem(item)"
            >
              <div class="item-title">
                {{ item.title }}
              </div>
              <div class="item-info">
                <div class="item-flex">
                  <div class="item-txt">
                    {{ item.description }}
                  </div>
                  <div v-if="!item.sku_list || item.sku_list.length === 0" class="tags h"></div>
                  <div v-else-if="item.sku_list.length <=2" class="tags h">
                    <el-tag v-for="(v, index) in item.sku_list" size="small" :key="index" type="gray">{{ v }}</el-tag>
                  </div>
                  <el-popover v-else placement="top" width="200px" trigger="hover">
                    <div class="tags h" slot="reference">
                      <el-tag v-for="(v, index) in item.sku_list.slice(0,2)" :key="index"  size="small" type="gray">{{ v }}</el-tag>
                      <span class="more">...</span>
                    </div>
                    <div class="tags w">
                      <el-tag v-for="(v, index) in item.sku_list" :key="index" type="gray">{{ v }}</el-tag>
                    </div>
                  </el-popover>
                </div>
                <img class="item-img" :src="item.cover_url" alt="" />
              </div>
              <img
                v-if="isChecked === item.id"
                class="item-select"
                src="@/assets/images/select-checked.png"
                alt=""
              />
            </div>
          </div>
          <div v-if="linkList.length > 0" class="fotter">
            <el-pagination
              @current-change="handleCurrentChangeFn"
              :current-page="meta.page"
              :page-size="meta.size"
              layout="prev, pager, next"
              :total="meta.total"
            >
            </el-pagination>
          </div>
        </div>
      </div>
      <span slot="footer" class="dialog-footer" v-if="defaultShow">
        <el-button @click="cancel" size="small">取 消</el-button>
        <el-button type="primary" @click="submitPar" size="small"
          >确 认</el-button
        >
      </span>
      <!-------------------新增链接模块 -------------------------->
      <div v-show="!defaultShow" class="add-link">
        <cy-form
          ref="cform"
          :rules="rules"
          :model="formData"
          label-width="90px"
        >
          <!-- <cy-form-item label="素材分组" prop="group_id">
            <cy-select
              v-model="formData.group_id"
              size="medium"
              placeholder="请选择素材分组"
              search-input
              :popperAppendToBody="true"
              v-if="!defaultShow"
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
              v-model="formData.name"
              ref="selectBlur"
            >
              <el-option
                :value="formData.group_id"
                :label="formData.name"
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
          <cy-form-item label="链接标题" prop="title" required>
            <cy-input
              ref="inp"
              placeholder="请输入链接标题"
              v-model.trim="formData.title"
              size="medium"
              maxlength="30"
              show-word-limit
            ></cy-input>
          </cy-form-item>
          <cy-form-item label="链接描述" prop="description" required>
            <cy-input
              ref="inp"
              placeholder="请输入链接描述"
              v-model.trim="formData.description"
              size="medium"
              maxlength="30"
              show-word-limit
            ></cy-input>
          </cy-form-item>
          <cy-form-item
            v-if="!defaultShow"
            label="封面图"
            required
            prop="cover_url"
            trigger="change"
          >
            <upload-image
              class="upload-template"
              :urls="formData.cover_url"
              accept="image/jpeg,image/png"
              @change="picSuccess"
              :before-upload="beforeUpload"
            ></upload-image>
            <div class="upload-tip">支持JPG、PNG格式，图片＜2MB</div>
          </cy-form-item>
          <cy-form-item label="链接内容" prop="link_type">
            <el-radio
              v-model="formData.link_type"
              :label="1"
              @change="radioChange"
              >URL
            </el-radio>
            <el-radio
              v-model="formData.link_type"
              :label="2"
              @change="radioChange"
              >自定义</el-radio
            >
          </cy-form-item>
          <cy-form-item
            label="URL"
            required
            v-if="formData.link_type == 1"
            prop="link_url"
          >
            <cy-input
              ref="inp"
              placeholder="请输入URL"
              v-model.trim="formData.link_url"
              size="medium"
            ></cy-input>
          </cy-form-item>
          <cy-form-item
            label="作者"
            required
            v-if="formData.link_type == 2"
            prop="author"
          >
            <cy-input
              ref="inp"
              placeholder="请输入作者名称"
              v-model.trim="formData.author"
              size="medium"
            ></cy-input>
          </cy-form-item>
          <cy-form-item
            label="正文"
            required
            v-if="formData.link_type == 2"
            prop="content"
          >
            <div id="editContent"></div>
          </cy-form-item>
        </cy-form>
      </div>
      <span slot="footer" class="dialog-footer" v-if="!defaultShow">
        <el-button type="primary" @click="addConfirm" size="small"
          >确 认</el-button
        >
      </span>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop} from "vue-property-decorator";
import UploadImage from "@/components/business/upload-image/index.vue";
import Api from "@/api/enterprise/material-library";
import { isNull } from "@/utils/tools";
// 引入 wangEditor富文本编辑器
import wangEditor from "wangeditor";
// todo:增加获取图片信息
import { getImageWh } from "@/utils/get-image-info";
import { getCosCompressUrl } from "@/utils/upload/utils";
@Component({
  components: {
    UploadImage,
  },
})
export default class extends Vue {
  // todo:是不是需要限制图片boolen
  @Prop({ default: false, type: Boolean }) isFriend?: boolean;
  @Prop({ default: '',type: Number}) groupSendType?: number;
  //------------------------------------------------------data-----------------------------------------------------------/
  private selectDialogVisible = true;
  // 搜索区域
  private search = {
    keyWord: "",
    sku_item: "",
    product_name: "",
    sku_color: "", // 商品颜色
  };
  // 是否全选
  private checked = false;
  private groupData: Array<any> = [];
  private originData: any[] = [
    { group_name: "全部分组", id: 0 },
    { group_name: "未分组", id: -1 },
  ];
  private newGroupData: Array<any> = [];
  private newOriginData: any[] = [{ group_name: "未分组", id: -1 }]; //左侧分组
  private linkList: any = [];
  private isActive = 0;
  private isChecked = 0;
  //分页
  private meta = {
    size: 9,
    page: 1,
    total: 10,
  };
  private emptyFlag = false;
  private defaultShow = true; //默认展示
  private selectedMaterial = {};
  // ---------新增链接data-----------------
  private group_id: any = 0;
  private formData = {
    title: "", //素材标题
    name: "",
    description: "", //素材描述
    group_id: "",
    link_type: 2, //(1 URL 2 自定义)
    cover_url: "", //封面图URL
    link_url: "", //链接URL
    author: "", //作者
    content: "", //正文内容
  };
  private rules = {
    link_url: [
      {
        validator: this.validateself,
        trigger: "blur",
      },
    ],
  };
  private validateself(rule: any, value: any, callback: any) {
    const match = /^(https?:\/\/)(([a-zA-Z0-9_-]+\.)+[a-zA-Z]{2,})(:\d+)?(\/[^\s?#]*)?(\?[^\s#]*)?(#[^\s]*)?$/;
    const flag = match.test(value);
    flag ? callback() : callback(new Error("请输入正确的URL"));
  }
  // 富文本编辑器
  private editorSms: any = null; //编辑器
  //----------------------------------------------------methods-----------------------------------------------------------/
  // 列表数据
  private getlist(type: number) {
    Api.consoleList({
      page: this.meta.page,
      size: this.meta.size,
      source_group_id: parseInt(this.group_id),
      just_me: this.checked ? 1 : 0, //仅显示自己  1是  0 否
      keyword: this.search.keyWord,
      sku_item: this.search.sku_item,
      sku_color: this.search.sku_color, // 商品颜色
      product_name: this.search.product_name,
      is_use: 1,
      type: 3,
    }).then((res: any) => {
      if (res.code == 200) {
        const data = res.data.list || [];
        const { size, page, total } = res.data;
        if (data.length <= 0 || isNull(data)) {
          if (type == 1) {
            this.emptyFlag = true;
            this.linkList = [];
          } else {
            this.emptyFlag = false;
            this.linkList = [];
          }
          return;
        }
        // data.map((item: any) => {
        //   item.checked = false;
        // });
        this.meta = { size, page, total };
        this.linkList = data;
      } else {
        this.$message.error(res.message);
      }
    });
  }
  selectNode(node: any) {
    this.formData.name = node.group_name;
    this.formData.group_id = node.id;
    this.group_id = node.id;
    (this.$refs.selectBlur as any).blur();
  }
  // 获取分组
  private grouplist() {
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
  // 选中素材
  private  async selectedItem(item: any) {
    console.log(item);
    item.cover_url = getCosCompressUrl({
      url: item.cover_url,
      maxWidth: 1080,
      maxHeight: 1440,
    });
    const { width, height } = await getImageWh(item.cover_url);
    // todo: groupSendType 为1时查询是否存在media_id
    if(this.groupSendType == 1) {
      const mediaRes: any = await Api.getMediaIdApi({
        source_id: item.id,
        source_type: 2,
        url:item.cover_url,
        type: this.isFriend ? 1 : 0,
      });
      if(mediaRes.code != 200){
        this.$message.error('该素材资源不符合要求，请重新选择');
        return;
      }
    }
    // todo: 判断是否超过限制
    if (this.isFriend) {
      // todo: 总像素不能超过1555200
      if (width * height > 1555200) {
        this.$message.error("图片像素不能超过1555200,请重新选择");
        return;
      }
    }
    this.isChecked = item.id;
    this.selectedMaterial = item;
  }
  // 取消
  private cancel() {
    this.$emit("submit");
    console.log("取消");
  }
  // 点击×号关闭弹框
  private cancelDialog() {
    this.cancel();
  }
  // 确认
  private submitPar() {
    if (Object.keys(this.selectedMaterial).length == 0) {
      this.$message.warning("请选择链接素材");
      return;
    }
    this.$emit("submit", this.selectedMaterial);
  }
  // 点击查询
  private touchSearch() {
    this.getlist(1);
  }
  // 点击重置
  private reset() {
    this.group_id = 0;
    this.isActive = 0;
    this.search.keyWord = "";
    this.search.sku_item = "";
    this.search.sku_color = ""; // 商品颜色
    this.search.product_name = "";
    this.checked = false;
    this.meta.page = 1;
    this.getlist(1);
  }
  // 点击新增
  private addBtn() {
    // this.radioChange();
    (this.$refs["cform"] as any).resetFields();
    // this.defaultShow = false;
    this.formData = {
      title: "", //素材标题
      name: this.formData.name,
      description: "", //素材描述
      group_id: this.group_id == 0 || this.group_id == -1 ? "" : this.group_id,
      // group_id: 0, //分组ID 全部=0  未分组=-1
      link_type: 1, //(1 URL 2 自定义)
      cover_url: "", //封面图URL
      link_url: "", //链接URL
      author: "", //作者
      content: "", //正文内容
    };

    const newLink = `${origin}/web-admin/myapp/enterprise/material-library/create-link?from=2`;
    window.open(newLink, "_blank");
  }
  //
  private picSuccess(res: any) {
    this.$set(this.formData, "cover_url", res);
    (this.$refs["cform"] as any).validateField("cover_url");
  }
  // 上传文件限制规范
  private beforeUpload(file: any) {
    if (file.type != "image/jpeg" && file.type != "image/png") {
      this.$message.error("请上传JPG、PNG格式的图片");
      return false;
    }
    if (file.size > 2 * 1024 * 1024) {
      this.$message.error("请上传2MB大小以内的图片");
      return false;
    }
  }
  private setMenuFn(data: any) {
    // this.group_id = id;
    this.group_id = data.id;
    // this.isActive = index;
    this.isChecked = 0;
    this.meta.page = 1;
    if (data.id != 0 && data.id != -1) {
      this.formData.group_id = data.id;
      this.formData.name = data.group_name;
    } else {
      this.formData.group_id = "";
      this.formData.name = "";
    }
    if (this.search.keyWord != "") {
      this.getlist(1);
    } else {
      this.getlist(0);
    }
  }
  // 页码选择回调
  private handleCurrentChangeFn(val: number) {
    this.meta.page = val;
    this.getlist(0);
  }
  private goback() {
    this.defaultShow = true;
    this.isChecked = 0;
  }
  private addConfirm() {
    if (this.editorSms) {
      this.$set(this.formData, "content", this.editorSms.txt.html());
    }
    const data = JSON.parse(JSON.stringify(this.formData));
    if (this.formData.link_type == 1) {
      data.author = "";
      data.content = "";
    } else {
      data.link_url = "";
      data.content = this.editorSms.txt.html();
    }
    data.group_id = parseInt(this.group_id) || -1;
    (this.$refs["cform"] as any).validate((valid: any) => {
      if (valid) {
        Api.linkSaveApi(data).then((res: any) => {
          if (res.code == 200) {
            this.$message.success("新增链接成功");
            this.$emit("submit", res.data);
          } else {
            this.$message.error(res.message);
          }
        });
      } else {
        return;
      }
    });
  }
  private radioChange() {
    if (this.formData.link_type == 2) {
      this.editorSms = new wangEditor(`#editContent`);
      this.editorSms.config.zIndex = 500;
      this.editorSms.create();
      if (this.formData.content) {
        this.editorSms.txt.html(this.formData.content);
      }
    }
    this.$nextTick(() => {
      (this.$refs["cform"] as any).clearValidate(["link_url"]);
      (this.$refs["cform"] as any).clearValidate(["content"]);
      (this.$refs["cform"] as any).clearValidate(["author"]);
    });
  }
  private mounted() {
    this.getlist(0);
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

/deep/.el-dialog {
  max-height: 100vh !important;
}

.tags {
  margin-top: 5px;
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
  &.h {
    height: 24px;
  }  
  &.w {
    width: 300px;
  }
  .more {
    color: #666;
    font-size: 12px;
  }
}
.el-tag { 
  display: block;
}

/deep/.el-input__inner {
  border-color:#DCDFE6 !important;
  &:focus {
    border-color: #000 !important;
  }
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
    height: 500px;
    .menu-left {
      background: $color-background;
      width: 200px;
      height: 100%;
      .sub-wrap {
        height: 100%;
        overflow-y: auto;
        dd {
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
        dd.check-color span {
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
      dd {
        margin-inline-start: 0px;
        height: 20px;
        line-height: 20px;
        // background-color: rgba(255, 255, 255, 1);
        color: $color-text-main;
        font-size: 14px;
        padding: 10px 0;
        cursor: pointer;
      }
      dd:hover {
        background: #ebebf1;
      }
    }
    .right-content {
      flex: 1;
      width: 731px;
      height: 289px;
      .right-item {
        width: 225px;
        min-height: 90px;
        position: relative;
        float: left;
        margin-left: 20px;
        padding: 10px;
        box-sizing: border-box;
        background: #ffffff;
        border: 1px solid #e4e7ed;
        border-radius: 4px;
        margin-bottom: 20px;
        &:hover {
          cursor: pointer;
        }
        .item-title {
          font-size: 14px;
          font-family: PingFangSC-Regular, PingFang SC;
          font-weight: 400;
          color: #303133;
          line-height: 20px;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
        .item-info {
          display: flex;
          margin-top: 10px;

          .item-flex {
            flex: 1;
            overflow: hidden;
          }
          .item-txt {
            width: 157px;
            height: 17px;
            font-size: 12px;
            font-family: PingFangSC-Regular, PingFang SC;
            font-weight: 400;
            color: #909399;
            line-height: 17px;
            overflow: hidden;
            text-overflow: ellipsis;
            display: -webkit-box;
            -webkit-line-clamp: 2; //行数需设置
            line-clamp: 2;
            -webkit-box-orient: vertical;
          }
          .item-img {
            margin-left: 10px;
            width: 40px;
            height: 40px;
            border-radius: 4px;
            img {
              width: 100%;
              height: 100%;
              border-radius: 4px;
            }
          }
        }
      }
      .selected-item {
        border: 1px solid #404040;
      }
      .item-select {
        position: absolute;
        right: 0;
        top: 0;
        width: 14px;
        height: 14px;
      }
    }
    .fotter {
      position: absolute;
      bottom: 81px;
      right: 20px;
      margin-top: 20px;
      text-align: right;
    }
    .empty-content {
      flex: 1;
      text-align: center;
      margin-top: 100px;
    }
  }
  /deep/.el-pagination {
    text-align: right;
  }
  // -------------------新增链接 ui-----------------
  .header {
    display: flex;
    justify-content: space-between;
    text-align: center;
    margin-bottom: 20px;
    font-size: 18px;
    font-family: PingFangSC-Semibold, PingFang SC;
    cursor: pointer;
    .chosoe-clas {
      color: #33363c;
      font-weight: 600;
    }
    .add-clas {
      color: #33363c;
      font-weight: 600;
    }
  }
  /deep/.el-dialog__header {
    display: none;
  }
  .add-link {
    max-height: 465px;
    overflow-y: auto;
  }
  .upload-tip {
    color: #909399;
    font-size: 12px;
    line-height: 18px;
    margin-top: 5px;
  }
  #editContent {
    width: 700px;
  }
}
</style>
