<template>
  <div class="select-materials-dialog">
    <el-dialog
      :visible.sync="selectDialogVisible"
      width="1020px"
      @close="cancel"
    >
      <div class="header">
        <div class="pointer-s">
          <span
            :class="{ 'chosoe-clas': currentAdd }"
            @click="goback"
          >选择素材<i
              v-if="currentAdd"
              class="el-icon-arrow-right m-l-5 m-r-5"
            ></i></span>
          <span
            v-if="currentAdd"
            class="add-clas"
          >新增{{ activTab | typeNameFilter }}</span>
        </div>
        <span
          class="pointer-s"
          @click="cancel"
        ><i
            aria-hidden="true"
            class="el-icon-close"
          ></i></span>
      </div>
      <div v-if="!currentAdd">
        <el-tabs
          v-model="activTab"
          @tab-click="handleTabClick(activTab)"
        >
          <el-tab-pane
            label="图片"
            name="1"
          ></el-tab-pane>
          <el-tab-pane
            label="小红书"
            name="2"
          ></el-tab-pane>
          <el-tab-pane
            label="小程序"
            name="3"
          ></el-tab-pane>
          <el-tab-pane
            label="视频"
            name="4"
          ></el-tab-pane>
          <el-tab-pane
            v-if="file"
            label="文件"
            name="5"
          ></el-tab-pane>
        </el-tabs>
      </div>
      <!-- 默认展示 -->
      <div v-if="!currentAdd">
        <div class="content">
          <div class="menu-left">
            <div class="add-new">
              <el-button
                type="primary"
                @click="addBtn"
                class="addClass"
                size="small"
              >新增{{ activTab | typeNameFilter }}</el-button>
            </div>
            <div class="sub-wrap">
              <el-tree
                :data="groupData"
                node-key="id"
                ref="treeRef"
                :default-expand-all="false"
                :expand-on-click-node="false"
                @node-click="handleNodeClick"
                :highlight-current="true"
                :current-node-key="group_id"
              >
                <template slot-scope="{ node, data }">
                  <span v-if="data.group_name.length <= 9">{{ data.group_name }}</span>
                  <el-popover
                    placement="left"
                    v-else
                    title=""
                    width="100%"
                    trigger="hover"
                    :content="data.group_name"
                  >
                    <p slot="reference" class="text-span">
                      {{ data.group_name }}
                    </p>
                  </el-popover>
                </template>
              </el-tree>
            </div>
          </div>
          <div class="right-content">
            <div class="search">
              <search-item text="搜索：">
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
              <search-item text="商品名称：">
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
              <search-item text="商品款号：">
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
              <search-item text="商品颜色：">
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
            </div>
            <div
              class="right-box"
              v-if="dataList.length > 0"
            >
              <div class="right-box-content"> 
                <div
                  class="right-content-content"
                  v-for="(item, index) in dataList"
                  :key="index"
                  @click="selectedItem(item, item.fileType)"
                >
                  <div
                    v-if="activTab == 1"
                    :class="checkGroup.indexOf(item.url) > -1
                        ? 'img-selected-item pic-item'
                        : 'pic-item'
                      "
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
                      <span class="text">{{ item.name}}</span>
                    </el-tooltip>
                    <el-checkbox
                      class="item-select"
                      v-if="checkGroup.indexOf(item.url) > -1"
                      :value="true"
                    ></el-checkbox>
                  </div>
                  <div
                    v-if="activTab == 2"
                    :class="checkGroup.indexOf(item.cover_url) > -1
                        ? 'selected-item link-item'
                        : 'link-item'
                      "
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
                          <el-tag v-for="(v, index) in item.sku_list" :key="index" size="small" type="gray">{{ v }}</el-tag>
                        </div>
                        <el-popover v-else placement="top" width="200px" trigger="hover">
                          <div class="tags h" slot="reference">
                            <el-tag v-for="(v, index) in item.sku_list.slice(0,2)" :key="index" size="small" type="gray">{{ v }}</el-tag>
                            <span class="more">...</span>
                          </div>
                          <div class="tags w">
                            <el-tag v-for="(v, index) in item.sku_list" :key="index" type="gray">{{ v }}</el-tag>
                          </div>
                        </el-popover>
                      </div>
                      <span class="item-img">
                        <img
                          :src="item.cover_url"
                          alt=""
                        />
                      </span>
                    </div>
                    <el-checkbox
                      class="item-select"
                      v-if="checkGroup.indexOf(item.cover_url) > -1"
                      :value="true"
                    ></el-checkbox>
                  </div>
                  <div
                    v-if="activTab == 3"
                    :class="checkGroup.indexOf(item.cover_url) > -1
                        ? 'selected-item small-item'
                        : 'small-item'
                      "
                  >
                    <div class="small-tit">
                      <img
                        class="minilogo"
                        src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/mini-applet-icon.png"
                        alt=""
                      /><span class="smallName">{{ item.app_name }}</span>
                    </div>
                    <p class="tit-name">
                      {{ item.title }}
                    </p>
                    <img
                      :src="item.cover_url"
                      class="miniimg"
                      alt=""
                    />
                    <el-checkbox
                      class="item-select"
                      v-if="checkGroup.indexOf(item.cover_url) > -1"
                      :value="true"
                    ></el-checkbox>
                  </div>
                  <div
                    v-if="activTab == 4"
                    :class="checkGroup.indexOf(item.url) > -1
                        ? 'selected-item video-item'
                        : 'video-item'
                      "
                  >
                    <div class="item-img">
                      <video :src="item.url" />
                    </div>
                    <div class="item-info">
                      <div class="item-title">
                        {{ item.name }}
                      </div>
                      <div class="item-time">
                        {{ (item.size / 1024 / 1024).toFixed(2) }}MB
                      </div>
                    </div>
                    <el-checkbox
                      class="item-select"
                      v-if="checkGroup.indexOf(item.url) > -1"
                      :value="true"
                    ></el-checkbox>
                  </div>
                  <div
                    v-if="activTab == 5"
                    :class="checkGroup.indexOf(item.url) > -1
                        ? 'selected-item file-item'
                        : 'file-item'
                      "
                  >
                    <div
                      class="file-type"
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
                    <div class="file-title">
                      <el-popover
                        placement="top-start"
                        v-if="item.name.length >= 8"
                        title=""
                        width="200px"
                        trigger="hover"
                        :content="item.name"
                      >
                        <h3 slot="reference">{{ item.name.slice(0, 8) }}...</h3>
                      </el-popover>
                      <h3 v-else>{{ item.name }}</h3>
                    </div>
                    <div class="date-time">
                      <span>{{
                        item.created_at | formatDate("yyyy-mm-dd hh:mm:ss")
                      }}</span>
                    </div>
                    <div class="btns">
                      <el-button
                        size="small"
                        type="text"
                        class="unload"
                        @click="downloadFn(item.url, item.name)"
                      >下载</el-button>
                      <el-popconfirm
                        title="你确定删除该素材吗？"
                        @onConfirm="getDeleteFn(item.id)"
                        @confirm="getDeleteFn(item.id)"
                        icon="el-icon-info"
                        confirm-button-text="确认"
                        cancel-button-type="Primary"
                      >
                        <el-button
                          slot="reference"
                          size="small"
                          type="text"
                          class="del-btn"
                        >删除</el-button>
                      </el-popconfirm>
                    </div>
                    <el-checkbox
                      class="item-select"
                      v-if="checkGroup.indexOf(item.url) > -1"
                      :value="true"
                    ></el-checkbox>
                  </div>
                </div>
              </div>
            </div>
            <div
              class="empty-content"
              v-if="dataList.length == 0"
            >
              <img
                class="empty-img"
                v-if="!emptyFlag"
                src="@/assets/images/empty-box.png"
                alt=""
              />
              <div
                class="empty-box"
                v-if="!emptyFlag"
              >
                <div class="desc">暂无素材</div>
              </div>
              <img
                class="empty-img"
                v-if="emptyFlag"
                src="@/assets/images/empty-search.png"
                alt=""
              />
              <div
                class="empty-box"
                v-if="emptyFlag"
              >
                <div class="desc">没有找到相关链接素材</div>
              </div>
            </div>
            <div
              v-if="dataList.length > 0"
              class="fotter"
            >
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
      </div>
      <span
        v-if="!currentAdd"
        slot="footer"
        class="dialog-footer"
      >
        <el-button
          @click="cancel"
          size="small"
        >取 消</el-button>
        <el-button
          type="primary"
          @click="submitPar"
          size="small"
        >确 认</el-button>
      </span>
      <!-------------------新增模块 -------------------------->
      <add-img
        v-if="currentAdd == 1"
        :newGroupData="newGroupData"
        @submit="addsure"
      ></add-img>
      <add-link
        v-if="currentAdd == 2"
        :newGroupData="newGroupData"
        @submit="addsure"
      ></add-link>
      <add-wxapp
        v-if="currentAdd == 3"
        :newGroupData="newGroupData"
        @submit="addsure"
      ></add-wxapp>
      <add-video
        v-if="currentAdd == 4"
        :newGroupData="newGroupData"
        @submit="addsure"
      ></add-video>
      <add-file
        v-if="currentAdd == 5"
        :newGroupData="newGroupData"
        @submit="addsure"
      ></add-file>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import Api from "@/api/enterprise/material-library";
