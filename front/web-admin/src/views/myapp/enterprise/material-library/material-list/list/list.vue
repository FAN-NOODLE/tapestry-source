<template>
  <div class="supernatant">
    <layout-list>
      <template slot="search">
        <div class="material-search-filters">
          <search-item text="素材搜索：">
            <cy-input
              size="small"
              prefix-icon="el-icon-search"
              placeholder="请输入素材名称"
              clearable
              v-model="search.keyWord"
              @keyup.enter.native="touchSearch"
            ></cy-input>
          </search-item>
          <search-item text="创建时间：">
            <el-date-picker
              v-model="search.datePicker"
              :default-time="['00:00:00', '23:59:59']"
              type="datetimerange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              clearable
              @change="touchSearch"
            >
            </el-date-picker>
          </search-item>
          <search-item text="素材标签：">
            <el-select
              v-model="search.category_ids"
              multiple
              :popper-append-to-body="false"
              placeholder="请选择"
              @change="touchSearch"
            >
              <el-option
                v-for="item in category_options"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              >
              </el-option>
            </el-select>
          </search-item>
          <search-item text="状态：">
            <el-select
              v-model="search.source_status"
              clearable
              placeholder="请选择"
              @change="touchSearch"
            >
              <el-option
                v-for="(item, index) in status_options"
                :key="index"
                :label="item.label"
                :value="item.value"
              >
              </el-option>
            </el-select>
          </search-item>
          <search-item v-if="ep.isCoach" text="商品名称：">
            <cy-input
              size="small"
              prefix-icon="el-icon-search"
              placeholder="请输入商品名称"
              clearable
              v-model="search.product_name"
              @keyup.enter.native="touchSearch"
            
            ></cy-input>
          </search-item>
          <search-item v-if="ep.isCoach" text="商品款号：">
            <cy-input
              size="small"
              prefix-icon="el-icon-search"
              placeholder="请输入商品款号"
              clearable
              v-model="search.sku_item"
              @keyup.enter.native="touchSearch"
            ></cy-input>
          </search-item>
          <search-item v-if="ep.isCoach" text="商品颜色：">
            <cy-input
              size="small"
              prefix-icon="el-icon-search"
              placeholder="请输入商品颜色"
              clearable
              v-model="search.sku_color"
              @keyup.enter.native="touchSearch"
            ></cy-input>
          </search-item>
          <search-item v-if="ep.isCoach" text="关联商品状态：">
          
            <el-select
              v-model="search.goods_relation"
              placeholder="请选择"
              @change="touchSearch"
            >
              <el-option
                v-for="(item, index) in relation_status_options"
                :key="index"
                :label="item.label"
                :value="item.value"
              >
              </el-option>
            </el-select>
          </search-item>
          <search-item
            v-if="ep.isCoach && activTab == '3'"
            text="素材来源："
          >
            <el-select
              v-model="search.source_origin"
              placeholder="请选择"
              @change="touchSearch"
            >
              <el-option
                v-for="(item, index) in source_origin_options"
                :key="index"
                :label="item.label"
                :value="item.value"
              >
              </el-option>
            </el-select>
          </search-item>
          <search-item text="素材搜索：">
            <el-checkbox v-model="search.checked" @change="touchSearch"
              >仅显示自己上传的素材</el-checkbox
            >
          </search-item>
        </div>
        <div v-if="showUploaderFilters" class="material-search-filters">
          <search-item text="上传人名称：">
            <cy-input
              size="small"
              placeholder="请输入上传人名称"
              clearable
              v-model="search.uploader_name"
              @keyup.enter.native="touchSearch"
            ></cy-input>
          </search-item>
          <search-item text="所属门店：">
            <cy-input
              size="small"
              placeholder="请输入门店名称"
              clearable
              v-model="search.department_name"
              @keyup.enter.native="touchSearch"
            ></cy-input>
          </search-item>
        </div>
        <div>
          <el-button type="primary" size="small" @click="touchSearch()"
            >查询</el-button
          >
          <el-button size="small" @click="reset()">重置</el-button>
        </div>
      </template>
      <template slot="column-left">
        <!-- 左侧分组 -->
        <add-edit-group
          v-if="pushOrEditGroup"
          :title="addtit"
          :filterObj="filterObj"
          :rowData="rowData"
          @closeDialog="closeaddgroup"
        ></add-edit-group>
        <dl class="menu-left">
          <dt @click="addgroup">
            <div>
              <i aria-hidden="true" class="el-icon-circle-plus-outline"></i>
              添加一级分组
            </div>
          </dt>
          <div class="sub-wrap">
            <el-tree
              :data="groupData"
              node-key="id"
              :default-expand-all="false"
              :expand-on-click-node="false"
              @node-click="setMenuFn"
              :allow-drop="allowDrop"
              @node-drag-start="handleDragStart"
              @node-drop="handleDrag"
              draggable
            >
              <div style="width:95%;" class="clearfix" slot-scope="{ data }">
                <div :title="data.group_name" class="group-name">
                  {{ data.group_name }}
                </div>
                <div
                  class="button-right"
                  v-if="!data.flag || (data.flag && data.id === 0)"
                >
                  <el-button
                    @click.stop="() => groupModel(data)"
                    style="padding-right:10px"
                    type="text"
                    size="mini"
                  >
                    <i class="el-icon-edit-outline"></i>
                  </el-button>
                </div>
              </div>
            </el-tree>
            <div class="group-li" @click="setMenuFn(recycleBin)">
              <p>回收站</p>
            </div>
          </div>
        </dl>
      </template>
      <!-- tabs -->
      <template slot="column-right">
        <el-tabs v-model="activTab" @tab-click="handleTabClick(activTab)">
          <el-tab-pane label="朋友圈" name="7" v-if="ep.isCoach">
            <friend-circle-materials
              v-if="activTab == '7'"
              :filterObj="filterObj"
              ref="friendCircle"
              @isupdate="isupdateFn"
            ></friend-circle-materials>
          </el-tab-pane>
          <el-tab-pane label="小红书" name="3">
            <link-materials
              v-if="activTab == '3'"
              :filterObj="filterObj"
              ref="link"
            ></link-materials>
          </el-tab-pane>
          <el-tab-pane label="图片" name="2">
            <picture-materials
              ref="picture"
              v-if="activTab == '2'"
              :filterObj="filterObj"
            ></picture-materials>
          </el-tab-pane>
          <el-tab-pane label="视频" name="5">
            <videoMaterials
              v-if="activTab == '5'"
              :filterObj="filterObj"
              ref="video"
            ></videoMaterials>
          </el-tab-pane>
          <el-tab-pane label="文本" name="1">
            <text-materials
              ref="text"
              v-if="activTab == '1'"
              :filterObj="filterObj"
            ></text-materials>
          </el-tab-pane>
          <el-tab-pane label="文件" name="6">
            <file-materials
              v-if="activTab == '6'"
              :filterObj="filterObj"
              ref="file"
              @isupdate="isupdateFn"
            ></file-materials>
          </el-tab-pane>
          <el-tab-pane label="小程序" name="4">
            <applets-materials
              :filterObj="filterObj"
              ref="applets"
              v-if="activTab == '4'"
            ></applets-materials>
          </el-tab-pane>
        </el-tabs>
      </template>
    </layout-list>
    <select-materials
      v-if="isShowDialog"
      @submit="closeMaterialDialog"
    ></select-materials>
    <text-material-dialog
      v-if="textDialog"
      @textData="closetextDialog"
    ></text-material-dialog>
    <pic-mterial-dialog
      v-if="picDialog"
      @submit="closepicDialog"
    ></pic-mterial-dialog>
    <small-mterial-dialog
      v-if="smallDialog"
      @submit="closesmallDialog"
    ></small-mterial-dialog>
    <edit-group
      :list="editList"
      @closeDialog="closeGroupEdit"
      v-if="showEditGroup"
    ></edit-group>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LayoutList from "@/components/public/layout-list/index.vue";
