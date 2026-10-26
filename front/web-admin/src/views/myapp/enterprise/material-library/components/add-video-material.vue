<template>
  <div class="select-materials-dialog">
    <el-dialog
      title="上传视频"
      :visible.sync="selectDialogVisible"
      width="1000px"
      @close="cancel"
    >
      <!-------------------新增视频模块 -------------------------->
      <div class="main-box">
        <cy-form
          ref="cform"
          :rules="rules"
          :model="formData"
          label-width="90px"
          @submit.native.prevent
        >
          <!-- <cy-form-item label="素材分组" prop="group_id">
            <cy-select
              v-model="formData.group_id"
              size="medium"
              placeholder="请选择素材分组"
              search-input
              :optionsList="groupList"
              :options="{ label: 'group_name', value: 'id' }"
            >
            </cy-select>
          </cy-form-item> -->
          <cy-form-item label="素材分组">
            <el-cascader
              v-model="formData.group_id"
              :options="groupList"
              :props="{
                label: 'group_name',
                value: 'id',
                checkStrictly: true
              }"
              :style="{
                width: '400px'
              }"
              :emitPath="false"
            ></el-cascader>
            
          </cy-form-item>
          <cy-form-item label="素材标签" prop="category_ids" trigger="change">
            <el-select
              v-model="formData.category_ids"
              multiple
              style="width:400px"
              placeholder="请选择素材标签"
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
              v-model="formData.datePicker"
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
          <cy-form-item label="文件名" prop="name" required>
            <cy-input
              ref="inp"
              width-size="medium"
              placeholder="请输入文件名"
              v-model.trim="formData.name"
              maxlength="50"
              show-word-limit
              :disabled="!formData.cover_url"
            ></cy-input>
          </cy-form-item>
          <cy-form-item
            class="vidio-up-item"
            label="本地视频"
            prop="cover_url"
            required
          >
            <div v-if="!sourceId">
              <large-file-upload
                v-model="formData.url"
                accept="video/mp4"
                :beforeUpload="beforeUploadFn"
                :onSuccess="onSuccessFn"
                :onRemove="removeFn"
                :onError="errorFn"
              >
                <el-button type="primary" size="small">上传视频</el-button>
              </large-file-upload>
            </div>
            <div
              class="select-video-box"
              @click="openVideo"
              v-if="formData.url"
            >
              <img
                class="video-play-icon"
                src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/play_icon.png"
                alt=""
              />
              <el-image :src="formData.cover_url" fit="contain" class="video-preview"></el-image>
            </div>
            <div class="tip-view">
              <p class="tip-text">1. 视频格式仅支持*.mp4、*.m4v</p>
              <p class="tip-text">2. 视频不得超过10M</p>
            </div>
          </cy-form-item>

          <cy-form-item label="关联商品" >
            <select-goods :value="[]" @confirm="handleSelectGoodsConfirm(formData, $event)"></select-goods>
            <relation-goods v-if="formData.url && formData.relation_goods && formData.relation_goods.length > 0" :item="formData" @confirm="handleSelectGoodsConfirm(formData, $event)" @delete="handleDeleteGoods(formData, $event)"></relation-goods>
          </cy-form-item>
        </cy-form>
      </div>
      <span slot="footer" class="dialog-footer">
        <el-button @click="cancel" size="small">取 消</el-button>
        <el-button type="primary" @click="submitFn" size="small"
          >确 认</el-button
        >
      </span>
    </el-dialog>
    <!-- 预览视频弹框 -->
    <open-video
      v-if="is_open_video"
      :video_url="formData.url"
      @closeModel="cancelVideo"
    ></open-video>
    <add-category
      v-if="addVisible"
      @closeVisible="closeVisible"
      @submit="addCateSubmit"
    ></add-category>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import LayoutList from "@/components/public/layout-list/index.vue";