import { downLoad } from "@/utils/down-load";
import { isNull } from "@/utils/tools";
import addImg from "./add-img.vue";
import addLink from "./add-link.vue";
import addWxapp from "./add-wxapp.vue";
import addVideo from "./add-video.vue";
import addFile from "./add-file.vue";
import {
  Credential,
  getFileCredentialApi
} from "@/api/enterprise/sensitive-word";
import COSSDK from "cos-js-sdk-v5";

import { ApiResult } from "@/types";
@Component({
  filters: {
    typeNameFilter(val: string) {
      switch (Number(val)) {
        case 1:
          return "图片";
        case 2:
          return "链接";
        case 3:
          return "小程序";
        case 4:
          return "视频";
        case 5:
          return "文件";
      }
    }
  },
  components: {
    addImg,
    addLink,
    addWxapp,
    addVideo,
    addFile
  }
})
export default class extends Vue {
  @Prop({ default: [] }) chooseInfo?: any;
  @Prop({ default: 1 }) type?: number;
  @Prop({ default: false }) file?: boolean;
  //------------------------------------------------------data-----------------------------------------------------------/
  private currentAdd = 0;
  private activTab = "1";
  private selectDialogVisible = true;
  // 搜索区域
  private search = {
    keyWord: "",
    sku_item: "",
    product_name: "",
    sku_color: ""
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
  private dataList: any = [];
  private isActive = 0;
  private isChecked: any = [];
  //分页
  private meta = {
    size: 12,
    page: 1,
    total: 0
  };
  private emptyFlag = false;
  private selectedMaterial: any = [];
  private customIcon = ["W", "X", "P", "PDF"]; //自定义icon
  private cos: any = null;
  private credential?: Credential;
  private checkGroup: any = [];
  // ---------新增链接data-----------------
  private group_id: any = 0;
  private formData = {
    title: "", //素材标题
    description: "", //素材描述
    group_id: "",
    link_type: 2, //(1 URL 2 自定义)
    cover_url: "", //封面图URL
    link_url: "", //链接URL
    author: "", //作者
    content: "" //正文内容
  };
  private rules = {
    link_url: [
      {
        validator: this.validateself,
        trigger: "blur"
      }
    ]
  };
  private validateself(rule: any, value: any, callback: any) {
    const match = /^(https?:\/\/)(([a-zA-Z0-9_-]+\.)+[a-zA-Z]{2,})(:\d+)?(\/[^\s?#]*)?(\?[^\s#]*)?(#[^\s]*)?$/;
    const flag = match.test(value);
    flag ? callback() : callback(new Error("请输入正确的URL"));
  }
  // 富文本编辑器
  private editorSms: any = null; //编辑器
  //----------------------------------------------------methods-----------------------------------------------------------/
  private addsure() {
    this.currentAdd = 0;
    this.getlist(0);
  }
  // 初始化腾讯对象存储Cos
  private async initCos() {
    const {
      code,
      message = "",
      data
    }: ApiResult<Credential> = await getFileCredentialApi();
    if (code !== 200) {
      this.$message(message);
      return false;
    }
    this.credential = data;
    this.cos = new COSSDK({
      getAuthorization: (options: any, callback: Function) => {
        console.log(options);
        callback({
          TmpSecretId: data?.secret_id,
          TmpSecretKey: data?.secret_key,
          XCosSecurityToken: data?.session_token,
          StartTime: data?.start_time,
          ExpiredTime: data?.expired_time
        });
      }
    });

    return true;
  }
  //下载
  private downloadFn(url: string, name: string) {
    const sliceUrl = url.slice(29);
    // url = uploads/2020/11/16/MDCRncszxDmRnkQxXzS4SE4fRteA43WW.mp4
    const Url = this.cos.getObjectUrl({
      Bucket: this.credential?.bucket,
      Region: this.credential?.region /* 存储桶所在地域，必须字段 */,
      Key: sliceUrl
    });
    downLoad(Url, name);
  }
  //删除
  private getDeleteFn(id?: number) {
    const ids: any = [];
    if (id) {
      ids.push(id);
    }
    Api.materialDeleteApi({ ids: ids }).then((res: any) => {
      if (res.code == 200) {
        this.getlist(1);
        if (this.isChecked.indexOf(id) > -1) {
          this.isChecked.splice(this.isChecked.indexOf(id), 1);
          this.selectedMaterial.splice(this.isChecked.indexOf(id), 1);
        }
        this.$message({
          type: "success",
          message: "删除成功!"
        });
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 列表数据
  private async getlist(type: number) {
    let filetype = 0;
    let type_name = "";
    if (this.activTab == "1") {
      filetype = 1;
      type_name = "url";
    } else if (this.activTab == "2") {
      filetype = 2;
      type_name = "title";
    } else if (this.activTab == "3") {
      filetype = 3;
      type_name = "title";
    } else if (this.activTab == "4") {
      filetype = 4;
      type_name = "name";
    } else if (this.activTab == "5") {
      filetype = 5;
      type_name = "name";
    }
    const res: any = await Api.consoleList({
      page: this.meta.page,
      size: this.meta.size,
      just_me: this.checked ? 1 : 0,
      // 分组ID  99 全部 -1 未分组
      source_group_id: parseInt(this.group_id),
      keyword: this.search.keyWord,
      sku_item: this.search.sku_item,
      product_name: this.search.product_name,
      sku_color: this.search.sku_color,
      type: filetype + 1,
      is_use: 1
    });
    console.log(res, type);
    if (res.code != 200) {
      this.$message.error(res.message);
      return;
    }
    const data = res.data.list;
    data.map((v: any) => {
      v.fileType = filetype;
      if (!v.title) {
        v.title = v[type_name];
      }
    });
    const { size, page, total } = res.data;
    if (data.length <= 0 || isNull(data)) {
      this.dataList = [];
      return;
    }
    this.meta = { size, page, total };
    this.dataList = data;
  }
  // 切换tab
  private handleTabClick(type: string) {
    this.activTab = type;
    this.meta.page = 1;
    this.meta.size = 12;
    this.group_id = 0;
    this.isActive = 0;
    (this.$refs.treeRef as any).setCurrentKey(0);
    this.getlist(0);
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
  private selectedItem(item: any, type: number) {
    console.log(item, type);
    if (type == 2 || type == 3) {
      if (this.checkGroup.indexOf(item.cover_url) > -1) {
        this.selectedMaterial.splice(
          this.checkGroup.indexOf(item.cover_url),
          1
        );
        this.checkGroup.splice(this.checkGroup.indexOf(item.cover_url), 1);
        return;
      }
    } else {
      if (this.checkGroup.indexOf(item.url) > -1) {
        this.selectedMaterial.splice(this.checkGroup.indexOf(item.url), 1);
        this.checkGroup.splice(this.checkGroup.indexOf(item.url), 1);
        return;
      }
    }
    if (this.checkGroup.length >= 9) {
      this.$message.warning("素材总数不能超过9个");
      return;
    }
    if (type == 2 || type == 3) {
      this.checkGroup.push(item.cover_url);
    } else {
      this.checkGroup.push(item.url);
    }
    this.selectedMaterial.push(item);
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
    if (this.selectedMaterial.length == 0) {
      this.$message.warning("请选择素材");
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
    this.search.product_name = "";
    this.search.sku_color = "";
    this.checked = false;
    this.meta.page = 1;
    this.getlist(1);
  }
  // 点击新增
  private addBtn() {
    const type = Number(this.activTab) + 1;
    switch (type) {
      case 2:
        window.open(`${origin}/web-admin/myapp/enterprise/material-library/create-image?from=1`, "_blank");
        break;
      case 3:
        window.open(`${origin}/web-admin/myapp/enterprise/material-library/create-link?from=1`, "_blank");
        break;
      case 5:
        window.open(`${origin}/web-admin/myapp/enterprise/material-library/create-video?from=1`, "_blank");
        break;
      default:
        window.open(`${origin}/web-admin/myapp/enterprise/material-library/material-list/list?isAddTask=1&type=${type}`, "_blank");
        break;
        
    }
  }
  //
  private picSuccess(res: any) {
    this.$set(this.formData, "cover_url", res);
    (this.$refs["cform"] as any).validateField("cover_url");
  }
  // 上传文件限制规范
  private beforeUpload(file: any) {
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
  // index: number, id: any
  private handleNodeClick(e) {
    const { id } = e;
    
    
    this.group_id = id;
    
    // this.isChecked = 0;
    this.meta.page = 1;
    if (id != 0 && id != -1) {
      this.formData.group_id = id;
    } else {
      this.formData.group_id = "";
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
    this.currentAdd = 0;
    // this.isChecked = 0;
  }
  private mounted() {
    this.activTab = String(this.type);
    this.initCos();
    this.getlist(0);
    this.grouplist();
    this.chooseInfo.forEach((v: any) => {
      if (v.fileType == 1 || v.fileType == 4 || v.fileType == 5) {
        this.checkGroup.push(v.url);
      } else if (v.fileType == 3 || v.fileType == 2) {
        this.checkGroup.push(v.cover_url);
      }
    });
    this.selectedMaterial = JSON.parse(JSON.stringify(this.chooseInfo));
  }
}
</script>
<style lang="scss" scoped>
/deep/ .el-tree {
  padding: 0 !important;

  .el-tree-node__content { 
    height: 34px;
  }
}
/deep/ .el-tree--highlight-current .el-tree-node.is-current > .el-tree-node__content {
  background-color: #ebebf1 !important;
  color: #000 !important;
  
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
.tags {
  margin-top: 5px;
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
  &.h {
    height: 24px;
  }
  &.w {
    max-width: 300px;
  }
  .more {
    color: #666;
    font-size: 12px;
  }
}
.el-tag { 
  display: block;
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

.text-span {
  text-overflow: ellipsis;
  white-space: nowrap;
  overflow: hidden;
}


.select-materials-dialog {
  .content {
    display: flex;
    height: 500px;

    .menu-left {
      background: $color-background;
      width: 200px;

      .add-new {
        background-color: #fff;
        padding-bottom: 20px;
      }

      .sub-wrap {
        height: 440px;
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
          border-left: 4px solid #000000;
          color: #000000;
          font-weight: 600;
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
      padding-left: 20px;

      .right-box {
        height: 340px;
        overflow: auto;

        .right-box-content {
          width: 100%;
          display: flex;
          flex-wrap: wrap;
          gap: 20px;

        }
      }

      .right-content-content {
        // float: left;
        // margin-right: 20px;
        // margin-bottom: 15px;
      }

      // 图片
      .pic-item {
        width: 128px;

        .text {
          display: inline-block;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
          width: 100%;
        }

        .picimg-box {
          position: relative;
          overflow: hidden;
          border-radius: 4px;
          margin-bottom: 10px;
        }

        .picimg {
          display: block;
          width: 128px;
          height: 128px;
          border-radius: 4px;
          cursor: pointer;
          border: 1px solid transparent;
          box-sizing: border-box;
          object-fit: cover;
        }
      }

      .img-selected-item {
        position: relative;

        .picimg {
          border: 1px solid #000000;
        }

        .item-select {
          position: absolute;
          right: 0;
          top: -3px;
        }
      }

      // 链接
      .link-item {
        width: 225px;
        min-height: 90px;
        position: relative;
        padding: 10px;
        box-sizing: border-box;
        background: #ffffff;
        border: 1px solid #e4e7ed;
        border-radius: 4px;

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

            .tags {
              margin-top: 5px;
              display: flex;
              flex-wrap: wrap;
              gap: 5px;
              height: 24px;
              

              .more {
                color: #666;
                font-size: 12px;
              }
            }
          }

          .item-txt {
            width: 157px;
            height: 17px;
            font-size: 12px;
            font-family: PingFangSC-Regular, PingFang SC;
            font-weight: 400;
            color: #909399;
            line-height: 17px;
            text-overflow: -o-ellipsis-lastline;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
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

        .item-select {
          position: absolute;
          right: -1px;
          top: -4px;
        }
      }

      // 小程序
      .small-item {
        width: 145px;
        // height: 191px;
        position: relative;
        background: #ffffff;
        border-radius: 4px;
        padding: 8px;
        border: 1px solid #dcdfe6;
        cursor: pointer;

        &:nth-child(4n) {
          margin-right: 0;
        }

        .small-tit {
          margin-bottom: 5px;
          display: flex;

          .minilogo {
            display: inline-block;
            margin-right: 5px;
            width: 20px;
            height: 20px;
            border-radius: 50%;
          }

          .smallName {
            color: #707277;
          }
        }

        span {
          width: 100%;
          display: inline-block;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }

        .miniimg {
          display: block;
          width: 145px;
          height: 116px;
          border-radius: 4px;
        }
      }

      // 视频
      .video-item {
        user-select: none;
        display: flex;
        position: relative;
        width: 225px;
        height: 83px;
        padding: 10px 10px 0;
        box-sizing: border-box;
        background: #ffffff;
        border: 1px solid #e4e7ed;
        border-radius: 4px;

        &:hover {
          cursor: pointer;
        }

        .item-img {
          width: 50px;
          height: 50px;
          border-radius: 4px;
          // background: skyblue;
          margin-right: 5px;

          video {
            width: 100%;
            height: 100%;
          }
        }

        .item-info {
          flex: 1;
          width: 152px;
          font-size: 12px;
          font-family: PingFangSC-Regular, PingFang SC;
          font-weight: 400;
          color: #909399;

          .item-title {
            height: 34px;
            line-height: 17px;
            text-overflow: -o-ellipsis-lastline;
            overflow: hidden;
            text-overflow: ellipsis;
            display: -webkit-box;
            -webkit-line-clamp: 2; //行数需设置
            line-clamp: 2;
            -webkit-box-orient: vertical;
          }

          .item-time {
            margin-top: 5px;
          }
        }

        .item-select {
          position: absolute;
          right: -1px;
          top: -4px;
        }
      }

      // 文件
      .file-item {
        float: left;
        width: 155px;
        padding: 0 10px;
        border-radius: 4px;
        border: 1px solid $color-second-line;
        position: relative;

        .file-title {
          // line-height: 20px;
          height: 20px;
          color: #707277;
          display: flex;
          align-items: center;
          justify-content: flex-end;

          span {
            flex: 0 0 14px;
            height: 20px;
            margin-right: 5px;
          }

          h3 {
            flex: 1;
            font-size: 14px;
            color: $color-text-main;
            text-align: right;
            overflow: hidden;
            white-space: nowrap;
            text-overflow: ellipsis;
            height: 20px;
            line-height: 20px;
          }
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

        .date-time {
          text-align: right;
          font-size: 12px;
          color: #707277;
          line-height: 20px;
          margin: 5px 0 10px 0;
        }

        .btns {
          border-top: 1px solid $color-second-line;
          text-align: right;
          margin: 2px 0;
        }

        .file-bg5 {
          background-color: #c9c9c9;
        }

        .del-btn {
          margin-left: 20px;
          color: #f67270;
        }

        .item-select {
          position: absolute;
          right: -1px;
          top: -4px;
        }
      }

      .item-select {
        position: absolute;
        right: -1px;
        top: -4px;
      }

      .selected-item {
        border: 1px solid #000000;
      }
    }

    .empty-content {
      flex: 1;
      text-align: center;
      margin-top: 100px;
    }
  }

  .fotter {
    margin-top: 20px;
    text-align: right;
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

    .chosoe-clas {
      font-size: 14px;
      color: #606266;
      font-weight: 400;
      line-height: 26px;
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

  /deep/.el-dialog {
    max-width: 1200px !important;
    max-height: none !important;
  }
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

.pointer-s {
  cursor: pointer;
}</style>