import textMaterials from "../../components/text-materials.vue";
import addEditGroup from "../../components/add-edit-group.vue";
import pictureMaterials from "../../components/picture-materials.vue";
import fileMaterials from "../../components/file-materials.vue";
import videoMaterials from "../../components/video-materials.vue";
//链接素材库
import linkMaterials from "../../components/link-materials.vue";
//小程序素材库
import appletsMaterials from "../../components/applets-materials.vue";
//文本素材库
import textMaterialDialog from "../../components/text-material-dialog.vue";
//图片素材库
import picMterialDialog from "../../components/pic-mterial-dialog.vue";
//新增小程序
import smallMterialDialog from "../../components/small-mterial-dialog.vue";

import selectMaterials from "../../components/link-material-dialog.vue";
// import selectMaterials from "../components/video-material-dialog.vue";
import editGroup from "../../components/edit-group.vue";
import Api from "@/api/enterprise/material-library";
import { shouldRequestSourceGroupMove } from "./source-group-drag";
// todo:新增朋友圈素材库
import friendCircleMaterials from "../../components/friend-circle-material.vue";
import PerModule from "@/store/modules/permission";
@Component({
  name: "MaterialLibrary",
  components: {
    editGroup,
    LayoutList,
    textMaterials,
    addEditGroup,
    pictureMaterials,
    fileMaterials,
    videoMaterials,
    linkMaterials,
    appletsMaterials,
    smallMterialDialog,
    textMaterialDialog,
    picMterialDialog,
    selectMaterials,
    friendCircleMaterials
  }
})
export default class extends Vue {
  private showUploaderFilters = false; // 上传人及门店筛选暂不展示
  // 素材title
  private materialTitle = "文本";
  private isShowDialog = false;
  private textDialog = false;
  private smallDialog = false;
  private picDialog = false;
  private showEditGroup = false;
  // 搜索区域
  private search = {
    keyWord: "",
    uploader_name: "",
    department_name: "",
    datePicker: "",
    checked: false, // 是否显示自己上传的素材
    category_ids: [], // 素材标签
    source_status: "", // 状态
    sku_item: "", // 商品款号
    product_name: "", // 商品名称
    sku_color: "", // 商品颜色
    goods_relation: 0, // 关联商品状态
    source_origin: 0, // 素材来源
  };
  private addtit = "添加一级分组";
  private rowData = {}; //分组每一项

