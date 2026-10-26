<template>
  <div class="select-materials-dialog">
    <el-dialog
      class="dialog-class"
      :visible.sync="selectDialogVisible"
      width="1000px"
      @close="cancel"
    >
      <div class="header">
        <div>
          <span
            @click="goback"
            :class="{ 'chosoe-clas': defaultShow }"
          >选择视频素材</span><span
            v-if="!defaultShow"
            :class="{ 'add-clas': !defaultShow }"
          >>新增视频</span>
        </div>
        <span @click="cancelDialog"><i
            aria-hidden="true"
            class="el-icon-close"
          ></i></span>
      </div>
      <!-- 默认展示 -->
      <div v-if="defaultShow">
        <div class="search">
          <search-item>
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
            @click="addBtn()"
            class="addClass"
          >新增视频</el-button>
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
          <div
            class="empty-content"
            v-if="linkList.length == 0"
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
              <div class="desc">暂无视频素材</div>
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
              <div class="desc">没有找到相关视频素材</div>
            </div>
          </div>
          <div
            class="right-content"
            v-else
          >
            <div
              :class="isChecked === item.id
                  ? 'right-item selected-item'
                  : 'right-item'
                "
              v-for="(item, index) in linkList"
              :key="index"
              @click="selectedItem(item)"
            >
              <div class="item-img">
                <video :src="item.url" />
              </div>
              <div class="item-info">
                <div class="item-title">
                  {{ item.name }}
                </div>
                <div class="item-time">
                  {{ item.h_size }}
                </div>
              </div>
              <img
                v-if="isChecked === item.id"
                class="item-select"
                src="@/assets/images/select-checked.png"
                alt=""
              />
            </div>
          </div>
          <div
            v-if="linkList.length > 0"
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
      <span
        slot="footer"
        class="dialog-footer"
        v-show="defaultShow"
      >
        <el-button
          size="small"
          @click="cancel"
        >取 消</el-button>
        <el-button
          size="small"
          type="primary"
          @click="submitPar"
        >确 认</el-button>
      </span>
      <!-------------------新增视频模块 -------------------------->
      <div v-if="!defaultShow">
        <cy-form
          ref="cform"
          :rules="rules"
          :model="formData"
          label-width="90px"
        >
          <search-item text="素材分组：">
            <el-select
              style="width:280px"
              v-model="formData.groupName"
              ref="selectBlur"
            >
              <el-option
                :value="formData.libraryGroup"
                :label="formData.groupName"
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
            label="本地视频"
            prop="cover_url"
            required
          >
            <large-file-upload
              v-model="formData.cover_url"
              accept="video/mp4"
              :beforeUpload="beforeUploadFn"
              :onSuccess="onSuccessFn"
              :onRemove="removeFn"
              :onError="errorFn"
            >
              <el-button
                type="primary"
                size="small"
              >上传视频</el-button>
            </large-file-upload>
          </cy-form-item>
        </cy-form>
      </div>
      <span
        slot="footer"
        class="dialog-footer"
        v-show="!defaultShow"
      >
        <el-button
          size="small"
          type="primary"
          @click="addConfirm"
        >确 认</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import LayoutList from "@/components/public/layout-list/index.vue";