import Api from "@/api/enterprise/material-library";
import LargeFileUpload from "@/components/business/large-file-upload/index-new.vue";
import addCategory from "./add-category.vue";
import openVideo from "./com-open-video.vue";
import SelectGoods from "./select-goods.vue";
import RelationGoods from "./relation-goods.vue";
import Times from "./datePicker";
@Component({
  components: {
    LayoutList,
    LargeFileUpload,
    addCategory,
    openVideo,
    SelectGoods,
    RelationGoods,
  }
})
export default class extends Vue {
  @Prop() groupId!: number;
  @Prop() groupName!: string;
  @Prop() sourceId?: any;
  //------------------------------------------------------data-----------------------------------------------------------/
  private selectDialogVisible = true;
  // -------------新增视频data--------------------
  private group_id: any = 0;
  private groupList: any = [];
  private datePicker: any = []; //有效期
  private formData: any = {
    // group_id: "",
    group_id: 0,
    group_name: "",
    cover_url: "", //视频封面
    name: "", //名称
    size: "", //大小
    url: "", //	视频
    ext: "", //后缀
    category_ids: [], //素材标签
    datePicker: Times.times,
    start_time: Times.start, //有效期开始时间
    end_time: Times.end, //有效期结束时间
    relation_goods: [] //关联商品
  };
  private rules = {};
  // 视频相关变量
  private fileObj: any = {}; //上传文件信息
  private fileSuccessObj: any = {}; //上传成功信息
  private isShowBtn = true; //是否显示按钮
  private isClosePop = true; //是否关闭弹框
  private isFail = false; //是否上传失败
  // 素材标签
  private category_options: any[] = [];
  private addVisible = false; // 控制添加素材标签弹窗
  private is_open_video = false; // 效果预览

  //----------------------------------------------------methods-----------------------------------------------------------/
  // 获取分组
  private grouplist() {
    Api.sourceGroupListApi({}).then((res: any) => {
      if (res.code == 200) {
        const arr = [{ group_name: "未分组", id: -1 }];
        this.groupList = [...arr, ...res.data];
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 获取素材详情
  private sourceDetail() {
    Api.sourceDetail({
      source_id: this.sourceId,
      type: 5
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
          this.formData.start_time = data.start_time;
          this.formData.end_time = data.end_time;
          this.formData.datePicker = [
            data.start_time * 1000,
            data.end_time * 1000
          ];
        }
        this.formData.cover_url = data.CoverURL || "";
        this.formData.url = data.URL;
        this.formData.category_ids = categoryIds;
        this.formData.group_id = data.group_id || -1;
        this.formData.name = data.Name;
        if (!data.group_id) {
          this.formData.group_name = "未分组";
        } else if (data.group_id === -1) {
          this.formData.group_name = "未分组";
        } else {
          this.formData.group_name = data.group_name;
        }
        this.formData.relation_goods = data.relation_goods || [];
      } else {
        this.$message.error(res.message);
      }
    });
  }

  // 取消
  private cancel() {
    this.$emit("closeDialog");
  }
  //----------------------------------------------------上传视频相关方法-----------------------------------------------------------/
  private submitFn() {
    (this.$refs["cform"] as any).validate((valid: any) => {
      if (valid) {
        if (this.sourceId) {
          this.sourceUpdate();
        } else {
          this.addConfirm();
        }
      }
    });
  }

  // 编辑接口
  private sourceUpdate() {
    (this.$refs["cform"] as any).validate((valid: any) => {
      if (valid) {
        Api.sourceUpdate({
          source_type: 5,
          id: this.sourceId,
          group_id: this.formData.group_id || 0,
          start_time: this.formData.start_time,
          end_time: this.formData.end_time,
          category_ids: this.formData.category_ids,
          name: this.formData.name,
          relation_goods: this.formData.relation_goods || []
        }).then((res: any) => {
          if (res.code == 200) {
            this.$message.success("编辑成功");
            this.$emit("closeDialog");
          } else {
            this.$message.error(res.message);
          }
        });
      }
    });
  } //点击完成
  private addConfirm() {
    if (this.isShowBtn) {
      this.$message.warning("请上传视频");
      return false;
    }
    if (!this.isClosePop && !this.isShowBtn) {
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
          group_id: Number(this.formData.group_id) || 0,
          cover_url: this.fileSuccessObj.poster,
          name: this.formData.name,
          size: Number(this.fileObj.size),
          url: this.fileSuccessObj.url,
          ext: "mp4",
          category_ids: this.formData.category_ids,
          start_time: this.formData.start_time,
          end_time: this.formData.end_time,
          relation_goods: this.formData.relation_goods || []
        }).then((res: any) => {
          if (res.code == 200) {
            this.$emit("closeDialog");
            this.$message.success("上传成功");
          } else {
            this.$message.error(res.message);
          }
        });
      } else {
        return;
      }
    });
  }
  selectNode(node: any) {
    this.formData.group_name = node.group_name;
    this.formData.group_id = node.id;
    (this.$refs.selectBlur as any).blur();
  }
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
  //上传成功
  private onSuccessFn(res: any) {
    console.log('onSuccessFn', res);
    
    this.isClosePop = true;
    this.fileSuccessObj = res.data;
    this.isFail = false;
    const { url, fileName, poster } = res.data;
    this.formData.name = fileName;
    this.$set(this.formData, "cover_url", poster);
    this.$set(this.formData, "url", url);

  }
  //移除上传的文件
  private removeFn() {
    this.isShowBtn = true;
  }
  // 打开添加素材标签弹窗
  private openAddCate() {
    this.addVisible = true;
  }
  // 关闭取消弹窗
  private closeVisible() {
    this.addVisible = false;
  }
  // 关闭预览视频
  private cancelVideo() {
    this.is_open_video = false;
  }
  // 关闭预览视频
  private openVideo() {
    this.is_open_video = true;
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
      this.formData.start_time = parseInt((val[0] / 1000).toString());
      this.formData.end_time = parseInt((val[1] / 1000).toString());
    } else {
      this.formData.start_time = 0;
      this.formData.end_time = 0;
    }
  }
  /**
   * 删除AI识图商品
   * @param item 
   * @param index 
   */
  handleDeleteGoods(item: any, index: number) {
    item.relation_goods.splice(index, 1)
    if(item.relation_goods.length === 0){
      item.is_ai = false
    }
  }
  /**
   * 选择商品
   * @param item 
   * @param goods 
   */
  handleSelectGoodsConfirm(item: any, goods: any) {
    console.log(item);
    console.log(goods);
    const upcList = item.relation_goods.map(v => v.upc)
    
    item.relation_goods = [
      ...item.relation_goods.filter(v => v.add_type === 2),
      ...goods.filter(v => !upcList.includes(v.upc))
    ]
    item.is_ai = item.relation_goods.filter(v => v.add_type === 2).length > 0
   }
  // ------------------生命周期-----------
  private mounted() {
    this.formData.group_id = this.groupId || -1;
    this.formData.group_name = this.groupId ? this.groupName : "未分组";
    this.grouplist();
    this.getCategory();
    if (this.sourceId) {
      this.sourceDetail();
    }
  }
}
</script>
<style lang="scss" scoped>
/deep/.el-dialog {
  max-height: 700px !important;
}