  // 是否打开添加修改分组弹框
  private pushOrEditGroup = false;
  // tab
  private activTab: any = "";
  // 素材分组
  private groupData: any[] = [];
  private originData: any[] = [
    { group_name: "全部分组", id: 0, children: {}, flag: true },
    { group_name: "未分组", id: -1, children: null, flag: true }
  ];
  private recycleBin: any = {
    group_name: "回收站",
    id: -2,
    children: null,
    flag: true
  };
  private isActive = 0;
  private filterObj: any = {
    keyword: "",
    uploader_name: "",
    department_name: "",
    create_end: 0,
    create_start: 0,
    children: {},
    just_me: 0, //查询条件
    source_status: 2, // 2全部 1 禁用  0 启用
    source_origin: 0
    
  };
  private dragStartContent: any = {};
  private dragStartSiblingIds: any[] = [];
  private total = 0;
  // 产品分类
  private category_options: any[] = [];
  private editList: any[] = [];
  private status_options: any[] = [
    { value: 0, label: "已启用" },
    { value: 1, label: "已禁用" }
  ];
  private relation_status_options: any[] = [
    { value: 0, label: "全部" },
    { value: 1, label: "已关联" },
    { value: 2, label: "未关联" }
  ];
  private source_origin_options: any[] = [
    { value: 0, label: "全部" },
    { value: 1, label: "本地创建" },
    { value: 2, label: "CH同步" }
  ];
  