import UploadImage from "@/components/business/upload-image/index.vue";
import { isNull } from "@/utils/tools";
import Api from "@/api/enterprise/material-library";
import LargeFileUpload from "@/components/business/large-file-upload/index.vue";
@Component({
  components: {
    LayoutList,
    UploadImage,
    LargeFileUpload,
  },
})
export default class extends Vue {
  // todo:是不是需要限制图片boolen
  @Prop({ default: false, type: Boolean }) isFriend?: boolean;
  @Prop({ default: '', type: Number }) groupSendType?: number;
  //------------------------------------------------------data-----------------------------------------------------------/
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
    { group_name: "未分组", id: -1 },
  ];
  private newGroupData: Array<any> = [];
  private newOriginData: any[] = [{ group_name: "未分组", id: -1 }]; //左侧分组
  private linkList: any = [];
  private isActive = 0;
  private isChecked = 0;
  private selectedMaterial = {};
  //分页
  private meta = {
    size: 9,
    page: 1,
    total: 10,
  };
  private emptyFlag = false;
  private defaultShow = true; //默认展示
  // -------------新增视频data--------------------
  private group_id: any = 0;
  private formData = {
    libraryGroup: "",
    groupName: "",
    // group_id: 0,
    cover_url: "", //视频封面
    name: "", //名称
    size: "", //大小
    url: "", //	视频
    ext: "", //后缀
  };
  private rules = {};
  // 视频相关变量
  private fileObj: any = {}; //上传文件信息
  private fileSuccessObj: any = {}; //上传成功信息
  private isShowBtn = true; //是否显示按钮
  private isClosePop = true; //是否关闭弹框
  private isFail = false; //是否上传失败

  // ------------------生命周期-----------
  private mounted() {
    this.getlist(0);
    this.grouplist();
  }
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
      product_name: this.search.product_name,
      sku_color: this.search.sku_color,
      is_use: 1,
      type: 5,
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
        this.meta = { size, page, total };
        this.linkList = data;
      } else {
        this.$message.error(res.message);
      }
    });
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
      this.$message.warning("请选择视频素材");
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
  private async selectedItem(item: any) {
    // todo: groupSendType 为1时查询是否存在media_id
    if (this.groupSendType == 1) {
      const mediaRes: any = await Api.getMediaIdApi({
        source_id: item.id,
        source_type: 5,
        url: item.url,
        type: this.isFriend ? 1 : 0,
      });
      if (mediaRes.code != 200) {
        this.$message.error('该素材资源不符合要求，请重新选择');
        return;
      }
    }
    if (this.isFriend) {
      const videoSize = item.size;
      if (videoSize > 1024 * 1024 * 10) {
        this.$message.error("视频大小不能超过10M,请重新选择");
        return;
      }
      // todo:需要创建一个video标签，然后获取视频的宽高
      const video = document.createElement("video");
      video.preload = 'metadata'
      video.src = item.url;
      // 定义一个promise
      const promise = new Promise((resolve) => {
        video.onloadedmetadata = () => {
          const { duration } = video;
          // 最长不超过30S
          if (duration > 30) {
            this.$message.error("视频时长不能超过30S,请重新选择");
            resolve(false);
          }
          resolve(true);
        };
        video.onerror = () => {
          this.$message.error("视频加载失败,请重新选择");
          resolve(false);
        };
      });
      // 等待promise执行完毕
      const result = await promise;
      video.remove();
      if (result) {
        // 需要将新增的标签删除掉
        this.isChecked = item.id;
        this.selectedMaterial = item;
      }
    } else {
      this.isChecked = item.id;
      this.selectedMaterial = item;
    }
  }
  // 点击新增
  private addBtn() {
    const newLink = `${origin}/web-admin/myapp/enterprise/material-library/create-video?from=2`;
    window.open(newLink, "_blank");
  }
  private picSuccess(res: any) {
    this.$set(this.formData, "cover_url", res);
  }
  //----------------------------------------------------上传视频相关方法-----------------------------------------------------------/
  // 上传文件前
  private beforeUploadFn(file: any) {
    // const fileMaxSize = 1073741824; //视频最大上传1G
    const fileMaxSize = 1024 * 1024 * 10; //视频最大上传1G 2020.11.30改为10M
    const isMP4 = file.type === "video/mp4";
    if (!isMP4) {
      this.$message.error("请上传MP4格式的视频");
      return false;
    }
    if (file.size > fileMaxSize) {
      this.$message.warning("请上传10M大小以内的视频");
      return false;
    }
    this.isShowBtn = false;
    this.isClosePop = false;
    this.fileObj = file;
    this.isFail = false; //上传失败
    return true;
  }
  // 上传错误
  private errorFn() {
    this.isClosePop = true;
    this.isFail = true; //上传失败
  }
  selectNode(node: any) {
    this.formData.groupName = node.group_name;
    this.formData.libraryGroup = node.id;
    this.group_id = node.id;
    (this.$refs.selectBlur as any).blur();
  }
  //上传成功
  private onSuccessFn(res: any) {
    this.isClosePop = true;
    this.fileSuccessObj = res.data;
    this.isFail = false;
    const { url } = res.data;
    this.$set(this.formData, "cover_url", url);
  }
  //移除上传的文件
  private removeFn() {
    this.isShowBtn = true;
  }
  //----------------------------------------------------上传视频相关方法-----------------------------------------------------------/
  private goback() {
    this.defaultShow = true;
    this.isChecked = 0;
    // this.isClosePop = false;
    this.isShowBtn = true;
  }
  private addConfirm() {
    if (this.isShowBtn) {
      this.$message.warning("请上传视频");
      return false;
    }
    if (!this.isClosePop) {
      this.$message.warning("视频正在上传中");
      return false;
    }
    if (this.isFail) {
      this.$message.error("视频上传失败");
      return false;
    }
    (this.$refs["cform"] as any).validate((valid: any) => {
      if (valid) {
        Api.videoSaveApi({
          group_id: parseInt(this.group_id) || -1,
          cover_url: "url",
          name: this.fileObj.name,
          size: Number(this.fileObj.size),
          url: this.fileSuccessObj.url,
          ext: "mp4",
        }).then((res: any) => {
          if (res.code == 200) {
            this.$message.success("新增视频成功");
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
  private setMenuFn(data: any) {
    // this.group_id = id;
    this.group_id = data.id;
    // this.isActive = index;
    this.isChecked = 0;
    this.meta.page = 1;
    if (data.id != 0 && data.id != -1) {
      this.formData.libraryGroup = data.id;
      this.formData.groupName = data.group_name;
    } else {
      this.formData.libraryGroup = "";
      this.formData.groupName = "";
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

  // 视频上传UI
  /deep/.el-progress-bar__outer {
    border-radius: 4px;
    width: 392px;
  }

  /deep/.el-progress-bar__inner {
    border-radius: 4px;
  }

  /deep/.el-progress__text {
    font-size: 14px !important;
    position: absolute;
    right: 65px;
    top: 11px;
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
    height: 332px;

    .menu-left {
      background: $color-background;
      width: 200px;

      .sub-wrap {
        height: 329px;
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
        float: left;
        display: flex;
        position: relative;
        width: 225px;
        height: 83px;
        margin-left: 20px;
        padding: 10px 10px 0;
        padding-top: 10px;
        box-sizing: border-box;
        background: #ffffff;
        border: 1px solid #e4e7ed;
        border-radius: 4px;
        margin-bottom: 20px;

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
          right: 0;
          top: 0;
          width: 14px;
          height: 14px;
        }
      }

      .selected-item {
        border: 1px solid #404040;
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

  .uploader-video {
    .uploader-btn {
      display: inline-block;
      background-color: #fbfdff;
      border: 1px dashed #c0ccda;
      border-radius: 6px;
      -webkit-box-sizing: border-box;
      box-sizing: border-box;
      width: 148px;
      height: 148px;
      line-height: 146px;
      vertical-align: top;
      cursor: pointer;
    }

    .uploader-btn:hover {
      border-color: #404040;
      color: #404040;
    }

    .upload-tip {
      color: #909399;
      font-size: 11px;
    }

    .uploader-drop {
      border: none;
      background: #fbfdff;
    }

    .el-icon-plus {
      position: relative;
      left: 45%;
      color: #404040;
    }

    .file-info {
      position: relative;
      width: 392px;

      .color-danger {
        color: #f56c6c;
      }

      .again-upload {
        position: absolute;
        right: 20%;
        top: -36px;
        z-index: 1;
        font-size: 14px;
      }

      .file-name {
        width: 112px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        position: absolute;
        left: 10px;
        top: -36px;
        z-index: 1;
      }

      .percent {
        position: absolute;
        right: 20%;
        top: -36px;
        z-index: 1;
      }

      .close {
        position: absolute;
        right: -7px;
        top: -40px;
        cursor: pointer;
        color: #c0c4cc;
        background: #ffffff;
        border-radius: 50%;
        font-size: 14px;
      }
    }
  }
}
</style>