.main-box{
  max-height: 500px;
  overflow-y: auto;
}

.select-materials-dialog {
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
      border-color: #409eff;
      color: #409eff;
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
      color: #409eff;
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
  .tip-view {
    font-size: 12px;
    color: $color-danger;
    line-height: 18px;
    margin-top: 10px;
  }
  .vidio-up-item {
    position: relative;
    .look-btn {
      position: absolute;
      top: 0;
      right: -80px;
      color: #6f3112;
      text-decoration: underline;
      cursor: pointer;
    }
  }
}
.file-show {
  width: 80%;
  height: 36px;
  background: #000;
  cursor: default;
  color: #ffffff;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  line-height: 36px;
  padding: 0 15px;
  border-radius: 4px;
}
.select-video-box {
  position: relative;
  width: 80px;
  height: 80px;
  border-radius: 4px;
  margin-top: 12px;
  border: 1px solid #eee;
  cursor: pointer;


  .video-preview {
    width: 100%;
    height: 100%;
  }

  video {
    border-radius: 4px;
  }

  .video-play-icon {
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;
    top: 0;
    margin: auto;
    width: 25px;
    height: 25px;
    z-index: 99;
    cursor: pointer;
  }

  .delete-video-box {
    width: 20px;
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    position: absolute;
    cursor: pointer;
    right: -10px;
    top: -10px;
    z-index: 2;

    img {
      width: 14px;
      height: 14px;
      display: block;
    }
  }
}
</style>