  /**
   * 当前品牌标识
   */
  get ep () {
    return this.$store.getters['enterprise']
  }
  //----------------------------------------------------------------methods-------------------------------------------------------------------
  isupdateFn(val: string) {
    console.log(val, "val");
    if (val === "video") {
      this.activTab = "5";
      setTimeout(() => {
        (this.$refs.video as any).getlist();
      });
    } else if (val === "img") {
      this.activTab = "2";
      setTimeout(() => {
        (this.$refs.picture as any).getlist(1);
      });
    } else if(val === "friendCircle") {
      this.activTab = "7";
      setTimeout(() => {
        (this.$refs.friendCircle as any).getlist(1);
      });
    }
  }
  closeGroupEdit() {
    this.showEditGroup = false;
    this.grouplist();
  }
  // ------------------------------------------------
  allowDrop(draggingNode: any, dropNode: any, type: any) {
    if (draggingNode.data.parent_id === dropNode.data.parent_id) {
      console.log("允许");
      return type === "prev";
    } else {
      // 不同级进行处理
      console.log("不允许");
      return false;
    }
  }
  // 节点开始拖拽时触发的事件
  handleDragStart(node: any) {
    this.dragStartContent = node;
    this.dragStartSiblingIds = node.parent.childNodes.map(
      (item: any) => item.data.id
    );
  }
  // 拖拽成功完成时触发的事件
  handleDrag(draggingNode: any, dropNode: any, dropType: string) {
    const currentSiblingIds = dropNode
      ? dropNode.parent.childNodes.map((item: any) => item.data.id)
      : this.dragStartSiblingIds;
    if (
      !shouldRequestSourceGroupMove(
        dropNode,
        dropType,
        this.dragStartSiblingIds,
        currentSiblingIds
      )
    ) {
      return;
    }

    Api.moveGroupListApi({
      parent_id: this.dragStartContent.data.parent_id,
      mover_id: this.dragStartContent.data.id,
      mover_order: this.dragStartContent.data.orders,
      swaper_id: dropNode.data.id,
      swaper_order: dropNode.data.orders
    }).then((res: any) => {
      if (res.code == 200) {
        this.$message.success("移动成功");
      } else {
        this.$message.error(res.message);
      }
      // this.grouplist();
    });
  }
  // ------------------------------------------------
  groupModel(command: any) {
    this.showEditGroup = true;
    this.editList = command;
  }
  // ------------------------------------------------------
  setLevel(arr, level) {
    arr.forEach(item => {
      item.level = level;
      Array.isArray(item.children) && this.setLevel(item.children, level + 1);
    });
    return arr;
  }
  // 分组接口
  grouplist() {
    Api.getSourceGroupListApi({}).then((res: any) => {
      if (res.code == 200) {
        const data = this.setLevel(res.data, 1);
        this.groupData = [...this.originData, ...data];
        // console.log(this.groupData, "所有分组！");
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 点击添加分组
  private addgroup() {
    this.addtit = "添加一级分组";
    this.pushOrEditGroup = true;
  }
  private closeaddgroup() {
    this.pushOrEditGroup = false;
    this.grouplist();
  }
  // 点击分组列表
  private setMenuFn(data: any) {
    this.filterObj.source_group_id = data.id;
    this.filterObj.name = data.group_name;
    this.filterObj.children = data.children;
    // this.isActive = index;
    this.resetFn();
  }
  // 点击三点
  private handleCommand(command: any) {
    if (command.id) {
      this.addtit = "修改分组";
      this.rowData = command;
      this.pushOrEditGroup = true;
      return;
    }
    if (command.indexOf("delete") > -1) {
      const id = Number(command.replace("delete_", ""));
      this.delGroup(id);
    }
  }
  private delGroup(data: any) {
    // const arr: any = [];
    // arr.push(data.id);
    this.filterObj.source_group_id = data.id;
    this.getlist();
    if (data.children !== null) {
      this.$message.error("该分组下有子级分组,不能删除！");
    } else {
      this.$confirm("你确认要删除该分组吗？", "删除分组", {
        confirmButtonText: "确定删除",
        cancelButtonText: "我再想想",
        type: "warning"
      }).then(() => {
        Api.sourceGroupDeleteApi({ id: data.id }).then((res: any) => {
          if (res.code == 200) {
            this.isActive = 1;
            this.filterObj.source_group_id = -1;
            this.resetFn();
            this.grouplist();
            this.$message({
              type: "success",
              message: "删除成功!"
            });
          } else {
            this.$message.error(res.message);
          }
        });
      });
    }
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
  // 点击查询
  touchSearch() {
    this.filterObj.just_me = this.search.checked ? 1 : 0; //仅显示自己 1 是 0 否
    this.filterObj.keyword = this.search.keyWord; //关键字
    this.filterObj.uploader_name = this.search.uploader_name.trim(); // 上传人名称
    this.filterObj.department_name = this.search.department_name.trim(); // 上传人所属门店
    this.filterObj.category_ids = this.search.category_ids; // 产品分类
    this.filterObj.source_status =
      this.search.source_status === "" ? 2 : this.search.source_status; // 状态
    this.filterObj.create_start =
      this.search.datePicker && this.search.datePicker.length
        ? (this.search.datePicker[0] as any).getTime() / 1000
        : 0; //开始时间
    this.filterObj.create_end =
      this.search.datePicker && this.search.datePicker.length
        ? (this.search.datePicker[1] as any).getTime() / 1000
        : 0; //结束时间

    this.filterObj.sku_item = this.search.sku_item; // 商品款号
    this.filterObj.sku_color = this.search.sku_color; // 商品颜色
    this.filterObj.product_name = this.search.product_name; // 商品名称
    this.filterObj.goods_relation = this.search.goods_relation; // 关联商品状态
    this.filterObj.source_origin = this.search.source_origin; // 素材来源
    this.resetFn();
  }
  // 点击重置
  reset() {
    this.search.keyWord = "";
    this.search.uploader_name = "";
    this.search.department_name = "";
    this.search.datePicker = "";
    this.search.checked = false;
    this.isActive = 0;
    this.filterObj.just_me = 0; //仅显示自己 1 是 0 否
    this.filterObj.keyword = ""; //关键字
    this.filterObj.uploader_name = "";
    this.filterObj.department_name = "";
    this.filterObj.create_start = 0; //开始时间
    this.filterObj.create_end = 0;
    this.filterObj.source_group_id = 0;
    this.filterObj.category_ids = [];

    // 删除状态 全部数据2  0 已启用，1已禁用
    this.filterObj.source_status = 2;
    this.search.category_ids = [];
    this.search.source_status = "";

    this.filterObj.sku_item = ""; // 商品款号
    this.filterObj.product_name = ""; // 商品名称
    this.filterObj.sku_color = ""; // 商品颜色
    this.filterObj.goods_relation = 0; // 关联商品状态
    this.filterObj.source_origin = 0; // 素材来源
    this.search.sku_item = "";
    this.search.sku_color = ""; // 商品颜色 
    this.search.goods_relation = 0; // 关联商品状态
    this.search.source_origin = 0; // 素材来源
    this.search.product_name = "";
    this.resetFn();
  }
  //重置调用方法
  private resetFn() {
    if (this.activTab == "1") {
      (this.$refs.text as any).getlist(1);
    } else if (this.activTab == "2") {
      (this.$refs.picture as any).getlist(1);
    } else if (this.activTab == "3") {
      (this.$refs.link as any).getlist(1);
    } else if (this.activTab == "4") {
      (this.$refs.applets as any).getAppletsListFn(true);
    } else if (this.activTab == "5") {
      (this.$refs.video as any).getlist(1);
    } else if (this.activTab == "6") {
      (this.$refs.file as any).fileListFn(true);
    }else if(this.activTab == "7") {
      (this.$refs.friendCircle as any).getlist(1);
    }
  }
  private closeMaterialDialog() {
    this.isShowDialog = false;
  }
  // 文本素材弹框关闭
  private closetextDialog() {
    this.textDialog = false;
  }
  // 图片素材弹框关闭
  private closepicDialog() {
    this.picDialog = false;
  }
  // 小程序素材弹框关闭
  private closesmallDialog() {
    this.smallDialog = false;
  }
  private handleTabClick(type: string) {
    this.activTab = type;
    // if (this.activTab == "4" || this.activTab == "6") {
    //   this.filterObj.source_group_id = this.groupData[this.isActive].id;
    // }
  }
  // 获取数据
  async getCategory() {
    const param: any = {
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
  mounted() {
    // this.materiallist.name = "全部";
    // this.materiallist.isShowTree = true;
    this.grouplist();
    this.getCategory();
    if (this.$route.query.type) {
      this.activTab = this.$route.query.type;
    } else {
      this.activTab = localStorage.getItem("activTab") || "7";
      localStorage.removeItem("activTab");
    }
  }

  activated() {
    this.resetFn();
  }
}
</script>
<style lang="scss" scoped>
.material-search-filters {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  column-gap: 20px;

  /deep/ .search-item {
    margin-right: 0;
    white-space: nowrap;
  }
}
/deep/.el-tooltip__popper.is-light {
  border: none !important;
}
/deep/.el-tooltip__popper.is-light {
  background: #ffffff;
  border: 1px solid white;
}
/deep/.el-tooltip__popper .popper__arrow {
  border: none !important;
}
/deep/.el-dropdown-menu__item.color-red {
  color: $color-danger;
}
/deep/.el-tree {
  background: #f7f7f7;
  padding-top: 0 !important;
  .is-current > .el-tree-node__content {
    background-color: #ebebeb;
    font-weight: bold;
  }
}
/deep/.el-tree-node__content:hover {
  background-color: #ffffff;
}
.menu-left {
  // margin-right: 6px !important;
  flex: 0 0 400px;
  margin: 0;
  background: $color-background;
  padding-top: 15px;
  dt {
    background: $color-background;
    font-size: 14px;
    // padding: 20px;
    width: 184px;
    height: 36px;
    line-height: 36px;
    margin: 0 auto 15px;
    text-align: center;
    cursor: pointer;
    color: #303133;
    border: 1px solid #000000;
    div {
      i {
        margin-right: 5px;
        font-size: 14px;
      }
    }
  }
  .sub-wrap {
    height: 610px;
    overflow-y: auto;
    .group-name {
      width: calc(100% - 30px);
      overflow: hidden;
      float: left;
      line-height: 35px;
      height: 35px;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
    .button-right {
      float: right;
      width: 30px;
      .el-button--text {
        color: #4d4d4d !important;
        i {
          font-weight: 600;
        }
      }
    }
    .group-li {
      // background: #ffffff;
      color: #606266;
      border-top: 1px solid #e4e7ed;
      cursor: pointer;
      padding: 10px 0;
      p {
        line-height: 30px;
        padding-left: 38px;
        &:hover {
          background-color: #f5f7fa;
        }
      }
    }
    dd {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 17px 0 20px;
      line-height: 20px;
      span {
        padding-left: 20px;
      }
      .menu-wrapper {
        height: 17px;
        .operation-box {
          width: 23px;
          height: 17px;
          img {
            width: 3px;
            height: 17px;
            margin: 0 10px;
          }
        }
      }
    }
    dd.check-color span {
      border-left: 4px solid #708df5;
      color: #708df5;
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
  .el-tree-node:focus > .el-tree-node__content {
    background-color: #ebebf1 !important;
  }
}
.btn-ul {
  .btn-li {
    line-height: 30px;
    cursor: pointer;
    text-align: center;
    &:hover {
      background: $color-background;
    }
  }
}
</style>
